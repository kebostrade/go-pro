// Package tools provides public access to AI agent tools.
package tools

import (
	"context"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/golang"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/tools/programming"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodeAnalysisTool wraps the internal programming.CodeAnalysisTool for public use.
type CodeAnalysisTool = programming.CodeAnalysisTool

// NewCodeAnalysisTool creates a new code analysis tool.
func NewCodeAnalysisTool(registry *common.LanguageRegistry) *CodeAnalysisTool {
	return programming.NewCodeAnalysisTool(registry)
}

// NewDefaultCodeAnalysisTool creates a new code analysis tool with Go language support.
func NewDefaultCodeAnalysisTool() *CodeAnalysisTool {
	registry := common.NewLanguageRegistry()
	registry.Register(golang.NewProvider())
	return programming.NewCodeAnalysisTool(registry)
}

// Execute runs code analysis with the given input.
func ExecuteCodeAnalysis(ctx context.Context, tool *CodeAnalysisTool, code, language string) (*types.ToolOutput, error) {
	input := types.NewToolInput(map[string]interface{}{
		"code":     code,
		"language": language,
	})
	return tool.Execute(ctx, input)
}
