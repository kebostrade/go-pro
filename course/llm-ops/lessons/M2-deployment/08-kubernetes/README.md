# IO-08: Kubernetes for AI

**Duration**: 3 hours
**Module**: 2 - Deployment & Serving

## Learning Objectives

- Deploy LLM services on Kubernetes
- Configure GPU scheduling and operators
- Implement auto-scaling for LLM workloads
- Use Helm charts for LLM deployment

## Kubernetes for LLM Serving

Kubernetes provides:

- **Orchestration**: Manage multi-container applications
- **Scaling**: Horizontal and vertical scaling
- **GPU Support**: NVIDIA GPU operator
- **High Availability**: Self-healing, rolling updates

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                            │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    Control Plane                             ││
│  │  API Server | Scheduler | Controller Manager | etcd        ││
│  └─────────────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Node 1    │  │   Node 2    │  │   Node 3    │              │
│  │  ┌───────┐  │  │  ┌───────┐  │  │  ┌───────┐  │              │
│  │  │ GPU 1 │  │  │  │ GPU 2 │  │  │  │ GPU 3 │  │              │
│  │  └───────┘  │  │  └───────┘  │  │  └───────┘  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

## 1. GPU Scheduling in Kubernetes

### Installing NVIDIA GPU Operator

```bash
# Add NVIDIA Helm repository
helm repo add nvidia https://nvidia.github.io/gpu-operator
helm repo update

# Install GPU Operator
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator \
  --create-namespace
```

### Verify GPU Access

```bash
# Check GPU nodes
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'

# Run a test pod with GPU
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test
spec:
  containers:
  - name: cuda
    image: nvidia/cuda:12.1.0-runtime-ubuntu22.04
    command: ["nvidia-smi"]
    resources:
      limits:
        nvidia.com/gpu: 1
EOF
```

## 2. Deploying vLLM on Kubernetes

### Deployment Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vllm-deployment
  namespace: llm
spec:
  replicas: 2
  selector:
    matchLabels:
      app: vllm
  template:
    metadata:
      labels:
        app: vllm
    spec:
      containers:
      - name: vllm
        image: vllm/vllm-openai:latest
        ports:
        - containerPort: 8000
        env:
        - name: MODEL
          value: "meta-llama/Llama-2-7b-hf"
        - name: GPU_MEMORY_UTILIZATION
          value: "0.9"
        - name: TENSOR_PARALLEL_SIZE
          value: "1"
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: "16Gi"
          requests:
            memory: "8Gi"
        volumeMounts:
        - name: model-cache
          mountPath: /root/.cache/huggingface
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 60
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
      volumes:
      - name: model-cache
        persistentVolumeClaim:
          claimName: model-cache-pvc
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: vllm-service
  namespace: llm
spec:
  selector:
    app: vllm
  ports:
  - port: 80
    targetPort: 8000
  type: ClusterIP
```

## 3. API Gateway on Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: llm
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: api-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: BACKEND_URL
          value: "http://vllm-service.llm.svc.cluster.local"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        resources:
          limits:
            memory: "512Mi"
            cpu: "500m"
          requests:
            memory: "256Mi"
            cpu: "250m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 4. Auto-Scaling

### Horizontal Pod Autoscaler (HPA)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: vllm-hpa
  namespace: llm
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: vllm-deployment
  minReplicas: 1
  maxReplicas: 10
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
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
      - type: Pods
        value: 4
        periodSeconds: 15
      selectPolicy: Max
```

### Pod Disruption Budget

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: vllm-pdb
  namespace: llm
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: vllm
```

## 5. Helm Charts for LLM Deployment

### Chart Structure

```
llm-deployment/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml
│   └── pdb.yaml
└── .helmignore
```

### values.yaml

```yaml
replicaCount: 2

image:
  repository: vllm/vllm-openai
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80
  targetPort: 8000

resources:
  limits:
    nvidia.com/gpu: 1
    memory: "16Gi"
  requests:
    memory: "8Gi"

env:
  MODEL: "meta-llama/Llama-2-7b-hf"
  GPU_MEMORY_UTILIZATION: "0.9"

autoscaling:
  enabled: true
  minReplicas: 1
  maxReplicas: 10

persistence:
  enabled: true
  size: 50Gi
```

### deployment.yaml template

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "llm.fullname" . }}
  labels:
    {{- include "llm.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "llm.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "llm.selectorLabels" . | nindent 8 }}
    spec:
      containers:
      - name: {{ .Chart.Name }}
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        ports:
        - name: http
          containerPort: 8000
          protocol: TCP
        {{- with .Values.env }}
        env:
          {{- range $key, $value := . }}
          - name: {{ $key }}
            value: {{ $value | quote }}
          {{- end }}
        {{- end }}
        {{- with .Values.resources }}
        resources:
          {{- toYaml . | nindent 10 }}
        {{- end }}
        {{- include "llm.probes" . | nindent 8 }}
```

## 6. Go Application Deployment

The [API Gateway in this repository](services/api-gateway/) can be deployed on Kubernetes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: api-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: BACKEND_URL
          valueFrom:
            configMapKeyRef:
              name: api-config
              key: backend-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        - name: RATE_LIMIT
          valueFrom:
            configMapKeyRef:
              name: api-config
              key: rate-limit
```

## 7. Production Considerations

### Resource Quotas

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: llm-quota
  namespace: llm
spec:
  hard:
    requests.cpu: "16"
    requests.memory: 64Gi
    limits.cpu: "32"
    limits.memory: 128Gi
    requests.nvidia.com/gpu: "4"
    pods: "20"
```

### Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: vllm-network-policy
spec:
  podSelector:
    matchLabels:
      app: vllm
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: api-gateway
    ports:
    - protocol: TCP
      port: 8000
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
```

## Key Takeaways

- ✅ NVIDIA GPU Operator enables GPU scheduling in Kubernetes
- ✅ HPA provides automatic scaling for LLM workloads
- ✅ Helm charts simplify deployment management
- ✅ Resource quotas prevent resource starvation

## Next Steps

Complete Module 2 and proceed to Module 3: Monitoring & Observability

## Additional Resources

- [Kubernetes GPU Documentation](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/)
- [NVIDIA GPU Operator](https://docs.nvidia.com/datacenter/driver-setup/)
- [Helm Documentation](https://helm.sh/docs/)
- [KEDA Scaling](https://keda.sh/) for event-driven scaling
