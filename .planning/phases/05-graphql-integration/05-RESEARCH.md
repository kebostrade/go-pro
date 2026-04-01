# Phase 5: GraphQL & Integration — Research

**Researched:** 2026-04-01
**Domain:** GraphQL APIs with gqlgen
**Confidence:** MEDIUM (verified against existing codebase patterns, Context7 unavailable)

---

## Summary

Phase 5 covers GraphQL API development with gqlgen, the most popular schema-first GraphQL library for Go. This is the smallest phase with only 1 topic, focusing on production-grade GraphQL patterns including schema design, resolver implementation, authentication, and subscriptions.

**Primary recommendation:** Use gqlgen 0.17+ with schema-first development, chi v5 router for HTTP handling, and relay preset for cursor-based pagination. The existing `advanced-topics/15-graphql-gqlgen/` provides excellent reference material to build upon.

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| github.com/99designs/gqlgen | v0.17.57 | GraphQL schema-first development | Most popular Go GraphQL library |
| github.com/go-chi/chi/v5 | v5.1.0 | HTTP routing | Consistent with Phase 1 REST API |
| github.com/gorilla/websocket | v1.5.3 | WebSocket for subscriptions | Industry standard |
| github.com/vektah/gqlparser/v2 | v2.5.19 | GraphQL schema parsing | Used by gqlgen |

**Installation:**
```bash
# gqlgen CLI
go install github.com/99designs/gqlgen@latest

# Server dependencies
go get github.com/99designs/gqlgen@v0.17.57
go get github.com/go-chi/chi/v5@v5.1.0
go get github.com/gorilla/websocket@v1.5.3
```

---

## Architecture Patterns

### Recommended Project Structure

```
basic/projects/graphql/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── graph/
│   │   ├── schema.graphqls     # GraphQL schema
│   │   ├── resolver.go         # Resolver interfaces
│   │   ├── resolver_impl.go    # Resolver implementations
│   │   └── middleware.go      # GraphQL middleware
│   ├── middleware/
│   │   └── auth.go            # JWT authentication
│   └── loader/
│       └── dataloader.go       # N+1 query prevention
├── pkg/
│   └── auth/
│       └── jwt.go              # JWT utilities
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── gqlgen.yml                 # gqlgen configuration
└── README.md
```

### gqlgen Configuration (gqlgen.yml)

```yaml
# gqlgen configuration
schema:
  - internal/graph/schema.graphqls

exec:
  filename: internal/graph/generated.go
  package: graph

model:
  filename: internal/graph/model_gen.go
  package: graph

resolver:
  layout: follow-schema
  dir: internal/graph
  package: graph

strict_validation: true

autobinding:
  - time.Time
```

### Schema-First Development Workflow

```
1. Define schema in schema.graphqls
2. Run: gqlgen generate
3. Implement resolver interfaces
4. Wire up server
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| GraphQL parsing | Custom parser | gqlparser/v2 | Edge cases handled, spec compliant |
| Schema validation | Manual checks | gqlgen strict mode | Built-in validation |
| Query cost analysis | Custom limiter | gqlgen complexity analysis | Proven implementation |
| N+1 queries | Manual batching | Data loaders | Industry pattern |

---

## Common Pitfalls

### Pitfall 1: N+1 Query Problem
**What goes wrong:** Resolving a list of users' posts triggers one query per user.
**How to avoid:** Use data loaders to batch database calls.

```go
// Data loader example
type Loaders struct {
    UserByID *dataloader.Loader
}

func NewLoaders(db *DB) *Loaders {
    return &Loaders{
        UserByID: dataloader.NewLoader(func(keys []string) []*dataloader.Result {
            // Batch fetch users
        }),
    }
}
```

### Pitfall 2: Missing Error Handling
**What goes wrong:** Returning raw errors exposes internal details.
**How to avoid:** Use gqlgen error helpers.

```go
return nil, gqlgen.Errorf("user not found: %s", id)
// Or with extensions:
return nil, gql.Errorf("not found")
```

### Pitfall 3: Resolver Too Complex
**What goes wrong:** Fat resolvers with business logic.
**How to avoid:** Keep resolvers thin, delegate to services.

---

## Code Examples

### Schema Definition (schema.graphqls)

```graphql
scalar Time

