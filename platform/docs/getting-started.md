# AI Platform CLI (`aictl`)

A CLI tool for deploying and managing AI agents on the self-service Kubernetes platform.

## Installation

```bash
# From source
go install github.com/your-org/go-pro/cmd/aictl@latest

# Download binary
curl -sL https://github.com/your-org/go-pro/releases/latest/download/aictl-linux-amd64 -o aictl
chmod +x aictl
sudo mv aictl /usr/local/bin/
```

## Quick Start

```bash
# Login to platform
aictl login --url https://platform.example.com

# Create a new agent
aictl agent create my-agent --type react --llm gpt-4

# Deploy agent
aictl agent deploy my-agent --env dev

# Check status
aictl agent status my-agent

# View logs
aictl agent logs my-agent --follow

# Scale agent
aictl agent scale my-agent --replicas 3
```

## Commands

### Authentication

```bash
aictl login [--url URL] [--token TOKEN]
aictl logout
aictl whoami
```

### Agent Management

```bash
# List agents
aictl agent list [--env ENV]

# Create agent from template
aictl agent create NAME --type TYPE [--llm MODEL]

# Deploy agent
aictl agent deploy NAME [--env ENV] [--values FILE]

# Update agent
aictl agent update NAME [--image IMAGE] [--values FILE]

# Delete agent
aictl agent delete NAME [--force]

# Get agent details
aictl agent get NAME [-o yaml|json]

# Scale agent
aictl agent scale NAME --replicas N

# Restart agent
aictl agent restart NAME
```

### Logs & Debugging

```bash
# Stream logs
aictl agent logs NAME [--follow] [--tail N]

# Execute command in pod
aictl agent exec NAME -- COMMAND

# Port forward
aictl agent port-forward NAME --local-port 8080

# Describe agent (K8s resources)
aictl agent describe NAME
```

### Templates

```bash
# List available templates
aictl template list

# Show template details
aictl template get NAME

# Create from template
aictl agent create NAME --from-template TEMPLATE
```

### Environment

```bash
# List environments
aictl env list

# Switch environment
aictl env use ENV

# Show current environment
aictl env current
```

### Secrets

```bash
# List secrets
aictl secret list

# Create secret
aictl secret create NAME --from-literal KEY=VALUE

# Update secret
aictl secret update NAME --from-literal KEY=VALUE

# Delete secret
aictl secret delete NAME
```

## Configuration

Config file: `~/.aictl/config.yaml`

```yaml
current_context: dev
contexts:
  dev:
    url: https://dev.platform.example.com
    token: ${AICTL_TOKEN}
  prod:
    url: https://platform.example.com
    token: ${AICTL_PROD_TOKEN}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AICTL_URL` | Platform URL |
| `AICTL_TOKEN` | Auth token |
| `AICTL_CONTEXT` | Current context |
| `AICTL_NAMESPACE` | Default namespace |

## Examples

### Deploy a ReAct Agent

```bash
# Create agent config
cat > agent-values.yaml << EOF
agent:
  type: react
  llmProvider: openai
  llmModel: gpt-4
  tools:
    - name: search
      type: web_search
    - name: calculator
      type: builtin
  memory:
    enabled: true
    type: redis
EOF

# Deploy
aictl agent deploy my-react-agent --values agent-values.yaml
```

### Deploy with Custom Tools

```bash
aictl agent create custom-agent \
  --type base \
  --llm gpt-4-turbo \
  --tools search,weather,stock

aictl agent deploy custom-agent --env staging
```

### Monitor Agent

```bash
# Watch logs
aictl agent logs my-agent --follow --tail 100

# Check metrics
aictl agent metrics my-agent

# View events
aictl agent events my-agent
```

## Shell Completion

```bash
# Bash
source <(aictl completion bash)

# Zsh
source <(aictl completion zsh)

# Fish
aictl completion fish | source
```
