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

// AssessmentHandler handles assessment-related HTTP requests.
type AssessmentHandler struct {
	repos *repository.Repositories
}

// NewAssessmentHandler creates a new assessment handler.
func NewAssessmentHandler(repos *repository.Repositories) *AssessmentHandler {
	return &AssessmentHandler{
		repos: repos,
	}
}

// CreateAssessment handles POST /api/cms/assessments
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	var req domain.CreateAssessmentRequest
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

	// Validate assessment type
	if !req.Type.IsValid() {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Invalid assessment type",
				Details: map[string]string{"type": string(req.Type)},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Create assessment
	assessment := &domain.Assessment{
		ID:               uuid.New().String(),
		LessonID:         req.LessonID,
		Type:             req.Type,
		Title:            req.Title,
		Description:      req.Description,
		Config:           req.Config,
		PassingScore:     req.PassingScore,
		TimeLimitMinutes: req.TimeLimitMinutes,
		OrderIndex:       req.OrderIndex,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.repos.Assessment.Create(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to create assessment",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.APIResponse{
		Success:   true,
		Data:      assessment,
		Message:   "Assessment created successfully",
		Timestamp: time.Now(),
	})
}

// ListAssessments handles GET /api/cms/assessments
func (h *AssessmentHandler) ListAssessments(c *gin.Context) {
	assessmentType := c.Query("type")
	page, pageSize := getPaginationParams(c)

	pagination := &domain.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	var assessments []*domain.Assessment
	var total int64
	var err error

	// Filter by type if provided
	if assessmentType != "" {
		assessmentTypeDomain := domain.AssessmentType(assessmentType)
		if !assessmentTypeDomain.IsValid() {
			c.JSON(http.StatusBadRequest, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Type:    "validation_error",
					Message: "Invalid assessment type",
					Details: map[string]string{"type": assessmentType},
				},
				Timestamp: time.Now(),
			})
			return
		}
		assessments, total, err = h.repos.Assessment.GetByType(c.Request.Context(), assessmentTypeDomain, pagination)
	} else {
		assessments, total, err = h.repos.Assessment.GetAll(c.Request.Context(), pagination)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to retrieve assessments",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data: domain.ListResponse{
			Items: assessments,
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

// GetAssessment handles GET /api/cms/assessments/:id
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
	id := c.Param("id")

	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"id": id},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      assessment,
		Timestamp: time.Now(),
	})
}

// UpdateAssessment handles PUT /api/cms/assessments/:id
func (h *AssessmentHandler) UpdateAssessment(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateAssessmentRequest
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

	// Get existing assessment
	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"id": id},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Update fields
	if req.Title != nil {
		assessment.Title = *req.Title
	}
	if req.Description != nil {
		assessment.Description = *req.Description
	}
	if req.Config != nil {
		assessment.Config = *req.Config
	}
	if req.PassingScore != nil {
		assessment.PassingScore = *req.PassingScore
	}
	if req.TimeLimitMinutes != nil {
		assessment.TimeLimitMinutes = req.TimeLimitMinutes
	}
	if req.OrderIndex != nil {
		assessment.OrderIndex = *req.OrderIndex
	}
	assessment.UpdatedAt = time.Now()

	if err := h.repos.Assessment.Update(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to update assessment",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      assessment,
		Message:   "Assessment updated successfully",
		Timestamp: time.Now(),
	})
}

// DeleteAssessment handles DELETE /api/cms/assessments/:id
func (h *AssessmentHandler) DeleteAssessment(c *gin.Context) {
	id := c.Param("id")

	if err := h.repos.Assessment.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to delete assessment",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Message:   "Assessment deleted successfully",
		Timestamp: time.Now(),
	})
}

