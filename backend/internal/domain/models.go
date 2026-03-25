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
	LastActivityDate *time.Time `json:"last_activity_date,omitempty"` // For streak tracking
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

// Streak represents user streak tracking data.
type Streak struct {
	UserID           string     `json:"user_id" validate:"required"`
	CurrentStreak    int        `json:"current_streak" validate:"min=0"`
	LongestStreak    int        `json:"longest_streak" validate:"min=0"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UpdateStreakRequest represents a request to update streak data.
type UpdateStreakRequest struct {
	LastActivityDate *time.Time `json:"last_activity_date,omitempty"`
}

// StreakResponse represents streak data in API responses.
type StreakResponse struct {
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty"`
}

// Assessment represents an assessment (quiz, coding exercise, or project).
type Assessment struct {
	ID             string                 `json:"id"`
	LessonID       string                 `json:"lesson_id" validate:"required,slug"`
	Type           AssessmentType         `json:"type" validate:"required"`
	Title          string                 `json:"title" validate:"required,min=3,max=200"`
	Description    string                 `json:"description" validate:"required,min=10,max=5000"`
	Config         map[string]interface{} `json:"config" validate:"required"` // Flexible config per type
	PassingScore   int                    `json:"passing_score" validate:"min=0,max=100"`
	TimeLimitMinutes *int                  `json:"time_limit_minutes,omitempty" validate:"omitempty,min=1,max=180"`
	OrderIndex     int                    `json:"order_index" validate:"required,min=1"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AssessmentType represents the type of assessment.
type AssessmentType string

const (
	AssessmentTypeQuiz           AssessmentType = "quiz"
	AssessmentTypeCodingExercise AssessmentType = "coding_exercise"
	AssessmentTypeProject        AssessmentType = "project"
)

// IsValid checks if the assessment type is valid.
func (a AssessmentType) IsValid() bool {
	switch a {
	case AssessmentTypeQuiz, AssessmentTypeCodingExercise, AssessmentTypeProject:
		return true
	default:
		return false
	}
}

// String returns the string representation of the assessment type.
func (a AssessmentType) String() string {
	return string(a)
}

// Submission represents a student assessment submission.
type Submission struct {
	ID           string             `json:"id"`
	AssessmentID string             `json:"assessment_id" validate:"required"`
	UserID       string             `json:"user_id" validate:"required"`
	Content      map[string]interface{} `json:"content" validate:"required"` // Flexible content
	Score        *int               `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Feedback     string             `json:"feedback,omitempty"`
	GradedBy     *string            `json:"graded_by,omitempty"` // User ID of grader
	GradedAt     *time.Time         `json:"graded_at,omitempty"`
	SubmittedAt  time.Time          `json:"submitted_at"`
	Status       SubmissionStatus   `json:"status" validate:"required"`
}

// SubmissionStatus represents the status of a submission.
type SubmissionStatus string

const (
	SubmissionStatusSubmitted SubmissionStatus = "submitted"
	SubmissionStatusGraded    SubmissionStatus = "graded"
	SubmissionStatusReturned  SubmissionStatus = "returned"
)

// IsValid checks if the submission status is valid.
func (s SubmissionStatus) IsValid() bool {
	switch s {
	case SubmissionStatusSubmitted, SubmissionStatusGraded, SubmissionStatusReturned:
		return true
	default:
		return false
	}
}

// String returns the string representation of the submission status.
func (s SubmissionStatus) String() string {
	return string(s)
}

// SubmissionComment represents an inline comment on a submission.
type SubmissionComment struct {
	ID           string    `json:"id"`
	SubmissionID string    `json:"submission_id" validate:"required"`
	AuthorID     string    `json:"author_id" validate:"required"`
	CommentText  string    `json:"comment_text" validate:"required,min=1,max=1000"`
	LineNumber   *int      `json:"line_number,omitempty" validate:"omitempty,min=1"`
	CreatedAt    time.Time `json:"created_at"`
}

// QuizSubmission represents quiz-specific submission content.
type QuizSubmission struct {
	Answers  map[string]interface{} `json:"answers"`  // Question ID -> answer
	Score    int                    `json:"score"`
	Attempts int                    `json:"attempts"`
}

