// Package clean provides clean architecture domain entities and repository interfaces.
package clean

import (
	"context"
	"time"
)

// User represents a user entity.
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Order represents an order entity.
type Order struct {
	ID        string
	UserID    string
	Total     float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// OrderRepository defines the interface for order data access.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
}

// InMemoryUserRepository implements UserRepository with in-memory storage.
type InMemoryUserRepository struct {
	users map[string]*User
}

// NewInMemoryUserRepository creates a new in-memory user repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

// Create creates a new user.
func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return ErrUserAlreadyExists
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID.
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// Update updates an existing user.
func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// Delete deletes a user by ID.
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

// List returns all users with pagination.
func (r *InMemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	if offset >= len(users) {
		return []*User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}
