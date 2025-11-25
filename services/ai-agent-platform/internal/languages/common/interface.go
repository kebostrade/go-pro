package common

import (
	"context"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// LanguageProvider combines analyzer and executor interfaces
type LanguageProvider interface {
	types.LanguageAnalyzer
	types.LanguageExecutor
}

// BaseLanguageProvider provides common functionality for language providers
type BaseLanguageProvider struct {
	language types.Language
}

// NewBaseLanguageProvider creates a new base language provider
func NewBaseLanguageProvider(language types.Language) *BaseLanguageProvider {
	return &BaseLanguageProvider{
		language: language,
	}
}

// GetLanguage returns the language this provider handles
func (b *BaseLanguageProvider) GetLanguage() types.Language {
	return b.language
}

// GetResourceLimits returns default resource limits
func (b *BaseLanguageProvider) GetResourceLimits() types.ResourceLimits {
	return types.ResourceLimits{
		MaxMemoryMB:      512,
		MaxCPUTime:       30,
		MaxProcesses:     10,
		MaxFileSize:      10 * 1024 * 1024, // 10MB
		MaxOutputSize:    1 * 1024 * 1024,  // 1MB
		NetworkAccess:    false,
		FileSystemAccess: false,
	}
}

// SupportsInteractive returns whether interactive execution is supported
func (b *BaseLanguageProvider) SupportsInteractive() bool {
	return false
}

// ValidateCode performs basic validation
func (b *BaseLanguageProvider) ValidateCode(ctx context.Context, code string) error {
	if code == "" {
		return &types.ToolError{
			Code:    "EMPTY_CODE",
			Message: "Code cannot be empty",
		}
	}
	return nil
}

// LanguageRegistry manages all language providers
type LanguageRegistry struct {
	providers map[string]LanguageProvider
}

// NewLanguageRegistry creates a new language registry
func NewLanguageRegistry() *LanguageRegistry {
	return &LanguageRegistry{
		providers: make(map[string]LanguageProvider),
	}
}

// Register adds a language provider to the registry
func (r *LanguageRegistry) Register(provider LanguageProvider) error {
	lang := provider.GetLanguage()
	if _, exists := r.providers[lang.Name]; exists {
		return &types.ToolError{
			Code:    "DUPLICATE_LANGUAGE",
			Message: "Language provider already registered: " + lang.Name,
		}
	}
	r.providers[lang.Name] = provider
	return nil
}

// Get retrieves a language provider by name
func (r *LanguageRegistry) Get(name string) (LanguageProvider, error) {
	provider, ok := r.providers[name]
	if !ok {
		return nil, &types.ToolError{
			Code:    "LANGUAGE_NOT_FOUND",
			Message: "Language provider not found: " + name,
		}
	}
	return provider, nil
}

// List returns all registered language providers
func (r *LanguageRegistry) List() []LanguageProvider {
	providers := make([]LanguageProvider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// GetSupportedLanguages returns all supported languages
func (r *LanguageRegistry) GetSupportedLanguages() []types.Language {
	languages := make([]types.Language, 0, len(r.providers))
	for _, provider := range r.providers {
		languages = append(languages, provider.GetLanguage())
	}
	return languages
}

// IsSupported checks if a language is supported
func (r *LanguageRegistry) IsSupported(name string) bool {
	_, ok := r.providers[name]
	return ok
}

