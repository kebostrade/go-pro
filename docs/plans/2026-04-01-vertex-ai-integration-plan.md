# Vertex AI Integration Implementation Plan

> **For Claude:** Use superpowers:subagent-driven-development to implement this plan task-by-task.

**Goal:** Integrate Vertex AI as primary LLM provider with OpenAI fallback into the AI Agent Platform.

**Architecture:** Add VertexAIProvider implementing LLMProvider interface, register in ProviderManager with configurable fallback order. Uses Google Cloud SDK for Vertex AI.

**Tech Stack:** Go, cloud.google.com/go/vertexai, Google Cloud SDK

---

## Task 1: Create Vertex AI Provider Skeleton

**Files:**
- Create: `services/ai-agent-platform/internal/llm/vertex.go`
- Test: `services/ai-agent-platform/internal/llm/vertex_test.go`

**Step 1: Create the test file**

```go
package llm

import (
    "context"
    "testing"
)

func TestVertexAIProvider_GetProviderName(t *testing.T) {
    // TODO: create test after we have mock setup
}
```

**Step 2: Create the vertex.go skeleton**

```go
package llm

import (
    "context"
    "fmt"
    "time"

    "github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
    "github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// VertexAIProvider implements LLMProvider for Google Vertex AI
type VertexAIProvider struct {
    *BaseProvider
    projectID string
    location  string
    model     string
}

// VertexAIConfig holds Vertex AI configuration
type VertexAIConfig struct {
    ProjectID string
    Location  string
    Model     string
    APIKey    string // optional, uses ADC if not set
    Timeout   time.Duration
}

// NewVertexAIProvider creates a new Vertex AI provider
func NewVertexAIProvider(config VertexAIConfig) (*VertexAIProvider, error) {
    if config.ProjectID == "" {
        return nil, errors.NewInvalidConfigError("Vertex AI project ID is required")
    }
    if config.Location == "" {
        config.Location = "us-central1"
    }
    if config.Model == "" {
        config.Model = "gemini-2.0-flash"
    }
    if config.Timeout == 0 {
        config.Timeout = 60 * time.Second
    }
    
    return &VertexAIProvider{
        BaseProvider: NewBaseProvider(ProviderConfig{
            Provider: "vertex",
            Timeout:  config.Timeout,
            MaxRetries: 3,
            RetryDelay: 2 * time.Second,
        }),
        projectID: config.ProjectID,
        location:  config.Location,
        model:     config.Model,
    }, nil
}
```

**Step 3: Add interface methods (empty implementations)**

```go
// Generate generates a completion for the given request
func (p *VertexAIProvider) Generate(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, error) {
    return nil, fmt.Errorf("not implemented")
}

// GenerateStream generates a completion and streams the response
func (p *VertexAIProvider) GenerateStream(ctx context.Context, request types.LLMRequest) (<-chan types.LLMStreamChunk, error) {
    return nil, fmt.Errorf("not implemented")
}

// GetModelInfo returns information about the model
func (p *VertexAIProvider) GetModelInfo() types.ModelInfo {
    return types.ModelInfo{}
}

// GetProviderName returns the name of the provider
func (p *VertexAIProvider) GetProviderName() string {
    return "vertex"
}

// SupportsStreaming returns whether the provider supports streaming
func (p *VertexAIProvider) SupportsStreaming() bool {
    return true
}

// SupportsFunctionCalling returns whether the provider supports function calling
func (p *VertexAIProvider) SupportsFunctionCalling() bool {
    return true
}
```

**Step 4: Verify it compiles**

Run: `cd services/ai-agent-platform && go build ./internal/llm/...`
Expected: Builds successfully

**Step 5: Commit**

```bash
git add services/ai-agent-platform/internal/llm/vertex.go
git commit -m "feat(llm): add vertex ai provider skeleton"
```

---

## Task 2: Implement Generate Method

**Files:**
- Modify: `services/ai-agent-platform/internal/llm/vertex.go`

**Step 1: Add Vertex AI SDK imports and client**

```go
import (
    "context"
    "fmt"
    "time"

    "cloud.google.com/go/vertexai/googlesAI"
    "github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
    "github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)
```

**Step 2: Add client field and init method**

```go
type VertexAIProvider struct {
    *BaseProvider
    projectID string
    location  string
    model     string
    apiKey    string
    client   *googlesAI.Model
}
```

**Step 3: Update NewVertexAIProvider to init client**

```go
func NewVertexAIProvider(config VertexAIConfig) (*VertexAIProvider, error) {
    // ... validation ...
    
    provider := &VertexAIProvider{
        BaseProvider: NewBaseProvider(ProviderConfig{
            Provider:   "vertex",
            Timeout:    config.Timeout,
            MaxRetries: 3,
            RetryDelay: 2 * time.Second,
        }),
        projectID: config.ProjectID,
        location:  config.Location,
        model:     config.Model,
        apiKey:    config.APIKey,
    }
    
    // Initialize client
    client, err := googlesAI.NewModel(ctx, config.ProjectID, config.Location, config.Model)
    if err != nil {
        return nil, fmt.Errorf("failed to create vertex client: %w", err)
    }
    provider.client = client
    
    return provider, nil
}
```

**Step 4: Implement Generate method**

