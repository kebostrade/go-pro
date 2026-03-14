package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TelemetryMessage struct {
	DeviceID  string                 `json:"device_id"`
	Type      string                 `json:"type"`
	Location  string                 `json:"location"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type DeviceState struct {
	DeviceID    string
	Type        string
	Location    string
	LastSeen    time.Time
	Status      string
	Temperature float64
	Humidity    float64
	Battery     int
}

type Alert struct {
	DeviceID  string    `json:"device_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type TelemetryProcessor struct {
	client        mqtt.Client
	devices       sync.Map
	alerts        []Alert
	thresholds    map[string]float64
	batchBuffer   []TelemetryMessage
	batchSize     int
	batchInterval time.Duration
	mu            sync.Mutex
}

func NewTelemetryProcessor(broker string) *TelemetryProcessor {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("telemetry-processor")
	opts.SetCleanSession(false)

	return &TelemetryProcessor{
		thresholds: map[string]float64{
			"temperature_max": 30.0,
			"temperature_min": 10.0,
			"humidity_max":    80.0,
			"battery_min":     20.0,
		},
		batchBuffer:   make([]TelemetryMessage, 0),
		batchSize:     10,
		batchInterval: 5 * time.Second,
	}
}

func (p *TelemetryProcessor) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("telemetry-processor")
	opts.SetCleanSession(false)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Telemetry processor connected")
		p.subscribe(c)
	}

	p.client = mqtt.NewClient(opts)
	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (p *TelemetryProcessor) subscribe(client mqtt.Client) {
	client.Subscribe("devices/+/telemetry", 1, p.handleTelemetry)
	client.Subscribe("devices/+/status", 1, p.handleStatus)
	client.Subscribe("devices/+/response", 1, p.handleResponse)
	log.Println("Subscribed to device topics")
}

func (p *TelemetryProcessor) handleTelemetry(c mqtt.Client, m mqtt.Message) {
	var msg TelemetryMessage
	if err := json.Unmarshal(m.Payload(), &msg); err != nil {
		log.Printf("Failed to parse telemetry: %v", err)
		return
	}

	p.updateDeviceState(msg)
	p.checkThresholds(msg)
	p.addToBatch(msg)
}

func (p *TelemetryProcessor) handleStatus(c mqtt.Client, m mqtt.Message) {
	deviceID := extractDeviceID(m.Topic())
	status := string(m.Payload())

	if state, ok := p.devices.Load(deviceID); ok {
		s := state.(*DeviceState)
		s.Status = status
		s.LastSeen = time.Now()
	}

	log.Printf("[STATUS] Device %s: %s", deviceID, status)
}

func (p *TelemetryProcessor) handleResponse(c mqtt.Client, m mqtt.Message) {
	deviceID := extractDeviceID(m.Topic())
	log.Printf("[RESPONSE] Device %s: %s", deviceID, string(m.Payload()))
}

func (p *TelemetryProcessor) updateDeviceState(msg TelemetryMessage) {
	state := &DeviceState{
		DeviceID: msg.DeviceID,
		Type:     msg.Type,
		Location: msg.Location,
		LastSeen: time.Now(),
		Status:   "online",
	}

	if temp, ok := msg.Data["temperature"].(float64); ok {
		state.Temperature = temp
	}
	if hum, ok := msg.Data["humidity"].(float64); ok {
		state.Humidity = hum
	}
	if batt, ok := msg.Data["battery"].(float64); ok {
		state.Battery = int(batt)
	}

	p.devices.Store(msg.DeviceID, state)
}

func (p *TelemetryProcessor) checkThresholds(msg TelemetryMessage) {
	now := time.Now()

	if temp, ok := msg.Data["temperature"].(float64); ok {
		if temp > p.thresholds["temperature_max"] {
			alert := Alert{
				DeviceID:  msg.DeviceID,
				Type:      "HIGH_TEMPERATURE",
				Value:     temp,
				Threshold: p.thresholds["temperature_max"],
				Message:   fmt.Sprintf("Temperature %.2f°C exceeds max threshold", temp),
				Timestamp: now,
			}
			p.addAlert(alert)
		}
		if temp < p.thresholds["temperature_min"] {
			alert := Alert{
				DeviceID:  msg.DeviceID,
				Type:      "LOW_TEMPERATURE",
				Value:     temp,
				Threshold: p.thresholds["temperature_min"],
				Message:   fmt.Sprintf("Temperature %.2f°C below min threshold", temp),
				Timestamp: now,
			}
			p.addAlert(alert)
		}
	}

	if hum, ok := msg.Data["humidity"].(float64); ok {
		if hum > p.thresholds["humidity_max"] {
			alert := Alert{
				DeviceID:  msg.DeviceID,
				Type:      "HIGH_HUMIDITY",
				Value:     hum,
				Threshold: p.thresholds["humidity_max"],
				Message:   fmt.Sprintf("Humidity %.2f%% exceeds max threshold", hum),
				Timestamp: now,
			}
			p.addAlert(alert)
		}
	}

	if batt, ok := msg.Data["battery"].(float64); ok {
		if batt < p.thresholds["battery_min"] {
			alert := Alert{
				DeviceID:  msg.DeviceID,
				Type:      "LOW_BATTERY",
				Value:     batt,
				Threshold: p.thresholds["battery_min"],
				Message:   fmt.Sprintf("Battery %.0f%% below min threshold", batt),
				Timestamp: now,
			}
			p.addAlert(alert)
		}
	}
}

