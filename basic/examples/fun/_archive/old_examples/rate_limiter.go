//go:build ignore

package main

import (
	"fmt"
	"time"
)

// Task: Implement a rate limiter using Go channels and goroutines.
// The rate limiter should allow a maximum number of requests per time window.

// RateLimiter controls the rate of operations
type RateLimiter struct {
	tokens   chan struct{}
	maxRate  int
	interval time.Duration
}

// NewRateLimiter creates a new rate limiter
// maxRate: maximum number of operations allowed per interval
// interval: time window for the rate limit
func NewRateLimiter(maxRate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(chan struct{}, maxRate),
		maxRate:  maxRate,
		interval: interval,
	}

	// Fill the token bucket initially
	for i := 0; i < maxRate; i++ {
		rl.tokens <- struct{}{}
	}

	// Start the token refill goroutine
	go rl.refillTokens()

	return rl
}

// refillTokens periodically adds tokens to the bucket
func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.interval / time.Duration(rl.maxRate))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.tokens <- struct{}{}:
			// Token added successfully
		default:
			// Bucket is full, skip
		}
	}
}

// Allow checks if an operation is allowed (non-blocking)
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until an operation is allowed
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// SimpleRateLimiter demonstrates a simpler rate limiting approach
func SimpleRateLimiter(requests int, duration time.Duration) {
	limiter := time.Tick(duration / time.Duration(requests))

	for i := 1; i <= requests*2; i++ {
		<-limiter
		fmt.Printf("Request %d processed at %s\n", i, time.Now().Format("15:04:05.000"))
	}
}

// processRequest simulates processing a request
func processRequest(id int, rl *RateLimiter) {
	if rl.Allow() {
		fmt.Printf("✓ Request %d: Allowed at %s\n", id, time.Now().Format("15:04:05.000"))
	} else {
		fmt.Printf("✗ Request %d: Rate limited at %s\n", id, time.Now().Format("15:04:05.000"))
	}
}

// processRequestWithWait simulates processing with waiting
func processRequestWithWait(id int, rl *RateLimiter) {
	rl.Wait()
	fmt.Printf("✓ Request %d: Processed at %s\n", id, time.Now().Format("15:04:05.000"))
}

func main() {
	fmt.Println("Rate Limiter Demo")
	fmt.Println("=".repeat(60))

	// Demo 1: Non-blocking rate limiter
	fmt.Println("\n1. Non-blocking Rate Limiter (5 requests per second)")
	fmt.Println("-".repeat(60))

	rl := NewRateLimiter(5, time.Second)

	// Send 10 requests rapidly
	for i := 1; i <= 10; i++ {
		processRequest(i, rl)
		time.Sleep(50 * time.Millisecond)
	}

	// Wait a bit and try again
	time.Sleep(1 * time.Second)
	fmt.Println("\nAfter 1 second delay:")
	for i := 11; i <= 15; i++ {
		processRequest(i, rl)
		time.Sleep(50 * time.Millisecond)
	}

	// Demo 2: Blocking rate limiter
	fmt.Println("\n\n2. Blocking Rate Limiter (3 requests per second)")
	fmt.Println("-".repeat(60))

	rl2 := NewRateLimiter(3, time.Second)

	// Process 6 requests with waiting
	for i := 1; i <= 6; i++ {
		go processRequestWithWait(i, rl2)
	}

	time.Sleep(3 * time.Second)

	// Demo 3: Simple ticker-based rate limiter
	fmt.Println("\n\n3. Simple Ticker-based Rate Limiter (2 requests per second)")
	fmt.Println("-".repeat(60))

	// This will process 4 requests at a rate of 2 per second
	go SimpleRateLimiter(2, time.Second)

	time.Sleep(3 * time.Second)
	fmt.Println("\nRate limiter demo completed!")
}

// Helper function to repeat strings
func (s string) repeat(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
