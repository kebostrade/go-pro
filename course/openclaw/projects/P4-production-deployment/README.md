# Project 4: Production Deployment

**Difficulty**: Advanced  
**Duration**: 6-8 hours  
**Prerequisites**: Lessons 9-11

## Project Overview

Deploy a production-ready OpenClaw instance on a VPS with HTTPS.

## Requirements

1. **VPS Setup**:
   - Choose provider (DO/Linode/Hetzner)
   - Configure firewall
   - Secure SSH

2. **OpenClaw Deployment**:
   - Docker Compose setup
   - Proper volume management
   - Automated backups

3. **Security**:
   - Reverse proxy with Nginx
   - SSL/TLS via Let's Encrypt
   - Channel access control

4. **Monitoring**:
   - Log aggregation
   - Health checks
   - Backup verification

## Deliverables

- Working production URL
- SSL certificate
- Backup script
- Security checklist

## Architecture

```
User → DNS → Nginx (SSL) → OpenClaw (Docker) → LLM API
```
