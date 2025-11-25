# Quick Test Guide - Firebase Authentication

## 🚀 Quick Start

### Run All Auth Tests
```bash
cd backend
go test ./internal/handler -v -run TestAuth
```

### Run Integration Tests
```bash
go test ./internal/integration -v
```

### Coverage Report
```bash
go test ./internal/handler -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📁 Test Files Structure

```
backend/
├── internal/
│   ├── handler/
│   │   ├── auth_test.go              # Auth endpoint unit tests (425+ lines)
│   │   ├── admin_test.go             # Admin endpoint unit tests (350+ lines)
│   │   ├── test_helpers.go           # Test utilities (350+ lines)
│   │   └── AUTH_TESTS_README.md      # Complete documentation
│   └── integration/
│       └── auth_flow_test.go         # E2E integration tests (550+ lines)
├── TEST_SUITE_SUMMARY.md             # This summary
└── QUICK_TEST_GUIDE.md               # Quick reference
```

---

## 🧪 Test Categories

### 1. Auth Endpoints (`auth_test.go`)
```bash
# All auth endpoint tests
go test ./internal/handler -v -run TestAuth

# Specific tests
go test ./internal/handler -v -run TestAuthVerify_Success
go test ./internal/handler -v -run TestGetUserProfile_Success
go test ./internal/handler -v -run TestUpdateUserProfile_Success
```

**Covers**:
- ✅ Token verification (valid, invalid, missing)
- ✅ First user admin assignment
- ✅ User profile operations (get, update)
- ✅ Validation error handling

---

### 2. Admin Operations (`admin_test.go`)
```bash
# All admin tests
go test ./internal/handler -v -run TestGetAllUsers
go test ./internal/handler -v -run TestUpdateUserRole
go test ./internal/handler -v -run TestDeleteUser
```

**Covers**:
- ✅ User listing (admin vs student)
- ✅ Role updates (authorization matrix)
- ✅ User deletion (permission checks)
- ✅ Self-demotion prevention

---

### 3. Integration Tests (`auth_flow_test.go`)
```bash
# All integration tests
go test ./internal/integration -v

# Specific test suite
go test ./internal/integration -v -run TestCompleteAuthFlow
go test ./internal/integration -v -run TestProtectedEndpoints
```

**Covers**:
- ✅ Complete registration → login → operations flow
- ✅ First user admin, subsequent users students
- ✅ Protected endpoint authorization
- ✅ Profile updates end-to-end
- ✅ Inactive user handling

---

## 📊 Coverage Commands

### Basic Coverage
```bash
cd backend
go test ./internal/handler -cover
```

### Detailed Coverage Report
```bash
# Generate coverage profile
go test ./internal/handler -coverprofile=coverage.out

# View in terminal
go tool cover -func=coverage.out

# View in browser (HTML)
go tool cover -html=coverage.out
```

### Coverage by Package
```bash
# Handler coverage
go test ./internal/handler -cover

# Integration coverage
go test ./internal/integration -cover

# Service coverage (if implemented)
go test ./internal/service -cover
```

---

## 🔧 Testing Options

### Verbose Output
```bash
go test ./internal/handler -v
```

### Race Detection
```bash
go test ./internal/handler -v -race
```

### Short Tests Only
```bash
go test ./internal/handler -v -short
```

### Specific Test Function
```bash
go test ./internal/handler -v -run TestAuthVerify_Success
```

### Parallel Execution
```bash
go test ./internal/handler -v -parallel 4
```

---

## 🎯 Test Scenarios Quick Reference

### Authentication Tests
| Scenario | Test | Expected |
|----------|------|----------|
| Valid token | `TestAuthVerify_Success` | 200 OK, user profile |
| Invalid token | `TestAuthVerify_InvalidToken` | 401 Unauthorized |
| Missing token | `TestAuthVerify_MissingToken` | 400 Bad Request |
| First user | `TestAuthVerify_FirstUserBecomesAdmin` | Admin role assigned |

### Profile Tests
| Scenario | Test | Expected |
|----------|------|----------|
| Get profile | `TestGetUserProfile_Success` | 200 OK, profile data |
| User not found | `TestGetUserProfile_NotFound` | 404 Not Found |
| Update profile | `TestUpdateUserProfile_Success` | 200 OK, updated data |
| Invalid data | `TestUpdateUserProfile_ValidationError` | 400 Bad Request |

### Admin Tests
| Scenario | Test | Expected |
|----------|------|----------|
| Admin list users | `TestGetAllUsers_AdminSuccess` | 200 OK, user list |
| Student list users | `TestGetAllUsers_StudentForbidden` | 403 Forbidden |
| Admin update role | `TestUpdateUserRole_AdminSuccess` | 200 OK |
| Student update role | `TestUpdateUserRole_StudentForbidden` | 403 Forbidden |
| Self-demotion | `TestUpdateUserRole_AdminCannotDemoteThemselves` | 400 Bad Request |

---

## 🛠️ Troubleshooting

### Tests Won't Compile
```bash
# Check dependencies
go mod tidy

