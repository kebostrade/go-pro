# ⚡ Go Concurrency Quick Reference

**One-page cheat sheet for Go concurrency**

---

## 🚀 Goroutines

### Launch a Goroutine
```go
go func() {
    fmt.Println("Hello from goroutine")
}()
```

### With Parameters
```go
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i) // Pass i as parameter!
}
```

---

## 📡 Channels

### Create
```go
ch := make(chan int)        // Unbuffered
ch := make(chan int, 100)   // Buffered (capacity 100)
```

### Send & Receive
```go
ch <- 42        // Send
value := <-ch   // Receive
```

### Close
```go
close(ch)       // Only sender should close
```

### Range
```go
for val := range ch {
    fmt.Println(val)
}
```

### Check if Closed
```go
val, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}
```

### Directional Channels
```go
func send(ch chan<- int) {  // Send-only
    ch <- 42
}

func receive(ch <-chan int) {  // Receive-only
    val := <-ch
}
```

---

## 🔄 WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Work here
    }(i)
}

wg.Wait()  // Block until all Done()
```

---

## 🎛️ Select

### Basic Select
```go
select {
case msg := <-ch1:
    fmt.Println(msg)
case msg := <-ch2:
    fmt.Println(msg)
}
```

### With Timeout
```go
select {
case res := <-ch:
    fmt.Println(res)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

### Non-Blocking
```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("No message")
}
```

---

## 🔒 Mutex

### Exclusive Lock
```go
var mu sync.Mutex

mu.Lock()
// Critical section
mu.Unlock()
```

### Read-Write Lock
```go
var mu sync.RWMutex

// Multiple readers
mu.RLock()
// Read data
mu.RUnlock()

// Single writer
mu.Lock()
// Write data
mu.Unlock()
```

---

## 🎯 Context

### With Cancel
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go worker(ctx)
// Later...
cancel()  // Cancel all workers
```

### With Timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

### With Deadline
```go
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### With Value
```go
ctx := context.WithValue(context.Background(), "userID", 123)
userID := ctx.Value("userID").(int)
```

---

## 🏭 Common Patterns

### Worker Pool
```go
jobs := make(chan int, 100)
results := make(chan int, 100)
var wg sync.WaitGroup

// Start workers
for w := 0; w < 5; w++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for job := range jobs {
            results <- process(job)
        }
    }()
}

// Send jobs
for j := 0; j < 100; j++ {
    jobs <- j
}
close(jobs)

// Wait and close
go func() {
    wg.Wait()
    close(results)
}()

// Collect results
for result := range results {
    fmt.Println(result)
}
```

### Pipeline
```go
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

// Use
nums := generate(1, 2, 3, 4, 5)
squares := square(nums)
for s := range squares {
    fmt.Println(s)
}
```

### Fan-Out/Fan-In
```go
// Fan-out: Distribute work
func fanOut(in <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = process(in)
    }
    return channels
}

// Fan-in: Merge results
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
```

---

## 🚨 Common Pitfalls

### ❌ Loop Variable Capture
```go
// WRONG
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // All print 5!
    }()
}

// CORRECT
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### ❌ Goroutine Leak
```go
// WRONG: Goroutine never exits
ch := make(chan int)
go func() {
    val := <-ch  // Blocks forever
}()
```

### ❌ Race Condition
```go
// WRONG: Data race
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++  // RACE!
    }()
}

// CORRECT: Use mutex
var mu sync.Mutex
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### ❌ Sending on Closed Channel
```go
// WRONG: Panic!
ch := make(chan int)
close(ch)
ch <- 1  // Panic!
```

---

## 🧪 Testing

### Race Detector
```bash
go run -race main.go
go build -race
go test -race
```

### Benchmarking
```go
func BenchmarkConcurrent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Concurrent code
    }
}
```

---

## 📊 Performance Tips

### Channel Buffer Sizing
```go
// Unbuffered: Synchronous
ch := make(chan int)

// Buffered: Asynchronous
ch := make(chan int, 100)

// Rule: Buffer size = expected burst
```

### Worker Pool Sizing
```go
// CPU-bound
workers := runtime.NumCPU()

// I/O-bound
workers := runtime.NumCPU() * 2

// Network-bound
workers := 100
```

---

## ✅ Best Practices

- [ ] Use `go run -race` during development
- [ ] Pass loop variables to goroutines
- [ ] Close channels only from sender
- [ ] Use `defer wg.Done()` immediately after `wg.Add(1)`
- [ ] Use buffered channels for known capacity
- [ ] Implement graceful shutdown with context
- [ ] Prefer channels over shared state
- [ ] Use `sync.RWMutex` for read-heavy workloads
- [ ] Profile before optimizing
- [ ] Document goroutine lifecycle

---

## 🔗 Resources

- [Concurrency Crash Course](CONCURRENCY_CRASH_COURSE.md)
- [Concurrency Deep Dive](concurrency-deep-dive.md)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)

---

**Print this page for quick reference! 📄**

