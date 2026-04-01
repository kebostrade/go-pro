package llm

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// ProviderConfig holds configuration for LLM providers
type ProviderConfig struct {
	// Provider name (openai, anthropic, ollama)
	Provider string `json:"provider"`

	// APIKey for authentication
	APIKey string `json:"api_key"`

	// BaseURL for the API (optional, for custom endpoints)
	BaseURL string `json:"base_url,omitempty"`

	// Model to use
	Model string `json:"model"`

	// Timeout for requests
	Timeout time.Duration `json:"timeout"`

	// MaxRetries for failed requests
	MaxRetries int `json:"max_retries"`

	// RetryDelay between retries
	RetryDelay time.Duration `json:"retry_delay"`

	// EnableCaching whether to cache responses
	EnableCaching bool `json:"enable_caching"`

	// CacheTTL time-to-live for cached responses
	CacheTTL time.Duration `json:"cache_ttl"`
}

// ProviderManager manages multiple LLM providers
type ProviderManager struct {
	providers map[string]types.LLMProvider
	mu        sync.RWMutex
	cache     *ResponseCache
}

// NewProviderManager creates a new provider manager
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]types.LLMProvider),
		cache:     NewResponseCache(),
	}
}

// Register registers a new provider
func (pm *ProviderManager) Register(name string, provider types.LLMProvider) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	pm.providers[name] = provider
	return nil
}

// Get retrieves a provider by name
func (pm *ProviderManager) Get(name string) (types.LLMProvider, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	provider, exists := pm.providers[name]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("provider: %s", name))
	}

	return provider, nil
}

// List returns all registered providers
func (pm *ProviderManager) List() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	names := make([]string, 0, len(pm.providers))
	for name := range pm.providers {
		names = append(names, name)
	}
	return names
}

// Generate delegates to the first registered provider
func (pm *ProviderManager) Generate(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.providers) == 0 {
		return nil, errors.NewLLMProviderError("provider manager", fmt.Errorf("no providers registered"))
	}

	for _, provider := range pm.providers {
		response, err := provider.Generate(ctx, request)
		if err == nil {
			return response, nil
		}
	}

	return nil, errors.NewLLMProviderError("all providers", fmt.Errorf("all providers failed"))
}

// GenerateStream delegates streaming to the first registered provider
func (pm *ProviderManager) GenerateStream(ctx context.Context, request types.LLMRequest) (<-chan types.LLMStreamChunk, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.providers) == 0 {
		return nil, errors.NewLLMProviderError("provider manager", fmt.Errorf("no providers registered"))
	}

	for _, provider := range pm.providers {
		if provider.SupportsStreaming() {
			return provider.GenerateStream(ctx, request)
		}
	}

	return nil, errors.NewLLMProviderError("all providers", fmt.Errorf("no provider supports streaming"))
}

// GetModelInfo returns info about the first registered provider
func (pm *ProviderManager) GetModelInfo() types.ModelInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, provider := range pm.providers {
		return provider.GetModelInfo()
	}

	return types.ModelInfo{}
}

// GetProviderName returns the provider manager name
func (pm *ProviderManager) GetProviderName() string {
	return "provider-manager"
}

// SupportsStreaming returns whether the first registered provider supports streaming
func (pm *ProviderManager) SupportsStreaming() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, provider := range pm.providers {
		return provider.SupportsStreaming()
	}
	return false
}

// SupportsFunctionCalling returns whether the first registered provider supports function calling
func (pm *ProviderManager) SupportsFunctionCalling() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, provider := range pm.providers {
		return provider.SupportsFunctionCalling()
	}
	return false
}

// GenerateWithFallback tries multiple providers in order
func (pm *ProviderManager) GenerateWithFallback(
	ctx context.Context,
	providerNames []string,
	request types.LLMRequest,
) (*types.LLMResponse, error) {
	var lastErr error

	for _, name := range providerNames {
		provider, err := pm.Get(name)
		if err != nil {
			lastErr = err
			continue
		}

		response, err := provider.Generate(ctx, request)
		if err == nil {
			return response, nil
		}

		lastErr = err
	}

	if lastErr != nil {
		return nil, errors.NewLLMProviderError("all providers", lastErr)
	}

	return nil, errors.NewLLMProviderError("all providers", fmt.Errorf("no providers available"))
}

// BaseProvider provides common functionality for LLM providers
type BaseProvider struct {
	config ProviderConfig
	cache  *ResponseCache
}

// NewBaseProvider creates a new base provider
func NewBaseProvider(config ProviderConfig) *BaseProvider {
	return &BaseProvider{
		config: config,
		cache:  NewResponseCache(),
	}
}

// GetConfig returns the provider configuration
func (bp *BaseProvider) GetConfig() ProviderConfig {
	return bp.config
}

// WithRetry executes a function with retry logic
func (bp *BaseProvider) WithRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	maxRetries := bp.config.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}

	for i := 0; i < maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(bp.config.RetryDelay):
				// Continue to next retry
			}
		}
	}

	return lastErr
}

// CalculateCost calculates the cost of a request
func CalculateCost(usage types.TokenUsage, costInfo types.CostInfo) float64 {
	inputCost := float64(usage.PromptTokens) / 1000.0 * costInfo.Input
	outputCost := float64(usage.CompletionTokens) / 1000.0 * costInfo.Output
	return inputCost + outputCost
}

// ValidateRequest validates an LLM request
func ValidateRequest(request types.LLMRequest) error {
	if len(request.Messages) == 0 {
		return errors.NewAgentInvalidInputError("messages cannot be empty")
	}

	if request.Temperature < 0 || request.Temperature > 2 {
		return errors.NewAgentInvalidInputError("temperature must be between 0 and 2")
	}

	if request.MaxTokens < 0 {
		return errors.NewAgentInvalidInputError("max_tokens must be positive")
	}

	return nil
}

// FormatMessages formats messages for display
func FormatMessages(messages []types.Message) string {
	var result string
	for i, msg := range messages {
		result += fmt.Sprintf("[%d] %s: %s\n", i, msg.Role, msg.Content)
	}
	return result
}

// TruncateMessages truncates messages to fit within token limit
func TruncateMessages(messages []types.Message, maxTokens int) []types.Message {
	// Simple implementation - keep system message and recent messages
	if len(messages) == 0 {
		return messages
	}

	var result []types.Message

	// Always keep system message if present
	if messages[0].Role == types.RoleSystem {
		result = append(result, messages[0])
		messages = messages[1:]
	}

	// Estimate tokens (rough approximation: 1 token ≈ 4 characters)
	estimatedTokens := 0
	for i := len(messages) - 1; i >= 0; i-- {
		msgTokens := len(messages[i].Content) / 4
		if estimatedTokens+msgTokens > maxTokens {
			break
		}
		result = append([]types.Message{messages[i]}, result...)
		estimatedTokens += msgTokens
	}

	return result
}

func getEnvOr(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

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
