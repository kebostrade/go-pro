```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                     📚 GO-PRO INTERACTIVE TUTORIALS                          ║
║                                                                              ║
║          Master Go Programming Through Hands-On Projects & Examples          ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**Welcome to Go-Pro!** This guide will take you from Go basics to building production-ready systems through step-by-step tutorials, real-world projects, and practical examples.

---

## 📑 TABLE OF CONTENTS

### 🌟 [Fundamentals](#-fundamentals-tutorials) (2-3 hours)
- [Tutorial 0.1](#-tutorial-01-hello-world--basic-syntax) - Hello World & Basic Syntax ⏱️ 10 min
- [Tutorial 0.2](#-tutorial-02-data-structures-basics) - Data Structures Basics ⏱️ 20 min
- [Tutorial 0.3](#%EF%B8%8F-tutorial-03-structs-and-interfaces) - Structs and Interfaces ⏱️ 25 min
- [Tutorial 0.4](#-tutorial-04-concurrency-basics) - Concurrency Basics ⏱️ 30 min
- [Tutorial 0.5](#-tutorial-05-testing-in-go) - Testing in Go ⏱️ 25 min
- [Tutorial 0.6](#-tutorial-06-file-io-operations) - File I/O Operations ⏱️ 20 min

### 🚀 [Projects](#-project-tutorials) (3-5 hours)
- [Tutorial 1](#-tutorial-1-your-first-go-project) - URL Shortener Service ⏱️ 15 min
- [Tutorial 2](#%EF%B8%8F-tutorial-2-building-a-cli-application) - Weather CLI Application ⏱️ 20 min
- [Tutorial 3](#-tutorial-3-file-encryption) - File Encryption Tool ⏱️ 15 min
- [Tutorial 4](#-tutorial-4-building-a-blog-api) - Blog API with Auth ⏱️ 30 min
- [Tutorial 5](#%EF%B8%8F-tutorial-5-job-queue-system) - Job Queue System ⏱️ 45 min
- [Tutorial 6](#-tutorial-6-rate-limiting) - Rate Limiting ⏱️ 30 min
- [Tutorial 7](#-tutorial-7-log-aggregation) - Log Aggregation ⏱️ 60 min
- [Tutorial 8](#%EF%B8%8F-tutorial-8-service-mesh) - Service Mesh ⏱️ 90 min
- [Tutorial 9](#-tutorial-9-time-series-database) - Time Series Database ⏱️ 120 min
- [Tutorial 10](#-tutorial-10-container-orchestrator) - Container Orchestrator ⏱️ 150 min

### 🎨 [Specialized Topics](#-specialized-tutorials) (4-6 hours)
- [Tutorial 11](#-tutorial-11-advanced-cryptography) - Advanced Cryptography ⏱️ 45 min
- [Tutorial 12](#-tutorial-12-websocket-real-time-communication) - WebSocket Real-Time ⏱️ 40 min
- [Tutorial 13](#-tutorial-13-algorithms--data-structures) - Algorithms & Data Structures ⏱️ 60 min
- [Tutorial 14](#-tutorial-14-performance-optimization) - Performance Optimization ⏱️ 50 min
- [Tutorial 15](#-tutorial-15-docker--deployment) - Docker & Deployment ⏱️ 45 min

### 📈 [Learning Paths](#-learning-path-summary)
- [Beginner Path](#beginner-path) - Start Here
- [Intermediate Path](#intermediate-path) - Next Steps
- [Advanced Path](#advanced-path) - Master Level

---

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                        🌟 FUNDAMENTALS TUTORIALS                             ║
║                                                                              ║
║                 Master Core Go Concepts Through Practice                     ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 👋 Tutorial 0.1: Hello World & Basic Syntax

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  10 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Learn Go basics and run your first program                    │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Go program structure                                              │
│     ✓ Package and imports                                               │
│     ✓ Functions and main entry point                                    │
│     ✓ Running Go programs                                               │
│     ✓ Basic syntax and formatting                                       │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Navigate to Examples
```bash
cd basic/examples/01_hello
```

#### Step 2: View the Code
```bash
cat main.go
```

**💡 Quick Tip:** Notice the `package main` declaration and `func main()` - these are required for executable programs.

#### Step 3: Run the Program
```bash
go run main.go
```

**📤 Expected Output:**
```
Hello, World!
Welcome to Go Programming!
```

#### Step 4: Modify and Experiment
```bash
# Edit main.go to print your name
# Change "World" to your name
# Run again to see the changes
```

**🎯 Try This Challenge:**
- Modify the program to print 5 different messages
- Add a variable to store your name
- Use formatted strings with `fmt.Printf()`

#### Step 5: Explore More Examples
```bash
# Variables and types
cd ../02_variables
go run main.go

