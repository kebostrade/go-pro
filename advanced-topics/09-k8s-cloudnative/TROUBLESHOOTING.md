# Kubernetes Troubleshooting Guide

Common issues and solutions for Kubernetes deployments.

## Common Issues

### 1. Image Pull Errors

**Error**: `ErrImagePull` or `ImagePullBackOff`

**Causes**:
- Image doesn't exist in registry
- Incorrect image name/tag
- Authentication issues with private registry

**Solutions**:

```bash
# Check pod status
kubectl describe pod <pod-name>

# If using Minikube, load image
minikube image load k8s-go-sample:v1

# Or use Minikube's Docker environment
eval $(minikube docker-env)
docker build -t k8s-go-sample:v1 .

# Verify image exists
docker images | grep k8s-go-sample
```

### 2. CrashLoopBackOff

**Error**: Pod keeps restarting

**Debug**:

```bash
# View pod logs
kubectl logs <pod-name>

# View logs from previous restart
kubectl logs <pod-name> --previous

# Describe pod for events
kubectl describe pod <pod-name>

# Execute into pod
kubectl exec -it <pod-name> -- /bin/sh
```

**Common Causes**:
- Application crashes on startup
- Missing environment variables
- Failed health checks
- Missing dependencies

### 3. Service Not Accessible

**Symptoms**: Can't reach service from external

**Debug**:

```bash
# Check service
kubectl get svc

# Describe service
kubectl describe svc <service-name>

# Check endpoints (should have pod IPs)
kubectl get endpoints <service-name>

# Test from within cluster
kubectl run test-pod --rm -it --image=busybox -- sh
wget -O- http://service-name:port
```

**Solutions**:

```bash
# Port-forward to test
kubectl port-forward svc/<service-name> 8080:80

# Check Ingress
kubectl get ingress
kubectl describe ingress <ingress-name>

# If using Minikube, run tunnel
minikube tunnel
```

### 4. Health Check Failures

**Symptoms**: Pod marked unhealthy

**Debug**:

```bash
# Describe pod to see probe failures
kubectl describe pod <pod-name>

# Check probe configuration
kubectl get pod <pod-name> -o yaml | grep -A 10 livenessProbe
kubectl get pod <pod-name> -o yaml | grep -A 10 readinessProbe
```

**Solutions**:

- Adjust `initialDelaySeconds` if app takes longer to start
- Increase `periodSeconds` for less frequent checks
- Verify endpoints are correct
- Check app logs for errors during health checks

### 5. Resource Issues

**Symptoms**: OOMKilled, CPU throttling

**Debug**:

```bash
# Check pod events
kubectl describe pod <pod-name> | grep -A 20 Events

# Check resource usage
kubectl top pods
kubectl top nodes

# View resource limits
kubectl get pod <pod-name> -o yaml | grep -A 5 resources
```

**Solutions**:

```yaml
# Increase limits in deployment.yaml
resources:
  requests:
    memory: "128Mi"   # Increase from 64Mi
    cpu: "500m"       # Increase from 250m
  limits:
    memory: "256Mi"   # Increase from 128Mi
    cpu: "1000m"      # Increase from 500m
```

## Diagnostic Commands

### Pod Status

```bash
# List all pods
kubectl get pods

# List pods with more info
kubectl get pods -o wide

# Watch pod status
kubectl get pods -w

# Get pod YAML
kubectl get pod <pod-name> -o yaml

# Get pod as JSON
kubectl get pod <pod-name> -o json
```

### Logs

```bash
# Follow logs
kubectl logs -f deployment/<app-name>

# Logs from all pods
kubectl logs -l app=<app-name> --all-containers=true

# Specific container
kubectl logs <pod-name> -c <container-name>

# Previous logs (before restart)
kubectl logs <pod-name> --previous
```

### Events

```bash
# List recent events
kubectl get events --sort-by=.metadata.creationTimestamp

# Watch events
kubectl get events -w

# Events for specific namespace
kubectl get events -n <namespace>
```

### Describe

```bash
# Describe pod
kubectl describe pod <pod-name>

# Describe service
kubectl describe svc <service-name>

# Describe deployment
kubectl describe deployment <deployment-name>
```

### Executing Commands

