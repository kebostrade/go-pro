// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides curriculum lessons 16-20
package service

import "go-pro-backend/internal/domain"

// Lesson 16: Advanced Concurrency Patterns
func (s *curriculumService) getComprehensiveLessonData16() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          16,
		Title:       "Advanced Concurrency Patterns",
		Description: "Master advanced concurrency patterns including context cancellation, worker pools, rate limiting, timeouts, and pipelines.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyAdvanced,
		Phase:       "Mastery & Production",
		Objectives: []string{
			"Use context for cancellation and deadlines",
			"Implement worker pool patterns",
			"Apply rate limiting techniques",
			"Handle timeouts properly",
			"Build concurrent pipelines",
			"Understand goroutine lifecycle management",
			"Design scalable concurrent systems",
		},
		Theory: `# Advanced Concurrency Patterns

## Context Package for Cancellation

The context package is fundamental for managing goroutine lifetimes and deadlines in Go. Context provides a way to cancel operations, set deadlines, and pass values across goroutine boundaries.

### Context Types

Go provides four main context types through constructor functions:

` + "```go" + `
// Create root context
ctx := context.Background()

// Create cancellable context
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Create context with deadline
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Create context with value (usually for request scoping)
ctx := context.WithValue(context.Background(), "userID", 123)
` + "```" + `

### Practical Cancellation Pattern

` + "```go" + `
func worker(ctx context.Context, tasks <-chan Task) {
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, gracefully exit
			fmt.Println("Worker cancelled:", ctx.Err())
			return
		case task := <-tasks:
			if err := processTask(task); err != nil {
				fmt.Printf("Task error: %v\n", err)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	tasks := make(chan Task, 10)
	for i := 0; i < 3; i++ {
		go worker(ctx, tasks)
	}

	// Send tasks
	for i := 0; i < 10; i++ {
		tasks <- Task{ID: i}
	}

	// Trigger cancellation after 5 seconds
	time.Sleep(5 * time.Second)
	cancel() // All workers will receive Done signal
	time.Sleep(1 * time.Second) // Let workers clean up
}
` + "```" + `

## Worker Pool Pattern

A worker pool maintains a fixed number of goroutines that process tasks from a shared channel. This prevents unbounded goroutine creation and provides backpressure.

` + "```go" + `
type WorkerPool struct {
	tasks    chan Task
	results  chan Result
	workers  int
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		tasks:   make(chan Task, workers*2),
		results: make(chan Result, workers*2),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.work()
	}
}

func (wp *WorkerPool) work() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			result := wp.processTask(task)
			wp.results <- result
		}
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.tasks)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	wp.wg.Wait()
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp *WorkerPool) processTask(task Task) Result {
	// Task processing logic
	time.Sleep(time.Millisecond * time.Duration(task.Duration))
	return Result{ID: task.ID, Value: task.ID * 2}
}
` + "```" + `

## Rate Limiting

Rate limiting controls the rate at which operations execute. Go's time.Ticker provides a simple rate limiting mechanism:

` + "```go" + `
// Ticker-based rate limiting
func RateLimitedLoop(limit int) {
	ticker := time.NewTicker(time.Second / time.Duration(limit))
	defer ticker.Stop()

	for i := 0; i < 100; i++ {
		<-ticker.C
		fmt.Printf("Request %d\n", i)
	}
}

// Semaphore-based rate limiting
type RateLimiter struct {
	sem chan struct{}
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (rl *RateLimiter) Do(ctx context.Context, fn func() error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case rl.sem <- struct{}{}:
		defer func() { <-rl.sem }()
		return fn()
	}
}
` + "```" + `

## Timeout Patterns

Timeouts ensure operations don't hang indefinitely:

` + "```go" + `
// Timeout for single operation
func CallWithTimeout(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// Timeout with fallback
func GetDataWithFallback(primary, fallback func() (string, error)) string {
	resultChan := make(chan string, 2)

	go func() {
		if data, err := primary(); err == nil {
			resultChan <- data
		}
	}()

	go func() {
		if data, err := fallback(); err == nil {
			resultChan <- data
		}
	}()

	select {
	case result := <-resultChan:
		return result
	case <-time.After(2 * time.Second):
		return "timeout"
	}
}
` + "```" + `

## Pipeline Pattern

Pipelines chain operations together, where each stage processes data and passes it to the next:

` + "```go" + `
// Generate numbers
func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}()
	return out
}

// Square numbers
func Square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n * n:
			}
		}
	}()
	return out
}

// Print results
func Print(ctx context.Context, in <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(n)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nums := Generate(ctx, 2, 3, 4)
	squared := Square(ctx, nums)
	Print(ctx, squared)
}
` + "```" + `

## Error Handling in Concurrent Code

` + "```go" + `
type ErrResult struct {
	Index int
	Value interface{}
	Err   error
}

func ConcurrentOperation(items []Item) ([]interface{}, error) {
	results := make([]ErrResult, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it Item) {
			defer wg.Done()
			val, err := process(it)
			results[idx] = ErrResult{
				Index: idx,
				Value: val,
				Err:   err,
			}
		}(i, item)
	}

	wg.Wait()

	// Check for errors
	for _, r := range results {
		if r.Err != nil {
			return nil, fmt.Errorf("operation %d failed: %w", r.Index, r.Err)
		}
	}

	// Extract values
	values := make([]interface{}, len(results))
	for i, r := range results {
		values[i] = r.Value
	}
	return values, nil
}
` + "```" + `
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "context-cancellation",
				Title:       "Context Cancellation",
				Description: "Implement a function that respects context cancellation.",
				Requirements: []string{
					"Create a function that accepts context",
					"Perform work in a loop",
					"Check context cancellation regularly",
					"Clean up resources on cancellation",
					"Demonstrate with timeout",
				},
				InitialCode: `package main

import (
	"context"
	"fmt"
	"time"
)

// TODO: Implement DownloadWithContext that:
// - Takes context and URL
// - Simulates download with time.Sleep
// - Checks context.Done() in loop
// - Returns error if cancelled

func main() {
	// Test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This should timeout
	err := DownloadWithContext(ctx, "https://example.com")
	fmt.Println("Result:", err)
}
`,
				Solution: `package main

import (
	"context"
	"fmt"
	"time"
)

