# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a **dual-purpose Go learning platform** containing:
1. **Learning Content**: Progressive Go tutorials, exercises, and projects for learners
2. **Production AI Agent Platform**: Financial services AI framework (FinAgent) as a real-world example

## Architecture & Module Structure

### Multi-Module Repository
This repository uses **multiple independent Go modules**, all standardized on **Go 1.23**:

```
go-pro/
├── basic/                       # Go 1.23 - Learning examples and exercises
├── basic/projects/*/            # Go 1.23 - Independent project modules
├── backend/                     # Go 1.23 - Learning platform API
├── services/ai-agent-platform/  # Go 1.23 - Production AI agent framework
├── services/api-gateway/        # Go 1.23 - API Gateway service
├── services/shared/             # Go 1.23 - Shared libraries
└── course/                      # Go 1.23 - Course content module
```

**Critical**: Each module has its own `go.mod`. Always run `go` commands from the correct module directory.

### Backend Architecture (Learning Platform API)

The backend follows **Clean Architecture** with strict separation:

```
backend/
├── cmd/server/              # Application entry point
├── internal/                # Private application code
│   ├── config/             # Configuration management
│   ├── domain/             # Business entities (models)
│   ├── handler/            # HTTP handlers (controllers)
│   ├── middleware/         # HTTP middleware (auth, rate limit, logging)
│   └── repository/         # Data access layer
│       ├── interfaces.go   # Repository contracts
│       ├── memory_simple.go # In-memory implementation
│       └── postgres/       # PostgreSQL implementation
└── pkg/                     # Public, reusable packages
```

**Key Pattern**: Repository interfaces in `internal/repository/interfaces.go` allow switching between in-memory and PostgreSQL implementations without changing business logic.

### AI Agent Platform Architecture

Production-ready framework with distinct layers:

```
services/ai-agent-platform/
├── internal/
│   ├── agent/              # Agent implementations (ReAct, Base)
│   ├── llm/                # LLM providers (OpenAI, caching)
│   └── tools/              # Tool system (financial, general)
├── pkg/
│   ├── types/              # Shared types and interfaces
│   └── errors/             # Error handling
└── examples/               # Usage examples (fraud detection, etc.)
```

**Design Philosophy**: Interface-driven design for extensibility. All agents implement `pkg/types/agent.go` interfaces.

## Development Workflows

### Working with Examples

Examples are **standalone, runnable Go files**:

```bash
# Run any example directly
go run basic/examples/fun/binary_search.go
go run basic/examples/fun/linked_list.go

# Or use the main.go dispatcher
cd basic/examples/fun && go run main.go
```

### Working with Exercises

Exercises follow a **solution pattern**:

```bash
# Student files: basic/exercises/01_basics/fizzbuzz.go
# Solutions: basic/exercises/01_basics/fizzbuzz_solution.go

# Test a solution
go run basic/exercises/01_basics/fizzbuzz_solution.go
```

### Working with Projects

Each project is an **independent module**:

```bash
# Weather CLI project
cd basic/projects/weather-cli
go mod tidy
go run .

# AI Agent Platform
cd services/ai-agent-platform
make deps
make run-example
```

### Backend Development

```bash
# Development with hot reload
make dev              # Uses Air for hot reload

# Testing
cd backend && go test ./...                    # All tests
cd backend && go test ./internal/handler/...   # Specific package
cd backend && go test -v -race -cover ./...    # Verbose with race detection

# Linting and formatting
make lint             # Run golangci-lint
make lint-fix         # Auto-fix issues
make fmt              # Format code with gofmt + goimports

# Security
make security         # Run gosec security scanner
make vuln-check       # Check for vulnerabilities with govulncheck

# Build
make build            # Build for Linux AMD64
make build-all        # Multi-platform build
```

### AI Agent Platform Development

```bash
cd services/ai-agent-platform

# Development
make deps             # Download dependencies
make run              # Run main agent server
make run-example      # Run fraud detection example

# Testing
make test             # Unit tests with race detection
make test-coverage    # Generate coverage report
make bench            # Run benchmarks

# Quality
make lint             # Run linter
make fmt              # Format code
make vet              # Run go vet

# Docker
make docker-build     # Build Docker image
make docker-compose-up # Start all services (Redis, PostgreSQL, etc.)
```

## Testing Strategy

### Comprehensive Test Suite

Run all examples, exercises, and projects:

```bash
./test-all.sh              # Tests everything
./test-all-examples.sh     # Tests only examples
```

### Testing Individual Components

```bash
# Backend API tests
cd backend && go test -v ./internal/handler/...
cd backend && go test -v -tags=integration ./...

# AI Agent tests
cd services/ai-agent-platform && go test -v ./internal/agent/...

# Example files (standalone)
go run basic/examples/fun/merge_sort.go
```

### Test Organization Patterns

- **Backend**: Tests live alongside code (`handler_test.go` next to `handler.go`)
- **Examples**: Self-contained, executable files with clear output
- **Exercises**: Solutions include comments explaining approach

## Key Development Patterns

### Repository Pattern (Backend)

All data access goes through repository interfaces:

```go
// Define interface
type CourseRepository interface {
    Create(ctx context.Context, course *domain.Course) error
    FindByID(ctx context.Context, id string) (*domain.Course, error)
}

// Implementations: memory_simple.go, postgres/course.go
```

