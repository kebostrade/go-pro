---
phase: advanced-topics-expansion
verified: 2026-04-01T14:30:00Z
status: gaps_found
score: 14/15 templates fully verified
gaps:
  - project: microservices
    issue: Missing Dockerfile
    status: failed
    missing:
      - "Dockerfile at basic/projects/microservices/Dockerfile"
---

# Advanced Topics Expansion - Project Templates Verification

**Milestone:** Advanced Topics Expansion
**Verified:** 2026-04-01
**Status:** gaps_found (14/15 templates passed)
**Total Projects:** 15

---

## Summary

| Category | Passed | Failed | Total |
|----------|--------|--------|-------|
| go build | 15 | 0 | 15 |
| go test | 15 | 0 | 15 |
| go vet | 15 | 0 | 15 |
| Dockerfile | 14 | 1 | 15 |

**Critical Issue:** The `microservices` project is missing a Dockerfile.

---

## Detailed Results by Template

### 1. rest-api ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 2 packages tested (handler, service) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/handler` - OK (0.005s)
- `internal/service` - OK (0.014s)

---

### 2. cli-app ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 2 packages tested (config, greeting) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/config` - OK (0.007s)
- `pkg/greeting` - OK (0.005s)

---

### 3. testing-patterns ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 3 packages tested |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/client` - OK (0.046s)
- `internal/handler` - OK (0.017s)
- `internal/service` - OK (0.028s)

---

### 4. gin-web ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (handler) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/handler` - OK (0.017s)

---

### 5. microservices ⚠️

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 3 packages tested |
| go vet | ✅ PASS | No issues |
| Dockerfile | ❌ MISSING | **Missing Dockerfile** |

**Test Results:**
- `cmd/service-a` - OK (cached)
- `cmd/service-b` - OK (cached)
- `internal/gateway` - OK (cached)

**Issue:** No Dockerfile present at `basic/projects/microservices/Dockerfile`

---

### 6. websocket-chat ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (hub) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/hub` - OK (cached)

---

### 7. grpc-service ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (service) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/service` - OK (cached)

---

### 8. kubernetes ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (controllers) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `controllers` - OK (0.035s)

---

### 9. nats-events ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | No test files (cmd packages only) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

---

### 10. serverless ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (handlers) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/handlers` - OK (0.012s)

---

### 11. ml-gorgonia ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 2 packages tested (model, tensor) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/model` - OK (cached)
- `internal/tensor` - OK (cached)

---

### 12. blockchain ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (ethereum) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/ethereum` - OK (cached)

---

### 13. iot-mqtt ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 2 packages tested (mqtt, processor) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/mqtt` - OK (cached)
- `internal/processor` - OK (cached)

---

### 14. system-design ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 3 packages tested (circuit, clean, concurrency) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `internal/circuit` - OK (cached)
- `internal/clean` - OK (cached)
- `internal/concurrency` - OK (cached)

---

### 15. graphql ✅

| Check | Status | Details |
|-------|--------|---------|
| go build | ✅ PASS | Build successful |
| go test | ✅ PASS | 1 package tested (models) |
| go vet | ✅ PASS | No issues |
| Dockerfile | ✅ EXISTS | Present |

**Test Results:**
- `pkg/models` - OK (cached)

---

## Gap Analysis

### Critical Gap: microservices Missing Dockerfile

**Project:** `basic/projects/microservices`
**Issue:** No Dockerfile present
**Impact:** Containerization/deployment not possible without Dockerfile
**Fix Required:** Add a Dockerfile to `basic/projects/microservices/`

### Recommended Dockerfile Template

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o service ./cmd/service

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service .
CMD ["./service"]
```

---

## Quality Gates Summary

| Gate | Threshold | Result |
|------|-----------|--------|
| All projects build | 15/15 | ✅ 15/15 |
| All projects pass tests | 15/15 | ✅ 15/15 |
| All projects pass vet | 15/15 | ✅ 15/15 |
| All projects have Dockerfile | 15/15 | ❌ 14/15 |

**Final Status:** PASS with 1 gap (Dockerfile missing for microservices)

---

_Verified: 2026-04-01_
_Verifier: gsd-verifier_
