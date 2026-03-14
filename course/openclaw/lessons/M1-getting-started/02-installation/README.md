# Lesson OC-02: Installation

## Overview

This lesson covers installing OpenClaw on macOS, Linux, and Windows (WSL2).

## Prerequisites

- API key from OpenAI or Anthropic
- Node.js 18+ or Docker installed
- Git installed

## Installation Methods

### Method 1: Homebrew (macOS Recommended)

```bash
# Install OpenClaw
brew install openclaw

# Verify installation
openclaw doctor
```

### Method 2: Docker (Cross-Platform Recommended)

```bash
# Create project directory
mkdir -p ~/openclaw && cd ~/openclaw

# Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
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
EOF

# Create .env file
cat > .env << 'EOF'
ANTHROPIC_API_KEY=sk-ant-your-key-here
EOF

# Create config directory
mkdir -p config

# Start OpenClaw
docker compose up -d

# View logs
docker compose logs -f
```

### Method 3: From Source (Linux/WSL2)

```bash
# Clone repository
git clone https://github.com/openclaw/openclaw.git
cd openclaw

# Install dependencies
npm install
# OR use Bun for faster installs
bun install

# Copy example config
cp config.example.json config.json

# Edit config with your API key
nano config.json

# Start OpenClaw
npm run start
```

## Verification

Run the doctor command to verify your setup:

```bash
openclaw doctor
```

Expected output:
```
✅ OpenClaw installed successfully
✅ Gateway running on http://localhost:18789
✅ Configuration valid
```

## Access the Dashboard

Open your browser to:
```
http://localhost:18789
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Port already in use | Change port in config or stop conflicting service |
| API key error | Verify key in .env file, no extra quotes/spaces |
| Permission denied | Check file permissions on config directory |

---

**Next**: [Lesson 3: Configuration](03-configuration/README.md)
