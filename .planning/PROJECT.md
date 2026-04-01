# Go Pro Learning Platform — Advanced Topics Expansion

## What This Is

A Go learning platform that teaches developers through progressive, production-grade project templates. The platform currently provides basic-to-intermediate Go tutorials, exercises, and an AI agent framework. This initiative expands it with 15 advanced Go project templates (REST APIs, microservices, blockchain, Kubernetes, etc.) and enhances the platform to serve, execute, and review code across all topics.

Target audience: anyone wanting to master Go by building real things across diverse domains.

## Core Value

Developers master Go through progressively harder, production-quality projects — each demonstrating real patterns used across 15 distinct Go application domains, with a platform that lets them study, run, and get feedback on their code.

## Requirements

### Validated

<!-- Existing capabilities confirmed by codebase -->

- ✓ Multi-module Go monorepo with independent modules (basic/, backend/, services/, course/) — existing
- ✓ Clean Architecture backend with repository pattern (in-memory + PostgreSQL) — existing
- ✓ REST API with CRUD for courses, lessons, exercises, progress, users — existing
- ✓ Firebase Authentication (email/password, Google, GitHub, phone) — existing
- ✓ Next.js 15 frontend with App Router, course pages, playground, interviews — existing
- ✓ Docker-based code execution sandbox (executor/) — existing
- ✓ Monaco editor in-browser code editing — existing
- ✓ AI Agent Platform with ReAct pattern, tool system, LLM providers — existing
- ✓ Multi-agent coordination (Executor, TestValidator, AIAnalyzer, StateManager) — existing
- ✓ Docker Compose dev environment (15 services) — existing
- ✓ API Gateway with JWT auth — existing
- ✓ Middleware chain (CORS, rate limiting, security, logging, metrics) — existing
- ✓ Progress tracking with streaks, assessments, submissions — existing
- ✓ CMS for lesson/assessment management — existing
- ✓ Infrastructure as Code (Terraform for AWS) — existing
- ✓ CI/CD with GitHub Actions — existing

### Active

<!-- All 15 advanced project templates are now complete -->

**Project Templates (15 topics) — ALL COMPLETE:**
- [x] ✅ Production-grade project template: RESTful APIs with Go (chi v5)
- [x] ✅ Production-grade project template: CLI Applications with Go (cobra v1.8.0)
- [x] ✅ Production-grade project template: Testing and Debugging in Go (testify)
- [x] ✅ Production-grade project template: Web Applications with Go and Gin (gin v1.12)
- [x] ✅ Production-grade project template: Microservices with Go and Docker (Docker Compose)
- [x] ✅ Production-grade project template: Real-time Applications with Go and WebSockets (gorilla/websocket)
- [x] ✅ Production-grade project template: Distributed Systems with Go and gRPC (protobuf)
- [x] ✅ Production-grade project template: Cloud-Native Applications with Go and Kubernetes (K8s/Helm)
- [x] ✅ Production-grade project template: Event-Driven Applications with Go and NATS (JetStream)
- [x] ✅ Production-grade project template: Machine Learning Applications with Go (gonum)
- [x] ✅ Production-grade project template: Blockchain Applications with Go and Ethereum (go-ethereum)
- [x] ✅ Production-grade project template: IoT Applications with Go and MQTT (paho.mqtt)
- [x] ✅ Production-grade project template: Serverless Applications with Go and AWS Lambda (SAM)
- [x] ✅ Production-grade project template: GraphQL APIs with Go and gqlgen (gqlgen v0.17+)
- [x] ✅ Production-grade project template: System Design with Golang (clean architecture)

**Platform Enhancements (Future Phases):**
- [ ] Course curriculum integration for all 15 new topics (lesson pages, exercises, progress tracking)
- [ ] In-browser code execution for each project (Go Playground-style)
- [ ] One-click Docker environment setup per project topic
- [ ] Code submission and review system for learner exercises

### Out of Scope

- Video content creation — platform supports text/code only, no media pipeline
- Live instructor features (chat, video calls) — beyond current scope
- Payment/billing integration — not needed for current phase
- Mobile native app — web-first platform
- Multi-language support (non-English) — English only for now

## Context