# Functions
cd ../03_functions
go run main.go

# Control flow
cd ../06_control_flow
go run main.go
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Can you write and run a Go program?                     │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I understand package declarations                                  │
│  [ ] I can write a main function                                        │
│  [ ] I know how to run Go programs                                      │
│  [ ] I can modify and experiment with code                              │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Congratulations!** You've written your first Go program!

---

## 🔢 Tutorial 0.2: Data Structures Basics

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  20 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Master arrays, slices, and maps                               │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Fixed-size arrays vs dynamic slices                              │
│     ✓ Slice operations (append, copy, slice)                           │
│     ✓ Working with maps (dictionaries)                                 │
│     ✓ Iterating over collections                                       │
│     ✓ Common data structure patterns                                   │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Arrays and Slices
```bash
cd basic/examples/05_arrays_slices
go run main.go
```

**📖 Key Concepts:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Arrays:  Fixed size, value type                               │
│  Slices:  Dynamic size, reference type                         │
│                                                                 │
│  var arr [5]int          // Array - size is part of type       │
│  var slice []int         // Slice - flexible size              │
│                                                                 │
│  slice = append(slice, 1, 2, 3)  // Dynamic growth             │
└─────────────────────────────────────────────────────────────────┘
```

**💡 Pro Tip:** Always use slices instead of arrays unless you need a fixed size.

#### Step 2: Maps (Key-Value Pairs)
```bash
cd basic/examples/07_maps
go run main.go
```

**📖 Key Operations:**
```go
// Creating maps
users := make(map[string]int)
users["alice"] = 30

// Checking existence
age, exists := users["alice"]
if exists {
    fmt.Println("Age:", age)
}

// Deleting elements
delete(users, "alice")

// Iterating
for key, value := range users {
    fmt.Printf("%s: %d\n", key, value)
}
```

#### Step 3: Practice with Real Implementations
```bash
cd basic/examples/fun

# Queue implementation (FIFO)
go run cmd/examples/datastructures/queue_demo.go

# Stack implementation (LIFO)
go run cmd/examples/datastructures/stack_demo.go

# Linked list
go run cmd/examples/datastructures/linked_list_demo.go
```

**🎯 Try This Challenge:**
- Implement a function to reverse a slice
- Create a map to count word frequencies in a string
- Build a simple cache using a map

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Do you understand Go data structures?                   │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I know when to use arrays vs slices                                │
│  [ ] I can use append, copy, and slicing operations                     │
│  [ ] I understand map operations (add, get, delete, check)              │
│  [ ] I can iterate over slices and maps                                 │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Success!** You understand Go data structures!

---

## 🏗️ Tutorial 0.3: Structs and Interfaces

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  25 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Learn object-oriented patterns in Go                          │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Defining and using structs                                       │
│     ✓ Methods on structs                                               │
│     ✓ Value receivers vs pointer receivers                             │
│     ✓ Interface definitions and implementations                        │
│     ✓ Type assertions and polymorphism                                 │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Structs Basics
```bash
cd basic/examples/08_structs
go run main.go
```

**📖 Key Concepts:**
```go
// Define a struct
type Person struct {
    Name string
    Age  int
}

// Value receiver - receives copy
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

