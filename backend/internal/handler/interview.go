// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP request handlers for interview API.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"

	apierrors "go-pro-backend/internal/errors"
)

// InterviewType represents the type of interview.
type InterviewType string

const (
	InterviewTypeCoding       InterviewType = "coding"
	InterviewTypeBehavioral   InterviewType = "behavioral"
	InterviewTypeSystemDesign InterviewType = "system_design"
)

// Difficulty represents interview difficulty level.
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

// InterviewSession represents an interview session.
//nolint:typecheck - using int64 for timestamps to work with repository
type InterviewSession struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	Type         InterviewType   `json:"type"`
	Difficulty   Difficulty     `json:"difficulty"`
	Questions    []Question      `json:"questions"`
	CurrentIndex int             `json:"current_index"`
	Answers      []Answer        `json:"answers"`
	Status       string         `json:"status"`
	Score        *int            `json:"score,omitempty"`
	CreatedAt    int64           `json:"created_at"`
	CompletedAt  *int64          `json:"completed_at,omitempty"`
}

// Question represents an interview question.
type Question struct {
	ID             string   `json:"id"`
	Content        string   `json:"content"`
	Type           string   `json:"type"`
	Difficulty     string   `json:"difficulty"`
	ExpectedPoints []string `json:"expected_points,omitempty"`
	TimeLimit      int      `json:"time_limit"`
}

// Answer represents a user's answer to a question.
type Answer struct {
	QuestionID string    `json:"question_id"`
	Content    string    `json:"content"`
	Score      *int      `json:"score,omitempty"`
	Feedback   string    `json:"feedback,omitempty"`
	CreatedAt  int64   `json:"created_at"`
}

// InterviewFeedback represents AI feedback on an interview.
type InterviewFeedback struct {
	SessionID        string             `json:"session_id"`
	OverallScore     int                `json:"overall_score"`
	Strengths        []string           `json:"strengths"`
	Improvements     []string           `json:"improvements"`
	DetailedFeedback []QuestionFeedback  `json:"detailed_feedback"`
}

// QuestionFeedback represents detailed feedback for a specific question.
type QuestionFeedback struct {
	QuestionID string   `json:"question_id"`
	Score      int      `json:"score"`
	Feedback   string   `json:"feedback"`
	Strengths  []string `json:"strengths"`
	Missed     []string `json:"missed"`
}

// InterviewHandler handles interview-related HTTP requests.
type InterviewHandler struct {
	interviewRepo repository.InterviewRepository
	logger        logger.Logger
	validator     validator.Validator
}

// NewInterviewHandler creates a new interview handler.
func NewInterviewHandler(interviewRepo repository.InterviewRepository, logger logger.Logger, validator validator.Validator) *InterviewHandler {
	return &InterviewHandler{
		interviewRepo: interviewRepo,
		logger:        logger,
		validator:     validator,
	}
}

// StartInterviewRequest represents request to start a new interview.
type StartInterviewRequest struct {
	Type       InterviewType `json:"type" validate:"required,oneof=coding behavioral system_design"`
	Difficulty Difficulty   `json:"difficulty" validate:"required,oneof=beginner intermediate advanced"`
}

// AnswerQuestionRequest represents request to answer a question.
type AnswerQuestionRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Answer    string `json:"answer" validate:"required"`
}

// handleStartInterview starts a new interview session.
func (h *InterviewHandler) handleStartInterview(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	var req StartInterviewRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	questions := h.generateQuestions(req.Type, req.Difficulty)

	session := &InterviewSession{
		ID:           generateID(),
		UserID:       user.ID,
		Type:         req.Type,
		Difficulty:   req.Difficulty,
		Questions:    questions,
		CurrentIndex: 0,
		Answers:      []Answer{},
		Status:       "in_progress",
		CreatedAt:    time.Now().Unix(),
	}

	if err := h.interviewRepo.Create(r.Context(), session); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Return the session with questions for the user to start answering
	h.writeSuccessResponse(w, r, session, "interview session started")
}

