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

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/server /app/

# Copy views and static files
COPY --from=builder /app/internal/views /app/internal/views
COPY --from=builder /app/static /app/static

# Set ownership
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/server"]
