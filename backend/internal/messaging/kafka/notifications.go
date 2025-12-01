// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"context"
	"fmt"
	"log"
	"time"
)

// NotificationHandler handles real-time notifications.
type NotificationHandler interface {
	SendNotification(ctx context.Context, notification *NotificationEvent) error
}

// NotificationService manages real-time notifications based on Kafka events.
type NotificationService struct {
	producer           *Producer
	eventConsumer      *EventConsumer
	handlers           []NotificationHandler
	progressThresholds map[string]int // lessonID -> threshold for milestone notifications
}

// NewNotificationService creates a new notification service.
func NewNotificationService(config *Config, topics *Topics) *NotificationService {
	producer := NewProducer(config, topics)
	eventConsumer := NewEventConsumer(config, topics)

	return &NotificationService{
		producer:           producer,
		eventConsumer:      eventConsumer,
		handlers:           make([]NotificationHandler, 0),
		progressThresholds: make(map[string]int),
	}
}

// RegisterHandler registers a notification handler.
func (ns *NotificationService) RegisterHandler(handler NotificationHandler) {
	ns.handlers = append(ns.handlers, handler)
}

// SetProgressThreshold sets a threshold for milestone notifications.
func (ns *NotificationService) SetProgressThreshold(lessonID string, threshold int) {
	ns.progressThresholds[lessonID] = threshold
}

// Start starts the notification service.
func (ns *NotificationService) Start(ctx context.Context) error {
	// Register handlers for various events
	ns.registerProgressEventHandler(ctx)
	ns.registerExerciseEventHandler(ctx)
	ns.registerLessonEventHandler(ctx)

	return ns.eventConsumer.Start(ctx)
}

// sendNotification sends a notification through registered handlers.
func (ns *NotificationService) sendNotification(ctx context.Context, notification *NotificationEvent) error {
	// First, publish to Kafka so it's stored
	if err := ns.producer.PublishNotificationEvent(ctx, notification); err != nil {
		log.Printf("Error publishing notification to Kafka: %v", err)
		// Continue anyway, try to send through handlers
	}

	// Send through all registered handlers
	var errs []error
	for _, handler := range ns.handlers {
		if err := handler.SendNotification(ctx, notification); err != nil {
			log.Printf("Error sending notification through handler: %v", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors sending notifications: %v", errs)
	}

	return nil
}

// registerProgressEventHandler sets up handlers for progress events.
func (ns *NotificationService) registerProgressEventHandler(ctx context.Context) {
	ns.eventConsumer.RegisterProgressEventHandler(func(ctx context.Context, event *ProgressEvent) error {
		// Create milestone notifications
		notification := &NotificationEvent{
			Event: Event{
				ID:        generateEventID(),
				Type:      NotificationSent,
				Source:    "go-pro-backend",
				Subject:   fmt.Sprintf("user:%s", event.UserID),
				Data:      make(map[string]interface{}),
				Metadata:  make(map[string]string),
				Timestamp: time.Now().Unix(),
				Version:   "1.0",
			},
			UserID:   event.UserID,
			Title:    "Progress Update",
			Message:  fmt.Sprintf("You completed lesson %s with score %d", event.LessonID, event.Score),
			Type:     "progress",
			Channel:  "in-app",
			Priority: "normal",
		}

		// Check for milestone
		if event.Score >= 80 {
			notification.Message = fmt.Sprintf("Great job! You scored %d on lesson %s", event.Score, event.LessonID)
			notification.Priority = "high"
		}

		if event.Completed {
			notification.Message = fmt.Sprintf("Congratulations! You completed lesson %s", event.LessonID)
			notification.Priority = "high"
		}

		return ns.sendNotification(ctx, notification)
	})
}

// registerExerciseEventHandler sets up handlers for exercise events.
func (ns *NotificationService) registerExerciseEventHandler(ctx context.Context) {
	ns.eventConsumer.RegisterExerciseEventHandler(func(ctx context.Context, event *ExerciseEvent) error {
		if event.Type != ExerciseSubmitted && event.Type != ExerciseCompleted {
			return nil
		}

		notification := &NotificationEvent{
			Event: Event{
				ID:        generateEventID(),
				Type:      NotificationSent,
				Source:    "go-pro-backend",
				Subject:   fmt.Sprintf("user:%s", event.UserID),
				Data:      make(map[string]interface{}),
				Metadata:  make(map[string]string),
				Timestamp: time.Now().Unix(),
				Version:   "1.0",
			},
			UserID:   event.UserID,
			Title:    "Exercise Results",
			Message:  fmt.Sprintf("Your exercise was graded. Score: %d", event.Score),
			Type:     "exercise",
			Channel:  "in-app",
			Priority: "normal",
		}

		if event.Score >= 90 {
			notification.Priority = "high"
			notification.Title = "Excellent Performance!"
		} else if event.Score < 50 {
			notification.Priority = "high"
			notification.Title = "Exercise Needs Review"
		}

		return ns.sendNotification(ctx, notification)
	})
}

// registerLessonEventHandler sets up handlers for lesson events.
func (ns *NotificationService) registerLessonEventHandler(ctx context.Context) {
	ns.eventConsumer.RegisterLessonEventHandler(func(ctx context.Context, event *LessonEvent) error {
		if event.Type != LessonCompleted {
			return nil
		}

		notification := &NotificationEvent{
			Event: Event{
				ID:        generateEventID(),
				Type:      NotificationSent,
				Source:    "go-pro-backend",
				Subject:   fmt.Sprintf("user:%s", event.UserID),
				Data:      make(map[string]interface{}),
				Metadata:  make(map[string]string),
				Timestamp: time.Now().Unix(),
				Version:   "1.0",
			},
			UserID:   event.UserID,
			Title:    "Lesson Completed",
			Message:  fmt.Sprintf("You have completed lesson %s. Keep up the great work!", event.LessonID),
			Type:     "lesson",
			Channel:  "in-app",
			Priority: "normal",
		}

		return ns.sendNotification(ctx, notification)
	})
}

// Close closes the notification service.
func (ns *NotificationService) Close() error {
	var errs []error

	if err := ns.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}

	if err := ns.eventConsumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close event consumer: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing notification service: %v", errs)
	}

	return nil
}

// LogNotificationHandler is a simple handler that logs notifications.
type LogNotificationHandler struct {
}

// SendNotification logs the notification.
func (h *LogNotificationHandler) SendNotification(ctx context.Context, notification *NotificationEvent) error {
	log.Printf("NOTIFICATION [%s]: %s - %s (Priority: %s)", notification.Type, notification.Title, notification.Message, notification.Priority)
	return nil
}

// PrintNotificationHandler prints notifications to stdout.
type PrintNotificationHandler struct {
}

// SendNotification prints the notification.
func (h *PrintNotificationHandler) SendNotification(ctx context.Context, notification *NotificationEvent) error {
	fmt.Printf("📬 [%s] %s\n   %s\n   User: %s | Priority: %s\n\n",
		notification.Type,
		notification.Title,
		notification.Message,
		notification.UserID,
		notification.Priority)
	return nil
}
