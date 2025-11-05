package programming

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// GitHubSearchTool searches GitHub for code examples and repositories
type GitHubSearchTool struct {
	baseURL string
}

// NewGitHubSearchTool creates a new GitHub search tool
func NewGitHubSearchTool() *GitHubSearchTool {
	return &GitHubSearchTool{
		baseURL: "https://github.com",
	}
}

// Name returns the tool name
func (t *GitHubSearchTool) Name() string {
	return "github_search"
}

// Description returns the tool description
func (t *GitHubSearchTool) Description() string {
	return "Searches GitHub for code examples, repositories, and open-source projects. Returns relevant repositories and code snippets."
}

// GetSchema returns the tool's input schema
func (t *GitHubSearchTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"query": {
				Type:        "string",
				Description: "Search query for GitHub",
			},
			"language": {
				Type:        "string",
				Description: "Programming language to filter by",
			},
			"search_type": {
				Type:        "string",
				Description: "Type of search: repositories, code, issues, users",
				Enum:        []string{"repositories", "code", "issues", "users"},
				Default:     "repositories",
			},
			"sort": {
				Type:        "string",
				Description: "Sort by: stars, forks, updated, best-match",
				Enum:        []string{"stars", "forks", "updated", "best-match"},
				Default:     "best-match",
			},
			"min_stars": {
				Type:        "number",
				Description: "Minimum number of stars",
				Default:     0,
			},
		},
		Required: []string{"query"},
	}
}

// Validate validates the input
func (t *GitHubSearchTool) Validate(input types.ToolInput) error {
	query, ok := input.GetString("query")
	if !ok || query == "" {
		return &types.ToolError{
			Code:    "MISSING_QUERY",
			Message: "Query parameter is required",
		}
	}

	return nil
}

// Execute searches GitHub
func (t *GitHubSearchTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
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
	searchType, ok := input.GetString("search_type")
	if !ok {
		searchType = "repositories"
	}
	sortBy, ok := input.GetString("sort")
	if !ok {
		sortBy = "best-match"
	}
	minStars, _ := input.GetInt("min_stars")

	// Build search URL
	searchURL := t.buildSearchURL(query, language, searchType, sortBy, minStars)

	// Generate results
	results := t.generateResults(query, language, searchType, minStars)

	// Convert to JSON
	resultJSON, _ := json.MarshalIndent(map[string]interface{}{
		"query":       query,
		"language":    language,
		"search_type": searchType,
		"sort_by":     sortBy,
		"min_stars":   minStars,
		"search_url":  searchURL,
		"results":     results,
		"count":       len(results),
	}, "", "  ")

	return &types.ToolOutput{
		Result: string(resultJSON),
		Metadata: types.ToolMetadata{
			ExecutionTime: time.Since(startTime).Milliseconds(),
			Success:       true,
			AdditionalInfo: map[string]interface{}{
				"results_count": len(results),
				"search_url":    searchURL,
				"search_type":   searchType,
			},
		},
	}, nil
}

// buildSearchURL constructs a GitHub search URL
func (t *GitHubSearchTool) buildSearchURL(query, language, searchType, sortBy string, minStars int) string {
	// Build search query
	searchQuery := query

	// Add language filter
	if language != "" {
		searchQuery += fmt.Sprintf(" language:%s", language)
	}

	// Add stars filter
	if minStars > 0 {
		searchQuery += fmt.Sprintf(" stars:>%d", minStars)
	}

	// URL encode
	encodedQuery := strings.ReplaceAll(searchQuery, " ", "+")

	// Build URL based on search type
	var url string
	switch searchType {
	case "code":
		url = fmt.Sprintf("%s/search?type=code&q=%s", t.baseURL, encodedQuery)
	case "issues":
		url = fmt.Sprintf("%s/search?type=issues&q=%s", t.baseURL, encodedQuery)
	case "users":
		url = fmt.Sprintf("%s/search?type=users&q=%s", t.baseURL, encodedQuery)
	default: // repositories
		url = fmt.Sprintf("%s/search?type=repositories&q=%s", t.baseURL, encodedQuery)
	}

	// Add sort parameter
	if sortBy != "best-match" {
		url += fmt.Sprintf("&s=%s&o=desc", sortBy)
	}

	return url
}

// generateResults generates mock search results
// In production, this would call the GitHub API
func (t *GitHubSearchTool) generateResults(query, language, searchType string, minStars int) []types.Reference {
	results := make([]types.Reference, 0)

	switch searchType {
	case "repositories":
		results = t.generateRepositoryResults(query, language, minStars)
	case "code":
		results = t.generateCodeResults(query, language)
	case "issues":
		results = t.generateIssueResults(query, language)
	default:
		results = t.generateRepositoryResults(query, language, minStars)
	}

	return results
}

