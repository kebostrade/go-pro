# Module Guide - Working with GO-PRO Modules

Complete reference for working with each independent Go module in GO-PRO.

## Module Overview

GO-PRO uses **multiple independent Go modules**, each with its own `go.mod` file.

```
go-pro/
├── backend/                    # go.mod v1.23 - Learning Platform API
├── course/                     # go.mod v1.23 - Course Content
├── basic/                      # go.mod v1.23 - Examples & Exercises
├── services/ai-agent-platform/ # go.mod v1.23 - Production AI Framework
├── services/api-gateway/       # go.mod v1.23 - API Gateway
├── services/shared/            # go.mod v1.23 - Shared Libraries
└── frontend/                   # package.json - Next.js Dashboard
```

**Critical Rule**: Always `cd` to the correct module directory before running Go commands.

---

## 1. Backend Module

**Location**: `/backend`
**Purpose**: REST API server for the learning platform
**Language**: Go 1.23
**Status**: Production-ready

### Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── domain/               # Business models and errors
│   ├── handler/              # HTTP handlers (controllers)
│   ├── middleware/           # HTTP middleware
│   ├── repository/           # Data access layer
│   │   ├── interfaces.go     # All repository contracts
│   │   ├── memory_simple.go  # In-memory implementation
│   │   └── postgres/         # PostgreSQL implementation
│   ├── service/              # Business logic
│   └── util/                 # Utilities
├── pkg/                      # Public packages
├── tests/                    # Test files
├── go.mod                    # Dependencies
├── Makefile                  # Build commands
└── .env.example              # Environment template
```

### Quick Commands

```bash
cd backend

# Run development server
go run ./cmd/server

# Run tests
go test ./...
go test -v ./internal/handler/...
go test -race ./...  # Detect race conditions

# Linting and formatting
make lint
make fmt
go vet ./...

# Build for production
go build -o ./bin/go-pro-api ./cmd/server

# Hot reload development
air  # Requires: go install github.com/cosmtrek/air@latest
```

### Environment Setup

Create `backend/.env`:
```env
PORT=8080
ENV=development
DATABASE_URL=postgres://user:password@localhost:5432/goprodb
JWT_SECRET=your-secret-key-here
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### Testing the API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all courses
curl http://localhost:8080/api/v1/courses

# Get specific exercise
curl http://localhost:8080/api/v1/exercises/1

# Submit exercise solution
curl -X POST http://localhost:8080/api/v1/exercises/1/submit \
  -H "Content-Type: application/json" \
  -d '{
    "code": "package main\nfunc main() { ... }",
    "userId": "user123"
  }'
```

### Key Patterns

#### Repository Pattern
```go
// Define interface (contracts)
type ExerciseRepository interface {
    Save(ctx context.Context, ex *Exercise) error
    FindByID(ctx context.Context, id string) (*Exercise, error)
}

// Multiple implementations can be used
// - memory_simple.go: Development/testing
// - postgres/exercise.go: Production
```

#### Service Layer
```go
// Business logic independent of HTTP
type ExerciseService struct {
    repo ExerciseRepository
}

func (s *ExerciseService) ValidateAndSubmit(ctx context.Context, code string) error {
    // Validate code
    // Compile and test
    // Store result
}
```

#### Handler Layer
```go
// HTTP interface
func (h *Handler) SubmitExercise(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call service
    // Write response
}
```

### Adding New Features

1. **Add Handler**: `internal/handler/new_feature.go`
2. **Add Service**: `internal/service/new_feature.go`
3. **Add Repository**: `internal/repository/interfaces.go` (add interface)
4. **Implement Repository**: `internal/repository/memory_simple.go` and `postgres/new_feature.go`
5. **Add Domain Models**: Update `internal/domain/models.go`
6. **Register Route**: In `cmd/server/main.go`
7. **Test Everything**: Write tests at each layer

---

## 2. Course Module

**Location**: `/course`
**Purpose**: Course content, lessons, exercises, and projects
**Language**: Go 1.23 (code examples and tests)
**Status**: Active

### Directory Structure

```
course/
├── lessons/
│   ├── lesson-01/
│   │   ├── README.md         # Lesson content
│   │   ├── examples.go       # Code examples
│   │   └── concepts.md       # Key concepts
│   ├── lesson-02/
│   └── ...
│
├── code/
│   ├── lesson-01/
│   │   ├── exercises/
│   │   │   ├── 01_fizzbuzz.go
│   │   │   ├── 02_palindrome.go
│   │   │   └── *_solution.go  # Solutions
│   │   ├── solutions/
│   │   │   └── main.go       # Reference implementation
│   │   ├── *_test.go         # Tests
│   │   └── README.md
│   └── ...
│
├── projects/
│   ├── cli-task-manager/
│   ├── rest-api/
│   └── microservices/
│
├── syllabus.md               # Complete curriculum
├── README.md                 # Course overview
└── go.mod
```

### Quick Commands

```bash
cd course

