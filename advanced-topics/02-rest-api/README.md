# 02 - Building RESTful APIs with Go

Create production-ready RESTful APIs using Go best practices, proper architecture, and industry standards.

## 📚 Overview

Learn to build robust, scalable REST APIs with Go's `net/http` package and popular frameworks. This guide covers authentication, validation, database integration, testing, and deployment.

## 🎯 Learning Objectives

- Design RESTful API endpoints following best practices
- Implement proper HTTP methods and status codes
- Add authentication and authorization (JWT, OAuth)
- Validate input and handle errors gracefully
- Integrate with databases (PostgreSQL, MongoDB)
- Write comprehensive tests for APIs
- Implement middleware for logging, CORS, rate limiting
- Document APIs with OpenAPI/Swagger
- Deploy to production with Docker

## 🏗️ API Architecture

### Recommended Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── domain/                  # Business entities
│   ├── handler/                 # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── auth_handler.go
│   │   └── middleware.go        # HTTP middleware
│   ├── repository/              # Data access layer
│   │   ├── interfaces.go        # Repository interfaces
│   │   └── postgres/            # PostgreSQL implementation
│   ├── service/                 # Business logic
│   └── pkg/                     # Public packages
│       ├── auth/                # Authentication utilities
│       ├── validator/           # Input validation
│       └── response/            # HTTP response helpers
├── migrations/                  # Database migrations
├── docs/                        # API documentation
├── tests/                       # Integration tests
├── Dockerfile
├── .env.example
└── go.mod
```

## 📋 RESTful Design Principles

### 1. Resource-Based URLs

Use nouns, not verbs. Use plural nouns for collections.

```
✅ GOOD                    ❌ BAD
GET    /users              /getUsers
GET    /users/123          /user
POST   /users              /createUser
PUT    /users/123          /updateUser
DELETE /users/123          /deleteUser
```

### 2. HTTP Methods

| Method | Operation | Safe? | Idempotent? |
|--------|-----------|-------|-------------|
| GET | Read resource | ✅ | ✅ |
| POST | Create resource | ❌ | ❌ |
| PUT | Update/Replace | ❌ | ✅ |
| PATCH | Partial update | ❌ | ❌ |
| DELETE | Delete resource | ❌ | ✅ |

### 3. Status Codes

Use appropriate HTTP status codes:

| Code | Meaning | Use Case |
|------|---------|----------|
| **200 OK** | Success | GET, PUT, PATCH |
| **201 Created** | Resource created | POST |
| **204 No Content** | Success, no response | DELETE |
| **400 Bad Request** | Invalid input | Validation errors |
| **401 Unauthorized** | Not authenticated | Missing/invalid token |
| **403 Forbidden** | Authenticated but not authorized | Permissions |
| **404 Not Found** | Resource doesn't exist | Invalid ID |
| **409 Conflict** | Conflict with existing state | Duplicate resource |
| **422 Unprocessable** | Semantic errors | Business logic errors |
| **429 Too Many Requests** | Rate limited | Too many requests |
| **500 Internal Server Error** | Server error | Unhandled errors |

### 4. Request/Response Format

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

**Response:**
```json
{
  "data": {
    "id": "123",
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2024-01-15T10:30:00Z"
  },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total_pages": 5
  }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  }
}
```

## 🔐 Authentication & Authorization

### JWT Authentication

```go
// Generate JWT token
func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Middleware to validate JWT
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Validate token
        claims, err := ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user ID to context
        ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 🧪 Input Validation

### Struct Validation

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=100"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,gte=18,lte=120"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    // Decode
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Validate
    if err := validator.New().Struct(&req); err != nil {
        // Format validation errors
        errors := formatValidationErrors(err)
        respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
            "error": map[string]interface{}{
                "code":    "VALIDATION_ERROR",
                "message": "Validation failed",
                "details": errors,
            },
        })
        return
    }

    // Process request
    user, err := h.service.CreateUser(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    respondWithJSON(w, http.StatusCreated, user)
}
```

## 🗄️ Database Integration

### Using sqlx with PostgreSQL

```go
type UserRepository struct {
    db *sqlx.DB
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (id, name, email, age, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Name, user.Email, user.Age, time.Now(),
    )

    return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    query := `SELECT * FROM users WHERE id = $1`

    var user User
    if err := r.db.GetContext(ctx, &user, query, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
    query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

    var users []User
    if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
        return nil, err
    }

    return users, nil
}
```

## 🚦 Middleware Examples

### Logging Middleware

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status code
        wrapped := &responseWrapper{ResponseWriter: w}

        // Call next handler
        next.ServeHTTP(wrapped, r)

        // Log request
        duration := time.Since(start)
        log.Printf(
            "%s %s %d %v",
            r.Method,
            r.URL.Path,
            wrapped.status,
            duration,
        )
    })
}
```