func DownloadWithContext(ctx context.Context, url string) error {
	// Simulate download with progress
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("download cancelled: %w", ctx.Err())
		default:
			fmt.Printf("Downloaded chunk %d\n", i+1)
			time.Sleep(500 * time.Millisecond)
		}
	}
	return nil
}

func main() {
	// Test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := DownloadWithContext(ctx, "https://example.com")
	fmt.Println("Result:", err)

	// Test with manual cancellation
	fmt.Println("\nWith manual cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1500 * time.Millisecond)
		cancel2()
	}()

	err = DownloadWithContext(ctx2, "https://example.com")
	fmt.Println("Result:", err)
}
`,
			},
			{
				ID:          "worker-pool",
				Title:       "Worker Pool Implementation",
				Description: "Build a worker pool that processes tasks concurrently.",
				Requirements: []string{
					"Create WorkerPool struct",
					"Maintain fixed number of workers",
					"Process tasks from channel",
					"Track completed tasks",
					"Graceful shutdown support",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Work func() string
}

// TODO: Implement WorkerPool with:
// - tasks channel
// - results channel
// - worker goroutines
// - WaitGroup for synchronization

func main() {
	pool := NewWorkerPool(3)
	pool.Start()

	// Submit tasks
	for i := 0; i < 10; i++ {
		task := Task{
			ID: i,
			Work: func(id int) func() string {
				return func() string {
					time.Sleep(time.Second)
					return fmt.Sprintf("Task %d completed", id)
				}
			}(i),
		}
		pool.Submit(task)
	}

	// Collect results
	pool.Wait()
	for result := range pool.Results() {
		fmt.Println(result)
	}
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Work func() string
}

type WorkerPool struct {
	tasks   chan Task
	results chan string
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	wp := &WorkerPool{
		tasks:   make(chan Task, workers*2),
		results: make(chan string, workers*2),
	}

	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		result := task.Work()
		wp.results <- result
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.tasks)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan string {
	return wp.results
}

func main() {
	pool := NewWorkerPool(3)

	// Submit tasks
	for i := 0; i < 10; i++ {
		task := Task{
			ID: i,
			Work: func(id int) func() string {
				return func() string {
					time.Sleep(time.Second)
					return fmt.Sprintf("Task %d completed", id)
				}
			}(i),
		}
		pool.Submit(task)
	}

	// Collect results
	pool.Wait()
	count := 0
	for result := range pool.Results() {
		fmt.Println(result)
		count++
	}
	fmt.Printf("Processed %d tasks\n", count)
}
`,
			},
			{
				ID:          "rate-limiting",
				Title:       "Rate Limiting Implementation",
				Description: "Implement a rate limiter using semaphore pattern.",
				Requirements: []string{
					"Create RateLimiter struct",
					"Use buffered channel as semaphore",
					"Limit concurrent operations",
					"Support context cancellation",
					"Measure rate limiting effectiveness",
				},
				InitialCode: `package main

import (
	"context"
	"fmt"
	"time"
)

// TODO: Implement RateLimiter struct with:
// - semaphore channel
// - Do(ctx, fn) method for rate-limited execution

func main() {
	limiter := NewRateLimiter(2) // Max 2 concurrent ops

	start := time.Now()
	for i := 0; i < 5; i++ {
		go func(id int) {
			err := limiter.Do(context.Background(), func() error {
				fmt.Printf("Op %d started\n", id)
				time.Sleep(1 * time.Second)
				fmt.Printf("Op %d completed\n", id)
				return nil
			})
			if err != nil {
				fmt.Printf("Op %d failed: %v\n", id, err)
			}
		}(i)
	}

	time.Sleep(5 * time.Second)
	fmt.Printf("Total time: %v\n", time.Since(start))
}
`,
				Solution: `package main

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter struct {
	sem chan struct{}
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (rl *RateLimiter) Do(ctx context.Context, fn func() error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case rl.sem <- struct{}{}:
		defer func() { <-rl.sem }()
		return fn()
	}
}

func main() {
	limiter := NewRateLimiter(2) // Max 2 concurrent ops

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := limiter.Do(context.Background(), func() error {
				fmt.Printf("Op %d started at %v\n", id, time.Since(start))
				time.Sleep(1 * time.Second)
				fmt.Printf("Op %d completed at %v\n", id, time.Since(start))
				return nil
			})
			if err != nil {
				fmt.Printf("Op %d failed: %v\n", id, err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Total time: %v\n", time.Since(start))
}
`,
			},
			{
				ID:          "timeout-pattern",
				Title:       "Timeout with Fallback",
				Description: "Implement timeout pattern with fallback mechanism.",
				Requirements: []string{
					"Create function with primary and fallback sources",
					"Use select with time.After",
					"Implement race between sources",
					"Return first successful result",
					"Handle timeout gracefully",
				},
				InitialCode: `package main

import (
	"fmt"
	"math/rand"
	"time"
)

// TODO: Implement GetWithFallback that:
// - Takes primary and fallback functions
// - Returns first successful result
// - Times out after 2 seconds
// - Handles errors properly

func SlowFetch() (string, error) {
	time.Sleep(3 * time.Second)
	return "slow result", nil
}

func FastFetch() (string, error) {
	time.Sleep(500 * time.Millisecond)
	return "fast result", nil
}

func main() {
	result := GetWithFallback(SlowFetch, FastFetch)
	fmt.Println("Got:", result)
}
`,
				Solution: `package main

import (
	"fmt"
	"time"
)

func GetWithFallback(primary, fallback func() (string, error)) string {
	resultChan := make(chan string, 2)

	// Try primary source
	go func() {
		if data, err := primary(); err == nil {
			resultChan <- data
		}
	}()

	// Try fallback source
	go func() {
		if data, err := fallback(); err == nil {
			resultChan <- data
		}
	}()

	select {
	case result := <-resultChan:
		return result
	case <-time.After(2 * time.Second):
		return "timeout - both sources failed"
	}
}

func SlowFetch() (string, error) {
	time.Sleep(3 * time.Second)
	return "slow result", nil
}

func FastFetch() (string, error) {
	time.Sleep(500 * time.Millisecond)
	return "fast result", nil
}

func main() {
	fmt.Println("Test 1 - fast source wins:")
	result := GetWithFallback(SlowFetch, FastFetch)
	fmt.Println("Got:", result)

	fmt.Println("\nTest 2 - timeout:")
	result = GetWithFallback(SlowFetch, SlowFetch)
	fmt.Println("Got:", result)
}
`,
			},
			{
				ID:          "pipeline",
				Title:       "Concurrent Pipeline",
				Description: "Build a multi-stage pipeline processing data concurrently.",
				Requirements: []string{
					"Create pipeline stages as generator functions",
					"Pass data through channels",
					"Use context for cancellation",
					"Compose multiple stages",
					"Process complete pipeline",
				},
				InitialCode: `package main

import (
	"context"
	"fmt"
)

// TODO: Implement pipeline functions:
// - Generate(ctx, numbers) -> channel of numbers
// - Filter(ctx, in) -> channel of odd numbers only
// - Transform(ctx, in) -> channel of squared numbers
// - Consume(ctx, in) -> prints results

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build pipeline: Generate -> Filter -> Transform -> Consume
	nums := Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	odd := Filter(ctx, nums)
	squared := Transform(ctx, odd)
	Consume(ctx, squared)
}
`,
				Solution: `package main

import (
	"context"
	"fmt"
)

func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}()
	return out
}

func Filter(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 != 0 { // Only odd numbers
				select {
				case <-ctx.Done():
					return
				case out <- n:
				}
			}
		}
	}()
	return out
}

func Transform(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n * n: // Square the number
			}
		}
	}()
	return out
}

func Consume(ctx context.Context, in <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(n)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nums := Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	odd := Filter(ctx, nums)
	squared := Transform(ctx, odd)
	Consume(ctx, squared)
}
`,
			},
			{
				ID:          "goroutine-lifecycle",
				Title:       "Goroutine Lifecycle Management",
				Description: "Manage goroutine lifecycle with proper resource cleanup.",
				Requirements: []string{
					"Create WaitGroup-based lifecycle",
					"Implement proper cleanup",
					"Handle panics safely",
					"Monitor goroutine state",
					"Demonstrate safe shutdown",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
	"time"
)

type Manager struct {
	// TODO: Add fields for managing goroutines
}

// TODO: Implement:
// - NewManager()
// - Start() to launch workers
// - Stop() to gracefully shutdown
// - Recover from panics

func main() {
	mgr := NewManager()
	mgr.Start()

	time.Sleep(2 * time.Second)

	mgr.Stop()
	fmt.Println("All workers stopped")
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

type Manager struct {
	wg      sync.WaitGroup
	done    chan struct{}
	active  int32
	mu      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		done: make(chan struct{}),
	}
}

func (m *Manager) Start() {
	m.wg.Add(1)
	go m.worker(1)
	m.wg.Add(1)
	go m.worker(2)
	m.wg.Add(1)
	go m.worker(3)
}

func (m *Manager) worker(id int) {
	defer m.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Worker %d panicked: %v\n", id, r)
		}
	}()

	for {
		select {
		case <-m.done:
			fmt.Printf("Worker %d shutting down\n", id)
			return
		default:
			fmt.Printf("Worker %d working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (m *Manager) Stop() {
	close(m.done)
	m.wg.Wait()
}

func main() {
	mgr := NewManager()
	mgr.Start()

	time.Sleep(2 * time.Second)

	fmt.Println("Stopping manager...")
	mgr.Stop()
	fmt.Println("All workers stopped")
}
`,
			},
			{
				ID:          "error-handling-concurrent",
				Title:       "Error Handling in Concurrent Operations",
				Description: "Handle errors properly in concurrent goroutines.",
				Requirements: []string{
					"Execute multiple operations concurrently",
					"Collect errors from all goroutines",
					"Return first error or aggregate",
					"Clean up on error",
					"Provide detailed error information",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

// TODO: Implement ProcessParallel that:
// - Takes slice of items and processor function
// - Executes processor concurrently
// - Collects all errors
// - Returns aggregated errors or results

func Process(item int) (int, error) {
	// Simulate processing that may fail
	if item%3 == 0 {
		return 0, fmt.Errorf("cannot process %d", item)
	}
	time.Sleep(100 * time.Millisecond)
	return item * 2, nil
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6}
	results, errs := ProcessParallel(items, Process)

	if len(errs) > 0 {
		fmt.Println("Errors occurred:", errs)
	}
	fmt.Println("Results:", results)
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

type ErrResult struct {
	Index int
	Value int
	Err   error
}

func ProcessParallel(items []int, processor func(int) (int, error)) ([]int, []error) {
	results := make([]ErrResult, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, val int) {
			defer wg.Done()
			result, err := processor(val)
			results[idx] = ErrResult{
				Index: idx,
				Value: result,
				Err:   err,
			}
		}(i, item)
	}

	wg.Wait()

	// Extract results and errors
	values := make([]int, 0)
	errs := make([]error, 0)

	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("item %d: %w", r.Index, r.Err))
		} else {
			values = append(values, r.Value)
		}
	}

	return values, errs
}

func Process(item int) (int, error) {
	if item%3 == 0 {
		return 0, fmt.Errorf("cannot process %d", item)
	}
	time.Sleep(100 * time.Millisecond)
	return item * 2, nil
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6}
	results, errs := ProcessParallel(items, Process)

	if len(errs) > 0 {
		fmt.Println("Errors occurred:")
		for _, err := range errs {
			fmt.Println("  -", err)
		}
	}
	fmt.Println("Results:", results)
}
`,
			},
		},
	}
}

