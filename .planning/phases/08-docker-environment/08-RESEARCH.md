# Phase 8: Docker Environment - Research

**Researched:** 2026-04-01
**Domain:** Docker Compose template patterns, Docker CLI management, container orchestration
**Confidence:** HIGH (verified against existing codebase, established Go patterns)

---

## Summary

Phase 8 delivers one-click Docker environment setup for 15 Go learning topics. This research covers template generation patterns, CLI management via Go's `exec.Command`, Docker SDK tradeoffs, status polling approaches, and cloud vs local deployment decisions.

**Primary recommendation:** Use the locked Hybrid approach (D-01, D-05) — pre-defined compose fragments per topic category, executed via `docker compose` CLI through backend API. This balances simplicity (CLI approach), reliability (proven compose files per topic), and flexibility (fragment composition).

---

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Hybrid approach — Pre-defined templates per topic category, copied and customized
- **D-02:** Topic categories: simple (std Go), database (postgres/redis), messaging (nats/kafka), cloud (k8s/lambda)
- **D-03:** Templates stored in `frontend/src/lib/docker-templates/` as composable fragments
- **D-04:** Each topic's docker-compose.yml in `basic/projects/[topic]/` serves as reference
- **D-05:** Docker CLI — use `docker compose up/down` commands for local management
- **D-06:** Execute CLI from frontend via backend API endpoint `/api/docker`
- **D-07:** Backend API spawns `docker compose` child process with timeout
- **D-08:** API endpoints: POST /api/docker/up, POST /api/docker/down, GET /api/docker/status
- **D-09:** Hybrid approach — simple topics run locally, complex topics (K8s, Lambda) use cloud or simulation
- **D-10:** Local topics: REST API, CLI, Testing, Gin Web, WebSocket, gRPC, Microservices, GraphQL, Blockchain, IoT, ML, System Design
- **D-11:** Complex/cloud topics: Kubernetes (cloud cluster or minikube fallback), AWS Lambda (localstack or mock), NATS (local jetstream)
- **D-12:** User's Docker Desktop must be running for local topics

### the agent's Discretion
- Template fragment composition strategy
- Status polling implementation details
- Frontend Docker control panel UI

### Deferred Ideas
None

---

## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| DOCK-01 | Generate docker-compose.yml tailored to each of 15 topics | Template fragments per category, existing compose files as reference |
| DOCK-02 | Start topic-specific Docker environment with one click | Backend API with `docker compose up`, frontend trigger |
| DOCK-03 | Show environment status (running/stopped) | `docker compose ps` command, polling implementation |

---

## Research Area 1: Docker Compose Template Patterns for Go Projects

### Findings

**Existing Pattern Analysis** (from 15 docker-compose.yml files):

| Category | Topics | Common Pattern |
|----------|--------|----------------|
| **Simple (std Go)** | rest-api, cli-tools, testing-patterns, gin-web, websocket-chat, grpc-service | Single service, build context, port exposure, health check via wget |
| **Database** | postgres-redis-go, microservices | Multiple services (Go + DB), named volumes, depends_on with condition: service_healthy |
| **Messaging** | nats-events | Single service (nats:2.10-alpine), JetStream enabled via command flags, monitoring port |
| **GraphQL** | graphql | Go service + postgres, similar to database pattern |
| **Complex** | kubernetes, distributed-systems | External services or simulation |

### Template Fragment Architecture

**Recommended structure:** `frontend/src/lib/docker-templates/`

```
docker-templates/
├── fragments/
│   ├── base.yml           # version, networks
│   ├── healthcheck.yml    # Reusable health check configs
│   ├── volumes.yml        # Named volumes
│   └── networks.yml       # Network definitions
├── services/
│   ├── go-service.yml     # Standard Go service template
│   ├── postgres.yml       # PostgreSQL with health check
│   ├── redis.yml           # Redis with health check
│   └── nats.yml            # NATS with JetStream
└── topics/
    ├── simple.yml          # Topic with 1 Go service
    ├── database.yml        # Topic with Go + postgres + redis
    └── messaging.yml        # Topic with NATS
```

### Template Composition Example

