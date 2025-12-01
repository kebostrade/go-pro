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

// DeadLetterQueueHandler handles messages that fail processing.
type DeadLetterQueueHandler struct {
	writer      *kafka.Writer
	topic       string
	maxRetries  int
	retryDelay  time.Duration
}

// DeadLetterMessage represents a message in the DLQ.
type DeadLetterMessage struct {
	OriginalTopic   string            `json:"original_topic"`
	OriginalMessage []byte            `json:"original_message"`
	OriginalHeaders map[string]string `json:"original_headers"`
	ErrorMessage    string            `json:"error_message"`
	FailureCount    int               `json:"failure_count"`
	FirstFailedAt   int64             `json:"first_failed_at"`
	LastFailedAt    int64             `json:"last_failed_at"`
	StackTrace      string            `json:"stack_trace,omitempty"`
}

// NewDeadLetterQueueHandler creates a new DLQ handler.
func NewDeadLetterQueueHandler(config *Config) *DeadLetterQueueHandler {
	dlqTopic := "go-pro.dlq"

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  dlqTopic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		ReadTimeout:            config.RequestTimeout,
		WriteTimeout:           config.RequestTimeout,
		AllowAutoTopicCreation: true,
	}

	return &DeadLetterQueueHandler{
		writer:     writer,
		topic:      dlqTopic,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryBackoff,
	}
}

// HandleFailedMessage sends a failed message to the DLQ.
func (dlq *DeadLetterQueueHandler) HandleFailedMessage(
	ctx context.Context,
	originalTopic string,
	message *kafka.Message,
	err error,
	failureCount int,
) error {
	// Create headers map from original message
	headers := make(map[string]string)
	for _, h := range message.Headers {
		headers[h.Key] = string(h.Value)
	}

	dlqMsg := &DeadLetterMessage{
		OriginalTopic:   originalTopic,
		OriginalMessage: message.Value,
		OriginalHeaders: headers,
		ErrorMessage:    err.Error(),
		FailureCount:    failureCount,
		FirstFailedAt:   time.Now().Unix(),
		LastFailedAt:    time.Now().Unix(),
	}

	payload, err := json.Marshal(dlqMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Topic: dlq.topic,
		Key:   message.Key,
		Value: payload,
		Headers: []kafka.Header{
			{Key: "original-topic", Value: []byte(originalTopic)},
			{Key: "error", Value: []byte(err.Error())},
			{Key: "failed-at", Value: []byte(fmt.Sprintf("%d", time.Now().Unix()))},
			{Key: "failure-count", Value: []byte(fmt.Sprintf("%d", failureCount))},
		},
	}

	return dlq.writer.WriteMessages(ctx, kafkaMsg)
}

// RetryFailedMessage retries a message from the DLQ.
func (dlq *DeadLetterQueueHandler) RetryFailedMessage(
	ctx context.Context,
	dlqMessage *DeadLetterMessage,
	producer *Producer,
) error {
	log.Printf("Retrying failed message from DLQ (failure count: %d, original topic: %s)",
		dlqMessage.FailureCount, dlqMessage.OriginalTopic)

	if dlqMessage.FailureCount >= dlq.maxRetries {
		log.Printf("Message has exceeded max retries (%d), keeping in DLQ", dlq.maxRetries)
		return fmt.Errorf("message exceeded max retries")
	}

	return nil
}

// Close closes the DLQ handler.
func (dlq *DeadLetterQueueHandler) Close() error {
	return dlq.writer.Close()
}

// RetryConsumer consumes messages from DLQ and retries them.
type RetryConsumer struct {
	reader  *kafka.Reader
	dlqHandler *DeadLetterQueueHandler
	producer *Producer
}

// NewRetryConsumer creates a new retry consumer.
func NewRetryConsumer(config *Config, dlqHandler *DeadLetterQueueHandler, producer *Producer) *RetryConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          "go-pro.dlq",
		GroupID:        "go-pro-dlq-consumer",
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        500 * time.Millisecond,
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &RetryConsumer{
		reader:     reader,
		dlqHandler: dlqHandler,
		producer:   producer,
	}
}

// Start starts consuming and retrying messages from DLQ.
func (rc *RetryConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := rc.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading DLQ message: %v", err)
				continue
			}

			var dlqMsg DeadLetterMessage
			if err := json.Unmarshal(message.Value, &dlqMsg); err != nil {
				log.Printf("Error unmarshaling DLQ message: %v", err)
				continue
			}

			// Attempt retry
			if err := rc.dlqHandler.RetryFailedMessage(ctx, &dlqMsg, rc.producer); err != nil {
				log.Printf("Error retrying message: %v", err)
			}
		}
	}
}

// Close closes the retry consumer.
func (rc *RetryConsumer) Close() error {
	return rc.reader.Close()
}

// SafeConsumer wraps a consumer with automatic DLQ handling on failures.
type SafeConsumer struct {
	consumer *Consumer
	dlqHandler *DeadLetterQueueHandler
	failureThreshold int
}

// NewSafeConsumer creates a new safe consumer with DLQ support.
func NewSafeConsumer(config *Config, topics *Topics, consumerTopics []string) *SafeConsumer {
	consumer := NewConsumer(config, topics, consumerTopics)
	dlqHandler := NewDeadLetterQueueHandler(config)

	return &SafeConsumer{
		consumer: consumer,
		dlqHandler: dlqHandler,
		failureThreshold: 3,
	}
}

// ProcessMessageSafely processes a message with automatic retry and DLQ handling.
func (sc *SafeConsumer) ProcessMessageSafely(
	ctx context.Context,
	originalTopic string,
	message *kafka.Message,
	handler func(context.Context, *kafka.Message) error,
) error {
	failureCount := 0
	retryDelay := 100 * time.Millisecond

	for failureCount < sc.failureThreshold {
		if err := handler(ctx, message); err == nil {
			return nil // Success
		} else {
			failureCount++
			if failureCount < sc.failureThreshold {
				log.Printf("Message processing failed (attempt %d/%d), retrying in %v: %v",
					failureCount, sc.failureThreshold, retryDelay, err)
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
			} else {
				// Max retries exceeded, send to DLQ
				log.Printf("Message processing failed after %d attempts, sending to DLQ", failureCount)
				if dlqErr := sc.dlqHandler.HandleFailedMessage(ctx, originalTopic, message, err, failureCount); dlqErr != nil {
					log.Printf("Error sending message to DLQ: %v", dlqErr)
					return fmt.Errorf("failed to process message and send to DLQ: %v", dlqErr)
				}
				return err
			}
		}
	}

	return fmt.Errorf("message processing failed after %d attempts", sc.failureThreshold)
}

// RegisterHandlerWithDLQ registers a handler with automatic DLQ support.
func (sc *SafeConsumer) RegisterHandlerWithDLQ(topic string, handler func(context.Context, *kafka.Message) error) {
	sc.consumer.RegisterHandlerFunc(topic, func(ctx context.Context, message *kafka.Message) error {
		return sc.ProcessMessageSafely(ctx, topic, message, handler)
	})
}

// Close closes the safe consumer.
func (sc *SafeConsumer) Close() error {
	var errs []error
	if err := sc.consumer.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := sc.dlqHandler.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing safe consumer: %v", errs)
	}

	return nil
}
