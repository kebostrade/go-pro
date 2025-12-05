# MCP Server Troubleshooting Guide

## Issue: "Error reading message: fetching message: multiple Read calls return no data or error"

### Root Causes Identified

1. **Too many MCP servers running** (109 processes → reduced to 16)
2. **Remote server timeouts**:
   - `mcp.api.coingecko.com` - Timing out
   - `docs.mcp.cloudflare.com` - Timing out
3. **High memory consumption** from multiple browser automation servers

### Applied Fixes

#### 1. Disabled Problematic Servers in `mcp.json`:

```json
{
  "coingecko": {
    "disabled": true  // ← Added: Server was timing out
  },
  "Cloudflare-docs": {
    "disabled": true  // ← Added: Server was timing out
  },
  "Cloudflare-observability": {
    "disabled": true  // ← Already disabled
  }
}
```

#### 2. Cleaned Up Excess Processes

Terminated 93 stale MCP processes using `scripts/cleanup-mcp.sh`

### Diagnostic Scripts Created

#### `/scripts/diagnose-mcp.sh`
- Counts running MCP processes
- Identifies high memory usage servers
- Tests remote server connectivity
- Checks for zombie processes

**Usage:**
```bash
cd /home/dima/Desktop/FUN/go-pro
./scripts/diagnose-mcp.sh
```

#### `/scripts/cleanup-mcp.sh`
- Gracefully terminates MCP servers
- Cleans up Docker MCP containers
- Reports cleanup statistics

**Usage:**
```bash
cd /home/dima/Desktop/FUN/go-pro
./scripts/cleanup-mcp.sh
```

### Recommendations

#### 1. **Disable Resource-Heavy Servers** (if not actively using)

Add `"disabled": true` to these servers in `mcp.json`:

```json
{
  "Bright Data": {
    "disabled": true  // Already disabled - web scraping
  },
  "playwright": {
    "disabled": true  // Consider if not actively testing
  },
  "puppeteer": {
    "disabled": true  // Consider if not actively testing
  },
  "Firebase MCP": {
    "disabled": true  // Consider if not using Firebase
  }
}
```

#### 2. **Monitor MCP Process Count**

```bash
# Check MCP process count
ps aux | grep -E "(mcp|context7|playwright)" | grep -v grep | wc -l

# Should be < 20 for normal usage
# If > 50, run cleanup script
```

#### 3. **Periodic Cleanup**

Add to your shell rc file (`~/.zshrc`):

```bash
alias mcp-status='ps aux | grep -E "(mcp|context7|playwright)" | grep -v grep | wc -l'
alias mcp-clean='~/Desktop/FUN/go-pro/scripts/cleanup-mcp.sh'
alias mcp-diagnose='~/Desktop/FUN/go-pro/scripts/diagnose-mcp.sh'
```

#### 4. **Restart Editor After Cleanup**

After running cleanup script:
1. Completely close Cursor/VS Code
2. Verify all processes terminated: `ps aux | grep mcp`
3. Restart editor
4. Only enabled MCP servers will start

### Monitoring

#### Check for Errors:
```bash
# Cursor logs
tail -f ~/.config/Cursor/logs/*.log | grep "Error reading message"

# VS Code logs
tail -f ~/.config/Code/logs/*.log | grep "Error reading message"
```

#### Check MCP Health:
```bash
# Run diagnostic
./scripts/diagnose-mcp.sh

# Expected output:
# - MCP processes: < 20
# - Remote servers: responding or disabled
# - No zombie processes
```

### When Errors Return

If you see "Error reading message" again:

1. **Run diagnostics:**
   ```bash
   ./scripts/diagnose-mcp.sh
   ```

2. **If process count > 50:**
   ```bash
   ./scripts/cleanup-mcp.sh
   # Then restart editor
   ```

3. **If remote servers failing:**
   - Edit `mcp.json`
   - Add `"disabled": true` to failing servers
   - Restart editor

4. **If specific server crashing:**
   - Check logs: `~/.config/Cursor/logs/`
   - Disable that server in `mcp.json`
   - Report issue to server maintainer

### Optimized MCP Configuration

For stable operation with minimal errors, keep only essential servers enabled:

**Essential (always on):**
- `filesystem` - File operations
- `memory` - Context retention
- `sequential-thinking` - Complex reasoning

**Development (enable when needed):**
- `playwright`/`puppeteer` - Browser testing
- `postgres`/`redis` - Database work
- `brave-search` - Web search
- `Context7` - Enhanced context

**Optional (disable unless actively using):**
- `Firebase MCP`
- `stripe`
- `coingecko`
- `Cloudflare-*`
- `Bright Data`
- `airtable`

### Current Status

✅ **Fixed:**
- Disabled timeout-prone remote servers (CoinGecko, Cloudflare Docs)
- Cleaned up 93 stale MCP processes
- Created diagnostic and cleanup tools

✅ **Active MCP Servers:** 16 (down from 109)

✅ **No zombie processes**

⚠️ **Next Steps:**
1. Restart your editor to apply mcp.json changes
2. Monitor for "Error reading message" in next session
3. Run diagnostics if issues return

### Support

If issues persist after applying fixes:

1. **Collect diagnostics:**
   ```bash
   ./scripts/diagnose-mcp.sh > mcp-diagnostic-report.txt
   ```

2. **Check specific server logs:**
   ```bash
   # Find failing server in logs
   grep -r "Error reading message" ~/.config/Cursor/logs/ | tail -50
   ```

3. **Test individual servers:**
   ```bash
   # Test remote connectivity
   curl -I https://mcp.api.coingecko.com/sse
   curl -I https://docs.mcp.cloudflare.com/sse
   ```

4. **Restart editor with clean state:**
   ```bash
   ./scripts/cleanup-mcp.sh
   # Close editor completely
   # Restart editor
   ```

---

**Created:** December 3, 2025  
**Issue:** Multiple "Error reading message" errors every 15 seconds  
**Resolution:** Disabled problematic remote servers, cleaned up 93 stale processes, created monitoring tools
