package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Get worker ID from args or use default
	workerID := "worker-1"
	if len(os.Args) > 1 {
		workerID = os.Args[1]
	}

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Printf("Worker %s - Connected to NATS\n", workerID)
	fmt.Printf("Worker %s - Subscribing to queue 'tasks'\n", workerID)

	// Subscribe to queue group
	// All workers with same queue name form a queue group
	// Each message goes to ONE worker only (load balancing)
	_, err = nc.QueueSubscribe("tasks", "task-workers", func(m *nats.Msg) {
		task := string(m.Data)
		fmt.Printf("Worker %s - Processing: %s\n", workerID, task)

		// Simulate work
		time.Sleep(1 * time.Second)

		fmt.Printf("Worker %s - Completed: %s\n", workerID, task)
		fmt.Println("---")
	})
	if err != nil {
		log.Fatal(err)
	}

	// Keep running
	fmt.Printf("Worker %s - Waiting for tasks...\n", workerID)
	select {}
}

/*
To test queue groups:

1. Start multiple workers in different terminals:
   terminal 1: go run worker.go worker-1
   terminal 2: go run worker.go worker-2
   terminal 3: go run worker.go worker-3

2. In another terminal, run the publisher:
   go run publisher.go

Expected output:
- Each task is processed by exactly ONE worker
- Work is distributed among all workers
- No task is processed twice
*/
