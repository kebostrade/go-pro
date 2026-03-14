# Quiz: Model Serving Architectures

## Question 1

What is the key innovation in vLLM that improves memory efficiency?

A) FlashAttention
B) PagedAttention
C) Tensor Parallelism
D) Continuous Batching

## Question 2

Which inference server is developed by NVIDIA?

A) vLLM
B) Text Generation Inference (TGI)
C) Triton Inference Server
D) Ray Serve

## Question 3

What type of batching does vLLM use to maximize throughput?

A) Static Batching
B) Dynamic Batching
C) Continuous Batching
D) Adaptive Batching

## Question 4

Which quantization method provides the highest compression ratio?

A) FP16
B) INT8
C) INT4
D) BF16

## Question 5

What is "Time to First Token" (TTFT)?

A) The total time to generate a complete response
B) The time from request to first token generation
C) The time between consecutive tokens
D) The average token generation speed

## Question 6

Which server provides native OpenAI-compatible API endpoints?

A) Triton only
B) vLLM only
C) Both vLLM and TGI
D) All three (vLLM, TGI, Triton)

## Question 7

What is speculative decoding?

A) Running multiple models in parallel
B) Using a smaller model to predict tokens for a larger model
C) Batch processing multiple requests
D) Caching previous responses

## Question 8

In the context of LLM serving, what does "continuous batching" mean?

A) Processing requests in fixed-size batches
B) Dynamically adding requests to batches during generation
C) Running multiple models simultaneously
D) Pre-loading model weights

## Question 9

Which file format does Triton use for model configuration?

A) JSON
B) YAML
C) Protocol Buffers (pbtxt)
D) TOML

## Question 10

What is the primary advantage of using OpenAI-compatible APIs?

A) Better performance than native APIs
B) Easier migration between different inference servers
C) Lower latency
D) More features

## Question 11

Which optimization technique is specifically used by TGI but not vLLM?

A) PagedAttention
B) FlashAttention 2
C) Speculative Decoding
D) Tensor Parallelism

## Question 12

What is the purpose of KV cache in LLM inference?

A) To store user authentication tokens
B) To cache computed attention keys and values for efficiency
C) To store model weights
D) To cache final responses
