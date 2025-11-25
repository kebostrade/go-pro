// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package testutil provides functionality for the GO-PRO Learning Platform.
package testutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/repository/postgres"
	"go-pro-backend/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const descriptionFormat = "Description for %s"

// TestDB represents a test database connection.
type TestDB struct {
	*postgres.DB
	t *testing.T
}

// NewTestDB creates a new test database connection.
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()
	config := &postgres.Config{
		Host: GetEnv("TEST_DB_HOST", "localhost"),
		Port: GetEnvAsInt("TEST_DB_PORT", 5432),
		User: GetEnv("TEST_DB_USER", "gopro_test"),
		// #nosec G101: Test environment password - uses env var with safe default
		Password:        GetEnv("TEST_DB_PASSWORD", "gopro_test"),
		Database:        GetEnv("TEST_DB_NAME", "gopro_test"),
		SSLMode:         "disable",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	db, err := postgres.NewConnection(config)
	require.NoError(t, err, "Failed to connect to test database")

	return &TestDB{
		DB: db,
		t:  t,
	}
}

// Cleanup closes the database connection and cleans up test data.
func (tdb *TestDB) Cleanup() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// TruncateTables truncates all tables for a clean test state.
func (tdb *TestDB) TruncateTables(ctx context.Context, tables ...string) {
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		_, err := tdb.ExecContext(ctx, query)
		require.NoError(tdb.t, err, "Failed to truncate table %s", table)
	}
}

// SeedData seeds test data into the database.
func (tdb *TestDB) SeedData(ctx context.Context, data interface{}) error {
	// Implementation depends on data type.
	// This is a placeholder for custom seeding logic.
	return nil
}

// TestLogger creates a test logger that doesn't output during tests.
type TestLogger struct {
	t *testing.T
}

// NewTestLogger creates a new test logger.
func NewTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	return &TestLogger{t: t}
}

func (l *TestLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// Silent in tests unless verbose mode.
}

func (l *TestLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// Silent in tests unless verbose mode.
}

func (l *TestLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.t.Logf("WARN: %s %v", msg, keysAndValues)
}

func (l *TestLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.t.Logf("ERROR: %s %v", msg, keysAndValues)
}

func (l *TestLogger) With(keysAndValues ...interface{}) logger.Logger {
	return l
}

// AssertError asserts that an error occurred.
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Error(t, err, msgAndArgs...)
}

// AssertNoError asserts that no error occurred.
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NoError(t, err, msgAndArgs...)
}

// AssertEqual asserts that two values are equal.
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertNotNil asserts that a value is not nil.
func AssertNotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotNil(t, object, msgAndArgs...)
}

// AssertNil asserts that a value is nil.
func AssertNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Nil(t, object, msgAndArgs...)
}

// RequireNoError requires that no error occurred.
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// RequireEqual requires that two values are equal.
func RequireEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	require.Equal(t, expected, actual, msgAndArgs...)
}

// CreateTestCourse creates a test course.
func CreateTestCourse(id, title string) *domain.Course {
	return &domain.Course{
		ID:          id,
		Title:       title,
		Description: fmt.Sprintf(descriptionFormat, title),
		Lessons:     []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestLesson creates a test lesson.
func CreateTestLesson(id, courseID, title string, order int) *domain.Lesson {
	return &domain.Lesson{
		ID:          id,
		CourseID:    courseID,
		Title:       title,
		Description: fmt.Sprintf(descriptionFormat, title),
		Content:     fmt.Sprintf("Content for %s", title),
		Order:       order,
		Exercises:   []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestExercise creates a test exercise.
func CreateTestExercise(id, lessonID, title string) *domain.Exercise {
	return &domain.Exercise{
		ID:          id,
		LessonID:    lessonID,
		Title:       title,
		Description: fmt.Sprintf(descriptionFormat, title),
		TestCases:   5,
		Difficulty:  domain.DifficultyIntermediate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestProgress creates a test progress record.
func CreateTestProgress(id, userID, lessonID string, status domain.Status) *domain.Progress {
	progress := &domain.Progress{
		ID:        id,
		UserID:    userID,
		LessonID:  lessonID,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if status == domain.StatusCompleted {
		now := time.Now()
		progress.CompletedAt = &now
	}

	return progress
}

// CreateTestUser creates a test user.
// Security Note: This is a test fixture with a dummy bcrypt hash for unit testing only.
// Real passwords are never stored in plain text; this hash is not used in production.
func CreateTestUser(id, username, email string) *domain.User {
	return &domain.User{
		ID:       id,
		Username: username,
		Email:    email,
		// #nosec G101: Test fixture only - not real credentials, for test environment only
		PasswordHash: "$2a$10$test.hash.value",
		FirstName:    "Test",
		LastName:     "User",
		Roles:        []string{"student"},
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// WithTransaction executes a function within a test transaction that is rolled back.
func WithTransaction(t *testing.T, db *sql.DB, fn func(*sql.Tx) error) {
	t.Helper()
	tx, err := db.Begin()
	require.NoError(t, err, "Failed to begin transaction")

	defer func() {
		// Always rollback test transactions.
		rbErr := tx.Rollback()
		if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			t.Logf("Failed to rollback transaction: %v", rbErr)
		}
	}()

	err = fn(tx)
	require.NoError(t, err, "Test function failed")
}

// GetEnv gets an environment variable with a default value.
func GetEnv(key, defaultValue string) string {
	// Implementation would check os.Getenv.
	return defaultValue
}

// GetEnvAsInt gets an environment variable as int with a default value.
func GetEnvAsInt(key string, defaultValue int) int {
	// Implementation would check os.Getenv and convert.
	return defaultValue
}

// WaitForCondition waits for a condition to be true with timeout.
func WaitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}

		<-ticker.C
		if time.Now().After(deadline) {
			t.Fatal("Timeout waiting for condition")
		}
	}
}

// RandomString generates a random string for testing.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}

	return string(b)
}

// RandomEmail generates a random email for testing.
func RandomEmail() string {
	return fmt.Sprintf("test_%s@example.com", RandomString(10))
}
