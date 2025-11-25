// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

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

// UserRepository implements the UserRepository interface for PostgreSQL.
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	query := `
		INSERT INTO gopro.users (
			id, firebase_uid, username, email, display_name, photo_url,
			role, is_active, created_at, updated_at, last_login_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Generate username from email if not provided
	username := user.Username
	if username == "" {
		username = generateUsernameFromEmail(user.Email)
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.FirebaseUID,
		username,
		user.Email,
		user.DisplayName,
		nullString(user.PhotoURL),
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLoginAt,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError("user with this Firebase UID or email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, username, email, display_name, photo_url,
		       role, is_active, created_at, updated_at, last_login_at
		FROM gopro.users
		WHERE id = $1
	`

	user := &domain.User{}
	var photoURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&photoURL,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// GetByFirebaseUID retrieves a user by their Firebase UID.
func (r *UserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, username, email, display_name, photo_url,
		       role, is_active, created_at, updated_at, last_login_at
		FROM gopro.users
		WHERE firebase_uid = $1
	`

	user := &domain.User{}
	var photoURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, firebaseUID).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&photoURL,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with Firebase UID %s not found", firebaseUID))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// GetByEmail retrieves a user by their email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, username, email, display_name, photo_url,
		       role, is_active, created_at, updated_at, last_login_at
		FROM gopro.users
		WHERE email = $1
	`

	user := &domain.User{}
	var photoURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&photoURL,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// GetAll retrieves all users with pagination.
func (r *UserRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error) {
	// Count total users
	var total int64
	countQuery := "SELECT COUNT(*) FROM gopro.users"
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated users
	query := `
		SELECT id, firebase_uid, username, email, display_name, photo_url,
		       role, is_active, created_at, updated_at, last_login_at
		FROM gopro.users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	limit := pagination.PageSize
	offset := (pagination.Page - 1) * pagination.PageSize

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var photoURL sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.FirebaseUID,
			&user.Username,
			&user.Email,
			&user.DisplayName,
			&photoURL,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&lastLoginAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		if photoURL.Valid {
			user.PhotoURL = photoURL.String
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// Update updates a user in the database.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE gopro.users
		SET username = $1, email = $2, display_name = $3, photo_url = $4,
		    role = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.DisplayName,
		nullString(user.PhotoURL),
		user.Role,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		if IsUniqueViolation(err) {
			return errors.NewConflictError("email already in use by another user")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", user.ID))
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user.
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE gopro.users
		SET last_login_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", userID))
	}

	return nil
}

// Delete deletes a user from the database.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM gopro.users WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	return nil
}

// Helper functions

// nullString converts an empty string to sql.NullString.
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// generateUsernameFromEmail generates a username from an email address.
func generateUsernameFromEmail(email string) string {
	// Split email at @ and take the first part
	parts := []rune(email)
	for i, char := range parts {
		if char == '@' {
			return string(parts[:i])
		}
	}
	return email
}
