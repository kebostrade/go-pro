// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package domain defines the core business entities and models.
package domain

import (
	"time"
)

// Course represents a learning course.
type Course struct {
	ID          string    `json:"id" validate:"required,slug"`
	Title       string    `json:"title" validate:"required,min=3,max=200"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Lessons     []string  `json:"lessons"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Lesson represents a course lesson.
type Lesson struct {
	ID          string    `json:"id" validate:"required,slug"`
	CourseID    string    `json:"course_id" validate:"required,slug"`
	Title       string    `json:"title" validate:"required,min=3,max=200"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Content     string    `json:"content" validate:"required"`
	Order       int       `json:"order" validate:"required,min=1"`
	Exercises   []string  `json:"exercises"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Exercise represents a coding exercise.
type Exercise struct {
	ID          string     `json:"id" validate:"required,slug"`
	LessonID    string     `json:"lesson_id" validate:"required,slug"`
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"required,min=10,max=1000"`
	TestCases   int        `json:"test_cases" validate:"required,min=1"`
	Difficulty  Difficulty `json:"difficulty" validate:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Status represents progress status levels.
type Status string

const (
	StatusNotStarted Status = "not_started"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

// IsValid checks if the status is valid.
func (s Status) IsValid() bool {
	switch s {
	case StatusNotStarted, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

// String returns the string representation of status.
func (s Status) String() string {
	return string(s)
}

// Progress represents user learning progress.
type Progress struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id" validate:"required"`
	LessonID    string     `json:"lesson_id" validate:"required,slug"`
	Status      Status     `json:"status" validate:"required"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserProgressSummary represents a summary of user progress.
type UserProgressSummary struct {
	UserID           string `json:"user_id"`
	LessonsStarted   int    `json:"lessons_started"`
	LessonsCompleted int    `json:"lessons_completed"`
	TotalLessons     int    `json:"total_lessons"`
	OverallProgress  int    `json:"overall_progress"`
}

// Difficulty represents exercise difficulty levels.
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
	DifficultyExpert       Difficulty = "expert"
)

// IsValid checks if the difficulty is valid.
func (d Difficulty) IsValid() bool {
	switch d {
	case DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced, DifficultyExpert:
		return true
	default:
		return false
	}
}

// String returns the string representation of difficulty.
func (d Difficulty) String() string {
	return string(d)
}

// CreateCourseRequest represents a request to create a new course.
type CreateCourseRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
}

// UpdateCourseRequest represents a request to update a course.
type UpdateCourseRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
}

// CreateLessonRequest represents a request to create a new lesson.
type CreateLessonRequest struct {
	CourseID    string `json:"course_id" validate:"required,slug"`
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
	Content     string `json:"content" validate:"required"`
	Order       int    `json:"order" validate:"required,min=1"`
}

// UpdateLessonRequest represents a request to update a lesson.
type UpdateLessonRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Content     *string `json:"content,omitempty"`
	Order       *int    `json:"order,omitempty" validate:"omitempty,min=1"`
}

// CreateExerciseRequest represents a request to create a new exercise.
type CreateExerciseRequest struct {
	LessonID    string     `json:"lesson_id" validate:"required,slug"`
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"required,min=10,max=1000"`
	TestCases   int        `json:"test_cases" validate:"required,min=1"`
	Difficulty  Difficulty `json:"difficulty" validate:"required"`
}

// UpdateExerciseRequest represents a request to update an exercise.
type UpdateExerciseRequest struct {
	Title       *string     `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string     `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	TestCases   *int        `json:"test_cases,omitempty" validate:"omitempty,min=1"`
	Difficulty  *Difficulty `json:"difficulty,omitempty"`
}

// CreateProgressRequest represents a request to create a new progress record.
type CreateProgressRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	LessonID string `json:"lesson_id" validate:"required,slug"`
	Status   Status `json:"status" validate:"required"`
}

// UpdateProgressRequest represents a request to update a progress record.
type UpdateProgressRequest struct {
	Status *Status `json:"status,omitempty"`
}

// SubmitExerciseRequest represents a request to submit an exercise solution.
type SubmitExerciseRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required,oneof=go python javascript"`
}

// ExerciseSubmissionResult represents the result of an exercise submission.
type ExerciseSubmissionResult struct {
	Success         bool         `json:"success"` // Always true for successful API call
	ExerciseID      string       `json:"exercise_id"`
	Score           int          `json:"score"`
	Passed          bool         `json:"passed"`
	Message         string       `json:"message"`
	TestResults     []TestResult `json:"results,omitempty"` // Changed from test_results to results
	ExecutionTimeMs int64        `json:"execution_time_ms"`
	SubmittedAt     time.Time    `json:"submitted_at"`
}

// TestResult represents the result of a single test case.
type TestResult struct {
	TestName string `json:"test_name"`
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Error    string `json:"error,omitempty"`
}

// User represents a user in the system.
type User struct {
	ID           string     `json:"id"`
	FirebaseUID  string     `json:"firebase_uid" validate:"required"` // Firebase User ID
	Username     string     `json:"username" validate:"required,min=3,max=50"`
	Email        string     `json:"email" validate:"required,email"`
	DisplayName  string     `json:"display_name,omitempty"` // Full name from Firebase
	PhotoURL     string     `json:"photo_url,omitempty"`    // Profile picture from Firebase
	PasswordHash string     `json:"-"`                      // Never expose password hash in JSON
	FirstName    string     `json:"first_name,omitempty"`
	LastName     string     `json:"last_name,omitempty"`
	Role         UserRole   `json:"role"`            // Single role: student or admin
	Roles        []string   `json:"roles,omitempty"` // Legacy: multiple roles support
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// UserRole represents user roles in the system.
type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleAdmin   UserRole = "admin"
)

// IsValid checks if the role is valid.
func (r UserRole) IsValid() bool {
	switch r {
	case RoleStudent, RoleAdmin:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role.
func (r UserRole) String() string {
	return string(r)
}

// ToProfileResponse converts a User to UserProfileResponse.
func (u *User) ToProfileResponse() *UserProfileResponse {
	lastLogin := time.Time{}
	if u.LastLoginAt != nil {
		lastLogin = *u.LastLoginAt
	}
	return &UserProfileResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PhotoURL:    u.PhotoURL,
		Role:        u.Role,
		IsActive:    u.IsActive,
		LastLoginAt: lastLogin,
		CreatedAt:   u.CreatedAt,
	}
}

// UserProfileResponse represents a user profile response.
type UserProfileResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	Role        UserRole  `json:"role"`
	IsActive    bool      `json:"is_active"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterRequest represents a user registration request.
type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FirstName string `json:"first_name,omitempty" validate:"max=50"`
	LastName  string `json:"last_name,omitempty" validate:"max=50"`
}

// LoginRequest represents a user login request.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents a token refresh request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest represents a password change request.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// UpdateProfileRequest represents a profile update request.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=50"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=50"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
}

// VerifyTokenRequest represents a request to verify a Firebase ID token.
type VerifyTokenRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

// VerifyTokenResponse represents the response after token verification.
type VerifyTokenResponse struct {
	User          *UserProfileResponse `json:"user"`
	IsNewUser     bool                 `json:"is_new_user"`
	TokenVerified bool                 `json:"token_verified"`
}

// FirebaseClaims represents custom claims from Firebase token.
type FirebaseClaims struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"name,omitempty"`
	Picture     string    `json:"picture,omitempty"`
	IssuedAt    time.Time `json:"iat"`
	ExpiresAt   time.Time `json:"exp"`
}

// UpdateUserRequest represents a request to update user profile.
type UpdateUserRequest struct {
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	PhotoURL    *string `json:"photo_url,omitempty" validate:"omitempty,url"`
}

// UpdateUserRoleRequest represents a request to update user role (admin only).
type UpdateUserRoleRequest struct {
	Role UserRole `json:"role" validate:"required"`
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

// APIResponse represents a standard API response.
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents an error in API responses.
type APIError struct {
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// PaginationRequest represents pagination parameters.
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata.
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ListResponse represents a paginated list response.
type ListResponse struct {
	Items      interface{}         `json:"items"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// Curriculum represents the complete learning curriculum.
type Curriculum struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    string            `json:"duration"`
	Phases      []CurriculumPhase `json:"phases"`
	Projects    []Project         `json:"projects"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CurriculumPhase represents a phase in the curriculum.
type CurriculumPhase struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Weeks       string             `json:"weeks"`
	Icon        string             `json:"icon"`
	Color       string             `json:"color"`
	Order       int                `json:"order"`
	Lessons     []CurriculumLesson `json:"lessons"`
	Progress    int                `json:"progress"`
}

// CurriculumLesson represents a lesson in the curriculum.
type CurriculumLesson struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Duration    string     `json:"duration"`
	Exercises   int        `json:"exercises"`
	Difficulty  Difficulty `json:"difficulty"`
	Completed   bool       `json:"completed"`
	Locked      bool       `json:"locked"`
	Order       int        `json:"order"`
}

// Project represents a hands-on project.
type Project struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Duration    string     `json:"duration"`
	Difficulty  Difficulty `json:"difficulty"`
	Skills      []string   `json:"skills"`
	Points      int        `json:"points"`
	Completed   bool       `json:"completed"`
	Locked      bool       `json:"locked"`
	Order       int        `json:"order"`
}

// LessonDetail represents detailed lesson content.
type LessonDetail struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Duration     string           `json:"duration"`
	Difficulty   Difficulty       `json:"difficulty"`
	Phase        string           `json:"phase"`
	Objectives   []string         `json:"objectives"`
	Theory       string           `json:"theory"`
	CodeExample  string           `json:"code_example"`
	Solution     string           `json:"solution"`
	Exercises    []LessonExercise `json:"exercises"`
	NextLessonID *int             `json:"next_lesson_id,omitempty"`
	PrevLessonID *int             `json:"prev_lesson_id,omitempty"`
}

// LessonExercise represents an exercise within a lesson.
type LessonExercise struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	InitialCode  string   `json:"initial_code"`
	Solution     string   `json:"solution"`
}
