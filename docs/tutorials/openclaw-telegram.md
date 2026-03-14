# OpenClaw Telegram Integration

Connect your OpenClaw agent to Telegram for instant messaging access.

## Step 1: Create a Telegram Bot

1. Open Telegram and search for @BotFather
2. Send `/newbot` command
3. Follow prompts to name your bot
4. Copy the HTTP API token provided

## Step 2: Configure OpenClaw

Add to your `config/openclaw.json5`:

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "YOUR_BOT_TOKEN_HERE",
      chatIds: []
    }
  }
}
```

## Step 3: Restart OpenClaw

```bash
# If using Docker
docker compose restart

# If using npm
# Restart your OpenClaw process
```

## Step 4: Find Your Chat ID

1. Message your bot on Telegram
2. Visit: `https://api.telegram.org/bot<TOKEN>/getUpdates`
3. Find your `chat` > `id` in the JSON response

## Step 5: Restrict Access (Production)

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

Only users with specified chat IDs can interact with your bot.

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