// Lesson 17: Web Development with HTTP
func (s *curriculumService) getComprehensiveLessonData17() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          17,
		Title:       "Web Development with HTTP",
		Description: "Build web applications with HTTP servers, routing, middleware, and JSON APIs.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyAdvanced,
		Phase:       "Mastery & Production",
		Objectives: []string{
			"Create HTTP servers with net/http",
			"Implement routing and handler chains",
			"Build middleware patterns",
			"Handle JSON requests and responses",
			"Serve static files",
			"Implement RESTful endpoints",
			"Handle errors and status codes properly",
		},
		Theory: `# Web Development with HTTP

## HTTP Server Basics

Go's net/http package provides a complete HTTP server implementation. Creating a server is straightforward:

` + "```go" + `
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
` + "```" + `

The http.HandlerFunc is a function signature that implements the http.Handler interface. Any function matching this signature can be a handler.

## Request Handling

Understanding the Request object is crucial:

` + "```go" + `
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// HTTP method
	fmt.Println("Method:", r.Method)

	// URL and path
	fmt.Println("Path:", r.URL.Path)
	fmt.Println("Query:", r.URL.Query())

	// Headers
	contentType := r.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	// Path parameters (with custom routing)
	// Body
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(body) // Echo back
}
` + "```" + `

## Routing Patterns

Go's default mux is basic. For production, use third-party routers, but understanding manual routing is valuable:

` + "```go" + `
type Router struct {
	routes map[string]http.Handler
}

func (r *Router) Register(path string, handler http.Handler) {
	r.routes[path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, exists := r.routes[req.URL.Path]; exists {
		handler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found")
	}
}

func main() {
	router := &Router{routes: make(map[string]http.Handler)}

	router.Register("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home page")
	}))

	http.ListenAndServe(":8080", router)
}
` + "```" + `

## Middleware Pattern

Middleware wraps handlers to add cross-cutting concerns:

` + "```go" + `
type Middleware func(http.Handler) http.Handler

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Compose middleware
func Chain(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Protected resource")
	})

	wrapped := Chain(handler, LoggingMiddleware, AuthMiddleware)
	http.ListenAndServe(":8080", wrapped)
}
` + "```" + `

## JSON APIs

` + "```go" + `
type User struct {
	ID    int    ` + "`json:\"id\"`" + `
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{ID: 1, Name: "John", Email: "john@example.com"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON: %v", err)
		return
	}

	user.ID = 2 // Simulate assignment

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/api/users/1", GetUserHandler)
	http.HandleFunc("/api/users", CreateUserHandler)
	http.ListenAndServe(":8080", nil)
}
` + "```" + `

## Static Files

` + "```go" + `
func main() {
	// Serve static files from public directory
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("public"))))

	// Also serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "public/index.html")
	})

	http.ListenAndServe(":8080", nil)
}
` + "```" + `
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "http-server-basics",
				Title:       "Basic HTTP Server",
				Description: "Create a simple HTTP server with multiple endpoints.",
				Requirements: []string{
					"Create GET /hello endpoint",
					"Create GET /user/:id endpoint",
					"Create POST /api/echo endpoint",
					"Handle JSON request/response",
					"Return appropriate status codes",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Implement handlers for:
// GET /hello - returns "Hello, World!"
// GET /user/:id - returns user info
// POST /api/echo - echoes back JSON

func main() {
	http.HandleFunc("/hello", helloHandler)
	// Register other handlers...

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, World!")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"name": "User " + id,
	})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/api/echo", echoHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
			},
			{
				ID:          "middleware-chain",
				Title:       "Middleware Implementation",
				Description: "Build a middleware chain for logging and authentication.",
				Requirements: []string{
					"Create LoggingMiddleware",
					"Create AuthMiddleware",
					"Compose middleware in chain",
					"Protect endpoints with auth",
					"Log all requests",
				},
				InitialCode: `package main

