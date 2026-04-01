package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"

	"cloud.google.com/go/vertexai/genai"
)

// VertexAIProvider implements LLMProvider for Google Vertex AI
type VertexAIProvider struct {
	*BaseProvider
	projectID string
	location  string
	model     string
	apiKey    string
	client    *genai.Client
	genModel  *genai.GenerativeModel
}

// VertexAIConfig holds Vertex AI configuration
type VertexAIConfig struct {
	ProjectID string
	Location  string
	Model     string
	APIKey    string
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
	ctx := context.Background()
	client, err := genai.NewClient(ctx, config.ProjectID, config.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to create vertex client: %w", err)
	}
	provider.client = client
	provider.genModel = client.GenerativeModel(config.Model)

	return provider, nil
}

// Generate generates a completion for the given request
func (p *VertexAIProvider) Generate(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, error) {
	if err := ValidateRequest(request); err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Convert messages to Vertex format
	contents := make([]*genai.Content, 0, len(request.Messages))
	for _, msg := range request.Messages {
		role := "user"
		if msg.Role == types.RoleSystem {
			role = "system"
		} else if msg.Role == types.RoleAssistant {
			role = "model"
		}
		content := &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		}
		contents = append(contents, content)
	}

	// Build request
	model := p.genModel
	if request.Model != "" {
		model = p.client.GenerativeModel(request.Model)
	}

	// Set generation config
	if request.Temperature > 0 {
		model.Temperature = genai.Ptr(request.Temperature)
	}
	if request.MaxTokens > 0 {
		model.MaxOutputTokens = genai.Ptr(int32(request.MaxTokens))
	}
	if request.TopP > 0 {
		model.TopP = genai.Ptr(request.TopP)
	}

	// Execute with retry
	var resp *genai.GenerateContentResponse
	err := p.WithRetry(ctx, func() error {
		var err error
		resp, err = model.GenerateContent(ctx, contents[0].Parts...)
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
		content = string(candidate.Content.Parts[0].(genai.Text))
	}

	finishReason := string(candidate.FinishReason)
	if finishReason == "FinishReasonUnspecified" {
		finishReason = ""
	}

	return &types.LLMResponse{
		Content:      content,
		FinishReason: finishReason,
		Usage: types.TokenUsage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		},
		Model:   p.model,
		Latency: time.Since(startTime),
	}, nil
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
