# Phase 3: Distributed & Cloud — Research

**Researched:** 2026-04-01  
**Domain:** Kubernetes Operators, NATS JetStream, AWS Lambda  
**Confidence:** MEDIUM (training data, could not verify with Context7/tools)

---

## Summary

Phase 3 covers three distributed and cloud-native patterns for Go: Kubernetes operators (managing K8s resources programmatically), NATS JetStream (event-driven messaging with persistence), and AWS Lambda (serverless function deployment).

**Primary recommendations:**
1. Use `controller-runtime` for Kubernetes operators (higher-level than raw client-go)
2. Use NATS JetStream for all new event-driven work (persistence, delivery guarantees)
3. Use SAM CLI with `template.yaml` for Lambda deployment (AWS native, good Go support)

---

## Standard Stack

### Topic 9: Kubernetes

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| client-go | v0.32+ | Kubernetes API client | Official Kubernetes Go client |
| controller-runtime | v0.20+ | Operator framework | Official operator development framework |
| kubebuilder | v3.17+ | Operator scaffolding | Official K8s project scaffolding |
| k8s.io/api | v0.32+ | K8s typed definitions | Official Kubernetes types |
| k8s.io/apimachinery | v0.32+ | K8s API machinery | K8s core infrastructure |

**Installation:**
```bash
# Install kubebuilder (operator scaffolding)
go install sigs.k8s.io/controller-tools/cmd/kubebuilder@latest

# Install controller-runtime
go get sigs.k8s.io/controller-runtime@v0.20.0
```

### Topic 10: NATS

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| nats-io/nats.go | v1.37+ | NATS client | Official NATS client, already in course/ |
| nats-io/nats-server/v2 | latest | NATS server (dev) | Official NATS server |

**Installation:**
```bash
go get github.com/nats-io/nats.go@v1.37.0
```

### Topic 11: AWS Lambda

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| aws-lambda-go | v1.47+ | Lambda runtime client | Official AWS Lambda Go SDK |
| aws-sdk-go-v2 | v1.32+ | AWS SDK | Modern AWS SDK for Go |

**Installation:**
```bash
go get github.com/aws/aws-lambda-go@v1.47.0
go get github.com/aws/aws-sdk-go-v2@v1.32.0
```

---

## Architecture Patterns

### Recommended Project Structure

```
basic/projects/
├── kubernetes/           # Topic 9
│   ├── api/             # Custom Resource Definitions (CRDs)
│   ├── controllers/     # Operator controllers
│   ├── manifests/        # K8s manifests (Deployment, Service, etc.)
│   ├── helm/            # Helm chart
│   ├── Dockerfile
│   └── go.mod
├── nats-events/         # Topic 10
│   ├── publisher/       # NATS publisher
│   ├── subscriber/      # NATS subscriber
│   ├── jetstream/       # JetStream examples
│   └── go.mod
└── serverless/          # Topic 11
    ├── handlers/        # Lambda handlers
    ├── events/          # Event sources
    ├── template.yaml    # SAM template
    └── go.mod
```

---

## Kubernetes Patterns

### Pattern 1: Kubernetes Operator with controller-runtime

**What:** Custom controller that manages custom Kubernetes resources  
**When to use:** Need to manage application lifecycle, state, or custom resources in K8s  

**Code Structure:**
```go
// api/v1alpha1/zz_generated.deepcopy.go (generated)
// controllers/suite_test.go
// controllers/myresource_controller.go

package controllers

import (
    "context"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    myv1alpha1 "github.com/go-pro/kubernetes/api/v1alpha1"
)

type MyResourceReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Reconciliation logic
    return ctrl.Result{}, nil
}

func (r *MyResourceReconciler) SetupWithManager(mgr *control.Manager) error {
    return r.ControllerManagedBy(mgr).
        For(&myv1alpha1.MyResource{}).
        Complete(r)
}
```

**Source:** Training knowledge — kubebuilder/controller-runtime documentation

### Pattern 2: Kubernetes Manifests (Deployment, Service, HPA)

**What:** Standard K8s resources for containerized applications  
**When to use:** Any K8s deployment  

**Key Elements:**
```yaml
# Deployment with probes
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    spec:
      containers:
      - name: app
        ports:
        - name: http
          containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

### Pattern 3: Helm Chart Structure

**What:** Package manager for Kubernetes  
**When to use:** Production K8s deployments, templated configs  

```
helm/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── _helpers.tpl
└── .helmignore
```

### Anti-Patterns to Avoid

- **Don't use raw client-go for new operators:** Use controller-runtime/kubebuilder (higher-level, follows K8s conventions)
- **Don't skip probes:** Liveness/readiness probes are essential for production
- **Don't use :latest tag:** Always use specific version tags
- **Don't run as root:** Use securityContext.runAsNonRoot: true

---

## NATS Patterns

### Pattern 1: JetStream Publisher/Subscriber

**What:** Event streaming with persistence and delivery guarantees  
**When to use:** Need durable messages, at-least-once delivery, or stream replay  

**Code Structure:**
```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

