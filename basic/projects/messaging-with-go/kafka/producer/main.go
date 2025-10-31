package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
}

func main() {
	printBanner()

	// Kafka configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionSnappy

	// Create producer
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("✅ Kafka producer connected")
	log.Println("📤 Sending messages to topic 'events'...")

	// Send messages
	topic := "events"
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
			event := Event{
				ID:        fmt.Sprintf("event-%d", messageCount),
				Type:      "user.created",
				Timestamp: time.Now(),
				Data:      fmt.Sprintf("User %d created", messageCount),
			}

			// Marshal to JSON
			value, err := json.Marshal(event)
			if err != nil {
				log.Printf("❌ Failed to marshal event: %v", err)
				continue
			}

			// Create message
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(event.ID),
				Value: sarama.ByteEncoder(value),
			}

			// Send message
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("❌ Failed to send message: %v", err)
				continue
			}

			log.Printf("✅ Message sent: partition=%d offset=%d id=%s", partition, offset, event.ID)

		case <-sigterm:
			log.Println("\n🛑 Shutting down producer...")
			return
		}
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  📤 Kafka Producer                          ║
║                                                              ║
║        Sending events to Kafka topic 'events'               ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