// handleAnswerQuestion submits an answer and returns the next question.
func (h *InterviewHandler) handleAnswerQuestion(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	var req AnswerQuestionRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	session, err := h.interviewRepo.GetByID(r.Context(), req.SessionID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Type assertion from interface{}
	sessionData, ok := session.(*InterviewSession)
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewInternalError("invalid session type", nil))
		return
	}

	if sessionData.UserID != user.ID {
		h.writeErrorResponse(w, r, apierrors.NewForbiddenError("access denied to this session"))
		return
	}

	if sessionData.Status != "in_progress" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("interview session is not in progress"))
		return
	}

	if sessionData.CurrentIndex >= len(sessionData.Questions) {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("all questions have been answered"))
		return
	}

	currentQuestion := sessionData.Questions[sessionData.CurrentIndex]
	answer := Answer{
		QuestionID: currentQuestion.ID,
		Content:    req.Answer,
		CreatedAt:  time.Now().Unix(),
		Score:      h.calculateAnswerScore(currentQuestion, req.Answer),
		Feedback:   h.generateFeedback(currentQuestion, req.Answer),
	}

	sessionData.Answers = append(sessionData.Answers, answer)
	sessionData.CurrentIndex++

	if sessionData.CurrentIndex >= len(sessionData.Questions) {
		sessionData.Status = "completed"
		nowUnix := time.Now().Unix()
		sessionData.CompletedAt = &nowUnix
		sessionData.Score = h.calculateOverallScore(sessionData.Answers)
	}

	if err := h.interviewRepo.Update(r.Context(), sessionData); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"answer":    answer,
		"completed": sessionData.Status == "completed",
	}

	if sessionData.Status == "completed" {
		response["next_question"] = nil
		response["session"] = sessionData
	} else {
		nextQuestion := sessionData.Questions[sessionData.CurrentIndex]
		response["next_question"] = nextQuestion
	}

	h.writeSuccessResponse(w, r, response, "answer submitted")
}

// handleGetFeedback retrieves AI feedback for a completed interview.
func (h *InterviewHandler) handleGetFeedback(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	sessionID := r.PathValue("id")
	if sessionID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("session ID is required"))
		return
	}

	sessionData, err := h.interviewRepo.GetByID(r.Context(), sessionID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Type assertion from interface{}
	session, ok := sessionData.(*InterviewSession)
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewInternalError("invalid session type", nil))
		return
	}

	if session.UserID != user.ID {
		h.writeErrorResponse(w, r, apierrors.NewForbiddenError("access denied to this session"))
		return
	}

	if session.Status != "completed" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("interview session must be completed to get feedback"))
		return
	}

	feedback := h.generateInterviewFeedback(session)

	h.writeSuccessResponse(w, r, feedback, "feedback retrieved")
}

// handleListSessions lists user's past interview sessions.
func (h *InterviewHandler) handleListSessions(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	sessionsData, err := h.interviewRepo.GetByUserID(r.Context(), user.ID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, sessionsData, "sessions retrieved")
}

// handleGetSession retrieves a specific interview session.
func (h *InterviewHandler) handleGetSession(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	sessionID := r.PathValue("id")
	if sessionID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("session ID is required"))
		return
	}

	sessionData, err := h.interviewRepo.GetByID(r.Context(), sessionID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Type assertion from interface{}
	session, ok := sessionData.(*InterviewSession)
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewInternalError("invalid session type", nil))
		return
	}

	if session.UserID != user.ID {
		h.writeErrorResponse(w, r, apierrors.NewForbiddenError("access denied to this session"))
		return
	}

	h.writeSuccessResponse(w, r, session, "session retrieved")
}

// RegisterInterviewRoutes registers interview API routes.
func (h *InterviewHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/interview/start", h.handleStartInterview)
	mux.HandleFunc("POST /api/v1/interview/answer", h.handleAnswerQuestion)
	mux.HandleFunc("GET /api/v1/interview/feedback/{id}", h.handleGetFeedback)
	mux.HandleFunc("GET /api/v1/interview/sessions", h.handleListSessions)
	mux.HandleFunc("GET /api/v1/interview/sessions/{id}", h.handleGetSession)
}

// Helper methods.

// generateQuestions generates questions based on interview type and difficulty.
func (h *InterviewHandler) generateQuestions(interviewType InterviewType, difficulty Difficulty) []Question {
	questions := []Question{}

	switch interviewType {
	case InterviewTypeCoding:
		questions = append(questions, h.getCodingQuestions(difficulty)...)
	case InterviewTypeBehavioral:
		questions = append(questions, h.getBehavioralQuestions(difficulty)...)
	case InterviewTypeSystemDesign:
		questions = append(questions, h.getSystemDesignQuestions(difficulty)...)
	}

	return questions
}

