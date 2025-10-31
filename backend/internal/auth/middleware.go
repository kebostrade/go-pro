// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package auth provides authentication and authorization middleware for the GO-PRO Learning Platform.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"go-pro-backend/pkg/logger"
	"go-pro-backend/security"

	apierrors "go-pro-backend/internal/errors"
)

// ContextKey represents a context key type.
type ContextKey string

const (
	// UserContextKey is the context key for user information.
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims.
	ClaimsContextKey ContextKey = "claims"
)

// AuthMiddleware provides JWT authentication middleware.
type AuthMiddleware struct {
	jwtManager *security.JWTManager
	logger     logger.Logger
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(jwtManager *security.JWTManager, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Authenticate middleware validates JWT tokens.
func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract token from Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.logger.Warn(ctx, "Missing Authorization header")
			a.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("missing authorization header"))

			return
		}

		// Validate Bearer token format.
		if !strings.HasPrefix(authHeader, "Bearer ") {
			a.logger.Warn(ctx, "Invalid Authorization header format")
			a.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("invalid authorization header format"))

			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			a.logger.Warn(ctx, "Empty JWT token")
			a.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("empty token"))

			return
		}

		// Validate token.
		claims, err := a.jwtManager.ValidateToken(token)
		if err != nil {
			a.logger.Warn(ctx, "Invalid JWT token", "error", err)
			a.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("invalid token"))

			return
		}

		// Add claims to context.
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, UserContextKey, &UserInfo{
			ID:    claims.UserID,
			Email: claims.Email,
			Roles: claims.Roles,
		})

		a.logger.Debug(ctx, "User authenticated successfully",
			"user_id", claims.UserID,
			"email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoles middleware checks if user has required roles.
func (a *AuthMiddleware) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get user from context.
			user, ok := GetUserFromContext(ctx)
			if !ok {
				a.logger.Warn(ctx, "User not found in context for role check")
				a.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("authentication required"))

				return
			}

			// Check if user has any of the required roles.
			if !hasAnyRole(user.Roles, roles) {
				a.logger.Warn(ctx, "User lacks required roles",
					"user_id", user.ID,
					"user_roles", user.Roles,
					"required_roles", roles)
				a.writeErrorResponse(w, r, apierrors.NewForbiddenError("insufficient permissions"))

				return
			}

			a.logger.Debug(ctx, "Role authorization successful",
				"user_id", user.ID,
				"required_roles", roles)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware checks if user has admin role.
func (a *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return a.RequireRoles("admin")(next)
}

// RequireInstructor middleware checks if user has instructor role.
func (a *AuthMiddleware) RequireInstructor(next http.Handler) http.Handler {
	return a.RequireRoles("instructor", "admin")(next)
}

// OptionalAuth middleware adds user info to context if token is present but doesn't require it.
func (a *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract token from Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication.
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without authentication.
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			// Empty token, continue without authentication.
			next.ServeHTTP(w, r)
			return
		}

		// Try to validate token.
		claims, err := a.jwtManager.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without authentication.
			a.logger.Debug(ctx, "Optional auth failed, continuing without authentication", "error", err)
			next.ServeHTTP(w, r)

			return
		}

		// Add claims to context.
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, UserContextKey, &UserInfo{
			ID:    claims.UserID,
			Email: claims.Email,
			Roles: claims.Roles,
		})

		a.logger.Debug(ctx, "Optional authentication successful", "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserInfo represents user information stored in context.
type UserInfo struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// GetUserFromContext extracts user information from context.
func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserInfo)
	return user, ok
}

// GetUserInfo extracts user information from context (returns nil if not found).
func GetUserInfo(ctx context.Context) *UserInfo {
	user, _ := GetUserFromContext(ctx)
	return user
}

// GetClaimsFromContext extracts JWT claims from context.
func GetClaimsFromContext(ctx context.Context) (*security.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*security.Claims)
	return claims, ok
}

// hasAnyRole checks if user has any of the required roles.
func hasAnyRole(userRoles, requiredRoles []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range userRoles {
		roleMap[role] = true
	}

	for _, required := range requiredRoles {
		if roleMap[required] {
			return true
		}
	}

	return false
}

// writeErrorResponse writes an error response.
func (a *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	var apiErr *apierrors.APIError
	if !errors.As(err, &apiErr) {
		apiErr = apierrors.NewInternalServerError("internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    apiErr.Code,
			"message": apiErr.Message,
			"details": apiErr.Details,
		},
		"success":   false,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error(ctx, "Failed to encode error response", "error", err)
	}
}

// Helper functions for testing and simple usage.

// userInfoKey is used for storing user info in context (for tests).
var userInfoKey = UserContextKey

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// JWTManager wraps security.JWTManager for testing convenience.
type JWTManager struct {
	*security.JWTManager
}

// NewJWTManager creates a JWT manager with simplified parameters for testing.
func NewJWTManager(secret []byte, ttl time.Duration, issuer string, logger logger.Logger) *JWTManager {
	config := security.JWTConfig{
		Secret:          secret,
		AccessTokenTTL:  ttl,
		RefreshTokenTTL: ttl * 10, // Refresh token lives longer
		Issuer:          issuer,
		Audience:        "test-audience",
	}

	return &JWTManager{
		JWTManager: security.NewJWTManager(config),
	}
}

// GenerateToken generates a token pair (convenience method for tests).
func (j *JWTManager) GenerateToken(userID, email, username string, roles []string) (*TokenPair, error) {
	accessToken, refreshToken, err := j.JWTManager.GenerateTokens(userID, email, roles)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Authenticate creates an authentication middleware.
func Authenticate(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract token from Authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Validate Bearer token format.
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "empty token", http.StatusUnauthorized)
				return
			}

			// Validate token.
			claims, err := jwtManager.JWTManager.ValidateToken(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to context.
			ctx = context.WithValue(ctx, UserContextKey, &UserInfo{
				ID:    claims.UserID,
				Email: claims.Email,
				Roles: claims.Roles,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles creates a middleware that checks if user has required roles.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			if !hasAnyRole(user.Roles, roles) {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin creates a middleware that checks if user has admin role.
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRoles("admin")
}

// OptionalAuth creates a middleware that adds user info to context if token is present.
func OptionalAuth(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract token from Authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Try to validate token.
			claims, err := jwtManager.JWTManager.ValidateToken(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Add user info to context.
			ctx = context.WithValue(ctx, UserContextKey, &UserInfo{
				ID:    claims.UserID,
				Email: claims.Email,
				Roles: claims.Roles,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
