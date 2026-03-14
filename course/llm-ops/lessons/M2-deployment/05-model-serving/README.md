# IO-05: Model Serving Architectures

**Duration**: 3 hours
**Module**: 2 - Deployment & Serving

## Learning Objectives

- Understand the architecture of modern LLM inference servers
- Compare vLLM, Text Generation Inference (TGI), and Triton Inference Server
- Implement OpenAI-compatible APIs for LLM serving
- Optimize inference performance with batching and caching

## What is Model Serving?

Model serving is the process of deploying and running trained machine learning models for inference. For LLMs, this involves:

- Loading model weights into memory (often requiring multiple GPUs)
- Processing inference requests with optimal throughput
- Managing context windows and token limits
- Handling streaming responses efficiently

```
┌─────────────────────────────────────────────────────────────────┐
│                      Client Requests                             │
│            (Chat completions, embeddings, streaming)           │
├─────────────────────────────────────────────────────────────────┤
│                    Model Serving Layer                           │
│        ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│        │    vLLM      │  │     TGI      │  │   Triton     │    │
│        │              │  │              │  │              │    │
│        │ PagedAttention│  │  Continuous │  │  Dynamic    │    │
│        │   AFC         │  │  Batching    │  │  Batching   │    │
│        └──────────────┘  └──────────────┘  └──────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│                      Model Layer                                 │
│              (Llama, Mistral, GPT, Claude weights)             │
├─────────────────────────────────────────────────────────────────┤
│                      Hardware Layer                              │
│                  (GPU Memory, Compute, NVLink)                 │
└─────────────────────────────────────────────────────────────────┘
```

## 1. vLLM - Virtual LLM

vLLM is an open-source inference server developed by UC Berkeley that achieves state-of-the-art throughput through PagedAttention.

### Key Features

- **PagedAttention**: Memory-efficient attention mechanism that reduces KV cache memory by up to 2x
- **Continuous Batching**: Dynamically adds requests to batches mid-generation
- **Tensor Parallelism**: Distributed inference across multiple GPUs
- **OpenAI-Compatible API**: Drop-in replacement for OpenAI API

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     vLLM Server                                  │
├─────────────────────────────────────────────────────────────────┤
│  API Server                                                     │
│  ├── /v1/chat/completions                                      │
│  ├── /v1/completions                                           │
│  └── /v1/embeddings                                           │
├─────────────────────────────────────────────────────────────────┤
│  Scheduler                                                      │
│  ├── Request Queue (async)                                     │
│  ├── Continuous Batching                                       │
│  └── Memory Manager (KV Cache)                                 │
├─────────────────────────────────────────────────────────────────┤
│  PagedAttention Engine                                          │
│  ├── Block Manager                                             │
│  ├── CUDA Kernels (FlashAttention)                             │
│  └── Tensor Parallelism (if multi-GPU)                         │
└─────────────────────────────────────────────────────────────────┘
```

### Running vLLM

```bash
# Single GPU inference
vllm serve meta-llama/Llama-2-7b-hf

# Multi-GPU with tensor parallelism
vllm serve meta-llama/Llama-2-70b-hf --tensor-parallel-size 4

# With custom port and GPU
vllm serve meta-llama/Llama-2-7b-hf \
  --host 0.0.0.0 \
  --port 8000 \
  --gpu-memory-utilization 0.95
```

### API Usage

```bash
# Chat completion
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d '{
    "model": "meta-llama/Llama-2-7b-hf",
    "messages": [
      {"role": "user", "content": "Explain quantum computing"}
    ],
    "temperature": 0.7,
    "max_tokens": 500
  }'

# Streaming response
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "meta-llama/Llama-2-7b-hf",
    "messages": [{"role": "user", "content": "Count to 5"}],
    "stream": true
  }'
