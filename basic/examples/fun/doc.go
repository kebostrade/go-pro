//go:build ignore

/*
Package main provides comprehensive examples of Go programming concepts.

# Variables
- Using := for short declaration (type inference)
- Using = for assignment
- Can be global (package level) or local (function level)

# Pointers
Variables that store memory addresses of other variables

# Constants
Immutable values defined at compile time

# Data Types and Structures
- Arrays: Fixed-size sequences
- Slices: Dynamic, flexible view into arrays
- Maps: Hash table implementation, key-value pairs
- Structs: Custom data types grouping related fields
- Custom Types: Define new types based on existing ones
- Type Embedding: Composition mechanism for types
- Interfaces: Define behavior through method signatures
- Generics: Write type-safe, reusable code (Go 1.18+)

# Concurrency Primitives
- Goroutines: Lightweight threads for concurrent execution
- Channels: Communication pipes between goroutines
- Mutexes: Mutual exclusion for shared resource access
- WaitGroups: Synchronization for multiple goroutines
- Select: Multi-way concurrent control structure
- Context: Cancellation, deadlines, and request-scoped values

# Memory and Safety
- Memory Management: Automatic memory allocation
- Garbage Collection: Automatic memory reclamation
- Unsafe: Low-level memory access (use with caution)
- Cgo: C code integration capabilities
- Error Handling: Explicit error handling with multiple return values
- Panic and Recover: Handling unexpected errors

# Project Organization
- Go Modules: Dependency management system
- Dependency Injection: Design pattern for loosely coupled code
- Composition: Building complex types from simpler ones
- Interfaces: Define behavior through method signatures
- Testing: Built-in testing framework

# Development Tools
- Testing: Built-in testing framework
- Benchmarking: Performance measurement tools
- Profiling: Runtime performance analysis
- Reflection: Runtime type inspection and manipulation
- Debugging: Tools for diagnosing issues
- Code Formatting: gofmt for consistent code style

# Control Structures
- For Loops: Single looping construct in Go
- Range: Iteration over arrays, slices, maps, etc.
- Break: Exit loops or switch statements
- Error Handling: Explicit error handling patterns
- Switch: Multi-way branching construct
- Defer: Delayed function execution
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Go Documentation Examples")
	fmt.Println("=========================")

	// Run various examples
	result := sum(1, 2)
	fmt.Println("Sum:", result)

	demoPointer()
	demoSlices()
	demoMaps()
	demoStructs()
	demoGoroutines()
}

// sum demonstrates basic function with parameters and return value
func sum(a int, b int) int {
	return a + b
}

// demoPointer demonstrates pointer usage
func demoPointer() {
	i := 42
	inc(&i)
	p := &i
	fmt.Println("Pointer value:", *p)
}

func inc(x *int) {
	*x++
}

// demoSlices demonstrates slice operations
func demoSlices() {
	arr := []int{1, 2, 3, 4, 5}
	arr = append(arr, 6)
	fmt.Println("Slice:", arr)
}

// demoMaps demonstrates map usage
func demoMaps() {
	vertices := make(map[string]int)
	vertices["a"] = 1
	vertices["b"] = 2
	vertices["c"] = 3
	fmt.Println("Map:", vertices)
	fmt.Println("Value at 'a':", vertices["a"])
}

// person struct demonstrates custom types
type person struct {
	name string
	age  int
}

// demoStructs demonstrates struct usage
func demoStructs() {
	p := person{name: "John", age: 30}
	fmt.Println("Person age:", p.age)
}

// demoGoroutines demonstrates concurrent execution
func demoGoroutines() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		countThings("sheep", &wg)
	}()

	go func() {
		countThings("fish", &wg)
	}()

	wg.Wait()
}

func countThings(thing string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i <= 3; i++ {
		fmt.Println(i, thing)
		time.Sleep(100 * time.Millisecond)
	}
}
