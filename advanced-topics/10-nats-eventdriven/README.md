# NATS & Event-Driven Go

Learn to build event-driven applications with NATS messaging system in Go.

## Overview

This module covers:
- **Basic Pub/Sub**: Simple publish-subscribe messaging
- **Queue Groups**: Work queue patterns for load balancing
- **Request/Reply**: Synchronous communication patterns
- **Wildcards**: Flexible subject-based routing
- **Connection Management**: Reliable connections and reconnection
- **Message Acknowledgment**: Reliable message delivery
- **JetStream Basics**: Persistence and message replay

## Prerequisites

### Install NATS Server

```bash
# Linux
curl -fsSL https://binaries.nats.dev/nats-io/nats-server/v2.10.18/nats-server-v2.10.18-linux-amd64.tar.gz | tar zx
cd nats-server-v2.10.18-linux-amd64
sudo cp nats-server /usr/local/bin/

# macOS
brew install nats-server

# Start NATS server
nats-server -js  # -js enables JetStream

# Or with Docker
docker run -p 4222:4222 -p 8222:8222 nats -js
```

### Install Go NATS Client

```bash
go get github.com/nats-io/nats.go/@latest
```

## Quick Start

### 1. Start NATS Server

```bash
# Basic server
nats-server

# With monitoring enabled
nats-server -m 8222

# With JetStream
nats-server -js

# Access monitoring dashboard
open http://localhost:8222
```

### 2. Run Examples

```bash
cd examples
go run nats_patterns.go
```

### 3. Try Specific Patterns

```bash
# Basic pub/sub
cd publisher
go run main.go

# In another terminal
cd subscriber
go run main.go

# Queue group
cd queue-group
go run subscriber.go & go run subscriber.go

# Request/reply
cd req-reply
go run responder.go &
go run requester.go
```

## NATS Concepts

### Subjects

Subjects are hierarchical, dot-separated strings:

```
orders                    # All orders
orders.new                # New orders
orders.shipped            # Shipped orders
orders.europe.de          # Orders from Germany
orders.>                 # All orders and sub-orders (wildcard)
orders.*                 # One level only
```

### Basic Pub/Sub

```go
// Connect to NATS
nc, _ := nats.Connect(nats.DefaultURL)

// Subscribe to subject
nc.Subscribe("greet", func(m *nats.Msg) {
    fmt.Printf("Received: %s\n", string(m.Data))
})

// Publish message
nc.Publish("greet", []byte("Hello NATS!"))
```

### Queue Groups

Queue groups distribute messages among subscribers:

```go
// Multiple subscribers with same queue name
// Only ONE receives each message (load balancing)
nc.QueueSubscribe("tasks", "workers", func(m *nats.Msg) {
    processTask(m)
})
```

### Request/Reply

Synchronous request-response pattern:

```go
// Responder
nc.Subscribe("help", func(m *nats.Msg) {
    m.Respond([]byte("I can help!"))
})

// Requester
msg, _ := nc.Request("help", []byte("help me"), time.Second)
fmt.Println(string(msg.Data))
```

### Wildcard Subscriptions

```
orders.*.created    # orders.new.created, orders.old.created
orders.>           # orders, orders.new, orders.new.created
_._.created        # Any two-level subject ending with .created
```

```go
nc.Subscribe("orders.>", func(m *nats.Msg) {
    // Receives orders.new, orders.new.created, etc.
})
```

## Message Patterns

### 1. Basic Pub/Sub

**Use case**: Broadcast messages to multiple consumers

```go
// Publisher
nc.Publish("updates", []byte("New update available"))

// Multiple subscribers
nc.Subscribe("updates", func(m *nats.Msg) {
    fmt.Println("Subscriber 1:", string(m.Data))
})

nc.Subscribe("updates", func(m *nats.Msg) {
    fmt.Println("Subscriber 2:", string(m.Data))
})
```

### 2. Queue Groups

**Use case**: Distribute work among workers

```go
// Multiple workers in same queue group
nc.QueueSubscribe("tasks", "worker-group", func(m *nats.Msg) {
    doWork(m)
})
```

**Result**: Each message goes to ONE worker only

### 3. Request/Reply

**Use case**: Synchronous queries

```go
// Service
nc.Subscribe("user.get", func(m *nats.Msg) {
    userId := string(m.Data)
    user := getUserFromDB(userId)
    m.Respond(user)
})

// Client
response, err := nc.Request("user.get", []byte("123"), time.Second)
```

### 4. Reply Subjects

**Use case**: Async responses to specific subjects

```go
nc.Subscribe("process", func(m *nats.Msg) {
    // Process and reply to specific inbox
    nc.Publish(m.Reply, []byte("done"))
})

inbox := nats.NewInbox()
nc.Subscribe(inbox, func(m *nats.Msg) {
    fmt.Println("Got response:", string(m.Data))
})

nc.PublishRequest("process", inbox, []byte("data"))
```

## Connection Management

### Basic Connection

```go
nc, err := nats.Connect("nats://localhost:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

### With Options

```go
nc, err := nats.Connect(
    "nats://localhost:4222",
    nats.Name("my-service"),
    nats.UserInfo("user", "pass"),
    nats.ReconnectWait(2*time.Second),
    nats.MaxReconnects(10),
    nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
        log.Printf("Disconnected: %v", err)
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected to %s", nc.ConnectedUrl())
    }),
    nats.ClosedHandler(func(nc *nats.Conn) {
        log.Println("Connection closed")
    }),
)
```

### Connection Status

```go
nc.Status()      // Current status (CONNECTED, CLOSED, etc.)
nc.IsConnected() // bool
nc.Stats()       // Connection statistics
```

## JetStream (Persistence)

### Enable JetStream

```bash
nats-server -js
```

### Create Stream

```go
js, _ := nc.JetStream()

