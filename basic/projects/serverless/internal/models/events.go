package models

import "time"

// Event represents a learning platform event
type Event struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}

// LessonCompletedPayload represents a completed lesson event
type LessonCompletedPayload struct {
	UserID    string `json:"user_id"`
	LessonID  string `json:"lesson_id"`
	Score     int    `json:"score"`
	Completed bool   `json:"completed"`
}

// CourseEnrolledPayload represents a course enrollment event
type CourseEnrolledPayload struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

// HealthResponse represents the Lambda health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}
