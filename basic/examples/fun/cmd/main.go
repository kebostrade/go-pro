package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Go Programming Examples - Interactive Menu")

	for {
		displayMenu()
		choice := getUserChoice()

		if choice == 0 {
			fmt.Println("\nThank you for using Go Programming Examples! Goodbye! 👋")
			break
		}

		executeChoice(choice)

		fmt.Println("\nPress Enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func displayMenu() {
	utils.PrintSubHeader("Main Menu")

	fmt.Println("1.  Data Structures")
	fmt.Println("2.  Algorithms")
	fmt.Println("3.  Concurrency Patterns")
	fmt.Println("4.  Cache Examples")
	fmt.Println("5.  Basic Go Concepts")
	fmt.Println("6.  Advanced Examples")
	fmt.Println("0.  Exit")
	fmt.Println()
	fmt.Print("Enter your choice: ")
}

func getUserChoice() int {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}

	return choice
}

func executeChoice(choice int) {
	switch choice {
	case 1:
		showDataStructuresMenu()
	case 2:
		showAlgorithmsMenu()
	case 3:
		showConcurrencyMenu()
	case 4:
		showCacheMenu()
	case 5:
		showBasicsMenu()
	case 6:
		showAdvancedMenu()
	default:
		fmt.Println("\n❌ Invalid choice. Please try again.")
	}
}

func showDataStructuresMenu() {
	utils.PrintSubHeader("Data Structures")

	fmt.Println("Available examples:")
	fmt.Println("  • Stack (LIFO)")
	fmt.Println("  • Queue (FIFO)")
	fmt.Println("  • Linked List")
	fmt.Println("  • Priority Queue")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/datastructures/stack_demo.go")
	fmt.Println("  go run cmd/examples/datastructures/queue_demo.go")
	fmt.Println("  go run cmd/examples/datastructures/linkedlist_demo.go")
}

func showAlgorithmsMenu() {
	utils.PrintSubHeader("Algorithms")

	fmt.Println("Available examples:")
	fmt.Println("  • Binary Search")
	fmt.Println("  • Merge Sort (Sequential & Concurrent)")
	fmt.Println("  • Prime Numbers")
	fmt.Println("  • Palindrome Detection")
	fmt.Println("  • Fibonacci Sequence")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/algorithms/search_demo.go")
	fmt.Println("  go run cmd/examples/algorithms/sort_demo.go")
	fmt.Println("  go run cmd/examples/algorithms/primes_demo.go")
}

func showConcurrencyMenu() {
	utils.PrintSubHeader("Concurrency Patterns")

	fmt.Println("Available examples:")
	fmt.Println("  • Goroutines & WaitGroups")
	fmt.Println("  • Fan-Out/Fan-In Pattern")
	fmt.Println("  • Thread-Safe Data Structures")
	fmt.Println("  • Producer-Consumer Pattern")
	fmt.Println("  • Worker Pool")
	fmt.Println("  • Pipeline Pattern")
	fmt.Println("  • Rate Limiters (Token Bucket, Sliding Window, Leaky Bucket)")
	fmt.Println("  • Context (Timeout, Cancellation, Task Groups)")
	fmt.Println("  • Parallel Map/Filter/Reduce")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/concurrency/goroutines_demo.go")
	fmt.Println("  go run cmd/examples/concurrency/producer_consumer_demo.go")
	fmt.Println("  go run cmd/examples/concurrency/ratelimiter_demo.go")
	fmt.Println("  go run cmd/examples/concurrency/context_demo.go")
	fmt.Println("  go run cmd/examples/concurrency/parallel_demo.go")
}

func showCacheMenu() {
	utils.PrintSubHeader("Cache Examples")

	fmt.Println("Available examples:")
	fmt.Println("  • Generic Cache with TTL and Statistics")
	fmt.Println("  • Loading Cache (Auto-Load on Miss)")
	fmt.Println("  • GetOrCompute Pattern")
	fmt.Println("  • LRU Cache (Least Recently Used)")
	fmt.Println("  • LFU Cache (Least Frequently Used)")
	fmt.Println("  • Cache Eviction Policies")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/cache/cache_demo.go")
	fmt.Println("  go run cmd/examples/cache/lru_demo.go")
}

func showBasicsMenu() {
	utils.PrintSubHeader("Basic Go Concepts")

	fmt.Println("Available examples:")
	fmt.Println("  • Variables & Types (int, float, string, bool, type inference)")
	fmt.Println("  • Functions (basic, multiple returns, variadic, closures)")
	fmt.Println("  • Pointers (addresses, dereferencing, pass by reference)")
	fmt.Println("  • Structs & Methods (value/pointer receivers, embedding)")
	fmt.Println("  • Interfaces (polymorphism, type assertions, type switches)")
	fmt.Println("  • Loops & Control Flow (for, range, break, continue)")
	fmt.Println("  • Iota & Constants (enums, bit flags, custom expressions)")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/basics/variables_demo.go")
	fmt.Println("  go run cmd/examples/basics/functions_demo.go")
	fmt.Println("  go run cmd/examples/basics/pointers_demo.go")
	fmt.Println("  go run cmd/examples/basics/structs_demo.go")
	fmt.Println("  go run cmd/examples/basics/interfaces_demo.go")
	fmt.Println("  go run cmd/examples/basics/loops_demo.go")
	fmt.Println("  go run cmd/examples/basics/iota_demo.go")
}

func showAdvancedMenu() {
	utils.PrintSubHeader("Advanced Examples")

	fmt.Println("Available examples:")
	fmt.Println("  • JSON Parsing & Serialization")
	fmt.Println("  • HTTP Client (Weather API)")
	fmt.Println("  • Word Counter")
	fmt.Println("  • Order Management System")
	fmt.Println("  • Error Handling Patterns")
	fmt.Println()
	fmt.Println("To run these examples, use:")
	fmt.Println("  go run cmd/examples/advanced/json_demo.go")
	fmt.Println("  go run cmd/examples/advanced/weather_demo.go")
	fmt.Println("  go run cmd/examples/advanced/wordcount_demo.go")
}
