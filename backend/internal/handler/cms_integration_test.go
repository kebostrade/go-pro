// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// +build integration

package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/repository/postgres"
	"go-pro-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gorilla/mux"
)

// CMSIntegrationTestHelper provides integration test utilities for CMS handlers.
type CMSIntegrationTestHelper struct {
	t               *testing.T
	db              *sql.DB
	router          *mux.Router
	testHelper      *TestHelper
	cmsTestHelper   *postgres.CMSTestHelper
}

// NewCMSIntegrationTestHelper creates a new CMS integration test helper.
func NewCMSIntegrationTestHelper(t *testing.T) *CMSIntegrationTestHelper {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Connect to test database
	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err, "Failed to connect to test database")

	// Setup router
	router := mux.NewRouter()
	cmsHandler := NewCMSHandler(nil) // TODO: Pass actual dependencies
	router.HandleFunc("/api/cms/content/{lessonId}", cmsHandler.GetContent).Methods("GET")
	router.HandleFunc("/api/cms/content/{lessonId}", cmsHandler.UpdateContent).Methods("PUT")
	router.HandleFunc("/api/cms/content/{lessonId}/publish", cmsHandler.PublishContent).Methods("POST")
	router.HandleFunc("/api/cms/content/{lessonId}/history", cmsHandler.GetVersionHistory).Methods("GET")
	router.HandleFunc("/api/cms/content/{lessonId}/compare", cmsHandler.CompareVersions).Methods("GET")
	router.HandleFunc("/api/cms/content/{lessonId}/rollback", cmsHandler.RollbackToVersion).Methods("POST")

	return &CMSIntegrationTestHelper{
		t:              t,
		db:             db,
		router:         router,
		testHelper:     NewTestHelper(t),
		cmsTestHelper:  postgres.NewCMSTestHelper(t, db),
	}
}

// Cleanup cleans up test resources.
func (h *CMSIntegrationTestHelper) Cleanup() {
	if h.db != nil {
		h.db.Close()
	}
}

// SetupTestUser creates and authenticates a test user.
func (h *CMSIntegrationTestHelper) SetupTestUser(role domain.UserRole) *domain.User {
	ctx := context.Background()
	h.cmsTestHelper.TruncateCMSTables(ctx)

	user := h.cmsTestHelper.CreateTestUser(ctx, role)
	return user
}

// CreateAuthMiddleware creates authentication middleware for testing.
func (h *CMSIntegrationTestHelper) CreateAuthMiddleware(user *domain.User) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MakeAuthenticatedRequest makes an authenticated request.
func (h *CMSIntegrationTestHelper) MakeAuthenticatedRequest(user *domain.User, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(h.t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add user to context
	ctx := middleware.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Apply auth middleware and record response
	recorder := httptest.NewRecorder()
	authMiddleware := h.CreateAuthMiddleware(user)
	authMiddleware(h.router).ServeHTTP(recorder, req)

	return recorder
}

func TestCMSHandler_GetContent(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Get latest published content as instructor", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create published version
		publishedVersion := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", instructor.ID)

		// Make request
		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s", lesson.ID), nil)

		// Assert response
		assert.Equal(t, http.StatusOK, recorder.Code, "Should return 200 OK")

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to unmarshal response")
		assert.True(t, response["success"].(bool), "Response should be successful")

		data := response["data"].(map[string]interface{})
		assert.Equal(t, lesson.ID, data["lesson_id"], "Should return correct lesson ID")
		assert.NotNil(t, data["content"], "Should return content")
	})

	t.Run("Get draft content as instructor", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create draft version
		helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "draft", instructor.ID)

		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s", lesson.ID), nil)

		assert.Equal(t, http.StatusOK, recorder.Code, "Should return 200 OK")

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool), "Response should be successful")
	})

	t.Run("Student cannot access CMS endpoints", func(t *testing.T) {
		student := helper.SetupTestUser(domain.RoleStudent)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		recorder := helper.MakeAuthenticatedRequest(student, "GET", fmt.Sprintf("/api/cms/content/%s", lesson.ID), nil)

		assert.Equal(t, http.StatusForbidden, recorder.Code, "Students should be forbidden from CMS")
	})
}

