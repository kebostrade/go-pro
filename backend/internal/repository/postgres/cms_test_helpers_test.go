// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"

	"github.com/stretchr/testify/require"
)

// CMSTestHelper provides CMS-specific test utilities.
type CMSTestHelper struct {
	t  *testing.T
	db *sql.DB
}

// NewCMSTestHelper creates a new CMS test helper.
func NewCMSTestHelper(t *testing.T, db *sql.DB) *CMSTestHelper {
	return &CMSTestHelper{
		t:  t,
		db: db,
	}
}

// CreateTestContentVersion creates a test content version.
func (h *CMSTestHelper) CreateTestContentVersion(ctx context.Context, lessonID string, version int, status string, authorID string) *domain.ContentVersion {
	contentStr := `{"sections": [{"type": "text", "content": "Test content section"}, {"type": "code", "content": "func main() { println('test') }"}]}`

	now := time.Now()
	versionID := rand.Int63()
	contentVersion := &domain.ContentVersion{
		ID:              versionID,
		ContentType:     "lesson",
		ContentID:       lessonID,
		VersionNumber:   version,
		Content:         contentStr,
		ChangedBy:       authorID,
		ChangeSummary:   "Test version",
		IsPublished:     status == "published",
		IsMajorRevision: version == 1,
		CreatedAt:       now,
	}

	query := `
		INSERT INTO content_versions (id, content_type, content_id, version_number, content, changed_by, change_summary, is_major_revision, is_published, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			is_published = EXCLUDED.is_published
	`

	_, err := h.db.ExecContext(ctx, query,
		contentVersion.ID,
		contentVersion.ContentType,
		contentVersion.ContentID,
		contentVersion.VersionNumber,
		contentVersion.Content,
		contentVersion.ChangedBy,
		contentVersion.ChangeSummary,
		contentVersion.IsMajorRevision,
		contentVersion.IsPublished,
		contentVersion.CreatedAt,
	)
	require.NoError(h.t, err, "Failed to create test content version")

	return contentVersion
}

