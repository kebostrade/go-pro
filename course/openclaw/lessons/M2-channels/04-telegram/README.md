# Lesson OC-04: Telegram Bot Integration

## Overview

This lesson covers creating and configuring a Telegram bot for OpenClaw.

## Why Telegram?

- Easy to set up
- No phone number required for bot
- Free and reliable
- Rich bot API

## Step 1: Create Your Bot

1. Open Telegram and search for **@BotFather**
2. Send `/newbot` command
3. Follow prompts:
   - Bot name: `My OpenClaw`
   - Username: `my_openclaw_bot` (must end in `bot`)
4. Copy the HTTP API token provided

## Step 2: Configure OpenClaw

Add to your `config/openclaw.json5`:

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "YOUR_BOT_TOKEN_HERE",
      chatIds: []  // Leave empty for now
    }
  }
}
```

## Step 3: Restart OpenClaw

```bash
# Docker
docker compose restart

# Source
# Restart your process
```

## Step 4: Find Your Chat ID

1. Message your bot on Telegram (send anything)
2. Visit: `https://api.telegram.org/bot<TOKEN>/getUpdates`
3. Look for `"chat": {"id": 123456789}`
4. Copy your chat ID

## Step 5: Secure Your Bot (Important!)

Restrict access to authorized users only:

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "YOUR_BOT_TOKEN",
      chatIds: ["123456789", "987654321"]
    }
  }
}
```

Now only users with these IDs can interact with your bot.

## Bot Commands

OpenClaw supports these commands:
- `/help` - Show available commands
- `/reset` - Clear conversation memory
- `/status` - Check agent status

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Bot not responding | Check token, restart OpenClaw |
| Can't find chat ID | Send message to bot first, then check getUpdates |
| "Forbidden" error | Bot was blocked by user |

---

**Next**: [Lesson 5: Discord](05-discord/README.md)
