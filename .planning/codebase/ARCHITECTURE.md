# Architecture

**Analysis Date:** 2026-03-31

## High-Level Architecture

**Overall Pattern:** Multi-module monorepo with microservice-oriented services and a learning content library.

The repository is a dual-purpose Go learning platform: (1) progressive Go tutorials, exercises, and projects; (2) a production AI agent platform and learning platform API. Each Go module is independent with its own `go.mod`, and the frontend is a separate Next.js application.

```
                         +---------------------+
                         |   Frontend (Next.js) |
                         |   Port 3000          |
                         +----------+-----------+
                                    |
                          API calls / Firebase Auth
                                    |
              +---------------------+---------------------+
              |                                           |
   +----------v----------+                    +-----------v-----------+
   |  Backend API (Go)    |                    |   API Gateway (Go)    |
   |  Port 8080           |                    |   (Proxy + JWT Auth)   |
   |  Clean Architecture  |                    +-----------+-----------+
   +----------+-----------+                                |
              |                          +-----------------+------------------+
              |                          |                 |                  |
     +--------v--------+        +-------v------+  +-------v------+  +-------v------+
     | In-Memory /     |        | User Service |  |Course Service|  |Progress Svc |
     | PostgreSQL      |        +--------------+  +--------------+  +-------------+
     | Repositories    |
     +--------+--------+
              |
     +--------v--------+
     | AI Agent Pool   |
     | Docker Executor |
     +-----------------+

   +---------------------------------------------------------------+
   |             Independent AI Agent Platform Service             |
   |  Agent Interface -> Tool Registry -> LLM Provider -> Sandbox |
   +---------------------------------------------------------------+
```

## Module Architecture

### Backend (Learning Platform API)

The backend follows **Clean Architecture** with strict layer separation using Go standard `net/http` (with `gin` imported but `http.ServeMux` used for routing in the main handler).

**Layers (innermost to outermost):**

1. **Domain Layer** (`backend/internal/domain/`)
   - Purpose: Core business entities with no external dependencies
   - Contains: `models.go` (Course, Lesson, Exercise, Progress, User, Streak, Assessment, Submission, etc.), request/response types, value objects (Difficulty, Status)
   - Depends on: Nothing (pure Go types and validation tags)
   - Used by: All other layers

2. **Repository Layer** (`backend/internal/repository/`)
   - Purpose: Data access abstraction with interface contracts
   - Contains: `interfaces.go` defines 11 repository interfaces (CourseRepository, LessonRepository, ExerciseRepository, ProgressRepository, UserRepository, StreakRepository, AssessmentRepository, QuestionRepository, SubmissionRepository, SubmissionCommentRepository, PeerReviewRepository, InterviewRepository)
   - Implementations:
     - `memory_simple.go` / `memory_interview.go` / `memory_lesson.go` - In-memory for dev/test
     - `postgres/` - PostgreSQL implementation with migrations
   - Depends on: Domain layer
   - Used by: Service layer
   - Key pattern: `Repositories` struct in `interfaces.go` aggregates all repositories as a single injectable dependency

3. **Service Layer** (`backend/internal/service/`)
   - Purpose: Business logic orchestration
   - Contains: `interfaces.go` defines 8 service interfaces (CourseService, LessonService, ExerciseService, ProgressService, CurriculumService, HealthService, ExecutorService, AuthService, UserService)
   - Implementations: `course.go`, `lesson.go`, `exercise.go`, `progress.go`, `curriculum.go`, `health.go`, `auth.go`, `user.go`, `executor.go`, `local_executor.go`, `exercise_evaluator.go`
   - Depends on: Repository layer, cache, messaging, logger, validator
   - Used by: Handler layer
   - Key pattern: `Services` struct in `interfaces.go` aggregates all services; `NewServices()` factory creates them all

4. **Handler Layer** (`backend/internal/handler/`)
   - Purpose: HTTP request handling, routing, request/response serialization
   - Contains: `handler.go` (main routes), `auth.go`, `admin.go`, `playground.go`, `playground_ai.go`, `interview.go`, `assessment.go`, `submission.go`, `gradebook.go`
   - Routes registered via `http.ServeMux` with Go 1.22+ pattern syntax (e.g., `GET /api/v1/courses/{id}`)
   - Depends on: Service layer, middleware, domain types
   - Used by: Main entry point wires handlers to mux