// Pointer receiver - can modify original
func (p *Person) Birthday() {
    p.Age++
}
```

**⚠️ Important:** Use pointer receivers when:
- You need to modify the receiver
- The struct is large (avoid copying)
- Consistency (if some methods use pointers, all should)

#### Step 2: Interfaces
```bash
cd basic/examples/09_interfaces
go run main.go
```

**📖 Interface Pattern:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Interfaces in Go are IMPLICIT                                  │
│                                                                 │
│  1. Define interface with method signatures                    │
│  2. Implement methods on your type                             │
│  3. No explicit "implements" keyword needed                    │
│  4. Use interface types for flexibility                        │
└─────────────────────────────────────────────────────────────────┘
```

**Example:**
```go
// Interface definition
type Speaker interface {
    Speak() string
}

// Dog implements Speaker (implicitly)
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

// Cat implements Speaker (implicitly)
type Cat struct {
    Name string
}

func (c Cat) Speak() string {
    return "Meow!"
}

// Polymorphism in action
func MakeSound(s Speaker) {
    fmt.Println(s.Speak())
}
```

#### Step 3: Real-World Example
```bash
cd basic/examples/fun
go run cmd/examples/basics/interfaces_demo.go
```

**🎯 Try This Challenge:**
- Create a `Shape` interface with `Area()` and `Perimeter()` methods
- Implement it for `Circle`, `Rectangle`, and `Triangle`
- Write a function that calculates total area of any shapes

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Have you mastered Go's type system?                     │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I can define and use structs                                       │
│  [ ] I understand value vs pointer receivers                            │
│  [ ] I can define and implement interfaces                              │
│  [ ] I understand type assertions and polymorphism                      │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Excellent!** You've mastered Go's type system!

---

## ⚡ Tutorial 0.4: Concurrency Basics

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  30 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Understand goroutines and channels                            │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Creating and using goroutines                                    │
│     ✓ Channel communication patterns                                   │
│     ✓ WaitGroups for synchronization                                   │
│     ✓ Context for cancellation and timeouts                            │
│     ✓ Common concurrency patterns                                      │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Goroutines
```bash
cd basic/examples/fun
go run cmd/examples/concurrency/goroutines_demo.go
```

**📖 Goroutines Overview:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Goroutines: Lightweight concurrent execution                  │
│                                                                 │
│  // Start a goroutine                                          │
│  go function()                                                 │
│  go func() { /* anonymous function */ }()                      │
│                                                                 │
│  💡 Goroutines are NOT threads - they're much lighter!         │
│     You can have thousands running simultaneously              │
└─────────────────────────────────────────────────────────────────┘
```

**⚠️ Important:** Always use WaitGroups or channels to ensure goroutines complete before main exits!

#### Step 2: Channels
```bash
go run cmd/examples/concurrency/channels_demo.go
```

**📖 Channel Patterns:**
```go
// Unbuffered channel - blocks until receiver ready
ch := make(chan int)

// Buffered channel - can hold N values
ch := make(chan int, 10)

// Send to channel
ch <- 42

// Receive from channel
value := <-ch

// Range over channel (until closed)
for value := range ch {
    fmt.Println(value)
}

// Close channel
close(ch)
```

**💡 Pro Tip:** Use channels to communicate between goroutines. "Don't communicate by sharing memory; share memory by communicating."

#### Step 3: Practical Patterns
```bash
# Worker pool pattern
go run cmd/examples/concurrency/worker_pool_demo.go

# Producer-consumer pattern
cd basic/examples
go run producer_consumer.go

# Rate limiter
go run rate_limiter.go
```

**📖 Worker Pool Pattern:**
```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│     Jobs Channel        Workers           Results Channel       │
│          ↓                 ↓                    ↑               │
│       [J][J][J]      [ W W W W ]          [R][R][R]            │
│          ↓                 ↓                    ↑               │
│      Queue jobs → Workers process → Collect results            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Step 4: Context and Timeouts
```bash
go run context_timeout.go
```

