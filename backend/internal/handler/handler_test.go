// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Services
type mockCourseService struct {
	mock.Mock
}

func (m *mockCourseService) CreateCourse(ctx context.Context, req *domain.CreateCourseRequest) (*domain.Course, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Course), args.Error(1)
}

func (m *mockCourseService) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Course), args.Error(1)
}

func (m *mockCourseService) GetAllCourses(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockCourseService) UpdateCourse(ctx context.Context, id string, req *domain.UpdateCourseRequest) (*domain.Course, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Course), args.Error(1)
}

func (m *mockCourseService) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockLessonService struct {
	mock.Mock
}

func (m *mockLessonService) CreateLesson(ctx context.Context, req *domain.CreateLessonRequest) (*domain.Lesson, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}

func (m *mockLessonService) GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}

func (m *mockLessonService) GetLessonsByCourseID(ctx context.Context, courseID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, courseID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockLessonService) GetAllLessons(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockLessonService) UpdateLesson(ctx context.Context, id string, req *domain.UpdateLessonRequest) (*domain.Lesson, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}

func (m *mockLessonService) DeleteLesson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockExerciseService struct {
	mock.Mock
}

func (m *mockExerciseService) CreateExercise(ctx context.Context, req *domain.CreateExerciseRequest) (*domain.Exercise, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *mockExerciseService) GetExerciseByID(ctx context.Context, id string) (*domain.Exercise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *mockExerciseService) GetExercisesByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, lessonID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockExerciseService) GetAllExercises(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockExerciseService) UpdateExercise(ctx context.Context, id string, req *domain.UpdateExerciseRequest) (*domain.Exercise, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *mockExerciseService) DeleteExercise(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockExerciseService) SubmitExercise(ctx context.Context, exerciseID string, req *domain.SubmitExerciseRequest) (*domain.ExerciseSubmissionResult, error) {
	args := m.Called(ctx, exerciseID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExerciseSubmissionResult), args.Error(1)
}

type mockProgressService struct {
	mock.Mock
}

func (m *mockProgressService) CreateProgress(ctx context.Context, req *domain.CreateProgressRequest) (*domain.Progress, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Progress), args.Error(1)
}

func (m *mockProgressService) GetProgressByID(ctx context.Context, id string) (*domain.Progress, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Progress), args.Error(1)
}

func (m *mockProgressService) GetProgressByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	args := m.Called(ctx, userID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListResponse), args.Error(1)
}

func (m *mockProgressService) GetProgressByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error) {
	args := m.Called(ctx, userID, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Progress), args.Error(1)
}

func (m *mockProgressService) UpdateProgress(ctx context.Context, userID, lessonID string, req *domain.UpdateProgressRequest) (*domain.Progress, error) {
	args := m.Called(ctx, userID, lessonID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Progress), args.Error(1)
}

func (m *mockProgressService) DeleteProgress(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockCurriculumService struct {
	mock.Mock
}

func (m *mockCurriculumService) GetCurriculum(ctx context.Context) (*domain.Curriculum, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Curriculum), args.Error(1)
}

func (m *mockCurriculumService) GetLessonDetail(ctx context.Context, lessonID int) (*domain.LessonDetail, error) {
	args := m.Called(ctx, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LessonDetail), args.Error(1)
}

type mockHealthService struct {
	mock.Mock
}

func (m *mockHealthService) GetHealthStatus(ctx context.Context) (*domain.HealthResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.HealthResponse), args.Error(1)
}

// Test helper functions
func setupTest() (*Handler, *mockCourseService, *mockLessonService, *mockExerciseService, *mockProgressService, *mockCurriculumService, *mockHealthService) {
	courseService := new(mockCourseService)
	lessonService := new(mockLessonService)
	exerciseService := new(mockExerciseService)
	progressService := new(mockProgressService)
	curriculumService := new(mockCurriculumService)
	healthService := new(mockHealthService)

	services := &service.Services{
		Course:     courseService,
		Lesson:     lessonService,
		Exercise:   exerciseService,
		Progress:   progressService,
		Curriculum: curriculumService,
		Health:     healthService,
	}

	log := logger.New("info", "json")
	v := validator.New()

	handler := New(services, log, v)

	return handler, courseService, lessonService, exerciseService, progressService, curriculumService, healthService
}

