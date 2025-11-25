package main

import (
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Rate Limiter Demo")

	demo1TokenBucket()
	demo2SlidingWindow()
	demo3LeakyBucket()
	demo4AdaptiveRateLimiter()
}

func demo1TokenBucket() {
	utils.PrintSubHeader("1. Token Bucket Rate Limiter")

	// Create rate limiter: 5 requests per second
	limiter := concurrency.NewRateLimiter(5, time.Second)
	defer limiter.Stop()

	fmt.Println("Rate limit: 5 requests per second")
	fmt.Println("Sending 10 requests rapidly...")

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("✓ Request %d: Allowed at %s\n", i, time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("✗ Request %d: Rate limited at %s\n", i, time.Now().Format("15:04:05.000"))
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for tokens to refill
	fmt.Println("\nWaiting for tokens to refill...")
	time.Sleep(1 * time.Second)

	fmt.Println("\nSending 3 more requests:")
	for i := 11; i <= 13; i++ {
		limiter.Wait() // Block until allowed
		fmt.Printf("✓ Request %d: Processed at %s\n", i, time.Now().Format("15:04:05.000"))
	}
}

func demo2SlidingWindow() {
	utils.PrintSubHeader("2. Sliding Window Rate Limiter")

	// Create sliding window limiter: 3 requests per 2 seconds
	limiter := concurrency.NewSlidingWindowRateLimiter(3, 2*time.Second)

	fmt.Println("Rate limit: 3 requests per 2-second window")
	fmt.Println("Sending requests...")

	for i := 1; i <= 8; i++ {
		if limiter.Allow() {
			fmt.Printf("✓ Request %d: Allowed (count: %d/3) at %s\n",
				i, limiter.Count(), time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("✗ Request %d: Rate limited (count: %d/3) at %s\n",
				i, limiter.Count(), time.Now().Format("15:04:05.000"))
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func demo3LeakyBucket() {
	utils.PrintSubHeader("3. Leaky Bucket Rate Limiter")

	// Create leaky bucket: capacity 5, leak rate 1 per 200ms
	bucket := concurrency.NewLeakyBucket(5, 200*time.Millisecond)
	defer bucket.Stop()

	fmt.Println("Bucket capacity: 5")
	fmt.Println("Leak rate: 1 item per 200ms")
	fmt.Println("Adding items rapidly...")

	for i := 1; i <= 10; i++ {
		if bucket.Add() {
			fmt.Printf("✓ Item %d: Added to bucket at %s\n", i, time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("✗ Item %d: Bucket full at %s\n", i, time.Now().Format("15:04:05.000"))
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("\nWaiting for bucket to drain...")
	time.Sleep(1 * time.Second)

	fmt.Println("\nAdding 2 more items:")
	for i := 11; i <= 12; i++ {
		bucket.Wait() // Block until space available
		fmt.Printf("✓ Item %d: Added at %s\n", i, time.Now().Format("15:04:05.000"))
	}
}

func demo4AdaptiveRateLimiter() {
	utils.PrintSubHeader("4. Adaptive Rate Limiter")

	// Create adaptive limiter: min 2, max 10 requests per second
	limiter := concurrency.NewAdaptiveRateLimiter(2, 10, time.Second)
	defer limiter.Stop()

	fmt.Println("Adaptive rate limiter: 2-10 requests/second")
	fmt.Printf("Initial rate: %d requests/second\n\n", limiter.CurrentRate())

	// Simulate successful requests
	fmt.Println("Simulating successful requests...")
	for i := 0; i < 15; i++ {
		if limiter.Allow() {
			limiter.RecordSuccess()
			if i%5 == 0 {
				fmt.Printf("  Success %d - Current rate: %d req/s\n", i, limiter.CurrentRate())
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Printf("\nAfter successes, rate increased to: %d req/s\n", limiter.CurrentRate())

	// Simulate failures
	fmt.Println("\nSimulating failures...")
	for i := 0; i < 5; i++ {
		limiter.RecordFailure()
		fmt.Printf("  Failure %d - Current rate: %d req/s\n", i, limiter.CurrentRate())
	}

	fmt.Printf("\nAfter failures, rate decreased to: %d req/s\n", limiter.CurrentRate())
}
