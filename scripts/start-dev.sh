#!/bin/bash

# GO-PRO Development Environment Startup Script
# Starts both backend and frontend servers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         GO-PRO Development Environment                     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if .env files exist
echo -e "${YELLOW}Checking environment configuration...${NC}"

if [ ! -f "backend/.env" ]; then
    echo -e "${YELLOW}⚠ backend/.env not found, creating from .env.example...${NC}"
    if [ -f "backend/.env.example" ]; then
        cp backend/.env.example backend/.env
        echo -e "${GREEN}✓ Created backend/.env${NC}"
    else
        echo -e "${RED}✗ backend/.env.example not found${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ backend/.env exists${NC}"
fi

if [ ! -f "frontend/.env.local" ]; then
    echo -e "${YELLOW}⚠ frontend/.env.local not found, creating...${NC}"
    cat > frontend/.env.local << EOF
# GO-PRO Frontend Configuration
# Development Environment

# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Environment
NODE_ENV=development
NEXT_PUBLIC_ENV=development

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_MONITORING=false
EOF
    echo -e "${GREEN}✓ Created frontend/.env.local${NC}"
else
    echo -e "${GREEN}✓ frontend/.env.local exists${NC}"
fi

echo ""

# Function to cleanup on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Shutting down servers...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo -e "${GREEN}✓ Backend stopped${NC}"
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo -e "${GREEN}✓ Frontend stopped${NC}"
    fi
    exit 0
}

trap cleanup SIGINT SIGTERM

# Build and start backend
echo -e "${YELLOW}Building backend...${NC}"
cd backend
go build -o bin/server ./cmd/server
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Backend build failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Backend built successfully${NC}"

echo -e "${YELLOW}Starting backend server...${NC}"
./bin/server > ../logs/backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# Wait for backend to start
sleep 2

# Check if backend is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo -e "${RED}✗ Backend failed to start${NC}"
    echo "Check logs/backend.log for details"
    exit 1
fi
echo -e "${GREEN}✓ Backend started (PID: $BACKEND_PID)${NC}"
echo ""

# Install frontend dependencies if needed
echo -e "${YELLOW}Checking frontend dependencies...${NC}"
cd frontend
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    npm install
fi
echo -e "${GREEN}✓ Frontend dependencies ready${NC}"

# Start frontend
echo -e "${YELLOW}Starting frontend server...${NC}"
npm run dev > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Wait for frontend to start
sleep 5

# Check if frontend is running
if ! curl -s http://localhost:3000 > /dev/null; then
    echo -e "${RED}✗ Frontend failed to start${NC}"
    echo "Check logs/frontend.log for details"
    cleanup
    exit 1
fi
echo -e "${GREEN}✓ Frontend started (PID: $FRONTEND_PID)${NC}"
echo ""

# Display status
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                 Servers Running                            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓ Backend:${NC}  http://localhost:8080"
echo -e "${GREEN}✓ Frontend:${NC} http://localhost:3000"
echo -e "${GREEN}✓ API:${NC}      http://localhost:8080/api/v1"
echo -e "${GREEN}✓ Health:${NC}   http://localhost:8080/api/v1/health"
echo ""
echo -e "${YELLOW}Logs:${NC}"
echo -e "  • Backend:  logs/backend.log"
echo -e "  • Frontend: logs/frontend.log"
echo ""
echo -e "${BLUE}Press Ctrl+C to stop all servers${NC}"
echo ""

# Keep script running
wait

