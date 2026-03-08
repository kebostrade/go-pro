// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides PostgreSQL repository implementations for assessments.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"

	"github.com/google/uuid"
)

// AssessmentRepository implements the AssessmentRepository interface for PostgreSQL.
type AssessmentRepository struct {
	db *DB
}

// NewAssessmentRepository creates a new PostgreSQL assessment repository.
func NewAssessmentRepository(db *DB) *AssessmentRepository {
	return &AssessmentRepository{db: db}
}

// Create creates a new assessment in the database.
func (r *AssessmentRepository) Create(ctx context.Context, assessment *domain.Assessment) error {
	if assessment.ID == "" {
		assessment.ID = uuid.New().String()
	}

	query := `
		INSERT INTO assessments (
			id, lesson_id, type, title, description, config,
			passing_score, time_limit_minutes, order_index, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	now := time.Now()
	assessment.CreatedAt = now
	assessment.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		assessment.ID,
		assessment.LessonID,
		assessment.Type,
		assessment.Title,
		assessment.Description,
		assessment.Config,
		assessment.PassingScore,
		assessment.TimeLimitMinutes,
		assessment.OrderIndex,
	).Scan(&assessment.CreatedAt, &assessment.UpdatedAt)

	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError("assessment with this order already exists in this lesson")
		}
		return fmt.Errorf("failed to create assessment: %w", err)
	}

	return nil
}

// GetByID retrieves an assessment by its ID.
func (r *AssessmentRepository) GetByID(ctx context.Context, id string) (*domain.Assessment, error) {
	query := `
		SELECT id, lesson_id, type, title, description, config,
			   passing_score, time_limit_minutes, order_index, created_at, updated_at
		FROM assessments
		WHERE id = $1
	`

	var assessment domain.Assessment
	var configJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&assessment.ID,
		&assessment.LessonID,
		&assessment.Type,
		&assessment.Title,
		&assessment.Description,
		&configJSON,
		&assessment.PassingScore,
		&assessment.TimeLimitMinutes,
		&assessment.OrderIndex,
		&assessment.CreatedAt,
		&assessment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("assessment not found")
		}
		return nil, fmt.Errorf("failed to get assessment: %w", err)
	}

	if err := json.Unmarshal(configJSON, &assessment.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assessment config: %w", err)
	}

	return &assessment, nil
}

// GetByLessonID retrieves all assessments for a specific lesson.
func (r *AssessmentRepository) GetByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) ([]*domain.Assessment, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM assessments WHERE lesson_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, lessonID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count assessments: %w", err)
	}

	// Get paginated assessments
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, lesson_id, type, title, description, config,
			   passing_score, time_limit_minutes, order_index, created_at, updated_at
		FROM assessments
		WHERE lesson_id = $1
		ORDER BY order_index ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, lessonID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query assessments: %w", err)
	}
	defer rows.Close()

	var assessments []*domain.Assessment
	for rows.Next() {
		var assessment domain.Assessment
		var configJSON []byte

		if err := rows.Scan(
			&assessment.ID,
			&assessment.LessonID,
			&assessment.Type,
			&assessment.Title,
			&assessment.Description,
			&configJSON,
			&assessment.PassingScore,
			&assessment.TimeLimitMinutes,
			&assessment.OrderIndex,
			&assessment.CreatedAt,
			&assessment.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan assessment: %w", err)
		}

		if err := json.Unmarshal(configJSON, &assessment.Config); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal assessment config: %w", err)
		}

		assessments = append(assessments, &assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating assessments: %w", err)
	}

	return assessments, total, nil
}

// GetAll retrieves all assessments with pagination.
func (r *AssessmentRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Assessment, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM assessments`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count assessments: %w", err)
	}

	// Get paginated assessments
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, lesson_id, type, title, description, config,
			   passing_score, time_limit_minutes, order_index, created_at, updated_at
		FROM assessments
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query assessments: %w", err)
	}
	defer rows.Close()

	var assessments []*domain.Assessment
	for rows.Next() {
		var assessment domain.Assessment
		var configJSON []byte

		if err := rows.Scan(
			&assessment.ID,
			&assessment.LessonID,
			&assessment.Type,
			&assessment.Title,
			&assessment.Description,
			&configJSON,
			&assessment.PassingScore,
			&assessment.TimeLimitMinutes,
			&assessment.OrderIndex,
			&assessment.CreatedAt,
			&assessment.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan assessment: %w", err)
		}

		if err := json.Unmarshal(configJSON, &assessment.Config); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal assessment config: %w", err)
		}

		assessments = append(assessments, &assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating assessments: %w", err)
	}

	return assessments, total, nil
}

// Update updates an existing assessment.
func (r *AssessmentRepository) Update(ctx context.Context, assessment *domain.Assessment) error {
	query := `
		UPDATE assessments
		SET type = $2, title = $3, description = $4, config = $5,
			passing_score = $6, time_limit_minutes = $7, order_index = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		assessment.ID,
		assessment.Type,
		assessment.Title,
		assessment.Description,
		assessment.Config,
		assessment.PassingScore,
		assessment.TimeLimitMinutes,
		assessment.OrderIndex,
	).Scan(&assessment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError("assessment not found")
		}
		if IsUniqueViolation(err) {
			return errors.NewConflictError("assessment with this order already exists")
		}
		return fmt.Errorf("failed to update assessment: %w", err)
	}

	return nil
}

// Delete deletes an assessment by its ID.
func (r *AssessmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM assessments WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete assessment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("assessment not found")
	}

	return nil
}
