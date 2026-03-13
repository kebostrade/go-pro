# Deploying AI Agents

## Prerequisites

- `aictl` CLI installed and configured
- Access to Kubernetes cluster
- LLM API keys configured

## Agent Types

| Type | Description | Use Case |
|------|-------------|----------|
| `react` | ReAct (Reasoning + Acting) | Complex tasks requiring tool use |
| `base` | Simple LLM agent | Basic Q&A, summarization |
| `workflow` | Multi-step orchestration | Complex pipelines |
| `custom` | User-defined | Specialized use cases |

## Quick Deploy

### 1. Using CLI

```bash
# Create and deploy in one command
aictl agent create my-agent --type react --llm gpt-4
aictl agent deploy my-agent --env dev
```

### 2. Using Helm

```bash
# Add repo
helm repo add ai-platform https://charts.example.com

# Install
helm install my-agent ai-platform/ai-agent \
  --set agent.type=react \
  --set agent.llmModel=gpt-4 \
  --set replicaCount=2
```

### 3. Using Kustomize

```bash
# Create overlay
mkdir -p overlays/my-agent
cat > overlays/my-agent/kustomization.yaml << EOF
resources:
  - ../../base
namePrefix: my-agent-
configMapGenerator:
  - name: agent-config
    behavior: merge
    literals:
      - AGENT_TYPE=react
      - LLM_MODEL=gpt-4
EOF

# Apply
kustomize build overlays/my-agent | kubectl apply -f -
```

## Configuration

### Agent Configuration (values.yaml)

```yaml
# Agent settings
agent:
  type: react                    # Agent type
  llmProvider: openai            # LLM provider
  llmModel: gpt-4                # Model name
  maxTokens: 4096                # Max tokens per request
  temperature: 0.7               # Temperature

  # Tool configuration
  tools:
    enabled: true
    registry: builtin            # builtin, custom, remote
    custom:
      - name: my_tool
        url: https://tools.example.com/my_tool

  # Memory settings
  memory:
    enabled: true
    type: redis                  # redis, postgres, inmemory
    ttl: 3600                    # Memory TTL in seconds

  # Vector store (for RAG)
  vectorStore:
    enabled: true
    type: qdrant                 # qdrant, pgvector
    collection: agent_memory

# Resource limits
resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 250m
    memory: 512Mi

# Autoscaling
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `LLM_API_KEY` | LLM provider API key | Yes |
| `DATABASE_URL` | PostgreSQL connection | Yes |
| `REDIS_URL` | Redis connection | Yes |
| `QDRANT_URL` | Vector DB URL | Optional |
| `LOG_LEVEL` | Logging level | No |

## Secrets Management

### Using Kubernetes Secrets

```bash
# Create LLM credentials
kubectl create secret generic llm-credentials \
  --from-literal=api-key=sk-xxx

# Create DB credentials
kubectl create secret generic postgres-credentials \
  --from-literal=username=agent \
  --from-literal=password=secret
```

### Using External Secrets Operator

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: llm-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: ClusterSecretStore
  target:
    name: llm-credentials
  data:
    - secretKey: api-key
      remoteRef:
        key: ai-platform/llm
        property: api_key
```

## Multi-Environment Deployments

### Directory Structure

```
agents/
├── base/
│   └── agent.yaml
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml
    │   └── patches/
    ├── staging/
    │   ├── kustomization.yaml
    │   └── patches/
    └── prod/
        ├── kustomization.yaml
        └── patches/
```

### Deploy to Environment

```bash
# Dev
kustomize build agents/overlays/dev | kubectl apply -f -

# Staging
kustomize build agents/overlays/staging | kubectl apply -f -

# Production
kustomize build agents/overlays/prod | kubectl apply -f -
```

## Monitoring

### Check Agent Health

```bash
# Get agent status
aictl agent status my-agent

# View pods
kubectl get pods -l app=my-agent

# Check logs
aictl agent logs my-agent --follow
```

### Metrics

Agents expose Prometheus metrics at `/metrics`:

| Metric | Description |
|--------|-------------|
| `agent_requests_total` | Total requests processed |
| `agent_request_duration_seconds` | Request latency |
| `agent_tokens_used_total` | LLM tokens consumed |
| `agent_errors_total` | Total errors |
| `agent_tool_calls_total` | Tool invocations |

### Dashboards

Access Grafana dashboards at:
- **Overview**: https://grafana.example.com/d/ai-agents
- **Cost Analysis**: https://grafana.example.com/d/ai-costs
- **Performance**: https://grafana.example.com/d/ai-perf

## Troubleshooting

### Agent Not Starting

```bash
# Check events
kubectl describe deployment my-agent

# Check logs
kubectl logs -l app=my-agent --tail=100

# Verify secrets
kubectl get secrets -l app=my-agent
```

### High Latency

```bash
# Check resource usage
kubectl top pods -l app=my-agent

# Check HPA status
kubectl get hpa my-agent

# View metrics
aictl agent metrics my-agent
```

### LLM Errors

```bash
# Verify API key
kubectl get secret llm-credentials -o jsonpath='{.data.api-key}' | base64 -d

# Test connectivity
kubectl run test --rm -it --image=curlimages/curl -- \
  curl -H "Authorization: Bearer $API_KEY" \
  https://api.openai.com/v1/models
```

## Best Practices

1. **Resource Limits**: Always set appropriate limits
2. **Health Checks**: Configure liveness/readiness probes
3. **Secrets**: Use External Secrets Operator for production
4. **Monitoring**: Enable metrics and set up alerts
5. **Cost Control**: Set token limits and budget alerts
6. **Testing**: Test in dev before promoting to prod
