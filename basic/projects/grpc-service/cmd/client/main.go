// Package main provides the gRPC client entry point with examples of all RPC types.
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/DimaJoyti/go-pro/basic/projects/grpc-service/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Connect to gRPC server with timeout.
	connCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	_ = connCtx // Suppress unused variable warning
	client := userpb.NewUserServiceClient(conn)

	log.Println("🔗 Connected to gRPC server at", address)
	log.Println("")

	// Run all examples.
	runUnaryExample(client)
	runServerStreamingExample(client)
	runClientStreamingExample(client)
	runBidirectionalStreamingExample(client)
}

// =============================================================================
// Unary RPC Example
// =============================================================================
func runUnaryExample(client userpb.UserServiceClient) {
	log.Println("📡 Example 1: Unary RPC (GetUser)")
	log.Println("   Sending: GetUserRequest{Id: \"1\"}")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
		return
	}

	log.Printf("   ✅ Response: User{Id: %q, Name: %q, Email: %q, Age: %d}",
		resp.User.Id, resp.User.Name, resp.User.Email, resp.User.Age)
	log.Println("")
}

// =============================================================================
// Server Streaming RPC Example
// =============================================================================
func runServerStreamingExample(client userpb.UserServiceClient) {
	log.Println("📡 Example 2: Server Streaming RPC (ListUsers)")
	log.Println("   Sending: ListUsersRequest{Limit: 5}")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.ListUsers(ctx, &userpb.ListUsersRequest{Limit: 5})
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
		return
	}

	count := 0
	for {
		user, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("   ❌ Stream error: %v", err)
			return
		}
		count++
		log.Printf("   ✅ User #%d: {Id: %q, Name: %q, Email: %q}",
			count, user.Id, user.Name, user.Email)
	}

	log.Printf("   📊 Total users received: %d", count)
	log.Println("")
}

// =============================================================================
// Client Streaming RPC Example
// =============================================================================
func runClientStreamingExample(client userpb.UserServiceClient) {
	log.Println("📡 Example 3: Client Streaming RPC (CreateUsers)")
	log.Println("   Sending 3 user creation requests...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateUsers(ctx)
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
		return
	}

	// Send 3 users.
	usersToCreate := []*userpb.CreateUserRequest{
		{Name: "Alice", Email: "alice@example.com", Age: 28, Tags: []string{"admin", "developer"}},
		{Name: "Bob", Email: "bob@example.com", Age: 34, Tags: []string{"developer"}},
		{Name: "Charlie", Email: "charlie@example.com", Age: 22, Tags: []string{"designer"}},
	}

	for _, req := range usersToCreate {
		log.Printf("   📤 Sending user: {Name: %q, Email: %q}",
			req.Name, req.Email)
		if err := stream.Send(req); err != nil {
			log.Printf("   ❌ Send error: %v", err)
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("   ❌ CloseAndRecv error: %v", err)
		return
	}

	log.Printf("   ✅ Created %d users:", resp.Count)
	for _, user := range resp.Users {
		log.Printf("      User: {Id: %q, Name: %q}",
			user.Id, user.Name)
	}
	log.Println("")
}

// =============================================================================
// Bidirectional Streaming RPC Example
// =============================================================================
func runBidirectionalStreamingExample(client userpb.UserServiceClient) {
	log.Println("📡 Example 4: Bidirectional Streaming RPC (Chat)")
	log.Println("   Starting chat session...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
		return
	}

	// Send and receive messages concurrently.
	done := make(chan struct{})

	// Receive messages.
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				select {
				case <-done:
					return
				default:
					log.Printf("   ❌ Receive error: %v", err)
					return
				}
			}
			log.Printf("   💬 [%s]: %s (at %s)",
				msg.UserId,
				msg.Message,
				time.Unix(msg.Timestamp, 0).Format("15:04:05"))
		}
	}()

	// Send messages.
	messages := []string{
		"Hello, everyone!",
		"How is the gRPC tutorial going?",
		"Streaming is awesome!",
		"Thanks for watching!",
	}

	for _, msg := range messages {
		log.Printf("   📤 Sending: %q", msg)
		timestamp := time.Now().Unix()
		if err := stream.Send(&userpb.ChatMessage{
			UserId:    "client-1",
			Message:   msg,
			Timestamp: timestamp,
		}); err != nil {
			log.Printf("   ❌ Send error: %v", err)
			return
		}
		time.Sleep(500 * time.Millisecond) // Brief pause between messages
	}

	// Close the stream.
	if err := stream.CloseSend(); err != nil {
		log.Printf("   ❌ CloseSend error: %v", err)
		return
	}

	close(done)
	time.Sleep(1 * time.Second) // Wait for final messages
	log.Println("")
}
