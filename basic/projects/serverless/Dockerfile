FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the handler
RUN CGO_ENABLED=0 GOOS=linux go build -o handler ./cmd/handler

# Final stage
FROM alpine:3.19

WORKDIR /var/task

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /build/handler /var/task/handler

# Set environment
ENV APP_VERSION=1.0.0

# Set the entrypoint
ENTRYPOINT ["/var/task/handler"]
