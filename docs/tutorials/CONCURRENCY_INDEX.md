# 🔄 Go Concurrency - Complete Index

**Your one-stop guide to all concurrency materials in GO-PRO**

---

## 🎯 Choose Your Path

### 🚀 I want to get started quickly (60-90 min)
→ **[Concurrency Crash Course](CONCURRENCY_CRASH_COURSE.md)**

### 📚 I want deep understanding (4-5 hours)
→ **[Concurrency Deep Dive](concurrency-deep-dive.md)**

### 📄 I need a quick reference
→ **[Quick Reference](CONCURRENCY_QUICK_REFERENCE.md)**

### 🗺️ I want to see all resources
→ **[Resources Guide](CONCURRENCY_RESOURCES.md)**

---

## 📖 All Materials

### Tutorials

| Resource | Duration | Level | Type |
|----------|----------|-------|------|
| [Crash Course](CONCURRENCY_CRASH_COURSE.md) | 60-90 min | Intermediate | Hands-on |
| [Deep Dive](concurrency-deep-dive.md) | 4-5 hours | Advanced | Comprehensive |
| [Quick Reference](CONCURRENCY_QUICK_REFERENCE.md) | 5 min | All | Cheat sheet |
| [Resources Guide](CONCURRENCY_RESOURCES.md) | 10 min | All | Navigation |

---

### Runnable Examples

| Location | Examples | Description |
|----------|----------|-------------|
| `basic/examples/concurrency-crash-course/` | 11 | Crash course examples |
| `basic/examples/11_concurrency/` | 10 | Basic patterns |
| `advanced/go_18_worker_pool/` | 1 | Worker pool |
| `advanced/go_21_goroutines_pipeline/` | 1 | Pipeline |
| `advanced/go_25_fan_out_fan_in/` | 1 | Fan-out/fan-in |
| `advanced/go_24_unbuffered_buffered_channels/` | 2 | Channels |
| `advanced/go_51_context/` | 4 | Context |
| `advanced/go_20_mutexes_and_confinement/` | 3 | Mutexes |

---

## 🎓 Learning Paths

### Beginner Path (2-3 hours)
1. **Read:** Crash Course (60 min)
2. **Code:** Run all crash course examples (30 min)
3. **Practice:** Complete 2 exercises (60 min)
4. **Reference:** Review quick reference (15 min)

### Intermediate Path (1 week)
1. **Day 1:** Crash Course + examples
2. **Day 2:** Basic concurrency examples
3. **Day 3-4:** Deep Dive tutorial
4. **Day 5:** Advanced examples
5. **Day 6-7:** Build a project

### Advanced Path (2 weeks)
1. **Week 1:** All tutorials + examples
2. **Week 2:** Production patterns + projects

---

## 📚 By Topic

### Goroutines
- **Tutorial:** Crash Course Section 1, Deep Dive Section 1
- **Examples:** `basic/examples/concurrency-crash-course/` (Example 1)
- **Reference:** Quick Reference - Goroutines

### Channels
- **Tutorial:** Crash Course Section 2, Deep Dive Section 2
- **Examples:** `advanced/go_24_unbuffered_buffered_channels/`
- **Reference:** Quick Reference - Channels

### Patterns
- **Tutorial:** Crash Course Sections 5-7, Deep Dive Section 5
- **Examples:** `advanced/go_18_worker_pool/`, `advanced/go_21_goroutines_pipeline/`
- **Reference:** Quick Reference - Common Patterns

### Synchronization
- **Tutorial:** Crash Course Sections 3, 9, Deep Dive Sections 3-4
- **Examples:** `advanced/go_20_mutexes_and_confinement/`
- **Reference:** Quick Reference - Mutex, WaitGroup

### Context
- **Tutorial:** Crash Course Section 8, Deep Dive Section 5
- **Examples:** `advanced/go_51_context/`
- **Reference:** Quick Reference - Context

---

## 🚀 Quick Start Commands

