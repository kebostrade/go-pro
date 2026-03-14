# 01 - System Design with Golang

Design scalable, maintainable systems using Go best practices and architectural patterns.

## 📚 Overview

This topic covers essential system design principles and patterns specifically applied to Go applications. You'll learn how to architect systems that are scalable, maintainable, and performant.

## 🎯 Learning Objectives

- Understand Go-specific architectural patterns
- Learn layer architecture and clean design principles
- Master concurrency patterns for scalable systems
- Implement caching strategies and data management
- Design fault-tolerant and observable systems

## 🏗️ Architecture Patterns

### 1. Layered Architecture

```
┌─────────────────────────────────────┐
│         API Layer (handlers)         │  ← HTTP/gRPC endpoints
├─────────────────────────────────────┤
│       Business Layer (services)      │  ← Business logic
├─────────────────────────────────────┤
│      Data Layer (repositories)       │  ← Data access
├─────────────────────────────────────┤
│         Infrastructure               │  ← DB, cache, external APIs
└─────────────────────────────────────┘
```

### 2. Clean Architecture

```
┌──────────────────────────────────────┐
│           Frameworks & Drivers        │  ← Go stdlib, external libs
├──────────────────────────────────────┤
│            Interface Adapters         │  ← Controllers, Presenters
├──────────────────────────────────────┤
│              Use Cases               │  ← Application business rules
├──────────────────────────────────────┤
│               Entities               │  ← Enterprise business rules
└──────────────────────────────────────┘
```

### 3. Microservices Architecture

```
                    ┌─────────────┐
                    │   Gateway   │
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
      ┌─────▼─────┐ ┌─────▼─────┐ ┌─────▼─────┐
      │  Service  │ │  Service  │ │  Service  │
      │     A     │ │     B     │ │     C     │
      └─────┬─────┘ └─────┬─────┘ └─────┬─────┘
            │             │             │
            └──────────────┼──────────────┘
                           │
                  ┌────────▼────────┐
                  │   Data Layer    │
                  │ (DBs, Cache)    │
                  └─────────────────┘
```

## 📁 Example Project Structure

```
system-design-example/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handler/                    # HTTP handlers
│   │   ├── user.go
│   │   └── middleware.go
│   ├── service/                    # Business logic
│   │   ├── user_service.go
│   │   └── auth_service.go
│   ├── repository/                 # Data access
│   │   ├── interfaces.go           # Repository interfaces
│   │   ├── user_repository.go      # Implementation
│   │   └── cache_repository.go     # Cache layer
│   ├── domain/                     # Business entities
│   │   ├── user.go
│   │   └── errors.go
│   └── config/                     # Configuration
│       └── config.go
├── pkg/                            # Public packages
│   ├── logger/
│   └── metrics/
├── go.mod
├── go.sum
└── README.md
```

## 🔑 Key Design Principles

### 1. Separation of Concerns
Each layer has a single responsibility:
- **Handlers**: Handle HTTP requests/responses only
- **Services**: Contain business logic
- **Repositories**: Handle data access

### 2. Dependency Inversion
Depend on abstractions (interfaces), not concrete implementations:

```go
// Good - depends on interface
type UserService struct {
    repo UserRepository  // Interface
}

// Bad - depends on concrete implementation
type UserService struct {
    repo *PostgresUserRepository  // Concrete type
}
```

### 3. Interface Segregation
Keep interfaces small and focused:

```go
// Good - focused interface
type UserReader interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

// Bad - bloated interface
type UserRepository interface {
    UserReader
    UserWriter
    DeleteAll() error
    ExportToCSV() string
    SendEmail() error
    // ... 20 more methods
}
```

### 4. Configuration Management

Use structured configuration with validation:

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Cache    CacheConfig    `mapstructure:"cache"`
}

type ServerConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`
    ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func Load(configPath string) (*Config, error) {
    cfg := &Config{}
    // Load and validate configuration
    return cfg, nil
}
```

## 🚀 Implementation Examples

### Example 1: Clean Layered Architecture

See `examples/clean_architecture.go` for a complete implementation.

### Example 2: Repository Pattern with Caching

See `examples/repository_pattern.go` for caching strategies.

### Example 3: Concurrent Processing Patterns

See `examples/concurrency_patterns.go` for worker pools and pipelines.

### Example 4: Circuit Breaker Pattern

See `examples/circuit_breaker.go` for fault tolerance.

### Example 5: Observability and Metrics

See `examples/observability.go` for logging, metrics, and tracing.

## 📊 Design Trade-offs

### Monolith vs Microservices

