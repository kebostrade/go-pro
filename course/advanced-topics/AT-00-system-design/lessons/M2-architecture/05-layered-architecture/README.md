# SD-05: Layered Architecture

Learn the classic layered architecture pattern and its implementation in Go.

## Overview

Layered architecture is one of the most common architectural patterns. It separates concerns into distinct layers, each with specific responsibilities.

## Learning Objectives

- Understand layered architecture principles
- Implement handlers, services, and repositories
- Define clear boundaries between layers
- Apply dependency injection

## Layer Structure

```
┌─────────────────────────────────────┐
│         Delivery Layer               │  ← HTTP/gRPC handlers
├─────────────────────────────────────┤
│          Service Layer              │  ← Business logic
├─────────────────────────────────────┤
│        Repository Layer             │  ← Data access
├─────────────────────────────────────┤
│         Infrastructure              │  ← DB, cache, APIs
└─────────────────────────────────────┘
```

## Implementation in Go

### 1. Domain/Entity Layer

```go
// internal/domain/user.go
type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
}

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### 2. Service Layer

```go
// internal/service/user_service.go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id string) (*User, error)
}

type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Business logic here
    user := &User{
        ID:    uuid.New().String(),
        Email: req.Email,
        Name:  req.Name,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 3. Handler Layer

```go
// internal/handler/user_handler.go
type UserHandler struct {
    service UserService
}

func NewUserHandler(service UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.service.CreateUser(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}
```

### 4. Wiring It Up

```go
// cmd/server/main.go
func main() {
    // Initialize dependencies
    db := connectDatabase()
    repo := NewPostgresUserRepository(db)
    service := NewUserService(repo)
    handler := NewUserHandler(service)
    
    // Setup routes
    http.HandleFunc("/users", handler.CreateUser)
    
    http.ListenAndServe(":8080", nil)
}
```

## Advantages

- **Separation of concerns**: Each layer has clear responsibility
- **Testability**: Mock layers independently
- **Maintainability**: Changes localized to one layer
- **Reusability**: Services can be used by different handlers

## Disadvantages

- **Performance**: Multiple layer traversals
- **Complexity**: More code to write initially
- **Overhead**: Abstraction can hide performance issues

## When to Use

- Medium to large applications
- Teams that need clear separation
- Applications with complex business logic

## Examples

See `examples/` directory for:
- `layered_architecture.go` - Full implementation
- `dependency_injection.go` - DI patterns

## Exercises

See `exercises/problems.md` for hands-on practice.

## Quiz

Test your knowledge with `quiz.md`.

## Summary

- Layer architecture separates concerns
- Each layer has specific responsibility
- Use interfaces for loose coupling
- Apply dependency injection

## Next Steps

Continue to [SD-06: Clean Architecture](06-clean-architecture/README.md)
