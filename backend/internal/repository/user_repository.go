// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package repository provides data access layer implementations.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// User represents a user entity in the database.
type User struct {
	ID           string     `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Roles        []string   `db:"roles" json:"roles"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	IsVerified   bool       `db:"is_verified" json:"is_verified"`
}

// PostgresUserRepository implements UserRepository using PostgreSQL.
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository.
func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, roles, created_at, updated_at, last_login_at, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Roles,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLoginAt,
		user.IsActive,
		user.IsVerified,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, roles, created_at, updated_at, last_login_at, is_active, is_verified
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, roles, created_at, updated_at, last_login_at, is_active, is_verified
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates an existing user in the database.
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, roles = $4, updated_at = $5, last_login_at = $6, is_active = $7, is_verified = $8
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Roles,
		user.UpdatedAt,
		user.LastLoginAt,
		user.IsActive,
		user.IsVerified,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user from the database.
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp for a user.
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, id string, loginTime time.Time) error {
	query := `UPDATE users SET last_login_at = $2, updated_at = $3 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, loginTime, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultDatabaseConfig returns default database configuration.
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "gopro",
		Password:        "",
		Database:        "gopro",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// NewDatabaseConnection creates a new database connection with proper pooling configuration.
func NewDatabaseConnection(config *DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return db, nil
}

// DatabaseUserStore adapts PostgresUserRepository to the security.UserStore interface.
type DatabaseUserStore struct {
	repo *PostgresUserRepository
}

// NewDatabaseUserStore creates a new database-backed user store.
func NewDatabaseUserStore(db *sqlx.DB) *DatabaseUserStore {
	return &DatabaseUserStore{
		repo: NewPostgresUserRepository(db),
	}
}

// CreateUser creates a new user in the database.
func (s *DatabaseUserStore) CreateUser(user *User) error {
	return s.repo.Create(context.Background(), user)
}

// GetUserByEmail retrieves a user by email.
func (s *DatabaseUserStore) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(context.Background(), email)
}

// GetUserByID retrieves a user by ID.
func (s *DatabaseUserStore) GetUserByID(id string) (*User, error) {
	return s.repo.GetByID(context.Background(), id)
}

// UpdateUser updates a user in the database.
func (s *DatabaseUserStore) UpdateUser(user *User) error {
	return s.repo.Update(context.Background(), user)
}

// DeleteUser deletes a user from the database.
func (s *DatabaseUserStore) DeleteUser(id string) error {
	return s.repo.Delete(context.Background(), id)
}
