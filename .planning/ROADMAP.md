# Roadmap: AI-Powered Mock Interviews v1.2

## Overview

Upgrade existing mock interview feature to use AI/LLM for question generation and personalized feedback. Users practice coding, behavioral, and system design interviews with AI that adapts to their skill level.

## Milestones

- ✅ **v1.0 Advanced Topics Expansion** - Phases 1-5 (shipped 2026-04-01)
- ✅ **v1.1 Platform Enhancements** - Phases 6-9 (shipped 2026-04-02)
- 🚧 **v1.2 AI-Powered Mock Interviews** - Phases 10-12 (in progress)

## Phases

- [ ] **Phase 10: AI Question Bank** - Database schema and question repository
- [ ] **Phase 11: AI Interview Session** - LLM-powered interview flow
- [ ] **Phase 12: AI Feedback & Progress** - Detailed feedback and progress tracking

---

## Phase Details

### Phase 10: AI Question Bank

**Goal**: Store and manage curated interview questions that the LLM can select from

**Depends on**: Nothing (first phase of new milestone)

**Requirements**: INTW-01a, INTW-01b, INTW-01c

**Success Criteria** (what must be TRUE):
1. User has access to a database of 50+ interview questions across types
2. Questions are tagged with difficulty, concepts, and categories
3. API can query questions by type, difficulty, and tags

**Plans**: TBD

**UI hint**: no

---

### Phase 11: AI Interview Session

**Goal**: LLM-powered interview experience with dynamic question selection and real-time analysis

**Depends on**: Phase 10

**Requirements**: INTW-02a, INTW-02b, INTW-02c, INTW-03a, INTW-03b, INTW-03c, INTW-04a, INTW-04b, INTW-04c

**Success Criteria** (what must be TRUE):
1. User starts interview and receives AI-selected questions
2. AI interviewer presents questions naturally with context
3. User answers and receives immediate analysis
4. AI asks follow-up questions based on answers

**Plans**: TBD

**UI hint**: yes

---

### Phase 12: AI Feedback & Progress

**Goal**: Comprehensive AI feedback after interviews with progress tracking

**Depends on**: Phase 11

**Requirements**: INTW-05a, INTW-05b, INTW-05c, INTW-05d, INTW-06a, INTW-06b, INTW-06c, INTW-06d, INTW-07a, INTW-07b

**Success Criteria** (what must be TRUE):
1. User receives detailed feedback with strengths and improvements
2. User sees personalized study recommendations
3. User can view progress over time with trends
4. User can see how they rank on leaderboard (nice to have)

**Plans**: TBD

**UI hint**: yes

---

## Progress

| Phase | Plans | Status | Completed |
|-------|-------|--------|-----------|
| 10. AI Question Bank | 0/3 | Not started | - |
| 11. AI Interview Session | 0/9 | Not started | - |
| 12. AI Feedback & Progress | 0/10 | Not started | - |

---

## Implementation Notes

### Phase 10: AI Question Bank
- Extend existing `memory_interview.go` repository
- Create new Question table/schema
- Add API endpoints for CRUD on questions
- Seed initial 50+ questions across types

### Phase 11: AI Interview Session
- Leverage existing interview handler (`backend/internal/handler/interview.go`)
- Integrate with AI Agent Platform for LLM calls
- Update question selection to use LLM
- Add real-time answer analysis

### Phase 12: AI Feedback & Progress
- Generate detailed feedback using LLM
- Extend session storage to include feedback
- Create progress dashboard UI
- Add leaderboard calculation

---

## Revision History

| Date | Phase | Change |
|------|-------|--------|
| 2026-04-01 | 1-5 | ✅ MILESTONE COMPLETE — 15 project templates |
| 2026-04-01 | 6-9 | 🚧 MILESTONE STARTED — Platform Enhancements |
| 2026-04-02 | 6-9 | ✅ MILESTONE COMPLETE — Platform Enhancements |
| 2026-04-02 | 10-12 | 🚧 MILESTONE STARTED — AI-Powered Mock Interviews |
