package service

import (
	"context"
	"errors"
)

// User represents a user entity
type User struct {
	ID    string
	Name  string
	Email string
}

// UserRepository interface for data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id string) error
}

// ErrNotFound when user not found
var ErrNotFound = errors.New("user not found")

// ErrValidation for invalid input
var ErrValidation = errors.New("validation error")

// UserService handles user business logic
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, ErrValidation
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List retrieves all users
func (s *UserService) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

// Delete removes a user
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
