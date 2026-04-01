# Phase 1: Foundation Patterns - Context

**Phase:** 1 - Foundation Patterns  
**Status:** Research Complete - Ready for Planning  
**Date:** 2026-04-01

---

## Decisions (Locked - Research These Deeply)

### Router Choice for REST API Template

After analyzing the Go REST API ecosystem in 2026, the recommendation is:

| Router | Status | Rationale |
|--------|--------|-----------|
| **chi v5** | ✅ **SELECTED** | Lightweight, idiomatic Go, stdlib-compatible, excellent for teaching |
| gorilla/mux | ⚠️ Maintenance | Last release 2023, in maintenance mode |
| gin | 🔄 Alternative | Full framework, used in Topic 4 (Gin Web Apps) |
| stdlib net/http | 🔄 Alternative | Too verbose for learning, no routing helpers |

**Key insight:** chi provides the best balance of simplicity and Go idioms for teaching REST API patterns. It's just `net/http` with routing sugar.

### CLI Framework Choice

| Framework | Status | Rationale |
|-----------|--------|-----------|
| **cobra** | ✅ **SELECTED** | Already used in `course/AT-02-cli-apps`, industry standard |
| urfave/cli | 🔄 Alternative | V2 is solid but cobra is more prevalent in Go ecosystem |
| stdlib flag | ❌ Avoid | Too basic for production CLIs |

### Testing Framework

| Framework | Status | Rationale |
|-----------|--------|-----------|
| **testify** | ✅ **SELECTED** | Already heavily used in backend, `assert` + `require` + `mock` |
| stdlib testing | 🔄 Base | Use as foundation, layer testify on top |
| ginkgo/gomega | 🔄 Alternative | BDD style, heavier weight |
| gock | 🔄 HTTP mocking | Add for HTTP client mocking |

### Gin Web App Patterns

The existing codebase already uses gin v1.12.0. Template will follow existing patterns from `advanced-topics/06-microservices-docker/service-a/main.go`.

---

## the agent's Discretion (Research Options, Make Recommendations)

### Project Structure Template

Recommend using the existing `basic/projects/` pattern with `cmd/` entry points and `internal/` organization:

```
basic/projects/rest-api/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
├── migrations/
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml
├── go.mod
├── go.sum
├── README.md
└── Makefile
```

### Error Handling Conventions

Recommend using the existing pattern from backend:
- Custom error types with `errors.NewNotFoundError()`, `errors.NewValidationError()`
- JSON error responses via `domain.APIResponse`
- HTTP status codes following REST conventions

### Middleware Patterns

For REST API template with chi:
- Use chi middleware package: `middleware.RequestID`, `middleware.Logger`, `middleware.Recoverer`
- Custom auth middleware following backend pattern

---

## Deferred Ideas (OUT OF SCOPE)

- **GraphQL APIs**: Phase 5 topic (gqlgen)
- **gRPC APIs**: Phase 2 topic
- **WebSocket APIs**: Phase 2 topic
- **Alternative routers** (echo, fiber): Not needed for Phase 1

---

## Phase 1 Topics

1. **RESTful APIs with Go** - Using chi router (this research)
2. **CLI Applications with Go** - Using cobra
3. **Testing and Debugging in Go** - Using testify
4. **Web Applications with Go and Gin** - Using gin framework

---

## Deliverables Checklist

- [x] Router choice: chi v5
- [x] CLI framework: cobra v1.8.0
- [x] Testing framework: testify + mock
- [x] Gin middleware patterns: from existing microservices examples
- [x] Common project structure: confirmed `cmd/` + `internal/` + `pkg/`
