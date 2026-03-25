// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
)

// PlaygroundExecuteRequest represents a playground code execution request.
type PlaygroundExecuteRequest struct {
	Code string `json:"code" validate:"required"`
}

// PlaygroundExecuteResult represents the result of playground code execution.
type PlaygroundExecuteResult struct {
	Output          string `json:"output"`
	Error           string `json:"error,omitempty"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
	Success         bool   `json:"success"`
}

// handlePlaygroundExecute executes Go code in the playground.
// POST /api/v1/playground/execute
func (h *Handler) handlePlaygroundExecute(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req PlaygroundExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	// Validate code is not empty
	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	// Check code size limit (64KB)
	if len(req.Code) > 65536 {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code too large (max 64KB)"))
		return
	}

	// Rate limit by IP
	userKey := getClientIP(r)
	if !h.checkSubmissionRateLimit(userKey) {
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded (10 requests per minute)"))
		return
	}

	// Create execution context with timeout (30s to allow for compilation)
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Execute code using the executor service
	// For playground, we use a single "output" test case
	result, err := h.services.Executor.ExecuteCode(ctx, &service.ExecuteRequest{
		Code:     req.Code,
		Language: "go",
		Timeout:  25 * time.Second, // Give time for compilation
		TestCases: []service.TestCase{
			{
				Name:     "Playground Output",
				Input:    "",
				Expected: "", // We don't care about expected, just want output
			},
		},
	})

	executionTimeMs := time.Since(startTime).Milliseconds()

	if err != nil {
		logger.LogError(h.logger, r.Context(), err, "playground execution failed")
		h.writeSuccessResponse(w, r, PlaygroundExecuteResult{
			Output:          "",
			Error:           err.Error(),
			ExecutionTimeMs: executionTimeMs,
			Success:         false,
		}, "execution failed")
		return
	}

	// Build output from results
	output := ""
	if len(result.Results) > 0 {
		// For playground, we just want the actual output
		output = result.Results[0].Actual
	}

	// The docker_executor returns an Error in the ExecuteResult only if there was an error in running the test (like timeout, etc.)
	// The TestResult's Error is set for test failures (including output mismatch) but we don't care about that for playground success.
	execErr := ""
	if result.Error != nil {
		execErr = result.Error.Error()
	}

	// For playground, success means the code compiled and ran without execution error (compilation/panic, timeout, etc.)
	// We don't care about the test passing or failing (since we have no expected output).
	// Check if any test result has an error that indicates a real execution problem (not just output mismatch)
	hasExecutionError := false
	if len(result.Results) > 0 {
		for _, testResult := range result.Results {
			if testResult.Error != "" {
				// Ignore "Output does not match expected result" as it's not a real execution error
				if !strings.Contains(testResult.Error, "Output does not match expected result") {
					hasExecutionError = true
					break
				}
			}
		}
	}

	// If the only error is output mismatch, clear the error field for display purposes
	if execErr == "Output does not match expected result" {
		execErr = ""
	}

	success := result.Error == nil && !hasExecutionError && execErr == ""

	h.writeSuccessResponse(w, r, PlaygroundExecuteResult{
		Output:          output,
		Error:           execErr,
		ExecutionTimeMs: executionTimeMs,
		Success:         success,
	}, "code executed successfully")
}
