# Testing and Debugging in Go

Comprehensive guide to testing strategies, debugging techniques, and best practices in Go.

## Table of Contents

1. [Testing Fundamentals](#testing-fundamentals)
2. [Testing Patterns](#testing-patterns)
3. [Integration Testing](#integration-testing)
4. [Benchmarking and Profiling](#benchmarking-and-profiling)
5. [Advanced Testing Techniques](#advanced-testing-techniques)
6. [Debugging Strategies](#debugging-strategies)
7. [Best Practices](#best-practices)

## Testing Fundamentals

### Go Testing Philosophy

Go's testing philosophy emphasizes:
- **Simplicity**: Tests are just Go functions
- **Speed**: Tests should run quickly
- **Integration**: Tests live alongside code
- **Visibility**: Use `go test` to discover and run tests

### Basic Test Structure

```go
func TestFunctionName(t *testing.T) {
    // t.Fatal() or t.Fatalf() - Stop test immediately
    // t.Error() or t.Errorf() - Mark as failed but continue
    // t.Log() or t.Logf() - Logging (shown with -v flag)
    // t.Skip() or t.Skipf() - Skip test (with reason)
}
```

### Running Tests

```bash
# Run all tests in current directory
go test

# Run with verbose output
go test -v

# Run specific test
go test -run TestFunctionName

# Run tests in all subdirectories
go test ./...

# Run with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Race detection
go test -race

# Run with specific timeout
go test -timeout 30s

# Run tests in parallel
go test -parallel 4
```

## Testing Patterns

### 1. Table-Driven Tests

**When to use**: Testing the same function with multiple inputs and expected outputs.

**Benefits**:
- Easy to add new test cases
- Clear test data separation
- Reduced code duplication

**Example**: See `examples/unit_tests.go` for complete implementation.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed numbers", -2, 3, 1},
        {"zeros", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 2. Test Setup and Teardown

**When to use**: When tests need shared setup or cleanup logic.

**Example**: See `examples/unit_tests.go` for complete implementation.

```go
func TestMain(m *testing.M) {
    fmt.Println("Setting up tests")
    // Setup: database connection, test fixtures, etc.

    code := m.Run() // Run all tests

    fmt.Println("Tearing down tests")
    // Teardown: close connections, cleanup, etc.

    os.Exit(code)
}

func TestWithSetup(t *testing.T) {
    // Setup for this specific test
    db := setupTestDB(t)
    defer db.Close()

    // Test code here
}

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper() // Marks this as a test helper
    // Setup code
    return db
}
```

### 3. Test Helpers and Helpers

**When to use**: To avoid duplication and improve test readability.

**Example**: See `examples/testing_tools.go` for complete implementation.

```go
// Helper function to assert equality
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// Helper to check error
func assertError(t *testing.T, err error, want bool) {
    t.Helper()
    if (err != nil) != want {
        t.Errorf("error = %v, want error? %v", err, want)
    }
}
```

### 4. Mocking and Interfaces

**When to use**: When testing code that depends on external services or complex dependencies.

**Example**: See `examples/unit_tests.go` for complete implementation.

```go
// Define interface
type Database interface {
    GetUser(id int) (*User, error)
}

// Mock implementation
type MockDatabase struct {
    mock.Mock
}

func (m *MockDatabase) GetUser(id int) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Test with mock
func TestUserService(t *testing.T) {
    mockDB := new(MockDatabase)
    mockDB.On("GetUser", 1).Return(&User{ID: 1, Name: "John"}, nil)

    service := NewUserService(mockDB)
    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, user.Name, "John")
    mockDB.AssertExpectations(t)
}
```

### 5. HTTP Testing with httptest

**When to use**: Testing HTTP handlers without starting a real server.

**Example**: See `examples/integration_tests.go` for complete implementation.

```go
func TestHandler(t *testing.T) {
    handler := MyHandler()

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }

    expected := `{"status":"ok"}`
    if w.Body.String() != expected {
        t.Errorf("expected body %s, got %s", expected, w.Body.String())
    }
}
```

## Integration Testing

### Database Integration Testing

**When to use**: Testing database interactions and queries.

**Example**: See `examples/integration_tests.go` for complete implementation.

**Best Practices**:
- Use separate test database
- Wrap tests in transactions and rollback
- Use test fixtures for consistent data
- Clean up after tests

```go
func TestUserRepository(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewUserRepository(db)

    // Create
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(user)
    assert.NoError(t, err)

    // Read
    found, err := repo.FindByID(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, found.Name, "John")
}
```

### API Integration Testing

**When to use**: Testing API endpoints and HTTP handlers.

**Example**: See `examples/integration_tests.go` for complete implementation.

**Best Practices**:
- Use `httptest` for HTTP testing
- Test all status codes and edge cases
- Validate response bodies
- Test authentication and authorization

## Benchmarking and Profiling

### Writing Benchmarks

**When to use**: Measuring and optimizing performance.

**Example**: See `examples/benchmarking.go` for complete implementation.

```go
func BenchmarkFunction(b *testing.B) {
    // Setup code that doesn't count toward benchmark time
    data := generateTestData()

    b.ResetTimer() // Reset timer after setup
    for i := 0; i < b.N; i++ {
        Function(data)
    }
}
```

**Running Benchmarks**:

```bash
# Run benchmarks
go test -bench=.

# Run with memory allocation stats
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkFunction

# Run benchmarks multiple times for accuracy
go test -bench=. -count=10

# Control CPU usage
go test -bench=. -cpu=1,2,4

# Run benchmark for specific time
go test -bench=. -benchtime=10s
```

### CPU Profiling

**When to use**: Identifying CPU bottlenecks and hotspots.

**Example**: See `examples/benchmarking.go` for complete implementation.

```bash
# Generate CPU profile
go test -cpuprofile=cpu.prof -bench=.

# Analyze with pprof tool
go tool pprof cpu.prof

# In pprof interactive mode:
# (pprof) top10        # Top 10 functions
# (pprof) list Function # Disassemble function
# (pprof) web          # Generate graph (requires graphviz)
# (pprof) pdf          # Generate PDF
```

### Memory Profiling

**When to use**: Identifying memory leaks and allocation patterns.

**Example**: See `examples/benchmarking.go` for complete implementation.

```bash
# Generate memory profile
go test -memprofile=mem.prof -bench=.

# Analyze with pprof
go tool pprof mem.prof

# Show top allocations
(pprof) top

# Show allocation by source location
(pprof) list Function

# Generate memory graph
(pprof) web
```

### Heap Profiling

**When to use**: Understanding memory allocation patterns.

```bash
# Capture heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Analyze
go tool pprof heap.prof

# Show top allocations
(pprof) top

# Compare two heap profiles
go tool pprof -base=heap1.prof heap2.prof
```

## Advanced Testing Techniques

### Race Detection

**When to use**: Detecting concurrent access to shared data.

**Example**: See `examples/unit_tests.go` for complete implementation.

```bash
# Run tests with race detection
go test -race

# Run with race detector and verbose output
go test -race -v

# Build with race detection
go build -race
```

**Common Race Conditions**:
- Concurrent map writes
- Unprotected shared variables
- Data races in goroutines

### Example-Based Tests

**When to use**: Providing usage examples that also serve as tests.

```go
func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}
```

**Running Example Tests**:
```bash
go test -run Examples
```

### Fuzzing

**When to use**: Finding edge cases and vulnerabilities with random inputs.

```go
func FuzzReverse(f *testing.F) {
    // Add seed corpus
    f.Add("hello")
    f.Add("world")

    // Fuzzing function
    f.Fuzz(func(t *testing.T, s string) {
        reversed := Reverse(s)
        if Reverse(reversed) != s {
            t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, reversed, s)
        }
    })
}
```

**Running Fuzz Tests**:
```bash
go test -fuzz=FuzzReverse
```

### Subtests and Sub-benchmarks

**When to use**: Organizing related tests and running subsets.

```go
func TestProcess(t *testing.T) {
    tests := []struct {
        name string
        input string
        expected string
    }{
        {"lowercase", "hello", "HELLO"},
        {"uppercase", "WORLD", "WORLD"},
        {"mixed", "Go", "GO"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Process(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

**Running specific subtest**:
```bash
go test -run TestProcess/lowercase
```

### Test Coverage

**When to use**: Ensuring adequate test coverage of codebase.

**Coverage Commands**:

```bash
# Run tests with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# View coverage percentage
go tool cover -func=coverage.out

# Set coverage threshold
go test -coverprofile=coverage.out -covermode=count
go tool cover -func=coverage.out | grep total
```

**Coverage Modes**:
- `set`: Did each statement run?
- `count`: How many times did each statement run?
- `atomic`: Atomic counting (accurate but slower)

## Debugging Strategies

### Delve Debugger

**Installation**:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Usage**:
```bash
# Debug a test
dlv test ./...

# Debug main package
dlv debug ./cmd/server

# Debug a test function
dlv test -run TestFunction

# Attach to running process
dlv attach <pid>

# Common delve commands:
# (dlv) break main.main       # Set breakpoint
# (dlv) breakpoints           # List breakpoints
# (dlv) next                  # Step over
# (dlv) step                  # Step into
# (dlv) continue              # Continue execution
# (dlv) print variable        # Print variable
# (dlv) locals                # Show local variables
# (dlv) goroutines            # List goroutines
# (dlv) goroutine 5           # Switch to goroutine 5
```

### Printf Debugging

**When to use**: Quick debugging without heavy tooling.

```go
fmt.Printf("Debug: value = %+v\n", value)
log.Printf("Debug: user = %+v\n", user)
```

### Runtime Profiling

**When to use**: Profiling running applications.

**Add to your code**:
```go
import (
    _ "net/http/pprof"
    "net/http"
)

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // Your application code
}
```

**Access profiles**:
```bash
# CPU profile
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

# Heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof

# Goroutine profile
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof

# Block profile
curl http://localhost:6060/debug/pprof/block > block.prof
```

### Logging Best Practices

**Structured Logging**:
```go
import "log/slog"

logger := slog.Default()
logger.Info("user logged in",
    "user_id", userID,
    "ip", remoteAddr,
    "timestamp", time.Now(),
)
```

**Context-Aware Logging**:
```go
func handleRequest(ctx context.Context, req Request) {
    logger := slog.With("request_id", req.ID)
    logger.Info("processing request")
    // ... processing
}
```

## Best Practices

### DO ✅

1. **Write tests alongside code**
   - Tests in same package: `package foo`
   - Tests in separate package: `package foo_test`

2. **Use table-driven tests for multiple cases**
   - Easy to add new test cases
   - Clear test data separation

3. **Use descriptive test names**
   ```go
   func TestUserRepository_ReturnsError_WhenUserNotFound(t *testing.T)
   ```

4. **Keep tests independent**
   - No shared state between tests
   - Each test should setup/teardown its own data

5. **Use t.Helper() in helper functions**
   - Improves error reporting
   - Shows correct line numbers

6. **Test edge cases and error conditions**
   - Empty inputs
   - Invalid inputs
   - Boundary conditions

7. **Use subtests for organization**
   ```go
   t.Run("success case", func(t *testing.T) { ... })
   t.Run("error case", func(t *testing.T) { ... })
   ```

8. **Run tests with race detector**
   ```bash
   go test -race ./...
   ```

9. **Maintain high test coverage**
   - Aim for >80% coverage
   - Focus on critical paths

10. **Use build tags for integration tests**
    ```go
    //go:build integration
    // +build integration
    ```

### DON'T ❌

1. **Don't skip tests without reason**
   ```go
   // Bad
   t.Skip("TODO: implement")

   // Good
   t.Skip("Skipping database test in CI environment")
   ```

2. **Don't ignore test failures**
   - Investigate flaky tests
   - Fix race conditions immediately

3. **Don't test external libraries**
   - Test your code, not dependencies
   - Mock external dependencies

4. **Don't write tests that are too brittle**
   - Avoid testing implementation details
   - Focus on behavior and outcomes

5. **Don't forget cleanup**
   ```go
   // Bad
   db := setupDB()

   // Good
   db := setupDB()
   defer db.Close()
   ```

6. **Don't use production databases in tests**
   - Use separate test database
   - Use in-memory databases when possible

7. **Don't ignore coverage reports**
   - Review uncovered code
   - Add tests for critical paths

8. **Don't write slow tests**
   - Tests should run quickly
   - Move slow tests to integration suite

## Running the Examples

All examples in this directory are runnable:

```bash
# Run all tests
cd /home/dima/Desktop/FUN/go-pro/advanced-topics/04-testing-debugging
go test -v ./...

# Run unit tests
go test -v ./examples/unit_tests.go

# Run integration tests
go test -v ./examples/integration_tests.go

# Run benchmarks
go test -v -bench=. ./examples/benchmarking.go

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...
```

## Additional Resources

- [Go Testing Blog](https://go.dev/blog/testing)
- [Testing Guidelines](https://go.dev/doc/add#testing)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [The Go Blog: Using Subtests and Sub-benchmarks](https://go.dev/blog/subtests)
- [Delve Debugger](https://github.com/go-delve/delve)
