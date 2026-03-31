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

<!-- New scope — 15 advanced project templates + platform enhancements -->

**Project Templates (15 topics):**
- [ ] Production-grade project template: RESTful APIs with Go
- [ ] Production-grade project template: CLI Applications with Go
- [ ] Production-grade project template: Testing and Debugging in Go
- [ ] Production-grade project template: Web Applications with Go and Gin
- [ ] Production-grade project template: Microservices with Go and Docker
- [ ] Production-grade project template: Real-time Applications with Go and WebSockets
- [ ] Production-grade project template: Distributed Systems with Go and gRPC
- [ ] Production-grade project template: Cloud-Native Applications with Go and Kubernetes
- [ ] Production-grade project template: Event-Driven Applications with Go and NATS
- [ ] Production-grade project template: Machine Learning Applications with Go and Gorgonia
- [ ] Production-grade project template: Blockchain Applications with Go and Ethereum
- [ ] Production-grade project template: IoT Applications with Go and MQTT
- [ ] Production-grade project template: Serverless Applications with Go and AWS Lambda
- [ ] Production-grade project template: GraphQL APIs with Go and gqlgen
- [ ] Production-grade project template: System Design with Golang

**Platform Enhancements:**
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
*Last updated: 2026-03-31 after initialization*
