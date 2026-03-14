// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// +build integration

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentVersionRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err, "Failed to connect to test database")
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)

	// Clean up before test
	helper.TruncateCMSTables(ctx)

	t.Run("Create successful content version", func(t *testing.T) {
		// Create test lesson
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		sections := []domain.ContentSection{
			{
				Type:    "text",
				Content: "Introduction to Go",
			},
			{
				Type:    "code",
				Content: "package main",
				Config: map[string]interface{}{
					"language": "go",
				},
			},
		}

		contentVersion := &domain.ContentVersion{
			LessonID:          lesson.ID,
			VersionNumber:    1,
			AuthorID:          author.ID,
			ChangeDescription: "Initial version",
			Status:            "draft",
		}

		contentVersion.SetSections(sections)

		err := repo.Create(ctx, contentVersion)
		require.NoError(t, err, "Failed to create content version")
		require.NotEmpty(t, contentVersion.ID, "ID should be generated")
		require.NotEmpty(t, contentVersion.CreatedAt, "CreatedAt should be set")
	})

	t.Run("Create with published status", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		contentVersion := &domain.ContentVersion{
			LessonID:          lesson.ID,
			VersionNumber:    1,
			AuthorID:          author.ID,
			ChangeDescription: "First published version",
			Status:            "published",
		}

		contentVersion.SetSections([]domain.ContentSection{
			{Type: "text", Content: "Published content"},
		})

		err := repo.Create(ctx, contentVersion)
		require.NoError(t, err, "Failed to create published content version")
		assert.NotNil(t, contentVersion.PublishedAt, "PublishedAt should be set for published versions")
	})
}

func TestContentVersionRepository_GetByVersionNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Get existing version by number", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		// Create test version
		created := helper.CreateTestContentVersion(ctx, lesson.ID, 2, "published", author.ID)

		// Retrieve by version number
		retrieved, err := repo.GetByVersionNumber(ctx, lesson.ID, 2)
		require.NoError(t, err, "Failed to get version by number")
		helper.AssertContentVersionEquals(created, retrieved)
	})

	t.Run("Get non-existent version", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		_, err := repo.GetByVersionNumber(ctx, lesson.ID, 999)
		assert.Error(t, err, "Should return error for non-existent version")
	})
}

func TestContentVersionRepository_GetVersionHistory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Get version history with multiple versions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		// Create multiple versions
		v1 := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", author.ID)
		v2 := helper.CreateTestContentVersion(ctx, lesson.ID, 2, "draft", author.ID)
		v3 := helper.CreateTestContentVersion(ctx, lesson.ID, 3, "published", author.ID)

		// Get history
		history, err := repo.GetVersionHistory(ctx, lesson.ID, nil)
		require.NoError(t, err, "Failed to get version history")
		require.Len(t, history, 3, "Should return 3 versions")

		// Verify ordering (most recent first)
		assert.Equal(t, v3.ID, history[0].ID, "First version should be most recent")
		assert.Equal(t, v2.ID, history[1].ID)
		assert.Equal(t, v1.ID, history[2].ID)
	})

	t.Run("Get version history with pagination", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		// Create 5 versions
		for i := 1; i <= 5; i++ {
			helper.CreateTestContentVersion(ctx, lesson.ID, i, "published", author.ID)
		}

		// Get first page
		pagination := &domain.PaginationRequest{
			Page:     1,
			PageSize: 2,
		}
		history, err := repo.GetVersionHistory(ctx, lesson.ID, pagination)
		require.NoError(t, err, "Failed to get paginated version history")
		require.Len(t, history, 2, "Should return 2 versions per page")
	})

	t.Run("Get empty history for lesson with no versions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		history, err := repo.GetVersionHistory(ctx, lesson.ID, nil)
		require.NoError(t, err, "Should not error for empty history")
		require.Empty(t, history, "Should return empty slice")
	})
}