```bash
# Execute into pod
kubectl exec -it <pod-name> -- /bin/sh

# Run single command
kubectl exec <pod-name> -- ls /

# Copy files from pod
kubectl cp <pod-name>:/path/to/file ./local-file

# Copy files to pod
kubectl cp ./local-file <pod-name>:/path/to/file
```

## Network Debugging

### Service Connectivity

```bash
# Check service endpoints
kubectl get endpoints <service-name>

# DNS resolution
kubectl run test-dns --rm -it --image=busybox -- nslookup <service-name>

# Test service from within cluster
kubectl run test-curl --rm -it --image=curlimages/curl -- curl http://<service-name>
```

### Ingress Issues

```bash
# Check ingress controller
kubectl get pods -n ingress-nginx

# Check ingress rules
kubectl describe ingress <ingress-name>

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

## Deployment Issues

### Rollout Problems

```bash
# Check rollout status
kubectl rollout status deployment/<deployment-name>

# View rollout history
kubectl rollout history deployment/<deployment-name>

# Pause rollout
kubectl rollout pause deployment/<deployment-name>

# Resume rollout
kubectl rollout resume deployment/<deployment-name>

# Undo rollout
kubectl rollout undo deployment/<deployment-name>

# Undo to specific revision
kubectl rollout undo deployment/<deployment-name> --to-revision=2
```

### Scaling Issues

```bash
# Check HPA status
kubectl get hpa

# Describe HPA
kubectl describe hpa <hpa-name>

# Check metrics server
kubectl get pods -n kube-system | grep metrics-server

# Manual scaling
kubectl scale deployment/<deployment-name> --replicas=5
```

## Configuration Issues

### ConfigMap/Secret Problems

```bash
# List configmaps
kubectl get configmaps

# Describe configmap
kubectl describe configmap <configmap-name>

# Get configmap value
kubectl get configmap <configmap-name> -o jsonpath='{.data.<key>}'

# List secrets
kubectl get secrets

# Decode secret
kubectl get secret <secret-name> -o jsonpath='{.data.<key>}' | base64 -d
```

### Environment Variables

```bash
# Check environment in pod
kubectl exec <pod-name> -- env

# Get specific env var
kubectl exec <pod-name> -- printenv <VAR_NAME>
```

## Performance Issues

### Resource Monitoring

```bash
# Top pods
kubectl top pods

# Top nodes
kubectl top nodes

# Resource usage by container
kubectl exec <pod-name> -- top
```

### Metrics Server

```bash
# Check if metrics-server is installed
kubectl get apiservice | grep metrics

# Install metrics-server (if not present)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Verify metrics
kubectl get --raw /apis/metrics.k8s.io/v1beta1/nodes
```

## Cleanup Commands

```bash
# Delete specific resources
kubectl delete deployment <deployment-name>
kubectl delete service <service-name>
kubectl delete ingress <ingress-name>
kubectl delete configmap <configmap-name>
kubectl delete secret <secret-name>
kubectl delete hpa <hpa-name>

# Delete all resources with label
kubectl delete all -l app=<app-name>

# Delete namespace and all resources
kubectl delete namespace <namespace>

# Force delete pod
kubectl delete pod <pod-name> --force --grace-period=0
```

## Best Practices

1. **Always use labels** for easy resource management
2. **Set resource limits** to prevent resource exhaustion
3. **Implement health checks** for automatic recovery
4. **Use liveness and readiness probes** appropriately
5. **Configure HPA** for auto-scaling
6. **Use namespaces** to organize resources
7. **Name resources consistently**
8. **Add annotations** for metadata
9. **Use multiple replicas** for high availability
10. **Monitor logs and metrics** regularly

## Getting Help

```bash
#kubectl help
kubectl --help

# Help for specific command
kubectl get --help
kubectl logs --help

# API documentation
kubectl api-resources
kubectl api-versions
```

## Useful Tools

### kubectl plugins

```bash
# Install krew (plugin manager)
# https://krew.sigs.k8s.io/

# List useful plugins
kubectl krew search
```

### Third-party tools

- **k9s**: Terminal UI for Kubernetes
- **kubectx**: Switch between clusters
- **kubens**: Switch between namespaces
- **stern**: Multi-pod log tailing
- **kubefwd**: Port forward multiple services