import (
	"fmt"
	"net/http"
	"time"
)

// TODO: Implement middleware functions:
// - LoggingMiddleware: logs request method and path
// - AuthMiddleware: checks Authorization header
// - Chain: composes middleware

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Secret data")
}

func main() {
	handler := http.HandlerFunc(protectedHandler)
	// Apply middleware...

	http.ListenAndServe(":8080", nil)
}
`,
				Solution: `package main

import (
	"fmt"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[%s] %s %s\n", start.Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("  Duration: %v\n", time.Since(start))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing Authorization header")
			return
		}
		if token != "Bearer valid-token" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Chain(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "Secret data")
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "Public data")
}

func main() {
	protected := Chain(http.HandlerFunc(protectedHandler), LoggingMiddleware, AuthMiddleware)
	public := Chain(http.HandlerFunc(publicHandler), LoggingMiddleware)

	http.Handle("/protected", protected)
	http.Handle("/public", public)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
			},
			{
				ID:          "json-api",
				Title:       "RESTful JSON API",
				Description: "Build a complete REST API with CRUD operations.",
				Requirements: []string{
					"Create User struct with JSON tags",
					"Implement GET /api/users endpoint",
					"Implement POST /api/users endpoint",
					"Handle JSON encoding/decoding",
					"Return appropriate status codes",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	// TODO: Add fields with JSON tags
}

var (
	users = make(map[int]User)
	mu    sync.RWMutex
	id    = 1
)

// TODO: Implement GetUsersHandler and CreateUserHandler

func main() {
	http.HandleFunc("/api/users", handleUsers)
	http.ListenAndServe(":8080", nil)
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	ID    int    ` + "`json:\"id\"`" + `
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}

var (
	users = make(map[int]User)
	mu    sync.RWMutex
	nextID = 1
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		userList := make([]User, 0, len(users))
		for _, u := range users {
			userList = append(userList, u)
		}
		mu.RUnlock()

		json.NewEncoder(w).Encode(userList)

	case http.MethodPost:
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, ` + "\"{error:\"Invalid JSON\"}\"" + `)
			return
		}

		mu.Lock()
		user.ID = nextID
		users[user.ID] = user
		nextID++
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/users", handleUsers)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
			},
			{
				ID:          "error-handling-http",
				Title:       "HTTP Error Handling",
				Description: "Implement comprehensive error handling for HTTP servers.",
				Requirements: []string{
					"Create error response wrapper",
					"Handle different error types",
					"Return appropriate status codes",
					"Provide detailed error messages",
					"Implement error middleware",
				},
				InitialCode: `package main

import (
	"fmt"
	"net/http"
)

// TODO: Create ErrorResponse struct
// TODO: Create error handler middleware
// TODO: Implement proper error status mapping

func dividHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate division that may fail
}

func main() {
	http.HandleFunc("/divide", dividHandler)
	http.ListenAndServe(":8080", nil)
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error   string ` + "`json:\"error\"`" + `
	Message string ` + "`json:\"message\"`" + `
	Status  int    ` + "`json:\"status\"`" + `
}

func errorHandler(status int, message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   http.StatusText(status),
			Message: message,
			Status:  status,
		})
	})
}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, err := strconv.Atoi(aStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "BadRequest",
			Message: "Invalid parameter a",
			Status:  http.StatusBadRequest,
		})
		return
	}

	b, err := strconv.Atoi(bStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "BadRequest",
			Message: "Invalid parameter b",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if b == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "BadRequest",
			Message: "Division by zero",
			Status:  http.StatusBadRequest,
		})
		return
	}

	result := a / b
	json.NewEncoder(w).Encode(map[string]int{"result": result})
}

func main() {
	http.HandleFunc("/divide", divideHandler)
	http.Handle("/notfound", errorHandler(http.StatusNotFound, "Endpoint not found"))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
			},
		},
	}
}

// Lesson 18: Reflection and Code Generation
func (s *curriculumService) getComprehensiveLessonData18() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          18,
		Title:       "Reflection and Code Generation",
		Description: "Master reflection for runtime type inspection and code generation with go generate.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyAdvanced,
		Phase:       "Mastery & Production",
		Objectives: []string{
			"Understand reflect package fundamentals",
			"Inspect types at runtime",
			"Work with dynamic values and calls",
			"Parse struct tags",
			"Implement custom marshallers",
			"Use go generate for code generation",
			"Build reflection-based tools",
		},
		Theory: `# Reflection and Code Generation

## Reflection Basics

Reflection allows you to inspect types and values at runtime. Go's reflect package provides the Type and Value interfaces:

` + "```go" + `
package main

import (
	"fmt"
	"reflect"
)

func InspectType(x interface{}) {
	t := reflect.TypeOf(x)
	fmt.Println("Type:", t)
	fmt.Println("Kind:", t.Kind())
}

func main() {
	InspectType(42)           // Kind: int
	InspectType("hello")      // Kind: string
	InspectType([]int{1, 2})  // Kind: slice
}
` + "```" + `

## Type and Value Inspection

` + "```go" + `
type Person struct {
	Name string
	Age  int
}

func InspectStruct(x interface{}) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	// Type information
	fmt.Println("Type:", t.Name())

	// Field information
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("  %s: %s = %v\n", field.Name, field.Type, value.Interface())
	}

	// Method information
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("  Method: %s\n", method.Name)
	}
}

