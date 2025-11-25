# 🔄 Concurrency Patterns in Go

A comprehensive guide to Go's concurrency primitives and patterns with practical examples and real-world use cases.

## 📚 Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Concurrency Primitives](#concurrency-primitives)
- [Concurrency Patterns](#concurrency-patterns)
- [Real-World Examples](#real-world-examples)
- [Best Practices](#best-practices)
- [Common Pitfalls](#common-pitfalls)

## 🎯 Overview

Go's concurrency model is based on **CSP (Communicating Sequential Processes)**. This project demonstrates:

- ✅ Goroutines - Lightweight threads
- ✅ Channels - Communication between goroutines
- ✅ Sync primitives - Mutexes, WaitGroups, etc.
- ✅ Concurrency patterns - Worker pools, pipelines, fan-out/fan-in
- ✅ Real-world applications - Web scrapers, parallel processors

## 🚀 Quick Start

### Installation

```bash
# Navigate to project
cd basic/projects/concurrency-patterns

# Download dependencies
make deps

# Run all examples
make run-all

# Run tests
make test
```

### Run Individual Examples

```bash
# Goroutine basics
make run-goroutines

# Channel patterns
make run-channels

# Worker pool
make run-worker-pool

# Pipeline
make run-pipeline
```

## 🔧 Concurrency Primitives

### 1. Goroutines

**Lightweight threads managed by Go runtime**

```go
// Launch goroutine
go func() {
    fmt.Println("Hello from goroutine")
}()

// With WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // Do work
}()
wg.Wait()
```

**Key Features**:
- Extremely lightweight (2KB stack)
- Managed by Go scheduler
- Can have thousands concurrently

### 2. Channels

**Typed conduits for communication**

```go
// Unbuffered channel
ch := make(chan int)

// Buffered channel
ch := make(chan int, 10)

// Send and receive
ch <- 42        // Send
value := <-ch   // Receive

// Close channel
close(ch)

// Range over channel
for val := range ch {
    fmt.Println(val)
}
```

**Channel Types**:
- **Unbuffered**: Synchronous (blocks until both ready)
- **Buffered**: Asynchronous (blocks when full/empty)
- **Send-only**: `chan<- T`
- **Receive-only**: `<-chan T`

### 3. Select Statement

**Wait on multiple channels**

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("No value ready")
}
```

### 4. Sync Primitives

**Synchronization tools**

```go
// WaitGroup - Wait for goroutines
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // Work
}()
wg.Wait()

// Mutex - Mutual exclusion
var mu sync.Mutex
mu.Lock()
// Critical section
mu.Unlock()

// RWMutex - Read-write lock
var rwmu sync.RWMutex
rwmu.RLock()  // Read lock
rwmu.RUnlock()
rwmu.Lock()   // Write lock
rwmu.Unlock()

// Once - Execute once
var once sync.Once
once.Do(func() {
    // Runs only once
})
```

## 🎨 Concurrency Patterns

### 1. Worker Pool

**Fixed number of workers processing jobs**

```go
pool := NewWorkerPool(numWorkers, queueSize)
pool.Start()

// Submit jobs
pool.Submit(Job{ID: 1, Data: data, Process: processFunc})

// Collect results
for result := range pool.Results() {
    fmt.Println(result)
}

pool.Close()
```

**Use Cases**:
- Image processing
- Data transformation
- API requests
- File processing

### 2. Pipeline

**Chain of processing stages**

```go
// Stage 1: Generate
numbers := generate(1, 2, 3, 4, 5)

// Stage 2: Transform
squares := square(numbers)

// Stage 3: Aggregate
sum := sum(squares)
```

**Use Cases**:
- Data processing pipelines
- Stream processing
- ETL operations

### 3. Fan-Out/Fan-In

**Distribute work, merge results**

```go
// Fan-out: Multiple workers
workers := make([]<-chan int, numWorkers)
for i := 0; i < numWorkers; i++ {
    workers[i] = worker(input)
}

// Fan-in: Merge results
results := merge(workers...)
```

**Use Cases**:
- Parallel processing
- Distributed computation
- Load balancing

### 4. Rate Limiting

**Control request rate**

```go
limiter := time.Tick(500 * time.Millisecond)

for req := range requests {
    <-limiter // Wait for rate limiter
    process(req)
}
```

**Use Cases**:
- API rate limiting
- Resource throttling
- Backpressure control

### 5. Semaphore

**Limit concurrent operations**

```go
semaphore := make(chan struct{}, maxConcurrent)

semaphore <- struct{}{}  // Acquire
// Do work
<-semaphore              // Release
```

**Use Cases**:
- Connection pooling
- Resource limits
- Concurrency control

## 💻 Real-World Examples

### Web Scraper

```go
// Concurrent web scraping with worker pool
scraper := NewScraper(numWorkers)
scraper.Start()

urls := []string{"url1", "url2", "url3"}
for _, url := range urls {
    scraper.Submit(url)
}

results := scraper.Results()
```

### Parallel Processor

```go
// Process large dataset in parallel
processor := NewParallelProcessor(numWorkers)
results := processor.Process(data)
```

## ✅ Best Practices

### 1. Always Close Channels

```go
// ✅ Good
go func() {
    defer close(ch)
    for _, item := range items {
        ch <- item
    }
}()

// ❌ Bad - channel never closed
go func() {
    for _, item := range items {
        ch <- item
    }
}()
```

### 2. Use WaitGroup

```go
// ✅ Good
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // Work
    }()
}
wg.Wait()

// ❌ Bad - no synchronization
for i := 0; i < 10; i++ {
    go func() {
        // Work
    }()
}
```

### 3. Pass Variables to Goroutines

```go
// ✅ Good
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}

// ❌ Bad - closure captures loop variable
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i) // May print 10, 10, 10...
    }()
}
```

### 4. Use Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case result := <-ch:
    return result
case <-ctx.Done():
    return ctx.Err()
}
```

## ⚠️ Common Pitfalls

### 1. Goroutine Leaks

```go
// ❌ Goroutine never exits
go func() {
    for {
        // No exit condition
    }
}()

// ✅ Use context for cancellation
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Work
        }
    }
}(ctx)
```

### 2. Race Conditions

```go
// ❌ Race condition
var counter int
for i := 0; i < 100; i++ {
    go func() {
        counter++ // Unsafe
    }()
}

// ✅ Use mutex or atomic
var mu sync.Mutex
for i := 0; i < 100; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### 3. Deadlocks

```go
// ❌ Deadlock - unbuffered channel
ch := make(chan int)
ch <- 42 // Blocks forever

// ✅ Use goroutine or buffered channel
go func() {
    ch <- 42
}()
```

## 📊 Performance Tips

1. **Reuse goroutines** with worker pools
2. **Use buffered channels** to reduce blocking
3. **Limit concurrency** with semaphores
4. **Profile with pprof** to find bottlenecks
5. **Use sync.Pool** for object reuse

---

**Built with ❤️ using Go's powerful concurrency primitives**

