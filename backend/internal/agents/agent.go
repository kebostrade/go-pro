// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package agents provides multi-agent debugging and code analysis system
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/service"
)

// AgentType represents the type of agent in the system
type AgentType string

const (
	AgentTypeExecutor      AgentType = "executor"
	AgentTypeTestValidator AgentType = "test_validator"
	AgentTypeAIAnalysis    AgentType = "ai_analysis"
	AgentTypeStateManager  AgentType = "state_manager"
	AgentTypeCollaborative AgentType = "collaborative"
)

// Agent interface defines the contract for all agents
type Agent interface {
	GetType() AgentType
	Process(ctx context.Context, request AgentRequest) (AgentResponse, error)
	GetCapabilities() []string
}

// AgentRequest represents a generic request to an agent
type AgentRequest struct {
	Type       AgentType              `json:"type"`
	SessionID  string                 `json:"session_id"`
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	Context    map[string]interface{} `json:"context"`
	Priority   int                    `json:"priority"`
	Cancelable bool                   `json:"cancelable"`
}

// AgentResponse represents the response from an agent
type AgentResponse struct {
	Success   bool                   `json:"success"`
	AgentType AgentType              `json:"agent_type"`
	Data      map[string]interface{} `json:"data"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
}

// AgentPool manages multiple agents and routes requests
type AgentPool struct {
	agents        map[AgentType]Agent
	executor      *ExecutionAgent
	testValidator *TestValidationAgent
	aiAnalyzer    *AIAnalysisAgent
	stateManager  *StateManagerAgent
	mu            sync.RWMutex
}

// NewAgentPool creates a new agent pool with all specialized agents
func NewAgentPool() *AgentPool {
	ap := &AgentPool{
		agents: make(map[AgentType]Agent),
	}

	// Initialize all agents
	ap.executor = NewExecutionAgent()
	ap.testValidator = NewTestValidationAgent()
	ap.aiAnalyzer = NewAIAnalysisAgent()
	ap.stateManager = NewStateManagerAgent()

	// Register agents
	ap.agents[AgentTypeExecutor] = ap.executor
	ap.agents[AgentTypeTestValidator] = ap.testValidator
	ap.agents[AgentTypeAIAnalysis] = ap.aiAnalyzer
	ap.agents[AgentTypeStateManager] = ap.stateManager

	return ap
}

// GetAgent returns an agent by type
func (ap *AgentPool) GetAgent(agentType AgentType) (Agent, error) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	agent, ok := ap.agents[agentType]
	if !ok {
		return nil, fmt.Errorf("agent type %s not found", agentType)
	}
	return agent, nil
}

// ProcessRequest routes a request to the appropriate agent
func (ap *AgentPool) ProcessRequest(ctx context.Context, req AgentRequest) (AgentResponse, error) {
	startTime := time.Now()

	agent, err := ap.GetAgent(req.Type)
	if err != nil {
		return AgentResponse{
			Success:   false,
			Error:     err.Error(),
			Timestamp: startTime,
			Duration:  time.Since(startTime),
		}, err
	}

	resp, err := agent.Process(ctx, req)
	if err != nil {
		return AgentResponse{
			Success:   false,
			AgentType: req.Type,
			Error:     err.Error(),
			Timestamp: startTime,
			Duration:  time.Since(startTime),
		}, err
	}

	resp.Duration = time.Since(startTime)
	return resp, nil
}

// ProcessCollaborativeRequest processes a request involving multiple agents
func (ap *AgentPool) ProcessCollaborativeRequest(ctx context.Context, req AgentRequest) ([]AgentResponse, error) {
	var responses []AgentResponse

	// Determine which agents are needed based on request
	agentTypes := determineRequiredAgents(req)

	for _, agentType := range agentTypes {
		req.Type = agentType
		resp, err := ap.ProcessRequest(ctx, req)
		if err != nil {
			// Log but continue with other agents
			continue
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// determineRequiredAgents determines which agents are needed for a request
func determineRequiredAgents(req AgentRequest) []AgentType {
	var agents []AgentType

	// Always include state manager for session management
	agents = append(agents, AgentTypeStateManager)

	// Include executor for code execution
	if req.Data["code"] != nil {
		agents = append(agents, AgentTypeExecutor)
	}

	// Include test validator if test cases are provided
	if req.Data["test_cases"] != nil {
		agents = append(agents, AgentTypeTestValidator)
	}

	// Always include AI analysis for suggestions
	agents = append(agents, AgentTypeAIAnalysis)

	return agents
}

// ExecutionAgent handles code execution with sandboxing
type ExecutionAgent struct {
	timeLimit   time.Duration
	memoryLimit string
	mu          sync.RWMutex
}

// NewExecutionAgent creates a new execution agent
func NewExecutionAgent() *ExecutionAgent {
	return &ExecutionAgent{
		timeLimit:   10 * time.Second,
		memoryLimit: "256m",
	}
}

// GetType returns the agent type
func (e *ExecutionAgent) GetType() AgentType {
	return AgentTypeExecutor
}

// Process handles code execution requests
func (e *ExecutionAgent) Process(ctx context.Context, req AgentRequest) (AgentResponse, error) {
	code, ok := req.Data["code"].(string)
	if !ok {
		return AgentResponse{Success: false, Error: "invalid code format"}, fmt.Errorf("invalid code format")
	}

	language, _ := req.Data["language"].(string)
	if language == "" {
		language = "go"
	}

	// Get test cases if provided
	var testCases []service.TestCase
	if tcData, ok := req.Data["test_cases"]; ok {
		if tcJSON, ok := tcData.(string); ok {
			json.Unmarshal([]byte(tcJSON), &testCases)
		}
	}

	// Execute with executor service
	execReq := &service.ExecuteRequest{
		Code:      code,
		Language:  language,
		Timeout:   e.timeLimit,
		TestCases: testCases,
	}

	result, err := executeCodeWithLimits(ctx, execReq)
	if err != nil {
		return AgentResponse{
			Success:   false,
			AgentType: AgentTypeExecutor,
			Error:     err.Error(),
			Data:      map[string]interface{}{"code": code},
		}, err
	}

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeExecutor,
		Data: map[string]interface{}{
			"output":         result.Output,
			"error":          result.Error,
			"execution_time": result.ExecutionTime,
			"memory_used":    result.MemoryUsed,
			"passed":         result.Passed,
			"score":          result.Score,
			"test_results":   result.Results,
			"breakpoints":    result.Breakpoints,
			"variables":      result.Variables,
			"call_stack":     result.CallStack,
			"stream_events":  result.StreamEvents,
		},
		Metadata: map[string]interface{}{
			"time_limit":   e.timeLimit.String(),
			"memory_limit": e.memoryLimit,
			"language":     language,
		},
	}, nil
}

// executeCodeWithLimits executes code with timeout and memory constraints
func executeCodeWithLimits(ctx context.Context, req *service.ExecuteRequest) (*ExecutionResult, error) {
	// This would integrate with the actual Docker executor
	// For now, return a placeholder result
	return &ExecutionResult{
		Output:        "",
		Error:         "",
		ExecutionTime: 0,
		MemoryUsed:    0,
		Passed:        false,
		Score:         0,
		Results:       []service.TestResult{},
		Breakpoints:   []BreakpointInfo{},
		Variables:     []VariableInfo{},
		CallStack:     []StackFrame{},
		StreamEvents:  []StreamEvent{},
	}, nil
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Output        string               `json:"output"`
	Error         string               `json:"error,omitempty"`
	ExecutionTime int64                `json:"execution_time_ms"`
	MemoryUsed    int64                `json:"memory_used_bytes"`
	Passed        bool                 `json:"passed"`
	Score         int                  `json:"score"`
	Results       []service.TestResult `json:"test_results"`
	Breakpoints   []BreakpointInfo     `json:"breakpoints,omitempty"`
	Variables     []VariableInfo       `json:"variables,omitempty"`
	CallStack     []StackFrame         `json:"call_stack,omitempty"`
	StreamEvents  []StreamEvent        `json:"stream_events,omitempty"`
}

// BreakpointInfo represents breakpoint information
type BreakpointInfo struct {
	ID        string `json:"id"`
	Line      int    `json:"line"`
	Condition string `json:"condition,omitempty"`
	Hits      int    `json:"hits"`
	Enabled   bool   `json:"enabled"`
}

// VariableInfo represents variable state during debugging
type VariableInfo struct {
	Name     string         `json:"name"`
	Value    interface{}    `json:"value"`
	Type     string         `json:"type"`
	Scope    string         `json:"scope"`
	Address  string         `json:"address,omitempty"`
	Children []VariableInfo `json:"children,omitempty"`
}

// StackFrame represents a call stack frame
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// StreamEvent represents a real-time execution event
type StreamEvent struct {
	Type string      `json:"type"`
	Time int64       `json:"timestamp"`
	Data interface{} `json:"data"`
	Line int         `json:"line,omitempty"`
}

// GetCapabilities returns the agent's capabilities
func (e *ExecutionAgent) GetCapabilities() []string {
	return []string{
		"code_execution",
		"sandboxed_execution",
		"breakpoint_debugging",
		"variable_inspection",
		"call_stack_tracing",
		"stream_output",
		"memory_profiling",
	}
}

// TestValidationAgent handles test case validation
type TestValidationAgent struct {
	mu sync.RWMutex
}

// NewTestValidationAgent creates a new test validation agent
func NewTestValidationAgent() *TestValidationAgent {
	return &TestValidationAgent{}
}

// GetType returns the agent type
func (t *TestValidationAgent) GetType() AgentType {
	return AgentTypeTestValidator
}

// Process handles test case validation requests
func (t *TestValidationAgent) Process(ctx context.Context, req AgentRequest) (AgentResponse, error) {
	testCases, ok := req.Data["test_cases"].([]interface{})
	if !ok {
		return AgentResponse{Success: false, Error: "invalid test cases format"}, fmt.Errorf("invalid test cases format")
	}

	// Parse and validate test cases
	validatedCases, validationResults := t.validateTestCases(testCases)

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeTestValidator,
		Data: map[string]interface{}{
			"validated_cases":    validatedCases,
			"validation_results": validationResults,
			"edge_cases":         t.detectEdgeCases(validatedCases),
			"coverage_hints":     t.suggestCoverage(validatedCases),
		},
	}, nil
}

// validateTestCases validates test case structure and content
func (t *TestValidationAgent) validateTestCases(cases []interface{}) ([]map[string]interface{}, []map[string]interface{}) {
	validated := make([]map[string]interface{}, 0, len(cases))
	results := make([]map[string]interface{}, 0, len(cases))

	for i, tc := range cases {
		tcMap, ok := tc.(map[string]interface{})
		if !ok {
			results = append(results, map[string]interface{}{
				"index": i, "valid": false, "error": "invalid test case format",
			})
			continue
		}

		// Validate required fields
		valid := true
		var validationErrors []string

		if tcMap["input"] == nil {
			valid = false
			validationErrors = append(validationErrors, "missing required field: input")
		}
		if tcMap["expected"] == nil {
			valid = false
			validationErrors = append(validationErrors, "missing required field: expected")
		}

		if valid {
			validated = append(validated, map[string]interface{}{
				"name":     tcMap["name"],
				"input":    tcMap["input"],
				"expected": tcMap["expected"],
				"hidden":   tcMap["hidden"],
			})
		}

		results = append(results, map[string]interface{}{
			"index":  i,
			"valid":  valid,
			"errors": validationErrors,
			"name":   tcMap["name"],
		})
	}

	return validated, results
}

// detectEdgeCases identifies potential edge cases from test cases
func (t *TestValidationAgent) detectEdgeCases(cases []map[string]interface{}) []map[string]interface{} {
	var edgeCases []map[string]interface{}

	// Detect common edge cases
	hasEmptyInput := false
	hasNullInput := false
	hasLargeInput := false
	hasNegativeNumbers := false

	for _, tc := range cases {
		input, ok := tc["input"].(string)
		if !ok {
			continue
		}

		if input == "" {
			hasEmptyInput = true
		}
		if input == "null" || input == "nil" {
			hasNullInput = true
		}
		if len(input) > 10000 {
			hasLargeInput = true
		}
		if len(input) > 0 && input[0] == '-' {
			hasNegativeNumbers = true
		}
	}

	if !hasEmptyInput {
		edgeCases = append(edgeCases, map[string]interface{}{
			"type":        "empty_input",
			"description": "Consider adding an empty input test case",
			"priority":    "high",
		})
	}

	if !hasNullInput {
		edgeCases = append(edgeCases, map[string]interface{}{
			"type":        "null_input",
			"description": "Consider adding a null/nil input test case",
			"priority":    "medium",
		})
	}

	if !hasLargeInput {
		edgeCases = append(edgeCases, map[string]interface{}{
			"type":        "large_input",
			"description": "Consider adding a large input test case",
			"priority":    "medium",
		})
	}

	if !hasNegativeNumbers {
		edgeCases = append(edgeCases, map[string]interface{}{
			"type":        "negative_numbers",
			"description": "Consider adding negative number test cases",
			"priority":    "medium",
		})
	}

	return edgeCases
}

// suggestCoverage provides coverage suggestions
func (t *TestValidationAgent) suggestCoverage(cases []map[string]interface{}) []string {
	suggestions := []string{}

	if len(cases) < 3 {
		suggestions = append(suggestions, "Add more test cases for better coverage")
	}

	// Check for boundary conditions
	hasBoundary := false
	for _, tc := range cases {
		input, ok := tc["input"].(string)
		if !ok {
			continue
		}
		if input == "0" || input == "1" || input == "-1" || input == "max" {
			hasBoundary = true
			break
		}
	}

	if !hasBoundary {
		suggestions = append(suggestions, "Add boundary condition test cases (0, 1, -1, max values)")
	}

	return suggestions
}

// GetCapabilities returns the agent's capabilities
func (t *TestValidationAgent) GetCapabilities() []string {
	return []string{
		"test_case_validation",
		"edge_case_detection",
		"coverage_suggestions",
		"test_case_generation",
	}
}

// AIAnalysisAgent handles AI-powered code analysis and suggestions
type AIAnalysisAgent struct {
	mu sync.RWMutex
}

// NewAIAnalysisAgent creates a new AI analysis agent
func NewAIAnalysisAgent() *AIAnalysisAgent {
	return &AIAnalysisAgent{}
}

// GetType returns the agent type
func (a *AIAnalysisAgent) GetType() AgentType {
	return AgentTypeAIAnalysis
}

// Process handles AI analysis requests
func (a *AIAnalysisAgent) Process(ctx context.Context, req AgentRequest) (AgentResponse, error) {
	code, ok := req.Data["code"].(string)
	if !ok {
		return AgentResponse{Success: false, Error: "invalid code format"}, fmt.Errorf("invalid code format")
	}

	// Perform AI analysis
	analysis := a.analyzeCode(code)

	// Get context from previous executions
	var executionContext map[string]interface{}
	if ctxData, ok := req.Context["execution"]; ok {
		executionContext, _ = ctxData.(map[string]interface{})
	}

	// Generate suggestions
	suggestions := a.generateSuggestions(code, analysis, executionContext)

	// Generate explanations
	explanations := a.explainErrors(code, req.Data["error"])

	// Generate optimized code
	optimizedCode := a.suggestOptimizations(code, analysis)

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeAIAnalysis,
		Data: map[string]interface{}{
			"analysis":        analysis,
			"suggestions":     suggestions,
			"explanations":    explanations,
			"optimized_code":  optimizedCode,
			"completion":      a.suggestCompletions(code),
			"test_generation": a.suggestTestCases(code),
		},
	}, nil
}

// CodeAnalysisResult represents the result of AI code analysis
type CodeAnalysisResult struct {
	Complexity      int            `json:"complexity"`
	Issues          []CodeIssue    `json:"issues"`
	Patterns        []string       `json:"patterns"`
	Strengths       []string       `json:"strengths"`
	LanguageVersion string         `json:"language_version"`
	Imports         []string       `json:"imports"`
	Functions       []FunctionInfo `json:"functions"`
	Variables       []VariableDecl `json:"variables"`
}

// CodeIssue represents a code issue found by AI
type CodeIssue struct {
	Type     string `json:"type"` // error, warning, suggestion
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Rule     string `json:"rule,omitempty"`
}

// FunctionInfo represents function information
type FunctionInfo struct {
	Name       string   `json:"name"`
	Parameters []string `json:"parameters"`
	Returns    []string `json:"returns"`
	Line       int      `json:"line"`
}

// VariableDecl represents a variable declaration
type VariableDecl struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Line int    `json:"line"`
}

// analyzeCode performs AI-powered code analysis
func (a *AIAnalysisAgent) analyzeCode(code string) *CodeAnalysisResult {
	result := &CodeAnalysisResult{
		Complexity:      1,
		Issues:          []CodeIssue{},
		Patterns:        []string{},
		Strengths:       []string{},
		LanguageVersion: "1.21",
		Imports:         []string{},
		Functions:       []FunctionInfo{},
		Variables:       []VariableDecl{},
	}

	// Analyze complexity (simplified)
	if len(code) > 500 {
		result.Complexity += 1
	}

	// Detect issues
	if !containsString(code, "func main()") {
		result.Issues = append(result.Issues, CodeIssue{
			Type:     "suggestion",
			Severity: "info",
			Message:  "Consider adding a main() function for executable code",
			Line:     1,
		})
	}

	if containsString(code, "fmt.Println") {
		result.Patterns = append(result.Patterns, "output_operation")
	}

	if containsString(code, "go func") {
		result.Patterns = append(result.Patterns, "goroutine")
		result.Issues = append(result.Issues, CodeIssue{
			Type:     "warning",
			Severity: "medium",
			Message:  "Goroutines detected - ensure proper synchronization",
			Rule:     "goroutine_warning",
		})
	}

	// Check for error handling
	if !containsString(code, "if err") && containsString(code, ",") {
		result.Issues = append(result.Issues, CodeIssue{
			Type:     "suggestion",
			Severity: "low",
			Message:  "Consider adding error handling",
			Rule:     "error_handling",
		})
	}

	// Detect strengths
	if containsString(code, "interface{}") {
		result.Strengths = append(result.Strengths, "Generic programming patterns")
	}

	return result
}

// generateSuggestions generates AI-powered code suggestions
func (a *AIAnalysisAgent) generateSuggestions(code string, analysis *CodeAnalysisResult, context map[string]interface{}) []map[string]interface{} {
	var suggestions []map[string]interface{}

	for _, issue := range analysis.Issues {
		if issue.Type == "suggestion" || issue.Type == "warning" {
			suggestions = append(suggestions, map[string]interface{}{
				"type":    issue.Type,
				"message": issue.Message,
				"line":    issue.Line,
				"rule":    issue.Rule,
				"action":  a.suggestFix(issue),
			})
		}
	}

	// Context-aware suggestions based on test cases
	if context != nil {
		if _, ok := context["test_cases"]; ok {
			suggestions = append(suggestions, map[string]interface{}{
				"type":    "context",
				"message": "Based on your test cases, consider handling edge cases",
				"action":  "add_edge_case_handling",
			})
		}
	}

	return suggestions
}

// suggestFix suggests a fix for an issue
func (a *AIAnalysisAgent) suggestFix(issue CodeIssue) string {
	switch issue.Rule {
	case "goroutine_warning":
		return "Use sync.WaitGroup or channels to manage goroutine lifecycle"
	case "error_handling":
		return "Add error checking: if err != nil { return err }"
	default:
		return "Review and fix the code at line " + strconv.Itoa(issue.Line)
	}
}

// explainErrors explains error messages in natural language
func (a *AIAnalysisAgent) explainErrors(code string, errData interface{}) []map[string]interface{} {
	var explanations []map[string]interface{}

	if errData == nil {
		return explanations
	}

	errStr, ok := errData.(string)
	if !ok {
		return explanations
	}

	// Common Go errors and their explanations
	errorExplanations := map[string]string{
		"undefined":            "This error indicates a variable or function that hasn't been declared or is out of scope. A variable or function name is not recognized in the current scope",
		"syntax error":         "There's a syntax error in your code - check for missing brackets, parentheses, or semicolons",
		"cannot use":           "Type mismatch - the types on both sides of the assignment are incompatible",
		"not enough arguments": "Function call is missing required arguments",
		"too many arguments":   "Function call has more arguments than expected",
		"panic:":               "A runtime panic occurred - this usually happens from nil pointer dereference or index out of bounds",
	}

	for errKey, explanation := range errorExplanations {
		if containsString(errStr, errKey) {
			explanations = append(explanations, map[string]interface{}{
				"error":          errStr,
				"explanation":    explanation,
				"severity":       "high",
				"fix_suggestion": a.suggestFixForError(errKey),
			})
		}
	}

	return explanations
}

// suggestFixForError suggests a fix for a specific error
func (a *AIAnalysisAgent) suggestFixForError(errorType string) string {
	fixes := map[string]string{
		"undefined":            "Check spelling and ensure the variable is declared before use",
		"syntax error":         "Review the line for missing or extra characters",
		"cannot use":           "Ensure types are compatible or use type conversion",
		"not enough arguments": "Add missing arguments to the function call",
		"too many arguments":   "Remove extra arguments from the function call",
		"panic":                "Add nil checks and boundary validation",
	}
	return fixes[errorType]
}

// suggestOptimizations suggests optimized code
func (a *AIAnalysisAgent) suggestOptimizations(code string, analysis *CodeAnalysisResult) string {
	var optimized strings.Builder
	optimized.WriteString(code)

	// Suggest string concatenation optimization
	if containsString(code, "+") && containsString(code, "string") {
		// Could suggest strings.Builder
	}

	return optimized.String()
}

// suggestCompletions suggests code completions
func (a *AIAnalysisAgent) suggestCompletions(code string) []map[string]interface{} {
	var completions []map[string]interface{}

	// Context-aware completions
	lastLine := getLastLine(code)

	if containsString(lastLine, "func") {
		completions = append(completions, map[string]interface{}{
			"label":      "main()",
			"insertText": "func main() {\n\t$0\n}",
			"detail":     "Entry point function",
			"kind":       "function",
		})
	}

	if containsString(lastLine, "for") {
		completions = append(completions, map[string]interface{}{
			"label":      "range loop",
			"insertText": "range ${1:collection} {\n\t$0\n}",
			"detail":     "Range-based for loop",
			"kind":       "snippet",
		})
	}

	return completions
}

// suggestTestCases suggests test cases based on code
func (a *AIAnalysisAgent) suggestTestCases(code string) []map[string]interface{} {
	var suggestions []map[string]interface{}

	// Analyze function signatures
	if containsString(code, "func add(") {
		suggestions = append(suggestions, map[string]interface{}{
			"name":     "test_add_basic",
			"input":    "2, 3",
			"expected": "5",
			"type":     "basic",
		}, map[string]interface{}{
			"name":     "test_add_negative",
			"input":    "-1, -2",
			"expected": "-3",
			"type":     "edge_case",
		}, map[string]interface{}{
			"name":     "test_add_zero",
			"input":    "0, 5",
			"expected": "5",
			"type":     "edge_case",
		})
	}

	return suggestions
}

// GetCapabilities returns the agent's capabilities
func (a *AIAnalysisAgent) GetCapabilities() []string {
	return []string{
		"code_analysis",
		"bug_detection",
		"optimization_suggestions",
		"error_explanation",
		"code_completion",
		"test_generation",
	}
}

// StateManagerAgent handles persistent state management
type StateManagerAgent struct {
	sessions    map[string]*Session
	history     map[string][]CodeHistoryEntry
	preferences map[string]UserPreferences
	mu          sync.RWMutex
}

// Session represents a user session
type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Data      map[string]interface{} `json:"data"`
	Agents    []AgentType            `json:"agents"`
}

// CodeHistoryEntry represents a code history entry
type CodeHistoryEntry struct {
	ID        string            `json:"id"`
	Code      string            `json:"code"`
	Language  string            `json:"language"`
	Timestamp time.Time         `json:"timestamp"`
	Output    string            `json:"output,omitempty"`
	Error     string            `json:"error,omitempty"`
	TestCases []domain.TestCase `json:"test_cases,omitempty"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	Theme           string   `json:"theme"`
	FontSize        int      `json:"font_size"`
	AutoSave        bool     `json:"auto_save"`
	Language        string   `json:"language"`
	DebuggerEnabled bool     `json:"debugger_enabled"`
	Collaborators   []string `json:"collaborators"`
}

