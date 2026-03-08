# Building IoT Applications with Go and MQTT

Develop IoT applications using Go and the MQTT protocol.

## Learning Objectives

- Understand MQTT protocol
- Implement MQTT publishers and subscribers
- Handle device connectivity
- Process sensor data
- Implement quality of service levels
- Scale IoT deployments

## Theory

### MQTT Client Setup

```go
import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
    client mqtt.Client
}

func NewMQTTClient(broker, clientID string) (*MQTTClient, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID(clientID)
    opts.SetAutoReconnect(true)
    opts.SetCleanSession(true)
    opts.SetOnConnectHandler(func(c mqtt.Client) {
        log.Println("Connected to MQTT broker")
    })
    opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
        log.Printf("Connection lost: %v", err)
    })

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &MQTTClient{client: client}, nil
}
```

### Publisher Implementation

```go
type SensorPublisher struct {
    client *MQTTClient
    topic  string
    qos    byte
}

func NewSensorPublisher(client *MQTTClient, deviceID string) *SensorPublisher {
    return &SensorPublisher{
        client: client,
        topic:  fmt.Sprintf("devices/%s/sensors", deviceID),
        qos:    1,
    }
}

type SensorData struct {
    DeviceID  string    `json:"device_id"`
    Sensor    string    `json:"sensor"`
    Value     float64   `json:"value"`
    Unit      string    `json:"unit"`
    Timestamp time.Time `json:"timestamp"`
}

func (p *SensorPublisher) Publish(data *SensorData) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("marshal: %w", err)
    }

    token := p.client.client.Publish(p.topic, p.qos, false, payload)
    token.Wait()
    return token.Error()
}
```

### Subscriber Implementation

```go
type SensorSubscriber struct {
    client    *MQTTClient
    processor DataProcessor
}

type DataProcessor interface {
    Process(data *SensorData) error
}

func NewSensorSubscriber(client *MQTTClient, processor DataProcessor) *SensorSubscriber {
    return &SensorSubscriber{
        client:    client,
        processor: processor,
    }
}

func (s *SensorSubscriber) Subscribe(pattern string) error {
    token := s.client.client.Subscribe(pattern, 1, func(c mqtt.Client, m mqtt.Message) {
        var data SensorData
        if err := json.Unmarshal(m.Payload(), &data); err != nil {
            log.Printf("unmarshal error: %v", err)
            return
        }

        if err := s.processor.Process(&data); err != nil {
            log.Printf("process error: %v", err)
        }

        m.Ack()
    })

    token.Wait()
    return token.Error()
}
```

### Device Management

```go
type DeviceManager struct {
    client   *MQTTClient
    devices  map[string]*Device
    mu       sync.RWMutex
}

type Device struct {
    ID        string
    Name      string
    Type      string
    Status    string
    LastSeen  time.Time
    Config    DeviceConfig
}

type DeviceConfig struct {
    SampleRate int    `json:"sample_rate"`
    Unit       string `json:"unit"`
}

func NewDeviceManager(client *MQTTClient) *DeviceManager {
    dm := &DeviceManager{
        client:  client,
        devices: make(map[string]*Device),
    }

    client.client.Subscribe("devices/+/status", 1, dm.handleStatus)
    client.client.Subscribe("devices/+/config", 1, dm.handleConfig)

    return dm
}

func (dm *DeviceManager) handleStatus(c mqtt.Client, m mqtt.Message) {
    topicParts := strings.Split(m.Topic(), "/")
    if len(topicParts) < 2 {
        return
    }
    deviceID := topicParts[1]

    dm.mu.Lock()
    defer dm.mu.Unlock()

    if device, ok := dm.devices[deviceID]; ok {
        device.Status = string(m.Payload())
        device.LastSeen = time.Now()
    } else {
        dm.devices[deviceID] = &Device{
            ID:       deviceID,
            Status:   string(m.Payload()),
            LastSeen: time.Now(),
        }
    }
}

func (dm *DeviceManager) SendConfig(deviceID string, config DeviceConfig) error {
    topic := fmt.Sprintf("devices/%s/config", deviceID)
    payload, _ := json.Marshal(config)

    token := dm.client.client.Publish(topic, 1, false, payload)
    token.Wait()
    return token.Error()
}
```

