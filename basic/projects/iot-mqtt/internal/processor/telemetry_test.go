package processor

import (
	"context"
	"testing"
)

func TestInMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	reading := &SensorReading{
		DeviceID:    "device-1",
		Timestamp:   1234567890,
		Temperature: 25.5,
		Humidity:    60.0,
		Battery:     85,
	}

	// Save reading
	err := storage.SaveReading(context.Background(), reading)
	if err != nil {
		t.Errorf("SaveReading() error = %v", err)
	}

	// Get readings
	readings, err := storage.GetReadings(context.Background(), "device-1", 10)
	if err != nil {
		t.Errorf("GetReadings() error = %v", err)
	}

	if len(readings) != 1 {
		t.Errorf("GetReadings() returned %d readings, want 1", len(readings))
	}

	if readings[0].Temperature != 25.5 {
		t.Errorf("Temperature = %v, want 25.5", readings[0].Temperature)
	}
}

func TestTelemetryProcessor_Process(t *testing.T) {
	storage := NewInMemoryStorage()
	// Note: We can't actually connect to MQTT in tests, so we test Process directly

	p := &TelemetryProcessor{
		storage:      storage,
		tempHigh:     30.0,
		tempLow:      10.0,
		humidityHigh: 80.0,
	}

	// Test temperature above threshold
	data := []byte(`{"device_id":"test","temperature":35.0,"humidity":60.0,"battery":80}`)
	alert, err := p.Process(data)
	if err != nil {
		t.Errorf("Process() error = %v", err)
	}
	if alert == nil {
		t.Error("Expected alert for high temperature")
	}
	if alert.Type != "HIGH_TEMPERATURE" {
		t.Errorf("Alert type = %s, want HIGH_TEMPERATURE", alert.Type)
	}

	// Test normal reading
	data = []byte(`{"device_id":"test","temperature":20.0,"humidity":50.0,"battery":80}`)
	alert, err = p.Process(data)
	if err != nil {
		t.Errorf("Process() error = %v", err)
	}
	if alert != nil {
		t.Error("Expected no alert for normal reading")
	}
}

func TestTelemetryProcessor_SetThresholds(t *testing.T) {
	p := &TelemetryProcessor{}

	p.SetThresholds(40.0, 5.0, 90.0)

	if p.tempHigh != 40.0 {
		t.Errorf("tempHigh = %v, want 40.0", p.tempHigh)
	}
	if p.tempLow != 5.0 {
		t.Errorf("tempLow = %v, want 5.0", p.tempLow)
	}
	if p.humidityHigh != 90.0 {
		t.Errorf("humidityHigh = %v, want 90.0", p.humidityHigh)
	}
}
