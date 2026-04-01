# Feature Research

**Domain:** Learning Platform — Code Execution & Curriculum Integration
**Researched:** 2026-04-01
**Confidence:** HIGH

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Code Editor** | Users need to write/edit code in-browser | LOW | Monaco editor already exists |
| **Run Code** | Execute code and see output | LOW | Playground already exists |
| **Submit Exercise** | Save work for review | LOW | Needs enhancement |
| **View Lessons** | Read content for each topic | LOW | Course module exists |
| **Track Progress** | See completion status | LOW | Progress API exists |
| **Docker Setup** | One-click environment | MEDIUM | Templates exist, need orchestration |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valuable.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **AI Code Review** | Personalized feedback without waiting | HIGH | Integrate existing AI agent |
| **Streaming Output** | Real-time execution feedback | MEDIUM | WebSocket already in stack |
| **Per-Topic Docker** | Tailored environments per project | MEDIUM | 15 different topics |
| **Interactive Exercises** | Fill-in-the-blank coding | MEDIUM | Parser + validation needed |
| **Side-by-Side Comparison** | Compare solution vs submitted | MEDIUM | Diff view component |

### Anti-Features (Commonly Requested, Often Problematic)

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| **Video content** | "Everyone does it" | Production overhead, no searchability | Keep text/code focused |
| **Live chat with instructor** | "Support" expectation | 1:1 scaling impossible | AI-powered Q&A |
| **Real-time collaboration** | "Study together" | Conflict resolution hard | Async code review |
| **Mobile app** | "On the go learning" | Development overhead | Responsive web |

## Feature Dependencies

```
[Curriculum Integration]
    └──requires──> [Lesson Pages for 15 topics]
    └──requires──> [Exercise Definitions]
    └──requires──> [Progress Tracking API]

[Code Execution]
    └──requires──> [Execution Sandbox]
    └──requires──> [Output Streaming]

[Docker Setup]
    └──requires──> [Curriculum Integration] (knows which topic)
    └──requires──> [Docker Compose Generator]

[Code Review System]
    └──requires──> [Code Execution] (run tests)
    └──requires──> [Submission Storage]
    └──enhances──> [AI Agent Integration]
```

## MVP Definition

### Launch With (v1.1)

Minimum viable product — what's needed to validate the concept.

- [x] **CURR-01**: Lesson pages for all 15 topics (link existing templates)
- [ ] **CURR-02**: Exercise definitions per topic (extensible format)
- [ ] **EXEC-01**: In-browser code execution for all 15 topics
- [ ] **DOCK-01**: One-click Docker setup (docker-compose per topic)
- [ ] **REVIEW-01**: Code submission storage and retrieval
- [ ] **REVIEW-02**: Basic AI-powered feedback (leverage existing agent)

### Add After Validation (v1.x)

Features to add once core is working.

- [ ] **EXEC-02**: Streaming execution output (WebSocket)
- [ ] **REVIEW-03**: Detailed diff view (solution vs submitted)
- [ ] **CURR-03**: Interactive fill-in-blank exercises

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] **EXEC-03**: Collaborative editing (real-time)
- [ ] **REVIEW-04**: Human code review scheduling
- [ ] **CURR-04**: Video content integration

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Lesson pages (15 topics) | HIGH | LOW | P1 |
| Code execution | HIGH | MEDIUM | P1 |
| Docker setup | MEDIUM | MEDIUM | P1 |
| Exercise definitions | HIGH | MEDIUM | P1 |
| Code submission | MEDIUM | LOW | P1 |
| AI code review | HIGH | HIGH | P2 |
| Streaming output | MEDIUM | MEDIUM | P2 |
| Diff view | MEDIUM | MEDIUM | P2 |

## Competitor Feature Analysis

| Feature | Exercism | Go.dev | Our Approach |
|---------|----------|--------|--------------|
| Code execution | Web terminal | Playground | Embed Monaco + execution API |
| Docker setup | None | None | docker-compose per topic |
| Code review | Human mentors | None | AI-first, human optional |
| Progress tracking | Yes | No | Existing API extend |

## Sources

- Exercism.io — code practice platform
- Go.dev — official Go learning
- LeetCode — coding challenges
- GitHub Codespaces — cloud development environments

---
*Feature research for: Platform Enhancements*
*Researched: 2026-04-01*
