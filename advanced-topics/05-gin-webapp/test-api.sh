#!/bin/bash

# Test Script for Gin Web Framework Examples

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     🧪 Gin Web Framework - API Test Suite                ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $2"
    else
        echo -e "${RED}❌ FAIL${NC}: $2"
    fi
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5

    echo -e "\n${BLUE}Testing:${NC} $description"
    echo "   $method $endpoint"

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi

    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    echo "   Status: $status_code"
    echo "   Response: $body"

    if [ "$status_code" = "$expected_status" ]; then
        print_result 0 "$description"
        return 0
    else
        print_result 1 "$description (expected $expected_status, got $status_code)"
        return 1
    fi
}

# Check if server is running
echo -e "${BLUE}Checking if server is running...${NC}"
if ! curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}❌ Server is not running. Please start the server first:${NC}"
    echo "   go run main.go"
    echo "   OR"
    echo "   ./quick-start.sh"
    exit 1
fi

echo -e "${GREEN}✅ Server is running${NC}"
echo ""

# Test counter
total_tests=0
passed_tests=0

# ═════════════════════════════════════════════════════════════
# Health Check Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ Health Check Tests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "GET" "/health" "" "200" "Health check endpoint"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# GET Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ GET Requests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "GET" "/api/users" "" "200" "Get all users"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "GET" "/api/users/1" "" "200" "Get user by ID"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "GET" "/api/users/999" "" "404" "Get non-existent user"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# POST Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ POST Requests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "POST" "/api/users" \
    '{"name":"Test User","email":"test@example.com","age":25}' \
    "401" \
    "Create user without authentication (should fail)"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "POST" "/api/users" \
    '{"name":"Test User","email":"test@example.com","age":25}' \
    "201" \
    "Create user with authentication" \
    -H "Authorization: Bearer valid-token"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "POST" "/api/users" \
    '{"name":"A","email":"invalid-email","age":25}' \
    "400" \
    "Create user with invalid data (should fail)"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# PUT Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ PUT Requests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "PUT" "/api/users/1" \
    '{"name":"Updated User","email":"updated@example.com","age":35}' \
    "401" \
    "Update user without authentication (should fail)"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "PUT" "/api/users/1" \
    '{"name":"Updated User","email":"updated@example.com","age":35}' \
    "200" \
    "Update user with authentication"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# DELETE Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ DELETE Requests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "DELETE" "/api/users/2" "" "401" \
    "Delete user without authentication (should fail)"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

total_tests=$((total_tests + 1))
test_endpoint "DELETE" "/api/users/2" "" "200" \
    "Delete user with authentication"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# 404 Tests
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ 404 Not Found Tests ═══${NC}"

total_tests=$((total_tests + 1))
test_endpoint "GET" "/nonexistent" "" "404" \
    "Access non-existent route"
if [ $? -eq 0 ]; then
    passed_tests=$((passed_tests + 1))
fi

# ═════════════════════════════════════════════════════════════
# Test Results Summary
# ═════════════════════════════════════════════════════════════
echo -e "\n${BLUE}═══ Test Results Summary ═══${NC}"
echo ""
echo "Total Tests:  $total_tests"
echo -e "Passed:       ${GREEN}$passed_tests${NC}"
echo -e "Failed:       ${RED}$((total_tests - passed_tests))${NC}"
echo ""

if [ $passed_tests -eq $total_tests ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Some tests failed${NC}"
    exit 1
fi
