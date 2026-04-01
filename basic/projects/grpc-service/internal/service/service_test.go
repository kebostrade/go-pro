package service

import (
	"testing"

	userpb "github.com/DimaJoyti/go-pro/basic/projects/grpc-service/proto"
)

func TestNewUserServiceServer(t *testing.T) {
	srv := NewUserServiceServer()
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if srv.UserCount() != 5 {
		t.Errorf("expected 5 users, got %d", srv.UserCount())
	}
}

func TestGetUserByID(t *testing.T) {
	srv := NewUserServiceServer()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"existing user 1", "1", false},
		{"existing user 2", "2", false},
		{"existing user 3", "3", false},
		{"existing user 4", "4", false},
		{"existing user 5", "5", false},
		{"nonexistent user", "999", true},
		{"empty id", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := srv.GetUserByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserByIDReturnsCorrectUser(t *testing.T) {
	srv := NewUserServiceServer()

	user, err := srv.GetUserByID("1")
	if err != nil {
		t.Fatalf("GetUserByID() error = %v", err)
	}

	if user.Name != "Alice Johnson" {
		t.Errorf("expected user name 'Alice Johnson', got %s", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected user email 'alice@example.com', got %s", user.Email)
	}
	if user.Age != 28 {
		t.Errorf("expected user age 28, got %d", user.Age)
	}
}

func TestGetAllUsers(t *testing.T) {
	srv := NewUserServiceServer()

	users := srv.GetAllUsers()
	if len(users) != 5 {
		t.Errorf("expected 5 users, got %d", len(users))
	}
}

func TestUserCount(t *testing.T) {
	srv := NewUserServiceServer()

	tests := []struct {
		name     string
		expected int
	}{
		{"initial count", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := srv.UserCount()
			if count != tt.expected {
				t.Errorf("UserCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestUserServiceServerImplementsInterface(t *testing.T) {
	srv := NewUserServiceServer()

	// Verify the server implements the protobuf interface
	var _ userpb.UserServiceServer = srv
}
