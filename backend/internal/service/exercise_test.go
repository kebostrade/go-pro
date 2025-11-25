// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package service

import (
	"context"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockExerciseRepository struct {
	mock.Mock
}

func (m *mockExerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *mockExerciseRepository) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *mockExerciseRepository) GetByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error) {
	args := m.Called(ctx, lessonID, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Exercise), args.Get(1).(int64), args.Error(2)
}

func (m *mockExerciseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Exercise), args.Get(1).(int64), args.Error(2)
}

func (m *mockExerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *mockExerciseRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock executor
type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) ExecuteCode(ctx context.Context, req *ExecuteRequest) (*ExecuteResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExecuteResult), args.Error(1)
}

// Test setup
func setupExerciseTest() (*exerciseService, *mockExerciseRepository, *mockExecutor) {
	repo := new(mockExerciseRepository)
	executor := new(mockExecutor)

	config := &Config{
		Logger:    logger.New("error", "text"),
		Validator: validator.New(),
		Cache:     nil,
		Messaging: nil,
	}

	service := &exerciseService{
		repo:      repo,
		executor:  executor,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}

	return service, repo, executor
}

// CreateExercise tests
func TestCreateExercise(t *testing.T) {
	tests := []struct {
		name      string
		request   *domain.CreateExerciseRequest
		mockError error
		wantError bool
	}{
		{
			name: "success",
			request: &domain.CreateExerciseRequest{
				LessonID:    "lesson-001",
				Title:       "Hello World Exercise",
				Description: "Write a program that prints Hello, World!",
				TestCases:   2,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "invalid difficulty",
			request: &domain.CreateExerciseRequest{
				LessonID:    "lesson-001",
				Title:       "Test Exercise",
				Description: "Test description",
				TestCases:   1,
				Difficulty:  "invalid",
			},
			mockError: nil,
			wantError: true,
		},
		{
			name: "repository error",
			request: &domain.CreateExerciseRequest{
				LessonID:    "lesson-001",
				Title:       "Test Exercise",
				Description: "Test description",
				TestCases:   1,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockError: assert.AnError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _ := setupExerciseTest()

			if !tt.wantError || tt.mockError != nil {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Exercise")).Return(tt.mockError)
			}

			result, err := service.CreateExercise(context.Background(), tt.request)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Title, result.Title)
				assert.Equal(t, tt.request.Description, result.Description)
				assert.Equal(t, tt.request.Difficulty, result.Difficulty)
			}

			if !tt.wantError || tt.mockError != nil {
				repo.AssertExpectations(t)
			}
		})
	}
}

// GetExerciseByID tests
func TestGetExerciseByID(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		mockExercise *domain.Exercise
		mockError    error
		wantError    bool
	}{
		{
			name:       "success",
			exerciseID: "exercise-001",
			mockExercise: &domain.Exercise{
				ID:          "exercise-001",
				LessonID:    "lesson-001",
				Title:       "Test Exercise",
				Description: "Test description",
				TestCases:   2,
				Difficulty:  domain.DifficultyBeginner,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:         "empty ID",
			exerciseID:   "",
			mockExercise: nil,
			mockError:    nil,
			wantError:    true,
		},
		{
			name:         "not found",
			exerciseID:   "nonexistent",
			mockExercise: nil,
			mockError:    errors.NewNotFoundError("exercise not found"),
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _ := setupExerciseTest()

			if tt.exerciseID != "" {
				repo.On("GetByID", mock.Anything, tt.exerciseID).Return(tt.mockExercise, tt.mockError)
			}

			result, err := service.GetExerciseByID(context.Background(), tt.exerciseID)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockExercise.ID, result.ID)
			}

			if tt.exerciseID != "" {
				repo.AssertExpectations(t)
			}
		})
	}
}

