# Fixed Configuration Files

Complete, production-ready configuration files with all security issues addressed.

---

## 1. Fixed Dockerfile

**File**: `docker/Dockerfile`

```dockerfile
# Multi-stage Docker build for Go application with security hardening
# Build arguments for version tracking
ARG GO_VERSION=1.21.13
ARG ALPINE_VERSION=3.19

# Stage 1: Build
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Build metadata
ARG APP_VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for runtime
RUN addgroup -g 65534 -S nonroot && \
    adduser -u 65534 -S nonroot -G nonroot

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download and verify dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

# Copy source code
COPY app/ ./app/

# Build the application with version info
# CGO_ENABLED=0 for static binary
# -ldflags for smaller binary size and version injection
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' \
              -X main.Version=${APP_VERSION} \
              -X main.GitCommit=${GIT_COMMIT} \
              -X main.BuildDate=${BUILD_DATE}" \
    -trimpath \
    -o /app \
    ./app

# Stage 2: Runtime (minimal scratch image)
FROM scratch

# Copy user/group files for non-root user
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app /app

# Use non-root user
USER 65534:65534

# Expose port
EXPOSE 8080

# Metadata labels
LABEL org.opencontainers.image.title="DevOps Go App" \
      org.opencontainers.image.description="Production-ready Go application" \
      org.opencontainers.image.version="${APP_VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${GIT_COMMIT}"

# Run the application
ENTRYPOINT ["/app"]
```

**Build command**:
```bash
docker build \
  --build-arg APP_VERSION=v1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t devops-go-app:v1.0.0 \
  -f docker/Dockerfile .
```

---

## 2. Fixed Docker Compose

**File**: `docker/docker-compose.yml`

```yaml
version: '3.8'

services:
  # Go Application
  app:
    build:
      context: ..
      dockerfile: docker/Dockerfile
      args:
        APP_VERSION: ${APP_VERSION:-v1.0.0}
        GIT_COMMIT: ${GIT_COMMIT:-dev}
        BUILD_DATE: ${BUILD_DATE:-unknown}
    container_name: devops-go-app
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - APP_VERSION=${APP_VERSION:-v1.0.0}
      - ENVIRONMENT=docker
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DB=${POSTGRES_DB:-devops_db}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - backend
      - frontend
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    healthcheck:
      test: ["CMD", "/app", "health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    security_opt:
      - no-new-privileges:true
    read_only: false
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE

  # PostgreSQL Database
  postgres:
    image: postgres:16.1-alpine
    container_name: devops-postgres
    # Remove port exposure for security - access via network only
    # ports:
    #   - "5432:5432"
    environment:
      - POSTGRES_DB=${POSTGRES_DB:-devops_db}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_INITDB_ARGS=--encoding=UTF-8 --lc-collate=C --lc-ctype=C
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - backend
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
    security_opt:
      - no-new-privileges:true
    tmpfs:
      - /tmp
      - /run

  # Redis Cache
  redis:
    image: redis:7.2-alpine
    container_name: devops-redis
    # Remove port exposure for security
    # ports:
    #   - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - backend
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.25'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 128M
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    security_opt:
      - no-new-privileges:true
    command: >
      redis-server
      --appendonly yes
      --maxmemory 200mb
      --maxmemory-policy allkeys-lru

  # Prometheus Monitoring
  prometheus:
    image: prom/prometheus:v2.48.0
    container_name: devops-prometheus
    ports:
      - "127.0.0.1:9090:9090"  # Bind to localhost only
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
      - '--web.enable-admin-api'
    networks:
      - monitoring
      - frontend
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    security_opt:
      - no-new-privileges:true
    user: "nobody"

  # Grafana Dashboards
  grafana:
    image: grafana/grafana:10.2.2
    container_name: devops-grafana
    ports:
      - "127.0.0.1:3000:3000"  # Bind to localhost only
    environment:
      - GF_SECURITY_ADMIN_USER=${GRAFANA_ADMIN_USER:-admin}
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=http://localhost:3000
      - GF_SECURITY_DISABLE_GRAVATAR=true
      - GF_ANALYTICS_REPORTING_ENABLED=false
      - GF_ANALYTICS_CHECK_FOR_UPDATES=false
    volumes:
      - grafana-data:/var/lib/grafana
    networks:
      - monitoring
      - frontend
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    depends_on:
      - prometheus
    security_opt:
      - no-new-privileges:true
    user: "472"  # Grafana user

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true  # No external access
  monitoring:
    driver: bridge

volumes:
  postgres-data:
    driver: local
  redis-data:
    driver: local
  prometheus-data:
    driver: local
  grafana-data:
    driver: local
```