```

## 2. Text Generation Inference (TGI)

Hugging Face's TGI is a production-ready inference server optimized for Hugging Face models.

### Key Features

- **Continuous Batching**: Maximum throughput with dynamic batching
- **FlashAttention 2**: Optimized attention computation
- **Quantization**: GPTQ, AWQ, and SqueezeLLM support
- **Kubernetes Integration**: Official Helm charts for deployment

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                   TGI Container                                 │
├─────────────────────────────────────────────────────────────────┤
│  HTTP Server (Rust)                                            │
│  ├── gRPC + REST API                                           │
│  ├── Prometheus metrics                                        │
│  └── Health checks                                            │
├─────────────────────────────────────────────────────────────────┤
│  Inference Engine                                              │
│  ├── Tokenizer (Rust)                                          │
│  ├── Model Runner (PyTorch)                                    │
│  └── Logits Processor                                          │
├─────────────────────────────────────────────────────────────────┤
│  Optimization Layer                                            │
│  ├── FlashAttention 2                                          │
│  ├── Speculative Decoding                                       │
│  └── Chunked Prefill                                          │
└─────────────────────────────────────────────────────────────────┘
```

### Running TGI

```bash
# Using Docker
docker run --gpus all \
  -p 8080:80 \
  -v $PWD/data:/data \
  ghcr.io/huggingface/text-generation-inference:1.4 \
  --model-id meta-llama/Llama-2-7b-hf

# With quantization
docker run --gpus all \
  -p 8080:80 \
  ghcr.io/huggingface/text-generation-inference:1.4 \
  --model-id meta-llama/Llama-2-7b-hf \
  --quantize bitsandbytes-nf4
```

### API Usage

```bash
# Generate endpoint
curl localhost:8080/generate \
  -X POST \
  -d '{
    "inputs": "Explain neural networks in simple terms",
    "parameters": {
      "max_new_tokens": 200,
      "temperature": 0.7
    }
  }'

# Stream generate
curl localhost:8080/generate_stream \
  -X POST \
  -d '{
    "inputs": "Write a short story",
    "parameters": {
      "max_new_tokens": 500,
      "do_sample": true,
      "temperature": 0.8
    }
  }'
```

## 3. Triton Inference Server

NVIDIA's Triton provides universal model serving with support for multiple frameworks.

### Key Features

