// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package service

import (
	"context"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/repository"
)

// UserService defines business logic for user management operations.
type UserService interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)

	// GetAllUsers retrieves all users with pagination
	GetAllUsers(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int, error)

	// DeleteUser deletes a user by ID
	DeleteUser(ctx context.Context, userID string) error
}

// userService implements UserService.
type userService struct {
	repo   repository.UserRepository
	config *Config
}

// NewUserService creates a new user service.
func NewUserService(userRepo repository.UserRepository, config *Config) UserService {
	return &userService{
		repo:   userRepo,
		config: config,
	}
}

// GetUserByID retrieves a user by ID.
func (s *userService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, apierrors.NewBadRequestError("user ID is required")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewNotFoundError("user not found")
	}

	return user, nil
}

// GetAllUsers retrieves all users with pagination.
func (s *userService) GetAllUsers(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int, error) {
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	users, total, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		return nil, 0, apierrors.NewInternalError("failed to retrieve users", err)
	}

	return users, int(total), nil
}

// DeleteUser deletes a user by ID.
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return apierrors.NewBadRequestError("user ID is required")
	}

	// Check if user exists.
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return apierrors.NewNotFoundError("user not found")
	}

	// Delete user.
	if err := s.repo.Delete(ctx, userID); err != nil {
		return apierrors.NewInternalError("failed to delete user", err)
	}

	s.config.Logger.Info(ctx, "user deleted", "user_id", userID)
	return nil
}
