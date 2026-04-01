# IoT-MQTT: MQTT Protocol with Eclipse Paho

Production-ready IoT MQTT template using eclipse/paho.mqtt.golang with device simulation and telemetry processing.

## Architecture

```
┌──────────┐   MQTT    ┌───────────┐   HTTP   ┌──────────┐
│  Device  │ ────────> │  Gateway  │ <──────> │  HTTP    │
│ (Sensor) │           │(Processor)│          │  Client  │
└──────────┘           └───────────┘          └──────────┘
       │                      │
       │                      │
       ▼                      ▼
┌──────────┐           ┌───────────┐
│ Mosquitto│           │ Telemetry  │
│  Broker  │           │  Storage   │
└──────────┘           └───────────┘
```

## Features

- **MQTT Client**: Robust connection management with auto-reconnect
- **Device Simulator**: Publishes telemetry at configurable intervals
- **Gateway Service**: Subscribes to telemetry and processes data
- **Alert System**: Threshold-based alerting
- **HTTP API**: Query devices, readings, and send commands
- **Docker Compose**: Full IoT stack with Mosquitto broker

## Quick Start

```bash
# Start the full IoT stack
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## MQTT Topics

| Topic | Direction | Description |
|-------|-----------|-------------|
| `devices/{id}/telemetry` | Device → Gateway | Sensor readings |
| `devices/{id}/commands` | Gateway → Device | Commands to device |
| `devices/{id}/status` | Device → Gateway | Online/offline status |

## HTTP API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/devices` | GET | List known devices |
| `/api/devices/{id}/readings` | GET | Get recent readings |
| `/api/devices/{id}/commands` | POST | Send command to device |
| `/api/alerts` | GET | Get current alerts |

## Environment Variables

### Device
| Variable | Description | Default |
|----------|-------------|---------|
| `DEVICE_ID` | Unique device identifier | `device-001` |
| `BROKER_URL` | MQTT broker URL | `tcp://localhost:1883` |
| `INTERVAL` | Telemetry publish interval | `30s` |

### Gateway
| Variable | Description | Default |
|----------|-------------|---------|
| `BROKER_URL` | MQTT broker URL | `tcp://localhost:1883` |
| `PORT` | HTTP API port | `8080` |

## Example Usage

### Publish telemetry manually

```bash
# Using mosquitto_pub (if installed)
mosquitto_pub -h localhost -t devices/test/telemetry -m '{"device_id":"test","temperature":25.5,"humidity":60}'
```

### Send command to device

```bash
curl -X POST http://localhost:8080/api/devices/device-001/commands \
  -H "Content-Type: application/json" \
  -d '{"action":"reboot"}'
```

### Get readings

```bash
curl http://localhost:8080/api/devices/device-001/readings
```

## Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Project Structure

```
iot-mqtt/
├── cmd/
│   ├── gateway/main.go     # MQTT gateway service
│   └── device/main.go      # Device simulator
├── internal/
│   ├── mqtt/client.go      # MQTT client wrapper
│   └── processor/telemetry.go # Telemetry processing
├── broker_config/
│   └── mosquitto.conf      # Broker configuration
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## License

MIT
