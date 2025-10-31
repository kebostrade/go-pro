package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func main() {
	printBanner()

	// Kafka configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create consumer
	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	log.Println("✅ Kafka consumer connected")

	// Subscribe to topic
	topic := "events"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer for partition: %v", err)
	}
	defer partitionConsumer.Close()

	log.Printf("📥 Consuming messages from topic '%s'...\n", topic)

	// Handle graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	messageCount := 0

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			messageCount++

			var event Event
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("❌ Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("✅ Message received [%d]: partition=%d offset=%d key=%s",
				messageCount, msg.Partition, msg.Offset, string(msg.Key))
			log.Printf("   Event: ID=%s Type=%s Data=%s", event.ID, event.Type, event.Data)

		case err := <-partitionConsumer.Errors():
			log.Printf("❌ Error: %v", err)

		case <-sigterm:
			log.Printf("\n🛑 Shutting down consumer... (processed %d messages)", messageCount)
			return
		}
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  📥 Kafka Consumer                          ║
║                                                              ║
║        Consuming events from Kafka topic 'events'           ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

