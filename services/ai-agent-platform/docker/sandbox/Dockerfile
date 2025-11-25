# Dockerfile for Go code execution sandbox
FROM golang:1.22-alpine AS builder

# Install security tools
RUN apk add --no-cache \
    ca-certificates \
    git

# Create non-root user
RUN addgroup -g 1000 sandbox && \
    adduser -D -u 1000 -G sandbox sandbox

# Set working directory
WORKDIR /workspace

# Switch to non-root user
USER sandbox

# Default command
CMD ["go", "version"]

