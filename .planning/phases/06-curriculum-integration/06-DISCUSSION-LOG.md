# Phase 6: Curriculum Integration - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-01
**Phase:** 06-curriculum-integration
**Areas discussed:** Lesson Structure, Exercise Format, Navigation, Progress Tracking

---

## Lesson Structure

| Option | Description | Selected |
|--------|-------------|----------|
| A | Tutorial-style with tabs — Overview → Content → Practice. TutorialViewer already exists with this structure. | ✓ |
| B | Reference-style — Project template with README, code browser, and inline exercises. More like GitHub documentation. | |
| C | Hybrid — Lesson overview first, then link to full project reference for deep dives. | |

**User's choice:** A — Tutorial-style with tabs
**Notes:** Reuse existing TutorialViewer component pattern

---

## Exercise Format

| Option | Description | Selected |
|--------|-------------|----------|
| A | YAML/JSON files — Exercise definitions in `basic/projects/[topic]/exercises.yaml`. Decoupled from code, easy to parse. | |
| B | Embedded in code comments — `// EXERCISE: [description]` markers in starter code. Simpler but harder to validate. | ✓ |
| C | Markdown with code blocks — Exercise descriptions in markdown files with fenced code. | |

**User's choice:** B — Embedded in code comments
**Notes:** Simpler approach, exercises embedded directly in starter code files

---

## Navigation

| Option | Description | Selected |
|--------|-------------|----------|
| A | Central Hub + Sidebar — Curriculum page shows topic cards. Sidebar shows 15 topics grouped by phase. Uses existing Tabs structure. | ✓ |
| B | Linear Path — Strict sequential order through all 15 topics like a course. | |
| C | Topic Browser — Grid of topic cards with filtering by difficulty/domain. Users pick what to learn. | |

**User's choice:** A — Central Hub + Sidebar
**Notes:** Existing curriculum page structure works well

---

## Progress Tracking

| Option | Description | Selected |
|--------|-------------|----------|
| A | Topic-level only — Mark topic complete when user views key files and runs project. Simple. | |
| B | Per-exercise tracking — Track each exercise individually within topics. More granular but more complex. | |
| C | Both topic + exercise — Topic shows aggregate progress, individual exercises tracked separately. Most complete but most work. | ✓ |

**User's choice:** C — Both topic-level and per-exercise tracking
**Notes:** Most comprehensive approach for tracking learner progress

---

## Phase Completion Decisions

**All 4 areas discussed and decisions captured in CONTEXT.md.**