### Run Crash Course Examples
```bash
cd basic/examples/concurrency-crash-course
go run main.go
```

### Run with Race Detector
```bash
go run -race main.go
```

### Run Tests
```bash
./test.sh
```

### Run Specific Example
```bash
# Edit main.go and uncomment the example you want
go run main.go
```

---

## 🧪 Testing & Debugging

### Race Detection
```bash
go run -race main.go
go build -race
go test -race ./...
```

### Profiling
```bash
go test -cpuprofile=cpu.prof -bench=.
go test -memprofile=mem.prof -bench=.
```

### Benchmarking
```bash
go test -bench=. -benchmem ./...
```

---

## 📊 Content Overview

### Total Materials
- **Tutorials:** 4 comprehensive guides
- **Examples:** 30+ runnable programs
- **Code Lines:** 2500+ lines
- **Patterns:** 7 production-ready patterns
- **Exercises:** 10+ practice problems

### Coverage
- ✅ Goroutines and scheduling
- ✅ Channels (buffered, unbuffered, directional)
- ✅ WaitGroups and synchronization
- ✅ Select and multiplexing
- ✅ Worker pools
- ✅ Pipelines
- ✅ Fan-out/fan-in
- ✅ Context and cancellation
- ✅ Mutexes and locks
- ✅ Race conditions
- ✅ Deadlock prevention
- ✅ Performance optimization

---

## 🎯 Recommended Sequence

### For Complete Beginners
1. Start: [Crash Course](CONCURRENCY_CRASH_COURSE.md)
2. Practice: Run all examples
3. Reference: [Quick Reference](CONCURRENCY_QUICK_REFERENCE.md)
4. Next: Build a simple project

### For Experienced Developers
1. Skim: [Crash Course](CONCURRENCY_CRASH_COURSE.md)
2. Deep dive: [Deep Dive](concurrency-deep-dive.md)
3. Study: Advanced examples
4. Build: Production system

### For Production Engineers
1. Review: All tutorials
2. Study: Race conditions and debugging
3. Practice: Profiling and optimization
4. Implement: Production patterns

---

## 🔗 External Resources

### Official
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go Blog - Concurrency](https://go.dev/blog/pipelines)

### Videos
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Advanced Concurrency](https://www.youtube.com/watch?v=QDDwwePbDtw)

### Books
- "Concurrency in Go" by Katherine Cox-Buday
- "The Go Programming Language" by Donovan & Kernighan

---

## 📝 Quick Reference

### Essential Commands
```bash
# Run with race detector
go run -race main.go

# Build with race detector
go build -race

# Test with race detector
go test -race ./...

# Benchmark
go test -bench=. ./...

# Profile goroutines
curl http://localhost:6060/debug/pprof/goroutine
```

### Essential Patterns
```go
// Goroutine
go func() { /* work */ }()

// Channel
ch := make(chan int)
ch <- value
value := <-ch

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()

// Select
select {
case msg := <-ch:
    // handle
case <-time.After(timeout):
    // timeout
}

// Context
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Mutex
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()
```

---

## ✅ Checklist

### Beginner
- [ ] Complete Crash Course
- [ ] Run all examples
- [ ] Complete 3 exercises
- [ ] Build simple concurrent program

### Intermediate
- [ ] Complete Deep Dive
- [ ] Run advanced examples
- [ ] Understand race conditions
- [ ] Build worker pool

### Advanced
- [ ] Master all patterns
- [ ] Debug complex issues
- [ ] Profile and optimize
- [ ] Build production system

---

## 🤝 Contributing

Want to improve these materials?

1. Fork the repository
2. Add examples or improvements
3. Test with `-race` flag
4. Submit pull request

---

## 📞 Support

- **Issues:** GitHub Issues
- **Discussions:** GitHub Discussions
- **Questions:** Stack Overflow (tag: go, concurrency)

---

**Start your concurrency journey today! 🚀**

*Choose your path above and begin learning!*

