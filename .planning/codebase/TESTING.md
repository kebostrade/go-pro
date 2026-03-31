# Testing Strategy

**Analysis Date:** 2026-03-31

## Test Framework & Tools

### Backend (Go)

**Runner:**
- Go standard `testing` package
- Config: `backend/.golangci.yml` (linters enforce test quality)

**Assertion Library:**
- `github.com/stretchr/testify` v1.11.1
  - `assert` - non-fatal assertions (test continues on failure)
  - `require` - fatal assertions (test stops on failure)
  - `mock` - mock object generation via `mock.Mock` embedded struct

**Run Commands:**
```bash
# From repository root
make test                     # Unit tests with race detection
make test-coverage            # Tests with coverage report (HTML)
make test-integration         # Integration tests (build tag required)
make quality                  # All checks: lint + vet + security + test

# From backend directory
cd backend && make test       # All unit tests
cd backend && make test-unit  # Unit tests only (skips integration)
cd backend && make test-integration   # Integration tests only
cd backend && make test-coverage      # Coverage with HTML report
cd backend && make test-race          # Race detector only
cd backend && make test-bench         # Benchmark tests
cd backend && make ci                 # Full CI pipeline (lint + test with coverage)

# Direct go test commands
cd backend && go test -v -race ./...
cd backend && go test -v -race -short ./...
cd backend && go test -v -tags=integration ./...
cd backend && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
cd backend && go test -bench=. -benchmem ./...
```

### AI Agent Platform

**Runner:**
- Go standard `testing` package (no testify in this module - uses `t.Errorf`/`t.Fatalf` directly)

**Run Commands:**
```bash
cd services/ai-agent-platform && make test
cd services/ai-agent-platform && make test-coverage
cd services/ai-agent-platform && make bench
```

### Frontend (TypeScript/React)

**Runner:**
- Jest (used in test files via `jest.fn()` and `describe`/`it` blocks)
- No dedicated test config file found (no `jest.config.*` or `vitest.config.*`)
- One test file exists: `frontend/src/lib/__tests__/api.test.ts`

**Run Commands:**
```bash
cd frontend && bun test      # If configured
```

**Note:** Frontend testing is minimal. No test runner is configured in `package.json` scripts. The existing test file uses Jest globals but there is no explicit Jest or Vitest dependency.

## Test Organization

### Backend Test File Locations

Tests are **co-located** alongside the code they test:

```
backend/
├── internal/
│   ├── auth/
│   │   └── middleware_test.go         # Auth middleware tests
│   ├── handler/
│   │   ├── handler_test.go            # Main handler tests (health, courses, exercises, progress)
│   │   ├── auth_test.go               # Auth handler tests (verify, profile)
│   │   ├── admin_test.go              # Admin handler tests (role management)
│   │   └── cms_integration_test.go    # CMS integration tests (build tag: integration)
│   ├── middleware/
│   │   └── auth_test.go               # Middleware auth tests (external test package)
│   ├── repository/
│   │   ├── memory_simple_test.go      # In-memory repository tests
│   │   ├── memory_progress_test.go    # In-memory progress repository tests
│   │   └── postgres/
│   │       ├── course_test.go         # PostgreSQL course repo tests
│   │       ├── streak_test.go         # PostgreSQL streak repo tests
│   │       ├── assessment_cms_test.go # CMS assessment tests
│   │       ├── content_version_test.go # Content versioning tests
│   │       └── cms_test_helpers_test.go # CMS test helper utilities
│   ├── service/
│   │   ├── course_test.go             # Course service tests
│   │   ├── exercise_test.go           # Exercise service tests
│   │   └── exercise_evaluator_test.go # Exercise evaluation tests
│   └── executor/
│       ├── docker_executor_test.go    # Docker executor tests
│       └── example_test.go            # Example tests
└── internal/testutil/
    ├── testutil.go                    # Test utilities (TestDB, TestLogger, factory functions)
    └── mocks.go                       # Mock implementations (MockCacheManager, MockCourseRepository)
```

### AI Agent Platform Test Locations

```
services/ai-agent-platform/
└── internal/
    ├── embeddings/
    │   └── openai_test.go             # OpenAI embedding tests
    ├── rag/
    │   └── pipeline_test.go           # RAG pipeline tests
    └── vectorstore/
        └── memory_test.go             # In-memory vector store tests
```

### Frontend Test Locations

```
frontend/
└── src/
    └── lib/
        └── __tests__/
            └── api.test.ts            # API client tests (only frontend test)
```

### Naming Pattern