**Required .env file** (`docker/.env`):
```env
# Application
APP_VERSION=v1.0.0
GIT_COMMIT=dev
BUILD_DATE=2025-10-31

# Database
POSTGRES_DB=devops_db
POSTGRES_USER=devops_user
POSTGRES_PASSWORD=change_this_strong_password

# Grafana
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=change_this_admin_password
```

---

## 3. Fixed Kubernetes Deployment

**File**: `kubernetes/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: devops-go-app
  namespace: devops-demo
  labels:
    app: devops-go-app
    version: v1
    component: backend
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
  template:
    metadata:
      labels:
        app: devops-go-app
        version: v1
        component: backend
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      # Pod security context
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
        runAsGroup: 65534
        fsGroup: 65534
        seccompProfile:
          type: RuntimeDefault
        supplementalGroups: []

      # Pod anti-affinity for HA
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - devops-go-app
              topologyKey: kubernetes.io/hostname

      # Topology spread for zone distribution
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: DoNotSchedule
        labelSelector:
          matchLabels:
            app: devops-go-app

      containers:
      - name: app
        image: registry.example.com/devops-go-app:v1.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP

        # Container security context
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: false
          runAsNonRoot: true
          runAsUser: 65534
          capabilities:
            drop:
              - ALL
            add:
              - NET_BIND_SERVICE

        # Environment variables from ConfigMap
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

        # Secrets
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
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: JWT_SECRET

        # Resource limits
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

        # Liveness probe
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 15
          periodSeconds: 15
          timeoutSeconds: 5
          failureThreshold: 3
          successThreshold: 1

        # Readiness probe
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
          successThreshold: 1

        # Startup probe for slow starting containers
        startupProbe:
          httpGet:
            path: /health/ready
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 12

      restartPolicy: Always
      terminationGracePeriodSeconds: 30
```

---

## 4. Fixed Kubernetes Service

**File**: `kubernetes/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: devops-go-service
  namespace: devops-demo
  labels:
    app: devops-go-app
    component: backend
spec:
  type: ClusterIP  # Changed from LoadBalancer (use Ingress instead)
  selector:
    app: devops-go-app
  ports:
  - name: http
    port: 80
    targetPort: 8080
    protocol: TCP
  sessionAffinity: None
```

---

## 5. Fixed Kubernetes Ingress

**File**: `kubernetes/ingress.yaml`

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: devops-go-ingress
  namespace: devops-demo
  annotations:
    # Basic routing
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"

    # Rate limiting
    nginx.ingress.kubernetes.io/limit-rps: "100"
    nginx.ingress.kubernetes.io/limit-connections: "50"
    nginx.ingress.kubernetes.io/limit-burst-multiplier: "5"

    # Backend configuration
    nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "60"

    # Security headers
    nginx.ingress.kubernetes.io/configuration-snippet: |
      more_set_headers "X-Frame-Options: DENY";
      more_set_headers "X-Content-Type-Options: nosniff";
      more_set_headers "X-XSS-Protection: 1; mode=block";
      more_set_headers "Strict-Transport-Security: max-age=31536000; includeSubDomains; preload";
      more_set_headers "Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'";
      more_set_headers "Referrer-Policy: strict-origin-when-cross-origin";
      more_set_headers "Permissions-Policy: geolocation=(), microphone=(), camera=()";

    # CORS (adjust as needed)
    nginx.ingress.kubernetes.io/enable-cors: "true"
    nginx.ingress.kubernetes.io/cors-allow-methods: "GET, POST, PUT, DELETE, OPTIONS"
    nginx.ingress.kubernetes.io/cors-allow-origin: "https://app.example.com"
    nginx.ingress.kubernetes.io/cors-allow-credentials: "true"
    nginx.ingress.kubernetes.io/cors-max-age: "3600"

