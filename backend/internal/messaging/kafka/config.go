// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"time"
)

// Config holds Kafka configuration.
type Config struct {
	Brokers             []string
	GroupID             string
	ClientID            string
	SecurityProtocol    string
	SASLMechanism       string
	SASLUsername        string
	SASLPassword        string
	EnableAutoCommit    bool
	AutoCommitInterval  time.Duration
	SessionTimeout      time.Duration
	HeartbeatInterval   time.Duration
	MaxPollRecords      int
	FetchMinBytes       int
	FetchMaxWait        time.Duration
	RetryBackoff        time.Duration
	ReconnectBackoff    time.Duration
	MaxRetries          int
	RequestTimeout      time.Duration
	MetadataRefreshFreq time.Duration
	EnableIdempotence   bool
	Acks                string
	CompressionType     string
	BatchSize           int
	LingerMs            time.Duration
	BufferMemory        int64
	MaxBlockMs          time.Duration
	MaxRequestSize      int
	ReceiveBufferBytes  int
	SendBufferBytes     int
}

// DefaultConfig returns a default Kafka configuration.
func DefaultConfig() *Config {
	return &Config{
		Brokers:             []string{"localhost:9092"},
		GroupID:             "go-pro-consumer-group",
		ClientID:            "go-pro-client",
		SecurityProtocol:    "PLAINTEXT",
		EnableAutoCommit:    true,
		AutoCommitInterval:  1 * time.Second,
		SessionTimeout:      30 * time.Second,
		HeartbeatInterval:   3 * time.Second,
		MaxPollRecords:      500,
		FetchMinBytes:       1,
		FetchMaxWait:        500 * time.Millisecond,
		RetryBackoff:        100 * time.Millisecond,
		ReconnectBackoff:    50 * time.Millisecond,
		MaxRetries:          3,
		RequestTimeout:      30 * time.Second,
		MetadataRefreshFreq: 5 * time.Minute,
		EnableIdempotence:   true,
		Acks:                "all",
		CompressionType:     "snappy",
		BatchSize:           16384,
		LingerMs:            5 * time.Millisecond,
		BufferMemory:        33554432, // 32MB
		MaxBlockMs:          60 * time.Second,
		MaxRequestSize:      1048576, // 1MB
		ReceiveBufferBytes:  65536,   // 64KB
		SendBufferBytes:     131072,  // 128KB
	}
}

// Topics defines the Kafka topics used by the application.
type Topics struct {
	UserEvents         string
	CourseEvents       string
	LessonEvents       string
	ExerciseEvents     string
	ProgressEvents     string
	NotificationEvents string
	AuditEvents        string
}

// DefaultTopics returns the default topic configuration.
func DefaultTopics() *Topics {
	return &Topics{
		UserEvents:         "go-pro.user.events",
		CourseEvents:       "go-pro.course.events",
		LessonEvents:       "go-pro.lesson.events",
		ExerciseEvents:     "go-pro.exercise.events",
		ProgressEvents:     "go-pro.progress.events",
		NotificationEvents: "go-pro.notification.events",
		AuditEvents:        "go-pro.audit.events",
	}
}

// EventType represents the type of event.
type EventType string

const (
	// User events.
	UserCreated   EventType = "user.created"
	UserUpdated   EventType = "user.updated"
	UserDeleted   EventType = "user.deleted"
	UserLoggedIn  EventType = "user.logged_in"
	UserLoggedOut EventType = "user.logged_out"

	// Course events.
	CourseCreated   EventType = "course.created"
	CourseUpdated   EventType = "course.updated"
	CourseDeleted   EventType = "course.deleted"
	CoursePublished EventType = "course.published"

	// Lesson events.
	LessonCreated   EventType = "lesson.created"
	LessonUpdated   EventType = "lesson.updated"
	LessonDeleted   EventType = "lesson.deleted"
	LessonCompleted EventType = "lesson.completed"

	// Exercise events.
	ExerciseCreated   EventType = "exercise.created"
	ExerciseUpdated   EventType = "exercise.updated"
	ExerciseDeleted   EventType = "exercise.deleted"
	ExerciseSubmitted EventType = "exercise.submitted"
	ExerciseCompleted EventType = "exercise.completed"

	// Progress events.
	ProgressStarted   EventType = "progress.started"
	ProgressUpdated   EventType = "progress.updated"
	ProgressCompleted EventType = "progress.completed"

	// Notification events.
	NotificationSent EventType = "notification.sent"
	NotificationRead EventType = "notification.read"

	// Audit events.
	AuditLog EventType = "audit.log"
)

// Event represents a generic event structure.
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Source    string                 `json:"source"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata"`
	Timestamp int64                  `json:"timestamp"`
	Version   string                 `json:"version"`
}

// UserEvent represents user-related events.
type UserEvent struct {
	Event
	UserID string `json:"user_id"`
}

// CourseEvent represents course-related events.
type CourseEvent struct {
	Event
	CourseID     string `json:"course_id"`
	InstructorID string `json:"instructor_id,omitempty"`
}

// LessonEvent represents lesson-related events.
type LessonEvent struct {
	Event
	LessonID string `json:"lesson_id"`
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id,omitempty"`
}

// ExerciseEvent represents exercise-related events.
type ExerciseEvent struct {
	Event
	ExerciseID string `json:"exercise_id"`
	LessonID   string `json:"lesson_id"`
	UserID     string `json:"user_id,omitempty"`
	Score      int    `json:"score,omitempty"`
}

// ProgressEvent represents progress-related events.
type ProgressEvent struct {
	Event
	UserID    string `json:"user_id"`
	LessonID  string `json:"lesson_id"`
	CourseID  string `json:"course_id"`
	Completed bool   `json:"completed"`
	Score     int    `json:"score"`
	TimeSpent int    `json:"time_spent_seconds"`
}

// NotificationEvent represents notification-related events.
type NotificationEvent struct {
	Event
	UserID   string `json:"user_id"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Type     string `json:"notification_type"`
	Channel  string `json:"channel"`
	Priority string `json:"priority"`
}

// AuditEvent represents audit log events.
type AuditEvent struct {
	Event
	UserID       string                 `json:"user_id,omitempty"`
	Action       string                 `json:"action"`
	Resource     string                 `json:"resource"`
	ResourceID   string                 `json:"resource_id,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	OldValues    map[string]interface{} `json:"old_values,omitempty"`
	NewValues    map[string]interface{} `json:"new_values,omitempty"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}
