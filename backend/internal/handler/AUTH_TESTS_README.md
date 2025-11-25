# Firebase Authentication Integration Tests

## Overview
Comprehensive test suite for Firebase authentication flow covering auth endpoints, admin operations, and end-to-end scenarios.

## Test Structure

### 1. Handler Unit Tests (`auth_test.go`)
Tests individual authentication endpoints with mocked dependencies:

#### Test Cases
- **TestAuthVerify_Success**: Valid Firebase token verification
- **TestAuthVerify_InvalidToken**: Invalid/expired token handling
- **TestAuthVerify_MissingToken**: Missing token validation
- **TestAuthVerify_FirstUserBecomesAdmin**: First user receives admin role
- **TestGetUserProfile_Success**: Profile retrieval for valid user
- **TestGetUserProfile_NotFound**: Non-existent user handling
- **TestUpdateUserProfile_Success**: Profile update with valid data
- **TestUpdateUserProfile_ValidationError**: Invalid data handling

#### Mock Implementations
- **MockAuthService**: Mocks `service.AuthService` interface
  - VerifyToken, GetOrCreateUser, VerifyAndSyncUser
  - GetUserProfile, UpdateUserRole, UpdateUserProfile

### 2. Admin Endpoint Tests (`admin_test.go`)
Tests admin-specific operations and authorization:

#### Test Cases
- **TestGetAllUsers_AdminSuccess**: Admin can list all users
- **TestGetAllUsers_StudentForbidden**: Students cannot access admin endpoints
- **TestUpdateUserRole_AdminSuccess**: Admin can update user roles
- **TestUpdateUserRole_StudentForbidden**: Students cannot update roles
- **TestUpdateUserRole_AdminCannotDemoteThemselves**: Self-demotion prevention
- **TestDeleteUser_AdminSuccess**: Admin can delete users
- **TestDeleteUser_StudentForbidden**: Students cannot delete users

#### Mock Implementations
- **MockUserRepository**: In-memory user repository for testing
  - CRUD operations: Create, GetByID, GetByEmail, GetByFirebaseUID
  - List operations: GetAll with pagination
  - Update operations: Update, UpdateLastLogin, Delete

### 3. Integration Tests (`integration/auth_flow_test.go`)
End-to-end tests covering complete authentication workflows:

#### Test Suites
- **AuthFlowTestSuite**: Comprehensive E2E test suite

#### Test Scenarios
1. **Complete Auth Flow**:
   - First user registration (becomes admin)
   - Second user registration (becomes student)
   - Student attempts admin endpoint (forbidden)
   - Admin updates student role (success)
   - Admin self-demotion (prevented)

2. **Protected Endpoints**:
   - Unauthorized access denied for all protected routes
   - Invalid token rejected
   - Valid token grants access

3. **User Profile Updates**:
   - Profile update with valid data
   - Validation error handling

4. **Inactive User Handling**:
   - Inactive users denied access to protected endpoints

## Test Helpers (`test_helpers.go`)

### TestHelper
Provides common test utilities:
- `CreateMockUser`: Generate test users with specified roles
- `CreateMockFirebaseClaims`: Generate Firebase token claims
- `AddUserToContext`: Inject user into request context
- `CreateAuthHeader`: Generate Authorization headers

### MockUserRepositorySimple
In-memory implementation for testing:
- No external dependencies
- Fast test execution
- Full CRUD operations

### TestScenarios
Common test scenarios:
- `CreateFirstUserScenario`: Admin registration scenario
- `CreateSubsequentUserScenario`: Student registration scenario

### AssertAPIResponse
Fluent assertion interface:
- `Success()`: Assert successful response
- `Error()`: Assert error response
- `HasErrorMessage(expected)`: Verify error message
- `HasData()`: Verify response has data

### TokenGenerator
Test token generation:
- `ValidAdminToken()`: Admin authentication token
- `ValidStudentToken()`: Student authentication token
- `InvalidToken()`: Invalid/malformed token
- `ExpiredToken()`: Expired token

