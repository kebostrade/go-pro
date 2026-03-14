# AI Infrastructure & LLM-Ops

A comprehensive, hands-on course on building and operating production AI infrastructure, covering foundations to advanced LLM operations and observability.

## Course Overview

**Duration**: 10 weeks (self-paced)
**Level**: Intermediate to Advanced
**Prerequisites**: Basic knowledge of cloud computing, containers, and API design. Familiarity with AI/LLMs helpful.

## Learning Objectives

By the end of this course, you will be able to:
- Design and provision AI-ready infrastructure on cloud platforms
- Deploy and serve LLMs using model serving frameworks
- Build and optimize RAG systems with vector databases
- Implement observability, monitoring, and logging for AI systems
- Optimize costs through caching, batching, and efficient resource allocation
- Secure AI infrastructure against common vulnerabilities
- Operationalize LLMs in production with proper CI/CD and governance

## Course Structure

### [Module 1: Foundations of AI Infrastructure](lessons/M1-foundations/README.md) (Weeks 1-2)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [IO-01](lessons/M1-foundations/01-introduction/README.md) | Introduction to AI Infrastructure | 2h |
| [IO-02](lessons/M1-foundations/02-llm-architecture/README.md) | LLM Architecture: Transformers, Attention, Quantization | 3h |
| [IO-03](lessons/M1-foundations/03-gpu-compute/README.md) | GPU Compute: CUDA, VRAM, Model Parallelism | 3h |
| [IO-04](lessons/M1-foundations/04-cloud-platforms/README.md) | Cloud AI Platforms: AWS, GCP, Azure, Lambda | 2h |

### [Module 2: LLM Deployment & Serving](lessons/M2-deployment/README.md) (Weeks 3-4)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [IO-05](lessons/M2-deployment/05-model-serving/README.md) | Model Serving: vLLM, TensorRT-LLM, TGI | 3h |
| [IO-06](lessons/M2-deployment/06-api-gateway/README.md) | API Gateway: Rate Limiting, Auth, Load Balancing | 2h |
| [IO-07](lessons/M2-deployment/07-containerization/README.md) | Containerization: Docker, CUDA Images, Multi-stage Builds | 3h |
| [IO-08](lessons/M2-deployment/08-kubernetes/README.md) | Kubernetes for AI: K8s, Kueue, KServe | 3h |

### [Module 3: Vector Databases & RAG](lessons/M3-vector-rag/README.md) (Weeks 5-6)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [IO-09](lessons/M3-vector-rag/09-vector-databases/README.md) | Vector Databases: pgvector, Milvus, Qdrant, Weaviate | 3h |
| [IO-10](lessons/M3-vector-rag/10-embedding-models/README.md) | Embedding Models: Selection, Fine-tuning, Evaluation | 2h |
| [IO-11](lessons/M3-vector-rag/11-rag-architecture/README.md) | RAG Architecture: Indexing, Retrieval, Fusion | 3h |
| [IO-12](lessons/M3-vector-rag/12-rag-optimization/README.md) | RAG Optimization: Hybrid Search, Re-ranking, Chunking | 2h |

### [Module 4: LLM Operations & Observability](lessons/M4-operations/README.md) (Week 7)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [IO-13](lessons/M4-operations/13-monitoring/README.md) | Monitoring: Prometheus, Grafana, LangSmith | 3h |
| [IO-14](lessons/M4-operations/14-prompt-caching/README.md) | Prompt Caching: Redis, Semantic Caching, KV Cache | 2h |
| [IO-15](lessons/M4-operations/15-cost-optimization/README.md) | Cost Optimization: Token Economics, Batching, Spot Instances | 2h |
| [IO-16](lessons/M4-operations/16-security/README.md) | Security: PII Redaction, Input Validation, RBAC | 2h |

### [Module 5: Advanced LLM-Ops & Production](lessons/M5-production/README.md) (Weeks 8-10)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [IO-17](lessons/M5-production/17-fine-tuning/README.md) | Fine-tuning: LoRA, QLoRA, PEFT, Distributed Training | 3h |
| [IO-18](lessons/M5-production/18-model-versioning/README.md) | Model Versioning: MLflow, Weights & Biases, Model Registry | 2h |
| [IO-19](lessons/M5-production/19-testing-evaluation/README.md) | Testing & Evaluation: LLM Benchmarks, A/B Testing | 3h |
| [IO-20](lessons/M5-production/20-production/README.md) | Production Systems: CI/CD, Multi-region, Disaster Recovery | 3h |

