---
phase: 08-docker-environment
plan: "01"
subsystem: backend
tags: [docker, api, environment-management]
dependency-graph:
  requires: []
  provides:
    - DOCK-01
    - DOCK-02
    - DOCK-03
  affects:
    - frontend (Docker control panel)
tech-stack:
  added:
    - Docker CLI via exec.CommandContext
    - JSON parsing of docker compose ps output
  patterns:
    - Service layer with timeout handling
    - HTTP handler with JSON request/response
key-files:
  created:
    - backend/internal/service/docker.go
    - backend/internal/handler/docker.go
  modified:
    - backend/internal/handler/handler.go
    - backend/cmd/server/main.go
decisions:
  - id: "DOCKER-EXEC"
    decision: "Use docker compose CLI via exec.CommandContext instead of Docker SDK"
    rationale: "Simpler implementation, existing codebase precedent (local_executor.go), sufficient for requirements"
  - id: "DOCKER-PATH"
    decision: "basePath set to '../..' to point from backend/ to repo root containing basic/projects/"
    rationale: "Docker compose commands need to run in directories containing docker-compose.yml"
  - id: "DOCKER-TIMEOUT"
    decision: "2 minute timeout for up/down, 10 second timeout for status"
    rationale: "docker compose up can take time for image pulls; status should be fast"
metrics:
  duration: "~2 minutes"
  completed: "2026-04-01T16:36:06Z"
  files-created: 2
  files-modified: 2
  lines-added: ~337
---

# Phase 08 Plan 01 Summary: Docker Environment Management API

## One-liner

Docker environment management API with POST /api/docker/up, POST /api/docker/down, and GET /api/docker/status endpoints that spawn docker compose commands for topic-based environments.

## What Was Built

Backend API endpoints for starting, stopping, and checking status of Docker environments for learning topics.

### Created Files

**backend/internal/service/docker.go**
- `DockerService` struct with `basePath` and `timeout` configuration
- `NewDockerService(basePath string)` constructor
- `StartEnvironment(ctx, topicID)` - runs `docker compose up -d --remove-orphans`
- `StopEnvironment(ctx, topicID)` - runs `docker compose down`
- `GetStatus(ctx, topicID)` - runs `docker compose ps --format json` and parses response
- `DockerStatus` and `ServiceStatus` types with JSON tags

**backend/internal/handler/docker.go**
- `DockerHandler` struct wrapping `*service.DockerService`
- `NewDockerHandler(dockerService)` constructor
- `handleDockerUp` - POST /api/docker/up with {topic_id} body
- `handleDockerDown` - POST /api/docker/down with {topic_id} body
- `handleDockerStatus` - GET /api/docker/status?topic_id=xxx
- `DockerRequest` and `DockerResponse` types

### Modified Files

**backend/internal/handler/handler.go**
- Added `dockerHandler *DockerHandler` field to `Handler` struct
- Added `SetDockerHandler(dockerHandler)` method
- Added Docker route registration in `RegisterRoutes()` when dockerHandler is set

**backend/cmd/server/main.go**
- Added Docker service initialization: `service.NewDockerService("../..")`
- Added Docker handler creation and injection via `SetDockerHandler`

## Verification

| Check | Status |
|-------|--------|
| `go build ./...` passes | ✅ |
| Handler has SetDockerHandler method | ✅ |
| Routes POST /api/docker/up registered | ✅ |
| Routes POST /api/docker/down registered | ✅ |
| Routes GET /api/docker/status registered | ✅ |
| DockerService uses exec.CommandContext with timeout | ✅ |
| basePath correctly points to repo root | ✅ |

## API Usage

### Start Docker environment for a topic
```bash
POST /api/docker/up
Content-Type: application/json

{"topic_id": "rest-api"}

# Response
{
  "success": true,
  "data": {
    "topic_id": "rest-api",
    "status": "running",
    "services": [
      {"name": "rest-api", "status": "running", "health": "healthy"}
    ],
    "ports": {"rest-api": "8080:8080"},
    "last_update": "2026-04-01T16:38:00Z"
  }
}
```

### Stop Docker environment for a topic
```bash
POST /api/docker/down
Content-Type: application/json

{"topic_id": "rest-api"}

# Response
{
  "success": true,
  "data": {
    "topic_id": "rest-api",
    "status": "stopped",
    "services": [],
    "ports": {},
    "last_update": "2026-04-01T16:40:00Z"
  }
}
```

### Get Docker environment status
```bash
GET /api/docker/status?topic_id=rest-api

# Response
{
  "success": true,
  "data": {
    "topic_id": "rest-api",
    "status": "running",
    "services": [
      {"name": "rest-api", "status": "running", "health": "healthy"}
    ],
    "ports": {"rest-api": "8080:8080"},
    "last_update": "2026-04-01T16:39:00Z"
  }
}
```

## Deviations from Plan

None - plan executed exactly as written.

## Commit

```
99ecd42 feat(08-docker-environment): add Docker environment management API
```

## Next Steps

- Plan 08-02: Create Docker control panel UI component in frontend
- Plan 08-03: Integrate Docker status polling and environment lifecycle management
