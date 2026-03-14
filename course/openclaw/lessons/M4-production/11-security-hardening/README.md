# Lesson OC-11: Security Best Practices

## Overview

This lesson covers securing your OpenClaw deployment.

## The Security Reality

⚠️ **Important**: Research found 93.4% of exposed OpenClaw instances are vulnerable. Security is critical!

## Essential Security Measures

### 1. Restrict Channel Access

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

- Never expose port 18789 directly to internet
- Always use reverse proxy with authentication
- Use firewall to restrict access

### 3. API Key Management

```bash
# ✅ Use environment variables
ANTHROPIC_API_KEY=sk-ant-xxx

# ❌ Never commit keys
# Add to .gitignore:
echo ".env" >> .gitignore
```

### 4. Keep Updated

```bash
# Regular updates
docker compose pull && docker compose up -d
```

### 5. Enable Audit Logging

```bash
# Monitor logs for suspicious activity
docker compose logs --tail 100 | grep -i error
docker compose logs --tail 100 | grep -i unauthorized
```

## Security Checklist

- [ ] Channel access restricted to authorized users
- [ ] Gateway behind reverse proxy
- [ ] API keys in environment variables only
- [ ] Firewall configured
- [ ] Regular backups
- [ ] Updates applied promptly
- [ ] Logs monitored

## Incident Response

If compromised:
1. Revoke API keys immediately
2. Restrict all channel access
3. Check for unauthorized actions in logs
4. Rotate all credentials
5. Update OpenClaw to latest version

---

**Congratulations! You've completed the OpenClaw course!**

## Next Steps

- Build your own custom Skills
- Explore advanced automation
- Contribute to the community
