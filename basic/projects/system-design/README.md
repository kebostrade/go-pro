# System Design: Clean Architecture, Patterns, and Case Studies

Production-ready System Design template demonstrating clean architecture, concurrency patterns, circuit breaker, and a URL shortener case study.

## Features

- **Clean Architecture**: Domain-driven design with repository pattern
- **Circuit Breaker**: Fault tolerance using sony/breaker library
- **Worker Pool**: Concurrent task processing
- **Caching**: In-memory cache with TTL
- **URL Shortener Case Study**: Complete implemented case study

## Patterns Demonstrated

### Clean Architecture

```
┌──────────────────────────────────────┐
│           Use Cases (Services)        │
├──────────────────────────────────────┤
│        Repository Interfaces           │
├──────────────────────────────────────┤
│         Concrete Implementations       │
└──────────────────────────────────────┘
```

### Circuit Breaker

The circuit breaker pattern prevents cascading failures:

```
Closed → (failure threshold) → Open → (timeout) → Half-Open → (success) → Closed
```

### Worker Pool

Concurrent task processing with bounded parallelism:

```
Jobs → [ Worker 1 ] → Results
     → [ Worker 2 ] →
     → [ Worker 3 ] →
```

## Quick Start

```bash
# Start infrastructure
docker-compose up -d redis postgres

# Run the URL shortener
go run ./examples/case_study_url_shortener.go

# Or build and run
go build -o server ./examples/case_study_url_shortener.go
./server
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/shorten` | POST | Create short URL |
| `/{code}` | GET | Redirect to original URL |
| `/api/stats/{code}` | GET | Get URL statistics |

### Create Short URL

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/very/long/path"}'
```

### Redirect

```bash
curl -I http://localhost:8080/abc123
```

## Architecture Components

### Clean Architecture

- **Entities**: User, Order
- **Repository Interfaces**: UserRepository, OrderRepository
- **Use Cases**: CreateUserUseCase, GetUserUseCase, etc.
- **In-Memory Implementation**: For development/testing

### Circuit Breaker

```go
cb := circuit.NewCircuitBreaker("external-service")

err := cb.Execute(ctx, func() error {
    return callExternalService()
})
```

### Worker Pool

```go
pool := concurrency.NewWorkerPool(10, processor)
pool.Start()
pool.Submit(concurrency.WorkItem{ID: "1", Payload: data})
pool.Stop()
```

## Project Structure

```
system-design/
├── cmd/server/main.go              # Main application entry
├── internal/
│   ├── clean/                    # Clean architecture
│   │   ├── user_repository.go
│   │   └── use_cases.go
│   ├── circuit/                 # Circuit breaker
│   │   └── breaker.go
│   ├── concurrency/            # Worker pool
│   │   └── worker_pool.go
│   └── cache/                  # Caching
│       └── cache.go
├── examples/
│   └── case_study_url_shortener.go
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_URL` | Redis connection URL | `redis://localhost:6379` |
| `DATABASE_URL` | PostgreSQL connection URL | `postgres://localhost:5432/db` |
| `PORT` | HTTP server port | `8080` |

## License

MIT
