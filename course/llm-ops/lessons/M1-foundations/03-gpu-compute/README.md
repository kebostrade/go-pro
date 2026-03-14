# IO-03: GPU Compute Fundamentals

**Duration**: 3 hours
**Module**: 1 - Foundations

## Learning Objectives

- Understand GPU architecture and why it's essential for AI
- Learn about CUDA and GPU programming fundamentals
- Explore model parallelism strategies
- Evaluate cloud GPU options and costs

## Why GPUs for AI?

Central Processing Units (CPUs) are optimized for sequential processing and complex logic. Graphics Processing Units (GPUs) are optimized for parallel computation - perfect for the matrix multiplications that power neural networks.

### CPU vs GPU for AI Workloads

```
CPU (Central Processing Unit):
┌─────────────────────────────────────────────────────────────────┐
│  Core 1  │  Core 2  │  Core 3  │  Core 4  │ ... │ Core N     │
│  ┌───┐   │  ┌───┐   │  ┌───┐   │  ┌───┐   │     │  ┌───┐     │
│  │ + │   │  │ + │   │  │ + │   │  │ + │   │     │  │ + │     │
│  └───┘   │  └───┘   │  └───┘   │  └───┘   │     │  └───┘     │
│                                                                  │
│  Optimized for: Sequential tasks, complex branching            │
│  Latency: Low per task, but limited parallelism                │
└─────────────────────────────────────────────────────────────────┘

GPU (Graphics Processing Unit):
┌─────────────────────────────────────────────────────────────────┐
│  SM 1   │  SM 2   │  SM 3   │  SM 4   │ ... │ SM N             │
│  ┌──┐┌──┐│ ┌──┐┌──┐│ ┌──┐┌──┐│ ┌──┐┌──┐│     │ ┌──┐┌──┐      │
│  │╔╗╗╔╗│ ││╔╗╗╔╗│ ││╔╗╗╔╗│ ││╔╗╗╔╗│     │ ││╔╗╗╔╗│      │
│  │╚╝╚╝│ ││╚╝╚╝│ ││╚╝╚╝│ ││╚╝╚╝│     │ ││╚╝╚╝│      │
│  └──┬──┘ │ └──┬──┘ │ └──┬──┘ │ └──┬──┘     │ └──┬──┘      │
│     │     │    │     │    │         │        │                 │
│  1000s of cores optimized for parallel matrix operations       │
│  Optimized for: Throughput, parallel math                      │
└─────────────────────────────────────────────────────────────────┘
```

## NVIDIA GPU Architecture

### Key GPU Series for AI

| Series | Model | VRAM | TDP | Best For |
|--------|-------|------|-----|----------|
| **Ampere** | A100 | 40/80GB | 400W | Data center training |
| **Ada Lovelace** | L40S | 48GB | 350W | Inference, training |
| **Hopper** | H100 | 80GB | 700W | Cutting-edge AI |
| **Consumer** | RTX 4090 | 24GB | 450W | Development, small models |

### GPU Memory (VRAM) Considerations

VRAM is the bottleneck for AI workloads:

```
VRAM Usage Breakdown (Inference):
┌─────────────────────────────────────────────────────────────────┐
│                        Total VRAM Allocation                    │
├─────────────────────────────────────────────────────────────────┤
│  Model Weights  ████████████████████████████████████  60-70%   │
│  KV Cache       ██████████                                 15-20% │
│  Activations    ███████                                   10-15% │
│  Temp Buffers   ████                                       5-10% │
└─────────────────────────────────────────────────────────────────┘
```

## CUDA Programming

CUDA (Compute Unified Device Architecture) is NVIDIA's platform for parallel computing.

### Basic CUDA Concepts

```c
// Kernel: Function that runs on GPU
__global__ void matrixMultiply(float* A, float* B, float* C, int N) {
    // Each thread computes one output element
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    
    if (row < N && col < N) {
        float sum = 0;
        for (int k = 0; k < N; k++) {
            sum += A[row * N + k] * B[k * N + col];
        }
        C[row * N + col] = sum;
    }
}

// Launch: 256 threads per block, N/256 blocks
matrixMultiply<<<blocks, 256>>>(d_A, d_B, d_C, N);
```

### CUDA in Practice

For AI workloads, you rarely write raw CUDA. Instead, use:

- **CUDA Libraries**: cuBLAS (linear algebra), cuDNN (deep learning)
- **Frameworks**: PyTorch, TensorFlow handle CUDA automatically
- **Model Servers**: vLLM, TGI handle GPU optimization

## Model Parallelism Strategies

When a model doesn't fit on a single GPU, you need to distribute it.

### 1. Data Parallelism

Each GPU has a full copy of the model, processes different data batches.

```
GPU 1                    GPU 2                    GPU 3
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│ Model Copy  │         │ Model Copy  │         │ Model Copy  │
│             │         │             │         │             │
│ Batch 1     │         │ Batch 2     │         │ Batch 3     │
│ Forward     │         │ Forward     │         │ Forward     │
│             │         │             │         │             │
│ ─────────── │         │ ─────────── │         │ ─────────── │
│ Gradients ──┼────────▶│ Gradients ──┼────────▶│ Gradients ──┤
│             │ Average │             │         │             │
└─────────────┘         └─────────────┘         └─────────────┘
          ◀────────────── AllReduce ──────────────▶
```

