# Codebase Structure

**Analysis Date:** 2026-03-31

## Directory Layout

```
go-pro/                                # Repository root
├── backend/                           # Learning Platform API (Go module)
│   ├── cmd/server/main.go             # Application entry point
│   ├── internal/                      # Private application code
│   │   ├── agents/                    # Multi-agent debugging system
│   │   ├── auth/                      # Auth utilities
│   │   ├── cache/                     # Redis cache manager
│   │   ├── circuitbreaker/            # Circuit breaker pattern
│   │   ├── config/                    # Environment config loading
│   │   ├── container/                 # Dependency injection container
│   │   ├── domain/                    # Business entities and models
│   │   ├── errors/                    # Typed API errors
│   │   ├── executor/                  # Docker-based code execution
│   │   ├── graceful/                  # Graceful shutdown utilities
│   │   ├── handler/                   # HTTP handlers (controllers)
│   │   ├── integration/               # Integration tests
│   │   ├── logger/                    # Internal logger
│   │   ├── messaging/                 # Kafka messaging service
│   │   ├── middleware/                # HTTP middleware chain
│   │   ├── observability/             # OpenTelemetry telemetry
│   │   ├── repository/                # Data access layer
│   │   │   ├── interfaces.go          # Repository contracts
│   │   │   ├── memory_simple.go       # In-memory implementation
│   │   │   ├── memory_interview.go    # In-memory interview store
│   │   │   ├── memory_lesson.go       # In-memory lesson store
│   │   │   ├── user_repository.go     # User repository helper
│   │   │   └── postgres/              # PostgreSQL implementation
│   │   │       ├── connection.go      # DB connection pool
│   │   │       ├── migration.go       # Schema migrations
│   │   │       ├── repositories.go    # PG repo factory
│   │   │       ├── querybuilder.go    # SQL query builder
│   │   │       ├── transaction.go     # Transaction helper
│   │   │       ├── monitor.go         # DB monitoring
│   │   │       └── *.go               # Entity-specific repos
│   │   ├── service/                   # Business logic layer
│   │   │   ├── interfaces.go          # Service contracts
│   │   │   └── *.go                   # Service implementations
│   │   └── testutil/                  # Test utilities
│   ├── pkg/                           # Public reusable packages
│   │   ├── logger/                    # Structured logger
│   │   └── validator/                 # Request validator
│   └── go.mod                         # go-pro-backend (Go 1.25)
│
├── frontend/                          # Next.js frontend application
│   ├── src/
│   │   ├── app/                       # Next.js App Router pages
│   │   │   ├── layout.tsx             # Root layout (AuthProvider wrapper)
│   │   │   ├── page.tsx               # Home page
│   │   │   ├── globals.css            # Global styles
│   │   │   ├── learn/                 # Learning content pages
│   │   │   ├── practice/              # Practice and exercises
│   │   │   ├── playground/            # Code playground
│   │   │   ├── interviews/            # Interview practice
│   │   │   ├── dashboard/             # User dashboard
│   │   │   ├── community/             # Community features
│   │   │   ├── projects/              # Project showcase
│   │   │   ├── cms/                   # Content management
│   │   │   ├── admin/                 # Admin panel
│   │   │   ├── auth/                  # Auth callbacks
│   │   │   ├── signin/                # Sign in page
│   │   │   ├── signup/                # Sign up page
│   │   │   ├── profile/               # User profile
│   │   │   ├── settings/              # User settings
│   │   │   ├── curriculum/            # Curriculum view
│   │   │   ├── exercises/             # Exercise pages
│   │   │   ├── tutorials/             # Tutorial pages
│   │   │   ├── algorithms/            # Algorithm pages
│   │   │   └── customers/             # Customer management
│   │   ├── components/                # React components
│   │   │   ├── layout/                # Header, Footer
│   │   │   ├── auth/                  # Auth UI components
│   │   │   ├── learning/              # Lesson display components
│   │   │   ├── practice/              # Exercise components
│   │   │   ├── ui/                    # Shared UI primitives
│   │   │   ├── dashboard/             # Dashboard widgets
│   │   │   ├── interviews/            # Interview components
│   │   │   ├── playground/            # Code editor (Monaco)
│   │   │   ├── admin/                 # Admin components
│   │   │   ├── community/             # Community components
│   │   │   ├── home/                  # Home page components
│   │   │   ├── examples/              # Code example components
│   │   │   ├── assessments/           # Assessment components
│   │   │   ├── workspace/             # Workspace layout
│   │   │   ├── algorithms/            # Algorithm components
│   │   │   └── customers/             # Customer management
│   │   ├── contexts/                  # React contexts
│   │   │   ├── auth-context.tsx        # Firebase auth context
│   │   │   └── auth-context-advanced.tsx
│   │   ├── lib/                       # Utilities and clients
│   │   │   ├── api.ts                 # Backend API client
│   │   │   ├── firebase.ts            # Firebase initialization
│   │   │   ├── auth-utils.ts          # Auth utility functions
│   │   │   └── api/services/          # Domain-specific API services
│   │   │       ├── community.ts
│   │   │       ├── practice.ts
│   │   │       └── projects.ts
│   │   ├── styles/                    # CSS styles
│   │   └── types/                     # TypeScript type definitions
│   ├── public/                        # Static assets
│   ├── types/                         # Global type definitions
│   └── package.json                   # Next.js 15, React 19, Firebase
│
├── services/                          # Production microservices
│   ├── ai-agent-platform/             # AI Agent Framework (Go module)
│   │   ├── cmd/coding-agent-server/   # Agent server entry point
│   │   ├── internal/
│   │   │   ├── agent/                 # Agent implementations
│   │   │   │   ├── base.go            # Base agent
│   │   │   │   ├── react.go           # ReAct pattern agent
│   │   │   │   └── coding_expert.go   # Coding specialist
│   │   │   ├── api/                   # HTTP server
│   │   │   ├── embeddings/            # OpenAI embeddings
│   │   │   ├── languages/             # Language support (golang/)
│   │   │   ├── llm/                   # LLM providers (OpenAI, cache)
│   │   │   ├── rag/                   # RAG pipeline
│   │   │   ├── sandbox/               # Docker sandbox
│   │   │   ├── tools/                 # Tool system
│   │   │   │   ├── registry.go        # Tool registry
│   │   │   │   ├── financial/         # Fraud check, transactions
│   │   │   │   ├── general/           # Calculator
│   │   │   │   └── programming/       # Code analysis, execution, search
│   │   │   └── vectorstore/           # In-memory vector DB
│   │   ├── pkg/
│   │   │   ├── types/                 # Shared interfaces and types
│   │   │   │   ├── agent.go           # Agent, AgentInput, AgentOutput
│   │   │   │   ├── tool.go            # Tool, ToolRegistry
│   │   │   │   ├── llm.go             # LLM provider types
│   │   │   │   ├── memory.go          # Memory system types
│   │   │   │   ├── vector.go          # Vector store types
│   │   │   │   ├── coding.go          # Coding-specific types
│   │   │   │   └── language.go        # Language types
│   │   │   └── errors/                # Error types
│   │   ├── examples/                  # Usage examples
│   │   │   ├── fraud_detection/       # Financial fraud detection
│   │   │   ├── coding_qa/             # Coding Q&A agent
│   │   │   └── rag_demo/              # RAG demonstration
│   │   └── go.mod                     # Go 1.23
│   │
│   ├── api-gateway/                   # API Gateway (Go module)
│   │   ├── cmd/main.go               # Entry point
│   │   ├── internal/
│   │   │   ├── auth/jwt.go           # JWT authentication
│   │   │   ├── config/config.go      # Environment config
│   │   │   ├── handler/handler.go    # Route setup + health checks
│   │   │   └── proxy/proxy.go        # Reverse proxy router
│   │   └── go.mod                     # Go 1.23, depends on shared
│   │
│   ├── shared/                        # Shared libraries (Go module)
│   │   ├── middleware/middleware.go   # Common HTTP middleware
│   │   ├── client/http_client.go     # HTTP client utilities
│   │   ├── events/events.go          # Event definitions
│   │   └── go.mod                     # Go 1.23
│   │
│   ├── llm-gateway/                   # LLM Gateway (skeleton)
│   │   ├── internal/config/           # Config only
│   │   └── go.mod
│   │
│   └── langchain/                     # LangChain-style Go library
│       ├── pkg/agent/react.go         # ReAct agent
│       ├── pkg/llm/provider.go        # LLM provider
│       ├── pkg/schema/schema.go       # Schema definitions
│       └── go.mod
│
├── basic/                             # Go learning content (Go module)
│   ├── go.mod                         # Go 1.23
│   ├── cmd/                           # Example entry points
│   ├── examples/                      # Runnable code examples
│   │   ├── fun/                       # Algorithm examples (own go.mod)
│   │   ├── hello_go/                  # Hello world (own go.mod)
│   │   └── concurrency-crash-course/  # Concurrency examples
│   ├── exercises/                     # Student exercises with solutions
│   │   └── 01_basics/                 # FizzBuzz, etc.
│   └── projects/                      # Independent project modules
│       ├── weather-cli/               # CLI project
│       ├── url-shortener/             # URL shortener project
│       ├── ai-engineering/chatbot-cli/ # AI chatbot project
│       └── [~18 other projects]/      # Each with own go.mod
│
├── course/                            # Course content module (Go module)
│   ├── go.mod                         # Go 1.23
│   ├── code/                          # Per-lesson code (each lesson own go.mod)
│   │   ├── lesson-01/ through lesson-15/
│   ├── lessons/                       # Lesson content
│   ├── algorithms/                    # Algorithm course material
│   ├── advanced-topics/               # Advanced topic examples
│   ├── prompt-engineering/            # Prompt engineering content
│   ├── openclaw/                      # OpenClaw content
│   ├── llm-ops/                       # LLM operations content
│   └── projects/                      # Course projects
│
├── advanced/                          # Advanced Go examples (collection)
│   ├── [~30 topic directories]/       # Each topic is self-contained
│   └── GoBootcamp/                    # Go Bootcamp exercises
│
├── advanced-topics/                   # Advanced topics (collection)
│   ├── 01-system-design/ through 15-graphql-gqlgen/
│   └── [Each topic has own go.mod(s)]
│
├── observability/                     # Observability infrastructure
│   ├── otel/                          # OpenTelemetry configs
│   ├── configs/                       # Collector configs
│   ├── dashboards/                    # Grafana dashboards
│   └── go.mod                         # Go 1.23
│
├── platform/                          # Platform/DevOps configurations
│   ├── argocd/                        # ArgoCD application manifests
│   ├── charts/                        # Helm charts (ai-agent)
│   ├── kustomize/                     # Kustomize overlays (dev/staging/prod)
│   ├── github-actions/                # CI/CD workflows
│   ├── terraform/                     # Infrastructure as code
│   │   ├── modules/                   # GKE, VPC, CloudSQL, Redis, GCS, IAM, Secrets
│   │   ├── environments/              # dev/ and prod/
│   │   └── shared/                    # Shared Terraform configs
│   ├── policies/                      # OPA Gatekeeper, Kyverno, network policies
│   └── docs/                          # Platform documentation
│
├── .planning/                         # GSD planning documents
│   └── codebase/                      # Codebase analysis documents
│
├── go.mod                             # Root go.mod (Go 1.23)
├── Makefile                           # Build automation
├── CLAUDE.md                          # Claude Code instructions
└── docker-compose*.yml               # Docker orchestration
```