// getCodingQuestions returns coding interview questions.
func (h *InterviewHandler) getCodingQuestions(difficulty Difficulty) []Question {
	baseQuestions := map[Difficulty][]Question{
		DifficultyBeginner: {
			{
				ID:             "c1",
				Content:        "Implement a function to reverse a linked list. You can assume a singly linked list node structure with `next` pointer.",
				Type:           "coding",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Iterative approach", "Pointer manipulation", "Edge cases (empty list, single node)"},
				TimeLimit:      900,
			},
			{
				ID:             "c2",
				Content:        "Write a function to check if a string is a palindrome. Ignore case and non-alphanumeric characters.",
				Type:           "coding",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Two-pointer technique", "String manipulation", "Edge case handling"},
				TimeLimit:      600,
			},
		},
		DifficultyIntermediate: {
			{
				ID:             "c3",
				Content:        "Implement a function to find the longest substring without repeating characters.",
				Type:           "coding",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"Sliding window technique", "Hash map usage", "Time complexity O(n)"},
				TimeLimit:      1200,
			},
			{
				ID:             "c4",
				Content:        "Implement a binary search tree with insert and search operations.",
				Type:           "coding",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"Recursive implementation", "BST properties", "Duplicate handling"},
				TimeLimit:      1500,
			},
		},
		DifficultyAdvanced: {
			{
				ID:             "c5",
				Content:        "Design an LRU (Least Recently Used) cache with get and put operations. It should support O(1) average time complexity.",
				Type:           "coding",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"Hash map + doubly linked list", "O(1) operations", "Capacity management"},
				TimeLimit:      1800,
			},
			{
				ID:             "c6",
				Content:        "Implement a function to serialize and deserialize a binary tree.",
				Type:           "coding",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"BFS/DFS traversal", "String building", "Tree reconstruction"},
				TimeLimit:      1800,
			},
		},
	}

	questions, ok := baseQuestions[difficulty]
	if !ok {
		questions = baseQuestions[DifficultyBeginner]
	}

	return questions
}

// getBehavioralQuestions returns behavioral interview questions.
func (h *InterviewHandler) getBehavioralQuestions(difficulty Difficulty) []Question {
	baseQuestions := map[Difficulty][]Question{
		DifficultyBeginner: {
			{
				ID:             "b1",
				Content:        "Tell me about a time you had to learn something new on the job. How did you approach it?",
				Type:           "behavioral",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Specific example", "Learning process", "Outcome"},
				TimeLimit:      300,
			},
			{
				ID:             "b2",
				Content:        "Describe a project you worked on where you collaborated with others. What was your role?",
				Type:           "behavioral",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Team role", "Collaboration skills", "Specific contribution"},
				TimeLimit:      300,
			},
		},
		DifficultyIntermediate: {
			{
				ID:             "b3",
				Content:        "Tell me about a time you had to deal with a difficult team member. How did you handle the situation?",
				Type:           "behavioral",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"Conflict resolution", "Professional approach", "Communication style"},
				TimeLimit:      300,
			},
			{
				ID:             "b4",
				Content:        "Describe a project where you had to learn a new technology quickly. What was your strategy?",
				Type:           "behavioral",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"Learning strategy", "Time management", "Resource utilization"},
				TimeLimit:      300,
			},
		},
		DifficultyAdvanced: {
			{
				ID:             "b5",
				Content:        "How do you handle tight deadlines and competing priorities? Give a specific example.",
				Type:           "behavioral",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"Prioritization framework", "Communication approach", "Example of successful delivery"},
				TimeLimit:      300,
			},
			{
				ID:             "b6",
				Content:        "Tell me about a time you made a technical decision that was unpopular but necessary. How did you handle the pushback?",
				Type:           "behavioral",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"Decision rationale", "Stakeholder management", "Results"},
				TimeLimit:      300,
			},
		},
	}

	questions, ok := baseQuestions[difficulty]
	if !ok {
		questions = baseQuestions[DifficultyBeginner]
	}

	return questions
}

