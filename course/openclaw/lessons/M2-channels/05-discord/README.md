# Lesson OC-05: Discord Bot Integration

## Overview

This lesson covers adding OpenClaw as a Discord bot.

## Step 1: Create Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click **New Application**
3. Name your application

## Step 2: Set Up Bot

1. Click **Bot** in left sidebar
2. Click **Reset Token** → Copy token
3. Under **Privileged Gateway Intents**, enable:
   - ✅ Message Content Intent (required!)

## Step 3: Invite Bot to Server

1. Go to **OAuth2** → **URL Generator**
2. Select scopes: `bot`
3. Select permissions:
   - ✅ Send Messages
   - ✅ Read Message History
   - ✅ Manage Messages
4. Copy generated URL
5. Open URL in browser → Select server → Authorize

## Step 4: Configure OpenClaw

```json5
{
  channels: {
    discord: {
      enabled: true,
      botToken: "YOUR_DISCORD_BOT_TOKEN",
      allowedChannelIds: []
    }
  }
}
```

## Step 5: Restart and Test

```bash
docker compose restart
```

Mention your bot in a channel to start a conversation:
```
@YourBot Hello!
```

## Restricting Access

```json5
{
  channels: {
    discord: {
      enabled: true,
      botToken: "YOUR_TOKEN",
      allowedChannelIds: ["123456789", "987654321"],
      allowedGuildIds: ["111222333"]
    }
  }
}
```

---

**Next**: [Lesson 6: WhatsApp & Signal](06-whatsapp-signal/README.md)
