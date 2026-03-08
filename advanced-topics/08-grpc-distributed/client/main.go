package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

const (
	serverAddr = "localhost:50051"
)

func main() {
	color.Cyan("🔌 Connecting to gRPC server at %s", serverAddr)

	// Create connection with TLS if available
	var opts []grpc.DialOption

	// Try TLS first
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // Only for development!
	})
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		color.Yellow("⚠️  TLS connection failed, trying insecure connection")
		conn, err = grpc.Dial(serverAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("❌ Failed to connect: %v", err)
		}
	}
	defer conn.Close()

	color.Green("✓ Connected to server")

	client := pb.NewUserServiceClient(conn)

	// Run examples
	runUnaryExample(client)
	runServerStreamingExample(client)
	runClientStreamingExample(client)
	runBidirectionalStreamingExample(client)
	runDeadlineExample(client)
	runErrorHandlingExample(client)

	color.Green("\n✓ All examples completed successfully!")
}

// Unary RPC Example
func runUnaryExample(client pb.UserServiceClient) {
	color.Blue("\n=== 📞 Unary RPC Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetUserRequest{Id: "1"}
	color.Yellow("Sending request for user ID: %s", req.Id)

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			color.Red("❌ gRPC error: %s (code: %s)", st.Message(), st.Code())
		} else {
			color.Red("❌ Error: %v", err)
		}
		return
	}

	user := resp.User
	color.Green("✓ Received user:")
	fmt.Printf("  ID: %s\n", user.Id)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Age: %d\n", user.Age)
	fmt.Printf("  Tags: %v\n", user.Tags)
}

// Server Streaming Example
func runServerStreamingExample(client pb.UserServiceClient) {
	color.Blue("\n=== 📡 Server Streaming Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.ListUsersRequest{Limit: 3}
	color.Yellow("Requesting user list (limit: %d)", req.Limit)

	stream, err := client.ListUsers(ctx, req)
	if err != nil {
		color.Red("❌ Error calling ListUsers: %v", err)
		return
	}

	count := 0
	for {
		user, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			color.Red("❌ Error receiving: %v", err)
			return
		}

		count++
		color.Green("✓ Received user #%d: %s", count, user.Name)
		fmt.Printf("  ID: %s, Email: %s, Age: %d\n", user.Id, user.Email, user.Age)
	}

	color.Green("✓ Received %d users", count)
}

// Client Streaming Example
func runClientStreamingExample(client pb.UserServiceClient) {
	color.Blue("\n=== 📤 Client Streaming Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.CreateUsers(ctx)
	if err != nil {
		color.Red("❌ Error creating stream: %v", err)
		return
	}

	// Send multiple users
	users := []*pb.CreateUserRequest{
		{Name: "David Wilson", Email: "david@example.com", Age: 30, Tags: []string{"premium"}},
		{Name: "Emma Davis", Email: "emma@example.com", Age: 25, Tags: []string{"standard"}},
		{Name: "Frank Miller", Email: "frank@example.com", Age: 38, Tags: []string{"vip"}},
	}

	for i, userReq := range users {
		color.Yellow("Sending user %d: %s", i+1, userReq.Name)
		if err := stream.Send(userReq); err != nil {
			color.Red("❌ Error sending: %v", err)
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Close and receive response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		color.Red("❌ Error closing stream: %v", err)
		return
	}

	color.Green("✓ Server response: %s", resp.Message)
	color.Green("✓ Total users created: %d", resp.Count)

	for i, user := range resp.Users {
		fmt.Printf("  %d. %s (%s) - Age: %d\n", i+1, user.Name, user.Email, user.Age)
	}
}

// Bidirectional Streaming Example
func runBidirectionalStreamingExample(client pb.UserServiceClient) {
	color.Blue("\n=== 💬 Bidirectional Streaming Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		color.Red("❌ Error creating chat stream: %v", err)
		return
	}

	// Send messages in goroutine
	go func() {
		messages := []string{
			"Hello, server!",
			"How are you?",
			"This is gRPC streaming!",
		}

		for i, msg := range messages {
			req := &pb.ChatMessage{
				UserId:    "client-1",
				Message:   msg,
				Timestamp: time.Now().Unix(),
				Type:      pb.MessageType_MESSAGE,
			}
			color.Yellow("📤 Sending: %s", msg)
			if err := stream.Send(req); err != nil {
				color.Red("❌ Error sending: %v", err)
				return
			}
			time.Sleep(time.Second)
		}

		if err := stream.CloseSend(); err != nil {
			color.Red("❌ Error closing send: %v", err)
		}
	}()

	// Receive messages
	count := 0
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			color.Red("❌ Error receiving: %v", err)
			break
		}

		count++
		color.Green("✓ Received #%d from %s: %s", count, resp.UserId, resp.Message)
	}

	color.Green("✓ Received %d messages", count)
}

// Deadline Example
func runDeadlineExample(client pb.UserServiceClient) {
	color.Blue("\n=== ⏰ Deadline/Timeout Example ===")

	// Very short deadline
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req := &pb.GetUserRequest{Id: "1"}
	color.Yellow("Sending request with 50ms deadline (likely to timeout)")

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			color.Red("❌ Request timed out (deadline exceeded)")
		} else {
			color.Red("❌ Error: %v", err)
		}
		return
	}

	color.Green("✓ Received user: %s", resp.User.Name)
}

// Error Handling Example
func runErrorHandlingExample(client pb.UserServiceClient) {
	color.Blue("\n=== ⚠️  Error Handling Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get non-existent user
	req := &pb.GetUserRequest{Id: "999"}
	color.Yellow("Requesting non-existent user ID: %s", req.Id)

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			color.Red("❌ gRPC Error:")
			color.Red("   Code: %s", st.Code())
			color.Red("   Message: %s", st.Message())
			return
		}
		color.Red("❌ Error: %v", err)
		return
	}

	color.Green("✓ Received user: %s", resp.User.Name)
}
