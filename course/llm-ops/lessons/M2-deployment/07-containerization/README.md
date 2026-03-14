# IO-07: Containerization with Docker

**Duration**: 3 hours
**Module**: 2 - Deployment & Serving

## Learning Objectives

- Containerize LLM applications using Docker
- Implement multi-stage builds for optimized images
- Configure GPU containers for inference
- Optimize Docker images for production deployment

## Why Containers for LLMs?

Containers provide consistent environments across development and production:

- **Reproducibility**: Same environment everywhere
- **Isolation**: Dependencies don't conflict
- **Portability**: Works locally, in cloud, on-prem
- **Resource Efficiency**: Better utilization than VMs

```
┌─────────────────────────────────────────────────────────────────┐
│                      Host System                                 │
│                    (Linux, NVIDIA Driver)                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   Docker Runtime                             ││
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐  ││
│  │  │   Container   │  │   Container   │  │   Container   │  ││
│  │  │   (vLLM)      │  │   (TGI)       │  │   (API GW)    │  ││
│  │  └───────────────┘  └───────────────┘  └───────────────┘  ││
│  │                         NVIDIA Container Runtime            ││
│  └─────────────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────────┤
│                    GPU (NVIDIA A100/H100)                       │
└─────────────────────────────────────────────────────────────────┘
```

## 1. Docker Basics for LLMs

### Dockerfile Structure

```dockerfile
# Base image with GPU support
FROM nvidia/cuda:12.1.0-runtime-ubuntu22.04

# Install Python and dependencies
RUN apt-get update && apt-get install -y \
    python3.11 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Expose port
EXPOSE 8000

# Run the application
CMD ["python3", "server.py"]
```

## 2. Multi-Stage Builds

Multi-stage builds reduce image size significantly:

```dockerfile
# Stage 1: Build stage
FROM python:3.11-slim AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --user -r requirements.txt

# Stage 2: Runtime stage
FROM nvidia/cuda:12.1.0-runtime-ubuntu22.04

WORKDIR /app

# Install runtime dependencies only
RUN apt-get update && apt-get install -y \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Copy built artifacts from builder
COPY --from=builder /root/.local /root/.local

# Copy application
COPY . .

# Set PATH
ENV PATH=/root/.local/bin:$PATH

EXPOSE 8000

CMD ["python3", "server.py"]
```

## 3. GPU Container Configuration

### Using NVIDIA Runtime

```dockerfile
# Use NVIDIA base image
FROM nvidia/cuda:12.1.0-base-ubuntu22.04

# Install vLLM
RUN pip install vllm

CMD ["vllm", "serve", "meta-llama/Llama-2-7b-hf"]
```

### Running with GPU

```bash
# Pull NVIDIA runtime image
docker pull nvidia/cuda:12.1.0-base-ubuntu22.04

# Run with GPU access
docker run --gpus all \
  --runtime nvidia \
  -p 8000:8000 \
  vllm-image

# Specify specific GPU
docker run --gpus '"device=0"' \
  -p 8000:8000 \
  vllm-image
```

### GPU Memory Management

```dockerfile
# Set environment variables for GPU memory
ENV CUDA_VISIBLE_DEVICES=0
ENV PYTORCH_CUDA_ALLOC_CONF=max_split_size_mb:128
```

## 4. vLLM Docker Image

```dockerfile
FROM nvidia/cuda:12.1.0-devel-ubuntu22.04 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    python3.11 \
    python3-pip \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install vLLM from source (for latest features)
RUN git clone https://github.com/vllm-project/vllm.git && \
    cd vllm && \
    pip install -e .

# Production runtime stage
FROM nvidia/cuda:12.1.0-runtime-ubuntu22.04

WORKDIR /app

RUN apt-get update && apt-get install -y \
    python3.11 \
    python3-pip \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/lib/python3.11/dist-packages /usr/local/lib/python3.11/dist-packages

# Create non-root user
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8000

ENV VLLM_WORKER_MULTIPROC_METHOD=spawn
ENV VLLM_LOGGING_LEVEL=INFO

CMD ["python3", "-m", "vllm.entrypoints.openai.api_server", \
     "--model", "meta-llama/Llama-2-7b-hf"]
```

## 5. API Gateway Container

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Build application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /api-gateway .

EXPOSE 8080

ENV PORT=8080
ENV BACKEND_URL=http://vllm:8000

ENTRYPOINT ["/api-gateway"]
```

## 6. Docker Compose for LLM Stack

```yaml
version: '3.8'

services:
  vllm:
    build:
      context: ./vllm
      dockerfile: Dockerfile
    image: vllm:latest
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
    ports:
      - "8000:8000"
    environment:
      - MODEL=meta-llama/Llama-2-7b-hf
      - GPU_MEMORY_UTILIZATION=0.9
    volumes:
      - model-cache:/root/.cache/huggingface
    shm_size: '10gb'
    restart: unless-stopped

  api-gateway:
    build:
      context: ./api-gateway
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - BACKEND_URL=http://vllm:8000
      - JWT_SECRET=${JWT_SECRET}
      - RATE_LIMIT=100
    depends_on:
      - vllm
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped

volumes:
  model-cache:
  redis-data:
```

## 7. Optimization Techniques

### Image Size Reduction

```dockerfile
# Use minimal base images
FROM python:3.11-slim

# Combine RUN commands to reduce layers
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
    && rm -rf /var/lib/apt/lists/*

# Use --no-cache-dir for pip
RUN pip install --no-cache-dir package1 package2

# Multi-stage build (as shown earlier)
```

### Layer Caching

```dockerfile
# Order from least to most frequently changing
# 1. Base image
FROM python:3.11-slim

# 2. System dependencies (rarely change)
RUN apt-get update && apt-get install -y ...

# 3. Python dependencies (change occasionally)
COPY requirements.txt .
RUN pip install -r requirements.txt

# 4. Application code (changes frequently)
COPY . .
```

### Security Best Practices

```dockerfile
# Use specific versions, not :latest
FROM python:3.11-slim@sha256:abc123...

# Create non-root user
RUN useradd -m -u 1000 appuser
USER appuser

# Use read-only file system where possible
# (requires careful volume configuration)

# Don't expose secrets in image
# Use environment variables or secrets at runtime
```

## 8. Go Application Docker Example

The [API Gateway in this repository](services/api-gateway/) demonstrates Go Docker best practices:

```dockerfile
# Multi-stage build for Go API Gateway
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o api-gateway ./cmd

# Minimal runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/api-gateway"]
```

Build and run:

```bash
cd services/api-gateway
docker build -t api-gateway:latest .
docker run -p 8080:8080 -e BACKEND_URL=http://localhost:8000 api-gateway:latest
```

## 9. Health Checks and Monitoring

```yaml
services:
  vllm:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s

  api-gateway:
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Key Takeaways

- ✅ Multi-stage builds reduce image size significantly
- ✅ NVIDIA runtime enables GPU access in containers
- ✅ Layer ordering optimizes build caching
- ✅ Non-root users improve security

## Next Steps

→ [IO-08: Kubernetes for AI](../08-kubernetes/README.md)

## Additional Resources

- [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/)
- [Docker Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [vLLM Docker](https://docs.vllm.ai/en/latest/serving/docker.html)