### Data Aggregation

```go
type DataAggregator struct {
    buffer    map[string][]*SensorData
    threshold int
    flushFunc func(deviceID string, data []*SensorData) error
    mu        sync.Mutex
}

func NewDataAggregator(threshold int, flushFunc func(string, []*SensorData) error) *DataAggregator {
    a := &DataAggregator{
        buffer:    make(map[string][]*SensorData),
        threshold: threshold,
        flushFunc: flushFunc,
    }

    go a.periodicFlush()
    return a
}

func (a *DataAggregator) Add(data *SensorData) error {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.buffer[data.DeviceID] = append(a.buffer[data.DeviceID], data)

    if len(a.buffer[data.DeviceID]) >= a.threshold {
        return a.flush(data.DeviceID)
    }

    return nil
}

func (a *DataAggregator) flush(deviceID string) error {
    data := a.buffer[deviceID]
    if len(data) == 0 {
        return nil
    }

    if err := a.flushFunc(deviceID, data); err != nil {
        return err
    }

    a.buffer[deviceID] = a.buffer[deviceID][:0]
    return nil
}

func (a *DataAggregator) periodicFlush() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        a.mu.Lock()
        for deviceID := range a.buffer {
            a.flush(deviceID)
        }
        a.mu.Unlock()
    }
}
```

### QoS Levels

```go
const (
    QoS0 = 0
    QoS1 = 1
    QoS2 = 2
)

func (c *MQTTClient) PublishWithQoS(topic string, payload []byte, qos byte, retained bool) error {
    token := c.client.Publish(topic, qos, retained, payload)
    token.Wait()

    if token.Error() != nil {
        return token.Error()
    }

    if qos > 0 {
        log.Printf("Published with QoS %d, awaiting acknowledgment", qos)
    }

    return nil
}
```

## Security Considerations

```go
func NewSecureMQTTClient(broker, clientID, username, password string) (*MQTTClient, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID(clientID)
    opts.SetUsername(username)
    opts.SetPassword(password)
    opts.SetAutoReconnect(true)

    tlsConfig := &tls.Config{
        InsecureSkipVerify: false,
        MinVersion:         tls.VersionTLS12,
    }
    opts.SetTLSConfig(tlsConfig)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &MQTTClient{client: client}, nil
}

type ACL struct {
    topics map[string][]string
}

func (a *ACL) CanPublish(clientID, topic string) bool {
    allowed, ok := a.topics[clientID]
    if !ok {
        return false
    }
    for _, t := range allowed {
        if matchTopic(t, topic) {
            return true
        }
    }
    return false
}
```

## Performance Tips

```go
type ConnectionPool struct {
    clients chan *MQTTClient
    size    int
}

func NewConnectionPool(broker string, size int) (*ConnectionPool, error) {
    pool := &ConnectionPool{
        clients: make(chan *MQTTClient, size),
        size:    size,
    }

    for i := 0; i < size; i++ {
        client, err := NewMQTTClient(broker, fmt.Sprintf("pool-%d", i))
        if err != nil {
            return nil, err
        }
        pool.clients <- client
    }

    return pool, nil
}

func (p *ConnectionPool) Get() *MQTTClient {
    return <-p.clients
}

func (p *ConnectionPool) Put(client *MQTTClient) {
    p.clients <- client
}
```

## Exercises

1. Build a temperature monitoring system
2. Implement device provisioning
3. Create alerting on threshold breach
4. Build a data visualization dashboard

## Validation

```bash
cd exercises
mosquitto &
go test -v ./...
```

## Key Takeaways

- Use QoS 1 for reliable delivery
- Implement reconnection handling
- Aggregate data before persistence
- Secure connections with TLS
- Use topic patterns for routing

## Next Steps

**[AT-14: GraphQL gqlgen](../AT-14-graphql-gqlgen/README.md)**

---

IoT: connect everything, everywhere. 🌐
