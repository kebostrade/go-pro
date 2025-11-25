package programming

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// DocumentationSearchTool searches official documentation
type DocumentationSearchTool struct {
	docSources map[string][]DocSource
}

// DocSource represents a documentation source
type DocSource struct {
	Name        string
	BaseURL     string
	Description string
}

// NewDocumentationSearchTool creates a new documentation search tool
func NewDocumentationSearchTool() *DocumentationSearchTool {
	return &DocumentationSearchTool{
		docSources: map[string][]DocSource{
			"go": {
				{
					Name:        "Go Official Documentation",
					BaseURL:     "https://go.dev/doc/",
					Description: "Official Go programming language documentation",
				},
				{
					Name:        "Go Package Documentation",
					BaseURL:     "https://pkg.go.dev/",
					Description: "Go package documentation and reference",
				},
				{
					Name:        "Go by Example",
					BaseURL:     "https://gobyexample.com/",
					Description: "Hands-on introduction to Go using annotated example programs",
				},
			},
			"python": {
				{
					Name:        "Python Official Documentation",
					BaseURL:     "https://docs.python.org/3/",
					Description: "Official Python documentation",
				},
				{
					Name:        "Python Package Index",
					BaseURL:     "https://pypi.org/",
					Description: "Python package repository",
				},
			},
			"javascript": {
				{
					Name:        "MDN Web Docs",
					BaseURL:     "https://developer.mozilla.org/en-US/docs/Web/JavaScript/",
					Description: "Comprehensive JavaScript documentation",
				},
				{
					Name:        "Node.js Documentation",
					BaseURL:     "https://nodejs.org/docs/",
					Description: "Official Node.js documentation",
				},
			},
			"typescript": {
				{
					Name:        "TypeScript Documentation",
					BaseURL:     "https://www.typescriptlang.org/docs/",
					Description: "Official TypeScript documentation",
				},
			},
			"rust": {
				{
					Name:        "The Rust Programming Language",
					BaseURL:     "https://doc.rust-lang.org/book/",
					Description: "The official Rust book",
				},
				{
					Name:        "Rust Standard Library",
					BaseURL:     "https://doc.rust-lang.org/std/",
					Description: "Rust standard library documentation",
				},
			},
		},
	}
}

// Name returns the tool name
func (t *DocumentationSearchTool) Name() string {
	return "documentation_search"
}

// Description returns the tool description
func (t *DocumentationSearchTool) Description() string {
	return "Searches official documentation for programming languages and frameworks. Returns relevant documentation links and snippets."
}

// GetSchema returns the tool's input schema
func (t *DocumentationSearchTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"query": {
				Type:        "string",
				Description: "Search query for documentation",
			},
			"language": {
				Type:        "string",
				Description: "Programming language (go, python, javascript, etc.)",
			},
			"topic": {
				Type:        "string",
				Description: "Specific topic or category (optional)",
			},
		},
		Required: []string{"query", "language"},
	}
}

// Validate validates the input
func (t *DocumentationSearchTool) Validate(input types.ToolInput) error {
	query, ok := input.GetString("query")
	if !ok || query == "" {
		return &types.ToolError{
			Code:    "MISSING_QUERY",
			Message: "Query parameter is required",
		}
	}

	language, ok := input.GetString("language")
	if !ok || language == "" {
		return &types.ToolError{
			Code:    "MISSING_LANGUAGE",
			Message: "Language parameter is required",
		}
	}

	return nil
}