- Unit tests: `<name>_test.go` in the same package
- Integration tests: `<name>_integration_test.go` with `// +build integration` tag
- Test package: mostly same package (`package handler`), some use external test package (`package middleware_test`)

## Test Structure

### Backend Test Suite Pattern

Tests use the **table-driven test** pattern with `t.Run` subtests:

```go
// From backend/internal/handler/handler_test.go
func TestHandleGetCourse(t *testing.T) {
    tests := []struct {
        name           string
        courseID       string
        mockCourse     *domain.Course
        mockError      error
        expectedStatus int
    }{
        {
            name:     "success",
            courseID: "go-basics",
            mockCourse: &domain.Course{...},
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "not found",
            courseID:       "nonexistent",
            mockCourse:     nil,
            mockError:      errors.NewNotFoundError("Exercise not found"),
            expectedStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler, courseService, _, _, _, _, _ := setupTest()
            courseService.On("GetCourseByID", mock.Anything, tt.courseID).Return(tt.mockCourse, tt.mockError)

            mux := http.NewServeMux()
            mux.HandleFunc("GET /api/v1/courses/{id}", handler.handleGetCourse)
            req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/"+tt.courseID, nil)
            req.SetPathValue("id", tt.courseID)
            w := httptest.NewRecorder()
            mux.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
            // ... additional assertions
            courseService.AssertExpectations(t)
        })
    }
}
```

### Arrange-Act-Assert Pattern

Service tests use explicit AAA comments:

```go
// From backend/internal/service/course_test.go
func TestCourseService_Create(t *testing.T) {
    tests := []struct { ... }{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange.
            mockRepo := testutil.NewMockCourseRepository()
            mockCache := testutil.NewMockCacheManager()
            logger := testutil.NewTestLogger(t)
            service := NewCourseService(mockRepo, config)
            ctx := context.Background()

            // Act.
            _, err := service.CreateCourse(ctx, req)

            // Assert.
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Setup/Teardown Patterns

**Handler test setup** (returns handler + mocks):
```go
// From backend/internal/handler/auth_test.go
func setupAuthTestHandler(_ *testing.T) (*Handler, *MockAuthService) {
    mockAuthService := new(MockAuthService)
    services := &service.Services{Auth: mockAuthService}
    log := logger.New("debug", "json")
    val := validator.New()
    handler := New(services, log, val)
    return handler, mockAuthService
}
```

**Integration test setup** (connects to real database, cleans up on defer):
```go
// From backend/internal/handler/cms_integration_test.go
func NewCMSIntegrationTestHelper(t *testing.T) *CMSIntegrationTestHelper {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "..."))
    require.NoError(t, err)
    // ... setup router, helpers
    return &CMSIntegrationTestHelper{t: t, db: db, ...}
}
// Usage:
helper := NewCMSIntegrationTestHelper(t)
defer helper.Cleanup()
```

**Transaction-based cleanup:**
```go
// From backend/internal/testutil/testutil.go
func WithTransaction(t *testing.T, db *sql.DB, fn func(*sql.Tx) error) {
    t.Helper()
    tx, err := db.Begin()
    require.NoError(t, err)
    defer func() {
        rbErr := tx.Rollback() // Always rollback test transactions
    }()
    err = fn(tx)
    require.NoError(t, err)
}
```

## Mocking

### Backend Mocking Framework

**Library:** `github.com/stretchr/testify/mock`

**Pattern:** Embed `mock.Mock` in a struct, implement interface methods with `Called()`:

```go
// From backend/internal/handler/handler_test.go
type mockCourseService struct {
    mock.Mock
}

func (m *mockCourseService) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Course), args.Error(1)
}
```

**Setting up mocks:**
```go
courseService.On("GetCourseByID", mock.Anything, "go-basics").
    Return(&domain.Course{ID: "go-basics", Title: "Go Basics"}, nil)
```

**Verifying mocks:**
```go
courseService.AssertExpectations(t)
mockRepo.GetCallCount("Create")  // Check call counts on custom mocks
```

### Custom Mock Implementations (testutil package)

`backend/internal/testutil/mocks.go` provides hand-rolled mocks with `sync.RWMutex`:

```go
// MockCourseRepository - thread-safe, tracks call counts
type MockCourseRepository struct {
    mu      sync.RWMutex
    courses map[string]*domain.Course
    calls   map[string]int
}

