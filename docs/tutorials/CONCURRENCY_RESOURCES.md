# 🔄 Go Concurrency Learning Resources

**Complete guide to all concurrency resources in GO-PRO**

---

## 📚 Learning Materials

### 1. ⚡ Concurrency Crash Course
**Location:** `docs/tutorials/CONCURRENCY_CRASH_COURSE.md`  
**Duration:** 60-90 minutes  
**Level:** Intermediate  
**Type:** Hands-on tutorial

**What's Included:**
- 11 runnable examples
- Real-world patterns
- Common pitfalls
- Practice exercises
- Quick reference guide

**Start Here If:**
- You want to get productive quickly
- You prefer learning by doing
- You need practical examples
- You're short on time

---

### 2. 📖 Concurrency Deep Dive
**Location:** `docs/tutorials/concurrency-deep-dive.md`  
**Duration:** 4-5 hours  
**Level:** Advanced  
**Type:** Comprehensive guide

**What's Included:**
- Goroutine lifecycle and scheduling
- Channel patterns and idioms
- Deadlock prevention
- Race condition detection
- Memory model principles
- Advanced patterns

**Use This When:**
- You want deep understanding
- You need to debug complex issues
- You're building production systems
- You want to master concurrency

---

### 3. 📄 Quick Reference
**Location:** `docs/tutorials/CONCURRENCY_QUICK_REFERENCE.md`  
**Duration:** 5 minutes  
**Level:** All levels  
**Type:** Cheat sheet

**What's Included:**
- Syntax reference
- Common patterns
- Best practices
- Pitfalls to avoid

**Use This When:**
- You need a quick reminder
- You're coding and need syntax
- You want a printable reference

---

## 💻 Runnable Examples

### Crash Course Examples
**Location:** `basic/examples/concurrency-crash-course/`

```bash
cd basic/examples/concurrency-crash-course
go run main.go
```

**11 Examples:**
1. Basic Goroutines
2. Channels
3. WaitGroups
4. Select
5. Worker Pool
6. Pipeline
7. Fan-Out/Fan-In
8. Context
9. Mutex
10. Web Scraper
11. Rate Limiter

---

### Basic Concurrency Examples
**Location:** `basic/examples/11_concurrency/`

```bash
cd basic/examples/11_concurrency
go run main.go
```

**10 Examples:**
- Basic goroutines
- WaitGroups
- Channels
- Buffered channels
- Channel directions
- Select
- Worker pool
- Producer-consumer
- Mutex
- Timeout patterns

---

### Advanced Examples

#### Worker Pool
**Location:** `advanced/go_18_worker_pool/`
```bash
cd advanced/go_18_worker_pool
go run main.go
```

#### Pipeline
**Location:** `advanced/go_21_goroutines_pipeline/`
```bash
cd advanced/go_21_goroutines_pipeline
go run main.go
```

#### Fan-Out/Fan-In
**Location:** `advanced/go_25_fan_out_fan_in/`
```bash
cd advanced/go_25_fan_out_fan_in/after
go run main.go
```

#### Channels (Buffered/Unbuffered)
**Location:** `advanced/go_24_unbuffered_buffered_channels/`
```bash
cd advanced/go_24_unbuffered_buffered_channels
go run buffered.go
go run unbuffered.go
```

#### Context
**Location:** `advanced/go_51_context/`
```bash
cd advanced/go_51_context/go_ctx_cancel
go run main.go
```

#### Mutexes
**Location:** `advanced/go_20_mutexes_and_confinement/`
```bash
cd advanced/go_20_mutexes_and_confinement/with_mutex
go run main.go
```

---

## 🎓 Learning Paths

### Path 1: Quick Start (1-2 hours)
1. Read: Concurrency Crash Course
2. Run: All crash course examples
3. Practice: Complete 2-3 exercises
4. Reference: Keep quick reference handy

### Path 2: Comprehensive (1 week)
1. **Day 1**: Concurrency Crash Course
2. **Day 2-3**: Basic examples + exercises
3. **Day 4-5**: Concurrency Deep Dive
4. **Day 6-7**: Advanced examples + projects

