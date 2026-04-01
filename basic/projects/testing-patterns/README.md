# Testing Patterns Template with Go and testify

A production-ready testing patterns project template using Go and testify for testing.

## Features

- **testify Suite**: Advanced testing toolkit with assertions and mocks
- **Mock Patterns**: Mock repository implementations for isolated testing
- **HTTP Handler Testing**: httptest with mock servers
- **Client Testing**: Mock HTTP servers for API client tests
- **Coverage**: >80% coverage target

## Project Structure

```
basic/projects/testing-patterns/
├── internal/
│   ├── service/             # Business logic with service tests
│   ├── handler/             # HTTP handlers with handler tests
│   └── client/              # API client with client tests
├── mocks/                   # Generated mock files
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml
├── go.mod
└── Makefile
```

## Prerequisites

- Go 1.23 or later
- Docker (optional)
- Make

## Quick Start

### Local Development

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Testing Patterns Covered

### 1. Service Layer Testing with Mocks

```go
// Mock repository embedding testify's mock.Mock
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// Usage in test
mockRepo := new(MockUserRepository)
svc := NewUserService(mockRepo)
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).Return(nil)
```

### 2. HTTP Handler Testing

```go
func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    svc := service.NewUserService(mockRepo)
    handler := NewUserHandler(svc)
    router := setupTestRouter(handler)

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).Return(nil)

    body := `{"name":"Alice","email":"alice@example.com"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### 3. HTTP Client Testing with Mock Server

```go
func TestGetUser(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := User{ID: "123", Name: "TestUser", Email: "test@example.com"}
        json.NewEncoder(w).Encode(user)
    }))
    defer ts.Close()

    client := NewAPIClient(ts.URL)
    user, err := client.GetUser(context.Background(), "123")

    assert.NoError(t, err)
    assert.Equal(t, "123", user.ID)
}
```

## Make Commands

```bash
make run       # Run tests (alias for test)
make build     # Build the project
make test      # Run all tests with coverage
make lint      # Run golangci-lint
make docker    # Build Docker image
make clean     # Clean up
```

## Docker

```bash
# Build and run with docker-compose
docker-compose up --build

# Build image manually
docker build -t testing-patterns .
```

## CI/CD

GitHub Actions workflow includes:
- Unit tests with race detection
- Code coverage checks (>70%)
- Linting with golangci-lint
- Docker build verification

## License

MIT