**Pros**: Easy to implement, scales well
**Cons**: Requires model to fit in each GPU

### 2. Tensor Parallelism

Split model layers across GPUs (for a single layer).

```
Layer N: Linear
┌──────────────────┐
│   Input Tensor   │
└────────┬─────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐ ┌───────┐
│ Col 1 │ │ Col 2 │   // Split weight matrix
│  of   │ │  of   │
│   W   │ │   W   │
└───┬───┘ └───┬───┘
    ▼         ▼
 AllReduce ◀──▶ AllReduce
    │         │
    └────┬────┘
         ▼
  Output Tensor
```

**Pros**: Enables huge models
**Cons**: Complex, high communication overhead

### 3. Pipeline Parallelism

Split model layers across GPUs sequentially.

```
GPU 0         GPU 1         GPU 2         GPU 3
┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐
│ Layers │   │ Layers │   │ Layers │   │ Layers │
│  1-8   │   │  9-16  │   │ 17-24  │   │ 25-32  │
└────────┘   └────────┘   └────────┘   └────────┘
    │            │            │            │
    ▼            ▼            ▼            ▼
Input      Pipeline      Pipeline      Output
```

**Pros**: Lower communication than tensor parallelism
**Cons**: Pipeline bubbles (idle time)

## Cloud GPU Options

### Comparison of GPU Cloud Providers

| Provider | GPU Options | Price (A100/hr) | Min. Commitment |
|----------|-------------|-----------------|-----------------|
| **AWS EC2** | P4d, P5 | ~$30-40 | None |
| **GCP** | A100, H100 | ~$35-45 | None |
| **Azure** | ND A100 v4 | ~$35 | None |
| **Lambda Labs** | A100, H100 | ~$1.50/hr spot | None |
| **RunPod** | A100, H100 | ~$2-3/hr | None |
| **Paperspace** | A100, H100 | ~$3-4/hr | None |

### Cost Optimization Strategies

1. **Spot/Preemptible Instances**: 60-80% discount, can be interrupted
2. **Savings Plans**: Commit to usage for discounts
3. **Multi-cloud**: Compare prices across providers
4. **Right-sizing**: Match GPU to actual workload

## GPU Selection for Different Workloads

| Workload | Recommended GPU | Reason |
|----------|-----------------|--------|
| Development/Testing | RTX 4090 | Cost-effective for small workloads |
| Production Inference (small models) | A100 40GB | Good balance of cost and performance |
| Production Inference (large models) | H100 80GB | Best inference performance |
| Fine-tuning (7B) | A100 80GB | Need VRAM for gradients |
| Fine-tuning (70B) | 8x H100 | Multi-GPU training required |

## Key Terminology

| Term | Definition |
|------|------------|
| **VRAM** | Video RAM - memory on the GPU |
| **CUDA** | NVIDIA's parallel computing platform |
| **cuDNN** | NVIDIA's deep learning library |
| **Tensor Core** | Specialized GPU cores for matrix multiplication |
| **SM (Streaming Multiprocessor)** | GPU compute unit |
| **Data Parallelism** | Same model, different data on each GPU |
| **Tensor Parallelism** | Split layer weights across GPUs |
| **Pipeline Parallelism** | Split model layers across GPUs |
| **AllReduce** | Collective operation to combine gradients |
| **KV Cache** | Cache for key-value attention computations |

## Exercise

### Exercise 3.1: GPU Memory Calculation

For a Llama 3 70B model running inference in INT8:

1. How much VRAM is needed for model weights?
2. If using a 128K context, approximately how much VRAM is needed for KV cache?
3. What's the total VRAM requirement?
4. How many A100 80GB GPUs are needed?

### Exercise 3.2: Parallelism Strategy

You're deploying a 70B parameter model. Compare:

| Strategy | GPUs Needed | Pros | Cons |
|----------|-------------|------|------|
| Data Parallelism | | | |
| Tensor Parallelism | | | |
| Pipeline Parallelism | | | |

### Exercise 3.3: Cost Analysis

Calculate the monthly cost for:

1. Running inference on a 7B model on a single RTX 4090 (24GB) for 24/7
2. Running the same workload on AWS p4d.24xlarge

Assume:
- RTX 4090: $0.50/hour (cloud)
- p4d.24xlarge: $32/hour
- 24/7 operation (730 hours/month)

### Exercise 3.4: Architecture Decision

A company needs to serve Llama 3 70B with <2 second latency. They have a $10K/month budget.

Design the infrastructure:
- Which GPUs?
- How many?
- What parallelism strategy?
- Estimated cost?

## Key Takeaways

- ✅ GPUs excel at parallel matrix operations essential for AI
- ✅ VRAM is the primary bottleneck for AI workloads
- ✅ Model parallelism enables running models too large for single GPU
- ✅ Cloud GPUs offer flexibility but require cost optimization strategies
- ✅ Choice of GPU depends on model size, latency requirements, and budget

## Next Steps

→ [IO-04: Cloud AI Platforms](../04-cloud-platforms/README.md)

## Additional Resources

- [NVIDIA CUDA Documentation](https://docs.nvidia.com/cuda/)
- [Deep Learning Performance Guide](https://docs.nvidia.com/deeplearning/performance/)
- [vLLM Architecture](https://docs.vllm.ai/)
- [Cloud GPU Pricing Comparison](https://www.runpod.io/gpu-cloud-comparison)
