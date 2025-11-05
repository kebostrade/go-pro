package programming

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// StackOverflowSearchTool searches Stack Overflow for programming questions
type StackOverflowSearchTool struct {
	baseURL string
}

// NewStackOverflowSearchTool creates a new Stack Overflow search tool
func NewStackOverflowSearchTool() *StackOverflowSearchTool {
	return &StackOverflowSearchTool{
		baseURL: "https://stackoverflow.com",
	}
}

// Name returns the tool name
func (t *StackOverflowSearchTool) Name() string {
	return "stackoverflow_search"
}

// Description returns the tool description
func (t *StackOverflowSearchTool) Description() string {
	return "Searches Stack Overflow for programming questions and answers. Returns relevant questions with high-quality answers."
}

// GetSchema returns the tool's input schema
func (t *StackOverflowSearchTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"query": {
				Type:        "string",
				Description: "Search query for Stack Overflow",
			},
			"language": {
				Type:        "string",
				Description: "Programming language tag (optional)",
			},
			"tags": {
				Type:        "array",
				Description: "Additional tags to filter by",
				Items: &types.PropertySchema{
					Type: "string",
				},
			},
			"min_score": {
				Type:        "number",
				Description: "Minimum question score (default: 0)",
				Default:     0,
			},
			"has_accepted_answer": {
				Type:        "boolean",
				Description: "Only return questions with accepted answers",
				Default:     false,
			},
		},
		Required: []string{"query"},
	}
}

// Validate validates the input
func (t *StackOverflowSearchTool) Validate(input types.ToolInput) error {
	query, ok := input.GetString("query")
	if !ok || query == "" {
		return &types.ToolError{
			Code:    "MISSING_QUERY",
			Message: "Query parameter is required",
		}
	}

	return nil
}