## Running Tests

### All Authentication Tests
```bash
cd backend
go test ./internal/handler -v -run TestAuth
```

### Specific Test Suite
```bash
# Handler tests only
go test ./internal/handler -v -run TestAuthVerify

# Admin tests only
go test ./internal/handler -v -run TestAdmin

# Integration tests
go test ./internal/integration -v
```

### With Coverage
```bash
# Handler tests with coverage
go test ./internal/handler -v -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

### Race Detection
```bash
go test ./internal/handler -v -race
```

## Integration with Actual Code

### Required Handler Methods
The following methods need to be implemented or exposed in `handler.go`:

```go
// Auth endpoints
func (h *Handler) HandleAuthVerify(w http.ResponseWriter, r *http.Request)
func (h *Handler) HandleGetUserProfile(w http.ResponseWriter, r *http.Request)
func (h *Handler) HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request)

// Admin endpoints
func (h *Handler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request)
func (h *Handler) HandleUpdateUserRole(w http.ResponseWriter, r *http.Request)
func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request)
```

### Required Service Interface
The `AuthService` interface must include:

```go
type AuthService interface {
    Initialize(ctx context.Context) error
    VerifyToken(ctx context.Context, idToken string) (*domain.FirebaseClaims, error)
    GetOrCreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error)
    VerifyAndSyncUser(ctx context.Context, idToken string) (*domain.VerifyTokenResponse, error)
    GetUserProfile(ctx context.Context, userID string) (*domain.UserProfileResponse, error)
    UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.UserRole) error
    UpdateUserProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.User, error)
}
```

## Test Patterns

### Setup Pattern
```go
func setupAuthTestHandler(t *testing.T) (*Handler, *MockAuthService) {
    mockAuthService := new(MockAuthService)
    services := &service.Services{Auth: mockAuthService}
    log := logger.New("debug", "json")
    val := validator.New()
    handler := New(services, log, val)
    return handler, mockAuthService
}
```

### Mock Pattern
```go
mockAuthService.On("VerifyAndSyncUser", mock.Anything, "valid-token").
    Return(&domain.VerifyTokenResponse{...}, nil)
```

### Request/Response Pattern
```go
req := httptest.NewRequest("POST", "/api/v1/auth/verify", body)
w := httptest.NewRecorder()
handler.HandleAuthVerify(w, req)
assert.Equal(t, http.StatusOK, w.Code)
```

## Coverage Goals

### Current Coverage
- Handler unit tests: Auth endpoints covered
- Admin unit tests: Authorization and role management covered
- Integration tests: Complete auth flow covered

### Target Coverage
- **>80% overall coverage** for authentication module
- **100% critical path coverage**: Auth verification, role management
- **Edge case coverage**: Invalid tokens, unauthorized access, self-demotion

## Future Enhancements

1. **Performance Tests**:
   - Load testing with concurrent authentication
   - Token verification performance benchmarks

2. **Security Tests**:
   - Token replay attack prevention
   - Session hijacking prevention
   - Rate limiting validation

3. **Integration Tests**:
   - Real Firebase integration (with test project)
   - Database persistence tests
   - Middleware integration tests

## Dependencies

### Testing Libraries
- `github.com/stretchr/testify/assert`: Assertions
- `github.com/stretchr/testify/mock`: Mocking
- `github.com/stretchr/testify/suite`: Test suites
- `github.com/stretchr/testify/require`: Required assertions

### Application Dependencies
- `go-pro-backend/internal/domain`: Data models
- `go-pro-backend/internal/service`: Business logic
- `go-pro-backend/internal/middleware`: HTTP middleware
- `go-pro-backend/internal/errors`: Error handling
- `go-pro-backend/pkg/logger`: Logging
- `go-pro-backend/pkg/validator`: Validation

## Notes

- Tests use in-memory implementations for speed
- Mock Firebase service simulates real Firebase behavior
- All tests are isolated and can run independently
- Context is properly propagated through all layers
- Error cases are comprehensively tested
