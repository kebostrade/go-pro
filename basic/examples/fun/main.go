package main

import (
	"fmt"
	"os"
)

func main() {
	// Dispatcher for running examples
	if len(os.Args) > 1 && os.Args[1] == "fizzbuzz" {
		runFizzBuzzDemo()
		return
	}

	// Default: show menu
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                            ║")
	fmt.Println("║        Go Programming Examples - Fun Directory             ║")
	fmt.Println("║                                                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📚 How to use:")
	fmt.Println()
	fmt.Println("1. FizzBuzz String Function:")
	fmt.Println("   go run . fizzbuzz")
	fmt.Println()
	fmt.Println("2. Basic Concepts:")
	fmt.Println("   go run cmd/examples/basics/variables_demo.go")
	fmt.Println("   go run cmd/examples/basics/functions_demo.go")
	fmt.Println("   go run cmd/examples/basics/pointers_demo.go")
	fmt.Println()
	fmt.Println("3. Data Structures:")
	fmt.Println("   go run cmd/examples/datastructures/stack_demo.go")
	fmt.Println("   go run cmd/examples/datastructures/queue_demo.go")
	fmt.Println()
	fmt.Println("4. Algorithms:")
	fmt.Println("   go run cmd/examples/algorithms/search_demo.go")
	fmt.Println("   go run cmd/examples/algorithms/sort_demo.go")
	fmt.Println()
	fmt.Println("5. Concurrency:")
	fmt.Println("   go run cmd/examples/concurrency/goroutines_demo.go")
	fmt.Println("   go run cmd/examples/concurrency/channels_demo.go")
	fmt.Println()
	fmt.Println("6. Cache:")
	fmt.Println("   go run cmd/examples/cache/cache_demo.go")
	fmt.Println()
	fmt.Println("📖 For more information, see README.md")
	fmt.Println()
}

func runFizzBuzzDemo() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"3", "3"},
		{"5", "5"},
		{"15", "15"},
		{"9", "3"},
		{"25", "5"},
		{"45", "15"},
		{"7", "7"},
		{"hello", "hello"},
		{"0", "0"},
		{"-3", "-3"},
		{"abc123", "abc123"},
		{"  15  ", "15"},
	}

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    FizzBuzz String Tests                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	for _, tc := range testCases {
		result := fizzbuzz(tc.input)
		status := "✓"
		if result != tc.expected {
			status = "✗"
			fmt.Printf("%s Input: %q → Output: %q (expected: %q)\n", status, tc.input, result, tc.expected)
		} else {
			fmt.Printf("%s Input: %q → Output: %q\n", status, tc.input, result)
		}
	}
	fmt.Println()
}
