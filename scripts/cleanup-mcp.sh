#!/bin/bash

# MCP Process Cleanup Script
# Safely terminates MCP server processes

echo "🧹 Cleaning up MCP processes..."
echo ""

# Count before
BEFORE=$(ps aux | grep -E "(mcp|context7|playwright|firebase-mcp|redis-mcp-server|airtable-mcp)" | grep -v grep | grep -v "cleanup-mcp" | wc -l)
echo "📊 MCP processes before cleanup: $BEFORE"

# Kill MCP server processes (not the editors)
echo ""
echo "🛑 Terminating MCP servers..."

# Gracefully terminate with SIGTERM first
pkill -15 -f "context7-mcp" 2>/dev/null
pkill -15 -f "mcp-server-playwright" 2>/dev/null
pkill -15 -f "mcp-server-puppeteer" 2>/dev/null
pkill -15 -f "firebase-mcp" 2>/dev/null
pkill -15 -f "redis-mcp-server" 2>/dev/null
pkill -15 -f "mcp-server-postgres" 2>/dev/null
pkill -15 -f "mcp-server-filesystem" 2>/dev/null
pkill -15 -f "mcp-server-everything" 2>/dev/null
pkill -15 -f "mcp-server-memory" 2>/dev/null
pkill -15 -f "mcp-server-sequential-thinking" 2>/dev/null
pkill -15 -f "mcp-server-brave-search" 2>/dev/null
pkill -15 -f "airtable-mcp-server" 2>/dev/null
pkill -15 -f "mcp-remote" 2>/dev/null
pkill -15 -f "next-devtools-mcp" 2>/dev/null
pkill -15 -f "svelte-mcp" 2>/dev/null

# Wait for graceful shutdown
sleep 2

# Force kill any remaining
pkill -9 -f "mcp-server" 2>/dev/null
pkill -9 -f "context7-mcp" 2>/dev/null
pkill -9 -f "mcp-remote" 2>/dev/null

# Clean up Docker MCP containers
echo ""
echo "🐳 Cleaning up Docker MCP containers..."
docker ps -a | grep "mcp/stripe" | awk '{print $1}' | xargs -r docker rm -f 2>/dev/null

# Count after
sleep 1
AFTER=$(ps aux | grep -E "(mcp|context7|playwright|firebase-mcp|redis-mcp-server)" | grep -v grep | grep -v "cleanup-mcp" | wc -l)
echo ""
echo "✅ Cleanup complete!"
echo "📊 MCP processes after cleanup: $AFTER"
echo "🗑️  Terminated: $((BEFORE - AFTER)) processes"
echo ""
echo "💡 Next steps:"
echo "   1. Restart your editor (Cursor/VS Code) to start fresh MCP servers"
echo "   2. Only enabled MCP servers in mcp.json will start"
echo "   3. Monitor: ps aux | grep mcp | wc -l"
