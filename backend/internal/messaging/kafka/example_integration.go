// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

// IntegrationExample demonstrates real-time data processing with Kafka.
type IntegrationExample struct {
	config                *Config
	topics                *Topics
	producer              *Producer
	analyticsConsumer     *AnalyticsConsumer
	notificationService   *NotificationService
	safeConsumer          *SafeConsumer
	dlqHandler            *DeadLetterQueueHandler
}

// NewIntegrationExample creates a new integration example.
func NewIntegrationExample() *IntegrationExample {
	config := DefaultConfig()
	topics := DefaultTopics()

	return &IntegrationExample{
		config: config,
		topics: topics,
	}
}

// Initialize sets up all Kafka components.
func (ie *IntegrationExample) Initialize(ctx context.Context) error {
	// Create producer
	ie.producer = NewProducer(ie.config, ie.topics)

	// Create analytics consumer
	ie.analyticsConsumer = NewAnalyticsConsumer(ie.config, ie.topics)

	// Create notification service
	ie.notificationService = NewNotificationService(ie.config, ie.topics)
	ie.notificationService.RegisterHandler(&PrintNotificationHandler{})
	ie.notificationService.RegisterHandler(&LogNotificationHandler{})

	// Create safe consumer with DLQ
	ie.safeConsumer = NewSafeConsumer(ie.config, ie.topics, []string{ie.topics.ProgressEvents})

	// Create DLQ handler
	ie.dlqHandler = NewDeadLetterQueueHandler(ie.config)

	// Register handler with DLQ support
	ie.safeConsumer.RegisterHandlerWithDLQ(
		ie.topics.ProgressEvents,
		func(ctx context.Context, msg *kafka.Message) error {
			// Simulate processing with 10% failure rate for demo
			// In production, this would be your actual business logic
			return nil
		})

	return nil
}

// StartConsumers starts all consumers.
func (ie *IntegrationExample) StartConsumers(ctx context.Context) error {
	// Start analytics consumer
	go func() {
		if err := ie.analyticsConsumer.Start(ctx); err != nil {
			log.Printf("Analytics consumer error: %v", err)
		}
	}()

	// Start notification service
	go func() {
		if err := ie.notificationService.Start(ctx); err != nil {
			log.Printf("Notification service error: %v", err)
		}
	}()

	return nil
}

// SimulateUserActivity simulates various user activities by publishing events.
func (ie *IntegrationExample) SimulateUserActivity(ctx context.Context, duration time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	eventCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("Simulation stopped after %d events", eventCount)
			return ctx.Err()

		case <-ticker.C:
			if time.Since(startTime) > duration {
				log.Printf("Simulation completed: %d events published", eventCount)
				return nil
			}

			// Publish a progress event
			userID := generateUserID()
			lessonID := generateLessonID()
			courseID := "course_101"
			score := generateScore()
			completed := score >= 70

			event := NewProgressEvent(
				ProgressUpdated,
				userID,
				lessonID,
				courseID,
				completed,
				score,
				generateTimeSpent(),
			)

			if err := ie.producer.PublishProgressEvent(ctx, event); err != nil {
				log.Printf("Error publishing progress event: %v", err)
				continue
			}

			eventCount++
			log.Printf("[%d] Published progress event: user=%s, lesson=%s, score=%d, completed=%v",
				eventCount, userID, lessonID, score, completed)

			// Also publish an exercise event occasionally
			if eventCount%3 == 0 {
				exerciseEvent := &ExerciseEvent{
					Event: Event{
						ID:        generateEventID(),
						Type:      ExerciseSubmitted,
						Source:    "simulation",
						Subject:   "exercise:" + lessonID,
						Data:      make(map[string]interface{}),
						Metadata:  make(map[string]string),
						Timestamp: time.Now().Unix(),
						Version:   "1.0",
					},
					ExerciseID: "exercise_" + lessonID,
					LessonID:   lessonID,
					UserID:     userID,
					Score:      score,
				}

				if err := ie.producer.PublishExerciseEvent(ctx, exerciseEvent); err != nil {
					log.Printf("Error publishing exercise event: %v", err)
				}
			}
		}
	}
}

// GetAnalyticsReport returns current analytics metrics.
func (ie *IntegrationExample) GetAnalyticsReport() map[string]interface{} {
	progressMetrics := ie.analyticsConsumer.GetMetrics()

	return map[string]interface{}{
		"total_lessons_completed": progressMetrics.TotalLessonsCompleted,
		"average_score":          progressMetrics.AverageScore,
		"total_time_spent_hours": float64(progressMetrics.TotalTimeSpent) / 3600,
		"active_users":           progressMetrics.ActiveUsers,
		"last_updated":           time.Unix(progressMetrics.LastUpdated, 0).Format(time.RFC3339),
	}
}

// GetCourseAnalytics returns analytics for a specific course.
func (ie *IntegrationExample) GetCourseAnalytics(courseID string) map[string]interface{} {
	courseMetrics := ie.analyticsConsumer.GetCourseMetrics(courseID)

	if courseMetrics == nil {
		return nil
	}

	return map[string]interface{}{
		"course_id":              courseMetrics.CourseID,
		"total_enrollments":      courseMetrics.TotalEnrollments,
		"completion_rate":        courseMetrics.CompletionRate,
		"average_progress":       courseMetrics.AverageProgress,
		"total_time_invested_hours": float64(courseMetrics.TotalTimeInvested) / 3600,
		"most_active_lesson":     courseMetrics.MostActiveLesson,
		"most_active_lesson_count": courseMetrics.MostActiveLessonCount,
		"last_updated":           time.Unix(courseMetrics.LastUpdated, 0).Format(time.RFC3339),
	}
}

// Close closes all resources.
func (ie *IntegrationExample) Close() error {
	var errs []error

	if err := ie.producer.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := ie.analyticsConsumer.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := ie.notificationService.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := ie.safeConsumer.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := ie.dlqHandler.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing integration: %v", errs)
	}

	return nil
}

// Helper functions for simulation

func generateUserID() string {
	users := []string{"user_001", "user_002", "user_003", "user_004", "user_005"}
	return users[int(time.Now().UnixNano())%len(users)]
}

func generateLessonID() string {
	lessons := []string{
		"lesson_001", "lesson_002", "lesson_003",
		"lesson_004", "lesson_005", "lesson_006",
	}
	return lessons[int(time.Now().UnixNano())%len(lessons)]
}

func generateScore() int {
	// Generate score with bias towards higher scores (realistic scenario)
	base := int(time.Now().UnixNano()%100) + 1
	if base < 30 {
		return 40 + int(time.Now().UnixNano()%30)
	}
	return base
}

func generateTimeSpent() int {
	// Time spent in seconds, between 5 minutes and 1 hour
	return 300 + int(time.Now().UnixNano()%3300)
}

func generateEventID() string {
	// Generate a unique event ID using timestamp and random component
	return fmt.Sprintf("evt_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
}

// MonitorMetrics continuously logs analytics metrics.
func (ie *IntegrationExample) MonitorMetrics(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			metrics := ie.GetAnalyticsReport()
			log.Printf("ANALYTICS SNAPSHOT: %+v", metrics)

			courseMetrics := ie.GetCourseAnalytics("course_101")
			if courseMetrics != nil {
				log.Printf("COURSE ANALYTICS: %+v", courseMetrics)
			}
		}
	}
}
