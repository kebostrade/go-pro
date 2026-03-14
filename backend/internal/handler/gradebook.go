// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP handlers for the GO-PRO Learning Platform.
package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/repository"
)

// GradebookHandler handles gradebook-related HTTP requests.
type GradebookHandler struct {
	repos *repository.Repositories
}

// NewGradebookHandler creates a new gradebook handler.
func NewGradebookHandler(repos *repository.Repositories) *GradebookHandler {
	return &GradebookHandler{
		repos: repos,
	}
}

// ExportGradebook handles GET /api/cms/gradebook
func (h *GradebookHandler) ExportGradebook(c *gin.Context) {
	assessmentID := c.Query("assessment_id")
	status := c.Query("status")

	// Build filters
	filters := map[string]interface{}{}
	if assessmentID != "" {
		filters["assessment_id"] = assessmentID
	}
	if status != "" {
		filters["status"] = domain.SubmissionStatus(status)
	}

	// Get all submissions with filters
	pagination := &domain.PaginationRequest{
		Page:     1,
		PageSize: 10000, // Large page size for export
	}

	submissions, _, err := h.repos.Submission.GetAll(c.Request.Context(), filters, pagination)
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

	// Build gradebook entries
	var gradebookEntries []domain.GradebookEntry
	for _, submission := range submissions {
		// Get user details
		user, err := h.repos.User.GetByID(c.Request.Context(), submission.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		// Get assessment details
		assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), submission.AssessmentID)
		if err != nil {
			continue // Skip if assessment not found
		}

		entry := domain.GradebookEntry{
			StudentID:       user.ID,
			StudentName:     user.DisplayName,
			StudentEmail:    user.Email,
			AssessmentID:    assessment.ID,
			AssessmentTitle: assessment.Title,
			Score:           submission.Score,
			SubmittedAt:     &submission.SubmittedAt,
			GradedAt:        submission.GradedAt,
			Status:          submission.Status,
		}

		gradebookEntries = append(gradebookEntries, entry)
	}

	// Generate CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="gradebook.csv"`)

	// Write UTF-8 BOM for Excel compatibility
	c.Writer.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	headers := []string{
		"Student ID",
		"Student Name",
		"Student Email",
		"Assessment ID",
		"Assessment Title",
		"Score",
		"Submitted At",
		"Graded At",
		"Status",
	}
	writer.Write(headers)

	// Write data rows
	for _, entry := range gradebookEntries {
		var score string
		if entry.Score != nil {
			score = strconv.Itoa(*entry.Score)
		} else {
			score = ""
		}

		var submittedAt, gradedAt string
		if entry.SubmittedAt != nil {
			submittedAt = entry.SubmittedAt.Format(time.RFC3339)
		}
		if entry.GradedAt != nil {
			gradedAt = entry.GradedAt.Format(time.RFC3339)
		}

		record := []string{
			entry.StudentID,
			entry.StudentName,
			entry.StudentEmail,
			entry.AssessmentID,
			entry.AssessmentTitle,
			score,
			submittedAt,
			gradedAt,
			string(entry.Status),
		}
		writer.Write(record)
	}
}

