# Kubernetes Cloud-Native Template

A production-grade Kubernetes deployment template with Helm charts and operator scaffolding using controller-runtime.

## Architecture Overview

This template provides:

- **Kubernetes Manifests**: Production-ready Deployment, Service, ConfigMap, and HPA
- **Helm Chart**: Templated Kubernetes configuration for reusable deployments
- **Operator Scaffolding**: CRD-based operator pattern using controller-runtime
- **Docker Image**: Multi-stage build optimized for small image size

## Prerequisites

- Go 1.23+
- Docker
- kubectl (for Kubernetes deployment)
- Helm 3.14+
- kubebuilder (for operator development)

## Project Structure

```
kubernetes/
├── manifests/                    # Plain Kubernetes manifests
│   ├── namespace.yaml          # Dedicated namespace
│   ├── deployment.yaml         # Production deployment with probes
│   ├── service.yaml           # ClusterIP service
│   ├── configmap.yaml         # Configuration
│   └── hpa.yaml               # Horizontal Pod Autoscaler
├── helm/
│   └── gopro-chart/           # Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       ├── .helmignore
│       └── templates/
│           ├── _helpers.tpl
│           ├── deployment.yaml
│           └── service.yaml
├── api/
│   └── v1alpha1/             # CRD types
│       └── gopro_types.go
├── controllers/                # Operator controllers
│   ├── gopro_controller.go
│   └── suite_test.go
├── cmd/
│   └── main.go               # Operator entry point
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```

## Quick Start

### 1. Using Plain Manifests

Apply the Kubernetes manifests directly:

```bash
# Create namespace and apply manifests
kubectl apply -f manifests/

# Verify deployment
kubectl get pods -n gopro
kubectl get services -n gopro

# Check deployment status
kubectl rollout status deployment/gopro-app -n gopro
```

### 2. Using Helm Chart

Install using Helm:

```bash
# Add the chart repository (if applicable)
helm repo add gopro https://charts.gopro.example.com

# Or install from local chart
helm install gopro-release ./helm/gopro-chart --namespace gopro --create-namespace

# Upgrade
helm upgrade gopro-release ./helm/gopro-chart --namespace gopro

# Uninstall
helm uninstall gopro-release --namespace gopro
```

### 3. Customizing Values

Override default values:

```bash
helm install gopro-release ./helm/gopro-chart \
  --namespace gopro \
  --create-namespace \
  --set replicaCount=5 \
  --set image.tag=v1.2.0 \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=3 \
  --set autoscaling.maxReplicas=15
```

### 4. Running the Operator

For operator development:

```bash
# Install kubebuilder (first time)
make install-kubebuilder

# Run operator locally
make run

# Or build and run in Docker
make docker-build
make docker-run
```

## Development

### Building

```bash
# Build operator binary
make build

# Run tests
make test

# Run go vet
make vet

# Tidy modules
make mod-tidy
```

### Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### Helm

```bash
# Lint chart
make helm-lint

# Render templates
make helm-template

# Deploy with Helm
make helm-deploy
```

### Kubernetes Manifests

```bash
# Validate manifests (dry-run)
make manifest-validate

# Apply manifests
make manifest-apply

# Delete manifests
make manifest-delete
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | production | Application environment |
| `LOG_LEVEL` | info | Logging level |
| `PORT` | 8080 | Container port |
| `HEALTH_LIVE_PATH` | /health/live | Liveness probe path |
| `HEALTH_READY_PATH` | /health/ready | Readiness probe path |

### Helm Values

See `helm/gopro-chart/values.yaml` for all configurable options:

- `replicaCount`: Number of pods (default: 3)
- `image.repository`: Container image
- `image.tag`: Image tag
- `service.type`: Service type (ClusterIP, LoadBalancer, etc.)
- `resources`: CPU/memory requests and limits
- `autoscaling`: HPA configuration
- `probes`: Liveness, readiness, and startup probes

## Custom Resources

The operator defines a `GoPro` CRD:

```yaml
apiVersion: gopro.example.com/v1alpha1
kind: GoPro
metadata:
  name: my-app
spec:
  image: gopro-app:latest
  replicas: 3
  env:
    APP_ENV: production
  port: 8080
```

## CI/CD

GitHub Actions workflow is configured in `.github/workflows/ci.yml`:

- Go 1.23 build and test
- Helm lint and template validation
- Docker build

## Health Checks

The deployment includes:

- **Liveness Probe**: `/health/live` - Checks if container is alive
- **Readiness Probe**: `/health/ready` - Checks if container can serve traffic
- **Startup Probe**: Initial delay for slow-starting applications

## Security

- Runs as non-root user (UID 1000)
- Read-only root filesystem recommended
- Security context configured

## Troubleshooting

### Pods not starting

```bash
kubectl describe pod <pod-name> -n gopro
kubectl logs <pod-name> -n gopro
```

### Service not accessible

```bash
kubectl get endpoints -n gopro
kubectl describe service gopro-service -n gopro
```

### Helm issues

```bash
helm history gopro-release -n gopro
helm rollback gopro-release -n gopro
```

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [kubebuilder](https://book.kubebuilder.io/)
