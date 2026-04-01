package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DimaJoyti/go-pro/basic/projects/testing-patterns/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation for handler tests
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *service.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.ID = "test-id"
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*service.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*service.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("GET /users", handler.ListUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
	return mux
}

func TestCreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).Return(nil)

		body := `{"name":"Alice","email":"alice@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var user service.User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alice@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		body := `{invalid}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		body := `{"name":"","email":"alice@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		expectedUser := &service.User{ID: "123", Name: "Test", Email: "test@example.com"}
		mockRepo.On("GetByID", mock.Anything, "123").Return(expectedUser, nil)

		req := httptest.NewRequest("GET", "/users/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var user service.User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, "123", user.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, service.ErrNotFound)

		req := httptest.NewRequest("GET", "/users/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestListUsers(t *testing.T) {
	t.Run("list users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		users := []*service.User{
			{ID: "1", Name: "User1", Email: "user1@example.com"},
			{ID: "2", Name: "User2", Email: "user2@example.com"},
		}
		mockRepo.On("List", mock.Anything).Return(users, nil)

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []*service.User
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("list error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("List", mock.Anything).Return(nil, assert.AnError)

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("Delete", mock.Anything, "123").Return(nil)

		req := httptest.NewRequest("DELETE", "/users/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("Delete", mock.Anything, "non-existent").Return(service.ErrNotFound)

		req := httptest.NewRequest("DELETE", "/users/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateUserInternalError(t *testing.T) {
	t.Run("repo create error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)
		handler := NewUserHandler(svc)
		router := setupTestRouter(handler)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).Return(assert.AnError)

		body := `{"name":"Alice","email":"alice@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