| Factor | Monolith | Microservices |
|--------|----------|---------------|
| **Development Speed** | ✅ Faster initially | ⚠️ Slower initially |
| **Deployment** | ✅ Simple | ⚠️ Complex orchestration |
| **Scalability** | ⚠️ Scale entire app | ✅ Scale independently |
| **Complexity** | ✅ Lower | ⚠️ Higher (distributed systems) |
| **Team Size** | ✅ Small teams | ✅ Large teams |
| **Technology Mix** | ❌ One stack | ✅ Polyglot |

### SQL vs NoSQL

| Factor | SQL | NoSQL |
|--------|-----|-------|
| **Schema** | Rigid, predefined | Flexible, dynamic |
| **Scaling** | Vertical (mostly) | Horizontal |
| **Consistency** | ACID guarantees | BASE / Tunable |
| **Query Complexity** | Complex joins | Simple lookups |
| **Use Case** | Structured data | Unstructured, high velocity |

### Caching Strategies

| Strategy | Use Case | Pros | Cons |
|----------|----------|------|------|
| **Cache-Aside** | Read-heavy | Simple, flexible | Cache stampede |
| **Write-Through** | Data consistency | Always fresh | Slower writes |
| **Write-Behind** | Write-heavy | Fast writes | Data loss risk |
| **Refresh-Ahead** | Predictable access | No misses | Wastes resources |

## 🧪 Testing Strategy

```go
// Unit test - mock repository
func TestUserService_Create(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    user := &User{Name: "John"}
    err := service.Create(context.Background(), user)

    assert.NoError(t, err)
    mockRepo.AssertCalled(t, "Create", user)
}

// Integration test - real database
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db := setupTestDB()
    repo := NewPostgresUserRepository(db)

    user := &User{Name: "John"}
    err := repo.Create(context.Background(), user)

    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

## 📈 Performance Optimization

### 1. Database Optimization
```go
// Use connection pooling
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)

// Use prepared statements
stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = ?")
defer stmt.Close()

// Batch operations
tx, _ := db.BeginTx(ctx, nil)
stmt, _ := tx.Prepare("INSERT INTO users (name) VALUES (?)")
for _, user := range users {
    stmt.Exec(user.Name)
}
tx.Commit()
```

### 2. Caching Strategy
```go
// Multi-level caching
type Cache struct {
    l1 *sync.Map          // In-memory cache
    l2 *redis.Client      // Redis cache
    db Database           // Database
}

func (c *Cache) Get(key string) (interface{}, error) {
    // L1: Check in-memory
    if val, ok := c.l1.Load(key); ok {
        return val, nil
    }

    // L2: Check Redis
    val, err := c.l2.Get(ctx, key).Result()
    if err == nil {
        c.l1.Store(key, val)
        return val, nil
    }

    // L3: Query database
    val, err = c.db.Query(key)
    if err != nil {
        return nil, err
    }

    // Populate cache
    c.l1.Store(key, val)
    c.l2.Set(ctx, key, val, time.Hour)
    return val, nil
}
```

### 3. Concurrent Processing
```go
// Worker pool for concurrent processing
func ProcessItems(items []Item, workers int) {
    itemsCh := make(chan Item, len(items))
    resultsCh := make(chan Result, len(items))

    // Start workers
    for i := 0; i < workers; i++ {
        go worker(itemsCh, resultsCh)
    }

    // Send work
    for _, item := range items {
        itemsCh <- item
    }
    close(itemsCh)

    // Collect results
    for i := 0; i < len(items); i++ {
        result := <-resultsCh
        // Handle result
    }
}
```

## 🔒 Security Considerations

1. **Input Validation**: Validate all inputs at handler layer
2. **Authentication/Authorization**: Implement in service layer
3. **SQL Injection**: Use parameterized queries
4. **Secrets Management**: Use environment variables or secret managers
5. **Rate Limiting**: Implement at handler layer
6. **TLS**: Always use HTTPS in production

## 📝 Best Practices Summary

✅ **DO**:
- Use clear layer separation
- Depend on interfaces, not implementations
- Keep functions small and focused
- Handle errors explicitly
- Use context for cancellation and timeouts
- Implement observability (logs, metrics, traces)
- Write tests for all business logic
- Use connection pooling
- Implement graceful shutdown

❌ **DON'T**:
- Mix business logic with handlers
- Create god objects (large structs with many responsibilities)
- Ignore errors
- Use global state
- Hardcode configuration
- Forget to close resources
- Block goroutines indefinitely
- Over-engineer simple problems

## 🚦 Next Steps

After mastering system design:
1. → **[02 - Building RESTful APIs with Go](../02-rest-api/)**
2. → **[06 - Building Microservices with Go and Docker](../06-microservices-docker/)**

## 📖 Further Reading

- [Go Cloud Development Kit](https://gocloud.dev/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Domain-Driven Design in Go](https://github.com/ThreeDotsLabs/watermill)
- [Microservices Patterns](https://microservices.io/patterns/)

---

**Ready to build?** Check out the [examples](./examples/) directory for working code! 🚀