// Health endpoint tests
func TestHandleHealth(t *testing.T) {
	tests := []struct {
		name           string
		mockHealth     *domain.HealthResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "success",
			mockHealth: &domain.HealthResponse{
				Status:    "healthy",
				Timestamp: time.Now(),
				Version:   "1.0.0",
				Uptime:    "5m",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "internal error",
			mockHealth:     nil,
			mockError:      errors.NewInternalError("database connection failed", nil),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _, _, _, _, _, healthService := setupTest()
			healthService.On("GetHealthStatus", mock.Anything).Return(tt.mockHealth, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			w := httptest.NewRecorder()

			handler.handleHealth(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.mockError == nil {
				var response domain.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
			}

			healthService.AssertExpectations(t)
		})
	}
}

// Course endpoint tests
func TestHandleGetCourse(t *testing.T) {
	tests := []struct {
		name           string
		courseID       string
		mockCourse     *domain.Course
		mockError      error
		expectedStatus int
	}{
		{
			name:     "success",
			courseID: "go-basics",
			mockCourse: &domain.Course{
				ID:          "go-basics",
				Title:       "Go Basics",
				Description: "Learn Go programming fundamentals",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			courseID:       "nonexistent",
			mockCourse:     nil,
			mockError:      errors.NewNotFoundError("Exercise not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, courseService, _, _, _, _, _ := setupTest()
			courseService.On("GetCourseByID", mock.Anything, tt.courseID).Return(tt.mockCourse, tt.mockError)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/v1/courses/{id}", handler.handleGetCourse)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/"+tt.courseID, nil)
			req.SetPathValue("id", tt.courseID)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response domain.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.mockError == nil {
				assert.True(t, response.Success)
			} else {
				assert.False(t, response.Success)
			}

			courseService.AssertExpectations(t)
		})
	}
}

// Exercise submission tests
func TestHandleSubmitExercise(t *testing.T) {
	tests := []struct {
		name           string
		exerciseID     string
		submitReq      domain.SubmitExerciseRequest
		mockResult     *domain.ExerciseSubmissionResult
		mockError      error
		expectedStatus int
	}{
		{
			name:       "success - all tests passed",
			exerciseID: "ex-001",
			submitReq: domain.SubmitExerciseRequest{
				Code:     "package main\nfunc main() {}",
				Language: "go",
			},
			mockResult: &domain.ExerciseSubmissionResult{
				Success:    true,
				ExerciseID: "ex-001",
				Score:      100,
				Passed:     true,
				Message:    "All tests passed!",
				TestResults: []domain.TestResult{
					{TestName: "Test 1", Passed: true, Expected: "output", Actual: "output"},
				},
				ExecutionTimeMs: 150,
				SubmittedAt:     time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "success - some tests failed",
			exerciseID: "ex-001",
			submitReq: domain.SubmitExerciseRequest{
				Code:     "package main\nfunc main() {}",
				Language: "go",
			},
			mockResult: &domain.ExerciseSubmissionResult{
				Success:    true,
				ExerciseID: "ex-001",
				Score:      50,
				Passed:     false,
				Message:    "Tests passed: 1/2",
				TestResults: []domain.TestResult{
					{TestName: "Test 1", Passed: true, Expected: "output", Actual: "output"},
					{TestName: "Test 2", Passed: false, Expected: "expected", Actual: "actual"},
				},
				ExecutionTimeMs: 200,
				SubmittedAt:     time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "invalid language",
			exerciseID: "ex-001",
			submitReq: domain.SubmitExerciseRequest{
				Code:     "print('hello')",
				Language: "ruby",
			},
			mockResult:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _, _, exerciseService, _, _, _ := setupTest()

			if tt.mockResult != nil || tt.mockError != nil {
				exerciseService.On("SubmitExercise", mock.Anything, tt.exerciseID, mock.AnythingOfType("*domain.SubmitExerciseRequest")).Return(tt.mockResult, tt.mockError)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /api/v1/exercises/{id}/submit", handler.handleSubmitExercise)

			body, _ := json.Marshal(tt.submitReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises/"+tt.exerciseID+"/submit", bytes.NewBuffer(body))
			req.SetPathValue("id", tt.exerciseID)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response domain.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response.Success)
			}
		})
	}
}

// Rate limiting tests
func TestSubmissionRateLimit(t *testing.T) {
	handler, _, _, exerciseService, _, _, _ := setupTest()

	mockResult := &domain.ExerciseSubmissionResult{
		Success:         true,
		ExerciseID:      "ex-001",
		Score:           100,
		Passed:          true,
		Message:         "All tests passed!",
		TestResults:     []domain.TestResult{},
		ExecutionTimeMs: 150,
		SubmittedAt:     time.Now(),
	}

	exerciseService.On("SubmitExercise", mock.Anything, "ex-001", mock.AnythingOfType("*domain.SubmitExerciseRequest")).Return(mockResult, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/exercises/{id}/submit", handler.handleSubmitExercise)

	submitReq := domain.SubmitExerciseRequest{
		Code:     "package main\nfunc main() {}",
		Language: "go",
	}

	// Submit 10 times - should all succeed
	for i := 0; i < 10; i++ {
		body, _ := json.Marshal(submitReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises/ex-001/submit", bytes.NewBuffer(body))
		req.SetPathValue("id", "ex-001")
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Submission %d should succeed", i+1)
	}

	// 11th submission should be rate limited
	body, _ := json.Marshal(submitReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises/ex-001/submit", bytes.NewBuffer(body))
	req.SetPathValue("id", "ex-001")
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var response domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error.Message, "rate limit")
}

// Progress endpoint tests
func TestHandleUpdateProgress(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		lessonID       string
		updateReq      domain.UpdateProgressRequest
		mockProgress   *domain.Progress
		mockError      error
		expectedStatus int
	}{
		{
			name:     "success",
			userID:   "user-123",
			lessonID: "lesson-001",
			updateReq: domain.UpdateProgressRequest{
				Status: func() *domain.Status { s := domain.StatusCompleted; return &s }(),
			},
			mockProgress: &domain.Progress{
				ID:        "progress-001",
				UserID:    "user-123",
				LessonID:  "lesson-001",
				Status:    domain.StatusCompleted,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:     "missing user ID",
			userID:   "",
			lessonID: "lesson-001",
			updateReq: domain.UpdateProgressRequest{
				Status: func() *domain.Status { s := domain.StatusCompleted; return &s }(),
			},
			mockProgress:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _, _, _, progressService, _, _ := setupTest()

			if tt.userID != "" && tt.lessonID != "" {
				progressService.On("UpdateProgress", mock.Anything, tt.userID, tt.lessonID, mock.AnythingOfType("*domain.UpdateProgressRequest")).Return(tt.mockProgress, tt.mockError)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /api/v1/progress/{userId}/lesson/{lessonId}", handler.handleUpdateProgress)

			body, _ := json.Marshal(tt.updateReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+tt.userID+"/lesson/"+tt.lessonID, bytes.NewBuffer(body))
			req.SetPathValue("userId", tt.userID)
			req.SetPathValue("lessonId", tt.lessonID)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Error handling tests
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockCourseService)
		expectedStatus int
		expectedType   string
	}{
		{
			name: "not found error",
			setupMock: func(m *mockCourseService) {
				m.On("GetCourseByID", mock.Anything, "nonexistent").Return(nil, errors.NewNotFoundError("User progress not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedType:   "NOT_FOUND",
		},
		{
			name: "bad request error",
			setupMock: func(m *mockCourseService) {
				m.On("GetCourseByID", mock.Anything, "invalid").Return(nil, errors.NewBadRequestError("invalid ID format"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedType:   "BAD_REQUEST",
		},
		{
			name: "internal error",
			setupMock: func(m *mockCourseService) {
				m.On("GetCourseByID", mock.Anything, "error").Return(nil, errors.NewInternalError("database error", nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedType:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, courseService, _, _, _, _, _ := setupTest()
			tt.setupMock(courseService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/v1/courses/{id}", handler.handleGetCourse)

			courseID := "nonexistent"
			if tt.name == "bad request error" {
				courseID = "invalid"
			} else if tt.name == "internal error" {
				courseID = "error"
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/"+courseID, nil)
			req.SetPathValue("id", courseID)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response domain.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.NotNil(t, response.Error)
			assert.Equal(t, tt.expectedType, response.Error.Type)

			courseService.AssertExpectations(t)
		})
	}
}

// Curriculum endpoint tests
func TestHandleGetCurriculum(t *testing.T) {
	handler, _, _, _, _, curriculumService, _ := setupTest()

	mockCurriculum := &domain.Curriculum{
		ID:          "go-pro-curriculum",
		Title:       "GO-PRO Learning Path",
		Description: "Comprehensive Go programming curriculum",
		Duration:    "12 weeks",
		Phases:      []domain.CurriculumPhase{},
		Projects:    []domain.Project{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	curriculumService.On("GetCurriculum", mock.Anything).Return(mockCurriculum, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/curriculum", nil)
	w := httptest.NewRecorder()

	handler.handleGetCurriculum(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	curriculumService.AssertExpectations(t)
}

func TestHandleGetLessonDetail(t *testing.T) {
	handler, _, _, _, _, curriculumService, _ := setupTest()

	mockLesson := &domain.LessonDetail{
		ID:          1,
		Title:       "Variables and Types",
		Description: "Learn about Go variables and types",
		Duration:    "30 minutes",
		Difficulty:  domain.DifficultyBeginner,
		Phase:       "Fundamentals",
		Objectives:  []string{"Understand variables", "Learn basic types"},
		Theory:      "Variables in Go...",
		CodeExample: "var x int = 10",
		Solution:    "Complete solution...",
		Exercises:   []domain.LessonExercise{},
	}

	curriculumService.On("GetLessonDetail", mock.Anything, 1).Return(mockLesson, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/curriculum/lesson/{id}", handler.handleGetLessonDetail)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/curriculum/lesson/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	curriculumService.AssertExpectations(t)
}
