package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goproject/iot-mqtt/internal/mqtt"
)

// SensorDevice simulates an IoT sensor device publishing telemetry.
type SensorDevice struct {
	client      *mqtt.Client
	deviceID    string
	topicBase   string
	interval    time.Duration
	temperature float64
	humidity    float64
	battery     int
}

func main() {
	log.Println("IoT Device Simulator Starting...")

	deviceID := os.Getenv("DEVICE_ID")
	if deviceID == "" {
		deviceID = "device-001"
	}

	broker := os.Getenv("BROKER_URL")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	intervalStr := os.Getenv("INTERVAL")
	interval := 30 * time.Second
	if intervalStr != "" {
		if d, err := time.ParseDuration(intervalStr); err == nil {
			interval = d
		}
	}

	device := &SensorDevice{
		deviceID:    deviceID,
		topicBase:   "devices/" + deviceID,
		interval:    interval,
		temperature: 20.0 + rand.Float64()*10, // 20-30°C
		humidity:    40.0 + rand.Float64()*40, // 40-80%
		battery:     80 + rand.Intn(20),       // 80-100%
	}

	// Create MQTT client
	client := mqtt.NewClient(broker, deviceID)

	// Configure LWT for graceful disconnect detection
	client.SetWill(device.topicBase+"/status", "offline", 1, true)

	// Connect
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	defer client.Disconnect()

	device.client = client

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start publishing
	go device.start(ctx)

	log.Printf("Device %s publishing to %s every %v", deviceID, device.topicBase, interval)

	// Wait for shutdown
	<-sigChan
	log.Println("Shutting down device...")
	cancel()

	// Publish offline status
	client.Publish(device.topicBase+"/status", "offline", 1)

	log.Println("Device stopped")
}

func (d *SensorDevice) start(ctx context.Context) {
	// Publish online status
	d.client.Publish(d.topicBase+"/status", "online", 1)

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.publishTelemetry()
			// Occasionally update simulated values
			d.temperature += (rand.Float64() - 0.5) * 2 // -1 to +1
			d.humidity += (rand.Float64() - 0.5) * 5    // -2.5 to +2.5
			if d.humidity < 0 {
				d.humidity = 0
			}
			if d.humidity > 100 {
				d.humidity = 100
			}
		}
	}
}

func (d *SensorDevice) publishTelemetry() {
	reading := map[string]interface{}{
		"device_id":   d.deviceID,
		"timestamp":   time.Now().Unix(),
		"temperature": d.temperature,
		"humidity":    d.humidity,
		"battery":     d.battery,
	}

	data, err := json.Marshal(reading)
	if err != nil {
		log.Printf("Failed to marshal telemetry: %v", err)
		return
	}

	err = d.client.Publish(d.topicBase+"/telemetry", data, 1)
	if err != nil {
		log.Printf("Failed to publish telemetry: %v", err)
	} else {
		log.Printf("Published telemetry: temperature=%.1f, humidity=%.1f", d.temperature, d.humidity)
	}
}