// GetGradebookData handles GET /api/cms/gradebook/data (returns JSON for UI)
func (h *GradebookHandler) GetGradebookData(c *gin.Context) {
	assessmentID := c.Query("assessment_id")
	status := c.Query("status")

	page, pageSize := getPaginationParams(c)
	pagination := &domain.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Build filters
	filters := map[string]interface{}{}
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
				Message: "Failed to retrieve gradebook data",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Build gradebook entries
	var gradebookEntries []domain.GradebookEntry
	for _, submission := range submissions {
		// Get user details
		user, err := h.repos.User.GetByID(c.Request.Context(), submission.UserID)
		if err != nil {
			continue
		}

		// Get assessment details
		assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), submission.AssessmentID)
		if err != nil {
			continue
		}

		entry := domain.GradebookEntry{
			StudentID:       user.ID,
			StudentName:     user.DisplayName,
			StudentEmail:    user.Email,
			AssessmentID:    assessment.ID,
			AssessmentTitle: assessment.Title,
			Score:           submission.Score,
			SubmittedAt:     &submission.SubmittedAt,
			GradedAt:        submission.GradedAt,
			Status:          submission.Status,
		}

		gradebookEntries = append(gradebookEntries, entry)
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data: domain.ListResponse{
			Items: gradebookEntries,
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

// GetStudentGrades handles GET /api/grades (returns grades for current student)
func (h *GradebookHandler) GetStudentGrades(c *gin.Context) {
	userID := c.GetString("user_id")

	page, pageSize := getPaginationParams(c)
	pagination := &domain.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Get user's submissions
	submissions, total, err := h.repos.Submission.GetByUserID(c.Request.Context(), userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to retrieve grades",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Build grade entries with assessment details
	var gradeEntries []map[string]interface{}
	for _, submission := range submissions {
		assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), submission.AssessmentID)
		if err != nil {
			continue
		}

		entry := map[string]interface{}{
			"submission_id":    submission.ID,
			"assessment_id":    assessment.ID,
			"assessment_title": assessment.Title,
			"assessment_type":  assessment.Type,
			"score":            submission.Score,
			"passing_score":    assessment.PassingScore,
			"passed":           false,
			"status":           submission.Status,
			"submitted_at":     submission.SubmittedAt,
			"graded_at":        submission.GradedAt,
			"feedback":         submission.Feedback,
		}

		// Check if passed
		if submission.Score != nil && *submission.Score >= assessment.PassingScore {
			entry["passed"] = true
		}

		gradeEntries = append(gradeEntries, entry)
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data: domain.ListResponse{
			Items: gradeEntries,
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

// GetGradeStatistics handles GET /api/cms/gradebook/statistics
func (h *GradebookHandler) GetGradeStatistics(c *gin.Context) {
	assessmentID := c.Query("assessment_id")

	pagination := &domain.PaginationRequest{
		Page:     1,
		PageSize: 10000,
	}

	// Build filters
	filters := map[string]interface{}{}
	if assessmentID != "" {
		filters["assessment_id"] = assessmentID
	}

	submissions, _, err := h.repos.Submission.GetAll(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Type:    "database_error",
				Message: "Failed to retrieve statistics",
				Details: map[string]string{"error": err.Error()},
			},
			Timestamp: time.Now(),
		})
		return
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_submissions": len(submissions),
		"graded_count":      0,
		"pending_count":     0,
		"returned_count":    0,
		"average_score":     0.0,
		"pass_rate":         0.0,
	}

	var totalScore int
	var passingCount int
	var scoredCount int

	for _, submission := range submissions {
		switch submission.Status {
		case domain.SubmissionStatusGraded:
			stats["graded_count"] = stats["graded_count"].(int) + 1
		case domain.SubmissionStatusSubmitted:
			stats["pending_count"] = stats["pending_count"].(int) + 1
		case domain.SubmissionStatusReturned:
			stats["returned_count"] = stats["returned_count"].(int) + 1
		}

		if submission.Score != nil {
			totalScore += *submission.Score
			scoredCount++

			// Get assessment to check passing score
			assessment, err := h.repos.Assessment.GetByID(c.Request.Context(), submission.AssessmentID)
			if err == nil && *submission.Score >= assessment.PassingScore {
				passingCount++
			}
		}
	}

	// Calculate averages
	if scoredCount > 0 {
		stats["average_score"] = float64(totalScore) / float64(scoredCount)
		stats["pass_rate"] = float64(passingCount) / float64(scoredCount) * 100
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success:   true,
		Data:      stats,
		Timestamp: time.Now(),
	})
}

// BulkGrade handles POST /api/cms/gradebook/bulk-grade
func (h *GradebookHandler) BulkGrade(c *gin.Context) {
	var req struct {
		SubmissionIDs []string `json:"submission_ids" validate:"required"`
		Score         int      `json:"score" validate:"required,min=0,max=100"`
		Feedback      string   `json:"feedback,omitempty"`
		Release       bool     `json:"release_immediately"`
	}

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

	graderID := c.GetString("user_id")
	gradedCount := 0
	now := time.Now()

	for _, submissionID := range req.SubmissionIDs {
		submission, err := h.repos.Submission.GetByID(c.Request.Context(), submissionID)
		if err != nil {
			continue // Skip if submission not found
		}

		submission.Score = &req.Score
		if req.Feedback != "" {
			submission.Feedback = req.Feedback
		}
		submission.GradedBy = &graderID
		submission.GradedAt = &now

		if req.Release {
			submission.Status = domain.SubmissionStatusReturned
		} else {
			submission.Status = domain.SubmissionStatusGraded
		}

		if err := h.repos.Submission.Update(c.Request.Context(), submission); err != nil {
			continue // Skip if update fails
		}

		gradedCount++
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"graded_count": gradedCount,
			"total_count":  len(req.SubmissionIDs),
		},
		Message:   fmt.Sprintf("Graded %d out of %d submissions", gradedCount, len(req.SubmissionIDs)),
		Timestamp: time.Now(),
	})
}