This allows easy switching between in-memory (development/testing) and PostgreSQL (production).

### Agent Pattern (AI Platform)

All agents implement common interfaces from `pkg/types/agent.go`:

```go
type Agent interface {
    Run(ctx context.Context, input AgentInput) (AgentOutput, error)
    Stream(ctx context.Context, input AgentInput) (<-chan AgentOutput, error)
}
```

Specific implementations: `internal/agent/base.go`, `internal/agent/react.go`

### Tool System (AI Platform)

Tools are registered in `internal/tools/registry.go`:

```go
// Financial tools: internal/tools/financial/
// General tools: internal/tools/general/
```

Each tool implements the `Tool` interface from `pkg/types/tool.go`.

## Module-Specific Commands

### Root Level (Platform Management)

```bash
make help              # Show all available commands
make build             # Build backend
make test              # Test backend
make docker-dev        # Start full development environment
make quality           # Run all quality checks (lint, vet, security, test)
```

### Backend Module

```bash
cd backend
go test ./...                           # Run tests
go test -v -race ./internal/handler/... # Test handlers with race detection
go run ./cmd/server                     # Run API server
```

### AI Agent Platform Module

```bash
cd services/ai-agent-platform
make test              # Run all tests
make example-fraud     # Run fraud detection example
make dev               # Development mode with hot reload (needs Air)
```

### Examples Module

```bash
cd basic/examples/fun
go run binary_search.go    # Run specific example
go test ./...              # Run tests if any
```

## Environment Setup

### Backend API (.env required)

```bash
cd backend
cp .env.example .env
# Edit .env with database credentials, JWT secrets, etc.
```

### AI Agent Platform (.env required)

```bash
cd services/ai-agent-platform
cp .env.example .env
# Edit .env with OpenAI API key, model settings, etc.
```

## Docker Development

### Full Stack Development

```bash
make docker-dev
# Starts: Backend API, Frontend, PostgreSQL, Redis, Adminer, Prometheus, Grafana
# Services:
#   - Frontend: http://localhost:3000
#   - Backend API: http://localhost:8080
#   - Adminer: http://localhost:8081
```

### AI Agent Platform

```bash
cd services/ai-agent-platform
make docker-compose-up
# Starts: PostgreSQL with pgvector, Redis, Qdrant
```

## Frontend Development

### Setup and Run

```bash
cd frontend
npm install
cp .env.example .env.local
# Edit .env.local with your configuration
npm run dev
# Frontend runs on http://localhost:3000
```

### Frontend Structure

```
frontend/
├── src/
│   ├── app/              # Next.js app router pages
│   ├── components/       # React components
│   ├── contexts/         # React contexts (auth, etc.)
│   ├── lib/              # Utilities and API client
│   └── styles/           # CSS and styling
├── public/               # Static assets
└── package.json          # Dependencies
```

### API Integration

The frontend uses a centralized API client (`src/lib/api.ts`) that connects to the backend:

```typescript
import { api } from '@/lib/api';

// Fetch curriculum
const curriculum = await api.getCurriculum();

// Get lesson details
const lesson = await api.getLessonDetail(1);

// Update progress
await api.updateProgress(userId, lessonId, { completed: true, score: 95 });
```

## Dependencies & Tools

### Required Tools

- **Go**: Version 1.23 (standardized across all modules)
- **Node.js**: Version 18+ (for frontend)
- **Make**: For build automation
- **Docker & Docker Compose**: For containerized development

### Development Tools (Auto-installed)

```bash
make install-tools     # Installs: golangci-lint, gosec, air, goimports, govulncheck
```

- **golangci-lint**: Linting with multiple linters
- **gosec**: Security vulnerability scanner
- **govulncheck**: Known vulnerability detection
- **air**: Hot reload for development
- **goimports**: Import management and formatting

## Common Pitfalls

1. **Wrong Module Directory**: Always `cd` to the correct module before running Go commands
2. **Missing .env Files**: Backend and AI platform require environment configuration
3. **Go Version**: All modules now use Go 1.23 for consistency
4. **Dependency Issues**: Run `go mod tidy` in the specific module directory, not root
5. **Frontend-Backend Connection**: Ensure `NEXT_PUBLIC_API_URL` is set correctly in frontend `.env.local`

## CI/CD Integration

### CI Pipeline Commands

```bash
make ci-test           # Run full CI test suite
make ci-build          # Run CI build pipeline (includes tests)
make release           # Create versioned release
```

### Quality Gates

The repository enforces quality through:
- **Pre-commit hooks**: Run via `make init-git`
- **Linting**: golangci-lint with strict rules
- **Security scanning**: gosec for vulnerability detection
- **Test coverage**: Coverage reports generated with `make test-coverage`

## Documentation Resources

- **Learning Paths**: See `LEARNING_PATHS.md` for structured curriculum
- **Projects Guide**: See `PROJECTS.md` for all project descriptions
- **Tutorials**: See `TUTORIALS.md` for step-by-step guides
- **API Reference**: Backend API docs at `http://localhost:8080` when running
- **AI Agent Docs**: See `services/ai-agent-platform/README.md`
- **Frontend-Backend Integration**: See `FRONTEND_BACKEND_INTEGRATION.md` for full stack setup
