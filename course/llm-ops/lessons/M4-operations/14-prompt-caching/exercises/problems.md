# Exercise: Prompt Caching Strategies

## Problem 1: Semantic Cache Implementation

Implement a semantic cache with similarity threshold:

```python
import redis
import numpy as np
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

class SemanticCache:
    def __init__(self, redis_url: str, threshold: float = 0.85):
        self.redis = redis.from_url(redis_url)
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        self.threshold = threshold
    
    def _embed(self, text: str) -> np.ndarray:
        return self.model.encode(text)
    
    def get(self, prompt: str):
        # 1. Generate embedding for prompt
        # 2. Iterate through cache keys
        # 3. Calculate cosine similarity
        # 4. Return if above threshold
        pass
    
    def set(self, prompt: str, response: str, ttl: int = 3600):
        # 1. Generate embedding
        # 2. Store prompt, response, embedding in Redis hash
        # 3. Set TTL
        pass

# Test it
cache = SemanticCache("redis://localhost:6379")
cache.set("How to reverse a string?", "Use string[::-1] in Python")
result = cache.get("How do I reverse a string in Python?")
print(result)  # Should return cached response
```

---

## Problem 2: Cache Hit Rate Calculation

Given the following data, calculate cache metrics:

| Day | Requests | Hits | Misses |
|-----|----------|------|--------|
| Monday | 5,000 | 1,500 | 3,500 |
| Tuesday | 6,000 | 2,400 | 3,600 |
| Wednesday | 7,000 | 3,500 | 3,500 |
| Thursday | 8,000 | 4,000 | 4,000 |
| Friday | 9,000 | 4,500 | 4,500 |

| Metric | Formula | Value |
|--------|---------|-------|
| Total Hit Rate | | |
| Total Requests | | |
| Total Hits | | |
| Total Misses | | |

---

## Problem 3: Cost Savings Analysis

Compare costs with and without caching:

Parameters:
- 100,000 API calls/day
- Average tokens per request: 1,000 input, 500 output
- Pricing: $0.0015/1k input tokens, $0.002/1k output tokens
- Cache hit rate: 40%
- Cache lookup cost: $0.0001/call

| Metric | Without Cache | With Cache |
|--------|--------------|------------|
| LLM API Cost | | |
| Cache Infrastructure | $0 | |
| Total Cost | | |
| Savings | | |
| Savings % | | |

---

## Problem 4: Prompt Optimization

Rewrite prompts for better cache efficiency:

### Before/After Comparison

| Original Prompt | Problem | Optimized Prompt |
|-----------------|---------|------------------|
| "Hello there! I hope you're having a wonderful day. I was wondering if you could help me understand how machine learning works. Thanks so much!" | | |
| "Acting as a Python expert with 20 years of experience, the best coding practices, deep knowledge of algorithms and data structures, please write a function to calculate factorial." | | |
| "Considering the context of our previous conversation about web development, and taking into account what we discussed yesterday about APIs, could you explain REST?" | | |

---

## Problem 5: Cache Invalidation Rules

Design invalidation strategies:

| Scenario | Invalidation Strategy | TTL Recommendation |
|----------|---------------------|-------------------|
| FAQ content updated | | |
| Product prices changed | | |
| User-specific data | | |
| News articles | | |
| Code documentation | | |

---

## Problem 6: KV Cache Configuration

Configure KV cache for vLLM:

```yaml
# vllm_config.yaml
# Your configuration
engine:
  # KV cache settings
  ...

# Calculate memory for KV cache
# Model: 70B parameters (Llama 2)
# KV cache per token: ~1.5KB
# Max tokens in context: 4096

# How much KV cache memory is needed?
# How many concurrent requests can you serve?
```

Calculate:

| Parameter | Value |
|-----------|-------|
| Model size (GB) | |
| KV cache per token (KB) | |
| Max context tokens | |
| Total KV cache needed | |
| Available GPU memory | |
| Memory for KV cache | |
| Concurrent requests | |

---

## Problem 7: Cache Middleware Pattern

Implement a caching decorator:

```python
import time
import hashlib
from functools import wraps
import json

def llm_cache(cache_backend, ttl: int = 3600):
    """Decorator for caching LLM responses"""
    def decorator(func):
        @wraps(func)
        def wrapper(prompt, **kwargs):
            # 1. Generate cache key from prompt + kwargs
            # 2. Check cache
            # 3. If miss, call function and cache result
            pass
        return wrapper
    return decorator

# Usage with Redis cache
redis_cache = redis.Redis(host='localhost', port=6379)

@llm_cache(redis_cache, ttl=1800)
def generate_response(prompt: str, model: str = "gpt-4") -> str:
    return call_llm_api(prompt, model)
```

---

## Problem 8: Multi-Layer Caching

Design a caching architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Request                                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Layer 1: In-Memory LRU Cache                                   │
│  - Fast access (<1ms)                                           │
│  - Limited size                                                 │
│  - Best for: Repeated exact requests                           │
└─────────────────────────────────────────────────────────────────┘
                              │ Miss
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Layer 2: Redis Semantic Cache                                 │
│  - Semantic matching                                            │
│  - Distributed across instances                                 │
│  - TTL: 1 hour                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │ Miss
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Layer 3: LLM API (KV Cache)                                   │
│  - Provider-side caching                                       │
│  - Automatic                                                   │
└─────────────────────────────────────────────────────────────────┘
```

| Layer | Cache Type | Hit Criteria | TTL |
|-------|-----------|--------------|-----|
| 1 | | | |
| 2 | | | |
| 3 | | | |

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.
