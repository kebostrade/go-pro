# Phase 2: Communication Patterns - Context

**Phase:** 2 - Communication Patterns  
**Status:** Research Complete - Ready for Planning  
**Date:** 2026-04-01

---

## Decisions (Locked - Research These Deeply)

### Microservices Docker Networking

The existing codebase already has a complete microservices example in `advanced-topics/06-microservices-docker/`. Template will follow these patterns:

| Pattern | Status | Rationale |
|---------|--------|-----------|
| **Docker Compose service discovery** | ✅ **SELECTED** | DNS-based service discovery via Docker networking |
| **Environment variables for config** | ✅ **SELECTED** | Already used in existing microservices |
| **API Gateway pattern** | ✅ **SELECTED** | gin-based gateway with proxy routes |
| **Health checks** | ✅ **SELECTED** | Docker healthcheck + /health endpoints |

**Key insight:** Docker Compose provides built-in DNS resolution. Services can reference each other by service name (e.g., `http://service-a:8001`) without manual service discovery infrastructure.

### WebSocket Library

| Library | Status | Rationale |
|---------|--------|-----------|
| **gorilla/websocket v1.5.3** | ✅ **SELECTED** | Already in codebase, stable, well-documented |
| gorilla/websocket (official) | ⚠️ Alternative | gorilla is the de facto standard |
| nhooyr/websocket | 🔄 Alternative | Lightweight, but gorilla is more established |

**Key insight:** gorilla/websocket is the standard. The codebase already uses it in backend, course, and advanced-topics. Template should follow existing patterns from `advanced-topics/07-websockets-realtime/examples/chat_server.go`.

### Hub Pattern for WebSocket

The hub pattern is the standard Go approach for WebSocket servers:

| Component | Purpose |
|-----------|---------|
| **Hub** | Central coordinator managing client connections |
| **Client** | Represents a single WebSocket connection |
| **Channels** | Buffered channels for registration, unregistration, broadcast |
| **Mutex** | Protects shared state (client map) |

**Key insight:** The hub pattern prevents race conditions. Uses channels for goroutine-safe communication.

### gRPC Library

| Library | Status | Rationale |
|---------|--------|-----------|
| **google.golang.org/grpc v1.72.x** | ✅ **SELECTED** | Already in backend/go.mod |
| **google.golang.org/protobuf v1.36.x** | ✅ **SELECTED** | Already in codebase |
| **protoc-gen-go** | ✅ **SELECTED** | Standard protobuf compiler plugin |

**Key insight:** gRPC ecosystem is stable. The codebase has complete examples in `advanced-topics/08-grpc-distributed/` with proto files, unary, and streaming patterns.

---

## the agent's Discretion (Research Options, Make Recommendations)

### Project Structure Template

Recommend creating three separate project templates:

```
basic/projects/microservices/           # Topic 5: Microservices
├── cmd/
│   ├── gateway/main.go                 # API Gateway entry
│   ├── service-a/main.go              # User service
│   └── service-b/main.go               # Order service
├── internal/
│   ├── gateway/                        # Gateway handlers
│   ├── service-a/                      # User service
│   │   ├── handler/
│   │   ├── repository/
│   │   └── models/
│   └── service-b/                      # Order service
├── pkg/
├── proto/                              # Shared proto definitions
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md

basic/projects/websocket-chat/           # Topic 6: WebSocket
├── cmd/
│   └── server/main.go
├── internal/
│   ├── hub/                           # Hub pattern implementation
│   ├── client/                        # Client handler
│   └── websocket/                     # WebSocket utilities
├── static/
│   └── index.html                     # Simple chat UI
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md

basic/projects/grpc-service/            # Topic 7: gRPC
├── cmd/
│   ├── server/main.go
│   └── client/main.go
├── internal/
│   └── service/                       # Service implementation
├── proto/
│   ├── user.proto                     # Proto definitions
│   └── user_grpc.pb.go               # Generated
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

### Error Handling Conventions

For microservices:
- Use HTTP status codes for REST inter-service communication
- Implement circuit breaker pattern (documented in existing ARCHITECTURE.md)
- Health endpoints returning JSON with service status

For WebSocket:
- JSON error messages with `type` field (e.g., `{"type": "error", "content": "..."}`)
- Graceful disconnection with close codes

For gRPC:
- Use gRPC status codes from `google.golang.org/grpc/codes`
- Implement error handling via status.Errorf()

### Streaming Patterns

gRPC streaming is a key learning objective:

| RPC Type | Use Case | Example |
|----------|----------|---------|
| **Unary** | Simple request-response | GetUser |
| **Server streaming** | Large responses | ListUsers (paginated) |
| **Client streaming** | Upload/batch operations | CreateUsers |
| **Bidirectional** | Real-time communication | Chat |

---

## Deferred Ideas (OUT OF SCOPE)

- **Service mesh (Istio/Linkerd)**: Too complex for learning templates
- **Consul/Etcd service discovery**: Docker DNS is sufficient for Phase 2
- **WebSocket over HTTP/2 (H2C)**: Requires TLS termination complexity
- **gRPC-gateway (REST-to-gRPC)**: Phase 5 GraphQL topic covers API translation
- **GraphQL subscriptions**: Separate Phase 5 topic

---

## Phase 2 Topics

1. **Microservices with Go and Docker** - Docker Compose, service discovery, API Gateway
2. **Real-time Applications with Go and WebSockets** - gorilla/websocket, hub pattern
3. **Distributed Systems with Go and gRPC** - protobuf, streaming RPC

---

## Deliverables Checklist

- [ ] Microservices template: 2 services + API gateway with Docker Compose
- [ ] WebSocket template: Hub pattern, chat server with browser client
- [ ] gRPC template: Unary + streaming examples with protoc
- [ ] All templates: Dockerfile, docker-compose.yml, Makefile
- [ ] All templates: Unit tests with >80% coverage
- [ ] All templates: GitHub Actions CI pipeline
- [ ] All templates: README with usage instructions

---

## Dependencies

- Requires Phase 1 CLI template (for project structure consistency)
- Requires Phase 1 REST API template (for HTTP patterns reference)
- Docker and Docker Compose must be available
- protoc must be available for gRPC template

---

## Locked Decisions Summary

| Topic | Decision | Library/Pattern |
|-------|----------|-----------------|
| Microservices | Docker Compose networking | Service discovery via DNS |
| WebSocket | gorilla/websocket | Hub pattern |
| gRPC | google.golang.org/grpc | protobuf + protoc |
