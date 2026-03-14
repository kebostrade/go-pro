# OpenClaw Production VPS Deployment

Deploy OpenClaw on a cloud VPS for 24/7 availability.

## Choose a VPS Provider

Recommended options:
- DigitalOcean (easy setup)
- Linode (good value)
- AWS Lightsail (AWS integration)
- Hetzner (best price/performance)

Minimum specs: 2GB RAM, 1 CPU, 20GB SSD

## Server Setup

### 1. Initialize Server

```bash
ssh root@your-server-ip
apt update && apt upgrade -y
```

### 2. Install Docker

```bash
curl -fsSL https://get.docker.com | sh
usermod -aG docker $USER
```

### 3. Configure Firewall

```bash
ufw allow 22/tcp   # SSH
ufw allow 18789/tcp  # OpenClaw (limit to your IP in production)
ufw enable
```

### 4. Deploy OpenClaw

```bash
mkdir -p ~/openclaw && cd ~/openclaw
# Create docker-compose.yml and config files
docker compose up -d
```

## Reverse Proxy with HTTPS

### Install Nginx and Certbot

```bash
apt install nginx certbot python3-certbot-nginx
```

### Nginx Configuration

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

### Get SSL Certificate

```bash
certbot --nginx -d openclaw.yourdomain.com
```

## Automated Backups

```bash
# Add to crontab
0 2 * * * tar -czf /backups/openclaw-$(date +\%Y\%m%d).tar.gz ~/openclaw/config ~/openclaw/data ~/openclaw/.env
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