**📖 Context Usage:**
```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use in goroutine
go func(ctx context.Context) {
    select {
    case <-ctx.Done():
        fmt.Println("Timeout or cancelled")
        return
    case <-time.After(10 * time.Second):
        fmt.Println("Work completed")
    }
}(ctx)
```

**🎯 Try This Challenge:**
- Create a program that fetches URLs concurrently
- Implement a timeout for each request
- Collect all results or cancel after timeout

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Do you understand Go concurrency?                       │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I can create and use goroutines                                    │
│  [ ] I understand channel communication                                 │
│  [ ] I can use WaitGroups for synchronization                           │
│  [ ] I understand context for cancellation                              │
│  [ ] I know common concurrency patterns                                 │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Amazing!** You understand Go concurrency!

---

## 🧪 Tutorial 0.5: Testing in Go

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  25 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Write and run tests like a pro                                │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Writing basic tests                                              │
│     ✓ Table-driven test pattern                                        │
│     ✓ Benchmarking performance                                         │
│     ✓ Measuring test coverage                                          │
│     ✓ Testing best practices                                           │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Basic Testing
```bash
cd basic/examples/13_testing/01_basic_test

# View the test file
cat math_test.go

# Run tests
go test -v
```

**📤 Expected Output:**
```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestSubtract
--- PASS: TestSubtract (0.00s)
PASS
ok      github.com/DimaJoyti/go-pro/basic/examples/13_testing/01_basic_test    0.004s
```

**📖 Test Anatomy:**
```go
// Test file: *_test.go
// Test function: TestXxx(t *testing.T)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5

    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

#### Step 2: Table-Driven Tests
```bash
cd ../02_table_driven_tests

# View the test
cat calculator_test.go

