# Project 2: Multi-Channel Support Bot

**Difficulty**: Intermediate  
**Duration**: 4-5 hours  
**Prerequisites**: Lessons 4-6

## Project Overview

Deploy an AI assistant that works across multiple messaging platforms simultaneously.

## Requirements

1. **Connect 3+ channels**:
   - Telegram
   - Discord
   - At least one more (WhatsApp/Slack)
2. **Unified identity**: Same conversation across channels
3. **Channel-specific config**: Different responses per channel

## Configuration Example

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "xxx",
      chatIds: ["user1"]
    },
    discord: {
      enabled: true,
      botToken: "xxx",
      allowedChannelIds: ["channel1"]
    }
  }
}
```

## Deliverables

- Working multi-channel bot
- Demo from each platform
- Configuration (sanitized)

## Challenges

- Handle different message formats
- Rate limiting per channel
- Session management
