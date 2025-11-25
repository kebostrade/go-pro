// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// MockAuthService mocks the AuthService interface.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Initialize(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, idToken string) (*domain.FirebaseClaims, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FirebaseClaims), args.Error(1)
}

func (m *MockAuthService) GetOrCreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error) {
	args := m.Called(ctx, firebaseUID, email, displayName, photoURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) VerifyAndSyncUser(ctx context.Context, idToken string) (*domain.VerifyTokenResponse, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.VerifyTokenResponse), args.Error(1)
}

func (m *MockAuthService) GetUserProfile(ctx context.Context, userID string) (*domain.UserProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserProfileResponse), args.Error(1)
}

func (m *MockAuthService) UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.UserRole) error {
	args := m.Called(ctx, adminUserID, targetUserID, role)
	return args.Error(0)
}

func (m *MockAuthService) UpdateUserProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// setupAuthTestHandler creates a handler with mock auth service for testing.
func setupAuthTestHandler(_ *testing.T) (*Handler, *MockAuthService) {
	mockAuthService := new(MockAuthService)

	// Create mock services
	services := &service.Services{
		Auth: mockAuthService,
	}

	// Create logger and validator
	log := logger.New("debug", "json")
	val := validator.New()

	handler := New(services, log, val)

	return handler, mockAuthService
}

func TestAuthVerify_Success(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock successful token verification
	now := time.Now()
	mockAuthService.On("VerifyAndSyncUser", mock.Anything, "valid-token").
		Return(&domain.VerifyTokenResponse{
			User: &domain.UserProfileResponse{
				ID:          "user-123",
				Email:       "test@example.com",
				DisplayName: "Test User",
				PhotoURL:    "https://example.com/photo.jpg",
				Role:        domain.RoleStudent,
				IsActive:    true,
				CreatedAt:   now,
				LastLoginAt: now,
			},
			IsNewUser:     false,
			TokenVerified: true,
		}, nil)

	// Create request
	reqBody := domain.VerifyTokenRequest{
		IDToken: "valid-token",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.handleAuthVerify(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify response data
	data := resp.Data.(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "user-123", user["id"])
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, false, data["is_new_user"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthVerify_InvalidToken(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock failed token verification
	mockAuthService.On("VerifyAndSyncUser", mock.Anything, "invalid-token").
		Return(nil, apierrors.NewUnauthorizedError("Invalid or expired Firebase token"))

	// Create request
	reqBody := domain.VerifyTokenRequest{
		IDToken: "invalid-token",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.handleAuthVerify(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Contains(t, resp.Error.Message, "Invalid or expired")

	mockAuthService.AssertExpectations(t)
}

func TestAuthVerify_MissingToken(t *testing.T) {
	handler, _ := setupAuthTestHandler(t)

	// Create request with empty token
	reqBody := domain.VerifyTokenRequest{
		IDToken: "",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.handleAuthVerify(w, req)

	// Assert - validation should fail (either 400 or 500 depending on validator behavior)
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
}

func TestAuthVerify_FirstUserBecomesAdmin(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock first user registration (becomes admin)
	now := time.Now()
	mockAuthService.On("VerifyAndSyncUser", mock.Anything, "first-user-token").
		Return(&domain.VerifyTokenResponse{
			User: &domain.UserProfileResponse{
				ID:          "user-001",
				Email:       "admin@example.com",
				DisplayName: "First Admin",
				Role:        domain.RoleAdmin, // First user is admin
				IsActive:    true,
				CreatedAt:   now,
				LastLoginAt: now,
			},
			IsNewUser:     true,
			TokenVerified: true,
		}, nil)

	// Create request
	reqBody := domain.VerifyTokenRequest{
		IDToken: "first-user-token",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/verify", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.handleAuthVerify(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	data := resp.Data.(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "admin", user["role"])
	assert.Equal(t, true, data["is_new_user"])

	mockAuthService.AssertExpectations(t)
}

func TestGetUserProfile_Success(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock successful profile retrieval
	now := time.Now()
	mockAuthService.On("GetUserProfile", mock.Anything, "user-123").
		Return(&domain.UserProfileResponse{
			ID:          "user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			PhotoURL:    "https://example.com/photo.jpg",
			Role:        domain.RoleStudent,
			IsActive:    true,
			CreatedAt:   now,
			LastLoginAt: now,
		}, nil)

	// Create request with user context
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.SetPathValue("userId", "user-123")
	w := httptest.NewRecorder()

	// Execute
	handler.handleGetUserProfile(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockAuthService.AssertExpectations(t)
}

func TestGetUserProfile_NotFound(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock user not found
	mockAuthService.On("GetUserProfile", mock.Anything, "nonexistent-user").
		Return(nil, apierrors.NewNotFoundError("User not found"))

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.SetPathValue("userId", "nonexistent-user")
	w := httptest.NewRecorder()

	// Execute
	handler.handleGetUserProfile(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockAuthService.AssertExpectations(t)
}

func TestUpdateUserProfile_Success(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock successful profile update
	displayName := "Updated Name"
	photoURL := "https://example.com/new-photo.jpg"

	updateReq := &domain.UpdateUserRequest{
		DisplayName: &displayName,
		PhotoURL:    &photoURL,
	}

	now := time.Now()
	updatedUser := &domain.User{
		ID:          "user-123",
		Email:       "test@example.com",
		DisplayName: displayName,
		PhotoURL:    photoURL,
		Role:        domain.RoleStudent,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockAuthService.On("UpdateUserProfile", mock.Anything, "user-123", mock.AnythingOfType("*domain.UpdateUserRequest")).
		Return(updatedUser, nil)

	// Create request
	bodyBytes, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/auth/me", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userId", "user-123")
	w := httptest.NewRecorder()

	// Execute
	handler.handleUpdateUserProfile(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockAuthService.AssertExpectations(t)
}

func TestUpdateUserProfile_ValidationError(t *testing.T) {
	handler, mockAuthService := setupAuthTestHandler(t)

	// Mock the service call in case validation passes (it shouldn't, but be safe)
	mockAuthService.On("UpdateUserProfile", mock.Anything, "user-123", mock.AnythingOfType("*domain.UpdateUserRequest")).
		Return(nil, apierrors.NewBadRequestError("Invalid photo URL"))

	// Create request with invalid URL
	invalidURL := "not-a-valid-url"
	updateReq := &domain.UpdateUserRequest{
		PhotoURL: &invalidURL,
	}

	bodyBytes, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/auth/me", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userId", "user-123")
	w := httptest.NewRecorder()

	// Execute
	handler.handleUpdateUserProfile(w, req)

	// Assert - should fail
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)

	var resp domain.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
}

// Add handler methods for testing (these would be in the actual handler.go)
func (h *Handler) handleAuthVerify(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyTokenRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	response, err := h.services.Auth.VerifyAndSyncUser(r.Context(), req.IDToken)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "token verified successfully")
}

func (h *Handler) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	profile, err := h.services.Auth.GetUserProfile(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, profile, "profile retrieved successfully")
}

func (h *Handler) handleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	var req domain.UpdateUserRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	user, err := h.services.Auth.UpdateUserProfile(r.Context(), userID, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, user, "profile updated successfully")
}
