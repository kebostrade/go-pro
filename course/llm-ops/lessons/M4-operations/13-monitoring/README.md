# IO-13: Monitoring AI Systems

**Duration**: 3 hours
**Module**: 4 - LLM Operations & Observability

## Learning Objectives

- Implement metrics collection for LLM applications
- Set up distributed tracing for request flows
- Create dashboards for AI system monitoring
- Use LangSmith for LLM-specific observability
- Build alerting rules for AI system anomalies

## The Need for AI Observability

LLM applications have unique monitoring challenges:
- Non-deterministic outputs make debugging harder
- Latency varies based on token generation
- Cost tracking is critical for budget control
- Quality issues may not be immediately visible

## Metrics Collection

### Key Metrics for LLM Systems

| Metric | Description | Importance |
|--------|-------------|------------|
| Request latency | Time from request to response | Performance |
| Token usage | Input/output tokens per request | Cost tracking |
| Error rate | Failed requests / total requests | Reliability |
| Queue depth | Pending requests waiting | Capacity planning |
| Cache hit rate | Cached responses / total | Cost optimization |
| Quality score | LLM output quality metric | User experience |

### Prometheus Metrics Implementation

```python
from prometheus_client import Counter, Histogram, Gauge

# Request metrics
llm_requests_total = Counter(
    'llm_requests_total',
    'Total LLM requests',
    ['model', 'status']
)

llm_request_duration = Histogram(
    'llm_request_duration_seconds',
    'LLM request duration',
    ['model', 'operation'],
    buckets=[0.1, 0.5, 1.0, 2.0, 5.0, 10.0]
)

# Token metrics
llm_tokens_used = Counter(
    'llm_tokens_used_total',
    'Total tokens used',
    ['model', 'type']  # type: input/output
)

# Queue metrics
llm_queue_depth = Gauge(
    'llm_queue_depth',
    'Number of pending requests',
    ['model']
)

# Cache metrics
llm_cache_hits = Counter(
    'llm_cache_hits_total',
    'Cache hits',
    ['cache_type']
)
```

### Go Metrics with Prometheus

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    LLMRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "llm_requests_total",
            Help: "Total number of LLM requests",
        },
        []string{"model", "status"},
    )

    LLMRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "llm_request_duration_seconds",
            Help:    "LLM request duration in seconds",
            Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
        },
        []string{"model", "operation"},
    )

    LLMTokensUsed = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "llm_tokens_used_total",
            Help: "Total tokens used",
        },
        []string{"model", "type"},
    )
)
```

## Logging for AI Systems

### Structured Logging

```python
import structlog
import time

structlog.configure(
    processors=[
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer()
    ]
)

logger = structlog.get_logger()

def log_llm_request(request_id, model, prompt, response, duration):
    logger.info(
        "llm_request",
        request_id=request_id,
        model=model,
        input_tokens=response.usage.prompt_tokens,
        output_tokens=response.usage.completion_tokens,
        duration_ms=duration * 1000,
        output_length=len(response.choices[0].message.content),
    )
```

### Log Levels for LLM Operations

| Level | Use Case |
|-------|----------|
| DEBUG | Token counts, intermediate steps |
| INFO | Request/response summaries |
| WARNING | Rate limits, retries, degraded quality |
| ERROR | API failures, validation errors |

## Distributed Tracing

### OpenTelemetry Setup

```python
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.otlp import OTLPSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor

# Setup
provider = TracerProvider()
processor = BatchSpanProcessor(OTLPSpanExporter(endpoint="http://localhost:4317"))
provider.add_span_processor(processor)
trace.set_tracer_provider(provider)

tracer = trace.get_tracer(__name__)

# Trace LLM call
@tracer.start_as_current_span("llm_call")
def call_llm(prompt: str, model: str):
    span = trace.get_current_span()
    span.set_attribute("llm.model", model)
    span.set_attribute("llm.prompt_tokens", estimate_tokens(prompt))
    
    with tracer.start_as_current_span("api_request") as api_span:
        response = client.chat.completions.create(
            model=model,
            messages=[{"role": "user", "content": prompt}]
        )
        api_span.set_attribute("llm.response_tokens", response.usage.completion_tokens)
    
    return response
