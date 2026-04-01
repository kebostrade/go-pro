# Architecture Research

**Domain:** Learning Platform — Code Execution & Curriculum Integration
**Researched:** 2026-04-01
**Confidence:** HIGH

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Next.js 15 Frontend                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  Monaco     │  │  Course     │  │  Progress   │          │
│  │  Editor     │  │  Pages      │  │  Dashboard   │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
│         │                │                │                   │
├─────────┴────────────────┴────────────────┴──────────────────┤
│                      Backend API (Go + chi)                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  Course     │  │  Execution  │  │  Review     │          │
│  │  Service    │  │  Service   │  │  Service    │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
│         │                │                │                   │
├─────────┴────────────────┴────────────────┴──────────────────┤
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Data Layer (PostgreSQL + Redis)          │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Execution Sandbox (gVisor/DinD)           │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| **Monaco Editor** | In-browser code editing | Already exists in frontend |
| **Course Service** | CRUD for lessons, exercises | Extend existing course module |
| **Execution Service** | Run code in sandbox | New service or extend executor/ |
| **Review Service** | Store submissions, trigger AI review | New service or extend backend/ |
| **Progress Service** | Track user advancement | Extend existing progress API |
| **AI Agent** | Analyze code, provide feedback | Extend existing ai-agent-platform/ |

## Project Structure

### Recommended Structure

```
backend/
├── cmd/server/              # Entry point
├── internal/
│   ├── course/             # Course/lesson management (existing)
│   ├── exercise/           # Exercise definitions
│   ├── execution/          # Code execution sandbox
│   ├── review/             # Submission & review
│   └── progress/           # Progress tracking (existing)
├── pkg/
│   └── sandbox/            # Sandbox implementation
└── services/
    └── ai-agent-platform/  # AI review (existing)

frontend/
├── src/
│   ├── app/
│   │   ├── courses/        # Course pages per topic
│   │   ├── playground/     # Code execution UI
│   │   └── exercises/      # Exercise pages
│   └── components/
│       ├── editor/         # Monaco wrapper
│       └── output/         # Execution output display
```

## Architectural Patterns

### Pattern 1: Service-Oriented with Clear Boundaries

**What:** Each capability is a separate service with defined interfaces
**When to use:** When different capabilities have different scaling needs
**Trade-offs:** More complex than monolith, but better isolation

### Pattern 2: Event-Driven for Async Processing

**What:** Code execution and AI review happen asynchronously
**When to use:** When operations take time (execution, AI analysis)
**Trade-offs:** Requires event infrastructure (already have Kafka/NATS)

### Pattern 3: Repository Pattern (Existing)

**What:** Data access through interfaces
**When to use:** Already in use, continue
**Trade-offs:** Good testability, flexibility

## Data Flow

### Code Execution Flow

```
[User writes code]
    ↓
[Frontend sends to Execution API]
    ↓
[Execution Service validates input]
    ↓
[Submit to sandbox (gVisor)]
    ↓
[Run code with timeout]
    ↓
[Capture output + errors]
    ↓
[Return results via WebSocket or polling]
    ↓
[Display to user]
```

### Submission Review Flow

```
[User submits code]
    ↓
[Review Service stores submission]
    ↓
[Trigger AI Agent analysis]
    ↓
[Agent runs tests + provides feedback]
    ↓
[Store review results]
    ↓
[Notify user of completion]
```

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| **Docker** | API or CLI | For building/running containers |
| **gVisor** | runsc runtime | For secure code execution |
| **AI Agent** | Internal service call | Already exists in services/ |
| **Redis** | Cache + sessions | Already in stack |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Frontend ↔ Backend | REST/WebSocket | Already established |
| Execution ↔ Sandbox | Process or RPC | Depends on sandbox choice |
| Review ↔ AI Agent | Internal service | Same host, direct call |

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|-------------------------|
| 0-100 users | Monolith fine, single executor |
| 100-1000 users | Add execution queue, multiple executors |
| 1000-10000 users | Kubernetes HPA for executors, Redis queue |
| 10000+ users | Multi-region, dedicated execution clusters |

### Scaling Priorities

1. **First bottleneck:** Execution queue backing up → Add queue + workers
2. **Second bottleneck:** AI review slow → Cache common patterns, async processing

## Anti-Patterns

### Anti-Pattern 1: Monolithic Execution

**What people do:** Try to run all code in one process
**Why it's wrong:** Security isolation, resource management, scaling
**Do this instead:** Separate sandboxed execution per request

### Anti-Pattern 2: Synchronous AI Review

**What people do:** Wait for AI response before returning submission
**Why it's wrong:** AI takes seconds to minutes, user experience suffers
**Do this instead:** Async with WebSocket notification

### Anti-Pattern 3: Tight Coupling to Docker

**What people do:** Assume Docker is always available
**Why it's wrong:** Not all environments have Docker (security restrictions)
**Do this instead:** Abstract sandbox behind interface, support multiple backends

## Sources

- Go backend architecture (existing codebase)
- Next.js 15 patterns
- Monaco editor integration
- gVisor security model

---
*Architecture research for: Platform Enhancements*
*Researched: 2026-04-01*
