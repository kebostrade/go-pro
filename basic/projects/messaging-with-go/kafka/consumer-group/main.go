package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	count int
	mu    sync.Mutex
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			consumer.mu.Lock()
			consumer.count++
			count := consumer.count
			consumer.mu.Unlock()

			var event Event
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("❌ Failed to unmarshal message: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			log.Printf("✅ Message received [%d]: partition=%d offset=%d",
				count, message.Partition, message.Offset)
			log.Printf("   Event: ID=%s Type=%s Data=%s", event.ID, event.Type, event.Data)

			// Mark message as processed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func main() {
	printBanner()

	// Kafka configuration
	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create consumer group
	brokers := []string{"localhost:9092"}
	group := "events-consumer-group"
	topics := []string{"events"}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	log.Printf("✅ Consumer group '%s' connected", group)
	log.Printf("📥 Consuming messages from topics: %v\n", topics)

	ctx, cancel := context.WithCancel(context.Background())
	consumer := &Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Printf("❌ Error from consumer: %v", err)
			}
			// Check if context was cancelled
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Wait till consumer is ready

	// Handle graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
	log.Println("\n🛑 Shutting down consumer group...")
	cancel()
	wg.Wait()

	consumer.mu.Lock()
	count := consumer.count
	consumer.mu.Unlock()
	log.Printf("✅ Consumer group stopped (processed %d messages)", count)
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║              📥 Kafka Consumer Group                        ║
║                                                              ║
║        Consuming events with load balancing                 ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

