# Deployment Guide

This guide covers various deployment strategies for the microservices architecture.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Compose Deployment](#docker-compose-deployment)
3. [Production Deployment](#production-deployment)
4. [CI/CD Pipeline](#cicd-pipeline)
5. [Monitoring and Logging](#monitoring-and-logging)
6. [Troubleshooting](#troubleshooting)

## Local Development

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.23+ (for running services locally)

### Running Services Locally

#### Option 1: Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check service status
docker-compose ps

# Stop all services
docker-compose down
```

#### Option 2: Running Services Natively

```bash
# Start databases only
docker-compose up -d users-db orders-db redis

# Run services locally (in separate terminals)
cd service-a && go run main.go
cd service-b && go run main.go
cd api-gateway && go run main.go
```

### Development Workflow

```bash
# Make code changes
vim service-a/handlers.go

# Rebuild and restart specific service
docker-compose up -d --build service-a

# View logs for specific service
docker-compose logs -f service-a
```

## Docker Compose Deployment

### Basic Deployment

```bash
# Build and start all services
docker-compose up -d

# Verify services are running
docker-compose ps

# Check service health
curl http://localhost:8080/health
```

### Environment Configuration

Create a `.env` file in the project root:

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
DB_PASSWORD=secure_password_here
DB_NAME=users_db
REDIS_HOST=redis
REDIS_PORT=6379

# Order Service
SERVICE_PORT=8002
DB_HOST=orders-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password_here
DB_NAME=orders_db
USER_SERVICE_URL=http://service-a:8001
```

Update `docker-compose.yml` to use the `.env` file:

```yaml
services:
  service-a:
    env_file:
      - .env
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
```

### Scaling Services

```bash
# Scale specific service
docker-compose up -d --scale service-a=3

# Scale multiple services
docker-compose up -d --scale service-a=3 --scale service-b=2
```

Note: Scaling requires a load balancer or service discovery mechanism.

### Persistent Data

Data is persisted in Docker volumes:

```bash
# List volumes
docker volume ls

# Backup volumes
docker run --rm -v users-db-data:/data -v $(pwd):/backup alpine tar czf /backup/users-db-backup.tar.gz /data

# Restore volumes
docker run --rm -v users-db-data:/data -v $(pwd):/backup alpine tar xzf /backup/users-db-backup.tar.gz -C /
```

## Production Deployment

### Using Docker Swarm

#### Initialize Swarm

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml microservices

# List services
docker service ls

# Scale services
docker service scale microservices_service-a=3

# Remove stack
docker stack rm microservices
```

#### Production docker-compose.yml

```yaml
version: '3.8'

services:
  service-a:
    image: your-registry/user-service:latest
    environment:
      - DB_PASSWORD_FILE=/run/secrets/db_password
    secrets:
      - db_password
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    networks:
      - microservices

secrets:
  db_password:
    external: true

networks:
  microservices:
    driver: overlay
```

### Using Kubernetes

#### Namespace and ConfigMaps

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: microservices
```

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
  namespace: microservices
data:
  SERVICE_PORT: "8001"
  DB_HOST: "users-db-service"
  DB_PORT: "5432"
  DB_NAME: "users_db"
```

#### Deployment

```yaml
# user-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: microservices
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: your-registry/user-service:latest
        ports:
        - containerPort: 8001
        envFrom:
        - configMapRef:
            name: user-service-config
        - secretRef:
            name: db-credentials
        livenessProbe:
          httpGet:
            path: /health
            port: 8001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8001
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Service

```yaml
# user-service-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: microservices
spec:
  selector:
    app: user-service
  ports:
  - protocol: TCP
    port: 8001
    targetPort: 8001
  type: ClusterIP
```

#### Deploy to Kubernetes

```bash
# Apply configurations
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f user-service-deployment.yaml
kubectl apply -f user-service-service.yaml

# Check deployment
kubectl get deployments -n microservices
kubectl get pods -n microservices
kubectl get services -n microservices

# View logs
kubectl logs -n microservices -f deployment/user-service

# Scale deployment
kubectl scale deployment/user-service --replicas=5 -n microservices
```

## CI/CD Pipeline

### GitHub Actions Example

```yaml
# .github/workflows/deploy.yml
name: Build and Deploy

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.23'

    - name: Build services
      run: |
        cd service-a && go build -o user-service
        cd ../service-b && go build -o order-service
        cd ../api-gateway && go build -o api-gateway

    - name: Run tests
      run: |
        cd service-a && go test ./...
        cd ../service-b && go test ./...
        cd ../api-gateway && go test ./...

    - name: Build Docker images
      run: |
        docker build -t user-service:${{ github.sha }} ./service-a
        docker build -t order-service:${{ github.sha }} ./service-b
        docker build -t api-gateway:${{ github.sha }} ./api-gateway

    - name: Push to registry
      run: |
        echo ${{ secrets.REGISTRY_PASSWORD }} | docker login -u ${{ secrets.REGISTRY_USER }} --password-stdin
        docker push user-service:${{ github.sha }}
        docker push order-service:${{ github.sha }}
        docker push api-gateway:${{ github.sha }}

    - name: Deploy to production
      run: |
        # Update deployment with new image
        kubectl set image deployment/user-service user-service=user-service:${{ github.sha }} -n microservices
        kubectl set image deployment/order-service order-service=order-service:${{ github.sha }} -n microservices
        kubectl set image deployment/api-gateway api-gateway=api-gateway:${{ github.sha }} -n microservices
```

### GitLab CI Example

```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - deploy

build:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE/user-service:$CI_COMMIT_SHA ./service-a
    - docker build -t $CI_REGISTRY_IMAGE/order-service:$CI_COMMIT_SHA ./service-b
    - docker build -t $CI_REGISTRY_IMAGE/api-gateway:$CI_COMMIT_SHA ./api-gateway
    - docker push $CI_REGISTRY_IMAGE/user-service:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE/order-service:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE/api-gateway:$CI_COMMIT_SHA

test:
  stage: test
  script:
    - cd service-a && go test ./...
    - cd ../service-b && go test ./...
    - cd ../api-gateway && go test ./...

deploy:
  stage: deploy
  script:
    - kubectl set image deployment/user-service user-service=$CI_REGISTRY_IMAGE/user-service:$CI_COMMIT_SHA
    - kubectl set image deployment/order-service order-service=$CI_REGISTRY_IMAGE/order-service:$CI_COMMIT_SHA
    - kubectl set image deployment/api-gateway api-gateway=$CI_REGISTRY_IMAGE/api-gateway:$CI_COMMIT_SHA
  only:
    - main
```

## Monitoring and Logging

### Health Checks

```bash
# Check all service health
curl http://localhost:8080/health

# Check specific service health
curl http://localhost:8001/health
curl http://localhost:8002/health
```

### Log Aggregation

```bash
# View all logs
docker-compose logs

# Follow logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f service-a

# Export logs
docker-compose logs > logs.txt
```

### Metrics Collection

Integrate Prometheus for metrics:

```go
// Add to service
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
}
```

### Alerting

Set up alerts for:
- High error rates (>5%)
- High latency (>500ms)
- Service downtime
- Database connection issues

## Troubleshooting

### Common Issues

#### Service Won't Start

```bash
# Check logs
docker-compose logs service-a

# Check container status
docker-compose ps

# Restart service
docker-compose restart service-a

# Rebuild service
docker-compose up -d --build service-a
```

#### Database Connection Issues

```bash
# Check database is running
docker-compose ps users-db

# Check database logs
docker-compose logs users-db

# Test connection
docker-compose exec service-a ping users-db

# Check network
docker network inspect microservices_default
```

#### Port Conflicts

```yaml
# Change external ports in docker-compose.yml
services:
  service-a:
    ports:
      - "8002:8001"  # Use different external port
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Check database performance
docker-compose exec users-db psql -U postgres -d users_db -c "SELECT * FROM pg_stat_activity;"

# Check Redis performance
docker-compose exec redis redis-cli INFO stats
```

### Debug Mode

Enable debug logging:

```yaml
# docker-compose.yml
services:
  service-a:
    environment:
      - LOG_LEVEL=debug
```

### Recovery Procedures

```bash
# Complete restart
docker-compose down
docker-compose up -d

# Reset data (WARNING: deletes all data)
docker-compose down -v
docker-compose up -d

# Restore from backup
docker run --rm -v users-db-data:/data -v $(pwd):/backup alpine tar xzf /backup/users-db-backup.tar.gz -C /
```

## Best Practices

1. **Use Environment Variables**: Never hardcode configuration
2. **Implement Health Checks**: Always include health endpoints
3. **Log Everything**: Use structured logging
4. **Monitor Performance**: Track latency, errors, throughput
5. **Automate Deployments**: Use CI/CD pipelines
6. **Test Before Deploying**: Run tests in CI pipeline
7. **Use Secrets Management**: Never commit credentials
8. **Implement Backups**: Regular database backups
9. **Document Everything**: Keep deployment docs updated
10. **Plan for Failures**: Implement circuit breakers, retries
