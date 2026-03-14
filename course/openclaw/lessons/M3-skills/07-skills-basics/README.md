# Lesson OC-07: Skills System Basics

## Overview

This lesson covers finding, installing, and using OpenClaw Skills.

## What Are Skills?

Skills are packages that give your agent capabilities beyond simple chat:
- Web search and fetching
- File system operations
- Code execution
- API integrations
- Custom automation

## Finding Skills

Official Skills repository:
- [OpenClaw Skills](https://github.com/openclaw/skills)
- npm: `npm search @openclaw skill`

## Installing a Skill

```bash
# Clone skill to skills directory
git clone https://github.com/openclaw/skill-web-search.git ./skills/web-search
```

## Configuring Skills

```json5
{
  skills: {
    "web-search": {
      enabled: true,
      config: {
        // Skill-specific config
      }
    },
    "fetch": {
      enabled: true
    }
  }
}
```

## Built-in Skills

| Skill | Description |
|-------|-------------|
| `fetch` | Make HTTP requests |
| `memory` | Persistent storage |
| `scheduler` | Cron-like scheduling |
| `bash` | Run shell commands |

## Skill Manifest

Each skill has a `skill.json5`:

```json5
{
  "name": "my-skill",
  "version": "1.0.0",
  "description": "What this skill does",
  "tools": [
    {
      "name": "tool_name",
      "description": "What the tool does",
      "parameters": {
        "type": "object",
        "properties": {
          "param1": {
            "type": "string",
            "description": "Parameter description"
          }
        },
        "required": ["param1"]
      }
    }
  ]
}
```

---

**Next**: [Lesson 8: Custom Skills](08-custom-skills/README.md)
