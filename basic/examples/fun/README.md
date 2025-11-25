# Go Programming Examples - Fun Directory

A comprehensive collection of Go programming examples, patterns, and best practices organized by category.

## 📁 Project Structure

```
fun/
├── cmd/examples/          # Runnable example programs
│   ├── basics/           # Basic Go concepts (variables, functions, pointers, structs, etc.)
│   ├── algorithms/       # Algorithm demonstrations (search, sort)
│   ├── concurrency/      # Concurrency patterns (goroutines, channels, worker pools)
│   ├── datastructures/   # Data structure examples (stack, queue, linked list)
│   ├── cache/            # Caching examples
│   └── advanced/         # Advanced topics
├── pkg/                  # Reusable library packages
│   ├── algorithms/       # Algorithm implementations (search, sort, math, strings)
│   ├── cache/            # Thread-safe cache with TTL, LRU, LFU
│   ├── concurrency/      # Concurrency utilities (rate limiters, worker pools, barriers)
│   ├── datastructures/   # Data structures (stack, queue, linked list)
│   ├── test/             # Testing utilities
│   └── utils/            # Common utility functions
├── cache-demo/           # Standalone cache demonstration
├── _archive/             # Archived old examples (not for use)
├── main.go               # Entry point with usage instructions
├── doc.go                # Package documentation
└── README.md             # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- Basic understanding of Go programming

### Installation

```bash
cd basic/examples/fun
go mod download
```

### Running Examples

#### Quick Start (Main Entry Point)
```bash
go run main.go
# Shows available examples and how to run them
```

#### Run Specific Examples
```bash
# Data structures
go run cmd/examples/datastructures/stack_demo.go
go run cmd/examples/datastructures/queue_demo.go

# Algorithms
go run cmd/examples/algorithms/search_demo.go
go run cmd/examples/algorithms/sort_demo.go

# Concurrency
go run cmd/examples/concurrency/goroutines_demo.go
go run cmd/examples/concurrency/channels_demo.go

# Cache
go run cmd/examples/cache/cache_demo.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## 📚 Categories

### 1. Data Structures (`pkg/datastructures`)
- **Stack (LIFO)**: Last-In-First-Out data structure
- **Queue (FIFO)**: First-In-First-Out data structure with priority queue
- **Linked List**: Singly and doubly linked list implementations

### 2. Algorithms (`pkg/algorithms`)
- **Search**: Binary search, linear search
- **Sort**: Merge sort, quick sort (sequential and concurrent)
- **Strings**: Palindrome detection, word counting
- **Math**: Prime number generation, Fibonacci sequence

### 3. Concurrency (`pkg/concurrency`) ✅
- **Goroutines**: Basic goroutine usage, WaitGroups, Fan-Out/Fan-In
- **Thread-Safe Structures**: SafeCounter, SafeMap, Semaphore, Barrier
- **Producer-Consumer**: Generic producer-consumer with worker pools
- **Worker Pool**: Efficient task distribution and processing
- **Rate Limiters**: Token Bucket, Sliding Window, Leaky Bucket, Adaptive
- **Context**: Timeout, cancellation, task groups, typed values
- **Parallel Processing**: Map, Filter, Reduce, Batch processing

### 4. Cache (`pkg/cache`) ✅
- **Generic Cache**: Thread-safe cache with TTL, statistics, and auto-eviction
- **Loading Cache**: Automatic value loading on cache miss
- **LRU Cache**: Least Recently Used eviction policy
- **LFU Cache**: Least Frequently Used eviction policy
- **Cache Statistics**: Hit rate, evictions, expirations tracking

### 5. Basics (`cmd/examples/basics`) ✅
- **Variables & Types**: int, float, string, bool, type inference, constants, zero values
- **Functions**: Basic, multiple returns, named returns, variadic, higher-order, closures
- **Pointers**: Addresses, dereferencing, pass by value vs reference, pointer receivers
- **Structs & Methods**: Value/pointer receivers, embedded structs, struct tags
- **Interfaces**: Polymorphism, type assertions, type switches, empty interface
- **Loops & Control Flow**: for, range, break, continue, labeled statements
- **Iota & Constants**: Enums, bit flags, custom expressions, real-world examples

### 6. Advanced (`cmd/examples/advanced`)
- JSON parsing and serialization
- HTTP clients and APIs
- Error handling patterns
- File I/O operations
- Testing and benchmarking

## 🎯 Key Features

- ✅ **Clean Architecture**: Well-organized, modular code
- ✅ **Best Practices**: Follows Go idioms and conventions
- ✅ **Type Safety**: Uses generics where appropriate
- ✅ **Thread Safety**: Proper synchronization in concurrent code
- ✅ **Error Handling**: Comprehensive error handling patterns
- ✅ **Testing**: Unit tests and benchmarks for all packages
- ✅ **Documentation**: GoDoc-style comments throughout
- ✅ **Examples**: Runnable examples for each concept

## 📖 Learning Path

### Beginner
1. Start with `cmd/examples/basics/` - Learn fundamental Go concepts
2. Explore `pkg/datastructures/` - Understand basic data structures
3. Try `pkg/algorithms/` - Practice with common algorithms

### Intermediate
1. Study `pkg/concurrency/` - Master goroutines and channels
2. Implement `pkg/cache/` - Build thread-safe systems
3. Review `cmd/examples/advanced/` - Work with real-world scenarios

### Advanced
1. Optimize with benchmarks - Use `go test -bench`
2. Add observability - Integrate OpenTelemetry
3. Build microservices - Apply patterns to distributed systems

## 🛠️ Development

### Adding New Examples

1. Create package in appropriate directory
2. Add implementation with GoDoc comments
3. Create test file with `_test.go` suffix
4. Add example in `cmd/examples/`
5. Update this README

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Write table-driven tests
- Add benchmarks for performance-critical code

## 📝 Examples

### Quick Start: Stack

```go
package main

import (
    "fmt"
    "github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/datastructures"
)

func main() {
    stack := datastructures.NewStack[int]()
    
    stack.Push(1)
    stack.Push(2)
    stack.Push(3)
    
    value, _ := stack.Pop()
    fmt.Println(value) // Output: 3
}
```

### Quick Start: Cache

```go
package main

import (
    "fmt"
    "time"
    "github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/cache"
)

func main() {
    c := cache.New[string, string](5 * time.Minute)
    
    c.Set("key", "value")
    
    if val, found := c.Get("key"); found {
        fmt.Println(val) // Output: value
    }
}
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Follow the existing code style
2. Add tests for new features
3. Update documentation
4. Run `go test ./...` before submitting

## 📄 License

This project is part of the go-pro learning repository.

## 🔗 Resources

- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Blog](https://blog.golang.org/)

---

**Happy Learning! 🎉**

