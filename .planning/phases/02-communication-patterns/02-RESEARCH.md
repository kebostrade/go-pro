# Phase 2 Research: Communication Patterns - Microservices, WebSockets, gRPC

**Phase:** 2 - Communication Patterns  
**Topic:** All Phase 2 Topics (Microservices, WebSockets, gRPC)  
**Research Date:** 2026-04-01  
**Confidence:** HIGH

## Summary

Phase 2 covers three communication patterns in Go: Microservices with Docker, WebSockets for real-time communication, and gRPC for high-performance RPC. The existing codebase has complete examples for all three patterns, providing a strong foundation for template development.

**Primary recommendations:**
- **Microservices**: Use existing `advanced-topics/06-microservices-docker/` patterns with Docker Compose service discovery
- **WebSockets**: Use `gorilla/websocket v1.5.3` with hub pattern from existing chat server examples
- **gRPC**: Use `google.golang.org/grpc v1.72.x` with `google.golang.org/protobuf v1.36.x` following existing `advanced-topics/08-grpc-distributed/` patterns

---

## Standard Stack

### Core Libraries for Phase 2

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| gorilla/websocket | 1.5.3 | WebSocket server | Already in codebase, stable, well-documented |
| google.golang.org/grpc | 1.72.x | gRPC framework | Official Google library, in backend |
| google.golang.org/protobuf | 1.36.x | Protocol buffers | Official, required for gRPC |
| github.com/grpc-ecosystem/grpc-gateway | 2.x | REST-to-gRPC gateway | For Phase 5 integration |
| docker/compose | 3.8 | Container orchestration | Industry standard |
| postgres (image) | 15-alpine | Database per service | Already in existing microservices |
| redis (image) | 7-alpine | Caching | Already in existing microservices |

### Installation

```bash
# WebSocket template
go get github.com/gorilla/websocket@v1.5.3

# gRPC template
go get google.golang.org/grpc@v1.72.0
go get google.golang.org/protobuf@v1.36.10
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# For Docker Compose networking (no install needed - Docker handles it)
```

---

## Architecture Patterns

### Recommended Project Structures

#### Microservices Template Structure

```
basic/projects/microservices/
├── cmd/
│   ├── gateway/main.go           # API Gateway entry point
│   ├── service-a/main.go        # User service entry
│   └── service-b/main.go        # Order service entry
├── internal/
│   ├── gateway/
│   │   ├── proxy.go             # Proxy logic
│   │   ├── routes.go           # Route configuration
│   │   └── registry.go          # Service registry
│   ├── service-a/
│   │   ├── handler/
│   │   ├── repository/
│   │   ├── models.go
│   │   └── docker-entrypoint.sh
│   └── service-b/
│       ├── handler/
│       ├── repository/
│       └── docker-entrypoint.sh
├── proto/                       # Shared proto definitions (for future gRPC)
├── docker-compose.yml           # Full stack orchestration
├── Dockerfile                   # Multi-stage build
├── Makefile
├── go.mod                       # go 1.23
└── README.md
```

#### WebSocket Template Structure

```
basic/projects/websocket-chat/
├── cmd/
│   └── server/main.go           # Hub + server setup
├── internal/
│   ├── hub/
│   │   ├── hub.go              # Hub implementation
│   │   └── hub_test.go
│   ├── client/
│   │   ├── client.go           # Client connection handler
│   │   └── client_test.go
│   └── websocket/
│       └── upgrader.go         # WebSocket upgrader config
├── static/
│   └── index.html              # Simple chat UI
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod                       # go 1.23
└── README.md
```

#### gRPC Template Structure

