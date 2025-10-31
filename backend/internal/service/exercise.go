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

// exerciseService implements the ExerciseService interface.
type exerciseService struct {
	repo      repository.ExerciseRepository
	logger    logger.Logger
	validator validator.Validator
	cache     cache.CacheManager
	messaging *messaging.Service
}

// NewExerciseService creates a new exercise service.
func NewExerciseService(
	repo repository.ExerciseRepository,
	config *Config,
) ExerciseService {
	return &exerciseService{
		repo:      repo,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}
}

// CreateExercise creates a new exercise.
func (s *exerciseService) CreateExercise(ctx context.Context, req *domain.CreateExerciseRequest) (*domain.Exercise, error) {
	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid exercise creation request", "error", err)
		return nil, errors.NewValidationError("invalid exercise data", err)
	}

	// Validate difficulty.
	if !req.Difficulty.IsValid() {
		return nil, errors.NewBadRequestError("invalid difficulty level")
	}

	// Create exercise entity.
	exercise := &domain.Exercise{
		ID:          generateID("exercise"),
		LessonID:    req.LessonID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		TestCases:   req.TestCases,
		Difficulty:  req.Difficulty,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository.
	if err := s.repo.Create(ctx, exercise); err != nil {
		s.logger.Error(ctx, "Failed to create exercise", "error", err, "exercise_id", exercise.ID)
		return nil, fmt.Errorf("failed to create exercise: %w", err)
	}

	// Publish exercise created event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		exerciseData := map[string]interface{}{
			"title":       exercise.Title,
			"description": exercise.Description,
			"test_cases":  exercise.TestCases,
			"difficulty":  exercise.Difficulty,
		}
		if err := s.messaging.PublishExerciseCreated(ctx, exercise.ID, exercise.LessonID, exerciseData); err != nil {
			s.logger.Warn(ctx, "Failed to publish exercise created event", "error", err, "exercise_id", exercise.ID)
		}
	}

	// Cache the exercise.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("exercise:%s", exercise.ID)
		if err := s.cache.Set(ctx, cacheKey, exercise, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache exercise", "error", err, "exercise_id", exercise.ID)
		}
	}

	s.logger.Info(ctx, "Exercise created successfully", "exercise_id", exercise.ID, "lesson_id", exercise.LessonID)

	return exercise, nil
}

// GetExerciseByID retrieves an exercise by ID.
func (s *exerciseService) GetExerciseByID(ctx context.Context, id string) (*domain.Exercise, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("exercise ID is required")
	}

	// Try cache first.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("exercise:%s", id)
		var cachedExercise domain.Exercise
		if err := s.cache.Get(ctx, cacheKey, &cachedExercise); err == nil {
			s.logger.Debug(ctx, "Exercise retrieved from cache", "exercise_id", id)
			return &cachedExercise, nil
		}
	}

	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise", "error", err, "exercise_id", id)
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// Cache the exercise.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("exercise:%s", id)
		if err := s.cache.Set(ctx, cacheKey, exercise, 15*time.Minute); err != nil {
			s.logger.Warn(ctx, "Failed to cache exercise", "error", err, "exercise_id", id)
		}
	}

	return exercise, nil
}