## Module Boundaries

Each module is an independent Go module with its own `go.mod`. Always run `go` commands from the correct module directory.

**Core production modules:**

| Module | go.mod Location | Purpose |
|--------|----------------|---------|
| Backend API | `backend/go.mod` (module `go-pro-backend`) | Learning platform REST API |
| AI Agent Platform | `services/ai-agent-platform/go.mod` | AI agent framework |
| API Gateway | `services/api-gateway/go.mod` | Reverse proxy + JWT auth |
| Shared Libraries | `services/shared/go.mod` | Common middleware, clients, events |
| LLM Gateway | `services/llm-gateway/go.mod` | LLM routing (skeleton) |
| Langchain | `services/langchain/go.mod` | LangChain-style Go library |
| Observability | `observability/go.mod` | Telemetry infrastructure |

**Learning content modules:**

| Module | go.mod Location | Purpose |
|--------|----------------|---------|
| Basic | `basic/go.mod` | Core Go examples and exercises |
| Course | `course/go.mod` | Course content and lesson code |
| Basic Examples | `basic/examples/fun/go.mod` | Algorithm examples |
| Basic Projects | `basic/projects/*/go.mod` | ~18 independent project modules |
| Course Lessons | `course/code/lesson-*/go.mod` | 15 per-lesson code modules |

