// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"go-pro-backend/internal/repository"
)

// Repositories implements the repository.Repositories interface for PostgreSQL.
type Repositories struct {
	db       *DB
	Course   *CourseRepository
	Lesson   *LessonRepository
	Exercise *ExerciseRepository
	Progress *ProgressRepository
	User     *UserRepository
}

// NewRepositories creates a new PostgreSQL repositories instance.
func NewRepositories(db *DB) *Repositories {
	return &Repositories{
		db:       db,
		Course:   NewCourseRepository(db),
		Lesson:   NewLessonRepository(db),
		Exercise: NewExerciseRepository(db),
		Progress: NewProgressRepository(db),
		User:     NewUserRepository(db),
	}
}

// NewRepositoriesFromEnv creates repositories from environment variables.
func NewRepositoriesFromEnv() (*Repositories, error) {
	config := &Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "gopro_user"),
		Password:        getEnv("DB_PASSWORD", ""),
		Database:        getEnv("DB_NAME", "gopro_dev"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	db, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return NewRepositories(db), nil
}

// Close closes the database connection.
func (r *Repositories) Close() error {
	return r.db.Close()
}

// Health checks the health of all repositories.
func (r *Repositories) Health(ctx context.Context) error {
	return r.db.Health(ctx)
}

// BeginTransaction starts a new database transaction.
func (r *Repositories) BeginTransaction(ctx context.Context) (*TransactionalRepositories, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TransactionalRepositories{
		tx:       tx,
		Course:   NewTransactionalCourseRepository(tx),
		Lesson:   NewTransactionalLessonRepository(tx),
		Exercise: NewTransactionalExerciseRepository(tx),
		Progress: NewTransactionalProgressRepository(tx),
	}, nil
}

// TransactionalRepositories provides transactional repository operations.
type TransactionalRepositories struct {
	tx       *sql.Tx
	Course   *TransactionalCourseRepository
	Lesson   *TransactionalLessonRepository
	Exercise *TransactionalExerciseRepository
	Progress *TransactionalProgressRepository
}

// Commit commits the transaction.
func (tr *TransactionalRepositories) Commit() error {
	return tr.tx.Commit()
}

// Rollback rolls back the transaction.
func (tr *TransactionalRepositories) Rollback() error {
	return tr.tx.Rollback()
}

// WithTransaction executes a function within a database transaction.
func (r *Repositories) WithTransaction(ctx context.Context, fn func(*TransactionalRepositories) error) error {
	return r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		txRepos := &TransactionalRepositories{
			tx:       tx,
			Course:   NewTransactionalCourseRepository(tx),
			Lesson:   NewTransactionalLessonRepository(tx),
			Exercise: NewTransactionalExerciseRepository(tx),
			Progress: NewTransactionalProgressRepository(tx),
		}

		return fn(txRepos)
	})
}

// Transactional repository implementations (simplified versions that use sql.Tx).
type TransactionalCourseRepository struct {
	tx *sql.Tx
}

func NewTransactionalCourseRepository(tx *sql.Tx) *TransactionalCourseRepository {
	return &TransactionalCourseRepository{tx: tx}
}

type TransactionalLessonRepository struct {
	tx *sql.Tx
}

func NewTransactionalLessonRepository(tx *sql.Tx) *TransactionalLessonRepository {
	return &TransactionalLessonRepository{tx: tx}
}

type TransactionalExerciseRepository struct {
	tx *sql.Tx
}

func NewTransactionalExerciseRepository(tx *sql.Tx) *TransactionalExerciseRepository {
	return &TransactionalExerciseRepository{tx: tx}
}

type TransactionalProgressRepository struct {
	tx *sql.Tx
}

func NewTransactionalProgressRepository(tx *sql.Tx) *TransactionalProgressRepository {
	return &TransactionalProgressRepository{tx: tx}
}

// RunMigrations runs database migrations.
func (r *Repositories) RunMigrations(ctx context.Context) error {
	migrations := []Migration{
		{
			Version:     1,
			Description: "Initial schema setup",
			Up: `
				-- This migration assumes the schema is already created by init-db.sql
				-- We'll add any additional migrations here as needed
				SELECT 1;
			`,
			Down: `
				-- Rollback would drop tables, but we'll keep it simple for now
				SELECT 1;
			`,
		},
	}

	return r.db.RunMigrations(ctx, migrations)
}

// Helper functions for environment variable parsing.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	return defaultValue
}

// Ensure our repositories implement the interfaces.
var (
	_ repository.CourseRepository   = (*CourseRepository)(nil)
	_ repository.LessonRepository   = (*LessonRepository)(nil)
	_ repository.ExerciseRepository = (*ExerciseRepository)(nil)
	_ repository.ProgressRepository = (*ProgressRepository)(nil)
	_ repository.UserRepository     = (*UserRepository)(nil)
)