// CodeSubmission represents code exercise submission content.
type CodeSubmission struct {
	Code       string `json:"code" validate:"required"`
	Language   string `json:"language" validate:"required"`
	TestPassed bool   `json:"test_passed"`
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
}

// ProjectSubmission represents project submission content.
type ProjectSubmission struct {
	SubmissionFormat string `json:"submission_format" validate:"required,oneof=repo_link zip_upload text"`
	RepoLink         string `json:"repo_link,omitempty" validate:"omitempty,url"`
	FileURL          string `json:"file_url,omitempty" validate:"omitempty,url"`
	TextContent      string `json:"text_content,omitempty"`
}

// Question represents a quiz question.
type Question struct {
	ID             string       `json:"id" validate:"required"`
	AssessmentID   string       `json:"assessment_id" validate:"required"`
	QuestionType   QuestionType `json:"question_type" validate:"required"`
	QuestionText   string       `json:"question_text" validate:"required,min=10,max=1000"`
	Options        []string     `json:"options,omitempty"` // For multiple choice
	CorrectAnswer  string       `json:"correct_answer" validate:"required"`
	Explanation    *string      `json:"explanation,omitempty"`
	Points         int          `json:"points" validate:"required,min=1,max=100"`
	OrderIndex     int          `json:"order_index" validate:"required,min=1"`
	Tags           []string     `json:"tags,omitempty"` // For question bank reuse
	Hints          []string     `json:"hints,omitempty"` // Progressive hints
	HintThreshold  int          `json:"hint_threshold,omitempty"` // Show after N failed attempts
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// QuestionType represents quiz question types.
type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse     QuestionType = "true_false"
	QuestionTypeShortAnswer    QuestionType = "short_answer"
	QuestionTypeCodeCompletion QuestionType = "code_completion"
)

// IsValid checks if the question type is valid.
func (q QuestionType) IsValid() bool {
	switch q {
	case QuestionTypeMultipleChoice, QuestionTypeTrueFalse, QuestionTypeShortAnswer, QuestionTypeCodeCompletion:
		return true
	default:
		return false
	}
}

// String returns the string representation of question type.
func (q QuestionType) String() string {
	return string(q)
}

// QuizConfig represents quiz-specific configuration.
type QuizConfig struct {
	Questions []Question   `json:"questions" validate:"required,min=1"`
	Settings  QuizSettings `json:"settings" validate:"required"`
}

// QuizSettings represents quiz settings.
type QuizSettings struct {
	MaxAttempts      int  `json:"max_attempts" validate:"min=1,max=10"`
	ShuffleQuestions bool `json:"shuffle_questions"`
	ShowExplanations  bool `json:"show_explanations"`
}

// CodingExerciseConfig represents coding exercise configuration.
type CodingExerciseConfig struct {
	StarterCode   string          `json:"starter_code" validate:"required"`
	TestCases     []TestCase      `json:"test_cases" validate:"required,min=1"`
	Limits        ExecutionLimits `json:"limits" validate:"required"`
	Hints         []string        `json:"hints,omitempty"`
	SolutionCode  string          `json:"solution_code,omitempty"` // Instructor's solution
}

// TestCase represents a test case for coding exercises.
type TestCase struct {
	ID             string    `json:"id" validate:"required"`
	AssessmentID   string    `json:"assessment_id" validate:"required"`
	Input          string    `json:"input" validate:"required"`
	ExpectedOutput string    `json:"expected_output" validate:"required"`
	IsVisible      bool      `json:"is_visible"` // Visible to students
	Points         int       `json:"points" validate:"required,min=1,max=100"`
	CreatedAt      time.Time `json:"created_at"`
}

// ExecutionLimits represents execution limits for code.
type ExecutionLimits struct {
	TimeLimitSeconds int `json:"time_limit_seconds" validate:"min=1,max=30"`
	MemoryLimitMB     int `json:"memory_limit_mb" validate:"min=64,max=256"`
}

// ExecuteRequest represents a code execution request.
type ExecuteRequest struct {
	Code      string        `json:"code" validate:"required"`
	Language string        `json:"language" validate:"required,oneof=go python javascript"`
	Timeout   time.Duration `json:"timeout"`
	TestCases []TestCase  `json:"test_cases" validate:"required,min=1"`
}