**Cross-module dependency:**
- `services/api-gateway` depends on `services/shared` via `replace` directive in go.mod

## Entry Points

**Backend API Server:**
- `backend/cmd/server/main.go` - Creates Container, wires handlers to mux, applies middleware chain, starts HTTP server

**AI Agent Platform Server:**
- `services/ai-agent-platform/cmd/coding-agent-server/main.go` - Agent server entry point
- `services/ai-agent-platform/examples/fraud_detection/main.go` - Fraud detection example
- `services/ai-agent-platform/examples/coding_qa/main.go` - Coding Q&A example
- `services/ai-agent-platform/examples/rag_demo/main.go` - RAG demonstration

**API Gateway:**
- `services/api-gateway/cmd/main.go` - Loads config, creates handler, starts HTTP server

**Frontend:**
- `frontend/src/app/layout.tsx` - Next.js root layout with AuthProvider
- `frontend/src/app/page.tsx` - Home page
- `frontend/package.json` scripts: `bun run dev` starts dev server on port 3000

**Root:**
- `Makefile` - Build, test, Docker, and quality commands for the entire platform

## Code Organization Within Modules

### Backend (`backend/`)

Follows Go standard layout with Clean Architecture:

- `cmd/` - Application entry points (only `cmd/server/main.go`)
- `internal/` - Private packages, never imported by other modules
  - `domain/` - Pure Go types, no dependencies
  - `repository/` - Data access interfaces + implementations
  - `service/` - Business logic, depends on repository interfaces only
  - `handler/` - HTTP handlers, depends on service interfaces
  - `middleware/` - HTTP middleware, depends on auth service interface
  - `container/` - DI container, wires everything together
  - Infrastructure packages (`cache/`, `messaging/`, `executor/`, etc.)
