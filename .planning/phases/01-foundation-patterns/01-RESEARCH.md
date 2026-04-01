# Phase 1 Research: Foundation Patterns - REST API Ecosystem

**Phase:** 1 - Foundation Patterns  
**Topic:** RESTful APIs with Go (Task 1)  
**Research Date:** 2026-04-01  
**Confidence:** MEDIUM-HIGH

## Summary

Phase 1 requires establishing router/framework patterns for 4 topics: RESTful APIs, CLI Applications, Testing & Debugging, and Gin Web Apps. After analyzing the existing codebase and Go ecosystem in 2026, **chi v5** is recommended for the REST API template, **cobra** for CLI, **testify** for testing, and **gin** for web apps (following existing patterns).

**Primary recommendation:** Use chi router for REST API template — it's `net/http` with routing sugar, idiomatic Go, minimal dependencies, and excellent for teaching.

---

## Standard Stack

### Core Libraries for Phase 1

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| chi v5 | 5.x | HTTP router | Lightweight, stdlib-compatible, idiomatic |
| cobra | 1.8.0 | CLI framework | Industry standard, used in course |
| testify | 1.11.x | Testing assertions/mocks | Already in backend |
| gin | 1.12.0 | Web framework | Already in backend, microservices examples |

### Installation

```bash
# REST API template (chi)
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware

# CLI template (cobra)
go get github.com/spf13/cobra@v1.8.0

# Testing template (testify)
go get github.com/stretchr/testify@v1.11.1
```

---

## Architecture Patterns

### Recommended Project Structure

```
basic/projects/rest-api/           # Topic 1: REST API
├── cmd/
│   └── server/main.go             # Entry point
├── internal/
│   ├── handler/                   # HTTP handlers
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── service/                   # Business logic
│   │   ├── service.go
│   │   └── service_test.go
│   └── repository/                # Data access
│       ├── repository.go
│       └── memory.go
├── pkg/
│   └── errors/                    # Custom errors
├── migrations/                    # DB migrations (future)
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml
├── go.mod                         # go 1.23
├── go.sum
├── README.md
└── Makefile

basic/projects/cli-app/            # Topic 2: CLI
├── cmd/
│   └── cli/main.go                # Cobra entry
├── internal/
│   ├── commands/                  # Cobra commands
│   └── config/                    # Config handling
├── pkg/
│   └── weather/                   # Core package
├── go.mod
└── ...

basic/projects/testing-patterns/    # Topic 3: Testing
├── internal/
│   ├── handler/                   # Code under test
│   └── service/
├── mocks/                         # Mock implementations
├── HTTP_TESTING.md                # HTTP testing patterns
├── go.mod
└── ...

basic/projects/gin-web/             # Topic 4: Gin
├── cmd/
│   └── server/main.go
├── internal/
│   ├── handler/
│   ├── middleware/
│   └── views/
├── static/                         # Static assets
├── go.mod
└── ...
```

### REST API Pattern with chi

```go
// cmd/server/main.go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // Routes
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/health", healthHandler)
        r.Route("/users", func(r chi.Router) {
            r.Post("/", createUser)
            r.Get("/", listUsers)
            r.Get("/{id}", getUser)
        })
    })
    
    http.ListenAndServe(":8080", r)
}
```

### Gin Web App Pattern (from existing code)

```go
// cmd/server/main.go (from advanced-topics/06-microservices-docker)
func main() {
    router := gin.New()
    router.Use(gin.Recovery())
    
    // Custom middleware
    router.Use(loggingMiddleware(log))
    router.Use(correlationMiddleware())
    
    // Routes
    v1 := router.Group("/api/v1")
    {
        v1.GET("/health", healthCheck)
    }
    
    router.Run(":8080")
}
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| HTTP routing | Custom router | chi v5 | Handles path params, method matching, middleware chain |
| CLI argument parsing | Manual os.Args | cobra | Subcommands, flags, help generation |
| Testing assertions | if conditions | testify assert/require | Better failure messages, rich matchers |
| HTTP error responses | ad-hoc maps | custom error types | Consistent JSON structure |
| JSON validation | manual parsing | gin binding | struct tags, built-in validators |

---

## Common Pitfalls

### Pitfall 1: Using gorilla/mux in 2026

**What goes wrong:** Gorilla/mux entered maintenance mode in 2023. No new features, security patches only.

**Why it happens:** It was the de facto standard for years, old tutorials still recommend it.

**How to avoid:** Use chi v5 instead — actively maintained, lighter, more idiomatic.

**Warning signs:** `github.com/gorilla/mux` in go.mod for new projects.

### Pitfall 2: Mixing Router Paradigms

**What goes wrong:** Using gin middleware with chi, or vice versa, causes confusion.

**Why it happens:** Each router has its own middleware signature.

**How to avoid:** Chi for REST API template, gin for Gin Web App template — keep separate.

### Pitfall 3: Skipping Middleware

**What goes wrong:** No logging, no recovery, no request ID — production debugging becomes impossible.

**How to avoid:** Always include at minimum: RequestID, Logger, Recoverer, Timeout.

---

## Code Examples

### chi Router with Middleware Chain

```go
// Source: Existing patterns + chi documentation
import (
    "net/http"
    chi "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
    r := chi.NewRouter()
    
    // Core middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.StripSlashes)
    r.Use(middleware.Timeout(30 * time.Second))
    
    // CORS (using chi/cors or manual)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"https://*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    }))
    
    return r
}
```

### testify Mock Pattern (from backend)

```go
// Source: backend/internal/handler/handler_test.go
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

