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

// NewRepositoriesSimple creates repository instances with simple approach.
func NewRepositoriesSimple() *Repositories {
	return &Repositories{
		Course: NewMemoryCourseRepository(),
		// TODO: Implement other repositories as needed.
	}
}
