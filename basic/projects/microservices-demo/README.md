# 🏗️ Microservices Architecture Demo

A production-ready microservices architecture built with Go, demonstrating service communication, API Gateway pattern, service discovery, and distributed system design.

## 🎯 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Gateway                             │
│                      (Port 8080)                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Request Routing                                       │  │
│  │  • Authentication & Authorization                        │  │
│  │  • Rate Limiting                                         │  │
│  │  • Load Balancing                                        │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────┬────────────────┬────────────────┬─────────────────┘
             │                │                │
    ┌────────▼────────┐  ┌───▼──────────┐  ┌─▼──────────────┐
    │  User Service   │  │   Product    │  │  Order Service │
    │   (Port 8081)   │  │   Service    │  │  (Port 8083)   │
    │                 │  │ (Port 8082)  │  │                │
    │  • User CRUD    │  │  • Product   │  │  • Order CRUD  │
    │  • Auth/Login   │  │    Catalog   │  │  • Status Mgmt │
    │  • JWT Tokens   │  │  • Inventory │  │  • Events      │
    └────────┬────────┘  └───┬──────────┘  └─┬──────────────┘
             │               │                │
    ┌────────▼────────┐  ┌──▼───────────┐   │
    │   PostgreSQL    │  │    Redis     │   │
    │   (Port 5432)   │  │  (Port 6379) │   │
    └─────────────────┘  └──────────────┘   │
                                             │
                         ┌───────────────────▼──────────┐
                         │  Shared Infrastructure       │
                         │  • Service Discovery         │
                         │  • Logging (Zap)             │
                         │  • Middleware (Auth, Logs)   │
                         └──────────────────────────────┘
```

## 🚀 Features

### Core Microservices Patterns
- ✅ **API Gateway** - Single entry point for all client requests
- ✅ **Service Discovery** - Dynamic service registration and discovery
- ✅ **Circuit Breaker** - Fault tolerance and resilience
- ✅ **Rate Limiting** - Token bucket algorithm for request throttling
- ✅ **Distributed Logging** - Structured logging with Zap
- ✅ **Health Checks** - Service health monitoring

### Services

#### 1. User Service (Port 8081)
- User registration and authentication
- JWT token generation and validation
- Password hashing with bcrypt
- User CRUD operations
- In-memory storage (easily replaceable with PostgreSQL)

#### 2. Product Service (Port 8082)
- Product catalog management
- Inventory tracking
- Product CRUD operations
- Redis caching support
- In-memory storage with seeded data

#### 3. Order Service (Port 8083)
- Order creation and management
- Order status tracking
- Event-driven architecture ready
- Integration with User and Product services

#### 4. API Gateway (Port 8080)
- Request routing to appropriate services
- Authentication middleware
- Rate limiting (200 req/s, burst 400)
- Service discovery integration
- Reverse proxy implementation

### Shared Infrastructure
- **Logger Package** - Structured logging with Zap
- **Discovery Package** - Service registry and discovery
- **Middleware Package** - Auth, logging, rate limiting
- **Proto Definitions** - gRPC service definitions

## 📦 Project Structure

```
microservices-demo/
├── services/
│   ├── api-gateway/          # API Gateway service
│   │   └── cmd/main.go
│   ├── user-service/         # User management service
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── models.go
│   │       ├── repository.go
│   │       └── handler.go
│   ├── product-service/      # Product catalog service
│   │   └── cmd/main.go
│   └── order-service/        # Order management service
│       └── cmd/main.go
├── pkg/                      # Shared packages
│   ├── logger/              # Logging utilities
│   ├── discovery/           # Service discovery
│   ├── middleware/          # HTTP middleware
│   └── proto/               # Protocol buffers
├── deployments/             # Deployment configs
│   ├── docker-compose.yml
│   ├── Dockerfile.*
│   └── k8s/                # Kubernetes manifests
├── Makefile                # Build automation
├── go.mod                  # Go module definition
└── README.md               # This file
```

## 🛠️ Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker & Docker Compose (optional)
- Make (optional, for convenience)

### Option 1: Run Locally

```bash
# 1. Navigate to project
cd basic/projects/microservices-demo

