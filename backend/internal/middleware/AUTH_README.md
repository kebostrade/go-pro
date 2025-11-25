# JWT Authentication Middleware

## Overview

Complete JWT authentication middleware for the Go-Pro learning platform with Firebase integration.

## Files

- **`auth.go`** (252 lines): Core authentication middleware implementation
- **`auth_test.go`** (271 lines): Comprehensive unit tests
- **`auth_example.go`** (209 lines): Usage examples and documentation

## Features

### ✅ AuthRequired Middleware
- Extracts and validates Bearer tokens from Authorization header
- Verifies Firebase ID tokens via AuthService
- Gets or creates users in backend database
- Adds authenticated user to request context
- Updates last login timestamp
- Checks user active status
- Returns standardized error responses (401/403/500)

### ✅ AdminRequired Middleware
- Checks if authenticated user has admin role
- Must be chained after AuthRequired
- Returns 403 if user lacks admin permissions

### ✅ Context Management
- `WithUser(ctx, user)`: Adds user to context
- `GetUserFromContext(ctx)`: Retrieves user from context
- Type-safe context key pattern

### ✅ Error Handling
- 401 Unauthorized: Missing/invalid/expired token, inactive account
- 403 Forbidden: Insufficient permissions
- 500 Internal Server Error: Database/authentication failures
- All errors logged with structured logging
- Follows existing `middleware.WriteErrorResponse` pattern

## Architecture

```
┌─────────────┐
│   Client    │
│  (Frontend) │
└──────┬──────┘
       │ Authorization: Bearer <firebase-token>
       ▼
┌─────────────────────────────────────┐
│     AuthRequired Middleware         │
├─────────────────────────────────────┤
│ 1. Extract Bearer token             │
│ 2. Verify with Firebase             │
│ 3. Get/Create user in backend       │
│ 4. Check active status              │
│ 5. Add user to context              │
│ 6. Update last login                │
└──────┬──────────────────────────────┘
       │ Context: user
       ▼
┌─────────────────────────────────────┐
│    AdminRequired (optional)         │
├─────────────────────────────────────┤
│ 1. Get user from context            │
│ 2. Check role == admin              │
│ 3. Return 403 if not admin          │
└──────┬──────────────────────────────┘
       │ Context: user (verified admin)
       ▼
┌─────────────────────────────────────┐
│         Handler                     │
│  GetUserFromContext(r.Context())    │
└─────────────────────────────────────┘
```

## Dependencies

- **AuthService**: Firebase token verification interface
- **UserRepository**: Backend user data access
- **Logger**: Structured logging
- **domain.User**: User entity model
- **internal/errors**: Standardized API errors

## Integration

### 1. Initialize Dependencies

```go
import (
    "go-pro-backend/internal/middleware"
    "go-pro-backend/internal/repository"
    "go-pro-backend/pkg/logger"
)

// Create logger
log := logger.New("info", "json")

// Initialize database and repositories
db := postgres.Connect(config)
userRepo := postgres.NewUserRepository(db)

// Initialize Firebase auth service (implement AuthService interface)
authService := firebase.NewAuthService(firebaseApp)

// Create auth middleware
authMiddleware := middleware.NewAuthMiddleware(authService, userRepo, log)
```

### 2. Protect Routes

```go
// Public route (no auth)
mux.Handle("/api/health", http.HandlerFunc(handleHealth))

// Protected route (auth required)
mux.Handle("/api/profile",
    middleware.Chain(
        http.HandlerFunc(handleProfile),
        authMiddleware.AuthRequired,
        middleware.Logging(log),
    ),
)

// Admin route (admin role required)
mux.Handle("/api/admin/users",
    middleware.Chain(
        http.HandlerFunc(handleAdminUsers),
        authMiddleware.AdminRequired,
        authMiddleware.AuthRequired,  // Must come after AdminRequired
        middleware.Logging(log),
    ),
)
```

### 3. Access User in Handler

```go
func handleProfile(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user
    user, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Use user data
    response := map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
        "role":  user.Role,
    }

    json.NewEncoder(w).Encode(response)
}
```

## Firebase AuthService Implementation

The middleware expects an `AuthService` interface for Firebase integration:

```go
type AuthService interface {
    VerifyToken(ctx context.Context, token string) (*FirebaseToken, error)
}

type FirebaseToken struct {
    UID         string
    Email       string
    DisplayName string
    PhotoURL    string
}
```

Example implementation using Firebase Admin SDK:

```go
import (
    "context"
    firebase "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/auth"
)

type FirebaseAuthService struct {
    client *auth.Client
}

func NewFirebaseAuthService(app *firebase.App) (*FirebaseAuthService, error) {
    client, err := app.Auth(context.Background())
    if err != nil {
        return nil, err
    }
    return &FirebaseAuthService{client: client}, nil
}

func (s *FirebaseAuthService) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
    // Verify ID token with Firebase
    firebaseToken, err := s.client.VerifyIDToken(ctx, token)
    if err != nil {
        return nil, fmt.Errorf("failed to verify token: %w", err)
    }

    // Extract user info from claims
    email, _ := firebaseToken.Claims["email"].(string)
    name, _ := firebaseToken.Claims["name"].(string)
    picture, _ := firebaseToken.Claims["picture"].(string)

    return &middleware.FirebaseToken{
        UID:         firebaseToken.UID,
        Email:       email,
        DisplayName: name,
        PhotoURL:    picture,
    }, nil
}
```