func main() {
	person := Person{Name: "Alice", Age: 30}
	InspectStruct(person)
}
` + "```" + `

## Struct Tags

Struct tags provide metadata that reflection can access:

` + "```go" + `
type User struct {
	Name  string ` + "`json:\"name\" db:\"user_name\"`" + `
	Email string ` + "`json:\"email\" db:\"user_email\" validate:\"required\"`" + `
	Age   int    ` + "`json:\"age\" db:\"user_age\" validate:\"min:0,max:150\"`" + `
}

func ParseTags(x interface{}, tagKey string) {
	t := reflect.TypeOf(x)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagKey)
		fmt.Printf("%s: %s\n", field.Name, tag)
	}
}

func main() {
	user := User{}
	ParseTags(user, "json")   // name, email, age
	ParseTags(user, "db")     // user_name, user_email, user_age
	ParseTags(user, "validate") // , required, min:0,max:150
}
` + "```" + `

## Dynamic Function Calls

` + "```go" + `
func Add(a, b int) int {
	return a + b
}

func CallFunction(fn interface{}, args ...interface{}) interface{} {
	f := reflect.ValueOf(fn)

	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	results := f.Call(params)
	if len(results) > 0 {
		return results[0].Interface()
	}
	return nil
}

func main() {
	result := CallFunction(Add, 5, 3)
	fmt.Println("5 + 3 =", result) // 8
}
` + "```" + `

## Custom Marshalling

` + "```go" + `
type Marshaller interface {
	MarshalCustom() map[string]interface{}
}

func MarshalWithReflection(x interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check for json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			result[jsonTag] = value.Interface()
		} else {
			result[field.Name] = value.Interface()
		}
	}

	return result
}

func main() {
	user := User{Name: "Bob", Email: "bob@example.com"}
	marshalled := MarshalWithReflection(user)
	fmt.Println(marshalled)
}
` + "```" + `

## Code Generation with go generate

Go provides a mechanism to run generation steps during build:

` + "```go" + `
//go:generate go run gen.go

// This comment above tells go generate to run gen.go
// gen.go contains code that generates other files

package main

// The generator creates code at build time
` + "```" + `

A typical generator example:

` + "```go" + `
// gen.go
//go:build ignore

package main

import (
	"fmt"
	"os"
)

func main() {
	code := ` + "`" + `package main

const GeneratedConstant = "generated at build time"
` + "`" + `

	os.WriteFile("generated.go", []byte(code), 0644)
	fmt.Println("Generated code")
}
` + "```" + `
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "reflection-inspection",
				Title:       "Type and Value Inspection",
				Description: "Implement a function that inspects types at runtime.",
				Requirements: []string{
					"Inspect struct fields",
					"Display field types and values",
					"Handle different field kinds",
					"Format output clearly",
				},
				InitialCode: `package main

import (
	"fmt"
	"reflect"
)

type Book struct {
	Title  string
	Author string
	Pages  int
	Price  float64
}

// TODO: Implement InspectStruct that prints all field information

func main() {
	book := Book{
		Title:  "Go Programming",
		Author: "John Doe",
		Pages:  500,
		Price:  49.99,
	}

	InspectStruct(book)
}
`,
				Solution: `package main

import (
	"fmt"
	"reflect"
)

type Book struct {
	Title  string
	Author string
	Pages  int
	Price  float64
}

func InspectStruct(x interface{}) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct!")
		return
	}

	fmt.Printf("Type: %s\n", t.Name())
	fmt.Printf("Fields:\n")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("  %d. %s (%s)\n", i+1, field.Name, field.Type)
		fmt.Printf("     Value: %v\n", value.Interface())
		fmt.Printf("     Kind: %s\n", value.Kind())
	}
}

