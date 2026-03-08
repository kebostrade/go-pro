// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package domain defines learning path domain models.
package domain

import "time"

// LearningPath represents a curated learning track.
type LearningPath struct {
	ID            string            `json:"id"`
	Title         string            `json:"title" validate:"required,min=3,max=200"`
	Description   string            `json:"description" validate:"required,min=10,max=1000"`
	TrackType     TrackType         `json:"track_type" validate:"required"`
	RoleTarget    string            `json:"role_target" validate:"required,max=100"`
	Difficulty    Difficulty        `json:"difficulty" validate:"required"`
	LessonIDs     []string          `json:"lesson_ids" validate:"required,min=1"`
	Prerequisites []PrerequisiteRule `json:"prerequisites"`
	IsPublished   bool              `json:"is_published"`
	EstimatedDuration string        `json:"estimated_duration"` // e.g., "12 weeks"
	CreatedBy     string            `json:"created_by" validate:"required"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	PublishedAt   *time.Time        `json:"published_at,omitempty"`
}

// TrackType represents different learning tracks.
type TrackType string

const (
	TrackWebDev      TrackType = "web_dev"
	TrackSystems     TrackType = "systems"
	TrackDistributed TrackType = "distributed"
	TrackCloudNative TrackType = "cloud_native"
	TrackFullStack   TrackType = "full_stack"
)

// IsValid checks if the track type is valid.
func (t TrackType) IsValid() bool {
	switch t {
	case TrackWebDev, TrackSystems, TrackDistributed, TrackCloudNative, TrackFullStack:
		return true
	default:
		return false
	}
}

// String returns the string representation of track type.
func (t TrackType) String() string {
	return string(t)
}

// PrerequisiteRule represents a prerequisite requirement.
type PrerequisiteRule struct {
	ID      string       `json:"id"`
	Type    PrereqType   `json:"type" validate:"required"`
	Config  interface{}  `json:"config" validate:"required"` // Type-specific config
}

// PrereqType represents the type of prerequisite rule.
type PrereqType string

const (
	PrereqSequential    PrereqType = "sequential"
	PrereqPhaseComplete PrereqType = "phase_complete"
	PrereqScoreBased    PrereqType = "score_based"
)

// SequentialPrereqConfig represents sequential prerequisite config.
type SequentialPrereqConfig struct {
	BeforeLessonID string `json:"before_lesson_id" validate:"required"`
	AfterLessonID  string `json:"after_lesson_id" validate:"required"`
}

// PhaseCompletePrereqConfig represents phase-based prerequisite config.
type PhaseCompletePrereqConfig struct {
	PhaseID          string `json:"phase_id" validate:"required"`
	UnlocksLessonID  string `json:"unlocks_lesson_id" validate:"required"`
}

// ScoreBasedPrereqConfig represents score-based prerequisite config.
type ScoreBasedPrereqConfig struct {
	LessonID        string  `json:"lesson_id" validate:"required"`
	MinScore        float64 `json:"min_score" validate:"required,min=0,max=100"`
	UnlocksLessonID string  `json:"unlocks_lesson_id" validate:"required"`
}

// UserPathProgress represents a user's progress in a learning path.
type UserPathProgress struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id" validate:"required"`
	PathID               string     `json:"path_id" validate:"required"`
	EnrolledAt           time.Time  `json:"enrolled_at"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	CurrentLessonIndex   int        `json:"current_lesson_index" validate:"min=0"`
	OverallProgress      float64    `json:"overall_progress"` // 0-100 percentage
	LastAccessAt         time.Time  `json:"last_access_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// LessonProgress represents progress on a specific lesson within a path.
type LessonProgress struct {
	LessonID       string     `json:"lesson_id"`
	IsCompleted    bool       `json:"is_completed"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Score          *float64   `json:"score,omitempty"`
	TimeSpent      time.Duration `json:"time_spent"`
}

// PathCertificate represents a completion certificate.
type PathCertificate struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	PathID          string    `json:"path_id"`
	PathTitle       string    `json:"path_title"`
	StudentName     string    `json:"student_name"`
	CompletionDate  time.Time `json:"completion_date"`
	VerificationURL string    `json:"verification_url"`
	CertificateURL  string    `json:"certificate_url"` // S3 URL
	CreatedAt       time.Time `json:"created_at"`
}

// CreateLearningPathRequest represents a request to create a learning path.
type CreateLearningPathRequest struct {
	Title              string            `json:"title" validate:"required,min=3,max=200"`
	Description        string            `json:"description" validate:"required,min=10,max=1000"`
	TrackType          TrackType         `json:"track_type" validate:"required"`
	RoleTarget         string            `json:"role_target" validate:"required,max=100"`
	Difficulty         Difficulty        `json:"difficulty" validate:"required"`
	LessonIDs          []string          `json:"lesson_ids" validate:"required,min=1"`
	Prerequisites      []PrerequisiteRule `json:"prerequisites"`
	EstimatedDuration  string            `json:"estimated_duration" validate:"required"`
}

// UpdateLearningPathRequest represents a request to update a learning path.
type UpdateLearningPathRequest struct {
	Title             *string            `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description       *string            `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	TrackType         *TrackType         `json:"track_type,omitempty"`
	RoleTarget        *string            `json:"role_target,omitempty" validate:"omitempty,max=100"`
	Difficulty        *Difficulty        `json:"difficulty,omitempty"`
	LessonIDs         *[]string          `json:"lesson_ids,omitempty" validate:"omitempty,min=1"`
	Prerequisites     *[]PrerequisiteRule `json:"prerequisites,omitempty"`
	EstimatedDuration *string            `json:"estimated_duration,omitempty"`
}

// PublishPathRequest represents a request to publish a learning path.
type PublishPathRequest struct {
	PathID      string `json:"path_id" validate:"required"`
	PublishNow  bool   `json:"publish_now"`
}

// EnrollInPathRequest represents a request to enroll in a learning path.
type EnrollInPathRequest struct {
	PathID             string `json:"path_id" validate:"required"`
	TransferProgress    bool   `json:"transfer_progress"` // Transfer progress from common lessons
}

// PathProgressResponse represents detailed progress in a learning path.
type PathProgressResponse struct {
	PathID             string                 `json:"path_id"`
	PathTitle          string                 `json:"path_title"`
	EnrolledAt         time.Time              `json:"enrolled_at"`
	OverallProgress    float64                `json:"overall_progress"`
	CurrentLessonIndex int                    `json:"current_lesson_index"`
	EstimatedCompletionDate *time.Time        `json:"estimated_completion_date,omitempty"`
	Lessons            []LessonProgress       `json:"lessons"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
	Certificate        *PathCertificate       `json:"certificate,omitempty"`
}

// PathComparison represents side-by-side path comparison.
type PathComparison struct {
	Paths      []LearningPathSummary `json:"paths"`
	Differences []PathDifference      `json:"differences"`
}

// LearningPathSummary represents a summary of a learning path.
type LearningPathSummary struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	TrackType          TrackType  `json:"track_type"`
	RoleTarget         string     `json:"role_target"`
	Difficulty         Difficulty `json:"difficulty"`
	LessonCount        int        `json:"lesson_count"`
	EstimatedDuration  string     `json:"estimated_duration"`
	EnrolledCount      int        `json:"enrolled_count"`
	IsPublished        bool       `json:"is_published"`
}

// PathDifference highlights differences between paths.
type PathDifference struct {
	Field      string      `json:"field"`
	PathValues []PathValue `json:"path_values"`
}

// PathValue represents a value for a specific path.
type PathValue struct {
	PathID string      `json:"path_id"`
	Value  interface{} `json:"value"`
}

// PrerequisiteCheckResult represents the result of prerequisite validation.
type PrerequisiteCheckResult struct {
	CanAccess       bool     `json:"can_access"`
	Reason          string   `json:"reason,omitempty"`
	MissingReqs     []string `json:"missing_requirements,omitempty"`
	SatisfiedReqs   []string `json:"satisfied_requirements,omitempty"`
}

// ValidatePrerequisitesRequest represents a request to validate prerequisites.
type ValidatePrerequisitesRequest struct {
	PathID        string `json:"path_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
	LessonIndex   int    `json:"lesson_index" validate:"required,min=0"`
}
