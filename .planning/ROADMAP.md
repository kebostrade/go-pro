# GSD Roadmap — Go Pro Advanced Topics

## Overview

15 advanced Go project templates organized into 5 phases of increasing complexity.

| Phase | Name | Topics | Focus |
|-------|------|--------|-------|
| 1 | Foundation Patterns | 4 | Core production patterns |
| 2 | Communication Patterns | 3 | IPC and networking |
| 3 | Distributed & Cloud | 3 | Scale and infrastructure |
| 4 | Specialized Domains | 4 | Niche application types |
| 5 | GraphQL & Integration | 1 | API query languages |

---

## Phase 1: Foundation Patterns 🟡 Active

**Rationale:** Establish core Go production patterns before introducing complexity.

### Topics
1. RESTful APIs with Go
2. CLI Applications with Go
3. Testing and Debugging in Go
4. Web Applications with Go and Gin

### Deliverables per Topic
- [ ] Production-grade project template (`basic/projects/[topic]`)
- [ ] `go.mod` with Go 1.23+
- [ ] Unit tests with >80% coverage
- [ ] Dockerfile and `docker-compose.yml`
- [ ] GitHub Actions CI pipeline
- [ ] README with usage instructions
- [ ] Integration with course curriculum

### Phase 1 Task Breakdown

```
Phase 1: Foundation Patterns
├── Task 1: Research Go REST API ecosystem (standard library vs frameworks)
├── Task 2: Template - RESTful API (chi v5)
├── Task 3: Template - CLI (cobra v1.8.0)
├── Task 4: Template - Testing (testify v1.11.x)
└── Task 5: Template - Gin Web App (gin v1.12.0)
```

### Plans

- [ ] 01-01-PLAN.md — RESTful API template (chi v5 router)
- [ ] 01-02-PLAN.md — CLI Application template (cobra)
- [ ] 01-03-PLAN.md — Testing patterns template (testify)
- [ ] 01-04-PLAN.md — Gin Web App template (gin v1.12.0)

### Exit Criteria
- All 4 templates pass CI
- All 4 templates are runnable locally with Docker
- Course module updated with lesson pages for all 4 topics

---

## Phase 2: Communication Patterns 🔘 Active

**Rationale:** After foundation, introduce modern Go communication patterns.

### Topics
5. Microservices with Go and Docker
6. Real-time Applications with Go and WebSockets
7. Distributed Systems with Go and gRPC

### Phase 2 Task Breakdown

```
Phase 2: Communication Patterns
├── Task 6: Template - Microservices (service discovery, Docker Compose)
├── Task 7: Template - WebSocket real-time (gorilla/websocket, hub pattern)
└── Task 8: Template - gRPC (protobuf, streaming)
```

### Plans

- [ ] 02-01-PLAN.md — Microservices template (Docker Compose DNS, API Gateway)
- [ ] 02-02-PLAN.md — WebSocket template (gorilla/websocket v1.5.3, hub pattern)
- [ ] 02-03-PLAN.md — gRPC template (protobuf, streaming RPC)

### Dependencies
- Requires Phase 1 CLI and REST API templates complete

### Exit Criteria
- All 3 templates pass CI
- All 3 templates are runnable locally with Docker
- Course module updated with lesson pages for all 3 topics

---

## Phase 3: Distributed & Cloud 🔘 Active

**Rationale:** Move to cloud-native patterns and event-driven architecture.

### Topics
9. Cloud-Native Applications with Go and Kubernetes
10. Event-Driven Applications with Go and NATS
11. Serverless Applications with Go and AWS Lambda

### Phase 3 Task Breakdown

```
Phase 3: Distributed & Cloud
├── Task 9: Template - Kubernetes (K8s manifests, Helm, operators)
├── Task 10: Template - NATS (JetStream, publish/subscribe)
└── Task 11: Template - AWS Lambda (SAM, serverless.yaml)
```

### Plans

- [ ] 03-01-PLAN.md — Kubernetes template (K8s manifests, Helm chart, operator)
- [ ] 03-02-PLAN.md — NATS Events template (JetStream, publisher/subscriber)
- [ ] 03-03-PLAN.md — AWS Lambda template (SAM, Lambda URLs)

### Dependencies
- Requires Phase 2 Microservices template
- Leverages existing NATS infrastructure in `course/` module

### Exit Criteria
- All 3 templates pass CI
- All 3 templates are deployable/runnable locally
- Course module updated with lesson pages for all 3 topics

---

## Phase 4: Specialized Domains 🔘 Pending

**Rationale:** Cover specialized application domains with unique requirements.

### Topics
12. Machine Learning Applications with Go and Gorgonia
13. Blockchain Applications with Go and Ethereum
14. IoT Applications with Go and MQTT
15. System Design with Golang

### Phase 4 Task Breakdown

```
Phase 4: Specialized Domains
├── Task 12: Template - ML with Gorgonia (tensor operations, model serving)
├── Task 13: Template - Blockchain with Ethereum (smart contracts, web3)
├── Task 14: Template - IoT with MQTT (mosquitto, sensors)
└── Task 15: Template - System Design (architecture patterns, case studies)
```

### Dependencies
- Requires Phase 3 cloud infrastructure patterns

---

## Phase 5: GraphQL & Integration 🔘 Pending

**Rationale:** Final API pattern completing the API design spectrum.

### Topics
16. GraphQL APIs with Go and gqlgen

### Phase 5 Task Breakdown

```
Phase 5: GraphQL & Integration
└── Task 16: Template - GraphQL (gqlgen, schema-first, relay)
```

### Dependencies
- Requires Phase 1 REST API patterns

---

## Platform Enhancements (Cross-Cutting)

All phases include platform integration tasks:

| Enhancement | Description | Phases |
|-------------|-------------|--------|
| Curriculum Integration | Lesson pages, exercises, progress tracking | 1-5 |
| Code Execution | In-browser Go execution (Go Playground-style) | 1-5 |
| Docker Setup | One-click environment per topic | 1-5 |
| Code Review System | Submission and feedback for exercises | 1-5 |

---

## Implementation Notes

### Module Structure
```
basic/projects/
├── rest-api/           # Phase 1
├── cli-app/            # Phase 1
├── testing-patterns/   # Phase 1
├── gin-web/            # Phase 1
├── microservices/      # Phase 2
├── websocket/         # Phase 2
├── grpc/              # Phase 2
├── kubernetes/        # Phase 3
├── nats-events/       # Phase 3
├── serverless/        # Phase 3
├── ml-gorgonia/       # Phase 4
├── blockchain/        # Phase 4
├── iot-mqtt/          # Phase 4
├── system-design/     # Phase 4
└── graphql/           # Phase 5
```

### Quality Gates (per template)
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes with >80% coverage
- [ ] `golangci-lint run` passes
- [ ] `docker build` succeeds
- [ ] `docker-compose up` runs without errors
- [ ] CI pipeline green on GitHub Actions

---

## Revision History

| Date | Phase | Change |
|------|-------|--------|
| 2026-04-01 | All | Initial roadmap created |
