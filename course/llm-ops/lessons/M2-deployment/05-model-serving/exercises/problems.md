# Exercise: Model Serving Architectures

## Problem 1: Serving Server Comparison

Compare the three model serving solutions (vLLM, TGI, Triton) for the following scenarios:

| Scenario | Recommended Server | Justification |
|----------|-------------------|---------------|
| High-volume chatbot serving 10k+ concurrent users | | |
| Running quantized models on limited GPU memory | | |
| Need to serve multiple model types (PyTorch, TensorRT) | | |
| Maximum throughput for a single model deployment | | |
| Quick prototyping with Hugging Face models | | |

---

## Problem 2: vLLM Deployment

Write a docker-compose.yml to deploy vLLM with the following requirements:

1. Single GPU deployment
2. Expose API on port 8000
3. Use Llama-2-7b model
4. Set GPU memory utilization to 90%
5. Include health check endpoint

### Your Answer:

```yaml

```

---

## Problem 3: OpenAI-Compatible Client

Write a Go program that:

1. Connects to a vLLM server at `http://localhost:8000/v1`
2. Sends a chat completion request
3. Handles streaming responses
4. Implements basic retry logic (3 retries with exponential backoff)

### Your Answer:

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// Implement the client

func main() {
	// Test your implementation
}
```

---

## Problem 4: Performance Optimization

Given a vLLM deployment with the following metrics:

- Average tokens/sec: 50
- Average prompt tokens: 200
- Average completion tokens: 150
- Concurrent requests: 100
- GPU memory utilization: 85%

Answer these questions:

1. What is the estimated time to first token (TTFT)?
2. What is the estimated time per output token?
3. What optimization would you recommend to increase throughput?
4. If you enable continuous batching and see 30% improvement, what are the new throughput metrics?

---

## Problem 5: API Gateway Integration

The [API Gateway in this repository](services/api-gateway/) provides authentication and rate limiting.

Design a Go middleware that:

1. Validates JWT tokens before forwarding to vLLM
2. Applies per-user rate limiting (e.g., 100 requests/minute)
3. Logs all requests with latency metrics
4. Handles errors gracefully

### Your Answer:

```go
package middleware

import (
	"net/http"
	"time"
)

// RateLimiter interface
type RateLimiter interface {
	Allow(key string) bool
}

// JWTAuth interface
type JWTAuth interface {
	Validate(token string) (string, error)
}

func NewLLMMiddleware(auth JWTAuth, limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement middleware logic
			
			next.ServeHTTP(w, r)
		})
	}
}
```

---

## Problem 6: Benchmarking LLM Performance

Write a Python script to benchmark an LLM inference server:

1. Measure throughput (tokens/second)
2. Measure latency (TTFT, total latency)
3. Measure error rate
4. Test with varying batch sizes

### Your Answer:

```python
import asyncio
import time
import aiohttp
from typing import List, Dict

async def benchmark_endpoint(
    url: str,
    num_requests: int,
    concurrent: int
) -> Dict:
    """Benchmark LLM endpoint"""
    # Implement
    
    return {
        "total_requests": 0,
        "successful": 0,
        "failed": 0,
        "avg_latency": 0.0,
        "tokens_per_second": 0.0
    }

if __name__ == "__main__":
    results = asyncio.run(benchmark_endpoint(
        "http://localhost:8000/v1/chat/completions",
        num_requests=100,
        concurrent=10
    ))
    print(results)
```

---

## Submission

Save your answers and be prepared to discuss them in the next lesson.
