# 13 - Building IoT Applications with Go and MQTT

Develop IoT applications using Go and the MQTT protocol for lightweight, efficient device communication.

## Overview

MQTT (Message Queuing Telemetry Transport) is a lightweight publish-subscribe protocol ideal for IoT:
- **Lightweight**: Minimal bandwidth, small code footprint
- **Reliable**: QoS levels, persistent sessions, last will testament
- **Bidirectional**: Device-to-cloud and cloud-to-device messaging
- **Scalable**: Millions of devices through topic hierarchy

## Learning Objectives

- Understand MQTT protocol fundamentals
- Implement MQTT publishers and subscribers in Go
- Handle QoS levels and message delivery guarantees
- Manage device connections and sessions
- Process telemetry data from multiple devices
- Build command and control systems for IoT devices
- Implement secure device authentication
- Design scalable IoT architectures

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      IoT Architecture                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐                │
│  │ Sensor 1 │   │ Sensor 2 │   │ Sensor N │   (Devices)    │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘                │
│       │              │              │                        │
│       └──────────────┼──────────────┘                        │
│                      │                                       │
│                      ▼                                       │
│              ┌───────────────┐                              │
│              │  MQTT Broker  │   (Mosquitto, EMQX, HiveMQ)  │
│              └───────┬───────┘                              │
│                      │                                       │
│       ┌──────────────┼──────────────┐                        │
│       │              │              │                        │
│       ▼              ▼              ▼                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│  │ Backend  │  │ Dashboard│  │ Alerting │  (Applications)   │
│  │ Service  │  │ Service  │  │ Service  │                   │
│  └──────────┘  └──────────┘  └──────────┘                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Topic Hierarchy

MQTT uses slash-separated topic paths:

```
sensors/temperature/livingroom     # Specific sensor
sensors/temperature/+              # All temperature sensors
sensors/#                          # All sensors (wildcard)
devices/{device_id}/telemetry      # Device telemetry
devices/{device_id}/commands       # Device commands
devices/{device_id}/status         # Device status
```

## Prerequisites

### Install MQTT Broker

```bash
# Using Docker (recommended)
docker run -it -p 1883:1883 -p 9001:9001 \
  eclipse-mosquitto:2

# With configuration
docker run -it -p 1883:1883 -p 9001:9001 \
  -v /path/to/mosquitto.conf:/mosquitto/config/mosquitto.conf \
  eclipse-mosquitto:2

# Linux
sudo apt-get install mosquitto mosquitto-clients

# macOS
brew install mosquitto

# Start broker
mosquitto -c /etc/mosquitto/mosquitto.conf
```

### Install Go MQTT Client

```bash
go get github.com/eclipse/paho.mqtt.golang
```

## Quick Start

```bash
# Start MQTT broker
docker run -it -p 1883:1883 eclipse-mosquitto:2

# Run examples
cd advanced-topics/13-iot-mqtt/examples
go run publisher.go
go run subscriber.go
go run iot_device.go
go run telemetry_processor.go
```

## MQTT Concepts

### QoS Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| **QoS 0** | At most once | Sensor data (frequent, loss acceptable) |
| **QoS 1** | At least once | Commands, important telemetry |
| **QoS 2** | Exactly once | Financial transactions, critical data |

### Retained Messages

```go
// Publish retained message (last known value)
token := client.Publish("device/status", 1, true, "online")
```

### Last Will and Testament (LWT)

```go
// Declare LWT on connect - published if device disconnects unexpectedly
opts.SetWill("device/status", "offline", 1, true)
```

### Persistent Sessions

```go
// Clean session = false keeps subscriptions across reconnects
opts.SetCleanSession(false)
opts.SetClientID("unique-device-id")
```

## Code Examples

### 1. Basic Publisher

```go
package main

import (
    "log"
    "time"
    
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("publisher")
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    defer client.Disconnect(250)
    
    // Publish message
    token := client.Publish("sensors/temperature", 1, false, "23.5")
    token.Wait()
    
    log.Println("Message published")
}
```

### 2. Basic Subscriber