```
basic/projects/grpc-service/
├── cmd/
│   ├── server/main.go           # gRPC server
│   └── client/main.go           # gRPC client
├── internal/
│   └── service/
│       ├── service.go           # Service implementation
│       └── service_test.go
├── proto/
│   ├── user.proto              # Proto definitions
│   ├── user.pb.go              # Generated (commit)
│   └── user_grpc.pb.go         # Generated (commit)
├── docker-compose.yml           # For local development
├── Dockerfile
├── Makefile                     # Includes protoc generation
├── go.mod                       # go 1.23
└── README.md
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| WebSocket server | Custom connection handling | gorilla/websocket | Handles HTTP upgrade, ping/pong, close handshake |
| WebSocket state management | Global variables | Hub pattern | Thread-safe via channels |
| Service discovery | Consul/Etcd | Docker DNS | Built into Docker Compose |
| Protocol serialization | JSON for internal comms | protobuf | 3-10x faster, smaller messages |
| gRPC code generation | Manual stubs | protoc + plugins | Standard tooling, fewer bugs |

---

## Common Pitfalls

### Pitfall 1: WebSocket Connection Leaks

**What goes wrong:** Connections not properly closed, goroutines leak, memory grows indefinitely.

**Why it happens:** Forgetting to remove clients from hub on disconnect, not closing WebSocket connections properly.

**How to avoid:**
```go
// Always handle close in readPump
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()  // Ensure connection is closed
    }()
    // ...
}
```

**Warning signs:** Goroutine count grows, file descriptors exhausted.

### Pitfall 2: Docker Service Network Isolation

**What goes wrong:** Services can't communicate across docker-compose services.

**Why it happens:** Docker Compose creates isolated networks by default. Services must be on same network.

**How to avoid:**
```yaml
# docker-compose.yml
services:
  service-a:
    networks:
      - microservices  # Must be same network
  service-b:
    networks:
      - microservices  # Same network here
networks:
  microservices:
    driver: bridge
```

### Pitfall 3: gRPC Connection Without Timeout

**What goes wrong:** gRPC client hangs indefinitely if server is down.

**Why it happens:** Default DialContext has no timeout.

**How to avoid:**
```go
conn, err := grpc.Dial(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithBlock(),
    grpc.WithTimeout(5*time.Second),  // Deprecated but clear
)
// Better:
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
conn, err := grpc.DialContext(ctx, "localhost:50051", ...)
```

### Pitfall 4: Protobuf Breaking Changes

**What goes wrong:** Adding fields breaks existing clients/servers.

**Why it happens:** Proto3 optional fields, missing backward compatibility.

**How to avoid:**
- Never rename fields (use deprecated注释)
- Never change field numbers
- Add new fields only as optional
- Use proto3 `reserved` for deprecated field numbers

---

## Code Examples

### Hub Pattern (WebSocket) - Source: advanced-topics/07-websockets-realtime/

```go
// Source: advanced-topics/07-websockets-realtime/examples/chat_server.go

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

### Docker Compose Service Discovery - Source: advanced-topics/06-microservices-docker/

```yaml
# Source: advanced-topics/06-microservices-docker/docker-compose.yml
services:
  service-a:
    environment:
      - SERVICE_PORT=8001
      - DB_HOST=users-db  # DNS name, not localhost
    depends_on:
      users-db:
        condition: service_healthy
    networks:
      - microservices

  service-b:
    environment:
      - USER_SERVICE_URL=http://service-a:8001  # Docker DNS
    depends_on:
      service-a:
        condition: service_healthy
    networks:
      - microservices

networks:
  microservices:
    driver: bridge
```

### gRPC Unary RPC - Source: advanced-topics/08-grpc-distributed/

```go
// Server implementation (advanced-topics/08-grpc-distributed/server/main.go)
type server struct {
    pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user := &pb.User{
        Id:    req.Id,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    return &pb.GetUserResponse{User: user}, nil
}

// Client (advanced-topics/08-grpc-distributed/examples/01-unary-rpc.go)
conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
client := pb.NewUserServiceClient(conn)
resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
```

### gRPC Streaming - Source: advanced-topics/08-grpc-distributed/

```protobuf
// Proto definition (advanced-topics/08-grpc-distributed/proto/user_service.proto)
service UserService {
    rpc ListUsers(ListUsersRequest) returns (stream User);  // Server streaming
    rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);  // Client streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);  // Bidirectional
}
```

