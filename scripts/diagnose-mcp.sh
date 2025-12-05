#!/bin/bash

# MCP Server Diagnostic Script
# Helps identify problematic MCP servers

echo "🔍 MCP Server Diagnostics"
echo "=========================="
echo ""

# Count running MCP processes
echo "📊 MCP Process Count:"
MCP_COUNT=$(ps aux | grep -E "(mcp|context7|playwright|firebase|stripe|redis-mcp)" | grep -v grep | wc -l)
echo "   Total MCP-related processes: $MCP_COUNT"
echo ""

# Show high memory usage MCP processes
echo "💾 High Memory MCP Processes (>100MB):"
ps aux | grep -E "(mcp|context7|playwright|firebase|stripe|redis-mcp)" | grep -v grep | awk '$6 > 100000 {printf "   %s: %dMB\n", $11, $6/1024}' | head -10
echo ""

# Check for zombie/defunct processes
echo "🧟 Zombie/Defunct Processes:"
ZOMBIE_COUNT=$(ps aux | grep -E "(mcp|context7|playwright)" | grep -E "(defunct|zombie)" | wc -l)
if [ $ZOMBIE_COUNT -gt 0 ]; then
    echo "   ⚠️  Found $ZOMBIE_COUNT zombie processes"
    ps aux | grep -E "(mcp|context7|playwright)" | grep -E "(defunct|zombie)" | head -5
else
    echo "   ✅ No zombie processes found"
fi
echo ""

# Test connectivity to remote MCP servers
echo "🌐 Testing Remote MCP Servers:"
REMOTE_SERVERS=(
    "https://mcp.api.coingecko.com/sse"
    "https://docs.mcp.cloudflare.com/sse"
    "https://bindings.mcp.cloudflare.com/sse"
    "https://observability.mcp.cloudflare.com/sse"
)

for server in "${REMOTE_SERVERS[@]}"; do
    echo -n "   Testing $server... "
    if timeout 5 curl -s -o /dev/null -w "%{http_code}" "$server" > /tmp/mcp_test 2>&1; then
        CODE=$(cat /tmp/mcp_test)
        if [ "$CODE" = "200" ] || [ "$CODE" = "401" ]; then
            echo "✅ ($CODE)"
        else
            echo "⚠️  HTTP $CODE"
        fi
    else
        echo "❌ Timeout/Failed"
    fi
done
rm -f /tmp/mcp_test
echo ""

# Check for errors in recent logs
echo "📋 Recent MCP Errors (last 20):"
grep -r "Error reading message" ~/.config/Cursor/logs/ 2>/dev/null | tail -20 || echo "   No error logs found in Cursor"
grep -r "Error reading message" ~/.config/Code/logs/ 2>/dev/null | tail -20 || echo "   No error logs found in VS Code"
echo ""

# Recommendations
echo "💡 Recommendations:"
echo ""
echo "1. If remote servers are timing out, consider disabling them:"
echo "   Edit mcp.json and add '\"disabled\": true' to problematic servers"
echo ""
echo "2. If too many MCP processes, restart your editor to clean up"
echo ""
echo "3. High memory usage? Consider disabling heavy servers like:"
echo "   - Bright Data (web scraping)"
echo "   - Playwright/Puppeteer (browser automation)"
echo "   - Firebase (if not actively using)"
echo ""
echo "4. Kill all MCP processes and restart editor:"
echo "   pkill -f 'mcp|context7|playwright|firebase-mcp|redis-mcp'"
echo ""
