package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
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

	// Declare exchange
	exchangeName := "events"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	log.Println("✅ RabbitMQ publisher connected")
	log.Printf("📤 Publishing messages to exchange '%s'...\n", exchangeName)

	// Publish messages
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Handle graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	messageCount := 0

	for {
		select {
		case <-ticker.C:
			messageCount++
			msg := Message{
				ID:        fmt.Sprintf("msg-%d", messageCount),
				Type:      "notification",
				Timestamp: time.Now(),
				Payload:   fmt.Sprintf("Notification %d", messageCount),
			}

			// Marshal to JSON
			body, err := json.Marshal(msg)
			if err != nil {
				log.Printf("❌ Failed to marshal message: %v", err)
				continue
			}

			// Routing key
			routingKey := "notification.email"

			// Publish message
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = ch.PublishWithContext(
				ctx,
				exchangeName, // exchange
				routingKey,   // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType:  "application/json",
					Body:         body,
					DeliveryMode: amqp.Persistent,
					Timestamp:    time.Now(),
					MessageId:    msg.ID,
				},
			)
			cancel()

			if err != nil {
				log.Printf("❌ Failed to publish message: %v", err)
				continue
			}

			log.Printf("✅ Message published: id=%s routing_key=%s", msg.ID, routingKey)

		case <-sigterm:
			log.Println("\n🛑 Shutting down publisher...")
			return
		}
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                📤 RabbitMQ Publisher                        ║
║                                                              ║
║        Publishing messages to RabbitMQ exchange             ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

