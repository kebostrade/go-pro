# Phase 9: Code Review System - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers AI-powered code submission and review for learner exercises. Users submit code and receive conversational feedback from an AI agent, with submission history stored per user.

</domain>

<decisions>
## Implementation Decisions

### Feedback Format
- **D-01:** Conversational feedback — AI explains in plain English like a code review mentor
- **D-02:** Uses existing CodeAnalysisTool from `services/ai-agent-platform/`
- **D-03:** Feedback covers: code quality, security issues, performance, best practices, improvement suggestions
- **D-04:** AI speaks encouragingly — acknowledges what's done well before suggesting improvements

### Submission Flow
- **D-05:** ExerciseCard submit button — Submit button on each exercise card
- **D-06:** Per-exercise review — each exercise has its own submission
- **D-07:** Visual feedback during review — "AI is reviewing your code..." spinner
- **D-08:** Feedback displayed in expandable section below the exercise

### History Storage
- **D-09:** Backend API — PostgreSQL with user associations
- **D-10:** Submission stored with: user_id, topic_id, exercise_id, code, feedback, submitted_at
- **D-11:** History API: GET /api/review/history (user's submissions)
- **D-12:** Unauthenticated users: feedback only, no history

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Existing AI Infrastructure
- `services/ai-agent-platform/internal/tools/programming/code_analysis.go` — CodeAnalysisTool
- `services/ai-agent-platform/internal/languages/golang/analyzer.go` — Go language analyzer
- `services/ai-agent-platform/pkg/types/coding.go` — Analysis result types
- `services/ai-agent-platform/internal/agent/react.go` — ReAct agent pattern

### Frontend Infrastructure
- `frontend/src/components/learning/exercise-card.tsx` — ExerciseCard with completion tracking
- `frontend/src/lib/api.ts` — API client pattern

### Backend Infrastructure
- `backend/internal/repository/interfaces.go` — Repository pattern
- `backend/internal/handler/` — HTTP handler pattern
- `backend/cmd/server/main.go` — Route registration

### Requirements
- REVIEW-01: User can submit code exercises for review
- REVIEW-02: AI agent analyzes and provides structured feedback
- REVIEW-03: User can view submission history

</canonical_refs>

<codebase_context>
## Existing Code Insights

### Reusable Assets
- CodeAnalysisTool with Analyze() method
- Go language analyzer with security/performance/best-practice checks
- ReAct agent for multi-step reasoning
- Backend repository pattern for persistence

### Established Patterns
- API client in frontend (api.ts)
- Handler registration in main.go
- Repository interfaces for data access
- Tool registration in agent

### Integration Points
- ExerciseCard → Submit button → Backend API
- Backend API → AI Agent Platform → CodeAnalysisTool
- History display in TopicViewer or dedicated page

</codebase_context>

<specifics>
## Specific Ideas

- New submit button on ExerciseCard
- POST /api/review/submit endpoint
- GET /api/review/history endpoint
- Feedback displayed in expandable panel
- AI feedback formatted conversationally

</specifics>

<deferred>
## Deferred Ideas

None — all decisions stayed within phase scope

</deferred>

---

*Phase: 09-code-review-system*
*Context gathered: 2026-04-01*
