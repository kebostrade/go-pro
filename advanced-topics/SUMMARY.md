# Advanced Topics Summary

This document summarizes the Kubernetes and NATS advanced topics modules.

## 09 - Kubernetes & Cloud-Native Go

### Learning Objectives
- Deploy Go applications in Kubernetes
- Implement health checks and auto-scaling
- Manage configuration and secrets
- Use rolling updates for zero-downtime deployments

### File Structure
```
09-k8s-cloudnative/
├── README.md              # Complete Kubernetes guide
├── TROUBLESHOOTING.md     # Debugging guide
├── deploy.sh              # Automated deployment script
├── Dockerfile             # Multi-stage build
├── deployment.yaml        # Deployment with probes
├── service.yaml           # Service discovery
├── configmap.yaml         # Configuration management
├── secret.yaml            # Secret management
├── ingress.yaml           # External access
├── hpa.yaml               # Horizontal Pod Autoscaler
└── sample-app/
    ├── main.go            # Go HTTP server
    └── go.mod             # Module definition
```

### Key Features Implemented

1. **Containerization**
   - Multi-stage Docker build
   - Minimal runtime image (Alpine)
   - Non-root user execution
   - Health check integration

2. **Kubernetes Deployment**
   - 3 replica deployment
   - Rolling update strategy
   - Resource limits (CPU/Memory)
   - Security context (non-root)

3. **Health Checks**
   - Liveness probe: `/health/live`
   - Readiness probe: `/health/ready`
   - Startup probe for slow starting apps
   - Appropriate timeouts and thresholds

4. **Service Discovery**
   - ClusterIP service
   - Internal DNS resolution
   - Load balancing across pods

5. **Configuration Management**
   - ConfigMap for environment vars
   - Secret for sensitive data
   - Volume mounting support

6. **Auto-scaling**
   - HPA based on CPU (80%)
   - HPA based on memory (80%)
   - Min 2 replicas, max 10 replicas
   - Scale down stabilization (5 min)

7. **External Access**
   - Ingress for HTTP routing
   - Host-based routing
   - NGINX ingress controller compatible

### Quick Start Commands

```bash
# Deploy everything
./deploy.sh

# Manual deployment
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Check status
kubectl get pods
kubectl get svc

# Access app
kubectl port-forward svc/k8s-go-sample 8080:80
curl http://localhost:8080/

# Cleanup
kubectl delete all -l app=k8s-go-sample
```

## 10 - NATS & Event-Driven Go

### Learning Objectives
- Build event-driven microservices with NATS
- Implement pub/sub, queue groups, request/reply
- Use JetStream for persistence
- Handle connection management and errors

### File Structure
```
10-nats-eventdriven/
├── README.md              # Complete NATS guide
├── TESTING.md             # Testing and debugging guide
├── publisher/             # Basic publisher example
│   ├── main.go
│   └── go.mod
├── subscriber/            # Basic subscriber example
│   ├── main.go
│   └── go.mod
├── queue-group/           # Queue group examples
│   ├── publisher.go
│   ├── worker.go
│   └── go.mod
├── req-reply/             # Request-reply examples
│   ├── requester.go
│   ├── responder.go
│   └── go.mod
└── examples/              # All patterns in one file
    ├── nats_patterns.go
    └── go.mod
```

### Key Patterns Implemented

1. **Basic Pub/Sub**
   - Simple publish/subscribe
   - Multiple subscribers receive all messages
   - Use case: Broadcast notifications

2. **Queue Groups**
   - Work queue pattern
   - Load balancing across workers
   - Each message to ONE worker only
   - Use case: Task distribution

3. **Request/Reply**
   - Synchronous communication
   - Request with timeout
   - Response handling
   - Use case: Queries and RPC

4. **Wildcards**
   - `*` - Single level wildcard
   - `>` - Multi-level wildcard
   - Hierarchical subject routing
   - Use case: Flexible routing

5. **Connection Management**
   - Auto-reconnection
   - Connection status handlers
   - Error handling
   - Graceful shutdown

6. **Message Acknowledgment**
   - Explicit ack for JetStream
   - Reliable delivery
   - Consumer configuration

### Quick Start Commands

