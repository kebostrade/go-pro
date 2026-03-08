# System Design with Golang

Design scalable, maintainable systems using Go principles and patterns.

## Learning Objectives

- Apply SOLID principles in Go
- Design for scalability and maintainability
- Choose appropriate architectural patterns
- Implement clean architecture
- Evaluate trade-offs in system design

## Theory

### Go Design Philosophy

Go favors simplicity over complexity:

```go
type UserService interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &userService{repo: repo}
}
```

### Clean Architecture Layers

```
┌─────────────────────────────────────┐
│           Delivery Layer            │  ← HTTP/gRPC/CLI
├─────────────────────────────────────┤
│           Use Case Layer            │  ← Business Logic
├─────────────────────────────────────┤
│           Entity Layer              │  ← Domain Models
├─────────────────────────────────────┤
│           Data Layer                │  ← Repository/DB
└─────────────────────────────────────┘
```

### Dependency Inversion

```go
type OrderProcessor struct {
    paymentGateway PaymentGateway
    inventory      InventoryService
    notifier       NotificationService
}

func NewOrderProcessor(
    pg PaymentGateway,
    inv InventoryService,
    notif NotificationService,
) *OrderProcessor {
    return &OrderProcessor{
        paymentGateway: pg,
        inventory:      inv,
        notifier:       notif,
    }
}
```

## Real-World Applications

### GO-PRO Backend Pattern

```go
type Handler struct {
    service CourseService
}

type CourseService interface {
    Create(ctx context.Context, req *CreateCourseReq) (*Course, error)
}

type courseService struct {
    repo   CourseRepository
    cache  CacheService
    events EventEmitter
}
```

## Performance Tips

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return bytes.NewBuffer(make([]byte, 0, 1024))
    },
}

func processLargeData(data []byte) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
}
```

## Exercises

1. Design a URL shortener architecture
2. Implement repository pattern with interfaces
3. Create a plugin-based system
4. Design for horizontal scaling

## Validation

```bash
cd exercises
go test -v ./...
```

## Key Takeaways

- Prefer composition over inheritance
- Design interfaces at boundaries
- Keep packages cohesive
- Optimize for readability first
- Use dependency injection

## Next Steps

**[AT-01: RESTful APIs](../AT-01-rest-apis/README.md)**

---

Design for simplicity, scale through patterns. 🏗️
