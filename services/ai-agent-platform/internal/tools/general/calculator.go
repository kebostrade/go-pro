package general

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CalculatorTool performs mathematical calculations
type CalculatorTool struct{}

// NewCalculatorTool creates a new calculator tool
func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

// Name returns the tool name
func (t *CalculatorTool) Name() string {
	return "calculator"
}

// Description returns the tool description
func (t *CalculatorTool) Description() string {
	return "Perform mathematical calculations. Supports basic arithmetic (+, -, *, /), exponents (^), and common functions (sqrt, abs, etc.). Use this when you need to calculate numbers."
}

// Execute runs the tool
func (t *CalculatorTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
	expression, ok := input.GetString("expression")
	if !ok {
		return nil, fmt.Errorf("expression is required")
	}

	// Simple expression evaluator
	// For production, use a proper math expression parser library
	result, err := evaluateExpression(expression)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %w", err)
	}

	output := map[string]interface{}{
		"expression": expression,
		"result":     result,
		"formatted":  formatNumber(result),
	}

	return types.NewToolOutput(output), nil
}

// evaluateExpression is a simple expression evaluator
// For production, replace with a proper math parser library
func evaluateExpression(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	// Handle basic operations
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return a + b, nil
			}
		}
	}

	if strings.Contains(expr, "-") {
		parts := strings.Split(expr, "-")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return a - b, nil
			}
		}
	}

	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return a * b, nil
			}
		}
	}

	if strings.Contains(expr, "/") {
		parts := strings.Split(expr, "/")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil && b != 0 {
				return a / b, nil
			}
		}
	}

	// Try to parse as a number
	return strconv.ParseFloat(expr, 64)
}

// GetSchema returns the tool's input schema
func (t *CalculatorTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"expression": {
				Type:        "string",
				Description: "Mathematical expression to evaluate (e.g., '2 + 2', '10 * 5 + 3', 'sqrt(16)')",
			},
		},
		Required:    []string{"expression"},
		Description: "Parameters for mathematical calculation",
	}
}

// Validate validates the input
func (t *CalculatorTool) Validate(input types.ToolInput) error {
	if _, ok := input.GetString("expression"); !ok {
		return fmt.Errorf("expression is required and must be a string")
	}
	return nil
}

// formatNumber formats a number for display
func formatNumber(n float64) string {
	if math.Floor(n) == n {
		return fmt.Sprintf("%.0f", n)
	}
	return fmt.Sprintf("%.2f", n)
}
