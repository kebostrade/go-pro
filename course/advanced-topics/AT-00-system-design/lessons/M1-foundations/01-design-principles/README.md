# SD-01: System Design Principles

Learn the fundamental principles that guide effective system design.

## Overview

This lesson covers the core principles that every software engineer needs to understand before designing systems. These principles form the foundation for all subsequent system design decisions.

## Learning Objectives

- Understand SOLID principles in system design
- Learn the DRY and KISS principles
- Apply separation of concerns
- Make informed trade-offs

## Key Design Principles

### 1. SOLID Principles

```
Single Responsibility: Each component has one reason to change
Open/Closed: Open for extension, closed for modification
Liskov Substitution: Subtypes must be substitutable for their base types
Interface Segregation: Many small interfaces > one large interface
Dependency Inversion: Depend on abstractions, not concretions
```

### 2. DRY (Don't Repeat Yourself)

Avoid code duplication by extracting common functionality:

```go
// Bad: Duplicated validation logic
func CreateUser(req *CreateUserRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    // ...
}

func UpdateUser(req *UpdateUserRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    // ...
}

// Good: Shared validation
func ValidateUser(user *User) error {
    if user.Name == "" {
        return errors.New("name is required")
    }
    return nil
}
```

### 3. KISS (Keep It Simple, Stupid)

Prefer simple solutions that work over complex ones:

```go
// Complex: Over-engineered solution
type UserService struct {
    validator *UserValidator
    normalizer *UserNormalizer
    transformer *UserTransformer
    // ... 10 more components
}

// Simple: Focused service
type UserService struct {
    repo UserRepository
    cache CacheService
}
```

### 4. Separation of Concerns

Each layer has distinct responsibilities:

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

### 5. YAGNI (You Aren't Gonna Need It)

Don't build features until they're necessary:

- Avoid speculative generalization
- Implement what's needed now
- Refactor when requirements change

## Design Trade-offs

| Principle | When to Apply | When to Avoid |
|-----------|---------------|---------------|
| Normalization | Write-heavy workloads | Read-heavy with complex queries |
| Denormalization | Read-heavy workloads | Frequently updated data |
| Caching | Expensive computations | Frequently changing data |
| Async Processing | Long-running operations | Need immediate feedback |

## Go-Specific Considerations

Go favors simplicity:

```go
// Go style: Simple interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose interfaces when needed
type ReadWriter interface {
    Reader
    Writer
}
```

## Examples

See `examples/` directory for:
- `solid_principles.go` - SOLID implementations
- `separation_of_concerns.go` - Layer separation
- `tradeoffs.go` - Design decision examples

## Exercises

See `exercises/problems.md` for hands-on practice.

## Quiz

Test your knowledge with `quiz.md`.

## Summary

- Apply SOLID principles for maintainable code
- Keep systems simple (KISS) until complexity is needed
- Separate concerns for better testability
- Make trade-offs consciously, not arbitrarily

## Next Steps

Continue to [SD-02: Requirements Gathering](02-requirements-gathering/README.md)
