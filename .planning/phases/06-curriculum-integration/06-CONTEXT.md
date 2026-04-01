# Phase 6: Curriculum Integration - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers lesson pages and exercise definitions for all 15 advanced Go project topics, integrated into the existing curriculum hub with progress tracking.

</domain>

<decisions>
## Implementation Decisions

### Lesson Structure
- **D-01:** Tutorial-style with tabs — Overview → Content → Practice
- **D-02:** Reuse existing `TutorialViewer` component pattern from `frontend/src/components/learning/tutorial-viewer.tsx`
- **D-03:** Each topic has hero section with icon, title, description, duration, difficulty
- **D-04:** Overview tab shows topics covered, learning outcomes, prerequisites, quick actions
- **D-05:** Content tab shows project path, download instructions, run commands
- **D-06:** Practice tab shows exercises (placeholder until Phase 7 code execution)

### Exercise Format
- **D-07:** Exercises defined as embedded comments in starter code files
- **D-08:** Exercise markers: `// EXERCISE: [description]` format in Go code comments
- **D-09:** Each exercise has starter code + expected solution pattern
- **D-10:** Exercises are self-contained within each project template's `exercises/` directory

### Navigation
- **D-11:** Central Hub page at `/curriculum` (already exists)
- **D-12:** Sidebar shows all 15 topics grouped by phase
- **D-13:** Topic cards on hub link to `/learn/[topic-slug]`
- **D-14:** Use existing `curriculum/page.tsx` Tabs structure for phase grouping

### Progress Tracking
- **D-15:** Topic-level progress tracking — marks complete when user views key files and runs project
- **D-16:** Per-exercise tracking — each exercise within a topic is tracked individually
- **D-17:** Progress persists across sessions via existing backend API
- **D-18:** ProgressTracker component already tracks Lessons, Exercises, Projects, XP, Streaks — extend for 15 topics

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Existing Components
- `frontend/src/components/learning/tutorial-viewer.tsx` — Reference implementation for lesson pages
- `frontend/src/components/learning/progress-tracker.tsx` — Reference for progress tracking UI
- `frontend/src/app/curriculum/page.tsx` — Reference for hub/navigation structure

### Backend API
- `backend/internal/repository/interfaces.go` — Repository contracts for progress tracking
- `backend/internal/handler/` — HTTP handlers for curriculum API

### Project Templates
- `basic/projects/*/` — 15 project templates to integrate as lessons

### Existing Patterns
- `frontend/src/components/learning/` — All learning-related components

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `TutorialViewer` component: Tab-based layout with Overview/Content/Practice — extend for topics
- `ProgressTracker` component: Already tracks Lessons, Exercises, Projects, XP, Streaks
- `LessonCard` component: Used in curriculum page for lesson display
- `Badge`, `Card`, `Button`, `Progress` UI components — already styled

### Established Patterns
- Tabs for switching between views (TutorialViewer, Curriculum page)
- Card-based layouts for lesson/topic display
- API client in `frontend/src/lib/api.ts` with typed responses
- Firebase Auth context for user identification

### Integration Points
- New topic pages at `/learn/[topic]` — new route
- Progress API updates via `api.updateProgress()` 
- Curriculum data from `api.getCurriculum()`

</code_context>

<specifics>
## Specific Ideas

- 15 topics map to existing project templates in `basic/projects/`
- Each topic page uses same TutorialViewer structure but with topic-specific content
- Progress tracking extends existing `ProgressStats` interface
- Exercise completion tracked via submission API (Phase 9 will add AI review)

</specifics>

<deferred>
## Deferred Ideas

None — all decisions stayed within phase scope

</deferred>

---

*Phase: 06-curriculum-integration*
*Context gathered: 2026-04-01*
