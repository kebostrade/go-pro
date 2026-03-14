# Exercise: GPU Compute Fundamentals

## Problem 1: VRAM Allocation Analysis

A 7B parameter model running inference with INT8 quantization processes a request with 4,000 input tokens and generates 500 output tokens.

### Questions:

a) How much VRAM is needed for model weights?

b) Estimate the KV cache size (use 128 bytes per token per layer, 32 layers):

c) What percentage of an A100 40GB is used?

---

## Problem 2: Parallelism Strategy Selection

Match each scenario to the best parallelism strategy:

**Scenarios:**
1. Running a 70B model on 8x A100 GPUs
2. Running the same model on 4x smaller consumer GPUs  
3. Need minimal changes to training code

**Strategies:**
- Data Parallelism
- Tensor Parallelism  
- Pipeline Parallelism

| Scenario | Best Strategy | Why |
|----------|--------------|-----|
| 1 | | |
| 2 | | |
| 3 | | |

---

## Problem 3: Cloud Cost Comparison

Compare running inference on a 7B model (FP16) 24/7:

### Option A: Lambda Labs (Spot)
- GPU: A100 80GB
- Price: $1.89/hour (spot)
- Availability: ~95%

### Option B: AWS p4d
- GPU: 8x A100 40GB (already includes overhead)
- Price: $32.77/hour

### Option C: RunPod
- GPU: A100 80GB
- Price: $2.50/hour (spot)

#### Calculate monthly cost for each option (assume 730 hours):

| Provider | Hourly Cost | Monthly Cost | Notes |
|----------|-------------|--------------|-------|
| A | | | |
| B | | | |
| C | | | |

#### Which is most cost-effective and why?

---

## Problem 4: GPU Bottleneck Analysis

For each component, identify if it's CPU-bound or GPU-bound:

| Operation | CPU or GPU Bound | Reason |
|-----------|-----------------|--------|
| Tokenization | | |
| Model inference | | |
| KV cache storage | | |
| Response sampling | | |
| De-tokenization | | |

---

## Problem 5: Hardware Selection

You're building an AI platform that needs to serve multiple models:

| Model | Parameters | Precision | Max Batch Size |
|-------|------------|-----------|----------------|
| Model A | 7B | FP16 | 16 |
| Model B | 70B | INT8 | 4 |
| Model C | 3B | FP16 | 32 |

### Requirements:
- Must handle 100 concurrent requests
- Latency < 2 seconds
- Budget: $5,000/month for GPU costs

### Design:

How many GPUs? Which ones? How would you distribute the models?

```
┌─────────────────────────────────────────────────────────────────┐
│                        Infrastructure Design                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  GPUs:                                                           │
│                                                                  │
│  Model Distribution:                                             │
│                                                                  │
│  Load Balancing Strategy:                                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Problem 6: Understanding GPU Utilization

You observe the following metrics on a GPU running inference:

- GPU Utilization: 85%
- VRAM Used: 35GB / 80GB (44%)
- GPU Temperature: 72°C
- Power Draw: 280W / 350W (80%)

### Questions:

a) Is the GPU bottlenecked? What type of bottleneck?

b) What would you investigate to improve throughput?

c) If you wanted to double throughput, what's the limiting factor?

---

## Problem 7: CUDA Programming Basics

Write pseudocode for a simple operation: element-wise addition of two vectors on GPU.

```cuda
__global__ void vectorAdd(float* A, float* B, float* C, int N) {
    // Your code here
}

// Kernel launch
int threadsPerBlock = 256;
int blocksPerGrid = (N + threadsPerBlock - 1) / threadsPerBlock;
vectorAdd<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, N);
```

Explain what each line does:

- `__global__`
- `threadsPerBlock`
- `blocksPerGrid`

---

## Submission

Complete all problems. For Problem 5, draw a simple diagram of your infrastructure design.
