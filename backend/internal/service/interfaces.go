// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// CourseService defines business logic for course operations.
type CourseService interface {
	CreateCourse(ctx context.Context, req *domain.CreateCourseRequest) (*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetAllCourses(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	UpdateCourse(ctx context.Context, id string, req *domain.UpdateCourseRequest) (*domain.Course, error)
	DeleteCourse(ctx context.Context, id string) error
}

// LessonService defines business logic for lesson operations.
type LessonService interface {
	CreateLesson(ctx context.Context, req *domain.CreateLessonRequest) (*domain.Lesson, error)
	GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error)
	GetLessonsByCourseID(ctx context.Context, courseID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	GetAllLessons(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	UpdateLesson(ctx context.Context, id string, req *domain.UpdateLessonRequest) (*domain.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
}

// ExerciseService defines business logic for exercise operations.
type ExerciseService interface {
	CreateExercise(ctx context.Context, req *domain.CreateExerciseRequest) (*domain.Exercise, error)
	GetExerciseByID(ctx context.Context, id string) (*domain.Exercise, error)
	GetExercisesByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	GetAllExercises(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	UpdateExercise(ctx context.Context, id string, req *domain.UpdateExerciseRequest) (*domain.Exercise, error)
	DeleteExercise(ctx context.Context, id string) error
	SubmitExercise(ctx context.Context, exerciseID string, req *domain.SubmitExerciseRequest) (*domain.ExerciseSubmissionResult, error)
}

// ProgressService defines business logic for progress tracking.
type ProgressService interface {
	CreateProgress(ctx context.Context, req *domain.CreateProgressRequest) (*domain.Progress, error)
	GetProgressByID(ctx context.Context, id string) (*domain.Progress, error)
	GetProgressByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
	GetProgressByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error)
	UpdateProgress(ctx context.Context, userID, lessonID string, req *domain.UpdateProgressRequest) (*domain.Progress, error)
	DeleteProgress(ctx context.Context, id string) error
}

// CurriculumService defines business logic for curriculum operations.
type CurriculumService interface {
	GetCurriculum(ctx context.Context) (*domain.Curriculum, error)
	GetLessonDetail(ctx context.Context, lessonID int) (*domain.LessonDetail, error)
}

// HealthService defines health check operations.
type HealthService interface {
	GetHealthStatus(ctx context.Context) (*domain.HealthResponse, error)
}

// Config holds configuration for service layer.
type Config struct {
	Logger    logger.Logger
	Validator validator.Validator
	Cache     cache.CacheManager
	Messaging *messaging.Service
}

// Services aggregates all service interfaces.
type Services struct {
	Course     CourseService
	Lesson     LessonService
	Exercise   ExerciseService
	Progress   ProgressService
	Curriculum CurriculumService
	Health     HealthService
	Executor   ExecutorService
	Auth       AuthService
	User       UserService
}

// NewServices creates a new Services instance with all dependencies.
func NewServices(repos *repository.Repositories, config *Config) (*Services, error) {
	if repos == nil {
		return nil, fmt.Errorf("repositories cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Services{
		Course:     NewCourseService(repos.Course, config),
		Lesson:     NewLessonService(repos.Lesson, config),
		Exercise:   NewExerciseService(repos.Exercise, config),
		Progress:   NewProgressService(repos.Progress, config),
		Curriculum: NewCurriculumService(config),
		Health:     NewHealthService("1.0.0"), // TODO: Get version from config
		Executor:   NewMockExecutorService(), // Mock executor (replaced at container level with Docker)
		Auth:       NewAuthService(repos.User, config),
		User:       NewUserService(repos.User, config),
	}, nil
}
