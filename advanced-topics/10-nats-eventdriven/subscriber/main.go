package main

import (
	"fmt"
	"log"

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
	fmt.Println("Subscribing to 'greetings'...")

	// Subscribe to subject
	_, err = nc.Subscribe("greetings", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
		fmt.Printf("Subject: %s\n", m.Subject)
		fmt.Printf("Reply: %s\n", m.Reply)
		fmt.Println("---")
	})
	if err != nil {
		log.Fatal(err)
	}

	// Keep running
	fmt.Println("Waiting for messages (Ctrl+C to stop)...")
	select {}
}
