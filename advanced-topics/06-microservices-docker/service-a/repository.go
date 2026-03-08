package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Repository handles data access operations
type Repository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logrus.Logger
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB, redis *redis.Client, log *logrus.Logger) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

// Initialize creates the users table if it doesn't exist
func (r *Repository) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	r.log.Info("Database initialized successfully")
	return nil
}

// Create inserts a new user into the database
func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (id, name, email, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation
				return fmt.Errorf("user with email %s already exists", user.Email)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Invalidate cache
	_ = r.redis.Del(ctx, "user:"+user.ID).Err()
	_ = r.redis.Del(ctx, "users:all").Err()

	r.log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created successfully")

	return nil
}

// FindByID retrieves a user by ID
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
	// Try cache first
	cached, err := r.redis.Get(ctx, "user:"+id).Result()
	if err == nil {
		var user User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			r.log.WithField("user_id", id).Debug("User retrieved from cache")
			return &user, nil
		}
	}

	// Cache miss, query database
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user User
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(user)
	_ = r.redis.Set(ctx, "user:"+id, data, 5*time.Minute).Err()

	r.log.WithField("user_id", id).Debug("User retrieved from database")
	return &user, nil
}

// FindAll retrieves all users
func (r *Repository) FindAll(ctx context.Context) ([]User, error) {
	// Try cache first
	cached, err := r.redis.Get(ctx, "users:all").Result()
	if err == nil {
		var users []User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			r.log.Debug("Users retrieved from cache")
			return users, nil
		}
	}

	// Cache miss, query database
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	// Cache the result
	data, _ := json.Marshal(users)
	_ = r.redis.Set(ctx, "users:all", data, 2*time.Minute).Err()

	r.log.WithField("count", len(users)).Debug("Users retrieved from database")
	return users, nil
}

// Update updates an existing user
func (r *Repository) Update(ctx context.Context, user *User) error {
	query := `
	UPDATE users
	SET name = $2, email = $3, updated_at = $4
	WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Invalidate cache
	_ = r.redis.Del(ctx, "user:"+user.ID).Err()
	_ = r.redis.Del(ctx, "users:all").Err()

	r.log.WithField("user_id", user.ID).Info("User updated successfully")
	return nil
}

// Delete removes a user from the database
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Invalidate cache
	_ = r.redis.Del(ctx, "user:"+id).Err()
	_ = r.redis.Del(ctx, "users:all").Err()

	r.log.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// CheckHealth verifies database and Redis connectivity
func (r *Repository) CheckHealth(ctx context.Context) (dbStatus, redisStatus string) {
	// Check database
	if err := r.db.PingContext(ctx); err != nil {
		r.log.WithError(err).Error("Database health check failed")
		dbStatus = "unhealthy"
	} else {
		dbStatus = "healthy"
	}

	// Check Redis
	if err := r.redis.Ping(ctx).Err(); err != nil {
		r.log.WithError(err).Error("Redis health check failed")
		redisStatus = "unhealthy"
	} else {
		redisStatus = "healthy"
	}

	return
}
