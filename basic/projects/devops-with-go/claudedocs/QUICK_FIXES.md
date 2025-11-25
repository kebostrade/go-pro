# Quick Fix Guide - DevOps Configuration Issues

**Project**: devops-with-go
**Priority**: CRITICAL issues first, then IMPORTANT

---

## CRITICAL FIXES (Do These First)

### 1. Remove Secrets from Git (C8)

```bash
# Immediate action required
cd /home/dima/Desktop/FUN/go-pro/basic/projects/devops-with-go

# Remove secret file from Git
git rm kubernetes/secret.yaml
echo "kubernetes/secret.yaml" >> .gitignore

# Create secret manually instead
kubectl create secret generic app-secrets \
  --from-literal=POSTGRES_USER=devops_user \
  --from-literal=POSTGRES_PASSWORD='change-this-password' \
  --from-literal=JWT_SECRET='change-this-jwt-secret' \
  -n devops-demo

# If already committed, purge from history
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch kubernetes/secret.yaml' \
  --prune-empty --tag-name-filter cat -- --all
```

---

### 2. Fix Docker Compose Secrets (C4, C5)

Create `docker/.env` file:
```env
POSTGRES_USER=devops_user
POSTGRES_PASSWORD=strong_password_here
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=change_this_password
```

Update `docker/docker-compose.yml`:
```yaml
services:
  app:
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

  postgres:
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

  grafana:
    environment:
      - GF_SECURITY_ADMIN_USER=${GRAFANA_ADMIN_USER:-admin}
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
```

Add to `.gitignore`:
```
docker/.env
```

---

### 3. Add Security Context to Deployment (C6)

Update `kubernetes/deployment.yaml` at line 23:
```yaml
spec:
  template:
    spec:
      # Add this pod security context
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
        runAsGroup: 65534
        fsGroup: 65534
        seccompProfile:
          type: RuntimeDefault

      containers:
      - name: app
        # Add this container security context
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: false  # Set to true if app doesn't write to disk
          runAsNonRoot: true
          runAsUser: 65534
          capabilities:
            drop:
              - ALL
            add:
              - NET_BIND_SERVICE
```

---

### 4. Fix Dockerfile Non-Root User (C1)

Update `docker/Dockerfile`:

**Builder stage** (add after line 5):
```dockerfile
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user
RUN addgroup -g 65534 -S nonroot && \
    adduser -u 65534 -S nonroot -G nonroot

WORKDIR /build
# ... rest of build
```

**Runtime stage** (update line 29 onwards):
```dockerfile
FROM scratch

# Copy user/group files
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app /app

# Switch to non-root user
USER 65534:65534

EXPOSE 8080

ENV PORT=8080

# Remove HEALTHCHECK (doesn't work with scratch image)

ENTRYPOINT ["/app"]
```

---

### 5. Pin Image Versions (C7)

Update `kubernetes/deployment.yaml` line 26:
```yaml
containers:
- name: app
  image: your-registry.example.com/devops-go-app:v1.0.0  # Use semantic versioning
  imagePullPolicy: Always  # Or IfNotPresent with versioned tags
```

---

## IMPORTANT FIXES (Do These Next)

### 6. Add Resource Limits to Docker Compose (I8)

Add to each service in `docker/docker-compose.yml`:
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  postgres:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M

  redis:
    deploy:
      resources:
        limits:
          cpus: '0.25'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 128M

  prometheus:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  grafana:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

---

### 7. Fix Kubernetes Resource Limits (I15)

Update `kubernetes/deployment.yaml` lines 58-64:
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

---

### 8. Add All ConfigMap Variables (I16)

Update `kubernetes/deployment.yaml` env section:
```yaml
env:
- name: PORT
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: PORT
- name: APP_VERSION
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: APP_VERSION
- name: ENVIRONMENT
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: ENVIRONMENT
- name: LOG_LEVEL
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: LOG_LEVEL
- name: POSTGRES_HOST
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: POSTGRES_HOST
- name: POSTGRES_PORT
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: POSTGRES_PORT
- name: POSTGRES_DB
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: POSTGRES_DB
- name: POSTGRES_USER
  valueFrom:
    secretKeyRef:
      name: app-secrets
      key: POSTGRES_USER
- name: POSTGRES_PASSWORD
  valueFrom:
    secretKeyRef:
      name: app-secrets
      key: POSTGRES_PASSWORD
- name: REDIS_HOST
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: REDIS_HOST
- name: REDIS_PORT
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: REDIS_PORT
- name: JWT_SECRET
  valueFrom:
    secretKeyRef:
      name: app-secrets
      key: JWT_SECRET
```