// Execute searches documentation
func (t *DocumentationSearchTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
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
	topic, _ := input.GetString("topic")

	// Get documentation sources for the language
	sources, ok := t.docSources[language]
	if !ok {
		sources = []DocSource{
			{
				Name:        fmt.Sprintf("%s Documentation", strings.Title(language)),
				BaseURL:     fmt.Sprintf("https://www.google.com/search?q=%s+%s+documentation", language, query),
				Description: fmt.Sprintf("Search for %s documentation", language),
			},
		}
	}

	// Build search results
	results := make([]types.Reference, 0)
	for _, source := range sources {
		// Create search URL
		searchURL := t.buildSearchURL(source.BaseURL, query, topic)

		ref := types.Reference{
			Title:     source.Name,
			URL:       searchURL,
			Type:      types.ReferenceTypeOfficialDocs,
			Relevance: t.calculateRelevance(source.Name, query, topic),
			Snippet:   source.Description,
		}
		results = append(results, ref)
	}

	// Add common documentation patterns
	results = append(results, t.getCommonDocPatterns(language, query)...)

	// Convert to JSON
	resultJSON, _ := json.MarshalIndent(map[string]interface{}{
		"query":     query,
		"language":  language,
		"topic":     topic,
		"results":   results,
		"count":     len(results),
	}, "", "  ")

	return &types.ToolOutput{
		Result: string(resultJSON),
		Metadata: types.ToolMetadata{
			ExecutionTime: time.Since(startTime).Milliseconds(),
			Success:       true,
			AdditionalInfo: map[string]interface{}{
				"language":      language,
				"results_count": len(results),
			},
		},
	}, nil
}

// buildSearchURL constructs a search URL
func (t *DocumentationSearchTool) buildSearchURL(baseURL, query, topic string) string {
	// Simple URL construction - in production, use proper URL encoding
	searchTerm := query
	if topic != "" {
		searchTerm = fmt.Sprintf("%s %s", topic, query)
	}
	
	// For Google search fallback
	if strings.Contains(baseURL, "google.com") {
		return baseURL
	}
	
	return fmt.Sprintf("%s?q=%s", baseURL, strings.ReplaceAll(searchTerm, " ", "+"))
}

// calculateRelevance calculates relevance score
func (t *DocumentationSearchTool) calculateRelevance(sourceName, query, topic string) float64 {
	relevance := 0.5 // Base relevance
	
	// Increase relevance for official docs
	if strings.Contains(strings.ToLower(sourceName), "official") {
		relevance += 0.3
	}
	
	// Increase relevance if topic matches
	if topic != "" && strings.Contains(strings.ToLower(sourceName), strings.ToLower(topic)) {
		relevance += 0.2
	}
	
	if relevance > 1.0 {
		relevance = 1.0
	}
	
	return relevance
}

// getCommonDocPatterns returns common documentation patterns
func (t *DocumentationSearchTool) getCommonDocPatterns(language, query string) []types.Reference {
	patterns := make([]types.Reference, 0)
	
	// Add language-specific patterns
	switch language {
	case "go":
		if strings.Contains(strings.ToLower(query), "package") || strings.Contains(strings.ToLower(query), "import") {
			patterns = append(patterns, types.Reference{
				Title:     "Go Packages",
				URL:       "https://pkg.go.dev/",
				Type:      types.ReferenceTypeOfficialDocs,
				Relevance: 0.9,
				Snippet:   "Search Go packages and their documentation",
			})
		}
	case "python":
		if strings.Contains(strings.ToLower(query), "pip") || strings.Contains(strings.ToLower(query), "install") {
			patterns = append(patterns, types.Reference{
				Title:     "Python Package Index",
				URL:       "https://pypi.org/",
				Type:      types.ReferenceTypeOfficialDocs,
				Relevance: 0.9,
				Snippet:   "Find and install Python packages",
			})
		}
	case "javascript":
		if strings.Contains(strings.ToLower(query), "npm") || strings.Contains(strings.ToLower(query), "package") {
			patterns = append(patterns, types.Reference{
				Title:     "npm Registry",
				URL:       "https://www.npmjs.com/",
				Type:      types.ReferenceTypeOfficialDocs,
				Relevance: 0.9,
				Snippet:   "Find and install JavaScript packages",
			})
		}
	}
	
	return patterns
}

