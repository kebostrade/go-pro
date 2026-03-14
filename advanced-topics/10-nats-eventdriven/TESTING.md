# NATS Testing Guide

How to test and debug NATS messaging applications.

## Testing Strategy

### 1. Start NATS Server

```bash
# Terminal 1: Start NATS server
nats-server -js

# With monitoring
nats-server -m 8222

# Access monitoring dashboard
open http://localhost:8222
```

### 2. Test Basic Pub/Sub

```bash
# Terminal 2: Start subscriber
cd subscriber
go mod tidy
go run main.go

# Terminal 3: Start publisher
cd publisher
go mod tidy
go run main.go
```

**Expected Output**:
- Subscriber receives all messages published
- Multiple subscribers each receive all messages (broadcast)

### 3. Test Queue Groups

```bash
# Terminal 2-4: Start multiple workers
cd queue-group
go mod tidy
go run worker.go worker-1 &
go run worker.go worker-2 &
go run worker.go worker-3

# Terminal 5: Publish tasks
go run publisher.go
```

**Expected Output**:
- Each task processed by exactly ONE worker
- Work distributed across all workers
- No task processed twice

### 4. Test Request/Reply

```bash
# Terminal 2: Start responder
cd req-reply
go mod tidy
go run responder.go

# Terminal 3: Start requester
go run requester.go
```

**Expected Output**:
- Requester sends requests for user IDs
- Responder responds with user data
- Each request gets exactly one response

### 5. Test All Patterns

```bash
cd examples
go mod tidy
go run nats_patterns.go
```

## Debugging NATS

### Server-Side Debugging

```bash
# Start with verbose logging
nats-server -DV

# Trace all messages
nats-server -trace

# Monitor connections
nats-server -SDV

# View server stats
curl http://localhost:8222/varz

# View connections
curl http://localhost:8222/connz

# View routes
curl http://localhost:8222/routez
```

### Client-Side Debugging

#### Connection Issues

```go
nc, err := nats.Connect(
    nats.DefaultURL,
    nats.Name("my-app"),
    nats.Verbose(),  // Enable verbose logging
    nats.PingInterval(20*time.Second),
    nats.MaxPingsOutstanding(5),
    nats.ReconnectWait(2*time.Second),
    nats.MaxReconnects(10),
    nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
        log.Printf("Disconnected: %v", err)
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected: %s", nc.ConnectedUrl())
    }),
    nats.ClosedHandler(func(nc *nats.Conn) {
        log.Println("Connection closed")
    }),
)
```

#### Message Inspection

```go
// Subscribe to all messages for debugging
nc.Subscribe(">", func(m *nats.Msg) {
    log.Printf("Subject: %s, Reply: %s, Data: %s",
        m.Subject, m.Reply, string(m.Data))
})
```

#### Connection Status

```go
// Check connection status
status := nc.Status()
fmt.Printf("Status: %v\n", status)

// Check if connected
if nc.IsConnected() {
    fmt.Println("Connected to NATS")
}

// Get connection stats
stats := nc.Stats()
fmt.Printf("Messages sent: %d\n", stats.MsgsSent)
fmt.Printf("Messages received: %d\n", stats.MsgsReceived)
```

## Common Issues

### 1. Connection Refused

**Error**: `dial tcp 127.0.0.1:4222: connect: connection refused`

**Solution**:
```bash
# Start NATS server
nats-server

# Or check if running
ps aux | grep nats-server
```

### 2. Timeout Errors

**Error**: `nats: timeout`

**Solutions**:
```go
// Increase timeout
msg, err := nc.Request("subject", data, 5*time.Second)

// Or use no timeout
sub, _ := nc.SubscribeSync("subject")
msg, err := sub.NextMsg(time.Minute)
```

### 3. No Messages Received

**Debug Steps**:

```bash
# Check server is running
curl http://localhost:8222/varz

# Check connection status
# Add Verbose() to connection options

# Check subject spelling
# Verify publisher and subscriber use same subject

# Add debug subscription
nc.Subscribe(">", func(m *nats.Msg) {
    log.Printf("DEBUG: %s -> %s", m.Subject, string(m.Data))
})
```

### 4. Queue Groups Not Load Balancing

**Problem**: All workers receive all messages

**Solution**:
```go
// Ensure same queue name
nc.QueueSubscribe("tasks", "workers", handler)  // Correct

// NOT different queue names
nc.QueueSubscribe("tasks", "worker-1", handler)  // Wrong
nc.QueueSubscribe("tasks", "worker-2", handler)  // Wrong
```

### 5. Request/Reply Hangs

**Problem**: Requester waits forever

**Debug**:

```go
// Add timeout
msg, err := nc.Request("subject", data, 2*time.Second)
if err != nil {
    if err == nats.ErrTimeout {
        log.Println("No responders available")
    }
    return
}

// Check if responders exist
// Add logging in responder handler
nc.Subscribe("subject", func(m *nats.Msg) {
    log.Println("Received request")
    m.Respond([]byte("response"))
})
```

## Performance Testing

### Load Testing

