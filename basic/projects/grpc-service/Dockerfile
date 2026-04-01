FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install protoc and plugins
RUN apk add --no-cache protobuf git

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source and proto
COPY proto/ ./proto/
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Generate protobuf code
RUN make generate

# Build server
RUN go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates protobuf

WORKDIR /app

COPY --from=builder /app/server /server

EXPOSE 50051

ENV GRPC_PORT=50051

CMD ["/server"]
