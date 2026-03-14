# OpenClaw WhatsApp Integration

Connect OpenClaw to WhatsApp using the WhatsApp Web protocol.

## Important Warning

WhatsApp may ban accounts that use unofficial clients, especially for automated messaging. Use this channel judiciously.

## Step 1: Enable WhatsApp Channel

Add to `config/openclaw.json5`:

```json5
{
  channels: {
    whatsapp: {
      enabled: true
    }
  }
}
```

## Step 2: Restart OpenClaw

```bash
docker compose restart
```

## Step 3: Scan QR Code

1. Access the OpenClaw dashboard at `http://localhost:18789`
2. Navigate to Channels > WhatsApp
3. Scan the QR code with your WhatsApp app
4. Wait for "Connected" status

## Session Management

WhatsApp sessions may expire. If disconnected:
1. Delete the old session from data directory
2. Restart OpenClaw
3. Re-scan the QR code

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