// NewStateManagerAgent creates a new state manager agent
func NewStateManagerAgent() *StateManagerAgent {
	return &StateManagerAgent{
		sessions:    make(map[string]*Session),
		history:     make(map[string][]CodeHistoryEntry),
		preferences: make(map[string]UserPreferences),
	}
}

// GetType returns the agent type
func (s *StateManagerAgent) GetType() AgentType {
	return AgentTypeStateManager
}

// Process handles state management requests
func (s *StateManagerAgent) Process(ctx context.Context, req AgentRequest) (AgentResponse, error) {
	action, _ := req.Data["action"].(string)

	switch action {
	case "create_session":
		return s.createSession(req)
	case "get_session":
		return s.getSession(req)
	case "update_session":
		return s.updateSession(req)
	case "save_code":
		return s.saveCode(req)
	case "get_history":
		return s.getHistory(req)
	case "save_preferences":
		return s.savePreferences(req)
	case "get_preferences":
		return s.getPreferences(req)
	case "share_session":
		return s.shareSession(req)
	default:
		return AgentResponse{Success: false, Error: "unknown action"}, fmt.Errorf("unknown action: %s", action)
	}
}

// createSession creates a new session
func (s *StateManagerAgent) createSession(req AgentRequest) (AgentResponse, error) {
	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Data:      make(map[string]interface{}),
		Agents:    []AgentType{},
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.history[sessionID] = []CodeHistoryEntry{}
	s.mu.Unlock()

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"session":    session,
		},
	}, nil
}

