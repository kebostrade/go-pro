# Building Event-Driven Applications with Go and NATS

Create scalable event-driven architectures using NATS messaging.

## Learning Objectives

- Understand pub/sub patterns
- Implement NATS publishers and subscribers
- Handle message acknowledgment
- Implement request-reply patterns
- Use JetStream for persistence
- Scale with queue groups

## Theory

### Basic Publisher/Subscriber

```go
type Publisher struct {
    nc *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
    nc, err := nats.Connect(url,
        nats.ReconnectWait(2*time.Second),
        nats.MaxReconnects(5),
    )
    if err != nil {
        return nil, fmt.Errorf("connect: %w", err)
    }
    return &Publisher{nc: nc}, nil
}

func (p *Publisher) Publish(subject string, data interface{}) error {
    msg, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("marshal: %w", err)
    }
    return p.nc.Publish(subject, msg)
}

type Subscriber struct {
    nc *nats.Conn
}

func (s *Subscriber) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
    return s.nc.Subscribe(subject, func(m *nats.Msg) {
        handler(m)
        m.Ack()
    })
}
```

### Queue Groups (Load Balancing)

```go
func (s *Subscriber) SubscribeQueue(subject, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
    return s.nc.QueueSubscribe(subject, queue, func(m *nats.Msg) {
        log.Printf("Processing message on %s", queue)
        handler(m)
        m.Ack()
    })
}
```

### JetStream for Persistence

```go
func setupJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
    js, err := nc.JetStream()
    if err != nil {
        return nil, fmt.Errorf("jetstream: %w", err)
    }

    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "ORDERS",
        Subjects: []string{"orders.*"},
        Retention: nats.LimitsPolicy,
        MaxMsgs:  10000,
        MaxBytes: 100 * 1024 * 1024,
    })
    if err != nil {
        return nil, fmt.Errorf("add stream: %w", err)
    }

    return js, nil
}

func (p *Publisher) PublishOrder(js nats.JetStreamContext, order *Order) error {
    data, _ := json.Marshal(order)
    _, err := js.Publish("orders.created", data)
    return err
}

func (s *Subscriber) ConsumeOrders(js nats.JetStreamContext, durable string, handler func(order *Order) error) (*nats.Subscription, error) {
    return js.Subscribe("orders.*", func(m *nats.Msg) {
        var order Order
        if err := json.Unmarshal(m.Data, &order); err != nil {
            m.Nak()
            return
        }

        if err := handler(&order); err != nil {
            m.Nak()
            return
        }

        m.Ack()
    }, nats.Durable(durable), nats.ManualAck())
}
```

### Request-Reply Pattern

```go
func (s *Subscriber) RegisterHandler(subject string, handler func(req []byte) ([]byte, error)) error {
    _, err := s.nc.Subscribe(subject, func(m *nats.Msg) {
        resp, err := handler(m.Data)
        if err != nil {
            m.Respond([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
            return
        }
        m.Respond(resp)
    })
    return err
}

func (p *Publisher) Request(ctx context.Context, subject string, req interface{}, resp interface{}) error {
    data, _ := json.Marshal(req)
    msg, err := p.nc.RequestWithContext(ctx, subject, data)
    if err != nil {
        return fmt.Errorf("request: %w", err)
    }
    return json.Unmarshal(msg.Data, resp)
}
```

### Event-Driven Service

```go
type OrderService struct {
    nc *nats.Conn
    js nats.JetStreamContext
    repo OrderRepository
}

func (s *OrderService) Start() error {
    _, err := s.js.Subscribe("orders.created", s.handleOrderCreated,
        nats.Durable("order-processor"),
        nats.DeliverAll(),
        nats.ManualAck(),
    )
    return err
}

func (s *OrderService) handleOrderCreated(m *nats.Msg) {
    var event OrderCreatedEvent
    if err := json.Unmarshal(m.Data, &event); err != nil {
        log.Printf("unmarshal error: %v", err)
        m.Nak()
        return
    }

    if err := s.processOrder(&event); err != nil {
        log.Printf("process error: %v", err)
        m.Nak()
        return
    }

    m.Ack()

    s.nc.Publish("orders.completed", m.Data)
}

func (s *OrderService) processOrder(event *OrderCreatedEvent) error {
    return s.repo.UpdateStatus(event.OrderID, "completed")
}
```

## Security Considerations

```go
func connectWithAuth(url, user, pass string) (*nats.Conn, error) {
    return nats.Connect(url,
        nats.UserInfo(user, pass),
        nats.Secure(&tls.Config{
            MinVersion: tls.VersionTLS12,
        }),
    )
}

type Event struct {
    Type      string          `json:"type"`
    Timestamp int64           `json:"timestamp"`
    Source    string          `json:"source"`
    Data      json.RawMessage `json:"data"`
    Signature string          `json:"signature,omitempty"`
}

func verifyEvent(e *Event, secret string) bool {
    expected := hmacSHA256(e.Data, secret)
    return subtle.ConstantTimeCompare([]byte(e.Signature), []byte(expected)) == 1
}
```

## Performance Tips

```go
type BatchPublisher struct {
    nc     *nats.Conn
    batch  [][]byte
    mu     sync.Mutex
    ticker *time.Ticker
}

func NewBatchPublisher(nc *nats.Conn, flushInterval time.Duration) *BatchPublisher {
    bp := &BatchPublisher{
        nc:     nc,
        batch:  make([][]byte, 0, 100),
        ticker: time.NewTicker(flushInterval),
    }
    go bp.flushLoop()
    return bp
}

func (bp *BatchPublisher) Add(subject string, data []byte) {
    bp.mu.Lock()
    bp.batch = append(bp.batch, data)
    if len(bp.batch) >= 100 {
        bp.flush()
    }
    bp.mu.Unlock()
}

func (bp *BatchPublisher) flushLoop() {
    for range bp.ticker.C {
        bp.mu.Lock()
        bp.flush()
        bp.mu.Unlock()
    }
}

func (bp *BatchPublisher) flush() {
    for _, msg := range bp.batch {
        bp.nc.Publish("batch.process", msg)
    }
    bp.batch = bp.batch[:0]
}
```

## Exercises

1. Build an order processing pipeline
2. Implement saga pattern with compensation
3. Create a notification service
4. Add dead letter handling

## Validation

```bash
cd exercises
nats-server &
go test -v ./...
```

## Key Takeaways

- Use queue groups for horizontal scaling
- Acknowledge messages after processing
- Use JetStream for guaranteed delivery
- Implement idempotent handlers
- Handle reconnection gracefully

## Next Steps

**[AT-09: Kubernetes Cloud](../AT-09-kubernetes-cloud/README.md)**

---

Events drive modern systems. 📡
