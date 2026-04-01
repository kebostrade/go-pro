# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o publisher ./cmd/publisher && \
    CGO_ENABLED=0 GOOS=linux go build -o subscriber ./cmd/subscriber && \
    CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Package stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Copy binaries from builder
COPY --from=builder /app/publisher /publisher
COPY --from=builder /app/subscriber /subscriber
COPY --from=builder /app/worker /worker

# Default command
CMD ["/publisher"]
