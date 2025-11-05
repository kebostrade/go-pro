// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// progressService implements the ProgressService interface.
type progressService struct {
	repo      repository.ProgressRepository
	logger    logger.Logger
	validator validator.Validator
	cache     cache.CacheManager
	messaging *messaging.Service
}

// NewProgressService creates a new progress service.
func NewProgressService(
	repo repository.ProgressRepository,
	config *Config,
) ProgressService {
	return &progressService{
		repo:      repo,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}
}

// CreateProgress creates a new progress record.
func (s *progressService) CreateProgress(ctx context.Context, req *domain.CreateProgressRequest) (*domain.Progress, error) {
	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid progress creation request", "error", err)
		return nil, errors.NewValidationError("invalid progress data", err)
	}

	// Validate status.
	if !req.Status.IsValid() {
		return nil, errors.NewBadRequestError("invalid progress status")
	}

	// Check if progress already exists for this user and lesson.
	existingProgress, err := s.repo.GetByUserAndLesson(ctx, req.UserID, req.LessonID)
	if err == nil && existingProgress != nil {
		return nil, errors.NewConflictError("progress already exists for this user and lesson")
	}

	// Create progress entity.
	progress := &domain.Progress{
		ID:          generateID("progress"),
		UserID:      req.UserID,
		LessonID:    req.LessonID,
		Status:      req.Status,
		CompletedAt: nil, // Will be set when status is completed
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set completion time if status is completed.
	if req.Status == domain.StatusCompleted {
		now := time.Now()
		progress.CompletedAt = &now
	}

	// Save to repository.
	if err := s.repo.Create(ctx, progress); err != nil {
		s.logger.Error(ctx, "Failed to create progress", "error", err, "progress_id", progress.ID)
		return nil, fmt.Errorf("failed to create progress: %w", err)
	}

	// Publish progress created event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		progressData := map[string]interface{}{
			"status":    string(progress.Status),
			"user_id":   progress.UserID,
			"lesson_id": progress.LessonID,
		}
		if err := s.messaging.PublishProgressCreated(ctx, progress.UserID, progress.LessonID, progressData); err != nil {
			s.logger.Warn(ctx, "Failed to publish progress created event", "error", err, "progress_id", progress.ID)
		}
	}

	// Cache the progress.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:%s", progress.ID)
		if err := s.cache.Set(ctx, cacheKey, progress, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache progress", "error", err, "progress_id", progress.ID)
		}

		// Also cache by user and lesson for quick lookup.
		userLessonKey := fmt.Sprintf("progress:user:%s:lesson:%s", progress.UserID, progress.LessonID)
		if err := s.cache.Set(ctx, userLessonKey, progress, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache progress by user-lesson", "error", err, "progress_id", progress.ID)
		}
	}

	s.logger.Info(ctx, "Progress created successfully", "progress_id", progress.ID, "user_id", progress.UserID, "lesson_id", progress.LessonID)

	return progress, nil
}

// GetProgressByID retrieves a progress record by ID.
func (s *progressService) GetProgressByID(ctx context.Context, id string) (*domain.Progress, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("progress ID is required")
	}

	// Try cache first.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:%s", id)
		var cachedProgress domain.Progress
		if err := s.cache.Get(ctx, cacheKey, &cachedProgress); err == nil {
			s.logger.Debug(ctx, "Progress retrieved from cache", "progress_id", id)
			return &cachedProgress, nil
		}
	}

	progress, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get progress", "error", err, "progress_id", id)
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Cache the progress.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:%s", id)
		if err := s.cache.Set(ctx, cacheKey, progress, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache progress", "error", err, "progress_id", id)
		}
	}

	return progress, nil
}

// GetProgressByUserID retrieves progress records for a specific user.
func (s *progressService) GetProgressByUserID(
	ctx context.Context,
	userID string,
	pagination *domain.PaginationRequest,
) (*domain.ListResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequestError("user ID is required")
	}

	// Set default pagination.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	progressRecords, total, err := s.repo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		s.logger.Error(ctx, "Failed to get progress by user ID", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Convert to interface slice.
	items := make([]interface{}, len(progressRecords))
	for i, progress := range progressRecords {
		items[i] = progress
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

// GetProgressByUserAndLesson retrieves progress for a specific user and lesson.
func (s *progressService) GetProgressByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error) {
	if userID == "" {
		return nil, errors.NewBadRequestError("user ID is required")
	}
	if lessonID == "" {
		return nil, errors.NewBadRequestError("lesson ID is required")
	}

	// Try cache first.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:user:%s:lesson:%s", userID, lessonID)
		var cachedProgress domain.Progress
		if err := s.cache.Get(ctx, cacheKey, &cachedProgress); err == nil {
			s.logger.Debug(ctx, "Progress retrieved from cache", "user_id", userID, "lesson_id", lessonID)
			return &cachedProgress, nil
		}
	}

	progress, err := s.repo.GetByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get progress by user and lesson", "error", err, "user_id", userID, "lesson_id", lessonID)
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Cache the progress.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:user:%s:lesson:%s", userID, lessonID)
		if err := s.cache.Set(ctx, cacheKey, progress, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache progress", "error", err, "user_id", userID, "lesson_id", lessonID)
		}
	}

	return progress, nil
}