## Request/Response Format

### Successful Authentication

**Request:**
```http
GET /api/profile HTTP/1.1
Host: api.gopro.dev
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "user-123",
    "email": "student@gopro.dev",
    "role": "student",
    "display_name": "John Doe"
  }
}
```

### Authentication Failures

**Missing Token (401):**
```json
{
  "success": false,
  "error": {
    "type": "UNAUTHORIZED",
    "message": "missing authorization header"
  },
  "request_id": "req-abc-123",
  "timestamp": "2025-01-24T12:00:00Z"
}
```

**Invalid Token (401):**
```json
{
  "success": false,
  "error": {
    "type": "UNAUTHORIZED",
    "message": "invalid or expired token"
  },
  "request_id": "req-abc-123",
  "timestamp": "2025-01-24T12:00:00Z"
}
```

**Insufficient Permissions (403):**
```json
{
  "success": false,
  "error": {
    "type": "FORBIDDEN",
    "message": "admin access required"
  },
  "request_id": "req-abc-123",
  "timestamp": "2025-01-24T12:00:00Z"
}
```

## User Creation Flow

1. **First Request**: User authenticates with Firebase on frontend
2. **Token Sent**: Frontend sends Firebase ID token to backend
3. **Token Verified**: Middleware verifies token with Firebase
4. **User Lookup**: Middleware checks if user exists by Firebase UID
5. **User Created**: If not exists, creates new user with:
   - Firebase UID (unique identifier)
   - Email, display name, photo URL from Firebase
   - Default role: `student`
   - Active status: `true`
6. **Context Added**: User object added to request context
7. **Login Updated**: Last login timestamp updated asynchronously

## Testing

Run tests with coverage:

```bash
go test -v -cover ./internal/middleware
```

Test output:
```
=== RUN   TestAuthRequired_Success
--- PASS: TestAuthRequired_Success (0.00s)
=== RUN   TestAuthRequired_MissingToken
--- PASS: TestAuthRequired_MissingToken (0.00s)
=== RUN   TestAdminRequired_Success
--- PASS: TestAdminRequired_Success (0.00s)
=== RUN   TestAdminRequired_Forbidden
--- PASS: TestAdminRequired_Forbidden (0.00s)
=== RUN   TestGetUserFromContext
--- PASS: TestGetUserFromContext (0.00s)
=== RUN   TestGetUserFromContext_NoUser
--- PASS: TestGetUserFromContext_NoUser (0.00s)
PASS
ok  	go-pro-backend/internal/middleware	0.006s
```

## Security Considerations

### ✅ Implemented
- Bearer token format strictly enforced
- Firebase token verification on every request (stateless)
- User active status checked on every request
- All authentication failures logged with context
- Atomic user creation with race condition handling
- Standardized error responses (no information leakage)
- Context-based user propagation (no global state)

### 🔒 Best Practices
- Tokens verified with Firebase (not stored in backend)
- Short-lived tokens (Firebase default: 1 hour)
- HTTPS required in production
- Rate limiting applied via separate middleware
- Logging excludes sensitive token data

### ⚠️ Important Notes
- Firebase Admin SDK credentials must be secured
- Database connection must be secure (SSL/TLS)
- Consider implementing token refresh endpoint
- Monitor Firebase quota and costs
- Implement proper CORS configuration

## Logging

The middleware provides structured logging for all authentication events:

```
INFO: user authenticated
  user_id=user-123
  email=student@gopro.dev
  role=student

INFO: admin access granted
  user_id=admin-456
  email=admin@gopro.dev

ERROR: invalid or expired token
  error="firebase token verification failed"

ERROR: insufficient permissions
  user_id=user-123
  role=student

WARN: failed to update last login
  user_id=user-123
  error="database connection timeout"
```

## Patterns Used

- **Middleware Pattern**: HTTP handler wrapper for cross-cutting concerns
- **Repository Pattern**: Abstract data access via UserRepository interface
- **Context Pattern**: Type-safe request-scoped user propagation
- **Interface Segregation**: AuthService interface for Firebase abstraction
- **Error Wrapping**: Structured errors via internal/errors package
- **Dependency Injection**: Constructor injection for testability

## Future Enhancements

- [ ] Token refresh endpoint
- [ ] Role-based permissions beyond admin/student
- [ ] API key authentication for service-to-service
- [ ] Rate limiting per user (not just per IP)
- [ ] Audit logging for admin actions
- [ ] User session management
- [ ] Multi-factor authentication support
- [ ] OAuth provider integration beyond Firebase

## Related Documentation

- See `auth_example.go` for complete usage examples
- See `middleware.go` for other available middleware
- See `internal/domain/models.go` for User entity definition
- See `internal/repository/interfaces.go` for UserRepository interface
- See `pkg/logger/` for logging utilities
