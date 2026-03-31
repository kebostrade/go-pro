# Technology Stack

**Analysis Date:** 2026-03-31

## Languages & Runtimes

**Primary:**
- Go 1.23 - All modules except `backend/` and `course/`
- Go 1.25 - Backend module (`backend/go.mod` specifies `go 1.25.0`)
- Go 1.24 - Course module (`course/go.mod` specifies `go 1.24.0`)

**Secondary:**
- TypeScript 5.9 - Frontend (`frontend/package.json`)
- React 19.1 - Frontend UI (`frontend/package.json`)
- SQL - PostgreSQL schema and migrations
- Terraform/HCL - Infrastructure as code (`terraform/`)

## Runtimes

- Go 1.23-1.25 runtime (backend services, AI platform, all learning projects)
- Node.js 18+ target (`frontend/tsconfig.json` target ES2017)
- Bun - Frontend package manager (`frontend/bun.lock` present)

## Package Managers

- **Go Modules** - All Go modules use `go mod` with individual `go.mod` files per module
- **Bun** - Frontend package manager (bun.lock present at `frontend/bun.lock`)
- **Make** - Build orchestration across all services

## Frameworks & Libraries

### Backend (`backend/go.mod`)

**Web Framework:**
- `github.com/gin-gonic/gin` v1.12.0 - HTTP router (primary)
- `github.com/gorilla/mux` v1.8.1 - Alternative HTTP router
- `github.com/gorilla/websocket` v1.5.3 - WebSocket support

**Database:**
- `github.com/jmoiron/sqlx` v1.4.0 - SQL extension for Go
- `github.com/lib/pq` v1.10.9 - PostgreSQL driver
- `go.mongodb.org/mongo-driver/v2` v2.5.0 - MongoDB driver (indirect)

**Caching & Messaging:**
- `github.com/go-redis/redis/v8` v8.11.5 - Redis client
- `github.com/segmentio/kafka-go` v0.4.47 - Kafka client

**Authentication:**
- `firebase.google.com/go/v4` v4.18.0 - Firebase Admin SDK
- `github.com/golang-jwt/jwt/v5` v5.3.0 - JWT handling
- `golang.org/x/crypto` v0.48.0 - Bcrypt and other crypto
- `google.golang.org/api` v0.231.0 - Google API client

**Configuration:**
- `github.com/joho/godotenv` v1.5.1 - .env file loading

**Testing:**
- `github.com/stretchr/testify` v1.11.1 - Assertion/mocking library

### AI Agent Platform (`services/ai-agent-platform/go.mod`)

**LLM Integration:**
- `github.com/sashabaranov/go-openai` v1.20.4 - OpenAI API client

**Core:**
- `github.com/google/uuid` v1.6.0 - UUID generation

### API Gateway (`services/api-gateway/go.mod`)

- `github.com/golang-jwt/jwt/v5` v5.2.2 - JWT auth
- `github.com/DimaJoyti/go-pro/services/shared` v0.0.0 - Shared library (local replace directive)
- `github.com/google/uuid` v1.5.0 - UUID generation (indirect)

### Shared Library (`services/shared/go.mod`)

- `github.com/google/uuid` v1.5.0 - UUID generation

### Course Module (`course/go.mod`)

- `github.com/aws/aws-lambda-go` v1.52.0 - AWS Lambda support
- `github.com/gorilla/websocket` v1.5.3 - WebSocket
- `github.com/lib/pq` v1.11.2 - PostgreSQL driver
- `github.com/nats-io/nats.go` v1.49.0 - NATS messaging
- `github.com/klauspost/compress` v1.18.2 - Compression
- `golang.org/x/crypto` v0.46.0 - Crypto utilities

### Frontend (`frontend/package.json`)

**Core Framework:**
- `next` v15.5.12 - Next.js App Router with Turbopack
- `react` v19.1.0 - React 19
- `react-dom` v19.1.0 - React DOM renderer

