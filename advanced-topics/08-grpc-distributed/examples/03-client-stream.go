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

// Client Streaming Example - Client sends multiple requests to one response
func main() {
	fmt.Println("📤 Client Streaming Example")
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client stream
	stream, err := client.CreateUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Send multiple requests
	users := []*pb.CreateUserRequest{
		{Name: "Alice Smith", Email: "alice@example.com", Age: 28, Tags: []string{"premium"}},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35, Tags: []string{"standard"}},
		{Name: "Carol Williams", Email: "carol@example.com", Age: 42, Tags: []string{"vip"}},
	}

	fmt.Println("\n📤 Sending users to server:")
	for i, userReq := range users {
		fmt.Printf("  %d. %s (%s) - Age: %d\n", i+1, userReq.Name, userReq.Email, userReq.Age)

		if err := stream.Send(userReq); err != nil {
			log.Fatalf("Failed to send user %d: %v", i+1, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Close stream and receive response
	fmt.Println("\n✓ All users sent. Closing stream...")
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close and receive: %v", err)
	}

	// Process response
	fmt.Println("\n✓ Server response:")
	fmt.Printf("  Message: %s\n", resp.Message)
	fmt.Printf("  Total created: %d\n", resp.Count)

	fmt.Println("\n  Created users:")
	for i, user := range resp.Users {
		fmt.Printf("    %d. %s (ID: %s, Email: %s)\n", i+1, user.Name, user.Id, user.Email)
	}

	fmt.Println("\n✅ Client streaming completed successfully!")
}
