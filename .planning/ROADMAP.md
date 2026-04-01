# Roadmap: Go Pro Platform Enhancements v1.1

## Overview

Platform enhancements to enable learners to study, run, and submit exercises for all 15 advanced Go project templates through an integrated web platform.

## Milestones

- ✅ **v1.0 Advanced Topics Expansion** - Phases 1-5 (shipped 2026-04-01)
- 🚧 **v1.1 Platform Enhancements** - Phases 6-9 (in progress)
- 📋 **v2.0** - Phases 10+ (planned)

## Phases

- [x] **Phase 6: Curriculum Integration** - Lesson pages and exercise definitions for all 15 topics (completed 2026-04-01)
- [ ] **Phase 7: Code Execution** - In-browser Go code execution with sandbox
- [ ] **Phase 8: Docker Environment** - One-click Docker setup per topic
- [ ] **Phase 9: Code Review System** - AI-powered code submission and review

---

## Phase Details

### Phase 6: Curriculum Integration

**Goal**: Users can access lesson pages and exercises for all 15 advanced Go project topics

**Depends on**: Nothing (first phase of new milestone)

**Requirements**: CURR-01, CURR-02, CURR-03, CURR-04

**Success Criteria** (what must be TRUE):
1. User can view a lesson page for each of the 15 topics (REST API, CLI, Microservices, etc.)
2. User can see structured exercise definitions with expected behavior and hints
3. User's progress (exercises completed) persists across sessions
4. User can navigate to any topic from a central course hub page

**Plans:**
2/2 plans complete
- [x] 06-02-PLAN.md — Topic pages with TutorialViewer pattern

---

### Phase 7: Code Execution

**Goal**: Users can write and execute Go code in-browser with secure sandboxing

**Depends on**: Phase 6

**Requirements**: EXEC-01, EXEC-02, EXEC-03, EXEC-04

**Success Criteria** (what must be TRUE):
1. User can write Go code in Monaco editor embedded in the platform
2. User can execute code and see output streamed in real-time
3. Code runs in an isolated sandbox (gVisor) with resource limits
4. Infinite loops and resource abuse are prevented with 30s timeout

**Plans**: TBD

---

### Phase 8: Docker Environment

**Goal**: Users can launch a complete Docker environment for each project topic with one click

**Depends on**: Phase 6 (needs curriculum context)

**Requirements**: DOCK-01, DOCK-02, DOCK-03

**Success Criteria** (what must be TRUE):
1. User can generate a docker-compose.yml tailored to the current topic
2. User can start the Docker environment with a single button click
3. User can see whether their environment is running or stopped

**Plans**: TBD

---

### Phase 9: Code Review System

**Goal**: Users can submit code exercises and receive AI-powered feedback

**Depends on**: Phase 7 (needs execution to run tests)

**Requirements**: REVIEW-01, REVIEW-02, REVIEW-03

**Success Criteria** (what must be TRUE):
1. User can submit their code exercise for review
2. AI agent analyzes the code and provides structured feedback
3. User can view history of past submissions and their review results

**Plans**: TBD

---

## Progress

| Phase | Plans | Status | Completed |
|-------|-------|--------|-----------|
| 6. Curriculum Integration | 2/2 | Complete   | 2026-04-01 |
| 7. Code Execution | TBD | Not started | - |
| 8. Docker Environment | TBD | Not started | - |
| 9. Code Review System | TBD | Not started | - |

---

## Implementation Notes

### Phase 6: Curriculum Integration
- Leverage existing `course/` module structure
- Define exercise schema in YAML/JSON, decoupled from template structure
- Extend existing progress tracking API

### Phase 7: Code Execution
- Use gVisor (runsc) for sandbox isolation
- Implement WebSocket for streaming output
- Build on existing `executor/` infrastructure

### Phase 8: Docker Environment
- Generate docker-compose.yml from template metadata
- Use existing Docker Compose infrastructure
- Add status polling for environment state

### Phase 9: Code Review System
- Integrate with existing `services/ai-agent-platform/`
- Implement async submission queue
- Create structured feedback format

---

## Revision History

| Date | Phase | Change |
|------|-------|--------|
| 2026-04-01 | 1-5 | ✅ MILESTONE COMPLETE — 15 project templates |
| 2026-04-01 | 6-9 | 🚧 MILESTONE STARTED — Platform Enhancements |