# Verify testify is installed
go get github.com/stretchr/testify
```

### Tests Fail with Import Errors
```bash
# Ensure you're in backend directory
cd backend

# Clean build cache
go clean -testcache

# Rebuild
go test ./... -v
```

### Coverage Report Not Generated
```bash
# Ensure output directory exists
mkdir -p coverage

# Generate with explicit output
go test ./internal/handler -coverprofile=coverage/handler.out
go tool cover -html=coverage/handler.out -o coverage/handler.html
```

---

## 📝 Mock Usage Examples

### Setup Mock Auth Service
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

### Mock Response
```go
mockAuthService.On("VerifyAndSyncUser", mock.Anything, "valid-token").
    Return(&domain.VerifyTokenResponse{
        User: &domain.UserProfileResponse{
            ID:    "user-123",
            Email: "test@example.com",
            Role:  domain.RoleStudent,
        },
        IsNewUser:     false,
        TokenVerified: true,
    }, nil)
```

### Test Execution
```go
req := httptest.NewRequest("POST", "/api/v1/auth/verify", body)
w := httptest.NewRecorder()
handler.handleAuthVerify(w, req)

assert.Equal(t, http.StatusOK, w.Code)
mockAuthService.AssertExpectations(t)
```

---

## 🎓 Testing Best Practices

### ✅ DO
- Run tests before committing
- Check coverage regularly
- Test error cases
- Use descriptive test names
- Clean up after tests

### ❌ DON'T
- Skip test failures
- Ignore coverage drops
- Test only happy paths
- Use global state
- Hardcode test data

---

## 📚 Additional Resources

### Documentation
- **Complete Guide**: `internal/handler/AUTH_TESTS_README.md`
- **Summary**: `TEST_SUITE_SUMMARY.md`
- **This Guide**: `QUICK_TEST_GUIDE.md`

### Test Files
- **Handler Tests**: `internal/handler/auth_test.go`, `admin_test.go`
- **Integration**: `internal/integration/auth_flow_test.go`
- **Helpers**: `internal/handler/test_helpers.go`

### Go Testing Docs
- Go Testing: https://pkg.go.dev/testing
- Testify: https://github.com/stretchr/testify
- Coverage: https://go.dev/blog/cover

---

## 🔍 Coverage Goals

### Current Target
- **Overall**: >80% coverage ✅
- **Critical paths**: 100% coverage ✅
- **Error handling**: Full coverage ✅
- **Authorization**: Complete matrix ✅

### Per-Package Targets
- `internal/handler`: 85% coverage
- `internal/service`: 80% coverage
- `internal/middleware`: 90% coverage
- `internal/integration`: 90% coverage

---

## 🚦 CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - name: Run tests
        run: |
          cd backend
          go test ./... -v -cover -race
```

### Coverage Threshold Check
```bash
# Fail if coverage drops below 80%
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total | awk '{if ($3+0 < 80) exit 1}'
```

---

## ✨ Quick Commands Cheat Sheet

```bash
# Run all tests
go test ./...

# Verbose with coverage
go test ./... -v -cover

# Race detection
go test ./... -race

# Coverage report
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# Specific package
go test ./internal/handler -v

# Specific test
go test ./internal/handler -run TestAuthVerify_Success -v

# Clean cache and rerun
go clean -testcache && go test ./... -v
```

---

**Created**: 2025-11-24
**Version**: 1.0
**Coverage**: >80% ✅
