# Phase 1 Research: Foundation Patterns

**Phase:** 1 - Foundation Patterns  
**Topics:** RESTful APIs, CLI Applications, Testing & Debugging, Web Apps with Gin  
**Research Date:** 2026-04-01  
**Status:** Initial research needed

## Research Required

Before implementing Phase 1 templates, research needed for:

### 1. Go REST API Ecosystem (Task 1)

**Key Questions:**
- Standard library `net/http` vs frameworks (chi, gorilla/mux, gin)?
- What router pattern to recommend?
- Error handling conventions?
- Middleware patterns?

**Research Sources to Check:**
- Context7: golang net/http
- Context7: chi router
- Context7: gorilla mux
- Context7: gin framework

### 2. Go CLI Ecosystem (Task 3)

**Key Questions:**
- Cobra vs urfave/cli vs clicumber?
- Structure for CLI projects (cmd/ pattern)?
- Argument parsing best practices?
- Testing CLI applications?

**Research Sources to Check:**
- Context7: cobra cli
- Context7: urfave cli

### 3. Go Testing Patterns (Task 4)

**Key Questions:**
- testify vs standard testing?
- Mocking frameworks (gock, mockery, minimock)?
- Integration testing patterns?
- Benchmarking best practices?

**Research Sources to Check:**
- Context7: golang testing
- Context7: testify

### 4. Gin Web Framework (Task 5)

**Key Questions:**
- Gin vs standard library?
- Middleware patterns?
- Binding and validation?
- Error handling?

**Research Sources to Check:**
- Context7: gin framework

## Template Structure (To Be Confirmed)

Based on existing `basic/projects/` pattern, each template should follow:

```
basic/projects/{topic}/
├── cmd/                    # Entry points
│   └── server/main.go
├── internal/               # Private code
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/                    # Public packages
├── migrations/             # DB migrations (if applicable)
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml
├── go.mod
├── go.sum
├── README.md
└── Makefile
```

## Deliverables Checklist

After research, confirm:

- [ ] Router choice for REST API (chi vs gin)
- [ ] CLI framework choice (cobra vs urfave)
- [ ] Testing framework (testify + mock strategy)
- [ ] Gin middleware patterns
- [ ] Common project structure across all 4 templates

## Next Steps

1. Run research on each topic area
2. Document decisions in STACK.md
3. Create first template (REST API) as reference
4. Clone structure for remaining 3 templates
