# IO-15: Cost Optimization

**Duration**: 3 hours
**Module**: 4 - LLM Operations & Observability

## Learning Objectives

- Understand LLM pricing models and token economics
- Implement token optimization strategies
- Apply batch processing for cost reduction
- Use spot instances and reserved capacity
- Build cost monitoring and budgeting systems

## LLM Cost Fundamentals

### Token-Based Pricing

LLMs are priced by token (typically ~4 characters = 1 token):

| Model | Input (/1M tokens) | Output (/1M tokens) |
|-------|------------------|-------------------|
| GPT-4o | $2.50 | $10.00 |
| GPT-4 Turbo | $10.00 | $30.00 |
| GPT-3.5 Turbo | $0.50 | $1.50 |
| Claude 3.5 Sonnet | $3.00 | $15.00 |
| Claude 3 Haiku | $0.25 | $1.25 |

### Cost Calculation

```
Total Cost = (Input Tokens × Input Price) + (Output Tokens × Output Price)
```

### Token Estimation

```python
def estimate_tokens(text: str) -> int:
    """Rough estimation: ~4 characters per token"""
    return len(text) // 4

def estimate_cost(prompt: str, completion: str, model: str = "gpt-4o") -> float:
    prices = {
        "gpt-4o": {"input": 0.0025, "output": 0.01},
        "gpt-3.5-turbo": {"input": 0.0005, "output": 0.0015},
    }
    
    input_tokens = estimate_tokens(prompt)
    output_tokens = estimate_tokens(completion)
    
    return (
        input_tokens * prices[model]["input"] / 1_000_000 +
        output_tokens * prices[model]["output"] / 1_000_000
    )
```

## Token Optimization Strategies

### Prompt Compression

```python
def compress_prompt(prompt: str) -> str:
    """Remove unnecessary tokens while preserving meaning"""
    # Remove filler words
    fillers = ["please", "could you", "would you mind", "kindly"]
    result = prompt.lower()
    for filler in fillers:
        result = result.replace(filler, "")
    
    # Remove extra whitespace
    result = " ".join(result.split())
    
    return result.strip()
```

### Context Truncation

```python
def truncate_context(context: str, max_tokens: int = 4000) -> str:
    """Truncate context to fit within token limit"""
    max_chars = max_tokens * 4  # Approximate
    
    if len(context) <= max_chars:
        return context
    
    # Truncate and add indicator
    return context[:max_chars] + "\n\n[... content truncated ...]"
```

### Smart Chunking

```python
def smart_chunk(text: str, chunk_size: int = 1000) -> list[str]:
    """Split text into chunks, respecting boundaries"""
    chunks = []
    
    # Try to split on paragraphs first
    paragraphs = text.split("\n\n")
    current_chunk = ""
    
    for para in paragraphs:
        if len(current_chunk) + len(para) <= chunk_size:
            current_chunk += para + "\n\n"
        else:
            if current_chunk:
                chunks.append(current_chunk.strip())
            current_chunk = para + "\n\n"
    
    if current_chunk:
        chunks.append(current_chunk.strip())
    
    return chunks
```

## Model Selection Strategy

### Routing Based on Query Complexity

```python
def select_model(query: str, user_tier: str = "free") -> str:
    """Route query to appropriate model based on complexity"""
    
    # Simple classification query → cheap model
    simple_patterns = ["what is", "define", "list", "find"]
    
    # Complex reasoning → premium model
    complex_patterns = ["analyze", "compare", "explain why", "design"]
    
    query_lower = query.lower()
    
    if any(p in query_lower for p in simple_patterns):
        return "gpt-3.5-turbo"  # $0.001/1k tokens
    elif any(p in query_lower for p in complex_patterns):
        return "gpt-4o"  # $0.0125/1k tokens
    else:
        return "gpt-4o-mini"  # Mid-tier option
```

### Model Selection Matrix

| Query Type | Recommended Model | Cost/1k Tokens |
|------------|-------------------|----------------|
| Simple Q&A | GPT-3.5 Turbo | $0.002 |
| Code generation | GPT-4o | $0.015 |
| Complex reasoning | GPT-4o | $0.015 |
| High volume, low latency | GPT-4o-mini | $0.003 |
| Large context | Claude 3.5 Sonnet | $0.018 |

## Batch Processing

### Queue-Based Batching

```python
import asyncio
from collections import defaultdict
import time

class BatchProcessor:
    def __init__(self, batch_size: int = 10, max_wait: float = 1.0):
        self.batch_size = batch_size
        self.max_wait = max_wait
        self.queue = asyncio.Queue()
        self.processing = False
    
    async def add(self, item):
        await self.queue.put(item)
        
        if self.queue.qsize() >= self.batch_size:
            await self.process_batch()
    
    async def process_batch(self):
        if self.processing:
            return
        
        self.processing = True
        batch = []
        
        # Collect batch
        while len(batch) < self.batch_size and not self.queue.empty():
            batch.append(await asyncio.wait_for(
                self.queue.get(), 
                timeout=self.max_wait
            ))
        
        if batch:
            # Process batch with LLM
            results = await self.llm_batch_process(batch)
            
            # Return results to callers
            for future, result in zip(batch, results):
                future.set_result(result)
        
        self.processing = False
```