func TestContentVersionRepository_GetLatestVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Get latest published version", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		// Create version sequence
		helper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", author.ID)
		helper.CreateTestContentVersion(ctx, lesson.ID, 2, "draft", author.ID)
		latestPublished := helper.CreateTestContentVersion(ctx, lesson.ID, 3, "published", author.ID)

		// Get latest published
		latest, err := repo.GetLatestVersion(ctx, lesson.ID)
		require.NoError(t, err, "Failed to get latest version")
		assert.Equal(t, latestPublished.ID, latest.ID, "Should return latest published version")
		assert.Equal(t, int32(3), latest.VersionNumber, "Should be version 3")
	})

	t.Run("Get latest when only draft exists", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		draftVersion := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "draft", author.ID)

		latest, err := repo.GetLatestVersion(ctx, lesson.ID)
		require.NoError(t, err, "Should return draft if no published version")
		assert.Equal(t, draftVersion.ID, latest.ID, "Should return draft version")
	})

	t.Run("Get latest for lesson with no versions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		_, err := repo.GetLatestVersion(ctx, lesson.ID)
		assert.Error(t, err, "Should return error when no versions exist")
	})
}

func TestContentVersionRepository_CompareVersions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Compare different versions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		// Create two versions with different content
		v1 := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", author.ID)
		v2 := helper.CreateTestContentVersion(ctx, lesson.ID, 2, "published", author.ID)

		// Update v2 content
		newSections := []domain.ContentSection{
			{Type: "text", Content: "Updated content"},
			{Type: "code", Content: "package updated"},
		}
		v2.SetSections(newSections)
		v2JSON, _ := json.Marshal(newSections)
		_, err := db.ExecContext(ctx, "UPDATE content_versions SET content = $1 WHERE id = $2", v2JSON, v2.ID)
		require.NoError(t, err)

		// Compare
		comparison, err := repo.CompareVersions(ctx, v1.ID, v2.ID)
		require.NoError(t, err, "Failed to compare versions")
		assert.NotNil(t, comparison, "Comparison should not be nil")
		assert.True(t, comparison.HasChanges, "Should detect changes")
		assert.NotEmpty(t, comparison.Diff, "Should have diff output")
	})

	t.Run("Compare identical versions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		v1 := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", author.ID)

		comparison, err := repo.CompareVersions(ctx, v1.ID, v1.ID)
		require.NoError(t, err, "Failed to compare identical versions")
		assert.False(t, comparison.HasChanges, "Should not detect changes in identical versions")
	})
}

func TestContentVersionRepository_PublishVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Publish draft version", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		draftVersion := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "draft", author.ID)
		require.Nil(t, draftVersion.PublishedAt, "Draft should not have PublishedAt")

		err := repo.PublishVersion(ctx, draftVersion.ID)
		require.NoError(t, err, "Failed to publish version")

		// Verify update
		var publishedAt time.Time
		var status string
		err = db.QueryRowContext(ctx, "SELECT published_at, status FROM content_versions WHERE id = $1", draftVersion.ID).Scan(&publishedAt, &status)
		require.NoError(t, err, "Failed to query published version")
		assert.Equal(t, "published", status, "Status should be published")
		assert.False(t, publishedAt.IsZero(), "PublishedAt should be set")
	})
}

func TestContentVersionRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewContentVersionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Delete existing version", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		author := helper.CreateTestUser(ctx, domain.RoleInstructor)

		version := helper.CreateTestContentVersion(ctx, lesson.ID, 1, "draft", author.ID)

		err := repo.Delete(ctx, version.ID)
		require.NoError(t, err, "Failed to delete version")

		// Verify deletion
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM content_versions WHERE id = $1", version.ID).Scan(&count)
		require.NoError(t, err, "Failed to check deletion")
		assert.Equal(t, 0, count, "Version should be deleted")
	})

	t.Run("Delete non-existent version", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent-id")
		assert.NoError(t, err, "Delete should be idempotent")
	})
}
