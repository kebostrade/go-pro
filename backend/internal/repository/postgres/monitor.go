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

// PoolMonitor monitors database connection pool health.
type PoolMonitor struct {
	db     *DB
	logger logger.Logger
	mu     sync.RWMutex
	stats  PoolStats
}

// PoolStats represents connection pool statistics.
type PoolStats struct {
	MaxOpenConnections int           `json:"max_open_connections"`
	OpenConnections    int           `json:"open_connections"`
	InUse              int           `json:"in_use"`
	Idle               int           `json:"idle"`
	WaitCount          int64         `json:"wait_count"`
	WaitDuration       time.Duration `json:"wait_duration"`
	MaxIdleClosed      int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed  int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed  int64         `json:"max_lifetime_closed"`
	Timestamp          time.Time     `json:"timestamp"`
}

// NewPoolMonitor creates a new connection pool monitor.
func NewPoolMonitor(db *DB, logger logger.Logger) *PoolMonitor {
	return &PoolMonitor{
		db:     db,
		logger: logger,
	}
}

// GetStats returns current pool statistics.
func (pm *PoolMonitor) GetStats() PoolStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.stats
}

// UpdateStats updates the pool statistics.
func (pm *PoolMonitor) UpdateStats(ctx context.Context) {
	dbStats := pm.db.Stats()

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.stats = PoolStats{
		MaxOpenConnections: dbStats.MaxOpenConnections,
		OpenConnections:    dbStats.OpenConnections,
		InUse:              dbStats.InUse,
		Idle:               dbStats.Idle,
		WaitCount:          dbStats.WaitCount,
		WaitDuration:       dbStats.WaitDuration,
		MaxIdleClosed:      dbStats.MaxIdleClosed,
		MaxIdleTimeClosed:  dbStats.MaxIdleTimeClosed,
		MaxLifetimeClosed:  dbStats.MaxLifetimeClosed,
		Timestamp:          time.Now(),
	}
}

// StartMonitoring starts periodic monitoring of the connection pool.
func (pm *PoolMonitor) StartMonitoring(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	pm.logger.Info(ctx, "Starting connection pool monitoring", "interval", interval)

	for {
		select {
		case <-ctx.Done():
			pm.logger.Info(ctx, "Stopping connection pool monitoring")
			return
		case <-ticker.C:
			pm.UpdateStats(ctx)
			stats := pm.GetStats()

			// Log warnings if pool is under stress.
			if stats.InUse >= int(float64(stats.MaxOpenConnections)*0.8) {
				pm.logger.Warn(ctx, "Connection pool usage is high",
					"in_use", stats.InUse,
					"max_open", stats.MaxOpenConnections,
					"usage_percent", float64(stats.InUse)/float64(stats.MaxOpenConnections)*100)
			}

			if stats.WaitCount > 0 {
				pm.logger.Warn(ctx, "Connections are waiting",
					"wait_count", stats.WaitCount,
					"wait_duration", stats.WaitDuration)
			}

			pm.logger.Debug(ctx, "Connection pool stats",
				"open", stats.OpenConnections,
				"in_use", stats.InUse,
				"idle", stats.Idle,
				"wait_count", stats.WaitCount)
		}
	}
}

// QueryPerformanceTracker tracks query performance.
type QueryPerformanceTracker struct {
	logger logger.Logger
	mu     sync.RWMutex
	stats  map[string]*QueryStats
}

// QueryStats represents statistics for a specific query.
type QueryStats struct {
	Query          string        `json:"query"`
	ExecutionCount int64         `json:"execution_count"`
	TotalDuration  time.Duration `json:"total_duration"`
	AvgDuration    time.Duration `json:"avg_duration"`
	MinDuration    time.Duration `json:"min_duration"`
	MaxDuration    time.Duration `json:"max_duration"`
	ErrorCount     int64         `json:"error_count"`
	LastExecuted   time.Time     `json:"last_executed"`
}

// NewQueryPerformanceTracker creates a new query performance tracker.
func NewQueryPerformanceTracker(logger logger.Logger) *QueryPerformanceTracker {
	return &QueryPerformanceTracker{
		logger: logger,
		stats:  make(map[string]*QueryStats),
	}
}

