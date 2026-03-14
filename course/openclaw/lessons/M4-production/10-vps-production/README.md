# Lesson OC-10: VPS Production Deployment

## Overview

This lesson covers deploying OpenClaw on a cloud VPS with HTTPS.

## Choose a VPS Provider

| Provider | Good For |
|----------|----------|
| DigitalOcean | Easy setup, good docs |
| Linode | Value, reliable |
| AWS Lightsail | AWS integration |
| Hetzner | Best price/performance |

**Minimum specs**: 2GB RAM, 1 CPU, 20GB SSD

## Server Setup

### 1. SSH to Server

```bash
ssh root@your-server-ip
```

### 2. Install Docker

```bash
apt update && apt upgrade -y
curl -fsSL https://get.docker.com | sh
usermod -aG docker $USER
```

### 3. Configure Firewall

```bash
ufw allow 22/tcp   # SSH
ufw allow 18789/tcp # OpenClaw (temporary)
ufw enable
```

## Deploy OpenClaw

```bash
mkdir -p ~/openclaw && cd ~/openclaw

# Create docker-compose.yml and .env
# (see Lesson 9)

docker compose up -d
```

## Reverse Proxy with Nginx

### Install Nginx & Certbot

```bash
apt install nginx certbot python3-certbot-nginx
```

### Nginx Config

```nginx
# /etc/nginx/sites-available/openclaw
server {
    listen 443 ssl http2;
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

### Enable and Get SSL

```bash
ln -s /etc/nginx/sites-available/openclaw /etc/nginx/sites-enabled/
nginx -t
certbot --nginx -d openclaw.yourdomain.com
```

### Update Firewall

```bash
ufw delete allow 18789/tcp  # Remove direct access
ufw allow 443/tcp  # HTTPS
```

## Automated Backups

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * tar -czf /backups/openclaw-$(date +\%Y\%m%d).tar.gz ~/openclaw/config ~/openclaw/data ~/openclaw/.env
```

---

**Next**: [Lesson 11: Security](11-security-hardening/README.md)
