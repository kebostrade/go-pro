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

	"go-pro-backend/internal/auth"
	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

const (
	exerciseCacheKeyFmt           = "exercise:%s"
	exerciseIDRequiredMsg         = "exercise ID is required"
	failedToGetExerciseFmt        = "failed to get exercise: %w"
	failedToUpdateExerciseFmt     = "failed to update exercise: %w"
	failedToDeleteExerciseFmt     = "failed to delete exercise: %w"
	failedToCreateExerciseFmt     = "failed to create exercise: %w"
	invalidDifficultyMsg          = "invalid difficulty level"
	codeTooLongMsg                = "code exceeds maximum length of 50KB"
	failedToExecuteCodeMsg        = "code execution failed"
	failedToGetExerciseForUpdateMsg = "failed to get exercise for update"
	failedToGetExerciseForDeleteMsg = "failed to get exercise for deletion"
	failedToGetExerciseForSubmissionMsg = "failed to get exercise for submission"
)

// exerciseService implements the ExerciseService interface.
type exerciseService struct {
	repo      repository.ExerciseRepository
	executor  ExecutorService
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
		executor:  NewMockExecutorService(), // Phase 2: Replace with real executor
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
		return nil, errors.NewBadRequestError(invalidDifficultyMsg)
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
		return nil, fmt.Errorf(failedToCreateExerciseFmt, err)
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
	s.cacheExercise(ctx, exercise)

	s.logger.Info(ctx, "Exercise created successfully", "exercise_id", exercise.ID, "lesson_id", exercise.LessonID)

	return exercise, nil
}

// GetExerciseByID retrieves an exercise by ID.
func (s *exerciseService) GetExerciseByID(ctx context.Context, id string) (*domain.Exercise, error) {
	if id == "" {
		return nil, errors.NewBadRequestError(exerciseIDRequiredMsg)
	}

	// Try cache first.
	if exercise := s.getExerciseFromCache(ctx, id); exercise != nil {
		s.logger.Debug(ctx, "Exercise retrieved from cache", "exercise_id", id)
		return exercise, nil
	}

	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise", "error", err, "exercise_id", id)
		return nil, fmt.Errorf(failedToGetExerciseFmt, err)
	}

	// Cache the exercise.
	s.cacheExercise(ctx, exercise)

	return exercise, nil
}

// GetExercisesByLessonID retrieves exercises for a specific lesson.
func (s *exerciseService) GetExercisesByLessonID(
	ctx context.Context,
	lessonID string,
	pagination *domain.PaginationRequest,
) (*domain.ListResponse, error) {
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
		return nil, errors.NewBadRequestError(exerciseIDRequiredMsg)
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
		return nil, fmt.Errorf(failedToGetExerciseFmt, err)
	}

	// Apply updates.
	s.applyExerciseUpdates(exercise, req)
	exercise.UpdatedAt = time.Now()

	// Save to repository.
	if err := s.repo.Update(ctx, exercise); err != nil {
		s.logger.Error(ctx, "Failed to update exercise", "error", err, "exercise_id", id)
		return nil, fmt.Errorf(failedToUpdateExerciseFmt, err)
	}

	// Invalidate cache and publish event.
	s.invalidateExerciseCache(ctx, id)
	s.publishExerciseUpdatedEvent(ctx, exercise)

	s.logger.Info(ctx, "Exercise updated successfully", "exercise_id", id)

	return exercise, nil
}

// DeleteExercise deletes an exercise.
func (s *exerciseService) DeleteExercise(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError(exerciseIDRequiredMsg)
	}

	// Get exercise for event publishing.
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise for deletion", "error", err, "exercise_id", id)
		return fmt.Errorf(failedToGetExerciseFmt, err)
	}

	// Delete from repository.
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete exercise", "error", err, "exercise_id", id)
		return fmt.Errorf(failedToDeleteExerciseFmt, err)
	}

	// Cleanup cache and publish event.
	s.invalidateExerciseCache(ctx, id)
	s.publishExerciseDeletedEvent(ctx, exercise)

	s.logger.Info(ctx, "Exercise deleted successfully", "exercise_id", id)

	return nil
}

// SubmitExercise handles exercise submission and evaluation.
func (s *exerciseService) SubmitExercise(
	ctx context.Context,
	exerciseID string,
	req *domain.SubmitExerciseRequest,
) (*domain.ExerciseSubmissionResult, error) {
	if exerciseID == "" {
		return nil, errors.NewBadRequestError(exerciseIDRequiredMsg)
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(ctx, "Invalid exercise submission request", "error", err, "exercise_id", exerciseID)
		return nil, errors.NewValidationError("invalid submission data", err)
	}

	// Validate code length.
	if len(req.Code) > 50000 { // 50KB limit
		return nil, errors.NewBadRequestError(codeTooLongMsg)
	}

	// Get the exercise.
	exercise, err := s.repo.GetByID(ctx, exerciseID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get exercise for submission", "error", err, "exercise_id", exerciseID)
		return nil, fmt.Errorf(failedToGetExerciseFmt, err)
	}

	// Execute code and build result.
	result := s.executeAndEvaluateCode(ctx, exerciseID, req)

	// Publish event.
	s.publishExerciseSubmittedEvent(ctx, exercise, req, result)

	s.logger.Info(ctx, "Exercise submitted successfully",
		"exercise_id", exerciseID,
		"score", result.Score,
		"passed", result.Passed,
		"lesson_id", exercise.LessonID)

	return result, nil
}

