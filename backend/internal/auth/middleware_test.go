package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pro-backend/internal/testutil"
)

func TestAuthenticate(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	jwtManager := NewJWTManager([]byte("test-secret"), 15*time.Minute, "test-issuer", logger)

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectUserInfo bool
	}{
		{
			name: "valid token",
			setupRequest: func(r *http.Request) {
				token, _ := jwtManager.GenerateToken("user-1", "test@example.com", "testuser", []string{"student"})
				r.Header.Set("Authorization", "Bearer "+token.AccessToken)
			},
			expectedStatus: http.StatusOK,
			expectUserInfo: true,
		},
		{
			name: "missing authorization header",
			setupRequest: func(r *http.Request) {
				// No authorization header
			},
			expectedStatus: http.StatusUnauthorized,
			expectUserInfo: false,
		},
		{
			name: "invalid token format",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat token")
			},
			expectedStatus: http.StatusUnauthorized,
			expectUserInfo: false,
		},
		{
			name: "invalid token",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			expectedStatus: http.StatusUnauthorized,
			expectUserInfo: false,
		},
		{
			name: "expired token",
			setupRequest: func(r *http.Request) {
				// Create a JWT manager with very short expiration
				shortJWT := NewJWTManager([]byte("test-secret"), 1*time.Nanosecond, "test-issuer", logger)
				token, _ := shortJWT.GenerateToken("user-1", "test@example.com", "testuser", []string{"student"})
				time.Sleep(10 * time.Millisecond) // Wait for token to expire
				r.Header.Set("Authorization", "Bearer "+token.AccessToken)
			},
			expectedStatus: http.StatusUnauthorized,
			expectUserInfo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectUserInfo {
					userInfo := GetUserInfo(r.Context())
					assert.NotNil(t, userInfo)
					assert.Equal(t, "user-1", userInfo.ID)
					assert.Equal(t, "test@example.com", userInfo.Email)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with authentication middleware
			middleware := Authenticate(jwtManager)
			wrappedHandler := middleware(handler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupRequest(req)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequireRoles(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name:           "user has required role",
			userRoles:      []string{"student", "instructor"},
			requiredRoles:  []string{"student"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user has one of multiple required roles",
			userRoles:      []string{"student"},
			requiredRoles:  []string{"student", "instructor"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user missing required role",
			userRoles:      []string{"student"},
			requiredRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "user has no roles",
			userRoles:      []string{},
			requiredRoles:  []string{"student"},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create user context
			userInfo := &UserInfo{
				ID:    "user-1",
				Email: "test@example.com",
				Roles: tt.userRoles,
			}

			// Wrap with role middleware
			middleware := RequireRoles(tt.requiredRoles...)
			wrappedHandler := middleware(handler)

			// Create request with user context
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), userInfoKey, userInfo)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		expectedStatus int
	}{
		{
			name:           "user is admin",
			userRoles:      []string{"admin"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user is not admin",
			userRoles:      []string{"student"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "user has multiple roles including admin",
			userRoles:      []string{"student", "instructor", "admin"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create user context
			userInfo := &UserInfo{
				ID:    "user-1",
				Email: "test@example.com",
				Roles: tt.userRoles,
			}

			// Wrap with admin middleware
			middleware := RequireAdmin()
			wrappedHandler := middleware(handler)

			// Create request with user context
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), userInfoKey, userInfo)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestOptionalAuth(t *testing.T) {
	logger := testutil.NewTestLogger(t)
	jwtManager := NewJWTManager([]byte("test-secret"), 15*time.Minute, "test-issuer", logger)

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectUserInfo bool
	}{
		{
			name: "valid token provided",
			setupRequest: func(r *http.Request) {
				token, _ := jwtManager.GenerateToken("user-1", "test@example.com", "testuser", []string{"student"})
				r.Header.Set("Authorization", "Bearer "+token.AccessToken)
			},
			expectedStatus: http.StatusOK,
			expectUserInfo: true,
		},
		{
			name: "no token provided",
			setupRequest: func(r *http.Request) {
				// No authorization header
			},
			expectedStatus: http.StatusOK,
			expectUserInfo: false,
		},
		{
			name: "invalid token provided",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			expectedStatus: http.StatusOK,
			expectUserInfo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userInfo := GetUserInfo(r.Context())
				if tt.expectUserInfo {
					assert.NotNil(t, userInfo)
				} else {
					assert.Nil(t, userInfo)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with optional auth middleware
			middleware := OptionalAuth(jwtManager)
			wrappedHandler := middleware(handler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupRequest(req)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() context.Context
		expected *UserInfo
	}{
		{
			name: "user info exists in context",
			setupCtx: func() context.Context {
				userInfo := &UserInfo{
					ID:    "user-1",
					Email: "test@example.com",
					Roles: []string{"student"},
				}
				return context.WithValue(context.Background(), userInfoKey, userInfo)
			},
			expected: &UserInfo{
				ID:    "user-1",
				Email: "test@example.com",
				Roles: []string{"student"},
			},
		},
		{
			name: "user info does not exist in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			userInfo := GetUserInfo(ctx)

			if tt.expected == nil {
				assert.Nil(t, userInfo)
			} else {
				require.NotNil(t, userInfo)
				assert.Equal(t, tt.expected.ID, userInfo.ID)
				assert.Equal(t, tt.expected.Email, userInfo.Email)
				assert.Equal(t, tt.expected.Roles, userInfo.Roles)
			}
		})
	}
}
