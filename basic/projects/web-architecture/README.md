# 🏗️ Web Architecture with Go

A comprehensive guide to building scalable web applications in Go using Clean Architecture principles, RESTful APIs, and production-ready patterns.

## 📚 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Design Patterns](#design-patterns)
- [Middleware](#middleware)
- [Testing](#testing)
- [Deployment](#deployment)

## 🎯 Overview

This project demonstrates **production-ready web architecture** in Go with:

- ✅ **Clean Architecture** - Separation of concerns
- ✅ **RESTful API** - Standard HTTP methods and status codes
- ✅ **Middleware Chain** - Logging, auth, recovery, CORS
- ✅ **Repository Pattern** - Data access abstraction
- ✅ **Service Layer** - Business logic isolation
- ✅ **JWT Authentication** - Secure token-based auth
- ✅ **Graceful Shutdown** - Proper server lifecycle management
- ✅ **Error Handling** - Standardized error responses

## 🏛️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│  PRESENTATION LAYER (HTTP Handlers)                         │
│  • User Handler                                             │
│  • Product Handler                                          │
│  • Request/Response DTOs                                    │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  BUSINESS LOGIC LAYER (Services)                            │
│  • User Service (Registration, Login, JWT)                 │
│  • Product Service (CRUD operations)                        │
│  • Validation & Business Rules                             │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  DATA ACCESS LAYER (Repositories)                           │
│  • User Repository Interface                                │
│  • Product Repository Interface                             │
│  • In-Memory Implementation                                 │
│  • PostgreSQL Implementation (optional)                     │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  DOMAIN LAYER (Models)                                      │
│  • User Model                                               │
│  • Product Model                                            │
│  • Request/Response DTOs                                    │
└─────────────────────────────────────────────────────────────┘
```

### Request Flow

```
HTTP Request
    ↓
Middleware Chain (Logging → Recovery → CORS → Auth)
    ↓
Router (chi)
    ↓
Handler (Parse & Validate)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database/Storage
    ↓
Response (JSON)
```

## 🚀 Quick Start

### Installation

```bash
# Navigate to project
cd basic/projects/web-architecture

# Download dependencies
go mod tidy

# Run server
make run
```

### First API Call

```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "username": "alice",
    "password": "password123",
    "first_name": "Alice",
    "last_name": "Smith"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'

# Use the token from login response
export TOKEN="your-jwt-token-here"

# Create a product (requires authentication)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "stock": 10,
    "category": "Electronics"
  }'

# List products
curl -X GET "http://localhost:8080/api/v1/products?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```

## 📡 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/register` | Register new user |
| POST | `/api/v1/login` | Login and get JWT token |
| GET | `/health` | Health check |

### Protected Endpoints (Require JWT)

#### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List all users |
| GET | `/api/v1/users/{id}` | Get user by ID |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Delete user |

#### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List products (with pagination) |
| GET | `/api/v1/products?category=Electronics` | Filter by category |
| GET | `/api/v1/products/{id}` | Get product by ID |
| POST | `/api/v1/products` | Create product |
| PUT | `/api/v1/products/{id}` | Update product |
| DELETE | `/api/v1/products/{id}` | Delete product |

## 📁 Project Structure

```
web-architecture/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handler/                 # HTTP handlers (controllers)
│   │   ├── user_handler.go
│   │   └── product_handler.go
│   ├── service/                 # Business logic
│   │   ├── user_service.go
│   │   └── product_service.go
│   ├── repository/              # Data access layer
│   │   ├── repository.go        # Interfaces
│   │   ├── memory_user.go       # In-memory implementation
│   │   └── memory_product.go
│   ├── model/                   # Domain models
│   │   ├── user.go
│   │   └── product.go
│   └── middleware/              # HTTP middleware
│       ├── auth.go              # JWT authentication
│       ├── logging.go           # Request logging
│       └── recovery.go          # Panic recovery
├── pkg/
│   └── response/                # Standardized responses
│       └── response.go
├── migrations/                  # Database migrations
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## 🎨 Design Patterns

### 1. Repository Pattern

**Abstracts data access logic**

```go
// Interface
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id int64) (*model.User, error)
    // ...
}

// Implementation
type MemoryUserRepository struct {
    users map[int64]*model.User
}
```

**Benefits:**
- Swap implementations (memory → PostgreSQL)
- Easy testing with mocks
- Decouples business logic from data storage

### 2. Service Layer Pattern

**Encapsulates business logic**

```go
type UserService struct {
    repo      UserRepository
    jwtSecret []byte
}

func (s *UserService) Register(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validation
    // Password hashing
    // Business rules
    // Call repository
}
```

**Benefits:**
- Reusable business logic
- Testable without HTTP layer
- Clear separation of concerns

### 3. Dependency Injection

**Inject dependencies through constructors**

```go
func NewUserService(repo UserRepository, jwtSecret string) *UserService {
    return &UserService{
        repo:      repo,
        jwtSecret: []byte(jwtSecret),
    }
}
```

**Benefits:**
- Loose coupling
- Easy testing
- Flexible configuration

## 🔐 Middleware

### Authentication Middleware

```go
func Auth(userService *UserService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            // Validate JWT token
            // Add user claims to context
            // Call next handler
        })
    }
}
```

### Logging Middleware

```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Call next handler
        // Log request details
    })
}
```

### Recovery Middleware

```go
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log panic
                // Return 500 error
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/service/...
```

## 🚢 Deployment

### Docker

```bash
# Build image
make docker-build

# Run with Docker Compose
make docker-run
```

### Environment Variables

```bash
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
DATABASE_URL=postgres://user:pass@localhost:5432/dbname
```

## 📊 Best Practices

### 1. Error Handling

```go
// Service layer
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// Handler layer
if err != nil {
    response.InternalServerError(w, "Failed to create user")
    return
}
```

### 2. Context Propagation

```go
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    // Pass context to repository
    return s.repo.GetByID(ctx, id)
}
```

### 3. Graceful Shutdown

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      r,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
}

// Handle shutdown signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

---

**Built with ❤️ using Go's powerful standard library and clean architecture principles**

