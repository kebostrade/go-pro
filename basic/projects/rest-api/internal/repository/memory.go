package repository

import (
	"context"
	"sync"

	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/domain"
	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/pkg/errors"
)

// UserRepository interface defines data access operations for users
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id, name, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

// MemoryUserStore implements UserRepository with in-memory storage
type MemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewMemoryUserStore creates a new in-memory user store with seed data
func NewMemoryUserStore() *MemoryUserStore {
	store := &MemoryUserStore{
		users: make(map[string]*domain.User),
	}

	// Seed with sample users for testing
	store.users["seed-1"] = &domain.User{
		ID:    "seed-1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	store.users["seed-2"] = &domain.User{
		ID:    "seed-2",
		Name:  "Jane Smith",
		Email: "jane@example.com",
	}
	store.users["seed-3"] = &domain.User{
		ID:    "seed-3",
		Name:  "Bob Wilson",
		Email: "bob@example.com",
	}

	return store
}

// Create adds a new user to the store
func (s *MemoryUserStore) Create(ctx context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate simple ID
	user.ID = generateID()
	s.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID
func (s *MemoryUserStore) GetByID(ctx context.Context, id string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if user, ok := s.users[id]; ok {
		return user, nil
	}
	return nil, errors.ErrNotFound
}

// List retrieves all users
func (s *MemoryUserStore) List(ctx context.Context) ([]*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*domain.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

// Update updates an existing user
func (s *MemoryUserStore) Update(ctx context.Context, id, name, email string) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, ok := s.users[id]; ok {
		user.Name = name
		user.Email = email
		return user, nil
	}
	return nil, errors.ErrNotFound
}

// Delete removes a user by ID
func (s *MemoryUserStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; ok {
		delete(s.users, id)
		return nil
	}
	return errors.ErrNotFound
}

// generateID creates a simple unique ID
func generateID() string {
	return "user-" + randomString(8)
}

// randomString generates a random alphanumeric string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
