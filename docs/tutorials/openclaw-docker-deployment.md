# OpenClaw Docker Deployment

This tutorial covers deploying OpenClaw using Docker Compose for reliable long-term operation.

## Why Docker?

- Environment isolation prevents dependency conflicts
- One-command start/stop
- Automatic container restarts
- Easy upgrades

## Setup

### 1. Create docker-compose.yml

```yaml
version: "3.8"

services:
  openclaw:
    image: openclaw/openclaw:latest
    container_name: openclaw
    restart: unless-stopped
    ports:
      - "18789:18789"
    volumes:
      - ./config:/root/.config/openclaw
      - ./data:/root/.openclaw
    env_file:
      - .env
```

### 2. Create .env file

```bash
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxxx
```

### 3. Create config/openclaw.json5

```json5
{
  providers: {
    anthropic: {
      enabled: true,
      defaultModel: "claude-sonnet-4-20250514"
    }
  },
  gateway: {
    port: 18789,
    host: "0.0.0.0"
  },
  channels: {}
}
```

### 4. Start Container

```bash
docker compose up -d
docker compose logs -f
```

## Day-to-Day Operations

```bash
docker compose down          # Stop
docker compose restart      # Restart
docker compose pull && docker compose up -d  # Upgrade
tar -czf backup-$(date +%Y%m%d).tar.gz config data .env  # Backup
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
