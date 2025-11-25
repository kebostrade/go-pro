#!/bin/bash

# Build script for coding sandbox Docker images
# This script builds all language-specific sandbox images

set -e

# Configuration
IMAGE_PREFIX="coding-sandbox"
DOCKER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🐳 Building Coding Sandbox Docker Images"
echo "=========================================="
echo ""

# Function to build an image
build_image() {
    local language=$1
    local dockerfile=$2
    local image_name="${IMAGE_PREFIX}-${language}:latest"
    
    echo "📦 Building ${image_name}..."
    
    if docker build -t "${image_name}" -f "${dockerfile}" "${DOCKER_DIR}"; then
        echo "✅ Successfully built ${image_name}"
    else
        echo "❌ Failed to build ${image_name}"
        return 1
    fi
    
    echo ""
}

# Build Go sandbox
build_image "go" "${DOCKER_DIR}/Dockerfile.go"

# Build Python sandbox
build_image "python" "${DOCKER_DIR}/Dockerfile.python"

# Build Node.js sandbox (for JavaScript and TypeScript)
build_image "node" "${DOCKER_DIR}/Dockerfile.node"

echo "🎉 All sandbox images built successfully!"
echo ""
echo "Available images:"
docker images | grep "${IMAGE_PREFIX}"
echo ""
echo "To test an image, run:"
echo "  docker run --rm ${IMAGE_PREFIX}-go:latest go version"
echo "  docker run --rm ${IMAGE_PREFIX}-python:latest python3 --version"
echo "  docker run --rm ${IMAGE_PREFIX}-node:latest node --version"