spec:
  ingressClassName: nginx
  rules:
  - host: devops-go-app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: devops-go-service
            port:
              number: 80
  tls:
  - hosts:
    - devops-go-app.example.com
    secretName: devops-go-tls
```

---

## 6. Fixed HPA

**File**: `kubernetes/hpa.yaml`

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: devops-go-hpa
  namespace: devops-demo
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: devops-go-app
  minReplicas: 3  # Increased from 2 for better HA
  maxReplicas: 20  # Increased from 10 for traffic spikes
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
      - type: Pods
        value: 1
        periodSeconds: 60
      selectPolicy: Min
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 2
        periodSeconds: 30
      selectPolicy: Max
```

---

## 7. New: NetworkPolicy

**File**: `kubernetes/networkpolicy.yaml`

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
  # Allow traffic from Ingress controller
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080

  # Allow traffic from Prometheus
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring
      podSelector:
        matchLabels:
          app: prometheus
    ports:
    - protocol: TCP
      port: 8080

  egress:
  # Allow DNS
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: UDP
      port: 53

  # Allow PostgreSQL
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432

  # Allow Redis
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379

  # Allow external HTTPS (for APIs, etc.)
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
```

---

## 8. New: PodDisruptionBudget

**File**: `kubernetes/pdb.yaml`

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: devops-go-pdb
  namespace: devops-demo
spec:
  minAvailable: 2  # Always keep at least 2 pods running
  selector:
    matchLabels:
      app: devops-go-app
```

---

## 9. New: ResourceQuota

**File**: `kubernetes/resourcequota.yaml`

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
    pods: "50"
    services: "10"
    services.loadbalancers: "0"  # Prevent expensive LoadBalancers
    configmaps: "20"
    secrets: "20"
```

---

## 10. Updated Namespace with Pod Security

**File**: `kubernetes/namespace.yaml`

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: devops-demo
  labels:
    name: devops-demo
    environment: production
    # Pod Security Standards (PSS)
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Deployment Instructions

### 1. Docker Compose Deployment

```bash
# Create .env file
cd docker
cp .env.example .env
# Edit .env with your values

# Build and start
docker-compose up -d

# Verify
docker-compose ps
docker-compose logs -f app
```

### 2. Kubernetes Deployment

```bash
cd kubernetes

# Create namespace
kubectl apply -f namespace.yaml

# Create secrets (DO NOT commit to Git)
kubectl create secret generic app-secrets \
  --from-literal=POSTGRES_USER=devops_user \
  --from-literal=POSTGRES_PASSWORD='your-strong-password' \
  --from-literal=JWT_SECRET='your-jwt-secret' \
  -n devops-demo

# Apply configurations
kubectl apply -f configmap.yaml
kubectl apply -f resourcequota.yaml

# Deploy application
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f pdb.yaml
kubectl apply -f hpa.yaml
kubectl apply -f networkpolicy.yaml
kubectl apply -f ingress.yaml

# Verify deployment
kubectl get all -n devops-demo
kubectl get pods -n devops-demo -w
```

### 3. Validation

```bash
# Check security contexts
kubectl get pods -n devops-demo -o jsonpath='{.items[*].spec.securityContext}'

# Verify non-root user
kubectl exec -it -n devops-demo $(kubectl get pod -n devops-demo -l app=devops-go-app -o jsonpath='{.items[0].metadata.name}') -- id

# Check resource limits
kubectl describe pod -n devops-demo -l app=devops-go-app

# Test health checks
kubectl exec -it -n devops-demo $(kubectl get pod -n devops-demo -l app=devops-go-app -o jsonpath='{.items[0].metadata.name}') -- wget -O- http://localhost:8080/health/ready

# Monitor HPA
kubectl get hpa -n devops-demo -w
```

---

**These configurations are production-ready with all security issues addressed!**