# Run tests
go test -v
```

**📖 Table-Driven Pattern:**
```go
func TestCalculator(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        op       string
        expected int
    }{
        {"add positive", 2, 3, "+", 5},
        {"subtract", 5, 3, "-", 2},
        {"multiply", 3, 4, "*", 12},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Calculate(tt.a, tt.b, tt.op)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

**💡 Pro Tip:** Table-driven tests make it easy to add new test cases - just add a row!

#### Step 3: Benchmarks
```bash
cd ../03_benchmarks

# Run benchmarks
go test -bench=.

# With memory allocation stats
go test -bench=. -benchmem
```

**📤 Expected Output:**
```
BenchmarkFibonacci-8    5000000    250 ns/op    0 B/op    0 allocs/op
```

**📖 Benchmark Pattern:**
```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        Function()
    }
}
```

#### Step 4: Test Coverage
```bash
# Run tests with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

**📊 Coverage Report:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Coverage Goals:                                                │
│                                                                 │
│  🟢 > 80%  - Excellent coverage                                 │
│  🟡 60-80% - Good coverage                                      │
│  🔴 < 60%  - Needs improvement                                  │
│                                                                 │
│  💡 Focus on critical paths first                               │
└─────────────────────────────────────────────────────────────────┘
```

**🎯 Try This Challenge:**
- Write tests for a function you created
- Use table-driven tests with at least 5 cases
- Achieve > 80% code coverage
- Benchmark different implementations

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Can you write professional tests?                       │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I can write basic test functions                                   │
│  [ ] I understand table-driven test pattern                             │
│  [ ] I can benchmark code performance                                   │
│  [ ] I know how to measure test coverage                                │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Great work!** You can write professional tests!

---

## 📁 Tutorial 0.6: File I/O Operations

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  20 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Master file operations in Go                                  │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Reading files (multiple methods)                                 │
│     ✓ Writing files safely                                             │
│     ✓ Line-by-line processing                                          │
│     ✓ Directory operations                                             │
│     ✓ File permissions and metadata                                    │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Reading Files
```bash
cd basic/examples/12_file_io/01_read_file

# Create a test file
echo "Hello from file!" > test.txt

# Run the example
go run main.go
```

**📖 Reading Methods:**
```go
// Method 1: Read entire file
data, err := os.ReadFile("test.txt")

// Method 2: Open and read manually
file, err := os.Open("test.txt")
defer file.Close()
scanner := bufio.NewScanner(file)

// Method 3: Read with buffer
buf := make([]byte, 1024)
n, err := file.Read(buf)
```

#### Step 2: Writing Files
```bash
cd ../02_write_file

# Run the example
go run main.go

# Check the created file
cat output.txt
```

**📖 Writing Patterns:**
```go
// Write entire file (overwrites)
err := os.WriteFile("output.txt", data, 0644)

// Append to file
file, err := os.OpenFile("output.txt",
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
defer file.Close()

// Buffered writing
writer := bufio.NewWriter(file)
writer.WriteString("Hello\n")
writer.Flush()
```

#### Step 3: Line-by-Line Reading
```bash
cd ../04_read_line_by_line

# Create a large file
for i in {1..100}; do echo "Line $i" >> large.txt; done

# Read it efficiently
go run main.go
```

**💡 Pro Tip:** Use `bufio.Scanner` for line-by-line reading of large files - it's memory efficient!

#### Step 4: Directory Operations
```bash
cd ../06_directory_operations

# Run the example
go run main.go
```

**📖 Directory Operations:**
```go
// Create directory
os.Mkdir("mydir", 0755)
os.MkdirAll("path/to/nested/dir", 0755)

// List files
files, err := os.ReadDir(".")
for _, file := range files {
    fmt.Println(file.Name())
}

// Walk directory tree
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
})

// Get file info
info, err := os.Stat("file.txt")
fmt.Println("Size:", info.Size())
fmt.Println("ModTime:", info.ModTime())
```

**🎯 Try This Challenge:**
- Create a program to copy a file
- Write a log file rotator (delete files older than N days)
- Build a simple file search tool
- Implement a word counter for text files

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ✅ CHECKPOINT: Can you work with files confidently?                    │
│                                                                          │
│  Self-Assessment:                                                        │
│  [ ] I can read files using different methods                           │
│  [ ] I know how to write and append to files                            │
│  [ ] I can process files line-by-line efficiently                       │
│  [ ] I understand directory operations                                  │
│  [ ] I know about file permissions and metadata                         │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Fantastic!** You can work with files in Go!

---

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                          🚀 PROJECT TUTORIALS                                ║
║                                                                              ║
║                  Build Production-Ready Applications                         ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 🚀 Tutorial 1: Your First Go Project

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  15 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: URL Shortener Service                                      │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ REST API with HTTP handlers                                      │
│     ✓ In-memory data storage                                           │
│     ✓ Short URL generation                                             │
│     ✓ Click analytics tracking                                         │
│     ✓ JSON request/response handling                                   │
│                                                                          │
│  🛠️ TECH STACK: Go standard library, HTTP server, JSON                  │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Navigate to the Project
```bash
cd basic/projects/url-shortener
```

#### Step 2: Explore the Structure
```bash
# View the README
cat README.md

# Check the project structure
tree -L 2
```

**📁 Project Structure:**
```
url-shortener/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── storage/         # Data storage
├── tests/               # Integration tests
├── Makefile            # Build automation
└── README.md           # Documentation
```

#### Step 3: Run Tests
```bash
# Run all tests
make test
```

**📤 Expected Output:**
```
✓ All tests passed
PASS
ok      github.com/DimaJoyti/go-pro/basic/projects/url-shortener/tests    0.004s
```

#### Step 4: Build and Run
```bash
# Build the application
make build

# Run the server
make run
```

**📤 Server Output:**
```
🚀 URL Shortener Server
═══════════════════════════════════════════════
✓ Server starting on :8080
✓ Ready to accept requests
```

#### Step 5: Test the API

**Create a Short URL:**
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/DimaJoyti/go-pro"
  }'
