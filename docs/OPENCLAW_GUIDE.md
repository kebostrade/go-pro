# 🦞 Guide to OpenClaw

The complete guide to building, deploying, and securing self-hosted AI agents with OpenClaw.

## What is OpenClaw?

OpenClaw is the fastest-growing open-source AI agent platform (150K+ GitHub stars). It lets you run a personal AI assistant on your own infrastructure and connect it to messaging apps you already use—Telegram, Discord, WhatsApp, Slack, and more.

## Why OpenClaw?

| Feature | Benefit |
|---------|---------|
| **Self-hosted** | Your data never leaves your infrastructure |
| **Multi-channel** | Connect to Telegram, Discord, WhatsApp, Slack, Signal, iMessage |
| **Extensible** | Build custom Skills for any automation |
| **Privacy-first** | Complete control over your AI assistant |
| **Open source** | Free to use, community-driven |

---

## 🚀 Quick Start

### 1. Install OpenClaw

```bash
# macOS
brew install openclaw

# Docker (cross-platform)
mkdir openclaw && cd openclaw
docker run -d -p 18789:18789 openclaw/openclaw
```

### 2. Configure API Key

```bash
# .env file
ANTHROPIC_API_KEY=sk-ant-xxx
# OR
OPENAI_API_KEY=sk-xxx
```

### 3. Connect a Channel

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "YOUR_TOKEN",
      chatIds: ["YOUR_CHAT_ID"]
    }
  }
}
```

### 4. Start Chatting

Message your bot and ask questions!

---

## 📚 Learning Paths

### Beginner Path

Perfect for those new to AI agents:

1. **[OpenClaw Full Tutorial](tutorials/openclaw-full-tutorial.md)** — Complete start-to-finish guide
2. **[Installation Tutorial](tutorials/openclaw-installation.md)** — Step-by-step setup
3. **[Configuration Guide](course/openclaw/lessons/M1-getting-started/03-configuration/README.md)** — Connect your LLM provider
4. **[Telegram Integration](tutorials/openclaw-telegram.md)** — Your first channel

**Estimated Time**: 2-3 hours

### Intermediate Path

For developers who want more:

1. **[Docker Deployment](tutorials/openclaw-docker-deployment.md)** — Production-ready container setup
2. **[Discord Integration](tutorials/openclaw-discord.md)** — Add more channels
3. **[Skills System](tutorials/openclaw-skills.md)** — Extend capabilities
4. **[Custom Skills](course/openclaw/lessons/M3-skills/08-custom-skills/README.md)** — Build your own

**Estimated Time**: 4-6 hours

### Advanced Path

For production deployments:

1. **[VPS Deployment](tutorials/openclaw-vps-deployment.md)** — Deploy on cloud
2. **[Security Hardening](tutorials/openclaw-security.md)** — Secure your instance
3. **[Ollama Setup](tutorials/openclaw-ollama.md)** — Local AI models

**Estimated Time**: 6-8 hours

---

## 📖 Course: Build Your Own AI Agent

A structured 4-week course with hands-on projects.

### Module 1: Getting Started
| Lesson | Topic |
|--------|-------|
| [OC-01](course/openclaw/lessons/M1-getting-started/01-introduction/README.md) | Architecture Deep Dive |
| [OC-02](course/openclaw/lessons/M1-getting-started/02-installation/README.md) | Installation |
| [OC-03](course/openclaw/lessons/M1-getting-started/03-configuration/README.md) | Configuration |

### Module 2: Channel Integration
| Lesson | Topic |
|--------|-------|
| [OC-04](course/openclaw/lessons/M2-channels/04-telegram/README.md) | Telegram |
| [OC-05](course/openclaw/lessons/M2-channels/05-discord/README.md) | Discord |
| [OC-06](course/openclaw/lessons/M2-channels/06-whatsapp-signal/README.md) | WhatsApp & Signal |

### Module 3: Skills & Extensibility
| Lesson | Topic |
|--------|-------|
| [OC-07](course/openclaw/lessons/M3-skills/07-skills-basics/README.md) | Skills Basics |
| [OC-08](course/openclaw/lessons/M3-skills/08-custom-skills/README.md) | Custom Skills |

### Module 4: Production
| Lesson | Topic |
|--------|-------|
| [OC-09](course/openclaw/lessons/M4-production/09-docker-deployment/README.md) | Docker Deployment |
| [OC-10](course/openclaw/lessons/M4-production/10-vps-production/README.md) | VPS Production |
| [OC-11](course/openclaw/lessons/M4-production/11-security-hardening/README.md) | Security |

---

## 🎯 Projects

| Project | Description | Difficulty |
|---------|-------------|------------|
| [P1](course/openclaw/projects/P1-personal-ai-assistant/README.md) | Personal AI Assistant | Beginner |
| [P2](course/openclaw/projects/P2-multi-channel-bot/README.md) | Multi-Channel Bot | Intermediate |
| [P3](course/openclaw/projects/P3-automation-workflows/README.md) | Custom Skills | Intermediate |
| [P4](course/openclaw/projects/P4-production-deployment/README.md) | Production VPS | Advanced |

---

## 🔧 Reference

### Configuration

```json5
{
  providers: {
    anthropic: { enabled: true, defaultModel: "claude-sonnet-4-20250514" },
    openai: { enabled: false, defaultModel: "gpt-4o" },
    ollama: { enabled: false, defaultModel: "llama3.2:3b" }
  },
  gateway: { port: 18789, host: "0.0.0.0" },
  channels: {}
}
```

### Common Commands

```bash
# Docker
docker compose up -d           # Start
docker compose logs -f         # View logs
docker compose restart         # Restart
docker compose pull            # Update

# CLI
openclaw doctor                # Verify setup
openclaw start                 # Start
```

### Ports

| Port | Service |
|------|---------|
| 18789 | OpenClaw Gateway (default) |

---

## ⚠️ Security Note

Research shows **93.4%** of exposed OpenClaw instances are vulnerable. Always:

- ✅ Restrict channel access to authorized users
- ✅ Never expose Gateway to public internet without proxy
- ✅ Keep API keys in environment variables
- ✅ Regular updates

See [Security Best Practices](tutorials/openclaw-security.md)

---

## 🔗 Resources

- [GitHub](https://github.com/openclaw/openclaw)
- [Official Docs](https://docs.openclaw.com)
- [Skills Library](https://github.com/openclaw/skills)
- [Community](https://github.com/openclaw/openclaw/discussions)

---

## Architecture Overview

```
User
  │
  ▼
Messaging App (Telegram/Discord)
  │
  ▼
Gateway (OpenClaw) ─── Port 18789
  │
  ▼
Agent (Brain)
  │
  ▼
Skills (Web/Memory)
  │
  ▼
LLM Provider (OpenAI/Anthropic)
```

---

*Last updated: March 2026*
