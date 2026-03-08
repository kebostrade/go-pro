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

// Bidirectional Streaming Example - Both client and server send messages asynchronously
func main() {
	fmt.Println("💬 Bidirectional Streaming Example")
	fmt.Println("===================================")

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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create bidirectional stream
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Start goroutine to receive messages
	receivedChan := make(chan *pb.ChatMessage, 10)
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				close(receivedChan)
				return
			}
			if err != nil {
				log.Printf("Receive error: %v", err)
				close(receivedChan)
				return
			}
			receivedChan <- msg
		}
	}()

	// Send messages
	messages := []string{
		"Hello, server!",
		"How are you doing?",
		"This is bidirectional streaming!",
		"Pretty cool, right?",
		"Okay, goodbye!",
	}

	fmt.Println("\n📤 Sending messages:")
	for i, msg := range messages {
		req := &pb.ChatMessage{
			UserId:    "client-1",
			Message:   msg,
			Timestamp: time.Now().Unix(),
			Type:      pb.MessageType_MESSAGE,
		}
		fmt.Printf("  %d. %s\n", i+1, msg)

		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send message %d: %v", i+1, err)
		}

		// Wait a bit between messages
		time.Sleep(500 * time.Millisecond)
	}

	// Close sending
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close send: %v", err)
	}
	fmt.Println("\n✓ Finished sending messages")

	// Receive all responses
	fmt.Println("\n✓ Receiving responses:")
	count := 0
	for msg := range receivedChan {
		count++
		fmt.Printf("  %d. From %s: %s\n", count, msg.UserId, msg.Message)
	}

	fmt.Printf("\n✅ Bidirectional streaming completed! Received %d responses\n", count)
}
