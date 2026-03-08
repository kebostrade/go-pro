# Microservices with Docker

A practical guide to building microservices architecture using Go and Docker. This comprehensive example demonstrates production-ready microservices with proper orchestration, communication patterns, and best practices.

## 🎯 Learning Objectives

- **Microservices Architecture**: Understand service decomposition and communication patterns
- **Docker Multi-Stage Builds**: Create optimized, production-ready Docker images
- **Docker Compose Orchestration**: Coordinate multiple services, databases, and networks
- **Service Discovery**: Implement basic service registration and discovery
- **Inter-Service Communication**: HTTP/REST communication between services
- **Configuration Management**: Environment-based configuration for each service
- **Structured Logging**: JSON logging with correlation IDs
- **Health Checks**: Implement and monitor service health
- **Database per Service**: Dedicated PostgreSQL/Redis instances per service
- **Network Isolation**: Docker networks for secure service communication

## 📊 Architecture Overview

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    │   :8080         │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              ▼                             ▼
    ┌─────────────────┐           ┌─────────────────┐
    │  User Service   │           │  Order Service  │
    │  :8001          │           │  :8002          │
    └────────┬────────┘           └────────┬────────┘
             │                             │
             │                             │
    ┌────────▼────────┐           ┌────────▼────────┐
    │  PostgreSQL     │           │  PostgreSQL     │
    │  users-db:5432  │           │  orders-db:5432 │
    └─────────────────┘           └─────────────────┘
             │                             │
             └──────────────┬──────────────┘
                            │
                   ┌────────▼────────┐
                   │  Redis Cache    │
                   │  redis:6379     │
                   └─────────────────┘
```

## 🏗️ Services

### 1. API Gateway (`api-gateway`)
- **Port**: 8080
- **Responsibility**: Route requests to appropriate services, handle CORS, rate limiting
- **Features**:
  - Request routing based on path patterns
  - Load balancing (round-robin)
  - Request/response logging
  - Health check aggregation

### 2. User Service (`service-a`)
- **Port**: 8001
- **Responsibility**: User management, authentication, profile data
- **Features**:
  - CRUD operations for users
  - Database persistence (PostgreSQL)
  - Caching layer (Redis)
  - Health checks
  - Structured logging

### 3. Order Service (`service-b`)
- **Port**: 8002
- **Responsibility**: Order management, order processing
- **Features**:
  - Create and retrieve orders
  - User validation via User Service
  - Database persistence (PostgreSQL)
  - Event logging
  - Health checks

## 🚀 Quick Start

### Prerequisites

- Docker (v20.10+)
- Docker Compose (v2.0+)
- Go 1.23+ (for local development)

### Start All Services

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check service health
curl http://localhost:8080/health
```

### Access Services

```bash
# API Gateway (Main entry point)
curl http://localhost:8080/health

# User Service (via gateway)
curl http://localhost:8080/api/users
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Order Service (via gateway)
curl http://localhost:8080/api/orders
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":"1","items":[{"product":"Widget","quantity":2}],"total":99.99}'
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 📁 Project Structure

```
06-microservices-docker/
├── README.md                 # This file
├── docker-compose.yml        # Orchestration configuration
├── api-gateway/              # API Gateway service
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── routes.go
├── service-a/                # User Service
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── handlers.go
│   ├── models.go
│   └── repository.go
├── service-b/                # Order Service
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── handlers.go
│   ├── models.go
│   └── repository.go
└── docs/
    ├── architecture.md       # Detailed architecture
    ├── deployment.md         # Deployment strategies
    └── monitoring.md         # Monitoring and observability
```

## 🔧 Configuration

### Environment Variables

Each service uses environment variables for configuration:

```bash
# API Gateway
GATEWAY_PORT=8080
USER_SERVICE_URL=http://service-a:8001
ORDER_SERVICE_URL=http://service-b:8002

# User Service
SERVICE_PORT=8001
DB_HOST=users-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=users_db
REDIS_HOST=redis
REDIS_PORT=6379