// Connect with JetStream
nc, _ := nats.Connect(nats.DefaultURL)
js, _ := jetstream.New(nc)

// Publish with persistence
js.Publish(ctx, "events.orders", data)

// Subscribe with consumer
sub, _ := js.Subscribe(ctx, "events.orders", func(msg jetstream.Msg) {
    msg.Ack()  // Acknowledge
})
```

**Source:** Training knowledge — nats.io documentation, existing codebase (`advanced-topics/10-nats-eventdriven/`)

### Pattern 2: Queue Groups (Load Balancing)

**What:** Multiple subscribers form a queue group; each message goes to ONE subscriber  
**When to use:** Distributed worker pools, task processing  

**Code Structure:**
```go
// Worker 1, 2, 3 all in same queue group "task-workers"
// Each task processed by exactly ONE worker
nc.QueueSubscribe("tasks", "task-workers", func(m *nats.Msg) {
    processTask(string(m.Data))
    m.Ack()
})
```

### Pattern 3: Request/Reply

**What:** RPC-style communication with response  
**When to use:** Synchronous service calls  

```go
// Request
resp, err := nc.Request("echo", []byte("hello"), time.Second)
fmt.Println(string(resp.Data))

// Reply
nc.Subscribe("echo", func(m *nats.Msg) {
    m.Respond([]byte("pong"))
})
```

### Anti-Patterns to Avoid

- **Don't use core NATS for new event-driven work:** Use JetStream for persistence
- **Don't forget to Ack:** Messages won't be redelivered without acknowledgment
- **Don't use sync=True in production:** Can cause deadlocks

---

## AWS Lambda Patterns

### Pattern 1: Lambda Handler Signature

**What:** Go function that handles Lambda invocations  
**When to use:** Any Lambda function  

**Handler Types:**
```go
// 1. Simple handler
func handler() error { return nil }

// 2. With return value
func handler() (string, error) { return "hello", nil }

// 3. With context and input
func handler(ctx context.Context, event MyEvent) (MyResponse, error) {
    return MyResponse{Message: "hello"}, nil
}

// 4. API Gateway proxy (most common for web)
func handler(ctx context.Context, event events.APIGatewayProxyRequest) 
    (events.APIGatewayProxyResponse, error)
```

### Pattern 2: SAM Template (template.yaml)

**What:** AWS SAM template for serverless application definition  
**When to use:** Any SAM deployment  

**Structure:**
```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Runtime: go1.x
    MemorySize: 256
    Timeout: 30

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: ./handler
      Handler: handler
      Events:
        Api:
          Type: Api
          Properties:
            Path: /hello
            Method: get
```

### Pattern 3: Lambda URL (Simple HTTPS Endpoint)

**What:** Direct HTTPS endpoint for Lambda without API Gateway  
**When to use:** Simple APIs, lower cost, no API Gateway features needed  

**SAM Configuration:**
```yaml
MyFunction:
  Type: AWS::Serverless::Function
  Properties:
    Handler: handler
    Runtime: go1.x
    AutoPublishAlias: live
    FunctionUrlConfig:
      AuthType: NONE  # Public access
      InvokeMode: BUFFERED
