# Phase 7: Code Execution - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers in-browser Go code execution with sandbox isolation for all 15 advanced Go project topics. Users can run their code and see output in real-time with resource limits and security guards.

</domain>

<decisions>
## Implementation Decisions

### Execution Flow
- **D-01:** Run All — Single "Run All" button executes the entire project code
- **D-02:** Execute button in Practice tab of TopicViewer
- **D-03:** User sees compilation status, then output/results

### Output Display
- **D-04:** Terminal-style console output for stdout/stderr
- **D-05:** Formatted test results showing pass/fail per test case
- **D-06:** Clear distinction between compilation errors and runtime output

### Error Presentation
- **D-07:** Compilation errors shown with line numbers
- **D-08:** Runtime panics shown with stack traces
- **D-09:** Timeout shown as "Execution timed out (15s limit)"

### Topic Packages
- **D-10:** Package allowlist per topic — pre-approved packages for each exercise
- **D-11:** Blocked packages: os, net, syscall, unsafe, runtime/debug (already implemented in docker_executor.go)
- **D-12:** Topic-specific allowlists (e.g., gRPC topic allows grpc packages, NATS topic allows nats packages)
- **D-13:** Execution uses existing Docker sandbox infrastructure (docker_executor.go)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Existing Infrastructure
- `backend/internal/executor/docker_executor.go` — Docker-based execution with sandboxing
- `backend/internal/executor/IMPLEMENTATION.md` — Implementation details
- `frontend/src/components/workspace/MonacoEditor.tsx` — Monaco editor with Go support
- `frontend/src/components/learning/topic-viewer.tsx` — TopicViewer (Phase 6)
- `frontend/src/components/learning/exercise-card.tsx` — ExerciseCard (Phase 6)
- `frontend/src/lib/topics-data.ts` — 15 topic definitions with exercise data

### Backend Service Layer
- `backend/internal/service/executor.go` — ExecutorService interface

### Requirements
- `backend/internal/requirements.go` — Code execution requirements

</canonical_refs>

<codebase_context>
## Existing Code Insights

### Reusable Assets
- `MonacoEditor` component with Go syntax highlighting and autocomplete
- `DockerExecutor` with 15s timeout, 256MB memory limit, dangerous import blocking
- `LocalExecutor` for local development without Docker
- TopicViewer with Practice tab (Phase 6) — where execution button will live

### Established Patterns
- Tab-based UI (TutorialViewer, TopicViewer)
- Card-based layouts for exercises
- Backend service pattern with ExecuteRequest/ExecuteResult

### Integration Points
- TopicViewer Practice tab → execution button
- MonacoEditor for code display
- Backend executor API endpoint

</codebase_context>

<specifics>
## Specific Ideas

- Execution button in TopicViewer Practice tab
- Output panel below MonacoEditor showing results
- Per-topic package allowlist configuration in topics-data.ts
- Reuse docker_executor.go ExecuteCode function

</specifics>

<deferred>
## Deferred Ideas

None — all decisions stayed within phase scope

</deferred>

---

*Phase: 07-code-execution*
*Context gathered: 2026-04-01*
