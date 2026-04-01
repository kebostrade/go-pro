package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/repository"
	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/service"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(handler *UserHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.ListUsers)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})
	return r
}

func TestHealthCheck(t *testing.T) {
	handler := NewHealthHandler()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestCreateUser(t *testing.T) {
	repo := repository.NewMemoryUserStore()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	router := setupTestRouter(handler)

	t.Run("valid user creation", func(t *testing.T) {
		body := `{"name":"Alice","email":"alice@example.com"}`
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Alice")
		assert.Contains(t, w.Body.String(), "alice@example.com")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		body := `{invalid json}`
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		body := `{"name":"Bob"}`
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid email format", func(t *testing.T) {
		body := `{"name":"Bob","email":"invalid-email"}`
		req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListUsers(t *testing.T) {
	repo := repository.NewMemoryUserStore()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	router := setupTestRouter(handler)

	// Seed users exist from NewMemoryUserStore
	req, _ := http.NewRequest("GET", "/api/v1/users/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Should contain seeded users
	assert.Contains(t, w.Body.String(), "John Doe")
	assert.Contains(t, w.Body.String(), "Jane Smith")
	assert.Contains(t, w.Body.String(), "Bob Wilson")
}

func TestGetUser(t *testing.T) {
	repo := repository.NewMemoryUserStore()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	router := setupTestRouter(handler)

	t.Run("existing user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/seed-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("non-existing user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	repo := repository.NewMemoryUserStore()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	router := setupTestRouter(handler)

	t.Run("successful update", func(t *testing.T) {
		body := `{"name":"UpdatedUser","email":"updated@example.com"}`
		req, _ := http.NewRequest("PUT", "/api/v1/users/seed-1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UpdatedUser")
	})

	t.Run("user not found", func(t *testing.T) {
		body := `{"name":"UpdatedUser","email":"updated@example.com"}`
		req, _ := http.NewRequest("PUT", "/api/v1/users/non-existent-id", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	repo := repository.NewMemoryUserStore()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	router := setupTestRouter(handler)

	t.Run("successful delete", func(t *testing.T) {
		// First create a user to delete
		newUser, _ := svc.Create(context.Background(), "ToDelete", "delete@example.com")

		req, _ := http.NewRequest("DELETE", "/api/v1/users/"+newUser.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Test for JSON response structure
func TestHealthCheckResponseFormat(t *testing.T) {
	handler := NewHealthHandler()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.NotEmpty(t, response["timestamp"])
}
