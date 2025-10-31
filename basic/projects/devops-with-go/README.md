# 🐳 DevOps with Go: Docker, Kubernetes, and Terraform

A comprehensive tutorial demonstrating modern DevOps practices with Go, including containerization with Docker, orchestration with Kubernetes, and infrastructure as code with Terraform.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Docker](#docker)
- [Kubernetes](#kubernetes)
- [Terraform](#terraform)
- [Monitoring](#monitoring)
- [Project Structure](#project-structure)

## 🎯 Overview

This project demonstrates:
- **Docker**: Multi-stage builds, Docker Compose, container networking
- **Kubernetes**: Deployments, Services, Ingress, ConfigMaps, Secrets, HPA
- **Terraform**: Infrastructure as Code for AWS ECS deployment
- **Observability**: Prometheus metrics, health checks, logging
- **Production Practices**: Graceful shutdown, resource limits, auto-scaling

## ✨ Features

### Application Features
- ✅ RESTful API with health checks
- ✅ Prometheus metrics endpoint
- ✅ Graceful shutdown
- ✅ Environment-based configuration
- ✅ Structured logging

### Docker Features
- ✅ Multi-stage builds for minimal image size
- ✅ Docker Compose with PostgreSQL, Redis, Prometheus, Grafana
- ✅ Health checks
- ✅ Development and production Dockerfiles

### Kubernetes Features
- ✅ Deployment with 3 replicas
- ✅ Liveness and readiness probes
- ✅ ConfigMaps and Secrets
- ✅ Horizontal Pod Autoscaler (HPA)
- ✅ Ingress with TLS
- ✅ Resource requests and limits

### Terraform Features
- ✅ AWS VPC with public subnets
- ✅ ECS Fargate cluster
- ✅ Application Load Balancer
- ✅ ECR repository
- ✅ CloudWatch logging
- ✅ Auto-scaling configuration

## 📦 Prerequisites

- **Go** 1.21 or higher
- **Docker** 20.10 or higher
- **Docker Compose** 2.0 or higher
- **Kubernetes** (minikube, kind, or cloud provider)
- **kubectl** 1.25 or higher
- **Terraform** 1.0 or higher
- **AWS CLI** (for Terraform deployment)

## 🚀 Quick Start

### 1. Clone and Navigate

```bash
cd basic/projects/devops-with-go
```

### 2. Run Locally

```bash
# Install dependencies
go mod tidy

# Run the application
make run

# Test endpoints
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## 🐳 Docker

### Build Docker Image

```bash
# Build production image
make docker-build

# Or manually
docker build -t devops-go-app:latest -f docker/Dockerfile .
```

### Run with Docker Compose

```bash
# Start all services
make docker-run

# Services available:
# - App: http://localhost:8080
# - Prometheus: http://localhost:9090
# - Grafana: http://localhost:3000 (admin/admin)
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379

# Stop services
make docker-stop
```

### Docker Image Details

**Multi-stage build:**
- **Stage 1 (Builder)**: Compiles Go application
- **Stage 2 (Runtime)**: Minimal scratch image with binary only

**Image size:** ~10MB (compared to ~800MB with full Go image)

## ☸️ Kubernetes

### Deploy to Kubernetes

```bash
# Deploy all resources
make k8s-deploy

# Or manually
kubectl apply -f kubernetes/namespace.yaml
kubectl apply -f kubernetes/configmap.yaml
kubectl apply -f kubernetes/secret.yaml
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml
kubectl apply -f kubernetes/ingress.yaml
kubectl apply -f kubernetes/hpa.yaml
```

### Check Deployment

```bash
# Check pods
kubectl get pods -n devops-demo

# Check services
kubectl get svc -n devops-demo

# Check HPA
kubectl get hpa -n devops-demo

# View logs
kubectl logs -f deployment/devops-go-app -n devops-demo

# Port forward to access locally
kubectl port-forward svc/devops-go-service 8080:80 -n devops-demo
```

### Kubernetes Resources

| Resource | Description |
|----------|-------------|
| **Namespace** | Isolated environment for resources |
| **ConfigMap** | Non-sensitive configuration |
| **Secret** | Sensitive data (passwords, tokens) |
| **Deployment** | 3 replicas with rolling updates |
| **Service** | LoadBalancer exposing port 80 |
| **Ingress** | HTTPS with TLS termination |
| **HPA** | Auto-scaling 2-10 pods based on CPU/memory |

### Health Checks

**Liveness Probe:**
- Endpoint: `/health/live`
- Checks if app is running
- Restarts pod if failing

**Readiness Probe:**
- Endpoint: `/health/ready`
- Checks if app can serve traffic
- Removes from load balancer if failing

## 🏗️ Terraform

### Initialize Terraform

```bash
make terraform-init

# Or manually
cd terraform && terraform init
```

### Plan Infrastructure

```bash
make terraform-plan
```

### Deploy to AWS

```bash
# Apply changes
make terraform-apply

# Infrastructure created:
# - VPC with 2 public subnets
# - Internet Gateway
# - Security Groups
# - Application Load Balancer
# - ECS Fargate Cluster
# - ECR Repository
# - CloudWatch Log Groups
# - IAM Roles and Policies
```

### Access Application

```bash
# Get load balancer DNS
terraform output application_url

# Example: http://devops-go-app-alb-123456789.us-east-1.elb.amazonaws.com
```

### Destroy Infrastructure

```bash
make terraform-destroy
```

### Terraform Resources

```
AWS Resources Created:
├── VPC (10.0.0.0/16)
│   ├── Public Subnet 1 (10.0.0.0/24)
│   ├── Public Subnet 2 (10.0.1.0/24)
│   └── Internet Gateway
├── Security Groups
│   └── App SG (80, 443, 8080)
├── Load Balancer
│   ├── Target Group (health checks)
│   └── Listener (port 80)
├── ECS
│   ├── Cluster
│   ├── Task Definition (Fargate)
│   └── Service (2 tasks)
├── ECR Repository
└── CloudWatch Logs
```

## 📊 Monitoring

### Prometheus Metrics

Access Prometheus at `http://localhost:9090` (Docker Compose)

**Available metrics:**
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request duration histogram
- `app_info` - Application version and environment

**Example queries:**
```promql
# Request rate
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
```

### Grafana Dashboards

Access Grafana at `http://localhost:3000` (admin/admin)

1. Add Prometheus data source: `http://prometheus:9090`
2. Import dashboard or create custom visualizations
3. Monitor request rates, latencies, error rates

## 📁 Project Structure

```
devops-with-go/
├── app/
│   └── main.go                 # Go application with metrics
├── docker/
│   ├── Dockerfile              # Production multi-stage build
│   ├── Dockerfile.dev          # Development with hot reload
│   ├── docker-compose.yml      # Full stack with monitoring
│   └── prometheus.yml          # Prometheus configuration
├── kubernetes/
│   ├── namespace.yaml          # Namespace definition
│   ├── configmap.yaml          # Configuration data
│   ├── secret.yaml             # Sensitive data
│   ├── deployment.yaml         # Application deployment
│   ├── service.yaml            # Service (LoadBalancer)
│   ├── ingress.yaml            # Ingress with TLS
│   └── hpa.yaml                # Horizontal Pod Autoscaler
├── terraform/
│   ├── main.tf                 # Main infrastructure
│   ├── variables.tf            # Input variables
│   └── outputs.tf              # Output values
├── Makefile                    # Build automation
├── go.mod                      # Go dependencies
└── README.md                   # This file
```

## 🎓 Learning Outcomes

After completing this tutorial, you'll understand:

### Docker
- ✅ Multi-stage builds for optimized images
- ✅ Docker Compose for local development
- ✅ Container networking and volumes
- ✅ Health checks and restart policies

### Kubernetes
- ✅ Deployments and replica management
- ✅ Services and load balancing
- ✅ ConfigMaps and Secrets
- ✅ Liveness and readiness probes
- ✅ Horizontal Pod Autoscaling
- ✅ Ingress and TLS termination

### Terraform
- ✅ Infrastructure as Code principles
- ✅ AWS resource provisioning
- ✅ State management
- ✅ Variables and outputs
- ✅ ECS Fargate deployment

### Observability
- ✅ Prometheus metrics collection
- ✅ Grafana visualization
- ✅ Health check endpoints
- ✅ Structured logging

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `APP_VERSION` | Application version | `1.0.0` |
| `ENVIRONMENT` | Environment name | `development` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Go Best Practices](https://golang.org/doc/effective_go)

## 🎯 Next Steps

1. **Add CI/CD**: GitHub Actions, GitLab CI, or Jenkins
2. **Add Database**: PostgreSQL with migrations
3. **Add Caching**: Redis integration
4. **Add Tracing**: OpenTelemetry or Jaeger
5. **Add Security**: Vault for secrets, RBAC
6. **Multi-region**: Deploy across multiple AWS regions
7. **Service Mesh**: Istio or Linkerd

## 📝 License

MIT License - feel free to use this project for learning!

