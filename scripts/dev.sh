#!/bin/bash

# GO-PRO Development Script
# Comprehensive development automation for the Go learning platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Emojis for better UX
ROCKET="🚀"
CHECK="✅"
CROSS="❌"
WARNING="⚠️"
INFO="ℹ️"
GEAR="⚙️"
BOOK="📚"
TEST="🧪"
BUILD="🔨"
CLEAN="🧹"

# Project directories
BACKEND_DIR="backend"
FRONTEND_DIR="frontend"
COURSE_DIR="course"

# Function to print colored output
print_status() {
    echo -e "${2}${1}${NC}"
}

print_header() {
    echo -e "\n${PURPLE}================================${NC}"
    echo -e "${PURPLE}${1}${NC}"
    echo -e "${PURPLE}================================${NC}\n"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_header "${GEAR} Checking Prerequisites"
    
    local missing_deps=()
    
    # Check Go
    if command_exists go; then
        GO_VERSION=$(go version | cut -d' ' -f3)
        print_status "${CHECK} Go installed: $GO_VERSION" $GREEN
    else
        print_status "${CROSS} Go not found" $RED
        missing_deps+=("go")
    fi
    
    # Check Node.js
    if command_exists node; then
        NODE_VERSION=$(node --version)
        print_status "${CHECK} Node.js installed: $NODE_VERSION" $GREEN
    else
        print_status "${CROSS} Node.js not found" $RED
        missing_deps+=("node")
    fi
    
    # Check bun
    if command_exists bun; then
        BUN_VERSION=$(bun --version)
        print_status "${CHECK} bun installed: $BUN_VERSION" $GREEN
    else
        print_status "${CROSS} bun not found" $RED
        missing_deps+=("bun")
    fi
    
    # Check optional tools
    if command_exists make; then
        print_status "${CHECK} Make available" $GREEN
    else
        print_status "${WARNING} Make not found (optional)" $YELLOW
    fi
    
    if command_exists curl; then
        print_status "${CHECK} curl available" $GREEN
    else
        print_status "${WARNING} curl not found (optional)" $YELLOW
    fi
    
    if command_exists jq; then
        print_status "${CHECK} jq available" $GREEN
    else
        print_status "${WARNING} jq not found (optional)" $YELLOW
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_status "${CROSS} Missing required dependencies: ${missing_deps[*]}" $RED
        echo "Please install the missing dependencies and try again."
        exit 1
    fi
    
    print_status "${CHECK} All prerequisites satisfied!" $GREEN
}

# Function to setup backend
setup_backend() {
    print_header "${BUILD} Setting up Backend"
    
    if [ ! -d "$BACKEND_DIR" ]; then
        print_status "${CROSS} Backend directory not found" $RED
        return 1
    fi
    
    cd "$BACKEND_DIR"
    
    print_status "${INFO} Installing Go dependencies..." $BLUE
    go mod tidy
    go mod verify
    
    print_status "${INFO} Building backend..." $BLUE
    if [ -f "Makefile" ]; then
        make build
    else
        go build -o bin/go-pro-backend main.go
    fi
    
    print_status "${CHECK} Backend setup complete!" $GREEN
    cd ..
}

# Function to setup frontend
setup_frontend() {
    print_header "${BUILD} Setting up Frontend"
    
    if [ ! -d "$FRONTEND_DIR" ]; then
        print_status "${CROSS} Frontend directory not found" $RED
        return 1
    fi
    
    cd "$FRONTEND_DIR"
    
    print_status "${INFO} Installing Node.js dependencies..." $BLUE
    bun install
    
    print_status "${CHECK} Frontend setup complete!" $GREEN
    cd ..
}

# Function to run tests
run_tests() {
    print_header "${TEST} Running Tests"
    
    # Test backend
    if [ -d "$BACKEND_DIR" ]; then
        print_status "${INFO} Testing backend..." $BLUE
        cd "$BACKEND_DIR"
        if [ -f "Makefile" ]; then
            make test
        else
            go test -v ./...
        fi
        cd ..
    fi
    
    # Test course exercises
    if [ -d "$COURSE_DIR" ]; then
        print_status "${INFO} Testing course exercises..." $BLUE
        cd "$COURSE_DIR"
        
        # Test lesson 1 exercises
        if [ -d "code/lesson-01" ]; then
            cd "code/lesson-01"
            print_status "${INFO} Testing Lesson 1 examples..." $BLUE
            go run main.go
            cd ../..
        fi
        
        cd ..
    fi
    
    print_status "${CHECK} All tests completed!" $GREEN
}

