package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
"

// Unary RPC Example - Simple request-response pattern
func main() {
	fmt.Println("📞 Unary RPC Example")
	fmt.Println("==================")

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call unary RPC
	fmt.Println("\n📤 Sending request for user ID: 1")
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	// Process response
	user := resp.User
	fmt.Println("\n✓ Received response:")
	fmt.Printf("  ID: %s\n", user.Id)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Age: %d\n", user.Age)
	fmt.Printf("  Tags: %v\n", user.Tags)

	fmt.Println("\n✅ Unary RPC completed successfully!")
}
