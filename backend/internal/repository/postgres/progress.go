// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"

	"github.com/google/uuid"
)

// ProgressRepository implements the ProgressRepository interface for PostgreSQL.
type ProgressRepository struct {
	db *DB
}

// NewProgressRepository creates a new PostgreSQL progress repository.
func NewProgressRepository(db *DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// Create creates a new progress record in the database.
func (r *ProgressRepository) Create(ctx context.Context, progress *domain.Progress) error {
	if progress.ID == "" {
		progress.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gopro.user_lesson_progress (id, user_id, lesson_id, status, started_at, completed_at, time_spent_seconds, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, lesson_id) 
		DO UPDATE SET
			status = EXCLUDED.status,
			completed_at = EXCLUDED.completed_at,
			time_spent_seconds = user_lesson_progress.time_spent_seconds + EXCLUDED.time_spent_seconds,
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	progress.CreatedAt = now
	progress.UpdatedAt = now

	status := string(progress.Status)
	var completedAt *time.Time
	var startedAt *time.Time

	if progress.Status == domain.StatusCompleted {
		completedAt = progress.CompletedAt
		startedAt = &now
	} else if progress.Status == domain.StatusInProgress {
		startedAt = &now
	}

	_, err := r.db.ExecContext(ctx, query,
		progress.ID,
		progress.UserID,
		progress.LessonID,
		status,
		startedAt,
		completedAt,
		0, // Default time spent
		progress.CreatedAt,
		progress.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create progress: %w", err)
	}

	return nil
}

// GetByID retrieves a progress record by its ID.
func (r *ProgressRepository) GetByID(ctx context.Context, id string) (*domain.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, status, completed_at, time_spent_seconds, created_at, updated_at
		FROM gopro.user_lesson_progress
		WHERE id = $1
	`

	var progress domain.Progress
	var status string
	var completedAt sql.NullTime
	var timeSpent int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&status,
		&completedAt,
		&timeSpent,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", id))
		}

		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Convert status to domain Status.
	progress.Status = domain.Status(status)
	if completedAt.Valid {
		completedAtTime := completedAt.Time
		progress.CompletedAt = &completedAtTime
	}

	return &progress, nil
}

// GetByUserID retrieves all progress records for a specific user.
func (r *ProgressRepository) GetByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) ([]*domain.Progress, int64, error) {
	// Count total progress records for the user.
	countQuery := "SELECT COUNT(*) FROM gopro.user_lesson_progress WHERE user_id = $1"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count progress records: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT id, user_id, lesson_id, status, completed_at, time_spent_seconds, created_at, updated_at
		FROM gopro.user_lesson_progress
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	args := []interface{}{userID}
	if pagination != nil {
		limit := pagination.PageSize
		offset := (pagination.Page - 1) * pagination.PageSize
		query += " LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query progress records: %w", err)
	}
	defer rows.Close()

	var progressRecords []*domain.Progress
	for rows.Next() {
		var progress domain.Progress
		var status string
		var completedAt sql.NullTime
		var timeSpent int

		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&status,
			&completedAt,
			&timeSpent,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan progress: %w", err)
		}

		// Convert status to domain Status.
		progress.Status = domain.Status(status)
		if completedAt.Valid {
			completedAtTime := completedAt.Time
			progress.CompletedAt = &completedAtTime
		}

		progressRecords = append(progressRecords, &progress)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating progress records: %w", err)
	}

	return progressRecords, total, nil
}

// GetByUserAndLesson retrieves a specific progress record for a user and lesson.
func (r *ProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, status, completed_at, time_spent_seconds, created_at, updated_at
		FROM gopro.user_lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	var progress domain.Progress
	var status string
	var completedAt sql.NullTime
	var timeSpent int

	err := r.db.QueryRowContext(ctx, query, userID, lessonID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&status,
		&completedAt,
		&timeSpent,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Sprintf("progress for user %s and lesson %s not found", userID, lessonID))
		}

		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	// Convert status to domain Status.
	progress.Status = domain.Status(status)
	if completedAt.Valid {
		completedAtTime := completedAt.Time
		progress.CompletedAt = &completedAtTime
	}

	return &progress, nil
}

// Update updates an existing progress record.
func (r *ProgressRepository) Update(ctx context.Context, progress *domain.Progress) error {
	query := `
		UPDATE gopro.user_lesson_progress 
		SET status = $2, completed_at = $3, updated_at = $4
		WHERE id = $1
	`

	progress.UpdatedAt = time.Now()

	status := string(progress.Status)
	var completedAt *time.Time

	if progress.Status == domain.StatusCompleted {
		completedAt = progress.CompletedAt
	}

	result, err := r.db.ExecContext(ctx, query,
		progress.ID,
		status,
		completedAt,
		progress.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", progress.ID))
	}

	return nil
}

// Delete deletes a progress record by ID.
func (r *ProgressRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM gopro.user_lesson_progress WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", id))
	}

	return nil
}

// GetUserProgressSummary retrieves a summary of user progress across all courses.
func (r *ProgressRepository) GetUserProgressSummary(ctx context.Context, userID string) (*domain.UserProgressSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT ulp.lesson_id) as lessons_started,
			COUNT(DISTINCT CASE WHEN ulp.status = 'completed' THEN ulp.lesson_id END) as lessons_completed,
			COUNT(DISTINCT l.id) as total_lessons_available,
			COALESCE(AVG(CASE WHEN ulp.status = 'completed' THEN 100 ELSE 0 END), 0) as overall_progress_percentage
		FROM gopro.lessons l
		LEFT JOIN gopro.user_lesson_progress ulp ON l.id = ulp.lesson_id AND ulp.user_id = $1
		WHERE l.status = 'published'
	`

	var summary domain.UserProgressSummary
	var progressPercentage float64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&summary.LessonsStarted,
		&summary.LessonsCompleted,
		&summary.TotalLessons,
		&progressPercentage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress summary: %w", err)
	}

	summary.UserID = userID
	summary.OverallProgress = int(progressPercentage)

	return &summary, nil
}