// getSystemDesignQuestions returns system design interview questions.
func (h *InterviewHandler) getSystemDesignQuestions(difficulty Difficulty) []Question {
	baseQuestions := map[Difficulty][]Question{
		DifficultyBeginner: {
			{
				ID:             "s1",
				Content:        "Design a URL shortener like bit.ly. What are the key components and how do they interact?",
				Type:           "system_design",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Core components", "Data flow", "Scalability considerations"},
				TimeLimit:      600,
			},
			{
				ID:             "s2",
				Content:        "Design a key-value store. What data structures would you use and why?",
				Type:           "system_design",
				Difficulty:     "beginner",
				ExpectedPoints: []string{"Data structure choice", "Hash function", "Collision handling"},
				TimeLimit:      600,
			},
		},
		DifficultyIntermediate: {
			{
				ID:             "s3",
				Content:        "Design a rate limiter for an API. How would you implement it at different scales?",
				Type:           "system_design",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"Rate limiting algorithms", "Distributed approach", "Redis/memcached usage"},
				TimeLimit:      900,
			},
			{
				ID:             "s4",
				Content:        "Design a real-time chat application. What are the main challenges and how would you solve them?",
				Type:           "system_design",
				Difficulty:     "intermediate",
				ExpectedPoints: []string{"WebSocket vs HTTP", "Message delivery guarantees", "Scaling strategy"},
				TimeLimit:      900,
			},
		},
		DifficultyAdvanced: {
			{
				ID:             "s5",
				Content:        "Design a distributed unique ID generator like Twitter's Snowflake. What are the trade-offs?",
				Type:           "system_design",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"Uniqueness guarantee", "Time synchronization", "Collision handling"},
				TimeLimit:      1200,
			},
			{
				ID:             "s6",
				Content:        "Design a real-time analytics pipeline processing 1M events per second. How would you handle backpressure?",
				Type:           "system_design",
				Difficulty:     "advanced",
				ExpectedPoints: []string{"Stream processing", "Windowing strategy", "Exactly-once semantics"},
				TimeLimit:      1200,
			},
		},
	}

	questions, ok := baseQuestions[difficulty]
	if !ok {
		questions = baseQuestions[DifficultyBeginner]
	}

	return questions
}

// calculateAnswerScore calculates a score for an answer (mock implementation).
func (h *InterviewHandler) calculateAnswerScore(question Question, answer string) *int {
	baseScore := 60

	lengthScore := 0
	answerLen := len(answer)
	if answerLen > 50 {
		lengthScore = 20
	} else if answerLen > 20 {
		lengthScore = 10
	}

	keywordScore := 0
	for _, expected := range question.ExpectedPoints {
		if containsKeyword(answer, expected) {
			keywordScore += 5
		}
	}

	totalScore := baseScore + lengthScore + keywordScore
	if totalScore > 100 {
		totalScore = 100
	}

	return &totalScore
}

// generateFeedback generates feedback for an answer (mock implementation).
func (h *InterviewHandler) generateFeedback(question Question, answer string) string {
	if len(answer) < 20 {
		return "Your answer is brief. Consider providing more detail and examples to strengthen your response."
	}

	if len(question.ExpectedPoints) > 0 {
		mentionedCount := 0
		for _, expected := range question.ExpectedPoints {
			if containsKeyword(answer, expected) {
				mentionedCount++
			}
		}
		if mentionedCount < len(question.ExpectedPoints)/2 {
			maxIdx := min(2, len(question.ExpectedPoints))
			pointsSlice := question.ExpectedPoints[:maxIdx]
			return fmt.Sprintf("You covered %d of %d expected points. Try to address: %s",
				mentionedCount, len(question.ExpectedPoints), joinExpectedPoints(pointsSlice))
		}
	}

	return "Good answer! You've covered the main points effectively."
}

// calculateOverallScore calculates the overall interview score.
func (h *InterviewHandler) calculateOverallScore(answers []Answer) *int {
	if len(answers) == 0 {
		return nil
	}

	totalScore := 0
	for _, ans := range answers {
		if ans.Score != nil {
			totalScore += *ans.Score
		}
	}

	avgScore := totalScore / len(answers)
	return &avgScore
}

