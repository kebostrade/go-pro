//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// WORKER POOL PATTERN
// ============================================================================

// Task represents a unit of work
type Task struct {
	ID   int
	Data string
}

// Result represents the output of a task
type Result struct {
	TaskID  int
	Output  string
	Elapsed time.Duration
}

// WorkerPool manages a pool of goroutines
type WorkerPool struct {
	workers   int
	taskQueue chan Task
	results   chan Result
	wg        sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan Task, queueSize),
		results:   make(chan Result, queueSize),
	}
}

// Start initializes the worker goroutines
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for task := range wp.taskQueue {
		start := time.Now()

		// Simulate work
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		output := fmt.Sprintf("Worker %d processed task %d: %s", id, task.ID, task.Data)

		wp.results <- Result{
			TaskID:  task.ID,
			Output:  output,
			Elapsed: time.Since(start),
		}
	}
}

// Submit adds a task to the queue
func (wp *WorkerPool) Submit(task Task) {
	wp.taskQueue <- task
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	wp.wg.Wait()
	close(wp.results)
}

// Results returns a channel for reading results
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// Example: Worker Pool
func exampleWorkerPool() {
	fmt.Println("\n=== Worker Pool Example ===")

	pool := NewWorkerPool(3, 10)
	pool.Start()

	// Submit tasks
	for i := 1; i <= 5; i++ {
		task := Task{
			ID:   i,
			Data: fmt.Sprintf("data-%d", i),
		}
		pool.Submit(task)
	}

	// Collect results
	go func() {
		for result := range pool.Results() {
			fmt.Printf("Task %d: %s (took %v)\n", result.TaskID, result.Output, result.Elapsed)
		}
	}()

	pool.Stop()
	time.Sleep(100 * time.Millisecond) // Wait for all results
}

// ============================================================================
// FAN-OUT/FAN-IN PATTERN
// ============================================================================

// fanOut distributes work to multiple workers
func fanOut(input <-chan int, workers int) []<-chan int {
	channels := make([]<-chan int, workers)

	for i := 0; i < workers; i++ {
		out := make(chan int)
		channels[i] = out

		go func(ch chan<- int) {
			for val := range input {
				// Process value
				result := val * val
				ch <- result
			}
			close(ch)
		}(out)
	}

	return channels
}

// fanIn combines multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Example: Fan-Out/Fan-In
func exampleFanOutFanIn() {
	fmt.Println("\n=== Fan-Out/Fan-In Example ===")

	// Input channel
	input := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			input <- i
		}
		close(input)
	}()

	// Fan out to 3 workers
	channels := fanOut(input, 3)

	// Fan in results
	results := fanIn(channels...)

	// Print results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// ============================================================================
// PIPELINE PATTERN
// ============================================================================

// Pipeline stage 1: Generate numbers
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Pipeline stage 2: Square numbers
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Pipeline stage 3: Filter even numbers
func filter(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

// Example: Pipeline
func examplePipeline() {
	fmt.Println("\n=== Pipeline Example ===")

	// Create pipeline: generator -> square -> filter
	numbers := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(numbers)
	filtered := filter(squared)

	// Consume results
	for result := range filtered {
		fmt.Printf("Pipeline result: %d\n", result)
	}
}

// ============================================================================
// RATE LIMITING PATTERN
// ============================================================================

// RateLimiter limits the rate of operations
type RateLimiter struct {
	ticker   *time.Ticker
	requests chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		ticker:   time.NewTicker(rate),
		requests: make(chan struct{}),
	}

	go rl.run()
	return rl
}

// run processes requests at the specified rate
func (rl *RateLimiter) run() {
	for range rl.ticker.C {
		// Allow one request per tick
		select {
		case <-rl.requests:
		default:
		}
	}
}

// Allow waits for permission to proceed
func (rl *RateLimiter) Allow() {
	rl.requests <- struct{}{}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
}

// Example: Rate Limiting
func exampleRateLimiting() {
	fmt.Println("\n=== Rate Limiting Example ===")

	limiter := NewRateLimiter(500 * time.Millisecond)
	defer limiter.Stop()

	// Try to do 5 things, but limited to 2 per second
	for i := 1; i <= 5; i++ {
		limiter.Allow()
		fmt.Printf("Operation %d at %v\n", i, time.Now().Format("15:04:05.000"))
	}
}

// ============================================================================
// ATOMIC COUNTERS PATTERN
// ============================================================================

// AtomicCounter uses atomic operations for thread-safe counting
type AtomicCounter struct {
	value int64
}

// Increment increments the counter atomically
func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Get returns the current value
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Example: Atomic Counters
func exampleAtomicCounters() {
	fmt.Println("\n=== Atomic Counters Example ===")

	counter := &AtomicCounter{}
	var wg sync.WaitGroup

	// Start 100 goroutines that all increment the counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Get())
}

// ============================================================================
// MAIN - Run all examples
// ============================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	// Run all examples
	exampleWorkerPool()
	exampleFanOutFanIn()
	examplePipeline()
	exampleRateLimiting()
	exampleAtomicCounters()

	fmt.Println("\n=== All examples completed ===")
}

// ============================================================================
// OUTPUT EXPLANATION
// ============================================================================

/*
Expected Output:
=== Worker Pool Example ===
Task 1: Worker 0 processed task 1: data-1 (took 50ms)
Task 2: Worker 1 processed task 2: data-2 (took 30ms)
...

=== Fan-Out/Fan-In Example ===
Result: 1
Result: 4
Result: 9
...

=== Pipeline Example ===
Pipeline result: 4
Pipeline result: 16
Pipeline result: 36
...

=== Rate Limiting Example ===
Operation 1 at 15:04:05.123
Operation 2 at 15:04:05.623
Operation 3 at 15:04:06.123
...

=== Atomic Counters Example ===
Final counter value: 100
*/