## Projects

| Project | Description | Difficulty |
|---------|-------------|------------|
| [P1](projects/P1-infra-provisioning/README.md) | Provision Cloud AI Infrastructure with Terraform | Intermediate |
| [P2](projects/P2-llm-serving-cluster/README.md) | Deploy LLM Serving Cluster with Kubernetes | Advanced |
| [P3](projects/P3-rag-system/README.md) | Build Production RAG System with Vector DB | Advanced |
| [P4](projects/P4-llm-observability/README.md) | Implement LLM Observability Platform | Expert |

## Quick Start

```bash
# Clone or navigate to the course
cd course/llm-ops

# Start with Module 1: Foundations
cat lessons/M1-foundations/README.md
cat lessons/M1-foundations/01-introduction/README.md

# Explore the AI Agent Platform in this repository
cd ../services/ai-agent-platform
cat README.md

# Run the fraud detection example
make run-example
```

## Tools & Resources

### GPU Cloud Providers
- [Lambda Labs](https://lambdalabs.com) - GPU cloud for AI training/inference
- [Paperspace](https://paperspace.com) - GPU-accelerated cloud
- [RunPod](https://runpod.io) - Serverless GPU infrastructure
- [Vast.ai](https://vast.ai) - GPU marketplace

### Vector Databases
- [pgvector](https://github.com/pgvector/pgvector) - PostgreSQL vector extension
- [Milvus](https://milvus.io) - Open source vector database
- [Qdrant](https://qdrant.io) - Vector similarity search engine
- [Weaviate](https://weaviate.io) - Open source vector database

### Model Serving
- [vLLM](https://github.com/vllm-project/vllm) - High-performance LLM serving
- [Text Generation Inference](https://github.com/huggingface/text-generation-inference)
- [TensorRT-LLM](https://developer.nvidia.com/tensorrt-llm) - NVIDIA inference engine
- [KServe](https://kserve.github.io/website) - Kubernetes model server

### Monitoring & Observability
- [LangSmith](https://smith.langchain.com) - LLM tracing and evaluation
- [Prometheus](https://prometheus.io) - Metrics collection
- [Grafana](https://grafana.com) - Visualization and dashboards
- [OpenTelemetry](https://opentelemetry.io) - Observability standard

### Infrastructure as Code
- [Terraform](https://terraform.io) - Infrastructure provisioning
- [Kubernetes](https://kubernetes.io) - Container orchestration
- [Helm](https://helm.sh) - Package manager for K8s

### Reference Implementation
This course references the [AI Agent Platform](../services/ai-agent-platform/README.md) in this repository, a production-ready financial services AI framework demonstrating:
- ReAct and Base agent implementations
- Tool system with financial and general tools
- LLM provider abstraction (OpenAI, caching)
- Production-grade patterns for AI systems

## File Structure

```
llm-ops/
├── README.md                 # This file
├── cheatsheet.md            # Quick reference
├── lessons/                 # Lesson content (organized by module)
│   ├── M1-foundations/
│   │   ├── 01-introduction/
│   │   ├── 02-llm-architecture/
│   │   ├── 03-gpu-compute/
│   │   ├── 04-cloud-platforms/
│   │   └── README.md
│   ├── M2-deployment/
│   │   ├── 05-model-serving/
│   │   ├── 06-api-gateway/
│   │   ├── 07-containerization/
│   │   ├── 08-kubernetes/
│   │   └── README.md
│   ├── M3-vector-rag/
│   │   ├── 09-vector-databases/
│   │   ├── 10-embedding-models/
│   │   ├── 11-rag-architecture/
│   │   ├── 12-rag-optimization/
│   │   └── README.md
│   ├── M4-operations/
│   │   ├── 13-monitoring/
│   │   ├── 14-prompt-caching/
│   │   ├── 15-cost-optimization/
│   │   ├── 16-security/
│   │   └── README.md
│   └── M5-production/
│       ├── 17-fine-tuning/
│       ├── 18-model-versioning/
│       ├── 19-testing-evaluation/
│       ├── 20-production/
│       └── README.md
├── exercises/               # Practice problems
├── examples/                # Code examples (reference to ai-agent-platform)
├── projects/                # Course projects
└── resources/              # Additional materials
```

## Certificate Requirements

Complete all modules + 2 projects to earn certificate:
- [ ] All 20 lessons reviewed
- [ ] 80% exercise completion
- [ ] 2 projects completed with documentation

---

*Last Updated: March 2026*