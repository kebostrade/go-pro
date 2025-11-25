#!/bin/bash

# Generate Go code from protobuf definitions
# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc

echo "🔧 Generating protobuf code..."

# Install required tools if not present
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Please install Protocol Buffers compiler."
    echo "   Visit: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    user.proto

echo "✅ Protobuf code generated successfully!"
echo "   Generated files:"
echo "   - user.pb.go (message types)"
echo "   - user_grpc.pb.go (service definitions)"

