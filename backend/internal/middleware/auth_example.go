// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package middleware

// Example usage of authentication middleware.
//
// Basic Setup:
//
//	import (
//		"go-pro-backend/internal/middleware"
//		"go-pro-backend/internal/repository"
//		"go-pro-backend/pkg/logger"
//	)
//
//	// Initialize dependencies
//	log := logger.New("info", "json")
//	userRepo := repository.NewUserRepository(db)
//	authService := firebase.NewAuthService() // Your Firebase auth service implementation
//
//	// Create auth middleware
//	authMiddleware := middleware.NewAuthMiddleware(authService, userRepo, log)
//
// Protected Routes (Authenticated Users Only):
//
//	// Chain middleware for protected endpoints
//	protectedHandler := middleware.Chain(
//		http.HandlerFunc(handleProtectedEndpoint),
//		authMiddleware.AuthRequired,
//		middleware.Logging(log),
//	)
//
//	mux.Handle("/api/profile", protectedHandler)
//
// Admin-Only Routes:
//
//	// Chain both AuthRequired and AdminRequired
//	adminHandler := middleware.Chain(
//		http.HandlerFunc(handleAdminEndpoint),
//		authMiddleware.AdminRequired,
//		authMiddleware.AuthRequired,  // AuthRequired must come after AdminRequired
//		middleware.Logging(log),
//	)
//
//	mux.Handle("/api/admin/users", adminHandler)
//
// Accessing User in Handler:
//
//	func handleProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
//		// Get authenticated user from context
//		user, ok := middleware.GetUserFromContext(r.Context())
//		if !ok {
//			// This should never happen if AuthRequired middleware is properly set up
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		// Use user information
//		fmt.Fprintf(w, "Welcome %s (ID: %s, Role: %s)", user.Email, user.ID, user.Role)
//	}
//
// Complete Example with All Middleware:
//
//	func setupRouter(authMiddleware *middleware.AuthMiddleware, log logger.Logger) *http.ServeMux {
//		mux := http.NewServeMux()
//
//		// Public endpoints (no auth required)
//		mux.Handle("/api/health", http.HandlerFunc(handleHealth))
//		mux.Handle("/api/curriculum", http.HandlerFunc(handleCurriculum))
//
//		// Protected endpoints (auth required)
//		mux.Handle("/api/profile",
//			middleware.Chain(
//				http.HandlerFunc(handleProfile),
//				authMiddleware.AuthRequired,
//				middleware.Logging(log),
//			),
//		)
//
//		mux.Handle("/api/progress",
//			middleware.Chain(
//				http.HandlerFunc(handleProgress),
//				authMiddleware.AuthRequired,
//				middleware.Logging(log),
//			),
//		)
//
//		// Admin endpoints (admin role required)
//		mux.Handle("/api/admin/users",
//			middleware.Chain(
//				http.HandlerFunc(handleAdminUsers),
//				authMiddleware.AdminRequired,
//				authMiddleware.AuthRequired,
//				middleware.Logging(log),
//			),
//		)
//
//		mux.Handle("/api/admin/courses",
//			middleware.Chain(
//				http.HandlerFunc(handleAdminCourses),
//				authMiddleware.AdminRequired,
//				authMiddleware.AuthRequired,
//				middleware.Logging(log),
//			),
//		)
//
//		return mux
//	}
//
// Error Responses:
//
// The middleware automatically returns standardized error responses:
//
//	401 Unauthorized:
//	  - Missing Authorization header
//	  - Invalid token format (not "Bearer <token>")
//	  - Expired or invalid Firebase token
//	  - User account disabled
//
//	403 Forbidden:
//	  - User lacks required admin role
//	  - User account is not active
//
//	500 Internal Server Error:
//	  - Database errors during user lookup/creation
//	  - Unexpected authentication errors
//
// Client Request Format:
//
// Clients must include Firebase ID token in Authorization header:
//
//	GET /api/profile HTTP/1.1
//	Host: api.gopro.dev
//	Authorization: Bearer <firebase-id-token>
//	Content-Type: application/json
//
// Response Format:
//
// Success responses include user data from context:
//	{
//	  "success": true,
//	  "data": {
//	    "id": "user-123",
//	    "email": "user@example.com",
//	    "role": "student"
//	  }
//	}
//
// Error responses follow standardized format:
//	{
//	  "success": false,
//	  "error": {
//	    "type": "UNAUTHORIZED",
//	    "message": "invalid or expired token"
//	  },
//	  "request_id": "req-abc-123",
//	  "timestamp": "2025-01-24T12:00:00Z"
//	}
//
// Firebase Integration:
//
// The AuthService interface expects Firebase token verification:
//
//	type AuthService interface {
//	  VerifyToken(ctx context.Context, token string) (*FirebaseToken, error)
//	}
//
// Example Firebase implementation:
//
//	type FirebaseAuthService struct {
//	  client *auth.Client
//	}
//
//	func (s *FirebaseAuthService) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
//	  // Verify token with Firebase Admin SDK
//	  firebaseToken, err := s.client.VerifyIDToken(ctx, token)
//	  if err != nil {
//	    return nil, err
//	  }
//
//	  return &middleware.FirebaseToken{
//	    UID:         firebaseToken.UID,
//	    Email:       firebaseToken.Claims["email"].(string),
//	    DisplayName: firebaseToken.Claims["name"].(string),
//	    PhotoURL:    firebaseToken.Claims["picture"].(string),
//	  }, nil
//	}
//
// User Creation Flow:
//
// 1. User signs in with Firebase on frontend
// 2. Frontend receives Firebase ID token
// 3. Frontend sends token to backend in Authorization header
// 4. AuthRequired middleware verifies token with Firebase
// 5. Middleware checks if user exists in backend database
// 6. If not exists, creates new user with Firebase UID
// 7. User is added to request context
// 8. Handler can access user via GetUserFromContext
//
// Security Considerations:
//
// - Tokens are verified with Firebase on every request (stateless)
// - User active status is checked on every request
// - Last login time is updated asynchronously (non-blocking)
// - All auth failures are logged with context
// - Bearer token format is strictly enforced
// - User creation is atomic and handles race conditions
//
