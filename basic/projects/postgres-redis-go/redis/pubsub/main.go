package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                    Redis Pub/Sub - Tutorial                                  ║
║                                                                              ║
║  Redis Pub/Sub enables message broadcasting:                                ║
║  • Publishers send messages to channels                                     ║
║  • Subscribers receive messages from channels                               ║
║  • Pattern-based subscriptions                                              ║
║  • Real-time messaging                                                       ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

func main() {
	fmt.Println("📡 Redis Pub/Sub Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	ctx := context.Background()

	// Example 1: Simple Pub/Sub
	fmt.Println("\n📌 Example 1: Simple Pub/Sub")
	simplePubSub(ctx)

	// Example 2: Pattern Subscription
	fmt.Println("\n📌 Example 2: Pattern Subscription")
	patternSubscription(ctx)

	// Example 3: Multiple Channels
	fmt.Println("\n📌 Example 3: Multiple Channels")
	multipleChannels(ctx)

	fmt.Println("\n✅ All Pub/Sub examples completed!")
}

func simplePubSub(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Subscribe to channel
	pubsub := rdb.Subscribe(ctx, "notifications")
	defer pubsub.Close()

	// Wait for subscription confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("❌ Subscribe failed: %v\n", err)
		return
	}
	fmt.Println("✅ Subscribed to 'notifications' channel")

	// Start goroutine to receive messages
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("📨 Received: %s (channel: %s)\n", msg.Payload, msg.Channel)
		}
	}()

	// Publish messages
	publisher := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer publisher.Close()

	messages := []string{
		"Hello, World!",
		"Welcome to Redis Pub/Sub",
		"This is a test message",
	}

	for _, msg := range messages {
		err := publisher.Publish(ctx, "notifications", msg).Err()
		if err != nil {
			log.Printf("❌ Publish failed: %v\n", err)
			continue
		}
		fmt.Printf("📤 Published: %s\n", msg)
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(1 * time.Second) // Wait for messages to be received
}

func patternSubscription(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// Subscribe to pattern
	pubsub := rdb.PSubscribe(ctx, "user:*:events")
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("❌ Pattern subscribe failed: %v\n", err)
		return
	}
	fmt.Println("✅ Subscribed to pattern 'user:*:events'")

	// Receive messages
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("📨 Pattern match: %s (channel: %s, pattern: %s)\n",
				msg.Payload, msg.Channel, msg.Pattern)
		}
	}()

	// Publish to matching channels
	publisher := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer publisher.Close()

	channels := []string{
		"user:123:events",
		"user:456:events",
		"user:789:events",
	}

	for _, channel := range channels {
		err := publisher.Publish(ctx, channel, fmt.Sprintf("Event on %s", channel)).Err()
		if err != nil {
			log.Printf("❌ Publish failed: %v\n", err)
			continue
		}
		fmt.Printf("📤 Published to: %s\n", channel)
		time.Sleep(300 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
}

func multipleChannels(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// Subscribe to multiple channels
	pubsub := rdb.Subscribe(ctx, "orders", "payments", "shipments")
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("❌ Subscribe failed: %v\n", err)
		return
	}
	fmt.Println("✅ Subscribed to multiple channels: orders, payments, shipments")

	// Receive messages
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("📨 [%s] %s\n", msg.Channel, msg.Payload)
		}
	}()

	// Publish to different channels
	publisher := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer publisher.Close()

	events := map[string]string{
		"orders":    "New order #12345",
		"payments":  "Payment received $99.99",
		"shipments": "Package shipped",
	}

	for channel, message := range events {
		err := publisher.Publish(ctx, channel, message).Err()
		if err != nil {
			log.Printf("❌ Publish failed: %v\n", err)
			continue
		}
		fmt.Printf("📤 [%s] %s\n", channel, message)
		time.Sleep(300 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
}

