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

// ExerciseRepository implements the ExerciseRepository interface for PostgreSQL.
type ExerciseRepository struct {
	db *DB
}

// NewExerciseRepository creates a new PostgreSQL exercise repository.
func NewExerciseRepository(db *DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// Create creates a new exercise in the database.
func (r *ExerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	if exercise.ID == "" {
		exercise.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gopro.exercises (
			id, lesson_id, slug, title, description, instructions, starter_code,
			solution_code, test_cases, difficulty, exercise_order, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	now := time.Now()
	exercise.CreatedAt = now
	exercise.UpdatedAt = now

	// Get the next exercise order for this lesson.
	order, err := r.getNextExerciseOrder(ctx, exercise.LessonID)
	if err != nil {
		return fmt.Errorf("failed to get next exercise order: %w", err)
	}

	// Convert difficulty to string.
	difficultyStr := string(exercise.Difficulty)
	if difficultyStr == "" {
		difficultyStr = "beginner"
	}

	_, err = r.db.ExecContext(ctx, query,
		exercise.ID,
		exercise.LessonID,
		generateSlug(exercise.Title),
		exercise.Title,
		exercise.Description,
		"Complete the exercise", // Default instructions
		"",                      // Default starter code
		"",                      // Default solution code
		"[]",                    // Default empty test cases JSON
		difficultyStr,
		order,
		exercise.CreatedAt,
		exercise.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("exercise with title '%s' already exists in this lesson", exercise.Title))
		}

		return fmt.Errorf("failed to create exercise: %w", err)
	}

	return nil
}

// GetByID retrieves an exercise by its ID.
func (r *ExerciseRepository) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	query := `
		SELECT id, lesson_id, title, description, difficulty, exercise_order, created_at, updated_at
		FROM gopro.exercises
		WHERE id = $1
	`

	var exercise domain.Exercise
	var difficultyStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.LessonID,
		&exercise.Title,
		&exercise.Description,
		&difficultyStr,
		&exercise.TestCases, // Using TestCases field to store order temporarily
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", id))
		}

		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// Convert difficulty string to enum.
	exercise.Difficulty = domain.Difficulty(difficultyStr)

	// Set a default test cases count (we'll improve this later)
	exercise.TestCases = 5

	return &exercise, nil
}

// GetByLessonID retrieves all exercises for a specific lesson.
func (r *ExerciseRepository) GetByLessonID(
	ctx context.Context,
	lessonID string,
	pagination *domain.PaginationRequest,
) ([]*domain.Exercise, int64, error) {
	// Count total exercises for the lesson.
	countQuery := "SELECT COUNT(*) FROM gopro.exercises WHERE lesson_id = $1"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, lessonID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count exercises: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT id, lesson_id, title, description, difficulty, exercise_order, created_at, updated_at
		FROM gopro.exercises
		WHERE lesson_id = $1
		ORDER BY exercise_order ASC
	`

	args := []interface{}{lessonID}
	if pagination != nil {
		limit := pagination.PageSize
		offset := (pagination.Page - 1) * pagination.PageSize
		query += " LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	var exercises []*domain.Exercise
	for rows.Next() {
		var exercise domain.Exercise
		var difficultyStr string

		err := rows.Scan(
			&exercise.ID,
			&exercise.LessonID,
			&exercise.Title,
			&exercise.Description,
			&difficultyStr,
			&exercise.TestCases, // Using TestCases field to store order temporarily
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan exercise: %w", err)
		}

		// Convert difficulty string to enum.
		exercise.Difficulty = domain.Difficulty(difficultyStr)

		// Set a default test cases count.
		exercise.TestCases = 5

		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating exercises: %w", err)
	}

	return exercises, total, nil
}

// Update updates an existing exercise.
func (r *ExerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	query := `
		UPDATE gopro.exercises 
		SET title = $2, description = $3, difficulty = $4, updated_at = $5
		WHERE id = $1
	`

	exercise.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		exercise.ID,
		exercise.Title,
		exercise.Description,
		string(exercise.Difficulty),
		exercise.UpdatedAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError(fmt.Sprintf("exercise with title '%s' already exists in this lesson", exercise.Title))
		}

		return fmt.Errorf("failed to update exercise: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", exercise.ID))
	}

	return nil
}

// Delete deletes an exercise by ID.
func (r *ExerciseRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM gopro.exercises WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", id))
	}

	return nil
}

// GetAll retrieves all exercises with pagination.
func (r *ExerciseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error) {
	// Count total exercises.
	countQuery := "SELECT COUNT(*) FROM gopro.exercises"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count exercises: %w", err)
	}

	// Build query with pagination.
	query := `
		SELECT id, lesson_id, title, description, difficulty, exercise_order, created_at, updated_at
		FROM gopro.exercises
		ORDER BY created_at DESC
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
		return nil, 0, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	var exercises []*domain.Exercise
	for rows.Next() {
		var exercise domain.Exercise
		var difficultyStr string

		err := rows.Scan(
			&exercise.ID,
			&exercise.LessonID,
			&exercise.Title,
			&exercise.Description,
			&difficultyStr,
			&exercise.TestCases, // Using TestCases field to store order temporarily
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan exercise: %w", err)
		}

		// Convert difficulty string to enum.
		exercise.Difficulty = domain.Difficulty(difficultyStr)

		// Set a default test cases count.
		exercise.TestCases = 5

		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating exercises: %w", err)
	}

	return exercises, total, nil
}

// getNextExerciseOrder gets the next exercise order number for a lesson.
func (r *ExerciseRepository) getNextExerciseOrder(ctx context.Context, lessonID string) (int, error) {
	query := "SELECT COALESCE(MAX(exercise_order), 0) + 1 FROM gopro.exercises WHERE lesson_id = $1"

	var nextOrder int
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(&nextOrder)
	if err != nil {
		return 0, fmt.Errorf("failed to get next exercise order: %w", err)
	}

	return nextOrder, nil
}
