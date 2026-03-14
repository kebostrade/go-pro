package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

// Server Streaming Example - Server sends multiple responses to one request
func main() {
	fmt.Println("📡 Server Streaming Example")
	fmt.Println("===========================")

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

	// Call server streaming RPC
	fmt.Println("\n📤 Requesting list of users (limit: 3)")
	stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{Limit: 3})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	// Receive stream of responses
	fmt.Println("\n✓ Receiving user stream:")
	count := 0
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			// Server finished sending
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}

		count++
		fmt.Printf("\n  User #%d:\n", count)
		fmt.Printf("    ID: %s\n", user.Id)
		fmt.Printf("    Name: %s\n", user.Name)
		fmt.Printf("    Email: %s\n", user.Email)
		fmt.Printf("    Age: %d\n", user.Age)
		fmt.Printf("    Tags: %v\n", user.Tags)
	}

	fmt.Printf("\n✅ Server streaming completed! Received %d users\n", count)
}
