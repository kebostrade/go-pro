# OpenClaw Discord Integration

Add OpenClaw as an AI assistant to your Discord server.

## Step 1: Create Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application"
3. Name your application

## Step 2: Set Up Bot

1. Click "Bot" in the left sidebar
2. Click "Reset Token" and copy your bot token
3. Enable "Message Content Intent" under "Privileged Gateway Intents"

## Step 3: Invite Bot to Server

1. Go to "OAuth2" > "URL Generator"
2. Select scopes: `bot`
3. Select permissions: `Send Messages`, `Read Message History`, `Manage Messages`
4. Copy the generated URL and open it in your browser
5. Select your server and authorize

## Step 4: Configure OpenClaw

Add to `config/openclaw.json5`:

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

## Step 5: Restart OpenClaw

```bash
docker compose restart
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
