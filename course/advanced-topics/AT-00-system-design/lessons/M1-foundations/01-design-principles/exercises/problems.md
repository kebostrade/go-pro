# Exercise: System Design Principles

## Problem 1: Identifying Violations

Identify which SOLID principle is violated in each case:

### Case A:
```go
type UserHandler struct {
    db *sql.DB
    cache *redis.Client
    email EmailService
    logger Logger
    metrics Metrics
    // 50 more dependencies
}
```

**Violated Principle:**
```
Answer: Interface Segregation - the handler has too many dependencies
```

### Case B:
```go
func ProcessPayment(order *Order) error {
    // Validate order
    // Calculate total
    // Charge card
    // Update inventory
    // Send email
    // Generate invoice
    // Log transaction
    // Notify warehouse
}
```

**Violated Principle:**
```
Answer: Single Responsibility - function does too many things
```

---

## Problem 2: Applying Separation of Concerns

You have a `UserService` that handles:
- User registration
- Password hashing
- Email validation
- Database operations
- Session management
- Sending welcome emails

### Task 1: Separate into distinct services:

| Concern | Service Name |
|---------|--------------|
| User registration | |
| Password hashing | |
| Email validation | |
| Database operations | |
| Session management | |
| Sending emails | |

### Task 2: Define interfaces for each:

```go
// Define interfaces for the separated concerns
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hash, password string) bool
}

type EmailValidator interface {
    // Add methods
}

type SessionManager interface {
    // Add methods
}
```

---

## Problem 3: Dependency Inversion

Convert this tightly-coupled code to use dependency injection:

### Before:
```go
type OrderService struct {
    postgresRepo *PostgresOrderRepository
    redisCache *RedisCache
    smtpEmail *SMTPService
}
```

### After (with interfaces):
```go
// Define interfaces
type OrderRepository interface {
    Create(ctx context.Context, order *Order) error
    GetByID(ctx context.Context, id string) (*Order, error)
}

type CacheService interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type EmailService interface {
    Send(ctx context.Context, to string, subject string, body string) error
}

// Refactored service
type OrderService struct {
    repo OrderRepository
    cache CacheService
    email EmailService
}
```

---

## Problem 4: Trade-off Analysis

For each scenario, identify the trade-off:

| Scenario | Trade-off |
|----------|-----------|
| Using a monolith vs microservices | |
| Normalized vs denormalized database | |
| Synchronous vs asynchronous processing | |
| Strong consistency vs eventual consistency | |
| In-memory cache vs disk storage | |

---

## Problem 5: KISS vs YAGNI

For each feature request, decide if you should implement it now or later:

1. **"We need a feature to export data to XML, CSV, and PDF"**
   - Current need: CSV only
   
2. **"The system should support multiple database backends (PostgreSQL, MySQL, MongoDB)"**
   - Current requirement: PostgreSQL only
   
3. **"We need authentication with OAuth2, SAML, and LDAP"**
   - Current requirement: Username/password only

Decide: KISS or YAGNI?