```

### Pattern 4: Lambda with S3 Events

**What:** Lambda triggered by S3 object creation/changes  
**When to use:** Image processing, file processing pipelines  

```go
func handler(ctx context.Context, s3Event events.S3Event) {
    for _, record := range s3Event.Records {
        bucket := record.S3.Bucket.Name
        key := record.S3.Object.Key
        // Process file
    }
}
```

### Anti-Patterns to Avoid

- **Don't use go1.x runtime:** Use provided.al2 or provided.al2023 for better performance
- **Don't forget context timeout:** Lambda has 15min max; use context.WithTimeout
- **Don't allocate large objects in handler:** Reuse connections, use initialization
- **Don't use Lambda for long-running tasks:** Max 15 minutes; use ECS/Fargate for longer

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| K8s client | Raw REST calls | client-go | Handles auth, watch, retry, rate limiting |
| K8s operator | Custom reconciliation | controller-runtime/kubebuilder | CRD generation, reconciliation loops |
| NATS connection | Manual reconnect | nats.go connection options | Handles reconnect, heartbeat, backpressure |
| Lambda invocation | Raw HTTP | aws-lambda-go | Proper event parsing, context propagation |
| SAM deployment | Manual AWS calls | SAM CLI | Template validation, local invoke, deployment |

---

## Common Pitfalls

### Kubernetes
1. **Missing RBAC:** ServiceAccount/RoleBinding not configured → permission denied
2. **Wrong image pull policy:** Using :latest → image never updates
3. **No resource limits:** OOMKilled in production
4. **Probe failure:** Wrong path/port → CrashLoopBackoff

### NATS
1. **No consumer acknowledgment:** Messages stuck in pending
2. **Consumer group misconfiguration:** Messages not distributed correctly
3. **No connection pooling:** Connection overload under high load
4. **Ignoring reconnected handler:** Silent message loss on reconnect

### AWS Lambda
1. **Cold start latency:** Initialize heavy dependencies outside handler
2. **Missing error handling:** Errors not logged or retried properly
3. **No provisioned concurrency:** Cold starts for infrequently called functions
4. **VPC configuration:** Lambda in VPC has cold start penalty

---

## Open Questions

1. **Operator SDK vs Raw kubebuilder:** Which scaffolding approach for the learning template?
   - What we know: Both use controller-runtime under the hood
   - What's unclear: Which is simpler for learners?
   - Recommendation: kubebuilder (simpler, generates CRDs easily)

2. **Helm vs Kustomize for Kubernetes:** Which to teach?
   - What we know: Both are popular, Helm has more ecosystem
   - What's unclear: Learning curve comparison
   - Recommendation: Helm (more tutorials, wider adoption)

3. **Lambda URL vs API Gateway:** Which for the template?
   - What we know: Lambda URLs are simpler, API GW has more features
   - What's unclear: Cost difference for learning scenarios
   - Recommendation: Lambda URLs for basic, API Gateway for advanced features

---

## Environment Availability

**Note:** Could not verify tools due to MCP connection issues. Verify before planning.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| kubebuilder | K8s operators | TBD | TBD | controller-runtime directly |
| kubectl | K8s manifests | TBD | TBD | Docker Desktop K8s |
| helm | K8s charts | TBD | TBD | Kustomize |
| NATS server | NATS events | TBD (via Docker) | TBD | nats-streaming |
| AWS SAM CLI | Lambda deployment | TBD | TBD | Terraform |

---

## Existing Codebase Patterns

### NATS (already exists in advanced-topics/10-nats-eventdriven/)
- Publisher/subscriber pattern: ✅ Covered
- Queue groups: ✅ Covered
- Request/reply: ✅ Covered
- JetStream: ❌ Not explicitly covered (only core NATS)

### Kubernetes (already exists in advanced-topics/09-k8s-cloudnative/)
- Deployment manifest: ✅ Covered
- Service manifest: ✅ Covered
- HPA manifest: ✅ Covered (but syntax outdated)
- ConfigMap/Secret: ✅ Covered
- Operators: ❌ Not covered

### Lambda (already exists in advanced-topics/14-serverless-lambda/)
- Basic handler: ✅ Covered
- API Gateway: ✅ Covered
- DynamoDB: ✅ Covered
- S3 events: ✅ Covered
- SAM template: ✅ Covered
- Lambda URLs: ❌ Not covered

---

## Sources

### Primary (HIGH confidence — training data, needs verification)
- Kubernetes documentation: https://kubernetes.io/docs/
- kubebuilder book: https://book.kubebuilder.io/
- controller-runtime: https://pkg.go.dev/sigs.k8s.io/controller-runtime
- NATS documentation: https://docs.nats.io/
- NATS Go client: https://pkg.go.dev/github.com/nats-io/nats.go
- AWS Lambda Go: https://github.com/aws/aws-lambda-go
- AWS SAM documentation: https://docs.aws.amazon.com/serverless-application-model/

### Secondary (MEDIUM confidence — training data)
- Go.K8s operator patterns — general knowledge
- JetStream patterns — documented in nats-io/nats.go README
- SAM template.yaml patterns — AWS SAM GitHub examples

### Tertiary (LOW confidence — training data only)
- Specific library versions (v0.32+, v1.37+) — recommend verification via `npm view` or `go list`
- HPA v2 API — verify kubernetes.io/docs for current API version

---

## Metadata

**Confidence breakdown:**
- Kubernetes patterns: MEDIUM — general patterns correct, specific API versions need verification
- NATS patterns: HIGH — matches existing codebase and training data
- Lambda patterns: MEDIUM — general patterns correct, SAM spec details need verification

**Research date:** 2026-04-01  
**Valid until:** 2026-05-01 (30 days for stable tech)
**Tools status:** Could not verify with Context7 — recommend verification before finalizing plans
