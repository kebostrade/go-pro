# Building Cloud-Native Applications with Go and Kubernetes

Deploy and manage Go applications on Kubernetes.

## Learning Objectives

- Containerize Go applications
- Write Kubernetes manifests
- Implement health probes
- Configure ConfigMaps and Secrets
- Set up horizontal pod autoscaling
- Deploy with rolling updates

## Theory

### Production Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /server ./cmd/server

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
        version: v1
    spec:
      containers:
      - name: user-service
        image: myorg/user-service:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: log-level
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: config
          mountPath: /etc/config
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: user-service-config
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Health Check Implementation

```go
type HealthChecker struct {
    db     *sql.DB
    redis  *redis.Client
    checks map[string]func() error
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    h := &HealthChecker{
        db:     db,
        redis:  redis,
        checks: make(map[string]func() error),
    }
    h.checks["database"] = h.checkDB
    h.checks["redis"] = h.checkRedis
    return h
}

func (h *HealthChecker) Liveness(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "alive",
        "time":   time.Now().Format(time.RFC3339),
    })
}

func (h *HealthChecker) Readiness(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    status := make(map[string]string)
    allHealthy := true

    for name, check := range h.checks {
        if err := check(); err != nil {
            status[name] = fmt.Sprintf("unhealthy: %v", err)
            allHealthy = false
        } else {
            status[name] = "healthy"
        }
    }

    code := http.StatusOK
    if !allHealthy {
        code = http.StatusServiceUnavailable
    }

    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":   map[bool]string{true: "ready", false: "not ready"}[allHealthy],
        "checks":   status,
        "time":     time.Now().Format(time.RFC3339),
    })
}

func (h *HealthChecker) checkDB() error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    return h.db.PingContext(ctx)
}

func (h *HealthChecker) checkRedis() error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    return h.redis.Ping(ctx).Err()
}
```

### Graceful Shutdown

```go
func runServer() error {
    srv := &http.Server{Addr: ":8080"}

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGTERM)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    <-done
    log.Println("shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("shutdown error: %v", err)
    }

    log.Println("server stopped")
    return nil
}
```

### ConfigMap and Secrets

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
data:
  log-level: "info"
  app.yaml: |
    server:
      port: 8080
      timeout: 30s
    cache:
      ttl: 300s
---
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
type: Opaque
stringData:
  database-url: "postgres://user:pass@db:5432/users"
  jwt-secret: "super-secret-key"
```

## Security Considerations

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: user-service
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
  containers:
  - name: user-service
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
```

## Performance Tips

```yaml
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "1000m"
    memory: "512Mi"
```

## Exercises

1. Deploy a multi-service application
2. Configure autoscaling
3. Set up service mesh with Istio
4. Implement canary deployments

## Validation

```bash
cd exercises
kubectl apply -f k8s/
kubectl get pods -w
kubectl port-forward svc/user-service 8080:80
curl http://localhost:8080/health/ready
```

## Key Takeaways

- Always implement health probes
- Use resource limits
- Store secrets in Kubernetes Secrets
- Implement graceful shutdown
- Use ConfigMaps for configuration

## Next Steps

**[AT-10: Serverless Lambda](../AT-10-serverless-lambda/README.md)**

---

Kubernetes: production at scale. 🚢
