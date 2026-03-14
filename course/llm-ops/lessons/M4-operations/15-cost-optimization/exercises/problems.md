# Exercise: Cost Optimization

## Problem 1: Token Cost Calculation

Calculate costs for each scenario:

| Parameter | A | B | C |
|-----------|---|---|---|
| Model | GPT-4o | GPT-4o | GPT-3.5 Turbo |
| Input tokens | 500 | 1000 | 200 |
| Output tokens | 200 | 500 | 100 |
| Price/1M input | $2.50 | $2.50 | $0.50 |
| Price/1M output | $10.00 | $10.00 | $1.50 |
| **Cost per request** | | | |

If you have 10,000 requests/day for each scenario:

| Metric | A | B | C |
|--------|---|---|---|
| Daily cost | | | |
| Monthly cost | | | |

---

## Problem 2: Model Routing Implementation

Implement a smart model router:

```python
def route_query(query: str, user_tier: str = "free") -> str:
    """
    Route query to appropriate model based on complexity
    """
    # Classify query complexity
    # Simple: what is, define, list, find
    # Complex: analyze, compare, design, explain why
    
    # Route based on:
    # 1. Query complexity
    # 2. User tier (free vs premium)
    
    return "gpt-3.5-turbo"  # Default

# Test your router
test_queries = [
    "What is Python?",
    "Compare React and Vue for enterprise apps",
    "Find all files modified today",
    "Design a microservices architecture",
]

for q in test_queries:
    print(f"{q} → {route_query(q)}")
```

---

## Problem 3: Batch Processing Savings

Compare batch vs. individual processing:

| Scenario | Individual | Batch (10) |
|----------|------------|------------|
| Requests | 100 | 100 (10 batches) |
| Avg input tokens | 500 | 500 |
| Avg output tokens | 200 | 200 |
| Price/1M tokens | $5.00 | $2.50 |
| Cost/request | | |
| Total cost | | |

**Calculate savings:**

| Metric | Value |
|--------|-------|
| Individual total | |
| Batch total | |
| Savings | |
| Savings % | |

---

## Problem 4: Token Optimization

Optimize this prompt to reduce tokens:

**Original (280 tokens):**
```
Hello! I hope you're having a wonderful day today. I'm reaching out because 
I need some help understanding how neural networks work under the hood. 

Could you please provide a clear, comprehensive explanation that covers:
1. What exactly is a neural network?
2. How does it learn from data?
3. What are the key components like layers, neurons, weights, and biases?

I've been trying to learn about this on my own but I'm still confused about 
the mathematical foundations. Any help would be greatly appreciated!

Thank you so much for your time and assistance.
```

**Optimized version (target: <100 tokens):**

```
[Write your optimized version here]
```

---

## Problem 5: Cost Budget Alert Setup

Configure Prometheus alerts for budget monitoring:

```yaml
# alerts.yml
groups:
  - name: llm_cost
    rules:
      # Daily budget warning at 80%
      - alert: LLMDailyBudgetWarning
        expr: 
        for: 5m
        labels:
          severity: warning
        
      # Daily budget exceeded
      - alert: LLMDailyBudgetExceeded
        expr: 
        for: 1m
        labels:
          severity: critical
```

---

## Problem 6: Reserved Capacity Analysis

Calculate savings with reserved capacity:

| Parameter | Value |
|-----------|-------|
| Monthly token volume | 50,000,000 |
| On-demand price/1M | $5.00 |
| Reserved price/1M (40% off) | |
| On-demand monthly cost | |
| Reserved monthly cost | |
| Monthly savings | |
| Annual savings | |

**At what usage level does reserved make sense?**

| Usage Tier | Monthly Tokens | On-Demand | Reserved | Break-Even |
|------------|---------------|-----------|----------|------------|
| Light | 5M | | | |
| Medium | 25M | | | |
| Heavy | 100M | | | |

---

## Problem 7: ROI Analysis

Compare cost optimization approaches:

| Approach | Monthly Cost | Quality | Savings |
|----------|-------------|---------|---------|
| GPT-4o baseline | $10,000 | Excellent | - |
| Route to GPT-3.5 | $3,000 | Good | 70% |
| Route + cache | $2,000 | Good | 80% |
| Fine-tuned small model | $1,500 | Good | 85% |

| Approach | Implementation Effort | When to Use |
|----------|---------------------|-------------|
| Route to GPT-3.5 | Low | Simple queries dominant |
| Route + cache | Medium | High repeat queries |
| Fine-tuned small model | High | Specific domain |

---

## Problem 8: Cost Per User Calculation

Calculate cost per user for SaaS application:

| Metric | Value |
|--------|-------|
| Total monthly users | 10,000 |
| Active users (≥1 query/month) | 3,000 |
| Queries per active user/month | 50 |
| Avg tokens per query | 700 |
| Price/1K tokens | $0.005 |
| Infrastructure overhead | 20% |

| Metric | Calculation | Value |
|--------|-------------|-------|
| Total queries/month | | |
| Total tokens/month | | |
| LLM cost/month | | |
| Infrastructure cost | | |
| Total cost/month | | |
| Cost per active user | | |
| Cost per all users | | |

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.