func (p *TelemetryProcessor) addAlert(alert Alert) {
	p.alerts = append(p.alerts, alert)
	p.publishAlert(alert)
	log.Printf("⚠️  ALERT: [%s] %s", alert.DeviceID, alert.Message)
}

func (p *TelemetryProcessor) publishAlert(alert Alert) {
	payload, _ := json.Marshal(alert)
	p.client.Publish("alerts", 1, false, payload)
}

func (p *TelemetryProcessor) addToBatch(msg TelemetryMessage) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.batchBuffer = append(p.batchBuffer, msg)

	if len(p.batchBuffer) >= p.batchSize {
		p.flushBatch()
	}
}

func (p *TelemetryProcessor) flushBatch() {
	if len(p.batchBuffer) == 0 {
		return
	}

	log.Printf("Processing batch of %d messages", len(p.batchBuffer))

	for _, msg := range p.batchBuffer {
		_ = msg
	}

	p.batchBuffer = p.batchBuffer[:0]
}

func (p *TelemetryProcessor) StartBatchProcessor() {
	ticker := time.NewTicker(p.batchInterval)
	go func() {
		for range ticker.C {
			p.mu.Lock()
			p.flushBatch()
			p.mu.Unlock()
		}
	}()
}

func (p *TelemetryProcessor) GetDeviceCount() int {
	count := 0
	p.devices.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (p *TelemetryProcessor) GetAlertCount() int {
	return len(p.alerts)
}

func (p *TelemetryProcessor) ListDevices() []DeviceState {
	var devices []DeviceState
	p.devices.Range(func(key, value interface{}) bool {
		devices = append(devices, *value.(*DeviceState))
		return true
	})
	return devices
}

func (p *TelemetryProcessor) SendCommand(deviceID, action string, params map[string]interface{}) error {
	cmd := struct {
		Action string                 `json:"action"`
		Params map[string]interface{} `json:"params"`
	}{
		Action: action,
		Params: params,
	}

	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("devices/%s/commands", deviceID)
	token := p.client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

func (p *TelemetryProcessor) Disconnect() {
	p.client.Disconnect(250)
}

func extractDeviceID(topic string) string {
	parts := splitString(topic, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func main() {
	processor := NewTelemetryProcessor("tcp://localhost:1883")

	if err := processor.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer processor.Disconnect()

	processor.StartBatchProcessor()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("\n╔══════════════════════════════════════════════╗")
	fmt.Println("║      Telemetry Processor Running             ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println("\nSubscribed to:")
	fmt.Println("  - devices/+/telemetry")
	fmt.Println("  - devices/+/status")
	fmt.Println("  - devices/+/response")
	fmt.Println("\nThresholds:")
	fmt.Printf("  - Temperature: max %.0f°C, min %.0f°C\n",
		processor.thresholds["temperature_max"],
		processor.thresholds["temperature_min"])
	fmt.Printf("  - Humidity: max %.0f%%\n", processor.thresholds["humidity_max"])
	fmt.Printf("  - Battery: min %.0f%%\n", processor.thresholds["battery_min"])
	fmt.Println("\nPress Ctrl+C to stop\n")

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\nShutting down...")
			fmt.Printf("Total devices seen: %d\n", processor.GetDeviceCount())
			fmt.Printf("Total alerts: %d\n", processor.GetAlertCount())
			return

		case <-ticker.C:
			fmt.Printf("\n[Stats] Devices: %d | Alerts: %d\n",
				processor.GetDeviceCount(),
				processor.GetAlertCount())

			for _, d := range processor.ListDevices() {
				fmt.Printf("  - %s (%s): temp=%.1f°C, hum=%.1f%%, status=%s\n",
					d.DeviceID, d.Location, d.Temperature, d.Humidity, d.Status)
			}
		}
	}
}