**UI Components:**
- `@radix-ui/react-*` (accordion, alert-dialog, avatar, dropdown-menu, label, navigation-menu, progress, select, separator, slot, switch, tabs, toast, tooltip) - Primitives
- `radix-ui` v1.4.3 - Radix UI core
- `lucide-react` v0.544.0 - Icon library
- `class-variance-authority` v0.7.1 - Variant styling
- `clsx` v2.1.1 - Class merging
- `tailwind-merge` v3.5.0 - Tailwind class merging
- `framer-motion` v12.36.0 - Animations

**Code Editor:**
- `@monaco-editor/react` v4.7.0 - Monaco editor (VS Code editor in browser)
- `monaco-editor` v0.53.0 - Editor core

**Rich Text:**
- `@tiptap/react` v3.20.1 - TipTap rich text editor
- `@tiptap/starter-kit` v3.20.1 - Editor extensions
- `@tiptap/extension-code-block-lowlight` v3.20.1 - Code highlighting
- `@tiptap/extension-link` v3.20.1 - Link extension
- `lowlight` v3.3.0 - Syntax highlighting

**Data Visualization:**
- `recharts` v3.8.0 - Chart library

**Drag & Drop:**
- `@dnd-kit/core` v6.3.1 - Drag and drop primitives
- `@dnd-kit/sortable` v10.0.0 - Sortable utilities
- `@dnd-kit/utilities` v3.2.2 - DnD utilities

**Authentication:**
- `firebase` v12.10.0 - Firebase client SDK
- `qrcode.react` v4.2.0 - QR code generation

**Deployment:**
- `@opennextjs/cloudflare` v1.17.1 - Cloudflare deployment adapter
- `wrangler` v4.73.0 - Cloudflare Workers CLI

**Styling:**
- `tailwindcss` v4.2.1 - Utility-first CSS
- `@tailwindcss/postcss` v4.2.1 - PostCSS plugin
- `tw-animate-css` v1.4.0 - Animation utilities

**TypeScript:**
- `typescript` v5.9.3 - Type checking

## Build Tools

**Primary Makefiles:**
- `Makefile` (root) - Full platform orchestration (build, test, docker, quality)
- `backend/Makefile` - Backend-specific build/test/lint
- `services/ai-agent-platform/Makefile` - AI platform build/test/docker
- `terraform/Makefile` - Infrastructure provisioning

**Hot Reload:**
- Air (`github.com/air-verse/air`) - Go hot reload (`backend/.air.toml`)
- Turbopack - Next.js dev bundler (`next dev --turbopack`)

**Linting & Formatting:**
- `golangci-lint` - Go linting (configurable, timeout 5m)
- `gosec` - Go security scanner
- `govulncheck` - Go vulnerability checker
- `goimports` - Import management
- `gofmt` - Go formatting
- `eslint` v9 + `eslint-config-next` v16 - Frontend linting
- `pre-commit` - Git hooks framework (` .pre-commit-config.yaml`)
- `hadolint` - Dockerfile linting
- `yamllint` - YAML linting
- `markdownlint` - Markdown linting
- `commitizen` v3.6 - Commit message formatting

**Build Optimization:**
- UPX binary compression in Docker (`upx --best --lzma`)
- Multi-stage Docker builds (builder + distroless runtime)
- Multi-platform builds (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)

## Infrastructure

**Container Images:**
- Backend: `golang:1.25-alpine` (build) + `gcr.io/distroless/static-debian12:nonroot` (runtime) (`backend/Dockerfile`)
- AI Platform: `golang:1.22-alpine` (build) + `alpine:latest` (runtime) (`services/ai-agent-platform/Dockerfile`)
- Frontend: Custom Dockerfile at `docker/frontend/Dockerfile`

