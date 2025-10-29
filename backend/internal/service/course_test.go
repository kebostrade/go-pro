package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"
	"go-pro-backend/pkg/validator"
)

func TestCourseService_Create(t *testing.T) {
	tests := []struct {
		name        string
		course      *domain.Course
		wantErr     bool
		errContains string
	}{
		{
			name: "valid course",
			course: &domain.Course{
				ID:          "course-1",
				Title:       "Go Programming",
				Description: "Learn Go programming from basics to advanced",
				Lessons:     []string{},
			},
			wantErr: false,
		},
		{
			name: "missing title",
			course: &domain.Course{
				ID:          "course-2",
				Description: "Test description that is long enough",
			},
			wantErr:     true,
			errContains: "Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := testutil.NewMockCourseRepository()
			mockCache := testutil.NewMockCacheManager()
			logger := testutil.NewTestLogger(t)
			v := validator.New()

			config := &Config{
				Cache:     mockCache,
				Messaging: nil, // Not used in current implementation
				Logger:    logger,
				Validator: v,
			}

			service := NewCourseService(mockRepo, config)
			ctx := context.Background()

			// Create request from course
			req := &domain.CreateCourseRequest{
				Title:       tt.course.Title,
				Description: tt.course.Description,
			}

			// Act
			_, err := service.CreateCourse(ctx, req)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, mockRepo.GetCallCount("Create"))
			}
		})
	}
}

func TestCourseService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		courseID    string
		setupMock   func(*testutil.MockCourseRepository, *testutil.MockCacheManager)
		wantErr     bool
		errContains string
	}{
		{
			name:     "course found in repository",
			courseID: "course-1",
			setupMock: func(repo *testutil.MockCourseRepository, cache *testutil.MockCacheManager) {
				course := testutil.CreateTestCourse("course-1", "Go Programming")
				repo.AddCourse(course)
			},
			wantErr: false,
		},
		{
			name:     "course not found",
			courseID: "course-999",
			setupMock: func(repo *testutil.MockCourseRepository, cache *testutil.MockCacheManager) {
				// No setup - course doesn't exist
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := testutil.NewMockCourseRepository()
			mockCache := testutil.NewMockCacheManager()
			logger := testutil.NewTestLogger(t)
			v := validator.New()

			tt.setupMock(mockRepo, mockCache)

			config := &Config{
				Cache:     mockCache,
				Messaging: nil, // Not used in current implementation
				Logger:    logger,
				Validator: v,
			}

			service := NewCourseService(mockRepo, config)
			ctx := context.Background()

			// Act
			course, err := service.GetCourseByID(ctx, tt.courseID)

			// Assert
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

func TestCourseService_Update(t *testing.T) {
	tests := []struct {
		name        string
		course      *domain.Course
		setupMock   func(*testutil.MockCourseRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			course: &domain.Course{
				ID:          "course-1",
				Title:       "Updated Go Programming",
				Description: "Updated description for the course",
				Lessons:     []string{},
			},
			setupMock: func(repo *testutil.MockCourseRepository) {
				course := testutil.CreateTestCourse("course-1", "Go Programming")
				repo.AddCourse(course)
			},
			wantErr: false,
		},
		{
			name: "course not found",
			course: &domain.Course{
				ID:          "course-999",
				Title:       "Non-existent Course",
				Description: "Description for non-existent course",
				Lessons:     []string{},
			},
			setupMock: func(repo *testutil.MockCourseRepository) {
				// No setup - course doesn't exist
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := testutil.NewMockCourseRepository()
			mockCache := testutil.NewMockCacheManager()
			logger := testutil.NewTestLogger(t)
			v := validator.New()

			tt.setupMock(mockRepo)

			config := &Config{
				Cache:     mockCache,
				Messaging: nil, // Not used in current implementation
				Logger:    logger,
				Validator: v,
			}

			service := NewCourseService(mockRepo, config)
			ctx := context.Background()

			// Create update request
			req := &domain.UpdateCourseRequest{
				Title:       &tt.course.Title,
				Description: &tt.course.Description,
			}

			// Act
			_, err := service.UpdateCourse(ctx, tt.course.ID, req)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, mockRepo.GetCallCount("GetByID"))
				assert.Equal(t, 1, mockRepo.GetCallCount("Update"))
			}
		})
	}
}

func TestCourseService_Delete(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockCourseRepository()
	mockCache := testutil.NewMockCacheManager()
	logger := testutil.NewTestLogger(t)
	v := validator.New()

	course := testutil.CreateTestCourse("course-1", "Go Programming")
	mockRepo.AddCourse(course)

	config := &Config{
		Cache:     mockCache,
		Messaging: nil, // Not used in current implementation
		Logger:    logger,
		Validator: v,
	}

	service := NewCourseService(mockRepo, config)
	ctx := context.Background()

	// Act
	err := service.DeleteCourse(ctx, "course-1")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, mockRepo.GetCallCount("GetByID"))
	assert.Equal(t, 1, mockRepo.GetCallCount("Delete"))

	// Verify course is deleted
	_, err = mockRepo.GetByID(ctx, "course-1")
	assert.Error(t, err)
}

// Benchmark tests
func BenchmarkCourseService_Create(b *testing.B) {
	mockRepo := testutil.NewMockCourseRepository()
	mockCache := testutil.NewMockCacheManager()
	logger := testutil.NewTestLogger(&testing.T{})
	v := validator.New()

	config := &Config{
		Cache:     mockCache,
		Messaging: nil, // Not used in current implementation
		Logger:    logger,
		Validator: v,
	}

	service := NewCourseService(mockRepo, config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &domain.CreateCourseRequest{
			Title:       "Benchmark Course",
			Description: "This is a benchmark course for testing performance",
		}
		_, _ = service.CreateCourse(ctx, req)
	}
}
