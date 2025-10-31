// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package kafka provides functionality for the GO-PRO Learning Platform.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageHandler defines the interface for handling Kafka messages.
type MessageHandler interface {
	HandleMessage(ctx context.Context, message *kafka.Message) error
}

// MessageHandlerFunc is a function adapter for MessageHandler.
type MessageHandlerFunc func(ctx context.Context, message *kafka.Message) error

// HandleMessage implements MessageHandler.
func (f MessageHandlerFunc) HandleMessage(ctx context.Context, message *kafka.Message) error {
	return f(ctx, message)
}

// Consumer wraps kafka-go reader with additional functionality.
type Consumer struct {
	reader   *kafka.Reader
	config   *Config
	topics   *Topics
	handlers map[string]MessageHandler
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(config *Config, topics *Topics, consumerTopics []string) *Consumer {
	if config == nil {
		config = DefaultConfig()
	}
	if topics == nil {
		topics = DefaultTopics()
	}

	// For kafka-go, we need to create separate readers for each topic.
	// or use the first topic if only one is specified.
	var topic string
	if len(consumerTopics) > 0 {
		topic = consumerTopics[0]
	} else {
		topic = "default"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          topic,
		GroupID:        config.GroupID,
		MinBytes:       config.FetchMinBytes,
		MaxBytes:       10e6, // 10MB
		MaxWait:        config.FetchMaxWait,
		CommitInterval: config.AutoCommitInterval,
		StartOffset:    kafka.LastOffset,
	})

	return &Consumer{
		reader:   reader,
		config:   config,
		topics:   topics,
		handlers: make(map[string]MessageHandler),
	}
}

// RegisterHandler registers a message handler for a specific topic.
func (c *Consumer) RegisterHandler(topic string, handler MessageHandler) {
	c.handlers[topic] = handler
}

// RegisterHandlerFunc registers a message handler function for a specific topic.
func (c *Consumer) RegisterHandlerFunc(topic string, handlerFunc MessageHandlerFunc) {
	c.handlers[topic] = handlerFunc
}

// Start starts consuming messages.
func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			if err := c.processMessage(ctx, &message); err != nil {
				log.Printf("Error processing message: %v", err)
				// In production, you might want to send to a dead letter queue.
				continue
			}
		}
	}
}

// Close closes the consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// processMessage processes a single message.
func (c *Consumer) processMessage(ctx context.Context, message *kafka.Message) error {
	handler, exists := c.handlers[message.Topic]
	if !exists {
		log.Printf("No handler registered for topic: %s", message.Topic)
		return nil
	}

	return handler.HandleMessage(ctx, message)
}

// EventConsumer provides high-level event consumption.
type EventConsumer struct {
	consumer *Consumer
}

// NewEventConsumer creates a new event consumer.
func NewEventConsumer(config *Config, topics *Topics) *EventConsumer {
	allTopics := []string{
		topics.UserEvents,
		topics.CourseEvents,
		topics.LessonEvents,
		topics.ExerciseEvents,
		topics.ProgressEvents,
		topics.NotificationEvents,
		topics.AuditEvents,
	}

	consumer := NewConsumer(config, topics, allTopics)

	return &EventConsumer{consumer: consumer}
}

// RegisterUserEventHandler registers a handler for user events.
func (ec *EventConsumer) RegisterUserEventHandler(handler func(ctx context.Context, event *UserEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.UserEvents, func(ctx context.Context, message *kafka.Message) error {
		var event UserEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterCourseEventHandler registers a handler for course events.
func (ec *EventConsumer) RegisterCourseEventHandler(handler func(ctx context.Context, event *CourseEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.CourseEvents, func(ctx context.Context, message *kafka.Message) error {
		var event CourseEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal course event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterLessonEventHandler registers a handler for lesson events.
func (ec *EventConsumer) RegisterLessonEventHandler(handler func(ctx context.Context, event *LessonEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.LessonEvents, func(ctx context.Context, message *kafka.Message) error {
		var event LessonEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal lesson event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterExerciseEventHandler registers a handler for exercise events.
func (ec *EventConsumer) RegisterExerciseEventHandler(handler func(ctx context.Context, event *ExerciseEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.ExerciseEvents, func(ctx context.Context, message *kafka.Message) error {
		var event ExerciseEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal exercise event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterProgressEventHandler registers a handler for progress events.
func (ec *EventConsumer) RegisterProgressEventHandler(handler func(ctx context.Context, event *ProgressEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.ProgressEvents, func(ctx context.Context, message *kafka.Message) error {
		var event ProgressEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal progress event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterNotificationEventHandler registers a handler for notification events.
func (ec *EventConsumer) RegisterNotificationEventHandler(handler func(ctx context.Context, event *NotificationEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.NotificationEvents, func(ctx context.Context, message *kafka.Message) error {
		var event NotificationEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal notification event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// RegisterAuditEventHandler registers a handler for audit events.
func (ec *EventConsumer) RegisterAuditEventHandler(handler func(ctx context.Context, event *AuditEvent) error) {
	ec.consumer.RegisterHandlerFunc(ec.consumer.topics.AuditEvents, func(ctx context.Context, message *kafka.Message) error {
		var event AuditEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal audit event: %w", err)
		}

		return handler(ctx, &event)
	})
}

// Start starts the event consumer.
func (ec *EventConsumer) Start(ctx context.Context) error {
	return ec.consumer.Start(ctx)
}

// Close closes the event consumer.
func (ec *EventConsumer) Close() error {
	return ec.consumer.Close()
}

// BatchConsumer provides batch message consumption.
type BatchConsumer struct {
	reader    *kafka.Reader
	batchSize int
	timeout   time.Duration
}

// NewBatchConsumer creates a new batch consumer.
func NewBatchConsumer(config *Config, topics []string, batchSize int, timeout time.Duration) *BatchConsumer {
	if config == nil {
		config = DefaultConfig()
	}

	// For kafka-go, use the first topic or default.
	var topic string
	if len(topics) > 0 {
		topic = topics[0]
	} else {
		topic = "default"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          topic,
		GroupID:        config.GroupID,
		MinBytes:       config.FetchMinBytes,
		MaxBytes:       10e6, // 10MB
		MaxWait:        config.FetchMaxWait,
		CommitInterval: config.AutoCommitInterval,
		StartOffset:    kafka.LastOffset,
	})

	return &BatchConsumer{
		reader:    reader,
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// ConsumeBatch consumes messages in batches.
func (bc *BatchConsumer) ConsumeBatch(ctx context.Context, handler func([]kafka.Message) error) error {
	messages := make([]kafka.Message, 0, bc.batchSize)
	timer := time.NewTimer(bc.timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(messages) > 0 {
				return handler(messages)
			}

			return ctx.Err()
		case <-timer.C:
			if len(messages) > 0 {
				if err := handler(messages); err != nil {
					return err
				}
				messages = messages[:0]
			}
			timer.Reset(bc.timeout)
		default:
			message, err := bc.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			messages = append(messages, message)
			if len(messages) >= bc.batchSize {
				if err := handler(messages); err != nil {
					return err
				}
				messages = messages[:0]
				timer.Reset(bc.timeout)
			}
		}
	}
}

// Close closes the batch consumer.
func (bc *BatchConsumer) Close() error {
	return bc.reader.Close()
}
