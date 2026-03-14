package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

// Metadata Example - Demonstrating how to send and receive metadata
//
// Metadata is used for:
// - Authentication tokens
// - Request tracing
// - Rate limiting information
// - Custom headers
func main() {
	fmt.Println("🏷️  Metadata Example")
	fmt.Println("====================")

	// Connect to server
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Example 1: Send authentication token
	fmt.Println("\n📝 Example 1: Send authentication token")
	sendAuthToken(client)

	// Example 2: Send tracing metadata
	fmt.Println("\n📝 Example 2: Send tracing metadata")
	sendTracingMetadata(client)

	// Example 3: Send custom metadata
	fmt.Println("\n📝 Example 3: Send custom metadata")
	sendCustomMetadata(client)

	// Example 4: Send multiple metadata pairs
	fmt.Println("\n📝 Example 4: Send multiple metadata pairs")
	sendMultipleMetadata(client)

	fmt.Println("\n✅ Metadata example completed!")
}

func sendAuthToken(client pb.UserServiceClient) {
	// Create metadata with auth token
	md := metadata.Pairs(
		"authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user-id", "12345",
	)

	// Create context with metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make RPC call
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s (with auth metadata)\n", resp.User.Name)
}

func sendTracingMetadata(client pb.UserServiceClient) {
	// Create tracing metadata
	md := metadata.Pairs(
		"trace-id", "abc-123-def-456",
		"span-id", "span-789",
		"request-id", "req-001",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s (with tracing metadata)\n", resp.User.Name)
}

func sendCustomMetadata(client pb.UserServiceClient) {
	// Create custom metadata
	md := metadata.Pairs(
		"client-version", "1.0.0",
		"device-id", "device-abc-123",
		"locale", "en-US",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s (with custom metadata)\n", resp.User.Name)
}

func sendMultipleMetadata(client pb.UserServiceClient) {
	// Create metadata with multiple values for the same key
	md := metadata.Pairs(
		"authorization", "Bearer token1",
		"authorization", "Bearer token2", // Multiple values
		"custom-header", "value1",
		"custom-header", "value2",
		"custom-header", "value3",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s (with multiple metadata values)\n", resp.User.Name)
}

// Example of merging metadata from multiple sources
func mergeMetadataExample() {
	fmt.Println("\n📝 Example: Merging metadata")

	// Existing metadata from incoming context
	existingMD := metadata.Pairs("trace-id", "existing-trace")

	// New metadata to add
	newMD := metadata.Pairs("user-id", "12345")

	// Merge metadata
	mergedMD := metadata.Join(existingMD, newMD)

	ctx := metadata.NewOutgoingContext(context.Background(), mergedMD)

	// Use merged context for RPC calls
	_ = ctx // In real usage, pass this to gRPC calls
	fmt.Println("  ✓ Metadata merged successfully")
}
