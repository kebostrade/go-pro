# Multi-stage Docker build for Go application
# Stage 1: Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies and security tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx \
    gcc \
    musl-dev \
    && update-ca-certificates

# Create app directory
WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Run go mod tidy to fix go.mod version format for Go 1.25
RUN go mod tidy

# Build arguments for metadata
ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Run security checks and tests
RUN go vet ./...
# Skip tests in Docker build for faster development - tests run via `go test` by developers
# RUN CGO_ENABLED=1 go test -race -short ./

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commit=${VCS_REF} -extldflags '-static'" \
    -a -installsuffix cgo \
    -tags netgo,osusergo \
    -trimpath \
    -o app ./cmd/server

# Compress binary with UPX (optional, can be disabled for debugging)
RUN upx --best --lzma app || true

# Stage 2: Runtime stage (distroless for maximum security)
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# Build arguments for metadata (need to be redeclared in each stage)
ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF

# Copy built application from builder stage
COPY --from=builder /app/app /app

# Copy timezone data and CA certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set working directory
WORKDIR /

# Use non-root user (already set in distroless image)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV TZ=UTC

# Alternative runtime stage with Alpine (uncomment if distroless causes issues)
# FROM alpine:3.19 AS runtime-alpine
#
# # Install runtime dependencies and security updates
# RUN apk --no-cache add \
#     ca-certificates \
#     tzdata \
#     wget \
#     && apk upgrade --available \
#     && addgroup -g 1001 appgroup \
#     && adduser -u 1001 -G appgroup -s /bin/sh -D appuser
#
# # Set working directory
# WORKDIR /app
#
# # Copy built application from builder stage
# COPY --from=builder /app/app .
#
# # Set ownership and permissions
# RUN chown -R appuser:appgroup /app && \
#     chmod +x app
#
# # Switch to non-root user
# USER appuser
#
# # Health check
# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
#     CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1
#
# # Expose port
# EXPOSE 8080
#
# # Environment variables
# ENV PORT=8080
# ENV TZ=UTC

# Metadata labels
LABEL maintainer="GO-PRO Team"
LABEL version="${VERSION}"
LABEL description="GO-PRO Learning Platform Backend API"
LABEL org.opencontainers.image.title="GO-PRO Backend"
LABEL org.opencontainers.image.description="RESTful API for Go programming learning platform"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.revision="${VCS_REF}"
LABEL org.opencontainers.image.source="https://github.com/your-org/go-pro"

# Run the application
ENTRYPOINT ["./app"]