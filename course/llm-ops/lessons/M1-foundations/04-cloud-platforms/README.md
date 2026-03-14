# IO-04: Cloud AI Platforms

**Duration**: 2 hours
**Module**: 1 - Foundations

## Learning Objectives

- Evaluate major cloud AI platforms (AWS, GCP, Azure)
- Understand managed inference APIs and their use cases
- Compare self-hosted vs. managed approaches
- Make infrastructure decisions based on requirements

## Major Cloud AI Platforms

### Overview Comparison

| Platform | AI Service Name | Strengths | Weaknesses |
|----------|---------------|-----------|-------------|
| **AWS** | SageMaker | Enterprise features, broad ecosystem | Complex pricing |
| **GCP** | Vertex AI | Strong ML tools, TPU access | Smaller market share |
| **Azure** | Azure ML | Microsoft integration, enterprise ready | Less optimized for ML |

## AWS SageMaker

### Key Services

```
┌─────────────────────────────────────────────────────────────────┐
│                        AWS AI Infrastructure                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  SageMaker     │  │  SageMaker     │  │  Bedrock      │  │
│  │  JumpStart      │  │  Endpoints     │  │  (Managed LLMs)│  │
│  │  (Pre-trained)  │  │  (Real-time)   │  │               │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  SageMaker     │  │  SageMaker     │  │  EC2          │  │
│  │  Batch         │  │  Processing    │  │  P4d/P5       │  │
│  │  Transform     │  │  (ETL)         │  │  (GPU VMs)    │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### SageMaker Features

- **JumpStart**: Pre-trained models ready for deployment
- **Real-time Endpoints**: Low-latency inference
- **Serverless Inference**: Pay-per-request, no infrastructure management
- **Multi-model endpoints**: Serve multiple models on one endpoint
- **SageMaker Canvas**: No-code ML

### Example: Deploying a Model on SageMaker

```python
import boto3
import sagemaker
from sagemaker.huggingface import HuggingFaceModel

# Create HuggingFace model
huggingface_model = HuggingFaceModel(
    model_data='s3://bucket/model.tar.gz',
    role=role,
    transformers_version='4.26',
    pytorch_version='1.13',
    py_version='py310'
)

# Deploy to real-time endpoint
predictor = huggingface_model.deploy(
    initial_instance_count=1,
    instance_type='ml.g5.xlarge'
)

# Make inference
response = predictor.predict({
    "inputs": "Analyze this transaction for fraud"
})
```

## Google Cloud Platform (GCP) Vertex AI

### Key Services

```
┌─────────────────────────────────────────────────────────────────┐
│                        GCP AI Infrastructure                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  Vertex AI    │  │  Vertex AI    │  │  Prediction   │  │
│  │  Workbench    │  │  Endpoints    │  │  API          │  │
│  │  (Notebooks)   │  │  (Managed)    │  │  (REST API)   │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  Model Garden │  │  AutoML        │  │  TPU          │  │
│  │  (Pre-trained) │  │  (No-code)     │  │  (Training)  │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Vertex AI Features

- **Model Garden**: Access to hundreds of pre-trained models
- **AutoML**: Train custom models without coding
- **Prediction**: Managed inference with autoscaling
- **TensorFlow Extended (TFX)**: ML pipelines
- **TPUs**: Google custom AI accelerators

### Example: Deploying on Vertex AI

```python
from google.cloud import aiplatform

# Initialize
aiplatform.init(project='my-project', location='us-central1')

# Deploy model
endpoint = aiplatform.Endpoint.create(
    display_name='fraud-detection-endpoint',
    deployed_model_id='fraud-model'
)

# Predict
response = endpoint.predict(instances=[{
    "transaction_id": "12345",
    "amount": 5000.00,
    "merchant": "electronics_store"
}])
```

## Microsoft Azure ML

### Key Services

```
┌─────────────────────────────────────────────────────────────────┐
│                      Azure AI Infrastructure                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  Azure ML     │  │  Azure ML     │  │  Azure OpenAI │  │
│  │  Compute      │  │  Endpoints    │  │  Service      │  │
│  │  (Clusters)   │  │  (Real-time)  │  │  (GPT access) │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌───────────────┐  │
│  │  Azure ML     │  │  Azure AI     │  │  Cognitive    │  │
│  │  Designer     │  │  Studio        │  │  Services     │  │
│  │  (No-code)    │  │  (Web portal) │  │  (Pre-built)  │  │
│  └─────────────────┘  └─────────────────┘  └───────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Azure ML Features

- **Azure OpenAI Service**: Access to GPT models with enterprise security
- **Designer**: Visual ML workflow builder
- **Azure AI Studio**: Unified AI development platform
- **MLOps**: CI/CD for ML workflows
- **Enterprise integration**: Active Directory, compliance

## Managed Inference APIs

### Comparison of LLM APIs

| Provider | API | Model | Strengths |
|----------|-----|-------|-----------|
| **OpenAI** | GPT API | GPT-4, GPT-4o | Best overall capability |
| **Anthropic** | Claude API | Claude 3 | Long context, helpful |
| **Google** | Gemini API | Gemini Pro | Multimodal, fast |
| **Meta** | Llama API | Llama 3 | Open weights, cost |
| **Mistral** | La Plateforme | Mistral | Open weights, fast |

### When to Use Managed APIs

- **Rapid prototyping**: Fast to start
- **Variable traffic**: Pay-per-use scales automatically
- **Limited DevOps**: No infrastructure management
- **Latest models**: Access to newest capabilities

### When to Self-Host

- **High volume**: When API costs exceed infrastructure costs
- **Data privacy**: Data cannot leave your environment
- **Customization**: Need fine-tuned models
- **Latency**: Need <100ms response times

## Cost Comparison

### Managed API Costs (Approximate)

| Model | Input/1K tokens | Output/1K tokens | Notes |
|-------|-----------------|------------------|-------|
| GPT-4o | $2.50 | $10.00 | Latest model |
| GPT-4 | $15.00 | $60.00 | Previous generation |
| Claude 3 Opus | $15.00 | $75.00 | Long context |
| Claude 3 Haiku | $0.25 | $1.25 | Fast, cheap |
| Gemini 1.5 Pro | $1.25 | $5.00 | 1M context |

### Self-Hosted vs. API Cost Break-even

```
Cost Analysis: Self-Hosted 7B vs. API

