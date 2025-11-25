//go:build ignore

package main

import (
	"fmt"
)

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// Start 4 worker goroutines
	for i := 0; i < 4; i++ {
		go worker(jobs, results)
	}

	// Send jobs
	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)

	// Collect results
	for j := 0; j < 100; j++ {
		fmt.Println(<-results)
	}
}

func worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- fib(n)
	}
}

func fib(n int) int {
	if n < 0 {
		return 0 // Return 0 for negative numbers
	}
	if n <= 1 {
		return n // Base cases: fib(0) = 0, fib(1) = 1
	}

	// Iterative implementation for better performance
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