5. **Middleware Layer** (`backend/internal/middleware/`)
   - Purpose: Cross-cutting HTTP concerns
   - Contains: `middleware.go` (chain, request ID, logging, recovery, CORS, security, timeout, rate limit, pagination, CSRF), `auth.go` (Firebase JWT verification), `ratelimit.go`, `metrics.go`
   - Key pattern: `Middleware` type is `func(http.Handler) http.Handler`; `Chain()` composes them

6. **Container / DI** (`backend/internal/container/`)
   - Purpose: Dependency injection and lifecycle management
   - Contains: `container.go` - wires all layers together
   - Initialization order: Validator -> Cache -> Messaging -> Repositories -> Services -> Agent Pool
   - Graceful shutdown in reverse order

**Infrastructure packages:**
- `backend/internal/cache/` - Redis-backed caching (cache, sessions, distributed locks, rate limiting, pub/sub)
- `backend/internal/messaging/` - Kafka-based messaging with `kafka/` and `realtime/` sub-packages
- `backend/internal/executor/` - Docker-based code execution sandbox
- `backend/internal/agents/` - Multi-agent system (Executor, TestValidator, AIAnalyzer, StateManager) with `AgentPool` for collaborative request processing
- `backend/internal/errors/` - Typed API errors with HTTP status codes
- `backend/internal/observability/` - OpenTelemetry telemetry
- `backend/internal/config/` - Environment-based configuration loading
- `backend/pkg/logger/` - Structured logging abstraction
- `backend/pkg/validator/` - Request validation abstraction

### AI Agent Platform

The AI agent platform is a separate, self-contained Go module at `services/ai-agent-platform/` implementing a production-ready AI agent framework.

**Architecture Layers:**

1. **Types Layer** (`pkg/types/`)
   - `agent.go` - `Agent` interface: `Run()`, `Stream()`, `GetMemory()`, `GetTools()`, `GetConfig()`; plus `AgentInput`, `AgentOutput`, `AgentStep`, `AgentEvent`, `AgentConfig`
   - `tool.go` - `Tool` interface: `Name()`, `Description()`, `Execute()`, `GetSchema()`, `Validate()`; plus `ToolInput`, `ToolOutput`, `ToolRegistry` interface
   - `llm.go` - LLM provider interfaces
   - `memory.go` - Memory system interfaces
   - `vector.go` - Vector store interfaces
   - `coding.go` - Coding-specific types
   - `language.go` - Language-specific types

2. **Agent Implementations** (`internal/agent/`)
   - `base.go` - Base agent with common functionality
   - `react.go` - ReAct (Reason-Act) pattern agent
   - `coding_expert.go` - Specialized coding assistant

3. **LLM Providers** (`internal/llm/`)
   - `provider.go` - Provider abstraction
   - `openai.go` - OpenAI integration
   - `cache.go` - LLM response caching

4. **Tool System** (`internal/tools/`)
   - `registry.go` - Thread-safe tool registry implementing `ToolRegistry` interface with category support
   - `financial/` - `fraud_check.go`, `transaction_lookup.go`
   - `general/` - `calculator.go`
   - `programming/` - `code_analysis.go`, `code_execution.go`, `doc_search.go`, `github_search.go`, `stackoverflow.go`

5. **Supporting Systems:**
   - `internal/embeddings/` - OpenAI embeddings for RAG
   - `internal/rag/` - Retrieval-Augmented Generation pipeline
   - `internal/vectorstore/` - In-memory vector store
   - `internal/sandbox/` - Docker-based code sandbox (`docker.go`, `limits.go`, `security.go`)
   - `internal/languages/` - Language-specific support (`golang/` with analyzer, executor, provider)
   - `internal/api/` - HTTP server for agent platform

### API Gateway

The API Gateway at `services/api-gateway/` is a lightweight reverse proxy with JWT authentication.

- Entry point: `services/api-gateway/cmd/main.go`
- Routing: `services/api-gateway/internal/proxy/proxy.go` - routes requests to backend services (User, Course, Progress)
- Auth: `services/api-gateway/internal/auth/jwt.go` - JWT validation with optional/auth-required middleware
- Config: `services/api-gateway/internal/config/config.go` - service URLs and JWT config from environment
- Uses `services/shared/` for common middleware
- Depends on `services/shared` via `replace` directive in `go.mod`

### Frontend (Next.js)

The frontend at `frontend/` is a Next.js 15 app using the App Router pattern.