### CORS Middleware

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Rate Limiting Middleware

```go
func RateLimitMiddleware(requests int, window time.Duration) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## 📊 Pagination

```go
type Pagination struct {
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    Total    int    `json:"total"`
}

func (p *Pagination) Offset() int {
    return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
    return p.PageSize
}

// Usage in handler
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
    // Parse pagination params
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }

    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    // Get users with pagination
    users, total, err := h.service.ListUsers(r.Context(), page, pageSize)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with paginated data
    respondWithJSON(w, http.StatusOK, map[string]interface{}{
        "data": users,
        "meta": Pagination{
            Page:     page,
            PageSize: pageSize,
            Total:    total,
        },
    })
}
```

## 🧪 Testing

### Unit Test Example

```go
func TestHandler_CreateUser(t *testing.T) {
    // Setup mock service
    mockService := &MockUserService{
        users: make(map[string]*User),
    }
    handler := NewUserHandler(mockService)

    // Create request body
    body := `{"name":"John Doe","email":"john@example.com","age":30}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    // Create response recorder
    w := httptest.NewRecorder()

    // Call handler
    handler.CreateUser(w, req)

    // Check response
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }

    var response User
    json.NewDecoder(w.Body).Decode(&response)

    if response.Name != "John Doe" {
        t.Errorf("Expected name 'John Doe', got '%s'", response.Name)
    }
}
```

### Integration Test Example

```go
func TestAPI_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Setup test database
    db := setupTestDB()
    defer db.Close()

    // Setup test server
    router := setupRouter(db)
    server := httptest.NewServer(router)
    defer server.Close()

    // Create user
    resp, err := http.Post(server.URL+"/users", "application/json",
        strings.NewReader(`{"name":"Test User","email":"test@example.com","age":25}`))

    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", resp.StatusCode)
    }

    // Get user
    resp, err = http.Get(server.URL + "/users/1")
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}
```

## 📝 API Documentation

### OpenAPI/Swagger

```go
// @title User API
// @version 1.0
// @description A user management service API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
    // Swagger setup
    http.HandleFunc("/swagger.json", swaggerHandler)
    http.Handle("/", http.FileServer(http.Dir("./swagger-ui")))

    // API routes
    http.HandleFunc("/api/v1/users", usersHandler)
    http.HandleFunc("/api/v1/users/", userByIDHandler)

    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## 🐳 Docker Deployment

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./server"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/mydb
      - JWT_SECRET=secret
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=mydb
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

## 🎯 Best Practices Checklist

### Design
- ✅ Use nouns for resource names
- ✅ Use plural nouns for collections
- ✅ Use kebab-case for URLs
- ✅ Implement pagination for list endpoints
- ✅ Provide filtering, sorting, and searching
- ✅ Version your API (`/api/v1/`)
- ✅ Use HTTP status codes correctly

### Security
- ✅ Always use HTTPS in production
- ✅ Implement authentication (JWT/OAuth)
- ✅ Validate all inputs
- ✅ Sanitize outputs
- ✅ Use parameterized queries
- ✅ Implement rate limiting
- ✅ Keep secrets in environment variables
- ✅ Use CORS carefully

### Performance
- ✅ Use connection pooling
- ✅ Implement caching (Redis)
- ✅ Compress responses (gzip)
- ✅ Use database indexes
- ✅ Implement lazy loading
- ✅ Monitor and profile

### Reliability
- ✅ Handle errors gracefully
- ✅ Implement retries with exponential backoff
- ✅ Use circuit breakers for external services
- ✅ Log important events
- ✅ Implement health checks
- ✅ Graceful shutdown

### Documentation
- ✅ Provide API documentation
- ✅ Include example requests/responses
- ✅ Document error codes
- ✅ Keep docs up to date

## 🚀 Quick Start

```bash
# Clone the examples
cd advanced-topics/02-rest-api/examples

# Run the example API
go run main.go

# Test the API
curl http://localhost:8080/api/v1/users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":30}'
```

## 📖 Further Reading

- [REST API Tutorial](https://restfulapi.net/)
- [Go net/http docs](https://pkg.go.dev/net/http)
- [JWT in Go](https://github.com/golang-jwt/jwt)
- [OpenAPI Specification](https://swagger.io/specification/)

---

**Ready to build?** Check out the [examples](./examples/) directory! 🚀
