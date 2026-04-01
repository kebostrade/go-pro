// Package service provides the UserService implementation.
package service

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/DimaJoyti/go-pro/basic/projects/grpc-service/proto"
)

// UserServiceServer implements the UserService gRPC server.
type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer

	// In-memory user store
	users map[string]*userpb.User
	mu    sync.RWMutex
}

// NewUserServiceServer creates a new UserServiceServer.
func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		users: initUserStore(),
	}
}

// initUserStore creates the initial user data.
func initUserStore() map[string]*userpb.User {
	return map[string]*userpb.User{
		"1": {
			Id:    "1",
			Name:  "Alice Johnson",
			Email: "alice@example.com",
			Age:   28,
			Tags:  []string{"admin", "developer"},
		},
		"2": {
			Id:    "2",
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Age:   34,
			Tags:  []string{"developer"},
		},
		"3": {
			Id:    "3",
			Name:  "Charlie Brown",
			Email: "charlie@example.com",
			Age:   22,
			Tags:  []string{"designer"},
		},
		"4": {
			Id:    "4",
			Name:  "Diana Prince",
			Email: "diana@example.com",
			Age:   30,
			Tags:  []string{"manager", "developer"},
		},
		"5": {
			Id:    "5",
			Name:  "Eve Williams",
			Email: "eve@example.com",
			Age:   26,
			Tags:  []string{"developer", "tester"},
		},
	}
}

// GetUser implements the Unary RPC - returns a single user by ID.
func (s *UserServiceServer) GetUser(
	ctx context.Context,
	req *userpb.GetUserRequest,
) (*userpb.GetUserResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	s.mu.RLock()
	user, ok := s.users[req.Id]
	s.mu.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %q not found", req.Id)
	}

	// Return a copy to avoid race conditions.
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Tags:  user.Tags,
		},
	}, nil
}

// ListUsers implements Server Streaming RPC - streams all users.
func (s *UserServiceServer) ListUsers(
	req *userpb.ListUsersRequest,
	stream userpb.UserService_ListUsersServer,
) error {
	if req == nil {
		req = &userpb.ListUsersRequest{Limit: 10}
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, user := range s.users {
		if count >= limit {
			break
		}

		// Send user to stream.
		if err := stream.Send(&userpb.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Tags:  user.Tags,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send user: %v", err)
		}
		count++
	}

	return nil
}

// CreateUsers implements Client Streaming RPC - receives multiple users, returns once.
func (s *UserServiceServer) CreateUsers(
	stream userpb.UserService_CreateUsersServer,
) error {
	var createdUsers []*userpb.User

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Finished receiving all users.
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive: %v", err)
		}

		// Generate ID for new user.
		s.mu.Lock()
		newId := fmt.Sprintf("%d", len(s.users)+1)
		newUser := &userpb.User{
			Id:    newId,
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
			Tags:  req.Tags,
		}
		s.users[newId] = newUser
		s.mu.Unlock()

		createdUsers = append(createdUsers, newUser)
	}

	// Send response with all created users.
	return stream.SendAndClose(&userpb.CreateUsersResponse{
		Users: createdUsers,
		Count: int32(len(createdUsers)),
	})
}

// Chat implements Bidirectional Streaming RPC - echo chat messages.
func (s *UserServiceServer) Chat(
	stream userpb.UserService_ChatServer,
) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive: %v", err)
		}

		// Echo the message back with timestamp.
		echoMsg := &userpb.ChatMessage{
			UserId:    msg.UserId,
			Message:   fmt.Sprintf("[Echo] %s", msg.Message),
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(echoMsg); err != nil {
			return status.Errorf(codes.Internal, "failed to send: %v", err)
		}
	}
}

// GetUserByID is exported for testing.
func (s *UserServiceServer) GetUserByID(id string) (*userpb.User, error) {
	s.mu.RLock()
	user, ok := s.users[id]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}

	return &userpb.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
		Tags:  user.Tags,
	}, nil
}

// GetAllUsers returns all users.
func (s *UserServiceServer) GetAllUsers() []*userpb.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*userpb.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, &userpb.User{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
			Age:   u.Age,
			Tags:  u.Tags,
		})
	}
	return users
}

// UserCount returns the number of users.
func (s *UserServiceServer) UserCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}