func main() {
	book := Book{
		Title:  "Go Programming",
		Author: "John Doe",
		Pages:  500,
		Price:  49.99,
	}

	InspectStruct(book)
}
`,
			},
			{
				ID:          "struct-tag-parsing",
				Title:       "Struct Tag Parsing",
				Description: "Parse struct tags for metadata.",
				Requirements: []string{
					"Define struct with custom tags",
					"Parse tags using reflection",
					"Extract tag values",
					"Handle missing tags",
				},
				InitialCode: `package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	Host     string ` + "`config:\"host\" required:\"true\"`" + `
	Port     int    ` + "`config:\"port\" required:\"true\" default:\"8080\"`" + `
	Debug    bool   ` + "`config:\"debug\" default:\"false\"`" + `
	Timeout  int    ` + "`config:\"timeout\" default:\"30\"`" + `
}

// TODO: Implement ParseConfigTags that returns tag information

type TagInfo struct {
	FieldName string
	ConfigKey string
	Required  bool
	Default   string
}

func main() {
	cfg := Config{}
	tags := ParseConfigTags(cfg)

	for _, t := range tags {
		fmt.Printf("%s: config=%s, required=%v, default=%s\n",
			t.FieldName, t.ConfigKey, t.Required, t.Default)
	}
}
`,
				Solution: `package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	Host    string ` + "`config:\"host\" required:\"true\"`" + `
	Port    int    ` + "`config:\"port\" required:\"true\" default:\"8080\"`" + `
	Debug   bool   ` + "`config:\"debug\" default:\"false\"`" + `
	Timeout int    ` + "`config:\"timeout\" default:\"30\"`" + `
}

type TagInfo struct {
	FieldName string
	ConfigKey string
	Required  bool
	Default   string
}

func ParseConfigTags(x interface{}) []TagInfo {
	t := reflect.TypeOf(x)
	var tags []TagInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		configTag := field.Tag.Get("config")
		requiredTag := field.Tag.Get("required")
		defaultTag := field.Tag.Get("default")

		required := requiredTag == "true"

		tags = append(tags, TagInfo{
			FieldName: field.Name,
			ConfigKey: configTag,
			Required:  required,
			Default:   defaultTag,
		})
	}

	return tags
}

func main() {
	cfg := Config{}
	tags := ParseConfigTags(cfg)

	for _, t := range tags {
		fmt.Printf("%s: config=%s, required=%v, default=%s\n",
			t.FieldName, t.ConfigKey, t.Required, t.Default)
	}
}
`,
			},
			{
				ID:          "dynamic-calls",
				Title:       "Dynamic Function Calls",
				Description: "Call functions dynamically using reflection.",
				Requirements: []string{
					"Accept function as interface{}",
					"Extract function signature",
					"Call with dynamic arguments",
					"Return results",
					"Handle errors",
				},
				InitialCode: `package main

import (
	"fmt"
	"reflect"
)

func Add(a, b int) int      { return a + b }
func Concat(a, b string) string { return a + b }
func Divide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// TODO: Implement DynamicCall that:
// - Takes function and arguments
// - Calls function dynamically
// - Returns result and error

func main() {
	result, err := DynamicCall(Add, 10, 5)
	fmt.Printf("Add(10, 5) = %v (err: %v)\n", result, err)

	result, err = DynamicCall(Concat, "Hello", " World")
	fmt.Printf("Concat = %v\n", result)

	result, err = DynamicCall(Divide, 10.0, 2.0)
	fmt.Printf("Divide = %v\n", result)
}
`,
				Solution: `package main

import (
	"fmt"
	"reflect"
)

func Add(a, b int) int      { return a + b }
func Concat(a, b string) string { return a + b }
func Divide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func DynamicCall(fn interface{}, args ...interface{}) (interface{}, error) {
	f := reflect.ValueOf(fn)

	if f.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}

	if len(args) != f.Type().NumIn() {
		return nil, fmt.Errorf("argument count mismatch: expected %d, got %d",
			f.Type().NumIn(), len(args))
	}

	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	results := f.Call(params)
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].Interface(), nil
}

func main() {
	result, err := DynamicCall(Add, 10, 5)
	fmt.Printf("Add(10, 5) = %v (err: %v)\n", result, err)

	result, err = DynamicCall(Concat, "Hello", " World")
	fmt.Printf("Concat = %v\n", result)

	result, err = DynamicCall(Divide, 10.0, 2.0)
	fmt.Printf("Divide = %v\n", result)

	// Error case
	result, err = DynamicCall(Add, "wrong", "types")
	fmt.Printf("Error case: %v\n", err)
}
`,
			},
			{
				ID:          "custom-marshaller",
				Title:       "Custom Marshalling",
				Description: "Implement reflection-based custom marshalling.",
				Requirements: []string{
					"Extract struct fields via reflection",
					"Use tags for field names",
					"Convert to map",
					"Handle nested structs",
				},
				InitialCode: `package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	FirstName string ` + "`marshal:\"first_name\"`" + `
	LastName  string ` + "`marshal:\"last_name\"`" + `
	Age       int    ` + "`marshal:\"age\"`" + `
}

type Address struct {
	Street string ` + "`marshal:\"street\"`" + `
	City   string ` + "`marshal:\"city\"`" + `
}

// TODO: Implement ReflectionMarshal that returns map[string]interface{}

func main() {
	person := Person{FirstName: "John", LastName: "Doe", Age: 30}
	marshalled := ReflectionMarshal(person)

	for key, val := range marshalled {
		fmt.Printf("%s: %v\n", key, val)
	}
}
`,
				Solution: `package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	FirstName string ` + "`marshal:\"first_name\"`" + `
	LastName  string ` + "`marshal:\"last_name\"`" + `
	Age       int    ` + "`marshal:\"age\"`" + `
}

type Address struct {
	Street string ` + "`marshal:\"street\"`" + `
	City   string ` + "`marshal:\"city\"`" + `
}

func ReflectionMarshal(x interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check for marshal tag
		key := field.Tag.Get("marshal")
		if key == "" {
			key = field.Name
		}

		result[key] = value.Interface()
	}

	return result
}

func main() {
	person := Person{FirstName: "John", LastName: "Doe", Age: 30}
	marshalled := ReflectionMarshal(person)

	for key, val := range marshalled {
		fmt.Printf("%s: %v\n", key, val)
	}

	fmt.Println()

	address := Address{Street: "123 Main St", City: "Springfield"}
	marshalled = ReflectionMarshal(address)

	for key, val := range marshalled {
		fmt.Printf("%s: %v\n", key, val)
	}
}
`,
			},
		},
	}
}