// TrackQuery tracks the execution of a query.
func (qpt *QueryPerformanceTracker) TrackQuery(query string, duration time.Duration, err error) {
	qpt.mu.Lock()
	defer qpt.mu.Unlock()

	stats, exists := qpt.stats[query]
	if !exists {
		stats = &QueryStats{
			Query:       query,
			MinDuration: duration,
			MaxDuration: duration,
		}
		qpt.stats[query] = stats
	}

	stats.ExecutionCount++
	stats.TotalDuration += duration
	stats.AvgDuration = stats.TotalDuration / time.Duration(stats.ExecutionCount)
	stats.LastExecuted = time.Now()

	if duration < stats.MinDuration {
		stats.MinDuration = duration
	}
	if duration > stats.MaxDuration {
		stats.MaxDuration = duration
	}

	if err != nil {
		stats.ErrorCount++
	}

	// Log slow queries.
	if duration > 1*time.Second {
		qpt.logger.Warn(context.Background(), "Slow query detected",
			"query", query,
			"duration", duration,
			"error", err)
	}
}

// GetStats returns statistics for all tracked queries.
func (qpt *QueryPerformanceTracker) GetStats() map[string]*QueryStats {
	qpt.mu.RLock()
	defer qpt.mu.RUnlock()

	// Create a copy to avoid race conditions.
	statsCopy := make(map[string]*QueryStats, len(qpt.stats))
	for k, v := range qpt.stats {
		statsCopy[k] = &QueryStats{
			Query:          v.Query,
			ExecutionCount: v.ExecutionCount,
			TotalDuration:  v.TotalDuration,
			AvgDuration:    v.AvgDuration,
			MinDuration:    v.MinDuration,
			MaxDuration:    v.MaxDuration,
			ErrorCount:     v.ErrorCount,
			LastExecuted:   v.LastExecuted,
		}
	}

	return statsCopy
}

// GetTopSlowQueries returns the N slowest queries by average duration.
func (qpt *QueryPerformanceTracker) GetTopSlowQueries(n int) []*QueryStats {
	qpt.mu.RLock()
	defer qpt.mu.RUnlock()

	// Convert map to slice.
	queries := make([]*QueryStats, 0, len(qpt.stats))
	for _, stats := range qpt.stats {
		queries = append(queries, stats)
	}

	// Sort by average duration (descending)
	for i := 0; i < len(queries)-1; i++ {
		for j := i + 1; j < len(queries); j++ {
			if queries[i].AvgDuration < queries[j].AvgDuration {
				queries[i], queries[j] = queries[j], queries[i]
			}
		}
	}

	// Return top N.
	if n > len(queries) {
		n = len(queries)
	}

	return queries[:n]
}

// Reset resets all query statistics.
func (qpt *QueryPerformanceTracker) Reset() {
	qpt.mu.Lock()
	defer qpt.mu.Unlock()
	qpt.stats = make(map[string]*QueryStats)
}

// TrackedDB wraps DB with query performance tracking.
type TrackedDB struct {
	*DB
	tracker *QueryPerformanceTracker
}

// NewTrackedDB creates a new tracked database connection.
func NewTrackedDB(db *DB, tracker *QueryPerformanceTracker) *TrackedDB {
	return &TrackedDB{
		DB:      db,
		tracker: tracker,
	}
}

// ExecContext executes a query with performance tracking.
func (tdb *TrackedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := tdb.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	tdb.tracker.TrackQuery(query, duration, err)

	return result, err
}

// QueryContext executes a query with performance tracking.
func (tdb *TrackedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := tdb.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	tdb.tracker.TrackQuery(query, duration, err)

	return rows, err
}

// QueryRowContext executes a query with performance tracking.
func (tdb *TrackedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := tdb.DB.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	// Note: We can't detect errors here as QueryRow doesn't return an error.
	tdb.tracker.TrackQuery(query, duration, nil)

	return row
}

// HealthReport generates a comprehensive health report.
type HealthReport struct {
	PoolStats       PoolStats     `json:"pool_stats"`
	TopSlowQueries  []*QueryStats `json:"top_slow_queries"`
	TotalQueries    int64         `json:"total_queries"`
	TotalErrors     int64         `json:"total_errors"`
	DatabaseVersion string        `json:"database_version"`
	Uptime          time.Duration `json:"uptime"`
	Timestamp       time.Time     `json:"timestamp"`
}

// GenerateHealthReport generates a comprehensive health report.
func GenerateHealthReport(ctx context.Context, db *DB, monitor *PoolMonitor, tracker *QueryPerformanceTracker) (*HealthReport, error) {
	report := &HealthReport{
		PoolStats:      monitor.GetStats(),
		TopSlowQueries: tracker.GetTopSlowQueries(10),
		Timestamp:      time.Now(),
	}

	// Get database version.
	var version string
	if err := db.QueryRowContext(ctx, "SELECT version()").Scan(&version); err != nil {
		return nil, fmt.Errorf("failed to get database version: %w", err)
	}
	report.DatabaseVersion = version

	// Calculate total queries and errors.
	stats := tracker.GetStats()
	for _, s := range stats {
		report.TotalQueries += s.ExecutionCount
		report.TotalErrors += s.ErrorCount
	}

	return report, nil
}
