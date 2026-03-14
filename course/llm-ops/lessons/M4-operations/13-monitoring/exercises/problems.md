# Exercise: Monitoring AI Systems

## Problem 1: Prometheus Metrics Implementation

Implement a metrics module for tracking LLM requests:

```python
from prometheus_client import Counter, Histogram, Gauge, Info
import time

# Define all required metrics
class LLMMetrics:
    def __init__(self, service_name: str):
        # Create metrics with appropriate labels
        pass
    
    def record_request(self, model: str, status: str, duration: float):
        """Record a completed LLM request"""
        pass
    
    def record_tokens(self, model: str, token_type: str, count: int):
        """Record token usage"""
        pass

# Test your implementation
metrics = LLMMetrics("chat-api")
metrics.record_request("gpt-4", "success", 1.5)
metrics.record_tokens("gpt-4", "input", 150)
metrics.record_tokens("gpt-4", "output", 200)

# Verify metrics are exposed
# curl http://localhost:8000/metrics | grep llm_
```

---

## Problem 2: Grafana Dashboard Query

Write PromQL queries for these dashboard panels:

### Panel 1: Requests per second by model

```promql
# Write your query
```

### Panel 2: Average token cost per request (assume $0.002/1k output tokens)

```promql
# Write your query
```

### Panel 3: Cache hit rate percentage

```promql
# Write your query
```

### Panel 4: Error rate over last 5 minutes

```promql
# Write your query
```

---

## Problem 3: Trace Analysis

You have this trace data from a slow request:

```json
{
  "trace_id": "abc123",
  "spans": [
    {"name": "total", "duration_ms": 5000, "children": ["validate", "retrieve", "llm", "format"]},
    {"name": "validate", "duration_ms": 10},
    {"name": "retrieve", "duration_ms": 3000, "children": ["embed", "search", "fetch"]},
    {"name": "embed", "duration_ms": 50},
    {"name": "search", "duration_ms": 2800},
    {"name": "fetch", "duration_ms": 150},
    {"name": "llm", "duration_ms": 1800},
    {"name": "format", "duration_ms": 190}
  ]
}
```

### Analysis Questions

| Question | Answer |
|----------|--------|
| What is the bottleneck? | |
| What percentage of total time is vector search? | |
| What would you investigate first? | |

### Optimization Recommendations

| Component | Current | Target | Approach |
|-----------|---------|--------|----------|
| Vector search | 2800ms | | |
| LLM call | 1800ms | | |
| Total | 5000ms | <2000ms | |

---

## Problem 4: Alert Rule Configuration

Create Prometheus alerting rules:

```yaml
# Write your alerting rules

groups:
  - name: llm_alerts
    rules:
      # 1. High error rate - alert if > 5% errors in 5 minutes
      - alert: 
        expr: 
        for: 
        
      # 2. High latency - alert if P99 > 10 seconds
      - alert: 
        expr: 
        for: 
        
      # 3. Cache miss rate - alert if < 40% hit rate
      - alert: 
        expr: 
        for: 
        
      # 4. Budget warning - alert if daily cost > $500
      - alert: 
        expr: 
        for: 
```

---

## Problem 5: OpenTelemetry Setup

Implement OpenTelemetry tracing for an LLM pipeline:

```python
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import SimpleSpanProcessor, ConsoleSpanExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource

def init_tracing(service_name: str):
    """Initialize OpenTelemetry with console and OTLP exporters"""
    # Your implementation
    pass

def create_llm_span(tracer, prompt: str, response):
    """Create a span for LLM call with attributes"""
    # Your implementation
    # Include: model, input_tokens, output_tokens
    pass

# Test
tracer = init_tracing("my-llm-service")
with tracer.start_as_current_span("llm_call") as span:
    # Simulate LLM call
    span.set_attribute("llm.model", "gpt-4")
    span.set_attribute("llm.input_tokens", 100)
    span.set_attribute("llm.output_tokens", 250)
```

---

## Problem 6: Grafana Dashboard JSON

Create a minimal Grafana dashboard JSON:

```json
{
  "dashboard": {
    "title": "LLM API Monitoring",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(llm_requests_total[5m])",
            "legendFormat": "{{model}}"
          }
        ]
      }
      // Add more panels:
      // - Latency P99
      // - Error Rate  
      // - Token Usage
      // - Cost Tracking
    ]
  }
}
```

---

## Problem 7: LangSmith Integration

Set up LangSmith for debugging:

```python
from langchain_openai import ChatOpenAI
from langchain.schema import HumanMessage
from langchain.callbacks import LangChainTracer
from langsmith import Client

# 1. Initialize LangSmith tracer
tracer = LangChainTracer(
    project_name="production-debug",
    client=Client(api_key="your-api-key")
)

# 2. Create traced LLM
llm = ChatOpenAI(model="gpt-4", callbacks=[tracer])

# 3. Make request and observe in LangSmith dashboard
response = llm.invoke([
    HumanMessage(content="Explain quantum computing in simple terms")
])

# 4. How would you add user feedback to improve the trace?
def add_feedback(run_id: str, score: float, comment: str):
    # Your code
    pass
```

---

## Problem 8: Cost Analysis Dashboard

Design a cost analysis dashboard:

### Metrics to Track

| Metric | Query | Unit |
|--------|-------|------|
| Cost per hour | | $ |
| Cost per day | | $ |
| Cost per model | | $ |
| Tokens per request | | tokens |
| Cost per 1k tokens | | $ |

### Calculate Monthly Spend

Given:
- 100,000 requests/day
- Average 500 input + 300 output tokens/request
- Price: $0.0015/1k input, $0.002/1k output

| Component | Calculation | Cost |
|-----------|-------------|------|
| Input tokens/day | | |
| Output tokens/day | | |
| Input cost/day | | |
| Output cost/day | | |
| **Total/day** | | |
| **Total/month** | | |

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.