# 2. Download dependencies
make deps

# 3. Build all services
make build

# 4. Run services (in separate terminals)
make run-user      # Terminal 1
make run-product   # Terminal 2
make run-order     # Terminal 3
make run-gateway   # Terminal 4
```

### Option 2: Run with Docker Compose

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop all services
make docker-down
```

## 📡 API Endpoints

### API Gateway (http://localhost:8080)

#### Health & Discovery
```bash
GET /health              # Gateway health check
GET /services            # List all registered services
```

#### User Service (via Gateway)
```bash
POST /api/users                    # Create user
POST /api/users/login              # Login user
GET  /api/users                    # List users (requires auth)
GET  /api/users/{id}               # Get user (requires auth)
PUT  /api/users/{id}               # Update user (requires auth)
DELETE /api/users/{id}             # Delete user (requires auth)
```

#### Product Service (via Gateway)
```bash
GET  /api/products                 # List products
POST /api/products                 # Create product
GET  /api/products/{id}            # Get product
PUT  /api/products/{id}            # Update product
DELETE /api/products/{id}          # Delete product
```

#### Order Service (via Gateway)
```bash
GET  /api/orders                   # List orders
POST /api/orders                   # Create order
GET  /api/orders/{id}              # Get order
PUT  /api/orders/{id}/status       # Update order status
```

## 💻 Usage Examples

### 1. Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

### 3. List Products
```bash
curl http://localhost:8080/api/products
```

### 4. Create an Order (with Auth)
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": "user-uuid",
    "product_id": "1",
    "quantity": 2,
    "total": 59.98
  }'
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run integration tests
make test-integration

# Run benchmarks
make bench

# Check test coverage
make test-coverage
```

## 🔧 Configuration

### Environment Variables

Each service supports the following environment variables:

```bash
# Service Configuration
SERVICE_PORT=8081              # Service port
LOG_LEVEL=debug                # Log level (debug, info, warn, error)

# Database (User Service, Order Service)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=microservices

# Redis (Product Service)
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT Secret (User Service)
JWT_SECRET=your-secret-key-change-in-production
```

## 🎓 Learning Outcomes

By studying this project, you'll learn:

1. **Microservices Architecture** - Service decomposition, boundaries
2. **API Gateway Pattern** - Request routing, authentication, rate limiting
3. **Service Discovery** - Dynamic service registration and lookup
4. **Inter-Service Communication** - HTTP/REST, gRPC (proto definitions)
5. **Distributed Logging** - Structured logging across services
6. **Authentication & Authorization** - JWT tokens, middleware
7. **Rate Limiting** - Token bucket algorithm implementation
8. **Containerization** - Docker, Docker Compose
9. **Graceful Shutdown** - Signal handling, context cancellation
10. **Clean Architecture** - Separation of concerns, dependency injection

## 🚧 Roadmap

- [ ] Add gRPC communication between services
- [ ] Implement event-driven architecture with message broker
- [ ] Add distributed tracing with OpenTelemetry
- [ ] Implement circuit breaker pattern
- [ ] Add Kubernetes deployment manifests
- [ ] Implement saga pattern for distributed transactions
- [ ] Add metrics and monitoring (Prometheus/Grafana)
- [ ] Implement service mesh (Istio)
- [ ] Add API versioning
- [ ] Implement CQRS pattern

## 📚 Resources

- [Microservices Patterns](https://microservices.io/patterns/index.html)
- [Go Microservices Blog](https://go.dev/blog)
- [Tutorial 13: Microservices in Go](../../docs/TUTORIALS.md)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

---

**Built with ❤️ using Go and Microservices Architecture**