### Path 3: Production Ready (2 weeks)
1. **Week 1**: All tutorials + examples
2. **Week 2**: Build real projects:
   - Concurrent web scraper
   - Worker pool system
   - Rate-limited API client
   - Pipeline processor

---

## 🧪 Testing & Debugging

### Race Detector
```bash
# Run with race detector
go run -race main.go

# Build with race detector
go build -race

# Test with race detector
go test -race ./...
```

### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Goroutine profiling
curl http://localhost:6060/debug/pprof/goroutine
```

### Benchmarking
```bash
# Run benchmarks
go test -bench=. ./...

# With memory stats
go test -bench=. -benchmem ./...
```

---

## 📖 External Resources

### Official Documentation
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go Blog - Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Blog - Advanced Patterns](https://go.dev/blog/io2013-talk-concurrency)

### Videos
- [Google I/O 2012 - Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Google I/O 2013 - Advanced Concurrency](https://www.youtube.com/watch?v=QDDwwePbDtw)
- [Rob Pike - Concurrency Is Not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg)

### Books
- "Concurrency in Go" by Katherine Cox-Buday
- "The Go Programming Language" by Donovan & Kernighan (Chapters 8-9)
- "Go in Action" by William Kennedy (Chapter 6)

### Interactive
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go by Example - Channels](https://gobyexample.com/channels)
- [Go Playground](https://go.dev/play/)

---

## 🎯 Quick Navigation

### By Topic

**Goroutines:**
- Crash Course: Section 1
- Deep Dive: Section 1
- Examples: `basic/examples/11_concurrency/`

**Channels:**
- Crash Course: Section 2
- Deep Dive: Section 2
- Examples: `advanced/go_24_unbuffered_buffered_channels/`

**Patterns:**
- Crash Course: Sections 5-7
- Deep Dive: Section 5
- Examples: `advanced/go_18_worker_pool/`, `advanced/go_21_goroutines_pipeline/`

**Context:**
- Crash Course: Section 8
- Deep Dive: Section 5
- Examples: `advanced/go_51_context/`

**Synchronization:**
- Crash Course: Sections 3, 9
- Deep Dive: Sections 3, 4
- Examples: `advanced/go_20_mutexes_and_confinement/`

---

## 🚀 Getting Started

### Complete Beginner
1. Start with: **Concurrency Crash Course**
2. Run: All crash course examples
3. Practice: Complete exercises
4. Next: Basic concurrency examples

### Experienced Developer
1. Skim: Concurrency Crash Course
2. Deep dive: Concurrency Deep Dive
3. Study: Advanced examples
4. Build: Real-world projects

### Production Developer
1. Review: All tutorials
2. Study: Memory model and race conditions
3. Practice: Debugging and profiling
4. Implement: Production patterns

---

## 📊 Progress Tracking

### Beginner Level
- [ ] Complete Concurrency Crash Course
- [ ] Run all 11 crash course examples
- [ ] Complete 3 practice exercises
- [ ] Build simple concurrent program

### Intermediate Level
- [ ] Complete Concurrency Deep Dive
- [ ] Run all advanced examples
- [ ] Understand race conditions
- [ ] Build worker pool system

### Advanced Level
- [ ] Master all patterns
- [ ] Debug complex concurrency issues
- [ ] Profile and optimize
- [ ] Build production systems

---

## 🤝 Contributing

Want to add concurrency examples or tutorials?

1. Follow the crash course format
2. Include runnable code
3. Add tests with `-race` flag
4. Document patterns clearly
5. Provide real-world use cases

---

## 📝 Summary

**Start Here:**
- ⚡ [Concurrency Crash Course](CONCURRENCY_CRASH_COURSE.md) - 60-90 minutes

**Go Deeper:**
- 📖 [Concurrency Deep Dive](concurrency-deep-dive.md) - 4-5 hours

**Quick Reference:**
- 📄 [Quick Reference](CONCURRENCY_QUICK_REFERENCE.md) - 5 minutes

**Practice:**
- 💻 [Runnable Examples](../../basic/examples/concurrency-crash-course/)

---

**Master concurrency, master Go! 🚀**