```go
package main

import (
    "log"
    
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("subscriber")
    opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
        log.Printf("Received: %s from %s", string(m.Payload()), m.Topic())
    })
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    defer client.Disconnect(250)
    
    // Subscribe to topic
    token := client.Subscribe("sensors/#", 1, nil)
    token.Wait()
    
    log.Println("Subscribed to sensors/#")
    select {} // Block forever
}
```

### 3. IoT Device Simulator

```go
type SensorDevice struct {
    client    mqtt.Client
    deviceID  string
    topicBase string
    interval  time.Duration
}

func (d *SensorDevice) Start(ctx context.Context) {
    ticker := time.NewTicker(d.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            d.client.Publish(d.topicBase+"/status", 1, true, "offline")
            return
        case <-ticker.C:
            d.publishTelemetry()
        }
    }
}

func (d *SensorDevice) publishTelemetry() {
    telemetry := map[string]interface{}{
        "device_id":  d.deviceID,
        "timestamp":  time.Now().Unix(),
        "temperature": 20 + rand.Float64()*10,
        "humidity":    40 + rand.Float64()*30,
        "battery":     80 + rand.Intn(20),
    }
    
    data, _ := json.Marshal(telemetry)
    token := d.client.Publish(d.topicBase+"/telemetry", 1, false, data)
    token.Wait()
}
```

### 4. Command and Control

```go
// Device subscribes to commands
func (d *SensorDevice) subscribeCommands() {
    d.client.Subscribe(d.topicBase+"/commands", 1, 
        func(c mqtt.Client, m mqtt.Message) {
            var cmd Command
            json.Unmarshal(m.Payload(), &cmd)
            
            switch cmd.Action {
            case "reboot":
                d.reboot()
            case "update_interval":
                d.updateInterval(cmd.Value)
            case "set_config":
                d.setConfig(cmd.Config)
            }
        })
}
```

### 5. Telemetry Processor

```go
type TelemetryProcessor struct {
    client   mqtt.Client
    storage  Storage
    alerts   AlertService
}

func (p *TelemetryProcessor) Start() {
    p.client.Subscribe("devices/+/telemetry", 1, 
        func(c mqtt.Client, m mqtt.Message) {
            var telemetry Telemetry
            json.Unmarshal(m.Payload(), &telemetry)
            
            // Store telemetry
            p.storage.Save(telemetry)
            
            // Check thresholds
            if telemetry.Temperature > 30 {
                p.alerts.Send(Alert{
                    Device: telemetry.DeviceID,
                    Type:   "HIGH_TEMPERATURE",
                    Value:  telemetry.Temperature,
                })
            }
        })
}
```

## Advanced Patterns

### 1. Device Registry

```go
type DeviceRegistry struct {
    devices sync.Map
    client  mqtt.Client
}

func (r *DeviceRegistry) Register(deviceID string, config DeviceConfig) {
    r.devices.Store(deviceID, &Device{
        ID:     deviceID,
        Config: config,
        Status: "registered",
    })
    
    // Subscribe to device topics
    r.client.Subscribe("devices/"+deviceID+"/+", 1, r.handleDeviceMessage)
}
```

### 2. Connection Management

```go
type ConnectionManager struct {
    client       mqtt.Client
    retryDelay   time.Duration
    maxRetries   int
    onConnect    func()
    onDisconnect func()
}

func (m *ConnectionManager) Connect() error {
    opts := mqtt.NewClientOptions()
    opts.OnConnect = func(c mqtt.Client) {
        log.Println("Connected to broker")
        if m.onConnect != nil {
            m.onConnect()
        }
    }
    opts.OnConnectionLost = func(c mqtt.Client, err error) {
        log.Printf("Connection lost: %v", err)
        if m.onDisconnect != nil {
            m.onDisconnect()
        }
    }
    opts.AutoReconnect = true
    opts.MaxReconnectInterval = 30 * time.Second
    
    m.client = mqtt.NewClient(opts)
    // ... connect logic
}
```

### 3. Batch Processing