**Existing codebase state:**
- Brownfield project with substantial foundation: Clean Architecture backend, Next.js 15 frontend, AI agent platform
- Backend: Go 1.25 with gin/gorilla, PostgreSQL, Redis, Kafka, Firebase Auth
- Frontend: React 19, Next.js 15, Radix UI, Tailwind 4, Monaco editor, TipTap rich text
- Dev environment: Docker Compose with 15 services (PostgreSQL, Redis, Kafka, Elasticsearch, etc.)
- CI/CD: GitHub Actions with backend, frontend, microservices, security, and Terraform pipelines
- The `course/` module already has dependencies on NATS, WebSocket, AWS Lambda — suggesting some infrastructure exists for these topics

**15 topics as listed (in order):**
1. System Design with Golang
2. RESTful APIs with Go
3. CLI Applications with Go
4. Testing and Debugging in Go
5. Web Applications with Go and Gin
6. Microservices with Go and Docker
7. Real-time Applications with Go and WebSockets
8. Distributed Systems with Go and gRPC
9. Cloud-Native Applications with Go and Kubernetes
10. Event-Driven Applications with Go and NATS
11. Machine Learning Applications with Go and Gorgonia
12. Blockchain Applications with Go and Ethereum
13. IoT Applications with Go and MQTT
14. Serverless Applications with Go and AWS Lambda
15. GraphQL APIs with Go and gqlgen

**Learner interaction model:** Study + extend — learners study the reference production code, understand patterns, then extend exercises within each project.

## Constraints

- **Tech Stack**: Go 1.23+ for all project templates; existing backend/frontend stack for platform
- **Module Independence**: Each project template must be its own Go module with its own `go.mod`
- **Production Quality**: Each template includes tests, Docker setup, CI config, proper error handling, documentation — not toy examples
- **Platform Integration**: All 15 topics must integrate with existing course curriculum, progress tracking, and authentication system

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| All 15 topics at once, as listed | User preference — build complete curriculum | — Pending |
| Production-grade templates (not minimal) | Learners study real patterns, not simplified demos | — Pending |
| Study + extend interaction model | Learners first understand reference code, then extend exercises | — Pending |
| All 4 platform features (curriculum, in-browser execution, Docker setup, code review) | Full-featured learning platform | — Pending |
| Multi-module Go layout for projects | Consistent with existing repo architecture | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---

## ✅ MILESTONE COMPLETE: Advanced Topics Expansion

**Completed:** 2026-04-01
**Total Plans:** 15/15 (100%)
**Total Project Templates:** 15

### Deliverables Summary

| Phase | Templates | Status |
|-------|-----------|--------|
| Phase 1: Foundation Patterns | 4 (REST API, CLI, Testing, Gin Web) | ✅ Complete |
| Phase 2: Communication Patterns | 3 (Microservices, WebSocket, gRPC) | ✅ Complete |
| Phase 3: Distributed & Cloud | 3 (Kubernetes, NATS, Lambda) | ✅ Complete |
| Phase 4: Specialized Domains | 4 (ML, Blockchain, IoT, System Design) | ✅ Complete |
| Phase 5: GraphQL & Integration | 1 (GraphQL API) | ✅ Complete |

### Key Achievements

1. **15 production-grade Go project templates** created in `basic/projects/`
2. **All templates follow Clean Architecture** with proper layering
3. **All templates include**: go.mod, tests, Docker support, CI configuration
4. **All templates are independently runnable** with `go mod tidy && go run .`
5. **Comprehensive documentation** for each template

### Project Locations

```
basic/projects/
├── rest-api/           # Phase 1 - chi v5 REST API
├── cli-app/            # Phase 1 - cobra CLI
├── testing-patterns/   # Phase 1 - testify testing
├── gin-web/            # Phase 1 - gin v1.12 web
├── microservices/      # Phase 2 - Docker Compose
├── websocket-chat/      # Phase 2 - gorilla/websocket
├── grpc-service/       # Phase 2 - protobuf/gRPC
├── kubernetes/         # Phase 3 - K8s/Helm
├── nats-events/        # Phase 3 - JetStream
├── serverless/         # Phase 3 - AWS Lambda
├── ml-gorgonia/        # Phase 4 - ML with gonum
├── blockchain/          # Phase 4 - go-ethereum
├── iot-mqtt/           # Phase 4 - paho.mqtt
├── system-design/      # Phase 4 - patterns
└── graphql/            # Phase 5 - gqlgen
```

---
*Last updated: 2026-04-01 — MILESTONE COMPLETE*
