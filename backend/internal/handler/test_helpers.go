// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/middleware"
)

// TestHelper provides common test utilities.
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper.
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateMockUser creates a mock user for testing.
func (h *TestHelper) CreateMockUser(role domain.UserRole, isActive bool) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:          "test-user-" + string(role),
		FirebaseUID: "firebase-" + string(role),
		Email:       string(role) + "@test.com",
		Username:    string(role),
		DisplayName: string(role) + " User",
		PhotoURL:    "https://example.com/" + string(role) + ".jpg",
		Role:        role,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}
}

// CreateMockFirebaseClaims creates mock Firebase claims for testing.
func (h *TestHelper) CreateMockFirebaseClaims(email string) *domain.FirebaseClaims {
	return &domain.FirebaseClaims{
		UserID:      "firebase-uid-" + email,
		Email:       email,
		DisplayName: "Test User",
		Picture:     "https://example.com/photo.jpg",
		IssuedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour),
	}
}

// AddUserToContext adds a user to the request context.
func (h *TestHelper) AddUserToContext(r *http.Request, user *domain.User) *http.Request {
	ctx := middleware.WithUser(r.Context(), user)
	return r.WithContext(ctx)
}

// CreateAuthHeader creates an Authorization header with Bearer token.
func (h *TestHelper) CreateAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// MockFirebaseAuthClient mocks the Firebase Auth Client for middleware testing.
type MockFirebaseAuthClient struct {
	mock.Mock
}

func (m *MockFirebaseAuthClient) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*middleware.FirebaseToken), args.Error(1)
}

// MockUserRepository provides a simple in-memory user repository for testing.
type MockUserRepositorySimple struct {
	users map[string]*domain.User // ID -> User
}

// NewMockUserRepositorySimple creates a new mock user repository.
func NewMockUserRepositorySimple() *MockUserRepositorySimple {
	return &MockUserRepositorySimple{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepositorySimple) Create(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = "user-" + user.FirebaseUID
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepositorySimple) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, apierrors.NewNotFoundError("User not found")
	}
	return user, nil
}

func (m *MockUserRepositorySimple) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, apierrors.NewNotFoundError("User not found")
}

func (m *MockUserRepositorySimple) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	for _, user := range m.users {
		if user.FirebaseUID == firebaseUID {
			return user, nil
		}
	}
	return nil, apierrors.NewNotFoundError("User not found")
}

func (m *MockUserRepositorySimple) Update(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return apierrors.NewNotFoundError("User not found")
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepositorySimple) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return apierrors.NewNotFoundError("User not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepositorySimple) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	total := int64(len(users))

	// Apply pagination
	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(users) {
			return []*domain.User{}, total, nil
		}

		if end > len(users) {
			end = len(users)
		}

		users = users[start:end]
	}

	return users, total, nil
}

func (m *MockUserRepositorySimple) UpdateLastLogin(ctx context.Context, id string) error {
	user, exists := m.users[id]
	if !exists {
		return apierrors.NewNotFoundError("User not found")
	}
	now := time.Now()
	user.LastLoginAt = &now
	return nil
}

// TestScenarios provides common test scenarios.
type TestScenarios struct{}

// NewTestScenarios creates a new test scenarios helper.
func NewTestScenarios() *TestScenarios {
	return &TestScenarios{}
}

// CreateFirstUserScenario sets up a scenario where the first user registers.
func (s *TestScenarios) CreateFirstUserScenario() (*domain.User, *domain.VerifyTokenResponse) {
	now := time.Now()
	user := &domain.User{
		ID:          "user-001",
		FirebaseUID: "firebase-001",
		Email:       "first@example.com",
		DisplayName: "First User",
		Role:        domain.RoleAdmin, // First user is admin
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}

	response := &domain.VerifyTokenResponse{
		User:          user.ToProfileResponse(),
		IsNewUser:     true,
		TokenVerified: true,
	}

	return user, response
}

// CreateSubsequentUserScenario sets up a scenario where a non-first user registers.
func (s *TestScenarios) CreateSubsequentUserScenario() (*domain.User, *domain.VerifyTokenResponse) {
	now := time.Now()
	user := &domain.User{
		ID:          "user-002",
		FirebaseUID: "firebase-002",
		Email:       "student@example.com",
		DisplayName: "Student User",
		Role:        domain.RoleStudent, // Subsequent users are students
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}

	response := &domain.VerifyTokenResponse{
		User:          user.ToProfileResponse(),
		IsNewUser:     true,
		TokenVerified: true,
	}

	return user, response
}

// AssertAPIResponse provides assertions for API responses.
type AssertAPIResponse struct {
	t        *testing.T
	response *domain.APIResponse
}

// NewAssertAPIResponse creates a new API response asserter.
func NewAssertAPIResponse(t *testing.T, resp *domain.APIResponse) *AssertAPIResponse {
	return &AssertAPIResponse{
		t:        t,
		response: resp,
	}
}

// Success asserts that the response is successful.
func (a *AssertAPIResponse) Success() *AssertAPIResponse {
	if !a.response.Success {
		a.t.Errorf("Expected successful response, got error: %v", a.response.Error)
	}
	return a
}

// Error asserts that the response is an error.
func (a *AssertAPIResponse) Error() *AssertAPIResponse {
	if a.response.Success {
		a.t.Error("Expected error response, got success")
	}
	return a
}

// HasErrorMessage asserts that the error message contains the expected text.
func (a *AssertAPIResponse) HasErrorMessage(expected string) *AssertAPIResponse {
	if a.response.Error == nil {
		a.t.Error("Expected error in response, got nil")
		return a
	}

	if !contains(a.response.Error.Message, expected) {
		a.t.Errorf("Expected error message to contain '%s', got '%s'", expected, a.response.Error.Message)
	}
	return a
}

// HasData asserts that the response has data.
func (a *AssertAPIResponse) HasData() *AssertAPIResponse {
	if a.response.Data == nil {
		a.t.Error("Expected data in response, got nil")
	}
	return a
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TokenGenerator generates test tokens.
type TokenGenerator struct{}

// NewTokenGenerator creates a new token generator.
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// ValidAdminToken generates a valid admin token.
func (g *TokenGenerator) ValidAdminToken() string {
	return "valid-admin-token-" + time.Now().Format("20060102150405")
}

// ValidStudentToken generates a valid student token.
func (g *TokenGenerator) ValidStudentToken() string {
	return "valid-student-token-" + time.Now().Format("20060102150405")
}

// InvalidToken generates an invalid token.
func (g *TokenGenerator) InvalidToken() string {
	return "invalid-token-" + time.Now().Format("20060102150405")
}

// ExpiredToken generates an expired token.
func (g *TokenGenerator) ExpiredToken() string {
	return "expired-token-" + time.Now().Format("20060102150405")
}
