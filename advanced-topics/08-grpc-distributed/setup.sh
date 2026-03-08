#!/bin/bash

# gRPC Distributed System Examples - Setup Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}gRPC Distributed System Examples${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Check if protoc is installed
echo -e "${YELLOW}Checking for protoc...${NC}"
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}❌ protoc is not installed${NC}"
    echo ""
    echo "Please install Protocol Buffer compiler:"
    echo "  macOS:   brew install protobuf"
    echo "  Ubuntu:  sudo apt-get install protobuf-compiler"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

PROTOC_VERSION=$(protoc --version | awk '{print $2}')
echo -e "${GREEN}✓ protoc version: $PROTOC_VERSION${NC}"
echo ""

# Install Go plugins
echo -e "${YELLOW}Installing protoc Go plugins...${NC}"
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo -e "${GREEN}✓ Plugins installed${NC}"
echo ""

# Generate Go code from .proto files
echo -e "${YELLOW}Generating Go code from .proto files...${NC}"
cd proto
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
echo -e "${GREEN}✓ Generated Go code${NC}"
cd ..
echo ""

# Initialize Go modules
echo -e "${YELLOW}Initializing Go modules...${NC}"
go mod tidy
cd server && go mod tidy && cd ..
cd client && go mod tidy && cd ..
cd examples && go mod tidy && cd ..
echo -e "${GREEN}✓ Go modules initialized${NC}"
echo ""

# Create certs directory
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p certs
echo -e "${GREEN}✓ Directories created${NC}"
echo ""

# Success message
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}Setup completed successfully!${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo ""
echo "1. Start the server:"
echo -e "   ${YELLOW}cd server && go run main.go${NC}"
echo ""
echo "2. In another terminal, run the client:"
echo -e "   ${YELLOW}cd client && go run main.go${NC}"
echo ""
echo "3. Or run individual examples:"
echo -e "   ${YELLOW}cd examples${NC}"
echo -e "   ${YELLOW}go run 01-unary-rpc.go${NC}"
echo ""
echo "4. Generate TLS certificates (optional):"
echo -e "   ${YELLOW}make certs${NC}"
echo ""
echo "5. Use grpcurl to inspect services:"
echo -e "   ${YELLOW}make grpcurl${NC}"
echo -e "   ${YELLOW}make reflect${NC}"
echo ""