// AddQuestion handles POST /api/cms/assessments/:id/questions
func (h *AssessmentHandler) AddQuestion(c *gin.Context) {
	assessmentID := c.Param("id")

	var question domain.Question
	if err := c.ShouldBindJSON(&question); err != nil {
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

	// Validate question type
	if !question.QuestionType.IsValid() {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Invalid question type",
				Details: map[string]string{"question_type": string(question.QuestionType)},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Set ID and assessment ID
	question.ID = uuid.New().String()
	question.AssessmentID = assessmentID
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	// Store question in assessment config
	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), assessmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"id": assessmentID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	if assessment.Type != domain.AssessmentTypeQuiz {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "validation_error",
				Message: "Questions can only be added to quiz assessments",
				Details: map[string]string{"assessment_type": string(assessment.Type)},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Get quiz config and add question
	var quizConfig domain.QuizConfig
	configBytes, _ := json.Marshal(assessment.Config)
	if err := json.Unmarshal(configBytes, &quizConfig); err != nil {
		quizConfig = domain.QuizConfig{
			Questions: []domain.Question{},
			Settings: domain.QuizSettings{
				MaxAttempts:      1,
				ShuffleQuestions: false,
				ShowExplanations: true,
			},
		}
	}

	quizConfig.Questions = append(quizConfig.Questions, question)
	updatedConfig, _ := json.Marshal(quizConfig)
	json.Unmarshal(updatedConfig, &assessment.Config)

	assessment.UpdatedAt = time.Now()
	if err := h.repos.Assessment.Update(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to add question",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.APIResponse{
		Success:   true,
		Data:      question,
		Message:   "Question added successfully",
		Timestamp: time.Now(),
	})
}

// UpdateQuestion handles PUT /api/cms/assessments/:id/questions/:qid
func (h *AssessmentHandler) UpdateQuestion(c *gin.Context) {
	assessmentID := c.Param("id")
	questionID := c.Param("qid")

	var updates domain.Question
	if err := c.ShouldBindJSON(&updates); err != nil {
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

	// Get assessment
	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), assessmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"id": assessmentID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Get quiz config and update question
	var quizConfig domain.QuizConfig
	configBytes, _ := json.Marshal(assessment.Config)
	if err := json.Unmarshal(configBytes, &quizConfig); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to parse quiz configuration",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Find and update question
	questionFound := false
	for i, q := range quizConfig.Questions {
		if q.ID == questionID {
			quizConfig.Questions[i].QuestionText = updates.QuestionText
			quizConfig.Questions[i].Options = updates.Options
			quizConfig.Questions[i].CorrectAnswer = updates.CorrectAnswer
			quizConfig.Questions[i].Explanation = updates.Explanation
			quizConfig.Questions[i].Points = updates.Points
			quizConfig.Questions[i].Tags = updates.Tags
			quizConfig.Questions[i].Hints = updates.Hints
			quizConfig.Questions[i].HintThreshold = updates.HintThreshold
			quizConfig.Questions[i].UpdatedAt = time.Now()
			questionFound = true
			break
		}
	}

	if !questionFound {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Question not found",
				Details: map[string]string{"question_id": questionID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Save updated config
	updatedConfig, _ := json.Marshal(quizConfig)
	json.Unmarshal(updatedConfig, &assessment.Config)
	assessment.UpdatedAt = time.Now()

	if err := h.repos.Assessment.Update(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to update question",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Message:   "Question updated successfully",
		Timestamp: time.Now(),
	})
}

// DeleteQuestion handles DELETE /api/cms/assessments/:id/questions/:qid
func (h *AssessmentHandler) DeleteQuestion(c *gin.Context) {
	assessmentID := c.Param("id")
	questionID := c.Param("qid")

	// Get assessment
	assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), assessmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Assessment not found",
				Details: map[string]string{"id": assessmentID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Get quiz config and delete question
	var quizConfig domain.QuizConfig
	configBytes, _ := json.Marshal(assessment.Config)
	if err := json.Unmarshal(configBytes, &quizConfig); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to parse quiz configuration",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Find and delete question
	questionFound := false
	updatedQuestions := []domain.Question{}
	for _, q := range quizConfig.Questions {
		if q.ID != questionID {
			updatedQuestions = append(updatedQuestions, q)
		} else {
			questionFound = true
		}
	}

	if !questionFound {
		c.JSON(http.StatusNotFound, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "not_found",
				Message: "Question not found",
				Details: map[string]string{"question_id": questionID},
			},
			Timestamp: time.Now(),
		})
		return
	}

	quizConfig.Questions = updatedQuestions

	// Save updated config
	updatedConfig, _ := json.Marshal(quizConfig)
	json.Unmarshal(updatedConfig, &assessment.Config)
	assessment.UpdatedAt = time.Now()

	if err := h.repos.Assessment.Update(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to delete question",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Message:   "Question deleted successfully",
		Timestamp: time.Now(),
	})
}
