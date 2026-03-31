# External Integrations

**Analysis Date:** 2026-03-31

## APIs & External Services

**OpenAI:**
- Purpose: LLM provider for AI Agent Platform
- SDK: `github.com/sashabaranov/go-openai` v1.20.4 (`services/ai-agent-platform/go.mod`)
- Auth: `OPENAI_API_KEY` env var
- Config: `DEFAULT_MODEL=gpt-4`, `LLM_TIMEOUT=60s`, `LLM_MAX_RETRIES=3` (`services/ai-agent-platform/.env.example`)

**Anthropic:**
- Purpose: Alternative LLM provider for AI Agent Platform
- Auth: `ANTHROPIC_API_KEY` env var
- Config: Referenced in `.env.example` but no explicit SDK dependency in `go.mod`

**Firebase:**
- Purpose: Authentication provider for the learning platform
- SDK (Backend): `firebase.google.com/go/v4` v4.18.0 (`backend/go.mod`)
- SDK (Frontend): `firebase` v12.10.0 (`frontend/package.json`)
- Auth: `FIREBASE_PROJECT_ID`, `NEXT_PUBLIC_FIREBASE_API_KEY`, `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN`, `NEXT_PUBLIC_FIREBASE_APP_ID` env vars
- Flow: Frontend gets Firebase ID token -> sends to backend -> backend verifies via Firebase Admin SDK
- Files: `backend/internal/service/auth.go` (verification), `frontend/src/lib/firebase.ts` (client init), `frontend/src/lib/api.ts` (token passing)

**Google Cloud Platform:**
- Purpose: Firebase Admin SDK dependencies, cloud storage
- SDK: `google.golang.org/api` v0.231.0, `cloud.google.com/go/storage` v1.53.0 (indirect via `backend/go.mod`)
- Used for: Firebase Auth token verification, potential GCS storage

**AWS Lambda:**
- Purpose: Serverless function support in course module
- SDK: `github.com/aws/aws-lambda-go` v1.52.0 (`course/go.mod`)

## Data Storage

**PostgreSQL:**
- Primary database for all services
- Driver: `github.com/lib/pq` v1.10.9 (backend), `github.com/jmoiron/sqlx` v1.4.0 (ORM extension)
- Connection: `DB_DSN` or individual `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` env vars
- Config: `backend/internal/config/config.go` (`DatabaseConfig` struct)
- Default: `localhost:5432`, database `gopro_dev`, user `gopro_user`
- Pool: 25 max open conns, 5 max idle conns, 5min lifetime
- Docker: `postgres:15-alpine` in all compose files

**PostgreSQL + pgvector:**
- Purpose: Vector storage for AI Agent Platform embeddings
- Image: `ankane/pgvector:latest` (`services/ai-agent-platform/docker-compose.yml`)
- Connection: `VECTOR_STORE_URL` env var, defaults to `postgresql://finagent:finagent_password@postgres:5432/finagent`
- Config: `EMBEDDING_DIMENSION=1536`, `SIMILARITY_METRIC=cosine` (`services/ai-agent-platform/.env.example`)

**Redis:**
- Purpose: Caching, session storage, rate limiting
- SDK: `github.com/go-redis/redis/v8` v8.11.5 (`backend/go.mod`)
- Connection: `REDIS_URL` env var, default `redis://localhost:6379/0`
- Docker: `redis:7-alpine` in all compose files

**Qdrant:**
- Purpose: Alternative vector database for AI Agent Platform
- Connection: `QDRANT_URL=http://localhost:6333` env var
- Docker: `qdrant/qdrant:latest` (`services/ai-agent-platform/docker-compose.yml`)

**MinIO:**
- Purpose: S3-compatible object storage (dev environment)
- Docker: `minio/minio:latest` (`docker/docker-compose.dev.yml`)
- Ports: API 9000, Console 9001
- Credentials: `minioadmin`/`minioadmin` (dev only)

## Authentication & Identity

**Firebase Authentication:**
- Provider: Google Firebase Auth
- Frontend implementation: `frontend/src/lib/firebase.ts` - Initializes Firebase app and exports auth instance
- Backend implementation: `backend/internal/service/auth.go` - Token verification via Firebase Admin SDK
- Middleware: `backend/internal/middleware/auth.go` - Extracts and validates Firebase ID tokens from requests
- Adapter: `backend/cmd/server/main.go` `firebaseAuthAdapter` struct - Bridges `service.AuthService` to `middleware.AuthService`
- JWT fallback: `github.com/golang-jwt/jwt/v5` v5.3.0 for non-Firebase token flows

