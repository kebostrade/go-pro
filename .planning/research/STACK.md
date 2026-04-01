# Stack Research

**Domain:** Learning Platform — Code Execution & Curriculum Integration
**Researched:** 2026-04-01
**Confidence:** HIGH

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Go Playground** | latest | In-browser Go code execution | Official Go team service, battle-tested sandbox |
| **WebContainer API** | latest | Browser-based Docker alternative | Lightweight, no server round-trip for simple cases |
| **Docker-in-Docker (DinD)** | latest | Full containerized execution | Complete isolation, supports all Go versions |
| **gVisor** | latest | Secure container runtime | Better security than DinD, faster than full VMs |
| **Kubernetes** | 1.28+ | Container orchestration | Already in project templates, natural fit |

### Supporting Libraries

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **go-playground/validator** | v10 | Input validation | Validate code execution requests |
| **gorilla/websocket** | v1.5 | Real-time output streaming | Stream execution results to browser |
| **chi** | v5 | HTTP routing | API endpoints (already in use) |

### Code Execution Strategies

| Approach | Security | Performance | Complexity | Best For |
|----------|----------|-------------|------------|----------|
| **Go Playground API** | ✅ High (Go team's sandbox) | Fast | Low | Simple code snippets |
| **WebContainer** | ✅ High (browser sandbox) | Very Fast | Medium | Frontend, simple Docker |
| **gVisor (runsc)** | ✅ High (kernel isolation) | Fast | Medium | Untrusted user code |
| **DinD** | ⚠️ Medium (full container) | Medium | Low | Full system access needed |
| **AWS Lambda** | ✅ High (managed) | Scalable | Medium | High volume, variable load |

### Go-Specific Considerations

- **Multiple Go versions**: Use containers with different Go versions mounted
- **CGO**: Disable CGO for security, enable for specific packages (sqlite, etc.)
- **Network access**: Restrict in sandbox, allow in exercises with firewall rules
- **File system**: Use tmpfs for ephemeral execution, persistent storage for submissions

### Integration Points

| Component | Purpose | Integration |
|-----------|---------|-------------|
| `frontend/playground` | Monaco editor + execution UI | Already exists |
| `backend/executor` | Code execution service | Already exists, needs enhancement |
| `backend/api` | Course/lesson management | Already exists |
| `course/` module | Curriculum content | Already exists |

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| **Raw process exec** | Security risk, no isolation | gVisor or containers |
| **Eval()** | Dangerous, no Go equivalent | N/A |
| **Full VMs (QEMU)** | Too slow for interactive use | gVisor for security |

## Version Compatibility

| Component | Compatible With | Notes |
|-----------|-----------------|-------|
| Go 1.23 | gVisor runsc | Works out of box |
| Docker | Kubernetes | Already integrated |
| Monaco Editor | Next.js 15 | Already integrated |

## Sources

- Go Playground: https://play.golang.org/
- WebContainer API: https://github.com/near/NearJS/wt-cli
- gVisor: https://gvisor.dev/docs/
- Docker security best practices

---
*Stack research for: Platform Enhancements*
*Researched: 2026-04-01*
