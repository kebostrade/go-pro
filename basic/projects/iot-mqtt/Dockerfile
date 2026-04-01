# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o device ./cmd/device

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Copy binaries from builder
COPY --from=builder /app/gateway ./gateway
COPY --from=builder /app/device ./device

# Copy broker config
COPY broker_config ./broker_config

# Run gateway by default
CMD ["./gateway"]
