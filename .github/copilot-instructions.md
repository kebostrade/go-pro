# Copilot Instructions for GO-PRO

GO-PRO is a comprehensive **full-stack Go learning platform** with both a **learning content system** and a **production AI agent framework**.

## Architecture Overview

### Repository Structure
- **Dual-module design**: Multiple independent Go modules (all Go 1.23)
  - `backend/` - Learning platform REST API (Clean Architecture)
  - `services/` - Production services (ai-agent-platform, api-gateway, shared)
  - `basic/`, `course/` - Learning content and examples
  - `frontend/` - Next.js dashboard (React + TypeScript)

### Backend Architecture (Clean Separation)
```
backend/cmd/server/           → Entry point with DI container
backend/internal/handler/      → HTTP handlers (controllers)
backend/internal/service/      → Business logic
backend/internal/repository/   → Data access (interfaces in interfaces.go)
  - memory_simple.go          → In-memory for dev/testing
  - postgres/                 → PostgreSQL for production
backend/pkg/                  → Public, reusable packages
```

**Key Pattern**: Repository interfaces (`CourseRepository`, `LessonRepository`, `ExerciseRepository`, `ProgressRepository`) allow runtime switching between in-memory and PostgreSQL implementations without code changes.

### Frontend Architecture (Next.js 15)
```
frontend/src/app/             → App Router pages
frontend/src/components/      → React components
frontend/src/contexts/        → React contexts (auth, etc.)
frontend/src/lib/api.ts       → Centralized API client
```

## Critical Development Workflows

### Module-First Operations
**Always run Go commands from the correct module directory** (not root):
```bash
cd backend && go test ./...         # Backend tests
cd services/ai-agent-platform && go test ./...  # AI platform
cd basic && go run examples/fun/binary_search.go  # Examples
```

### Backend Development Cycle
```bash
# Development with hot reload
cd backend && air -c .air.toml        # Uses Air for file watching

# Quick test + lint + security
make test lint security vuln-check

# Full quality gate (what CI runs)
make quality  # deps → lint → vet → security → test
```

### Frontend Development
```bash
cd frontend
npm run dev          # Starts at http://localhost:3000
```

### Full Stack Development (with Docker)
```bash
make docker-dev      # Starts backend, frontend, PostgreSQL, Redis, Grafana
# Services: Backend (8080), Frontend (3000), Adminer (8081), Prometheus (9090)
```

## Project-Specific Patterns

### Error Handling
- Go idiomatic: Return errors as values, check `if err != nil`
- No exceptions/panics in production code
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`

### Testing Strategy
- **Unit tests**: Live alongside code (`handler_test.go` next to `handler.go`)
- **Test flags**: Use `-race` for concurrency bugs: `go test -race ./...`
- **Integration tests**: Build tag `// +build integration` or `-tags=integration`
- **Coverage target**: Aim for >70% coverage on business logic

### Code Executor Component
Backend includes a **Docker-based code executor** (`internal/executor/docker_executor.go`):
- Executes untrusted student Go code safely
- Validates: `package main` and `func main()` required
- Returns stdout/stderr to API
- See `executor_test.go` for validation patterns

### Handler Pattern
Handlers receive dependency-injected services, use context for cancellation:
```go
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Use request context
    var req CreateCourseRequest
    // Decode → validate → call service → respond with errors
}
```

### Repository Pattern
All data access through interfaces for testability:
```go
// Declare interface
type CourseRepository interface {
    Create(ctx context.Context, course *domain.Course) error
    GetByID(ctx context.Context, id string) (*domain.Course, error)
}
// Inject implementation (memory or postgres)
```

## Build & Quality Commands

### Quick Development
```bash
make dev              # Hot reload backend
cd frontend && npm run dev   # Frontend dev
make docker-dev       # Full stack with services
```

### Testing
```bash
make test             # All tests with race detection
make test-coverage    # With HTML report (coverage.html)
make test-integration # Integration tests only (tags=integration)
```

### Code Quality
```bash
make fmt              # Format + organize imports (gofmt + goimports)
make lint             # golangci-lint with 5min timeout
make lint-fix         # Auto-fix linting issues
make vet              # go vet (deadcode, etc.)
make security         # gosec vulnerability scan
make vuln-check       # govulncheck for known CVEs
```

### CI/CD Gates
```bash
make quality          # Full quality gate: deps → lint → vet → security → test
make ci-build         # CI build: quality → build → docker-build
```

## Frontend-Backend Integration
- Backend runs on `:8080` (`internal/handler` registers routes)
- Frontend API client: `frontend/src/lib/api.ts`
- Frontend connects via `NEXT_PUBLIC_API_URL` in `.env.local`
- CORS enabled for localhost development

## Dependencies to Know
- **HTTP**: Standard `net/http` + `github.com/gorilla/mux` for routing
- **Database**: `github.com/lib/pq` (PostgreSQL driver)
- **Caching**: `github.com/go-redis/redis/v8`
- **Auth**: `github.com/golang-jwt/jwt/v5` for JWT tokens
- **Crypto**: `golang.org/x/crypto` for password hashing
- **Frontend**: Next.js 15 (App Router), React 19, Tailwind CSS, Radix UI

## Environment Configuration
- Backend: `backend/.env` (copy from `.env.example`)
  - Database credentials, JWT secrets, Redis URL
- Frontend: `frontend/.env.local` (copy from `.env.example`)
  - `NEXT_PUBLIC_API_URL` must match backend
- AI Agent Platform: `services/ai-agent-platform/.env`
  - OpenAI API key, model settings

## Common Pitfalls
1. **Wrong working directory**: Always `cd` to module before running `go` commands
2. **Stale tests**: Run `go clean -testcache` if tests behave unexpectedly
3. **Import issues**: Run `make fmt` to organize imports (goimports)
4. **Docker networking**: `make docker-dev` includes all services on custom network
5. **Go version**: All modules require Go 1.23+ (check with `go version`)

## Directory Reference
- **Backend API**: `backend/cmd/server/main.go` (DI container at `internal/container.go`)
- **Key Interfaces**: `backend/internal/repository/interfaces.go`
- **Domain Models**: `backend/internal/domain/`
- **Frontend API Client**: `frontend/src/lib/api.ts`
- **Makefile Commands**: Root `Makefile` for build automation
- **Docker Setup**: `docker-compose.dev.yml` (dev) and `docker-compose.prod.yml` (prod)

## When to Ask for Context
This codebase has:
- 4+ independent Go modules with different purposes
- Production AI agent framework alongside learning platform
- Custom code executor component
- Full-stack integration (Go backend + Next.js frontend)
- Docker-based development environment

If working in a specific module or feature area, request clarification on dependencies and existing patterns in that subsystem.
