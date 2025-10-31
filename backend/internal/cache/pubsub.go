// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

// RedisPubSub implements pub/sub using Redis.
type RedisPubSub struct {
	client      *redis.Client
	prefix      string
	subscribers map[string]*redis.PubSub
	mu          sync.RWMutex
}

// NewRedisPubSub creates a new Redis pub/sub instance.
func NewRedisPubSub(client *redis.Client, prefix string) *RedisPubSub {
	if prefix == "" {
		prefix = "pubsub:"
	}

	return &RedisPubSub{
		client:      client,
		prefix:      prefix,
		subscribers: make(map[string]*redis.PubSub),
	}
}

// channelKey generates a channel key with prefix.
func (r *RedisPubSub) channelKey(channel string) string {
	return r.prefix + channel
}

// Publish publishes a message to a channel.
func (r *RedisPubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	channelKey := r.channelKey(channel)

	// Serialize message to JSON.
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message.
	result, err := r.client.Publish(ctx, channelKey, data).Result()
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// result is the number of subscribers that received the message.
	_ = result

	return nil
}

// Subscribe subscribes to one or more channels.
func (r *RedisPubSub) Subscribe(ctx context.Context, channels ...string) (<-chan Message, error) {
	if len(channels) == 0 {
		return nil, fmt.Errorf("no channels specified")
	}

	// Add prefix to channels.
	prefixedChannels := make([]string, len(channels))
	for i, channel := range channels {
		prefixedChannels[i] = r.channelKey(channel)
	}

	// Create subscription.
	pubsub := r.client.Subscribe(ctx, prefixedChannels...)

	// Store subscription for cleanup.
	subscriptionKey := fmt.Sprintf("sub_%p", pubsub)
	r.mu.Lock()
	r.subscribers[subscriptionKey] = pubsub
	r.mu.Unlock()

	// Create message channel.
	messageChan := make(chan Message, 100) // Buffer to prevent blocking

	// Start goroutine to handle messages.
	go func() {
		defer func() {
			close(messageChan)
			// Remove from subscribers map.
			r.mu.Lock()
			delete(r.subscribers, subscriptionKey)
			r.mu.Unlock()
		}()

		// Wait for subscription confirmation.
		_, err := pubsub.Receive(ctx)
		if err != nil {
			return
		}

		// Handle messages.
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}

				// Remove prefix from channel name.
				originalChannel := msg.Channel
				if len(originalChannel) > len(r.prefix) {
					originalChannel = originalChannel[len(r.prefix):]
				}

				// Create message.
				message := Message{
					Channel: originalChannel,
					Payload: msg.Payload,
				}

				// Send to message channel (non-blocking)
				select {
				case messageChan <- message:
				default:
					// Channel is full, drop message.
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return messageChan, nil
}

// Unsubscribe unsubscribes from channels.
func (r *RedisPubSub) Unsubscribe(ctx context.Context, channels ...string) error {
	// Add prefix to channels.
	prefixedChannels := make([]string, len(channels))
	for i, channel := range channels {
		prefixedChannels[i] = r.channelKey(channel)
	}

	// Find and unsubscribe from matching subscriptions.
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, pubsub := range r.subscribers {
		if err := pubsub.Unsubscribe(ctx, prefixedChannels...); err != nil {
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
	}

	return nil
}

// PSubscribe subscribes to channels matching patterns.
func (r *RedisPubSub) PSubscribe(ctx context.Context, patterns ...string) (<-chan Message, error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no patterns specified")
	}

	// Add prefix to patterns.
	prefixedPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		prefixedPatterns[i] = r.channelKey(pattern)
	}

	// Create pattern subscription.
	pubsub := r.client.PSubscribe(ctx, prefixedPatterns...)

	// Store subscription for cleanup.
	subscriptionKey := fmt.Sprintf("psub_%p", pubsub)
	r.mu.Lock()
	r.subscribers[subscriptionKey] = pubsub
	r.mu.Unlock()

	// Create message channel.
	messageChan := make(chan Message, 100)

	// Start goroutine to handle messages.
	go func() {
		defer func() {
			close(messageChan)
			// Remove from subscribers map.
			r.mu.Lock()
			delete(r.subscribers, subscriptionKey)
			r.mu.Unlock()
		}()

		// Wait for subscription confirmation.
		_, err := pubsub.Receive(ctx)
		if err != nil {
			return
		}

		// Handle messages.
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}

				// Remove prefix from channel and pattern names.
				originalChannel := msg.Channel
				if len(originalChannel) > len(r.prefix) {
					originalChannel = originalChannel[len(r.prefix):]
				}

				originalPattern := msg.Pattern
				if len(originalPattern) > len(r.prefix) {
					originalPattern = originalPattern[len(r.prefix):]
				}

				// Create message.
				message := Message{
					Channel: originalChannel,
					Pattern: originalPattern,
					Payload: msg.Payload,
				}

				// Send to message channel (non-blocking)
				select {
				case messageChan <- message:
				default:
					// Channel is full, drop message.
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return messageChan, nil
}

// PUnsubscribe unsubscribes from pattern subscriptions.
func (r *RedisPubSub) PUnsubscribe(ctx context.Context, patterns ...string) error {
	// Add prefix to patterns.
	prefixedPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		prefixedPatterns[i] = r.channelKey(pattern)
	}

	// Find and unsubscribe from matching pattern subscriptions.
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, pubsub := range r.subscribers {
		if err := pubsub.PUnsubscribe(ctx, prefixedPatterns...); err != nil {
			return fmt.Errorf("failed to pattern unsubscribe: %w", err)
		}
	}

	return nil
}

