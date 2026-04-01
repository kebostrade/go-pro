# Phase 5: GraphQL & Integration — Context

**Phase:** 5
**Name:** GraphQL & Integration
**Topics:** GraphQL APIs with Go and gqlgen
**Status:** Pending Research & Planning

---

## Decisions (Locked)

### Topic 16: GraphQL APIs with Go and gqlgen
- **Template:** `basic/projects/graphql/`
- **Focus:** Schema-first development, resolver implementation, query/mutation patterns
- **Library:** github.com/99designs/gqlgen (v0.17+)
- **Router:** github.com/go-chi/chi/v5 (consistent with Phase 1 REST API)
- **Deliverable:** Production-grade GraphQL API template with authentication, middleware, and subscriptions

### Standard Stack Decisions
| Component | Decision | Rationale |
|-----------|----------|-----------|
| Go Version | 1.23+ | Standardized across all modules |
| GraphQL Library | gqlgen 0.17+ | Schema-first, most popular Go GraphQL library |
| HTTP Router | chi v5 | Consistent with Phase 1 REST API template |
| Auth Pattern | JWT middleware | Consistent with backend platform auth |
| Subscriptions | WebSocket via gorilla/websocket | Industry standard for GraphQL subscriptions |

---

## the agent's Discretion

The following are **NOT locked** — the planner researches and recommends:

1. **Relay preset vs standard schema**: Should we use gqlgen's relay preset for pagination?
2. **Data loader pattern**: Implement N+1 query prevention with data loaders?
3. **Subscription transport**: WebSocket only, or also Server-Sent Events (SSE)?
4. **File upload handling**: graphql-upload pattern vs base64 encoding?

---

## Deferred Ideas (Out of Scope for Phase 5)

- ~~Federation (Apollo Federation) — separate gateway concern~~
- ~~Real-time subscriptions with NATS — Phase 3 already covers NATS~~
- ~~GraphQL schema stitching — too advanced for template~~
- ~~Production GraphQL security (Persisted queries, depth limiting) — add-on topic~~

---

## Dependencies

- **Requires:** Phase 1 REST API patterns (chi router, middleware patterns)
- **Leverages:** existing `advanced-topics/15-graphql-gqlgen/` content
- **Optional:** Phase 2 WebSocket patterns for subscription transport

---

## Phase 5 Task Breakdown

```
Phase 5: GraphQL & Integration
└── Task 16: Template - GraphQL (gqlgen, schema-first, relay)
```

---

## Quality Gates

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes with >80% coverage
- [ ] `golangci-lint run` passes
- [ ] `docker build` succeeds
- [ ] `docker-compose up` runs without errors
- [ ] CI pipeline green on GitHub Actions

---

## Notes

### Phase 5 Scale
This is the **smallest phase** with only 1 topic. The GraphQL template should:
- Be production-grade but focused
- Leverage existing `advanced-topics/15-graphql-gqlgen/` content
- Follow patterns established in Phase 1 (chi router, middleware chain)

### Existing GraphQL Content
The `advanced-topics/15-graphql-gqlgen/` directory already has:
- Complete schema definition (`schema.graphqls`)
- Server setup with chi router (`server.go`)
- Resolver implementations (`resolvers.go`)
- Models and mock database (`models.go`)
- Examples with authentication (`examples/`)
- gqlgen configuration (`gqlgen.yml`)

The template should be a **refined, production-ready version** of this content.

### gqlgen Schema-First Workflow
```
1. Define schema.graphqls
2. Run `gqlgen generate`
3. Implement resolvers in generated interfaces
4. Add custom middleware and extensions
```