// SubmitExercise tests
func TestSubmitExercise(t *testing.T) {
	tests := []struct {
		name          string
		exerciseID    string
		request       *domain.SubmitExerciseRequest
		mockExercise  *domain.Exercise
		mockExecResult *ExecuteResult
		mockExecError error
		wantError     bool
		expectedScore int
		expectedPass  bool
	}{
		{
			name:       "success - all tests passed",
			exerciseID: "exercise-001",
			request: &domain.SubmitExerciseRequest{
				Code:     "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"Hello, World!\") }",
				Language: "go",
			},
			mockExercise: &domain.Exercise{
				ID:          "exercise-001",
				LessonID:    "lesson-001",
				Title:       "Hello World",
				Description: "Print Hello, World!",
				TestCases:   2,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockExecResult: &ExecuteResult{
				Passed:        true,
				Score:         100,
				ExecutionTime: 150 * time.Millisecond,
				Results: []TestResult{
					{TestName: "Test 1", Passed: true, Expected: "Hello, World!", Actual: "Hello, World!"},
					{TestName: "Test 2", Passed: true, Expected: "Edge case handled", Actual: "Edge case handled"},
				},
			},
			mockExecError: nil,
			wantError:     false,
			expectedScore: 100,
			expectedPass:  true,
		},
		{
			name:       "success - some tests failed",
			exerciseID: "exercise-001",
			request: &domain.SubmitExerciseRequest{
				Code:     "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"Wrong output\") }",
				Language: "go",
			},
			mockExercise: &domain.Exercise{
				ID:          "exercise-001",
				LessonID:    "lesson-001",
				Title:       "Hello World",
				Description: "Print Hello, World!",
				TestCases:   2,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockExecResult: &ExecuteResult{
				Passed:        false,
				Score:         50,
				ExecutionTime: 200 * time.Millisecond,
				Results: []TestResult{
					{TestName: "Test 1", Passed: false, Expected: "Hello, World!", Actual: "Wrong output"},
					{TestName: "Test 2", Passed: true, Expected: "Edge case handled", Actual: "Edge case handled"},
				},
			},
			mockExecError: nil,
			wantError:     false,
			expectedScore: 50,
			expectedPass:  false,
		},
		{
			name:       "code too long",
			exerciseID: "exercise-001",
			request: &domain.SubmitExerciseRequest{
				Code:     string(make([]byte, 60000)),
				Language: "go",
			},
			mockExercise:   nil,
			mockExecResult: nil,
			mockExecError:  nil,
			wantError:      true,
		},
		{
			name:       "exercise not found",
			exerciseID: "nonexistent",
			request: &domain.SubmitExerciseRequest{
				Code:     "package main\nfunc main() {}",
				Language: "go",
			},
			mockExercise:   nil,
			mockExecResult: nil,
			mockExecError:  nil,
			wantError:      true,
		},
		{
			name:       "execution error",
			exerciseID: "exercise-001",
			request: &domain.SubmitExerciseRequest{
				Code:     "package main\nfunc main() { panic(\"error\") }",
				Language: "go",
			},
			mockExercise: &domain.Exercise{
				ID:          "exercise-001",
				LessonID:    "lesson-001",
				Title:       "Test",
				Description: "Test",
				TestCases:   1,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockExecResult: nil,
			mockExecError:  assert.AnError,
			wantError:      false, // Service returns result with error message, not error
			expectedScore:  0,
			expectedPass:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, executor := setupExerciseTest()

			if tt.name == "code too long" {
				// No mocks needed - validation happens before repo call
			} else if tt.name == "exercise not found" {
				repo.On("GetByID", mock.Anything, tt.exerciseID).Return(nil, errors.NewNotFoundError("exercise not found"))
			} else {
				repo.On("GetByID", mock.Anything, tt.exerciseID).Return(tt.mockExercise, nil)
				if tt.mockExercise != nil {
					executor.On("ExecuteCode", mock.Anything, mock.AnythingOfType("*service.ExecuteRequest")).Return(tt.mockExecResult, tt.mockExecError)
				}
			}

			result, err := service.SubmitExercise(context.Background(), tt.exerciseID, tt.request)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedScore, result.Score)
				assert.Equal(t, tt.expectedPass, result.Passed)
				assert.Equal(t, tt.exerciseID, result.ExerciseID)
				// For execution errors, Success should be false
				if tt.mockExecError != nil {
					assert.False(t, result.Success)
				} else {
					assert.True(t, result.Success)
				}
			}
		})
	}
}

