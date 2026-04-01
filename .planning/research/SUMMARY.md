# Project Research Summary

**Project:** Go Pro Learning Platform — Platform Enhancements v1.1
**Domain:** Learning Platform — Code Execution & Curriculum Integration
**Researched:** 2026-04-01
**Confidence:** HIGH

## Executive Summary

The Go Pro Learning Platform needs 4 key enhancements to deliver on its vision of letting learners study, run, and get feedback on code across all 15 advanced Go project templates.

**What type of product this is:** A code-learning platform combining curated reference code (project templates), in-browser execution (playground), and AI-powered review. Similar to Exercism but with deeper integration into production-grade templates.

**Recommended approach based on research:**
1. **Curriculum Integration (Phase 6):** Create lesson pages that reference existing project templates in `basic/projects/`. Define exercise schema that is decoupled from template structure.
2. **Code Execution (Phase 7):** Implement sandboxed execution using gVisor (runsc) for security. Support streaming output via WebSocket. Use existing executor/ infrastructure as foundation.
3. **Docker Setup (Phase 8):** Generate docker-compose.yml per topic from templates. Provide one-click "Start Environment" buttons that launch topic-specific containers.
4. **Code Review (Phase 9):** Store submissions, trigger AI analysis using existing ai-agent-platform, provide structured feedback with categories and scores.

**Key risks and mitigation:**
- **Security risk:** Arbitrary code execution → Mitigate with gVisor sandboxing, resource limits, timeouts
- **Template coupling:** Lessons breaking when templates change → Abstract via exercise schema
- **AI quality inconsistency:** Unpredictable feedback → Structured prompts, feedback templates

## Key Findings

### Recommended Stack

**Core technologies:**
- **gVisor (runsc)**: Secure container runtime for code execution — recommended over raw Docker for security
- **WebSocket**: Real-time streaming of execution output — already in stack via gorilla/websocket
- **Existing infrastructure**: Next.js 15, Go backend with chi, PostgreSQL, Redis, Monaco editor — leverage existing

**Integration points:**
- `frontend/` already has Monaco editor → extend for execution UI
- `backend/executor/` already exists → enhance for gVisor integration
- `services/ai-agent-platform/` already exists → extend for code review
- `course/` module exists → add lesson/exercise content

### Expected Features

**Must have (table stakes):**
- Code editor with syntax highlighting — Monaco exists
- Run code and see output — needs enhancement for sandbox
- Lesson pages for 15 topics — create new content
- Exercise definitions per topic — define schema
- Progress tracking per user — extend existing API
- Docker environment setup — generate docker-compose per topic

**Should have (competitive):**
- AI-powered code review — leverage existing agent
- Streaming execution output — WebSocket implementation
- Side-by-side diff view — compare solution vs submission
- Interactive exercises — fill-in-blank parsing

**Defer (v2+):**
- Video content integration
- Real-time collaborative editing
- Human code review scheduling

### Architecture Approach

**Major components:**
1. **Course Service** (backend/internal/course/) — Manages lessons, exercises, links to templates
2. **Execution Service** (backend/internal/execution/) — Sandboxed code running with gVisor
3. **Review Service** (backend/internal/review/) — Submission storage, AI agent orchestration
4. **Frontend Components** (frontend/src/components/) — Editor, output, progress UI

**Key architectural decisions:**
- Exercise schema decoupled from template structure (JSON/YAML format)
- Execution uses gVisor for security isolation
- AI review is async (queue + notify) to avoid blocking
- Docker setup generates docker-compose from template metadata

### Critical Pitfalls

1. **Arbitrary Code Execution Security** — Malicious code execution, resource abuse → Use gVisor, rate limits, timeouts
2. **Long Execution Blocking UI** — Infinite loops freeze UI → 30s timeout, streaming output, stop button
3. **Curriculum-template Coupling** — Template changes break lessons → Abstract exercise schema, version format
4. **No Progress Persistence** — Lost progress erodes trust → Auto-save to backend, visual confirmation
5. **AI Feedback Quality Inconsistency** — Random quality → Structured prompts, feedback templates, test with known code

## Implications for Roadmap

Based on research, suggested phase structure for v1.1:

### Phase 6: Curriculum Integration
**Rationale:** Foundation — everything else depends on knowing what exercises exist
**Delivers:** Lesson pages, exercise schema, progress tracking hooks
**Addresses:** Pitfall 3 (template coupling), Pitfall 4 (progress persistence)
**Avoids:** Lessons tightly coupled to template structure

### Phase 7: Code Execution
**Rationale:** Core value prop — learners need to run code in-browser
**Delivers:** Sandboxed execution API, streaming output, timeout handling
**Addresses:** Pitfall 1 (security), Pitfall 2 (blocking UI)
**Uses:** gVisor, WebSocket streaming

### Phase 8: Docker Environment Setup
**Rationale:** Bridges browser and local development — natural next after execution
**Delivers:** One-click docker-compose generation, environment status UI
**Addresses:** User need for "it works on my machine" parity

### Phase 9: Code Review System
**Rationale:** Completes the learning loop — study, run, get feedback
**Delivers:** Submission storage, AI agent integration, feedback display
**Addresses:** Pitfall 5 (AI quality), user expectation of feedback

### Phase Ordering Rationale

- Phase 6 must come first (lesson/exercise definitions needed)
- Phase 7 (execution) enables Phase 8 (Docker can test Phase 7 too)
- Phase 9 depends on Phase 7 (need execution to run tests for review)

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 7 (Code Execution)**: Complex security sandboxing — needs API research for gVisor integration
- **Phase 9 (Code Review)**: AI prompt engineering for code analysis — needs experimentation

Phases with standard patterns (skip research-phase):
- **Phase 6**: Content management — extend existing course module patterns
- **Phase 8**: Docker Compose generation — straightforward templating

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Existing infrastructure well-understood |
| Features | HIGH | Clear table stakes from competitor analysis |
| Architecture | HIGH | Well-architected brownfield, clear extension points |
| Pitfalls | MEDIUM | General web security + domain-specific |

**Overall confidence:** HIGH

### Gaps to Address

- **Gap**: Exact gVisor integration API — How to resolve: Research gVisor runsc API during Phase 7 planning
- **Gap**: AI prompt engineering quality — How to handle: Create feedback templates, test with sample code

## Sources

### Primary (HIGH confidence)
- Existing codebase — architecture patterns
- gVisor documentation — security sandboxing
- Exercism.io — code learning platform patterns

### Secondary (MEDIUM confidence)
- Go Playground model — execution patterns
- LeetCode — coding challenge features

### Tertiary (LOW confidence)
- WebContainer API — alternative execution (may not support Go)

---
*Research completed: 2026-04-01*
*Ready for roadmap: yes*
