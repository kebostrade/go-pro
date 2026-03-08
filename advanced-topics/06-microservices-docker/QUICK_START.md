# Quick Start Guide

Get the microservices architecture running in 5 minutes.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- curl (for testing)

## Start Services

```bash
# Navigate to project directory
cd /home/dima/Desktop/FUN/go-pro/advanced-topics/06-microservices-docker

# Start all services
docker-compose up -d

# Wait for services to be healthy (30-60 seconds)
docker-compose ps
```

## Verify Deployment

```bash
# Check API Gateway health (includes all service health)
curl http://localhost:8080/health

# Expected response:
{
  "status": "healthy",
  "service": "api-gateway",
  "timestamp": "2024-01-08T10:30:00Z",
  "services": {
    "user-service": "healthy",
    "order-service": "healthy"
  }
}
```

## Test the Services

### 1. Create a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }'

# Response:
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "created_at": "2024-01-08T10:30:00Z",
  "updated_at": "2024-01-08T10:30:00Z"
}
```

### 2. Get User by ID

```bash
# Replace {user_id} with the ID from step 1
curl http://localhost:8080/api/users/{user_id}
```

### 3. List All Users

```bash
curl http://localhost:8080/api/users
```

### 4. Create an Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "{user_id}",
    "items": [
      {"product": "Widget", "quantity": 2, "price": 24.99},
      {"product": "Gadget", "quantity": 1, "price": 49.99}
    ],
    "total": 99.97
  }'

# Response:
{
  "id": "order-id-here",
  "user_id": "user-id-here",
  "items": [...],
  "total": 99.97,
  "status": "pending",
  "created_at": "2024-01-08T10:31:00Z",
  "updated_at": "2024-01-08T10:31:00Z"
}
```

### 5. List All Orders

```bash
curl http://localhost:8080/api/orders
```

### 6. Get Orders by User

```bash
curl http://localhost:8080/api/orders/user/{user_id}
```

### 7. Update Order Status

```bash
curl -X PUT http://localhost:8080/api/orders/{order_id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "processing"
  }'
```

## View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway
docker-compose logs -f service-a
docker-compose logs -f service-b
```

## Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove data volumes
docker-compose down -v
```

## Architecture Diagram

```
Client → API Gateway (8080) → User Service (8001) → users-db
                               ↓
                          Order Service (8002) → orders-db
                               ↓
                            Redis (6379)
```

## Service Endpoints

### API Gateway (Port 8080)
- `GET /health` - Health check for all services
- `GET /api/users` - List all users
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `GET /api/orders` - List all orders
- `POST /api/orders` - Create order
- `GET /api/orders/:id` - Get order by ID
- `GET /api/orders/user/:user_id` - Get orders by user
- `PUT /api/orders/:id/status` - Update order status
- `DELETE /api/orders/:id` - Delete order

### User Service (Port 8001)
- `GET /health` - Service health check
- Same endpoints as API Gateway but without `/api` prefix

### Order Service (Port 8002)
- `GET /health` - Service health check
- Same endpoints as API Gateway but without `/api` prefix

## Run Examples Script

```bash
# Run all examples
./examples.sh
```

## Troubleshooting

### Services not starting

```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache

# Restart services
docker-compose restart
```

### Port conflicts

Edit `docker-compose.yml` and change external ports:

```yaml
services:
  api-gateway:
    ports:
      - "8081:8080"  # Use port 8081 instead
```

### Database connection issues

```bash
# Check database status
docker-compose ps users-db orders-db

# View database logs
docker-compose logs users-db
docker-compose logs orders-db

# Restart databases
docker-compose restart users-db orders-db
```

## Next Steps

1. Read [README.md](README.md) for detailed documentation
2. Review [ARCHITECTURE.md](ARCHITECTURE.md) for architecture details
3. Check [DEPLOYMENT.md](DEPLOYMENT.md) for deployment strategies
4. Modify services and rebuild with `docker-compose up -d --build <service>`

## Key Concepts Demonstrated

- **Microservices Architecture**: Separate services for different domains
- **API Gateway Pattern**: Single entry point for all requests
- **Database per Service**: Each service has its own database
- **Service Communication**: HTTP/REST between services
- **Docker Orchestration**: Docker Compose for multi-container apps
- **Health Checks**: Health endpoints for monitoring
- **Structured Logging**: JSON logs with correlation IDs
- **Configuration Management**: Environment variables
- **Caching**: Redis for performance
- **Network Isolation**: Docker networks for security

## Support

For issues or questions:
1. Check logs: `docker-compose logs -f`
2. Verify health: `curl http://localhost:8080/health`
3. Review documentation in this directory
