# Building Microservices with Go and Docker

Develop containerized microservices with proper communication patterns.

## Learning Objectives

- Design microservice architecture
- Containerize Go applications
- Implement service discovery
- Handle inter-service communication
- Implement health checks
- Deploy with Docker Compose

## Theory

### Microservice Structure

```
user-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
│   └── models/
├── Dockerfile
└── docker-compose.yml
```

### Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Service Implementation

```go
type UserService struct {
    repo   UserRepository
    client *http.Client
}

func (s *UserService) GetWithOrders(ctx context.Context, userID string) (*UserWithOrders, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("find user: %w", err)
    }

    orders, err := s.fetchOrders(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("fetch orders: %w", err)
    }

    return &UserWithOrders{
        User:   user,
        Orders: orders,
    }, nil
}

func (s *UserService) fetchOrders(ctx context.Context, userID string) ([]Order, error) {
    url := fmt.Sprintf("http://order-service:8081/orders?user_id=%s", userID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var orders []Order
    if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
        return nil, err
    }

    return orders, nil
}
```

### Health Checks

```go
func (s *Server) setupHealthChecks() {
    http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
    })

    http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
        if err := s.db.Ping(r.Context()); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
    })
}
```

### Docker Compose

```yaml
version: '3.8'
services:
  user-service:
    build: ./user-service
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/users
      - ORDER_SERVICE_URL=http://order-service:8081
    depends_on:
      - db
      - order-service
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3

  order-service:
    build: ./order-service
    ports:
      - "8081:8081"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/orders
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

## Security Considerations

```go
func serviceAuthMiddleware(expectedToken string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("X-Service-Token")
            if token == "" || !constantTimeCompare(token, expectedToken) {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Performance Tips

```go
var httpClient = &http.Client{
    Timeout: 5 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

## Exercises

1. Create two communicating services
2. Implement circuit breaker pattern
3. Add distributed tracing
4. Deploy with Docker Compose

## Validation

```bash
cd exercises
docker-compose up --build
curl http://localhost:8080/health/ready
```

## Key Takeaways

- Keep services small and focused
- Use multi-stage Docker builds
- Implement health endpoints
- Handle service failures gracefully
- Use environment variables for config

## Next Steps

**[AT-07: gRPC Distributed](../AT-07-grpc-distributed/README.md)**

---

Microservices: divide and conquer. 🐳
