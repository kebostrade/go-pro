# ⚡ Go Concurrency Crash Course

**Duration:** 60-90 minutes | **Level:** Intermediate | **Hands-On:** 100%

A fast-paced, practical guide to mastering Go concurrency. Learn by doing with real-world examples you can run immediately.

---

## 🎯 What You'll Learn

In the next hour, you'll master:
- ✅ Goroutines and how they work
- ✅ Channels for communication
- ✅ Common concurrency patterns
- ✅ Avoiding deadlocks and race conditions
- ✅ Real-world applications

---

## 🚀 Quick Setup

```bash
# Create a workspace
mkdir go-concurrency-crash && cd go-concurrency-crash
go mod init crash

# You're ready! Copy and run examples below
```

---

## 📖 Table of Contents

1. [Goroutines: Lightweight Threads](#1-goroutines-lightweight-threads)
2. [Channels: Communication](#2-channels-communication)
3. [WaitGroups: Synchronization](#3-waitgroups-synchronization)
4. [Select: Multiplexing](#4-select-multiplexing)
5. [Patterns: Worker Pool](#5-patterns-worker-pool)
6. [Patterns: Pipeline](#6-patterns-pipeline)
7. [Patterns: Fan-Out/Fan-In](#7-patterns-fan-outfan-in)
8. [Context: Cancellation](#8-context-cancellation)
9. [Mutex: Shared State](#9-mutex-shared-state)
10. [Common Pitfalls](#10-common-pitfalls)

---

## 1. Goroutines: Lightweight Threads

### The Basics

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Sequential execution
    sayHello("Alice")
    sayHello("Bob")
    
    // Concurrent execution
    go sayHello("Charlie")
    go sayHello("Diana")
    
    // Wait for goroutines (bad way - we'll fix this)
    time.Sleep(time.Second)
}

func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
```

**Output:**
```
Hello, Alice!
Hello, Bob!
Hello, Charlie!
Hello, Diana!
```

### Key Points
- `go` keyword launches a goroutine
- Goroutines are cheap (2KB stack)
- Main goroutine exit = program exit
- Never use `time.Sleep` for synchronization in production!

---

## 2. Channels: Communication

### Unbuffered Channels

```go
package main

import "fmt"

func main() {
    // Create channel
    ch := make(chan string)
    
    // Send in goroutine (must be concurrent!)
    go func() {
        ch <- "Hello from goroutine"
    }()
    
    // Receive (blocks until data available)
    msg := <-ch
    fmt.Println(msg)
}
```

### Buffered Channels

```go
package main

import "fmt"

func main() {
    // Buffered channel (capacity 2)
    ch := make(chan int, 2)
    
    // Can send without blocking (up to capacity)
    ch <- 1
    ch <- 2
    
    // Receive
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
}
```

### Channel Directions

```go
// Send-only channel
func sendOnly(ch chan<- int) {
    ch <- 42
}

// Receive-only channel
func receiveOnly(ch <-chan int) {
    val := <-ch
    fmt.Println(val)
}

func main() {
    ch := make(chan int)
    go sendOnly(ch)
    receiveOnly(ch)
}
```

### Closing Channels

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 3)
    
    // Send values
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch) // Signal no more values
    
    // Range over channel (stops when closed)
    for val := range ch {
        fmt.Println(val)
    }
}
```

**Rule:** Only the sender should close a channel.

---

## 3. WaitGroups: Synchronization

### The Right Way to Wait

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1) // Increment counter
        
        go func(id int) {
            defer wg.Done() // Decrement when done
            
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Second)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    
    wg.Wait() // Block until counter is 0
    fmt.Println("All workers completed")
}
```

**Critical:** Pass loop variable as parameter to avoid closure issues!

---

## 4. Select: Multiplexing

### Waiting on Multiple Channels

```go
package main

import (
    "fmt"
    "time"
)

func main() {
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
            fmt.Println("Received", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received", msg2)
        }
    }
}
```

### Timeout Pattern

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case res := <-ch:
        fmt.Println("Got:", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

### Non-Blocking Operations

```go
select {
case msg := <-ch:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message available")
}
```

---

## 🎓 Real-World Examples

### Example 1: Concurrent Web Scraper

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

type Result struct {
    URL    string
    Status int
    Error  error
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

    results <- Result{URL: url, Status: resp.StatusCode}
}

func main() {
    urls := []string{
        "https://golang.org",
        "https://github.com",
        "https://stackoverflow.com",
        "https://reddit.com",
        "https://news.ycombinator.com",
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
    for result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.URL, result.Error)
        } else {
            fmt.Printf("✅ %s: %d\n", result.URL, result.Status)
        }
    }
}
```

### Example 2: Rate-Limited API Client

```go
package main

import (
    "fmt"
    "time"
)

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

func main() {
    limiter := NewRateLimiter(2) // 2 requests per second

    for i := 1; i <= 10; i++ {
        limiter.Wait()
        go func(id int) {
            fmt.Printf("Request %d at %s\n", id, time.Now().Format("15:04:05"))
        }(i)
    }

    time.Sleep(6 * time.Second)
}
```

### Example 3: Concurrent File Processor

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileInfo struct {
    Path string
    Size int64
}

func processFiles(root string, workers int) ([]FileInfo, error) {
    files := make(chan string, 100)
    results := make(chan FileInfo, 100)
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range files {
                info, err := os.Stat(path)
                if err == nil {
                    results <- FileInfo{Path: path, Size: info.Size()}
                }
            }
        }()
    }

    // Walk directory
    go func() {
        filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err == nil && !info.IsDir() {
                files <- path
            }
            return nil
        })
        close(files)
    }()

    // Collect results
    go func() {
        wg.Wait()
        close(results)
    }()

    var fileInfos []FileInfo
    for info := range results {
        fileInfos = append(fileInfos, info)
    }

    return fileInfos, nil
}

func main() {
    files, err := processFiles(".", 4)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    var totalSize int64
    for _, f := range files {
        totalSize += f.Size
    }

    fmt.Printf("Processed %d files, total size: %d bytes\n", len(files), totalSize)
}
```

---

## 🧪 Practice Exercises

### Exercise 1: Parallel Sum

Calculate sum of numbers using multiple goroutines.

```go
package main

import (
    "fmt"
    "sync"
)

func parallelSum(numbers []int, workers int) int {
    // TODO: Implement parallel sum
    // Hint: Divide numbers among workers
    // Each worker sums its portion
    // Combine results

    return 0
}

func main() {
    numbers := make([]int, 1000000)
    for i := range numbers {
        numbers[i] = i + 1
    }

    result := parallelSum(numbers, 4)
    fmt.Println("Sum:", result)
}
```

<details>
<summary>Solution</summary>

```go
func parallelSum(numbers []int, workers int) int {
    chunkSize := len(numbers) / workers
    results := make(chan int, workers)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        start := i * chunkSize
        end := start + chunkSize
        if i == workers-1 {
            end = len(numbers)
        }

        go func(nums []int) {
            defer wg.Done()
            sum := 0
            for _, n := range nums {
                sum += n
            }
            results <- sum
        }(numbers[start:end])
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    total := 0
    for sum := range results {
        total += sum
    }

    return total
}
```
</details>

### Exercise 2: Concurrent Cache

Implement a thread-safe cache with expiration.

```go
package main

import (
    "sync"
    "time"
)

type Cache struct {
    // TODO: Add fields
}

func NewCache() *Cache {
    // TODO: Implement
    return &Cache{}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    // TODO: Implement
}

func (c *Cache) Get(key string) (interface{}, bool) {
    // TODO: Implement
    return nil, false
}

func main() {
    cache := NewCache()
    cache.Set("user:1", "Alice", 5*time.Second)

    if val, ok := cache.Get("user:1"); ok {
        println("Found:", val.(string))
    }
}
```

<details>
<summary>Solution</summary>

```go
type item struct {
    value      interface{}
    expiration time.Time
}

type Cache struct {
    mu    sync.RWMutex
    items map[string]item
}

func NewCache() *Cache {
    c := &Cache{
        items: make(map[string]item),
    }

    // Cleanup expired items
    go func() {
        ticker := time.NewTicker(time.Minute)
        defer ticker.Stop()

        for range ticker.C {
            c.mu.Lock()
            now := time.Now()
            for key, item := range c.items {
                if now.After(item.expiration) {
                    delete(c.items, key)
                }
            }
            c.mu.Unlock()
        }
    }()

    return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = item{
        value:      value,
        expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, ok := c.items[key]
    if !ok || time.Now().After(item.expiration) {
        return nil, false
    }

    return item.value, true
}
```
</details>

### Exercise 3: Job Queue

Build a job queue with priority.

```go
package main

import (
    "fmt"
    "sync"
)

type Job struct {
    ID       int
    Priority int
    Task     func()
}

type JobQueue struct {
    // TODO: Implement priority queue
}

func NewJobQueue(workers int) *JobQueue {
    // TODO: Implement
    return &JobQueue{}
}

func (jq *JobQueue) Submit(job Job) {
    // TODO: Implement
}

func (jq *JobQueue) Start() {
    // TODO: Implement
}

func main() {
    queue := NewJobQueue(3)
    queue.Start()

    for i := 1; i <= 10; i++ {
        id := i
        queue.Submit(Job{
            ID:       id,
            Priority: id % 3,
            Task: func() {
                fmt.Printf("Processing job %d\n", id)
            },
        })
    }
}
```

---

## 🔍 Debugging Concurrency

### 1. Race Detector

```bash
# Build with race detector
go build -race

# Run with race detector
go run -race main.go

# Test with race detector
go test -race ./...
```

### 2. Deadlock Detection

Go runtime detects deadlocks automatically:

```go
package main

func main() {
    ch := make(chan int)
    ch <- 1 // Deadlock! No receiver
}
```

**Output:**
```
fatal error: all goroutines are asleep - deadlock!
```

### 3. Profiling Goroutines

```go
package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "runtime"
)

func main() {
    // Start pprof server
    go http.ListenAndServe("localhost:6060", nil)

    // Your concurrent code here

    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}
```

Visit `http://localhost:6060/debug/pprof/goroutine` to see goroutine profiles.

---

## 📊 Performance Tips

### 1. Channel Buffer Sizing

```go
// Unbuffered: Synchronous, slower
ch := make(chan int)

// Buffered: Asynchronous, faster
ch := make(chan int, 100)

// Rule of thumb: Buffer size = expected burst
```

### 2. Worker Pool Sizing

```go
// CPU-bound tasks
workers := runtime.NumCPU()

// I/O-bound tasks
workers := runtime.NumCPU() * 2

// Network-bound tasks
workers := 100 // or more
```

### 3. Avoid Goroutine Overhead

```go
// BAD: Too many goroutines
for i := 0; i < 1000000; i++ {
    go process(i)
}

// GOOD: Worker pool
jobs := make(chan int, 1000000)
for w := 0; w < 100; w++ {
    go worker(jobs)
}
for i := 0; i < 1000000; i++ {
    jobs <- i
}
```

---

## 🎯 Quick Reference

### Channel Operations

| Operation | Syntax | Blocks? |
|-----------|--------|---------|
| Create | `ch := make(chan T)` | No |
| Send | `ch <- value` | Yes (unbuffered) |
| Receive | `value := <-ch` | Yes (until data) |
| Close | `close(ch)` | No |
| Range | `for v := range ch` | Yes (until closed) |

### Synchronization Primitives

| Type | Use Case | Example |
|------|----------|---------|
| `sync.WaitGroup` | Wait for goroutines | `wg.Add(1); wg.Done(); wg.Wait()` |
| `sync.Mutex` | Exclusive access | `mu.Lock(); mu.Unlock()` |
| `sync.RWMutex` | Read/write lock | `mu.RLock(); mu.RUnlock()` |
| `sync.Once` | Run once | `once.Do(func() {...})` |
| `context.Context` | Cancellation | `ctx.Done()` |

### Common Patterns

```go
// Worker Pool
for w := 0; w < workers; w++ {
    go worker(jobs, results)
}

// Pipeline
stage1 := generate(data)
stage2 := process(stage1)
stage3 := output(stage2)

// Fan-Out/Fan-In
workers := fanOut(input, 5)
results := fanIn(workers...)

// Timeout
select {
case <-ch:
case <-time.After(timeout):
}

// Rate Limiting
limiter := time.Tick(time.Second / rate)
<-limiter
```

---

## ✅ Best Practices Checklist

- [ ] Use `go run -race` during development
- [ ] Always pass loop variables to goroutines
- [ ] Close channels only from sender
- [ ] Use `defer wg.Done()` immediately after `wg.Add(1)`
- [ ] Use buffered channels for known capacity
- [ ] Implement graceful shutdown with context
- [ ] Avoid shared state; prefer channels
- [ ] Use `sync.RWMutex` for read-heavy workloads
- [ ] Profile before optimizing
- [ ] Document goroutine lifecycle

---

## 🚀 Next Steps

### Beginner → Intermediate
1. ✅ Complete this crash course
2. 📖 Read [Concurrency Deep Dive](concurrency-deep-dive.md)
3. 🔨 Build: [Concurrent Web Crawler](../../basic/projects/)
4. 📚 Study: [Go Memory Model](https://go.dev/ref/mem)

### Intermediate → Advanced
1. 📖 Read: [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
2. 🎥 Watch: [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
3. 🔨 Build: Distributed task queue
4. 📚 Study: [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)

### Production-Ready
1. 📖 Implement observability (OpenTelemetry)
2. 🔨 Add circuit breakers and retries
3. 📊 Profile and optimize
4. 🧪 Load test concurrent systems

---

## 📚 Additional Resources

### Official Documentation
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go Blog - Concurrency](https://go.dev/blog/pipelines)

### Books
- "Concurrency in Go" by Katherine Cox-Buday
- "The Go Programming Language" by Donovan & Kernighan (Chapter 8-9)

### Videos
- [Google I/O 2012 - Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Google I/O 2013 - Advanced Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)

### Code Examples
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Concurrency Patterns](https://github.com/lotusirous/go-concurrency-patterns)

---

## 🎓 Summary

**You've learned:**
- ✅ Goroutines are lightweight threads
- ✅ Channels enable safe communication
- ✅ WaitGroups synchronize goroutines
- ✅ Select multiplexes channels
- ✅ Common patterns: Worker Pool, Pipeline, Fan-Out/Fan-In
- ✅ Context for cancellation
- ✅ Mutex for shared state
- ✅ How to avoid common pitfalls

**Key Takeaways:**
1. **"Don't communicate by sharing memory; share memory by communicating"**
2. **Always use `-race` flag during development**
3. **Channels are first-class citizens in Go**
4. **Goroutines are cheap, but not free**
5. **Context is your friend for cancellation**

---

## 🤝 Contributing

Found an issue or want to add an example? Contributions welcome!

1. Fork the repository
2. Create your feature branch
3. Add examples with tests
4. Submit a pull request

---

## 📝 License

This crash course is part of the GO-PRO learning platform.

---

**Happy Concurrent Programming! 🚀**

*Master concurrency, master Go!*

## 5. Patterns: Worker Pool

**Problem:** Process many tasks with limited workers.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Second) // Simulate work
        results <- job * 2
    }
}

func main() {
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
    for result := range results {
        fmt.Println("Result:", result)
    }
}
```

**Use Case:** API rate limiting, batch processing, web scraping.

---

## 6. Patterns: Pipeline

**Problem:** Chain processing stages.

```go
package main

import "fmt"

// Stage 1: Generate numbers
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

// Stage 2: Square numbers
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

// Stage 3: Print results
func main() {
    // Chain stages
    numbers := generate(1, 2, 3, 4, 5)
    squares := square(numbers)
    
    // Consume
    for result := range squares {
        fmt.Println(result)
    }
}
```

**Output:** `1, 4, 9, 16, 25`

---

## 7. Patterns: Fan-Out/Fan-In

**Problem:** Parallelize work, then merge results.

```go
package main

import (
    "fmt"
    "sync"
)

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
            out <- n * n // Process
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

func main() {
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
    for result := range results {
        fmt.Println(result)
    }
}
```

---

## 8. Context: Cancellation

### Basic Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx, 1)
    go worker(ctx, 2)
    
    time.Sleep(2 * time.Second)
    cancel() // Cancel all workers
    
    time.Sleep(time.Second) // Let them finish
}
```

### Timeout Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Operation completed")
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    }
}
```

---

## 9. Mutex: Shared State

### The Problem: Race Condition

```go
// ❌ BROKEN CODE
package main

import (
    "fmt"
    "sync"
)

func main() {
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
    fmt.Println("Counter:", counter) // Not 1000!
}
```

### The Solution: Mutex

```go
// ✅ CORRECT CODE
package main

import (
    "fmt"
    "sync"
)

func main() {
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
    fmt.Println("Counter:", counter) // Always 1000
}
```

### Detect Race Conditions

```bash
go run -race main.go
```

---

## 10. Common Pitfalls

### ❌ Pitfall 1: Goroutine Leaks

```go
// BAD: Goroutine never exits
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // Blocks forever
        fmt.Println(val)
    }()
    // Channel never receives data
}
```

**Fix:** Always ensure goroutines can exit.

### ❌ Pitfall 2: Closing Receive-Only Channel

```go
// BAD: Receiver closes channel
func bad(ch <-chan int) {
    close(ch) // Panic!
}
```

**Fix:** Only sender closes channels.

### ❌ Pitfall 3: Loop Variable Capture

```go
// BAD
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // All print 5!
    }()
}

// GOOD
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### ❌ Pitfall 4: Sending on Closed Channel

```go
ch := make(chan int)
close(ch)
ch <- 1 // Panic!
```

**Fix:** Never send after close.

---


