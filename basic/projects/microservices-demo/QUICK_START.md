# 🚀 Quick Start - Microservices Demo

Get your microservices architecture running in 5 minutes!

## ⚡ 1-Minute Docker Setup

```bash
# Navigate to project
cd basic/projects/microservices-demo

# Start all services with Docker
make docker-up
```

**That's it!** All services are now running:
- API Gateway: http://localhost:8080
- User Service: http://localhost:8081
- Product Service: http://localhost:8082
- Order Service: http://localhost:8083

## 📱 Test the System

### 1. Check Health
```bash
curl http://localhost:8080/health
# Response: OK
```

### 2. List Services
```bash
curl http://localhost:8080/services
# Shows all registered services
```

### 3. Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

### 4. Login
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'
```

**Save the token from the response!**

### 5. List Products
```bash
curl http://localhost:8080/api/products
```

### 6. Create an Order
```bash
# Replace YOUR_TOKEN with the JWT from step 4
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_id": "user-id-from-step-3",
    "product_id": "1",
    "quantity": 2,
    "total": 59.98
  }'
```

## 💻 Local Development Setup

### Prerequisites
- Go 1.21+
- Make (optional)

### Run Services Locally

```bash
# Terminal 1: User Service
make run-user

# Terminal 2: Product Service
make run-product

# Terminal 3: Order Service
make run-order

# Terminal 4: API Gateway
make run-gateway
```

## 🔧 Common Commands

```bash
# Build all services
make build

# Run tests
make test

# View Docker logs
make docker-logs

# Stop all services
make docker-down

# Clean everything
make docker-clean
```

## 🎯 Quick Challenges

Try these to learn the system:

1. **User Flow**: Create 3 users, login with each
2. **Product Management**: Add 5 new products
3. **Order Processing**: Create orders for different users
4. **Service Discovery**: Check which services are registered
5. **Rate Limiting**: Send 300 requests quickly and see rate limiting in action
6. **Authentication**: Try accessing protected endpoints without a token

## 🐛 Troubleshooting

### Services won't start
```bash
# Check if ports are in use
lsof -i :8080
lsof -i :8081
lsof -i :8082
lsof -i :8083

# Kill processes if needed
kill -9 <PID>
```

### Docker issues
```bash
# Clean everything and restart
make docker-clean
make docker-up
```

### Can't connect to services
```bash
# Check service health
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Product Service
curl http://localhost:8083/health  # Order Service
curl http://localhost:8080/health  # API Gateway
```

## 📚 Next Steps

1. **Read the full README**: [README.md](README.md)
2. **Study the tutorial**: [Tutorial 13 in TUTORIALS.md](../../docs/TUTORIALS.md)
3. **Explore the code**: Start with `services/api-gateway/cmd/main.go`
4. **Customize**: Add new endpoints, services, or features
5. **Deploy**: Try Kubernetes deployment

## 🔗 Resources

- [Full Documentation](README.md)
- [Tutorial 13](../../docs/TUTORIALS.md)
- [Microservices Patterns](https://microservices.io/)
- [Go Documentation](https://go.dev/doc/)

---

**Happy Coding! 🎉**

