# Kafka Real-Time Data Processing Guide

## Overview

This guide explains how to use Kafka for real-time data processing in the GO-PRO Learning Platform. The system includes:

- **Real-Time Analytics**: Stream aggregation of progress events
- **Notifications**: Automatic notifications based on learner actions
- **Dead Letter Queue (DLQ)**: Failed message handling and retry logic
- **WebSocket Integration**: Live dashboard updates
- **Docker Ecosystem**: Complete Kafka stack with monitoring

## Quick Start

### 1. Start the Kafka Ecosystem

```bash
# Start all services (Kafka, Redis, PostgreSQL, Prometheus, Grafana, etc.)
docker-compose -f docker-compose.kafka.yml up -d

# Verify all services are running
docker-compose -f docker-compose.kafka.yml ps

# View logs
docker-compose -f docker-compose.kafka.yml logs -f kafka
```

**Service URLs:**
- Kafka: `localhost:9092`
- Kafka UI: `http://localhost:8080`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)
- PostgreSQL: `localhost:5432` (gopro/gopro_password)
- Redis: `localhost:6379`
- Adminer: `http://localhost:8081`
- Qdrant (Vector DB): `http://localhost:6333`

### 2. Run the Backend with Kafka

```bash
cd backend

# Set environment variables
export KAFKA_BROKERS=localhost:9092
export KAFKA_GROUP_ID=go-pro-consumer-group
export DATABASE_URL=postgres://gopro:gopro_password@localhost:5432/gopro_db
export REDIS_URL=redis://localhost:6379
export MESSAGING_ENABLED=true

# Run the server
go run ./cmd/server
```

## Architecture

### Kafka Topics

Topics are automatically created and follow this naming convention:

```
go-pro.user.events          # User login, registration, updates
go-pro.course.events        # Course creation, publication
go-pro.lesson.events        # Lesson completion, updates
go-pro.exercise.events      # Exercise submission, grading
go-pro.progress.events      # Progress tracking
go-pro.notification.events  # User notifications
go-pro.audit.events         # Audit logs
go-pro.dlq                  # Dead letter queue for failed messages
```

### Event Flow

```
┌─────────────────┐
│  Backend API    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Event Producer (Kafka Writer)  │
└────────┬────────────────────────┘
         │
         ▼
    ┌──────────────────────────────────────────────────────────┐
    │                    Kafka Topics                          │
    │ ┌──────────────┬──────────────┬──────────────────────┐  │
    │ │ user.events  │ course.events│ progress.events      │  │
    │ └──────────────┴──────────────┴──────────────────────┘  │
    └──────────────────────────────────────────────────────────┘
         │
         ├──────────────────┬────────────────────┬──────────────┐
         ▼                  ▼                    ▼              ▼
    ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────┐
    │  Analytics  │  │Notifications │  │Audit Logger  │  │Dashboard │
    │  Aggregator │  │  Service     │  │   Handler    │  │Broadcast │
    └─────────────┘  └──────────────┘  └──────────────┘  └──────────┘
         │                  │
         ▼                  ▼
    ┌──────────────────────────────┐
    │  Real-Time Metrics & Alerts  │
    └──────────────────────────────┘
```

## Core Components

### 1. Analytics Aggregator

Real-time aggregation of progress metrics.

```go
import "go-pro-backend/internal/messaging/kafka"

// Create aggregator
aggregator := kafka.NewAnalyticsAggregator()

// Process progress events
aggregator.ProcessProgressEvent(&event)

// Get current metrics
metrics := aggregator.GetProgressMetrics()
fmt.Printf("Active Users: %d\n", metrics.ActiveUsers)
fmt.Printf("Avg Score: %.2f\n", metrics.AverageScore)
fmt.Printf("Total Time: %d seconds\n", metrics.TotalTimeSpent)
```

**Tracked Metrics:**
- `TotalLessonsCompleted`: Count of completed lessons
- `AverageScore`: Average score across all users
- `TotalTimeSpent`: Aggregate time spent learning
- `ActiveUsers`: Currently active users
- `CompletionRate`: % of lessons completed per course
- `MostActiveLesson`: Most engaged lesson

### 2. Notification Service

Automatic notifications based on learner actions.

