//go:build ignore

package main

import (
	"fmt"
	"sync"
)

// Task: Write a function that finds all prime numbers up to a given limit
// using concurrent goroutines for better performance.
// Divide the range into chunks and process each chunk in a separate goroutine.

// Example input: 30
// Expected output: [2 3 5 7 11 13 17 19 23 29]

// Function signature to implement:
// func findPrimes(limit int, workers int) []int

func findPrimes(limit int, workers int) []int {
	if limit < 2 {
		return []int{}
	}

	// Channel to collect prime numbers
	primesChan := make(chan int, limit)
	var wg sync.WaitGroup

	// Calculate chunk size for each worker
	chunkSize := limit / workers
	if chunkSize < 1 {
		chunkSize = 1
	}

	// Launch workers to check ranges
	for i := 0; i < workers; i++ {
		wg.Add(1)
		start := i*chunkSize + 2
		end := start + chunkSize
		if i == workers-1 {
			end = limit + 1 // Last worker handles remaining numbers
		}

		go func(s, e int) {
			defer wg.Done()
			for n := s; n < e; n++ {
				if isPrime(n) {
					primesChan <- n
				}
			}
		}(start, end)
	}

	// Close channel when all workers are done
	go func() {
		wg.Wait()
		close(primesChan)
	}()

	// Collect results
	var primes []int
	for prime := range primesChan {
		primes = append(primes, prime)
	}

	// Sort the results (since goroutines may finish in any order)
	sortInts(primes)
	return primes
}

// isPrime checks if a number is prime
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Check odd divisors up to sqrt(n)
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// sortInts performs a simple insertion sort on integers
func sortInts(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func main() {
	limit := 100
	workers := 4

	fmt.Printf("Finding prime numbers up to %d using %d workers...\n", limit, workers)
	primes := findPrimes(limit, workers)

	fmt.Printf("Found %d prime numbers:\n", len(primes))
	fmt.Println(primes)

	// Test with different limits
	fmt.Println("\nPrimes up to 30:")
	fmt.Println(findPrimes(30, 2))
}