- **Entry point:** `frontend/src/app/layout.tsx` - Root layout wrapping all pages with `AuthProvider`
- **Auth:** `frontend/src/contexts/auth-context.tsx` - Firebase Auth with email/password, Google, GitHub, phone; session persistence; Firestore user profiles
- **API Client:** `frontend/src/lib/api.ts` - Centralized API client connecting to backend, with typed interfaces for curriculum, progress, users
- **API Services:** `frontend/src/lib/api/services/` - `community.ts`, `practice.ts`, `projects.ts`
- **Firebase:** `frontend/src/lib/firebase.ts` - Firebase initialization for auth and Firestore

**App Router Pages:**
- `/` - Home page (`frontend/src/app/page.tsx`)
- `/learn/` - Learning content with dynamic lesson routes
- `/learn/lesson-[id]/` - Individual lesson pages
- `/learn/prompt-engineering/`, `/learn/openclaw/`, `/learn/ai-platform-engineering/`, `/learn/devops-engineering/` - Specialized learning tracks
- `/practice/` - Practice exercises with assessments and challenges
- `/playground/` - Code playground
- `/interviews/` - Interview practice (sessions, history, system design, tools)
- `/dashboard/`, `/profile/`, `/settings/` - User management
- `/community/` - Community posts
- `/projects/` - Project showcase
- `/cms/` - Content management (lessons, assessments, grading)
- `/admin/` - Admin analytics
- `/auth/` - Auth callbacks (magic-link, email verification)
- `/signin/`, `/signup/` - Authentication pages

**Component Organization:**
- `frontend/src/components/layout/` - Header, Footer
- `frontend/src/components/auth/` - Auth-related UI
- `frontend/src/components/learning/` - Lesson display
- `frontend/src/components/practice/` - Exercise UI
- `frontend/src/components/ui/` - Shared UI primitives
- `frontend/src/components/dashboard/` - Dashboard widgets
- `frontend/src/components/interviews/` - Interview UI
- `frontend/src/components/playground/` - Code editor (Monaco)
- `frontend/src/components/admin/` - Admin UI
- `frontend/src/components/customers/` - Customer management
- `frontend/src/components/workspace/` - Workspace layout

## Key Design Patterns

**Repository Pattern (Backend):**
- Interfaces defined in `backend/internal/repository/interfaces.go`
- Two implementations: in-memory (`memory_simple.go`) and PostgreSQL (`postgres/`)
- Selection at startup via `DB_DRIVER` env var in container initialization
- `Repositories` aggregate struct enables dependency injection of the full data layer

**Service Pattern (Backend):**
- Interfaces defined in `backend/internal/service/interfaces.go`
- `Services` aggregate struct groups all business logic
- `NewServices()` factory creates services with repository and config dependencies

**Dependency Injection Container (Backend):**
- `backend/internal/container/container.go` manages all component lifecycles
- Initialization order: Validator -> Cache -> Messaging -> Repositories -> Services -> AgentPool
- Graceful shutdown in reverse order via `shutdownFuncs`
- Falls back to no-op implementations when optional infra (Redis, Kafka) is unavailable

**Agent Pattern (AI Platform):**
- `Agent` interface in `services/ai-agent-platform/pkg/types/agent.go`
- Implementations: `base.go`, `react.go`, `coding_expert.go`
- ReAct pattern: Thought -> Action -> Observation loop

**Tool Registry Pattern (AI Platform):**
- `ToolRegistry` interface in `pkg/types/tool.go`
- Thread-safe `Registry` in `internal/tools/registry.go`
- Tools registered by name with optional category grouping

**Agent Pool Pattern (Backend):**
- `backend/internal/agents/agent.go` - Multi-agent coordination
- Agent types: Executor, TestValidator, AIAnalyzer, StateManager
- `ProcessCollaborativeRequest()` routes to multiple agents for complex tasks

**Middleware Chain Pattern:**
- `Middleware` type: `func(http.Handler) http.Handler`
- `Chain()` function composes middleware in reverse order
- Used identically in both `backend/internal/middleware/middleware.go` and `services/shared/middleware/middleware.go`

**Error Hierarchy (Backend):**
- `backend/internal/errors/errors.go` - `APIError` struct with Type, Code, Message, StatusCode, Cause
- Factory functions: `NewNotFoundError()`, `NewValidationError()`, `NewUnauthorizedError()`, etc.
- `IsAPIError()` for type assertion; handlers use this to determine HTTP status codes

## Data Flow

**Typical API Request Flow:**

