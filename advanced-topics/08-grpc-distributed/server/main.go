package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

func main() {
	color.Cyan("🚀 Starting gRPC Server on port %s", port)

	// Initialize server with sample data
	s := &server{
		users: make(map[string]*pb.User),
	}
	s.initializeData()

	// Create listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	// Create gRPC server with options
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor),
		grpc.StreamInterceptor(streamLoggingInterceptor),
	}

	// Check if TLS is enabled
	if _, err := os.Stat("certs/server.crt"); err == nil {
		color.Yellow("🔒 TLS enabled")

		// Load TLS certificates
		creds, err := credentials.LoadTLSCredentials("certs/server.crt", "certs/server.key")
		if err != nil {
			log.Fatalf("❌ Failed to load TLS credentials: %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
	} else {
		color.Yellow("⚠️  Running without TLS (use 'make certs' to enable)")
	}

	grpcServer := grpc.NewServer(opts...)

	// Register services
	pb.RegisterUserServiceServer(grpcServer, s)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	color.Green("✓ Server registered and reflection enabled")
	color.Blue("📡 Listening on %s", port)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		color.Yellow("\n⏸️  Shutting down server...")
		grpcServer.GracefulStop()
		color.Green("✓ Server stopped gracefully")
		os.Exit(0)
	}()

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}

func (s *server) initializeData() {
	s.users["1"] = &pb.User{
		Id:    "1",
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Age:   28,
		Tags:  []string{"premium", "active"},
	}
	s.users["2"] = &pb.User{
		Id:    "2",
		Name:  "Bob Smith",
		Email: "bob@example.com",
		Age:   35,
		Tags:  []string{"standard"},
	}
	s.users["3"] = &pb.User{
		Id:    "3",
		Name:  "Charlie Brown",
		Email: "charlie@example.com",
		Age:   42,
		Tags:  []string{"premium", "vip"},
	}
}

// Unary RPC: GetUser
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	color.Yellow("📥 GetUser request: ID=%s", req.Id)

	// Check deadline
	if dl, ok := ctx.Deadline(); ok {
		color.Cyan("⏰ Request deadline: %v", dl)
	}

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user with ID %s not found", req.Id))
	}

	return &pb.GetUserResponse{User: user}, nil
}

// Server Streaming: ListUsers
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
	color.Yellow("📥 ListUsers request: limit=%d, filter=%s", req.Limit, req.Filter)

	count := 0
	for _, user := range s.users {
		// Check if context is cancelled
		if stream.Context().Err() != nil {
			color.Red("❌ Context cancelled")
			return stream.Context().Err()
		}

		// Apply filter if specified
		if req.Filter != "" {
			found := false
			for _, tag := range user.Tags {
				if tag == req.Filter {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Send user
		if err := stream.Send(user); err != nil {
			color.Red("❌ Failed to send user: %v", err)
			return err
		}

		color.Green("✓ Sent user: %s", user.Name)
		count++

		// Check limit
		if req.Limit > 0 && count >= int(req.Limit) {
			break
		}

		// Simulate delay between sends
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// Client Streaming: CreateUsers
func (s *server) CreateUsers(stream pb.UserService_CreateUsersServer) error {
	color.Yellow("📥 CreateUsers stream started")

	var createdUsers []*pb.User
	count := 0

	for {
		req, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			color.Red("❌ Failed to receive: %v", err)
			return err
		}

		// Validate input
		if req.Name == "" {
			return status.Error(codes.InvalidArgument, "name is required")
		}
		if req.Email == "" {
			return status.Error(codes.InvalidArgument, "email is required")
		}

		// Create user
		user := &pb.User{
			Id:    fmt.Sprintf("%d", len(s.users)+count+1),
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
			Tags:  req.Tags,
		}
		createdUsers = append(createdUsers, user)
		count++

		color.Green("✓ Created user: %s (%s)", user.Name, user.Email)
	}

	// Send response
	resp := &pb.CreateUsersResponse{
		Users:   createdUsers,
		Count:   int32(count),
		Message: fmt.Sprintf("Successfully created %d users", count),
	}

	if err := stream.SendAndClose(resp); err != nil {
		color.Red("❌ Failed to send response: %v", err)
		return err
	}

	color.Green("✓ Created %d users", count)
	return nil
}

// Bidirectional Streaming: Chat
func (s *server) Chat(stream pb.UserService_ChatServer) error {
	color.Yellow("💬 Chat session started")

	// Handle incoming messages in goroutine
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					color.Yellow("💬 Client closed connection")
					return
				}
				color.Red("❌ Failed to receive: %v", err)
				return
			}

			color.Cyan("📨 Received from %s: %s", req.UserId, req.Message)

			// Echo back with timestamp
			resp := &pb.ChatMessage{
				UserId:    "server",
				Message:   fmt.Sprintf("Echo: %s", req.Message),
				Timestamp: time.Now().Unix(),
				Type:      pb.MessageType_MESSAGE,
			}

			if err := stream.Send(resp); err != nil {
				color.Red("❌ Failed to send: %v", err)
				return
			}

			color.Green("✓ Sent response to %s", req.UserId)
		}
	}()

	// Keep stream open
	<-stream.Context().Done()
	color.Yellow("💬 Chat session ended")
	return nil
}

// Logging interceptor for unary RPCs
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	color.Cyan("🔍 [UNARY] Method: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	elapsed := time.Since(start)
	if err != nil {
		color.Red("❌ [UNARY] %s failed: %v (took %v)", info.FullMethod, err, elapsed)
	} else {
		color.Green("✓ [UNARY] %s succeeded (took %v)", info.FullMethod, elapsed)
	}

	return resp, err
}

// Logging interceptor for streaming RPCs
func streamLoggingInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()

	color.Cyan("🔍 [STREAM] Method: %s", info.FullMethod)

	err := handler(srv, stream)

	elapsed := time.Since(start)
	if err != nil {
		color.Red("❌ [STREAM] %s failed: %v (took %v)", info.FullMethod, err, elapsed)
	} else {
		color.Green("✓ [STREAM] %s succeeded (took %v)", info.FullMethod, elapsed)
	}

	return err
}
