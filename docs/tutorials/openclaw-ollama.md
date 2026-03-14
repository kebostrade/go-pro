# OpenClaw + Ollama: Local AI Models

Run OpenClaw entirely locally with Ollama for complete privacy.

## Why Ollama?

- **Privacy**: Your data never leaves your machine
- **No API costs**: Pay once for hardware
- **Offline capability**: Works without internet
- **Custom models**: Fine-tune for your needs

## Install Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh
```

## Download Models

```bash
# Lightweight model (fast, less capable)
ollama pull llama3.2:1b

# Balanced model (recommended)
ollama pull llama3.2:3b

# Powerful model (requires more RAM)
ollama pull llama3.1:8b
```

Check available RAM:
```bash
# macOS
system_profiler SPHardwareDataType | grep Memory

# Linux
free -h
```

## Configure OpenClaw

Update `config/openclaw.json5`:

```json5
{
  providers: {
    ollama: {
      enabled: true,
      defaultModel: "llama3.2:3b",
      baseURL: "http://localhost:11434"
    }
  },
  gateway: {
    port: 18789,
    host: "0.0.0.0"
  },
  channels: {}
}
```

## Performance Tips

1. **Keep model loaded**: Ollama loads model into memory on first request
2. **Use appropriate model**: Smaller models for faster responses
3. **Sufficient RAM**: Ensure enough memory for model + OpenClaw
4. **GPU acceleration**: Install CUDA for NVIDIA GPUs

## Verify Setup

```bash
# Test Ollama
ollama list

# Test OpenClaw
openclaw doctor
```

---

*See [openclaw-full-tutorial.md](./openclaw-full-tutorial.md) for the complete guide.*