// In test:
courseService.On("GetCourseByID", mock.Anything, "course-123").Return(course, nil)
```

### Cobra CLI Pattern (from course/AT-02-cli-apps)

```go
// Source: course/advanced-topics/AT-02-cli-apps/examples/cobra_cli.go
var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "A brief description",
    Run:   func(cmd *cobra.Command, args []string) {},
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Run:   runServe,
}

func init() {
    rootCmd.AddCommand(serveCmd)
    serveCmd.Flags().StringVar(&port, "port", "8080", "Port to listen on")
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| gorilla/mux | chi v5 | 2023-2024 | Lighter, idiomatic, active maintenance |
| manual testing | testify assertions | 2018+ | Better failure messages |
| stdlib flag | cobra/urfave | 2015+ | Subcommands, composable |
| gin-only | chi + gin selection | 2024+ | chi for APIs, gin for web apps |

**Deprecated/outdated:**
- `github.com/gorilla/sessions` — Use chi sessions or go-jwt
- `github.com/gorilla/schema` — Use gin binding or standard json

---

## Open Questions

1. **Should REST API template use chi or stdlib + chi?**
   - What we know: chi is just net/http with routing helpers
   - What's unclear: Whether to show stdlib first or jump to chi
   - Recommendation: Show stdlib basics in course, chi in template

2. **Database integration in templates?**
   - What we know: Phase 1 topics don't require DB per ROADMAP
   - What's unclear: Whether to include in-memory or skip persistence
   - Recommendation: In-memory for REST API, PostgreSQL integration as exercise

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | All templates | ✓ | 1.26.1 | Go 1.23+ required |
| Docker | All templates | ✓ | 24.x | — |
| cobra | CLI template | ✓ (via go get) | 1.8.0 | — |
| chi | REST template | ✓ (via go get) | 5.x | — |
| testify | Testing template | ✓ (via go get) | 1.11.x | — |

**Missing dependencies with no fallback:** None identified.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` + testify |
| Config file | N/A (per-module) |
| Quick run command | `go test ./... -v -short` |
| Full suite command | `go test ./... -v -race -cover` |

### Phase Requirements → Test Map
| Topic | Test Type | Quick Command | File Exists |
|-------|-----------|---------------|-------------|
| REST API | Unit + Handler | `go test ./internal/handler/... -v` | Will create |
| CLI | Integration | `go test ./... -v` | Will create |
| Testing | Unit | `go test ./... -v` | Will create |
| Gin Web | Unit + Handler | `go test ./internal/... -v` | Will create |

### Wave 0 Gaps
- [ ] `basic/projects/rest-api/internal/handler/handler_test.go` — REQ-1 REST handlers
- [ ] `basic/projects/cli-app/cmd/cli/main_test.go` — REQ-2 CLI commands
- [ ] `basic/projects/testing-patterns/internal/service/service_test.go` — REQ-3 mocks
- [ ] `basic/projects/gin-web/internal/handler/handler_test.go` — REQ-4 Gin handlers

---

## Sources

### Primary (HIGH confidence)
- Existing codebase: `backend/go.mod` — gin 1.12.0, testify 1.11.1, gorilla/mux 1.8.1
- Existing code: `advanced-topics/15-graphql-gqlgen/server.go` — chi v5 patterns
- Existing code: `advanced-topics/06-microservices-docker/` — gin patterns
- Existing tests: `backend/internal/handler/handler_test.go` — testify patterns
- Course example: `course/advanced-topics/AT-02-cli-apps/examples/` — cobra patterns

### Secondary (MEDIUM confidence)
- Go package repositories — chi v5 active development confirmed
- Cobra release history — v1.8.0 current, stable

### Tertiary (LOW confidence)
- Market usage statistics — training data, not verified

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — based on existing codebase + active package repos
- Architecture: HIGH — based on existing patterns in advanced-topics
- Pitfalls: MEDIUM — based on known ecosystem evolution

**Research date:** 2026-04-01
**Valid until:** 2026-07-01 (30 days for stable ecosystem)