// Close closes all subscriptions.
func (r *RedisPubSub) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var lastErr error
	for _, pubsub := range r.subscribers {
		if err := pubsub.Close(); err != nil {
			lastErr = err
		}
	}

	// Clear subscribers map.
	r.subscribers = make(map[string]*redis.PubSub)

	return lastErr
}

// GetSubscriberCount returns the number of active subscribers.
func (r *RedisPubSub) GetSubscriberCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.subscribers)
}

// PublishJSON publishes a JSON message to a channel.
func (r *RedisPubSub) PublishJSON(ctx context.Context, channel string, message interface{}) error {
	return r.Publish(ctx, channel, message)
}

// PublishString publishes a string message to a channel.
func (r *RedisPubSub) PublishString(ctx context.Context, channel, message string) error {
	channelKey := r.channelKey(channel)
	_, err := r.client.Publish(ctx, channelKey, message).Result()
	if err != nil {
		return fmt.Errorf("failed to publish string message: %w", err)
	}

	return nil
}

// PublishBytes publishes a byte array message to a channel.
func (r *RedisPubSub) PublishBytes(ctx context.Context, channel string, message []byte) error {
	channelKey := r.channelKey(channel)
	_, err := r.client.Publish(ctx, channelKey, message).Result()
	if err != nil {
		return fmt.Errorf("failed to publish bytes message: %w", err)
	}

	return nil
}

// GetChannelSubscribers returns the number of subscribers for a channel.
func (r *RedisPubSub) GetChannelSubscribers(ctx context.Context, channel string) (int64, error) {
	channelKey := r.channelKey(channel)

	// Use PUBSUB NUMSUB command.
	result, err := r.client.PubSubNumSub(ctx, channelKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get channel subscribers: %w", err)
	}

	if count, ok := result[channelKey]; ok {
		return count, nil
	}

	return 0, nil
}

// ListChannels returns all active channels matching a pattern.
func (r *RedisPubSub) ListChannels(ctx context.Context, pattern string) ([]string, error) {
	prefixedPattern := r.channelKey(pattern)

	channels, err := r.client.PubSubChannels(ctx, prefixedPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	// Remove prefix from channel names.
	result := make([]string, len(channels))
	for i, channel := range channels {
		if len(channel) > len(r.prefix) {
			result[i] = channel[len(r.prefix):]
		} else {
			result[i] = channel
		}
	}

	return result, nil
}

// BulkPublish publishes multiple messages to different channels.
func (r *RedisPubSub) BulkPublish(ctx context.Context, messages map[string]interface{}) error {
	if len(messages) == 0 {
		return nil
	}

	// Use pipeline for bulk publishing.
	pipe := r.client.Pipeline()

	for channel, message := range messages {
		channelKey := r.channelKey(channel)

		// Serialize message.
		data, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message for channel %s: %w", channel, err)
		}

		pipe.Publish(ctx, channelKey, data)
	}

	// Execute pipeline.
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to bulk publish messages: %w", err)
	}

	return nil
}

// SubscribeWithHandler subscribes to a channel and handles messages with a callback.
func (r *RedisPubSub) SubscribeWithHandler(ctx context.Context, channel string, handler func(Message) error) error {
	messageChan, err := r.Subscribe(ctx, channel)
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	// Handle messages in a goroutine.
	go func() {
		for {
			select {
			case message, ok := <-messageChan:
				if !ok {
					return
				}

				// Call handler.
				if err := handler(message); err != nil {
					// Log error (in production, use proper logging)
					fmt.Printf("Message handler error for channel %s: %v\n", channel, err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
