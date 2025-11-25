package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the LLMProvider interface for OpenAI
type OpenAIProvider struct {
	*BaseProvider
	client *openai.Client
	model  string
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config OpenAIConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, errors.NewInvalidConfigError("OpenAI API key is required")
	}

	if config.Model == "" {
		config.Model = openai.GPT4TurboPreview
	}

	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	baseConfig := ProviderConfig{
		Provider:   "openai",
		APIKey:     config.APIKey,
		Model:      config.Model,
		Timeout:    config.Timeout,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
	}

	return &OpenAIProvider{
		BaseProvider: NewBaseProvider(baseConfig),
		client:       client,
		model:        config.Model,
	}, nil
}

// Generate generates a completion for the given request
func (p *OpenAIProvider) Generate(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, error) {
	if err := ValidateRequest(request); err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Convert messages to OpenAI format
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		}
	}

	// Build request
	model := request.Model
	if model == "" {
		model = p.model
	}

	chatRequest := openai.ChatCompletionRequest{
		Model:            model,
		Messages:         messages,
		Temperature:      request.Temperature,
		MaxTokens:        request.MaxTokens,
		TopP:             request.TopP,
		FrequencyPenalty: request.FrequencyPenalty,
		PresencePenalty:  request.PresencePenalty,
		Stop:             request.Stop,
		User:             request.User,
	}

	// Add function calling if provided
	if len(request.Functions) > 0 {
		functions := make([]openai.FunctionDefinition, len(request.Functions))
		for i, fn := range request.Functions {
			functions[i] = openai.FunctionDefinition{
				Name:        fn.Name,
				Description: fn.Description,
				Parameters:  fn.Parameters,
			}
		}
		chatRequest.Tools = make([]openai.Tool, len(functions))
		for i, fn := range functions {
			chatRequest.Tools[i] = openai.Tool{
				Type:     openai.ToolTypeFunction,
				Function: &fn,
			}
		}
	}

	// Execute with retry
	var resp openai.ChatCompletionResponse
	var err error

	err = p.WithRetry(ctx, func() error {
		resp, err = p.client.CreateChatCompletion(ctx, chatRequest)
		return err
	})

	if err != nil {
		return nil, errors.NewLLMProviderError("openai", err)
	}

	// Parse response
	if len(resp.Choices) == 0 {
		return nil, errors.NewLLMProviderError("openai", fmt.Errorf("no choices in response"))
	}

	choice := resp.Choices[0]
	response := &types.LLMResponse{
		Content:      choice.Message.Content,
		FinishReason: string(choice.FinishReason),
		Usage: types.TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Model:   resp.Model,
		Latency: time.Since(startTime),
	}

	// Handle function calls
	if len(choice.Message.ToolCalls) > 0 {
		toolCall := choice.Message.ToolCalls[0]
		response.FunctionCall = &types.FunctionCall{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		}
	}

	return response, nil
}

// GenerateStream generates a completion and streams the response
func (p *OpenAIProvider) GenerateStream(ctx context.Context, request types.LLMRequest) (<-chan types.LLMStreamChunk, error) {
	if err := ValidateRequest(request); err != nil {
		return nil, err
	}

	// Convert messages to OpenAI format
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		}
	}

	// Build request
	model := request.Model
	if model == "" {
		model = p.model
	}

	chatRequest := openai.ChatCompletionRequest{
		Model:            model,
		Messages:         messages,
		Temperature:      request.Temperature,
		MaxTokens:        request.MaxTokens,
		TopP:             request.TopP,
		FrequencyPenalty: request.FrequencyPenalty,
		PresencePenalty:  request.PresencePenalty,
		Stop:             request.Stop,
		User:             request.User,
		Stream:           true,
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, chatRequest)
	if err != nil {
		return nil, errors.NewLLMProviderError("openai", err)
	}

	// Create output channel
	chunks := make(chan types.LLMStreamChunk, 10)

	// Stream responses in goroutine
	go func() {
		defer close(chunks)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				chunks <- types.LLMStreamChunk{
					Done:  true,
					Error: err,
				}
				return
			}

			if len(response.Choices) == 0 {
				continue
			}

			choice := response.Choices[0]
			chunk := types.LLMStreamChunk{
				Content:      choice.Delta.Content,
				FinishReason: string(choice.FinishReason),
				Done:         choice.FinishReason != "",
			}

			chunks <- chunk

			if chunk.Done {
				return
			}
		}
	}()

	return chunks, nil
}

// GetModelInfo returns information about the model
func (p *OpenAIProvider) GetModelInfo() types.ModelInfo {
	// Pricing as of 2024 (approximate)
	costInfo := types.CostInfo{
		Input:  0.01,  // $0.01 per 1K tokens
		Output: 0.03,  // $0.03 per 1K tokens
	}

	if p.model == openai.GPT4TurboPreview || p.model == openai.GPT4 {
		costInfo.Input = 0.03
		costInfo.Output = 0.06
	}

	return types.ModelInfo{
		Name:              p.model,
		Provider:          "openai",
		MaxTokens:         8192,
		CostPer1KTokens:   costInfo,
		SupportsStreaming: true,
		SupportsFunctions: true,
		SupportsVision:    p.model == openai.GPT4VisionPreview,
	}
}

// GetProviderName returns the name of the provider
func (p *OpenAIProvider) GetProviderName() string {
	return "openai"
}

// SupportsStreaming returns whether the provider supports streaming
func (p *OpenAIProvider) SupportsStreaming() bool {
	return true
}

// SupportsFunctionCalling returns whether the provider supports function calling
func (p *OpenAIProvider) SupportsFunctionCalling() bool {
	return true
}