// getSession retrieves a session
func (s *StateManagerAgent) getSession(req AgentRequest) (AgentResponse, error) {
	sessionID, _ := req.Data["session_id"].(string)

	s.mu.RLock()
	session, ok := s.sessions[sessionID]
	s.mu.RUnlock()

	if !ok {
		return AgentResponse{Success: false, Error: "session not found"}, fmt.Errorf("session not found")
	}

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"session": session,
		},
	}, nil
}

// updateSession updates a session
func (s *StateManagerAgent) updateSession(req AgentRequest) (AgentResponse, error) {
	sessionID, _ := req.Data["session_id"].(string)
	data, _ := req.Data["data"].(map[string]interface{})

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return AgentResponse{Success: false, Error: "session not found"}, fmt.Errorf("session not found")
	}

	for k, v := range data {
		session.Data[k] = v
	}
	session.UpdatedAt = time.Now()

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"session": session,
		},
	}, nil
}

// saveCode saves code to history
func (s *StateManagerAgent) saveCode(req AgentRequest) (AgentResponse, error) {
	sessionID, _ := req.Data["session_id"].(string)
	code, _ := req.Data["code"].(string)
	language, _ := req.Data["language"].(string)
	output, _ := req.Data["output"].(string)
	testCasesJSON, _ := req.Data["test_cases"].(string)

	entry := CodeHistoryEntry{
		ID:        generateEntryID(),
		Code:      code,
		Language:  language,
		Timestamp: time.Now(),
		Output:    output,
	}

	if testCasesJSON != "" {
		json.Unmarshal([]byte(testCasesJSON), &entry.TestCases)
	}

	s.mu.Lock()
	if history, ok := s.history[sessionID]; ok {
		// Keep only last 50 entries
		if len(history) >= 50 {
			s.history[sessionID] = append([]CodeHistoryEntry{entry}, history[:49]...)
		} else {
			s.history[sessionID] = append([]CodeHistoryEntry{entry}, history...)
		}
	}
	s.mu.Unlock()

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"entry_id": entry.ID,
		},
	}, nil
}