```go
// Create notification service
notificationService := kafka.NewNotificationService(config, topics)

// Register handlers
notificationService.RegisterHandler(&kafka.PrintNotificationHandler{})
notificationService.RegisterHandler(&kafka.LogNotificationHandler{})

// Set progress thresholds for milestones
notificationService.SetProgressThreshold("lesson_1", 80)

// Start consuming and processing events
go notificationService.Start(ctx)

// Notifications are sent automatically for:
// - Lesson completion
// - High scores (>=80)
// - Exercise submission/grading
// - Progress milestones
```

**Notification Types:**
- `progress`: Progress updates
- `exercise`: Exercise submission results
- `lesson`: Lesson completion
- `milestone`: Achievement unlocked
- `alert`: System alerts

### 3. Dead Letter Queue (DLQ)

Handles failed message processing with automatic retry.

```go
// Create DLQ handler
dlqHandler := kafka.NewDeadLetterQueueHandler(config)

// Create safe consumer with DLQ support
safeConsumer := kafka.NewSafeConsumer(config, topics, consumerTopics)

// Register handler with automatic DLQ on failure
safeConsumer.RegisterHandlerWithDLQ(topic, func(ctx context.Context, msg *kafka.Message) error {
    // Your processing logic
    // If it returns an error, message is automatically sent to DLQ
    return processMessage(msg)
})

// Retry messages from DLQ
retryConsumer := kafka.NewRetryConsumer(config, dlqHandler, producer)
go retryConsumer.Start(ctx)
```

**DLQ Features:**
- Automatic message routing on failures
- Exponential backoff retry strategy (100ms → 200ms → 400ms)
- Failure tracking and history
- Manual review capability via Kafka UI

### 4. WebSocket Real-Time Updates

Live dashboard and notifications over WebSocket.

```go
import "go-pro-backend/internal/messaging/realtime"

// Create WebSocket hub
hub := realtime.NewHub()
go hub.Start(ctx)

// Create dashboard adapter
adapter := realtime.NewDashboardAdapter(hub)

// Broadcast analytics updates to all clients
go func() {
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        metrics := analyticsConsumer.GetMetrics()
        adapter.BroadcastAnalyticsUpdate(metrics)
    }
}()

// Send notifications to specific users
adapter.SendUserNotification(userID, "Great Job!", "You completed lesson X", "high")

// Broadcast course updates
adapter.BroadcastCourseUpdate(courseID, courseMetrics)
```

**WebSocket Events:**
- `analytics_update`: Real-time metrics
- `notification`: User notifications
- `progress_update`: Student progress
- `course_update`: Course metrics
- `echo`: Echo received messages

## Event Publishing

### Publishing from Handlers

```go
// In your handler, after a user completes a lesson:
func (h *Handler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
    // ... business logic ...

    // Publish event to Kafka
    event := kafka.NewProgressEvent(
        kafka.ProgressCompleted,
        userID,
        lessonID,
        courseID,
        true,    // completed
        95,      // score
        3600,    // timeSpent seconds
    )

    if err := h.messagingService.PublishProgressUpdated(
        r.Context(),
        userID,
        lessonID,
        courseID,
        true,
        95,
        3600,
    ); err != nil {
        log.Printf("Error publishing progress: %v", err)
    }
}
```

### Batch Publishing

```go
messages := []kafka.EventMessage{
    {
        Topic:     topics.ProgressEvents,
        Key:       "user1:lesson1",
        Data:      progressEvent1,
        EventType: "progress.completed",
    },
    {
        Topic:     topics.ProgressEvents,
        Key:       "user2:lesson1",
        Data:      progressEvent2,
        EventType: "progress.completed",
    },
}

if err := producer.PublishBatch(ctx, messages); err != nil {
    log.Printf("Error publishing batch: %v", err)
}
```

## Monitoring & Observability

### Kafka UI

Access dashboard at `http://localhost:8080`:
- View all topics and messages
- Monitor broker health
- Inspect message content
- Consumer group lag
- Topic metrics

### Prometheus Metrics

Metrics available at `http://localhost:9090`:

```promql
# Kafka metrics
kafka_brokers
kafka_topics
kafka_consumer_lag

# Application metrics
go_pro_events_published_total
go_pro_events_processed_total
go_pro_events_failed_total
go_pro_dlq_messages_total
```

### Grafana Dashboards