**Docker Compose Environments:**
- Dev: `docker/docker-compose.dev.yml` - 15 services (backend, frontend, PostgreSQL, Redis, Kafka, Elasticsearch, Kibana, MinIO, RabbitMQ, Prometheus, Grafana, Jaeger, Mailhog, Adminer, Redis Commander)
- Prod: `docker/docker-compose.prod.yml` - Backend, Redis, PostgreSQL, Nginx, Prometheus, Node Exporter, Loki, Promtail
- Services: `services/docker-compose.yml` - Microservices stack (API gateway, user-service, course-service, progress-service, Kafka, PostgreSQL, Redis, Jaeger, Prometheus, Grafana)
- AI Platform: `services/ai-agent-platform/docker-compose.yml` - PostgreSQL+pgvector, Redis, Qdrant, Jaeger, Prometheus, Grafana

**Container Images Used:**
- `postgres:15-alpine` - Primary database
- `redis:7-alpine` - Caching and sessions
- `ankane/pgvector:latest` - PostgreSQL with vector extension (AI platform)
- `qdrant/qdrant:latest` - Vector database (AI platform)
- `confluentinc/cp-kafka:7.4.0` / `7.5.0` - Kafka
- `confluentinc/cp-zookeeper:7.4.0` / `7.5.0` - Zookeeper
- `prom/prometheus` - Metrics
- `grafana/grafana` - Dashboards
- `jaegertracing/all-in-one` - Distributed tracing
- `docker.elastic.co/elasticsearch/elasticsearch:8.11.0` - Log aggregation
- `docker.elastic.co/kibana/kibana:8.11.0` - Log visualization
- `minio/minio` - S3-compatible object storage
- `rabbitmq:3-management-alpine` - Message queue
- `mailhog/mailhog` - Email testing
- `adminer:4.8.1` - Database management UI
- `nginx:alpine` - Reverse proxy (production)
- `grafana/loki:2.9.0` + `grafana/promtail:2.9.0` - Log shipping (production)
- `prom/node-exporter` - System metrics (production)

**Infrastructure as Code:**
- Terraform - AWS infrastructure (`terraform/`) with environment-specific tfvars (dev, staging, production)
- Target: AWS EKS, RDS, ElastiCache, S3, DynamoDB, CloudWatch

**CI/CD:**
- GitHub Actions - `.github/workflows/` containing:
  - `backend-ci.yml` - Backend CI pipeline
  - `frontend-ci.yml` - Frontend CI pipeline
  - `microservices-ci.yml` - Microservices CI
  - `security.yml` - Security scanning
  - `terraform-ci.yml` - Infrastructure CI
  - `dependency-update.yml` - Dependency updates
- Docker registry: `ghcr.io` (GitHub Container Registry)

**Deployment Targets:**
- Cloudflare Workers/Edge - Frontend (`@opennextjs/cloudflare`)
- Firebase Hosting - Alternative frontend deployment (`firebase deploy`)
- Docker containers - Backend and services
- AWS (EKS, RDS) - Production infrastructure via Terraform

## Development Tools

**Auto-installed via `make install-tools`:**
- `golangci-lint` - Multi-linter aggregation
- `gosec` - Security vulnerability scanner
- `air` - Hot reload for Go
- `goimports` - Import management
- `govulncheck` - Known vulnerability detection
- `pre-commit` - Git hooks management

**Pre-commit Hooks** (`.pre-commit-config.yaml`):
- `pre-commit-hooks` v4.4.0 - Trailing whitespace, end-of-file, YAML/JSON/TOML checks, merge conflict detection, large file checks
- `pre-commit-golang` v0.5.1 - go-fmt, go-imports, go-vet, go-mod-tidy, go-unit-tests, go-build, golangci-lint
- `detect-secrets` v1.4.0 - Secret detection with baseline
- `hadolint` v2.12.0 - Dockerfile linting
- `yamllint` v1.32.0 - YAML linting
- `markdownlint-cli` v0.35.0 - Markdown linting
- `commitizen` v3.6.0 - Commit message formatting
- `insert-license` v1.5.1 - License header enforcement

**Monitoring & Observability Stack:**
- Prometheus - Metrics collection
- Grafana - Dashboarding
- Jaeger - Distributed tracing
- Elasticsearch + Kibana - Log aggregation and search (dev)
- Loki + Promtail - Log aggregation (production)

---

*Stack analysis: 2026-03-31*
