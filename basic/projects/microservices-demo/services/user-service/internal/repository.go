package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for user storage
type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id string, user *User) error
	Delete(id string) error
	List(page, pageSize int) ([]*User, int, error)
}

// InMemoryRepository implements Repository using in-memory storage
type InMemoryRepository struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[string]*User),
	}
}

// Create creates a new user
func (r *InMemoryRepository) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return fmt.Errorf("email already exists")
		}
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryRepository) GetByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryRepository) GetByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// Update updates a user
func (r *InMemoryRepository) Update(id string, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}

	// Update fields
	if user.Username != "" {
		existing.Username = user.Username
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	existing.UpdatedAt = time.Now()

	return nil
}

// Delete deletes a user
func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	return nil
}

// List returns a paginated list of users
func (r *InMemoryRepository) List(page, pageSize int) ([]*User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.users)
	start := (page - 1) * pageSize
	end := start + pageSize

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	if start >= len(users) {
		return []*User{}, total, nil
	}

	if end > len(users) {
		end = len(users)
	}

	return users[start:end], total, nil
}

