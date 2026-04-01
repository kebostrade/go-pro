# NATS Event-Driven Template

A production-ready NATS JetStream event-driven architecture template with publish/subscribe patterns and queue workers.

## Architecture Overview

This template provides:

- **JetStream Publisher**: Persistent message publishing with acknowledgment
- **JetStream Subscriber**: Subscribe with automatic acknowledgment
- **Queue Workers**: Load-balanced worker pool pattern
- **Docker Compose**: Local NATS server with JetStream enabled

## Prerequisites

- Go 1.23+
- Docker and Docker Compose
- NATS CLI (optional, for testing)

## Project Structure

```
nats-events/
├── docker-compose.yml           # NATS server with JetStream
├── cmd/
│   ├── publisher/main.go       # Event publisher
│   ├── subscriber/main.go      # Event subscriber
│   └── worker/main.go          # Queue group worker
├── internal/
│   ├── jetstream/
│   │   ├── publisher.go       # JetStream publish helper
│   │   └── subscriber.go      # JetStream subscribe with ack
│   ├── queue/
│   │   └── worker.go          # Queue group worker pattern
│   └── models/
│       └── events.go          # Event type definitions
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```

## Quick Start

### 1. Start NATS Server

```bash
# Using Docker Compose
docker-compose up -d nats

# Verify NATS is running
docker-compose ps

# View NATS logs
docker-compose logs -f nats

# Or using nats-server directly
nats-server -js
```

### 2. Run the Publisher

```bash
go run ./cmd/publisher
```

The publisher will send order events to `events.orders` subject every 5 seconds.

### 3. Run the Subscriber

```bash
go run ./cmd/subscriber
```

The subscriber will receive and log all order events from `events.orders`.

### 4. Run Queue Workers

```bash
# Run multiple workers (each gets exclusive messages)
go run ./cmd/worker &
go run ./cmd/worker &
```

Workers in the same queue group will distribute messages among themselves.

## Using NATS CLI

Install the NATS CLI:

```bash
# Linux
curl -fsSL https://binaries.nats.dev/nats-io/natscli/releases/latest/download/nats-$(uname -s)-$(uname -m).tar.gz | tar zx

# macOS
brew install nats-io/nats-tools/nats
```

Test with NATS CLI:

```bash
# Subscribe to events
nats sub events.orders

# Publish a test event
nats pub events.orders '{"event_type":"created","order_id":"ORD-001","user_id":"USER-1","amount":99.99}'
```

## Event Types

### OrderEvent

```go
type OrderEvent struct {
    EventType string    `json:"event_type"` // "created", "updated", "deleted"
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    Timestamp time.Time `json:"timestamp"`
}
```

### UserEvent

```go
type UserEvent struct {
    EventType string    `json:"event_type"` // "registered", "updated", "deleted"
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}
```

### TaskEvent

```go
type TaskEvent struct {
    TaskID     string    `json:"task_id"`
    Type       string    `json:"type"` // "process", "email", "webhook"
    Payload    string    `json:"payload"`
    Priority   int       `json:"priority"`
    RetryCount int      `json:"retry_count"`
    Timestamp  time.Time `json:"timestamp"`
}
```

## JetStream Concepts

### Why JetStream?

- **Persistence**: Messages are stored and can be replayed
- **Delivery Guarantees**: At-least-once delivery with acknowledgments
- **Consumer Groups**: Multiple consumers can process the same stream
- **Replay Policies**: Replay from start, last, or by timestamp

### Key Patterns

1. **Publish with Persistence**: `js.Publish()` stores messages
2. **Subscribe with Ack**: Auto-acknowledgment after processing
3. **Queue Groups**: Load-balanced distribution across workers

## Development

### Building

```bash
# Build all binaries
make build

# Or individually
go build -o publisher ./cmd/publisher
go build -o subscriber ./cmd/subscriber
go build -o worker ./cmd/worker
```

### Testing

```bash
# Run tests
make test

# Run with race detection
go test -race ./...
```

### Docker

```bash
# Build Docker image
make docker-build

# Run NATS
make run-nats
```

## Configuration

### NATS Connection

Default connection URL: `nats://localhost:4222`

Change via environment variable or code:

```go
nc, err := nats.Connect(os.Getenv("NATS_URL"))
```

### Subjects

| Subject | Description |
|---------|-------------|
| `events.orders` | Order event stream |
| `events.users` | User event stream |
| `tasks` | Task processing queue |

## CI/CD

GitHub Actions workflow is configured in `.github/workflows/ci.yml`:

- Go 1.23 build and test
- Docker Compose validation
- Docker build

## Troubleshooting

### NATS not starting

```bash
# Check if port is already in use
lsof -i :4222

# View NATS logs
docker-compose logs nats
```

### Connection refused

```bash
# Verify NATS is running
curl http://localhost:8222/healthz

# Check NATS version
nats server --version
```

## References

- [NATS Documentation](https://docs.nats.io/)
- [JetStream Documentation](https://docs.nats.io/using-nats/developer/core_operations/jetstream)
- [nats.go Client](https://github.com/nats-io/nats.go)
