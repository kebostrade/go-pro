#!/bin/bash

# Quick Start Script for Gin Web Framework Examples

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     🚀 Gin Web Framework - Quick Start Guide              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
echo ""

# Ask which example to run
echo "📋 Available Examples:"
echo ""
echo "  1. Complete Application (main.go)"
echo "     - Full-featured web application"
echo "     - REST API with authentication"
echo "     - HTML templates"
echo ""
echo "  2. Gin Basics (examples/gin_basics.go)"
echo "     - Basic routing and handlers"
echo "     - Middleware examples"
echo "     - Route groups"
echo ""
echo "  3. REST API (examples/gin_rest_api.go)"
echo "     - Complete REST API"
echo "     - CRUD operations"
echo "     - Request validation"
echo ""
echo "  4. Templates (examples/gin_templates.go)"
echo "     - HTML template rendering"
echo "     - Static file serving"
echo ""
echo "  5. Advanced Features (examples/gin_advanced.go)"
echo "     - File uploads"
echo "     - Session management"
echo "     - Streaming responses"
echo ""
echo "  6. Run All Examples (Sequential)"
echo ""
read -p "Select an example (1-6): " choice

case $choice in
    1)
        echo ""
        echo "🚀 Starting Complete Application..."
        echo "   Visit: http://localhost:8080"
        echo ""
        go run main.go
        ;;
    2)
        echo ""
        echo "🚀 Starting Gin Basics..."
        echo "   Visit: http://localhost:8080"
        echo ""
        go run examples/gin_basics.go
        ;;
    3)
        echo ""
        echo "🚀 Starting REST API..."
        echo "   API available at: http://localhost:8080/api"
        echo ""
        go run examples/gin_rest_api.go
        ;;
    4)
        echo ""
        echo "🚀 Starting Templates Example..."
        echo "   Visit: http://localhost:8080"
        echo ""
        go run examples/gin_templates.go
        ;;
    5)
        echo ""
        echo "🚀 Starting Advanced Features..."
        echo "   Visit: http://localhost:8080"
        echo ""
        go run examples/gin_advanced.go
        ;;
    6)
        echo ""
        echo "🚀 Running all examples sequentially..."
        echo ""

        echo "═══ Example 1: Gin Basics ═══"
        timeout 5 go run examples/gin_basics.go &
        sleep 2
        echo "✅ Basics example completed"
        echo ""

        echo "═══ Example 2: REST API ═══"
        timeout 5 go run examples/gin_rest_api.go &
        sleep 2
        echo "✅ REST API example completed"
        echo ""

        echo "═══ Example 3: Templates ═══"
        timeout 5 go run examples/gin_templates.go &
        sleep 2
        echo "✅ Templates example completed"
        echo ""

        echo "═══ Example 4: Advanced Features ═══"
        timeout 5 go run examples/gin_advanced.go &
        sleep 2
        echo "✅ Advanced features example completed"
        echo ""

        echo "✅ All examples completed!"
        ;;
    *)
        echo "❌ Invalid selection"
        exit 1
        ;;
esac
