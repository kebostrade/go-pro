package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS")

	// Publish messages
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message #%d", i)

		// Publish to subject
		err := nc.Publish("greetings", []byte(message))
		if err != nil {
			log.Printf("Error publishing: %v", err)
			continue
		}

		fmt.Printf("Published: %s\n", message)
		time.Sleep(500 * time.Millisecond)
	}

	// Flush to ensure messages are sent
	nc.Flush()

	// Wait for messages to be received
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All messages published!")
}