// UpdateExercise tests
func TestUpdateExercise(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		request      *domain.UpdateExerciseRequest
		mockExercise *domain.Exercise
		mockGetError error
		mockUpdError error
		wantError    bool
	}{
		{
			name:       "success - update title",
			exerciseID: "exercise-001",
			request: &domain.UpdateExerciseRequest{
				Title: func() *string { s := "Updated Title"; return &s }(),
			},
			mockExercise: &domain.Exercise{
				ID:          "exercise-001",
				LessonID:    "lesson-001",
				Title:       "Original Title",
				Description: "Description",
				TestCases:   2,
				Difficulty:  domain.DifficultyBeginner,
			},
			mockGetError: nil,
			mockUpdError: nil,
			wantError:    false,
		},
		{
			name:       "exercise not found",
			exerciseID: "nonexistent",
			request: &domain.UpdateExerciseRequest{
				Title: func() *string { s := "Title"; return &s }(),
			},
			mockExercise: nil,
			mockGetError: errors.NewNotFoundError("exercise not found"),
			mockUpdError: nil,
			wantError:    true,
		},
		{
			name:       "repository update error",
			exerciseID: "exercise-001",
			request: &domain.UpdateExerciseRequest{
				Title: func() *string { s := "Title"; return &s }(),
			},
			mockExercise: &domain.Exercise{
				ID:       "exercise-001",
				LessonID: "lesson-001",
			},
			mockGetError: nil,
			mockUpdError: assert.AnError,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _ := setupExerciseTest()

			repo.On("GetByID", mock.Anything, tt.exerciseID).Return(tt.mockExercise, tt.mockGetError)
			if tt.mockExercise != nil {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Exercise")).Return(tt.mockUpdError)
			}

			result, err := service.UpdateExercise(context.Background(), tt.exerciseID, tt.request)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.request.Title != nil {
					assert.Equal(t, *tt.request.Title, result.Title)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

// DeleteExercise tests
func TestDeleteExercise(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		mockExercise *domain.Exercise
		mockGetError error
		mockDelError error
		wantError    bool
	}{
		{
			name:       "success",
			exerciseID: "exercise-001",
			mockExercise: &domain.Exercise{
				ID:       "exercise-001",
				LessonID: "lesson-001",
			},
			mockGetError: nil,
			mockDelError: nil,
			wantError:    false,
		},
		{
			name:         "empty ID",
			exerciseID:   "",
			mockExercise: nil,
			mockGetError: nil,
			mockDelError: nil,
			wantError:    true,
		},
		{
			name:         "exercise not found",
			exerciseID:   "nonexistent",
			mockExercise: nil,
			mockGetError: errors.NewNotFoundError("exercise not found"),
			mockDelError: nil,
			wantError:    true,
		},
		{
			name:       "delete error",
			exerciseID: "exercise-001",
			mockExercise: &domain.Exercise{
				ID:       "exercise-001",
				LessonID: "lesson-001",
			},
			mockGetError: nil,
			mockDelError: assert.AnError,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _ := setupExerciseTest()

			if tt.exerciseID != "" {
				repo.On("GetByID", mock.Anything, tt.exerciseID).Return(tt.mockExercise, tt.mockGetError)
				if tt.mockExercise != nil {
					repo.On("Delete", mock.Anything, tt.exerciseID).Return(tt.mockDelError)
				}
			}

			err := service.DeleteExercise(context.Background(), tt.exerciseID)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.exerciseID != "" {
				repo.AssertExpectations(t)
			}
		})
	}
}

// Score calculation tests
func TestCountPassedTests(t *testing.T) {
	tests := []struct {
		name     string
		results  []domain.TestResult
		expected int
	}{
		{
			name: "all passed",
			results: []domain.TestResult{
				{Passed: true},
				{Passed: true},
				{Passed: true},
			},
			expected: 3,
		},
		{
			name: "some passed",
			results: []domain.TestResult{
				{Passed: true},
				{Passed: false},
				{Passed: true},
			},
			expected: 2,
		},
		{
			name: "none passed",
			results: []domain.TestResult{
				{Passed: false},
				{Passed: false},
			},
			expected: 0,
		},
		{
			name:     "empty results",
			results:  []domain.TestResult{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countPassedTests(tt.results)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Validation tests
func TestSubmitExerciseValidation(t *testing.T) {
	tests := []struct {
		name      string
		request   *domain.SubmitExerciseRequest
		wantError bool
	}{
		{
			name: "valid Go code",
			request: &domain.SubmitExerciseRequest{
				Code:     "package main\nfunc main() {}",
				Language: "go",
			},
			wantError: false,
		},
		{
			name: "valid Python code",
			request: &domain.SubmitExerciseRequest{
				Code:     "print('hello')",
				Language: "python",
			},
			wantError: false,
		},
		{
			name: "empty code",
			request: &domain.SubmitExerciseRequest{
				Code:     "",
				Language: "go",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			err := v.Validate(tt.request)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Cache integration tests
func TestExerciseServiceWithCache(t *testing.T) {
	repo := new(mockExerciseRepository)
	executor := new(mockExecutor)

	// Note: Using nil cache for now - Redis cache requires connection
	config := &Config{
		Logger:    logger.New("error", "text"),
		Validator: validator.New(),
		Cache:     nil, // No in-memory cache available, would need Redis
		Messaging: nil,
	}

	service := &exerciseService{
		repo:      repo,
		executor:  executor,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}

	mockExercise := &domain.Exercise{
		ID:          "exercise-001",
		LessonID:    "lesson-001",
		Title:       "Test Exercise",
		Description: "Test",
		TestCases:   1,
		Difficulty:  domain.DifficultyBeginner,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Note: Without actual cache (nil cache), both calls will hit repository
	// In production with Redis, the second call would hit cache
	repo.On("GetByID", mock.Anything, "exercise-001").Return(mockExercise, nil).Times(2)

	result1, err := service.GetExerciseByID(context.Background(), "exercise-001")
	assert.NoError(t, err)
	assert.NotNil(t, result1)

	// Second call - would hit cache if Redis was configured
	result2, err := service.GetExerciseByID(context.Background(), "exercise-001")
	assert.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, result1.ID, result2.ID)

	repo.AssertExpectations(t)
}

// Messaging integration tests
// Note: Skipped - requires Kafka infrastructure
func TestExerciseServiceWithMessaging(t *testing.T) {
	t.Skip("Skipping messaging integration test - requires Kafka infrastructure")

	repo := new(mockExerciseRepository)
	executor := new(mockExecutor)

	messagingService, err := messaging.NewService(&messaging.Config{
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create messaging service: %v", err)
	}

	config := &Config{
		Logger:    logger.New("error", "text"),
		Validator: validator.New(),
		Cache:     nil,
		Messaging: messagingService,
	}

	service := &exerciseService{
		repo:      repo,
		executor:  executor,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}

	request := &domain.CreateExerciseRequest{
		LessonID:    "lesson-001",
		Title:       "Test Exercise",
		Description: "Test description",
		TestCases:   1,
		Difficulty:  domain.DifficultyBeginner,
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Exercise")).Return(nil)

	result, err := service.CreateExercise(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	repo.AssertExpectations(t)
}