```yaml
# Generated for gin-web topic
version: '3.8'

services:
  gin-web:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Key Observations from Existing Compose Files

1. **Health checks are essential** — All production-ready compose files include health checks
2. **Port collisions possible** — microservices uses ports 5433, 5434 (not standard 5432) to avoid conflicts
3. **Volume naming** — postgres-redis-go uses `postgres_data`, `redis_data` with underscore (not hyphen)
4. **Network isolation** — Multi-service setups use dedicated bridge networks
5. **Build context** — All use `build: .` with local Dockerfile

---

## Research Area 2: Docker CLI Management via Go (exec.Command Patterns)

### Findings

**Backend already uses exec.Command** (see `backend/internal/service/local_executor.go`):

```go
// Pattern from local_executor.go - lines 162-194
func (e *LocalExecutor) runCode(ctx context.Context, codeDir string, input string) (string, error) {
    args := []string{"run", filepath.Join(codeDir, "main.go")}
    cmd := exec.CommandContext(ctx, "go", args...)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    cmd.Stdin = strings.NewReader(input)

    err := cmd.Run()
    // ... error handling
}
```

### Recommended Docker Handler Pattern

```go
// backend/internal/handler/docker.go

type DockerHandler struct {
    basePath string  // base path for all projects (e.g., "basic/projects")
    timeout  time.Duration
}

type DockerRequest struct {
    TopicID string `json:"topic_id"`  // e.g., "rest-api", "microservices"
}

type DockerStatus struct {
    TopicID   string `json:"topic_id"`
    Status    string `json:"status"`     // "running", "stopped", "not_created"
    Services  []string `json:"services"`
    Ports    map[string]string `json:"ports"`  // service -> port mapping
}

// POST /api/docker/up
func (h *DockerHandler) StartEnvironment(w http.ResponseWriter, r *http.Request) {
    var req DockerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    composePath := filepath.Join(h.basePath, req.TopicID, "docker-compose.yml")
    
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
    defer cancel()

    cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "up", "-d")
    cmd.Dir = filepath.Dir(composePath)

    output, err := cmd.CombinedOutput()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to start: %s", output), http.StatusInternalServerError)
        return
    }

    // Return success with status
    status, _ := h.getStatus(ctx, req.TopicID)
    json.NewEncoder(w).Encode(status)
}

