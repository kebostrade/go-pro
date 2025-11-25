package errors

import (
	"fmt"
)

// ErrorCode represents an error code
type ErrorCode string

const (
	// Agent errors
	ErrCodeAgentExecution     ErrorCode = "AGENT_EXECUTION_ERROR"
	ErrCodeAgentTimeout       ErrorCode = "AGENT_TIMEOUT"
	ErrCodeAgentMaxSteps      ErrorCode = "AGENT_MAX_STEPS_EXCEEDED"
	ErrCodeAgentInvalidInput  ErrorCode = "AGENT_INVALID_INPUT"
	ErrCodeAgentNotFound      ErrorCode = "AGENT_NOT_FOUND"

	// LLM errors
	ErrCodeLLMProvider        ErrorCode = "LLM_PROVIDER_ERROR"
	ErrCodeLLMRateLimit       ErrorCode = "LLM_RATE_LIMIT"
	ErrCodeLLMInvalidResponse ErrorCode = "LLM_INVALID_RESPONSE"
	ErrCodeLLMTimeout         ErrorCode = "LLM_TIMEOUT"
	ErrCodeLLMQuotaExceeded   ErrorCode = "LLM_QUOTA_EXCEEDED"

	// Tool errors
	ErrCodeToolNotFound       ErrorCode = "TOOL_NOT_FOUND"
	ErrCodeToolExecution      ErrorCode = "TOOL_EXECUTION_ERROR"
	ErrCodeToolInvalidInput   ErrorCode = "TOOL_INVALID_INPUT"
	ErrCodeToolTimeout        ErrorCode = "TOOL_TIMEOUT"

	// Memory errors
	ErrCodeMemoryNotFound     ErrorCode = "MEMORY_NOT_FOUND"
	ErrCodeMemoryStorage      ErrorCode = "MEMORY_STORAGE_ERROR"
	ErrCodeMemoryRetrieval    ErrorCode = "MEMORY_RETRIEVAL_ERROR"

	// Vector store errors
	ErrCodeVectorStoreConn    ErrorCode = "VECTOR_STORE_CONNECTION_ERROR"
	ErrCodeVectorStoreQuery   ErrorCode = "VECTOR_STORE_QUERY_ERROR"
	ErrCodeVectorStoreInsert  ErrorCode = "VECTOR_STORE_INSERT_ERROR"

	// Workflow errors
	ErrCodeWorkflowExecution  ErrorCode = "WORKFLOW_EXECUTION_ERROR"
	ErrCodeWorkflowInvalid    ErrorCode = "WORKFLOW_INVALID"
	ErrCodeWorkflowTimeout    ErrorCode = "WORKFLOW_TIMEOUT"

	// General errors
	ErrCodeInvalidConfig      ErrorCode = "INVALID_CONFIG"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"
)

// AgentError represents an error in the agent platform
type AgentError struct {
	Code    ErrorCode
	Message string
	Details string
	Err     error
}

// Error implements the error interface
func (e *AgentError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AgentError) Unwrap() error {
	return e.Err
}

// New creates a new AgentError
func New(code ErrorCode, message string) *AgentError {
	return &AgentError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *AgentError {
	return &AgentError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails adds details to an error
func (e *AgentError) WithDetails(details string) *AgentError {
	e.Details = details
	return e
}

// Agent execution errors
func NewAgentExecutionError(message string, err error) *AgentError {
	return Wrap(err, ErrCodeAgentExecution, message)
}

func NewAgentTimeoutError(message string) *AgentError {
	return New(ErrCodeAgentTimeout, message)
}

func NewAgentMaxStepsError(maxSteps int) *AgentError {
	return New(ErrCodeAgentMaxSteps, fmt.Sprintf("maximum steps exceeded: %d", maxSteps))
}

func NewAgentInvalidInputError(message string) *AgentError {
	return New(ErrCodeAgentInvalidInput, message)
}

// LLM errors
func NewLLMProviderError(provider string, err error) *AgentError {
	return Wrap(err, ErrCodeLLMProvider, fmt.Sprintf("LLM provider error: %s", provider))
}

func NewLLMRateLimitError(provider string) *AgentError {
	return New(ErrCodeLLMRateLimit, fmt.Sprintf("rate limit exceeded for provider: %s", provider))
}

func NewLLMTimeoutError(provider string) *AgentError {
	return New(ErrCodeLLMTimeout, fmt.Sprintf("LLM request timeout: %s", provider))
}

// Tool errors
func NewToolNotFoundError(toolName string) *AgentError {
	return New(ErrCodeToolNotFound, fmt.Sprintf("tool not found: %s", toolName))
}

func NewToolExecutionError(toolName string, err error) *AgentError {
	return Wrap(err, ErrCodeToolExecution, fmt.Sprintf("tool execution failed: %s", toolName))
}

func NewToolInvalidInputError(toolName, message string) *AgentError {
	return New(ErrCodeToolInvalidInput, fmt.Sprintf("invalid input for tool %s: %s", toolName, message))
}

// Memory errors
func NewMemoryStorageError(err error) *AgentError {
	return Wrap(err, ErrCodeMemoryStorage, "failed to store in memory")
}

func NewMemoryRetrievalError(err error) *AgentError {
	return Wrap(err, ErrCodeMemoryRetrieval, "failed to retrieve from memory")
}

// Vector store errors
func NewVectorStoreConnectionError(provider string, err error) *AgentError {
	return Wrap(err, ErrCodeVectorStoreConn, fmt.Sprintf("vector store connection failed: %s", provider))
}

func NewVectorStoreQueryError(err error) *AgentError {
	return Wrap(err, ErrCodeVectorStoreQuery, "vector store query failed")
}

// Workflow errors
func NewWorkflowExecutionError(workflowName string, err error) *AgentError {
	return Wrap(err, ErrCodeWorkflowExecution, fmt.Sprintf("workflow execution failed: %s", workflowName))
}

func NewWorkflowInvalidError(message string) *AgentError {
	return New(ErrCodeWorkflowInvalid, message)
}

// General errors
func NewInvalidConfigError(message string) *AgentError {
	return New(ErrCodeInvalidConfig, message)
}

func NewUnauthorizedError(message string) *AgentError {
	return New(ErrCodeUnauthorized, message)
}

func NewNotFoundError(resource string) *AgentError {
	return New(ErrCodeNotFound, fmt.Sprintf("resource not found: %s", resource))
}

func NewInternalError(message string, err error) *AgentError {
	return Wrap(err, ErrCodeInternal, message)
}

