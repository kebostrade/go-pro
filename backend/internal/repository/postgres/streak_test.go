//go:build integration
// +build integration

// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"
)

const skipIntegrationTest = "Skipping integration test"

// TestStreakIncrement tests streak increment logic (yesterday activity).
func TestStreakIncrement(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTest)
	}

	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	repo := NewStreakRepository(db.DB)
	ctx := context.Background()
	db.TruncateTables(ctx, "streaks")
	userID := "test-user-1"

	// Create initial streak
	yesterday := time.Now().UTC().Add(-24 * time.Hour)
	streak := &domain.Streak{
		UserID:           userID,
		CurrentStreak:    5,
		LongestStreak:    10,
		LastActivityDate: &yesterday,
	}

	err := repo.Upsert(ctx, streak)
	require.NoError(t, err)

	// Update streak with today's activity
	today := time.Now().UTC()
	err = repo.UpdateStreak(ctx, userID, &today)
	require.NoError(t, err)

	// Verify streak incremented
	updated, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 6, updated.CurrentStreak, "streak should increment from 5 to 6")
	assert.Equal(t, 10, updated.LongestStreak, "longest streak should remain 10")
}

// TestStreakReset tests streak reset logic (gap > 1 day).
func TestStreakReset(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTest)
	}

	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	repo := NewStreakRepository(db.DB)
	ctx := context.Background()
	db.TruncateTables(ctx, "streaks")
	userID := "test-user-2"

	// Create initial streak with activity 3 days ago
	threeDaysAgo := time.Now().UTC().Add(-72 * time.Hour)
	streak := &domain.Streak{
		UserID:           userID,
		CurrentStreak:    8,
		LongestStreak:    15,
		LastActivityDate: &threeDaysAgo,
	}

	err := repo.Upsert(ctx, streak)
	require.NoError(t, err)

	// Update streak with today's activity
	today := time.Now().UTC()
	err = repo.UpdateStreak(ctx, userID, &today)
	require.NoError(t, err)

	// Verify streak reset to 1
	updated, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1, updated.CurrentStreak, "streak should reset to 1")
	assert.Equal(t, 15, updated.LongestStreak, "longest streak should remain 15")
}

// TestLongestStreakPreservation tests that longest streak is preserved.
func TestLongestStreakPreservation(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTest)
	}

	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	repo := NewStreakRepository(db.DB)
	ctx := context.Background()
	db.TruncateTables(ctx, "streaks")
	userID := "test-user-3"

	// Start with streak of 1
	today := time.Now().UTC()
	streak := &domain.Streak{
		UserID:           userID,
		CurrentStreak:    1,
		LongestStreak:    20,
		LastActivityDate: &today,
	}

	err := repo.Upsert(ctx, streak)
	require.NoError(t, err)

	// Update multiple times to build streak
	for i := 0; i < 25; i++ {
		yesterday := time.Now().UTC().Add(-time.Duration(i) * 24 * time.Hour)
		err := repo.UpdateStreak(ctx, userID, &yesterday)
		require.NoError(t, err)
	}

	// Verify longest streak updated if current exceeded it
	updated, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, updated.LongestStreak, 20, "longest streak should be preserved or updated")
}

// TestNewUserStreakInitialization tests new user streak initialization.
func TestNewUserStreakInitialization(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTest)
	}

	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	repo := NewStreakRepository(db.DB)
	ctx := context.Background()
	db.TruncateTables(ctx, "streaks")
	userID := "test-user-4"

	// Initialize streak for new user
	err := repo.CreateInitialStreak(ctx, userID)
	require.NoError(t, err)

	// Retrieve and verify initialization
	streak, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, streak)
	assert.Equal(t, 0, streak.CurrentStreak, "new user should have current streak of 0")
	assert.Equal(t, 0, streak.LongestStreak, "new user should have longest streak of 0")
	assert.Nil(t, streak.LastActivityDate, "new user should have no last activity date")
}
