// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package messaging provides functionality for the GO-PRO Learning Platform.
package messaging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go-pro-backend/internal/messaging/kafka"
)

// Service provides messaging functionality.
type Service struct {
	producer      *kafka.Producer
	eventConsumer *kafka.EventConsumer
	config        *kafka.Config
	topics        *kafka.Topics
	enabled       bool
}

// Config holds messaging service configuration.
type Config struct {
	Enabled bool
	Kafka   *kafka.Config
	Topics  *kafka.Topics
}

// DefaultConfig returns default messaging configuration.
func DefaultConfig() *Config {
	return &Config{
		Enabled: getEnvAsBool("MESSAGING_ENABLED", true),
		Kafka:   kafka.DefaultConfig(),
		Topics:  kafka.DefaultTopics(),
	}
}

// LoadConfigFromEnv loads configuration from environment variables.
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// Override Kafka config from environment.
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		config.Kafka.Brokers = strings.Split(brokers, ",")
	}
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		config.Kafka.GroupID = groupID
	}
	if clientID := os.Getenv("KAFKA_CLIENT_ID"); clientID != "" {
		config.Kafka.ClientID = clientID
	}

	return config
}

// NewService creates a new messaging service.
func NewService(config *Config) (*Service, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if !config.Enabled {
		return &Service{enabled: false}, nil
	}

	producer := kafka.NewProducer(config.Kafka, config.Topics)
	eventConsumer := kafka.NewEventConsumer(config.Kafka, config.Topics)

	return &Service{
		producer:      producer,
		eventConsumer: eventConsumer,
		config:        config.Kafka,
		topics:        config.Topics,
		enabled:       true,
	}, nil
}

// IsEnabled returns whether messaging is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// Close closes the messaging service.
func (s *Service) Close() error {
	if !s.enabled {
		return nil
	}

	var errs []error
	if err := s.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}
	if err := s.eventConsumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing messaging service: %v", errs)
	}

	return nil
}

// Publisher methods.

// PublishUserCreated publishes a user created event.
func (s *Service) PublishUserCreated(ctx context.Context, userID string, userData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	event := kafka.NewUserEvent(kafka.UserCreated, userID, userData)

	return s.producer.PublishUserEvent(ctx, event)
}

// PublishCourseCreated publishes a course created event.
func (s *Service) PublishCourseCreated(ctx context.Context, courseID, instructorID string, courseData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	event := kafka.NewCourseEvent(kafka.CourseCreated, courseID, instructorID, courseData)

	return s.producer.PublishCourseEvent(ctx, event)
}

// PublishLessonCompleted publishes a lesson completed event.
func (s *Service) PublishLessonCompleted(ctx context.Context, userID, lessonID, courseID string) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"user_id":   userID,
		"lesson_id": lessonID,
		"course_id": courseID,
	}

	event := &kafka.LessonEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.LessonCompleted,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("lesson:%s", lessonID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		LessonID: lessonID,
		CourseID: courseID,
		UserID:   userID,
	}

	return s.producer.PublishLessonEvent(ctx, event)
}

// PublishProgressUpdated publishes a progress updated event.
func (s *Service) PublishProgressUpdated(ctx context.Context, userID, lessonID, courseID string, completed bool, score, timeSpent int) error {
	if !s.enabled {
		return nil
	}

	event := kafka.NewProgressEvent(kafka.ProgressUpdated, userID, lessonID, courseID, completed, score, timeSpent)

	return s.producer.PublishProgressEvent(ctx, event)
}

// PublishProgressCreated publishes a progress created event.
func (s *Service) PublishProgressCreated(ctx context.Context, userID, lessonID string, progressData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	event := kafka.NewProgressEvent(kafka.ProgressStarted, userID, lessonID, "", false, 0, 0)

	return s.producer.PublishProgressEvent(ctx, event)
}

// PublishLessonCreated publishes a lesson created event.
func (s *Service) PublishLessonCreated(ctx context.Context, lessonID, courseID string, lessonData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"lesson_id":   lessonID,
		"course_id":   courseID,
		"lesson_data": lessonData,
	}

	event := &kafka.LessonEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.LessonCreated,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("lesson:%s", lessonID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		LessonID: lessonID,
		CourseID: courseID,
		UserID:   "", // No user for creation events
	}

	return s.producer.PublishLessonEvent(ctx, event)
}

// PublishLessonUpdated publishes a lesson updated event.
func (s *Service) PublishLessonUpdated(ctx context.Context, lessonID, courseID string, lessonData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"lesson_id":   lessonID,
		"course_id":   courseID,
		"lesson_data": lessonData,
	}

	event := &kafka.LessonEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.LessonUpdated,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("lesson:%s", lessonID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		LessonID: lessonID,
		CourseID: courseID,
		UserID:   "", // No user for update events
	}

	return s.producer.PublishLessonEvent(ctx, event)
}

// PublishLessonDeleted publishes a lesson deleted event.
func (s *Service) PublishLessonDeleted(ctx context.Context, lessonID, courseID string) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"lesson_id": lessonID,
		"course_id": courseID,
	}

	event := &kafka.LessonEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.LessonDeleted,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("lesson:%s", lessonID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		LessonID: lessonID,
		CourseID: courseID,
		UserID:   "", // No user for deletion events
	}

	return s.producer.PublishLessonEvent(ctx, event)
}