// MockCacheManager - full cache interface implementation
type MockCacheManager struct {
    mu    sync.RWMutex
    data  map[string]interface{}
    calls map[string]int
}
```

These mocks are used in service-layer tests where testify mock generation would be overkill or where state tracking is needed.

### Function-Based Mocks (Middleware Tests)

Middleware tests use function fields for flexible mocking:

```go
// From backend/internal/middleware/auth_test.go
type MockAuthService struct {
    verifyTokenFunc func(ctx context.Context, token string) (*middleware.FirebaseToken, error)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
    if m.verifyTokenFunc != nil {
        return m.verifyTokenFunc(ctx, token)
    }
    return &middleware.FirebaseToken{UID: "test-uid"}, nil
}
```

### What to Mock

- **Service interfaces** in handler tests (via testify mock)
- **Repository interfaces** in service tests (via testutil hand-rolled mocks)
- **External services** (Firebase, OpenAI, Docker) in all tests
- **Cache** via `MockCacheManager`
- **Logger** via `testutil.NewTestLogger(t)` (silent unless warnings/errors)

### What NOT to Mock

- **Domain entities** - use real structs and factory functions from `testutil`
- **Validators** - use real `validator.New()`
- **In-memory repositories** - use real `NewMemoryCourseRepository()` etc. for simple tests

### AI Agent Platform Mocking

The AI agent platform uses mock implementations for external dependencies:

```go
// From services/ai-agent-platform/internal/rag/pipeline_test.go
func createTestPipeline(t *testing.T) *RAGPipeline {
    embedder := embeddings.NewMockEmbedder(1536)
    vectorStore := vectorstore.NewMemoryVectorStore()
    config := types.RAGConfig{
        VectorStore: vectorStore,
        Embedder:    embedder,
    }
    pipeline, err := NewRAGPipeline(config)
    if err != nil {
        t.Fatalf("Failed to create test pipeline: %v", err)
    }
    return pipeline
}
```

Uses `t.Fatalf`/`t.Errorf` directly (no testify).

### Frontend Mocking

```typescript
// From frontend/src/lib/__tests__/api.test.ts
const mockFetch = jest.fn();
global.fetch = mockFetch;

// Mock responses
mockFetch.mockResolvedValueOnce({
    ok: true,
    json: async () => ({ success: true, data: mockHealth }),
});
```

## Fixtures and Factories

### Test Data Factory Functions

Located in `backend/internal/testutil/testutil.go`:

```go
CreateTestCourse(id, title string) *domain.Course
CreateTestLesson(id, courseID, title string, order int) *domain.Lesson
CreateTestExercise(id, lessonID, title string) *domain.Exercise
CreateTestProgress(id, userID, lessonID string, status domain.Status) *domain.Progress
CreateTestUser(id, username, email string) *domain.User
```

### Test Utilities

```go
// Test database
NewTestDB(t *testing.T) *TestDB                    // Connect to test PostgreSQL
tdb.TruncateTables(ctx, "courses", "lessons")       // Clean tables
WithTransaction(t, db, func(tx *sql.Tx) error)      // Auto-rollback transactions

// Helpers
RandomString(length int) string                     // Random test string
RandomEmail() string                                // Random test email
WaitForCondition(t, timeout, condition)              // Polling helper
GetEnv(key, defaultValue string) string              // Env with fallback
```

### Integration Test Helpers

`backend/internal/repository/postgres/cms_test_helpers_test.go` provides:
- `CMSTestHelper` for creating test users, courses, lessons, content versions
- Table truncation for clean test state
- Integration test DB setup from `TEST_DATABASE_URL`

## Coverage

**Requirements:** No enforced minimum, but coverage reports are generated.

**View Coverage:**
```bash
# Backend
cd backend && make test-coverage        # Generates coverage/coverage.html
cd backend && make test-coverage-view   # Opens HTML report
make test-coverage                      # From root (generates backend/coverage.html)

# AI Agent Platform
cd services/ai-agent-platform && make test-coverage
```

**Coverage Configuration:**
- Cover mode: `atomic` (for race-safe coverage)
- Output: `coverage/coverage.out` + `coverage/coverage.html`
- CI: coverage generated during `ci-test` target

## Test Types

### Unit Tests

**Scope:** Individual functions, methods, handlers, services, repositories.

**Location:** Co-located `*_test.go` files.

**Dependencies:** All external dependencies mocked. Tests run without database or network.

**Pattern:** Table-driven tests with testify assertions. Each test case is a struct in a slice, iterated with `t.Run`.

**Handler unit tests** use `httptest.NewRecorder()` and `httptest.NewRequest()`:
```go
req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/"+tt.courseID, nil)
req.SetPathValue("id", tt.courseID)
w := httptest.NewRecorder()
mux.ServeHTTP(w, req)
```

**Service unit tests** use `testutil` mocks and `context.Background()`:
```go
mockRepo := testutil.NewMockCourseRepository()
service := NewCourseService(mockRepo, config)
result, err := service.GetCourseByID(context.Background(), "course-1")
```

**Repository unit tests** use in-memory implementations:
```go
repo := NewMemoryCourseRepository()
err := repo.Create(ctx, course)
retrieved, err := repo.GetByID(ctx, "test-course-1")
```

### Integration Tests

**Scope:** Full request-response cycles against real database (PostgreSQL).

**Location:** Files with `_integration_test.go` suffix.

**Build Tag:** `// +build integration` at top of file.

