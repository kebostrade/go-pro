# Pitfalls Research

**Domain:** Learning Platform — Code Execution & Curriculum Integration
**Researched:** 2026-04-01
**Confidence:** MEDIUM

## Critical Pitfalls

### Pitfall 1: Arbitrary Code Execution Security

**What goes wrong:** Users can execute malicious code, access system resources, mine cryptocurrency, etc.

**Why it happens:**
- Underestimating creativity of malicious users
- Insufficient sandboxing
- Missing resource limits

**How to avoid:**
1. Use gVisor (runsc) or equivalent for isolation
2. Apply resource limits (CPU, memory, time)
3. Disable network access in sandbox
4. Scan for known malicious patterns before execution

**Warning signs:**
- Unexpected processes on host
- High CPU usage from executor
- Network connections from execution container
- Disk I/O spikes

**Phase to address:** Phase 7 (Code Execution)

---

### Pitfall 2: Long Execution Times Blocking UI

**What goes wrong:** Code takes 30+ seconds to run, UI freezes, user thinks it's broken

**Why it happens:**
- Infinite loops in user code
- Heavy computation
- Network calls with long timeouts

**How to avoid:**
1. Set strict timeout (10-30 seconds max)
2. Show real-time progress/status
3. Implement streaming output so user sees activity
4. Provide "Stop" button

**Warning signs:**
- Users complaining "it's stuck"
- Memory usage climbing
- WebSocket connections timing out

**Phase to address:** Phase 7 (Code Execution)

---

### Pitfall 3: Curriculum-template Coupling

**What goes wrong:** Lesson pages tightly coupled to template structure, changes break content

**Why it happens:**
- Hard-coded paths to template files
- Template changes require content updates
- No abstraction layer

**How to avoid:**
1. Define exercise schema (JSON/YAML) separate from templates
2. Use relative paths and stable interfaces
3. Version the exercise format
4. Provide migration tools for format changes

**Warning signs:**
- "I changed the template and lessons broke"
- Hard-coded topic IDs
- Exercise validation tied to specific file structure

**Phase to address:** Phase 6 (Curriculum Integration)

---

### Pitfall 4: No Progress Persistence

**What goes wrong:** User completes exercises but progress isn't saved, they have to redo it

**Why it happens:**
- LocalStorage cleared
- Not connected to user account
- Save failures silently ignored

**How to avoid:**
1. Auto-save on significant actions
2. Sync to backend on auth
3. Visual confirmation of save
4. Offline-capable with eventual consistency

**Warning signs:**
- Users reporting "my progress disappeared"
- Save errors in logs
- Disconnect between frontend state and backend

**Phase to address:** Phase 6 (Curriculum Integration)

---

### Pitfall 5: AI Feedback Quality Inconsistency

**What goes wrong:** AI gives wildly different quality feedback, sometimes useless or wrong

**Why it happens:**
- No prompt engineering
- Varying code complexity confuses model
- No validation of feedback quality

**How to avoid:**
1. Create structured feedback format (categories, scores)
2. Implement feedback templates for common issues
3. Test AI with known code samples
4. Allow user to regenerate or get human review

**Warning signs:**
- Users complaining about AI feedback
- Feedback doesn't match code
- Inconsistent tone/quality

**Phase to address:** Phase 9 (Code Review System)

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Skip sandbox security | Faster dev | Security breach, data leak | Never |
| Hard-code exercise paths | Simpler | Break on template changes | Never |
| Sync AI calls | Simpler | UI freeze, timeout | v1 only |
| One executor instance | Low resource | Slow under load | v1 only |

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Docker | Expect Docker in all envs | Check Docker availability, fallback |
| Monaco | Not resizing properly | Use ResizeObserver |
| WebSocket | Not handling reconnect | Implement exponential backoff |
| AI Agent | No timeout | Set 30s max, queue for later |

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Large code in Monaco | Editor lag | Limit to 10k chars | Files > 500 lines |
| Long output streaming | Memory buildup | Chunk output, truncate display | Infinite print loop |
| Many concurrent executions | Resource exhaustion | Queue + worker pool | > 10 simultaneous |
| AI agent not cached | Slow repeated reviews | Cache common patterns | Repeated errors |

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| No rate limiting on execution | DoS, crypto mining | IP + user rate limits |
| Execute with elevated privileges | Container escape | Run as non-root user |
| Store secrets in execution env | Credential leak | Use ephemeral environments |
| Allow network access | Data exfiltration | Block all outbound except needed |

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| No save indicator | "Did my code save?" | Show save status always |
| Cryptic error messages | "What does E001 mean?" | Human-readable errors |
| Force re-login | Lost work, frustration | Keep session alive, refresh token |
| Hidden progress | "Am I doing well?" | Show completion %, badges |

## "Looks Done But Isn't" Checklist

- [ ] **Execution**: Verified with infinite loop — times out properly
- [ ] **Execution**: Verified with network call — blocked properly
- [ ] **Curriculum**: Verified after template change — still works
- [ ] **Progress**: Verified after browser close — persists
- [ ] **Review**: Verified with syntax error code — gives useful feedback

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Security breach | HIGH | Stop execution, audit logs, notify users |
| Data loss | HIGH | Restore from backup, check integrity |
| AI feedback wrong | MEDIUM | Provide human review option, fix prompt |
| Progress lost | MEDIUM | Restore from last sync, check frontend logs |

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Security | Phase 7 | Pentest, timeout tests |
| Long execution | Phase 7 | Load test with infinite loop |
| Curriculum coupling | Phase 6 | Template change test |
| Progress loss | Phase 6 | Clear browser data test |
| AI quality | Phase 9 | User feedback survey |

## Sources

- OWASP Web Security
- Go Playground security model
- Docker security best practices
- gVisor documentation

---
*Pitfalls research for: Platform Enhancements*
*Researched: 2026-04-01*