// Helper methods to reduce cognitive complexity

func (s *exerciseService) getExerciseFromCache(ctx context.Context, id string) *domain.Exercise {
	if s.cache == nil {
		return nil
	}

	cacheKey := fmt.Sprintf(exerciseCacheKeyFmt, id)
	var cachedExercise domain.Exercise
	if err := s.cache.Get(ctx, cacheKey, &cachedExercise); err == nil {
		return &cachedExercise
	}
	return nil
}

func (s *exerciseService) cacheExercise(ctx context.Context, exercise *domain.Exercise) {
	if s.cache == nil || exercise == nil {
		return
	}

	cacheKey := fmt.Sprintf(exerciseCacheKeyFmt, exercise.ID)
	if err := s.cache.Set(ctx, cacheKey, exercise, 15*time.Minute); err != nil {
		s.logger.Warn(ctx, "Failed to cache exercise", "error", err, "exercise_id", exercise.ID)
	}
}

func (s *exerciseService) invalidateExerciseCache(ctx context.Context, id string) {
	if s.cache == nil {
		return
	}

	cacheKey := fmt.Sprintf(exerciseCacheKeyFmt, id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn(ctx, "Failed to invalidate exercise cache", "error", err, "exercise_id", id)
	}
}

func (s *exerciseService) applyExerciseUpdates(exercise *domain.Exercise, req *domain.UpdateExerciseRequest) {
	if req.Title != nil {
		exercise.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		exercise.Description = strings.TrimSpace(*req.Description)
	}
	if req.TestCases != nil {
		exercise.TestCases = *req.TestCases
	}
	if req.Difficulty != nil && req.Difficulty.IsValid() {
		exercise.Difficulty = *req.Difficulty
	}
}

func (s *exerciseService) publishExerciseUpdatedEvent(ctx context.Context, exercise *domain.Exercise) {
	if s.messaging == nil || !s.messaging.IsEnabled() {
		return
	}

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

func (s *exerciseService) publishExerciseDeletedEvent(ctx context.Context, exercise *domain.Exercise) {
	if s.messaging == nil || !s.messaging.IsEnabled() {
		return
	}

	if err := s.messaging.PublishExerciseDeleted(ctx, exercise.ID, exercise.LessonID); err != nil {
		s.logger.Warn(ctx, "Failed to publish exercise deleted event", "error", err, "exercise_id", exercise.ID)
	}
}

func (s *exerciseService) executeAndEvaluateCode(
	ctx context.Context,
	exerciseID string,
	req *domain.SubmitExerciseRequest,
) *domain.ExerciseSubmissionResult {
	// Phase 1: Mock test cases (Phase 2 will load from database).
	testCases := []TestCase{
		{
			Name:     "Test basic output",
			Input:    "",
			Expected: "Hello, World!",
		},
		{
			Name:     "Test edge case",
			Input:    "edge",
			Expected: "Edge case handled",
		},
	}

	// Execute code against test cases.
	execReq := &ExecuteRequest{
		Code:      req.Code,
		Language:  req.Language,
		TestCases: testCases,
		Timeout:   30 * time.Second,
	}

	execResult, err := s.executor.ExecuteCode(ctx, execReq)
	if err != nil {
		s.logger.Error(ctx, "Failed to execute code", "error", err, "exercise_id", exerciseID)
		return &domain.ExerciseSubmissionResult{
			Success:     false,
			ExerciseID:  exerciseID,
			Message:     failedToExecuteCodeMsg,
			SubmittedAt: time.Now(),
		}
	}

	// Convert executor results to domain results.
	domainResults := make([]domain.TestResult, len(execResult.Results))
	for i, r := range execResult.Results {
		domainResults[i] = domain.TestResult{
			TestName: r.TestName,
			Passed:   r.Passed,
			Expected: r.Expected,
			Actual:   r.Actual,
			Error:    r.Error,
		}
	}

	// Build result message.
	message := "All tests passed!"
	if !execResult.Passed {
		message = fmt.Sprintf("Tests passed: %d/%d", countPassedTests(domainResults), len(domainResults))
	}

	return &domain.ExerciseSubmissionResult{
		Success:         true,
		ExerciseID:      exerciseID,
		Score:           execResult.Score,
		Passed:          execResult.Passed,
		Message:         message,
		TestResults:     domainResults,
		ExecutionTimeMs: execResult.ExecutionTime.Milliseconds(),
		SubmittedAt:     time.Now(),
	}
}

func (s *exerciseService) publishExerciseSubmittedEvent(
	ctx context.Context,
	exercise *domain.Exercise,
	req *domain.SubmitExerciseRequest,
	result *domain.ExerciseSubmissionResult,
) {
	if s.messaging == nil || !s.messaging.IsEnabled() {
		return
	}

	userInfo := auth.GetUserInfo(ctx)
	userID := "anonymous"
	if userInfo != nil {
		userID = userInfo.ID
	}

	submissionData := map[string]interface{}{
		"code":     req.Code,
		"language": req.Language,
		"passed":   result.Passed,
	}

	if err := s.messaging.PublishExerciseSubmitted(ctx, userID, exercise.ID, exercise.LessonID, result.Score, submissionData); err != nil {
		s.logger.Warn(ctx, "Failed to publish exercise submitted event", "error", err, "exercise_id", exercise.ID)
	}
}

// countPassedTests counts the number of passed tests.
func countPassedTests(results []domain.TestResult) int {
	count := 0
	for _, r := range results {
		if r.Passed {
			count++
		}
	}
	return count
}
