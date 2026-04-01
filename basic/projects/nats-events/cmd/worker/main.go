package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"

	"basic/projects/nats-events/internal/models"
	"basic/projects/nats-events/internal/queue"
)

const (
	natsURL    = "nats://localhost:4222"
	groupName  = "task-workers"
	tasksTopic = "tasks"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS")

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Task handler
	taskHandler := func(data []byte) error {
		var task models.TaskEvent
		if err := json.Unmarshal(data, &task); err != nil {
			log.Printf("Failed to unmarshal task: %v", err)
			return err
		}

		log.Printf("Processing task: %s (type: %s)", task.TaskID, task.Type)

		// Simulate work
		// In real implementation, this would do actual work
		return nil
	}

	// Create and start worker
	worker := queue.NewWorker(nc, groupName, tasksTopic, taskHandler)
	if err := worker.Start(ctx); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	log.Printf("Worker started in group %s. Waiting for tasks...", groupName)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down worker...")
	cancel()
	log.Println("Worker stopped")
}