// generateInterviewFeedback generates comprehensive feedback for a completed interview.
func (h *InterviewHandler) generateInterviewFeedback(session *InterviewSession) *InterviewFeedback {
	strengths := []string{}
	improvements := []string{}
	detailedFeedback := []QuestionFeedback{}

	totalScore := 0
	for _, ans := range session.Answers {
		if ans.Score != nil {
			totalScore += *ans.Score
		}
	}
	overallScore := 0
	if len(session.Answers) > 0 {
		overallScore = totalScore / len(session.Answers)
	}

	if overallScore >= 80 {
		strengths = append(strengths, "Strong overall performance across all questions")
		strengths = append(strengths, "Good communication and structure in answers")
	} else if overallScore >= 60 {
		strengths = append(strengths, "Solid foundational knowledge")
		improvements = append(improvements, "Consider providing more detailed examples")
	} else {
		improvements = append(improvements, "Focus on strengthening foundational concepts")
		improvements = append(improvements, "Practice explaining technical concepts clearly")
	}

	switch session.Type {
	case InterviewTypeCoding:
		if overallScore >= 70 {
			strengths = append(strengths, "Good algorithmic thinking")
		} else {
			improvements = append(improvements, "Practice more algorithmic problems")
		}
	case InterviewTypeBehavioral:
		if overallScore >= 70 {
			strengths = append(strengths, "Good situational awareness")
		} else {
			improvements = append(improvements, "Use STAR method for behavioral questions")
		}
	case InterviewTypeSystemDesign:
		if overallScore >= 70 {
			strengths = append(strengths, "Good system thinking")
		} else {
			improvements = append(improvements, "Study common system design patterns")
		}
	}

	for i, ans := range session.Answers {
		if i >= len(session.Questions) {
			break
		}
		question := session.Questions[i]

		qf := QuestionFeedback{
			QuestionID: question.ID,
			Score:      getScoreValue(ans.Score),
			Feedback:   ans.Feedback,
		}

		if ans.Score != nil && *ans.Score >= 70 {
			qf.Strengths = question.ExpectedPoints
		} else if ans.Score != nil {
			missed := []string{}
			for _, ep := range question.ExpectedPoints {
				if !containsKeyword(ans.Content, ep) {
					missed = append(missed, ep)
				}
			}
			qf.Missed = missed
		}

		detailedFeedback = append(detailedFeedback, qf)
	}

	return &InterviewFeedback{
		SessionID:        session.ID,
		OverallScore:     overallScore,
		Strengths:        strengths,
		Improvements:     improvements,
		DetailedFeedback: detailedFeedback,
	}
}

// writeSuccessResponse writes a successful API response.
func (h *InterviewHandler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	response := &domain.APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: logger.GetRequestID(r.Context()),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error API response.
func (h *InterviewHandler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierrors.APIError
	var statusCode int

	if errAs(err, &apiErr) {
		statusCode = apiErr.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
		apiErr = apierrors.NewInternalError("internal server error", err)
	}

	response := &domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Type:    apiErr.Type,
			Message: apiErr.Message,
		},
		RequestID: logger.GetRequestID(r.Context()),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)

	logger.LogError(h.logger, r.Context(), apiErr, "HTTP request error",
		"status_code", statusCode,
		"error_type", apiErr.Type,
	)
}

// Utility functions.

// generateID generates a unique ID for interview sessions.
func generateID() string {
	return fmt.Sprintf("int_%d", time.Now().UnixNano())
}

// containsKeyword checks if answer contains a keyword (case-insensitive).
func containsKeyword(answer, keyword string) bool {
	answerLower := toLower(answer)
	keywordLower := toLower(keyword)

	for i := 0; i <= len(answerLower)-len(keywordLower); i++ {
		if answerLower[i:i+len(keywordLower)] == keywordLower {
			return true
		}
	}
	return false
}

// toLower converts string to lowercase (simplified).
func toLower(s string) string {
	result := []rune{}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+('a'-'A'))
		} else if r >= 'a' && r <= 'z' {
			result = append(result, r)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// joinExpectedPoints joins expected points for feedback.
func joinExpectedPoints(points []string) string {
	result := ""
	for i, p := range points {
		if i > 0 {
			result += ", "
		}
		result += "`" + p + "`"
	}
	return result
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getScoreValue returns score value or 0.
func getScoreValue(score *int) int {
	if score == nil {
		return 0
	}
	return *score
}

// errAs checks if error can be asserted to APIError.
func errAs(err error, target interface{}) bool {
	return false
}
