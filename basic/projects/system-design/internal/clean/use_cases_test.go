package clean

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateUserUseCase(t *testing.T) {
	repo := NewInMemoryUserRepository()
	uc := NewCreateUserUseCase(repo)

	user, err := uc.Execute(context.Background(), "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("CreateUserUseCase failed: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", user.Email)
	}

	if user.Name != "Test User" {
		t.Errorf("Name = %s, want Test User", user.Name)
	}
}

func TestCreateUserUseCase_Duplicate(t *testing.T) {
	repo := NewInMemoryUserRepository()
	uc := NewCreateUserUseCase(repo)

	_, err := uc.Execute(context.Background(), "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("First creation failed: %v", err)
	}

	_, err = uc.Execute(context.Background(), "test@example.com", "Another User")
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestGetUserUseCase(t *testing.T) {
	repo := NewInMemoryUserRepository()
	createUC := NewCreateUserUseCase(repo)
	getUC := NewGetUserUseCase(repo)

	created, _ := createUC.Execute(context.Background(), "test@example.com", "Test User")

	retrieved, err := getUC.Execute(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetUserUseCase failed: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("ID = %s, want %s", retrieved.ID, created.ID)
	}
}

func TestGetUserUseCase_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()
	getUC := NewGetUserUseCase(repo)

	_, err := getUC.Execute(context.Background(), "nonexistent")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUserUseCase(t *testing.T) {
	repo := NewInMemoryUserRepository()
	createUC := NewCreateUserUseCase(repo)
	updateUC := NewUpdateUserUseCase(repo)

	user, _ := createUC.Execute(context.Background(), "test@example.com", "Original Name")

	user.Name = "Updated Name"
	err := updateUC.Execute(context.Background(), user)
	if err != nil {
		t.Fatalf("UpdateUserUseCase failed: %v", err)
	}

	updated, _ := repo.GetByID(context.Background(), user.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("Name = %s, want Updated Name", updated.Name)
	}
}

func TestDeleteUserUseCase(t *testing.T) {
	repo := NewInMemoryUserRepository()
	createUC := NewCreateUserUseCase(repo)
	deleteUC := NewDeleteUserUseCase(repo)

	user, _ := createUC.Execute(context.Background(), "test@example.com", "Test User")

	err := deleteUC.Execute(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("DeleteUserUseCase failed: %v", err)
	}

	_, err = repo.GetByID(context.Background(), user.ID)
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound after deletion, got %v", err)
	}
}

func TestListUsersUseCase(t *testing.T) {
	repo := NewInMemoryUserRepository()
	createUC := NewCreateUserUseCase(repo)
	listUC := NewListUsersUseCase(repo)

	// Create 5 users with unique emails
	for i := 0; i < 5; i++ {
		_, err := createUC.Execute(context.Background(), fmt.Sprintf("test%d@example.com", i), "Test User")
		if err != nil {
			t.Fatalf("CreateUserUseCase %d failed: %v", i, err)
		}
	}

	users, err := listUC.Execute(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("ListUsersUseCase failed: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Got %d users, want 5", len(users))
	}
}

func TestListUsersUseCase_Pagination(t *testing.T) {
	repo := NewInMemoryUserRepository()
	createUC := NewCreateUserUseCase(repo)
	listUC := NewListUsersUseCase(repo)

	// Create 5 users with unique emails
	for i := 0; i < 5; i++ {
		_, err := createUC.Execute(context.Background(), fmt.Sprintf("test%d@example.com", i), "Test User")
		if err != nil {
			t.Fatalf("CreateUserUseCase %d failed: %v", i, err)
		}
	}

	users, err := listUC.Execute(context.Background(), 2, 0)
	if err != nil {
		t.Fatalf("ListUsersUseCase failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Got %d users, want 2", len(users))
	}

	users, err = listUC.Execute(context.Background(), 2, 2)
	if err != nil {
		t.Fatalf("ListUsersUseCase failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Got %d users, want 2", len(users))
	}
}