// generateRepositoryResults generates repository search results
func (t *GitHubSearchTool) generateRepositoryResults(query, language string, minStars int) []types.Reference {
	results := make([]types.Reference, 0)

	// Add popular repositories based on language
	if language != "" {
		results = append(results, t.getPopularRepos(language)...)
	}

	// Add query-specific results
	queryLower := strings.ToLower(query)
	
	if strings.Contains(queryLower, "example") || strings.Contains(queryLower, "tutorial") {
		results = append(results, types.Reference{
			Title:     fmt.Sprintf("%s Examples and Tutorials", language),
			URL:       fmt.Sprintf("%s/search?q=%s+examples+language:%s", t.baseURL, query, language),
			Type:      types.ReferenceTypeGitHub,
			Relevance: 0.9,
			Snippet:   fmt.Sprintf("Collection of %s examples and tutorials", language),
		})
	}

	if strings.Contains(queryLower, "framework") || strings.Contains(queryLower, "library") {
		results = append(results, types.Reference{
			Title:     fmt.Sprintf("Popular %s Frameworks", language),
			URL:       fmt.Sprintf("%s/search?q=%s+framework+language:%s&s=stars", t.baseURL, query, language),
			Type:      types.ReferenceTypeGitHub,
			Relevance: 0.85,
			Snippet:   fmt.Sprintf("Most starred %s frameworks and libraries", language),
		})
	}

	return results
}

// generateCodeResults generates code search results
func (t *GitHubSearchTool) generateCodeResults(query, language string) []types.Reference {
	results := make([]types.Reference, 0)

	results = append(results, types.Reference{
		Title:     fmt.Sprintf("Code examples for: %s", query),
		URL:       fmt.Sprintf("%s/search?type=code&q=%s+language:%s", t.baseURL, query, language),
		Type:      types.ReferenceTypeGitHub,
		Relevance: 0.9,
		Snippet:   fmt.Sprintf("Real-world code examples of %s in %s", query, language),
	})

	return results
}

// generateIssueResults generates issue search results
func (t *GitHubSearchTool) generateIssueResults(query, language string) []types.Reference {
	results := make([]types.Reference, 0)

	results = append(results, types.Reference{
		Title:     fmt.Sprintf("Issues and discussions: %s", query),
		URL:       fmt.Sprintf("%s/search?type=issues&q=%s+language:%s", t.baseURL, query, language),
		Type:      types.ReferenceTypeGitHub,
		Relevance: 0.8,
		Snippet:   fmt.Sprintf("GitHub issues and discussions about %s", query),
	})

	return results
}

// getPopularRepos returns popular repositories for a language
func (t *GitHubSearchTool) getPopularRepos(language string) []types.Reference {
	repos := make([]types.Reference, 0)

	switch language {
	case "go":
		repos = append(repos, types.Reference{
			Title:     "golang/go - The Go Programming Language",
			URL:       "https://github.com/golang/go",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 1.0,
			Snippet:   "Official Go programming language repository",
		})
		repos = append(repos, types.Reference{
			Title:     "avelino/awesome-go - Awesome Go",
			URL:       "https://github.com/avelino/awesome-go",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 0.95,
			Snippet:   "Curated list of awesome Go frameworks, libraries and software",
		})

	case "python":
		repos = append(repos, types.Reference{
			Title:     "python/cpython - Python Programming Language",
			URL:       "https://github.com/python/cpython",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 1.0,
			Snippet:   "Official Python programming language repository",
		})
		repos = append(repos, types.Reference{
			Title:     "vinta/awesome-python - Awesome Python",
			URL:       "https://github.com/vinta/awesome-python",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 0.95,
			Snippet:   "Curated list of awesome Python frameworks, libraries and software",
		})

	case "javascript":
		repos = append(repos, types.Reference{
			Title:     "nodejs/node - Node.js",
			URL:       "https://github.com/nodejs/node",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 1.0,
			Snippet:   "Official Node.js JavaScript runtime",
		})
		repos = append(repos, types.Reference{
			Title:     "sorrycc/awesome-javascript - Awesome JavaScript",
			URL:       "https://github.com/sorrycc/awesome-javascript",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 0.95,
			Snippet:   "Curated list of awesome JavaScript libraries and resources",
		})

	case "rust":
		repos = append(repos, types.Reference{
			Title:     "rust-lang/rust - The Rust Programming Language",
			URL:       "https://github.com/rust-lang/rust",
			Type:      types.ReferenceTypeGitHub,
			Relevance: 1.0,
			Snippet:   "Official Rust programming language repository",
		})
	}

	return repos
}