```

**📤 Response:**
```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://github.com/DimaJoyti/go-pro"
}
```

**Visit the Short URL:**
```bash
curl -L http://localhost:8080/abc123
# Redirects to original URL
```

**Get Analytics:**
```bash
curl http://localhost:8080/api/analytics/abc123
```

**📤 Analytics Response:**
```json
{
  "short_code": "abc123",
  "original_url": "https://github.com/DimaJoyti/go-pro",
  "clicks": 5,
  "created_at": "2024-01-15T10:30:00Z",
  "last_accessed": "2024-01-15T11:45:00Z"
}
```

**🎯 Quick Wins:**
```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Shorten 5 different URLs                                     │
│  ✓ Visit each short link                                        │
│  ✓ Check analytics for click counts                             │
│  ✓ Test with invalid URLs                                       │
│  ✓ Explore the code in internal/ directory                      │
└──────────────────────────────────────────────────────────────────┘
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  🎓 WHAT YOU LEARNED:                                                    │
│                                                                          │
│  • HTTP server setup with net/http                                      │
│  • JSON encoding/decoding                                               │
│  • REST API design patterns                                             │
│  • In-memory data storage                                               │
│  • Project structure and organization                                   │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Congratulations!** You've built your first Go web service!

---

## 🌤️ Tutorial 2: Building a CLI Application

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  20 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Weather CLI Application                                    │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ Command-line interface with cobra                                │
│     ✓ External API integration                                         │
│     ✓ Response caching strategy                                        │
│     ✓ Formatted table output                                           │
│     ✓ Configuration management                                         │
│                                                                          │
│  🛠️ TECH STACK: Cobra, OpenWeatherMap API, Cache                        │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Get an API Key
```
┌──────────────────────────────────────────────────────────────────┐
│  1. Visit https://openweathermap.org/api                        │
│  2. Sign up for a free account                                  │
│  3. Get your API key from dashboard                             │
│  4. Free tier: 60 calls/minute, 1,000,000 calls/month          │
└──────────────────────────────────────────────────────────────────┘
```

#### Step 2: Setup the Project
```bash
cd basic/projects/weather-cli

# Set your API key
export WEATHER_API_KEY="your-api-key-here"

# Or create .env file
echo "WEATHER_API_KEY=your-api-key-here" > .env
```

#### Step 3: Build the CLI
```bash
make build

# Binary created at: bin/weather
```

#### Step 4: Get Current Weather
```bash
# Get weather for a city
./bin/weather current --city "London"
```

**📤 Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              ☀️  Weather in London                        ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  🌡️  Temperature:        15°C                            ║
║  🤚 Feels Like:          13°C                            ║
║  ☁️  Conditions:         Clear sky                       ║
║  💧 Humidity:            65%                             ║
║  💨 Wind:                12 km/h NW                      ║
║  👁️  Visibility:         10 km                           ║
║  🌅 Sunrise:             06:42 AM                        ║
║  🌇 Sunset:              07:15 PM                        ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

#### Step 5: Get Forecast
```bash
# Get 5-day forecast
./bin/weather forecast --city "Tokyo"

# Get detailed forecast
./bin/weather forecast --city "Paris" --detailed

# Get JSON output
./bin/weather current --city "New York" --format json
```

**📤 Forecast Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║            📅 5-Day Forecast for Tokyo                    ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Day 1  │  ☀️  Sunny        │  High: 22°C  │  Low: 15°C  ║
║  Day 2  │  ⛅ Partly Cloudy │  High: 20°C  │  Low: 14°C  ║
║  Day 3  │  🌧️  Rainy        │  High: 18°C  │  Low: 13°C  ║
║  Day 4  │  ☁️  Cloudy       │  High: 19°C  │  Low: 14°C  ║
║  Day 5  │  ☀️  Clear        │  High: 23°C  │  Low: 16°C  ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

#### Step 6: Explore Caching
```bash
# First request (hits API)
time ./bin/weather current --city "Berlin"
# Takes ~300ms

# Second request (uses cache - much faster!)
time ./bin/weather current --city "Berlin"
# Takes ~5ms
```

**💡 Pro Tip:** Cache expires after 10 minutes. Use `--no-cache` flag to force fresh data.

