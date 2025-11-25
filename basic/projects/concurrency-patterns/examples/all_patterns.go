package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/concurrency-patterns/patterns"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║           🔄 Concurrency Patterns in Go - Demo              ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 1. Worker Pool
	demoWorkerPool()

	// 2. Pipeline
	demoPipeline()

	// 3. Fan-Out/Fan-In
	demoFanOutFanIn()

	// 4. Rate Limiting
	demoRateLimiting()

	// 5. Semaphore Pattern
	demoSemaphore()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    ✅ Demo Complete!                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

// 1. Worker Pool Pattern
func demoWorkerPool() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("1️⃣  WORKER POOL PATTERN")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	patterns.DemoWorkerPool()
}

// 2. Pipeline Pattern
func demoPipeline() {
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("2️⃣  PIPELINE PATTERN")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Stage 1: Generate numbers
	numbers := generate(1, 2, 3, 4, 5)

	// Stage 2: Square numbers
	squares := square(numbers)

	// Stage 3: Sum squares
	sum := sum(squares)

	fmt.Printf("Sum of squares: %d\n", sum)
}

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func sum(in <-chan int) int {
	sum := 0
	for n := range in {
		sum += n
	}
	return sum
}

// 3. Fan-Out/Fan-In Pattern
func demoFanOutFanIn() {
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("3️⃣  FAN-OUT/FAN-IN PATTERN")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Input channel
	input := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	// Fan-out: Multiple workers process input
	numWorkers := 3
	workers := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = fanOutWorker(i+1, input)
	}

	// Fan-in: Merge results
	results := fanIn(workers...)

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

func fanOutWorker(id int, input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range input {
			fmt.Printf("Worker %d processing %d\n", id, num)
			time.Sleep(50 * time.Millisecond)
			out <- num * 2
		}
	}()
	return out
}

func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
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

// 4. Rate Limiting Pattern
func demoRateLimiting() {
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("4️⃣  RATE LIMITING PATTERN")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Allow 2 requests per second
	limiter := time.Tick(500 * time.Millisecond)

	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	for req := range requests {
		<-limiter // Wait for rate limiter
		fmt.Printf("⏱️  Request %d processed at %s\n", req, time.Now().Format("15:04:05.000"))
	}
}

// 5. Semaphore Pattern
func demoSemaphore() {
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("5️⃣  SEMAPHORE PATTERN")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Limit to 2 concurrent operations
	semaphore := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			fmt.Printf("🔒 Task %d acquired semaphore\n", id)

			// Do work
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("✅ Task %d completed\n", id)

			// Release semaphore
			<-semaphore
			fmt.Printf("🔓 Task %d released semaphore\n", id)
		}(i)
	}

	wg.Wait()
}

/*
CONCURRENCY PATTERNS SUMMARY:

1. WORKER POOL
   - Fixed number of workers
   - Process jobs from queue
   - Reuse goroutines

2. PIPELINE
   - Chain of processing stages
   - Each stage is a goroutine
   - Data flows through channels

3. FAN-OUT/FAN-IN
   - Fan-out: Distribute work to multiple workers
   - Fan-in: Merge results from workers
   - Parallel processing

4. RATE LIMITING
   - Control request rate
   - Use time.Tick or time.NewTicker
   - Prevent resource exhaustion

5. SEMAPHORE
   - Limit concurrent operations
   - Use buffered channel
   - Control resource access

BEST PRACTICES:
- Always close channels when done sending
- Use WaitGroup to wait for goroutines
- Handle errors properly
- Use context for cancellation
- Avoid goroutine leaks
*/

