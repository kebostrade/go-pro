// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"

	"github.com/google/uuid"
)

// CourseRepository implements the CourseRepository interface for PostgreSQL.
type CourseRepository struct {
	db *DB
}

// NewCourseRepository creates a new PostgreSQL course repository.
func NewCourseRepository(db *DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course in the database.
func (r *CourseRepository) Create(ctx context.Context, course *domain.Course) error {
	if course.ID == "" {
		course.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gopro.courses (
			id, slug, title, description, long_description, difficulty,
			estimated_duration_hours, is_published, instructor_id, thumbnail_url,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	course.CreatedAt = now
	course.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		generateSlug(course.Title),
		course.Title,
		course.Description,
		course.Description, // Using description as long_description for now
		"beginner",         // Default difficulty
		0,                  // Default duration
		true,               // Default published
		nil,                // No instructor for now
		nil,                // No thumbnail
		course.CreatedAt,
		course.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("course with title '%s' already exists", course.Title))
		}

		return fmt.Errorf("failed to create course: %w", err)
	}

	return nil
}

// GetByID retrieves a course by its ID.
func (r *CourseRepository) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	query := `
		SELECT c.id, c.title, c.description, c.created_at, c.updated_at,
		       COALESCE(array_agg(l.id) FILTER (WHERE l.id IS NOT NULL), '{}') as lesson_ids
		FROM gopro.courses c
		LEFT JOIN gopro.lessons l ON c.id = l.course_id AND l.status = 'published'
		WHERE c.id = $1
		GROUP BY c.id, c.title, c.description, c.created_at, c.updated_at
	`

	var course domain.Course
	var lessonIDs []string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.CreatedAt,
		&course.UpdatedAt,
		&lessonIDs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
		}

		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	course.Lessons = lessonIDs

	return &course, nil
}

// GetAll retrieves all courses with pagination.
func (r *CourseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error) {
	// Count total courses.
	countQuery := "SELECT COUNT(*) FROM gopro.courses WHERE is_published = true"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT c.id, c.title, c.description, c.created_at, c.updated_at,
		       COALESCE(array_agg(l.id) FILTER (WHERE l.id IS NOT NULL), '{}') as lesson_ids
		FROM gopro.courses c
		LEFT JOIN gopro.lessons l ON c.id = l.course_id AND l.status = 'published'
		WHERE c.is_published = true
		GROUP BY c.id, c.title, c.description, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
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
		return nil, 0, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()

	var courses []*domain.Course
	for rows.Next() {
		var course domain.Course
		var lessonIDs []string

		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.CreatedAt,
			&course.UpdatedAt,
			&lessonIDs,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}

		course.Lessons = lessonIDs
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating courses: %w", err)
	}

	return courses, total, nil
}

// Update updates an existing course.
func (r *CourseRepository) Update(ctx context.Context, course *domain.Course) error {
	query := `
		UPDATE gopro.courses 
		SET title = $2, description = $3, long_description = $4, updated_at = $5
		WHERE id = $1
	`

	course.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.Title,
		course.Description,
		course.Description, // Using description as long_description
		course.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("course with title '%s' already exists", course.Title))
		}

		return fmt.Errorf("failed to update course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", course.ID))
	}

	return nil
}

// Delete deletes a course by ID.
func (r *CourseRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM gopro.courses WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
	}

	return nil
}

// Search searches for courses by title or description.
func (r *CourseRepository) Search(ctx context.Context, query string, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error) {
	searchQuery := `
		SELECT c.id, c.title, c.description, c.created_at, c.updated_at,
		       COALESCE(array_agg(l.id) FILTER (WHERE l.id IS NOT NULL), '{}') as lesson_ids
		FROM gopro.courses c
		LEFT JOIN gopro.lessons l ON c.id = l.course_id AND l.status = 'published'
		WHERE c.is_published = true 
		  AND (c.title ILIKE $1 OR c.description ILIKE $1)
		GROUP BY c.id, c.title, c.description, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
	`

	searchTerm := "%" + query + "%"
	args := []interface{}{searchTerm}

	if pagination != nil {
		limit := pagination.PageSize
		offset := (pagination.Page - 1) * pagination.PageSize
		searchQuery += " LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search courses: %w", err)
	}
	defer rows.Close()

	var courses []*domain.Course
	for rows.Next() {
		var course domain.Course
		var lessonIDs []string

		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.CreatedAt,
			&course.UpdatedAt,
			&lessonIDs,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}

		course.Lessons = lessonIDs
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating course rows: %w", err)
	}

	// Count total matching courses.
	countQuery := `
		SELECT COUNT(*) 
		FROM gopro.courses 
		WHERE is_published = true 
		  AND (title ILIKE $1 OR description ILIKE $1)
	`
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return courses, total, nil
}

// generateSlug creates a URL-friendly slug from a title.
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ":", "")
	slug = strings.ReplaceAll(slug, "'", "")

	return slug
}
