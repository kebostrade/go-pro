# Exercise: Introduction to AI Infrastructure

## Problem 1: Infrastructure Component Identification

For each component below, identify which layer of the AI infrastructure stack it belongs to:

1. NVIDIA H100 GPU
2. Kubernetes cluster
3. PostgreSQL with pgvector
4. vLLM inference server
5. Redis cache
6. S3 object storage
7. OpenAI API
8. Prometheus monitoring

### Your Answer:

| Component | Infrastructure Layer |
|-----------|---------------------|
| | |
| | |
| | |
| | |
| | |
| | |
| | |
| | |

---

## Problem 2: Self-Hosted vs Managed Decision

A fintech company needs to build a fraud detection system using LLMs. They have:
- Strict data privacy requirements (cannot send data to external APIs)
- Need to process 1 million transactions per day
- Have a team of 3 engineers with Kubernetes experience
- Budget of $5,000/month for infrastructure

Should they use self-hosted or managed services? Justify your answer with specific reasons.

### Your Justification:

```

```

---

## Problem 3: Design Infrastructure for RAG System

Design the infrastructure for a RAG (Retrieval-Augmented Generation) system that:
- Stores company documents (approximately 10GB)
- Supports semantic search on documents
- Uses an LLM to answer questions based on retrieved context
- Needs to handle 100 concurrent users

Create a simple architecture diagram and list the components needed.

### Architecture:

```

(Use text to describe your architecture)

```

### Components:

| Component | Purpose | Recommended Technology |
|-----------|---------|----------------------|
| | | |
| | | |
| | | |
| | | |
| | | |

---

## Problem 4: Cost Estimation

Estimate the monthly cost for self-hosting an LLM inference service with the following specs:
- 1x NVIDIA A100 GPU (PCIe, 80GB)
- 16 vCPU, 64GB RAM
- 500GB SSD storage
- 10TB bandwidth

Use current cloud pricing from AWS, GCP, or Lambda Labs.

### Your Estimate:

| Resource | Provider | Price/Month |
|----------|----------|-------------|
| | | |
| | | |
| | | |
| | | |
| **Total** | | |

---

## Problem 5: Infrastructure Scalability Analysis

The AI Agent Platform in this repository uses the following architecture:
- API Gateway for auth and rate limiting
- Agent Service for orchestration
- LLM Service with provider abstraction
- Tool Registry for business logic
- Vector Store for semantic search

What happens to this architecture when:
1. Traffic increases 10x
2. One LLM provider goes down
3. The vector database becomes the bottleneck

Describe how the infrastructure would need to evolve for each scenario.

### Scenario 1 - Traffic 10x:

```

```

### Scenario 2 - LLM Provider Down:

```

```

### Scenario 3 - Vector DB Bottleneck:

```

```

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.
