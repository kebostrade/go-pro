# REST API Template with Go and chi v5

A production-ready REST API project template using Go and the chi v5 router.

## Features

- **chi v5 Router**: Lightweight, idiomatic HTTP router
- **Middleware Stack**: RequestID, Logger, Recoverer, Timeout, RealIP
- **Clean Architecture**: Handler → Service → Repository layers
- **In-Memory Storage**: Thread-safe with sync.RWMutex
- **Docker Support**: Multi-stage build for small images
- **CI/CD**: GitHub Actions with test, lint, and docker build

## Project Structure

```
basic/projects/rest-api/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── handler/handler.go      # HTTP handlers
│   ├── service/service.go      # Business logic
│   └── repository/memory.go    # In-memory data store
├── pkg/errors/errors.go         # Custom error types
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
# Run the server
go run ./cmd/server

# Or use make
make run
```

### Docker

```bash
# Build and run with docker-compose
docker-compose up --build

# Or build and run manually
docker build -t rest-api .
docker run -p 8080:8080 rest-api
```

## API Endpoints

| Method | Endpoint           | Description          |
|--------|-------------------|----------------------|
| GET    | /health           | Health check         |
| POST   | /api/v1/users     | Create a new user    |
| GET    | /api/v1/users     | List all users       |
| GET    | /api/v1/users/{id}| Get user by ID       |
| PUT    | /api/v1/users/{id}| Update user          |
| DELETE | /api/v1/users/{id}| Delete user          |

## Example Requests

```bash
# Health check
curl http://localhost:8080/health

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

# List users
curl http://localhost:8080/api/v1/users

# Get user
curl http://localhost:8080/api/v1/users/seed-1

# Update user
curl -X PUT http://localhost:8080/api/v1/users/seed-1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/seed-1
```

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Make Commands

```bash
make run       # Run the server
make build     # Build the binary
make test      # Run tests
make lint      # Run golangci-lint
make docker    # Build Docker image
make clean     # Clean up
```

## License

MIT
