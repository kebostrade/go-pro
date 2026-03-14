# Building RESTful APIs with Go

Create robust, production-ready REST APIs using Go's net/http and best practices.

## Learning Objectives

- Design RESTful API endpoints
- Implement proper HTTP handling
- Handle JSON request/response
- Apply middleware patterns
- Implement authentication/authorization
- Write API documentation

## Theory

### REST Principles

- **Stateless**: Each request contains all needed information
- **Resource-based**: URLs represent resources
- **HTTP Methods**: GET, POST, PUT, PATCH, DELETE
- **Status Codes**: Proper HTTP status usage

### Basic Server Structure

```go
type Server struct {
    router *http.ServeMux
    db     *sql.DB
}

func NewServer(db *sql.DB) *Server {
    s := &Server{
        router: http.NewServeMux(),
        db:     db,
    }
    s.routes()
    return s
}

func (s *Server) routes() {
    s.router.HandleFunc("GET /api/users", s.listUsers)
    s.router.HandleFunc("POST /api/users", s.createUser)
    s.router.HandleFunc("GET /api/users/{id}", s.getUser)
    s.router.HandleFunc("PUT /api/users/{id}", s.updateUser)
    s.router.HandleFunc("DELETE /api/users/{id}", s.deleteUser)
}
```

### JSON Handling

```go
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    user, err := s.userService.Create(r.Context(), &req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, user)
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}
```

### Middleware Pattern

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            respondError(w, http.StatusUnauthorized, "unauthorized")
            return
        }
        ctx := context.WithValue(r.Context(), "user_id", validateToken(token))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Security Considerations

```go
type ValidationMiddleware struct {
    maxBodySize int64
}

func (v *ValidationMiddleware) Wrap(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, v.maxBodySize)
        next.ServeHTTP(w, r)
    })
}

func sanitizeInput(input string) string {
    return html.EscapeString(strings.TrimSpace(input))
}
```

## Performance Tips

```go
var responsePool = sync.Pool{
    New: func() interface{} {
        return &APIResponse{}
    },
}

func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    users, err := s.userService.List(ctx)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, users)
}
```

## Exercises

1. Build a CRUD API for products
2. Implement JWT authentication
3. Add rate limiting middleware
4. Create API documentation

## Validation

```bash
cd exercises
go test -v ./...
curl -X GET http://localhost:8080/api/users
```

## Key Takeaways

- Use proper HTTP methods and status codes
- Validate all inputs
- Implement middleware for cross-cutting concerns
- Always use context for cancellation
- Document your API

## Next Steps

**[AT-04: Gin Framework](../AT-04-gin-framework/README.md)**

---

Build APIs that scale and last. 🔌