// GetExercisesByLessonID retrieves exercises for a specific lesson.
func (s *exerciseService) GetExercisesByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	if lessonID == "" {
		return nil, errors.NewBadRequestError("lesson ID is required")
	}

	// Set default pagination.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	exercises, total, err := s.repo.GetByLessonID(ctx, lessonID, pagination)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercises by lesson ID", "error", err, "lesson_id", lessonID)
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	// Convert to interface slice.
	items := make([]interface{}, len(exercises))
	for i, exercise := range exercises {
		items[i] = exercise
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

// GetAllExercises retrieves all exercises with pagination.
func (s *exerciseService) GetAllExercises(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	// Set default pagination.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	exercises, total, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all exercises", "error", err)
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	// Convert to interface slice.
	items := make([]interface{}, len(exercises))
	for i, exercise := range exercises {
		items[i] = exercise
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

// UpdateExercise updates an existing exercise.
func (s *exerciseService) UpdateExercise(ctx context.Context, id string, req *domain.UpdateExerciseRequest) (*domain.Exercise, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("exercise ID is required")
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid exercise update request", "error", err, "exercise_id", id)
		return nil, errors.NewValidationError("invalid exercise data", err)
	}

	// Get existing exercise.
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise for update", "error", err, "exercise_id", id)
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// Update fields.
	if req.Title != nil {
		exercise.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		exercise.Description = strings.TrimSpace(*req.Description)
	}
	if req.TestCases != nil {
		exercise.TestCases = *req.TestCases
	}
	if req.Difficulty != nil {
		if !req.Difficulty.IsValid() {
			return nil, errors.NewBadRequestError("invalid difficulty level")
		}
		exercise.Difficulty = *req.Difficulty
	}
	exercise.UpdatedAt = time.Now()

	// Save to repository.
	if err := s.repo.Update(ctx, exercise); err != nil {
		s.logger.Error(ctx, "Failed to update exercise", "error", err, "exercise_id", id)
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("exercise:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate exercise cache", "error", err, "exercise_id", id)
		}
	}

	// Publish exercise updated event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		exerciseData := map[string]interface{}{
			"title":       exercise.Title,
			"description": exercise.Description,
			"test_cases":  exercise.TestCases,
			"difficulty":  exercise.Difficulty,
		}
		if err := s.messaging.PublishExerciseUpdated(ctx, exercise.ID, exercise.LessonID, exerciseData); err != nil {
			s.logger.Warn(ctx, "Failed to publish exercise updated event", "error", err, "exercise_id", exercise.ID)
		}
	}

	s.logger.Info(ctx, "Exercise updated successfully", "exercise_id", id)

	return exercise, nil
}

// DeleteExercise deletes an exercise.
func (s *exerciseService) DeleteExercise(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError("exercise ID is required")
	}

	// Get exercise for event publishing.
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise for deletion", "error", err, "exercise_id", id)
		return fmt.Errorf("failed to get exercise: %w", err)
	}

	// Delete from repository.
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete exercise", "error", err, "exercise_id", id)
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	// Invalidate cache.
	if s.cache != nil {
		cacheKey := fmt.Sprintf("exercise:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn(ctx, "Failed to invalidate exercise cache", "error", err, "exercise_id", id)
		}
	}

	// Publish exercise deleted event.
	if s.messaging != nil && s.messaging.IsEnabled() {
		if err := s.messaging.PublishExerciseDeleted(ctx, exercise.ID, exercise.LessonID); err != nil {
			s.logger.Warn(ctx, "Failed to publish exercise deleted event", "error", err, "exercise_id", exercise.ID)
		}
	}

	s.logger.Info(ctx, "Exercise deleted successfully", "exercise_id", id)

	return nil
}

// SubmitExercise handles exercise submission and evaluation.
func (s *exerciseService) SubmitExercise(ctx context.Context, exerciseID string, req *domain.SubmitExerciseRequest) (*domain.ExerciseSubmissionResult, error) {
	if exerciseID == "" {
		return nil, errors.NewBadRequestError("exercise ID is required")
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid exercise submission request", "error", err, "exercise_id", exerciseID)
		return nil, errors.NewValidationError("invalid submission data", err)
	}

	// Get the exercise to validate it exists.
	exercise, err := s.repo.GetByID(ctx, exerciseID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise for submission", "error", err, "exercise_id", exerciseID)
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// For now, implement a simple submission result.
	// In a real implementation, this would run the code against test cases.
	result := &domain.ExerciseSubmissionResult{
		ExerciseID: exerciseID,
		Score:      85, // Mock score
		Passed:     true,
		Message:    "Good solution! Consider optimizing for edge cases.",
		TestResults: []domain.TestResult{
			{
				Name:     "Basic functionality",
				Passed:   true,
				Expected: "expected output",
				Actual:   "expected output",
			},
			{
				Name:     "Edge case handling",
				Passed:   false,
				Expected: "edge case output",
				Actual:   "different output",
				Error:    "Edge case not handled properly",
			},
		},
		SubmittedAt: time.Now(),
	}

	s.logger.Info(ctx, "Exercise submitted successfully",
		"exercise_id", exerciseID,
		"score", result.Score,
		"lesson_id", exercise.LessonID)

	return result, nil
}
