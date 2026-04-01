package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"basic/projects/nats-events/internal/jetstream"
	"basic/projects/nats-events/internal/models"
)

const (
	natsURL     = "nats://localhost:4222"
	ordersTopic = "events.orders"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS")

	// Create JetStream publisher
	publisher, err := jetstream.NewPublisher(nc)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	log.Println("Publisher created")

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start publishing events every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		orderNum := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				orderNum++
				event := models.OrderEvent{
					EventType: "created",
					OrderID:   "ORD-" + string(rune('0'+orderNum%10)),
					UserID:    "USER-" + string(rune('0'+(orderNum%5+1))),
					Amount:    float64(orderNum * 100),
					Timestamp: time.Now(),
				}

				if err := publisher.PublishJSON(ctx, ordersTopic, event); err != nil {
					log.Printf("Failed to publish event: %v", err)
				} else {
					log.Printf("Published order event: %+v", event)
				}
			}
		}
	}()

	log.Println("Publisher started. Press Ctrl+C to exit.")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down publisher...")
	cancel()
	log.Println("Publisher stopped")
}
