package clean

import (
	"context"
	"time"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = &DomainError{Message: "user not found"}

// ErrUserAlreadyExists is returned when attempting to create a duplicate user.
var ErrUserAlreadyExists = &DomainError{Message: "user already exists"}

// DomainError represents a domain-level error.
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// CreateUserUseCase handles user creation.
type CreateUserUseCase struct {
	repo UserRepository
}

// NewCreateUserUseCase creates a new CreateUserUseCase.
func NewCreateUserUseCase(repo UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

// Execute creates a new user.
func (uc *CreateUserUseCase) Execute(ctx context.Context, email, name string) (*User, error) {
	// Check if user already exists
	existing, _ := uc.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &User{
		ID:    generateID(),
		Email: email,
		Name:  name,
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserUseCase handles retrieving a user.
type GetUserUseCase struct {
	repo UserRepository
}

// NewGetUserUseCase creates a new GetUserUseCase.
func NewGetUserUseCase(repo UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

// Execute retrieves a user by ID.
func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*User, error) {
	return uc.repo.GetByID(ctx, id)
}

// UpdateUserUseCase handles user updates.
type UpdateUserUseCase struct {
	repo UserRepository
}

// NewUpdateUserUseCase creates a new UpdateUserUseCase.
func NewUpdateUserUseCase(repo UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

// Execute updates an existing user.
func (uc *UpdateUserUseCase) Execute(ctx context.Context, user *User) error {
	return uc.repo.Update(ctx, user)
}

// DeleteUserUseCase handles user deletion.
type DeleteUserUseCase struct {
	repo UserRepository
}

// NewDeleteUserUseCase creates a new DeleteUserUseCase.
func NewDeleteUserUseCase(repo UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

// Execute deletes a user by ID.
func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// ListUsersUseCase handles listing users.
type ListUsersUseCase struct {
	repo UserRepository
}

// NewListUsersUseCase creates a new ListUsersUseCase.
func NewListUsersUseCase(repo UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

// Execute returns all users with pagination.
func (uc *ListUsersUseCase) Execute(ctx context.Context, limit, offset int) ([]*User, error) {
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.List(ctx, limit, offset)
}

// generateID generates a simple ID (in production, use UUID).
func generateID() string {
	return "user-" + randomString(8)
}

// randomString generates a random string of given length.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[fastRandIntn(len(charset))]
	}
	return string(result)
}

// fastRandIntn returns a pseudo-random number in [0, n).
func fastRandIntn(n int) int {
	// Simple pseudo-random generator for test stability
	return (int(time.Now().UnixNano())%n + n) % n
}
