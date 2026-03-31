# Coding Conventions

**Analysis Date:** 2026-03-31

## Naming Conventions

### Go Files (Backend)

**File naming:**
- Use `snake_case` for all Go files: `course_test.go`, `docker_executor.go`
- Test files: append `_test.go` suffix to the file being tested: `handler.go` -> `handler_test.go`
- Integration tests: append `_integration_test.go`: `cms_integration_test.go`
- Mock files: use `mock.go` or `mocks.go` suffix: `backend/internal/testutil/mocks.go`
- Package names: single lowercase word, no underscores: `handler`, `service`, `repository`, `domain`, `middleware`

**Functions and methods:**
- Use `PascalCase` for exported functions: `NewCourseService`, `GetCourseByID`
- Use `camelCase` for unexported functions: `setupTest`, `countPassedTests`, `writeSuccessResponse`
- Handler methods follow `handle` prefix pattern: `handleHealth`, `handleCreateCourse`, `handleGetCourse`
- Constructor functions use `New` prefix: `New()`, `NewCourseService()`, `NewTestDB()`
- Test helper functions use `setup` or `create` prefix: `setupTest()`, `setupAuthTestHandler()`, `createTestPipeline()`
- Test functions follow `Test` + `MethodName` + `Scenario` pattern: `TestAuthVerify_Success`, `TestCourseService_Create`

**Variables:**
- Use `camelCase` for local variables: `mockRepo`, `statusCode`, `courseID`
- Use `PascalCase` for exported constants: `StatusCompleted`, `RoleAdmin`
- Use `camelCase` for unexported constants: `contentTypeJSON`, `adminUserID`
- Acronyms stay uppercase: `userID`, `courseID`, `apiErr`, not `userId`, `courseId`
- Error sentinel variables use `Err` prefix: `ErrNotFound`, `ErrValidation`

**Types:**
- Use `PascalCase` for all types (exported and unexported): `Handler`, `MockAuthService`, `rateLimitState`
- Interfaces: `CourseService`, `ExerciseRepository`, `AuthService`
- Struct implementations: `courseService` (unexported), `exerciseService` (unexported)
- Request/Response types: `CreateCourseRequest`, `ExerciseSubmissionResult`, `APIResponse`
- Error types: `APIError` with factory constructors: `NewNotFoundError()`, `NewBadRequestError()`

### TypeScript Files (Frontend)

**File naming:**
- Use `camelCase` for utilities: `api.ts`, `firebase.ts`
- Use `PascalCase` for components: directories like `components/`, `contexts/`
- Test files: `__tests__/api.test.ts` in co-located `__tests__` directory
- Path alias: `@/*` maps to `./src/*` (configured in `frontend/tsconfig.json`)

**Variables and functions:**
- Use `camelCase` for variables and functions: `mockFetch`, `api.health()`
- Use `PascalCase` for types and interfaces: `APIResponse`, `APIError`, `BackendUser`
- Test descriptions: use `describe`/`it` blocks with descriptive strings: `it('should fetch health status successfully')`

## Code Style

### Go Formatting (Backend)

**Formatter:** `gofmt` with `gofumpt` (stricter rules) configured in `backend/.golangci.yml`
- Line length limit: 150 characters (configured via `lll` linter)
- Tab width: 4 spaces
- Use `goimports` for import management
- Local prefix for imports: `go-pro-backend` (configured in `goimports` linter settings)

**Linting:** `golangci-lint` with extensive configuration at `backend/.golangci.yml`
- 40+ linters enabled including: `errcheck`, `staticcheck`, `gosec`, `revive`, `gocritic`, `funlen`, `cyclop`
- Cyclomatic complexity max: 15 (package average: 10)
- Function length max: 100 lines / 60 statements
- Cognitive complexity max: 15
- Test files exempt from: `gocyclo`, `errcheck`, `dupl`, `gosec`, `funlen`, `goconst`, `gocognit`, `cyclop`

**License header (required):**
```go
// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License
```
Enforced by `goheader` linter in `backend/.golangci.yml`.

**Package comments:**
```go
// Package handler provides HTTP request handlers for the API.
package handler
```

### TypeScript Formatting (Frontend)