```go
package main

import (
    "fmt"
    "sync"
    "time"

    "github.com/nats-io/nats.go"
)

func benchmarkPublish() {
    nc, _ := nats.Connect(nats.DefaultURL)
    defer nc.Close()

    numMessages := 10000
    var wg sync.WaitGroup

    start := time.Now()

    for i := 0; i < numMessages; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            nc.Publish("test", []byte(fmt.Sprintf("msg-%d", idx)))
        }(i)
    }

    wg.Wait()
    nc.Flush()

    elapsed := time.Since(start)
    fmt.Printf("Published %d messages in %v\n", numMessages, elapsed)
    fmt.Printf("Messages/sec: %.2f\n", float64(numMessages)/elapsed.Seconds())
}

func benchmarkSubscribe() {
    nc, _ := nats.Connect(nats.DefaultURL)
    defer nc.Close()

    received := 0
    sub, _ := nc.Subscribe("test", func(m *nats.Msg) {
        received++
    })

    start := time.Now()

    // Publish messages
    for i := 0; i < 10000; i++ {
        nc.Publish("test", []byte("data"))
    }
    nc.Flush()

    // Wait for all messages
    for received < 10000 {
        time.Sleep(10 * time.Millisecond)
    }

    elapsed := time.Since(start)
    fmt.Printf("Received %d messages in %v\n", received, elapsed)

    sub.Unsubscribe()
}
```

## Monitoring

### Application Metrics

```go
type Metrics struct {
    MessagesSent     int64
    MessagesReceived int64
    Errors           int64
    Latency          time.Duration
}

func (m *Metrics) RecordLatency(start time.Time) {
    m.Latency += time.Since(start)
}

func (m *Metrics) Report() {
    fmt.Printf("Sent: %d, Received: %d, Errors: %d\n",
        m.MessagesSent, m.MessagesReceived, m.Errors)
    fmt.Printf("Avg Latency: %v\n",
        m.Latency/time.Duration(m.MessagesReceived))
}
```

### NATS Monitoring

```bash
# Server statistics
curl -s http://localhost:8222/varz | jq

# Connection details
curl -s http://localhost:8222/connz | jq

# View subscriptions
curl -s http://localhost:8222/subz | jq
```

## Testing Checklist

- [ ] NATS server running
- [ ] Connection successful
- [ ] Publish to subject
- [ ] Subscribe to subject
- [ ] Messages received
- [ ] Queue group load balancing
- [ ] Request/Reply working
- [ ] Wildcard subscriptions working
- [ ] Connection recovery on restart
- [ ] Performance acceptable

## Test Scripts

### test_pubsub.sh

```bash
#!/bin/bash

echo "Testing NATS Pub/Sub"

# Start server
nats-server -p 4222 &
SERVER_PID=$!
sleep 2

# Run subscriber in background
cd subscriber
go run main.go &
SUB_PID=$!
cd ..

sleep 1

# Run publisher
cd publisher
go run main.go
cd ..

# Wait for subscriber
sleep 2
kill $SUB_PID

# Stop server
kill $SERVER_PID

echo "Pub/Sub test complete"
```

### test_queue_group.sh

```bash
#!/bin/bash

echo "Testing NATS Queue Groups"

# Start server
nats-server -p 4222 &
SERVER_PID=$!
sleep 2

# Start 3 workers
cd queue-group
for i in 1 2 3; do
    go run worker.go worker-$i &
    WORKER_PIDS[$i]=$!
done
cd ..

sleep 1

# Publish tasks
cd queue-group
go run publisher.go
cd ..

# Wait for completion
sleep 5

# Cleanup
for pid in "${WORKER_PIDS[@]}"; do
    kill $pid 2>/dev/null
done
kill $SERVER_PID

echo "Queue group test complete"
```

## Integration Testing

```go
package main

import (
    "testing"
    "time"

    "github.com/nats-io/nats.go"
)

func TestPubSub(t *testing.T) {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        t.Fatal(err)
    }
    defer nc.Close()

    received := false
    sub, _ := nc.Subscribe("test", func(m *nats.Msg) {
        received = true
    })

    nc.Publish("test", []byte("data"))
    nc.Flush()

    time.Sleep(100 * time.Millisecond)

    if !received {
        t.Error("Message not received")
    }

    sub.Unsubscribe()
}

func TestRequestReply(t *testing.T) {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        t.Fatal(err)
    }
    defer nc.Close()

    nc.Subscribe("echo", func(m *nats.Msg) {
        m.Respond(m.Data)
    })

    msg, err := nc.Request("echo", []byte("hello"), time.Second)
    if err != nil {
        t.Fatal(err)
    }

    if string(msg.Data) != "hello" {
        t.Errorf("Expected 'hello', got '%s'", string(msg.Data))
    }
}
```

## Best Practices

1. **Always handle connection errors**
2. **Use timeouts for requests**
3. **Implement reconnection logic**
4. **Add logging for debugging**
5. **Test connection recovery**
6. **Monitor message rates**
7. **Handle errors gracefully**
8. **Use appropriate subjects**
9. **Clean up subscriptions**
10. **Test under load**