type User {
  id: ID!
  username: String!
  email: String!
  role: Role!
  posts: [Post!]!
  createdAt: Time!
}

type Post {
  id: ID!
  title: String!
  content: String!
  author: User!
  comments: [Comment!]!
  createdAt: Time!
}

enum Role {
  ADMIN
  USER
  GUEST
}

input CreateUserInput {
  username: String!
  email: String!
  password: String!
  role: Role = USER
}

type Query {
  user(id: ID!): User
  users: [User!]!
  me: User!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  deleteUser(id: ID!): Boolean!
}

type Subscription {
  userCreated: User!
}
```

### Resolver Implementation

```go
// Resolver struct holds dependencies
type resolver struct {
    db *DB
}

// Query resolvers
func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
    user, err := r.db.UserByID(id)
    if err != nil {
        return nil, gqlgen.Errorf("user not found: %s", id)
    }
    return user, nil
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
    user := auth.ForContext(ctx)
    if user == nil {
        return nil, gql.Errorf("not authenticated")
    }
    return user, nil
}

// Mutation resolvers
func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Validate
    if input.Username == "" {
        return nil, gqlgen.Errorf("username required")
    }
    // Create
    return r.db.CreateUser(input)
}
```

### GraphQL Middleware

```go
type loggingMiddleware struct{}

func (loggingMiddleware) ExtensionName() string {
    return "LoggingMiddleware"
}

func (loggingMiddleware) Validate(schema graphql.ExecutableSchema) error {
    return nil
}

func (m *loggingMiddleware) Middleware(ctx context.Context, next graphql.Resolver) (interface{}, error) {
    operation := graphql.GetOperationContext(ctx)
    start := time.Now()
    
    res, err := next(ctx)
    
    log.Printf("GraphQL: %s %s took %v", 
        operation.Operation, 
        operation.OperationName, 
        time.Since(start))
    
    return res, err
}
```

### JWT Authentication Middleware

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        claims, err := validateJWT(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        ctx := auth.ContextWithUser(r.Context(), claims.User)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Code-first GraphQL | Schema-first (gqlgen) | 2018+ | Type safety, codegen |
| Manual resolvers | Generated from schema | gqlgen 0.13+ | Less boilerplate |
| REST for everything | GraphQL for complex APIs | 2019+ | Flexible queries |
| Basic auth | JWT + GraphQL middleware | Industry standard | Unified auth |

---

## Open Questions

1. **Relay preset vs standard pagination?**
   - What we know: gqlgen supports relay preset with cursor-based pagination
   - What's unclear: Whether the learning template should use relay conventions
   - Recommendation: Use relay preset for production-ready pagination patterns

2. **Subscription transport?**
   - What we know: WebSocket is standard for GraphQL subscriptions
   - What's unclear: SSE fallback for environments that block WebSocket
   - Recommendation: WebSocket primary, document SSE as alternative

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies beyond Go modules)

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| gqlgen | Code generation | ✓ | 0.17.57 | — |

---

## Sources

### Primary (HIGH confidence)
- Existing `advanced-topics/15-graphql-gqlgen/` codebase — verified patterns
- gqlgen official docs (gqlgen.com) — configuration reference

### Secondary (MEDIUM confidence)
- gqlgen GitHub repository — examples, issues
- GraphQL best practices — official spec

### Tertiary (LOW confidence)
- WebSearch gqlgen tutorials — various quality, need verification

---

## Metadata

**Confidence breakdown:**
- Standard stack: MEDIUM — existing codebase verified, Context7 unavailable
- Architecture: HIGH — follows established project patterns
- Pitfalls: MEDIUM — based on common GraphQL issues

**Research date:** 2026-04-01
**Valid until:** 2026-05-01 (30 days for stable topic)