// Lesson 19: Performance and Profiling
func (s *curriculumService) getComprehensiveLessonData19() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          19,
		Title:       "Performance and Profiling",
		Description: "Master benchmarking, profiling, and optimization techniques in Go.",
		Duration:    "6-7 hours",
		Difficulty:  domain.DifficultyAdvanced,
		Phase:       "Mastery & Production",
		Objectives: []string{
			"Write effective benchmarks",
			"Profile CPU usage",
			"Analyze memory allocation",
			"Detect race conditions",
			"Optimize hot paths",
			"Interpret profiling results",
			"Apply optimization techniques",
		},
		Theory: `# Performance and Profiling

## Benchmarking

Go's testing package provides built-in benchmarking support. Benchmarks measure how long operations take:

` + "```go" + `
package main

import "testing"

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 5 + 3
	}
}

func BenchmarkAppend(b *testing.B) {
	var s []int
	b.ResetTimer() // Reset timer after setup

	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}
}

// Run with: go test -bench=. -benchmem
` + "```" + `

Benchmark output shows:
- BenchmarkAdd-4: runs 4 times, 1000000000 iterations, 0.63 ns/op
- Memory allocations per operation

## Memory Profiling

Memory profiling reveals allocation patterns:

` + "```go" + `
func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := "a" + "b" + "c" + "d"
		_ = s
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.WriteString("a")
		sb.WriteString("b")
		sb.WriteString("c")
		sb.WriteString("d")
		_ = sb.String()
	}
}

// BenchmarkStringConcat: 3 allocs/op
// BenchmarkStringBuilder: 1 alloc/op
` + "```" + `

## CPU Profiling

CPU profiles show where time is spent:

` + "```go" + `
import (
	"os"
	"runtime/pprof"
)

func main() {
	// Create CPU profile
	f, _ := os.Create("cpu.prof")
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Your code here
	expensiveOperation()
}

// Analyze: go tool pprof cpu.prof
` + "```" + `

## Race Detection

The race detector finds data races:

` + "```bash" + `
// Run tests with race detector
go test -race ./...

// Run program with race detector
go run -race main.go
` + "```" + `

Detects concurrent access to shared memory:

` + "```go" + `
var counter int

func main() {
	go func() {
		counter++ // Race!
	}()

	go func() {
		counter++ // Race!
	}()
}
` + "```" + `

## Optimization Techniques

### 1. Reduce Allocations

` + "```go" + `
// Bad: allocates each iteration
func Bad(n int) int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return len(result)
}

// Good: pre-allocate
func Good(n int) int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return len(result)
}
` + "```" + `

### 2. Use Efficient Data Structures

` + "```go" + `
// Use sync.Pool for temporary objects
pool := sync.NewPool(func() interface{} {
	return make([]byte, 0, 4096)
})

// Use sync.Map for concurrent access instead of mutex+map
var m sync.Map

// Use channels with buffers to reduce context switches
ch := make(chan int, 100) // Better than unbuffered
` + "```" + `

### 3. Inline-Friendly Code

` + "```go" + `
// Compiler can inline simple functions
func Add(a, b int) int {
	return a + b
}

// Large functions won't inline
func ComplexCalculation(items []int) int {
	var result int
	for _, item := range items {
		result += item * item
		result /= 2
		result -= 1
	}
	return result
}
` + "```" + `
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "benchmarking",
				Title:       "Writing Benchmarks",
				Description: "Create benchmarks to measure operation performance.",
				Requirements: []string{
					"Write benchmark functions",
					"Use b.N for iterations",
					"Report memory allocations",
					"Compare alternatives",
					"Analyze results",
				},
				InitialCode: `package main

import (
	"strings"
	"testing"
)

// TODO: Write benchmarks comparing:
// 1. String concatenation with +
// 2. String concatenation with strings.Builder
// 3. String concatenation with bytes.Buffer

func BenchmarkStringConcat(b *testing.B) {
	// Benchmark using +
}

func BenchmarkStringBuilder(b *testing.B) {
	// Benchmark using Builder
}

func BenchmarkBytesBuffer(b *testing.B) {
	// Benchmark using Buffer
}

// Run with: go test -bench=. -benchmem
`,
				Solution: `package main

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := "hello"
		s = s + " " + "world"
		_ = s
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.WriteString("hello")
		sb.WriteString(" ")
		sb.WriteString("world")
		_ = sb.String()
	}
}

func BenchmarkBytesBuffer(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.WriteString("hello")
		buf.WriteString(" ")
		buf.WriteString("world")
		_ = buf.String()
	}
}
`,
			},
			{
				ID:          "allocation-optimization",
				Title:       "Allocation Optimization",
				Description: "Optimize code to reduce memory allocations.",
				Requirements: []string{
					"Identify allocation-heavy code",
					"Pre-allocate slices",
					"Reuse buffers",
					"Measure improvement",
				},
				InitialCode: `package main

import "testing"

// Inefficient: allocates many times
func SlowSum(n int) int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}

	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

// TODO: Implement FastSum that pre-allocates

func BenchmarkSlowSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowSum(1000)
	}
}

func BenchmarkFastSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FastSum(1000)
	}
}
`,
				Solution: `package main

import "testing"

func SlowSum(n int) int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}

	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

func FastSum(n int) int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}

	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

func BenchmarkSlowSum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SlowSum(1000)
	}
}

func BenchmarkFastSum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FastSum(1000)
	}
}

// Output shows FastSum is much faster with fewer allocations
`,
			},
		},
	}
}

// Lesson 20: Building Production Applications
func (s *curriculumService) getComprehensiveLessonData20() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          20,
		Title:       "Building Production Applications",
		Description: "Master production-ready Go application development with configuration, logging, monitoring, and deployment.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyAdvanced,
		Phase:       "Mastery & Production",
		Objectives: []string{
			"Organize code for production",
			"Implement configuration management",
			"Set up logging properly",
			"Add health checks and metrics",
			"Handle graceful shutdown",
			"Prepare for deployment",
			"Build resilient systems",
		},
		Theory: `# Building Production Applications

## Project Structure

A production Go application should follow standard organization patterns:

` + "```" + `
myapp/
├── cmd/
│   └── myapp/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration loading
│   ├── handler/                 # HTTP handlers
│   ├── service/                 # Business logic
│   ├── repository/              # Data access
│   └── middleware/              # HTTP middleware
├── pkg/                         # Public, reusable packages
│   ├── logger/
│   ├── errors/
│   └── client/
├── migrations/                  # Database migrations
├── scripts/                     # Build and deployment scripts
├── Dockerfile                   # Container configuration
├── docker-compose.yml           # Local development
├── Makefile                     # Build targets
├── go.mod
├── go.sum
├── README.md
└── .env.example
` + "```" + `

## Configuration Management

Production apps need flexible configuration:

