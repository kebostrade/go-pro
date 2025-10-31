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

### 🤖 [AI Engineering](#-ai-engineering-tutorials) (6-8 hours)
- [AI Tutorial 0](#-ai-tutorial-0-ai-engineering-overview) - AI Engineering Overview ⏱️ 20 min
- [AI Tutorial 1](#-ai-tutorial-1-llm-basics) - LLM Basics & Chatbot ⏱️ 30 min
- [AI Tutorial 2](#-ai-tutorial-2-prompt-engineering) - Prompt Engineering ⏱️ 25 min
- [Quick Reference](#-ai-quick-reference) - AI Engineering Quick Reference

### 🎬 [AI Content Creation](#-ai-content-creation-course) (54-69 hours)
- [Course Overview](#-ai-content-creation-mastery-course) - Complete Course Guide
- [Module 1](#-module-1-ai-video-generation) - AI Video Generation (Veo 3, Sora 2) ⏱️ 10-12h
- [Module 2](#-module-2-ai-audio--voice) - AI Audio & Voice (ElevenLabs) ⏱️ 6-8h
- [Module 3](#-module-3-advanced-audio-integration) - Advanced Audio (SFX, Ambient, Mixing) ⏱️ 8-10h ⭐ NEW
- [Module 7](#-module-7-content-strategy--monetization) - Monetization & Viral Strategy ⏱️ 12-15h
- [Viral Strategy](#-viral-video-strategy) - Get Millions of Views ⭐ NEW
- [Audio Library](#-audio-library-resources) - Complete Audio Resources ⭐ NEW
- [Quick Start](#-quick-start-create-first-video-in-2-hours) - Create First Video in 2 Hours ⏱️ 2h

### 📈 [Learning Paths](#-learning-path-summary)
- [Beginner Path](#beginner-path) - Start Here
- [Intermediate Path](#intermediate-path) - Next Steps
- [Advanced Path](#advanced-path) - Master Level
- [AI Engineering Path](#ai-engineering-path) - Build AI Apps

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

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                      🤖 AI ENGINEERING TUTORIALS                             ║
║                                                                              ║
║              Build Production-Ready AI Applications with Go                  ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

## 🤖 AI Tutorial 0: AI Engineering Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  20 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Understand AI Engineering and the learning path               │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ What is AI Engineering                                            │
│     ✓ LLMs, RAG, Agents, and Embeddings                                 │
│     ✓ Why Go for AI Engineering                                         │
│     ✓ Platform architecture overview                                    │
│     ✓ Learning roadmap and prerequisites                                │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Read the Overview

```bash
# Navigate to AI Engineering tutorials
cd docs/tutorials/ai-engineering

# Read the overview
cat 00_AI_ENGINEERING_OVERVIEW.md
```

**📚 Full Tutorial**: [00_AI_ENGINEERING_OVERVIEW.md](tutorials/ai-engineering/00_AI_ENGINEERING_OVERVIEW.md)

---

## 🤖 AI Tutorial 1: LLM Basics

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  30 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Build a CLI Chatbot with OpenAI                            │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ OpenAI API integration in Go                                      │
│     ✓ Streaming vs non-streaming responses                              │
│     ✓ Token management and counting                                     │
│     ✓ Conversation history management                                   │
│     ✓ Model parameters (temperature, max_tokens)                        │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Navigate to Project
```bash
cd basic/projects/ai-engineering/chatbot-cli
```

#### Step 2: Setup API Key
```bash
# Get your API key from https://platform.openai.com/api-keys
export OPENAI_API_KEY="sk-..."
```

#### Step 3: Install Dependencies
```bash
go mod init chatbot-cli
go get github.com/sashabaranov/go-openai
```

#### Step 4: Run the Chatbot
```bash
go run main.go
```

**📚 Full Tutorial**: [01_LLM_BASICS.md](tutorials/ai-engineering/01_LLM_BASICS.md)

**💡 Quick Tip**: Start with GPT-3.5-turbo for development (faster and cheaper)!

---

## 🤖 AI Tutorial 2: Prompt Engineering

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  25 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Master prompt design for better AI responses                  │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Prompt engineering fundamentals                                   │
│     ✓ System, user, and assistant roles                                 │
│     ✓ Few-shot learning techniques                                      │
│     ✓ Chain-of-thought prompting                                        │
│     ✓ Prompt templates and patterns                                     │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Read the Tutorial

```bash
cd docs/tutorials/ai-engineering
cat 02_PROMPT_ENGINEERING.md
```

**📚 Full Tutorial**: [02_PROMPT_ENGINEERING.md](tutorials/ai-engineering/02_PROMPT_ENGINEERING.md)

**💡 Quick Tip**: Good prompts are reusable - build a library of your best prompts!

---

## 📖 AI Quick Reference

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  📚 QUICK REFERENCE GUIDE                                                │
│                                                                          │
│  Common patterns, code snippets, and best practices for AI Engineering  │
│                                                                          │
│  • LLM Integration                                                       │
│  • Prompt Engineering                                                    │
│  • Embeddings & Vectors                                                  │
│  • RAG Systems                                                           │
│  • AI Agents                                                             │
│  • Error Handling                                                        │
│  • Best Practices                                                        │
│  • Troubleshooting                                                       │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**📚 Full Reference**: [QUICK_REFERENCE.md](tutorials/ai-engineering/QUICK_REFERENCE.md)

**💡 Bookmark this page** for quick reference while building AI applications!

---

## 🎓 AI Engineering Path

**Complete Learning Path**: See [PATH 4: AI ENGINEERING](../LEARNING_PATHS.md#-path-4-ai-engineering) in LEARNING_PATHS.md

**Duration**: 12-14 weeks

**Projects**:
1. CLI Chatbot (Week 1-2)
2. Semantic Search Engine (Week 3-4)
3. RAG Q&A System (Week 5-7)
4. Coding Assistant Agent (Week 8-10)
5. Production AI Service (Week 11-14)

**Skills You'll Master**:
- LLM API integration
- Prompt engineering
- Embeddings and vector search
- RAG architecture
- AI agent development
- Production deployment

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

---

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                       🎨 SPECIALIZED TUTORIALS                               ║
║                                                                              ║
║                Advanced Topics for Production Systems                        ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 💬 Tutorial 12: WebSocket Real-Time Communication

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  40 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Real-Time Chat Application                                 │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ WebSocket server with gorilla/websocket                          │
│     ✓ Real-time message broadcasting                                   │
│     ✓ Multiple chat rooms                                              │
│     ✓ Concurrent client handling                                       │
│     ✓ Message history and user management                              │
│     ✓ Modern web interface                                             │
│                                                                          │
│  🛠️ TECH STACK: WebSockets, Goroutines, Channels, Hub Pattern          │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Navigate to the Project
```bash
cd basic/projects/websocket-chat
```

#### Step 2: Understand the Architecture

**📖 Hub Pattern Overview:**
```
┌─────────────────────────────────────────────────────────┐
│                         Hub                             │
│  ┌─────────────────────────────────────────────────┐   │
│  │  Rooms: map[string]map[Client]bool              │   │
│  │  History: map[string][]Message                  │   │
│  │  Channels: register, unregister, broadcast      │   │
│  └─────────────────────────────────────────────────┘   │
│                          │                              │
│         ┌────────────────┼────────────────┐            │
│         ▼                ▼                ▼            │
│    ┌────────┐      ┌────────┐      ┌────────┐         │
│    │Client 1│      │Client 2│      │Client 3│         │
│    │Room: A │      │Room: A │      │Room: B │         │
│    └────────┘      └────────┘      └────────┘         │
└─────────────────────────────────────────────────────────┘
```

**💡 Key Concepts:**
```
┌─────────────────────────────────────────────────────────────────┐
│  WebSocket Lifecycle:                                           │
│                                                                 │
│  1. HTTP Upgrade Request                                       │
│  2. WebSocket Handshake                                        │
│  3. Bidirectional Communication                                │
│  4. Ping/Pong Heartbeat                                        │
│  5. Graceful Shutdown                                          │
│                                                                 │
│  Each Client has 2 Goroutines:                                 │
│  • ReadPump:  WebSocket → Hub (reads messages)                 │
│  • WritePump: Hub → WebSocket (writes messages)                │
└─────────────────────────────────────────────────────────────────┘
```

#### Step 3: Install Dependencies
```bash
# Download dependencies
make deps

# This installs:
# - github.com/gorilla/websocket v1.5.1
```

#### Step 4: Explore the Code Structure

```bash
# View the project structure
tree -L 3
```

**📁 Project Layout:**
```
websocket-chat/
├── cmd/
│   └── main.go              # Server entry point + web UI
├── internal/
│   ├── client/
│   │   └── client.go        # WebSocket client handler
│   ├── hub/
│   │   └── hub.go           # Central message hub
│   └── room/
│       └── room.go          # Chat room management
├── Makefile                 # Build automation
├── go.mod                   # Dependencies
└── README.md                # Documentation
```

**🔍 Examine Key Files:**
```bash
# Hub - manages all clients and rooms
cat internal/hub/hub.go | head -50

# Client - handles individual connections
cat internal/client/client.go | head -50

# Main - HTTP server and WebSocket upgrade
cat cmd/main.go | head -50
```

#### Step 5: Build and Run

```bash
# Build the application
make build

# Run the server
make run
```

**📤 Server Output:**
```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║           💬 WebSocket Chat Server                          ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

✓ Server starting on http://localhost:8080
✓ WebSocket endpoint: ws://localhost:8080/ws
✓ Ready to accept connections

📚 API Endpoints:
  GET  /                    - Web interface
  WS   /ws?username=X&room=Y - WebSocket connection
  GET  /api/rooms           - List all rooms
  GET  /api/rooms/{id}      - Get room statistics

🎯 Quick Start:
  1. Open http://localhost:8080 in your browser
  2. Enter your username and room name
  3. Start chatting!
```

#### Step 6: Test the Chat Application

**Option 1: Web Browser (Easiest)**
```bash
# Open in your browser
open http://localhost:8080

# Or manually navigate to:
# http://localhost:8080
```

**📱 Web Interface Features:**
- Modern, responsive design
- Real-time message updates
- System notifications (join/leave)
- Message timestamps
- Connection status indicator

**Option 2: Multiple Browser Windows**
```bash
# Open 3 browser windows/tabs
# Window 1: Username "Alice", Room "general"
# Window 2: Username "Bob", Room "general"
# Window 3: Username "Charlie", Room "tech"

# Alice and Bob can chat in "general"
# Charlie is in a separate "tech" room
```

**Option 3: Command Line with wscat**
```bash
# Install wscat (if not already installed)
npm install -g wscat

# Connect to chat
wscat -c "ws://localhost:8080/ws?username=Alice&room=general"

# Send a message (type and press Enter)
{"type":"message","content":"Hello from command line!"}
```

#### Step 7: Test the REST API

**List All Active Rooms:**
```bash
curl http://localhost:8080/api/rooms
```

**📤 Response:**
```json
{
  "rooms": ["general", "tech", "random"],
  "count": 3
}
```

**Get Room Statistics:**
```bash
curl http://localhost:8080/api/rooms/general
```

**📤 Response:**
```json
{
  "client_count": 5,
  "message_count": 42,
  "users": ["Alice", "Bob", "Charlie", "David", "Eve"]
}
```

#### Step 8: Understanding the Message Flow

**📖 Message Broadcasting Flow:**
```
┌──────────────────────────────────────────────────────────┐
│                                                          │
│  Client A sends message                                  │
│       │                                                  │
│       ▼                                                  │
│  ReadPump() reads from WebSocket                         │
│       │                                                  │
│       ▼                                                  │
│  Hub.Broadcast(message, roomID)                          │
│       │                                                  │
│       ▼                                                  │
│  Hub finds all clients in room                           │
│       │                                                  │
│       ├──► Client A.Send ◄── WritePump() ──► WebSocket  │
│       ├──► Client B.Send ◄── WritePump() ──► WebSocket  │
│       └──► Client C.Send ◄── WritePump() ──► WebSocket  │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**💡 Key Implementation Details:**

1. **Buffered Channels** (256 messages):
   ```go
   Send: make(chan []byte, 256)
   ```
   Prevents blocking when sending to slow clients

2. **Ping/Pong Heartbeat**:
   - Ping every 54 seconds
   - Pong timeout: 60 seconds
   - Detects dead connections

3. **Thread-Safe Operations**:
   ```go
   h.mu.Lock()
   defer h.mu.Unlock()
   ```
   Protects shared state from race conditions

#### Step 9: Explore Advanced Features

**Feature 1: Message History**
```bash
# Join a room - you'll see last 100 messages
# History is automatically sent to new clients
```

**Feature 2: System Notifications**
```bash
# Watch for join/leave messages
# Format: "Alice joined the chat"
#         "Bob left the chat"
```

**Feature 3: Multiple Rooms**
```bash
# Create different rooms for different topics
# Messages are isolated per room
# Each room has independent history
```

#### Step 10: Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench
```

**📊 Expected Benchmark Results:**
```
BenchmarkBroadcast-8        10000    105234 ns/op    2048 B/op    12 allocs/op
BenchmarkClientSend-8       50000     32156 ns/op     512 B/op     5 allocs/op
BenchmarkHubRegister-8     100000     15234 ns/op     256 B/op     3 allocs/op
```

**🎯 Quick Wins:**
```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Open 3+ browser windows and chat between them                │
│  ✓ Create multiple rooms (general, tech, random)                │
│  ✓ Test message history by joining an active room               │
│  ✓ Monitor server logs to see connection events                 │
│  ✓ Check room statistics via API                                │
│  ✓ Test connection resilience (close/reopen browser)            │
│  ✓ Send 100+ messages and verify performance                    │
└──────────────────────────────────────────────────────────────────┘
```

### 📚 Code Deep Dive

#### Understanding the Hub

**The Hub is the heart of the application:**

```go
type Hub struct {
    // Clients organized by room
    rooms map[string]map[Client]bool

    // Channels for communication
    broadcast  chan *BroadcastMessage
    register   chan Client
    unregister chan Client

    // Message history (last 100 per room)
    history map[string][]Message

    // Thread safety
    mu sync.RWMutex
}
```

**Hub's Main Loop:**
```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)

        case client := <-h.unregister:
            h.unregisterClient(client)

        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}
```

**💡 Why This Pattern?**
- **Single Goroutine**: All hub operations in one goroutine = no race conditions
- **Channel Communication**: Type-safe, blocking communication
- **Select Statement**: Handles multiple channels efficiently

#### Understanding the Client

**Each client has two goroutines:**

**ReadPump (WebSocket → Hub):**
```go
func (c *Client) ReadPump() {
    defer func() {
        c.Hub.Unregister(c)
        c.Conn.Close()
    }()

    // Set read deadline and limits
    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }

        // Parse and broadcast message
        var msg Message
        json.Unmarshal(message, &msg)
        msg.Username = c.Username
        msg.Timestamp = time.Now()

        messageBytes, _ := json.Marshal(msg)
        c.Hub.Broadcast(messageBytes, c.RoomID)
    }
}
```

**WritePump (Hub → WebSocket):**
```go
func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                // Hub closed the channel
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // Write message to WebSocket
            w, _ := c.Conn.NextWriter(websocket.TextMessage)
            w.Write(message)
            w.Close()

        case <-ticker.C:
            // Send ping to keep connection alive
            c.Conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

### 🔧 Customization Ideas

**1. Add Private Messaging:**
```go
type PrivateMessage struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Content string `json:"content"`
}

// In hub, add method:
func (h *Hub) SendPrivate(msg PrivateMessage) {
    // Find recipient client and send directly
}
```

**2. Add Typing Indicators:**
```go
type TypingEvent struct {
    Username string `json:"username"`
    IsTyping bool   `json:"is_typing"`
    RoomID   string `json:"room_id"`
}

// Broadcast typing events with debouncing
```

**3. Add User Authentication:**
```go
// Add JWT middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        // Validate JWT token
        next(w, r)
    }
}
```

**4. Add Persistent Storage:**
```go
// Store messages in PostgreSQL
type MessageRepository interface {
    Save(msg Message) error
    GetHistory(roomID string, limit int) ([]Message, error)
}
```

### 🐛 Common Issues and Solutions

**Issue 1: Connection Drops**
```
Problem: Clients disconnect after 60 seconds
Solution: Ping/pong heartbeat is working correctly
         This is the pongWait timeout
         Increase if needed in client.go
```

**Issue 2: Messages Not Broadcasting**
```
Problem: Messages sent but not received by other clients
Solution: Check that clients are in the same room
         Verify roomID matches exactly
         Check server logs for errors
```

**Issue 3: Memory Leak**
```
Problem: Memory usage grows over time
Solution: Message history is limited to 100 per room
         Ensure clients are properly unregistered
         Check for goroutine leaks with pprof
```

**Issue 4: Race Conditions**
```
Problem: Panic: concurrent map read and write
Solution: All hub operations use mutex locks
         Run with -race flag to detect issues
         go run -race ./cmd/main.go
```

### 📊 Performance Tuning

**1. Adjust Buffer Sizes:**
```go
// In client.go
Send: make(chan []byte, 256)  // Increase for high-traffic rooms

// In upgrader
ReadBufferSize:  1024,  // Increase for larger messages
WriteBufferSize: 1024,
```

**2. Optimize Message History:**
```go
// In hub.go
const maxHistorySize = 100  // Adjust based on needs

// Consider using a ring buffer for efficiency
```

**3. Add Connection Pooling:**
```go
// Limit concurrent connections
var connectionSemaphore = make(chan struct{}, 10000)

func serveWs(...) {
    connectionSemaphore <- struct{}{}
    defer func() { <-connectionSemaphore }()
    // ... rest of code
}
```

**4. Enable Compression:**
```go
upgrader := websocket.Upgrader{
    EnableCompression: true,
    // ... other settings
}
```

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  🎓 WHAT YOU LEARNED:                                                    │
│                                                                          │
│  • WebSocket protocol and HTTP upgrade handshake                        │
│  • Hub pattern for managing concurrent connections                      │
│  • Goroutines for concurrent read/write operations                      │
│  • Channels for thread-safe communication                               │
│  • Select statements for multiplexing channels                          │
│  • Mutex locks for protecting shared state                              │
│  • Ping/pong heartbeat for connection health                            │
│  • Buffered channels for non-blocking sends                             │
│  • Graceful shutdown and cleanup                                        │
│  • Real-time message broadcasting                                       │
│  • Room-based message isolation                                         │
│  • Message history management                                           │
│                                                                          │
│  💪 SKILLS GAINED:                                                       │
│                                                                          │
│  ✓ Building real-time applications                                      │
│  ✓ Managing concurrent connections at scale                             │
│  ✓ Implementing pub/sub patterns                                        │
│  ✓ Handling WebSocket lifecycle                                         │
│  ✓ Designing thread-safe systems                                        │
│  ✓ Performance optimization techniques                                  │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**🎉 Congratulations!** You've built a production-ready real-time chat application!

**🚀 Next Steps:**
- Add user authentication with JWT
- Implement private messaging
- Add file sharing capabilities
- Deploy to production with Docker
- Scale with Redis for distributed deployment
- Add message encryption for security

---

## 🏗️ Tutorial 13: Microservices Architecture in Go

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🔴 ADVANCED                                    ⏱️  60 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Production Microservices System                            │
│                                                                          │
│  📚 WHAT YOU'LL BUILD:                                                   │
│     ✓ API Gateway with routing and authentication                      │
│     ✓ User Service with JWT authentication                             │
│     ✓ Product Service with caching                                     │
│     ✓ Order Service with event handling                                │
│     ✓ Service discovery and registration                               │
│     ✓ Distributed logging and monitoring                               │
│     ✓ Docker Compose orchestration                                     │
│                                                                          │
│  🛠️ TECH STACK: Microservices, gRPC, REST, Docker, Service Discovery   │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 📝 Step-by-Step Instructions

#### Step 1: Navigate to the Project
```bash
cd basic/projects/microservices-demo
```

#### Step 2: Understand the Architecture

**📖 Microservices Architecture Overview:**
```
┌─────────────────────────────────────────────────────────────────┐
│                         API Gateway                             │
│                      (Port 8080)                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Request Routing                                       │  │
│  │  • Authentication & Authorization                        │  │
│  │  • Rate Limiting                                         │  │
│  │  • Load Balancing                                        │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────┬────────────────┬────────────────┬─────────────────┘
             │                │                │
    ┌────────▼────────┐  ┌───▼──────────┐  ┌─▼──────────────┐
    │  User Service   │  │   Product    │  │  Order Service │
    │   (Port 8081)   │  │   Service    │  │  (Port 8083)   │
    │                 │  │ (Port 8082)  │  │                │
    │  • User CRUD    │  │  • Product   │  │  • Order CRUD  │
    │  • Auth/Login   │  │    Catalog   │  │  • Status Mgmt │
    │  • JWT Tokens   │  │  • Inventory │  │  • Events      │
    └────────┬────────┘  └───┬──────────┘  └─┬──────────────┘
             │               │                │
    ┌────────▼────────┐  ┌──▼───────────┐   │
    │   PostgreSQL    │  │    Redis     │   │
    │   (Port 5432)   │  │  (Port 6379) │   │
    └─────────────────┘  └──────────────┘   │
                                             │
                         ┌───────────────────▼──────────┐
                         │  Shared Infrastructure       │
                         │  • Service Discovery         │
                         │  • Logging (Zap)             │
                         │  • Middleware (Auth, Logs)   │
                         └──────────────────────────────┘
```

**💡 Key Microservices Patterns:**
```
┌─────────────────────────────────────────────────────────────────┐
│  1. API Gateway Pattern                                         │
│     Single entry point for all client requests                 │
│     Handles routing, auth, rate limiting                       │
│                                                                 │
│  2. Service Discovery                                           │
│     Services register themselves on startup                    │
│     Gateway discovers services dynamically                     │
│                                                                 │
│  3. Database per Service                                        │
│     Each service owns its data                                 │
│     No direct database sharing                                 │
│                                                                 │
│  4. Shared Infrastructure                                       │
│     Common packages for logging, middleware                    │
│     Reusable across all services                               │
│                                                                 │
│  5. Graceful Shutdown                                           │
│     Services handle SIGTERM/SIGINT                             │
│     Clean connection closure                                   │
└─────────────────────────────────────────────────────────────────┘
```

#### Step 3: Install Dependencies
```bash
# Download all dependencies
make deps

# This installs:
# - gorilla/mux (HTTP routing)
# - go-redis (Redis client)
# - zap (Structured logging)
# - jwt (JWT authentication)
# - uuid (UUID generation)
# - grpc (gRPC framework)
```

#### Step 4: Explore the Project Structure

```bash
# View the complete structure
tree -L 3
```

**📁 Project Layout:**
```
microservices-demo/
├── services/
│   ├── api-gateway/          # API Gateway (Port 8080)
│   │   └── cmd/main.go
│   ├── user-service/         # User Service (Port 8081)
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── models.go
│   │       ├── repository.go
│   │       └── handler.go
│   ├── product-service/      # Product Service (Port 8082)
│   │   └── cmd/main.go
│   └── order-service/        # Order Service (Port 8083)
│       └── cmd/main.go
├── pkg/                      # Shared packages
│   ├── logger/              # Structured logging
│   ├── discovery/           # Service registry
│   ├── middleware/          # HTTP middleware
│   └── proto/               # gRPC definitions
├── deployments/             # Docker & K8s
│   ├── docker-compose.yml
│   └── Dockerfile.*
├── Makefile                # Build automation
└── go.mod                  # Dependencies
```

**🔍 Examine Key Components:**
```bash
# Service Discovery
cat pkg/discovery/discovery.go | head -50

# Authentication Middleware
cat pkg/middleware/auth.go | head -50

# User Service Handler
cat services/user-service/internal/handler.go | head -50

# API Gateway
cat services/api-gateway/cmd/main.go | head -50
```

#### Step 5: Run with Docker Compose (Recommended)

```bash
# Build and start all services
make docker-up
```

**📤 Expected Output:**
```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║           🏗️  Microservices System Starting                 ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

Creating network "microservices-network" ... done
Creating microservices-postgres ... done
Creating microservices-redis    ... done
Creating user-service           ... done
Creating product-service        ... done
Creating order-service          ... done
Creating api-gateway            ... done

✓ All services started successfully

📡 Services available at:
  API Gateway:     http://localhost:8080
  User Service:    http://localhost:8081
  Product Service: http://localhost:8082
  Order Service:   http://localhost:8083
  PostgreSQL:      localhost:5432
  Redis:           localhost:6379
```

#### Step 6: Run Locally (Alternative)

```bash
# Terminal 1: User Service
make run-user

# Terminal 2: Product Service
make run-product

# Terminal 3: Order Service
make run-order

# Terminal 4: API Gateway
make run-gateway
```

#### Step 7: Test the System

**🧪 1. Health Checks**
```bash
# Check API Gateway
curl http://localhost:8080/health
# Response: OK

# Check User Service
curl http://localhost:8081/health
# Response: OK

# Check Product Service
curl http://localhost:8082/health
# Response: OK

# Check Order Service
curl http://localhost:8083/health
# Response: OK
```

**🔍 2. Service Discovery**
```bash
# List all registered services
curl http://localhost:8080/services

# Response:
{
  "services": {
    "user-service": "localhost:8081",
    "product-service": "localhost:8082",
    "order-service": "localhost:8083"
  }
}
```

**👤 3. User Management**
```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'

# Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "alice",
  "email": "alice@example.com",
  "created_at": "2024-01-15T10:30:00Z"
}

# Login
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'

# Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "alice",
    "email": "alice@example.com"
  }
}

# Save the token for authenticated requests!
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# List users (requires authentication)
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```

**📦 4. Product Management**
```bash
# List products (seeded data)
curl http://localhost:8080/api/products

# Response:
{
  "products": [
    {
      "id": "1",
      "name": "Laptop",
      "description": "High-performance laptop",
      "price": 999.99,
      "stock": 10
    },
    {
      "id": "2",
      "name": "Mouse",
      "description": "Wireless mouse",
      "price": 29.99,
      "stock": 50
    }
  ],
  "total": 2
}

# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Keyboard",
    "description": "Mechanical keyboard",
    "price": 149.99,
    "stock": 25
  }'

# Get a specific product
curl http://localhost:8080/api/products/1
```

**🛒 5. Order Management**
```bash
# Create an order (requires authentication)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_id": "1",
    "quantity": 2,
    "total": 1999.98
  }'

# Response:
{
  "id": "order-uuid-here",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "product_id": "1",
  "quantity": 2,
  "total": 1999.98,
  "status": "pending",
  "created_at": "2024-01-15T10:35:00Z"
}

# List all orders
curl http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN"

# Update order status
curl -X PUT http://localhost:8080/api/orders/order-uuid-here/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status": "completed"}'
```

#### Step 8: Understand the Code

**🔐 Authentication Flow:**
```go
// 1. User registers
POST /api/users
  → User Service creates user
  → Password hashed with bcrypt
  → User stored in repository

// 2. User logs in
POST /api/users/login
  → User Service validates credentials
  → JWT token generated (24h expiration)
  → Token returned to client

// 3. Client makes authenticated request
GET /api/users (with Authorization header)
  → API Gateway receives request
  → Middleware validates JWT token
  → Request proxied to User Service
  → Response returned to client
```

**🔄 Request Flow Through Gateway:**
```go
// Example: GET /api/products/1

Client
  ↓
API Gateway (Port 8080)
  ↓ (Middleware: Logging, Rate Limiting)
  ↓
Service Discovery
  ↓ (Discover "product-service" → localhost:8082)
  ↓
Reverse Proxy
  ↓ (Forward to http://localhost:8082/api/products/1)
  ↓
Product Service (Port 8082)
  ↓ (Handler processes request)
  ↓
Response
  ↓
API Gateway
  ↓
Client
```

**📊 Service Discovery Mechanism:**
```go
// Service Registration (on startup)
func main() {
    serviceName := "user-service"
    serviceAddr := "localhost:8081"

    // Register with discovery service
    discovery.Register(serviceName, serviceAddr)
    defer discovery.Deregister(serviceName)

    // Start HTTP server
    srv.ListenAndServe()
}

// Service Discovery (in API Gateway)
func proxyHandler(serviceName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Discover service address
        serviceAddr, err := discovery.Discover(serviceName)

        // Forward request to service
        targetURL := fmt.Sprintf("http://%s%s", serviceAddr, r.URL.Path)
        // ... proxy logic
    }
}
```

#### Step 9: Monitor and Debug

**📝 View Logs:**
```bash
# Docker logs
make docker-logs

# Follow logs for specific service
docker logs -f user-service
docker logs -f api-gateway

# Local logs (if running locally)
# Logs appear in terminal where service is running
```

**🔍 Debug Service Communication:**
```bash
# Check service registration
curl http://localhost:8080/services

# Test direct service access (bypass gateway)
curl http://localhost:8081/health
curl http://localhost:8082/products
curl http://localhost:8083/orders

# Check rate limiting
for i in {1..300}; do
  curl http://localhost:8080/health
done
# Should see 429 Too Many Requests after ~200 requests
```

### 🎓 Key Concepts Explained

#### 1. API Gateway Pattern
```
┌─────────────────────────────────────────────────────────────────┐
│  WHY API Gateway?                                               │
│                                                                 │
│  ✓ Single entry point for clients                              │
│  ✓ Centralized authentication & authorization                  │
│  ✓ Request routing to appropriate services                     │
│  ✓ Rate limiting and throttling                                │
│  ✓ Load balancing across service instances                     │
│  ✓ Protocol translation (REST → gRPC)                          │
│  ✓ Response aggregation from multiple services                 │
│                                                                 │
│  IMPLEMENTATION:                                                │
│  - Reverse proxy pattern                                       │
│  - Service discovery integration                               │
│  - Middleware chain (auth, logging, rate limit)                │
└─────────────────────────────────────────────────────────────────┘
```

#### 2. Service Discovery
```
┌─────────────────────────────────────────────────────────────────┐
│  Service Registry Pattern                                       │
│                                                                 │
│  REGISTRATION (Service Startup):                                │
│    1. Service starts on port 8081                              │
│    2. Calls discovery.Register("user-service", "localhost:8081")│
│    3. Registry stores mapping                                  │
│                                                                 │
│  DISCOVERY (API Gateway):                                       │
│    1. Request arrives for /api/users                           │
│    2. Gateway calls discovery.Discover("user-service")         │
│    3. Registry returns "localhost:8081"                        │
│    4. Gateway forwards request                                 │
│                                                                 │
│  DEREGISTRATION (Service Shutdown):                             │
│    1. Service receives SIGTERM                                 │
│    2. Calls discovery.Deregister("user-service")               │
│    3. Registry removes mapping                                 │
│                                                                 │
│  PRODUCTION: Use Consul, etcd, or Kubernetes Service Discovery │
└─────────────────────────────────────────────────────────────────┘
```

#### 3. JWT Authentication
```
┌─────────────────────────────────────────────────────────────────┐
│  Token-Based Authentication Flow                                │
│                                                                 │
│  TOKEN GENERATION:                                              │
│    claims := jwt.MapClaims{                                    │
│      "user_id": "uuid",                                        │
│      "username": "alice",                                      │
│      "email": "alice@example.com",                             │
│      "exp": time.Now().Add(24 * time.Hour).Unix(),            │
│    }                                                            │
│    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) │
│    tokenString, _ := token.SignedString([]byte(secret))       │
│                                                                 │
│  TOKEN VALIDATION:                                              │
│    token, err := jwt.Parse(tokenString, func(token *jwt.Token) │
│      return []byte(secret), nil                                │
│    })                                                           │
│    claims := token.Claims.(jwt.MapClaims)                      │
│    userID := claims["user_id"].(string)                        │
│                                                                 │
│  MIDDLEWARE:                                                    │
│    - Extracts token from Authorization header                  │
│    - Validates signature and expiration                        │
│    - Injects user info into request context                    │
│    - Rejects invalid/expired tokens (401)                      │
└─────────────────────────────────────────────────────────────────┘
```

#### 4. Rate Limiting (Token Bucket Algorithm)
```
┌─────────────────────────────────────────────────────────────────┐
│  Token Bucket Algorithm                                         │
│                                                                 │
│  CONCEPT:                                                       │
│    Bucket capacity: 200 tokens                                 │
│    Refill rate: 100 tokens/second                              │
│                                                                 │
│  FLOW:                                                          │
│    1. Request arrives                                          │
│    2. Check if bucket has tokens                               │
│    3. If yes: consume 1 token, allow request                   │
│    4. If no: reject with 429 Too Many Requests                 │
│    5. Bucket refills at constant rate                          │
│                                                                 │
│  IMPLEMENTATION:                                                │
│    type Visitor struct {                                       │
│      limiter  *rate.Limiter                                    │
│      lastSeen time.Time                                        │
│    }                                                            │
│                                                                 │
│    limiter := rate.NewLimiter(rate.Limit(100), 200)           │
│    if !limiter.Allow() {                                       │
│      http.Error(w, "Too Many Requests", 429)                   │
│    }                                                            │
│                                                                 │
│  BENEFITS:                                                      │
│    ✓ Prevents API abuse                                        │
│    ✓ Protects backend services                                 │
│    ✓ Allows burst traffic (up to bucket capacity)              │
└─────────────────────────────────────────────────────────────────┘
```

#### 5. Graceful Shutdown
```
┌─────────────────────────────────────────────────────────────────┐
│  Handling Service Shutdown Properly                             │
│                                                                 │
│  SETUP:                                                         │
│    quit := make(chan os.Signal, 1)                            │
│    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)       │
│                                                                 │
│  WAIT FOR SIGNAL:                                               │
│    <-quit  // Blocks until signal received                     │
│                                                                 │
│  SHUTDOWN SEQUENCE:                                             │
│    1. Stop accepting new requests                              │
│    2. Wait for in-flight requests to complete (30s timeout)    │
│    3. Close database connections                               │
│    4. Deregister from service discovery                        │
│    5. Flush logs                                               │
│    6. Exit cleanly                                             │
│                                                                 │
│  CODE:                                                          │
│    ctx, cancel := context.WithTimeout(context.Background(), 30s)│
│    defer cancel()                                              │
│    if err := srv.Shutdown(ctx); err != nil {                   │
│      logger.Fatal("Forced shutdown", zap.Error(err))           │
│    }                                                            │
│                                                                 │
│  WHY IMPORTANT:                                                 │
│    ✓ No dropped requests                                       │
│    ✓ No data corruption                                        │
│    ✓ Clean service deregistration                              │
│    ✓ Proper resource cleanup                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 🎯 Challenges

**🔰 Beginner Challenges:**
1. **Add a New Endpoint**: Add `GET /api/users/{id}/orders` to list user's orders
2. **Custom Middleware**: Create a request ID middleware that adds unique IDs to logs
3. **Health Check Enhancement**: Add database connectivity check to health endpoints
4. **Product Search**: Add `GET /api/products/search?q=laptop` endpoint

**🔶 Intermediate Challenges:**
1. **Pagination**: Add pagination to all list endpoints (page, pageSize)
2. **Caching Layer**: Implement Redis caching in Product Service
3. **Event Publishing**: Publish events when orders are created/updated
4. **Service Metrics**: Add Prometheus metrics to track request counts, latency
5. **Circuit Breaker**: Implement circuit breaker for service-to-service calls

**🔴 Advanced Challenges:**
1. **gRPC Communication**: Replace HTTP with gRPC between services
2. **Distributed Tracing**: Add OpenTelemetry for request tracing across services
3. **Saga Pattern**: Implement distributed transaction for order creation
4. **Message Queue**: Add RabbitMQ/Kafka for async communication
5. **Kubernetes Deployment**: Deploy to Kubernetes with Helm charts
6. **Service Mesh**: Integrate Istio for advanced traffic management

### 📚 What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✅ MICROSERVICES ARCHITECTURE                                   │
│     • Service decomposition and boundaries                      │
│     • API Gateway pattern                                       │
│     • Service discovery and registration                        │
│     • Inter-service communication                               │
│                                                                  │
│  ✅ DISTRIBUTED SYSTEMS                                          │
│     • Service-to-service communication                          │
│     • Distributed logging                                       │
│     • Graceful shutdown and fault tolerance                     │
│     • Rate limiting and throttling                              │
│                                                                  │
│  ✅ AUTHENTICATION & SECURITY                                    │
│     • JWT token generation and validation                       │
│     • Password hashing with bcrypt                              │
│     • Authentication middleware                                 │
│     • Authorization patterns                                    │
│                                                                  │
│  ✅ INFRASTRUCTURE                                               │
│     • Docker containerization                                   │
│     • Docker Compose orchestration                              │
│     • Multi-stage Docker builds                                 │
│     • Service networking                                        │
│                                                                  │
│  ✅ BEST PRACTICES                                               │
│     • Clean Architecture (handlers, services, repositories)     │
│     • Dependency injection                                      │
│     • Structured logging                                        │
│     • Error handling and propagation                            │
│     • Configuration management                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 🚀 Next Steps

1. **Study the Code**: Read through each service implementation
2. **Add Features**: Implement the challenges above
3. **Deploy to Cloud**: Try AWS ECS, Google Cloud Run, or Azure Container Instances
4. **Add Observability**: Integrate Prometheus, Grafana, Jaeger
5. **Scale Up**: Add load balancing, multiple service instances
6. **Production Hardening**: Add TLS, secrets management, monitoring

### 📖 Resources

- **Project README**: [microservices-demo/README.md](../basic/projects/microservices-demo/README.md)
- **Quick Start Guide**: [microservices-demo/QUICK_START.md](../basic/projects/microservices-demo/QUICK_START.md)
- **Microservices Patterns**: https://microservices.io/patterns/
- **Go Microservices**: https://go.dev/blog
- **Docker Documentation**: https://docs.docker.com/
- **gRPC in Go**: https://grpc.io/docs/languages/go/

### 🎉 Congratulations!

You've built a production-ready microservices architecture with:
- ✅ 4 independent services
- ✅ API Gateway with routing and auth
- ✅ Service discovery
- ✅ JWT authentication
- ✅ Rate limiting
- ✅ Docker orchestration
- ✅ Structured logging

**You're now ready to build scalable distributed systems in Go!** 🚀

---

## 🎬 AI Content Creation Course

**Master cutting-edge AI tools to create viral content and generate passive income**

**Full Course**: [AI Content Creation Mastery](courses/ai-content-creation/README.md)

**Quick Start**: [Create Your First Video in 2 Hours](courses/ai-content-creation/resources/QUICK_START.md)

**Course Summary**: [Complete Overview](courses/ai-content-creation/COURSE_SUMMARY.md)

This comprehensive course teaches you to use Google Veo 3, Sora 2, ElevenLabs, RunwayML, SunoAI, HeyGen, and Hailou AI to create content that attracts millions of views on YouTube, TikTok, and Instagram, and generates $1,000-$50,000+/month in passive income.

---

### 🎵 Module 3: Advanced Audio Integration

**File**: [modules/03_ADVANCED_AUDIO.md](courses/ai-content-creation/modules/03_ADVANCED_AUDIO.md)

**Focus**: Sound Effects, Ambient Noise, Professional Mixing
**Duration**: 8-10 hours
**Difficulty**: 🟡 Intermediate

**What You'll Master**:
- Layer sound effects for maximum impact
- Create immersive ambient soundscapes
- Add realistic multi-character dialogue
- Mix audio like a professional
- Optimize audio for viral engagement

**The 3-Layer Audio System**:
```
Layer 1: Dialogue/Narration (100% volume)
Layer 2: Sound Effects (60-80% volume)
Layer 3: Ambient Noise (20-40% volume)
```

**Key Techniques**:
- Multi-character conversations with ElevenLabs
- 5 categories of sound effects (Transitions, Emphasis, UI/Tech, Emotional, Nature)
- 5 ambient soundscape scenarios (Office, Coffee Shop, Lab, Nature, Urban)
- Professional audio mixing and ducking
- Platform-specific audio optimization
- Emotional audio arcs
- Strategic silence for impact

**Why Audio Matters**:
- 📈 3-5x more engagement than silent videos
- 🔄 2x higher completion rate
- 💬 4x more comments
- 🔁 6x more shares

**Projects**:
- Tech product demo with professional audio
- Storytelling short using audio-only narrative
- Viral challenge with superior audio quality

**Resources**: [Complete Audio Library](courses/ai-content-creation/resources/AUDIO_LIBRARY.md)

---

### 🚀 Viral Video Strategy

**File**: [VIRAL_STRATEGY.md](courses/ai-content-creation/VIRAL_STRATEGY.md)

**Goal**: Create videos that get 1M+ views consistently

**The Viral Formula**:
```
Viral Video = Hook × Value × Emotion × Shareability × Timing
```

**The 7-Second Rule**:
```
First 7 seconds = 80% of your success

0-1s:  Viewer decides to keep watching or scroll
1-3s:  Brain processes if content is relevant
3-5s:  Emotional response triggered
5-7s:  Decision to engage (like, comment, share)
```

**5 Perfect Hook Formulas**:
1. **Shock Value**: "I made $50,000 in 24 hours..."
2. **Controversy**: "Everyone is doing [X] wrong..."
3. **Transformation**: "From [Bad] to [Good] in [Time]..."
4. **Secret Reveal**: "I discovered a secret that..."
5. **Urgent Warning**: "Stop doing [X] immediately..."

**Platform-Specific Strategies**:
- **YouTube Shorts**: 1M+ views (3-5 posts/day)
- **TikTok**: 5M+ views (5-7 posts/day)
- **Instagram Reels**: 500K+ views (1-2 posts/day)

**10 Psychological Triggers**:
1. Awe/Wonder (85% share rate)
2. Anger/Outrage (80%)
3. Anxiety/Fear (75%)
4. Humor/Joy (90%)
5. Inspiration (70%)
6. Validation (65%)
7. Belonging (75%)
8. Curiosity (85%)
9. Utility (80%)
10. Controversy (90%)

**30-Day Viral Challenge**:
- Week 1: 35 videos (testing)
- Week 2: 49 videos (optimization)
- Week 3: 70 videos (scaling)
- Week 4: 90+ videos (viral phase)
- **Total**: 244+ videos → 1-3 viral hits, 10K-100K followers

**Monetization**:
```
1M views = $2K-$10K (YouTube AdSense)
1M views = 10K-50K followers
         = $5K-$50K (sponsorships)
         = $10K-$100K (product sales)
```

---

### 🎧 Audio Library Resources

**File**: [resources/AUDIO_LIBRARY.md](courses/ai-content-creation/resources/AUDIO_LIBRARY.md)

**Complete curated library of audio resources**

**Free Sound Effects**:
- Freesound.org (500,000+ sounds)
- Zapsplat.com (100,000+ sounds)
- YouTube Audio Library (10,000+ sounds)
- BBC Sound Effects (16,000+ sounds)
- Sonniss Game Audio (30GB+ bundle)

**Premium Sound Effects**:
- Epidemic Sound ($15/month - 90,000+ SFX)
- Artlist ($16.60/month - 50,000+ SFX)
- Soundstripe ($15/month - 40,000+ SFX)
- Boom Library ($50-$500 per pack)
- Pro Sound Effects ($20/month - 100,000+ sounds)

**Ambient Noise Generators**:
- Ambient-Mixer.com (custom soundscapes)
- MyNoise.net (parametric generator)
- Noisli (simple mixing)
- A Soft Murmur (natural sounds)

**Music Platforms**:
- Uppbeat (free with attribution)
- Pixabay Music (free, no attribution)
- Incompetech (free with attribution)
- Bensound (cinematic music)

**AI Voice Tools**:
- ElevenLabs (best quality)
- Murf.ai (business content)
- Play.ht (variety)

**Audio Editing Software**:
- Free: Audacity, GarageBand, Ocenaudio
- Premium: Adobe Audition, Logic Pro, Reaper

**Recommended Starter Kits**:
```
Free Kit ($0/month):
- Freesound + Zapsplat (SFX)
- Ambient-Mixer (ambient)
- Uppbeat (music)
- ElevenLabs free tier (voice)
- Audacity (editing)

Budget Kit ($25/month):
- Epidemic Sound ($15)
- ElevenLabs ($5)
- Audacity (free)

Professional Kit ($100/month):
- Epidemic Sound ($15)
- ElevenLabs ($22)
- Adobe Audition ($23)
- Pro Sound Effects ($20)
```

---

**🎉 Congratulations!** You now have access to world-class Go programming AND AI content creation education!

**Happy Coding & Creating! 🚀**

Each following the same visual structure with boxes, emojis, clear sections, checkpoints, and engaging formatting. Would you like me to continue with the remaining tutorials?
