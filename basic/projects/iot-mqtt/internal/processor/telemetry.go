// Package processor provides telemetry data processing capabilities.
package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goproject/iot-mqtt/internal/mqtt"
)

// SensorReading represents telemetry data from a sensor.
type SensorReading struct {
	DeviceID    string  `json:"device_id"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Battery     int     `json:"battery"`
}

// Alert represents a threshold breach alert.
type Alert struct {
	DeviceID  string  `json:"device_id"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	Threshold float64 `json:"threshold"`
	Timestamp int64   `json:"timestamp"`
}

// Storage defines the interface for storing telemetry data.
type Storage interface {
	SaveReading(ctx context.Context, reading *SensorReading) error
	GetReadings(ctx context.Context, deviceID string, limit int) ([]*SensorReading, error)
}

// InMemoryStorage implements Storage with in-memory storage.
type InMemoryStorage struct {
	mu       sync.RWMutex
	readings map[string][]*SensorReading
}

// NewInMemoryStorage creates a new in-memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		readings: make(map[string][]*SensorReading),
	}
}

// SaveReading saves a sensor reading.
func (s *InMemoryStorage) SaveReading(ctx context.Context, reading *SensorReading) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.readings[reading.DeviceID] = append(s.readings[reading.DeviceID], reading)
	return nil
}

// GetReadings retrieves recent readings for a device.
func (s *InMemoryStorage) GetReadings(ctx context.Context, deviceID string, limit int) ([]*SensorReading, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	readings := s.readings[deviceID]
	if len(readings) <= limit {
		return readings, nil
	}

	return readings[len(readings)-limit:], nil
}

// TelemetryProcessor processes incoming telemetry data.
type TelemetryProcessor struct {
	client       *mqtt.Client
	storage      Storage
	alerts       []Alert
	tempHigh     float64
	tempLow      float64
	humidityHigh float64
	mu           sync.RWMutex
}

// NewTelemetryProcessor creates a new telemetry processor.
func NewTelemetryProcessor(client *mqtt.Client) *TelemetryProcessor {
	return &TelemetryProcessor{
		client:       client,
		storage:      NewInMemoryStorage(),
		tempHigh:     30.0, // High temperature threshold
		tempLow:      10.0, // Low temperature threshold
		humidityHigh: 80.0, // High humidity threshold
	}
}

// SetStorage sets the storage backend.
func (p *TelemetryProcessor) SetStorage(storage Storage) {
	p.storage = storage
}

// SetThresholds sets alert thresholds.
func (p *TelemetryProcessor) SetThresholds(tempHigh, tempLow, humidityHigh float64) {
	p.tempHigh = tempHigh
	p.tempLow = tempLow
	p.humidityHigh = humidityHigh
}

// Process parses and processes telemetry data.
func (p *TelemetryProcessor) Process(data []byte) (*Alert, error) {
	var reading SensorReading
	if err := json.Unmarshal(data, &reading); err != nil {
		return nil, fmt.Errorf("failed to parse telemetry: %w", err)
	}

	// Set timestamp if not present
	if reading.Timestamp == 0 {
		reading.Timestamp = time.Now().Unix()
	}

	// Save to storage
	if err := p.storage.SaveReading(context.Background(), &reading); err != nil {
		log.Printf("Failed to save reading: %v", err)
	}

	// Check thresholds
	alert := p.checkThresholds(&reading)
	if alert != nil {
		p.mu.Lock()
		p.alerts = append(p.alerts, *alert)
		p.mu.Unlock()
	}

	return alert, nil
}

// checkThresholds checks if any threshold is breached.
func (p *TelemetryProcessor) checkThresholds(reading *SensorReading) *Alert {
	if reading.Temperature > p.tempHigh {
		return &Alert{
			DeviceID:  reading.DeviceID,
			Type:      "HIGH_TEMPERATURE",
			Value:     reading.Temperature,
			Threshold: p.tempHigh,
			Timestamp: reading.Timestamp,
		}
	}

	if reading.Temperature < p.tempLow {
		return &Alert{
			DeviceID:  reading.DeviceID,
			Type:      "LOW_TEMPERATURE",
			Value:     reading.Temperature,
			Threshold: p.tempLow,
			Timestamp: reading.Timestamp,
		}
	}

	if reading.Humidity > p.humidityHigh {
		return &Alert{
			DeviceID:  reading.DeviceID,
			Type:      "HIGH_HUMIDITY",
			Value:     reading.Humidity,
			Threshold: p.humidityHigh,
			Timestamp: reading.Timestamp,
		}
	}

	return nil
}

// Start begins processing telemetry from the specified topic.
func (p *TelemetryProcessor) Start(ctx context.Context, topic string) error {
	handler := func(topic string, payload []byte) {
		alert, err := p.Process(payload)
		if err != nil {
			log.Printf("Failed to process telemetry: %v", err)
			return
		}

		if alert != nil {
			log.Printf("Alert: %s for device %s", alert.Type, alert.DeviceID)
		}
	}

	return p.client.Subscribe(topic, 1, handler)
}

// GetAlerts returns all generated alerts.
func (p *TelemetryProcessor) GetAlerts() []Alert {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.alerts
}

// GetReadings retrieves readings for a device.
func (p *TelemetryProcessor) GetReadings(deviceID string, limit int) ([]*SensorReading, error) {
	return p.storage.GetReadings(context.Background(), deviceID, limit)
}
