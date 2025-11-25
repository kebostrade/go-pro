package financial

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// TransactionLookupTool looks up transaction details
type TransactionLookupTool struct{}

// NewTransactionLookupTool creates a new transaction lookup tool
func NewTransactionLookupTool() *TransactionLookupTool {
	return &TransactionLookupTool{}
}

// Name returns the tool name
func (t *TransactionLookupTool) Name() string {
	return "transaction_lookup"
}

// Description returns the tool description
func (t *TransactionLookupTool) Description() string {
	return "Look up details of a financial transaction by transaction ID. Returns transaction amount, merchant, timestamp, status, and other relevant information."
}

// Execute runs the tool
func (t *TransactionLookupTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
	transactionID, ok := input.GetString("transaction_id")
	if !ok {
		return nil, fmt.Errorf("transaction_id is required")
	}

	// Simulate database lookup
	// In production, this would query your actual transaction database
	transaction := map[string]interface{}{
		"transaction_id": transactionID,
		"amount":         1250.50,
		"currency":       "USD",
		"merchant":       "Amazon.com",
		"merchant_id":    "AMZN_12345",
		"category":       "E-commerce",
		"timestamp":      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"status":         "completed",
		"payment_method": "credit_card",
		"card_last4":     "4242",
		"description":    "Online purchase - Electronics",
		"location": map[string]interface{}{
			"country": "US",
			"city":    "Seattle",
			"ip":      "192.168.1.1",
		},
		"risk_score": 0.15,
		"flags":      []string{},
	}

	return types.NewToolOutput(transaction), nil
}

// GetSchema returns the tool's input schema
func (t *TransactionLookupTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"transaction_id": {
				Type:        "string",
				Description: "The unique identifier of the transaction to look up",
			},
		},
		Required:    []string{"transaction_id"},
		Description: "Parameters for looking up a transaction",
	}
}

// Validate validates the input
func (t *TransactionLookupTool) Validate(input types.ToolInput) error {
	if _, ok := input.GetString("transaction_id"); !ok {
		return fmt.Errorf("transaction_id is required and must be a string")
	}
	return nil
}

