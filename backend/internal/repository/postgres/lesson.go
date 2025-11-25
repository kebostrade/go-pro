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

// LessonRepository implements the LessonRepository interface for PostgreSQL.
type LessonRepository struct {
	db *DB
}

// NewLessonRepository creates a new PostgreSQL lesson repository.
func NewLessonRepository(db *DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create creates a new lesson in the database.
func (r *LessonRepository) Create(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.ID == "" {
		lesson.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gopro.lessons (
			id, course_id, slug, title, description, content, lesson_order,
			estimated_duration_minutes, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	lesson.CreatedAt = now
	lesson.UpdatedAt = now

	// Get the next lesson order for this course.
	order, err := r.getNextLessonOrder(ctx, lesson.CourseID)
	if err != nil {
		return fmt.Errorf("failed to get next lesson order: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		lesson.ID,
		lesson.CourseID,
		generateSlug(lesson.Title),
		lesson.Title,
		lesson.Description,
		lesson.Content,
		order,
		60, // Default 60 minutes
		"published",
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("lesson with title '%s' already exists in this course", lesson.Title))
		}

		return fmt.Errorf("failed to create lesson: %w", err)
	}

	return nil
}

// GetByID retrieves a lesson by its ID.
func (r *LessonRepository) GetByID(ctx context.Context, id string) (*domain.Lesson, error) {
	query := `
		SELECT l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at,
		       COALESCE(array_agg(e.id) FILTER (WHERE e.id IS NOT NULL), '{}') as exercise_ids
		FROM gopro.lessons l
		LEFT JOIN gopro.exercises e ON l.id = e.lesson_id
		WHERE l.id = $1 AND l.status = 'published'
		GROUP BY l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at
	`

	var lesson domain.Lesson
	var exerciseIDs []string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		&lesson.Description,
		&lesson.Content,
		&lesson.Order,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
		&exerciseIDs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", id))
		}

		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	lesson.Exercises = exerciseIDs

	return &lesson, nil
}

// GetByCourseID retrieves all lessons for a specific course.
func (r *LessonRepository) GetByCourseID(
	ctx context.Context,
	courseID string,
	pagination *domain.PaginationRequest,
) ([]*domain.Lesson, int64, error) {
	// Count total lessons for the course.
	countQuery := "SELECT COUNT(*) FROM gopro.lessons WHERE course_id = $1 AND status = 'published'"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at,
		       COALESCE(array_agg(e.id) FILTER (WHERE e.id IS NOT NULL), '{}') as exercise_ids
		FROM gopro.lessons l
		LEFT JOIN gopro.exercises e ON l.id = e.lesson_id
		WHERE l.course_id = $1 AND l.status = 'published'
		GROUP BY l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at
		ORDER BY l.lesson_order ASC
	`

	args := []interface{}{courseID}
	if pagination != nil {
		limit := pagination.PageSize
		offset := (pagination.Page - 1) * pagination.PageSize
		query += " LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query lessons: %w", err)
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		var lesson domain.Lesson
		var exerciseIDs []string

		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Content,
			&lesson.Order,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
			&exerciseIDs,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan lesson: %w", err)
		}

		lesson.Exercises = exerciseIDs
		lessons = append(lessons, &lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating lessons: %w", err)
	}

	return lessons, total, nil
}

// Update updates an existing lesson.
func (r *LessonRepository) Update(ctx context.Context, lesson *domain.Lesson) error {
	query := `
		UPDATE gopro.lessons 
		SET title = $2, description = $3, content = $4, updated_at = $5
		WHERE id = $1
	`

	lesson.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		lesson.ID,
		lesson.Title,
		lesson.Description,
		lesson.Content,
		lesson.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("lesson with title '%s' already exists in this course", lesson.Title))
		}

		return fmt.Errorf("failed to update lesson: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", lesson.ID))
	}

	return nil
}

// Delete deletes a lesson by ID.
func (r *LessonRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM gopro.lessons WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("lesson with id %s not found", id))
	}

	return nil
}

// GetAll retrieves all lessons with pagination.
func (r *LessonRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Lesson, int64, error) {
	// Count total lessons.
	countQuery := "SELECT COUNT(*) FROM gopro.lessons WHERE status = 'published'"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at,
		       COALESCE(array_agg(e.id) FILTER (WHERE e.id IS NOT NULL), '{}') as exercise_ids
		FROM gopro.lessons l
		LEFT JOIN gopro.exercises e ON l.id = e.lesson_id
		WHERE l.status = 'published'
		GROUP BY l.id, l.course_id, l.title, l.description, l.content, l.lesson_order, l.created_at, l.updated_at
		ORDER BY l.created_at DESC
	`

	args := []interface{}{}
	if pagination != nil {
		limit := pagination.PageSize
		offset := (pagination.Page - 1) * pagination.PageSize
		query += " LIMIT $1 OFFSET $2"
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query lessons: %w", err)
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		var lesson domain.Lesson
		var exerciseIDs []string

		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Content,
			&lesson.Order,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
			&exerciseIDs,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan lesson: %w", err)
		}

		lesson.Exercises = exerciseIDs
		lessons = append(lessons, &lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating lessons: %w", err)
	}

	return lessons, total, nil
}

// getNextLessonOrder gets the next lesson order number for a course.
func (r *LessonRepository) getNextLessonOrder(ctx context.Context, courseID string) (int, error) {
	query := "SELECT COALESCE(MAX(lesson_order), 0) + 1 FROM gopro.lessons WHERE course_id = $1"

	var nextOrder int
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&nextOrder)
	if err != nil {
		return 0, fmt.Errorf("failed to get next lesson order: %w", err)
	}

	return nextOrder, nil
}
