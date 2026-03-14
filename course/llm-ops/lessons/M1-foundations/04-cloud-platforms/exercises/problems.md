# Exercise: Cloud AI Platforms

## Problem 1: Platform Selection Matrix

For each scenario, select the most appropriate cloud platform and service:

| Scenario | Platform | Service | Why |
|----------|----------|---------|-----|
| Need access to GPT-4 with enterprise compliance | | | |
| Want to run Llama 3 at scale with minimal DevOps | | | |
| Have existing Azure infrastructure | | | |
| Need TPU access for training | | | |
| Want to deploy custom PyTorch model | | | |
| Building a no-code ML solution | | | |

---

## Problem 2: Total Cost of Ownership

Compare two deployment options for a customer support chatbot:

### Option A: Managed API (OpenAI GPT-4o)
- Monthly requests: 2,000,000
- Average input tokens: 500
- Average output tokens: 200

### Option B: Self-Hosted (AWS SageMaker with g5.xlarge)
- Instance cost: $1.01/hour
- Can handle 100 requests/minute per instance
- Need 4 instances for peak load

### Calculate:

| Metric | Option A (API) | Option B (Self-Hosted) |
|--------|---------------|----------------------|
| Monthly Compute Cost | | |
| Other Costs (storage, etc.) | | |
| **Total Monthly Cost** | | |

### Which is more cost-effective?

---

## Problem 3: Multi-Cloud Architecture

Design a system that uses:
- AWS for primary inference
- GCP for data processing/embedding generation
- Azure for compliance-required workloads

Draw the architecture and explain data flow:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Architecture Diagram                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  User Request                                                    │
│     │                                                            │
│     ▼                                                            │
│  ┌──────────┐                                                   │
│  │  AWS     │  ← Primary inference (OpenAI or self-hosted)    │
│  │  API GW  │                                                   │
│  └────┬─────┘                                                   │
│       │                                                         │
│       ├──────────────────┐                                       │
│       ▼                  ▼                                       │
│  ┌──────────┐      ┌──────────┐                                 │
│  │ GCP      │      │ Azure    │  ← Compliance workloads        │
│  │ Embedding│      │ Audit    │                                 │
│  │ Service  │      │ Logging  │                                 │
│  └──────────┘      └──────────┘                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Problem 4: API vs Self-Hosted Decision

A startup is building an AI-powered feature for their product:

**Current State:**
- 50,000 monthly active users
- Average 10 AI requests per user per month (500,000 requests)
- Average 200 tokens input, 150 tokens output

**Growth Projection:**
- 2x users in 6 months
- 3x requests per user in 12 months

### Questions:

a) What is the current monthly cost using OpenAI API?

b) At what user/request volume does self-hosting become cost-effective?

c) What would you recommend for their 12-month projection?

---

## Problem 5: Compliance Requirements

A healthcare company needs to process patient data with AI:

**Requirements:**
- HIPAA compliance required
- Data must stay within US borders
- Audit logging for all AI decisions
- Need to support 1 million requests/month

### Questions:

a) Which cloud platform would you recommend and why?

b) What specific services would you use?

c) What additional security measures would you implement?

---

## Problem 6: Performance Benchmarking

Design a benchmark to compare:
- AWS SageMaker endpoint
- GCP Vertex AI endpoint
- Azure ML endpoint

For each metric, describe how you would measure it:

| Metric | Measurement Approach |
|--------|---------------------|
| Latency (p50) | |
| Latency (p99) | |
| Throughput (req/sec) | |
| Cold start time | |
| Cost per 1K requests | |

---

## Problem 7: Vendor Lock-in Mitigation

Design a multi-cloud abstraction layer that allows switching between providers:

```go
// Example: Abstract LLM provider interface
type LLMProvider interface {
    Generate(ctx context.Context, req Request) (*Response, error)
    GetProviderName() string
}

// Implementations for each provider
type OpenAIProvider struct{}
type AnthropicProvider struct{}
type VertexAIProvider struct{}
type AzureProvider struct{}
```

### Questions:

a) What are the benefits of this abstraction?

b) What challenges might you face when switching providers?

c) How does this relate to the AI Agent Platform architecture?

---

## Submission

Complete all problems. For Problem 3, draw a more detailed architecture diagram showing data flows.
