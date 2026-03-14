# Quiz: GPU Compute Fundamentals

## Question 1

Why are GPUs more effective than CPUs for deep learning?

A) GPUs have more cache memory
B) GPUs are optimized for parallel matrix operations
C) GPUs use less power
D) GPUs have faster single-threaded performance

## Question 2

What is VRAM?

A) Virtual RAM that extends system memory
B) Video RAM on the GPU used for model storage
C) A type of CPU memory
D) Network-attached storage

## Question 3

In data parallelism, each GPU:

A) Stores a portion of the model weights
B) Has a complete copy of the model
C) Processes different layers of the model
D) Handles different types of computations

## Question 4

What is tensor parallelism used for?

A) Splitting training data across GPUs
B) Splitting model weights across GPUs within a single layer
C) Reducing communication between GPUs
D) Increasing batch size

## Question 5

Which NVIDIA GPU is designed specifically for data center AI workloads?

A) RTX 4090
B) GTX 1080
C) A100
D) All of the above

## Question 6

What is the KV cache in LLM inference?

A) A cache for storing model weights
B) A cache for key and value tensors in attention
C) A CPU cache for tokenization
D) A disk-based cache for responses

## Question 7

Which parallelism strategy has the highest communication overhead?

A) Data parallelism
B) Pipeline parallelism
C) Tensor parallelism
D) None, they all have similar overhead

## Question 8

Spot/preemptible instances typically offer what discount compared to on-demand?

A) 10-20%
B) 30-40%
C) 60-80%
D) 90-95%

## Question 9

What is AllReduce in distributed training?

A) A reduction operation that combines gradients from all GPUs
B) A method to reduce model size
C) A memory allocation technique
D) A type of GPU kernel

## Question 10

Which GPU would be most appropriate for fine-tuning a 70B parameter model?

A) Single RTX 4090 (24GB)
B) Single A100 (80GB)
C) 8x H100 (80GB each)
D) CPU-only

## Question 11

What does "GPU utilization" measure?

A) The percentage of time the GPU is running
B) The percentage of VRAM being used
C) The temperature of the GPU
D) The power consumption

## Question 12

What is the primary bottleneck when running inference on large models?

A) CPU processing
B) Network bandwidth
C) VRAM capacity
D) Disk I/O

## Question 13

In tensor parallelism, when you split a linear layer column-wise:

A) Each GPU computes partial dot products
B) Each GPU needs the full input
C) Communication is minimized
D) The model must be retrained

## Question 14

What is pipeline bubble in pipeline parallelism?

A) A memory buffer overflow
B) Idle time when some GPUs wait for others
C) A type of GPU error
D) The initial warm-up phase
