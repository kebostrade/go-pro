package golang

import (
	"context"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// Provider combines Go analyzer and executor
type Provider struct {
	*common.BaseLanguageProvider
	analyzer *Analyzer
	executor *Executor
}

// NewProvider creates a new Go language provider
func NewProvider() *Provider {
	lang, _ := types.GetLanguage("go")
	return &Provider{
		BaseLanguageProvider: common.NewBaseLanguageProvider(lang),
		analyzer:             NewAnalyzer(),
		executor:             NewExecutor(),
	}
}

// Analyze performs static analysis on Go code
func (p *Provider) Analyze(ctx context.Context, code string) (*types.CodeAnalysis, error) {
	return p.analyzer.Analyze(ctx, code)
}

// Lint runs linting on Go code
func (p *Provider) Lint(ctx context.Context, code string) ([]types.CodeIssue, error) {
	return p.analyzer.Lint(ctx, code)
}

// Format formats Go code
func (p *Provider) Format(ctx context.Context, code string) (string, error) {
	return p.analyzer.Format(ctx, code)
}

// Validate checks if Go code is valid
func (p *Provider) Validate(ctx context.Context, code string) error {
	return p.analyzer.Validate(ctx, code)
}

// ExtractImports extracts imports from Go code
func (p *Provider) ExtractImports(ctx context.Context, code string) ([]string, error) {
	return p.analyzer.ExtractImports(ctx, code)
}

// GetComplexity calculates cyclomatic complexity
func (p *Provider) GetComplexity(ctx context.Context, code string) (int, error) {
	return p.analyzer.GetComplexity(ctx, code)
}

// Execute runs Go code
func (p *Provider) Execute(ctx context.Context, request types.ExecutionRequest) (*types.ExecutionResult, error) {
	return p.executor.Execute(ctx, request)
}

// ValidateCode checks if code is safe to execute
func (p *Provider) ValidateCode(ctx context.Context, code string) error {
	return p.executor.ValidateCode(ctx, code)
}

