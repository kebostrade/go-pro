# ⚡ Go Concurrency Crash Course - Runnable Examples

**Companion code for the [Concurrency Crash Course](../../../docs/tutorials/CONCURRENCY_CRASH_COURSE.md)**

All examples are self-contained and ready to run!

---

## 🚀 Quick Start

```bash
# Navigate to this directory
cd basic/examples/concurrency-crash-course

# Run all examples
go run main.go

# Run with race detector (recommended!)
go run -race main.go

# Run interactive menu
# (uncomment runInteractive() in main.go)
```

---

## 📚 Available Examples

### Basic Concepts
1. **Basic Goroutines** - Launch concurrent functions
2. **Channels** - Communicate between goroutines
3. **WaitGroups** - Synchronize goroutine completion
4. **Select** - Multiplex channel operations

### Patterns
5. **Worker Pool** - Process jobs with limited workers
6. **Pipeline** - Chain processing stages
7. **Fan-Out/Fan-In** - Parallelize and merge results

### Advanced
8. **Context** - Cancellation and timeouts
9. **Mutex** - Protect shared state

### Real-World
10. **Web Scraper** - Concurrent HTTP requests
11. **Rate Limiter** - Token bucket implementation

### Pitfalls
- **Goroutine Leak** - Detect leaked goroutines
- **Race Condition** - See data races in action
- **Loop Variable Capture** - Common closure mistake

---

## 🎯 How to Use

### Method 1: Run Individual Examples

Edit `main.go` and uncomment the example you want:

```go
func main() {
    // Uncomment the example you want to run:
    example1_BasicGoroutines()
    // example2_Channels()
    // example3_WaitGroups()
    // ...
}
```

### Method 2: Interactive Menu

Uncomment `runInteractive()` in `main.go`:

```go
func main() {
    runInteractive()
}
```

Then run:
```bash
go run main.go
```

---

## 🧪 Testing with Race Detector

**Always use the race detector during development:**

```bash
# Run with race detection
go run -race main.go

# Build with race detection
go build -race

# Test with race detection
go test -race
```

The race detector will catch data races like this:

```
==================
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.pitfall2_RaceCondition.func1()
      /path/to/main.go:567 +0x3e

Previous write at 0x... by goroutine 6:
  main.pitfall2_RaceCondition.func1()
      /path/to/main.go:567 +0x54
==================
```

---

## 📖 Example Walkthrough

### Example 1: Basic Goroutines

```go
// Sequential execution
sayHello("Alice")
sayHello("Bob")

// Concurrent execution
go sayHello("Charlie")
go sayHello("Diana")
```

**Output:**
```
Hello, Alice!
Hello, Bob!
Hello, Charlie!
Hello, Diana!
```

### Example 5: Worker Pool

```go
const numJobs = 10
const numWorkers = 3

jobs := make(chan int, numJobs)
results := make(chan int, numJobs)

// Start workers
for w := 1; w <= numWorkers; w++ {
    go worker(w, jobs, results)
}

// Send jobs
for j := 1; j <= numJobs; j++ {
    jobs <- j
}
```

**Output:**
```
Worker 1 processing job 1
Worker 2 processing job 2
Worker 3 processing job 3
...
```

---

## 🔍 Debugging Tips

### 1. Check Goroutine Count

```go
import "runtime"

fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
```

### 2. Use Race Detector

```bash
go run -race main.go
```

### 3. Add Logging

```go
import "log"

log.Printf("Worker %d processing job %d", id, job)
```

### 4. Profile with pprof

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
```

Visit: `http://localhost:6060/debug/pprof/goroutine`

---

## 🎓 Learning Path

1. **Start Here**: Run examples 1-4 (basics)
2. **Patterns**: Run examples 5-7 (common patterns)
3. **Advanced**: Run examples 8-9 (context, mutex)
4. **Real-World**: Run examples 10-11 (practical applications)
5. **Pitfalls**: Run pitfall examples to see common mistakes

---

## 💡 Pro Tips

### Tip 1: Always Pass Loop Variables

```go
// ❌ BAD
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // All print 5!
    }()
}

// ✅ GOOD
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### Tip 2: Use defer with WaitGroup

```go
wg.Add(1)
go func() {
    defer wg.Done() // Always defer!
    // ... work ...
}()
```

### Tip 3: Close Channels from Sender

```go
// Sender
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch) // Sender closes
}()

// Receiver
for val := range ch {
    fmt.Println(val)
}
```

### Tip 4: Use Buffered Channels for Known Capacity

```go
// Unbuffered: Synchronous
ch := make(chan int)

// Buffered: Can send without blocking
ch := make(chan int, 100)
```

---

## 🚨 Common Mistakes

### Mistake 1: Goroutine Leak

```go
// ❌ BAD: Goroutine never exits
ch := make(chan int)
go func() {
    val := <-ch // Blocks forever
    fmt.Println(val)
}()
// Channel never receives data
```

**Fix:** Always ensure goroutines can exit.

### Mistake 2: Race Condition

```go
// ❌ BAD: Data race
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++ // RACE!
    }()
}
```

**Fix:** Use mutex or channels.

### Mistake 3: Sending on Closed Channel

```go
// ❌ BAD: Panic!
ch := make(chan int)
close(ch)
ch <- 1 // Panic!
```

**Fix:** Never send after close.

---

## 📊 Performance Benchmarks

Run benchmarks to compare approaches:

```bash
go test -bench=. -benchmem
```

Example benchmark:

```go
func BenchmarkSequential(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Sequential processing
    }
}

func BenchmarkConcurrent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Concurrent processing
    }
}
```

---

## 🔗 Related Resources

### Documentation
- [Main Crash Course](../../../docs/tutorials/CONCURRENCY_CRASH_COURSE.md)
- [Concurrency Deep Dive](../../../docs/tutorials/concurrency-deep-dive.md)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)

### Examples in This Repo
- [Basic Concurrency](../11_concurrency/)
- [Worker Pool](../../../advanced/go_18_worker_pool/)
- [Pipeline](../../../advanced/go_21_goroutines_pipeline/)
- [Fan-Out/Fan-In](../../../advanced/go_25_fan_out_fan_in/)

### External Resources
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Concurrency Patterns (Video)](https://www.youtube.com/watch?v=f6kdp27TYZs)

---

## 🤝 Contributing

Found a bug or want to add an example?

1. Fork the repository
2. Create your feature branch
3. Add your example with clear comments
4. Test with `-race` flag
5. Submit a pull request

---

## 📝 License

Part of the GO-PRO learning platform.

---

**Happy Learning! 🚀**

*Master concurrency, master Go!*

