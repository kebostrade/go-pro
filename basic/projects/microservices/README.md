# Microservices Template

A microservices architecture template demonstrating Docker Compose service discovery, inter-service HTTP communication, and API Gateway pattern.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway (:8080)                      │
│                         ┌─────────────┐                          │
│   ┌────────────────────→│  service-a  │ (User Service :8001)   │
│   │  /api/users/*        │   (Go HTTP) │                        │
│   │                      └─────────────┘                        │
│   │                                                                
│   │  /api/orders/*      ┌─────────────┐                         │
│   └────────────────────→│  service-b  │ (Order Service :8002)    │
│                         │   (Go HTTP) │                         │
│                         └─────────────┘                         │
│                              │                                   
│                              │ USER_SERVICE_URL                 
│                              ↓                                   
│                         ┌─────────────┐                         │
└────────────────────────→│  service-a  │◄───────────────────────┘
                          └─────────────┘
                               ▲
                               │
                    ┌──────────┴──────────┐
                    │   users-db (5433)    │   orders-db (5434)
                    │   PostgreSQL 15       │   PostgreSQL 15
                    └─────────────────────┘
                               ▲
                               │
                         ┌─────┴─────┐
                         │   redis   │
                         │   (:6379) │
                         └───────────┘
```

## Services

| Service | Port | Purpose |
|---------|------|---------|
| api-gateway | 8080 | Routes requests to backend services |
| service-a | 8001 | User service - manages users |
| service-b | 8002 | Order service - manages orders |
| users-db | 5433 | PostgreSQL for user data |
| orders-db | 5434 | PostgreSQL for order data |
| redis | 6379 | Redis cache |

## Docker DNS Service Discovery

Services discover each other via Docker DNS:
- `service-a` is reachable at `http://service-a:8001` from other containers
- `service-b` is reachable at `http://service-b:8002` from other containers
- Gateway uses `USER_SERVICE_URL=http://service-a:8001` and `ORDER_SERVICE_URL=http://service-b:8002`

## Quick Start

```bash
# Start all services
docker-compose up -d

# Check health
curl http://localhost:8080/health

# Get users via gateway
curl http://localhost:8080/api/users

# Get orders via gateway
curl http://localhost:8080/api/orders

# Direct access to service-a
curl http://localhost:8001/health
curl http://localhost:8001/api/users

# Direct access to service-b
curl http://localhost:8002/health
curl http://localhost:8002/api/orders
```

## Development

```bash
# Build without Docker
go build ./cmd/service-a
go build ./cmd/service-b
go build ./cmd/gateway

# Run directly (requires services running)
SERVICE_PORT=8001 go run ./cmd/service-a
SERVICE_PORT=8002 USER_SERVICE_URL=http://localhost:8001 go run ./cmd/service-b
GATEWAY_PORT=8080 USER_SERVICE_URL=http://localhost:8001 ORDER_SERVICE_URL=http://localhost:8002 go run ./cmd/gateway

# Run tests
go test -v ./...

# Validate docker-compose
docker-compose config
```

## Environment Variables

### service-a
- `SERVICE_PORT` - Port to listen on (default: 8001)
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `REDIS_HOST` - Redis host

### service-b
- `SERVICE_PORT` - Port to listen on (default: 8002)
- `USER_SERVICE_URL` - URL to user service (default: http://service-a:8001)

### api-gateway
- `GATEWAY_PORT` - Port to listen on (default: 8080)
- `USER_SERVICE_URL` - URL to user service
- `ORDER_SERVICE_URL` - URL to order service

## Project Structure

```
microservices/
├── cmd/
│   ├── gateway/          # API Gateway
│   ├── service-a/        # User Service
│   └── service-b/        # Order Service
├── internal/
│   └── gateway/          # Gateway internals (proxy, registry, routes)
├── docker-compose.yml    # Full stack orchestration
├── Dockerfile.*          # Multi-stage Dockerfiles per service
├── Makefile
├── go.mod
└── README.md
```

## API Endpoints

### Gateway (:8080)
- `GET /health` - Gateway health
- `GET /api/users` - List users (proxied to service-a)
- `GET /api/users/{id}` - Get user by ID (proxied to service-a)
- `GET /api/orders` - List orders (proxied to service-b)
- `GET /api/orders/{id}` - Get order by ID (proxied to service-b)

### Service-A (:8001)
- `GET /health` - Service health
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID

### Service-B (:8002)
- `GET /health` - Service health
- `GET /api/orders` - List all orders
- `GET /api/orders/{id}` - Get order by ID
- `GET /api/orders/user/{user_id}` - Get orders for a user

## License

MIT