# Function to start development servers
start_dev() {
    print_header "${ROCKET} Starting Development Servers"
    
    # Start backend
    if [ -d "$BACKEND_DIR" ]; then
        print_status "${INFO} Starting backend server..." $BLUE
        cd "$BACKEND_DIR"
        
        if [ -f "Makefile" ]; then
            make run &
        else
            go run main.go &
        fi
        
        BACKEND_PID=$!
        print_status "${CHECK} Backend started (PID: $BACKEND_PID)" $GREEN
        cd ..
        
        # Wait for backend to start
        sleep 3
        
        # Test backend health
        if command_exists curl; then
            print_status "${INFO} Testing backend health..." $BLUE
            if curl -s http://localhost:8080/api/v1/health > /dev/null; then
                print_status "${CHECK} Backend is healthy!" $GREEN
            else
                print_status "${WARNING} Backend health check failed" $YELLOW
            fi
        fi
    fi
    
    # Start frontend
    if [ -d "$FRONTEND_DIR" ]; then
        print_status "${INFO} Starting frontend server..." $BLUE
        cd "$FRONTEND_DIR"
        bun run dev &
        FRONTEND_PID=$!
        print_status "${CHECK} Frontend started (PID: $FRONTEND_PID)" $GREEN
        cd ..
    fi
    
    print_status "${ROCKET} Development servers are running!" $GREEN
    print_status "${INFO} Backend: http://localhost:8080" $BLUE
    print_status "${INFO} Frontend: http://localhost:3000" $BLUE
    print_status "${INFO} Press Ctrl+C to stop all servers" $YELLOW
    
    # Wait for interrupt
    trap 'kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit' INT
    wait
}

# Function to clean project
clean_project() {
    print_header "${CLEAN} Cleaning Project"
    
    # Clean backend
    if [ -d "$BACKEND_DIR" ]; then
        print_status "${INFO} Cleaning backend..." $BLUE
        cd "$BACKEND_DIR"
        if [ -f "Makefile" ]; then
            make clean
        else
            rm -rf bin/
            go clean
        fi
        cd ..
    fi
    
    # Clean frontend
    if [ -d "$FRONTEND_DIR" ]; then
        print_status "${INFO} Cleaning frontend..." $BLUE
        cd "$FRONTEND_DIR"
        rm -rf .next/
        rm -rf node_modules/.cache/
        cd ..
    fi
    
    print_status "${CHECK} Project cleaned!" $GREEN
}

# Function to show project status
show_status() {
    print_header "${INFO} Project Status"
    
    # Backend status
    if [ -d "$BACKEND_DIR" ]; then
        print_status "${BOOK} Backend:" $CYAN
        cd "$BACKEND_DIR"
        if [ -f "go.mod" ]; then
            GO_MOD_NAME=$(grep "^module" go.mod | cut -d' ' -f2)
            print_status "  Module: $GO_MOD_NAME" $NC
        fi
        if [ -f "bin/go-pro-backend" ]; then
            print_status "  ${CHECK} Binary built" $GREEN
        else
            print_status "  ${CROSS} Binary not built" $RED
        fi
        cd ..
    fi
    
    # Frontend status
    if [ -d "$FRONTEND_DIR" ]; then
        print_status "${BOOK} Frontend:" $CYAN
        cd "$FRONTEND_DIR"
        if [ -f "package.json" ]; then
            PACKAGE_NAME=$(grep '"name"' package.json | cut -d'"' -f4)
            print_status "  Package: $PACKAGE_NAME" $NC
        fi
        if [ -d "node_modules" ]; then
            print_status "  ${CHECK} Dependencies installed" $GREEN
        else
            print_status "  ${CROSS} Dependencies not installed" $RED
        fi
        cd ..
    fi
    
    # Course status
    if [ -d "$COURSE_DIR" ]; then
        print_status "${BOOK} Course:" $CYAN
        LESSON_COUNT=$(find "$COURSE_DIR/lessons" -name "lesson-*" -type d 2>/dev/null | wc -l)
        print_status "  Lessons: $LESSON_COUNT" $NC
        
        CODE_COUNT=$(find "$COURSE_DIR/code" -name "lesson-*" -type d 2>/dev/null | wc -l)
        print_status "  Code examples: $CODE_COUNT" $NC
    fi
}

# Main function
main() {
    print_header "${ROCKET} GO-PRO Development Environment"
    
    case "${1:-help}" in
        "setup")
            check_prerequisites
            setup_backend
            setup_frontend
            print_status "${CHECK} Setup complete! Run './dev.sh start' to begin development." $GREEN
            ;;
        "start"|"dev")
            start_dev
            ;;
        "test")
            run_tests
            ;;
        "clean")
            clean_project
            ;;
        "status")
            show_status
            ;;
        "help"|*)
            echo -e "${PURPLE}GO-PRO Development Script${NC}"
            echo ""
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  setup    - Setup development environment"
            echo "  start    - Start development servers"
            echo "  dev      - Alias for start"
            echo "  test     - Run all tests"
            echo "  clean    - Clean build artifacts"
            echo "  status   - Show project status"
            echo "  help     - Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 setup     # Initial setup"
            echo "  $0 start     # Start development"
            echo "  $0 test      # Run tests"
            ;;
    esac
}

# Run main function with all arguments
main "$@"