Access at `http://localhost:3000`:
1. Add Prometheus as data source
2. Import Kafka dashboards
3. Create custom dashboards for:
   - Event throughput
   - Consumer lag
   - DLQ monitoring
   - Real-time metrics

## Performance Tuning

### Producer Configuration

```go
config := kafka.DefaultConfig()

// Batch settings for throughput
config.BatchSize = 16384           // 16KB batches
config.LingerMs = 5 * time.Millisecond // Wait up to 5ms for batches

// Compression for bandwidth
config.CompressionType = "snappy"  // or "gzip", "lz4", "zstd"

// Acks for durability
config.Acks = "all"                // Wait for all replicas
```

### Consumer Configuration

```go
// Throughput optimization
config.MaxPollRecords = 500
config.FetchMinBytes = 1024
config.FetchMaxWait = 500 * time.Millisecond

// Auto-commit for simplicity
config.EnableAutoCommit = true
config.AutoCommitInterval = 1 * time.Second
```

## Troubleshooting

### Check Kafka Cluster Health

```bash
# Get broker list
docker-compose -f docker-compose.kafka.yml exec kafka kafka-broker-api-versions.sh --bootstrap-server=kafka:29092

# List topics
docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics.sh --list --bootstrap-server=kafka:29092

# Describe topic
docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics.sh --describe --topic go-pro.progress.events --bootstrap-server=kafka:29092

# Check consumer groups
docker-compose -f docker-compose.kafka.yml exec kafka kafka-consumer-groups.sh --list --bootstrap-server=kafka:29092

# Check consumer lag
docker-compose -f docker-compose.kafka.yml exec kafka kafka-consumer-groups.sh --describe --group go-pro-consumer-group --bootstrap-server=kafka:29092
```

### View Kafka Logs

```bash
# Broker logs
docker-compose -f docker-compose.kafka.yml logs kafka

# Zookeeper logs
docker-compose -f docker-compose.kafka.yml logs zookeeper
```

### Common Issues

**Messages not flowing:**
- Check broker connectivity: `docker-compose -f docker-compose.kafka.yml exec kafka kafka-broker-api-versions.sh`
- Verify topic exists: `docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics.sh --list`
- Check consumer group lag

**High latency:**
- Reduce batch size for lower latency
- Check network connectivity between containers
- Monitor broker CPU/memory

**DLQ messages stuck:**
- Check DLQ consumer group status
- Manual message replay from Kafka UI
- Check logs for processing errors

## Example: Complete Real-Time Integration

```go
package main

import (
    "context"
    "log"
    "time"

    "go-pro-backend/internal/messaging/kafka"
    "go-pro-backend/internal/messaging/realtime"
)

func main() {
    // Setup
    config := kafka.DefaultConfig()
    topics := kafka.DefaultTopics()
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 1. Start WebSocket hub
    hub := realtime.NewHub()
    go hub.Start(ctx)

    // 2. Create analytics consumer
    analyticsConsumer := kafka.NewAnalyticsConsumer(config, topics)
    go analyticsConsumer.Start(ctx)

    // 3. Create notification service
    notificationService := kafka.NewNotificationService(config, topics)
    notificationService.RegisterHandler(&kafka.PrintNotificationHandler{})
    go notificationService.Start(ctx)

    // 4. Broadcast metrics every 5 seconds
    adapter := realtime.NewDashboardAdapter(hub)
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        for range ticker.C {
            metrics := analyticsConsumer.GetMetrics()
            adapter.BroadcastAnalyticsUpdate(metrics)
        }
    }()

    // 5. Wait for signals
    select {}
}
```

## Best Practices

1. **Always use context**: Pass context through all async operations
2. **Handle errors**: DLQ is automatic, but monitor DLQ messages
3. **Batch when possible**: Higher throughput, lower latency variance
4. **Use compression**: Save bandwidth for large payloads
5. **Monitor lag**: Alert when consumer lag increases
6. **Test in CI/CD**: Include Kafka in test pipeline
7. **Version events**: Include version field for future compatibility
8. **Idempotent processing**: Design handlers to handle duplicate messages

## References

- [Segmentio Kafka-Go Docs](https://pkg.go.dev/github.com/segmentio/kafka-go)
- [Kafka Official Docs](https://kafka.apache.org/documentation/)
- [Confluent Platform Docs](https://docs.confluent.io/)