// ExecuteResult represents the result of code execution.
type ExecuteResult struct {
	Passed        bool              `json:"passed"`
	Score         int              `json:"score"`
	Results       []TestResult  `json:"results"`
	ExecutionTime time.Duration `json:"execution_time"`
	Error         error           `json:"error,omitempty"`
}

// TestCaseForExecution represents a test case for code execution.
type TestCaseForExecution struct {
	Name     string `json:"name" validate:"required"`
	Input    string `json:"input"`
	Expected string `json:"expected" validate:"required"`
}

// TestResultForExecution represents the result of a single test execution.
type TestResultForExecution struct {
	TestName string `json:"test_name"`
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Error    string `json:"error,omitempty"`
}

// ProjectAssignmentConfig represents project assignment configuration.
type ProjectAssignmentConfig struct {
	Deliverables     []string          `json:"deliverables" validate:"required,min=1"`
	SubmissionFormat SubmissionFormat  `json:"submission_format" validate:"required"`
	StarterCodeURL   *string           `json:"starter_code_url,omitempty" validate:"omitempty,url"`
	Rubric           []RubricCriterion `json:"rubric" validate:"required,min=1"`
	RequirePeerReview bool             `json:"require_peer_review"`
	PeerReviewCount  int               `json:"peer_review_count,omitempty" validate:"min=1,max=5"`
}

// SubmissionFormat represents project submission formats.
type SubmissionFormat string

const (
	SubmissionFormatRepoLink  SubmissionFormat = "repo_link"
	SubmissionFormatZipUpload SubmissionFormat = "zip_upload"
	SubmissionFormatText      SubmissionFormat = "text"
)

// IsValid checks if the submission format is valid.
func (s SubmissionFormat) IsValid() bool {
	switch s {
	case SubmissionFormatRepoLink, SubmissionFormatZipUpload, SubmissionFormatText:
		return true
	default:
		return false
	}
}

// String returns the string representation of submission format.
func (s SubmissionFormat) String() string {
	return string(s)
}

// RubricCriterion represents a grading rubric criterion.
type RubricCriterion struct {
	ID                string             `json:"id" validate:"required"`
	AssessmentID      string             `json:"assessment_id" validate:"required"`
	Description       string             `json:"description" validate:"required,min=10,max=500"`
	MaxPoints         int                `json:"max_points" validate:"required,min=1,max=100"`
	PerformanceLevels []PerformanceLevel `json:"performance_levels" validate:"required,min=2,max=5"`
	OrderIndex        int                `json:"order_index" validate:"required,min=1"`
}

// PerformanceLevel represents a rubric performance level.
type PerformanceLevel struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"required,min=10,max=500"`
	Points      int    `json:"points" validate:"required,min=0"`
}

// GradebookEntry represents a single gradebook entry.
type GradebookEntry struct {
	StudentID       string           `json:"student_id"`
	StudentName     string           `json:"student_name"`
	StudentEmail    string           `json:"student_email"`
	AssessmentID    string           `json:"assessment_id"`
	AssessmentTitle string           `json:"assessment_title"`
	Score           *int             `json:"score,omitempty"`
	SubmittedAt     *time.Time       `json:"submitted_at,omitempty"`
	GradedAt        *time.Time       `json:"graded_at,omitempty"`
	Status          SubmissionStatus `json:"status"`
}

// CreateAssessmentRequest represents a request to create an assessment.
type CreateAssessmentRequest struct {
	LessonID        string   `json:"lesson_id" validate:"required,slug"`
	Type            AssessmentType `json:"type" validate:"required"`
	Title           string   `json:"title" validate:"required,min=3,max=200"`
	Description     string   `json:"description" validate:"required,min=10,max=5000"`
	Config          map[string]interface{} `json:"config" validate:"required"`
	PassingScore    int      `json:"passing_score" validate:"min=0,max=100"`
	TimeLimitMinutes *int     `json:"time_limit_minutes,omitempty" validate:"omitempty,min=1,max=180"`
	OrderIndex      int      `json:"order_index" validate:"required,min=1"`
}

// UpdateAssessmentRequest represents a request to update an assessment.
type UpdateAssessmentRequest struct {
	Title           *string  `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description     *string  `json:"description,omitempty" validate:"omitempty,min=10,max=5000"`
	Config          *map[string]interface{} `json:"config,omitempty"`
	PassingScore    *int     `json:"passing_score,omitempty" validate:"omitempty,min=0,max=100"`
	TimeLimitMinutes *int     `json:"time_limit_minutes,omitempty" validate:"omitempty,min=1,max=180"`
	OrderIndex      *int     `json:"order_index,omitempty" validate:"omitempty,min=1"`
}

