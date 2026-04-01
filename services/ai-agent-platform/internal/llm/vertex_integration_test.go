//go:build integration
// +build integration

package llm

import (
	"context"
	"testing"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
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
