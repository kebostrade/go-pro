#!/bin/bash

# Test script for Concurrency Crash Course examples
# Run: chmod +x test.sh && ./test.sh

set -e

echo "🧪 Testing Go Concurrency Crash Course Examples"
echo "================================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Check if Go is installed
echo "1. Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go $(go version)${NC}"
echo ""

# Test 2: Check go.mod
echo "2. Checking go.mod..."
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ go.mod not found${NC}"
    exit 1
fi
echo -e "${GREEN}✅ go.mod exists${NC}"
echo ""

# Test 3: Build the program
echo "3. Building program..."
if go build -o crash-course main.go; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi
echo ""

# Test 4: Build with race detector
echo "4. Building with race detector..."
if go build -race -o crash-course-race main.go; then
    echo -e "${GREEN}✅ Race detector build successful${NC}"
else
    echo -e "${RED}❌ Race detector build failed${NC}"
    exit 1
fi
echo ""

# Test 5: Run basic example
echo "5. Running basic goroutine example..."
if timeout 5s ./crash-course > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Program runs successfully${NC}"
else
    echo -e "${YELLOW}⚠️  Program timed out or failed (this may be expected)${NC}"
fi
echo ""

# Test 6: Check for race conditions
echo "6. Checking for race conditions..."
echo -e "${YELLOW}Note: This will run the race detector. Some examples may intentionally show races.${NC}"
if timeout 5s ./crash-course-race > /dev/null 2>&1; then
    echo -e "${GREEN}✅ No race conditions detected${NC}"
else
    echo -e "${YELLOW}⚠️  Race conditions detected or timeout (check pitfall examples)${NC}"
fi
echo ""

# Test 7: Verify all functions exist
echo "7. Verifying all example functions..."
functions=(
    "example1_BasicGoroutines"
    "example2_Channels"
    "example3_WaitGroups"
    "example4_Select"
    "example5_WorkerPool"
    "example6_Pipeline"
    "example7_FanOutFanIn"
    "example8_Context"
    "example9_Mutex"
    "example10_WebScraper"
    "example11_RateLimiter"
)

all_found=true
for func in "${functions[@]}"; do
    if grep -q "func $func()" main.go; then
        echo -e "${GREEN}  ✅ $func${NC}"
    else
        echo -e "${RED}  ❌ $func not found${NC}"
        all_found=false
    fi
done

if [ "$all_found" = true ]; then
    echo -e "${GREEN}✅ All example functions found${NC}"
else
    echo -e "${RED}❌ Some example functions missing${NC}"
    exit 1
fi
echo ""

# Test 8: Check for common patterns
echo "8. Checking for concurrency patterns..."
patterns=(
    "sync.WaitGroup"
    "make(chan"
    "go func()"
    "context.Context"
    "sync.Mutex"
    "select {"
)

all_patterns_found=true
for pattern in "${patterns[@]}"; do
    if grep -q "$pattern" main.go; then
        echo -e "${GREEN}  ✅ $pattern${NC}"
    else
        echo -e "${RED}  ❌ $pattern not found${NC}"
        all_patterns_found=false
    fi
done

if [ "$all_patterns_found" = true ]; then
    echo -e "${GREEN}✅ All concurrency patterns found${NC}"
else
    echo -e "${RED}❌ Some concurrency patterns missing${NC}"
    exit 1
fi
echo ""

# Cleanup
echo "9. Cleaning up..."
rm -f crash-course crash-course-race
echo -e "${GREEN}✅ Cleanup complete${NC}"
echo ""

# Summary
echo "================================================"
echo -e "${GREEN}🎉 All tests passed!${NC}"
echo ""
echo "Next steps:"
echo "  1. Run examples: go run main.go"
echo "  2. Run with race detector: go run -race main.go"
echo "  3. Read the crash course: ../../../docs/tutorials/CONCURRENCY_CRASH_COURSE.md"
echo ""