// GET /api/docker/status
func (h *DockerHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
    topicID := r.URL.Query().Get("topic_id")
    
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

    status, err := h.getStatus(ctx, topicID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(status)
}

// getStatus runs `docker compose ps --format json`
func (h *DockerHandler) getStatus(ctx context.Context, topicID string) (*DockerStatus, error) {
    composePath := filepath.Join(h.basePath, topicID, "docker-compose.yml")
    
    cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "ps", "--format", "json")
    cmd.Dir = filepath.Dir(composePath)

    output, err := cmd.Output()
    if err != nil {
        // docker compose ps returns exit 1 when no containers
        if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
            return &DockerStatus{
                TopicID: topicID,
                Status:  "stopped",
                Services: []string{},
            }, nil
        }
        return nil, err
    }

    // Parse JSON output
    var services []struct {
        Service   string `json:"Service"`
        Status    string `json:"Status"`
        Ports     string `json:"Ports"`
    }
    
    if err := json.Unmarshal(output, &services); err != nil {
        return nil, err
    }

    // Determine overall status
    allRunning := len(services) > 0
    for _, s := range services {
        if !strings.Contains(strings.ToLower(s.Status), "running") {
            allRunning = false
            break
        }
    }

    return &DockerStatus{
        TopicID:  topicID,
        Status:   map[bool]string{true: "running", false: "stopped"}[allRunning],
        Services: extractServiceNames(services),
    }, nil
}
```

### Key Implementation Details

1. **Working directory** — Set `cmd.Dir` to the project directory containing docker-compose.yml
2. **Context timeout** — 2 minutes for up/down, 10 seconds for status
3. **Exit code handling** — `docker compose ps` returns exit 1 when no containers exist
4. **JSON parsing** — Use `--format json` for reliable parsing
5. **Combined output for errors** — `cmd.CombinedOutput()` captures both stdout and stderr

---

## Research Area 3: Docker SDK vs CLI for Production Use

### Comparison

| Aspect | Docker CLI (exec.Command) | Docker SDK (dockersdk) |
|--------|---------------------------|------------------------|
| **Setup complexity** | None — CLI already installed | Requires Go module, connection handling |
| **Reliability** | HIGH — CLI is battle-tested | HIGH — official SDK |
| **Error handling** | String parsing of CLI output | Typed errors |
| **Feature access** | All CLI features | All SDK features |
| **Connection management** | N/A | Docker daemon socket |
| **Windows compatibility** | Docker Desktop CLI | Docker Desktop SDK |
| **Learning curve** | Low — familiar CLI | Medium — API patterns |

### Recommendation

**Use CLI approach** (locked by D-05) for this phase because:

1. **Existing codebase precedent** — `local_executor.go` already uses `exec.Command` for `go run`
2. **Simplicity** — No SDK connection management needed
3. **Familiarity** — Docker CLI is well-documented and widely understood
4. **Sufficient for requirements** — DOCK-01/02/03 only need up/down/status
5. **Debugging ease** — CLI output is human-readable

**SDK would be preferred for:**
- High-frequency operations (SDK is faster, no fork overhead)
- Complex orchestration (swarm, stack deploy)
- Real-time events (container logs streaming)
- Secure environments (SDK supports certificate-based auth)

---

## Research Area 4: Container Status Polling Approaches

### Findings

**Three viable approaches:**

| Approach | Mechanism | Pros | Cons |
|----------|-----------|------|------|
| **1. Polling (selected)** | Frontend polls GET /api/docker/status every 5s | Simple, works with CLI | Slight delay, extra requests |
| **2. WebSocket push** | Backend pushes status on change | Real-time, efficient | More complex, new endpoint |
| **3. SSE** | Server-Sent Events | Real-time, simpler than WS | Less flexible than WebSocket |

### Recommended Implementation

**Polling approach** with smart debouncing:

```typescript
// frontend/src/lib/docker-api.ts
export class DockerEnvironment {
  private pollInterval: number = 5000; // 5 seconds
  private pollTimer?: NodeJS.Timeout;
  private statusCallback?: (status: DockerStatus) => void;

  async start(topicId: string): Promise<DockerStatus> {
    const response = await api.post('/api/docker/up', { topic_id: topicId });
    this.startPolling(topicId);
    return response.data;
  }

  async stop(topicId: string): Promise<DockerStatus> {
    const response = await api.post('/api/docker/down', { topic_id: topicId });
    this.stopPolling();
    return response.data;
  }

  async getStatus(topicId: string): Promise<DockerStatus> {
    return (await api.get('/api/docker/status', { 
      params: { topic_id: topicId } 
    })).data;
  }

  private startPolling(topicId: string) {
    this.stopPolling(); // Clear existing
    this.pollTimer = setInterval(async () => {
      const status = await this.getStatus(topicId);
      this.statusCallback?.(status);
    }, this.pollInterval);
  }

  private stopPolling() {
    if (this.pollTimer) {
      clearInterval(this.pollTimer);
      this.pollTimer = undefined;
    }
  }