// CreateTestAssessment creates a test assessment.
func (h *CMSTestHelper) CreateTestAssessment(ctx context.Context, lessonID string, assessmentType string) *domain.Assessment {
	assessmentID := testutil.RandomString(10)
	now := time.Now()
	assessment := &domain.Assessment{
		ID:              assessmentID,
		LessonID:        lessonID,
		Type:            domain.AssessmentType(assessmentType),
		Title:           "Test Assessment",
		Description:     "Test assessment description",
		Config:          map[string]interface{}{"time_limit_minutes": 30},
		PassingScore:    80,
		TimeLimitMinutes: func() *int { t := 30; return &t }(),
		OrderIndex:      1,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	query := `
		INSERT INTO assessments (id, lesson_id, type, title, description, config, passing_score, time_limit_minutes, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO NOTHING
	`

	configJSON, err := json.Marshal(assessment.Config)
	require.NoError(h.t, err, "Failed to marshal config")

	_, err = h.db.ExecContext(ctx, query,
		assessment.ID,
		assessment.LessonID,
		assessment.Type,
		assessment.Title,
		assessment.Description,
		configJSON,
		assessment.PassingScore,
		30, // time_limit_minutes
		assessment.OrderIndex,
		assessment.CreatedAt,
		assessment.UpdatedAt,
	)
	require.NoError(h.t, err, "Failed to create test assessment")

	return assessment
}

// CreateTestQuestion creates a test question.
func (h *CMSTestHelper) CreateTestQuestion(ctx context.Context, assessmentID string, questionType string) *domain.Question {
	questionID := testutil.RandomString(10)
	now := time.Now()
	explanation := "2+2=4"

	question := &domain.Question{
		ID:            questionID,
		AssessmentID:  assessmentID,
		QuestionType:  domain.QuestionType(questionType),
		QuestionText:  "What is 2+2?",
		Options:       []string{"1", "2", "3", "4"},
		CorrectAnswer: "4",
		Explanation:   &explanation,
		Points:        1,
		OrderIndex:    1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	query := `
		INSERT INTO quiz_questions (id, assessment_id, question_type, question_text, options, correct_answer, explanation, points, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`

	optionsJSON, err := json.Marshal(question.Options)
	require.NoError(h.t, err, "Failed to marshal options")

	_, err = h.db.ExecContext(ctx, query,
		question.ID,
		question.AssessmentID,
		question.QuestionType,
		question.QuestionText,
		optionsJSON,
		question.CorrectAnswer,
		question.Explanation,
		question.Points,
		question.OrderIndex,
	)
	require.NoError(h.t, err, "Failed to create test question")

	return question
}

// CreateTestUser creates a test user with a specific role.
func (h *CMSTestHelper) CreateTestUser(ctx context.Context, role domain.UserRole) *domain.User {
	userID := "user-" + testutil.RandomString(10)
	now := time.Now()

	user := &domain.User{
		ID:          userID,
		FirebaseUID: "firebase-" + userID,
		Email:       testutil.RandomEmail(),
		Username:    "testuser_" + testutil.RandomString(5),
		DisplayName: "Test User",
		PhotoURL:    "https://example.com/photo.jpg",
		Role:        role,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}

	query := `
		INSERT INTO users (id, firebase_uid, email, username, display_name, photo_url, role, is_active, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := h.db.ExecContext(ctx, query,
		user.ID,
		user.FirebaseUID,
		user.Email,
		user.Username,
		user.DisplayName,
		user.PhotoURL,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLoginAt,
	)
	require.NoError(h.t, err, "Failed to create test user")

	return user
}

// CreateTestLesson creates a test lesson.
func (h *CMSTestHelper) CreateTestLesson(ctx context.Context, courseID string) *domain.Lesson {
	lessonID := "lesson-" + testutil.RandomString(10)
	now := time.Now()

	lesson := &domain.Lesson{
		ID:          lessonID,
		CourseID:    courseID,
		Title:       "Test Lesson",
		Description: "Test lesson description",
		Content:     "Test lesson content",
		Order:       1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO lessons (id, course_id, title, description, content, "order", created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := h.db.ExecContext(ctx, query,
		lesson.ID,
		lesson.CourseID,
		lesson.Title,
		lesson.Description,
		lesson.Content,
		lesson.Order,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)
	require.NoError(h.t, err, "Failed to create test lesson")

	return lesson
}

// CreateTestCourse creates a test course.
func (h *CMSTestHelper) CreateTestCourse(ctx context.Context) *domain.Course {
	courseID := "course-" + testutil.RandomString(10)
	now := time.Now()

	course := &domain.Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "Test course description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO courses (id, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := h.db.ExecContext(ctx, query,
		course.ID,
		course.Title,
		course.Description,
		course.CreatedAt,
		course.UpdatedAt,
	)
	require.NoError(h.t, err, "Failed to create test course")

	return course
}

// TruncateCMSTables truncates all CMS-related tables for clean test state.
func (h *CMSTestHelper) TruncateCMSTables(ctx context.Context) {
	tables := []string{
		"submission_comments",
		"submissions",
		"quiz_questions",
		"assessments",
		"content_versions",
		"lessons",
		"courses",
	}

	for _, table := range tables {
		query := "TRUNCATE TABLE " + table + " CASCADE"
		_, err := h.db.ExecContext(ctx, query)
		require.NoError(h.t, err, "Failed to truncate table %s", table)
	}
}

// AssertContentVersionEquals asserts that two content versions are equal.
func (h *CMSTestHelper) AssertContentVersionEquals(expected, actual *domain.ContentVersion) {
	require.Equal(h.t, expected.ID, actual.ID, "ContentVersion ID mismatch")
	require.Equal(h.t, expected.ContentID, actual.ContentID, "ContentID mismatch")
	require.Equal(h.t, expected.VersionNumber, actual.VersionNumber, "VersionNumber mismatch")
	require.Equal(h.t, expected.IsPublished, actual.IsPublished, "IsPublished mismatch")
	require.Equal(h.t, expected.ChangedBy, actual.ChangedBy, "ChangedBy mismatch")

	// Compare content
	require.Equal(h.t, expected.Content, actual.Content, "Content mismatch")
}