# Order Service
SERVICE_PORT=8002
DB_HOST=orders-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=orders_db
USER_SERVICE_URL=http://service-a:8001
```

## 🏛️ Architecture Patterns

### 1. API Gateway Pattern

The API Gateway acts as a single entry point for all client requests:

```go
// Route: /api/users/* -> User Service
// Route: /api/orders/* -> Order Service
// Route: /health -> Health check aggregator
```

**Benefits**:
- Single entry point for clients
- Centralized cross-cutting concerns (auth, logging, rate limiting)
- Service abstraction (clients don't need to know service locations)

### 2. Database per Service Pattern

Each service has its own database:

```
User Service -> users-db (PostgreSQL)
Order Service -> orders-db (PostgreSQL)
```

**Benefits**:
- Service autonomy
- Independent scaling
- Technology diversity (could use different DBs per service)
- Fault isolation

### 3. Shared Cache Pattern

Redis is shared across services for caching:

```go
// User Service caches user data
// Order Service caches validation results
```

### 4. Service-to-Service Communication

Services communicate via HTTP/REST:

```go
// Order Service validates users via User Service
resp, err := http.Get("http://service-a:8001/api/users/" + userID)
```

## 🔍 Key Concepts

### Docker Multi-Stage Builds

Each service uses multi-stage builds for optimized images:

```dockerfile
# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o service

# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /app/service /service
ENTRYPOINT ["/service"]
```

**Benefits**:
- Smaller final images (no build tools)
- Better security (fewer attack vectors)
- Faster deployment

### Docker Compose Orchestration

`docker-compose.yml` defines:
- **Services**: API Gateway, User Service, Order Service
- **Networks**: Isolated networks for service communication
- **Volumes**: Persistent data storage for databases
- **Health Checks**: Container health monitoring
- **Depends_on**: Service startup ordering

### Service Health Checks

Each service exposes a `/health` endpoint:

```go
func (h *HealthHandler) HandleHealth(c *gin.Context) {
    status := map[string]string{
        "status": "healthy",
        "service": h.serviceName,
        "timestamp": time.Now().Format(time.RFC3339),
    }
    c.JSON(http.StatusOK, status)
}
```

**Health Check Types**:
1. **Liveness**: Is the service running?
2. **Readiness**: Can the service handle requests?
3. **Startup**: Did the service start successfully?

### Structured Logging

All services use structured JSON logging:

```go
log.WithFields(log.Fields{
    "method": r.Method,
    "path": r.URL.Path,
    "status": statusCode,
    "duration": duration.Milliseconds(),
    "correlation_id": correlationID,
}).Info("Request completed")
```

**Benefits**:
- Parseable logs for log aggregators
- Consistent structure across services
- Easy filtering and searching

## 📊 Monitoring and Observability

### Logs

View logs for all services:

```bash
docker-compose logs -f
docker-compose logs -f api-gateway
docker-compose logs -f service-a
```

### Metrics

Each service logs request metrics:

```json
{
  "level": "info",
  "service": "user-service",
  "method": "GET",
  "path": "/api/users/1",
  "status": 200,
  "duration": 45,
  "correlation_id": "abc123"
}
```

### Health Monitoring

Check health status:

```bash
# API Gateway health (includes downstream service health)
curl http://localhost:8080/health

# Individual service health
curl http://localhost:8080/health/service-a
curl http://localhost:8080/health/service-b
```

## 🧪 Testing

### Test Endpoints

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }'

# Get user by ID
curl http://localhost:8080/api/users/1

# List all users
curl http://localhost:8080/api/users

# Create an order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "1",
    "items": [
      {"product": "Widget", "quantity": 2},
      {"product": "Gadget", "quantity": 1}
    ],
    "total": 149.99
  }'

# Get order by ID
curl http://localhost:8080/api/orders/1

# List all orders
curl http://localhost:8080/api/orders
```

## 🛠️ Development

### Local Development (Without Docker)

```bash
# Start databases with Docker
docker-compose up -d users-db orders-db redis

# Run services locally
cd service-a && go run main.go
cd service-b && go run main.go
cd api-gateway && go run main.go
```

### Build Images

```bash
# Build individual service image
docker-compose build service-a

# Build all images
docker-compose build
```

### Rebuild After Code Changes

