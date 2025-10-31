package main

import (
	"context"
	"fmt"
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

func (s *UserStore) GetByID(id int32) (*pb.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

func (s *UserStore) GetAll(role string, active bool) []*pb.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*pb.User, 0)
	for _, user := range s.users {
		if role != "" && user.Role != role {
			continue
		}
		if !active && user.Active {
			continue
		}
		users = append(users, user)
	}
	return users
}

func (s *UserStore) Update(id int32, username, email, role string, active bool) (*pb.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, false
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if role != "" {
		user.Role = role
	}
	user.Active = active
	user.UpdatedAt = timestamppb.Now()
	return user, true
}

func (s *UserStore) Delete(id int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[id]
	if exists {
		delete(s.users, id)
	}
	return exists
}

// UserServiceServer implements the gRPC UserService
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	store *UserStore
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		store: NewUserStore(),
	}
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	log.Printf("CreateUser called: username=%s, email=%s, role=%s", req.Username, req.Email, req.Role)

	// Validation
	if req.Username == "" || req.Email == "" || req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email, and role are required")
	}

	user := s.store.Create(req.Username, req.Email, req.Role)
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Printf("GetUser called: id=%d", req.Id)

	user, exists := s.store.GetByID(req.Id)
	if !exists {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	log.Printf("UpdateUser called: id=%d", req.Id)

	user, exists := s.store.Update(req.Id, req.Username, req.Email, req.Role, req.Active)
	if !exists {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteUser called: id=%d", req.Id)

	if !s.store.Delete(req.Id) {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &emptypb.Empty{}, nil
}

// ListUsers lists all users with optional filtering
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("ListUsers called: role=%s, active=%v", req.Role, req.Active)

	users := s.store.GetAll(req.Role, req.Active)

	// Simple pagination
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	page := req.Page
	if page == 0 {
		page = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > int32(len(users)) {
		start = int32(len(users))
	}
	if end > int32(len(users)) {
		end = int32(len(users))
	}

	return &pb.ListUsersResponse{
		Users:    users[start:end],
		Total:    int32(len(users)),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Server streaming and bidirectional streaming are in the streaming example
func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
	return status.Error(codes.Unimplemented, "see streaming example")
}

func (s *UserServiceServer) CreateUsers(stream pb.UserService_CreateUsersServer) error {
	return status.Error(codes.Unimplemented, "see streaming example")
}

func (s *UserServiceServer) Chat(stream pb.UserService_ChatServer) error {
	return status.Error(codes.Unimplemented, "see streaming example")
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

	fmt.Println("🚀 gRPC server starting on :9090")
	fmt.Println("📚 Available methods:")
	fmt.Println("  - CreateUser")
	fmt.Println("  - GetUser")
	fmt.Println("  - UpdateUser")
	fmt.Println("  - DeleteUser")
	fmt.Println("  - ListUsers")

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("\n🔧 Testing gRPC client...")

	// 1. Create user
	fmt.Println("\n1️⃣ Creating user...")
	user, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: "david",
		Email:    "david@example.com",
		Role:     "user",
	})
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
	} else {
		fmt.Printf("✅ Created user: %+v\n", user)
	}

	// 2. Get user
	fmt.Println("\n2️⃣ Getting user...")
	user, err = client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Printf("GetUser failed: %v", err)
	} else {
		fmt.Printf("✅ Got user: %+v\n", user)
	}

	// 3. List users
	fmt.Println("\n3️⃣ Listing users...")
	resp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("ListUsers failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d users (total: %d)\n", len(resp.Users), resp.Total)
		for _, u := range resp.Users {
			fmt.Printf("   - %s (%s) - %s\n", u.Username, u.Email, u.Role)
		}
	}

	// 4. Update user
	fmt.Println("\n4️⃣ Updating user...")
	user, err = client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       1,
		Username: "alice-updated",
		Active:   true,
	})
	if err != nil {
		log.Printf("UpdateUser failed: %v", err)
	} else {
		fmt.Printf("✅ Updated user: %+v\n", user)
	}

	// 5. Delete user
	fmt.Println("\n5️⃣ Deleting user...")
	_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: 3})
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
	} else {
		fmt.Println("✅ User deleted successfully")
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
	select {}
}

