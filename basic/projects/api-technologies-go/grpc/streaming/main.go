package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/DimaJoyti/go-pro/api-technologies-go/grpc/proto"
)

// UserStore - In-memory storage
type UserStore struct {
	mu     sync.RWMutex
	users  map[int32]*pb.User
	nextID int32
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int32]*pb.User),
		nextID: 1,
	}
}

func (s *UserStore) Create(username, email, role string) *pb.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := timestamppb.Now()
	user := &pb.User{
		Id:        s.nextID,
		Username:  username,
		Email:     email,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserStore) GetAll(role string) []*pb.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*pb.User, 0)
	for _, user := range s.users {
		if role != "" && user.Role != role {
			continue
		}
		users = append(users, user)
	}
	return users
}

// UserServiceServer implements streaming gRPC methods
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	store      *UserStore
	chatRooms  map[int32]chan *pb.ChatMessage
	chatMutex  sync.RWMutex
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		store:     NewUserStore(),
		chatRooms: make(map[int32]chan *pb.ChatMessage),
	}
}

// Unary methods (minimal implementation)
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	return s.store.Create(req.Username, req.Email, req.Role), nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "see unary example")
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "see unary example")
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "see unary example")
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "see unary example")
}

// SERVER STREAMING: Stream users one by one
func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
	log.Printf("StreamUsers called: role=%s", req.Role)

	users := s.store.GetAll(req.Role)

	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return status.Errorf(codes.Internal, "failed to send user: %v", err)
		}
		// Simulate delay for demonstration
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("StreamUsers completed: sent %d users", len(users))
	return nil
}

// CLIENT STREAMING: Receive multiple users and create them
func (s *UserServiceServer) CreateUsers(stream pb.UserService_CreateUsersServer) error {
	log.Println("CreateUsers called (client streaming)")

	var users []*pb.User
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Client finished sending
			log.Printf("CreateUsers completed: created %d users", count)
			return stream.SendAndClose(&pb.CreateUsersResponse{
				Users: users,
				Count: int32(count),
			})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive: %v", err)
		}

		user := s.store.Create(req.Username, req.Email, req.Role)
		users = append(users, user)
		count++
		log.Printf("Created user %d: %s", count, user.Username)
	}
}

// BIDIRECTIONAL STREAMING: Chat room
func (s *UserServiceServer) Chat(stream pb.UserService_ChatServer) error {
	log.Println("Chat called (bidirectional streaming)")

	// Create a channel for this chat session
	chatChan := make(chan *pb.ChatMessage, 10)
	defer close(chatChan)

	// Goroutine to receive messages from client
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}

			log.Printf("Received chat message from user %d: %s", msg.UserId, msg.Message)

			// Broadcast to all connected clients (simplified - just echo back)
			msg.Timestamp = timestamppb.Now()
			chatChan <- msg
		}
	}()

	// Send messages to client
	for msg := range chatChan {
		if err := stream.Send(msg); err != nil {
			return status.Errorf(codes.Internal, "failed to send message: %v", err)
		}
	}

	return nil
}

// Start gRPC server
func startServer() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userService := NewUserServiceServer()

	// Seed data
	userService.store.Create("alice", "alice@example.com", "admin")
	userService.store.Create("bob", "bob@example.com", "user")
	userService.store.Create("charlie", "charlie@example.com", "guest")

	pb.RegisterUserServiceServer(grpcServer, userService)

	fmt.Println("🚀 gRPC Streaming server starting on :9090")
	fmt.Println("📚 Available streaming methods:")
	fmt.Println("  - StreamUsers (server streaming)")
	fmt.Println("  - CreateUsers (client streaming)")
	fmt.Println("  - Chat (bidirectional streaming)")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Example client
func runClient() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	fmt.Println("\n🔧 Testing gRPC streaming client...")

	// 1. Server Streaming
	fmt.Println("\n1️⃣ Server Streaming: StreamUsers")
	streamCtx, streamCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer streamCancel()

	stream, err := client.StreamUsers(streamCtx, &pb.StreamUsersRequest{Role: ""})
	if err != nil {
		log.Printf("StreamUsers failed: %v", err)
	} else {
		fmt.Println("📥 Receiving users from server...")
		for {
			user, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving: %v", err)
				break
			}
			fmt.Printf("   ✅ Received: %s (%s) - %s\n", user.Username, user.Email, user.Role)
		}
	}

	// 2. Client Streaming
	fmt.Println("\n2️⃣ Client Streaming: CreateUsers")
	createCtx, createCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer createCancel()

	createStream, err := client.CreateUsers(createCtx)
	if err != nil {
		log.Printf("CreateUsers failed: %v", err)
	} else {
		newUsers := []*pb.CreateUserRequest{
			{Username: "david", Email: "david@example.com", Role: "user"},
			{Username: "eve", Email: "eve@example.com", Role: "admin"},
			{Username: "frank", Email: "frank@example.com", Role: "guest"},
		}

		fmt.Println("📤 Sending users to server...")
		for _, req := range newUsers {
			if err := createStream.Send(req); err != nil {
				log.Printf("Error sending: %v", err)
				break
			}
			fmt.Printf("   ✅ Sent: %s\n", req.Username)
			time.Sleep(300 * time.Millisecond)
		}

		resp, err := createStream.CloseAndRecv()
		if err != nil {
			log.Printf("Error closing stream: %v", err)
		} else {
			fmt.Printf("✅ Server created %d users\n", resp.Count)
		}
	}

	// 3. Bidirectional Streaming
	fmt.Println("\n3️⃣ Bidirectional Streaming: Chat")
	chatCtx, chatCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer chatCancel()

	chatStream, err := client.Chat(chatCtx)
	if err != nil {
		log.Printf("Chat failed: %v", err)
	} else {
		// Goroutine to receive messages
		go func() {
			for {
				msg, err := chatStream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("Error receiving chat: %v", err)
					return
				}
				fmt.Printf("   📨 Received: User %d: %s\n", msg.UserId, msg.Message)
			}
		}()

		// Send messages
		messages := []string{"Hello!", "How are you?", "Goodbye!"}
		for i, text := range messages {
			msg := &pb.ChatMessage{
				UserId:  1,
				Message: text,
			}
			if err := chatStream.Send(msg); err != nil {
				log.Printf("Error sending chat: %v", err)
				break
			}
			fmt.Printf("   📤 Sent: %s\n", text)
			time.Sleep(1 * time.Second)

			// Close after last message
			if i == len(messages)-1 {
				chatStream.CloseSend()
			}
		}

		time.Sleep(2 * time.Second) // Wait for responses
	}
}

func main() {
	// Start server in goroutine
	go startServer()

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Run client examples
	runClient()

	// Keep server running
	fmt.Println("\n✅ All streaming examples completed!")
	fmt.Println("Press Ctrl+C to exit...")
	select {}
}

