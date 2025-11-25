package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Goroutines & Concurrency Basics Demo")

	demo1BasicGoroutines()
	demo2WaitGroup()
	demo3FanOutFanIn()
	demo4SafeCounter()
	demo5SafeMap()
}

func demo1BasicGoroutines() {
	utils.PrintSubHeader("1. Basic Goroutines")

	tasks := []func(){
		func() {
			fmt.Println("Task 1: Processing...")
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Task 1: Done!")
		},
		func() {
			fmt.Println("Task 2: Processing...")
			time.Sleep(150 * time.Millisecond)
			fmt.Println("Task 2: Done!")
		},
		func() {
			fmt.Println("Task 3: Processing...")
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Task 3: Done!")
		},
	}

	fmt.Println("Running tasks concurrently...")
	start := time.Now()
	concurrency.RunConcurrent(tasks)
	duration := time.Since(start)

	fmt.Printf("\nAll tasks completed in %v\n", duration)
	fmt.Println("(Notice they run concurrently, not sequentially)")
}

func demo2WaitGroup() {
	utils.PrintSubHeader("2. WaitGroup Pattern")

	ctx := context.Background()

	tasks := []func(context.Context){
		func(ctx context.Context) {
			fmt.Println("Worker 1: Starting...")
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Worker 1: Finished")
		},
		func(ctx context.Context) {
			fmt.Println("Worker 2: Starting...")
			time.Sleep(200 * time.Millisecond)
			fmt.Println("Worker 2: Finished")
		},
		func(ctx context.Context) {
			fmt.Println("Worker 3: Starting...")
			time.Sleep(150 * time.Millisecond)
			fmt.Println("Worker 3: Finished")
		},
	}

	fmt.Println("Running workers with context...")
	start := time.Now()
	err := concurrency.RunConcurrentWithContext(ctx, tasks)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("\nAll workers completed in %v\n", duration)
	}
}

func demo3FanOutFanIn() {
	utils.PrintSubHeader("3. Fan-Out/Fan-In Pattern")

	// Input data
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("Input: %v\n", numbers)

	// Process function: square the number
	process := func(n int) int {
		time.Sleep(50 * time.Millisecond) // Simulate work
		return n * n
	}

	// Fan-out to 4 workers
	fmt.Println("\nFanning out to 4 workers...")
	start := time.Now()
	results := concurrency.FanOut(numbers, 4, process)
	duration := time.Since(start)

	fmt.Printf("Results: %v\n", results)
	fmt.Printf("Completed in %v (with 4 workers)\n", duration)

	// Compare with sequential
	fmt.Println("\nSequential processing for comparison...")
	start = time.Now()
	sequential := make([]int, len(numbers))
	for i, n := range numbers {
		sequential[i] = process(n)
	}
	seqDuration := time.Since(start)

	fmt.Printf("Sequential completed in %v\n", seqDuration)
	fmt.Printf("Speedup: %.2fx\n", float64(seqDuration)/float64(duration))
}

func demo4SafeCounter() {
	utils.PrintSubHeader("4. Thread-Safe Counter")

	counter := concurrency.NewSafeCounter()

	// Increment counter concurrently
	tasks := make([]func(), 100)
	for i := 0; i < 100; i++ {
		tasks[i] = func() {
			counter.Increment()
		}
	}

	fmt.Println("Incrementing counter 100 times concurrently...")
	concurrency.RunConcurrent(tasks)

	fmt.Printf("Final counter value: %d\n", counter.Value())
	fmt.Println("(Should be exactly 100 due to thread-safety)")

	// Add values
	counter.Reset()
	addTasks := make([]func(), 10)
	for i := 0; i < 10; i++ {
		val := int64(i + 1)
		addTasks[i] = func() {
			counter.Add(val)
		}
	}

	concurrency.RunConcurrent(addTasks)
	fmt.Printf("\nAfter adding 1+2+3+...+10: %d\n", counter.Value())
}

func demo5SafeMap() {
	utils.PrintSubHeader("5. Thread-Safe Map")

	safeMap := concurrency.NewSafeMap[string, int]()

	// Set values concurrently
	tasks := make([]func(), 10)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		value := i * 10
		tasks[i] = func() {
			safeMap.Set(key, value)
		}
	}

	fmt.Println("Setting 10 key-value pairs concurrently...")
	concurrency.RunConcurrent(tasks)

	fmt.Printf("Map size: %d\n", safeMap.Len())
	fmt.Printf("Map contents: %s\n", safeMap.String())

	// Get values
	fmt.Println("\nRetrieving values:")
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key%d", i)
		if val, ok := safeMap.Get(key); ok {
			fmt.Printf("  %s = %d\n", key, val)
		}
	}

	// ForEach
	fmt.Println("\nIterating with ForEach:")
	safeMap.ForEach(func(k string, v int) {
		fmt.Printf("  %s: %d\n", k, v)
	})

	// Delete
	safeMap.Delete("key0")
	fmt.Printf("\nAfter deleting 'key0', size: %d\n", safeMap.Len())
}
