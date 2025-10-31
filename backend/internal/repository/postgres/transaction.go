// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go-pro-backend/pkg/logger"
)

// TransactionManager manages database transactions with advanced features.
type TransactionManager struct {
	db     *DB
	logger logger.Logger
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(db *DB, logger logger.Logger) *TransactionManager {
	return &TransactionManager{
		db:     db,
		logger: logger,
	}
}

// TxOptions represents transaction options.
type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
	Timeout   time.Duration
}

// DefaultTxOptions returns default transaction options.
func DefaultTxOptions() *TxOptions {
	return &TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   30 * time.Second,
	}
}

// Transaction represents a database transaction with savepoint support.
type Transaction struct {
	tx         *sql.Tx
	db         *DB
	logger     logger.Logger
	savepoints []string
	mu         sync.Mutex
	committed  bool
	rolledBack bool
}

// BeginTransaction starts a new transaction with options.
func (tm *TransactionManager) BeginTransaction(ctx context.Context, opts *TxOptions) (*Transaction, error) {
	if opts == nil {
		opts = DefaultTxOptions()
	}

	// Apply timeout if specified.
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	sqlOpts := &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	}

	tx, err := tm.db.BeginTx(ctx, sqlOpts)
	if err != nil {
		tm.logger.Error(ctx, "Failed to begin transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	tm.logger.Debug(ctx, "Transaction started",
		"isolation", opts.Isolation,
		"read_only", opts.ReadOnly)

	return &Transaction{
		tx:         tx,
		db:         tm.db,
		logger:     tm.logger,
		savepoints: make([]string, 0),
	}, nil
}

// Commit commits the transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed {
		return fmt.Errorf("transaction already committed")
	}
	if t.rolledBack {
		return fmt.Errorf("transaction already rolled back")
	}

	if err := t.tx.Commit(); err != nil {
		t.logger.Error(ctx, "Failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	t.committed = true
	t.logger.Debug(ctx, "Transaction committed successfully")

	return nil
}

// Rollback rolls back the transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed {
		return fmt.Errorf("transaction already committed")
	}
	if t.rolledBack {
		return fmt.Errorf("transaction already rolled back")
	}

	if err := t.tx.Rollback(); err != nil {
		t.logger.Error(ctx, "Failed to rollback transaction", "error", err)
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	t.rolledBack = true
	t.logger.Debug(ctx, "Transaction rolled back successfully")

	return nil
}

// Savepoint creates a savepoint within the transaction.
func (t *Transaction) Savepoint(ctx context.Context, name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed || t.rolledBack {
		return fmt.Errorf("transaction is not active")
	}

	query := fmt.Sprintf("SAVEPOINT %s", name)
	if _, err := t.tx.ExecContext(ctx, query); err != nil {
		t.logger.Error(ctx, "Failed to create savepoint", "error", err, "savepoint", name)
		return fmt.Errorf("failed to create savepoint %s: %w", name, err)
	}

	t.savepoints = append(t.savepoints, name)
	t.logger.Debug(ctx, "Savepoint created", "savepoint", name)

	return nil
}

// RollbackToSavepoint rolls back to a specific savepoint.
func (t *Transaction) RollbackToSavepoint(ctx context.Context, name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed || t.rolledBack {
		return fmt.Errorf("transaction is not active")
	}

	// Check if savepoint exists.
	found := false
	for _, sp := range t.savepoints {
		if sp == name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("savepoint %s not found", name)
	}

	query := fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name)
	if _, err := t.tx.ExecContext(ctx, query); err != nil {
		t.logger.Error(ctx, "Failed to rollback to savepoint", "error", err, "savepoint", name)
		return fmt.Errorf("failed to rollback to savepoint %s: %w", name, err)
	}

	t.logger.Debug(ctx, "Rolled back to savepoint", "savepoint", name)

	return nil
}

// ReleaseSavepoint releases a savepoint.
func (t *Transaction) ReleaseSavepoint(ctx context.Context, name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed || t.rolledBack {
		return fmt.Errorf("transaction is not active")
	}

	query := fmt.Sprintf("RELEASE SAVEPOINT %s", name)
	if _, err := t.tx.ExecContext(ctx, query); err != nil {
		t.logger.Error(ctx, "Failed to release savepoint", "error", err, "savepoint", name)
		return fmt.Errorf("failed to release savepoint %s: %w", name, err)
	}

	// Remove savepoint from list.
	for i, sp := range t.savepoints {
		if sp == name {
			t.savepoints = append(t.savepoints[:i], t.savepoints[i+1:]...)
			break
		}
	}

	t.logger.Debug(ctx, "Savepoint released", "savepoint", name)

	return nil
}

// Exec executes a query within the transaction.
func (t *Transaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows.
func (t *Transaction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns at most one row.
func (t *Transaction) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// Prepare creates a prepared statement within the transaction.
func (t *Transaction) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}

// WithTransaction executes a function within a transaction with automatic rollback on error.
func (tm *TransactionManager) WithTransaction(ctx context.Context, opts *TxOptions, fn func(*Transaction) error) error {
	tx, err := tm.BeginTransaction(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			tm.logger.Error(ctx, "Transaction panicked", "panic", p)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			tm.logger.Error(ctx, "Failed to rollback after error",
				"original_error", err,
				"rollback_error", rbErr)

			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rbErr)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// WithSavepoint executes a function within a savepoint.
func (tm *TransactionManager) WithSavepoint(ctx context.Context, tx *Transaction, name string, fn func() error) error {
	if err := tx.Savepoint(ctx, name); err != nil {
		return err
	}

	if err := fn(); err != nil {
		if rbErr := tx.RollbackToSavepoint(ctx, name); rbErr != nil {
			tm.logger.Error(ctx, "Failed to rollback to savepoint",
				"savepoint", name,
				"original_error", err,
				"rollback_error", rbErr)

			return fmt.Errorf("savepoint error: %w, rollback error: %w", err, rbErr)
		}

		return err
	}

	if err := tx.ReleaseSavepoint(ctx, name); err != nil {
		return err
	}

	return nil
}
