package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// ============================================================================
// GO CONCURRENCY CRASH COURSE - RUNNABLE EXAMPLES
// ============================================================================
// Run: go run main.go
// Run with race detector: go run -race main.go
// ============================================================================

func main() {
	fmt.Println("🚀 Go Concurrency Crash Course - Interactive Examples")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Uncomment the example you want to run:

	// Basic Examples
	example1_BasicGoroutines()
	// example2_Channels()
	// example3_WaitGroups()
	// example4_Select()

	// Patterns
	// example5_WorkerPool()
	// example6_Pipeline()
	// example7_FanOutFanIn()

	// Advanced
	// example8_Context()
	// example9_Mutex()

	// Real-World
	// example10_WebScraper()
	// example11_RateLimiter()

	// Pitfalls (uncomment to see errors)
	// pitfall1_GoroutineLeak()
	// pitfall2_RaceCondition()
	// pitfall3_LoopVariableCapture()
}

// ============================================================================
// EXAMPLE 1: Basic Goroutines
// ============================================================================

func example1_BasicGoroutines() {
	fmt.Println("📌 Example 1: Basic Goroutines")
	fmt.Println("---")

	// Sequential execution
	fmt.Println("Sequential:")
	sayHello("Alice")
	sayHello("Bob")

	// Concurrent execution
	fmt.Println("\nConcurrent:")
	go sayHello("Charlie")
	go sayHello("Diana")

	// Wait for goroutines (we'll improve this in example 3)
	time.Sleep(time.Second)
	fmt.Println()
}

func sayHello(name string) {
	fmt.Printf("  Hello, %s!\n", name)
}

// ============================================================================
// EXAMPLE 2: Channels
// ============================================================================

func example2_Channels() {
	fmt.Println("📌 Example 2: Channels")
	fmt.Println("---")

	// Unbuffered channel
	fmt.Println("Unbuffered channel:")
	ch1 := make(chan string)
	go func() {
		ch1 <- "Hello from goroutine"
	}()
	msg := <-ch1
	fmt.Printf("  Received: %s\n", msg)

	// Buffered channel
	fmt.Println("\nBuffered channel:")
	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	fmt.Printf("  Received: %d\n", <-ch2)
	fmt.Printf("  Received: %d\n", <-ch2)

	// Range over channel
	fmt.Println("\nRange over channel:")
	ch3 := make(chan int, 3)
	ch3 <- 1
	ch3 <- 2
	ch3 <- 3
	close(ch3)

	for val := range ch3 {
		fmt.Printf("  Value: %d\n", val)
	}
	fmt.Println()
}

// ============================================================================
// EXAMPLE 3: WaitGroups
// ============================================================================

func example3_WaitGroups() {
	fmt.Println("📌 Example 3: WaitGroups")
	fmt.Println("---")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  Worker %d starting\n", id)
			time.Sleep(time.Second)
			fmt.Printf("  Worker %d done\n", id)
		}(i) // Pass i as parameter!
	}

	wg.Wait()
	fmt.Println("  All workers completed")
	fmt.Println()
}

// ============================================================================
// EXAMPLE 4: Select
// ============================================================================

func example4_Select() {
	fmt.Println("📌 Example 4: Select Statement")
	fmt.Println("---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "from channel 2"
	}()

	// Wait for both
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("  Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("  Received: %s\n", msg2)
		}
	}

	// Timeout example
	fmt.Println("\nTimeout example:")
	ch3 := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch3 <- "result"
	}()

	select {
	case res := <-ch3:
		fmt.Printf("  Got: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("  Timeout!")
	}
	fmt.Println()
}

// ============================================================================
// EXAMPLE 5: Worker Pool Pattern
// ============================================================================

func example5_WorkerPool() {
	fmt.Println("📌 Example 5: Worker Pool Pattern")
	fmt.Println("---")

	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Wait for workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("  Results:")
	for result := range results {
		fmt.Printf("    Result: %d\n", result)
	}
	fmt.Println()
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("  Worker %d processing job %d\n", id, job)
		time.Sleep(500 * time.Millisecond)
		results <- job * 2
	}
}

// ============================================================================
// EXAMPLE 6: Pipeline Pattern
// ============================================================================

func example6_Pipeline() {
	fmt.Println("📌 Example 6: Pipeline Pattern")
	fmt.Println("---")

	// Chain stages
	numbers := generate(1, 2, 3, 4, 5)
	squares := square(numbers)

	// Consume
	fmt.Println("  Results:")
	for result := range squares {
		fmt.Printf("    %d\n", result)
	}
	fmt.Println()
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

// ============================================================================
// EXAMPLE 7: Fan-Out/Fan-In Pattern
// ============================================================================

func example7_FanOutFanIn() {
	fmt.Println("📌 Example 7: Fan-Out/Fan-In Pattern")
	fmt.Println("---")

	// Input
	in := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
		}
		close(in)
	}()

	// Fan-out to 3 workers
	workers := fanOut(in, 3)

	// Fan-in results
	results := fanIn(workers...)

	// Consume
	fmt.Println("  Results:")
	for result := range results {
		fmt.Printf("    %d\n", result)
	}
	fmt.Println()
}

func fanOut(in <-chan int, workers int) []<-chan int {
	channels := make([]<-chan int, workers)
	for i := 0; i < workers; i++ {
		channels[i] = process(in)
	}
	return channels
}

