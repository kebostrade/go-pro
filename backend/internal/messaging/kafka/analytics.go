// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ProgressMetrics tracks real-time progress statistics.
type ProgressMetrics struct {
	TotalLessonsCompleted int     `json:"total_lessons_completed"`
	AverageScore         float64 `json:"average_score"`
	TotalTimeSpent       int     `json:"total_time_spent_seconds"`
	LastUpdated          int64   `json:"last_updated"`
	ActiveUsers          int     `json:"active_users"`
}

// CourseMetrics tracks real-time course statistics.
type CourseMetrics struct {
	CourseID              string  `json:"course_id"`
	TotalEnrollments      int     `json:"total_enrollments"`
	CompletionRate        float64 `json:"completion_rate"`
	AverageProgress       float64 `json:"average_progress"`
	TotalTimeInvested     int     `json:"total_time_invested_seconds"`
	LastUpdated           int64   `json:"last_updated"`
	MostActiveLesson      string  `json:"most_active_lesson"`
	MostActiveLessonCount int     `json:"most_active_lesson_count"`
}

// AnalyticsAggregator aggregates real-time metrics from Kafka streams.
type AnalyticsAggregator struct {
	mu                    sync.RWMutex
	progressMetrics       *ProgressMetrics
	courseMetrics         map[string]*CourseMetrics
	userProgressMap       map[string]*ProgressEvent
	lessonActivityMap     map[string]int
	metricsUpdateChan     chan interface{}
	ctx                   context.Context
	cancel                context.CancelFunc
}

// NewAnalyticsAggregator creates a new analytics aggregator.
func NewAnalyticsAggregator() *AnalyticsAggregator {
	ctx, cancel := context.WithCancel(context.Background())

	return &AnalyticsAggregator{
		progressMetrics:   &ProgressMetrics{},
		courseMetrics:     make(map[string]*CourseMetrics),
		userProgressMap:   make(map[string]*ProgressEvent),
		lessonActivityMap: make(map[string]int),
		metricsUpdateChan: make(chan interface{}, 100),
		ctx:               ctx,
		cancel:            cancel,
	}
}

// ProcessProgressEvent processes a progress event and updates metrics.
func (aa *AnalyticsAggregator) ProcessProgressEvent(event *ProgressEvent) {
	aa.mu.Lock()
	defer aa.mu.Unlock()

	// Update user progress map
	userProgressKey := fmt.Sprintf("%s:%s", event.UserID, event.LessonID)
	aa.userProgressMap[userProgressKey] = event

	// Update lesson activity
	aa.lessonActivityMap[event.LessonID]++

	// Update overall metrics
	if event.Completed {
		aa.progressMetrics.TotalLessonsCompleted++
	}

	aa.progressMetrics.TotalTimeSpent += event.TimeSpent
	aa.progressMetrics.LastUpdated = time.Now().Unix()
	aa.progressMetrics.ActiveUsers = len(aa.getUserSet())

	// Recalculate average score
	aa.recalculateAverageScore()

	// Update course metrics
	aa.updateCourseMetrics(event)

	// Send update notification
	select {
	case aa.metricsUpdateChan <- event:
	default:
		log.Println("Metrics update channel full, dropping update")
	}
}

// ProcessUserEvent processes user events for analytics.
func (aa *AnalyticsAggregator) ProcessUserEvent(event *UserEvent) {
	aa.mu.Lock()
	defer aa.mu.Unlock()

	aa.progressMetrics.ActiveUsers = len(aa.getUserSet())
	aa.progressMetrics.LastUpdated = time.Now().Unix()

	select {
	case aa.metricsUpdateChan <- event:
	default:
		log.Println("Metrics update channel full, dropping update")
	}
}

// GetProgressMetrics returns current progress metrics.
func (aa *AnalyticsAggregator) GetProgressMetrics() *ProgressMetrics {
	aa.mu.RLock()
	defer aa.mu.RUnlock()

	return &ProgressMetrics{
		TotalLessonsCompleted: aa.progressMetrics.TotalLessonsCompleted,
		AverageScore:         aa.progressMetrics.AverageScore,
		TotalTimeSpent:       aa.progressMetrics.TotalTimeSpent,
		LastUpdated:          aa.progressMetrics.LastUpdated,
		ActiveUsers:          aa.progressMetrics.ActiveUsers,
	}
}

// GetCourseMetrics returns current course metrics.
func (aa *AnalyticsAggregator) GetCourseMetrics(courseID string) *CourseMetrics {
	aa.mu.RLock()
	defer aa.mu.RUnlock()

	if metrics, ok := aa.courseMetrics[courseID]; ok {
		return &CourseMetrics{
			CourseID:              metrics.CourseID,
			TotalEnrollments:      metrics.TotalEnrollments,
			CompletionRate:        metrics.CompletionRate,
			AverageProgress:       metrics.AverageProgress,
			TotalTimeInvested:     metrics.TotalTimeInvested,
			LastUpdated:           metrics.LastUpdated,
			MostActiveLesson:      metrics.MostActiveLesson,
			MostActiveLessonCount: metrics.MostActiveLessonCount,
		}
	}

	return nil
}