- **Multi-Framework**: PyTorch, TensorFlow, ONNX, TensorRT
- **Dynamic Batching**: Groups requests for batched inference
- **Model Ensembles**: Pipeline multiple models
- **Concurrent Execution**: Multiple models on same GPU

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Triton Inference Server                       │
├─────────────────────────────────────────────────────────────────┤
│  HTTP/REST & gRPC Endpoints                                     │
│  ├── /v2/models/{model_name}/infer                             │
│  ├── /v2/health/ready                                          │
│  └── /v2/metrics                                               │
├─────────────────────────────────────────────────────────────────┤
│  Model Repository                                               │
│  ├── model1/ (PyTorch)                                         │
│  ├── model2/ (TensorRT)                                        │
│  └── ensemble/ (Pipeline)                                      │
├─────────────────────────────────────────────────────────────────┤
│  Backends                                                       │
│  ├── PyTorch Backend                                           │
│  ├── TensorRT Backend                                          │
│  ├── ONNX Runtime Backend                                      │
│  └── Python Backend                                            │
├─────────────────────────────────────────────────────────────────┤
│  Scheduling & Batching                                         │
│  ├── Dynamic Batching                                          │
│  ├── Sequence Batching                                        │
│  └── Priority Queuing                                         │
└─────────────────────────────────────────────────────────────────┘
```

### Configuration

```yaml
# config.pbtxt for LLM
name: llama_model
platform: pytorch_libtorch
max_batch_size: 32
input [
  {
    name: "input_ids"
    data_type: INT32
    dims: [-1, -1]
  },
  {
    name: "attention_mask"
    data_type: INT32
    dims: [-1, -1]
  }
]
output [
  {
    name: "logits"
    data_type: FP16
    dims: [-1, -1, 32000]
  }
]
parameters: {
  key: "num_heads"
  value: { string_value: "32" }
}
instance_group {
  count: 2
  kind: KIND_GPU
}
```

## 4. OpenAI-Compatible APIs

All major inference servers provide OpenAI-compatible APIs, enabling:

- Easy migration from OpenAI API
- Consistent client code across providers
- Standardized monitoring and testing

### Go Client Example

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	client := openai.NewClient(
		openai.WithBaseURL("http://localhost:8000/v1"),
		openai.WithAPIKey("dummy-key"),
	)

	ctx := context.Background()

	resp, err := client.ChatCompletions(ctx, openai.ChatCompletionRequest{
		Model: "meta-llama/Llama-2-7b-hf",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "What is the capital of France?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

### Streaming with Go

```go
func streamChat(ctx context.Context, client *openai.Client, prompt string) {
	req := openai.ChatCompletionRequest{
		Model: "meta-llama/Llama-2-7b-hf",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()

	fmt.Print("Response: ")
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		if len(resp.Choices) > 0 {
			fmt.Print(resp.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}
```

## Performance Comparison

| Feature | vLLM | TGI | Triton |
|---------|------|-----|--------|
| **PagedAttention** | Yes | No | No |
| **Continuous Batching** | Yes | Yes | Yes |
| **Quantization** | AWQ, GPTQ | GPTQ, AWQ, NF4 | TensorRT |
| **Multi-GPU** | Tensor Parallelism | Tensor Parallelism | Ensemble |
| **OpenAI API** | Yes | Partial | No |
| **Ease of Use** | High | High | Medium |

## Optimization Techniques

### 1. Batching Strategies

```
Static Batching (Traditional):
Request 1: [====]         Request 2:    [====]
Request 3:       [====]   Request 4:           [====]

Continuous Batching (vLLM, TGI):
Request 1: [====]            
Request 2:    [====]              
Request 3:       [====]
Request 4:            [====]
         ↑ Add new requests mid-generation
```

### 2. Quantization

| Method | Compression | Quality Loss | Use Case |
|--------|-------------|--------------|----------|
| FP16 | 2x | Minimal | Production default |
| INT8 | 4x | Low | Memory constrained |
| INT4 | 8x | Moderate | High volume, cost-sensitive |

### 3. Caching

```go
// Simple token caching example
type TokenCache struct {
	cache *lru.Cache
	mu    sync.RWMutex
}

func (c *TokenCache) Get(prompt string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	hash := sha256.Sum256([]byte(prompt))
	key := string(hash[:])
	
	if cached, ok := c.cache.Get(key); ok {
		return cached.(string), true
	}
	return "", false
}

func (c *TokenCache) Set(prompt, response string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	hash := sha256.Sum256([]byte(prompt))
	key := string(hash[:])
	c.cache.Add(key, response)
}
```

## Real-World Example: AI Agent Platform

The [API Gateway in this repository](services/api-gateway/README.md) demonstrates production patterns:

```go
// Simplified proxy to LLM backend
func (p *Proxy) handleChatCompletion(w http.ResponseWriter, r *http.Request) {
    // Rate limiting
    if !p.rateLimiter.Allow(r) {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    // Auth check
    apiKey := r.Header.Get("Authorization")
    if !p.auth.Validate(apiKey) {
        http.Error(w, "Invalid API key", http.StatusUnauthorized)
        return
    }

    // Forward to inference server
    backendURL := p.getBackendURL()
    proxyReq, _ := http.NewRequest(r.Method, backendURL+r.URL.Path, r.Body)
    proxyReq.Header = r.Header

    resp, err := p.httpClient.Do(proxyReq)
    // Handle response...
}
```

## Key Takeaways

- ✅ vLLM offers highest throughput with PagedAttention
- ✅ TGI provides excellent Hugging Face integration
- ✅ Triton offers multi-framework support
- ✅ OpenAI-compatible APIs enable easy migrations
- ✅ Continuous batching significantly improves throughput

## Next Steps

→ [IO-06: API Gateway Patterns](../06-api-gateway/README.md)

## Additional Resources

- [vLLM Documentation](https://docs.vllm.ai)
- [TGI GitHub](https://github.com/huggingface/text-generation-inference)
- [Triton Documentation](https://docs.nvidia.com/deeplearning/triton-inference-server/)
- [OpenAI API Compatibility](https://platform.openai.com/docs/api-reference)
