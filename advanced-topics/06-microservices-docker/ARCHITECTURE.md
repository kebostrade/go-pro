# Microservices Architecture

This document provides a detailed explanation of the microservices architecture implemented in this example.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Design Principles](#design-principles)
3. [Service Communication](#service-communication)
4. [Data Management](#data-management)
5. [Configuration Management](#configuration-management)
6. [Security Considerations](#security-considerations)
7. [Scalability Strategies](#scalability-strategies)
8. [Deployment Patterns](#deployment-patterns)

## Architecture Overview

### Component Diagram

```
┌─────────────────────────────────────────────────────────┐
│                      Client                             │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│                   API Gateway                           │
│  Responsibilities:                                      │
│  - Request routing                                       │
│  - Load balancing                                        │
│  - CORS handling                                         │
│  - Request logging                                       │
│  - Health aggregation                                    │
└─────────┬───────────────────────────────┬───────────────┘
          │                               │
          ▼                               ▼
┌─────────────────────┐         ┌─────────────────────┐
│   User Service      │         │   Order Service     │
│   (service-a)       │         │   (service-b)       │
│                     │         │                     │
│  - User management  │         │  - Order processing │
│  - CRUD operations  │         │  - User validation  │
│  - Redis caching    │         │  - Order tracking   │
└──────────┬──────────┘         └──────────┬──────────┘
           │                               │
           │                               │
    ┌──────▼──────┐                ┌──────▼──────┐
    │ users-db    │                │ orders-db   │
    │ (PostgreSQL)│                │ (PostgreSQL)│
    └─────────────┘                └─────────────┘
           │
           └──────────┬──────────────────────┘
                      │
               ┌──────▼──────┐
               │ Redis Cache │
               └─────────────┘
```

## Design Principles

### 1. Single Responsibility Principle

Each service has a single, well-defined responsibility:

- **User Service**: Manages user data and authentication
- **Order Service**: Manages order processing and tracking
- **API Gateway**: Routes requests and handles cross-cutting concerns

### 2. Database per Service

Each service owns its database:

```go
// User Service Database
users-db:
  - users table
  - User-specific data only

// Order Service Database
orders-db:
  - orders table
  - Order-specific data only
```

**Benefits**:
- Service autonomy
- Independent scaling
- Fault isolation
- Technology flexibility

### 3. Loose Coupling

Services communicate through well-defined APIs:

```go
// Order Service calls User Service via HTTP
GET http://service-a:8001/api/users/{user_id}
```

**Benefits**:
- Independent deployment
- Technology diversity
- Team autonomy

### 4. High Cohesion

Related functionality is grouped together:

```go
// User Service endpoints
POST   /api/users
GET    /api/users
GET    /api/users/:id
PUT    /api/users/:id
DELETE /api/users/:id
```

## Service Communication

### Synchronous Communication (HTTP/REST)

**Current Implementation**:

```go
// Order Service validates users via User Service
func (r *Repository) ValidateUser(ctx context.Context, userID string) error {
    resp, err := http.Get(r.userServiceURL + "/api/users/" + userID)
    if err != nil || resp.StatusCode != http.StatusOK {
        return fmt.Errorf("user not found")
    }
    return nil
}
```

**Pros**:
- Simple to implement
- Easy to debug
- Real-time response

**Cons**:
- Tight coupling
- Cascading failures
- Blocking operations

### Communication Flow

```
1. Client Request
   ↓
2. API Gateway
   ↓
3. Route to Service
   ↓
4. Service Processing
   ↓
5. Database Operations
   ↓
6. Response (back through chain)
```

### Example Request Flow

```bash
# Client creates order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123",
    "items": [{"product": "Widget", "quantity": 2}],
    "total": 99.99
  }'

# Flow:
# 1. API Gateway receives request
# 2. Routes to Order Service
# 3. Order Service validates user via User Service
# 4. Order Service saves to orders-db
# 5. Response returned through API Gateway
```

## Data Management

### Database Design

#### User Service Database

```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

#### Order Service Database

```sql
CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    items JSONB NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
```

### Caching Strategy

```go
// User Service uses Redis caching
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
    // Try cache first
    cached, err := r.redis.Get(ctx, "user:"+id).Result()
    if err == nil {
        var user User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }

    // Cache miss, query database
    user, err := r.queryDatabase(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cache the result
    data, _ := json.Marshal(user)
    r.redis.Set(ctx, "user:"+id, data, 5*time.Minute)

    return user, nil
}
```

**Cache Keys**:
- `user:{id}` - Individual user data (5 min TTL)
- `users:all` - List of all users (2 min TTL)

**Cache Invalidation**:
- Write operations invalidate relevant cache keys
- Time-based expiration

## Configuration Management

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

### Configuration Best Practices

1. **Externalize Configuration**: All config via environment variables
2. **Sane Defaults**: Provide defaults for development
3. **Validation**: Validate configuration on startup
4. **No Secrets in Code**: Never hardcode credentials

## Security Considerations

### Current Implementation

1. **Network Isolation**: Docker networks for service isolation
2. **Database Security**: Separate databases per service
3. **CORS**: API Gateway handles CORS

### Recommendations for Production

1. **Authentication**: Implement JWT-based authentication
2. **Authorization**: Add role-based access control
3. **TLS/SSL**: Use HTTPS for all communication
4. **Secrets Management**: Use Docker secrets or external vault
5. **Rate Limiting**: Add rate limiting to API Gateway
6. **Input Validation**: Validate all inputs
7. **SQL Injection**: Use parameterized queries (already implemented)

## Scalability Strategies

### Horizontal Scaling

```yaml
# docker-compose.yml
services:
  service-a:
    deploy:
      replicas: 3  # Run 3 instances

  api-gateway:
    deploy:
      replicas: 2  # Run 2 instances
```

### Database Scaling

1. **Read Replicas**: Scale read operations
2. **Connection Pooling**: Reuse database connections
3. **Caching**: Reduce database load

### Load Balancing

Current: Round-robin via Docker Compose
Production: Use dedicated load balancer (HAProxy, Nginx)

## Deployment Patterns

### Blue-Green Deployment

``# Step 1: Deploy blue environment
docker-compose -f docker-compose.blue.yml up -d

# Step 2: Test blue environment
curl http://blue.example.com/health

# Step 3: Switch traffic to blue
# Update DNS or load balancer

# Step 4: Deploy green environment
docker-compose -f docker-compose.green.yml up -d

# Step 5: Test green environment
curl http://green.example.com/health

# Step 6: Switch traffic to green or rollback to blue
````

### Rolling Update

```bash
# Update service one instance at a time
docker-compose up -d --no-deps --build service-a
```

### Canary Deployment

```yaml
# Route 10% of traffic to new version
api-gateway:
  environment:
    CANARY_PERCENTAGE: 10
    NEW_SERVICE_URL: http://service-a-v2:8001
```

## Monitoring and Observability

### Logging

All services use structured JSON logging:

```json
{
  "level": "info",
  "service": "user-service",
  "method": "GET",
  "path": "/api/users/123",
  "status": 200,
  "duration_ms": 45,
  "correlation_id": "abc123",
  "client_ip": "192.168.1.1"
}
```

### Health Checks

Each service exposes a `/health` endpoint:

```go
GET /health

Response:
{
  "status": "healthy",
  "service": "user-service",
  "timestamp": "2024-01-08T10:30:00Z",
  "db": "healthy",
  "redis": "healthy"
}
```

### Metrics to Track

1. **Request Latency**: Response time per endpoint
2. **Error Rate**: Failed requests percentage
3. **Throughput**: Requests per second
4. **Database Connections**: Active connections
5. **Cache Hit Rate**: Cache effectiveness

## Fault Tolerance

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures     int
    resetTimeout    time.Duration
    failureCount    int
    lastFailureTime time.Time
    state           string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.state = "half-open"
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }

    err := fn()
    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()
        if cb.failureCount >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    cb.failureCount = 0
    cb.state = "closed"
    return nil
}
```

### Retry Logic

```go
func retry(attempts int, sleep time.Duration, fn func() error) error {
    for i := 0; i < attempts; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if i < attempts-1 {
            time.Sleep(sleep)
            sleep *= 2 // Exponential backoff
        }
    }
    return fmt.Errorf("failed after %d attempts", attempts)
}
```

## Conclusion

This microservices architecture demonstrates:

- **Service Decomposition**: Logical separation of concerns
- **Independent Deployment**: Services can be deployed independently
- **Technology Diversity**: Each service can use different technologies
- **Fault Isolation**: Failure in one service doesn't cascade
- **Scalability**: Services can be scaled independently

The implementation follows best practices while remaining simple enough for learning purposes.