```bash
# Rebuild and restart
docker-compose up -d --build

# Rebuild specific service
docker-compose up -d --build service-a
```

## 📚 Advanced Topics

### 1. Service Discovery

Current implementation uses hardcoded URLs. For production, consider:

- **Consul**: Service registration and discovery
- **etcd**: Distributed key-value store
- **Kubernetes**: Built-in service discovery

### 2. Load Balancing

Current implementation uses round-robin in Docker Compose. For production:

- **HAProxy**: Software load balancer
- **Nginx**: Reverse proxy and load balancer
- **Kubernetes Services**: Cluster-level load balancing

### 3. Inter-Service Communication

**Synchronous (HTTP/REST)**:
- ✅ Simple to implement
- ✅ Easy to debug
- ❌ Tight coupling
- ❌ Cascading failures

**Asynchronous (Message Queues)**:
- ✅ Loose coupling
- ✅ Fault tolerance
- ✅ Scalability
- ❌ Complexity
- ❌ Eventual consistency

For production, consider:
- **RabbitMQ**: Message broker
- **Kafka**: Distributed event streaming
- **NATS**: Lightweight messaging

### 4. Authentication and Authorization

Current implementation doesn't include auth. For production, add:

- **JWT**: JSON Web Tokens for authentication
- **OAuth 2.0**: Authorization framework
- **API Keys**: Simple API authentication
- **Service Mesh**: Istio, Linkerd for advanced security

### 5. Deployment Strategies

**Rolling Update**:
```bash
docker-compose up -d --no-deps --build service-a
```

**Blue-Green Deployment**:
- Run two environments (blue and green)
- Switch traffic when new version is ready

**Canary Deployment**:
- Route small percentage of traffic to new version
- Gradually increase if no issues detected

## 🔒 Best Practices

### 1. Security

- ✅ Use secrets management for sensitive data
- ✅ Implement rate limiting
- ✅ Add authentication/authorization
- ✅ Use HTTPS in production
- ✅ Regularly update base images
- ✅ Scan images for vulnerabilities

### 2. Configuration

- ✅ Use environment variables for configuration
- ✅ Never hardcode credentials
- ✅ Provide default values
- ✅ Validate configuration on startup

### 3. Logging

- ✅ Use structured logging (JSON)
- ✅ Include correlation IDs
- ✅ Log at appropriate levels
- ✅ Avoid logging sensitive data

### 4. Error Handling

- ✅ Use consistent error responses
- ✅ Include error codes and messages
- ✅ Log errors with context
- ✅ Don't expose internal details

### 5. Performance

- ✅ Use connection pooling for databases
- ✅ Implement caching where appropriate
- ✅ Set appropriate timeouts
- ✅ Monitor resource usage

## 🐛 Troubleshooting

### Service won't start

```bash
# Check logs
docker-compose logs [service-name]

# Check container status
docker-compose ps

# Restart service
docker-compose restart [service-name]
```

### Database connection issues

```bash
# Check database is running
docker-compose ps users-db

# Check database logs
docker-compose logs users-db

# Verify network
docker network inspect microservices_default
```

### Port conflicts

Edit `docker-compose.yml` to change exposed ports:

```yaml
ports:
  - "8081:8080"  # Change external port
```

### Rebuild from scratch

```bash
# Stop and remove everything
docker-compose down -v

# Rebuild images
docker-compose build --no-cache

# Start fresh
docker-compose up -d
```

## 📖 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Building Microservices](https://www.oreilly.com/library/view/building-microservices/9781491950340/)
- [Go Microservices Tutorial](https://golangbot.com/websockets/)

## 🎓 Exercises

1. **Add a new service**: Create a Product Service with its own database
2. **Implement caching**: Add Redis caching to User Service
3. **Add authentication**: Implement JWT-based authentication in API Gateway
4. **Implement rate limiting**: Add rate limiting to API Gateway
5. **Add metrics**: Integrate Prometheus metrics collection
6. **Implement circuit breaker**: Add circuit breaker for service-to-service calls
7. **Add distributed tracing**: Integrate OpenTelemetry or Jaeger
8. **Implement async communication**: Add message queue for order processing

## 📝 License

This is educational code. Feel free to use and modify as needed.