- `pkg/` - Public reusable packages (`logger/`, `validator/`)

**New feature pattern:** Domain model -> Repository interface -> Repository impl -> Service interface -> Service impl -> Handler -> Route registration

### AI Agent Platform (`services/ai-agent-platform/`)

- `pkg/types/` - Public interfaces (Agent, Tool, ToolRegistry, LLM, Memory, Vector)
- `internal/agent/` - Agent implementations
- `internal/tools/` - Tool implementations organized by category
- `internal/llm/` - LLM provider implementations
- `internal/sandbox/` - Docker sandbox
- `examples/` - Standalone usage examples

**New tool pattern:** Define struct implementing `Tool` interface in `internal/tools/<category>/` -> Register in `Registry`

### Frontend (`frontend/`)

- `src/app/` - Next.js App Router (file-system based routing)
- `src/components/` - React components organized by feature
- `src/contexts/` - React contexts for global state (auth)
- `src/lib/` - Utilities, API client, Firebase config
- `src/types/` - TypeScript type definitions

**New page pattern:** Create `src/app/<route>/page.tsx` -> Add components in `src/components/<feature>/`

### Learning Content (`basic/`, `course/`, `advanced/`)

- Each project/lesson has its own `go.mod` for independence
- Examples are standalone runnable `.go` files
- Exercises follow `name.go` (student) + `name_solution.go` (solution) pattern

## Shared Code

### `services/shared/` (Go module)

Common microservice utilities shared across the `services/` modules:

- `services/shared/middleware/middleware.go` - HTTP middleware (RequestID, Logger, Recovery, CORS, Timeout, RateLimit, Chain, ServiceInfo, HealthCheck, Metrics, Tracing)
- `services/shared/client/http_client.go` - HTTP client utilities for inter-service communication
- `services/shared/events/events.go` - Event definitions for messaging

Currently only imported by `services/api-gateway` via `replace` directive.

### `backend/pkg/` (Backend-internal public packages)

Reusable within the backend module:

- `backend/pkg/logger/` - Structured logger interface
- `backend/pkg/validator/` - Request validation wrapping `go-playground/validator`

### Cross-cutting patterns duplicated between modules

