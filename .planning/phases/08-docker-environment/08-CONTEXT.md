# Phase 8: Docker Environment - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers one-click Docker environment setup for each of the 15 project topics. Users can generate, launch, and monitor Docker environments tailored to their current learning topic.

</domain>

<decisions>
## Implementation Decisions

### Generation Approach
- **D-01:** Hybrid approach — Pre-defined templates per topic category, copied and customized
- **D-02:** Topic categories: simple (std Go), database (postgres/redis), messaging (nats/kafka), cloud (k8s/lambda)
- **D-03:** Templates stored in `frontend/src/lib/docker-templates/` as composable fragments
- **D-04:** Each topic's docker-compose.yml in `basic/projects/[topic]/` serves as reference

### Container Management
- **D-05:** Docker CLI — use `docker compose up/down` commands for local management
- **D-06:** Execute CLI from frontend via backend API endpoint `/api/docker`
- **D-07:** Backend API spawns `docker compose` child process with timeout
- **D-08:** API endpoints: POST /api/docker/up, POST /api/docker/down, GET /api/docker/status

### Environment Scope
- **D-09:** Hybrid approach — simple topics run locally, complex topics (K8s, Lambda) use cloud or simulation
- **D-10:** Local topics: REST API, CLI, Testing, Gin Web, WebSocket, gRPC, Microservices, GraphQL, Blockchain, IoT, ML, System Design
- **D-11:** Complex/cloud topics: Kubernetes (cloud cluster or minikube fallback), AWS Lambda (localstack or mock), NATS (local jetstream)
- **D-12:** User's Docker Desktop must be running for local topics

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Existing Docker Compose Files
- `basic/projects/rest-api/docker-compose.yml` — Simple Go service
- `basic/projects/microservices/docker-compose.yml` — Multi-service
- `basic/projects/nats-events/docker-compose.yml` — NATS JetStream
- `basic/projects/grpc-service/docker-compose.yml` — gRPC services
- `basic/projects/websocket-chat/docker-compose.yml` — Real-time with WebSocket

### Topic Categories (from REQUIREMENTS.md)
- DOCK-01: Generate docker-compose.yml per topic
- DOCK-02: Start environment with one click
- DOCK-03: Show environment status (running/stopped)

### Frontend Infrastructure
- `frontend/src/lib/topics-data.ts` — Topic definitions with id, projectPath
- `frontend/src/components/learning/topic-viewer.tsx` — Topic page component

### Backend Infrastructure
- `backend/internal/handler/` — Existing handler pattern
- `backend/cmd/server/main.go` — Route registration

</canonical_refs>

<codebase_context>
## Existing Code Insights

### Reusable Assets
- All 15 topics have existing docker-compose.yml files
- Backend handler pattern for API endpoints
- TopicViewer for Docker control button placement

### Established Patterns
- CLI execution via backend child process
- Status polling for environment state
- Topic-specific configuration

### Integration Points
- Docker controls in TopicViewer Content tab or dedicated panel
- Status indicator on topic cards in curriculum hub
- API client for Docker management endpoints

</codebase_context>

<specifics>
## Specific Ideas

- Docker control panel in TopicViewer with Start/Stop/Status
- Topic category detection from topic.id or projectPath
- Graceful fallback for cloud topics (prompt user or show mock)
- Docker Compose fragment composition per category

</specifics>

<deferred>
## Deferred Ideas

None — all decisions stayed within phase scope

</deferred>

---

*Phase: 08-docker-environment*
*Context gathered: 2026-04-01*