// Create stream
js.AddStream(&nats.StreamConfig{
    Name:     "ORDERS",
    Subjects: []string{"orders.>"},
    Storage:  nats.FileStorage,
})
```

### Publish with Acknowledgment

```go
ack, err := js.Publish("orders.new", []byte(orderData))
if err != nil {
    log.Fatal(err)
}
fmt.Println("Stream Seq:", ack.Sequence)
```

### Durable Consumer

```go
// Create durable consumer
js.AddConsumer("ORDERS", &nats.ConsumerConfig{
    Durable:   "worker-1",
    AckPolicy: nats.AckExplicit,
})

// Subscribe as consumer
sub, _ := js.PullSubscribe("orders.new", "worker-1", nats.AckExplicit())

for {
    msg, err := sub.NextMsg(time.Second)
    if err != nil {
        continue
    }

    processOrder(msg)
    msg.Ack() // Acknowledge processing
}
```

### Consumer Types

```go
// Pull consumer - fetch messages on demand
sub, _ := js.PullSubscribe("orders.>", "pull-consumer")
msg, _ := sub.NextMsg(time.Second)

// Push consumer - automatic delivery
sub, _ := js.Subscribe("orders.>", nats.Durable("push-consumer"), nats.ManualAck())
```

## Best Practices

### 1. Connection Handling
- Use meaningful connection names
- Implement reconnect handlers
- Close connections gracefully

### 2. Subject Naming
- Use hierarchical naming
- Be specific but flexible
- Use wildcards appropriately

```go
// Good
orders.created
orders.europe.de.created

// Avoid too broad
orders
>
```

### 3. Error Handling
```go
nc, err := nats.Connect(nats.DefaultURL)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Subscribe with error handler
nc.Subscribe("subject", func(m *nats.Msg) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic in handler: %v", r)
        }
    }()

    handleMessage(m)
})
```

### 4. Message Validation
```go
type Order struct {
    ID     string  `json:"id"`
    Amount float64 `json:"amount"`
}

nc.Subscribe("orders.new", func(m *nats.Msg) {
    var order Order
    if err := json.Unmarshal(m.Data, &order); err != nil {
        log.Printf("Invalid message: %v", err)
        return
    }

    processOrder(order)
})
```

### 5. Timeouts
```go
msg, err := nc.Request("subject", data, 2*time.Second)
if err != nil {
    if err == nats.ErrTimeout {
        log.Println("Request timed out")
    }
    return
}
```

## Monitoring

### NATS Monitoring

```bash
# Start with monitoring
nats-server -m 8222

# View metrics
open http://localhost:8222

# Connection stats
curl http://localhost:8222/varz

# Route stats
curl http://localhost:8222/routez

# Subscription info
curl http://localhost:8222/connz
```

### Application Metrics

```go
nc.Subscribe("metrics.>", func(m *nats.Msg) {
    // Process and aggregate metrics
})

// Publish metrics
nc.Publish("metrics.requests", []byte("100"))
nc.Publish("metrics.errors", []byte("5"))
```

## Debugging

```bash
# View server logs
nats-server -DV  # Debug and verbose

# Trace messages
nats-server -trace

# Monitor connections
nats-server -SDV

# Using CLI tools
nats server info
nats server check
nats stream ls
nats consumer ls
```

## Common Patterns

### CQRS (Command Query Responsibility Segregation)

```go
// Write side (Command)
nc.Publish("commands.createOrder", orderData)

// Read side (Query)
msg, _ := nc.Request("queries.getOrder", orderId, time.Second)
```

### Event Sourcing

```go
// Stream all events
stream, _ := js.AddStream(&nats.StreamConfig{
    Name:     "EVENTS",
    Subjects: []string{"events.>"},
    Retention: nats.LimitsPolicy,
    MaxAge:   time.Hour * 24 * 30, // 30 days
})

// Replay events
sub, _ := js.PullSubscribe("events.>", "", nats.OptStartAtSequence(1))
```

### Saga Pattern

```go
// Orchestration
func executeOrder(order Order) {
    // Step 1: Reserve inventory
    msg, _ := nc.Request("inventory.reserve", order, time.Second)
    if !reserveSuccess(msg) {
        publishEvent("order.failed", order)
        return
    }

    // Step 2: Process payment
    msg, _ = nc.Request("payment.process", order, time.Second)
    if !paymentSuccess(msg) {
        nc.Publish("inventory.cancel", order)
        publishEvent("order.failed", order)
        return
    }

    // Step 3: Confirm order
    publishEvent("order.completed", order)
}
```

## Troubleshooting

### Connection Issues
```go
// Add connection handlers
nc, err := nats.Connect(
    nats.DefaultURL,
    nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
        log.Printf("Disconnected: %v", err)
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected to %s", nc.ConnectedUrl())
    }),
)
```

### Message Loss
- Use JetStream for persistence
- Implement acknowledgments
- Check consumer lag

```go
js.GetStreamInfo("ORDERS")
js.GetConsumerInfo("ORDERS", "consumer-name")
```

## Additional Resources

- [NATS Documentation](https://docs.nats.io/)
- [Go Client Guide](https://github.com/nats-io/nats.go)
- [JetStream Guide](https://docs.nats.io/nats-concepts/jetstream)
- [NATS by Example](https://natsbyexample.com/)

## Next Steps

After completing this module:
1. Study service mesh with NATS (NATS Service)
2. Learn about NATS KV (Key-Value store)
3. Explore object store for large data
4. Build microservices with NATS