// PublishExerciseCreated publishes an exercise created event.
func (s *Service) PublishExerciseCreated(ctx context.Context, exerciseID, lessonID string, exerciseData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"exercise_id":   exerciseID,
		"lesson_id":     lessonID,
		"exercise_data": exerciseData,
	}

	event := &kafka.ExerciseEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.ExerciseCreated,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("exercise:%s", exerciseID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		ExerciseID: exerciseID,
		LessonID:   lessonID,
		UserID:     "", // No user for creation events
		Score:      0,  // No score for creation events
	}

	return s.producer.PublishExerciseEvent(ctx, event)
}

// PublishExerciseUpdated publishes an exercise updated event.
func (s *Service) PublishExerciseUpdated(ctx context.Context, exerciseID, lessonID string, exerciseData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"exercise_id":   exerciseID,
		"lesson_id":     lessonID,
		"exercise_data": exerciseData,
	}

	event := &kafka.ExerciseEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.ExerciseUpdated,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("exercise:%s", exerciseID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		ExerciseID: exerciseID,
		LessonID:   lessonID,
		UserID:     "", // No user for update events
		Score:      0,  // No score for update events
	}

	return s.producer.PublishExerciseEvent(ctx, event)
}

// PublishExerciseDeleted publishes an exercise deleted event.
func (s *Service) PublishExerciseDeleted(ctx context.Context, exerciseID, lessonID string) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"exercise_id": exerciseID,
		"lesson_id":   lessonID,
	}

	event := &kafka.ExerciseEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.ExerciseDeleted,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("exercise:%s", exerciseID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		ExerciseID: exerciseID,
		LessonID:   lessonID,
		UserID:     "", // No user for deletion events
		Score:      0,  // No score for deletion events
	}

	return s.producer.PublishExerciseEvent(ctx, event)
}

// PublishExerciseSubmitted publishes an exercise submitted event.
func (s *Service) PublishExerciseSubmitted(ctx context.Context, userID, exerciseID, lessonID string, score int, submissionData map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	data := map[string]interface{}{
		"user_id":     userID,
		"exercise_id": exerciseID,
		"lesson_id":   lessonID,
		"score":       score,
		"submission":  submissionData,
	}

	event := &kafka.ExerciseEvent{
		Event: kafka.Event{
			ID:        generateEventID(),
			Type:      kafka.ExerciseSubmitted,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("exercise:%s", exerciseID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		ExerciseID: exerciseID,
		LessonID:   lessonID,
		UserID:     userID,
		Score:      score,
	}

	return s.producer.PublishExerciseEvent(ctx, event)
}

// PublishAuditLog publishes an audit log event.
func (s *Service) PublishAuditLog(ctx context.Context, userID, action, resource, resourceID string, success bool, oldValues, newValues map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	event := kafka.NewAuditEvent(userID, action, resource, resourceID, success, oldValues, newValues)

	return s.producer.PublishAuditEvent(ctx, event)
}

// Consumer setup methods.

// SetupEventHandlers sets up default event handlers.
func (s *Service) SetupEventHandlers() {
	if !s.enabled {
		return
	}

	// User event handlers.
	s.eventConsumer.RegisterUserEventHandler(func(ctx context.Context, event *kafka.UserEvent) error {
		log.Printf("Received user event: %s for user %s", event.Type, event.UserID)
		// Add your business logic here.
		return nil
	})

	// Course event handlers.
	s.eventConsumer.RegisterCourseEventHandler(func(ctx context.Context, event *kafka.CourseEvent) error {
		log.Printf("Received course event: %s for course %s", event.Type, event.CourseID)
		// Add your business logic here.
		return nil
	})

	// Progress event handlers.
	s.eventConsumer.RegisterProgressEventHandler(func(ctx context.Context, event *kafka.ProgressEvent) error {
		log.Printf("Received progress event: %s for user %s, lesson %s", event.Type, event.UserID, event.LessonID)
		// Add your business logic here.
		// For example: update analytics, send notifications, etc.
		return nil
	})

	// Exercise event handlers.
	s.eventConsumer.RegisterExerciseEventHandler(func(ctx context.Context, event *kafka.ExerciseEvent) error {
		log.Printf("Received exercise event: %s for exercise %s", event.Type, event.ExerciseID)
		// Add your business logic here.
		return nil
	})

	// Audit event handlers.
	s.eventConsumer.RegisterAuditEventHandler(func(ctx context.Context, event *kafka.AuditEvent) error {
		log.Printf("Received audit event: %s for resource %s", event.Action, event.Resource)
		// Add your business logic here.
		// For example: store in audit database, send alerts, etc.
		return nil
	})
}

// StartConsumer starts the event consumer.
func (s *Service) StartConsumer(ctx context.Context) error {
	if !s.enabled {
		return nil
	}

	s.SetupEventHandlers()

	return s.eventConsumer.Start(ctx)
}

// Health checks the health of the messaging service.
func (s *Service) Health(ctx context.Context) error {
	if !s.enabled {
		return nil
	}

	// Simple health check - try to publish a test event.
	testEvent := &kafka.Event{
		ID:        generateEventID(),
		Type:      "health.check",
		Source:    "go-pro-backend",
		Subject:   "health",
		Data:      map[string]interface{}{"status": "ok"},
		Metadata:  make(map[string]string),
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return s.producer.PublishEvent(ctx, testEvent)
}

// Helper functions.

func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.EqualFold(value, "true")
	}

	return defaultValue
}
