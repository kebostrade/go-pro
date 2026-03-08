package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

// Deadlines and Timeouts Example - Demonstrating proper deadline handling
func main() {
	fmt.Println("⏰ Deadlines and Timeouts Example")
	fmt.Println("==================================")

	// Connect to server
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Example 1: Request with sufficient deadline
	fmt.Println("\n📝 Example 1: Sufficient deadline (5 seconds)")
	callWithDeadline(client, "1", 5*time.Second)

	// Example 2: Request with insufficient deadline (will timeout)
	fmt.Println("\n📝 Example 2: Insufficient deadline (50ms)")
	callWithDeadline(client, "1", 50*time.Millisecond)

	// Example 3: Cancel request mid-flight
	fmt.Println("\n📝 Example 3: Cancellation")
	callWithCancellation(client)

	// Example 4: Server deadline checking
	fmt.Println("\n📝 Example 4: Multiple requests with varying deadlines")
	callWithMultipleDeadlines(client)

	fmt.Println("\n✅ Deadlines example completed!")
}

func callWithDeadline(client pb.UserServiceClient, userID string, deadline time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	fmt.Printf("  📤 Calling GetUser with %v deadline\n", deadline)
	start := time.Now()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
	elapsed := time.Since(start)

	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("  ❌ Request timed out after %v (deadline exceeded)\n", elapsed)
		} else {
			fmt.Printf("  ❌ Error: %v (took %v)\n", err, elapsed)
		}
		return
	}

	fmt.Printf("  ✓ Received: %s (took %v)\n", resp.User.Name, elapsed)
}

func callWithCancellation(client pb.UserServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  ⚠️  Cancelling request...")
		cancel()
	}()

	fmt.Println("  📤 Calling GetUser (will be cancelled)")
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})

	if err != nil {
		if err == context.Canceled {
			fmt.Println("  ❌ Request was cancelled")
		} else {
			fmt.Printf("  ❌ Error: %v\n", err)
		}
		return
	}

	fmt.Printf("  ✓ Received: %s\n", resp.User.Name)
}

func callWithMultipleDeadlines(client pb.UserServiceClient) {
	deadlines := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		200 * time.Millisecond,
		1 * time.Second,
	}

	for i, deadline := range deadlines {
		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		fmt.Printf("  📤 Request %d: %v deadline\n", i+1, deadline)

		resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
		cancel()

		if err != nil {
			if err == context.DeadlineExceeded {
				fmt.Printf("    ❌ Timed out\n")
			} else {
				fmt.Printf("    ❌ Error: %v\n", err)
			}
			continue
		}

		fmt.Printf("    ✓ Success: %s\n", resp.User.Name)
	}
}
