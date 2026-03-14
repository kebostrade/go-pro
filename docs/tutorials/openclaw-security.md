# OpenClaw Security Best Practices

Secure your OpenClaw deployment against common vulnerabilities.

## The Security Reality

Research shows 93.4% of exposed OpenClaw instances are vulnerable. Understanding risks is critical.

## Essential Security Measures

### 1. Restrict Channel Access

Never leave channels open to everyone:

```json5
// ❌ BAD - Anyone can access
channels: {
  telegram: { enabled: true, botToken: "xxx", chatIds: [] }
}

// ✅ GOOD - Only authorized users
channels: {
  telegram: { 
    enabled: true, 
    botToken: "xxx", 
    chatIds: ["123456789"] 
  }
}
```

### 2. Protect Your Gateway

- Never expose to public internet without authentication
- Use reverse proxy with authentication
- Configure firewall rules

### 3. API Key Management

```bash
# ✅ Use environment variables
ANTHROPIC_API_KEY=sk-ant-xxx

# ❌ Never commit keys
# Remove .env from Git
```

### 4. Keep Updated

```bash
# Regularly pull updates
docker compose pull && docker compose up -d
```

### 5. Enable Audit Logging

Monitor for suspicious activity in logs:

```bash
docker compose logs --tail 100 | grep -i error
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