func TestCMSHandler_UpdateContent(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Update lesson content and create new version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create initial version
		helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", instructor.ID)

		// Update content
		updateRequest := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type":    "text",
					"content": "Updated content",
				},
			},
			"change_description": "Update introduction",
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "PUT", fmt.Sprintf("/api/cms/content/%s", lesson.ID), updateRequest)

		assert.Equal(t, http.StatusOK, recorder.Code, "Should return 200 OK")

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool), "Update should be successful")

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(2), data["version_number"], "Should create version 2")
	})

	t.Run("Auto-save creates draft version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		updateRequest := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type":    "code",
					"content": "package main",
				},
			},
			"auto_save": true,
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "PUT", fmt.Sprintf("/api/cms/content/%s", lesson.ID), updateRequest)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "draft", data["status"], "Auto-save should create draft")
	})
}

func TestCMSHandler_PublishContent(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Publish draft version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create draft version
		helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "draft", instructor.ID)

		publishRequest := map[string]interface{}{
			"version_id": nil, // Publish latest draft
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "POST", fmt.Sprintf("/api/cms/content/%s/publish", lesson.ID), publishRequest)

		assert.Equal(t, http.StatusOK, recorder.Code, "Should publish successfully")

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool), "Publish should be successful")
	})

	t.Run("Cannot publish non-existent version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		publishRequest := map[string]interface{}{
			"version_id": "non-existent-id",
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "POST", fmt.Sprintf("/api/cms/content/%s/publish", lesson.ID), publishRequest)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Should return 400 for non-existent version")
	})
}

func TestCMSHandler_GetVersionHistory(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Get full version history", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create multiple versions
		v1 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", instructor.ID)
		v2 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 2, "draft", instructor.ID)
		v3 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 3, "published", instructor.ID)

		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s/history", lesson.ID), nil)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		versions := data["versions"].([]interface{})
		assert.Len(t, versions, 3, "Should return all versions")
	})

	t.Run("Get paginated version history", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create 5 versions
		for i := 1; i <= 5; i++ {
			helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, i, "published", instructor.ID)
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s/history?page=1&page_size=2", lesson.ID), nil)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		versions := data["versions"].([]interface{})
		assert.Len(t, versions, 2, "Should return paginated results")
		assert.Equal(t, float64(5), data["total"], "Total should reflect all versions")
	})
}

func TestCMSHandler_CompareVersions(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Compare two different versions", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		v1 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", instructor.ID)
		v2 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 2, "published", instructor.ID)

		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s/compare?version1=%s&version2=%s", lesson.ID, v1.ID, v2.ID), nil)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["diff"], "Should return diff output")
		assert.True(t, data["has_changes"].(bool), "Should detect changes")
	})

	t.Run("Missing version parameters", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		recorder := helper.MakeAuthenticatedRequest(instructor, "GET", fmt.Sprintf("/api/cms/content/%s/compare", lesson.ID), nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Should require version parameters")
	})
}

func TestCMSHandler_RollbackToVersion(t *testing.T) {
	helper := NewCMSIntegrationTestHelper(t)
	defer helper.Cleanup()

	ctx := context.Background()

	t.Run("Rollback to previous version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		// Create version sequence
		v1 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 1, "published", instructor.ID)
		v2 := helper.cmsTestHelper.CreateTestContentVersion(ctx, lesson.ID, 2, "published", instructor.ID)

		rollbackRequest := map[string]interface{}{
			"version_id": v1.ID,
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "POST", fmt.Sprintf("/api/cms/content/%s/rollback", lesson.ID), rollbackRequest)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(3), data["version_number"], "Rollback should create new version")
	})

	t.Run("Cannot rollback to non-existent version", func(t *testing.T) {
		instructor := helper.SetupTestUser(domain.RoleInstructor)
		course := helper.cmsTestHelper.CreateTestCourse(ctx)
		lesson := helper.cmsTestHelper.CreateTestLesson(ctx, course.ID)

		rollbackRequest := map[string]interface{}{
			"version_id": "non-existent-id",
		}

		recorder := helper.MakeAuthenticatedRequest(instructor, "POST", fmt.Sprintf("/api/cms/content/%s/rollback", lesson.ID), rollbackRequest)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Should return 400 for non-existent version")
	})
}