```go
// Server streaming (advanced-topics/08-grpc-distributed/examples/02-server-stream.go)
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users := []*pb.User{/* ... */}
    for _, user := range users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| gorilla/mux for WebSocket | gorilla/websocket | 2016+ | Proper WebSocket protocol |
| Manual service URLs | Docker DNS discovery | 2017+ | Zero-config networking |
| JSON-REST internal | protobuf + gRPC | 2018+ | 3-10x performance |
| XML/SOAP | JSON + REST/gRPC | 2015+ | Simplified APIs |
| Point-to-point | API Gateway pattern | 2015+ | Centralized routing |

**Deprecated/outdated:**
- `github.com/gorilla/websocket` (old import path) — Use `github.com/gorilla/websocket`
- `github.com/golang/protobuf` — Use `google.golang.org/protobuf` (new module)
- `grpc.WithTimeout()` — Use `context.WithTimeout()` instead

---

## Open Questions

1. **Should templates include TLS?**
   - What we know: TLS adds security but complexity for beginners
   - What's unclear: Whether to show TLS setup or keep simple
   - Recommendation: Show TLS in production-ready variants, use insecure for learning

2. **Circuit breaker implementation in microservices?**
   - What we know: ARCHITECTURE.md documents circuit breaker pattern
   - What's unclear: Whether to implement or defer to library
   - Recommendation: Show pattern conceptually, recommend `sony/gobreaker`

3. **gRPC reflection for debugging?**
   - What we know: Reflection enables `grpcurl` debugging
   - What's unclear: Security implications in production
   - Recommendation: Enable reflection in development templates

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | All templates | ✓ | 1.26.1 | Go 1.23+ |
| Docker | All templates | ✓ | 24.x | — |
| Docker Compose | Microservices | ✓ | 2.x | — |
| protoc | gRPC template | ✓ (via install) | 3.x | — |
| protoc-gen-go | gRPC template | ✓ (via go install) | latest | — |
| protoc-gen-go-grpc | gRPC template | ✓ (via go install) | latest | — |

**Missing dependencies with no fallback:** None identified.

**Missing dependencies with fallback:**
- protoc plugins — Can install via `go install`

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
| Microservices | Integration | `docker-compose up && curl localhost:8080/health` | Will create |
| WebSocket | Unit + Integration | `go test ./internal/... -v` | Will create |
| gRPC | Unit | `go test ./internal/... -v` | Will create |

### Wave 0 Gaps
- [ ] `basic/projects/microservices/internal/hub/hub_test.go` — Hub pattern tests
- [ ] `basic/projects/websocket-chat/internal/service/service_test.go` — WebSocket service tests
- [ ] `basic/projects/grpc-service/internal/service/service_test.go` — gRPC service tests
- [ ] Framework install: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` — if not present

---

## Sources

### Primary (HIGH confidence)
- Existing code: `advanced-topics/06-microservices-docker/` — Docker Compose, service patterns
- Existing code: `advanced-topics/07-websockets-realtime/examples/chat_server.go` — Hub pattern
- Existing code: `advanced-topics/08-grpc-distributed/` — gRPC patterns, proto definitions
- Existing code: `backend/go.mod` — gorilla/websocket 1.5.3, grpc 1.72.0, protobuf 1.36.10
- Existing code: `course/advanced-topics/AT-06-websockets/` — WebSocket teaching examples

### Secondary (MEDIUM confidence)
- gorilla/websocket GitHub — v1.5.3 current, stable
- grpc.io official docs — Streaming patterns verified

### Tertiary (LOW confidence)
- Training data for gRPC streaming best practices — verify against official docs

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — based on existing codebase + active package repos
- Architecture: HIGH — based on existing patterns in advanced-topics
- Pitfalls: MEDIUM — based on known ecosystem evolution

**Research date:** 2026-04-01
**Valid until:** 2026-07-01 (30 days for stable ecosystem)

---

## Phase 1 Locked Decisions (Reference)

For consistency with Phase 1:

| Phase 1 Decision | Value |
|-------------------|-------|
| REST Router | chi v5 |
| CLI Framework | cobra v1.8.0 |
| Testing | testify v1.11.x |
| Web Framework | gin v1.12.0 |
| Go Version | 1.23+ |

Phase 2 templates should maintain these same standards where applicable.