func process(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

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

// ============================================================================
// EXAMPLE 8: Context for Cancellation
// ============================================================================

func example8_Context() {
	fmt.Println("📌 Example 8: Context for Cancellation")
	fmt.Println("---")

	ctx, cancel := context.WithCancel(context.Background())

	go contextWorker(ctx, 1)
	go contextWorker(ctx, 2)

	time.Sleep(2 * time.Second)
	fmt.Println("  Cancelling workers...")
	cancel()

	time.Sleep(time.Second)
	fmt.Println()
}

func contextWorker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  Worker %d cancelled\n", id)
			return
		default:
			fmt.Printf("  Worker %d working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// ============================================================================
// EXAMPLE 9: Mutex for Shared State
// ============================================================================

func example9_Mutex() {
	fmt.Println("📌 Example 9: Mutex for Shared State")
	fmt.Println("---")

	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("  Final counter: %d (expected: 1000)\n", counter)
	fmt.Println()
}

// ============================================================================
// EXAMPLE 10: Real-World Web Scraper
// ============================================================================

type Result struct {
	URL    string
	Status int
	Error  error
}

func example10_WebScraper() {
	fmt.Println("📌 Example 10: Concurrent Web Scraper")
	fmt.Println("---")

	urls := []string{
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
	}

	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	// Launch concurrent fetches
	for _, url := range urls {
		wg.Add(1)
		go fetch(url, results, &wg)
	}

	// Wait and close
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("  Results:")
	for result := range results {
		if result.Error != nil {
			fmt.Printf("    ❌ %s: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("    ✅ %s: %d\n", result.URL, result.Status)
		}
	}
	fmt.Println()
}

func fetch(url string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		results <- Result{URL: url, Error: err}
		return
	}
	defer resp.Body.Close()

	// Read and discard body
	io.Copy(io.Discard, resp.Body)

	results <- Result{URL: url, Status: resp.StatusCode}
}

// ============================================================================
// EXAMPLE 11: Rate Limiter
// ============================================================================

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, requestsPerSecond),
	}

	// Refill tokens
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer ticker.Stop()

		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

func example11_RateLimiter() {
	fmt.Println("📌 Example 11: Rate Limiter")
	fmt.Println("---")

	limiter := NewRateLimiter(2) // 2 requests per second

	fmt.Println("  Making 6 requests (rate: 2/sec):")
	for i := 1; i <= 6; i++ {
		limiter.Wait()
		go func(id int) {
			fmt.Printf("    Request %d at %s\n", id, time.Now().Format("15:04:05"))
		}(i)
	}

	time.Sleep(4 * time.Second)
	fmt.Println()
}

// ============================================================================
// COMMON PITFALLS (Uncomment to see errors)
// ============================================================================

// Pitfall 1: Goroutine Leak
func pitfall1_GoroutineLeak() {
	fmt.Println("⚠️  Pitfall 1: Goroutine Leak")
	fmt.Println("---")

	ch := make(chan int)
	go func() {
		val := <-ch // This will block forever!
		fmt.Println(val)
	}()

	// Channel never receives data
	fmt.Printf("  Goroutines before: %d\n", runtime.NumGoroutine())
	time.Sleep(time.Second)
	fmt.Printf("  Goroutines after: %d (leaked!)\n", runtime.NumGoroutine())
	fmt.Println()
}

// Pitfall 2: Race Condition (run with -race flag)
func pitfall2_RaceCondition() {
	fmt.Println("⚠️  Pitfall 2: Race Condition")
	fmt.Println("---")
	fmt.Println("  Run with: go run -race main.go")

	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // RACE CONDITION!
		}()
	}

	wg.Wait()
	fmt.Printf("  Counter: %d (expected: 1000, but may vary!)\n", counter)
	fmt.Println()
}

// Pitfall 3: Loop Variable Capture
func pitfall3_LoopVariableCapture() {
	fmt.Println("⚠️  Pitfall 3: Loop Variable Capture")
	fmt.Println("---")

	fmt.Println("  Bad (all print same value):")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("    %d ", i) // Wrong!
		}()
	}
	wg.Wait()
	fmt.Println()

	fmt.Println("  Good (pass as parameter):")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("    %d ", id) // Correct!
		}(i)
	}
	wg.Wait()
	fmt.Println()
}

// ============================================================================
// BONUS: Interactive Menu
// ============================================================================

func runInteractive() {
	examples := map[string]func(){
		"1":  example1_BasicGoroutines,
		"2":  example2_Channels,
		"3":  example3_WaitGroups,
		"4":  example4_Select,
		"5":  example5_WorkerPool,
		"6":  example6_Pipeline,
		"7":  example7_FanOutFanIn,
		"8":  example8_Context,
		"9":  example9_Mutex,
		"10": example10_WebScraper,
		"11": example11_RateLimiter,
	}

	for {
		fmt.Println("\n🚀 Go Concurrency Crash Course")
		fmt.Println("=" + string(make([]byte, 40)))
		fmt.Println("1.  Basic Goroutines")
		fmt.Println("2.  Channels")
		fmt.Println("3.  WaitGroups")
		fmt.Println("4.  Select")
		fmt.Println("5.  Worker Pool")
		fmt.Println("6.  Pipeline")
		fmt.Println("7.  Fan-Out/Fan-In")
		fmt.Println("8.  Context")
		fmt.Println("9.  Mutex")
		fmt.Println("10. Web Scraper")
		fmt.Println("11. Rate Limiter")
		fmt.Println("0.  Exit")
		fmt.Print("\nSelect example: ")

		var choice string
		fmt.Scanln(&choice)

		if choice == "0" {
			fmt.Println("Goodbye! 👋")
			break
		}

		if fn, ok := examples[choice]; ok {
			fmt.Println()
			fn()
		} else {
			fmt.Println("Invalid choice!")
		}
	}
}