**🎯 Quick Wins:**
```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Check weather for 5 different cities                         │
│  ✓ Get both current and forecast data                           │
│  ✓ Compare cached vs non-cached request speed                   │
│  ✓ Try different output formats (table, json)                   │
│  ✓ Check multiple cities in quick succession                    │
└──────────────────────────────────────────────────────────────────┘
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  🎓 WHAT YOU LEARNED:                                                    │
│                                                                          │
│  • Building CLI apps with Cobra framework                               │
│  • Making HTTP requests to external APIs                                │
│  • Implementing response caching                                        │
│  • Formatting output with tables                                        │
│  • Environment variable configuration                                   │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Awesome!** You've built a production-ready CLI tool!

---

## 🔐 Tutorial 3: File Encryption

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  15 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: File Encryption Tool                                       │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ AES-256-GCM encryption                                           │
│     ✓ Secure key derivation (PBKDF2)                                   │
│     ✓ Password-based encryption                                        │
│     ✓ Progress bars for user feedback                                 │
│     ✓ CLI with encrypt/decrypt commands                               │
│                                                                          │
│  🛠️ TECH STACK: crypto/aes, PBKDF2, CLI                                 │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Setup
```bash
cd basic/projects/file-encryptor
make build
```

#### Step 2: Create a Test File
```bash
echo "This is a secret message!" > secret.txt
cat secret.txt
```

#### Step 3: Encrypt the File
```bash
./bin/encrypt encrypt --input secret.txt
```

**📤 Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              🔐 File Encryption Tool                      ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

📁 Encrypting: secret.txt
📊 Size: 27 B
🔑 Algorithm: AES-256-GCM

Enter password: ********
Confirm password: ********

⚙️  Deriving key from password...
[████████████████████████████] 100%

✓ Encryption complete!
  Output: secret.txt.enc
  Size: 75 B (includes nonce and salt)
```

#### Step 4: Decrypt the File
```bash
./bin/encrypt decrypt --input secret.txt.enc
```

**📤 Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              🔓 File Decryption Tool                      ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

📁 Decrypting: secret.txt.enc
🔑 Algorithm: AES-256-GCM

Enter password: ********

⚙️  Deriving key from password...
[████████████████████████████] 100%

✓ Decryption complete!
  Output: secret.txt.dec
  Size: 27 B
```

#### Step 5: Verify
```bash
# Compare original and decrypted
diff secret.txt secret.txt.dec

# No output means files are identical! ✓
```

#### Step 6: Run the Demo
```bash
# Automated demo with examples
make demo
```

**📖 Security Features:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Security Measures Implemented:                                 │
│                                                                 │
│  🔐 AES-256-GCM encryption (industry standard)                  │
│  🔑 PBKDF2 key derivation (100,000 iterations)                  │
│  🎲 Secure random salt generation                               │
│  🔒 Authenticated encryption (prevents tampering)               │
│  💪 Strong password requirements                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**🎯 Quick Wins:**
```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Encrypt multiple files                                       │
│  ✓ Try wrong password (should fail)                             │
│  ✓ Encrypt a large file (see progress bar)                      │
│  ✓ Inspect encrypted file (looks random)                        │
│  ✓ Check file size increase (overhead)                          │
└──────────────────────────────────────────────────────────────────┘
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  🎓 WHAT YOU LEARNED:                                                    │
│                                                                          │
│  • AES-256-GCM authenticated encryption                                 │
│  • PBKDF2 key derivation from passwords                                 │
│  • Secure random number generation                                      │
│  • Binary file handling                                                 │
│  • CLI progress indicators                                              │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Excellent!** You've mastered cryptography in Go!

---

## 📝 Tutorial 4: Building a Blog API

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🔴 ADVANCED                                    ⏱️  30 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Blog Engine with Authentication                            │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ REST API with PostgreSQL                                         │
│     ✓ JWT authentication                                               │
│     ✓ User registration and login                                      │
│     ✓ CRUD operations for blog posts                                   │
│     ✓ Database migrations                                              │
│                                                                          │
│  🛠️ TECH STACK: PostgreSQL, JWT, Gorilla Mux, GORM                      │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Setup Database
```bash
cd basic/projects/blog-engine

# Create PostgreSQL database
make db-setup

# Run migrations
make db-migrate
```

