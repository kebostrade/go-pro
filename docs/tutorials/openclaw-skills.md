# OpenClaw Skills System

Skills extend OpenClaw's capabilities with pre-built automation packages.

## Finding Skills

Browse the official Skills library:
- [OpenClaw Skills Repository](https://github.com/openclaw/skills)
- Community-contributed Skills

## Installing a Skill

```bash
git clone https://github.com/openclaw/skill-example.git ./skills/example
```

Add to your configuration:

```json5
{
  skills: {
    example: {
      enabled: true,
      config: {}
    }
  }
}
```

## Creating Custom Skills

### Skill Manifest (skill.json5)

```json5
{
  "name": "my-skill",
  "version": "1.0.0",
  "description": "My custom skill",
  "tools": [
    {
      "name": "do_something",
      "description": "Performs an action",
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

### Tool Implementation

Create `tools/do_something.ts`:

```typescript
export const do_something = async ({ input }: { input: string }) => {
  // Your logic here
  return { result: `Processed: ${input}` };
};
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
