# 🦞 OpenClaw Tutorials

Welcome to the OpenClaw tutorials hub! OpenClaw is the fastest-growing open-source AI agent platform, with over 150,000 GitHub stars. These tutorials will guide you from installation to advanced production deployments.

---

## 📖 Available Tutorials

### 🦞 [OpenClaw Full Tutorial](openclaw-full-tutorial.md) **← START HERE!**

**Duration:** 2-3 hours | **Level:** Beginner to Intermediate | **Hands-On:** 100%

The comprehensive guide to getting started with OpenClaw, from installation to production deployment.

**What You'll Learn:**
- ✅ Understanding OpenClaw architecture (Gateway, Nodes, Skills)
- ✅ Installation via Homebrew, source, or Docker
- ✅ Configuration with LLM providers (OpenAI, Anthropic, Ollama)
- ✅ Connecting messaging channels (Telegram, Discord, WhatsApp, Slack)
- ✅ Using and creating Skills
- ✅ Security best practices
- ✅ Docker deployment
- ✅ Production VPS deployment with reverse proxy
- ✅ Troubleshooting common issues

**Perfect for:** Developers who want a complete, production-ready OpenClaw setup

**Prerequisites:**
- API key from OpenAI, Anthropic, or Ollama
- 4GB RAM minimum (8GB recommended)
- Git installed
- Node.js 20+ or Bun 1.1+

---

### 📦 Installation Tutorials

#### [OpenClaw Installation Guide](./openclaw-installation.md)
**Duration:** 30 minutes | **Level:** Beginner

Step-by-step installation for macOS, Linux, and Windows (WSL2).

**Covers:**
- Homebrew installation (macOS)
- Source installation (all platforms)
- Initial configuration
- Verification with `openclaw doctor`

---

#### [OpenClaw Docker Deployment](./openclaw-docker-deployment.md)
**Duration:** 45 minutes | **Level:** Intermediate

Complete Docker Compose deployment for reliable long-term operation.

**Covers:**
- Docker and Docker Compose setup
- Volume mounting for persistence
- Environment variable management
- Day-to-day Docker operations
- Backup strategies

---

### 🔌 Channel Integration Tutorials

#### [Telegram Integration](./openclaw-telegram.md)
**Duration:** 20 minutes | **Level:** Beginner

Connect OpenClaw to Telegram for instant messaging access.

**Covers:**
- Creating a bot via BotFather
- Configuring the Telegram channel
- Restricting access to authorized users
- Handling bot commands

---

#### [Discord Integration](./openclaw-discord.md)
**Duration:** 30 minutes | **Level:** Beginner

Add OpenClaw to your Discord server as an AI assistant.

**Covers:**
- Creating a Discord Application
- Bot setup and permissions
- OAuth2 invitation
- Channel configuration

---

#### [WhatsApp Integration](./openclaw-whatsapp.md)
**Duration:** 20 minutes | **Level:** Beginner

Connect OpenClaw to WhatsApp using the WhatsApp Web protocol.

**Covers:**
- QR code pairing
- Session management
- Automation considerations

---

### 🔧 Configuration & Advanced

#### [OpenClaw Skills System](./openclaw-skills.md)
**Duration:** 1 hour | **Level:** Intermediate

Extend OpenClaw with community Skills and create custom ones.

**Covers:**
- Finding and installing community Skills
- Skill manifest format
- Creating custom Skills
- Tool definitions

---

#### [Security Best Practices](./openclaw-security.md)
**Duration:** 45 minutes | **Level:** Intermediate

Secure your OpenClaw deployment against common vulnerabilities.

**Covers:**
- Understanding the "lethal trifecta" risk
- Channel access control
- API key management
- Gateway protection
- Regular security updates

---

### 🚀 Deployment Tutorials

#### [Production VPS Deployment](./openclaw-vps-deployment.md)
**Duration:** 1 hour | **Level:** Advanced

Deploy OpenClaw on a cloud VPS for 24/7 availability.

**Covers:**
- VPS provider selection
- Server setup and hardening
- Docker installation
- Firewall configuration
- Reverse proxy with HTTPS
- Automated backups

---

#### [OpenClaw + Ollama Setup](./openclaw-ollama.md)
**Duration:** 45 minutes | **Level:** Intermediate

Run OpenClaw with local AI models for complete privacy.

**Covers:**
- Installing Ollama
- Downloading local models
- Configuring OpenClaw for Ollama
- Performance optimization
- Privacy benefits

---

## 🎯 Quick Start

```bash
# 1. Install OpenClaw (macOS)
brew install openclaw

# 2. Verify installation
openclaw doctor

# 3. Configure your API key
# Edit config/openclaw.json5

# 4. Start OpenClaw
openclaw start

# 5. Access dashboard
# Open http://localhost:18789
```

---

## 🔗 Useful Links

- [OpenClaw GitHub](https://github.com/openclaw/openclaw)
- [Official Documentation](https://docs.openclaw.com)
- [Skills Library](https://github.com/openclaw/skills)
- [Community Forum](https://github.com/openclaw/openclaw/discussions)

---

## 📚 Learning Path

```
┌─────────────────────────────────────────────────────────┐
│                    OpenClaw Learning Path               │
└─────────────────────────────────────────────────────────┘

    ┌──────────────────┐
    │  1. Full Tutorial │ ◄── Start here
    │  (openclaw-full)  │
    └────────┬─────────┘
             │
    ┌────────┴─────────┐
    ▼                 ▼
┌────────────┐   ┌──────────────┐
│ 2. Docker  │   │ 2. Channel   │
│ Deployment │   │ Integration  │
└────────────┘   └──────────────┘
    │                 │
    └────────┬────────┘
             │
    ┌────────┴─────────┐
    ▼                 ▼
┌──────────────┐  ┌─────────────┐
│ 3. VPS      │  │ 3. Skills   │
│ Production  │  │ Development │
└──────────────┘  └─────────────┘
             │
    ┌────────┴─────────┐
    ▼                 ▼
┌─────────────────────────────┐
│     4. Advanced Topics      │
│  • Ollama Local Models      │
│  • Custom Skill Development │
│  • Security Hardening       │
└─────────────────────────────┘
```

---

*Last updated: March 2026*