# Read lesson content
cat lessons/lesson-01/README.md

# Run all tests
go test ./...

# Test specific lesson
go test ./code/lesson-01/...

# Run solution examples
go run ./code/lesson-01/solutions/main.go

# Check your progress
go test -v ./code/lesson-01/exercises/...
```

### Course Structure

**20 Progressive Lessons**:
- **Lessons 1-5**: Foundations (syntax, types, functions)
- **Lessons 6-10**: Intermediate (structs, interfaces, concurrency)
- **Lessons 11-15**: Advanced (testing, HTTP, databases)
- **Lessons 16-20**: Expert (performance, security, systems)

### Adding New Content

1. Create lesson directory: `lessons/lesson-XX/`
2. Create lesson README: `lessons/lesson-XX/README.md`
3. Create exercises: `code/lesson-XX/exercises/`
4. Create solutions: `code/lesson-XX/solutions/`
5. Create tests: `code/lesson-XX/*_test.go`
6. Update syllabus: `syllabus.md`
7. Update course README: `README.md`

### Exercise Format

```go
// exercises/01_fizzbuzz.go
// Student starts with this empty signature:
func FizzBuzz(n int) string {
    // TODO: Implement
}

// solutions/main.go shows the reference:
func FizzBuzz(n int) string {
    if n%15 == 0 {
        return "FizzBuzz"
    }
    // ...
}

// Tests verify the implementation:
func TestFizzBuzz(t *testing.T) {
    tests := []struct {
        input    int
        expected string
    }{
        {15, "FizzBuzz"},
        {3, "Fizz"},
        {5, "Buzz"},
        {1, "1"},
    }
    // ...
}
```

---

## 3. Basic Module

**Location**: `/basic`
**Purpose**: Learning examples, standalone exercises, projects
**Language**: Go 1.23
**Status**: Active

### Directory Structure

```
basic/
├── examples/
│   ├── fun/
│   │   ├── binary_search.go
│   │   ├── linked_list.go
│   │   └── main.go           # Dispatcher
│   └── data-structures/
│
├── exercises/
│   ├── 01_basics/
│   ├── 02_functions/
│   └── ...
│
└── projects/
    ├── weather-cli/          # Go 1.23 module
    │   ├── go.mod
    │   ├── main.go
    │   └── ...
    ├── task-manager/         # Go 1.23 module
    └── rest-api/             # Go 1.23 module
```

### Quick Commands

```bash
cd basic

# Run specific example
go run examples/fun/binary_search.go

# Use example dispatcher
cd examples/fun && go run main.go

# Run tests
go test ./...

# Work on exercises
cd exercises/01_basics
go run fibonacci.go
```

### Working with Projects

Each project is an **independent module**:

```bash
# Weather CLI
cd basic/projects/weather-cli
go mod tidy
go build
./weather-cli --help

# Task Manager
cd basic/projects/task-manager
go run main.go
```

### Example Format

```go
// examples/fun/binary_search.go
// Standalone, runnable example with comments

package main

import "fmt"

func BinarySearch(arr []int, target int) int {
    // Implementation with comments
}

func main() {
    arr := []int{1, 3, 5, 7, 9}
    result := BinarySearch(arr, 5)
    fmt.Printf("Found at index: %d\n", result)
}

// Run with: go run binary_search.go
```

---

## 4. AI Agent Platform Module

**Location**: `/services/ai-agent-platform`
**Purpose**: Production-ready AI agent framework
**Language**: Go 1.23
**Status**: Production (example use case)

### Directory Structure

```
services/ai-agent-platform/
├── internal/
│   ├── agent/
│   │   ├── base.go           # Base agent implementation
│   │   └── react.go          # ReAct agent pattern
│   ├── llm/
│   │   ├── openai.go         # OpenAI provider
│   │   └── cache.go          # LLM response caching
│   └── tools/
│       ├── financial/        # Financial tools
│       └── general/          # General tools
│
├── pkg/
│   ├── types/                # Shared interfaces
│   └── errors/               # Error types
│
├── examples/
│   ├── fraud_detection.go
│   ├── sentiment_analysis.go
│   └── main.go
│
├── go.mod
├── Makefile
└── .env.example
```

### Quick Commands

```bash
cd services/ai-agent-platform

# Install dependencies
make deps

# Run fraud detection example
make run-example

# Run tests
make test

# View coverage
make test-coverage

# Run benchmarks
make bench

# Docker support
make docker-build
```

### Environment Setup

Create `.env`:
```env
OPENAI_API_KEY=your-key
MODEL=gpt-4
LOG_LEVEL=info
REDIS_URL=redis://localhost:6379
```

### Key Components

#### Agent Interface
```go
// pkg/types/agent.go
type Agent interface {
    Run(ctx context.Context, input AgentInput) (AgentOutput, error)
    Stream(ctx context.Context, input AgentInput) (<-chan AgentOutput, error)
}
```

#### Tool System
```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (string, error)
}
```

---

## 5. API Gateway Module

**Location**: `/services/api-gateway`
**Purpose**: Route requests to microservices
**Language**: Go 1.23
**Status**: Experimental

### Quick Start

```bash
cd services/api-gateway

go run main.go
# Gateway running at :9000
```

---

## 6. Shared Module

**Location**: `/services/shared`
**Purpose**: Shared libraries for microservices
**Language**: Go 1.23
**Status**: Experimental

### Usage

```go
import "github.com/DimaJoyti/go-pro/services/shared/pkg/log"

logger := log.NewLogger()
logger.Info("Message")
```

---

## 7. Frontend Module

**Location**: `/frontend`
**Purpose**: Next.js learning dashboard
**Language**: TypeScript/React
**Status**: Production

### Quick Commands

```bash
cd frontend

# Install dependencies
bun install

# Development server
bun run dev
# Open: http://localhost:3000

# Build for production
bun run build

# Production server
bun start

# Testing
bun test

# Linting
bun run lint

# Type checking
bun run type-check
```

### Environment Setup

Create `frontend/.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_ENV=development
```

---

## Common Tasks Across Modules

### Run All Tests

```bash
# Backend tests
cd backend && go test ./...

# Course tests
cd course && go test ./...

# Basic tests
cd basic && go test ./...

# Frontend tests
cd frontend && bun test
```

### Update Dependencies

```bash
# Go modules
cd <module-dir>
go get -u ./...
go mod tidy

# Frontend
cd frontend
bun pm update
```

### Build for Production

```bash
# Backend
cd backend && go build -o bin/go-pro-api ./cmd/server

# Frontend
cd frontend && bun run build

# AI Platform
cd services/ai-agent-platform && go build
```

### Clean Build

```bash
# Remove build artifacts
cd <module-dir>
go clean

# Remove module cache
go clean -modcache

# Frontend
cd frontend && rm -rf .next
```

---

## Dependency Management

### Viewing Dependencies

```bash
cd <module-dir>

# List all dependencies
go list -m all

# Check for outdated packages
go list -u -m all

# Check for security vulnerabilities
govulncheck ./...  # Requires: go install golang.org/x/vuln/cmd/govulncheck@latest
```

### Adding Dependencies

```bash
cd <module-dir>
go get github.com/username/package@v1.0.0
go mod tidy
```

### Removing Unused Dependencies

```bash
cd <module-dir>
go mod tidy  # Automatic removal of unused imports
```

---

## Module-Specific Configuration

### Environment Variables

Each module can have its own `.env` file:

```
backend/.env
course/.env (if needed)
services/ai-agent-platform/.env
frontend/.env.local
```

### Configuration Files

```
backend/config.yaml
services/ai-agent-platform/config.yaml
```

---

## Troubleshooting Module Issues

### Issue: "go.mod: not found"
```bash
# You're in wrong directory
# Check current module: cat go.mod | head -1
cd path/to/correct/module
```

### Issue: Dependency conflicts
```bash
cd <module-dir>
go mod tidy
go mod verify
```

### Issue: Build fails
```bash
# Clear cache and rebuild
go clean -modcache
go mod download
go build ./...
```

### Issue: Tests fail
```bash
# Run with verbose output
go test -v ./...

# Check for race conditions
go test -race ./...

# Check code coverage
go test -cover ./...
```

---

## Module Communication

Modules communicate through:

1. **APIs** (backend exposes HTTP API)
2. **Shared packages** (services/shared)
3. **Event queues** (future: message-driven architecture)

**Example**: Frontend calls Backend API
```
Frontend (port 3000)
  └─ HTTP Request
     └─ Backend API (port 8080)
        └─ Service Layer
           └─ Repository Layer
              └─ Database
```

---

## Best Practices

1. ✅ **Always cd to module directory** before running Go commands
2. ✅ **Use go mod tidy** after adding/removing dependencies
3. ✅ **Run tests before committing** code
4. ✅ **Use Makefile commands** for common tasks
5. ✅ **Keep .env files** out of version control
6. ✅ **Document new modules** with README.md

---

## Related Documentation

- [Getting Started](GETTING_STARTED.md) - Quick start guide
- [Architecture](ARCHITECTURE.md) - System architecture
- [Testing Guide](TESTING_GUIDE.md) - Testing strategies
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues

---

**Questions?** Check the module's README.md or CLAUDE.md in project root.