```
Client Request
    |
    v
Middleware Chain (RequestID -> Logging -> Recovery -> CORS -> Security -> ContentType -> Timeout -> RateLimit -> Pagination -> CSRF -> Auth)
    |
    v
Handler (Parse + Validate JSON -> Call Service Method)
    |
    v
Service (Business Logic -> Call Repository -> Optional Cache Check/Set)
    |
    v
Repository (In-Memory or PostgreSQL -> SQL via sqlx/pq)
    |
    v
Response (APIResponse{success, data, error, request_id, timestamp})
```

**Code Execution Flow:**

```
POST /api/v1/playground/execute
    |
    v
handler.go -> handlePlaygroundExecute()
    |
    v
Docker Executor (executor/docker_executor.go)
    |  (or)
    v
AI Handler -> Agent Pool -> Specialized Agent
    |
    v
Container execution -> Output capture -> Response
```

**AI Agent Platform Flow:**

```
User Query
    |
    v
Agent.Run() / Agent.Stream()
    |
    v
ReAct Loop: Thought -> Select Tool -> Execute Tool -> Observe -> Repeat
    |
    v
Tool Registry -> Tool.Execute()
    |
    v
LLM Provider (OpenAI) -> Generate response
    |
    v
AgentOutput (with Steps, ToolCalls, Metadata, TokenUsage)
```

**Frontend Auth Flow:**

```
User Login -> Firebase Auth (email/password, Google, GitHub, phone)
    |
    v
Auth Context receives Firebase User
    |
    v
api.ts -> GET /api/v1/auth/verify with Firebase ID token
    |
    v
Backend middleware verifies Firebase token -> Creates/finds backend user
    |
    v
Backend user linked to Firebase UID -> Session established
```

## Error Handling Strategy

**Backend Error Handling:**

1. **Domain Errors:** Sentinel errors in `backend/internal/errors/errors.go` (ErrNotFound, ErrValidation, etc.)

2. **APIError Struct:** Typed errors with HTTP status mapping:
   - `APIError.Type` - Error category (e.g., "NOT_FOUND", "VALIDATION_ERROR")
   - `APIError.StatusCode` - HTTP status code
   - `APIError.Cause` - Wrapped underlying error

3. **Error Propagation:** Repository -> Service -> Handler chain. Each layer can wrap errors with context.

4. **Handler Error Response:** All handlers use `writeErrorResponse()` which detects `APIError` via `IsAPIError()` and maps to appropriate HTTP status. Non-API errors default to 500.

5. **Standard Response Format:**
   ```json
   {
     "success": false,
     "error": { "type": "...", "message": "...", "details": {...} },
     "request_id": "...",
     "timestamp": "..."
   }
   ```

6. **Panic Recovery:** `middleware.Recovery` catches panics and returns 500 with request ID.

**AI Platform Error Handling:**

- `pkg/errors/errors.go` - Custom error types for tool and agent errors
- `AgentError` struct in output with Code, Message, Details
- `ToolError` struct implementing Go's `error` interface
- Errors are captured per-step in `AgentStep` and in final `AgentOutput.Error`

## Cross-Cutting Concerns

**Logging:**
- Backend: `backend/pkg/logger/` - Structured logger abstraction, configured via `LOG_LEVEL` and `LOG_FORMAT` env vars
- Used throughout via `applog.Info(ctx, msg, key, value...)` pattern (slog-style)
- Shared services: `fmt.Printf` for simple output in `services/shared/middleware/middleware.go`

**Validation:**
- Backend: `backend/pkg/validator/` wraps `go-playground/validator`
- Domain models use `validate` struct tags (e.g., `validate:"required,min=3,max=200"`)
- Handler uses `validator.ValidateJSON()` to parse and validate in one step

**Authentication:**
- Frontend: Firebase Auth (`frontend/src/contexts/auth-context.tsx`)
- Backend: Firebase ID token verification via `backend/internal/middleware/auth.go`
- API Gateway: JWT verification via `services/api-gateway/internal/auth/jwt.go`
- Backend user records linked to Firebase UIDs

**Caching:**
- `backend/internal/cache/` - `CacheManager` interface with Redis implementation
- Supports: key-value, sessions, distributed locks, rate limiting, pub/sub
- Falls back to `NoOpCacheManager` when Redis is unavailable

**Observability:**
- `backend/internal/observability/telemetry.go` - OpenTelemetry integration
- `observability/` directory at repo root - configs, dashboards, OTel collector setup
- Middleware adds request IDs (`X-Request-ID` header) for request tracing

---

*Architecture analysis: 2026-03-31*
