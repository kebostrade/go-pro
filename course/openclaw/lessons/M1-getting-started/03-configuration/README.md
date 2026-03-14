# Lesson OC-03: Configuration

## Overview

This lesson covers configuring OpenClaw with LLM providers and understanding the configuration file.

## Configuration File Location

- **Docker**: `./config/openclaw.json5`
- **Source**: `./config.json`

## Basic Configuration

```json5
{
  // AI model provider settings
  providers: {
    anthropic: {
      enabled: true,
      defaultModel: "claude-sonnet-4-20250514"
    },
    openai: {
      enabled: false,
      defaultModel: "gpt-4o"
    }
  },

  // Gateway settings
  gateway: {
    port: 18789,
    host: "0.0.0.0"
  },

  // Channel settings (covered in next module)
  channels: {}
}
```

## Supported LLM Providers

### Anthropic (Claude)

```json5
providers: {
  anthropic: {
    enabled: true,
    defaultModel: "claude-sonnet-4-20250514",
    // Options: claude-3-5-sonnet-20241022, claude-3-opus-20240229
  }
}
```

Environment variable: `ANTHROPIC_API_KEY`

### OpenAI

```json5
providers: {
  openai: {
    enabled: true,
    defaultModel: "gpt-4o"
    // Options: gpt-4o, gpt-4o-mini, gpt-4-turbo
  }
}
```

Environment variable: `OPENAI_API_KEY`

### Ollama (Local)

```json5
providers: {
  ollama: {
    enabled: true,
    defaultModel: "llama3.2:3b",
    baseURL: "http://localhost:11434"
  }
}
```

## Environment Variables

Create a `.env` file:

```bash
# Required: At least one provider
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxxx

# Optional: Custom settings
OPENCLAW_LOG_LEVEL=info
TZ=America/New_York
```

## Gateway Settings

| Setting | Default | Description |
|---------|---------|-------------|
| port | 18789 | Gateway listening port |
| host | "0.0.0.0" | Bind address (use 127.0.0.1 for local only) |

## Testing Your Configuration

1. Restart OpenClaw after changes
2. Check logs for provider connection
3. Send a test message via dashboard

```bash
# Docker
docker compose restart
docker compose logs -f

# Source
npm run start
```

---

## Quiz

1. Which configuration file format does OpenClaw use?
2. What environment variable is needed for Anthropic?
3. Why set host to "0.0.0.0" in Docker?

---

**Next**: [Module 2: Channel Integration](M2-channels/README.md)
