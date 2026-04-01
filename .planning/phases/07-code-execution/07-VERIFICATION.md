---
phase: 07-code-execution
verified: 2026-04-01T19:30:00Z
status: gaps_found
score: 3/4 must-haves verified
gaps:
  - truth: "User sees real-time streaming output during execution"
    status: partial
    reason: "CodeEditor uses simple fetch() that waits for complete response. No streaming via SSE or WebSocket."
    artifacts:
      - path: frontend/src/components/workspace/CodeEditor.tsx
        issue: "handleExecute uses await fetch() without streaming. Output only appears after execution completes."
    missing:
      - "Server-Sent Events (SSE) or WebSocket for real-time streaming output"
      - "Frontend handler for streaming chunks"
  - truth: "User code runs with enforced CPU and memory limits"
    status: failed
    reason: "DockerExecutor struct has memory/cpuLimit fields but runContainer() never applies them to docker run command."
    artifacts:
      - path: backend/internal/executor/docker_executor.go
        issue: "args array in runContainer (line 454-461) does not include --memory or --cpus flags despite fields existing"
    missing:
      - "Add --memory flag to docker run args"
      - "Add --cpus flag to docker run args"
---

# Phase 07: Code Execution Verification Report

**Phase Goal:** In-browser code execution with sandbox isolation
**Verified:** 2026-04-01
**Status:** gaps_found
**Score:** 3/4 must-haves verified

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | User sees Monaco editor in Practice tab | ✓ VERIFIED | `topic-viewer.tsx` line 271 renders `<CodeEditor>` which wraps `MonacoEditor` |
| 2 | User can type Go code in the editor | ✓ VERIFIED | `CodeEditor.tsx` line 118: `onChange={(v) => setCode(v \|\| '')}` |
| 3 | User can click Run button to execute code | ✓ VERIFIED | `CodeEditor.tsx` lines 103-110: Run button with `handleExecute` |
| 4 | User sees terminal-style output with results | ✓ VERIFIED | `OutputConsole.tsx` renders terminal-style output with pass/fail icons |
| 5 | User sees real-time streaming output | ✗ PARTIAL | Simple fetch without streaming - output only after completion |
| 6 | User code runs in sandbox with resource limits | ✗ FAILED | Memory/CPU defined but not enforced in docker run |

**Score:** 4/6 truths verified

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `frontend/src/components/workspace/CodeEditor.tsx` | Monaco editor with Run button | ✓ VERIFIED | 131 lines, wraps MonacoEditor, has handleExecute and handleReset |
| `frontend/src/components/workspace/OutputConsole.tsx` | Terminal-style output display | ✓ VERIFIED | 101 lines, dark theme, pass/fail icons, expected vs actual |
| `frontend/src/components/learning/topic-viewer.tsx` | Practice tab with CodeEditor | ✓ VERIFIED | Lines 271-280 render CodeEditor with topicId |
| `backend/internal/handler/execute.go` | /api/execute endpoint | ✓ VERIFIED | Lines 25-98 implement ExecuteCode handler |
| `backend/internal/executor/docker_executor.go` | Sandbox with limits | ⚠️ PARTIAL | Fields exist but memory/cpu not enforced in runContainer |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|----|--------|---------|
| topic-viewer.tsx | CodeEditor.tsx | PracticeTab renders CodeEditor | ✓ WIRED | Line 271: `<CodeEditor topicId={topic.id} ...>` |
| CodeEditor.tsx | MonacoEditor.tsx | import and render | ✓ WIRED | Line 4: import, line 116: `<MonacoEditor>` |
| CodeEditor.tsx | /api/execute | fetch POST | ✓ WIRED | Line 60: `fetch('/api/execute', ...)` |
| CodeEditor.tsx | OutputConsole.tsx | import and render | ✓ WIRED | Line 5: import, line 127: `<OutputConsole>` |

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|--------------|--------|-------------------|--------|
| CodeEditor.tsx | code, output, isRunning | useState | N/A (UI state) | ✓ VERIFIED |
| OutputConsole.tsx | result, error | Props from parent | N/A (display only) | ✓ VERIFIED |

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| Backend builds without errors | `cd backend && go build ./...` | Build succeeded | ✓ PASS |
| Frontend builds without errors | `cd frontend && npx tsc --noEmit 2>&1 | head -20` | TypeScript errors only | ⚠️ SKIP (needs full deps install) |

### Requirements Coverage

| Requirement | Source | Description | Status | Evidence |
|------------|--------|-------------|--------|----------|
| EXEC-01 | Plan | Monaco editor in Practice tab | ✓ SATISFIED | MonacoEditor wrapped in CodeEditor, integrated in topic-viewer.tsx PracticeTab |
| EXEC-02 | Plan | Real-time execution streaming | ⚠️ PARTIAL | fetch POST works but no SSE/streaming - output only after completion |
| EXEC-03 | Plan | Secure sandbox with resource limits | ✗ BLOCKED | Blocked packages enforced (os, net, syscall, unsafe, runtime/debug), but memory/cpu limits NOT applied to docker run |
| EXEC-04 | Plan | Topic-specific package allowlists | ✓ SATISFIED | getTopicAllowlist() returns topic-specific allowlists (grpc, nats, kafka, prometheus) |

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| backend/internal/executor/docker_executor.go | 454-461 | Missing --memory/--cpus flags | 🛑 Blocker | Memory/CPU limits not enforced despite being defined |
| frontend/src/components/workspace/CodeEditor.tsx | 54-82 | Simple fetch without streaming | ⚠️ Warning | Output only available after full execution completes |

### Human Verification Required

1. **End-to-end execution flow**
   - Test: Open TopicViewer → Practice tab → type Go code → click Run → see output
   - Expected: Output appears within seconds for simple "Hello World"
   - Why human: Need browser interaction to verify full flow

2. **Blocked package enforcement**
   - Test: Try importing "os" or "net" in code → click Run → see security error
   - Expected: Error message "Package 'os' is not allowed"
   - Why human: Need to verify error propagates correctly from backend to UI

3. **Topic allowlist enforcement**
   - Test: Use grpc-specific package on a non-grpc topic → Run → see error
   - Expected: "Package 'google.golang.org/grpc' is not allowed for this topic"
   - Why human: Need to verify topic-specific filtering works

### Gaps Summary

**Gap 1: Missing Memory/CPU Enforcement (EXEC-03 partial failure)**

The `DockerExecutor` struct defines `memory` and `cpuLimit` fields:
```go
type DockerExecutor struct {
    image     string
    timeout   time.Duration
    memory    string     // "256m"
    cpuLimit  string     // "1.0"
    tmpfsSize string
}
```

But `runContainer()` never uses them - the docker run command is built without `--memory` or `--cpus` flags:
```go
args := []string{
    "run",
    "--rm",
    "-i",
    e.image,
    "sh", "-c",
    fmt.Sprintf("printf '%%s' '%s' > /tmp/main.go && go run /tmp/main.go", ...),
}
```

**Impact:** User code could consume unlimited memory or CPU despite requirements stating resource limits.

**Gap 2: No Streaming Output (EXEC-02 partial failure)**

`handleExecute` uses a simple `await fetch()`:
```typescript
const response = await fetch('/api/execute', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code, topic_id: topicId, test_cases: testCases }),
});
const result: ExecuteResult = await response.json();
setOutput(result);
```

The requirement says "real-time (streaming)" but the current implementation waits for the entire response before displaying anything. Long-running code would show no output until completion.

**Impact:** Poor UX for longer executions - user sees no feedback during execution.

---

_Verified: 2026-04-01_
_Verifier: the agent (gsd-verifier)_
