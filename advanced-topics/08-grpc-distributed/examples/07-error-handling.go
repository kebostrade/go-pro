package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
	"google.golang.org/grpc"
)

// Error Handling Example - Demonstrating proper gRPC error handling
func main() {
	fmt.Println("⚠️  Error Handling Example")
	fmt.Println("==========================")

	// Connect to server
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Example 1: Handle NotFound error
	fmt.Println("\n📝 Example 1: NotFound error")
	handleUserNotFound(client)

	// Example 2: Handle InvalidArgument error
	fmt.Println("\n📝 Example 2: InvalidArgument error")
	handleInvalidArgument(client)

	// Example 3: Handle DeadlineExceeded
	fmt.Println("\n📝 Example 3: Deadline exceeded")
	handleDeadlineExceeded(client)

	// Example 4: Generic error handling
	fmt.Println("\n📝 Example 4: Generic error handling")
	handleGenericError(client)

	fmt.Println("\n✅ Error handling example completed!")
}

func handleUserNotFound(client pb.UserServiceClient) {
	ctx := context.Background()
	req := &pb.GetUserRequest{Id: "999"}

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Printf("  ❌ Non-gRPC error: %v\n", err)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			fmt.Printf("  ✓ Handled NotFound: %s\n", st.Message())
		default:
			fmt.Printf("  ❌ Unexpected code: %s, message: %s\n", st.Code(), st.Message())
		}
		return
	}

	fmt.Printf("  ✓ Found user: %s\n", resp.User.Name)
}

func handleInvalidArgument(client pb.UserServiceClient) {
	fmt.Println("  Note: This would require server validation to demonstrate")
	fmt.Println("  In a real scenario, the server would return codes.InvalidArgument")
	fmt.Println("  when client sends invalid data (e.g., empty email)")
}

func handleDeadlineExceeded(client pb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 1)
	defer cancel()

	req := &pb.GetUserRequest{Id: "1"}
	resp, err := client.GetUser(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Printf("  ❌ Non-gRPC error: %v\n", err)
			return
		}

		switch st.Code() {
		case codes.DeadlineExceeded:
			fmt.Printf("  ✓ Handled DeadlineExceeded: Request took too long\n")
		default:
			fmt.Printf("  ❌ Unexpected code: %s\n", st.Code())
		}
		return
	}

	fmt.Printf("  ✓ Received: %s\n", resp.User.Name)
}

func handleGenericError(client pb.UserServiceClient) {
	fmt.Println("  Generic error handling pattern:")

	ctx := context.Background()
	req := &pb.GetUserRequest{Id: "999"}
	resp, err := client.GetUser(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Printf("  ❌ Non-gRPC error: %v\n", err)
			return
		}

		fmt.Printf("  Code: %s\n", st.Code())
		fmt.Printf("  Message: %s\n", st.Message())

		// Handle different error codes
		switch st.Code() {
		case codes.NotFound:
			fmt.Println("  Action: Show user not found message")
		case codes.InvalidArgument:
			fmt.Println("  Action: Show validation error to user")
		case codes.DeadlineExceeded:
			fmt.Println("  Action: Show timeout message, offer retry")
		case codes.Unauthenticated:
			fmt.Println("  Action: Redirect to login")
		case codes.PermissionDenied:
			fmt.Println("  Action: Show access denied message")
		default:
			fmt.Println("  Action: Show generic error message")
		}

		return
	}

	fmt.Printf("  ✓ Success: %s\n", resp.User.Name)
}