` + "```go" + `
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func LoadConfig() (*Config, error) {
	// Load from environment variables
	cfg := &Config{
		Server: ServerConfig{
			Host:    os.Getenv("SERVER_HOST"),
			Port:    parseInt(os.Getenv("SERVER_PORT")),
			Timeout: parseDuration(os.Getenv("SERVER_TIMEOUT")),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     parseInt(os.Getenv("DB_PORT")),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
	}

	return cfg, validate(cfg)
}
` + "```" + `

## Logging

Structured logging is essential for production:

` + "```go" + `
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

// Using a package like zap or logrus
type Field struct {
	Key   string
	Value interface{}
}

func main() {
	logger := NewLogger("myapp")

	logger.Info("application started",
		Field{"version", "1.0.0"},
		Field{"port", 8080})

	logger.Error("database connection failed",
		Field{"error", err},
		Field{"retry_count", 3})
}
` + "```" + `

## Graceful Shutdown

Shut down cleanly without losing requests:

` + "```go" + `
func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for shutdown signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done

		// Give requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
` + "```" + `

## Health Checks

Provide endpoints for monitoring:

` + "```go" + `
type HealthStatus struct {
	Status   string ` + "`json:\"status\"`" + `
	Uptime   int64  ` + "`json:\"uptime\"`" + `
	Database bool   ` + "`json:\"database\"`" + `
}

var startTime = time.Now()

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	dbOk := checkDatabase()

	status := "healthy"
	statusCode := http.StatusOK
	if !dbOk {
		status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(HealthStatus{
		Status:   status,
		Uptime:   int64(time.Since(startTime).Seconds()),
		Database: dbOk,
	})
}
` + "```" + `

## Deployment Best Practices

### 1. Multi-stage Docker builds

` + "```dockerfile" + `
FROM golang:1.23 AS builder
WORKDIR /build
COPY . .
RUN go build -o app ./cmd/myapp

FROM alpine:latest
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
` + "```" + `

### 2. Environment-specific configuration

` + "```bash" + `
# development
export LOG_LEVEL=debug
export DATABASE_HOST=localhost
export SERVER_PORT=8080

# production
export LOG_LEVEL=info
export DATABASE_HOST=prod-db.internal
export SERVER_PORT=8080
` + "```" + `

### 3. Monitoring and metrics

` + "```go" + `
// Record metrics during operation
func (h *Handler) recordMetrics(method string, statusCode int, duration time.Duration) {
	metricsRegistry.RecordRequest(method, statusCode, duration)
}

// Expose metrics endpoint
http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, metricsRegistry.String())
})
` + "```" + `
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "project-structure",
				Title:       "Production Project Setup",
				Description: "Set up a production-ready project structure.",
				Requirements: []string{
					"Create standard directory structure",
					"Implement configuration loading",
					"Set up basic handlers",
					"Create main entry point",
					"Document structure",
				},
				InitialCode: `// cmd/app/main.go
package main

import "fmt"

func main() {
	// TODO: Load configuration
	// TODO: Initialize database connection
	// TODO: Set up logger
	// TODO: Create router
	// TODO: Start server
	fmt.Println("Application started")
}

// This file should load config and start the application
// Create supporting files in internal/ directory
`,
				Solution: `// cmd/app/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load configuration
	cfg := LoadConfigFromEnv()

	// Initialize logger
	logger := NewLogger(cfg.LogLevel)

	// Create router
	router := SetupRouter()

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	logger.Info("Server starting", map[string]interface{}{
		"port": cfg.Port,
		"env":  cfg.Environment,
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server error", err)
	}
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Port:        getIntEnv("PORT", 8080),
		Environment: os.Getenv("ENVIRONMENT"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}
}

func getIntEnv(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return defaultVal
}

type Config struct {
	Port        int
	Environment string
	LogLevel    string
}

type Logger struct {
	level string
}

func NewLogger(level string) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s %v\n", msg, fields)
}

func (l *Logger) Fatal(msg string, err error) {
	fmt.Printf("[FATAL] %s: %v\n", msg, err)
	log.Fatal(err)
}

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Production!")
	})
	return mux
}
`,
			},
			{
				ID:          "graceful-shutdown",
				Title:       "Graceful Shutdown Implementation",
				Description: "Implement proper server shutdown with cleanup.",
				Requirements: []string{
					"Handle OS signals",
					"Complete in-flight requests",
					"Close resources gracefully",
					"Timeout long operations",
					"Report shutdown status",
				},
				InitialCode: `package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: Implement graceful shutdown that:
// - Listens for SIGINT/SIGTERM
// - Allows 30 seconds for requests to complete
// - Closes database connections
// - Logs shutdown process

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
	}

	// TODO: Implement shutdown logic

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}
`,
				Solution: `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OK")
		}),
	}

	// Channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, gracefully stopping...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
`,
			},
			{
				ID:          "health-checks",
				Title:       "Health Check Endpoints",
				Description: "Add monitoring endpoints for application health.",
				Requirements: []string{
					"Create /health endpoint",
					"Check database connectivity",
					"Track uptime",
					"Return appropriate status codes",
					"Include detailed status info",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status     string ` + "`json:\"status\"`" + `
	Timestamp  string ` + "`json:\"timestamp\"`" + `
	Uptime     int64  ` + "`json:\"uptime_seconds\"`" + `
	Database   bool   ` + "`json:\"database\"`" + `
	Memory     uint64 ` + "`json:\"memory_mb\"`" + `
}

// TODO: Implement HealthHandler
// TODO: Track start time
// TODO: Check database status

var startTime = time.Now()

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Return health status as JSON
}

func main() {
	http.HandleFunc("/health", HealthHandler)
	http.ListenAndServe(":8080", nil)
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type HealthStatus struct {
	Status     string ` + "`json:\"status\"`" + `
	Timestamp  string ` + "`json:\"timestamp\"`" + `
	Uptime     int64  ` + "`json:\"uptime_seconds\"`" + `
	Database   bool   ` + "`json:\"database\"`" + `
	Memory     uint64 ` + "`json:\"memory_mb\"`" + `
}

var startTime = time.Now()

func CheckDatabase() bool {
	// Simulate database check
	// In real app, would ping database
	return true
}

func GetMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc / 1024 / 1024 // Convert to MB
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	dbOk := CheckDatabase()
	uptime := int64(time.Since(startTime).Seconds())
	memory := GetMemoryUsage()

	status := "healthy"
	statusCode := http.StatusOK

	if !dbOk {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	health := HealthStatus{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    uptime,
		Database:  dbOk,
		Memory:    memory,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

func main() {
	http.HandleFunc("/health", HealthHandler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
`,
			},
		},
	}
}