// GetAllCourseMetrics returns all course metrics.
func (aa *AnalyticsAggregator) GetAllCourseMetrics() map[string]*CourseMetrics {
	aa.mu.RLock()
	defer aa.mu.RUnlock()

	result := make(map[string]*CourseMetrics)
	for k, v := range aa.courseMetrics {
		result[k] = v
	}

	return result
}

// MetricsUpdatesChan returns a channel for metrics updates.
func (aa *AnalyticsAggregator) MetricsUpdatesChan() <-chan interface{} {
	return aa.metricsUpdateChan
}

// Close closes the aggregator.
func (aa *AnalyticsAggregator) Close() error {
	aa.cancel()
	close(aa.metricsUpdateChan)
	return nil
}

// Private helper methods

func (aa *AnalyticsAggregator) getUserSet() map[string]bool {
	userSet := make(map[string]bool)
	for key := range aa.userProgressMap {
		// Parse "userID:lessonID" format
		parts := make([]byte, len(key))
		copy(parts, key)
		for i, b := range parts {
			if b == ':' {
				userSet[string(parts[:i])] = true
				break
			}
		}
	}
	return userSet
}

func (aa *AnalyticsAggregator) recalculateAverageScore() {
	if len(aa.userProgressMap) == 0 {
		aa.progressMetrics.AverageScore = 0
		return
	}

	total := 0
	for _, progress := range aa.userProgressMap {
		total += progress.Score
	}

	aa.progressMetrics.AverageScore = float64(total) / float64(len(aa.userProgressMap))
}

func (aa *AnalyticsAggregator) updateCourseMetrics(event *ProgressEvent) {
	if event.CourseID == "" {
		return
	}

	if _, ok := aa.courseMetrics[event.CourseID]; !ok {
		aa.courseMetrics[event.CourseID] = &CourseMetrics{
			CourseID: event.CourseID,
		}
	}

	metrics := aa.courseMetrics[event.CourseID]

	// Update course time invested
	metrics.TotalTimeInvested += event.TimeSpent

	// Update most active lesson
	if aa.lessonActivityMap[event.LessonID] > metrics.MostActiveLessonCount {
		metrics.MostActiveLesson = event.LessonID
		metrics.MostActiveLessonCount = aa.lessonActivityMap[event.LessonID]
	}

	metrics.LastUpdated = time.Now().Unix()

	// Recalculate completion rate and average progress
	aa.recalculateCourseMetrics(metrics)
}

func (aa *AnalyticsAggregator) recalculateCourseMetrics(metrics *CourseMetrics) {
	// Count unique users and completed lessons for this course
	userCount := 0
	completedCount := 0
	totalScore := 0

	userSet := make(map[string]bool)
	for key, progress := range aa.userProgressMap {
		if progress.CourseID != metrics.CourseID {
			continue
		}

		parts := make([]byte, len(key))
		copy(parts, key)
		var userID string
		for i, b := range parts {
			if b == ':' {
				userID = string(parts[:i])
				break
			}
		}

		if !userSet[userID] {
			userSet[userID] = true
			userCount++
		}

		if progress.Completed {
			completedCount++
		}

		totalScore += progress.Score
	}

	metrics.TotalEnrollments = userCount
	if userCount > 0 {
		metrics.CompletionRate = float64(completedCount) / float64(userCount)
		metrics.AverageProgress = float64(totalScore) / float64(userCount)
	}
}

// AnalyticsConsumer consumes progress events and feeds them to the aggregator.
type AnalyticsConsumer struct {
	eventConsumer *EventConsumer
	aggregator    *AnalyticsAggregator
}

// NewAnalyticsConsumer creates a new analytics consumer.
func NewAnalyticsConsumer(config *Config, topics *Topics) *AnalyticsConsumer {
	eventConsumer := NewEventConsumer(config, topics)
	aggregator := NewAnalyticsAggregator()

	return &AnalyticsConsumer{
		eventConsumer: eventConsumer,
		aggregator:    aggregator,
	}
}

// Start starts consuming events for analytics.
func (ac *AnalyticsConsumer) Start(ctx context.Context) error {
	// Register handlers
	ac.eventConsumer.RegisterProgressEventHandler(func(ctx context.Context, event *ProgressEvent) error {
		ac.aggregator.ProcessProgressEvent(event)
		log.Printf("Analytics: Processed progress event for user %s, lesson %s", event.UserID, event.LessonID)
		return nil
	})

	ac.eventConsumer.RegisterUserEventHandler(func(ctx context.Context, event *UserEvent) error {
		ac.aggregator.ProcessUserEvent(event)
		return nil
	})

	return ac.eventConsumer.Start(ctx)
}

// GetMetrics returns current metrics.
func (ac *AnalyticsConsumer) GetMetrics() *ProgressMetrics {
	return ac.aggregator.GetProgressMetrics()
}

// GetCourseMetrics returns course-specific metrics.
func (ac *AnalyticsConsumer) GetCourseMetrics(courseID string) *CourseMetrics {
	return ac.aggregator.GetCourseMetrics(courseID)
}

// Close closes the consumer.
func (ac *AnalyticsConsumer) Close() error {
	var errs []error
	if err := ac.aggregator.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := ac.eventConsumer.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing analytics consumer: %v", errs)
	}

	return nil
}
