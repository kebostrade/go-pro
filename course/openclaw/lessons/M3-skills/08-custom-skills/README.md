# Lesson OC-08: Building Custom Skills

## Overview

This lesson covers creating your own OpenClaw Skills.

## Skill Structure

```
my-skill/
├── skill.json5      # Manifest
├── src/
│   └── index.ts     # Implementation
└── package.json
```

## Step 1: Create Manifest

```json5
// skill.json5
{
  "name": "weather",
  "version": "1.0.0",
  "description": "Get weather information",
  "tools": [
    {
      "name": "get_weather",
      "description": "Get current weather for a city",
      "parameters": {
        "type": "object",
        "properties": {
          "city": {
            "type": "string",
            "description": "City name"
          }
        },
        "required": ["city"]
      }
    }
  ]
}
```

## Step 2: Implement Tool

```typescript
// src/index.ts
interface ToolParams {
  city: string;
}

export const get_weather = async ({ city }: ToolParams) => {
  const response = await fetch(
    `https://api.weather.com/v3/wx/conditions/current?city=${city}&key=YOUR_API_KEY`
  );
  
  const data = await response.json();
  
  return {
    city: data.location,
    temperature: data.temp,
    conditions: data.conditions,
    humidity: data.humidity
  };
};
```

## Step 3: Register Tools

```typescript
// src/index.ts (continued)
export const tools = {
  get_weather
};
```

## Step 4: Install and Configure

```bash
# Copy to skills directory
cp -r ./weather ../skills/weather
```

```json5
{
  skills: {
    weather: {
      enabled: true,
      config: {
        apiKey: "YOUR_WEATHER_API_KEY"
      }
    }
  }
}
```

## Testing Your Skill

1. Restart OpenClaw
2. Ask: "What's the weather in Tokyo?"
3. Agent should use your tool automatically

---

**Next**: [Module 4: Production Deployment](M4-production/README.md)
