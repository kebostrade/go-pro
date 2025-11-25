//go:build ignore

package main

import "fmt"

// Task: Implement a FIFO (First-In-First-Out) queue data structure.
// This complements the LIFO stack implementation.

// Queue represents a FIFO (First-In-First-Out) data structure
type Queue struct {
	items []string
}

// Enqueue adds an item to the back of the queue
func (q *Queue) Enqueue(item string) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue
func (q *Queue) Dequeue() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the front item without removing it
func (q *Queue) Peek() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}
	return q.items[0], true
}

// IsEmpty checks if the queue has no items
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q *Queue) Size() int {
	return len(q.items)
}

// Clear removes all items from the queue
func (q *Queue) Clear() {
	q.items = []string{}
}

// PriorityQueue represents a queue with priority levels
type PriorityQueue struct {
	high   []string
	medium []string
	low    []string
}

// Enqueue adds an item with a priority level (1=high, 2=medium, 3=low)
func (pq *PriorityQueue) Enqueue(item string, priority int) {
	switch priority {
	case 1:
		pq.high = append(pq.high, item)
	case 2:
		pq.medium = append(pq.medium, item)
	case 3:
		pq.low = append(pq.low, item)
	default:
		pq.medium = append(pq.medium, item)
	}
}

// Dequeue removes and returns the highest priority item
func (pq *PriorityQueue) Dequeue() (string, bool) {
	// Check high priority first
	if len(pq.high) > 0 {
		item := pq.high[0]
		pq.high = pq.high[1:]
		return item, true
	}

	// Then medium priority
	if len(pq.medium) > 0 {
		item := pq.medium[0]
		pq.medium = pq.medium[1:]
		return item, true
	}

	// Finally low priority
	if len(pq.low) > 0 {
		item := pq.low[0]
		pq.low = pq.low[1:]
		return item, true
	}

	return "", false
}

// IsEmpty checks if all priority queues are empty
func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.high) == 0 && len(pq.medium) == 0 && len(pq.low) == 0
}

func main() {
	fmt.Println("Queue (FIFO) Data Structure Demo")
	fmt.Println("=".repeat(60))

	// Demo 1: Basic Queue Operations
	fmt.Println("\n1. Basic Queue Operations")
	fmt.Println("-".repeat(60))

	queue := &Queue{}

	// Enqueue items
	queue.Enqueue("First")
	queue.Enqueue("Second")
	queue.Enqueue("Third")

	fmt.Printf("Queue size: %d\n", queue.Size())

	// Peek at front
	if item, ok := queue.Peek(); ok {
		fmt.Printf("Front item (peek): %s\n", item)
	}

	// Dequeue all items
	fmt.Println("\nDequeuing all items:")
	for !queue.IsEmpty() {
		if item, ok := queue.Dequeue(); ok {
			fmt.Printf("  Dequeued: %s\n", item)
		}
	}

	// Demo 2: Customer Service Queue
	fmt.Println("\n2. Customer Service Queue Simulation")
	fmt.Println("-".repeat(60))

	customerQueue := &Queue{}

	// Customers arrive
	customers := []string{"Alice", "Bob", "Carol", "David", "Eve"}
	fmt.Println("Customers arriving:")
	for _, customer := range customers {
		customerQueue.Enqueue(customer)
		fmt.Printf("  %s joined the queue\n", customer)
	}

	fmt.Printf("\nQueue size: %d customers waiting\n", customerQueue.Size())

	// Serve customers
	fmt.Println("\nServing customers:")
	servedCount := 0
	for !customerQueue.IsEmpty() && servedCount < 3 {
		if customer, ok := customerQueue.Dequeue(); ok {
			fmt.Printf("  Now serving: %s\n", customer)
			servedCount++
		}
	}

	fmt.Printf("\nRemaining customers in queue: %d\n", customerQueue.Size())

	// Demo 3: Priority Queue
	fmt.Println("\n3. Priority Queue")
	fmt.Println("-".repeat(60))

	pq := &PriorityQueue{}

	// Add tasks with different priorities
	pq.Enqueue("Fix critical bug", 1)     // High priority
	pq.Enqueue("Update documentation", 3) // Low priority
	pq.Enqueue("Code review", 2)          // Medium priority
	pq.Enqueue("Security patch", 1)       // High priority
	pq.Enqueue("Refactor code", 2)        // Medium priority
	pq.Enqueue("Write tests", 3)          // Low priority

	fmt.Println("Processing tasks by priority:")
	for !pq.IsEmpty() {
		if task, ok := pq.Dequeue(); ok {
			fmt.Printf("  Processing: %s\n", task)
		}
	}

	// Demo 4: Circular Buffer (Ring Buffer) concept
	fmt.Println("\n4. Fixed-Size Queue (Circular Buffer Concept)")
	fmt.Println("-".repeat(60))

	maxSize := 3
	circularQueue := &Queue{}

	items := []string{"A", "B", "C", "D", "E"}
	for _, item := range items {
		// If queue is full, remove oldest item
		if circularQueue.Size() >= maxSize {
			if removed, ok := circularQueue.Dequeue(); ok {
				fmt.Printf("Queue full, removed: %s\n", removed)
			}
		}
		circularQueue.Enqueue(item)
		fmt.Printf("Added: %s (size: %d)\n", item, circularQueue.Size())
	}

	fmt.Println("\nFinal queue contents:")
	for !circularQueue.IsEmpty() {
		if item, ok := circularQueue.Dequeue(); ok {
			fmt.Printf("  %s\n", item)
		}
	}

	fmt.Println("\nQueue demo completed!")
}

// Helper function to repeat strings
func (s string) repeat(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
