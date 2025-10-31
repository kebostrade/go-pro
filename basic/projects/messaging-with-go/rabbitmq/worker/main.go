package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Task struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
	Data     string `json:"data"`
}

func main() {
	printBanner()

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare queue
	queueName := "tasks"
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Set QoS - only process one task at a time
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("✅ RabbitMQ worker connected")
	log.Printf("👷 Processing tasks from queue '%s'...\n", queueName)

	// Handle graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	taskCount := 0

	for {
		select {
		case msg := <-msgs:
			taskCount++

			var task Task
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Printf("❌ Failed to unmarshal task: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("🔨 Processing task [%d]: id=%s type=%s priority=%d",
				taskCount, task.ID, task.Type, task.Priority)

			// Simulate task processing
			processingTime := time.Duration(rand.Intn(3)+1) * time.Second
			time.Sleep(processingTime)

			log.Printf("✅ Task completed: id=%s (took %s)", task.ID, processingTime)

			// Acknowledge message
			msg.Ack(false)

		case <-sigterm:
			log.Printf("\n🛑 Shutting down worker... (processed %d tasks)", taskCount)
			return
		}
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                 👷 RabbitMQ Worker                          ║
║                                                              ║
║        Processing tasks from RabbitMQ queue                 ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

