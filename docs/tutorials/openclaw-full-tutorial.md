# OpenClaw Full Tutorial

OpenClaw is an open-source AI agent platform that has exploded in popularity in early 2026, reaching over 150,000 GitHub stars. It is a self-hosted personal AI assistant that runs entirely on your local machine and can be accessed through messaging apps you already use—including WhatsApp, Telegram, Discord, Slack, Signal, iMessage, Teams, and Google Chat. Unlike cloud-based AI assistants, OpenClaw gives you complete control over your data while enabling powerful automation capabilities.

This comprehensive tutorial covers everything you need to go from zero to a fully functional OpenClaw deployment, whether you are a beginner or an experienced developer.

---

## Table of Contents

1. [Understanding OpenClaw](#1-understanding-openclaw)
2. [Prerequisites](#2-prerequisites)
3. [Installation Methods](#3-installation-methods)
4. [Initial Configuration](#4-initial-configuration)
5. [Connecting Messaging Channels](#5-connecting-messaging-channels)
6. [Understanding and Using Skills](#6-understanding-and-using-skills)
7. [Security Best Practices](#7-security-best-practices)
8. [Docker Deployment](#8-docker-deployment)
9. [Production Deployment on a VPS](#9-production-deployment-on-a-vps)
10. [Troubleshooting Common Issues](#10-troubleshooting-common-issues)

---

## 1. Understanding OpenClaw

### 1.1 What Makes OpenClaw Different

OpenClaw distinguishes itself from other AI agent frameworks in several important ways. First, it operates entirely locally, meaning your conversations and data never leave your infrastructure unless you explicitly choose to share them. Second, its multi-channel support means you can interact with your AI assistant through whichever messaging platform you prefer, without being locked into a single interface. Third, OpenClaw's Skills system allows you to extend functionality through a powerful plugin architecture that can automate complex workflows.

The architecture consists of three main components working together. The Gateway serves as the central hub that handles incoming and outgoing communications across all connected channels. The Node system enables you to connect additional devices and services through a pairing mechanism, creating a distributed network of capabilities. The Skills framework provides pre-built automation packages that your agent can use to perform tasks ranging from simple information retrieval to complex multi-step workflows.

### 1.2 Core Features

OpenClaw provides an impressive array of features that make it a versatile AI assistant platform. The heartbeat scheduler enables 24/7 monitoring and automated responses based on time or event triggers. Cron job integration allows you to schedule recurring tasks with precision. Memory management ensures your agent maintains context across conversations while giving you control over what information is retained. The tool execution system lets your AI agent interact with external services, APIs, and your local filesystem to accomplish real work.

---

## 2. Prerequisites

Before installing OpenClaw, ensure your system meets the following requirements. You will need an API key from at least one supported LLM provider—OpenAI, Anthropic, or local models via Ollama. For hardware, a minimum of 4GB RAM is required, though 8GB is recommended for optimal performance. You must have Git installed and configured on your system. Finally, Node.js version 20 or higher is required, though Bun 1.1+ is also supported as an alternative runtime.

### 2.1 Supported LLM Providers

OpenClaw supports multiple LLM backends, giving you flexibility in your AI provider selection. OpenAI's GPT-4o and GPT-4o Mini models work out of the box with an API key. Anthropic's Claude models, including the latest Sonnet and Opus versions, are fully supported. For privacy-conscious users, Ollama enables running local models entirely on your own hardware, though this requires a sufficiently powerful GPU for acceptable performance.

---

## 3. Installation Methods

### 3.1 Installation via Homebrew (macOS)

The fastest way to get OpenClaw running on macOS is through Homebrew. Open your terminal and run the following command to install OpenClaw:

```bash
brew install openclaw
```

After installation, verify that everything is working correctly by running the diagnosis command:

```bash
openclaw doctor
```

If you encounter any issues, the doctor command will provide specific guidance on what needs to be fixed.

### 3.2 Installation from Source

For users who prefer installing from source or who are on Linux or Windows (via WSL2), the manual installation process gives you more control over the setup. First, clone the official repository:

```bash
git clone https://github.com/openclaw/openclaw.git
cd openclaw
```

Next, install the dependencies using your preferred package manager:

```bash
npm install
```

For faster installs, you can use Bun instead:

```bash
bun install
```

After dependencies are installed, create your initial configuration by copying the example file:

```bash
cp config.example.json config.json
```

Finally, start OpenClaw in development mode:

```bash
npm run start
```

### 3.3 Installation via Docker

Docker installation provides the most consistent experience across different operating systems and simplifies deployment to production environments. First, create a dedicated directory for your OpenClaw deployment:

```bash
mkdir -p ~/openclaw-docker && cd ~/openclaw-docker
```

Create a `docker-compose.yml` file with the following configuration:

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
```

This configuration maps port 18789 for the Gateway, mounts configuration and data directories for persistence, and reads sensitive environment variables from a separate file. Adjust the timezone setting to match your location.

---

## 4. Initial Configuration

### 4.1 Environment Variables Setup

Create a `.env` file in your OpenClaw directory to store your API keys and sensitive configuration. The file should contain only the provider keys you intend to use:

```bash
# .env file
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxxx
```

Never commit this file to version control. Add it to your `.gitignore` if you are using Git. Each key corresponds to a specific AI provider, so only include the ones you actually plan to use.

### 4.2 Configuration File

Create the main configuration file at `config/openclaw.json5`. This file controls all aspects of your OpenClaw deployment:

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

  // Channel settings - enable as needed
  channels: {}
}
```

When running inside Docker, the host must be set to `0.0.0.0` to allow connections from outside the container. For local development outside Docker, `127.0.0.1` or `localhost` is more secure.

### 4.3 First Run Checklist

After completing your configuration, run through this checklist to verify your setup. First, start OpenClaw using your chosen installation method. Next, access the admin dashboard at `http://localhost:18789` to verify the Gateway is running. Run the diagnostic command to check for any issues:

```bash
openclaw doctor
```

Ensure your API keys are correctly configured by checking the provider status in the dashboard. Finally, verify that at least one channel is properly connected before proceeding.

---

## 5. Connecting Messaging Channels

### 5.1 Telegram Setup

Telegram is often the easiest channel to set up and serves as an excellent starting point. Begin by creating a new bot through BotFather, Telegram's official bot management interface. Open Telegram and search for @BotFather, then send the command `/newbot` to create a new bot. Follow the prompts to choose a name and username for your bot. BotFather will provide you with an HTTP API token—copy this token as you will need it later.

Now add the token to your OpenClaw configuration. Edit your `config/openclaw.json5` to include the Telegram channel:

```json5
{
  channels: {
    telegram: {
      enabled: true,
      botToken: "YOUR_TELEGRAM_BOT_TOKEN",
      chatIds: []
    }
  }
}
```

The `chatIds` array can remain empty initially, which allows any user who messages your bot to interact with OpenClaw. For production use, you should add specific chat IDs to restrict access to authorized users only.

Restart OpenClaw after making configuration changes. To find your chat ID, message your bot and then visit `https://api.telegram.org/bot<TOKEN>/getUpdates` in your browser. Your chat ID will appear in the response.

### 5.2 Discord Setup

Setting up Discord requires creating a Discord Application and adding it to your server. First, navigate to the Discord Developer Portal and create a new application. Under the Bot section, reset and copy your bot token. Enable the Message Content Intent under the Bot settings—this is required for OpenClaw to read messages.

Invite the bot to your server using the OAuth2 URL generator. In the Developer Portal, select the bot scope and the required permissions (at minimum, Send Messages, Read Message History, and Manage Messages). Copy the generated URL and open it in your browser to complete the authorization process.

Configure OpenClaw to use Discord:

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

### 5.3 WhatsApp Setup

WhatsApp requires a slightly more complex setup due to WhatsApp's restrictions on unofficial clients. OpenClaw supports pairing with an existing WhatsApp account through the WhatsApp Web protocol. In your configuration, enable the WhatsApp channel and restart OpenClaw. The dashboard will display a QR code that you scan with your WhatsApp app to establish the connection.

```json5
{
  channels: {
    whatsapp: {
      enabled: true
    }
  }
}
```

Note that WhatsApp may ban accounts that use unofficial clients, particularly if they are used for automated messaging. Use this channel judiciously and consider the risks before enabling automated responses.

### 5.4 Additional Channels

OpenClaw supports numerous other channels including Slack, Signal, iMessage, Microsoft Teams, and Google Chat. Each channel has its own specific setup process, but the general pattern remains consistent. You obtain necessary credentials or tokens from the respective platform's developer portal, add the configuration to your channel settings, and restart OpenClaw to activate the connection.

---

## 6. Understanding and Using Skills

### 6.1 What Are Skills

Skills are the building blocks of OpenClaw's extensibility. Each Skill is a pre-configured package that gives your AI agent specific capabilities—from searching the web to executing code, managing files, or interacting with APIs. The Skills system is what transforms OpenClaw from a simple chatbot into a powerful autonomous agent.

### 6.2 Installing Skills

OpenClaw maintains a library of community-contributed Skills that you can install and enable. Browse the official Skills repository to find capabilities that match your needs. To install a Skill, you typically clone its repository into your Skills directory and add its configuration to your main config file.

```bash
# Example: Installing a skill
git clone https://github.com/openclaw/skill-example.git ./skills/example
```

### 6.3 Creating Custom Skills

For advanced users, creating custom Skills allows you to tailor OpenClaw to your specific requirements. A Skill consists of a manifest file defining its capabilities and one or more implementation files that execute the actual logic. The manifest uses a standard format that OpenClaw's Skills system understands:

```json5
{
  "name": "my-custom-skill",
  "version": "1.0.0",
  "description": "A custom skill for specific tasks",
  "tools": [
    {
      "name": "do_something",
      "description": "Performs a specific action",
      "parameters": {
        "type": "object",
        "properties": {
          "input": {
            "type": "string",
            "description": "Input for the action"
          }
        },
        "required": ["input"]
      }
    }
  ]
}
```

---

## 7. Security Best Practices

### 7.1 The Security Reality

Security should be your primary concern when deploying OpenClaw. Recent security research has revealed concerning statistics: in January 2026, security researcher Maor Dayan found 42,665 exposed OpenClaw instances online, with 93.4% being vulnerable to exploitation. Cisco's analysis found that 26% of agent Skills contain vulnerabilities. Simon Willison's "lethal trifecta" concept explains why AI agents pose particular security risks—agents that access private data, are exposed to untrusted content, and can communicate externally are inherently dangerous without proper safeguards.

### 7.2 Essential Security Measures

Always restrict channel access to authorized users only. Rather than leaving channel configurations open to anyone, explicitly specify allowed chat IDs, user IDs, or server IDs in your configuration. Never expose your Gateway to the public internet without authentication—use a reverse proxy with authentication or ensure proper firewall rules are in place.

Keep your installation updated. OpenClaw releases security patches regularly, and running outdated versions may leave known vulnerabilities exposed. Subscribe to the GitHub security advisories for the OpenClaw repository to stay informed about critical updates.

### 7.3 API Key Management

Your API keys are the most sensitive component of your OpenClaw deployment. Never hardcode keys in configuration files that are committed to version control. Use environment variables or secrets management tools. For Docker deployments, use Docker secrets or external secret management services. Regularly rotate your API keys, especially if you suspect they may have been compromised.

---

## 8. Docker Deployment

### 8.1 Starting Your Container

With your configuration files prepared, bring up the OpenClaw service:

```bash
docker compose up -d
```

Verify the container is running correctly:

```bash
docker compose ps
```

You should see the openclaw container in a healthy state.

### 8.2 Viewing Logs

Logs provide crucial information for troubleshooting and monitoring. To watch logs in real time:

```bash
docker compose logs -f openclaw
```

To view only the most recent entries:

```bash
docker compose logs --tail 100 openclaw
```

If you see API connection errors during startup, double-check that your keys in the `.env` file are correct and properly formatted.

### 8.3 Day-to-Day Operations

Common maintenance tasks are straightforward with Docker. To stop the service, use `docker compose down`. To restart after configuration changes, run `docker compose restart`. To upgrade to a newer version, pull the latest image and recreate the container:

```bash
docker compose pull
docker compose up -d
```

Your configuration and data remain safe in the mounted volumes during upgrades. To access a shell inside the container for advanced troubleshooting:

```bash
docker exec -it openclaw sh
```

### 8.4 Backup Strategy

Because configuration and data are stored in mounted host directories, backups are straightforward. Create a backup script that archives your critical files:

```bash
tar -czf openclaw-backup-$(date +%Y%m%d).tar.gz config data .env
```

Schedule regular backups using cron and consider storing backups in a separate location from your server. Always test your backup restoration process before you actually need it.

---

## 9. Production Deployment on a VPS

### 9.1 Choosing a VPS Provider

For reliable 24/7 operation, a Virtual Private Server provides better uptime and accessibility than a home computer. Popular options include DigitalOcean, Linode, AWS Lightsail, and Hetzner. For OpenClaw, a server with 2GB RAM and a single CPU core provides adequate resources for moderate usage.

### 9.2 Server Setup

Start with a fresh Ubuntu 22.04 or 20.04 installation. Update system packages and install Docker:

```bash
sudo apt update && sudo apt upgrade -y
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

Logout and log back in for group permissions to take effect, or use `newgrp docker` to apply them immediately.

### 9.3 Firewall Configuration

Configure UFW to allow only necessary ports:

```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 18789/tcp # OpenClaw Gateway (restrict to your IP if using reverse proxy)
sudo ufw enable
```

For production, it is strongly recommended to place OpenClaw behind a reverse proxy with HTTPS termination.

### 9.4 Reverse Proxy with HTTPS

Nginx serves as an effective reverse proxy. Install Nginx and obtain SSL certificates using Certbot:

```bash
sudo apt install nginx certbot python3-certbot-nginx
```

Create an Nginx configuration file:

```nginx
server {
    listen 443 ssl;
    server_name openclaw.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:18789;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

The WebSocket upgrade headers are essential—OpenClaw's real-time communication depends on them. After creating the configuration, enable it:

```bash
sudo ln -s /etc/nginx/sites-available/openclaw /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 10. Troubleshooting Common Issues

### 10.1 API Connection Errors

If OpenClaw fails to connect to your AI provider, first verify your API key is correct and has not expired. Check that the key is properly formatted in your `.env` file without extra spaces or quotation marks. Ensure your account has available credits or quota. For Anthropic, verify that you have accepted the latest API terms.

### 10.2 Channel Connection Problems

Each channel has specific requirements that may cause connection failures. For Telegram, double-check your bot token and ensure the bot has not been banned. For Discord, confirm the Message Content Intent is enabled in the Developer Portal. For WhatsApp, your session may have been logged out—re-scan the QR code in the dashboard.

### 10.3 Gateway Inaccessible

If you cannot access the Gateway, verify the service is running with `docker compose ps`. Check that the port is correctly mapped in your docker-compose.yml and that any firewall allows incoming connections. When running inside Docker, ensure the gateway host is set to `0.0.0.0`.

### 10.4 Performance Issues

Slow responses may indicate network issues with your AI provider or insufficient server resources. Monitor your memory usage—if you are running low on RAM, consider upgrading your server or reducing concurrent session limits. For Ollama users, ensure your model is loaded into memory before expecting fast responses.

---

## Conclusion

OpenClaw represents a significant advancement in personal AI assistants, offering the power of modern AI agents with the privacy and control of self-hosting. This tutorial has covered the essential knowledge needed to install, configure, secure, and maintain an OpenClaw deployment. As you become more comfortable with the platform, explore the Skills system to extend your agent's capabilities and customize it to your specific needs.

For continued learning, consult the official OpenClaw documentation, engage with the community on GitHub discussions, and stay updated on security best practices as the platform evolves.

---

## Quick Reference Commands

```bash
# Installation verification
openclaw doctor

# Docker operations
docker compose up -d
docker compose logs -f
docker compose restart
docker compose pull && docker compose up -d

# Backup
tar -czf openclaw-backup-$(date +%Y%m%d).tar.gz config data .env

# Access container shell
docker exec -it openclaw sh
```

---

*Last updated: March 2026*