**📤 Output:**
```
✓ Database created: blogdb
✓ Running migrations...
✓ Table 'users' created
✓ Table 'posts' created
✓ Migrations complete
```

#### Step 2: Start the Server
```bash
# Set environment variables
export DATABASE_URL="postgres://localhost/blogdb?sslmode=disable"
export JWT_SECRET="your-secret-key-change-in-production"

# Run the server
make run
```

**📤 Server Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              📝 Blog Engine API Server                    ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

✓ Database connected
✓ Migrations applied
✓ Server starting on :8080
✓ Ready to accept requests

API Endpoints:
  POST   /api/auth/register    - Register new user
  POST   /api/auth/login       - Login
  GET    /api/posts            - List posts
  POST   /api/posts            - Create post (auth required)
  GET    /api/posts/:id        - Get post
  PUT    /api/posts/:id        - Update post (auth required)
  DELETE /api/posts/:id        - Delete post (auth required)
```

#### Step 3: Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

**📤 Response:**
```json
{
  "id": 1,
  "username": "john",
  "email": "john@example.com",
  "full_name": "John Doe",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Step 4: Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**📤 Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "full_name": "John Doe"
  },
  "expires_at": "2024-01-15T18:30:00Z"
}
```

**💡 Save the token for next steps:**
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### Step 5: Create a Post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "My First Blog Post",
    "content": "# Hello World\n\nThis is my first post!",
    "status": "published",
    "tags": ["golang", "tutorial"]
  }'
```

**📤 Response:**
```json
{
  "id": 1,
  "title": "My First Blog Post",
  "slug": "my-first-blog-post",
  "content": "# Hello World\n\nThis is my first post!",
  "status": "published",
  "tags": ["golang", "tutorial"],
  "author": {
    "id": 1,
    "username": "john",
    "full_name": "John Doe"
  },
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

#### Step 6: Get Posts
```bash
# Get all posts (public)
curl http://localhost:8080/api/posts

# Get specific post
curl http://localhost:8080/api/posts/1

# Get post by slug
curl http://localhost:8080/api/posts/slug/my-first-blog-post

# Filter by status
curl "http://localhost:8080/api/posts?status=published"

# Filter by author
curl "http://localhost:8080/api/posts?author=john"
```

#### Step 7: Update and Delete
```bash
# Update post
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content"
  }'

# Delete post
curl -X DELETE http://localhost:8080/api/posts/1 \
  -H "Authorization: Bearer $TOKEN"
```

**🎯 Quick Wins:**
```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Register 3 users                                              │
│  ✓ Each user creates 2 posts                                     │
│  ✓ Try to create post without auth (should fail)                 │
│  ✓ Try to edit another user's post (should fail)                 │
│  ✓ List all posts and filter by author                           │
│  ✓ Check auto-generated slugs                                    │
└──────────────────────────────────────────────────────────────────┘
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  🎓 WHAT YOU LEARNED:                                                    │
│                                                                          │
│  • REST API design best practices                                       │
│  • JWT authentication and authorization                                 │
│  • Database integration with GORM                                       │
│  • Database migrations                                                  │
│  • CRUD operations and filtering                                        │
│  • Middleware for authentication                                        │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Incredible!** You've built a complete blog API with authentication!

---

*[Continues with remaining tutorials 5-15 in the same beautiful format...]*

Due to length constraints, I've shown the transformed format for the first several tutorials. The complete file would continue with:

- Tutorial 5: Job Queue System
- Tutorial 6: Rate Limiting
- Tutorial 7: Log Aggregation
- Tutorial 8: Service Mesh
- Tutorial 9: Time Series Database
- Tutorial 10: Container Orchestrator
- Specialized Tutorials (11-15)
- Learning paths and resources

Each following the same visual structure with boxes, emojis, clear sections, checkpoints, and engaging formatting. Would you like me to continue with the remaining tutorials?
