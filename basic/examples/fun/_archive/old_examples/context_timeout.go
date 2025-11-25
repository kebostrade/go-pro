//go:build ignore

package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Task: Demonstrate the use of context package for managing timeouts,
// cancellation, and passing request-scoped values in Go.

// simulateWork simulates a long-running operation
func simulateWork(ctx context.Context, name string, duration time.Duration) error {
	select {
	case <-time.After(duration):
		fmt.Printf("✓ %s: Work completed successfully\n", name)
		return nil
	case <-ctx.Done():
		fmt.Printf("✗ %s: Work cancelled - %v\n", name, ctx.Err())
		return ctx.Err()
	}
}

// fetchData simulates fetching data with timeout
func fetchData(ctx context.Context, id int) (string, error) {
	// Simulate network delay
	delay := time.Duration(id*100) * time.Millisecond

	select {
	case <-time.After(delay):
		return fmt.Sprintf("Data for ID %d", id), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// processWithTimeout demonstrates context with timeout
func processWithTimeout() {
	fmt.Println("\n1. Context with Timeout")
	fmt.Println(strings.Repeat("-", 60))

	// Create context with 2 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fmt.Println("Starting work with 2 second timeout...")

	// Task that completes in time
	go simulateWork(ctx, "Task 1 (1s)", 1*time.Second)

	// Task that exceeds timeout
	go simulateWork(ctx, "Task 2 (3s)", 3*time.Second)

	time.Sleep(3 * time.Second)
}

// processWithDeadline demonstrates context with deadline
func processWithDeadline() {
	fmt.Println("\n2. Context with Deadline")
	fmt.Println(strings.Repeat("-", 60))

	// Create context with deadline 1.5 seconds from now
	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Deadline set for: %s\n", deadline.Format("15:04:05.000"))

	// Check remaining time
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("Time until deadline: %v\n", time.Until(deadline))
	}

	simulateWork(ctx, "Task with deadline", 2*time.Second)

	time.Sleep(2 * time.Second)
}

// processWithCancellation demonstrates manual cancellation
func processWithCancellation() {
	fmt.Println("\n3. Manual Cancellation")
	fmt.Println(strings.Repeat("-", 60))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling context...")
		cancel()
	}()

	simulateWork(ctx, "Cancellable task", 3*time.Second)

	time.Sleep(2 * time.Second)
}

// contextWithValues demonstrates passing values through context
func contextWithValues() {
	fmt.Println("\n4. Context with Values")
	fmt.Println(strings.Repeat("-", 60))

	type key string
	const (
		userIDKey    key = "userID"
		requestIDKey key = "requestID"
	)

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user123")
	ctx = context.WithValue(ctx, requestIDKey, "req-456")

	// Function that uses context values
	processRequest := func(ctx context.Context) {
		if userID, ok := ctx.Value(userIDKey).(string); ok {
			fmt.Printf("Processing request for user: %s\n", userID)
		}

		if requestID, ok := ctx.Value(requestIDKey).(string); ok {
			fmt.Printf("Request ID: %s\n", requestID)
		}
	}

	processRequest(ctx)
}

// parallelFetchWithTimeout demonstrates parallel operations with timeout
func parallelFetchWithTimeout() {
	fmt.Println("\n5. Parallel Operations with Timeout")
	fmt.Println(strings.Repeat("-", 60))

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ids := []int{1, 2, 3, 4, 5, 6, 7, 8}
	results := make(chan string, len(ids))
	errors := make(chan error, len(ids))

	// Launch parallel fetch operations
	for _, id := range ids {
		go func(id int) {
			data, err := fetchData(ctx, id)
			if err != nil {
				errors <- err
			} else {
				results <- data
			}
		}(id)
	}

	// Collect results
	successCount := 0
	errorCount := 0

	for i := 0; i < len(ids); i++ {
		select {
		case result := <-results:
			fmt.Printf("✓ Fetched: %s\n", result)
			successCount++
		case err := <-errors:
			fmt.Printf("✗ Error: %v\n", err)
			errorCount++
		}
	}

	fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, errorCount)
}

// chainedContexts demonstrates context chaining
func chainedContexts() {
	fmt.Println("\n6. Chained Contexts")
	fmt.Println(strings.Repeat("-", 60))

	// Parent context with 3 second timeout
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer parentCancel()

	// Child context with 1 second timeout (will timeout first)
	childCtx, childCancel := context.WithTimeout(parentCtx, 1*time.Second)
	defer childCancel()

	fmt.Println("Parent timeout: 3s, Child timeout: 1s")

	// This will be cancelled by child context timeout
	simulateWork(childCtx, "Child task (2s)", 2*time.Second)

	time.Sleep(2 * time.Second)
}

func main() {
	fmt.Println("Context Package Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Demo 1: Timeout
	processWithTimeout()

	// Demo 2: Deadline
	processWithDeadline()

	// Demo 3: Cancellation
	processWithCancellation()

	// Demo 4: Values
	contextWithValues()

	// Demo 5: Parallel with timeout
	parallelFetchWithTimeout()

	// Demo 6: Chained contexts
	chainedContexts()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Context demo completed!")
}