### Batch API Usage

```python
# OpenAI Batch API
from openai import OpenAI

client = OpenAI()

# Create batch job
batch_job = client.batches.create(
    input_file_id=file.id,
    endpoint="/v1/chat/completions",
    completion_window="24h",
)

# Check status
status = client.batches.retrieve(batch_job.id)
```

## Infrastructure Cost Optimization

### Spot Instances

```yaml
# Terraform AWS spot instance
resource "aws_instance" "llm_inference" {
  instance_type = "g5.xlarge"
  
  # Spot configuration
  instance_market_options {
    market_type = "spot"
    spot_options {
      instance_interruption_behavior = "terminate"
      max_price = "0.50"  # Max hourly price
    }
  }
  
  tags = {
    Name = "llm-inference-spot"
  }
}
```

### Reserved Capacity

```python
class ReservedCapacityManager:
    def __init__(self, provider: str):
        self.provider = provider
        self.reserved_tokens = 0
        self.used_tokens = 0
    
    def purchase_reserved(self, tokens_per_month: int, discount: float = 0.5):
        """Purchase reserved capacity"""
        self.reserved_tokens = tokens_per_month
        
        # Calculate savings
        reserved_cost = tokens_per_month * 0.01 * (1 - discount)
        on_demand_cost = tokens_per_month * 0.01
        
        return {
            "reserved": reserved_cost,
            "on_demand": on_demand_cost,
            "savings": on_demand_cost - reserved_cost,
            "savings_percent": discount * 100
        }
    
    def use_tokens(self, tokens: int):
        """Track token usage against reserved"""
        self.used_tokens += tokens
        
        if self.used_tokens > self.reserved_tokens:
            return "overage"  # Pay premium for overage
        return "reserved"
```

## Cost Budgeting and Alerts

### Budget Configuration

```python
from datetime import datetime, timedelta

class CostBudget:
    def __init__(self, daily_limit: float, monthly_limit: float):
        self.daily_limit = daily_limit
        self.monthly_limit = monthly_limit
        self.daily_spend = 0
        self.monthly_spend = 0
        self.last_reset = datetime.now()
    
    def track_cost(self, tokens: int, model: str):
        """Track cost and check budgets"""
        cost = self.calculate_cost(tokens, model)
        
        self.daily_spend += cost
        self.monthly_spend += cost
        
        # Check thresholds
        alerts = []
        
        if self.daily_spend > self.daily_limit * 0.8:
            alerts.append(f"WARNING: 80% of daily budget used: ${self.daily_spend:.2f}")
        
        if self.daily_spend > self.daily_limit:
            alerts.append(f"CRITICAL: Daily budget exceeded: ${self.daily_spend:.2f}")
        
        if self.monthly_spend > self.monthly_limit * 0.9:
            alerts.append(f"WARNING: 90% of monthly budget used: ${self.monthly_spend:.2f}")
        
        return cost, alerts
    
    def calculate_cost(self, tokens: int, model: str) -> float:
        prices = {
            "gpt-4o": 0.0125,
            "gpt-4o-mini": 0.003,
            "gpt-3.5-turbo": 0.002,
        }
        return tokens / 1_000_000 * prices.get(model, 0.01)
```

### Alert Integration

```python
def send_alert(message: str, severity: str):
    """Send budget alert"""
    if severity == "critical":
        # Page on-call
        send_pagerduty(message)
    else:
        # Send to Slack
        send_slack(message)
    
    # Log for tracking
    log_alert(message, severity)
```

## Cost Monitoring Dashboard

### Key Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Cost per request | Average cost per LLM call | < $0.01 |
| Token efficiency | Output tokens / Input tokens | > 0.5 |
| Cache savings | % requests from cache | > 30% |
| Model distribution | % requests by model tier | Tiered |
| Cost vs. value | Cost per successful request | Mapped to ROI |

### Dashboard Queries

```promql
# Daily spend
sum(llm_tokens_used_total[24h]) * 0.0125 / 1e6

# Cost by model
sum by (model) (llm_tokens_used_total[24h] * 0.0125) / 1e6

# Requests by model
sum by (model) (llm_requests_total[24h])

# Cost per request
sum(llm_tokens_used_total[1h]) / sum(llm_requests_total[1h]) * 0.0125 / 1000
```

## Optimization Techniques Summary

### Quick Wins

| Technique | Savings | Implementation |
|-----------|---------|----------------|
| Use GPT-4o-mini for simple tasks | 60-80% | Model routing |
| Enable KV cache | 30-50% | Provider config |
| Semantic caching | 30-40% | Redis layer |
| Prompt compression | 10-20% | Pre-processing |

### Advanced Optimization

| Technique | Savings | Complexity |
|-----------|---------|------------|
| Batch API | 50% | Medium |
| Reserved capacity | 40-60% | Low |
| Spot instances | 60-90% | High |
| Fine-tuned small models | 70-90% | High |

## Cost Optimization in AI Agent Platform

