# Exercise: Kubernetes for AI

## Problem 1: Kubernetes Resources

Match each Kubernetes resource with its purpose:

| Resource | Purpose |
|----------|---------|
| Deployment | |
| Service | |
| Pod | |
| ConfigMap | |
| Secret | |
| HorizontalPodAutoscaler | |
| PodDisruptionBudget | |
| ResourceQuota | |

---

## Problem 2: Create vLLM Deployment

Write a complete Kubernetes Deployment for vLLM:

```yaml

```

Requirements:
- 2 replicas
- GPU allocation (nvidia.com/gpu: 1)
- Memory limits (16Gi)
- Health probes
- Environment variables for model configuration

---

## Problem 3: HPA Configuration

Configure HPA for the vLLM deployment:

1. Scale between 1 and 10 replicas
2. Target 70% CPU utilization
3. Add scale-up/scale-down policies

```yaml

```

---

## Problem 4: Helm Chart Values

Create a values.yaml for an LLM deployment Helm chart:

```yaml
# Default configuration values

replicaCount:

image:
  repository:
  tag:
  pullPolicy:

service:
  type:
  port:
  targetPort:

resources:
  limits:
    nvidia.com/gpu:
    memory:
  requests:

env:
  MODEL:
  GPU_MEMORY_UTILIZATION:

autoscaling:
  enabled:
  minReplicas:
  maxReplicas:
  targetCPUUtilizationPercentage:
```

---

## Problem 5: API Gateway Deployment

Write Kubernetes manifests to deploy the API Gateway from `services/api-gateway/`:

1. ConfigMap for configuration
2. Secret for sensitive data
3. Deployment with 3 replicas
4. Service
5. HPA

```yaml
# ConfigMap

---

# Secret

---

# Deployment

---

# Service

---

# HPA

```

---

## Problem 6: Troubleshooting

A vLLM pod is failing to start with the following error:

```
0/3 nodes are available: 1 insufficient nvidia.com/gpu, 2 node(s) didn't match pod anti-affinity
```

1. What are the possible causes?
2. How would you troubleshoot each?
3. What changes would you make to the manifest?

---

## Problem 7: Production Checklist

Create a production readiness checklist for LLM deployment on Kubernetes:

| Check | Status |
|-------|--------|
| GPU resources configured | |
| Memory limits set | |
| Health probes configured | |
| Liveness probe | |
| Readiness probe | |
| PodDisruptionBudget | |
| ResourceQuota | |
| NetworkPolicy | |
| Auto-scaling configured | |
| Secrets management | |
| Logging configured | |
| Monitoring configured | |

---

## Submission

Save your answers and be prepared to discuss them in the next lesson.
