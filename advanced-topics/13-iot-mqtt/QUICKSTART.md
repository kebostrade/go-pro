# IoT with MQTT - Quick Start Guide

## Prerequisites

1. **Install MQTT Broker (Mosquitto)**
```bash
docker run -it -p 1883:1883 -p 9001:9001 eclipse-mosquitto:2
```

2. **Verify broker is running**
```bash
docker ps | grep mosquitto
```

## Running Examples

### 1. Basic Publisher (Sensor Simulator)
```bash
cd examples/01-publisher
go run .
```

Output:
```
2024/01/15 10:30:00 Connected to MQTT broker
2024/01/15 10:30:00 Starting publisher...
2024/01/15 10:30:02 Published: temp=28.45°C, humidity=65.23%, battery=89%
```

### 2. Basic Subscriber (Data Consumer)
```bash
cd examples/02-subscriber
go run .
```

Output:
```
2024/01/15 10:30:01 Connected to MQTT broker
2024/01/15 10:30:01 Subscribed to: devices/+/telemetry
2024/01/15 10:30:01 Subscribed successfully

[10:30:02] Device: sensor-001
  Temperature: 28.45°C
  Humidity:    65.23%
  Battery:     89%
```

### 3. IoT Device Simulator (Multiple Sensors)
```bash
cd examples/03-iot-device
go run .
```

Output:
```
2024/01/15 10:30:00 [temp-sensor-001] Connected to broker
2024/01/15 10:30:00 [motion-sensor-001] Connected to broker
2024/01/15 10:30:00 [air-quality-001] Connected to broker

IoT Devices Running:
  - temp-sensor-001 (temperature, living-room)
  - motion-sensor-001 (motion, entrance)
  - air-quality-001 (air_quality, bedroom)
```

### 4. Telemetry Processor (Backend Service)
```bash
cd examples/04-telemetry-processor
go run .
```

Output:
```
╔══════════════════════════════════════════════╗
║      Telemetry Processor Running             ║
╚══════════════════════════════════════════════╝

Subscribed to:
  - devices/+/telemetry
  - devices/+/status
  - devices/+/response

Thresholds:
  - Temperature: max 30°C, min 10°C
  - Humidity: max 80%
  - Battery: min 20%

[Stats] Devices: 3 | Alerts: 0
```

### 5. Command Controller (Device Management)
```bash
cd examples/05-command-control
go run .
```

Output:
```
╔══════════════════════════════════════════════╗
║        IoT Command Controller                ║
╚══════════════════════════════════════════════╝

Command Controller - Interactive Mode
Commands:
  list              - List all known devices
  online            - List online devices
  ping <device_id>  - Ping a specific device
```

## Testing with CLI Tools

### Install Mosquitto Clients
```bash
# Ubuntu/Debian
sudo apt-get install mosquitto-clients

# macOS
brew install mosquitto
```

### Manual Testing
```bash
# Subscribe to all device telemetry
mosquitto_sub -h localhost -t "devices/+/telemetry" -v

# Subscribe to all topics
mosquitto_sub -h localhost -t "#" -v

# Publish a test message
mosquitto_pub -h localhost -t "test/hello" -m "Hello MQTT"

# Send command to device
mosquitto_pub -h localhost -t "devices/sensor-001/commands" \
  -m '{"action":"ping"}'
```

## Common Patterns

### Topic Structure
```
devices/{device_id}/telemetry   # Sensor data
devices/{device_id}/status      # Online/offline status
devices/{device_id}/commands    # Commands to device
devices/{device_id}/response    # Device responses
alerts                          # System alerts
```

### Message Format (JSON)
```json
{
  "device_id": "sensor-001",
  "timestamp": 1705312200,
  "temperature": 25.5,
  "humidity": 60.2,
  "battery": 85
}
```

### QoS Levels
- **QoS 0**: Fire and forget (sensors with frequent updates)
- **QoS 1**: At least once (most IoT use cases)
- **QoS 2**: Exactly once (critical data)

## Troubleshooting

### Connection Refused
```bash
# Check if broker is running
docker ps | grep mosquitto

# Check broker logs
docker logs <container_id>
```

### Messages Not Received
1. Check topic spelling
2. Verify subscriber is running before publisher
3. Check QoS settings match
4. Verify broker is accessible

### Debug Mode
Add to your Go code:
```go
mqtt.DEBUG = log.New(os.Stdout, "MQTT: ", log.LstdFlags)
mqtt.ERROR = log.New(os.Stdout, "MQTT ERROR: ", log.LstdFlags)
```

## Architecture Diagram

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Publisher   │     │ IoT Devices  │     │   Sensors    │
│  (01-pub)    │     │  (03-device) │     │  (simulated) │
└──────┬───────┘     └──────┬───────┘     └──────────────┘
       │                    │
       │    MQTT Protocol   │
       └────────┬───────────┘
                │
                ▼
        ┌───────────────┐
        │  MQTT Broker  │
        │  (Mosquitto)  │
        └───────┬───────┘
                │
       ┌────────┼────────┐
       │        │        │
       ▼        ▼        ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│Subscriber│ │Processor │ │Controller│
│  (02-sub)│ │   (04)   │ │   (05)   │
└──────────┘ └──────────┘ └──────────┘
```

## Next Steps

1. Add TLS encryption for secure communication
2. Implement device authentication (username/password or certificates)
3. Add persistence with a database (PostgreSQL, InfluxDB)
4. Create a web dashboard for visualization
5. Implement OTA firmware updates
6. Add edge computing capabilities