┌─────────────────────────────────────────────────────────────────┐
│                     Requests per Month                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Cost ($)                                                       │
│  │                                                               │
│  │                    ╲                                         │
│  │               ╱────╲  ← Self-hosted (fixed cost)            │
│  │          ╱───         ╲                                       │
│  │     ╱──               ╲───── Managed API (per-request)        │
│  │───                                                          │
│  └─────────────────────────────────────────────────────────────▶ │
│       0          1M         2M         3M     Requests       │
│                                                                  │
│  At ~1.5M requests/month, self-hosted becomes cheaper           │
└─────────────────────────────────────────────────────────────────┘
```

## Infrastructure Decision Framework

```
┌─────────────────────────────────────────────────────────────────┐
│                  Infrastructure Decision Tree                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Start: What's your primary constraint?                          │
│           │                                                     │
│     ┌─────┴─────┐                                               │
│     ▼           ▼                                               │
│  Privacy    Performance                                         │
│     │           │                                               │
│     ▼           ▼                                               │
│  Self-Host   Latency < 500ms?                                   │
│               │     │                                           │
│               ▼     ▼                                           │
│              Yes   No                                           │
│                │     │                                           │
│                ▼     ▼                                         │
│            Self-Host  Budget < $5K/mo?                           │
│                           │     │                                │
│                           ▼     ▼                               │
│                         Yes    No                               │
│                           │     │                                │
│                           ▼     ▼                               │
│                       Managed   Managed API                      │
│                       (vLLM)    or Self-Host                                            │
└─────────────────────────────────────────────────────────────────┘
```

## Key Terminology

| Term | Definition |
|------|------------|
| **SageMaker** | AWS's comprehensive ML platform |
| **Vertex AI** | GCP's ML platform with Model Garden |
| **Azure ML** | Microsoft's ML platform |
| **Managed Endpoint** | Cloud-hosted inference with autoscaling |
| **Serverless Inference** | Pay-per-request without provisioning |
| **Model Registry** | Centralized model versioning and management |
| **MLOps** | Machine learning DevOps practices |
| **Fine-tuning** | Adapting pre-trained models to specific tasks |

## Exercise

### Exercise 4.1: Platform Comparison

For each use case, recommend the best platform and explain why:

| Use Case | Recommended Platform | Justification |
|----------|---------------------|---------------|
| Need GPT-4 access with enterprise compliance | | |
| Running open-source models at scale | | |
| TPU access for training | | |
| Deep Microsoft ecosystem integration | | |

### Exercise 4.2: Cost Analysis

A company expects 500,000 requests per month, with average 1,000 input tokens and 500 output tokens per request.

Calculate:

1. **Using OpenAI GPT-4o API**:
   - Input cost: 
   - Output cost:
   - Total monthly cost:

2. **Self-hosted on Lambda Labs** (A100 80GB at $1.89/hr):
   - Hours needed: (estimate 50 req/hour per GPU)
   - Monthly cost:

3. **Break-even point**: At what request volume does self-hosted become cheaper?

### Exercise 4.3: Architecture Design

Design a multi-tier inference architecture:

- Tier 1: Simple queries → Cache or smaller model
- Tier 2: Complex queries → Full model
- Tier 3: Compliance-required → Self-hosted

What cloud services would you use?

---

## Key Takeaways

- ✅ AWS SageMaker, GCP Vertex AI, and Azure ML offer comprehensive ML platforms
- ✅ Managed APIs are best for prototyping and variable workloads
- ✅ Self-hosting becomes cost-effective at high volumes
- ✅ Platform choice depends on ecosystem, compliance, and cost requirements

## Next Steps

Complete Module 1 exercises and review. Then proceed to Module 2: LLM Deployment & Serving.

## Additional Resources

- [AWS SageMaker Documentation](https://docs.aws.amazon.com/sagemaker/)
- [GCP Vertex AI Documentation](https://cloud.google.com/vertex-ai)
- [Azure ML Documentation](https://learn.microsoft.com/azure/machine-learning/)
- [OpenAI Pricing](https://openai.com/pricing)
- [Cloud Cost Calculator](https://cloudcostcalculator.com)
