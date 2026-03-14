# Testing and Debugging in Go

Master Go's testing framework, debugging tools, and best practices.

## Learning Objectives

- Write unit and integration tests
- Use table-driven tests effectively
- Benchmark and profile code
- Debug with delve
- Mock dependencies
- Achieve high test coverage

## Theory

### Table-Driven Tests

```go
func TestCalculate(t *testing.T) {
    tests := []struct {
        name    string
        input   int
        want    int
        wantErr bool
    }{
        {"positive", 5, 25, false},
        {"negative", -3, 9, false},
        {"zero", 0, 0, false},
        {"large", 1000000, 1000000000000, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Calculate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Calculate() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Mocking with Interfaces

```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

type MockUserRepository struct {
    users map[string]*User
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    if u, ok := m.users[id]; ok {
        return u, nil
    }
    return nil, ErrNotFound
}

func TestUserService_Get(t *testing.T) {
    mock := &MockUserRepository{
        users: map[string]*User{"1": {ID: "1", Name: "Test"}},
    }
    svc := NewUserService(mock)

    user, err := svc.Get(context.Background(), "1")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Test" {
        t.Errorf("got %q, want %q", user.Name, "Test")
    }
}
```

### Benchmarks

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

func BenchmarkProcessParallel(b *testing.B) {
    data := generateTestData(1000)

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Process(data)
        }
    })
}
```

### Fuzzing

```go
func FuzzParse(f *testing.F) {
    testcases := []string{"hello", "world", "123"}
    for _, tc := range testcases {
        f.Add(tc)
    }

    f.Fuzz(func(t *testing.T, orig string) {
        parsed, err := Parse(orig)
        if err != nil {
            return
        }
        result := parsed.String()
        if result != orig {
            t.Errorf("roundtrip failed: got %q, want %q", result, orig)
        }
    })
}
```

## Debugging with Delve

```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a test
dlv test ./... -- -test.run TestMyFunction

# Debug with breakpoints
dlv debug ./cmd/myapp
(dlv) break main.go:42
(dlv) continue
(dlv) print variableName
(dlv) next
```

## Profiling

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
}
```

```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# In browser
go tool pprof -http=:8080 cpu.prof
```

## Test Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Exercises

1. Write table-driven tests for a calculator
2. Create mocks for database tests
3. Benchmark string operations
4. Profile a memory leak

## Validation

```bash
cd exercises
go test -v -race -cover ./...
go test -bench=. -benchmem ./...
go test -fuzz=. ./...
```

## Key Takeaways

- Use table-driven tests for clarity
- Mock external dependencies
- Always check error conditions
- Profile before optimizing
- Aim for >80% coverage on critical paths

## Next Steps

**[AT-04: Gin Framework](../AT-04-gin-framework/README.md)**

---

Test early, test often, test everything. 🧪