```go
func (p *VertexAIProvider) Generate(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, error) {
    if err := ValidateRequest(request); err != nil {
        return nil, err
    }

    startTime := time.Now()

    // Convert messages to Vertex format
    vertexMessages := make([]googlesAI.Message, len(request.Messages))
    for i, msg := range request.Messages {
        role := "user"
        if msg.Role == types.RoleSystem {
            role = "system"
        } else if msg.Role == types.RoleAssistant {
            role = "model"
        }
        vertexMessages[i] = googlesAI.Message{
            Role:  role,
            Parts: []googlesAI.Part{{Text: msg.Content}},
        }
    }

    // Build GenerateContentRequest
    genRequest := googlesAI.GenerateContentRequest{
        Messages: vertexMessages,
    }

    // Execute with retry
    var resp *googlesAI.GenerateContentResponse
    err := p.WithRetry(ctx, func() error {
        var err error
        resp, err = p.client.GenerateContent(ctx, genRequest)
        return err
    })
    if err != nil {
        return nil, errors.NewLLMProviderError("vertex", err)
    }

    // Parse response
    if len(resp.Candidates) == 0 {
        return nil, errors.NewLLMProviderError("vertex", fmt.Errorf("no candidates in response"))
    }

    candidate := resp.Candidates[0]
    content := ""
    if len(candidate.Content.Parts) > 0 {
        content = candidate.Content.Parts[0].Text
    }

    return &types.LLMResponse{
        Content:      content,
        FinishReason: string(candidate.FinishReason),
        Usage: types.TokenUsage{
            PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
            CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
            TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
        },
        Model:   p.model,
        Latency: time.Since(startTime),
    }, nil
}
```

**Step 5: Verify it compiles**

Run: `cd services/ai-agent-platform && go build ./internal/llm/...`
Expected: Builds successfully

**Step 6: Commit**

```bash
git add services/ai-agent-platform/internal/llm/vertex.go
git commit -m "feat(llm): implement vertex ai Generate method"
```

---

## Task 3: Add Provider Registration

**Files:**
- Modify: `services/ai-agent-platform/internal/llm/provider.go`
- Modify: `services/ai-agent-platform/cmd/coding-agent-server/main.go`

**Step 1: Add environment variable reading to provider.go**

Add to `ProviderConfig` struct or create a helper:

```go
// GetVertexConfigFromEnv reads Vertex AI config from environment
func GetVertexConfigFromEnv() (VertexAIConfig, bool) {
    projectID := os.Getenv("VERTEX_PROJECT_ID")
    if projectID == "" {
        return VertexAIConfig{}, false
    }
    
    return VertexAIConfig{
        ProjectID: projectID,
        Location:  getEnvOr("VERTEX_LOCATION", "us-central1"),
        Model:     getEnvOr("VERTEX_MODEL", "gemini-2.0-flash"),
        APIKey:    os.Getenv("VERTEX_API_KEY"),
        Timeout:   60 * time.Second,
    }, true
}

func getEnvOr(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}
```

**Step 2: Add initialization in main.go**

```go
// Add Vertex AI provider if configured
if vertexConfig, ok := llm.GetVertexConfigFromEnv(); ok {
    vertexProvider, err := llm.NewVertexAIProvider(vertexConfig)
    if err != nil {
        log.Printf("Warning: failed to create Vertex AI provider: %v", err)
    } else {
        providerManager.Register("vertex", vertexProvider)
        log.Println("Registered Vertex AI provider")
    }
}
```

**Step 3: Verify it compiles**

Run: `cd services/ai-agent-platform && go build ./...`
Expected: Builds successfully

**Step 4: Commit**

```bash
git add services/ai-agent-platform/internal/llm/provider.go services/ai-agent-platform/cmd/coding-agent-server/main.go
git commit -m "feat(llm): register vertex ai provider"
```

---

## Task 4: Add Integration Test (Manual Verification)

**Files:**
- Create: `services/ai-agent-platform/internal/llm/vertex_integration_test.go`

**Step 1: Create integration test skeleton**

```go
//go:build integration
// +build integration

package llm

import (
    "context"
    "testing"
)

func TestVertexAIProvider_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    ctx := context.Background()
    
    config, ok := GetVertexConfigFromEnv()
    if !ok {
        t.Skip("VERTEX_PROJECT_ID not set")
    }
    
    provider, err := NewVertexAIProvider(config)
    if err != nil {
        t.Fatalf("Failed to create provider: %v", err)
    }
    
    resp, err := provider.Generate(ctx, types.LLMRequest{
        Messages: []types.Message{
            {Role: types.RoleUser, Content: "Say 'Hello, Vertex!'"},
        },
        MaxTokens: 100,
    })
    if err != nil {
        t.Fatalf("Generate failed: %v", err)
    }
    
    if resp.Content == "" {
        t.Error("Expected non-empty content")
    }
    
    t.Logf("Response: %s", resp.Content)
}
```

**Step 2: Run short test**

Run: `cd services/ai-agent-platform && go test ./internal/llm/... -short`
Expected: Test skipped (integration test)

**Step 3: Commit**

```bash
git add services/ai-agent-platform/internal/llm/vertex_integration_test.go
git commit -m "test(llm): add vertex ai integration test"
```

---

## Verification Checklist

- [ ] `go build ./...` passes for entire module
- [ ] `go test ./internal/llm/... -short` passes
- [ ] Vertex AI provider registered when VERTEX_PROJECT_ID set
- [ ] Falls back to OpenAI when VERTEX not configured
- [ ] Environment variables documented in README or .env.example
