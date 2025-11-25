package types

import (
	"time"
)

// CodingRequest represents a request to a coding agent
type CodingRequest struct {
	// Query is the programming question or request
	Query string `json:"query"`

	// Language specifies the programming language (go, python, javascript, etc.)
	Language string `json:"language,omitempty"`

	// Code is the code snippet to analyze/debug/review
	Code string `json:"code,omitempty"`

	// Context provides additional context (error messages, requirements, etc.)
	Context map[string]interface{} `json:"context,omitempty"`

	// RequestType specifies the type of request
	RequestType CodingRequestType `json:"request_type"`

	// Difficulty level for generated examples
	Difficulty DifficultyLevel `json:"difficulty,omitempty"`

	// IncludeTests whether to include test cases
	IncludeTests bool `json:"include_tests,omitempty"`

	// IncludeDocs whether to include documentation
	IncludeDocs bool `json:"include_docs,omitempty"`

	// SessionID for tracking conversation
	SessionID string `json:"session_id,omitempty"`

	// UserID for personalization
	UserID string `json:"user_id,omitempty"`
}

// CodingRequestType defines types of coding requests
type CodingRequestType string

const (
	RequestTypeQuestion    CodingRequestType = "question"     // General Q&A
	RequestTypeDebug       CodingRequestType = "debug"        // Debug code
	RequestTypeReview      CodingRequestType = "review"       // Code review
	RequestTypeRefactor    CodingRequestType = "refactor"     // Refactoring suggestions
	RequestTypeGenerate    CodingRequestType = "generate"     // Generate code
	RequestTypeExplain     CodingRequestType = "explain"      // Explain code
	RequestTypeOptimize    CodingRequestType = "optimize"     // Performance optimization
	RequestTypeArchitect   CodingRequestType = "architect"    // Architecture guidance
	RequestTypeTest        CodingRequestType = "test"         // Generate tests
	RequestTypeDocument    CodingRequestType = "document"     // Generate documentation
)

// DifficultyLevel represents code complexity level
type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyExpert       DifficultyLevel = "expert"
)

// CodingResponse represents a response from a coding agent
type CodingResponse struct {
	// Answer is the main response
	Answer string `json:"answer"`

	// Code contains generated or modified code
	Code *CodeSnippet `json:"code,omitempty"`

	// Explanation provides detailed explanation
	Explanation string `json:"explanation,omitempty"`

	// Examples contains code examples
	Examples []CodeExample `json:"examples,omitempty"`

	// Suggestions for improvements or alternatives
	Suggestions []string `json:"suggestions,omitempty"`

	// References to documentation or resources
	References []Reference `json:"references,omitempty"`

	// ExecutionResult if code was executed
	ExecutionResult *ExecutionResult `json:"execution_result,omitempty"`

	// Analysis results (for review/debug requests)
	Analysis *CodeAnalysis `json:"analysis,omitempty"`

	// Metadata about the response
	Metadata CodingMetadata `json:"metadata"`
}

// CodeSnippet represents a code snippet
type CodeSnippet struct {
	// Language of the code
	Language string `json:"language"`

	// Code content
	Code string `json:"code"`

	// FileName suggested filename
	FileName string `json:"file_name,omitempty"`

	// Description of what the code does
	Description string `json:"description,omitempty"`

	// Tests associated test cases
	Tests string `json:"tests,omitempty"`

	// Documentation inline documentation
	Documentation string `json:"documentation,omitempty"`
}

// CodeExample represents a code example
type CodeExample struct {
	// Title of the example
	Title string `json:"title"`

	// Description of what it demonstrates
	Description string `json:"description"`

	// Code the example code
	Code string `json:"code"`

	// Language programming language
	Language string `json:"language"`

	// Output expected output
	Output string `json:"output,omitempty"`

	// Difficulty level
	Difficulty DifficultyLevel `json:"difficulty"`
}

// Reference represents a documentation or resource reference
type Reference struct {
	// Title of the reference
	Title string `json:"title"`

	// URL to the resource
	URL string `json:"url"`

	// Type of reference (official_docs, stackoverflow, github, blog, etc.)
	Type ReferenceType `json:"type"`

	// Relevance score (0-1)
	Relevance float64 `json:"relevance,omitempty"`

	// Snippet relevant excerpt
	Snippet string `json:"snippet,omitempty"`
}

// ReferenceType defines types of references
type ReferenceType string

const (
	ReferenceTypeOfficialDocs ReferenceType = "official_docs"
	ReferenceTypeStackOverflow ReferenceType = "stackoverflow"
	ReferenceTypeGitHub        ReferenceType = "github"
	ReferenceTypeBlog          ReferenceType = "blog"
	ReferenceTypeTutorial      ReferenceType = "tutorial"
	ReferenceTypeBook          ReferenceType = "book"
)

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	// Success whether execution was successful
	Success bool `json:"success"`

	// Output from stdout
	Output string `json:"output,omitempty"`

	// Error from stderr
	Error string `json:"error,omitempty"`

	// ExitCode process exit code
	ExitCode int `json:"exit_code"`

	// ExecutionTime in milliseconds
	ExecutionTime int64 `json:"execution_time_ms"`

	// MemoryUsed in bytes
	MemoryUsed int64 `json:"memory_used_bytes,omitempty"`

	// CPUTime in milliseconds
	CPUTime int64 `json:"cpu_time_ms,omitempty"`
}

