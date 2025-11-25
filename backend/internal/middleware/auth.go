// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package middleware provides functionality for the GO-PRO Learning Platform.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"

	apierrors "go-pro-backend/internal/errors"
)

// Context keys for user information.
type authContextKey string

const (
	userContextKey authContextKey = "user"
)

// Error messages for authentication.
const (
	errMissingAuthHeader       = "missing authorization header"
	errInvalidAuthHeaderFormat = "invalid authorization header format"
	errUserAccountDisabled     = "user account is disabled"
)

// AuthService defines the interface for Firebase authentication operations.
type AuthService interface {
	VerifyToken(ctx context.Context, token string) (*FirebaseToken, error)
}

// FirebaseToken represents a verified Firebase token.
type FirebaseToken struct {
	UID         string
	Email       string
	DisplayName string
	PhotoURL    string
}

// AuthMiddleware handles authentication and authorization.
type AuthMiddleware struct {
	authService AuthService
	userRepo    repository.UserRepository
	log         logger.Logger
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(authService AuthService, userRepo repository.UserRepository, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		userRepo:    userRepo,
		log:         log,
	}
}

// AuthRequired verifies the Firebase token and adds the user to the request context.
func (am *AuthMiddleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.LogError(am.log, r.Context(),
				apierrors.NewUnauthorizedError(errMissingAuthHeader),
				errMissingAuthHeader,
			)
			WriteErrorResponse(w, r, apierrors.NewUnauthorizedError(errMissingAuthHeader))
			return
		}

		// Extract Bearer token.
		token := extractBearerToken(authHeader)
		if token == "" {
			logger.LogError(am.log, r.Context(),
				apierrors.NewUnauthorizedError(errInvalidAuthHeaderFormat),
				errInvalidAuthHeaderFormat,
			)
			WriteErrorResponse(w, r, apierrors.NewUnauthorizedError(errInvalidAuthHeaderFormat))
			return
		}

		// Verify Firebase token.
		firebaseToken, err := am.authService.VerifyToken(r.Context(), token)
		if err != nil {
			logger.LogError(am.log, r.Context(),
				apierrors.NewUnauthorizedError("invalid or expired token"),
				"firebase token verification failed",
				"error", err.Error(),
			)
			WriteErrorResponse(w, r, apierrors.NewUnauthorizedError("invalid or expired token"))
			return
		}

		// Get or create user in backend.
		user, err := am.getOrCreateUser(r.Context(), firebaseToken)
		if err != nil {
			logger.LogError(am.log, r.Context(),
				apierrors.NewInternalError("failed to get or create user", err),
				"failed to get or create user",
				"firebase_uid", firebaseToken.UID,
				"error", err.Error(),
			)
			WriteErrorResponse(w, r, apierrors.NewInternalError("authentication failed", err))
			return
		}

		// Check if user is active.
		if !user.IsActive {
			logger.LogError(am.log, r.Context(),
				apierrors.NewForbiddenError(errUserAccountDisabled),
				errUserAccountDisabled,
				"user_id", user.ID,
			)
			WriteErrorResponse(w, r, apierrors.NewForbiddenError(errUserAccountDisabled))
			return
		}

		// Update last login time.
		if err := am.userRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
			// Log error but don't fail the request.
			am.log.Warn(r.Context(), "failed to update last login",
				"user_id", user.ID,
				"error", err.Error(),
			)
		}

		// Add user to request context.
		ctx := WithUser(r.Context(), user)

		// Log successful authentication.
		am.log.Info(r.Context(), "user authenticated",
			"user_id", user.ID,
			"email", user.Email,
			"role", user.Role,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired checks if the authenticated user has admin role.
// Must be used after AuthRequired middleware.
func (am *AuthMiddleware) AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context.
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			logger.LogError(am.log, r.Context(),
				apierrors.NewUnauthorizedError("user not authenticated"),
				"user not found in context",
			)
			WriteErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
			return
		}

		// Check if user has admin role.
		if user.Role != domain.RoleAdmin {
			logger.LogError(am.log, r.Context(),
				apierrors.NewForbiddenError("admin access required"),
				"insufficient permissions",
				"user_id", user.ID,
				"role", user.Role,
			)
			WriteErrorResponse(w, r, apierrors.NewForbiddenError("admin access required"))
			return
		}

		// Log admin access.
		am.log.Info(r.Context(), "admin access granted",
			"user_id", user.ID,
			"email", user.Email,
		)

		next.ServeHTTP(w, r)
	})
}

// getOrCreateUser retrieves existing user or creates new one from Firebase token.
func (am *AuthMiddleware) getOrCreateUser(ctx context.Context, token *FirebaseToken) (*domain.User, error) {
	// Try to get existing user by Firebase UID.
	user, err := am.userRepo.GetByFirebaseUID(ctx, token.UID)
	if err == nil {
		return user, nil
	}

	// Check if error is "not found".
	apiErr, ok := apierrors.IsAPIError(err)
	if ok && apiErr.StatusCode == http.StatusNotFound {
		// User doesn't exist, create new one.
		return am.createUser(ctx, token)
	}

	// Other error occurred.
	return nil, err
}

// createUser creates a new user from Firebase token.
func (am *AuthMiddleware) createUser(ctx context.Context, token *FirebaseToken) (*domain.User, error) {
	user := &domain.User{
		FirebaseUID: token.UID,
		Email:       token.Email,
		DisplayName: token.DisplayName,
		PhotoURL:    token.PhotoURL,
		Username:    generateUsernameFromEmail(token.Email),
		Role:        domain.RoleStudent, // Default role
		IsActive:    true,
	}

	if err := am.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	am.log.Info(ctx, "new user created",
		"user_id", user.ID,
		"email", user.Email,
		"firebase_uid", user.FirebaseUID,
	)

	return user, nil
}

// Context helper functions.

// WithUser adds a user to the context.
func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext retrieves the user from the context.
func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}

// Utility functions.

// extractBearerToken extracts the token from the Authorization header.
func extractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// generateUsernameFromEmail generates a username from email address.
func generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}