// Execute searches Stack Overflow
func (t *StackOverflowSearchTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
	startTime := time.Now()

	// Validate input
	if err := t.Validate(input); err != nil {
		return &types.ToolOutput{
			Error: err.(*types.ToolError),
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Get parameters
	query, _ := input.GetString("query")
	language, _ := input.GetString("language")
	minScore, _ := input.GetInt("min_score")
	hasAcceptedAnswer, _ := input.GetBool("has_accepted_answer")

	// Build tags
	tags := make([]string, 0)
	if language != "" {
		tags = append(tags, language)
	}
	if tagSlice, ok := input.GetSlice("tags"); ok {
		for _, tag := range tagSlice {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			}
		}
	}

	// Build search URL
	searchURL := t.buildSearchURL(query, tags, minScore, hasAcceptedAnswer)

	// Simulate search results (in production, use Stack Overflow API)
	results := t.generateMockResults(query, language, tags)

	// Convert to JSON
	resultJSON, _ := json.MarshalIndent(map[string]interface{}{
		"query":               query,
		"language":            language,
		"tags":                tags,
		"search_url":          searchURL,
		"results":             results,
		"count":               len(results),
		"has_accepted_answer": hasAcceptedAnswer,
		"min_score":           minScore,
	}, "", "  ")

	return &types.ToolOutput{
		Result: string(resultJSON),
		Metadata: types.ToolMetadata{
			ExecutionTime: time.Since(startTime).Milliseconds(),
			Success:       true,
			AdditionalInfo: map[string]interface{}{
				"results_count": len(results),
				"search_url":    searchURL,
			},
		},
	}, nil
}

// buildSearchURL constructs a Stack Overflow search URL
func (t *StackOverflowSearchTool) buildSearchURL(query string, tags []string, minScore int, hasAcceptedAnswer bool) string {
	// Build search query
	searchQuery := query

	// Add tags
	if len(tags) > 0 {
		for _, tag := range tags {
			searchQuery += fmt.Sprintf(" [%s]", tag)
		}
	}

	// Add filters
	if minScore > 0 {
		searchQuery += fmt.Sprintf(" score:%d", minScore)
	}

	if hasAcceptedAnswer {
		searchQuery += " hasaccepted:yes"
	}

	// URL encode
	encodedQuery := strings.ReplaceAll(searchQuery, " ", "+")

	return fmt.Sprintf("%s/search?q=%s", t.baseURL, encodedQuery)
}

// generateMockResults generates mock search results
// In production, this would call the Stack Overflow API
func (t *StackOverflowSearchTool) generateMockResults(query, language string, tags []string) []types.Reference {
	results := make([]types.Reference, 0)

	// Generate relevant mock results based on query
	if strings.Contains(strings.ToLower(query), "error") || strings.Contains(strings.ToLower(query), "fix") {
		results = append(results, types.Reference{
			Title:     fmt.Sprintf("How to fix %s in %s", query, language),
			URL:       fmt.Sprintf("%s/questions/12345/how-to-fix-%s", t.baseURL, strings.ReplaceAll(query, " ", "-")),
			Type:      types.ReferenceTypeStackOverflow,
			Relevance: 0.9,
			Snippet:   fmt.Sprintf("Common solutions for %s errors in %s programming", query, language),
		})
	}

	if strings.Contains(strings.ToLower(query), "how to") || strings.Contains(strings.ToLower(query), "best way") {
		results = append(results, types.Reference{
			Title:     fmt.Sprintf("Best practices for %s in %s", query, language),
			URL:       fmt.Sprintf("%s/questions/23456/best-practices-%s", t.baseURL, strings.ReplaceAll(query, " ", "-")),
			Type:      types.ReferenceTypeStackOverflow,
			Relevance: 0.85,
			Snippet:   fmt.Sprintf("Community-recommended approaches for %s", query),
		})
	}

	// Add general result
	results = append(results, types.Reference{
		Title:     fmt.Sprintf("%s - Stack Overflow", query),
		URL:       t.buildSearchURL(query, tags, 0, false),
		Type:      types.ReferenceTypeStackOverflow,
		Relevance: 0.8,
		Snippet:   fmt.Sprintf("Search results for '%s' on Stack Overflow", query),
	})

	// Add language-specific common questions
	if language != "" {
		results = append(results, t.getCommonQuestions(language, query)...)
	}

	return results
}

// getCommonQuestions returns common questions for a language
func (t *StackOverflowSearchTool) getCommonQuestions(language, query string) []types.Reference {
	common := make([]types.Reference, 0)

	switch language {
	case "go":
		if strings.Contains(strings.ToLower(query), "goroutine") || strings.Contains(strings.ToLower(query), "concurrent") {
			common = append(common, types.Reference{
				Title:     "Understanding Goroutines and Channels",
				URL:       fmt.Sprintf("%s/questions/18058164/understanding-goroutines", t.baseURL),
				Type:      types.ReferenceTypeStackOverflow,
				Relevance: 0.95,
				Snippet:   "Comprehensive guide to Go concurrency patterns",
			})
		}
		if strings.Contains(strings.ToLower(query), "interface") {
			common = append(common, types.Reference{
				Title:     "How to use interfaces in Go",
				URL:       fmt.Sprintf("%s/questions/23148812/go-interfaces", t.baseURL),
				Type:      types.ReferenceTypeStackOverflow,
				Relevance: 0.9,
				Snippet:   "Best practices for using interfaces in Go",
			})
		}

	case "python":
		if strings.Contains(strings.ToLower(query), "list") || strings.Contains(strings.ToLower(query), "array") {
			common = append(common, types.Reference{
				Title:     "Python List Comprehensions",
				URL:       fmt.Sprintf("%s/questions/34835951/list-comprehension", t.baseURL),
				Type:      types.ReferenceTypeStackOverflow,
				Relevance: 0.9,
				Snippet:   "Understanding Python list comprehensions",
			})
		}

	case "javascript":
		if strings.Contains(strings.ToLower(query), "async") || strings.Contains(strings.ToLower(query), "promise") {
			common = append(common, types.Reference{
				Title:     "Understanding JavaScript Promises and Async/Await",
				URL:       fmt.Sprintf("%s/questions/14220321/promises-async", t.baseURL),
				Type:      types.ReferenceTypeStackOverflow,
				Relevance: 0.95,
				Snippet:   "Complete guide to asynchronous JavaScript",
			})
		}
	}

	return common
}

