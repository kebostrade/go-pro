package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// MemoryUserRepository implements UserRepository using in-memory storage
type MemoryUserRepository struct {
	mu      sync.RWMutex
	users   map[int64]*model.User
	nextID  int64
	byEmail map[string]*model.User
	byUsername map[string]*model.User
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:      make(map[int64]*model.User),
		byEmail:    make(map[string]*model.User),
		byUsername: make(map[string]*model.User),
		nextID:     1,
	}
}

// Create creates a new user
func (r *MemoryUserRepository) Create(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.byEmail[user.Email]; exists {
		return ErrUserAlreadyExists
	}

	// Check if username already exists
	if _, exists := r.byUsername[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	// Assign ID and timestamps
	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Store user
	r.users[user.ID] = user
	r.byEmail[user.Email] = user
	r.byUsername[user.Username] = user

	return nil
}

// GetByID retrieves a user by ID
func (r *MemoryUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.byEmail[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *MemoryUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.byUsername[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// List retrieves a list of users with pagination
func (r *MemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Apply pagination
	start := offset
	if start > len(users) {
		start = len(users)
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// Update updates a user
func (r *MemoryUserRepository) Update(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.users[user.ID]
	if !exists {
		return ErrUserNotFound
	}

	// Update email index if changed
	if existing.Email != user.Email {
		delete(r.byEmail, existing.Email)
		r.byEmail[user.Email] = user
	}

	// Update username index if changed
	if existing.Username != user.Username {
		delete(r.byUsername, existing.Username)
		r.byUsername[user.Username] = user
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// Delete deletes a user
func (r *MemoryUserRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.byEmail, user.Email)
	delete(r.byUsername, user.Username)

	return nil
}