// CodeAnalysis represents code analysis results
type CodeAnalysis struct {
	// Issues found in the code
	Issues []CodeIssue `json:"issues,omitempty"`

	// Metrics code quality metrics
	Metrics CodeMetrics `json:"metrics"`

	// SecurityIssues security vulnerabilities
	SecurityIssues []SecurityIssue `json:"security_issues,omitempty"`

	// PerformanceIssues performance problems
	PerformanceIssues []PerformanceIssue `json:"performance_issues,omitempty"`

	// BestPractices violations of best practices
	BestPractices []BestPracticeViolation `json:"best_practices,omitempty"`

	// OverallScore quality score (0-100)
	OverallScore int `json:"overall_score"`
}

// CodeIssue represents a code issue
type CodeIssue struct {
	// Severity of the issue
	Severity IssueSeverity `json:"severity"`

	// Type of issue
	Type IssueType `json:"type"`

	// Message description of the issue
	Message string `json:"message"`

	// Line line number where issue occurs
	Line int `json:"line,omitempty"`

	// Column column number
	Column int `json:"column,omitempty"`

	// Suggestion how to fix
	Suggestion string `json:"suggestion,omitempty"`

	// Code snippet showing the issue
	CodeSnippet string `json:"code_snippet,omitempty"`
}

// IssueSeverity defines severity levels
type IssueSeverity string

const (
	SeverityError   IssueSeverity = "error"
	SeverityWarning IssueSeverity = "warning"
	SeverityInfo    IssueSeverity = "info"
	SeverityHint    IssueSeverity = "hint"
)

// IssueType defines types of issues
type IssueType string

const (
	IssueTypeSyntax       IssueType = "syntax"
	IssueTypeLogic        IssueType = "logic"
	IssueTypeSecurity     IssueType = "security"
	IssueTypePerformance  IssueType = "performance"
	IssueTypeStyle        IssueType = "style"
	IssueTypeBestPractice IssueType = "best_practice"
	IssueTypeDeprecated   IssueType = "deprecated"
)

// CodeMetrics represents code quality metrics
type CodeMetrics struct {
	// LinesOfCode total lines
	LinesOfCode int `json:"lines_of_code"`

	// Complexity cyclomatic complexity
	Complexity int `json:"complexity"`

	// Maintainability index (0-100)
	Maintainability int `json:"maintainability"`

	// TestCoverage percentage (0-100)
	TestCoverage int `json:"test_coverage,omitempty"`

	// CommentRatio percentage of comments
	CommentRatio float64 `json:"comment_ratio"`
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	// Type of vulnerability
	Type string `json:"type"`

	// Severity level
	Severity IssueSeverity `json:"severity"`

	// Description of the vulnerability
	Description string `json:"description"`

	// Recommendation how to fix
	Recommendation string `json:"recommendation"`

	// CWE Common Weakness Enumeration ID
	CWE string `json:"cwe,omitempty"`

	// Line where vulnerability occurs
	Line int `json:"line,omitempty"`
}

// PerformanceIssue represents a performance problem
type PerformanceIssue struct {
	// Type of performance issue
	Type string `json:"type"`

	// Description of the issue
	Description string `json:"description"`

	// Impact estimated impact
	Impact string `json:"impact"`

	// Suggestion how to improve
	Suggestion string `json:"suggestion"`

	// Line where issue occurs
	Line int `json:"line,omitempty"`
}

// BestPracticeViolation represents a best practice violation
type BestPracticeViolation struct {
	// Practice that was violated
	Practice string `json:"practice"`

	// Description of the violation
	Description string `json:"description"`

	// Recommendation how to follow best practice
	Recommendation string `json:"recommendation"`

	// Line where violation occurs
	Line int `json:"line,omitempty"`
}

// CodingMetadata contains metadata about the coding response
type CodingMetadata struct {
	// RequestID unique identifier
	RequestID string `json:"request_id"`

	// AgentType which agent handled the request
	AgentType string `json:"agent_type"`

	// Language detected or specified language
	Language string `json:"language"`

	// ProcessingTime in milliseconds
	ProcessingTime int64 `json:"processing_time_ms"`

	// TokensUsed LLM tokens consumed
	TokensUsed TokenUsage `json:"tokens_used"`

	// ToolsUsed which tools were invoked
	ToolsUsed []string `json:"tools_used,omitempty"`

	// Timestamp when response was generated
	Timestamp time.Time `json:"timestamp"`

	// Confidence score (0-1)
	Confidence float64 `json:"confidence,omitempty"`
}

