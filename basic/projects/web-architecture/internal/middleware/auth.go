package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/service"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/pkg/response"
)

type contextKey string

const UserContextKey contextKey = "user"

// Auth is a middleware that validates JWT tokens
func Auth(userService *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing authorization header")
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "Invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			claims, err := userService.ValidateToken(token)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole is a middleware that checks if the user has a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(map[string]interface{})
			if !ok {
				response.Unauthorized(w, "Unauthorized")
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				response.Forbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

