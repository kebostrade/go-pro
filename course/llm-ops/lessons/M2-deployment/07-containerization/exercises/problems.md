# Exercise: Containerization with Docker

## Problem 1: Docker Image Layers

Given the following Dockerfile, identify the layers and explain the build caching:

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python3", "app.py"]
```

1. What happens when you modify `requirements.txt`?
2. What happens when you modify `app.py`?
3. How would you reorder to optimize caching?

---

## Problem 2: Multi-Stage Build

Write a multi-stage Dockerfile for a Python LLM inference server:

```dockerfile
# Stage 1: Builder
# ...

# Stage 2: Runtime
# ...
```

Requirements:
- Use Python 3.11
- Install dependencies from requirements.txt
- Use slim base image in final stage
- Include security best practices

---

## Problem 3: GPU Container Configuration

Create a docker-compose.yml that:

1. Runs vLLM with GPU access
2. Runs API Gateway (from services/api-gateway)
3. Runs Redis for caching
4. Includes proper networking
5. Adds health checks

```yaml

```

---

## Problem 4: Optimize Dockerfile

Optimize the following Dockerfile for build time and image size:

```dockerfile
FROM python:3.11

WORKDIR /app

# Install system dependencies
RUN apt-get update
RUN apt-get install -y build-essential git wget curl

# Install Python dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy all files
COPY . .

# Install more packages
RUN apt-get install -y vim nano

CMD ["python3", "server.py"]
```

### Issues identified:
1.
2.
3.

### Optimized version:
```dockerfile

```

---

## Problem 5: API Gateway Container

Build a Docker container for the API Gateway in `services/api-gateway/`:

1. Write a Dockerfile
2. Write a docker-compose service definition
3. Test the build

### Dockerfile:

```dockerfile

```

### docker-compose:

```yaml

```

---

## Problem 6: Container Security

Add security configurations to the following Docker Compose:

```yaml
version: '3.8'

services:
  api:
    image: my-api:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/mydb
```

Add:
- Non-root user
- Read-only filesystem
- Resource limits
- Remove secrets from environment variables

---

## Problem 7: Health Check Implementation

Add a proper health check to this service:

```yaml
services:
  vllm:
    image: vllm/vllm-openai:latest
    ports:
      - "8000:8000"
```

---

## Submission

Save your answers and be prepared to discuss them in the next lesson.