```

### Trace Span Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                      LLM Request Span                           │
├─────────────────────────────────────────────────────────────────┤
│  Attributes:                                                    │
│    - model: gpt-4                                               │
│    - prompt_tokens: 150                                         │
│    - temperature: 0.7                                            │
│                                                                 │
│  Child Spans:                                                  │
│    ├─ validate_input (5ms)                                     │
│    ├─ retrieve_context (45ms)                                   │
│    ├─ build_prompt (2ms)                                       │
│    ├─ api_request (1200ms) ◄── Main LLM call                  │
│    ├─ parse_response (3ms)                                      │
│    └─ log_and_cache (8ms)                                      │
│                                                                 │
│  Events:                                                        │
│    - cache_hit: false                                           │
│    - rate_limit_retry: 0                                        │
└─────────────────────────────────────────────────────────────────┘
```

## Grafana Dashboards

### Dashboard Components

```yaml
# dashboard.yml
panels:
  - title: "LLM Request Rate"
    type: graph
    targets:
      - expr: rate(llm_requests_total[5m])
        legend: "{{model}}"

  - title: "Token Usage by Model"
    type: graph
    targets:
      - expr: sum by (model) (rate(llm_tokens_used_total[1h]))

  - title: "Request Latency P99"
    type: graph
    targets:
      - expr: histogram_quantile(0.99, rate(llm_request_duration_bucket[5m]))

  - title: "Cost per Hour"
    type: stat
    targets:
      - expr: sum(rate(llm_tokens_used_total[1h])) * 0.03
```

### Key Dashboard Panels

1. **Request Volume**: Requests per minute by model
2. **Latency Distribution**: P50, P95, P99 response times
3. **Token Consumption**: Input/output tokens over time
4. **Error Rate**: Failed requests percentage
5. **Cache Performance**: Hit rate and savings
6. **Cost Tracking**: Spend per hour/day/month
7. **Queue Health**: Pending requests and wait times

## LangSmith Integration

### Setup and Configuration

```python
from langsmith import Client

client = Client()

# Trace runs
from langchain.callbacks import LangChainTracer

tracer = LangChainTracer(project_name="production-agent")

chain = prompt | llm | output_parser
chain.invoke(
    user_input,
    config={"callbacks": [tracer]}
)
```

### LangSmith Features

| Feature | Description |
|---------|-------------|
| Trace Visualization | See full request flow with all steps |
| Latency Analysis | Identify slow components |
| Token Tracking | Monitor usage per run |
| Dataset Evaluation | Test against ground truth |
| Feedback Collection | User thumbs up/down |

### Custom Evaluators

```python
from langsmith.evaluation import evaluate

def correctness_evaluator(run, example):
    """Evaluate if output matches expected."""
    predicted = run.outputs["output"]
    expected = example.outputs["expected"]
    return {
        "key": "correctness",
        "score": predicted.lower().strip() == expected.lower().strip()
    }

results = evaluate(
    my_agent.run,
    data="my-evaluation-dataset",
    evaluators=[correctness_evaluator],
)
```

## Custom Dashboards for AI

### Latency Breakdown

```sql
-- Query for latency percentiles
SELECT
    model,
    quantile(0.50) AS p50_latency,
    quantile(0.95) AS p95_latency,
    quantile(0.99) AS p99_latency,
    count(*) AS requests
FROM llm_requests
WHERE timestamp > now() - interval '1 hour'
GROUP BY model
```

### Cost Analysis

```sql
-- Daily cost by model
SELECT
    date(timestamp) AS day,
    model,
    SUM(input_tokens) * 0.0015 + SUM(output_tokens) * 0.002 AS cost
FROM llm_requests
GROUP BY date(timestamp), model
ORDER BY day DESC
```

### Quality Metrics

```python
# Track quality scores over time
quality_scores = Gauge(
    'llm_quality_score',
    'Quality score from user feedback',
    ['model', 'endpoint']
)

# Record feedback
def record_feedback(request_id, is_positive):
    quality_scores.labels(
        model=get_model(request_id),
        endpoint=get_endpoint(request_id)
    ).set(1 if is_positive else 0)
```

## Alerting Rules

### Critical Alerts

