// GO-PRO Learning Platform Backend
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

type MemoryLessonRepository struct {
	lessons  map[string]*domain.Lesson
	byCourse map[string][]*domain.Lesson
	mu       sync.RWMutex
}

func NewMemoryLessonRepository() *MemoryLessonRepository {
	return &MemoryLessonRepository{
		lessons:  make(map[string]*domain.Lesson),
		byCourse: make(map[string][]*domain.Lesson),
	}
}

func (r *MemoryLessonRepository) Create(ctx context.Context, lesson *domain.Lesson) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.lessons[lesson.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("lesson with id %s already exists", lesson.ID))
	}

	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()
	r.lessons[lesson.ID] = lesson
	r.byCourse[lesson.CourseID] = append(r.byCourse[lesson.CourseID], lesson)

	return nil
}

func (r *MemoryLessonRepository) GetByID(ctx context.Context, id string) (*domain.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lesson, exists := r.lessons[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", id))
	}

	lessonCopy := *lesson
	return &lessonCopy, nil
}

func (r *MemoryLessonRepository) GetByCourseID(ctx context.Context, courseID string, pagination *domain.PaginationRequest) ([]*domain.Lesson, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lessons := r.byCourse[courseID]
	if lessons == nil {
		return []*domain.Lesson{}, 0, nil
	}

	var result []*domain.Lesson
	for _, lesson := range lessons {
		lessonCopy := *lesson
		result = append(result, &lessonCopy)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})

	total := int64(len(result))

	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(result) {
			return []*domain.Lesson{}, total, nil
		}

		if end > len(result) {
			end = len(result)
		}

		result = result[start:end]
	}

	return result, total, nil
}

func (r *MemoryLessonRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Lesson, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var lessons []*domain.Lesson
	for _, lesson := range r.lessons {
		lessonCopy := *lesson
		lessons = append(lessons, &lessonCopy)
	}

	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Order < lessons[j].Order
	})

	total := int64(len(lessons))

	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(lessons) {
			return []*domain.Lesson{}, total, nil
		}

		if end > len(lessons) {
			end = len(lessons)
		}

		lessons = lessons[start:end]
	}

	return lessons, total, nil
}

func (r *MemoryLessonRepository) Update(ctx context.Context, lesson *domain.Lesson) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.lessons[lesson.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", lesson.ID))
	}

	lesson.UpdatedAt = time.Now()
	r.lessons[lesson.ID] = lesson

	if lessons := r.byCourse[lesson.CourseID]; lessons != nil {
		for i, l := range lessons {
			if l.ID == lesson.ID {
				lessons[i] = lesson
				break
			}
		}
	}

	return nil
}

func (r *MemoryLessonRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lesson, exists := r.lessons[id]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", id))
	}

	courseID := lesson.CourseID
	delete(r.lessons, id)

	if lessons := r.byCourse[courseID]; lessons != nil {
		for i, l := range lessons {
			if l.ID == id {
				r.byCourse[courseID] = append(lessons[:i], lessons[i+1:]...)
				break
			}
		}
	}

	return nil
}
