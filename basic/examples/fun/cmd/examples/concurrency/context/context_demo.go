package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Context Patterns Demo")

	demo1Timeout()
	demo2Cancellation()
	demo3TaskGroup()
	demo4Retry()
	demo5ContextValues()
}

func demo1Timeout() {
	utils.PrintSubHeader("1. Context with Timeout")

	// Fast operation (should succeed)
	fmt.Println("Running fast operation with 2s timeout...")
	err := concurrency.WithTimeout(2*time.Second, func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("  ✓ Fast operation completed")
		return nil
	})

	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	// Slow operation (should timeout)
	fmt.Println("\nRunning slow operation with 1s timeout...")
	err = concurrency.WithTimeout(1*time.Second, func(ctx context.Context) error {
		time.Sleep(2 * time.Second)
		fmt.Println("  This won't print")
		return nil
	})

	if err != nil {
		fmt.Printf("  ✗ Error: %v\n", err)
	}
}

func demo2Cancellation() {
	utils.PrintSubHeader("2. Cancellable Task")

	// Create a long-running task
	task := concurrency.NewCancellableTask(func(ctx context.Context) error {
		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Printf("  Working... step %d/10\n", i)
				time.Sleep(200 * time.Millisecond)
			}
		}
		return nil
	})

	// Cancel after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("\n  Cancelling task...")
		task.Cancel()
	}()

	// Wait for task
	err := task.Wait()
	if err != nil {
		fmt.Printf("  Task cancelled: %v\n", err)
	}
}

func demo3TaskGroup() {
	utils.PrintSubHeader("3. Task Group (Fail Fast)")

	ctx := context.Background()
	group := concurrency.NewTaskGroup(ctx)

	// Add successful tasks
	group.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  Task 1: Success")
		return nil
	})

	group.Go(func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("  Task 2: Success")
		return nil
	})

	// Add failing task
	group.Go(func(ctx context.Context) error {
		time.Sleep(150 * time.Millisecond)
		fmt.Println("  Task 3: Failed!")
		return errors.New("task 3 failed")
	})

	// This task should be cancelled
	group.Go(func(ctx context.Context) error {
		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("  Task 4: Cancelled due to other task failure")
				return ctx.Err()
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
		return nil
	})

	// Wait for all tasks
	err := group.Wait()
	if err != nil {
		fmt.Printf("\n  Group failed: %v\n", err)
	}

	// Show all errors
	allErrors := group.Errors()
	fmt.Printf("  Total errors: %d\n", len(allErrors))
}

func demo4Retry() {
	utils.PrintSubHeader("4. Retry with Exponential Backoff")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	attempt := 0

	err := concurrency.RetryWithContext(
		ctx,
		5,                    // max attempts
		100*time.Millisecond, // initial delay
		func() error {
			attempt++
			fmt.Printf("  Attempt %d...\n", attempt)

			// Fail first 3 attempts
			if attempt < 4 {
				return fmt.Errorf("temporary error")
			}

			fmt.Println("  ✓ Success!")
			return nil
		},
	)

	if err != nil {
		fmt.Printf("  Failed after retries: %v\n", err)
	}
}

func demo5ContextValues() {
	utils.PrintSubHeader("5. Typed Context Values")

	// Create typed context values
	userIDKey := concurrency.NewContextValue[int]("userID")
	usernameKey := concurrency.NewContextValue[string]("username")

	// Create context with values
	ctx := context.Background()
	ctx = userIDKey.WithValue(ctx, 12345)
	ctx = usernameKey.WithValue(ctx, "alice")

	// Retrieve values
	if userID, ok := userIDKey.Value(ctx); ok {
		fmt.Printf("User ID: %d\n", userID)
	}

	if username, ok := usernameKey.Value(ctx); ok {
		fmt.Printf("Username: %s\n", username)
	}

	// Pass context to function
	processRequest(ctx, userIDKey, usernameKey)
}

func processRequest(ctx context.Context, userIDKey *concurrency.ContextValue[int], usernameKey *concurrency.ContextValue[string]) {
	userID := userIDKey.MustValue(ctx)
	username := usernameKey.MustValue(ctx)

	fmt.Printf("\nProcessing request for user %s (ID: %d)\n", username, userID)
}