// CreateSubmissionRequest represents a request to create a submission.
type CreateSubmissionRequest struct {
	AssessmentID string                 `json:"assessment_id" validate:"required"`
	Content      map[string]interface{} `json:"content" validate:"required"`
}

// GradeSubmissionRequest represents a request to grade a submission.
type GradeSubmissionRequest struct {
	Score            *int                       `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Feedback         *string                    `json:"feedback,omitempty" validate:"omitempty,max=5000"`
	RubricScores     map[string]int             `json:"rubric_scores,omitempty"` // Criterion ID -> score
	ReleaseImmediately bool                    `json:"release_immediately"`
}

// CreateSubmissionCommentRequest represents a request to add a comment to a submission.
type CreateSubmissionCommentRequest struct {
	CommentText string `json:"comment_text" validate:"required,min=1,max=1000"`
	LineNumber  *int   `json:"line_number,omitempty" validate:"omitempty,min=1"`
}

// ContentVersion represents a version of content in the CMS.
type ContentVersion struct {
	ID             int64      `json:"id"`
	ContentType    string     `json:"content_type" validate:"required,oneof=lesson exercise"` // lesson or exercise
	ContentID      string     `json:"content_id" validate:"required"`
	VersionNumber  int        `json:"version_number" validate:"required,min=1"`
	Title          string     `json:"title,omitempty"`
	Content        string     `json:"content,omitempty"`
	Difficulty     Difficulty `json:"difficulty,omitempty"`
	Objectives     []string   `json:"objectives,omitempty"`
	Theory         string     `json:"theory,omitempty"`
	CodeExample    string     `json:"code_example,omitempty"`
	Solution       string     `json:"solution,omitempty"`
	Exercises      []string   `json:"exercises,omitempty"`
	ChangeSummary  string     `json:"change_summary,omitempty"`
	ChangedBy      string     `json:"changed_by" validate:"required"`
	IsMajorRevision bool      `json:"is_major_revision"`
	IsPublished    bool       `json:"is_published"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ContentVersionRequest represents a request to create a content version.
type ContentVersionRequest struct {
	ContentType    string     `json:"content_type" validate:"required,oneof=lesson exercise"`
	ContentID      string     `json:"content_id" validate:"required"`
	ChangeSummary  string     `json:"change_summary" validate:"required,min=5,max=500"`
	IsMajorRevision bool     `json:"is_major_revision"`
}

// ContentVersionHistory represents the version history of content.
type ContentVersionHistory struct {
	ContentID      string            `json:"content_id"`
	ContentType    string            `json:"content_type"`
	TotalVersions  int               `json:"total_versions"`
	Versions       []ContentVersion  `json:"versions"`
}

// PublishContentRequest represents a request to publish a content version.
type PublishContentRequest struct {
	ContentType   string `json:"content_type" validate:"required,oneof=lesson exercise"`
	ContentID     string `json:"content_id" validate:"required"`
	VersionNumber int    `json:"version_number" validate:"required,min=1"`
}

// PeerReview represents a peer review assignment.
type PeerReview struct {
	ID             string         `json:"id"`
	SubmissionID   string         `json:"submission_id" validate:"required"`
	ReviewerID     string         `json:"reviewer_id" validate:"required"`
	RubricScores   map[string]int `json:"rubric_scores" validate:"required"` // Criterion ID -> Score
	Feedback       string         `json:"feedback" validate:"required,min=1,max=5000"`
	InlineComments []InlineComment `json:"inline_comments"`
	IsAnonymous    bool           `json:"is_anonymous"`
	KarmaPoints    int            `json:"karma_points" validate:"min=0,max=5"`
	Status         PeerReviewStatus `json:"status" validate:"required"`
	SubmittedAt    *time.Time     `json:"submitted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
}

// PeerReviewStatus represents the status of a peer review.
type PeerReviewStatus string

const (
	PeerReviewStatusPending   PeerReviewStatus = "pending"
	PeerReviewStatusCompleted PeerReviewStatus = "completed"
	PeerReviewStatusLate      PeerReviewStatus = "late"
)

// IsValid checks if the peer review status is valid.
func (p PeerReviewStatus) IsValid() bool {
	switch p {
	case PeerReviewStatusPending, PeerReviewStatusCompleted, PeerReviewStatusLate:
		return true
	default:
		return false
	}
}

// String returns the string representation of peer review status.
func (p PeerReviewStatus) String() string {
	return string(p)
}

// InlineComment represents an inline comment on code.
type InlineComment struct {
	LineNumber  int    `json:"line_number" validate:"required,min=1"`
	CommentText string `json:"comment_text" validate:"required,min=1,max=1000"`
	FilePath    string `json:"file_path,omitempty"`
}

// CreatePeerReviewRequest represents a request to create a peer review.
type CreatePeerReviewRequest struct {
	SubmissionID string         `json:"submission_id" validate:"required"`
	ReviewerID   string         `json:"reviewer_id" validate:"required"`
	Deadline     *time.Time     `json:"deadline,omitempty"`
}

// UpdatePeerReviewRequest represents a request to update a peer review.
type UpdatePeerReviewRequest struct {
	RubricScores   *map[string]int  `json:"rubric_scores,omitempty" validate:"omitempty,min=1"`
	Feedback       *string          `json:"feedback,omitempty" validate:"omitempty,min=1,max=5000"`
	InlineComments *[]InlineComment `json:"inline_comments,omitempty"`
	IsAnonymous    *bool            `json:"is_anonymous,omitempty"`
}

// CompletePeerReviewRequest represents a request to complete a peer review.
type CompletePeerReviewRequest struct {
	RubricScores   map[string]int  `json:"rubric_scores" validate:"required,min=1"`
	Feedback       string          `json:"feedback" validate:"required,min=1,max=5000"`
	InlineComments []InlineComment `json:"inline_comments"`
	IsAnonymous    bool            `json:"is_anonymous"`
}

// AwardKarmaRequest represents a request to award karma points.
type AwardKarmaRequest struct {
	KarmaPoints int `json:"karma_points" validate:"required,min=1,max=5"`
}

// PeerReviewAssignmentConfig represents configuration for auto-assignment.
type PeerReviewAssignmentConfig struct {
	ReviewersPerAssignment int        `json:"reviewers_per_assignment" validate:"required,min=1,max=5"`
	AssignmentMethod        string     `json:"assignment_method" validate:"required,oneof=round_robin random"`
	DeadlineHours          int        `json:"deadline_hours" validate:"required,min=24,max=168"`
	RemindersEnabled       bool       `json:"reminders_enabled"`
	ReminderHours          []int      `json:"reminder_hours" validate:"omitempty,min=1,dive,min=1"`
	LatePenaltyPerDay      float64    `json:"late_penalty_per_day" validate:"required,min=0,max=0.3"`
	MaxLatePenalty         float64    `json:"max_late_penalty" validate:"required,min=0,max=0.5"`
}

// PeerReviewSummary represents a summary of peer reviews for a submission.
type PeerReviewSummary struct {
	SubmissionID        string             `json:"submission_id"`
	TotalReviews        int                `json:"total_reviews"`
	CompletedReviews    int                `json:"completed_reviews"`
	AverageScore        float64            `json:"average_score"`
	ReviewCount         int                `json:"review_count"`
	RubricAverageScores map[string]float64 `json:"rubric_average_scores,omitempty"`
	Reviews             []*PeerReview      `json:"reviews,omitempty"`
}

// ReviewerStats represents statistics for a reviewer.
type ReviewerStats struct {
	ReviewerID        string    `json:"reviewer_id"`
	TotalAssigned     int       `json:"total_assigned"`
	TotalCompleted    int       `json:"total_completed"`
	TotalPending      int       `json:"total_pending"`
	AverageKarma      float64   `json:"average_karma"`
	OnTimeRate        float64   `json:"on_time_rate"`
	AverageReviewTime float64   `json:"average_review_time_hours"`
	LastReviewDate    *time.Time `json:"last_review_date,omitempty"`
}