The platform demonstrates cost optimization:

### Token Budgeting

```go
type TokenBudget struct {
    MaxTokens     int
    ReservedFor   int  // reserved for response
    UsedForPrompt int
}

func (b *TokenBudget) CanProcess(promptTokens int) bool {
    return promptTokens <= b.MaxTokens - b.ReservedFor
}
```

### Cost Tracking

```go
type CostTracker struct {
    DailySpend    float64
    MonthlySpend  float64
    DailyLimit    float64
    MonthlyLimit  float64
}

func (c *CostTracker) Record(tokens int, model string) error {
    cost := CalculateCost(tokens, model)
    c.DailySpend += cost
    c.MonthlySpend += cost
    
    if c.DailySpend > c.DailyLimit {
        return ErrBudgetExceeded
    }
    return nil
}
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Token** | Basic unit of LLM text (~4 chars) |
| **Input tokens** | Tokens in the prompt |
| **Output tokens** | Tokens in the response |
| **Batch processing** | Grouping requests for efficiency |
| **Spot instances** | Excess cloud capacity at discount |
| **Reserved capacity** | Pre-paid compute at discount |

## Exercise

### Exercise 15.1: Cost Analysis

Calculate monthly costs for different scenarios:

| Scenario | Requests/Day | Avg Input | Avg Output | Model | Daily Cost | Monthly Cost |
|----------|--------------|-----------|------------|-------|------------|--------------|
| A | 1,000 | 500 | 200 | GPT-4o | | |
| B | 10,000 | 200 | 100 | GPT-3.5 Turbo | | |
| C | 5,000 | 1000 | 500 | GPT-4o-mini | | |

---

### Exercise 15.2: Model Routing Implementation

Implement a model router:

```python
def route_model(query: str, user_tier: str = "free") -> str:
    """Route query to appropriate model"""
    # Classify query complexity
    # Consider user tier
    # Return model name
    
    # Simple queries → cheap model
    # Complex queries → premium model
    pass

# Test cases
print(route_model("What is Python?"))  # Should be gpt-3.5-turbo
print(route_model("Analyze the trade implications of..."))  # Should be gpt-4o
```

---

### Exercise 15.3: Batch Processing

Design a batch processing system:

```
Request 1: "What is AI?"         ┐
Request 2: "What is ML?"         │
Request 3: "What is DL?"         ├─→ Batch → Single API Call
...                              │
Request 10: "What is NLP?"      ┘
```

| Aspect | Value |
|--------|-------|
| Batch size | |
| Max wait time | |
| Expected savings | |
| Trade-off (latency vs cost) | |

---

### Exercise 15.4: Budget Alert Configuration

Set up budget alerts:

| Alert | Threshold | Action |
|-------|-----------|--------|
| Daily warning | 80% of daily | Slack notification |
| Daily critical | 100% of daily | Page on-call |
| Monthly warning | 80% of monthly | Email finance |
| Monthly critical | 100% of monthly | Disable service |

---

### Exercise 15.5: Token Optimization

Optimize this prompt:

**Original (250 tokens):**
```
Hello! I hope you're doing great today. I was wondering if you could help me 
with something. You see, I'm trying to understand how machine learning works 
in the context of modern software development. Could you please provide a 
clear and concise explanation? That would be really helpful. Thank you so much!
```

**Optimized:**

```
[Write optimized version - target: 100 tokens or less while preserving meaning]
```

---

### Exercise 15.6: Cost-Benefit Analysis

Compare costs and benefits:

| Approach | Cost/Month | Benefit | ROI |
|----------|------------|---------|-----|
| GPT-4o for all | $3,000 | High quality | Baseline |
| GPT-4o-mini for simple | $1,200 | Good quality | ? |
| Hybrid (route + cache) | $800 | Good quality | ? |

| Approach | Savings vs Baseline | Quality Impact |
|----------|---------------------|----------------|
| GPT-4o-mini for simple | | Minimal |
| Hybrid | | None |

---

### Exercise 15.7: Reserved Capacity Calculator

Calculate reserved capacity savings:

| Parameter | Value |
|-----------|-------|
| Expected monthly tokens | 10,000,000 |
| On-demand rate | $10/1M tokens |
| Reserved rate (50% off) | |
| Monthly on-demand cost | |
| Monthly reserved cost | |
| Annual savings | |

## Key Takeaways

- ✅ Token optimization can reduce costs by 20-40%
- ✅ Model routing matches query complexity to model cost
- ✅ Batch processing provides ~50% savings on large volumes
- ✅ Spot instances and reserved capacity offer 40-90% infrastructure savings
- ✅ Budget alerts prevent cost overruns

## Module Summary

This module covered:
- Token economics and pricing models
- Token optimization strategies
- Model selection and routing
- Batch processing techniques
- Infrastructure cost optimization
- Budget monitoring and alerting

## Additional Resources

- [OpenAI Pricing](https://openai.com/pricing)
- [Anthropic Pricing](https://www.anthropic.com/pricing)
- [AWS Spot Instances](https://aws.amazon.com/ec2/spot/)
- [vLLM Optimization](https://vllm.ai/)