**Formatter:** Next.js defaults (no separate Prettier config detected)
**Linting:** `next lint` (Next.js built-in ESLint)
**TypeScript:** Strict mode enabled in `frontend/tsconfig.json`
**Package manager:** Bun

## Go-Specific Conventions

### Error Handling

Use custom `APIError` type from `backend/internal/errors/errors.go`:

```go
// Typed error constructors
apierrors.NewNotFoundError("course not found")
apierrors.NewBadRequestError("invalid ID format")
apierrors.NewValidationError("validation failed", err)
apierrors.NewInternalError("database error", err)
apierrors.NewUnauthorizedError("invalid token")
apierrors.NewForbiddenError("admin access required")
apierrors.NewConflictError("resource already exists")
```

**Error checking in handlers:**
```go
if err != nil {
    h.writeErrorResponse(w, r, err)
    return
}
```

**Error type assertion:**
```go
var apiErr *apierrors.APIError
if errors.As(err, &apiErr) {
    statusCode = apiErr.StatusCode
}
```

### Interface Definitions

Repository interfaces in `backend/internal/repository/interfaces.go`:
```go
type CourseRepository interface {
    Create(ctx context.Context, course *domain.Course) error
    GetByID(ctx context.Context, id string) (*domain.Course, error)
    GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error)
    Update(ctx context.Context, course *domain.Course) error
    Delete(ctx context.Context, id string) error
}
```

Service interfaces in `backend/internal/service/interfaces.go`:
```go
type CourseService interface {
    CreateCourse(ctx context.Context, req *domain.CreateCourseRequest) (*domain.Course, error)
    GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
    GetAllCourses(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error)
    ...
}
```

**Pattern:** Always accept `context.Context` as first parameter. Return `(result, error)` tuple.

### Context Usage

- Always pass `context.Context` as first parameter to all service/repository methods
- Use `r.Context()` in handlers to propagate request context
- Use `context.Background()` in tests and background operations
- Store user info in context via `middleware.WithUser(ctx, user)`

### Dependency Injection

Services are aggregated in a `Services` struct:
```go
// backend/internal/service/interfaces.go
type Services struct {
    Course     CourseService
    Lesson     LessonService
    Exercise   ExerciseService
    Progress   ProgressService
    Curriculum CurriculumService
    Health     HealthService
    Auth       AuthService
}
```

Handler receives `Services` through constructor:
```go
handler := New(services, log, val)
```

### Validation

Use struct tags for validation rules (via `backend/pkg/validator/`):
```go
type Course struct {
    ID          string `json:"id" validate:"required,slug"`
    Title       string `json:"title" validate:"required,min=3,max=200"`
    Description string `json:"description" validate:"required,min=10,max=1000"`
}
```

Validate JSON in handlers:
```go
if err := h.validator.ValidateJSON(r, &req); err != nil {
    h.writeErrorResponse(w, r, err)
    return
}
```

### Logging

Use the `logger.Logger` interface from `backend/pkg/logger/`:
```go
// Production logging
logger.New("debug", "json") // level, format
logger.New("info", "text")

// Test logging (silent)
testutil.NewTestLogger(t)
```

## TypeScript/Frontend Conventions

### Component Patterns

- Framework: Next.js 15 with App Router (`frontend/src/app/`)
- UI Library: Radix UI primitives + custom components
- Styling: Tailwind CSS with `class-variance-authority` + `clsx`
- State: React Context (`frontend/src/contexts/`)
- Editor: Monaco Editor for code editing (`@monaco-editor/react`)
- Rich Text: TipTap editor (`@tiptap/react`)

### API Client Usage

Centralized API client at `frontend/src/lib/api.ts`:
```typescript
import { api } from '@/lib/api';

// All API calls go through the centralized client
const curriculum = await api.getCurriculum();
const lesson = await api.getLessonDetail(1);
const result = await api.submitExercise('exercise-001', code, authToken);
```

### Error Handling Pattern

Custom `APIError` class wraps all API errors:
```typescript
class APIError extends Error {
  status: number;
  type: string;
  details?: Record<string, string>;
}
```

## API Design Conventions

### REST Endpoint Patterns

All endpoints under `/api/v1/` prefix. Route registration in `backend/internal/handler/handler.go`:

