# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (test binary for coverage)
RUN CGO_ENABLED=0 GOOS=linux go test -c -o tests ./...

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy test binary from builder
COPY --from=builder /app/tests /app/

# Set ownership
RUN chown -R appuser:appuser /app

USER appuser

CMD ["/app/tests"]