```yaml
groups:
  - name: llm-critical
    rules:
      - alert: LLMHighErrorRate
        expr: rate(llm_requests_total{status="error"}[5m]) > 0.05
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "LLM error rate above 5%"

      - alert: LLMHighLatency
        expr: histogram_quantile(0.99, rate(llm_request_duration_bucket[5m])) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "LLM P99 latency above 10s"

      - alert: LLMBudgetExceeded
        expr: llm_daily_cost > 1000
        for: 1m
        labels:
          severity: critical
```

### Warning Alerts

```yaml
      - alert: LLMCacheLowHitRate
        expr: rate(llm_cache_hits_total[10m]) / rate(llm_requests_total[10m]) < 0.3
        for: 10m
        labels:
          severity: warning

      - alert: LLMQueueBacklog
        expr: llm_queue_depth > 100
        for: 5m
        labels:
          severity: warning
```

## Observability in AI Agent Platform

The AI Agent Platform demonstrates observability patterns:

### Metrics Middleware

```go
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Track request
        rr := httptest.NewRecorder()
        next.ServeHTTP(rr, r)
        
        // Record metrics
        duration := time.Since(start)
        status := rr.Code
        
        metrics.LLMRequestsTotal.WithLabelValues(
            getModel(r),
            strconv.Itoa(status),
        ).Inc()
        
        metrics.LLMRequestDuration.WithLabelValues(
            getModel(r),
            "http",
        ).Observe(duration.Seconds())
    })
}
```

### Structured Logging

```go
func (s *AgentService) logRequest(ctx context.Context, req *AgentRequest, resp *AgentResponse) {
    logger := logging.FromContext(ctx)
    
    logger.Info("agent_request_completed",
        "request_id", req.ID,
        "model", req.Model,
        "input_tokens", resp.Usage.PromptTokens,
        "output_tokens", resp.Usage.CompletionTokens,
        "duration_ms", resp.Duration.Milliseconds(),
        "cache_hit", resp.Cached,
    )
}
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Prometheus** | Time-series database for metrics |
| **Grafana** | Visualization platform for metrics |
| **OpenTelemetry** | Standard for distributed tracing |
| **LangSmith** | LLM-specific observability platform |
| **P99 Latency** | 99th percentile response time |
| **Observability** | Ability to understand system state from outputs |

## Exercise

### Exercise 13.1: Prometheus Metrics

Add metrics to an LLM application:

1. Create a metrics module with request count, latency, and token usage
2. Expose `/metrics` endpoint
3. Add middleware to track all LLM requests

```python
# Your implementation
from prometheus_client import Counter, Histogram, Gauge

# Define metrics
llm_requests = Counter(...)
llm_latency = Histogram(...)
llm_tokens = Counter(...)

# Add to your LLM client
def call_llm(prompt):
    start = time.time()
    response = client.chat.completions.create(...)
    duration = time.time() - start
    
    # Record metrics
    llm_requests.labels(model="gpt-4").inc()
    llm_latency.observe(duration)
    llm_tokens.labels(type="input").inc(response.usage.prompt_tokens)
    
    return response
```

---

### Exercise 13.2: Grafana Dashboard Design

Design a Grafana dashboard for an LLM API:

| Panel | Type | Query |
|-------|------|-------|
| Request Rate | Graph | `rate(llm_requests_total[5m])` |
| Latency P99 | Graph | `histogram_quantile(0.99, rate(llm_request_duration_bucket[5m]))` |
| Error Rate | Stat | `rate(llm_requests_total{status="error"}[5m]) / rate(llm_requests_total[5m])` |
| Cost Today | Stat | `sum(llm_tokens_used_total{type="output"}[24h]) * 0.002` |

Create mockup:

```
┌─────────────────────────────────────────────────────────────────┐
│                      LLM API Dashboard                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────────┐  ┌──────────────────────────┐    │
│  │   Requests/min: 1,234    │  │   Error Rate: 0.5%      │    │
│  │   ████████████░░░        │  │   ██░░░░░░░░░░░░        │    │
│  └──────────────────────────┘  └──────────────────────────┘    │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Latency (P50/P95/P99)                                  │   │
│  │  ─────────────────────────────────                      │   │
│  │  P50: 450ms  ████████░░                                 │   │
│  │  P95: 1.2s   ██████████████░░                           │   │
│  │  P99: 3.5s   ██████████████████████░                    │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Token Usage (24h)                                       │   │
│  │  ─────────────────────────────────                      │   │
│  │  Input:  2.5M tokens ██████████████████████████████    │   │
│  │  Output: 1.8M tokens ████████████████████░░░░░░        │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