**Skip Condition:** `testing.Short()` check skips integration tests during `go test -short`.

**Run Command:**
```bash
cd backend && go test -v -tags=integration ./...
# or
cd backend && make test-integration
```

**Required Infrastructure:**
- PostgreSQL test database (configured via `TEST_DATABASE_URL`)
- Test credentials: `gopro_test:gopro_test@localhost:5432/gopro_test`

**Example:** `backend/internal/handler/cms_integration_test.go` - Tests CMS CRUD, versioning, publishing, rollback against real PostgreSQL.

### End-to-End Tests

Not currently implemented. The test infrastructure exists for load testing via k6:

```bash
make test-load    # Runs k6 load test from backend/scripts/load-test.js
```

### Benchmark Tests

Benchmarks follow Go standard pattern with `b *testing.B`:

```go
// From backend/internal/service/course_test.go
func BenchmarkCourseService_Create(b *testing.B) {
    service := NewCourseService(mockRepo, config)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.CreateCourse(ctx, req)
    }
}
```

**Run:**
```bash
cd backend && go test -bench=. -benchmem ./...
cd backend && make test-bench
cd services/ai-agent-platform && make bench
```

## CI Testing

### CI Pipeline Commands

```bash
# Full CI pipeline
make ci-test                          # deps + lint + vet + security + test-coverage
make ci-build                         # ci-test + build + docker-build

# From backend
cd backend && make ci                 # ci-lint + ci-test
cd backend && make ci-lint            # golangci-lint run --timeout=5m
cd backend && make ci-test            # tests with coverage profiling
```

### Quality Gates

From root `Makefile`:
```bash
make quality    # deps + lint + vet + security + test
```

Individual quality checks:
- **Linting:** `golangci-lint run --timeout=5m` (backend only, extensive linter config)
- **Vetting:** `go vet ./...`
- **Security:** `gosec ./...`
- **Vulnerabilities:** `govulncheck ./...`
- **Formatting:** `gofmt -s -w .` + `goimports -w .`
- **Race detection:** `-race` flag on all test commands

### Linter Exemptions for Tests

From `backend/.golangci.yml`, these linters are disabled for `*_test.go` files:
- `gocyclo`, `errcheck`, `dupl`, `gosec`, `funlen`, `goconst`, `gocognit`, `cyclop`, `thelper`, `nlreturn`, `errorlint`

This allows more flexibility in test code (longer functions, duplicated test setup, relaxed error checking).

## Common Patterns

### HTTP Handler Testing

```go
func TestHandleHealth(t *testing.T) {
    handler, _, _, _, _, _, healthService := setupTest()
    healthService.On("GetHealthStatus", mock.Anything).Return(&domain.HealthResponse{...}, nil)

    req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
    w := httptest.NewRecorder()
    handler.handleHealth(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response domain.APIResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success)

    healthService.AssertExpectations(t)
}
```

### Error Path Testing

```go
// Test various error types from the same endpoint
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name           string
        setupMock      func(*mockCourseService)
        expectedStatus int
        expectedType   string
    }{
        {
            name: "not found error",
            setupMock: func(m *mockCourseService) {
                m.On("GetCourseByID", mock.Anything, "nonexistent").
                    Return(nil, errors.NewNotFoundError("not found"))
            },
            expectedStatus: http.StatusNotFound,
            expectedType:   "NOT_FOUND",
        },
        // ... more error cases
    }
}
```

### Async Testing

```go
// From backend/internal/testutil/testutil.go
func WaitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
    t.Helper()
    deadline := time.Now().Add(timeout)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    for {
        if condition() { return }
        <-ticker.C
        if time.Now().After(deadline) { t.Fatal("Timeout waiting for condition") }
    }
}
```

### Skipping Tests Requiring Infrastructure

```go
// Integration test skip pattern
if testing.Short() {
    t.Skip("Skipping integration test")
}

// Infrastructure-dependent test skip
func TestExerciseServiceWithMessaging(t *testing.T) {
    t.Skip("Skipping messaging integration test - requires Kafka infrastructure")
}
```

---

*Testing analysis: 2026-03-31*
