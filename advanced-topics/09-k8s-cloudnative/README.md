# Kubernetes & Cloud-Native Go

Learn to deploy and manage Go applications in Kubernetes with production-ready configurations.

## Overview

This module covers:
- **Containerization**: Building Docker images for Go apps
- **Kubernetes Deployment**: Deployments, Services, ConfigMaps, Secrets
- **Health Checks**: Liveness and readiness probes
- **Scaling**: Horizontal Pod Autoscaling (HPA)
- **Traffic Management**: Ingress for external access
- **Configuration Management**: ConfigMaps and Secrets
- **Resource Management**: CPU and memory limits

## Prerequisites

```bash
# Install kubectl (Kubernetes CLI)
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Minikube (local K8s cluster)
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube-linux-amd64
sudo mv minikube-linux-amd64 /usr/local/bin/minikube

# Start Minikube
minikube start
```

## Quick Start

### 1. Build the Go Application

```bash
cd sample-app
go mod init k8s-sample
go mod tidy
go build -o app
```

### 2. Build Docker Image

```bash
# Build image
docker build -t k8s-go-sample:v1 .

# If using Minikube, load image into Minikube
minikube image load k8s-go-sample:v1
```

### 3. Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f ingress.yaml
kubectl apply -f hpa.yaml

# Check deployment status
kubectl get pods
kubectl get services
kubectl get ingress
```

### 4. Access the Application

```bash
# If using Minikube
minikube tunnel  # Run in separate terminal for Ingress

# Access via Ingress (after tunnel)
curl http://k8s-go-sample.local/

# Or port-forward directly
kubectl port-forward svc/k8s-go-sample 8080:80
curl http://localhost:8080/
```

## Kubernetes Concepts

### Deployment

**deployment.yaml** manages:
- **ReplicaSet**: Ensures specified number of pod replicas
- **Pod Rolling Updates**: Zero-downtime deployments
- **Rollback**: Easy rollback to previous versions

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-go-sample
spec:
  replicas: 3
  selector:
    matchLabels:
      app: k8s-go-sample
  template:
    metadata:
      labels:
        app: k8s-go-sample
    spec:
      containers:
      - name: app
        image: k8s-go-sample:v1
        ports:
        - containerPort: 8080
```

### Service

**service.yaml** provides:
- **Service Discovery**: Internal DNS for pods
- **Load Balancing**: Distributes traffic across pods
- **Stable Endpoints**: Abstracts dynamic pod IPs

```yaml
apiVersion: v1
kind: Service
metadata:
  name: k8s-go-sample
spec:
  selector:
    app: k8s-go-sample
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### ConfigMap

**configmap.yaml** stores:
- **Configuration Data**: Environment-specific settings
- **Feature Flags**: Toggle features without rebuilding
- **Application Config**: Non-sensitive configuration

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
```

### Secret

**secret.yaml** manages:
- **Sensitive Data**: Passwords, API keys, tokens
- **Encoded Data**: Base64 encoding
- **Mount as Files**: Or inject as environment variables

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  API_KEY: YWRtaW4tc2VjcmV0LWtleQ==
```

### Ingress

**ingress.yaml** enables:
- **External Access**: HTTP/HTTPS routing to services
- **Host-based Routing**: Multiple services on same IP
- **TLS Termination**: SSL/TLS at ingress level

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8s-go-sample-ingress
spec:
  rules:
  - host: k8s-go-sample.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k8s-go-sample
            port:
              number: 80
```

### Horizontal Pod Autoscaler

**hpa.yaml** provides:
- **Auto-scaling**: Scale pods based on CPU/memory
- **Metric-based**: CPU, memory, or custom metrics
- **Cost Optimization**: Scale down when load decreases

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: k8s-go-sample-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: k8s-go-sample
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 80
```

## Health Checks

Kubernetes uses two types of probes:

### Liveness Probe
- **Purpose**: Detect and restart deadlocked pods
- **Action**: Restarts pod if fails
- **Endpoint**: `/health/live`

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 20
```

### Readiness Probe
- **Purpose**: Check if pod can receive traffic
- **Action**: Removes from service if fails
- **Endpoint**: `/health/ready`

```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## Rolling Updates

```bash
# Update image
kubectl set image deployment/k8s-go-sample app=k8s-go-sample:v2

# Watch rollout status
kubectl rollout status deployment/k8s-go-sample

# Check rollout history
kubectl rollout history deployment/k8s-go-sample

# Rollback if needed
kubectl rollout undo deployment/k8s-go-sample
kubectl rollout undo deployment/k8s-go-sample --to-revision=2
```

## Scaling

```bash
# Manual scaling
kubectl scale deployment/k8s-go-sample --replicas=5

# Check HPA status
kubectl get hpa

# Check resource usage
kubectl top pods
kubectl top nodes
```

## Debugging

```bash
# View pod logs
kubectl logs -f deployment/k8s-go-sample

# Execute command in pod
kubectl exec -it <pod-name> -- /bin/sh

# Describe pod for details
kubectl describe pod <pod-name>

# Get events
kubectl get events --sort-by=.metadata.creationTimestamp
```

## Best Practices

### 1. Resource Limits
Always set resource requests and limits:

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "250m"
  limits:
    memory: "128Mi"
    cpu: "500m"
```

### 2. Health Checks
Implement both liveness and readiness probes

### 3. Configuration
Use ConfigMaps for config, Secrets for sensitive data

### 4. Images
Use specific version tags (not `latest`)

### 5. Security
Run as non-root user when possible

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
```

### 6. Probes
Set appropriate periods and thresholds

### 7. Replicas
Run at least 2 replicas for high availability

## Cleanup

```bash
# Delete all resources
kubectl delete -f ingress.yaml
kubectl delete -f hpa.yaml
kubectl delete -f secret.yaml
kubectl delete -f configmap.yaml
kubectl delete -f service.yaml
kubectl delete -f deployment.yaml

# Or delete by label
kubectl delete all -l app=k8s-go-sample
```

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
- [Docker for Go](https://docs.docker.com/language/golang/)

## Next Steps

After completing this module:
1. Explore Helm for package management
2. Learn about service meshes (Istio, Linkerd)
3. Study monitoring (Prometheus, Grafana)
4. Understand CI/CD with Kubernetes (ArgoCD, Flux)
