//go:build integration
// +build integration

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"
)

func TestCourseRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)
	course := testutil.CreateTestCourse("course-1", "Go Programming")

	// Act.
	err := repo.Create(ctx, course)

	// Assert.
	require.NoError(t, err)

	// Verify course was created.
	retrieved, err := repo.GetByID(ctx, "course-1")
	require.NoError(t, err)
	assert.Equal(t, course.ID, retrieved.ID)
	assert.Equal(t, course.Title, retrieved.Title)
	assert.Equal(t, course.Slug, retrieved.Slug)
}

func TestCourseRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name        string
		courseID    string
		setupData   func(*CourseRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "course exists",
			courseID: "course-1",
			setupData: func(repo *CourseRepository) {
				course := testutil.CreateTestCourse("course-1", "Go Programming")
				_ = repo.Create(context.Background(), course)
			},
			wantErr: false,
		},
		{
			name:        "course not found",
			courseID:    "course-999",
			setupData:   func(repo *CourseRepository) {},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			db := testutil.NewTestDB(t)
			defer db.Cleanup()

			ctx := context.Background()
			db.TruncateTables(ctx, "courses")

			repo := NewCourseRepository(db.DB)
			tt.setupData(repo)

			// Act.
			course, err := repo.GetByID(ctx, tt.courseID)

			// Assert.
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, course)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, course)
				assert.Equal(t, tt.courseID, course.ID)
			}
		})
	}
}

func TestCourseRepository_GetBySlug(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)
	course := testutil.CreateTestCourse("course-1", "Go Programming")
	err := repo.Create(ctx, course)
	require.NoError(t, err)

	// Act.
	retrieved, err := repo.GetBySlug(ctx, course.Slug)

	// Assert.
	require.NoError(t, err)
	assert.Equal(t, course.ID, retrieved.ID)
	assert.Equal(t, course.Slug, retrieved.Slug)
}

func TestCourseRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)

	// Create multiple courses.
	courses := []*domain.Course{
		testutil.CreateTestCourse("course-1", "Go Basics"),
		testutil.CreateTestCourse("course-2", "Advanced Go"),
		testutil.CreateTestCourse("course-3", "Go Concurrency"),
	}

	for _, course := range courses {
		err := repo.Create(ctx, course)
		require.NoError(t, err)
	}

	// Act.
	req := &domain.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	result, err := repo.List(ctx, req)

	// Assert.
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.Items))
	assert.Equal(t, 3, result.TotalItems)
	assert.Equal(t, 1, result.TotalPages)
}

func TestCourseRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)
	course := testutil.CreateTestCourse("course-1", "Go Programming")
	err := repo.Create(ctx, course)
	require.NoError(t, err)

	// Modify course.
	course.Title = "Updated Go Programming"
	course.Difficulty = domain.DifficultyAdvanced

	// Act.
	err = repo.Update(ctx, course)

	// Assert.
	require.NoError(t, err)

	// Verify update.
	retrieved, err := repo.GetByID(ctx, "course-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Go Programming", retrieved.Title)
	assert.Equal(t, domain.DifficultyAdvanced, retrieved.Difficulty)
}

func TestCourseRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)
	course := testutil.CreateTestCourse("course-1", "Go Programming")
	err := repo.Create(ctx, course)
	require.NoError(t, err)

	// Act.
	err = repo.Delete(ctx, "course-1")

	// Assert.
	require.NoError(t, err)

	// Verify deletion.
	_, err = repo.GetByID(ctx, "course-1")
	assert.Error(t, err)
}

func TestCourseRepository_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange.
	db := testutil.NewTestDB(t)
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)

	// Create 25 courses.
	for i := 1; i <= 25; i++ {
		course := testutil.CreateTestCourse(
			testutil.RandomString(10),
			testutil.RandomString(20),
		)
		err := repo.Create(ctx, course)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		page          int
		pageSize      int
		expectedItems int
		expectedPages int
	}{
		{
			name:          "first page",
			page:          1,
			pageSize:      10,
			expectedItems: 10,
			expectedPages: 3,
		},
		{
			name:          "second page",
			page:          2,
			pageSize:      10,
			expectedItems: 10,
			expectedPages: 3,
		},
		{
			name:          "last page",
			page:          3,
			pageSize:      10,
			expectedItems: 5,
			expectedPages: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			req := &domain.PaginationRequest{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			result, err := repo.List(ctx, req)

			// Assert.
			require.NoError(t, err)
			assert.Equal(t, tt.expectedItems, len(result.Items))
			assert.Equal(t, 25, result.TotalItems)
			assert.Equal(t, tt.expectedPages, result.TotalPages)
			assert.Equal(t, tt.page, result.Page)
		})
	}
}

// Benchmark tests.
func BenchmarkCourseRepository_Create(b *testing.B) {
	db := testutil.NewTestDB(&testing.T{})
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		course := testutil.CreateTestCourse(
			testutil.RandomString(10),
			testutil.RandomString(20),
		)
		_ = repo.Create(ctx, course)
	}
}

func BenchmarkCourseRepository_GetByID(b *testing.B) {
	db := testutil.NewTestDB(&testing.T{})
	defer db.Cleanup()

	ctx := context.Background()
	db.TruncateTables(ctx, "courses")

	repo := NewCourseRepository(db.DB)
	course := testutil.CreateTestCourse("course-bench", "Benchmark Course")
	_ = repo.Create(ctx, course)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByID(ctx, "course-bench")
	}
}
