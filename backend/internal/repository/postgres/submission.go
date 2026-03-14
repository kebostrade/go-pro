// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides PostgreSQL repository implementations for submissions.
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

// SubmissionRepository implements the SubmissionRepository interface for PostgreSQL.
type SubmissionRepository struct {
	db *DB
}

// NewSubmissionRepository creates a new PostgreSQL submission repository.
func NewSubmissionRepository(db *DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// Create creates a new submission in the database.
func (r *SubmissionRepository) Create(ctx context.Context, submission *domain.Submission) error {
	if submission.ID == "" {
		submission.ID = uuid.New().String()
	}

	if submission.SubmittedAt.IsZero() {
		submission.SubmittedAt = time.Now()
	}

	query := `
		INSERT INTO submissions (
			id, assessment_id, user_id, content, score, feedback,
			graded_by, graded_at, submitted_at, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	contentJSON, err := json.Marshal(submission.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal submission content: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		submission.ID,
		submission.AssessmentID,
		submission.UserID,
		contentJSON,
		submission.Score,
		submission.Feedback,
		submission.GradedBy,
		submission.GradedAt,
		submission.SubmittedAt,
		submission.Status,
	)

	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError("submission already exists for this user and assessment")
		}
		return fmt.Errorf("failed to create submission: %w", err)
	}

	return nil
}

// GetByID retrieves a submission by its ID.
func (r *SubmissionRepository) GetByID(ctx context.Context, id string) (*domain.Submission, error) {
	query := `
		SELECT id, assessment_id, user_id, content, score, feedback,
			   graded_by, graded_at, submitted_at, status
		FROM submissions
		WHERE id = $1
	`

	var submission domain.Submission
	var contentJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&submission.ID,
		&submission.AssessmentID,
		&submission.UserID,
		&contentJSON,
		&submission.Score,
		&submission.Feedback,
		&submission.GradedBy,
		&submission.GradedAt,
		&submission.SubmittedAt,
		&submission.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	if err := json.Unmarshal(contentJSON, &submission.Content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal submission content: %w", err)
	}

	return &submission, nil
}

// GetByUserID retrieves all submissions for a specific user.
func (r *SubmissionRepository) GetByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM submissions WHERE user_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Get paginated submissions
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, assessment_id, user_id, content, score, feedback,
			   graded_by, graded_at, submitted_at, status
		FROM submissions
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*domain.Submission
	for rows.Next() {
		submission, err := r.scanSubmission(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, total, nil
}

// GetByAssessmentID retrieves all submissions for a specific assessment.
func (r *SubmissionRepository) GetByAssessmentID(ctx context.Context, assessmentID string, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM submissions WHERE assessment_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, assessmentID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Get paginated submissions
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, assessment_id, user_id, content, score, feedback,
			   graded_by, graded_at, submitted_at, status
		FROM submissions
		WHERE assessment_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, assessmentID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*domain.Submission
	for rows.Next() {
		submission, err := r.scanSubmission(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, total, nil
}

// GetByUserAndAssessment retrieves a submission by user ID and assessment ID.
func (r *SubmissionRepository) GetByUserAndAssessment(ctx context.Context, userID, assessmentID string) (*domain.Submission, error) {
	query := `
		SELECT id, assessment_id, user_id, content, score, feedback,
			   graded_by, graded_at, submitted_at, status
		FROM submissions
		WHERE user_id = $1 AND assessment_id = $2
	`

	var submission domain.Submission
	var contentJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID, assessmentID).Scan(
		&submission.ID,
		&submission.AssessmentID,
		&submission.UserID,
		&contentJSON,
		&submission.Score,
		&submission.Feedback,
		&submission.GradedBy,
		&submission.GradedAt,
		&submission.SubmittedAt,
		&submission.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	if err := json.Unmarshal(contentJSON, &submission.Content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal submission content: %w", err)
	}

	return &submission, nil
}

// GetAll retrieves all submissions with optional filtering.
func (r *SubmissionRepository) GetAll(ctx context.Context, filters map[string]interface{}, pagination *domain.PaginationRequest) ([]*domain.Submission, int64, error) {
	// Build dynamic query with filters
	baseQuery := `FROM submissions`
	whereClause := ""
	args := []interface{}{}
	argCount := 0

	if len(filters) > 0 {
		whereConditions := []string{}

		if assessmentID, ok := filters["assessment_id"]; ok {
			argCount++
			whereConditions = append(whereConditions, fmt.Sprintf("assessment_id = $%d", argCount))
			args = append(args, assessmentID)
		}

		if userID, ok := filters["user_id"]; ok {
			argCount++
			whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argCount))
			args = append(args, userID)
		}

		if status, ok := filters["status"]; ok {
			argCount++
			whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argCount))
			args = append(args, status)
		}

		if len(whereConditions) > 0 {
			whereClause = " WHERE " + joinConditions(whereConditions, " AND ")
		}
	}

	// Get total count
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Get paginated submissions
	offset := (pagination.Page - 1) * pagination.PageSize
	argCount++
	args = append(args, pagination.PageSize)
	argCount++
	args = append(args, offset)

	query := `
		SELECT id, assessment_id, user_id, content, score, feedback,
			   graded_by, graded_at, submitted_at, status
	` + baseQuery + whereClause + `
		ORDER BY submitted_at DESC
		LIMIT $` + fmt.Sprintf("%d", argCount-1) + ` OFFSET $` + fmt.Sprintf("%d", argCount)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*domain.Submission
	for rows.Next() {
		submission, err := r.scanSubmission(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, total, nil
}

// Update updates an existing submission.
func (r *SubmissionRepository) Update(ctx context.Context, submission *domain.Submission) error {
	query := `
		UPDATE submissions
		SET content = $2, score = $3, feedback = $4, graded_by = $5,
			graded_at = $6, status = $7
		WHERE id = $1
	`

	contentJSON, err := json.Marshal(submission.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal submission content: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		submission.ID,
		contentJSON,
		submission.Score,
		submission.Feedback,
		submission.GradedBy,
		submission.GradedAt,
		submission.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("submission not found")
	}

	return nil
}

// Delete deletes a submission by its ID.
func (r *SubmissionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM submissions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete submission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("submission not found")
	}

	return nil
}

// scanSubmission scans a submission from a database row.
func (r *SubmissionRepository) scanSubmission(rows *sql.Rows) (*domain.Submission, error) {
	var submission domain.Submission
	var contentJSON []byte

	err := rows.Scan(
		&submission.ID,
		&submission.AssessmentID,
		&submission.UserID,
		&contentJSON,
		&submission.Score,
		&submission.Feedback,
		&submission.GradedBy,
		&submission.GradedAt,
		&submission.SubmittedAt,
		&submission.Status,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(contentJSON, &submission.Content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal submission content: %w", err)
	}

	return &submission, nil
}

// joinConditions joins SQL conditions with a separator.
func joinConditions(conditions []string, separator string) string {
	if len(conditions) == 0 {
		return ""
	}
	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result += separator + conditions[i]
	}
	return result
}

// SubmissionCommentRepository implements the SubmissionCommentRepository interface for PostgreSQL.
type SubmissionCommentRepository struct {
	db *DB
}

// NewSubmissionCommentRepository creates a new PostgreSQL submission comment repository.
func NewSubmissionCommentRepository(db *DB) *SubmissionCommentRepository {
	return &SubmissionCommentRepository{db: db}
}

// Create creates a new submission comment.
func (r *SubmissionCommentRepository) Create(ctx context.Context, comment *domain.SubmissionComment) error {
	if comment.ID == "" {
		comment.ID = uuid.New().String()
	}

	query := `
		INSERT INTO submission_comments (id, submission_id, author_id, comment_text, line_number, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING created_at
	`

	now := time.Now()
	comment.CreatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		comment.ID,
		comment.SubmissionID,
		comment.AuthorID,
		comment.CommentText,
		comment.LineNumber,
	).Scan(&comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create submission comment: %w", err)
	}

	return nil
}

// GetBySubmissionID retrieves all comments for a submission.
func (r *SubmissionCommentRepository) GetBySubmissionID(ctx context.Context, submissionID string) ([]*domain.SubmissionComment, error) {
	query := `
		SELECT id, submission_id, author_id, comment_text, line_number, created_at
		FROM submission_comments
		WHERE submission_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, submissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submission comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.SubmissionComment
	for rows.Next() {
		var comment domain.SubmissionComment
		if err := rows.Scan(
			&comment.ID,
			&comment.SubmissionID,
			&comment.AuthorID,
			&comment.CommentText,
			&comment.LineNumber,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan submission comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating submission comments: %w", err)
	}

	return comments, nil
}

// GetByID retrieves a comment by its ID.
func (r *SubmissionCommentRepository) GetByID(ctx context.Context, id string) (*domain.SubmissionComment, error) {
	query := `
		SELECT id, submission_id, author_id, comment_text, line_number, created_at
		FROM submission_comments
		WHERE id = $1
	`

	var comment domain.SubmissionComment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.SubmissionID,
		&comment.AuthorID,
		&comment.CommentText,
		&comment.LineNumber,
		&comment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("submission comment not found")
		}
		return nil, fmt.Errorf("failed to get submission comment: %w", err)
	}

	return &comment, nil
}

// Update updates an existing comment.
func (r *SubmissionCommentRepository) Update(ctx context.Context, comment *domain.SubmissionComment) error {
	query := `
		UPDATE submission_comments
		SET comment_text = $2, line_number = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		comment.ID,
		comment.CommentText,
		comment.LineNumber,
	)

	if err != nil {
		return fmt.Errorf("failed to update submission comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("submission comment not found")
	}

	return nil
}

// Delete deletes a comment by its ID.
func (r *SubmissionCommentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM submission_comments WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete submission comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("submission comment not found")
	}

	return nil
}