---

### 9. Add PodDisruptionBudget (I18)

Create `kubernetes/pdb.yaml`:
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: devops-go-pdb
  namespace: devops-demo
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: devops-go-app
```

---

### 10. Add Rolling Update Strategy (I19)

Update `kubernetes/deployment.yaml` after line 10:
```yaml
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: devops-go-app
```

---

### 11. Add Rate Limiting to Ingress (I24)

Update `kubernetes/ingress.yaml` annotations:
```yaml
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/limit-rps: "100"
    nginx.ingress.kubernetes.io/limit-connections: "50"
    nginx.ingress.kubernetes.io/limit-burst-multiplier: "5"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
```

---

### 12. Add Security Headers (I25)

Update `kubernetes/ingress.yaml` annotations:
```yaml
metadata:
  annotations:
    # ... existing annotations ...
    nginx.ingress.kubernetes.io/configuration-snippet: |
      more_set_headers "X-Frame-Options: DENY";
      more_set_headers "X-Content-Type-Options: nosniff";
      more_set_headers "X-XSS-Protection: 1; mode=block";
      more_set_headers "Strict-Transport-Security: max-age=31536000; includeSubDomains";
      more_set_headers "Content-Security-Policy: default-src 'self'";
```

---

### 13. Remove Exposed Ports (I10)

Update `docker/docker-compose.yml`:
```yaml
postgres:
  # Remove these lines:
  # ports:
  #   - "5432:5432"
  # Access via docker network only

redis:
  # Remove these lines:
  # ports:
  #   - "6379:6379"

prometheus:
  ports:
    - "127.0.0.1:9090:9090"  # Bind to localhost only

grafana:
  ports:
    - "127.0.0.1:3000:3000"  # Bind to localhost only
```

---

## NEW FILES TO CREATE

### 14. NetworkPolicy

Create `kubernetes/networkpolicy.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: devops-go-netpol
  namespace: devops-demo
spec:
  podSelector:
    matchLabels:
      app: devops-go-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: UDP
      port: 53
```

---

### 15. Pod Security Standards

Update `kubernetes/namespace.yaml`:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: devops-demo
  labels:
    name: devops-demo
    environment: production
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

### 16. ResourceQuota

Create `kubernetes/resourcequota.yaml`:
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: devops-demo-quota
  namespace: devops-demo
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
    persistentvolumeclaims: "5"
    services.loadbalancers: "1"
```

---

## VALIDATION COMMANDS

After applying fixes, run these commands:

```bash
# Test Docker build
cd docker
docker build -t devops-go-app:test -f Dockerfile ..

# Validate Docker Compose
docker-compose config

# Validate Kubernetes manifests
kubectl apply --dry-run=client -f kubernetes/

# Check security contexts
kubectl get pods -n devops-demo -o json | \
  jq '.items[].spec.securityContext'

# Verify non-root user
kubectl exec -it -n devops-demo <pod-name> -- id

# Test resource limits
kubectl describe pod -n devops-demo <pod-name> | grep -A 5 "Limits"

# Validate network policies
kubectl get networkpolicies -n devops-demo

# Check HPA status
kubectl get hpa -n devops-demo

# Verify secrets not in Git
git log --all --full-history -- kubernetes/secret.yaml
```

---

## TESTING CHECKLIST

After implementing fixes:

- [ ] Docker image builds successfully
- [ ] Container runs as non-root (id shows uid=65534)
- [ ] Docker Compose starts all services
- [ ] Secrets loaded from .env file (not hardcoded)
- [ ] Kubernetes deployment applies without errors
- [ ] Pods start and pass health checks
- [ ] HPA scales based on load
- [ ] Network policies restrict traffic
- [ ] Resource limits prevent runaway processes
- [ ] Ingress routing works with rate limiting
- [ ] No secrets visible in `kubectl get secret -o yaml`
- [ ] Security scan passes (trivy/snyk)

---

## IMMEDIATE ACTION ITEMS

**Before next deployment:**

1. **CRITICAL**: Remove `kubernetes/secret.yaml` from Git
2. **CRITICAL**: Change all default passwords
3. **CRITICAL**: Add security contexts to all pods
4. **HIGH**: Add resource limits to all containers
5. **HIGH**: Create NetworkPolicy

**Time estimate**: 2-4 hours for critical fixes

---

**Need Help?**

- Full analysis: See `DEVOPS_ANALYSIS.md`
- Docker issues: Check Docker section
- Kubernetes issues: Check Kubernetes section
- Security questions: See Security Recommendations section
