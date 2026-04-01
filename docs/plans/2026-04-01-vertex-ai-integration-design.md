# Vertex AI Integration Design

**Date:** 2026-04-01
**Status:** Approved

## Overview

Integrate Google Cloud Vertex AI as the primary LLM provider for the AI Agent Platform, with OpenAI as a fallback. This provides a production-grade, enterprise-ready AI infrastructure on GCP.

## Architecture

```
services/ai-agent-platform/internal/llm/
├── provider.go    (existing - unchanged)
├── openai.go      (existing - fallback)
└── vertex.go      (NEW - primary)
```

## Decisions

| Aspect | Decision |
|--------|----------|
| **Fallback order** | Vertex AI first → OpenAI fallback |
| **Model selection** | Runtime config via env var |
| **Auth** | ADC first → API key fallback |
| **Interface** | Implements `types.LLMProvider` |

## Provider Interface

```go
// VertexAIProvider implements LLMProvider interface
type VertexAIProvider struct {
    *BaseProvider
    client *vertexai.Client
    model  string
}
```

Methods:
- `Generate(ctx, request) -> response`
- `GenerateStream(ctx, request) -> stream`
- `GetModelInfo() -> modelInfo`
- `GetProviderName() -> "vertex"`
- `SupportsStreaming() -> true`
- `SupportsFunctionCalling() -> true`

## Environment Variables

```bash
# Vertex AI (primary)
VERTEX_PROJECT_ID=my-gcp-project
VERTEX_LOCATION=us-central1
VERTEX_MODEL=gemini-2.0-flash
VERTEX_API_KEY=optional-if-adc-unavailable

# OpenAI (fallback)
OPENAI_API_KEY=sk-...

# Provider order
LLM_FALLBACK_ORDER=vertex,openai
```

## Auth Flow

```
1. Check VERTEX_API_KEY
2. If set, use API key auth
3. If not, try Application Default Credentials (ADC)
   - GCP service account
   - Workload identity
   - `gcloud auth application-default login`
4. If neither available, provider init fails
```

## Fallback Flow

```
GenerateWithFallback(["vertex", "openai"], request)
  → Try vertex.Generate()
    → Success: return response
    → Failure: log warning, continue
  → Try openai.Generate()
    → Success: return response
    → Failure: return error
```

## Files to Create

| File | Description |
|------|-------------|
| `internal/llm/vertex.go` | Vertex AI provider implementation |

## Files to Modify

| File | Change |
|------|--------|
| `internal/llm/provider.go` | Register vertex provider |
| `cmd/coding-agent-server/main.go` | Wire vertex + fallback config |

## Vertex AI SDK

Use `cloud.google.com/go/vertexai` package:

```go
import (
    "cloud.google.com/go/vertexai/googlesAI"
    "context"
)
```

## Testing

- Mock LLM responses for unit tests
- Integration test with real Vertex AI (requires GCP credentials)
- Fallback test: mock vertex failure, verify openai called

## Out of Scope

- Gemini for Workspace (not needed)
- Vertex AI Agent Builder (use existing ReAct agent)
- Cloud Run deployment (handled separately)
