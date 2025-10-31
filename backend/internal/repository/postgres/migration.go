// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"go-pro-backend/pkg/logger"
)

// MigrationV2 represents a database migration with function-based up/down.
type MigrationV2 struct {
	Version     int64
	Description string
	Up          func(*sql.Tx) error
	Down        func(*sql.Tx) error
}

// MigrationManager manages database migrations.
type MigrationManager struct {
	db         *DB
	logger     logger.Logger
	migrations []MigrationV2
}

// NewMigrationManager creates a new migration manager.
func NewMigrationManager(db *DB, logger logger.Logger) *MigrationManager {
	return &MigrationManager{
		db:         db,
		logger:     logger,
		migrations: make([]MigrationV2, 0),
	}
}

// Register registers a migration.
func (mm *MigrationManager) Register(migration MigrationV2) {
	mm.migrations = append(mm.migrations, migration)
}

// RegisterMultiple registers multiple migrations.
func (mm *MigrationManager) RegisterMultiple(migrations []MigrationV2) {
	mm.migrations = append(mm.migrations, migrations...)
}

// ensureMigrationTable creates the migration tracking table if it doesn't exist.
func (mm *MigrationManager) ensureMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			execution_time_ms INTEGER NOT NULL
		)
	`

	if _, err := mm.db.ExecContext(ctx, query); err != nil {
		mm.logger.Error(ctx, "Failed to create migration table", "error", err)
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration versions.
func (mm *MigrationManager) getAppliedMigrations(ctx context.Context) (map[int64]bool, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := mm.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migration rows: %w", err)
	}

	return applied, nil
}

// recordMigration records a migration as applied.
func (mm *MigrationManager) recordMigration(ctx context.Context, tx *sql.Tx, version int64, description string, executionTimeMs int64) error {
	query := `
		INSERT INTO schema_migrations (version, description, execution_time_ms)
		VALUES ($1, $2, $3)
	`

	if _, err := tx.ExecContext(ctx, query, version, description, executionTimeMs); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// removeMigration removes a migration record.
func (mm *MigrationManager) removeMigration(ctx context.Context, tx *sql.Tx, version int64) error {
	query := "DELETE FROM schema_migrations WHERE version = $1"

	if _, err := tx.ExecContext(ctx, query, version); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return nil
}

// Up applies all pending migrations.
func (mm *MigrationManager) Up(ctx context.Context) error {
	mm.logger.Info(ctx, "Starting database migrations")

	// Ensure migration table exists.
	if err := mm.ensureMigrationTable(ctx); err != nil {
		return err
	}

	// Get applied migrations.
	applied, err := mm.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// Sort migrations by version.
	sort.Slice(mm.migrations, func(i, j int) bool {
		return mm.migrations[i].Version < mm.migrations[j].Version
	})

	// Apply pending migrations.
	appliedCount := 0
	for _, migration := range mm.migrations {
		if applied[migration.Version] {
			mm.logger.Debug(ctx, "Migration already applied",
				"version", migration.Version,
				"description", migration.Description)

			continue
		}

		mm.logger.Info(ctx, "Applying migration",
			"version", migration.Version,
			"description", migration.Description)

		start := time.Now()

		// Execute migration in a transaction.
		err := mm.db.WithTransaction(ctx, func(tx *sql.Tx) error {
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration up failed: %w", err)
			}

			executionTimeMs := time.Since(start).Milliseconds()
			if err := mm.recordMigration(ctx, tx, migration.Version, migration.Description, executionTimeMs); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			mm.logger.Error(ctx, "Migration failed",
				"version", migration.Version,
				"description", migration.Description,
				"error", err)

			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		executionTime := time.Since(start)
		mm.logger.Info(ctx, "Migration applied successfully",
			"version", migration.Version,
			"description", migration.Description,
			"execution_time", executionTime)

		appliedCount++
	}

	if appliedCount == 0 {
		mm.logger.Info(ctx, "No pending migrations to apply")
	} else {
		mm.logger.Info(ctx, "Migrations completed successfully",
			"applied_count", appliedCount)
	}

	return nil
}

// Down rolls back the last N migrations.
func (mm *MigrationManager) Down(ctx context.Context, steps int) error {
	mm.logger.Info(ctx, "Rolling back migrations", "steps", steps)

	// Ensure migration table exists.
	if err := mm.ensureMigrationTable(ctx); err != nil {
		return err
	}

	// Get applied migrations in reverse order.
	query := "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT $1"
	rows, err := mm.db.QueryContext(ctx, query, steps)
	if err != nil {
		return fmt.Errorf("failed to query migrations to rollback: %w", err)
	}
	defer rows.Close()

	var versionsToRollback []int64
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		versionsToRollback = append(versionsToRollback, version)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating migration rows: %w", err)
	}

	if len(versionsToRollback) == 0 {
		mm.logger.Info(ctx, "No migrations to rollback")
		return nil
	}

	// Create a map of migrations by version.
	migrationMap := make(map[int64]MigrationV2)
	for _, migration := range mm.migrations {
		migrationMap[migration.Version] = migration
	}

	// Rollback migrations.
	rolledBackCount := 0
	for _, version := range versionsToRollback {
		migration, exists := migrationMap[version]
		if !exists {
			mm.logger.Warn(ctx, "Migration not found in code",
				"version", version)

			continue
		}

		if migration.Down == nil {
			mm.logger.Warn(ctx, "Migration has no down function",
				"version", version,
				"description", migration.Description)

			continue
		}

		mm.logger.Info(ctx, "Rolling back migration",
			"version", version,
			"description", migration.Description)

		start := time.Now()

		// Execute rollback in a transaction.
		err := mm.db.WithTransaction(ctx, func(tx *sql.Tx) error {
			if err := migration.Down(tx); err != nil {
				return fmt.Errorf("migration down failed: %w", err)
			}

			if err := mm.removeMigration(ctx, tx, version); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			mm.logger.Error(ctx, "Migration rollback failed",
				"version", version,
				"description", migration.Description,
				"error", err)

			return fmt.Errorf("migration %d rollback failed: %w", version, err)
		}

		executionTime := time.Since(start)
		mm.logger.Info(ctx, "Migration rolled back successfully",
			"version", version,
			"description", migration.Description,
			"execution_time", executionTime)

		rolledBackCount++
	}

	mm.logger.Info(ctx, "Migrations rolled back successfully",
		"rolled_back_count", rolledBackCount)

	return nil
}

// Status returns the current migration status.
func (mm *MigrationManager) Status(ctx context.Context) ([]MigrationStatus, error) {
	// Ensure migration table exists.
	if err := mm.ensureMigrationTable(ctx); err != nil {
		return nil, err
	}

	// Get applied migrations.
	applied, err := mm.getAppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}

	// Sort migrations by version.
	sort.Slice(mm.migrations, func(i, j int) bool {
		return mm.migrations[i].Version < mm.migrations[j].Version
	})

	// Build status list.
	statuses := make([]MigrationStatus, 0, len(mm.migrations))
	for _, migration := range mm.migrations {
		status := MigrationStatus{
			Version:     migration.Version,
			Description: migration.Description,
			Applied:     applied[migration.Version],
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// MigrationStatus represents the status of a migration.
type MigrationStatus struct {
	Version     int64  `json:"version"`
	Description string `json:"description"`
	Applied     bool   `json:"applied"`
}
