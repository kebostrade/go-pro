package service

import (
	"context"
	"regexp"

	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/domain"
	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/repository"
	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/pkg/errors"
)

// UserService interface defines the business operations for users
type UserService interface {
	Create(ctx context.Context, name, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id, name, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// emailRegex validates email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Create creates a new user
func (s *userService) Create(ctx context.Context, name, email string) (*domain.User, error) {
	// Validate email format
	if !emailRegex.MatchString(email) {
		return nil, errors.ErrValidation
	}

	// Create user
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List retrieves all users
func (s *userService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id, name, email string) (*domain.User, error) {
	// Validate email format
	if !emailRegex.MatchString(email) {
		return nil, errors.ErrValidation
	}

	// Update user
	user, err := s.repo.Update(ctx, id, name, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user by ID
func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
