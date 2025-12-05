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
)

// StreakRepository implements the repository.StreakRepository interface using PostgreSQL.
type StreakRepository struct {
	db *sql.DB
}

// NewStreakRepository creates a new StreakRepository.
func NewStreakRepository(db *sql.DB) *StreakRepository {
	return &StreakRepository{
		db: db,
	}
}

// GetByUserID retrieves streak data for a user.
func (r *StreakRepository) GetByUserID(ctx context.Context, userID string) (*domain.Streak, error) {
	query := `
		SELECT
			user_id,
			current_streak,
			longest_streak,
			last_activity_date,
			created_at,
			updated_at
		FROM streaks
		WHERE user_id = $1
	`

	streak := &domain.Streak{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&streak.UserID,
		&streak.CurrentStreak,
		&streak.LongestStreak,
		&streak.LastActivityDate,
		&streak.CreatedAt,
		&streak.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No streak record yet
		}
		return nil, fmt.Errorf("failed to get streak: %w", err)
	}

	return streak, nil
}

// Upsert creates or updates a streak record.
func (r *StreakRepository) Upsert(ctx context.Context, streak *domain.Streak) error {
	query := `
		INSERT INTO streaks (user_id, current_streak, longest_streak, last_activity_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = $2,
			longest_streak = $3,
			last_activity_date = $4,
			updated_at = $6
	`

	now := time.Now().UTC()
	if streak.CreatedAt.IsZero() {
		streak.CreatedAt = now
	}
	streak.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		streak.UserID,
		streak.CurrentStreak,
		streak.LongestStreak,
		streak.LastActivityDate,
		streak.CreatedAt,
		streak.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert streak: %w", err)
	}

	return nil
}

// UpdateStreak updates streak logic based on activity.
// Increments if last activity was yesterday, resets if gap > 1 day.
func (r *StreakRepository) UpdateStreak(ctx context.Context, userID string, lastActivityDate *time.Time) error {
	// Get current streak
	streak, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Create new streak record if it doesn't exist
	if streak == nil {
		streak = &domain.Streak{
			UserID:           userID,
			CurrentStreak:    1,
			LongestStreak:    1,
			LastActivityDate: lastActivityDate,
		}
		return r.Upsert(ctx, streak)
	}

	// If no previous activity, start streak at 1
	if streak.LastActivityDate == nil {
		streak.CurrentStreak = 1
		streak.LongestStreak = 1
		streak.LastActivityDate = lastActivityDate
		return r.Upsert(ctx, streak)
	}

	// Calculate days between last activity and today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastActivity := streak.LastActivityDate.Truncate(24 * time.Hour)
	daysSinceLastActivity := today.Sub(lastActivity).Hours() / 24

	// If activity today, no change
	if daysSinceLastActivity < 1 {
		return nil
	}

	// If activity yesterday, increment streak
	if daysSinceLastActivity >= 1 && daysSinceLastActivity < 2 {
		streak.CurrentStreak++
		// Update longest streak if needed
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}
	} else {
		// Gap > 1 day, reset current streak
		streak.CurrentStreak = 1
	}

	streak.LastActivityDate = lastActivityDate
	return r.Upsert(ctx, streak)
}

// CreateInitialStreak creates an initial streak record for a new user.
func (r *StreakRepository) CreateInitialStreak(ctx context.Context, userID string) error {
	return r.Upsert(ctx, &domain.Streak{
		UserID:        userID,
		CurrentStreak: 0,
		LongestStreak: 0,
	})
}
