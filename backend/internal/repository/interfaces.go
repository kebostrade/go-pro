// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package repository provides functionality for the GO-PRO Learning Platform.
package repository

import (
	"context"
	"time"

	"go-pro-backend/internal/domain"
)

// CourseRepository defines the interface for course data operations.
type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetByID(ctx context.Context, id string) (*domain.Course, error)
	GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error)
	Update(ctx context.Context, course *domain.Course) error
	Delete(ctx context.Context, id string) error
}

// LessonRepository defines the interface for lesson data operations.
type LessonRepository interface {
	Create(ctx context.Context, lesson *domain.Lesson) error
	GetByID(ctx context.Context, id string) (*domain.Lesson, error)
	GetByCourseID(ctx context.Context, courseID string, pagination *domain.PaginationRequest) ([]*domain.Lesson, int64, error)
	GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Lesson, int64, error)
	Update(ctx context.Context, lesson *domain.Lesson) error
	Delete(ctx context.Context, id string) error
}

// ExerciseRepository defines the interface for exercise data operations.
type ExerciseRepository interface {
	Create(ctx context.Context, exercise *domain.Exercise) error
	GetByID(ctx context.Context, id string) (*domain.Exercise, error)
	GetByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error)
	GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error)
	Update(ctx context.Context, exercise *domain.Exercise) error
	Delete(ctx context.Context, id string) error
}

// ProgressRepository defines the interface for progress data operations.
type ProgressRepository interface {
	Create(ctx context.Context, progress *domain.Progress) error
	GetByID(ctx context.Context, id string) (*domain.Progress, error)
	GetByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) ([]*domain.Progress, int64, error)
	GetByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error)
	Update(ctx context.Context, progress *domain.Progress) error
	Delete(ctx context.Context, id string) error
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	UpdateLastActivity(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

// StreakRepository defines the interface for streak data operations.
type StreakRepository interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Streak, error)
	Upsert(ctx context.Context, streak *domain.Streak) error
	UpdateStreak(ctx context.Context, userID string, lastActivityDate *time.Time) error
}

// AssessmentRepository defines the interface for assessment data operations.
type AssessmentRepository interface {
	Create(ctx context.Context, assessment *domain.Assessment) error
	GetByID(ctx context.Context, id string) (*domain.Assessment, error)
	GetByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) ([]*domain.Assessment, int64, error)
	GetByType(ctx context.Context, assessmentType domain.AssessmentType, pagination *domain.PaginationRequest) ([]*domain.Assessment, int64, error)
	GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Assessment, int64, error)
	Update(ctx context.Context, assessment *domain.Assessment) error
	Delete(ctx context.Context, id string) error
}

// QuestionRepository defines the interface for quiz question data operations.
type QuestionRepository interface {
	Create(ctx context.Context, question *domain.Question) error
	GetByID(ctx context.Context, id string) (*domain.Question, error)
	GetByAssessmentID(ctx context.Context, assessmentID string) ([]*domain.Question, error)
	GetByTag(ctx context.Context, tag string, pagination *domain.PaginationRequest) ([]*domain.Question, int64, error)
	Update(ctx context.Context, question *domain.Question) error
	Delete(ctx context.Context, id string) error
}

// SubmissionRepository defines the interface for submission data operations.
type SubmissionRepository interface {
	Create(ctx context.Context, submission *domain.Submission) error
	GetByID(ctx context.Context, id string) (*domain.Submission, error)
	GetByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error)
	GetByAssessmentID(ctx context.Context, assessmentID string, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error)
	GetByUserAndAssessment(ctx context.Context, userID, assessmentID string) (*domain.Submission, error)
	GetAll(ctx context.Context, filters map[string]interface{}, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error)
	Update(ctx context.Context, submission *domain.Submission) error
	Delete(ctx context.Context, id string) error
}

// SubmissionCommentRepository defines the interface for submission comment data operations.
type SubmissionCommentRepository interface {
	Create(ctx context.Context, comment *domain.SubmissionComment) error
	GetBySubmissionID(ctx context.Context, submissionID string) ([]*domain.SubmissionComment, error)
	GetByID(ctx context.Context, id string) (*domain.SubmissionComment, error)
	Update(ctx context.Context, comment *domain.SubmissionComment) error
	Delete(ctx context.Context, id string) error
}

// PeerReviewRepository defines the interface for peer review data operations.
type PeerReviewRepository interface {
	Create(ctx context.Context, review *domain.PeerReview) error
	GetByID(ctx context.Context, id string) (*domain.PeerReview, error)
	GetBySubmissionID(ctx context.Context, submissionID string) ([]*domain.PeerReview, error)
	GetByReviewerID(ctx context.Context, reviewerID string, filters map[string]interface{}, pagination *domain.PaginationRequest) ([]*domain.PeerReview, int64, error)
	Update(ctx context.Context, review *domain.PeerReview) error
	AssignReviewers(ctx context.Context, submissionID string, reviewerIDs []string, deadline *time.Time) error
	GetPendingReviews(ctx context.Context, reviewerID string, beforeDeadline time.Time) ([]*domain.PeerReview, error)
	Delete(ctx context.Context, id string) error
}

// Repositories aggregates all repository interfaces.
type Repositories struct {
	Course             CourseRepository
	Lesson             LessonRepository
	Exercise           ExerciseRepository
	Progress           ProgressRepository
	User               UserRepository
	Streak             StreakRepository
	Assessment         AssessmentRepository
	Question           QuestionRepository
	Submission         SubmissionRepository
	SubmissionComment  SubmissionCommentRepository
	PeerReview         PeerReviewRepository
}