// getHistory retrieves code history
func (s *StateManagerAgent) getHistory(req AgentRequest) (AgentResponse, error) {
	sessionID, _ := req.Data["session_id"].(string)
	limit, _ := req.Data["limit"].(float64)
	if limit == 0 {
		limit = 10
	}

	s.mu.RLock()
	history := s.history[sessionID]
	s.mu.RUnlock()

	if len(history) > int(limit) {
		history = history[:int(limit)]
	}

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"history": history,
		},
	}, nil
}

// savePreferences saves user preferences
func (s *StateManagerAgent) savePreferences(req AgentRequest) (AgentResponse, error) {
	prefsJSON, _ := req.Data["preferences"].(string)
	var prefs UserPreferences
	json.Unmarshal([]byte(prefsJSON), &prefs)

	s.mu.Lock()
	s.preferences[req.UserID] = prefs
	s.mu.Unlock()

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
	}, nil
}

// getPreferences retrieves user preferences
func (s *StateManagerAgent) getPreferences(req AgentRequest) (AgentResponse, error) {
	s.mu.RLock()
	prefs, ok := s.preferences[req.UserID]
	s.mu.RUnlock()

	if !ok {
		prefs = UserPreferences{
			Theme:    "dark",
			FontSize: 14,
			AutoSave: true,
		}
	}

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"preferences": prefs,
		},
	}, nil
}

// shareSession shares a session with collaborators
func (s *StateManagerAgent) shareSession(req AgentRequest) (AgentResponse, error) {
	sessionID, _ := req.Data["session_id"].(string)
	collaborators, _ := req.Data["collaborators"].([]interface{})

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return AgentResponse{Success: false, Error: "session not found"}, fmt.Errorf("session not found")
	}

	for _, c := range collaborators {
		if collab, ok := c.(string); ok {
			session.Data["collaborator_"+collab] = true
		}
	}

	return AgentResponse{
		Success:   true,
		AgentType: AgentTypeStateManager,
		Data: map[string]interface{}{
			"shared_with": collaborators,
		},
	}, nil
}

// GetCapabilities returns the agent's capabilities
func (s *StateManagerAgent) GetCapabilities() []string {
	return []string{
		"session_management",
		"code_history",
		"user_preferences",
		"collaboration",
	}
}

// Helper functions
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || containsStringHelper(s, substr)))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func getLastLine(code string) string {
	lines := strings.Split(code, "\n")
	if len(lines) > 0 {
		return lines[len(lines)-1]
	}
	return ""
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

func generateEntryID() string {
	return fmt.Sprintf("entry_%d", time.Now().UnixNano())
}
