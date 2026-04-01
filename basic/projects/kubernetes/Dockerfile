# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the operator
RUN CGO_ENABLED=0 GOOS=linux go build -o operator ./cmd

# Package stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S operator && adduser -S operator -G operator

WORKDIR /app

# Copy operator binary from builder
COPY --from=builder /app/operator /operator

# Copy manifests directory
COPY --from=builder /app/manifests /manifests

# Change ownership
RUN chown -R operator:operator /app

# Switch to non-root user
USER operator

# Expose ports
EXPOSE 8080 8081

# Set entrypoint
ENTRYPOINT ["/operator"]