```bash
# Start NATS server
nats-server -js

# Test basic pub/sub (2 terminals)
cd subscriber && go run main.go
cd publisher && go run main.go

# Test queue groups (4 terminals)
cd queue-group
go run worker.go worker-1 &
go run worker.go worker-2 &
go run worker.go worker-3 &
go run publisher.go

# Test request-reply (2 terminals)
cd req-reply
go run responder.go
go run requester.go

# Test all patterns
cd examples
go run nats_patterns.go
```

## Common Patterns

### Kubernetes + NATS Integration

```go
// Deploy NATS in Kubernetes
kubectl apply -f nats-deployment.yaml

// Environment-based configuration
natsURL := getEnv("NATS_URL", "nats://nats:4222")
nc, _ := nats.Connect(natsURL)

// Kubernetes-style readiness probe
func isReady() bool {
    return nc.IsConnected() && nc.Status() == nats.CONNECTED
}
```

### Health Checks for NATS

```go
// Kubernetes liveness probe
func healthHandler(w http.ResponseWriter, r *http.Request) {
    if nc.IsConnected() {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Disconnected"))
    }
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    if nc.Status() == nats.CONNECTED {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Ready"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Not Ready"))
    }
}
```

## Deployment Scenarios

### Scenario 1: Microservices with NATS

```
[Frontend] -> [API Gateway (K8s)]
                |
                v
[Service A] ----> [NATS (K8s)] <---- [Service B]
     |                   |
     v                   v
[Database]        [Service C]
```

### Scenario 2: Event-Driven Processing

```
[Producer Pod] -> [NATS] -> [Consumer Queue Group]
                           -> [Worker 1] -> [Database]
                           -> [Worker 2] -> [Database]
                           -> [Worker 3] -> [Database]
```

### Scenario 3: Hybrid Cloud

```
[On-Prem K8s] --NATS Tunnel--> [Cloud K8s]
       |                            |
    [Services]                  [NATS Cluster]
       |                            |
       v                            v
[Local Database]            [Cloud Services]
```

## Best Practices

### Kubernetes
1. Always set resource limits
2. Use specific image versions (not `latest`)
3. Implement health checks
4. Run as non-root user
5. Use ConfigMaps/Secrets for config
6. Deploy with multiple replicas
7. Use rolling updates
8. Enable auto-scaling

### NATS
1. Handle connection errors gracefully
2. Use timeouts for requests
3. Implement reconnection handlers
4. Use queue groups for load balancing
5. Add logging for debugging
6. Test connection recovery
7. Use JetStream for persistence
8. Clean up subscriptions

## Monitoring & Observability

### Kubernetes Monitoring

```bash
# View logs
kubectl logs -f deployment/app

# Check resource usage
kubectl top pods
kubectl top nodes

# View events
kubectl get events --sort-by=.metadata.creationTimestamp
```

### NATS Monitoring

```bash
# NATS monitoring
nats-server -m 8222
curl http://localhost:8222/varz

# Application metrics
# Subscribe to metrics.> subjects
nc.Subscribe("metrics.>", func(m *nats.Msg) {
    processMetric(m)
})
```

## Troubleshooting

### Kubernetes
- **ImagePullBackOff**: Image doesn't exist or auth error
- **CrashLoopBackOff**: App crashing on startup
- **Pending**: Resource constraints or scheduling issues
- **Service unreachable**: Network policies or wrong port

### NATS
- **Connection refused**: NATS server not running
- **Timeout**: No responders or slow processing
- **No messages**: Wrong subject or no publisher
- **All workers receive**: Wrong queue group configuration

## Next Steps

After these modules:

1. **Kubernetes**
   - Learn Helm for packaging
   - Study service meshes (Istio, Linkerd)
   - Implement CI/CD (ArgoCD, Flux)
   - Add monitoring (Prometheus, Grafana)

2. **NATS**
   - Explore NATS JetStream in depth
   - Learn NATS KV and Object Store
   - Study NATS Service (mesh)
   - Build complex event-driven systems

3. **Combined**
   - Deploy NATS on Kubernetes
   - Build event-driven microservices
   - Implement service-to-service communication
   - Create distributed systems

## Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [NATS Documentation](https://docs.nats.io/)
- [Docker for Go](https://docs.docker.com/language/golang/)
- [NATS by Example](https://natsbyexample.com/)
