package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	DeviceID    string  `json:"device_id"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Battery     int     `json:"battery"`
}

type Alert struct {
	DeviceID string  `json:"device_id"`
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
	Message  string  `json:"message"`
}

var (
	temperatureThreshold float64 = 30.0
	humidityThreshold    float64 = 75.0
	alerts               []Alert
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("subscriber-001")
	opts.SetCleanSession(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
		log.Println("Subscribing to: devices/+/telemetry")

		if token := c.Subscribe("devices/+/telemetry", 1, handleMessage); token.Wait() && token.Error() != nil {
			log.Printf("Subscribe failed: %v", token.Error())
		}

		if token := c.Subscribe("devices/+/status", 1, handleStatus); token.Wait() && token.Error() != nil {
			log.Printf("Subscribe failed: %v", token.Error())
		}

		log.Println("Subscribed successfully")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect: %v", token.Error())
	}
	defer client.Disconnect(250)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Subscriber running. Press Ctrl+C to exit.")
	<-sigChan

	log.Println("\nShutting down...")
	log.Printf("Total alerts generated: %d", len(alerts))
}

func handleMessage(c mqtt.Client, m mqtt.Message) {
	var data SensorData
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	fmt.Printf("\n[%s] Device: %s\n", timeNow(), data.DeviceID)
	fmt.Printf("  Temperature: %.2f°C\n", data.Temperature)
	fmt.Printf("  Humidity:    %.2f%%\n", data.Humidity)
	fmt.Printf("  Battery:     %d%%\n", data.Battery)

	checkThresholds(data)
}

func handleStatus(c mqtt.Client, m mqtt.Message) {
	deviceID := extractDeviceID(m.Topic())
	status := string(m.Payload())
	log.Printf("[STATUS] Device %s: %s", deviceID, status)
}

func checkThresholds(data SensorData) {
	if data.Temperature > temperatureThreshold {
		alert := Alert{
			DeviceID: data.DeviceID,
			Type:     "HIGH_TEMPERATURE",
			Value:    data.Temperature,
			Message:  fmt.Sprintf("Temperature %.2f°C exceeds threshold %.2f°C", data.Temperature, temperatureThreshold),
		}
		alerts = append(alerts, alert)
		log.Printf("  ⚠️  ALERT: %s", alert.Message)
	}

	if data.Humidity > humidityThreshold {
		alert := Alert{
			DeviceID: data.DeviceID,
			Type:     "HIGH_HUMIDITY",
			Value:    data.Humidity,
			Message:  fmt.Sprintf("Humidity %.2f%% exceeds threshold %.2f%%", data.Humidity, humidityThreshold),
		}
		alerts = append(alerts, alert)
		log.Printf("  ⚠️  ALERT: %s", alert.Message)
	}

	if data.Battery < 20 {
		alert := Alert{
			DeviceID: data.DeviceID,
			Type:     "LOW_BATTERY",
			Value:    float64(data.Battery),
			Message:  fmt.Sprintf("Battery level %d%% is critically low", data.Battery),
		}
		alerts = append(alerts, alert)
		log.Printf("  ⚠️  ALERT: %s", alert.Message)
	}
}

func extractDeviceID(topic string) string {
	parts := stringsSplit(topic, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

func timeNow() string {
	return time.Now().Format("15:04:05")
}

func stringsSplit(s, sep string) []string {
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
