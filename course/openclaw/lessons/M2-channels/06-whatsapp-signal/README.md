# Lesson OC-06: WhatsApp & Signal Integration

## ⚠️ Important Warning

WhatsApp may ban accounts using unofficial clients. Use this for personal use only, not for automated messaging to others.

## WhatsApp Setup

### Step 1: Enable WhatsApp Channel

```json5
{
  channels: {
    whatsapp: {
      enabled: true
    }
  }
}
```

### Step 2: Restart OpenClaw

```bash
docker compose restart
```

### Step 3: Pair Device

1. Open dashboard: `http://localhost:18789`
2. Go to **Channels** → **WhatsApp**
3. Scan QR code with WhatsApp app
4. Wait for "Connected" status

## Session Management

WhatsApp sessions expire. If disconnected:
1. Delete session from data directory
2. Restart OpenClaw
3. Re-scan QR code

## Signal Setup

Signal requires more setup due to privacy features:

```json5
{
  channels: {
    signal: {
      enabled: true,
      phoneNumber: "+1234567890"
    }
  }
}
```

Signal registration requires additional verification steps.

## Best Practices

- Don't use for bulk automated messaging
- Keep sessions active
- Monitor for account bans

---

**Next**: [Module 3: Skills & Extensibility](M3-skills/README.md)