  onStatusChange(callback: (status: DockerStatus) => void) {
    this.statusCallback = callback;
  }
}
```

### Status Response Shape

```json
{
  "topic_id": "microservices",
  "status": "running",
  "services": [
    { "name": "users-db", "status": "running", "health": "healthy" },
    { "name": "orders-db", "status": "running", "health": "healthy" },
    { "name": "redis", "status": "running", "health": "healthy" },
    { "name": "service-a", "status": "running", "health": "healthy" },
    { "name": "service-b", "status": "running", "health": "healthy" },
    { "name": "api-gateway", "status": "running", "health": "unhealthy" }
  ]
}
```

---

## Research Area 5: Cloud vs Local Docker for Learning Platforms

### Analysis (from D-09, D-10, D-11)

| Topic Category | Examples | Deployment | Rationale |
|----------------|----------|------------|-----------|
| **Simple** | rest-api, cli-tools, testing | Local Docker | Minimal infra, easy local setup |
| **Database** | postgres-redis-go, microservices | Local Docker | Standard databases, well-supported |
| **Messaging** | nats-events | Local JetStream | NATS has excellent local support |
| **Cloud-native** | kubernetes, distributed-systems | Cloud/simulation | K8s/minikube overhead too high |

### Hybrid Strategy

**Local Topics (11):**
- REST API, CLI Tools, Concurrent Patterns, Error Handling
- gRPC Services, Message Queues, WebSocket
- Microservices, GraphQL, Blockchain, IoT, ML, System Design
- Testing Patterns

**Cloud/Simulation Topics (4):**
- **Kubernetes** → minikube fallback OR cloud sandbox
- **AWS Lambda** → LocalStack OR mock simulation
- **Distributed Systems** → Simulation mode (no real cluster)
- **Observability** → Can use local Prometheus/Grafana via docker-compose

### Fallback Strategy

For topics requiring cloud:

```typescript
// When user clicks "Start" on kubernetes topic
async function startKubernetesEnvironment() {
  // Check if minikube is available
  const hasMinikube = await checkCommand('minikube');
  
  if (hasMinikube) {
    // Use minikube
    await dockerCommand('minikube docker-env');
    await dockerCommand('docker compose up');
  } else {
    // Show modal with options:
    // 1. "Install minikube" (link to instructions)
    // 2. "Use cloud sandbox" (if available)
    // 3. "View simulated environment"
    showKubernetesFallbackModal();
  }
}
```

---

## Architecture Patterns

### Recommended Project Structure

```
backend/
├── internal/
│   ├── handler/
│   │   └── docker.go           # Docker API handlers
│   └── service/
│       └── docker.go           # Docker business logic
└── cmd/
    └── server/
        └── main.go             # Register /api/docker routes

frontend/src/
├── lib/
│   ├── docker-templates/        # Compose file fragments
│   │   ├── fragments/
│   │   │   ├── base.yml
│   │   │   ├── postgres.yml
│   │   │   ├── redis.yml
│   │   │   └── nats.yml
│   │   └── topics/
│   │       ├── simple.yml
│   │       ├── database.yml
│   │       ├── messaging.yml
│   │       └── cloud.yml
│   ├── docker-api.ts           # API client
│   └── docker-hooks.ts         # React hooks for Docker state
├── components/
│   └── learning/
│       └── topic-viewer.tsx    # Add Docker control panel
└── app/
    └── api/docker/             # (if needed for streaming)
```

### API Design

| Endpoint | Method | Request | Response |
|----------|--------|---------|----------|
| `/api/docker/up` | POST | `{ topic_id: string }` | `DockerStatus` |
| `/api/docker/down` | POST | `{ topic_id: string }` | `DockerStatus` |
| `/api/docker/status` | GET | `?topic_id=string` | `DockerStatus` |

### DockerStatus Schema

```go
type DockerStatus struct {
    TopicID     string           `json:"topic_id"`
    Status      string           `json:"status"`        // "running" | "stopped" | "error"
    Services    []ServiceStatus  `json:"services"`
    Ports       map[string]string `json:"ports"`        // service -> "host:container"
    Error      string           `json:"error,omitempty"`
    LastUpdate  time.Time        `json:"last_update"`
}

type ServiceStatus struct {
    Name     string `json:"name"`
    Status   string `json:"status"`
    Health   string `json:"health"`    // "healthy" | "unhealthy" | "starting"
}
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Container status detection | Parse `docker ps` text output | `docker compose ps --format json` | Structured, reliable parsing |
| Health check polling | Custom goroutine per container | Rely on compose health checks | Built into Docker, well-tested |
| Compose file generation | Generate from scratch | Compose existing + fragments | Existing files are production-ready |
| Port allocation | Hardcoded ports | Dynamic port allocation | Avoid conflicts with other services |

---

## Common Pitfalls

### Pitfall 1: Port Conflicts
**What goes wrong:** Multiple topics expose same port (8080), causing startup failures.
**How to avoid:** Use topic-specific ports in compose files (e.g., rest-api: 8080, gin-web: 8081, graphql: 8082).
**Verification:** Check `docker compose ps` for duplicate port bindings before returning success.

### Pitfall 2: Orphaned Containers
**What goes wrong:** `docker compose down` fails to remove containers from previous runs.
**How to avoid:** Use `docker compose up -d --remove-orphans` to clean up.

