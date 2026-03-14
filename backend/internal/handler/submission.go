// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP handlers for the GO-PRO Learning Platform.
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/repository"
)

// SubmissionHandler handles submission-related HTTP requests.
type SubmissionHandler struct {
	repos *repository.Repositories
}

// NewSubmissionHandler creates a new submission handler.
func NewSubmissionHandler(repos *repository.Repositories) *SubmissionHandler {
	return &SubmissionHandler{
		repos: repos,
	}
}

// CreateSubmission handles POST /api/submissions
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	userID := c.GetString("user_id")

	var req domain.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Invalid request payload",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Get assessment to check if it exists
	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), req.AssessmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"assessment_id": req.AssessmentID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Check previous attempts (for quiz with max_attempts)
	if assessment.Type == domain.AssessmentTypeQuiz {
		var quizConfig domain.QuizConfig
		configBytes, _ := json.Marshal(assessment.Config)
		json.Unmarshal(configBytes, &quizConfig)

		// Get user's previous submissions for this assessment
		prevSubmission, _ := h.repos.Submission.GetByUserAndAssessment(c.Request.Context(), userID, req.AssessmentID)

		// Check if already submitted (one submission per assessment)
		if prevSubmission != nil {
			c.JSON(http.StatusBadRequest, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Type:    "validation_error",
					Message: "Already submitted this assessment",
				},
				Timestamp: time.Now(),
			})
			return
		}

		// Auto-grade quiz
		score := h.autoGradeQuiz(req.Content, quizConfig)

		// Create submission with score
		submission := &domain.Submission{
			ID:           uuid.New().String(),
			AssessmentID: req.AssessmentID,
			UserID:       userID,
			Content:      req.Content,
			Score:        &score,
			Status:       domain.SubmissionStatusSubmitted,
			SubmittedAt:  time.Now(),
		}

		if err := h.repos.Submission.Create(c.Request.Context(), submission); err != nil {
			c.JSON(http.StatusInternalServerError, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Type:    "database_error",
					Message: "Failed to create submission",
					Details: map[string]string{"error": err.Error()},
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusCreated, domain.APIResponse{
			Success:   true,
			Data:      submission,
			Message:   "Quiz submitted successfully",
			Timestamp: time.Now(),
		})
		return
	}

	// For coding exercises and projects, store without auto-grading
	submission := &domain.Submission{
		ID:           uuid.New().String(),
		AssessmentID: req.AssessmentID,
		UserID:       userID,
		Content:      req.Content,
		Status:       domain.SubmissionStatusSubmitted,
		SubmittedAt:  time.Now(),
	}

	if err := h.repos.Submission.Create(c.Request.Context(), submission); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to create submission",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.APIResponse{
		Success:   true,
		Data:      submission,
		Message:   "Submission created successfully",
		Timestamp: time.Now(),
	})
}

// autoGradeQuiz automatically grades quiz submissions
func (h *SubmissionHandler) autoGradeQuiz(content map[string]interface{}, quizConfig domain.QuizConfig) int {
	answers, ok := content["answers"].(map[string]interface{})
	if !ok {
		return 0
	}

	totalScore := 0
	maxScore := 0

	for _, question := range quizConfig.Questions {
		maxScore += question.Points

		userAnswer, exists := answers[question.ID]
		if !exists {
			continue
		}

		// Auto-grade based on question type
		switch question.QuestionType {
		case domain.QuestionTypeMultipleChoice, domain.QuestionTypeTrueFalse, domain.QuestionTypeCodeCompletion:
			// Direct comparison
			if userAnswer == question.CorrectAnswer {
				totalScore += question.Points
			}
		case domain.QuestionTypeShortAnswer:
			// Short answer requires manual grading
			// Don't add points, instructor will grade manually
		}
	}

	// Calculate percentage
	if maxScore > 0 {
		return (totalScore * 100) / maxScore
	}
	return 0
}

// ListSubmissions handles GET /api/submissions
func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	userID := c.GetString("user_id")
	assessmentID := c.Query("assessment_id")
	status := c.Query("status")

	page, pageSize := getPaginationParams(c)
	pagination := &domain.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Build filters
	filters := map[string]interface{}{}
	if userID != "" {
		filters["user_id"] = userID
	}
	if assessmentID != "" {
		filters["assessment_id"] = assessmentID
	}
	if status != "" {
		filters["status"] = domain.SubmissionStatus(status)
	}

	submissions, total, err := h.repos.Submission.GetAll(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to retrieve submissions",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data: domain.ListResponse{
			Items: submissions,
			Pagination: &domain.PaginationResponse{
				Page:       page,
				PageSize:   pageSize,
				TotalItems: total,
				TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
				HasNext:    page*pageSize < int(total),
				HasPrev:    page > 1,
			},
		},
		Timestamp: time.Now(),
	})
}

// GetSubmission handles GET /api/submissions/:id
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	id := c.Param("id")

	submission, err := h.repos.Submission.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Submission not found",
				Details: map[string]string{"id": id},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      submission,
		Timestamp: time.Now(),
	})
}

// GradeSubmission handles POST /api/cms/submissions/:id/grade
func (h *SubmissionHandler) GradeSubmission(c *gin.Context) {
	id := c.Param("id")
	graderID := c.GetString("user_id")

	var req domain.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Invalid request payload",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Get submission
	submission, err := h.repos.Submission.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Submission not found",
				Details: map[string]string{"id": id},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Update submission with grade
	now := time.Now()
	submission.Score = req.Score
	if req.Feedback != nil {
		submission.Feedback = *req.Feedback
	}
	submission.GradedBy = &graderID
	submission.GradedAt = &now

	if req.ReleaseImmediately {
		submission.Status = domain.SubmissionStatusReturned
	} else {
		submission.Status = domain.SubmissionStatusGraded
	}

	if err := h.repos.Submission.Update(c.Request.Context(), submission); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to grade submission",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// TODO: Send email notification if release_immediately is true

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      submission,
		Message:   "Submission graded successfully",
		Timestamp: time.Now(),
	})
}

// AddSubmissionComment handles POST /api/cms/submissions/:id/comments
func (h *SubmissionHandler) AddSubmissionComment(c *gin.Context) {
	id := c.Param("id")
	authorID := c.GetString("user_id")

	var req domain.CreateSubmissionCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Invalid request payload",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Verify submission exists
	_, err := h.repos.Submission.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Submission not found",
				Details: map[string]string{"id": id},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Create comment
	comment := &domain.SubmissionComment{
		ID:           uuid.New().String(),
		SubmissionID: id,
		AuthorID:     authorID,
		CommentText:  req.CommentText,
		LineNumber:   req.LineNumber,
		CreatedAt:    time.Now(),
	}

	if err := h.repos.SubmissionComment.Create(c.Request.Context(), comment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to create comment",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.APIResponse{
		Success:   true,
		Data:      comment,
		Message:   "Comment added successfully",
		Timestamp: time.Now(),
	})
}

// GetSubmissionComments handles GET /api/cms/submissions/:id/comments
func (h *SubmissionHandler) GetSubmissionComments(c *gin.Context) {
	id := c.Param("id")

	comments, err := h.repos.SubmissionComment.GetBySubmissionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to retrieve comments",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      comments,
		Timestamp: time.Now(),
	})
}
