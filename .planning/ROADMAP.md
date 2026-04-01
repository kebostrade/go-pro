# GSD Roadmap — Go Pro Advanced Topics

## Overview

15 advanced Go project templates organized into 5 phases of increasing complexity.

| Phase | Name | Topics | Focus | Status |
|-------|------|--------|-------|--------|
| 1 | Foundation Patterns | 4 | Core production patterns | ✅ Complete |
| 2 | Communication Patterns | 3 | IPC and networking | ✅ Complete |
| 3 | Distributed & Cloud | 3 | Scale and infrastructure | ✅ Complete |
| 4 | Specialized Domains | 4 | Niche application types | ✅ Complete |
| 5 | GraphQL & Integration | 1 | API query languages | ✅ Complete |

---

## Phase 1: Foundation Patterns ✅ Complete

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

## Phase 2: Communication Patterns ✅ Complete

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

## Phase 3: Distributed & Cloud ✅ Complete

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

## Phase 4: Specialized Domains ✅ Complete

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

### Plans

- [ ] 04-01-PLAN.md — ML with Gorgonia template (tensor ops, ONNX inference)
- [ ] 04-02-PLAN.md — Blockchain with Ethereum template (smart contracts, wallet)
- [ ] 04-03-PLAN.md — IoT with MQTT template (device, gateway, broker)
- [ ] 04-04-PLAN.md — System Design template (clean architecture, patterns)

### Dependencies
- Requires Phase 3 cloud infrastructure patterns

### Exit Criteria
- All 4 templates pass CI
- All 4 templates are runnable locally with Docker
- Course module updated with lesson pages for all 4 topics

---

## Phase 5: GraphQL & Integration ✅ Complete

**Rationale:** Final API pattern completing the API design spectrum.

### Topics
16. GraphQL APIs with Go and gqlgen

### Phase 5 Task Breakdown

```
Phase 5: GraphQL & Integration
└── Task 16: Template - GraphQL (gqlgen, schema-first, relay) ✅ COMPLETE
```

### Plans

- [x] 05-01-PLAN.md — GraphQL API template (gqlgen v0.17+, chi v5, WebSocket subscriptions)

### Dependencies
- Requires Phase 1 REST API patterns

---

## 🎉 MILESTONE COMPLETE: Advanced Topics Expansion (2026-04-01)

**Total Progress:** 15/15 plans (100%)

| Metric | Count |
|--------|-------|
| Total Phases | 5 |
| Total Plans | 15 |
| Total Templates | 15 |
| Templates in basic/projects/ | 15 |
| Lines of Documentation | 1000+ |

### Template Inventory

| # | Template | Phase | Key Tech |
|---|----------|-------|----------|
| 1 | rest-api | 1 | chi v5, clean architecture |
| 2 | cli-app | 1 | cobra v1.8.0 |
| 3 | testing-patterns | 1 | testify |
| 4 | gin-web | 1 | gin v1.12 |
| 5 | microservices | 2 | Docker Compose, chi |
| 6 | websocket-chat | 2 | gorilla/websocket |
| 7 | grpc-service | 2 | protobuf, grpc |
| 8 | kubernetes | 3 | K8s, Helm |
| 9 | nats-events | 3 | JetStream |
| 10 | serverless | 3 | AWS Lambda, SAM |
| 11 | ml-gorgonia | 4 | gonum |
| 12 | blockchain | 4 | go-ethereum |
| 13 | iot-mqtt | 4 | paho.mqtt |
| 14 | system-design | 4 | clean architecture |
| 15 | graphql | 5 | gqlgen v0.17+ |

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
| 2026-04-01 | All | ✅ MILESTONE COMPLETE — 15/15 templates created |
