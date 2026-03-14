package llm

import (
	"context"
	"os"
	"time"

	"github.com/DimaJoyti/go-pro/services/langchain/pkg/schema"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMProvider struct {
	client      *openai.LLM
	model       string
	provider    string
	temperature float64
	maxTokens   int
}

type Options struct {
	Model       string
	APIKey      string
	BaseURL     string
	Temperature float64
	MaxTokens   int
}

func NewLLMProvider(provider string, opts Options) (*LLMProvider, error) {
	var apiKey string
	var baseURL string

	switch provider {
	case "openai":
		apiKey = opts.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if opts.BaseURL != "" {
			baseURL = opts.BaseURL
		}
	case "anthropic":
		apiKey = opts.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		baseURL = "https://api.anthropic.com/v1"
	default:
		apiKey = opts.APIKey
		baseURL = opts.BaseURL
	}

	temperature := opts.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	maxTokens := opts.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048
	}

	model := opts.Model
	if model == "" {
		model = "gpt-4"
	}

	var client *openai.LLM
	var err error

	options := []openai.Option{
		openai.WithToken(apiKey),
		openai.WithModel(model),
	}

	if baseURL != "" {
		options = append(options, openai.WithBaseURL(baseURL))
	}

	client, err = openai.New(options...)

	if err != nil {
		return nil, err
	}

	p := &LLMProvider{
		client:      client,
		model:       model,
		provider:    provider,
		temperature: temperature,
		maxTokens:   maxTokens,
	}

	return p, nil
}

func (p *LLMProvider) Generate(ctx context.Context, messages []schema.Message) (schema.Message, error) {
	llmMessages := convertToLLMMessages(messages)

	resp, err := p.client.GenerateContent(ctx, llmMessages,
		llms.WithTemperature(p.temperature),
		llms.WithMaxTokens(p.maxTokens),
	)
	if err != nil {
		return schema.Message{}, err
	}

	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Content
	}

	return schema.NewAssistantMessage(content), nil
}

func (p *LLMProvider) GenerateWithTools(ctx context.Context, messages []schema.Message, tools []schema.Tool) (schema.Message, error) {
	llmMessages := convertToLLMMessages(messages)
	functions := convertToFunctions(tools)

	resp, err := p.client.GenerateContent(ctx, llmMessages,
		llms.WithFunctions(functions),
		llms.WithTemperature(p.temperature),
		llms.WithMaxTokens(p.maxTokens),
	)
	if err != nil {
		return schema.Message{}, err
	}

	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Content
	}

	return schema.Message{
		Role:      schema.RoleAssistant,
		Content:   content,
		Timestamp: time.Now(),
	}, nil
}

func convertToLLMMessages(messages []schema.Message) []llms.MessageContent {
	llmMessages := make([]llms.MessageContent, 0, len(messages))

	for _, msg := range messages {
		var role llms.ChatMessageType
		switch msg.Role {
		case schema.RoleSystem:
			role = llms.ChatMessageTypeSystem
		case schema.RoleUser:
			role = llms.ChatMessageTypeHuman
		case schema.RoleAssistant:
			role = llms.ChatMessageTypeAI
		case schema.RoleTool:
			role = llms.ChatMessageTypeTool
		default:
			role = llms.ChatMessageTypeHuman
		}

		llmMessages = append(llmMessages, llms.TextParts(role, msg.Content))
	}

	return llmMessages
}

func convertToFunctions(tools []schema.Tool) []llms.FunctionDefinition {
	functions := make([]llms.FunctionDefinition, 0, len(tools))

	for _, tool := range tools {
		functions = append(functions, llms.FunctionDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.ArgsSchema,
		})
	}

	return functions
}