**Auth Flow:**
1. Frontend authenticates user via Firebase client SDK
2. Firebase returns an ID token (JWT)
3. Frontend sends ID token in `Authorization: Bearer <token>` header to backend
4. Backend middleware (`middleware/auth.go`) extracts token
5. Backend verifies token via Firebase Admin SDK
6. User record synced to local PostgreSQL database

## Monitoring & Observability

**Distributed Tracing:**
- Jaeger: All environments
- OTLP collector enabled (`COLLECTOR_OTLP_ENABLED=true` in dev compose)
- Config: `JAEGER_ENDPOINT=http://localhost:14268/api/traces`
- Feature flags: `ENABLE_TRACING=true`

**Metrics:**
- Prometheus: Metrics collection
- Grafana: Dashboards with Redis datasource plugin
- Node Exporter: System-level metrics (production)
- Config: `ENABLE_METRICS=true`

**Logging:**
- Go `log/slog` via custom `pkg/logger` (`backend/pkg/logger/`)
- JSON format by default (`LOG_FORMAT=json`)
- Elasticsearch + Kibana for log aggregation (dev): `docker/docker-compose.dev.yml`
- Loki + Promtail for log aggregation (prod): `docker/docker-compose.prod.yml`

**Error Tracking:**
- Not detected - No Sentry, Rollbar, or similar service configured

## Message Queues & Event Streaming

**Apache Kafka:**
- Purpose: Event streaming between microservices
- SDK: `github.com/segmentio/kafka-go` v0.4.47 (`backend/go.mod`)
- Docker: `confluentinc/cp-kafka:7.4.0` (dev), `7.5.0` (services)
- Config: `KAFKA_BROKERS=kafka:9092` in services compose
- UI: Kafka UI available at `http://localhost:8083` (dev)

**RabbitMQ:**
- Purpose: Message queuing (dev environment)
- Docker: `rabbitmq:3-management-alpine` (`docker/docker-compose.dev.yml`)
- Ports: AMQP 5672, Management UI 15672

**NATS:**
- Purpose: Lightweight messaging in course module
- SDK: `github.com/nats-io/nats.go` v1.49.0 (`course/go.mod`)

## Internal Service Communication

**Microservices Architecture** (`services/docker-compose.yml`):
- API Gateway (`services/api-gateway/`) on port 8080
  - Routes to: User Service (8081), Course Service (8082), Progress Service (8083)
  - Depends on: `services/shared` module (local `replace` directive in `go.mod`)
  - JWT validation at gateway level
- User Service on port 8081 - User management, schema: `users`
- Course Service on port 8082 - Course content, schema: `courses`
- Progress Service on port 8083 - Learning progress, schema: `progress`
- All services share PostgreSQL (different schemas), Redis, and Kafka

**Backend Monolith** (`backend/`):
- Single Go binary with Clean Architecture layers
- HTTP handlers at `backend/internal/handler/`
- Repository pattern with interfaces at `backend/internal/repository/interfaces.go`
- Switchable implementations: in-memory (`memory_simple.go`) and PostgreSQL (`postgres/`)
- Dependency injection via `backend/internal/container/`

**Frontend-Backend Communication:**
- REST API over HTTP
- Frontend API client: `frontend/src/lib/api.ts`
- Base URL: `NEXT_PUBLIC_API_URL` env var (default: `http://localhost:8080`)
- Auth: Firebase ID token in Authorization header
- WebSocket: Supported via `github.com/gorilla/websocket` v1.5.3

## API Contracts

**REST API Endpoints** (Backend monolith):
- Base path: `/api/v1/`
- Health check: `GET /api/v1/health`
- Authentication: `POST /api/v1/auth/*` (Firebase token verification, user sync)
- Courses: `GET /api/v1/courses/*`
- Lessons: `GET /api/v1/lessons/*`
- Progress: `GET/POST /api/v1/progress/*`
- Admin: `/api/v1/admin/*` (admin-only endpoints)
- Playground AI: `/api/v1/playground/*` (AI code execution, conditional on AgentPool)
- Interview: `/api/v1/interview/*`

