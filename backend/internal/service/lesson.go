// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// lessonService implements the LessonService interface.
type lessonService struct {
	repo      repository.LessonRepository
	logger    logger.Logger
	validator validator.Validator
	cache     cache.CacheManager
	messaging *messaging.Service
}

// NewLessonService creates a new lesson service.
func NewLessonService(
	repo repository.LessonRepository,
	config *Config,
) LessonService {
	return &lessonService{
		repo:      repo,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}
}

// CreateLesson creates a new lesson.
func (s *lessonService) CreateLesson(ctx context.Context, req *domain.CreateLessonRequest) (*domain.Lesson, error) {
	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid lesson creation request", "error", err)
		return nil, errors.NewValidationError("invalid lesson data", err)
	}

	// Create lesson entity.
	lesson := &domain.Lesson{
		ID:          generateID("lesson"),
		CourseID:    req.CourseID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Content:     req.Content,
		Order:       req.Order,
		Exercises:   []string{}, // Initialize empty exercises slice
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository.
	if err := s.repo.Create(ctx, lesson); err != nil {
		s.logger.Error(ctx, "Failed to create lesson", "error", err, "lesson_id", lesson.ID)
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}

	// Publish lesson created event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		lessonData := map[string]interface{}{
			"title":       lesson.Title,
			"description": lesson.Description,
			"content":     lesson.Content,
			"order":       lesson.Order,
		}
		if err := s.messaging.PublishLessonCreated(ctx, lesson.ID, lesson.CourseID, lessonData); err != nil {
			s.logger.Warn(ctx, "Failed to publish lesson created event", "error", err, "lesson_id", lesson.ID)
		}
	}

	// Cache the lesson.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("lesson:%s", lesson.ID)
		if err := s.cache.Set(ctx, cacheKey, lesson, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache lesson", "error", err, "lesson_id", lesson.ID)
		}
	}

	s.logger.Info(ctx, "Lesson created successfully", "lesson_id", lesson.ID, "course_id", lesson.CourseID)

	return lesson, nil
}

// GetLessonByID retrieves a lesson by ID.
func (s *lessonService) GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("lesson ID is required")
	}

	// Try cache first.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("lesson:%s", id)
		var cachedLesson domain.Lesson
		if err := s.cache.Get(ctx, cacheKey, &cachedLesson); err == nil {
			s.logger.Debug(ctx, "Lesson retrieved from cache", "lesson_id", id)
			return &cachedLesson, nil
		}
	}

	lesson, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get lesson", "error", err, "lesson_id", id)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	// Cache the lesson.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("lesson:%s", id)
		if err := s.cache.Set(ctx, cacheKey, lesson, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache lesson", "error", err, "lesson_id", id)
		}
	}

	return lesson, nil
}

// GetLessonsByCourseID retrieves lessons for a specific course.
func (s *lessonService) GetLessonsByCourseID(
	ctx context.Context,
	courseID string,
	pagination *domain.PaginationRequest,
) (*domain.ListResponse, error) {
	if courseID == "" {
		return nil, errors.NewBadRequestError("course ID is required")
	}

	// Set default pagination.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	lessons, total, err := s.repo.GetByCourseID(ctx, courseID, pagination)
	if err != nil {
		s.logger.Error(ctx, "Failed to get lessons by course ID", "error", err, "course_id", courseID)
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	// Convert to interface slice.
	items := make([]interface{}, len(lessons))
	for i, lesson := range lessons {
		items[i] = lesson
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	response := &domain.ListResponse{
		Items: items,
		Pagination: &domain.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    pagination.Page < totalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	return response, nil
}

// GetAllLessons retrieves all lessons with pagination.
func (s *lessonService) GetAllLessons(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	// Set default pagination.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	lessons, total, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all lessons", "error", err)
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	// Convert to interface slice.
	items := make([]interface{}, len(lessons))
	for i, lesson := range lessons {
		items[i] = lesson
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	response := &domain.ListResponse{
		Items: items,
		Pagination: &domain.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    pagination.Page < totalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	return response, nil
}

// UpdateLesson updates an existing lesson.
func (s *lessonService) UpdateLesson(ctx context.Context, id string, req *domain.UpdateLessonRequest) (*domain.Lesson, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("lesson ID is required")
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid lesson update request", "error", err, "lesson_id", id)
		return nil, errors.NewValidationError("invalid lesson data", err)
	}

	// Get existing lesson.
	lesson, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get lesson for update", "error", err, "lesson_id", id)
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	// Update fields.
	if req.Title != nil {
		lesson.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		lesson.Description = strings.TrimSpace(*req.Description)
	}
	if req.Content != nil {
		lesson.Content = *req.Content
	}
	if req.Order != nil {
		lesson.Order = *req.Order
	}
	lesson.UpdatedAt = time.Now()

	// Save to repository.
	if err := s.repo.Update(ctx, lesson); err != nil {
		s.logger.Error(ctx, "Failed to update lesson", "error", err, "lesson_id", id)
		return nil, fmt.Errorf("failed to update lesson: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("lesson:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate lesson cache", "error", err, "lesson_id", id)
		}
	}

	// Publish lesson updated event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		lessonData := map[string]interface{}{
			"title":       lesson.Title,
			"description": lesson.Description,
			"content":     lesson.Content,
			"order":       lesson.Order,
		}
		if err := s.messaging.PublishLessonUpdated(ctx, lesson.ID, lesson.CourseID, lessonData); err != nil {
			s.logger.Warn(ctx, "Failed to publish lesson updated event", "error", err, "lesson_id", lesson.ID)
		}
	}

	s.logger.Info(ctx, "Lesson updated successfully", "lesson_id", id)

	return lesson, nil
}

// DeleteLesson deletes a lesson.
func (s *lessonService) DeleteLesson(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError("lesson ID is required")
	}

	// Get lesson for event publishing.
	lesson, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get lesson for deletion", "error", err, "lesson_id", id)
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	// Delete from repository.
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete lesson", "error", err, "lesson_id", id)
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("lesson:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate lesson cache", "error", err, "lesson_id", id)
		}
	}

	// Publish lesson deleted event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		if err := s.messaging.PublishLessonDeleted(ctx, lesson.ID, lesson.CourseID); err != nil {
			s.logger.Warn(ctx, "Failed to publish lesson deleted event", "error", err, "lesson_id", lesson.ID)
		}
	}

	s.logger.Info(ctx, "Lesson deleted successfully", "lesson_id", id)

	return nil
}

// Helper function to generate IDs (should be moved to a utility package).
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
