package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/datastructures"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Queue (FIFO) Data Structure Demo")

	// Demo 1: Basic Queue Operations
	demo1BasicOperations()

	// Demo 2: Customer Service Queue
	demo2CustomerService()

	// Demo 3: Priority Queue
	demo3PriorityQueue()

	// Demo 4: Task Scheduler
	demo4TaskScheduler()
}

func demo1BasicOperations() {
	utils.PrintSubHeader("1. Basic Queue Operations")

	queue := datastructures.NewQueue[int]()

	// Enqueue elements
	fmt.Println("Enqueueing elements: 1, 2, 3, 4, 5")
	for i := 1; i <= 5; i++ {
		queue.Enqueue(i)
		fmt.Printf("  Enqueued %d, size: %d\n", i, queue.Size())
	}

	// Peek
	if front, err := queue.Peek(); err == nil {
		fmt.Printf("\nFront element (peek): %d\n", front)
	}

	// Dequeue elements
	fmt.Println("\nDequeueing elements:")
	for !queue.IsEmpty() {
		if val, err := queue.Dequeue(); err == nil {
			fmt.Printf("  Dequeued %d, remaining size: %d\n", val, queue.Size())
		}
	}

	// Try to dequeue from empty queue
	if _, err := queue.Dequeue(); err != nil {
		fmt.Printf("\nError dequeueing from empty queue: %v\n", err)
	}
}

func demo2CustomerService() {
	utils.PrintSubHeader("2. Customer Service Queue Simulation")

	type Customer struct {
		ID   int
		Name string
	}

	queue := datastructures.NewQueue[Customer]()

	// Customers arrive
	customers := []Customer{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
		{ID: 4, Name: "Diana"},
	}

	fmt.Println("Customers arriving:")
	for _, customer := range customers {
		queue.Enqueue(customer)
		fmt.Printf("  %s (ID: %d) joined the queue\n", customer.Name, customer.ID)
	}

	fmt.Printf("\nQueue size: %d customers waiting\n", queue.Size())

	// Serve customers
	fmt.Println("\nServing customers:")
	servedCount := 0
	for !queue.IsEmpty() && servedCount < 3 {
		if customer, err := queue.Dequeue(); err == nil {
			fmt.Printf("  Now serving: %s (ID: %d)\n", customer.Name, customer.ID)
			servedCount++
		}
	}

	fmt.Printf("\nRemaining customers in queue: %d\n", queue.Size())
}

func demo3PriorityQueue() {
	utils.PrintSubHeader("3. Priority Queue")

	type Task struct {
		Name string
		ID   int
	}

	pq := datastructures.NewPriorityQueue[Task]()

	// Add tasks with different priorities
	tasks := []struct {
		task     Task
		priority datastructures.Priority
	}{
		{Task{Name: "Fix critical bug", ID: 1}, datastructures.PriorityHigh},
		{Task{Name: "Update documentation", ID: 2}, datastructures.PriorityLow},
		{Task{Name: "Code review", ID: 3}, datastructures.PriorityMedium},
		{Task{Name: "Security patch", ID: 4}, datastructures.PriorityHigh},
		{Task{Name: "Refactor code", ID: 5}, datastructures.PriorityMedium},
		{Task{Name: "Write tests", ID: 6}, datastructures.PriorityLow},
	}

	fmt.Println("Adding tasks to priority queue:")
	for _, t := range tasks {
		pq.Enqueue(t.task, t.priority)
		priorityName := "Medium"
		if t.priority == datastructures.PriorityHigh {
			priorityName = "High"
		} else if t.priority == datastructures.PriorityLow {
			priorityName = "Low"
		}
		fmt.Printf("  [%s] %s\n", priorityName, t.task.Name)
	}

	fmt.Printf("\nTotal tasks: %d\n", pq.Size())

	// Process tasks by priority
	fmt.Println("\nProcessing tasks (by priority):")
	for !pq.IsEmpty() {
		if task, err := pq.Dequeue(); err == nil {
			fmt.Printf("  Processing: %s (ID: %d)\n", task.Name, task.ID)
		}
	}
}

func demo4TaskScheduler() {
	utils.PrintSubHeader("4. Task Scheduler with Round-Robin")

	queue := datastructures.NewQueue[string]()

	// Add tasks
	tasks := []string{"Task A", "Task B", "Task C", "Task D"}
	fmt.Println("Adding tasks:")
	for _, task := range tasks {
		queue.Enqueue(task)
		fmt.Printf("  Added: %s\n", task)
	}

	// Simulate round-robin scheduling (3 rounds)
	fmt.Println("\nRound-robin execution (3 rounds):")
	for round := 1; round <= 3; round++ {
		fmt.Printf("\nRound %d:\n", round)

		// Process each task once
		size := queue.Size()
		for i := 0; i < size; i++ {
			if task, err := queue.Dequeue(); err == nil {
				fmt.Printf("  Executing: %s\n", task)
				// Re-enqueue for next round
				queue.Enqueue(task)
			}
		}
	}

	// Final cleanup
	fmt.Println("\nFinal queue state:")
	fmt.Printf("  Size: %d\n", queue.Size())
	fmt.Printf("  Contents: %s\n", queue.String())
}