```go
type BatchProcessor struct {
    buffer    []Telemetry
    ticker    *time.Ticker
    batchSize int
    storage   Storage
    mu        sync.Mutex
}

func (p *BatchProcessor) Add(t Telemetry) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.buffer = append(p.buffer, t)
    
    if len(p.buffer) >= p.batchSize {
        p.flush()
    }
}

func (p *BatchProcessor) flush() {
    if len(p.buffer) == 0 {
        return
    }
    
    p.storage.BatchSave(p.buffer)
    p.buffer = p.buffer[:0]
}
```

### 4. Security with TLS

```go
func createTLSConfig() *tls.Config {
    certPool := x509.NewCertPool()
    ca, _ := os.ReadFile("ca.crt")
    certPool.AppendCertsFromPEM(ca)
    
    cert, _ := tls.LoadX509KeyPair("client.crt", "client.key")
    
    return &tls.Config{
        RootCAs:      certPool,
        Certificates: []tls.Certificate{cert},
    }
}

opts := mqtt.NewClientOptions()
opts.AddBroker("tls://localhost:8883")
opts.SetTLSConfig(createTLSConfig())
```

## Project Structure

```
13-iot-mqtt/
├── README.md
├── examples/
│   ├── publisher.go           # Basic publisher
│   ├── subscriber.go          # Basic subscriber
│   ├── iot_device.go          # Device simulator
│   ├── telemetry_processor.go # Backend processor
│   ├── command_control.go     # Command handling
│   └── broker_config/
│       └── mosquitto.conf     # Broker configuration
```

## Best Practices

### 1. Topic Design
```
# Good: Hierarchical, specific
devices/{device_id}/telemetry/temperature
devices/{device_id}/telemetry/humidity
devices/{device_id}/commands
devices/{device_id}/status

# Avoid: Flat, ambiguous
device123
data
command
```

### 2. Message Payload
```go
type TelemetryMessage struct {
    DeviceID  string      `json:"device_id"`
    Timestamp int64       `json:"timestamp"`
    Type      string      `json:"type"`
    Value     interface{} `json:"value"`
    Metadata  Metadata    `json:"metadata,omitempty"`
}
```

### 3. Error Handling
```go
token := client.Publish(topic, qos, retained, payload)
go func() {
    token.Wait()
    if token.Error() != nil {
        log.Printf("Publish failed: %v", token.Error())
        // Retry logic
    }
}()
```

### 4. Graceful Shutdown
```go
func main() {
    // Setup...
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go device.Start(ctx)
    
    <-sigChan
    cancel()
    
    // Wait for cleanup
    time.Sleep(time.Second)
    client.Disconnect(250)
}
```

## Testing

```go
func TestTelemetryProcessing(t *testing.T) {
    // Use in-memory broker or mock
    broker := NewMockBroker()
    defer broker.Close()
    
    processor := NewTelemetryProcessor(broker.Client())
    
    // Simulate device publishing
    broker.Publish("devices/test/telemetry", telemetryData)
    
    // Verify processing
    assert.Equal(t, expectedValue, processor.GetLastValue())
}
```

## Performance Considerations

| Aspect | Recommendation |
|--------|----------------|
| QoS | Use QoS 0 for high-frequency sensors |
| Payload | Keep messages small (< 1KB) |
| Connections | Reuse connections, don't reconnect per message |
| Batching | Batch process telemetry on backend |
| Topics | Use wildcards wisely, avoid `#` where possible |

## Troubleshooting

### Connection Issues
```bash
# Test broker connectivity
mosquitto_pub -h localhost -t test -m "hello"
mosquitto_sub -h localhost -t test

# Verbose logging
opts.SetOrderMatters(false)
mqtt.DEBUG = log.New(os.Stdout, "", 0)
```

### Message Not Received
- Check topic spelling and wildcards
- Verify QoS settings
- Ensure subscriber is connected before publisher
- Check retained message flag

## Resources

- [MQTT Specification](https://mqtt.org/mqtt-specification/)
- [Paho Go Client](https://github.com/eclipse/paho.mqtt.golang)
- [Mosquitto Documentation](https://mosquitto.org/documentation/)
- [MQTT 5.0 Features](https://www.hivemq.com/mqtt-5/)

## Next Steps

1. Add JWT authentication for devices
2. Implement device shadow (state sync)
3. Build fleet management dashboard
4. Add edge computing capabilities
5. Implement OTA firmware updates
