#!/bin/bash

# GO-PRO Learning Platform - Quick Start Script
# This script helps you quickly start the learning platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print banner
print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║                                                           ║"
    echo "║           🚀 GO-PRO Learning Platform 🚀                  ║"
    echo "║                                                           ║"
    echo "║     Complete Go Programming Learning Suite                ║"
    echo "║                                                           ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Check prerequisites
check_prerequisites() {
    echo -e "${YELLOW}Checking prerequisites...${NC}"
    
    # Check Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed. Please install Go 1.21+ from https://go.dev/dl/${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Go $(go version | awk '{print $3}') found${NC}"
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}❌ Node.js is not installed. Please install Node.js 18+ from https://nodejs.org/${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Node.js $(node --version) found${NC}"
    
    # Check bun
    if ! command -v bun &> /dev/null; then
        echo -e "${RED}❌ bun is not installed${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ bun $(bun --version) found${NC}"
    
    echo ""
}

# Setup backend
setup_backend() {
    echo -e "${YELLOW}Setting up backend...${NC}"
    cd backend
    
    if [ ! -d "bin" ]; then
        mkdir -p bin
    fi
    
    echo -e "${BLUE}Installing Go dependencies...${NC}"
    go mod download
    go mod tidy
    
    echo -e "${GREEN}✅ Backend setup complete${NC}"
    cd ..
    echo ""
}

# Setup frontend
setup_frontend() {
    echo -e "${YELLOW}Setting up frontend...${NC}"
    cd frontend
    
    if [ ! -d "node_modules" ]; then
        echo -e "${BLUE}Installing bun dependencies...${NC}"
        bun install
    else
        echo -e "${BLUE}Dependencies already installed${NC}"
    fi
    
    echo -e "${GREEN}✅ Frontend setup complete${NC}"
    cd ..
    echo ""
}

# Start backend
start_backend() {
    echo -e "${YELLOW}Starting backend server...${NC}"
    cd backend
    echo -e "${CYAN}Backend API will be available at: ${GREEN}http://localhost:8080${NC}"
    go run ./cmd/server &
    BACKEND_PID=$!
    cd ..
    sleep 2
    echo -e "${GREEN}✅ Backend started (PID: $BACKEND_PID)${NC}"
    echo ""
}

# Start frontend
start_frontend() {
    echo -e "${YELLOW}Starting frontend server...${NC}"
    cd frontend
    echo -e "${CYAN}Frontend dashboard will be available at: ${GREEN}http://localhost:3000${NC}"
    bun run dev &
    FRONTEND_PID=$!
    cd ..
    sleep 2
    echo -e "${GREEN}✅ Frontend started (PID: $FRONTEND_PID)${NC}"
    echo ""
}

# Print success message
print_success() {
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║                                                           ║"
    echo "║              🎉 Platform Started Successfully! 🎉         ║"
    echo "║                                                           ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo ""
    echo -e "${CYAN}📚 Access your learning platform:${NC}"
    echo -e "   ${GREEN}Frontend Dashboard:${NC} http://localhost:3000"
    echo -e "   ${GREEN}Backend API:${NC}        http://localhost:8080"
    echo -e "   ${GREEN}API Health Check:${NC}   http://localhost:8080/api/v1/health"
    echo ""
    echo -e "${CYAN}📖 Quick Start Guide:${NC}"
    echo -e "   1. Open ${GREEN}http://localhost:3000${NC} in your browser"
    echo -e "   2. Start with Lesson 1 in the ${GREEN}course/lessons/${NC} directory"
    echo -e "   3. Try exercises in ${GREEN}course/code/${NC}"
    echo -e "   4. Run tests with ${GREEN}go test${NC}"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
    echo ""
}

# Cleanup on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Shutting down services...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo -e "${GREEN}✅ Backend stopped${NC}"
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo -e "${GREEN}✅ Frontend stopped${NC}"
    fi
    echo -e "${GREEN}Goodbye! Happy coding! 🚀${NC}"
    exit 0
}

# Main execution
main() {
    print_banner
    check_prerequisites
    
    # Parse arguments
    SETUP_ONLY=false
    if [ "$1" == "--setup" ]; then
        SETUP_ONLY=true
    fi
    
    setup_backend
    setup_frontend
    
    if [ "$SETUP_ONLY" = true ]; then
        echo -e "${GREEN}Setup complete! Run './start.sh' to start the platform.${NC}"
        exit 0
    fi
    
    # Set up cleanup trap
    trap cleanup SIGINT SIGTERM
    
    start_backend
    start_frontend
    print_success
    
    # Keep script running
    while true; do
        sleep 1
    done
}

# Run main function
main "$@"

