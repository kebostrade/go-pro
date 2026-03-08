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

	fmt.Println("Task Publisher")
	fmt.Println("Connected to NATS")

	// Publish 10 tasks
	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("Task #%d", i)

		err := nc.Publish("tasks", []byte(task))
		if err != nil {
			log.Printf("Error publishing: %v", err)
			continue
		}

		fmt.Printf("Published task: %s\n", task)
		time.Sleep(500 * time.Millisecond)
	}

	// Flush and wait
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAll tasks published!")
	fmt.Println("Note: In a queue group, each message goes to ONE worker only")
}
