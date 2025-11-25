// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package repository provides functionality for the GO-PRO Learning Platform.
package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
)

// MemoryCourseRepository implements CourseRepository using in-memory storage.
type MemoryCourseRepository struct {
	courses map[string]*domain.Course
	mu      sync.RWMutex
}

// NewMemoryCourseRepository creates a new in-memory course repository.
func NewMemoryCourseRepository() *MemoryCourseRepository {
	return &MemoryCourseRepository{
		courses: make(map[string]*domain.Course),
	}
}

// Create implements CourseRepository.Create.
func (r *MemoryCourseRepository) Create(ctx context.Context, course *domain.Course) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[course.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("course with id %s already exists", course.ID))
	}

	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	r.courses[course.ID] = course

	return nil
}

// GetByID implements CourseRepository.GetByID.
func (r *MemoryCourseRepository) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	course, exists := r.courses[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
	}

	// Return a copy to prevent modification.
	courseCopy := *course

	return &courseCopy, nil
}

// GetAll implements CourseRepository.GetAll.
func (r *MemoryCourseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var courses []*domain.Course
	for _, course := range r.courses {
		courseCopy := *course
		courses = append(courses, &courseCopy)
	}

	// Sort by creation time (newest first)
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].CreatedAt.After(courses[j].CreatedAt)
	})

	total := int64(len(courses))

	// Apply pagination if provided.
	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(courses) {
			return []*domain.Course{}, total, nil
		}

		if end > len(courses) {
			end = len(courses)
		}

		courses = courses[start:end]
	}

	return courses, total, nil
}

// Update implements CourseRepository.Update.
func (r *MemoryCourseRepository) Update(ctx context.Context, course *domain.Course) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[course.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", course.ID))
	}

	course.UpdatedAt = time.Now()
	r.courses[course.ID] = course

	return nil
}

// Delete implements CourseRepository.Delete.
func (r *MemoryCourseRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[id]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
	}

	delete(r.courses, id)

	return nil
}

// MemoryUserRepository implements UserRepository using in-memory storage.
type MemoryUserRepository struct {
	users           map[string]*domain.User            // Keyed by user ID
	usersByFirebase map[string]*domain.User            // Keyed by Firebase UID
	usersByEmail    map[string]*domain.User            // Keyed by email
	mu              sync.RWMutex
}

// NewMemoryUserRepository creates a new in-memory user repository.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:           make(map[string]*domain.User),
		usersByFirebase: make(map[string]*domain.User),
		usersByEmail:    make(map[string]*domain.User),
	}
}

// Create implements UserRepository.Create.
func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with id %s already exists", user.ID))
	}

	if _, exists := r.usersByFirebase[user.FirebaseUID]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with Firebase UID %s already exists", user.FirebaseUID))
	}

	if _, exists := r.usersByEmail[user.Email]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with email %s already exists", user.Email))
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	r.usersByFirebase[user.FirebaseUID] = user
	r.usersByEmail[user.Email] = user

	return nil
}

// GetByID implements UserRepository.GetByID.
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetByFirebaseUID implements UserRepository.GetByFirebaseUID.
func (r *MemoryUserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByFirebase[firebaseUID]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with Firebase UID %s not found", firebaseUID))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetByEmail implements UserRepository.GetByEmail.
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetAll implements UserRepository.GetAll.
func (r *MemoryUserRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var users []*domain.User
	for _, user := range r.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	// Sort by creation time (newest first)
	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})

	total := int64(len(users))

	// Apply pagination
	if pagination != nil && pagination.PageSize > 0 {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(users) {
			return []*domain.User{}, total, nil
		}

		if end > len(users) {
			end = len(users)
		}

		users = users[start:end]
	}

	return users, total, nil
}

// Update implements UserRepository.Update.
func (r *MemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.users[user.ID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", user.ID))
	}

	// Update indexes if email or Firebase UID changed
	if existing.Email != user.Email {
		delete(r.usersByEmail, existing.Email)
		r.usersByEmail[user.Email] = user
	}

	if existing.FirebaseUID != user.FirebaseUID {
		delete(r.usersByFirebase, existing.FirebaseUID)
		r.usersByFirebase[user.FirebaseUID] = user
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// UpdateLastLogin implements UserRepository.UpdateLastLogin.
func (r *MemoryUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", userID))
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	return nil
}

// Delete implements UserRepository.Delete.
func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	delete(r.users, id)
	delete(r.usersByFirebase, user.FirebaseUID)
	delete(r.usersByEmail, user.Email)

	return nil
}

// NewRepositoriesSimple creates repository instances with simple approach.
func NewRepositoriesSimple() *Repositories {
	return &Repositories{
		Course: NewMemoryCourseRepository(),
		User:   NewMemoryUserRepository(),
		// TODO: Implement other repositories as needed.
	}
}
