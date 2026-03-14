# Lesson OC-09: Docker Deployment

## Overview

This lesson covers production-ready Docker deployment for OpenClaw.

## Why Docker?

- Environment isolation
- Easy updates
- Consistent deployments
- Simple backups

## docker-compose.yml

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
    environment:
      - NODE_ENV=production
      - TZ=America/New_York
    env_file:
      - .env
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:18789/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Environment File

```bash
# .env
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxxx
```

## Operations

### Start

```bash
docker compose up -d
```

### View Logs

```bash
docker compose logs -f
docker compose logs --tail 100
```

### Restart

```bash
docker compose restart
```

### Update

```bash
docker compose pull
docker compose up -d
```

### Backup

```bash
tar -czf openclaw-backup-$(date +%Y%m%d).tar.gz config data .env
```

## Troubleshooting

| Issue | Solution |
|--------|----------|
| Container won't start | Check logs: `docker compose logs` |
| API errors | Verify .env keys are correct |
| Port conflict | Change port mapping |

---

**Next**: [Lesson 10: VPS Production](10-vps-production/README.md)
