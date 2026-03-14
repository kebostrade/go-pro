# Lesson OC-01: What is OpenClaw? Architecture Deep Dive

## Overview

This lesson introduces OpenClaw, the fastest-growing open-source AI agent platform, and explains its architecture.

## What You'll Learn

- What makes OpenClaw different from other AI assistants
- The three core components: Gateway, Nodes, Skills
- Understanding the agent loop and tool execution

## OpenClaw: The Big Picture

OpenClaw is a self-hosted AI agent platform that gives you:
- **Privacy**: Your conversations stay on your infrastructure
- **Control**: Full customization of behavior and capabilities
- **Multi-channel**: Connect to Telegram, Discord, WhatsApp, Slack, and more
- **Automation**: Scheduled tasks, webhooks, and custom Skills

### Why Self-Hosted?

| Cloud AI Assistants | OpenClaw |
|-------------------|----------|
| Data leaves your machine | Data stays local |
| Limited customization | Full control |
| Subscription costs | One-time infrastructure cost |
| Dependent on uptime | Run on your hardware |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        OpenClaw                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌────────────┐  │
│  │   Gateway    │────▶│    Agent     │────▶│   Skills   │  │
│  │  (Channels)  │     │   (Brain)    │     │  (Tools)   │  │
│  └──────────────┘     └──────────────┘     └────────────┘  │
│         │                    │                    │          │
│         ▼                    ▼                    ▼          │
│  ┌──────────────┐     ┌──────────────┐     ┌────────────┐  │
│  │ Telegram     │     │   Memory     │     │  HTTP      │  │
│  │ Discord      │     │   Session    │     │  Files     │  │
│  │ WhatsApp     │     │   Context    │     │  Custom    │  │
│  └──────────────┘     └──────────────┘     └────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1. Gateway

The Gateway is the entry point for all communications:
- Handles incoming messages from all channels
- Routes outgoing messages back to users
- Manages authentication and access control
- Runs on port 18789 by default

### 2. Agent (The Brain)

The Agent orchestrates the conversation:
- Maintains conversation history and context
- Decides which tools to use
- Manages the reasoning loop
- Integrates with LLM providers

### 3. Skills (The Tools)

Skills provide capabilities:
- Pre-built integrations (web search, code execution)
- Custom automation workflows
- File system operations
- API interactions

## The Agent Loop

When a user sends a message:

```
User Message
    │
    ▼
┌─────────────┐
│   Gateway   │ ─── Receive message
└─────────────┘
    │
    ▼
┌─────────────┐
│    Agent    │ ─── 1. Add to context
└─────────────┘    │
    │              ▼
    │         ┌─────────────┐
    │         │  LLM Call   │ ─── "What should I do?"
    │         └─────────────┘
    │              │
    │              ▼ (LLM decides to use a tool)
    │         ┌─────────────┐
    │         │    Tool     │ ─── Execute action
    │         │   Execution │
    │         └─────────────┘
    │              │
    │              ▼
    │         ┌─────────────┐
    │         │  LLM Call   │ ─── "Now respond to user"
    │         └─────────────┘
    │
    ▼
┌─────────────┐
│   Gateway   │ ─── Send response to user
└─────────────┘
```

## Understanding Skills

Skills are what make OpenClaw powerful:

| Skill | Capability |
|-------|------------|
| `web-search` | Search the internet |
| `fetch` | HTTP requests |
| `bash` | Run shell commands |
| `memory` | Persistent storage |
| `scheduler` | Cron-like scheduling |

## Quiz

1. What are the three main components of OpenClaw?
2. What port does the Gateway run on by default?
3. Why might someone choose OpenClaw over ChatGPT?

---

**Next**: [Lesson 2: Installation](02-installation/README.md)
