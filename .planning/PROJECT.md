# Go Pro Learning Platform — Platform Enhancements

## What This Is

A Go learning platform that teaches developers through progressive, production-grade project templates. The platform provides basic-to-advanced Go tutorials, exercises, 15 advanced project templates, and an AI agent framework. This milestone enhances the platform with AI-powered mock interviews — enabling learners to practice technical interviews with dynamic AI-generated questions and personalized feedback.

Target audience: anyone wanting to master Go by building real things across diverse domains.

## Core Value

Developers master Go through progressively harder, production-quality projects — each demonstrating real patterns used across 15 distinct Go application domains, with a platform that lets them study reference code, execute it in-browser, get AI-powered feedback on their exercises, and practice technical interviews with AI.

## Current Milestone: v1.2 AI-Powered Mock Interviews

**Goal:** Upgrade existing mock interview feature to use AI/LLM for question generation and personalized feedback

**Status:** 🚧 In Progress (2026-04-02)

**Planned features:**
- AI-powered question bank with 50+ curated questions
- LLM question selection based on user profile
- AI interviewer persona with follow-up questions
- Real-time answer analysis
- Detailed AI feedback with improvement suggestions
- Progress tracking dashboard with trends

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

**Platform Enhancements (v1.1 — COMPLETE):**
- [x] **CURR-01**: Course curriculum integration for all 15 topics (lesson pages, exercises)
- [x] **CURR-02**: User can access structured exercise definitions per topic
- [x] **CURR-03**: User progress tracked per topic and per exercise (UI exists, no persistence)
- [x] **CURR-04**: User can navigate between all 15 topics from a central hub
- [x] **EXEC-01**: In-browser code execution for each project (Go Playground-style)
- [x] **EXEC-02**: User can execute code and see output (simple fetch, no streaming)
- [x] **EXEC-03**: User code runs in a secure sandbox with resource limits
- [x] **EXEC-04**: Execution supports topic-specific requirements (external packages per template)
- [x] **DOCK-01**: One-click Docker environment setup per project topic (templates exist)
- [x] **DOCK-02**: User can start environment with one click
- [x] **DOCK-03**: User can see environment status
- [x] **REVIEW-01**: Code submission and review system for learner exercises
- [x] **REVIEW-02**: AI agent analyzes submitted code and provides structured feedback
- [x] **REVIEW-03**: User can view submission history and past review feedback

**AI-Powered Mock Interviews (v1.2 — IN PROGRESS):**
- [ ] **INTW-01a**: Create database schema for interview questions
- [ ] **INTW-01b**: Store questions with tags, difficulty levels, and expected concepts
- [ ] **INTW-01c**: API endpoint to query questions by type, difficulty, and tags
- [ ] **INTW-02a**: LLM selects 3-5 questions per interview session
- [ ] **INTW-02b**: Selection considers user's past performance and target role
- [ ] **INTW-02c**: Questions have clear expected concepts for scoring
- [ ] **INTW-03a**: Generate contextual intro explaining interview type and expectations
- [ ] **INTW-03b**: Present questions in natural language with hints if needed
- [ ] **INTW-03c**: Ask follow-up questions based on user answers
- [ ] **INTW-04a**: Immediate validation of answer completeness
- [ ] **INTW-04b**: Identify mentioned concepts vs expected concepts
- [ ] **INTW-04c**: Generate instant feedback on answer quality
- [ ] **INTW-05a**: Overall score with breakdown by question
- [ ] **INTW-05b**: Strengths identified in user's answers
- [ ] **INTW-05c**: Areas for improvement with specific suggestions
- [ ] **INTW-05d**: Personalized study recommendations based on gaps
- [ ] **INTW-06a**: Store all interview sessions with scores
- [ ] **INTW-06b**: Calculate improvement trends over time
- [ ] **INTW-06c**: Show breakdown by interview type
- [ ] **INTW-06d**: Display weak areas needing more practice

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
| All 15 topics at once, as listed | User preference — build complete curriculum | ✅ Complete |
| Production-grade templates (not minimal) | Learners study real patterns, not simplified demos | ✅ Complete |
| Study + extend interaction model | Learners first understand reference code, then extend exercises | ✅ Complete |
| All 4 platform features (curriculum, in-browser execution, Docker setup, code review) | Full-featured learning platform | ✅ Complete |
| Multi-module Go layout for projects | Consistent with existing repo architecture | ✅ Complete |

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

