// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package testutils provides testing utilities and helpers.
package testutils

import "time"

// Domain models for testing (duplicated from main package to avoid circular imports).
type Course struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Lessons     []string `json:"lessons"`
	CreatedAt   string   `json:"created_at"`
}

type Lesson struct {
	ID          string `json:"id"`
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Order       int    `json:"order"`
	CreatedAt   string `json:"created_at"`
}

type Exercise struct {
	ID          string `json:"id"`
	LessonID    string `json:"lesson_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TestCases   int    `json:"test_cases"`
	Difficulty  string `json:"difficulty"`
}

type Progress struct {
	UserID      string    `json:"user_id"`
	LessonID    string    `json:"lesson_id"`
	Completed   bool      `json:"completed"`
	Score       int       `json:"score"`
	CompletedAt time.Time `json:"completed_at"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