The middleware chain pattern (`func(http.Handler) http.Handler` + `Chain()`) is implemented independently in:
- `backend/internal/middleware/middleware.go`
- `services/shared/middleware/middleware.go`

Both provide RequestID, Logging, Recovery, CORS, Timeout functionality but with different implementations (backend uses structured logger, shared uses `fmt.Printf`).

## Static Assets

### Course Content

- `course/lessons/` - Lesson markdown/content files
- `course/code/lesson-*/` - Per-lesson code examples (15 lessons)
- `course/algorithms/` - Algorithm course material
- `course/advanced-topics/` - Advanced topic examples with go.mod files
- `course/prompt-engineering/` - Prompt engineering content
- `course/openclaw/` - OpenClaw course content
- `course/llm-ops/` - LLM operations content

### Learning Examples

- `basic/examples/` - Standalone runnable Go examples
- `basic/exercises/` - Student exercises with solution files
- `basic/projects/` - ~18 independent project modules
- `advanced/` - ~30 advanced Go topic directories
- `advanced-topics/` - 15 advanced topic directories (REST APIs, gRPC, K8s, NATS, etc.)

### Frontend Static Assets

- `frontend/public/` - Static web assets (favicons, avatars)
- `frontend/src/styles/` - CSS stylesheets

### Infrastructure Configs

- `platform/argocd/` - ArgoCD manifests
- `platform/charts/` - Helm charts
- `platform/kustomize/` - Kustomize overlays (base, dev, staging, prod)
- `platform/terraform/` - Terraform IaC (GKE, VPC, CloudSQL, Redis, GCS, IAM, Secrets)
- `platform/policies/` - OPA Gatekeeper, Kyverno, network policies, pod security
- `observability/configs/` - OpenTelemetry collector configs
- `observability/dashboards/` - Grafana dashboard definitions

## Where to Add New Code

**New backend API endpoint:**
1. Add domain model in `backend/internal/domain/models.go`
2. Add repository interface in `backend/internal/repository/interfaces.go`
3. Implement in `backend/internal/repository/memory_simple.go` and `backend/internal/repository/postgres/<entity>.go`
4. Add service interface + impl in `backend/internal/service/`
5. Add handler in `backend/internal/handler/`
6. Register route in `handler.RegisterRoutes()` in `backend/internal/handler/handler.go`

**New AI agent tool:**
1. Define tool struct implementing `types.Tool` interface
2. Place in `services/ai-agent-platform/internal/tools/<category>/`
3. Register with `tools.Registry.Register()` or `RegisterWithCategory()`

**New AI agent type:**
1. Define agent struct implementing `types.Agent` interface
2. Place in `services/ai-agent-platform/internal/agent/`

**New frontend page:**
1. Create `frontend/src/app/<route>/page.tsx` for route
2. Add components in `frontend/src/components/<feature>/`
3. Add API calls in `frontend/src/lib/api.ts` or `frontend/src/lib/api/services/`

**New learning content:**
1. Add lesson directory in `course/code/lesson-<N>/` with own `go.mod`
2. Add content in `course/lessons/`
3. Add exercises in `basic/exercises/`

**New shared middleware:**
1. Add to `services/shared/middleware/middleware.go` for cross-service use
2. Or `backend/internal/middleware/middleware.go` for backend-only

## Special Directories

**`.planning/`:**
- Purpose: GSD workflow planning documents
- Generated: Yes (by GSD mapping agents)
- Committed: Yes (part of repo)

**`platform/`:**
- Purpose: DevOps and infrastructure configuration
- Generated: No (manually maintained)
- Committed: Yes

**`frontend/.next/`:**
- Purpose: Next.js build output
- Generated: Yes (by `next build`)
- Committed: No (in .gitignore)

**`frontend/node_modules/`:**
- Purpose: Frontend dependencies
- Generated: Yes (by `bun install`)
- Committed: No (in .gitignore)

**`backend/internal/`:**
- Purpose: Private backend code, not importable by other Go modules
- Go convention: `internal/` packages cannot be imported outside the module

**`backend/pkg/`:**
- Purpose: Public backend packages, importable within the module
- Note: Cannot be imported by other modules since `go-pro-backend` is not published

---

*Structure analysis: 2026-03-31*
