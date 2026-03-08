package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID(fmt.Sprintf("publisher-%d", rand.Intn(10000)))
	opts.SetCleanSession(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Println("Starting publisher...")
	log.Println("Publishing to: sensors/temperature, sensors/humidity")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	deviceID := "sensor-001"

	for {
		select {
		case <-sigChan:
			log.Println("Shutting down...")
			retained := true
			client.Publish("devices/"+deviceID+"/status", 1, retained, "offline")
			return

		case <-ticker.C:
			publishSensorData(client, deviceID)
		}
	}
}

func publishSensorData(client mqtt.Client, deviceID string) {
	data := SensorData{
		DeviceID:    deviceID,
		Timestamp:   time.Now().Unix(),
		Temperature: 20 + rand.Float64()*15,
		Humidity:    40 + rand.Float64()*40,
		Battery:     75 + rand.Intn(25),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return
	}

	telemetryTopic := fmt.Sprintf("devices/%s/telemetry", deviceID)
	token := client.Publish(telemetryTopic, 1, false, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Failed to publish: %v", token.Error())
		return
	}

	log.Printf("Published: temp=%.2f°C, humidity=%.2f%%, battery=%d%%",
		data.Temperature, data.Humidity, data.Battery)

	client.Publish(fmt.Sprintf("devices/%s/status", deviceID), 1, true, "online")
}