### Exercise 13.3: Tracing Implementation

Add OpenTelemetry tracing to an LLM pipeline:

1. Setup tracer provider with OTLP exporter
2. Create spans for each pipeline stage
3. Add span attributes for model, tokens, etc.

```python
# Your implementation
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider

# Setup (write your code)
def setup_tracing():
    # Create provider
    provider = TracerProvider()
    # Add OTLP exporter
    # Set as global provider
    pass

# Use in pipeline
def process_request(user_input):
    with tracer.start_as_current_span("process_request") as span:
        span.set_attribute("user.id", get_user_id())
        
        # Validation
        with tracer.start_as_current_span("validate") as span:
            validate_input(user_input)
        
        # Retrieval
        with tracer.start_as_current_span("retrieve") as span:
            context = retrieve_context(user_input)
            span.set_attribute("chunks.retrieved", len(context))
        
        # LLM call
        with tracer.start_as_current_span("llm_call") as span:
            response = call_llm(build_prompt(user_input, context))
            span.set_attribute("tokens.output", response.usage.completion_tokens)
        
        return response
```

---

### Exercise 13.4: Alert Configuration

Create Prometheus alerting rules:

| Alert | Condition | Severity | Action |
|-------|-----------|----------|--------|
| High Error Rate | error_rate > 5% | Critical | Page on-call |
| High Latency | p99 > 10s | Warning | Investigate |
| Low Cache Hit | hit_rate < 30% | Warning | Review caching |
| Budget Warning | daily_cost > 80% budget | Warning | Notify finance |

Write alert rules:

```yaml
# alerts.yml (write your rules)
- alert: 
  expr: 
  for: 
  labels:
    severity: 
  annotations:
    summary: 
```

---

### Exercise 13.5: Debug with Tracing

Given this trace, identify the bottleneck:

```
Span Name                    Duration    % of Total
─────────────────────────────────────────────────────
llm_request                  2500ms      100%
  ├─ validate_input         5ms         0.2%
  ├─ retrieve_context       1800ms      72%
  │   ├─ embedding_query    50ms
  │   ├─ vector_search     1700ms      ← BOTTLENECK
  │   └─ fetch_docs        50ms
  ├─ build_prompt           10ms        0.4%
  ├─ api_request            650ms       26%
  └─ parse_response         35ms        1.4%
```

**Answer**: Vector search takes 1700ms (68% of total). The issue is likely:
- Index not optimized
- Too many results being retrieved
- Network latency to vector database

---

### Exercise 13.6: LangSmith Setup

Configure LangSmith for an agent:

```python
# Your implementation
from langchain_openai import ChatOpenAI
from langchain.callbacks import LangChainTracer

# 1. Setup tracer
tracer = LangChainTracer(project_name="my-agent")

# 2. Create agent with tracing
llm = ChatOpenAI(model="gpt-4")
agent = create_agent(llm, tools, tracer)

# 3. Run and view in LangSmith
result = agent.invoke(
    {"input": "What is the weather?"},
    config={"callbacks": [tracer]}
)

# 4. Add feedback
tracer.end_feedback(
    run_id=result.run_id,
    feedback={"correctness": 1}  # or use user feedback
)
```

## Key Takeaways

- ✅ Prometheus + Grafana provides comprehensive metrics visibility
- ✅ OpenTelemetry enables distributed tracing across services
- ✅ LangSmith offers LLM-specific debugging and evaluation
- ✅ Alerting rules should cover error rates, latency, and costs
- ✅ Structured logging is essential for debugging LLM issues

## Module Summary

This module covered:
- Metrics collection with Prometheus
- Distributed tracing with OpenTelemetry
- Visualization with Grafana
- LLM-specific observability with LangSmith
- Alert configuration for AI systems

## Additional Resources

- [Prometheus Metrics Best Practices](https://prometheus.io/docs/practices/naming/)
- [OpenTelemetry Python SDK](https://opentelemetry.io/docs/python/)
- [LangSmith Documentation](https://docs.langchain.com/langsmith)
- [Grafana LLM Dashboard](https://grafana.com/blog/2024/02/15/llm-observability-dashboard/)