// UpdateProgress updates an existing progress record by user and lesson.
func (s *progressService) UpdateProgress(ctx context.Context, userID, lessonID string, req *domain.UpdateProgressRequest) (*domain.Progress, error) {
	if userID == "" || lessonID == "" {
		return nil, errors.NewBadRequestError("user ID and lesson ID are required")
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid progress update request", "error", err, "user_id", userID, "lesson_id", lessonID)
		return nil, errors.NewValidationError("invalid progress data", err)
	}

	// Get existing progress by user and lesson.
	progress, err := s.repo.GetByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get progress for update", "error", err, "user_id", userID, "lesson_id", lessonID)
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	oldStatus := progress.Status

	// Update fields.
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, errors.NewBadRequestError("invalid progress status")
		}
		progress.Status = *req.Status

		// Set or clear completion time based on status.
		if *req.Status == domain.StatusCompleted && oldStatus != domain.StatusCompleted {
			now := time.Now()
			progress.CompletedAt = &now
		} else if *req.Status != domain.StatusCompleted {
			progress.CompletedAt = nil
		}
	}
	progress.UpdatedAt = time.Now()

	// Save to repository.
	if err := s.repo.Update(ctx, progress); err != nil {
		s.logger.Error(ctx, "Failed to update progress", "error", err, "progress_id", progress.ID)
		return nil, fmt.Errorf("failed to update progress: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:%s", progress.ID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate progress cache", "error", err, "progress_id", progress.ID)
		}

		// Also invalidate user-lesson cache.
		userLessonKey := fmt.Sprintf("progress:user:%s:lesson:%s", progress.UserID, progress.LessonID)
		if err := s.cache.Delete(ctx, userLessonKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate progress user-lesson cache", "error", err, "progress_id", progress.ID)
		}
	}

	// Publish progress updated event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		completed := progress.Status == domain.StatusCompleted
		score := 0
		if completed {
			score = 100 // Simple scoring for now
		}
		if err := s.messaging.PublishProgressUpdated(ctx, progress.UserID, progress.LessonID, "", completed, score, 0); err != nil {
			s.logger.Warn(ctx, "Failed to publish progress updated event", "error", err, "progress_id", progress.ID)
		}

		// If lesson was completed, publish lesson completed event.
		if progress.Status == domain.StatusCompleted && oldStatus != domain.StatusCompleted {
			if err := s.messaging.PublishLessonCompleted(ctx, progress.UserID, progress.LessonID, ""); err != nil {
				s.logger.Warn(ctx, "Failed to publish lesson completed event", "error", err, "progress_id", progress.ID)
			}
		}
	}

	s.logger.Info(ctx, "Progress updated successfully", "progress_id", progress.ID, "old_status", oldStatus, "new_status", progress.Status)

	return progress, nil
}

// DeleteProgress deletes a progress record.
func (s *progressService) DeleteProgress(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError("progress ID is required")
	}

	// Get progress for event publishing and cache invalidation.
	progress, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get progress for deletion", "error", err, "progress_id", id)
		return fmt.Errorf("failed to get progress: %w", err)
	}

	// Delete from repository.
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete progress", "error", err, "progress_id", id)
		return fmt.Errorf("failed to delete progress: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("progress:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate progress cache", "error", err, "progress_id", id)
		}

		// Also invalidate user-lesson cache.
		userLessonKey := fmt.Sprintf("progress:user:%s:lesson:%s", progress.UserID, progress.LessonID)
		if err := s.cache.Delete(ctx, userLessonKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate progress user-lesson cache", "error", err, "progress_id", id)
		}
	}

	// Note: Progress deletion events are not typically published as they're administrative actions.

	s.logger.Info(ctx, "Progress deleted successfully", "progress_id", id)

	return nil
}