**Middleware Chain** (`backend/cmd/server/main.go`):
1. RequestID
2. Logging
3. Recovery (panic handling)
4. CORS
5. Security headers
6. Content-Type enforcement
7. Timeout (30s)
8. Rate limiting (100 req/min)
9. Pagination defaults
10. CSRF protection

## Webhooks & Callbacks

**Incoming:**
- Not detected

**Outgoing:**
- Not detected

## Environment Configuration

**Backend Required Env Vars** (`backend/internal/config/config.go`):
- `SERVER_HOST` (default: `localhost`)
- `SERVER_PORT` (default: `8080`)
- `DB_DRIVER` (default: `postgres`)
- `DB_DSN` or `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `LOG_LEVEL` (default: `info`)
- `LOG_FORMAT` (default: `json`)
- `CORS_ALLOWED_ORIGINS` (default: `http://localhost:3000`)
- `FIREBASE_PROJECT_ID` - Required for Firebase Auth
- `DEV_MODE` - Set to `true` to skip Firebase initialization

**AI Agent Platform Required Env Vars** (`services/ai-agent-platform/.env.example`):
- `OPENAI_API_KEY` - OpenAI API access
- `ANTHROPIC_API_KEY` - Anthropic API access (optional)
- `DEFAULT_LLM_PROVIDER` (default: `openai`)
- `DEFAULT_MODEL` (default: `gpt-4`)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `VECTOR_STORE_URL` - pgvector connection string
- `JWT_SECRET` - JWT signing key
- `APP_ENV` - Environment name
- `APP_PORT` (default: `8080`)
- `FRAUD_DETECTION_THRESHOLD` (default: `0.7`)
- `COMPLIANCE_MODE` (default: `strict`)

**Frontend Required Env Vars** (`frontend/.env.example`):
- `NEXT_PUBLIC_API_URL` - Backend API base URL (default: `http://localhost:8080`)
- `NEXT_PUBLIC_FIREBASE_API_KEY`
- `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN`
- `NEXT_PUBLIC_FIREBASE_PROJECT_ID`
- `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET`
- `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID`
- `NEXT_PUBLIC_FIREBASE_APP_ID`
- `NEXT_PUBLIC_APP_NAME` (default: `GO-PRO Learning Platform`)
- `NEXT_PUBLIC_ENABLE_CODE_EXECUTION` (default: `true`)
- `NEXT_PUBLIC_ENABLE_AI_ASSISTANT` (default: `false`)
- `NEXT_PUBLIC_ENABLE_ANALYTICS` (default: `false`)
- `NEXT_PUBLIC_DEBUG_MODE` (default: `false`)

**Config Pattern:**
- Backend uses `github.com/joho/godotenv` to load `.env` files, with env var fallbacks defined in `backend/internal/config/config.go`
- AI platform uses `.env` file with explicit configuration
- Frontend uses Next.js `NEXT_PUBLIC_*` convention for client-side env vars

## CI/CD Pipeline

**GitHub Actions Workflows** (`.github/workflows/`):
- `backend-ci.yml` - Backend build, lint, test pipeline
- `frontend-ci.yml` - Frontend build, lint, test pipeline
- `microservices-ci.yml` - Microservices build and test
- `security.yml` - Security scanning across all services
- `terraform-ci.yml` - Infrastructure validation and planning
- `dependency-update.yml` - Automated dependency updates

**Quality Gates:**
- `golangci-lint` must pass
- `go vet` must pass
- `gosec` security scan must pass
- `govulncheck` vulnerability check must pass
- Unit tests with race detection must pass
- Pre-commit hooks enforce formatting, linting, secret detection, license headers

## Package Dependencies at Risk

**Version Inconsistencies:**
- Backend uses Go 1.25 (`backend/go.mod`) while all other modules use Go 1.23, and course uses Go 1.24
- API Gateway uses `golang-jwt/jwt` v5.2.2 while backend uses v5.3.0
- Shared module uses `google/uuid` v1.5.0 while backend and AI platform use v1.6.0

**Indirect Dependencies:**
- MongoDB driver (`go.mongodb.org/mongo-driver/v2` v2.5.0) appears as indirect in backend but no MongoDB service is configured in any compose file

---

*Integration audit: 2026-03-31*
