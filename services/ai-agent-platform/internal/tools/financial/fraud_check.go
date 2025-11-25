package financial

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// FraudCheckTool checks a transaction for fraud indicators
type FraudCheckTool struct{}

// NewFraudCheckTool creates a new fraud check tool
func NewFraudCheckTool() *FraudCheckTool {
	return &FraudCheckTool{}
}

// Name returns the tool name
func (t *FraudCheckTool) Name() string {
	return "fraud_check"
}

// Description returns the tool description
func (t *FraudCheckTool) Description() string {
	return "Analyze a transaction for fraud indicators. Returns risk score (0-1), fraud indicators, and recommended action. Use this when you need to assess if a transaction is potentially fraudulent."
}

// Execute runs the tool
func (t *FraudCheckTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
	transactionID, ok := input.GetString("transaction_id")
	if !ok {
		return nil, fmt.Errorf("transaction_id is required")
	}

	// Simulate fraud detection analysis
	// In production, this would use ML models and rule engines
	riskScore := rand.Float64() * 0.5 // Simulate risk score 0-0.5

	indicators := []string{}
	if riskScore > 0.3 {
		indicators = append(indicators, "unusual_transaction_amount")
	}
	if riskScore > 0.2 {
		indicators = append(indicators, "new_merchant")
	}

	var action string
	var confidence float64

	if riskScore < 0.2 {
		action = "approve"
		confidence = 0.95
	} else if riskScore < 0.4 {
		action = "review"
		confidence = 0.75
	} else {
		action = "decline"
		confidence = 0.85
		indicators = append(indicators, "high_risk_score")
	}

	result := map[string]interface{}{
		"transaction_id": transactionID,
		"risk_score":     riskScore,
		"risk_level":     getRiskLevel(riskScore),
		"indicators":     indicators,
		"action":         action,
		"confidence":     confidence,
		"analysis": map[string]interface{}{
			"velocity_check":     "passed",
			"location_check":     "passed",
			"device_check":       "passed",
			"pattern_match":      len(indicators) > 0,
			"blacklist_check":    "passed",
			"amount_threshold":   riskScore > 0.3,
		},
		"recommendation": getRecommendation(action, riskScore),
	}

	return types.NewToolOutput(result), nil
}

// GetSchema returns the tool's input schema
func (t *FraudCheckTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"transaction_id": {
				Type:        "string",
				Description: "The unique identifier of the transaction to check for fraud",
			},
		},
		Required:    []string{"transaction_id"},
		Description: "Parameters for fraud detection analysis",
	}
}

// Validate validates the input
func (t *FraudCheckTool) Validate(input types.ToolInput) error {
	if _, ok := input.GetString("transaction_id"); !ok {
		return fmt.Errorf("transaction_id is required and must be a string")
	}
	return nil
}

// getRiskLevel converts risk score to level
func getRiskLevel(score float64) string {
	if score < 0.2 {
		return "low"
	} else if score < 0.4 {
		return "medium"
	} else if score < 0.7 {
		return "high"
	}
	return "critical"
}

// getRecommendation provides action recommendation
func getRecommendation(action string, score float64) string {
	switch action {
	case "approve":
		return "Transaction appears legitimate. Safe to approve."
	case "review":
		return fmt.Sprintf("Transaction has moderate risk (%.2f). Manual review recommended.", score)
	case "decline":
		return fmt.Sprintf("Transaction has high risk (%.2f). Recommend declining and contacting customer.", score)
	default:
		return "Unable to determine recommendation."
	}
}

