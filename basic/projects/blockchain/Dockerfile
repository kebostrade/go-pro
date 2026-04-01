# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the wallet tool
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet ./cmd/wallet

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/wallet .

# Copy ABI files
COPY --from=builder /app/abi ./abi

# Run the wallet
CMD ["./wallet"]
