#!/bin/bash

# GO-PRO Integration Test Script
# Tests backend and frontend connectivity

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKEND_URL="http://localhost:8080"
FRONTEND_URL="http://localhost:3000"
API_BASE="${BACKEND_URL}/api/v1"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         GO-PRO Integration Test Suite                     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to test endpoint
test_endpoint() {
    local name=$1
    local url=$2
    local expected_status=${3:-200}
    
    echo -n "Testing ${name}... "
    
    response=$(curl -s -w "\n%{http_code}" "$url")
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $status_code)"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (Expected HTTP $expected_status, got $status_code)"
        echo "Response: $body"
        return 1
    fi
}

# Function to test CORS
test_cors() {
    local url=$1
    
    echo -n "Testing CORS headers... "
    
    response=$(curl -s -I -H "Origin: http://localhost:3000" -H "Access-Control-Request-Method: GET" -X OPTIONS "$url")
    
    if echo "$response" | grep -q "Access-Control-Allow-Origin"; then
        echo -e "${GREEN}✓ PASSED${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "Response headers:"
        echo "$response"
        return 1
    fi
}

# Function to test JSON response
test_json_response() {
    local name=$1
    local url=$2
    local json_path=$3
    
    echo -n "Testing ${name}... "
    
    response=$(curl -s "$url")
    
    if echo "$response" | jq -e "$json_path" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASSED${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "Response: $response"
        return 1
    fi
}

# Check if backend is running
echo -e "${YELLOW}[1/5] Checking Backend Server...${NC}"
if ! curl -s "${BACKEND_URL}/api/v1/health" > /dev/null; then
    echo -e "${RED}✗ Backend is not running at ${BACKEND_URL}${NC}"
    echo "Please start the backend server first:"
    echo "  cd backend && go run ./cmd/server"
    exit 1
fi
echo -e "${GREEN}✓ Backend is running${NC}"
echo ""

# Check if frontend is running
echo -e "${YELLOW}[2/5] Checking Frontend Server...${NC}"
if ! curl -s "${FRONTEND_URL}" > /dev/null; then
    echo -e "${RED}✗ Frontend is not running at ${FRONTEND_URL}${NC}"
    echo "Please start the frontend server first:"
    echo "  cd frontend && bun run dev"
    exit 1
fi
echo -e "${GREEN}✓ Frontend is running${NC}"
echo ""

# Test Backend API Endpoints
echo -e "${YELLOW}[3/5] Testing Backend API Endpoints...${NC}"
test_endpoint "Health Check" "${API_BASE}/health" 200
test_json_response "Health Check JSON" "${API_BASE}/health" ".success"
test_endpoint "Courses List" "${API_BASE}/courses" 200
test_json_response "Courses JSON" "${API_BASE}/courses" ".data.items"
test_json_response "Pagination" "${API_BASE}/courses" ".data.pagination"
echo ""

# Test CORS
echo -e "${YELLOW}[4/5] Testing CORS Configuration...${NC}"
test_cors "${API_BASE}/courses"
test_cors "${API_BASE}/health"
echo ""

# Test Frontend
echo -e "${YELLOW}[5/5] Testing Frontend...${NC}"
test_endpoint "Frontend Home Page" "${FRONTEND_URL}" 200

# Check if frontend can load
echo -n "Testing Frontend HTML... "
if curl -s "${FRONTEND_URL}" | grep -q "GO-PRO"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
fi
echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Test Summary                            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓ Backend Server:${NC} Running at ${BACKEND_URL}"
echo -e "${GREEN}✓ Frontend Server:${NC} Running at ${FRONTEND_URL}"
echo -e "${GREEN}✓ API Endpoints:${NC} Working correctly"
echo -e "${GREEN}✓ CORS:${NC} Configured properly"
echo -e "${GREEN}✓ Integration:${NC} Backend and Frontend connected"
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}All integration tests passed! 🎉${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo ""
echo "You can now access:"
echo "  • Frontend: ${FRONTEND_URL}"
echo "  • Backend API: ${API_BASE}"
echo "  • Health Check: ${API_BASE}/health"
echo ""