| Resource | GET (list) | GET (single) | POST (create) | PUT (update) | DELETE |
|----------|-----------|-------------|---------------|-------------|--------|
| Courses | `/api/v1/courses` | `/api/v1/courses/{id}` | `/api/v1/courses` | `/api/v1/courses/{id}` | `/api/v1/courses/{id}` |
| Lessons | `/api/v1/lessons` | `/api/v1/lessons/{id}` | `/api/v1/lessons` | `/api/v1/lessons/{id}` | `/api/v1/lessons/{id}` |
| Exercises | - | `/api/v1/exercises/{id}` | - | - | - |
| Progress | `/api/v1/users/{userId}/progress` | - | - | - | - |
| Curriculum | `/api/v1/curriculum` | `/api/v1/curriculum/lesson/{id}` | - | - | - |

**Special actions:**
- `POST /api/v1/exercises/{id}/submit` - Submit exercise solution
- `POST /api/v1/playground/execute` - Execute code in playground
- `POST /api/v1/progress/{userId}/lesson/{lessonId}` - Update progress

### Request/Response Format

**Standard response envelope** (`backend/internal/domain/models.go`):
```json
{
  "success": true,
  "data": {},
  "message": "operation successful",
  "request_id": "req-123",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

**Error response:**
```json
{
  "success": false,
  "error": {
    "type": "NOT_FOUND",
    "message": "Course not found",
    "details": {}
  },
  "request_id": "req-456",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

**Error types:** `NOT_FOUND`, `BAD_REQUEST`, `VALIDATION_ERROR`, `INTERNAL_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`, `CONFLICT`, `RATE_LIMIT_EXCEEDED`

### Pagination

Query parameters: `page` and `page_size` (default 10, max 100):
```
GET /api/v1/courses?page=1&page_size=10
```

Response includes pagination metadata:
```json
{
  "items": [],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 42,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

### Status Codes

- `200 OK` - Success
- `400 Bad Request` - Validation error, bad input
- `401 Unauthorized` - Missing or invalid auth token
- `403 Forbidden` - Insufficient permissions (e.g., student accessing admin)
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded (exercise submissions)
- `500 Internal Server Error` - Unexpected errors

## Configuration Conventions

### Environment Variables (Backend)

- `.env` file present in `backend/` (gitignored)
- `.env.test` file for test environment
- `.env.example` in `services/ai-agent-platform/`
- Load with `godotenv` (`github.com/joho/godotenv`)
- Test DB connection: `TEST_DATABASE_URL` or individual `TEST_DB_HOST`, `TEST_DB_PORT`, `TEST_DB_USER`, `TEST_DB_PASSWORD`, `TEST_DB_NAME`

### Environment Variables (Frontend)

- `NEXT_PUBLIC_API_URL` - Backend API base URL (required for backend integration)
- Firebase config for authentication
- Configuration loaded from `.env.local`

### Build Configuration

- Backend module: `go-pro-backend` (Go 1.23+)
- AI Agent module: `github.com/DimaJoyti/go-pro/services/ai-agent-platform`
- Hot reload: Air (`backend/.air.toml`)
- Frontend build: `next build` with Turbopack (`next dev --turbopack`)

## Git & Commit Conventions

### Commit Message Format

Conventional Commits style (from `git log`):

```
<type>(<scope>): <description>
```

**Types used:** `feat`, `fix`, `chore`, `refactor`, `docs`, `style`, `test`, `perf`, `ci`, `build`

**Examples:**
```
feat: Implement in-memory repositories for Progress and Exercise
feat(playground): add playground code execution and AI features
fix: restore Playground page and PlaygroundEditor component
refactor: streamline footer and header components
chore: add lesson 11 (Advanced Concurrency) code
feat(docs): Add MCP Server Troubleshooting Guide
feat(graphql): add GraphQL schema and models for blog API
```

### Branch Naming

- `main` - Production branch
- `dev` - Development branch (current)
- `merge-<feature>` - Feature merge branches: `merge-customer-management`
- `dependabot/<type>/<path>` - Automated dependency updates
- Feature branches appear to use descriptive names without strict prefix conventions

### PR Conventions

- PRs merge to `main` from feature branches
- Example: `Merge pull request #73 from DimaJoyti:merge-customer-management`

---

*Convention analysis: 2026-03-31*
