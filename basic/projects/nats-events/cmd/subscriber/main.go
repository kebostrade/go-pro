package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"

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
	log.Println("Subscriber created")

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Subscribe
	sub, err := nc.Subscribe(ordersTopic, func(msg *nats.Msg) {
		var order models.OrderEvent
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}
		log.Printf("Received order event: %+v", order)

		// Acknowledge the message
		if msg.Reply != "" {
			if err := msg.Respond([]byte("ack")); err != nil {
				log.Printf("Failed to respond: %v", err)
			}
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	log.Printf("Subscribed to %s. Waiting for messages...", ordersTopic)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down subscriber...")
	log.Println("Subscriber stopped")
}
