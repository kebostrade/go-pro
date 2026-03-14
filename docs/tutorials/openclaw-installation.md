# OpenClaw Installation Guide

This tutorial covers installing OpenClaw on macOS, Linux, and Windows via WSL2.

## Prerequisites

- API key from OpenAI, Anthropic, or Ollama
- 4GB RAM minimum (8GB recommended)
- Git installed
- Node.js 20+ or Bun 1.1+

## Installation Methods

### macOS (Homebrew)

```bash
brew install openclaw
openclaw doctor
```

### Linux / Windows WSL2

```bash
git clone https://github.com/openclaw/openclaw.git
cd openclaw
npm install
cp config.example.json config.json
npm run start
```

### Docker

```bash
mkdir -p ~/openclaw-docker && cd ~/openclaw-docker
# See docker-compose.yml in main tutorial
```

## Verification

Run `openclaw doctor` to verify your installation.

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
