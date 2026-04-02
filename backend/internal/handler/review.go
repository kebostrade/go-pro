// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP request handlers for the API.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/tools"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
)

// ReviewHandler handles code review requests.
type ReviewHandler struct {
	codeAnalysisTool *tools.CodeAnalysisTool
	repo             repository.ReviewRepository
	logger           logger.Logger
}

// NewReviewHandler creates a new review handler.
func NewReviewHandler(codeAnalysisTool *tools.CodeAnalysisTool, repo repository.ReviewRepository, logger logger.Logger) *ReviewHandler {
	return &ReviewHandler{
		codeAnalysisTool: codeAnalysisTool,
		repo:             repo,
		logger:           logger,
	}
}

// SubmitReviewRequest represents a request to submit code for review.
type SubmitReviewRequest struct {
	UserID     string `json:"user_id"`
	TopicID    string `json:"topic_id"`
	ExerciseID string `json:"exercise_id"`
	Code       string `json:"code"`
}

// ReviewResponse represents a code review response.
type ReviewResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TopicID     string    `json:"topic_id"`
	ExerciseID  string    `json:"exercise_id"`
	Code        string    `json:"code"`
	Feedback    string    `json:"feedback"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// RegisterRoutes registers the review routes.
func (h *ReviewHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/review/submit", h.handleSubmitReview)
	mux.HandleFunc("GET /api/review/history", h.handleGetReviewHistory)
}

// handleSubmitReview handles POST /api/review/submit
func (h *ReviewHandler) handleSubmitReview(w http.ResponseWriter, r *http.Request) {
	var req SubmitReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	if req.UserID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user_id is required"))
		return
	}

	ctx := r.Context()

	feedback, err := h.analyzeCode(ctx, req.Code)
	if err != nil {
		h.logger.Error(ctx, "Code analysis failed", "error", err)
		h.writeErrorResponse(w, r, apierrors.NewInternalError("code analysis failed", err))
		return
	}

	// Create review record
	review := &domain.Review{
		UserID:     req.UserID,
		TopicID:    req.TopicID,
		ExerciseID: req.ExerciseID,
		Code:       req.Code,
		Feedback:   feedback,
	}

	// Store in repository
	if h.repo != nil {
		if err := h.repo.Create(ctx, review); err != nil {
			h.logger.Error(ctx, "Failed to save review", "error", err)
			h.writeErrorResponse(w, r, apierrors.NewInternalError("failed to save review", err))
			return
		}
	}

	response := ReviewResponse{
		ID:          review.ID,
		UserID:      review.UserID,
		TopicID:     review.TopicID,
		ExerciseID:  review.ExerciseID,
		Code:        review.Code,
		Feedback:    review.Feedback,
		SubmittedAt: review.SubmittedAt,
	}

	h.writeSuccessResponse(w, r, response, "code review submitted successfully")
}

// handleGetReviewHistory handles GET /api/review/history
func (h *ReviewHandler) handleGetReviewHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user_id query parameter is required"))
		return
	}

	// Fetch from repository if available
	if h.repo != nil {
		reviews, err := h.repo.GetByUserID(r.Context(), userID)
		if err != nil {
			h.logger.Error(r.Context(), "Failed to fetch review history", "error", err)
			h.writeErrorResponse(w, r, apierrors.NewInternalError("failed to fetch review history", err))
			return
		}

		response := make([]ReviewResponse, 0, len(reviews))
		for _, review := range reviews {
			response = append(response, ReviewResponse{
				ID:          review.ID,
				UserID:      review.UserID,
				TopicID:     review.TopicID,
				ExerciseID:  review.ExerciseID,
				Code:        review.Code,
				Feedback:    review.Feedback,
				SubmittedAt: review.SubmittedAt,
			})
		}

		h.writeSuccessResponse(w, r, response, "review history retrieved successfully")
		return
	}

	// Fallback if no repository
	response := []ReviewResponse{}
	h.writeSuccessResponse(w, r, response, "review history retrieved successfully")
}

// analyzeCode analyzes code using the CodeAnalysisTool.
func (h *ReviewHandler) analyzeCode(ctx context.Context, code string) (string, error) {
	input := types.NewToolInput(map[string]interface{}{
		"code":     code,
		"language": "go",
	})
	output, err := h.codeAnalysisTool.Execute(ctx, input)
	if err != nil {
		return "", err
	}
	if output.Error != nil {
		return "", fmt.Errorf("analysis error: %s", output.Error.Message)
	}
	return formatConversationalFeedback(output.Result), nil
}

// formatConversationalFeedback converts analysis result to encouraging feedback.
func formatConversationalFeedback(result interface{}) string {
	builder := &feedbackBuilder{}

	builder.WriteLine("Great job submitting your code!")

	var analysis struct {
		Issues            []interface{} `json:"issues"`
		SecurityIssues    []interface{} `json:"security_issues"`
		PerformanceIssues []interface{} `json:"performance_issues"`
		BestPractices     []interface{} `json:"best_practices"`
		OverallScore      int           `json:"overall_score"`
	}

	if result != nil {
		resultStr, ok := result.(string)
		if ok {
			json.Unmarshal([]byte(resultStr), &analysis)
		}
	}

	if len(analysis.Issues) == 0 && len(analysis.SecurityIssues) == 0 &&
		len(analysis.PerformanceIssues) == 0 && len(analysis.BestPractices) == 0 {
		builder.WriteLine("Your code looks clean and follows good practices.")
	} else {
		if len(analysis.SecurityIssues) > 0 {
			builder.WriteLine("\n## Security Considerations")
			for _, issue := range analysis.SecurityIssues {
				if issueMap, ok := issue.(map[string]interface{}); ok {
					if desc, ok := issueMap["description"].(string); ok {
						builder.WriteLine(fmt.Sprintf("- %s", desc))
					}
				}
			}
		}

		if len(analysis.PerformanceIssues) > 0 {
			builder.WriteLine("\n## Performance Tips")
			for _, issue := range analysis.PerformanceIssues {
				if issueMap, ok := issue.(map[string]interface{}); ok {
					if desc, ok := issueMap["description"].(string); ok {
						builder.WriteLine(fmt.Sprintf("- %s", desc))
					}
				}
			}
		}

		if len(analysis.BestPractices) > 0 {
			builder.WriteLine("\n## Best Practices")
			for _, bp := range analysis.BestPractices {
				if bpMap, ok := bp.(map[string]interface{}); ok {
					if practice, ok := bpMap["practice"].(string); ok {
						if recommendation, ok := bpMap["recommendation"].(string); ok {
							builder.WriteLine(fmt.Sprintf("- %s: %s", practice, recommendation))
						}
					}
				}
			}
		}
	}

	builder.WriteLine(fmt.Sprintf("\n## Overall Score: %d/100", analysis.OverallScore))

	if analysis.OverallScore >= 90 {
		builder.WriteLine("Excellent work! Your code is production-ready.")
	} else if analysis.OverallScore >= 70 {
		builder.WriteLine("Good job! A few minor improvements could make it even better.")
	} else if analysis.OverallScore >= 50 {
		builder.WriteLine("Nice effort! Consider addressing the issues above to improve your code quality.")
	} else {
		builder.WriteLine("Keep practicing! Review the feedback above and try again.")
	}

	return builder.String()
}

type feedbackBuilder struct {
	content string
}

func (fb *feedbackBuilder) WriteLine(s string) {
	fb.content += s + "\n"
}

func (fb *feedbackBuilder) String() string {
	return fb.content
}

func generateReviewID() string {
	return "review_" + time.Now().Format("20060102150405")
}

func (h *ReviewHandler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success":   true,
		"data":      data,
		"message":   message,
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "Failed to encode response", "error", err)
	}
}

func (h *ReviewHandler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierrors.APIError
	var statusCode int

	if e, ok := err.(*apierrors.APIError); ok {
		apiErr = e
		statusCode = e.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
		apiErr = apierrors.NewInternalError("internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": apiErr.Message,
			"type":    apiErr.Type,
		},
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "Failed to encode error response", "error", err)
	}
}
