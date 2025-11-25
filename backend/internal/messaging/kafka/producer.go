// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Producer wraps kafka-go writer with additional functionality.
type Producer struct {
	writer *kafka.Writer
	config *Config
	topics *Topics
}

// NewProducer creates a new Kafka producer.
func NewProducer(config *Config, topics *Topics) *Producer {
	if config == nil {
		config = DefaultConfig()
	}
	if topics == nil {
		topics = DefaultTopics()
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		BatchSize:              config.BatchSize,
		BatchTimeout:           config.LingerMs,
		ReadTimeout:            config.RequestTimeout,
		WriteTimeout:           config.RequestTimeout,
		Compression:            getCompressionCodec(config.CompressionType),
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
		config: config,
		topics: topics,
	}
}

// Close closes the producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// PublishEvent publishes a generic event to the appropriate topic.
func (p *Producer) PublishEvent(ctx context.Context, event *Event) error {
	topic := p.getTopicForEventType(event.Type)
	return p.publishToTopic(ctx, topic, event.ID, event)
}

// PublishUserEvent publishes a user event.
func (p *Producer) PublishUserEvent(ctx context.Context, event *UserEvent) error {
	return p.publishToTopic(ctx, p.topics.UserEvents, event.ID, event)
}

// PublishCourseEvent publishes a course event.
func (p *Producer) PublishCourseEvent(ctx context.Context, event *CourseEvent) error {
	return p.publishToTopic(ctx, p.topics.CourseEvents, event.ID, event)
}

// PublishLessonEvent publishes a lesson event.
func (p *Producer) PublishLessonEvent(ctx context.Context, event *LessonEvent) error {
	return p.publishToTopic(ctx, p.topics.LessonEvents, event.ID, event)
}

// PublishExerciseEvent publishes an exercise event.
func (p *Producer) PublishExerciseEvent(ctx context.Context, event *ExerciseEvent) error {
	return p.publishToTopic(ctx, p.topics.ExerciseEvents, event.ID, event)
}

// PublishProgressEvent publishes a progress event.
func (p *Producer) PublishProgressEvent(ctx context.Context, event *ProgressEvent) error {
	return p.publishToTopic(ctx, p.topics.ProgressEvents, event.ID, event)
}

// PublishNotificationEvent publishes a notification event.
func (p *Producer) PublishNotificationEvent(ctx context.Context, event *NotificationEvent) error {
	return p.publishToTopic(ctx, p.topics.NotificationEvents, event.ID, event)
}

// PublishAuditEvent publishes an audit event.
func (p *Producer) PublishAuditEvent(ctx context.Context, event *AuditEvent) error {
	return p.publishToTopic(ctx, p.topics.AuditEvents, event.ID, event)
}

// PublishBatch publishes multiple events in a single batch.
func (p *Producer) PublishBatch(ctx context.Context, events []EventMessage) error {
	messages := make([]kafka.Message, len(events))

	for i, event := range events {
		data, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal event %d: %w", i, err)
		}

		messages[i] = kafka.Message{
			Topic: event.Topic,
			Key:   []byte(event.Key),
			Value: data,
			Headers: []kafka.Header{
				{Key: "event-type", Value: []byte(event.EventType)},
				{Key: "source", Value: []byte("go-pro-backend")},
				{Key: "timestamp", Value: []byte(fmt.Sprintf("%d", time.Now().Unix()))},
			},
		}
	}

	return p.writer.WriteMessages(ctx, messages...)
}

// publishToTopic publishes data to a specific topic.
func (p *Producer) publishToTopic(ctx context.Context, topic, key string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "source", Value: []byte("go-pro-backend")},
			{Key: "timestamp", Value: []byte(fmt.Sprintf("%d", time.Now().Unix()))},
		},
	}

	return p.writer.WriteMessages(ctx, message)
}

// getTopicForEventType returns the appropriate topic for an event type.
func (p *Producer) getTopicForEventType(eventType EventType) string {
	switch eventType {
	case UserCreated, UserUpdated, UserDeleted, UserLoggedIn, UserLoggedOut:
		return p.topics.UserEvents
	case CourseCreated, CourseUpdated, CourseDeleted, CoursePublished:
		return p.topics.CourseEvents
	case LessonCreated, LessonUpdated, LessonDeleted, LessonCompleted:
		return p.topics.LessonEvents
	case ExerciseCreated, ExerciseUpdated, ExerciseDeleted, ExerciseSubmitted, ExerciseCompleted:
		return p.topics.ExerciseEvents
	case ProgressStarted, ProgressUpdated, ProgressCompleted:
		return p.topics.ProgressEvents
	case NotificationSent, NotificationRead:
		return p.topics.NotificationEvents
	case AuditLog:
		return p.topics.AuditEvents
	default:
		return p.topics.AuditEvents // Default to audit events
	}
}

// getCompressionCodec returns the compression codec for the given type.
func getCompressionCodec(compressionType string) kafka.Compression {
	switch compressionType {
	case "gzip":
		return kafka.Gzip
	case "snappy":
		return kafka.Snappy
	case "lz4":
		return kafka.Lz4
	case "zstd":
		return kafka.Zstd
	default:
		return kafka.Snappy
	}
}

// EventMessage represents a message to be published.
type EventMessage struct {
	Topic     string
	Key       string
	Data      interface{}
	EventType string
}

// Helper functions to create events.

// NewUserEvent creates a new user event.
func NewUserEvent(eventType EventType, userID string, data map[string]interface{}) *UserEvent {
	return &UserEvent{
		Event: Event{
			ID:        uuid.New().String(),
			Type:      eventType,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("user:%s", userID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		UserID: userID,
	}
}

// NewCourseEvent creates a new course event.
func NewCourseEvent(eventType EventType, courseID, instructorID string, data map[string]interface{}) *CourseEvent {
	return &CourseEvent{
		Event: Event{
			ID:        uuid.New().String(),
			Type:      eventType,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("course:%s", courseID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		CourseID:     courseID,
		InstructorID: instructorID,
	}
}

// NewProgressEvent creates a new progress event.
func NewProgressEvent(eventType EventType, userID, lessonID, courseID string, completed bool, score, timeSpent int) *ProgressEvent {
	data := map[string]interface{}{
		"completed":  completed,
		"score":      score,
		"time_spent": timeSpent,
	}

	return &ProgressEvent{
		Event: Event{
			ID:        uuid.New().String(),
			Type:      eventType,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("progress:%s:%s", userID, lessonID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		UserID:    userID,
		LessonID:  lessonID,
		CourseID:  courseID,
		Completed: completed,
		Score:     score,
		TimeSpent: timeSpent,
	}
}

// NewAuditEvent creates a new audit event.
func NewAuditEvent(userID, action, resource, resourceID string, success bool, oldValues, newValues map[string]interface{}) *AuditEvent {
	data := map[string]interface{}{
		"action":   action,
		"resource": resource,
		"success":  success,
	}

	if oldValues != nil {
		data["old_values"] = oldValues
	}
	if newValues != nil {
		data["new_values"] = newValues
	}

	return &AuditEvent{
		Event: Event{
			ID:        uuid.New().String(),
			Type:      AuditLog,
			Source:    "go-pro-backend",
			Subject:   fmt.Sprintf("audit:%s:%s", resource, resourceID),
			Data:      data,
			Metadata:  make(map[string]string),
			Timestamp: time.Now().Unix(),
			Version:   "1.0",
		},
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		OldValues:  oldValues,
		NewValues:  newValues,
		Success:    success,
	}
}
