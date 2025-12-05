// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/pkg/logger"
)

// MockAuthService is a mock implementation of AuthService for testing.
type MockAuthService struct {
	verifyTokenFunc func(ctx context.Context, token string) (*middleware.FirebaseToken, error)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
	if m.verifyTokenFunc != nil {
		return m.verifyTokenFunc(ctx, token)
	}
	return &middleware.FirebaseToken{
		UID:         "test-uid",
		Email:       "test@example.com",
		DisplayName: "Test User",
		PhotoURL:    "",
	}, nil
}

// MockUserRepository is a mock implementation of UserRepository for testing.
type MockUserRepository struct {
	getByFirebaseUIDFunc func(ctx context.Context, firebaseUID string) (*domain.User, error)
	createFunc           func(ctx context.Context, user *domain.User) error
	updateLastLoginFunc  func(ctx context.Context, userID string) error
}

func (m *MockUserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	if m.getByFirebaseUIDFunc != nil {
		return m.getByFirebaseUIDFunc(ctx, firebaseUID)
	}
	return &domain.User{
		ID:          "test-user-id",
		FirebaseUID: firebaseUID,
		Email:       "test@example.com",
		DisplayName: "Test User",
		Role:        domain.RoleStudent,
		IsActive:    true,
	}, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	user.ID = "new-user-id"
	return nil
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	if m.updateLastLoginFunc != nil {
		return m.updateLastLoginFunc(ctx, userID)
	}
	return nil
}

// Implement other UserRepository methods as no-ops for testing.
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *MockUserRepository) UpdateLastActivity(ctx context.Context, userID string) error {
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func TestAuthRequired_Success(t *testing.T) {
	// Setup
	mockAuth := &MockAuthService{}
	mockUserRepo := &MockUserRepository{}
	log := logger.New("debug", "text")

	authMiddleware := middleware.NewAuthMiddleware(mockAuth, mockUserRepo, log)

	// Create test handler
	handler := authMiddleware.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify user is in context
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			t.Error("Expected user in context")
			return
		}

		if user.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", user.Email)
		}

		w.WriteHeader(http.StatusOK)
	}))

	// Create request with valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Verify
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAuthRequired_MissingToken(t *testing.T) {
	// Setup
	mockAuth := &MockAuthService{}
	mockUserRepo := &MockUserRepository{}
	log := logger.New("debug", "text")

	authMiddleware := middleware.NewAuthMiddleware(mockAuth, mockUserRepo, log)

	handler := authMiddleware.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	}))

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Verify
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAdminRequired_Success(t *testing.T) {
	// Setup
	mockAuth := &MockAuthService{}
	mockUserRepo := &MockUserRepository{
		getByFirebaseUIDFunc: func(ctx context.Context, firebaseUID string) (*domain.User, error) {
			return &domain.User{
				ID:          "admin-user-id",
				FirebaseUID: firebaseUID,
				Email:       "admin@example.com",
				Role:        domain.RoleAdmin,
				IsActive:    true,
			}, nil
		},
	}
	log := logger.New("debug", "text")

	authMiddleware := middleware.NewAuthMiddleware(mockAuth, mockUserRepo, log)

	// Chain AuthRequired and AdminRequired
	handler := authMiddleware.AuthRequired(
		authMiddleware.AdminRequired(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		),
	)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/admin/test", nil)
	req.Header.Set("Authorization", "Bearer valid-admin-token")

	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Verify
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAdminRequired_Forbidden(t *testing.T) {
	// Setup
	mockAuth := &MockAuthService{}
	mockUserRepo := &MockUserRepository{
		getByFirebaseUIDFunc: func(ctx context.Context, firebaseUID string) (*domain.User, error) {
			return &domain.User{
				ID:          "student-user-id",
				FirebaseUID: firebaseUID,
				Email:       "student@example.com",
				Role:        domain.RoleStudent, // Not admin
				IsActive:    true,
			}, nil
		},
	}
	log := logger.New("debug", "text")

	authMiddleware := middleware.NewAuthMiddleware(mockAuth, mockUserRepo, log)

	// Chain AuthRequired and AdminRequired
	handler := authMiddleware.AuthRequired(
		authMiddleware.AdminRequired(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Handler should not be called for non-admin user")
			}),
		),
	)

	// Create request with valid token but non-admin user
	req := httptest.NewRequest("GET", "/admin/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rr, req)

	// Verify
	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestGetUserFromContext(t *testing.T) {
	// Create a user
	user := &domain.User{
		ID:    "test-id",
		Email: "test@example.com",
		Role:  domain.RoleStudent,
	}

	// Create context with user
	ctx := middleware.WithUser(context.Background(), user)

	// Retrieve user from context
	retrievedUser, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		t.Error("Expected user in context")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}
}

func TestGetUserFromContext_NoUser(t *testing.T) {
	// Create empty context
	ctx := context.Background()

	// Try to retrieve user
	_, ok := middleware.GetUserFromContext(ctx)
	if ok {
		t.Error("Expected no user in context")
	}
}
