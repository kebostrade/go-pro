// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
		Code:      req.Code,
		Language:  "go",
		Timeout:   25 * time.Second, // Give time for compilation
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
	var execErr string
	if result.Error != nil {
		execErr = fmt.Sprintf("Error: %v", result.Error)
	} else if len(result.Results) > 0 {
		// For playground, we just want the actual output
		output = result.Results[0].Actual
		if result.Results[0].Error != "" {
			execErr = result.Results[0].Error
		}
	}

	h.writeSuccessResponse(w, r, PlaygroundExecuteResult{
		Output:          output,
		Error:           execErr,
		ExecutionTimeMs: executionTimeMs,
		Success:         result.Error == nil && execErr == "",
	}, "code executed successfully")
}