## ✅ MILESTONE COMPLETE: Platform Enhancements v1.1

**Completed:** 2026-04-02
**Total Plans:** 11/11 (100%)
**Total Phases:** 4 (Phases 6-9)

### Deliverables Summary

| Phase | Plans | Status | Description |
|-------|-------|--------|-------------|
| Phase 6: Curriculum Integration | 2/2 | ✅ Complete | 15 topics with lesson pages and exercises |
| Phase 7: Code Execution | 2/2 | ✅ Complete | Monaco editor + /api/execute endpoint |
| Phase 8: Docker Environment | 3/3 | ✅ Complete | Docker panel UI + template registry |
| Phase 9: Code Review System | 4/4 | ✅ Complete | Submit for review + AI feedback + history |

### Key Features Delivered

1. **Curriculum Hub** (`/curriculum`) - All 15 topics organized by phase tabs
2. **Topic Pages** (`/learn/[topic]`) - Three-tab layout (Overview/Content/Practice)
3. **Code Editor** - Monaco editor with Run/Reset functionality
4. **Output Console** - Terminal-style test results display
5. **Docker Panel** - Start/Stop/Status controls with auto-polling
6. **Code Review** - Submit code for AI analysis and receive feedback
7. **Review History** - View past submissions and feedback

### Known Gaps (Tech Debt)

| ID | Description | Severity | Fix Complexity |
|----|-------------|----------|----------------|
| EXEC-02 | No streaming output (uses simple fetch) | Medium | Medium |
| DOCK-01 | Templates not integrated into user flow | Medium | Medium |
| CURR-03 | Progress persistence | Medium | Low |

### Project Locations

```
frontend/src/
├── app/
│   ├── curriculum/           # Curriculum hub page
│   └── learn/[topic]/        # Dynamic topic pages
├── components/learning/
│   ├── topic-viewer.tsx     # Main topic viewer with 3 tabs
│   ├── exercise-card.tsx    # Exercise display with code editor
│   ├── docker-panel.tsx      # Docker environment controls
│   └── review-history.tsx   # Submission history
├── components/workspace/
│   ├── code-editor.tsx      # Monaco editor wrapper
│   └── output-console.tsx    # Test results display
└── lib/
    ├── topics-data.ts        # 15 topic definitions
    ├── docker-api.ts         # Docker API client
    ├── docker-hooks.ts       # Docker environment hook
    └── docker-templates/     # Template registry

backend/internal/
├── handler/
│   ├── execute.go           # POST /api/execute
│   ├── docker.go            # Docker environment API
│   └── review.go            # Code review API
└── repository/
    └── memory_simple.go     # In-memory storage (ReviewRepository)
```

---

## ✅ MILESTONE COMPLETE: Advanced Topics Expansion

**Completed:** 2026-04-01
**Total Plans:** 15/15 (100%)
**Total Project Templates:** 15

### Deliverables Summary (v1.0)

| Phase | Templates | Status |
|-------|-----------|--------|
| Phase 1: Foundation Patterns | 4 (REST API, CLI, Testing, Gin Web) | ✅ Complete |
| Phase 2: Communication Patterns | 3 (Microservices, WebSocket, gRPC) | ✅ Complete |
| Phase 3: Distributed & Cloud | 3 (Kubernetes, NATS, Lambda) | ✅ Complete |
| Phase 4: Specialized Domains | 4 (ML, Blockchain, IoT, System Design) | ✅ Complete |
| Phase 5: GraphQL & Integration | 1 (GraphQL API) | ✅ Complete |

### Key Achievements (v1.0)

1. **15 production-grade Go project templates** created in `basic/projects/`
2. **All templates follow Clean Architecture** with proper layering
3. **All templates include**: go.mod, tests, Docker support, CI configuration
4. **All templates are independently runnable** with `go mod tidy && go run .`
5. **Comprehensive documentation** for each template

### Project Locations (v1.0)

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

*Last updated: 2026-04-02 — Milestone v1.1 Complete*
