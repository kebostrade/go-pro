# 🎨 Design Patterns in Go

A comprehensive guide to implementing classic design patterns in Go with practical examples, tests, and real-world use cases.

## 📚 Table of Contents

- [Overview](#overview)
- [Pattern Categories](#pattern-categories)
- [Quick Start](#quick-start)
- [Patterns Implemented](#patterns-implemented)
- [Usage Examples](#usage-examples)
- [Testing](#testing)
- [Learning Resources](#learning-resources)

## 🎯 Overview

This project demonstrates **23 classic design patterns** adapted for Go's idioms and best practices. Each pattern includes:

- ✅ Clear explanation of purpose and use cases
- ✅ Go-specific implementation details
- ✅ Comprehensive tests with examples
- ✅ Real-world usage scenarios
- ✅ Performance benchmarks

## 📦 Pattern Categories

### 🏗️ Creational Patterns (Object Creation)
Create objects in a manner suitable to the situation.

| Pattern | Purpose | Use Case |
|---------|---------|----------|
| **Singleton** | Ensure only one instance exists | Database connections, Config managers |
| **Factory** | Create objects without specifying exact class | Payment processors, Notifications |
| **Builder** | Construct complex objects step by step | HTTP requests, SQL queries |
| **Prototype** | Clone existing objects | Document templates, Game objects |

### 🔧 Structural Patterns (Object Composition)
Compose objects to form larger structures.

| Pattern | Purpose | Use Case |
|---------|---------|----------|
| **Adapter** | Convert interface to another interface | Third-party library integration |
| **Decorator** | Add responsibilities dynamically | Logging, Caching, Middleware |
| **Proxy** | Provide surrogate for another object | Lazy loading, Access control |
| **Facade** | Provide simplified interface | Complex subsystem access |

### 🎭 Behavioral Patterns (Object Interaction)
Define how objects communicate and distribute responsibility.

| Pattern | Purpose | Use Case |
|---------|---------|----------|
| **Strategy** | Define family of algorithms | Payment methods, Sorting |
| **Observer** | Notify dependents of state changes | Event systems, Pub/Sub |
| **Command** | Encapsulate requests as objects | Undo/Redo, Task queues |
| **Chain of Responsibility** | Pass request along chain | Middleware, Logging |

## 🚀 Quick Start

### Installation

```bash
# Navigate to project
cd basic/projects/design-patterns

# Download dependencies
make deps

# Run all examples
make run-all-examples

# Run tests
make test

# Run benchmarks
make bench
```

### Run Individual Patterns

```bash
# Singleton
make run-singleton

# Factory
make run-factory

# Builder
make run-builder

# Strategy
make run-strategy

# Observer
make run-observer
```

## 💻 Patterns Implemented

### 1. Singleton Pattern

**Purpose**: Ensure a class has only one instance.

```go
// Thread-safe singleton using sync.Once
db := creational.GetDatabase()
db.Connect()

// Configuration singleton
config := creational.GetConfig()
config.Set("theme", "dark")
```

**Key Features**:
- Thread-safe with `sync.Once`
- Lazy initialization
- Global access point

### 2. Factory Pattern

**Purpose**: Create objects without specifying exact class.

```go
// Create different notification types
email, _ := creational.NotificationFactory("email")
sms, _ := creational.NotificationFactory("sms")
push, _ := creational.NotificationFactory("push")

email.Send("Welcome!")
sms.Send("Verification code: 123456")
push.Send("New message")
```

**Key Features**:
- Interface-based design
- Extensible with new types
- Centralized object creation

### 3. Builder Pattern

**Purpose**: Construct complex objects step by step.

```go
// Build HTTP request
request := creational.NewHTTPRequestBuilder().
    Method("POST").
    URL("https://api.example.com/users").
    Header("Content-Type", "application/json").
    Body(`{"name":"John"}`).
    Timeout(60).
    Build()

// Build SQL query
query := creational.NewSQLQueryBuilder().
    Select("id", "name", "email").
    From("users").
    Where("age > 18").
    OrderBy("created_at DESC").
    Limit(10).
    Build()
```

**Key Features**:
- Fluent interface (method chaining)
- Immutable final object
- Readable construction code

### 4. Adapter Pattern

**Purpose**: Convert interface to another interface.

```go
// Adapt different media players to common interface
var player structural.MediaPlayer

player = structural.NewVLCAdapter()
player.Play("movie.vlc")

player = structural.NewMP4Adapter()
player.Play("video.mp4")
```

**Key Features**:
- Interface compatibility
- Wrapper pattern
- Third-party integration

### 5. Decorator Pattern

**Purpose**: Add responsibilities dynamically.

```go
// Build coffee with decorators
coffee := &structural.SimpleCoffee{}
coffeeWithMilk := structural.NewMilkDecorator(coffee)
coffeeWithMilkAndSugar := structural.NewSugarDecorator(coffeeWithMilk)
fancyCoffee := structural.NewWhipDecorator(coffeeWithMilkAndSugar)

fmt.Printf("%s: $%.2f\n", fancyCoffee.Description(), fancyCoffee.Cost())
// Output: Simple Coffee, Milk, Sugar, Whipped Cream: $3.40
```

**Key Features**:
- Composition over inheritance
- Runtime behavior modification
- Stackable decorators

### 6. Strategy Pattern

**Purpose**: Define family of interchangeable algorithms.

```go
// Payment strategies
payment := behavioral.NewPaymentContext(&behavioral.CreditCardStrategy{
    CardNumber: "1234567890123456",
})
payment.ExecutePayment(100.00)

// Change strategy at runtime
payment.SetStrategy(&behavioral.PayPalStrategy{
    Email: "user@example.com",
})
payment.ExecutePayment(50.00)
```

**Key Features**:
- Algorithm encapsulation
- Runtime strategy switching
- Open/Closed principle

### 7. Observer Pattern

**Purpose**: Notify dependents of state changes.

```go
// Create event manager
eventManager := behavioral.NewEventManager()

// Attach observers
eventManager.Attach(&behavioral.EmailObserver{ID: "email-1"})
eventManager.Attach(&behavioral.SMSObserver{ID: "sms-1"})
eventManager.Attach(&behavioral.LogObserver{ID: "log-1"})

// Notify all observers
eventManager.Notify("user.registered", map[string]string{
    "username": "johndoe",
})
```

**Key Features**:
- One-to-many dependency
- Loose coupling
- Event-driven architecture

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench

# Verbose output
make test-verbose
```

### Test Coverage

All patterns include comprehensive tests:
- Unit tests for each pattern
- Concurrency tests (where applicable)
- Benchmark tests for performance
- Example-based tests

## 📊 Benchmarks

```bash
make bench
```

Example output:
```
BenchmarkGetDatabase-8          1000000000    0.25 ns/op
BenchmarkGetConfig-8            1000000000    0.24 ns/op
BenchmarkFactory-8              5000000       250 ns/op
BenchmarkBuilder-8              2000000       600 ns/op
```

## 🎓 Learning Outcomes

By studying this project, you'll learn:

1. **Design Principles**
   - SOLID principles
   - Composition over inheritance
   - Program to interfaces

2. **Go-Specific Patterns**
   - Interface-based design
   - Functional options
   - Channel-based patterns

3. **Best Practices**
   - Thread-safe implementations
   - Idiomatic Go code
   - Performance optimization

## 📚 Resources

- **Tutorial**: [Tutorial 14 in TUTORIALS.md](../../docs/TUTORIALS.md)
- **Go Design Patterns**: https://refactoring.guru/design-patterns/go
- **Gang of Four Book**: Design Patterns: Elements of Reusable Object-Oriented Software
- **Go Proverbs**: https://go-proverbs.github.io/

## 🏆 Pattern Selection Guide

**When to use each pattern:**

- **Singleton**: Global state, resource pooling
- **Factory**: Multiple product types, runtime selection
- **Builder**: Complex object construction, many parameters
- **Adapter**: Interface incompatibility, legacy integration
- **Decorator**: Dynamic behavior addition, middleware
- **Strategy**: Algorithm selection, runtime behavior change
- **Observer**: Event handling, state change notifications

---

**Built with ❤️ using Go and Design Patterns**