### Pitfall 3: Missing Working Directory
**What goes wrong:** `docker compose` runs in wrong directory, can't find docker-compose.yml.
**How to avoid:** Always set `cmd.Dir` to project directory.

### Pitfall 4: Health Check Timeout
**What goes wrong:** Service marked unhealthy even when starting.
**How to avoid:** Set reasonable `start_period` (40s for Go services with DB dependencies).

### Pitfall 5: Context Timeout Too Short
**What goes wrong:** `docker compose up` times out before containers finish starting.
**How to avoid:** Use 2-minute timeout for up/down, rely on background health checks for actual readiness.

---

## Code Examples

### Docker Handler Registration (backend/cmd/server/main.go)

```go
// Add to RegisterRoutes method
func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
    // ... existing routes ...

    // Docker management
    mux.HandleFunc("POST /api/docker/up", h.handleDockerUp)
    mux.HandleFunc("POST /api/docker/down", h.handleDockerDown)
    mux.HandleFunc("GET /api/docker/status", h.handleDockerStatus)
}
```

### React Hook for Docker Status

```typescript
// frontend/src/lib/docker-hooks.ts
import { useState, useEffect, useCallback } from 'react';

export function useDockerEnvironment(topicId: string | null) {
  const [status, setStatus] = useState<DockerStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const start = useCallback(async () => {
    if (!topicId) return;
    setLoading(true);
    try {
      const result = await dockerApi.start(topicId);
      setStatus(result);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, [topicId]);

  const stop = useCallback(async () => {
    if (!topicId) return;
    setLoading(true);
    try {
      const result = await dockerApi.stop(topicId);
      setStatus(result);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, [topicId]);

  // Poll for status when running
  useEffect(() => {
    if (!topicId || status?.status !== 'running') return;

    const interval = setInterval(async () => {
      const current = await dockerApi.getStatus(topicId);
      setStatus(current);
    }, 5000);

    return () => clearInterval(interval);
  }, [topicId, status?.status]);

  return { status, loading, error, start, stop };
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `docker-compose` (V1) | `docker compose` (V2) | Docker Desktop 2.4+ (2020) | Built-in, no install needed |
| Text parsing | `--format json` | Docker CLI 2.0+ (2020) | Reliable parsing |
| Manual health checks | Compose healthcheck | Docker Compose 2.0+ | Declarative, consistent |
| Port mapping conflicts | Dynamic port allocation | Ongoing | Better multi-topic support |

---

## Open Questions

1. **Should Docker environments persist across sessions?**
   - What we know: Current compose files use named volumes
   - What's unclear: Whether to auto-cleanup on user logout
   - Recommendation: Keep volumes for 7 days, then cleanup

2. **How to handle topic switching with running environments?**
   - What we know: Users may work on multiple topics
   - What's unclear: Auto-stop old environment or run multiple
   - Recommendation: Stop previous before starting new, warn user

3. **Resource limits for learning environments?**
   - What we know: Docker can enforce memory/CPU limits
   - What's unclear: What limits per topic
   - Recommendation: No limits initially, add if resource issues arise

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies beyond Docker CLI)

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Docker CLI | All DOCK requirements | System-installed | 20.10+ | Require Docker Desktop |
| docker compose | All DOCK requirements | Bundled with Docker | 2.0+ | — |

**Missing dependencies with no fallback:**
- Docker Desktop — Required for local topics. Users must install before use.

**Missing dependencies with fallback:**
- None identified.

---

## Sources

### Primary (HIGH confidence)
- Existing 15 docker-compose.yml files in `basic/projects/*/docker-compose.yml` — verified patterns
- `backend/internal/service/local_executor.go` — exec.Command patterns
- Docker Compose official documentation — V2 CLI reference
- Docker Compose file reference — healthcheck, depends_on patterns

### Secondary (MEDIUM confidence)
- Docker SDK Go documentation — SDK vs CLI comparison
- Docker Desktop for Mac/Windows — local development requirements

### Tertiary (LOW confidence)
- Community compose templates — various patterns (need verification)

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — existing codebase verified, Docker CLI is standard
- Architecture: HIGH — follows established handler pattern
- Pitfalls: MEDIUM — based on common Docker Compose issues

**Research date:** 2026-04-01
**Valid until:** 2026-05-01 (30 days for stable Docker patterns)
