package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceConfig struct {
	DeviceID      string
	SensorType    string
	Location      string
	PublishPeriod time.Duration
}

type Telemetry struct {
	DeviceID  string                 `json:"device_id"`
	Type      string                 `json:"type"`
	Location  string                 `json:"location"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type Command struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type IoTDevice struct {
	config   DeviceConfig
	client   mqtt.Client
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewIoTDevice(config DeviceConfig) *IoTDevice {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID(config.DeviceID)
	opts.SetCleanSession(false)
	opts.SetWill(
		fmt.Sprintf("devices/%s/status", config.DeviceID),
		"offline",
		1,
		true,
	)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("[%s] Connected to broker", config.DeviceID)
		c.Publish(fmt.Sprintf("devices/%s/status", config.DeviceID), 1, true, "online")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("[%s] Connection lost: %v", config.DeviceID, err)
	}

	return &IoTDevice{
		config:   config,
		client:   mqtt.NewClient(opts),
		stopChan: make(chan struct{}),
	}
}

func (d *IoTDevice) Connect() error {
	if token := d.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	d.subscribeCommands()
	return nil
}

func (d *IoTDevice) subscribeCommands() {
	topic := fmt.Sprintf("devices/%s/commands", d.config.DeviceID)
	d.client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		var cmd Command
		if err := json.Unmarshal(m.Payload(), &cmd); err != nil {
			log.Printf("[%s] Invalid command: %v", d.config.DeviceID, err)
			return
		}
		d.handleCommand(cmd)
	})
	log.Printf("[%s] Subscribed to commands", d.config.DeviceID)
}

func (d *IoTDevice) handleCommand(cmd Command) {
	log.Printf("[%s] Received command: %s", d.config.DeviceID, cmd.Action)

	switch cmd.Action {
	case "set_interval":
		if interval, ok := cmd.Params["interval"].(float64); ok {
			d.config.PublishPeriod = time.Duration(interval) * time.Second
			log.Printf("[%s] Interval updated to %v", d.config.DeviceID, d.config.PublishPeriod)
		}

	case "reboot":
		log.Printf("[%s] Rebooting...", d.config.DeviceID)
		d.Stop()
		time.Sleep(2 * time.Second)
		d.Connect()
		d.Start()

	case "ping":
		d.client.Publish(
			fmt.Sprintf("devices/%s/response", d.config.DeviceID),
			1, false,
			`{"status":"pong","timestamp":`+fmt.Sprintf("%d", time.Now().Unix())+`}`,
		)

	default:
		log.Printf("[%s] Unknown command: %s", d.config.DeviceID, cmd.Action)
	}
}

func (d *IoTDevice) Start() {
	d.wg.Add(1)
	go d.publishLoop()
	log.Printf("[%s] Device started", d.config.DeviceID)
}

func (d *IoTDevice) Stop() {
	close(d.stopChan)
	d.wg.Wait()
	d.client.Publish(fmt.Sprintf("devices/%s/status", d.config.DeviceID), 1, true, "offline")
	d.client.Disconnect(250)
	log.Printf("[%s] Device stopped", d.config.DeviceID)
}

func (d *IoTDevice) publishLoop() {
	defer d.wg.Done()

	ticker := time.NewTicker(d.config.PublishPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			d.publishTelemetry()
		}
	}
}

func (d *IoTDevice) publishTelemetry() {
	telemetry := Telemetry{
		DeviceID:  d.config.DeviceID,
		Type:      d.config.SensorType,
		Location:  d.config.Location,
		Timestamp: time.Now().Unix(),
		Data:      d.generateSensorData(),
	}

	payload, err := json.Marshal(telemetry)
	if err != nil {
		log.Printf("[%s] Failed to marshal telemetry: %v", d.config.DeviceID, err)
		return
	}

	topic := fmt.Sprintf("devices/%s/telemetry", d.config.DeviceID)
	token := d.client.Publish(topic, 1, false, payload)

	go func() {
		token.Wait()
		if token.Error() != nil {
			log.Printf("[%s] Publish failed: %v", d.config.DeviceID, token.Error())
		}
	}()

	log.Printf("[%s] Published telemetry: %v", d.config.DeviceID, telemetry.Data)
}

func (d *IoTDevice) generateSensorData() map[string]interface{} {
	switch d.config.SensorType {
	case "temperature":
		return map[string]interface{}{
			"temperature": 20 + rand.Float64()*15,
			"humidity":    40 + rand.Float64()*40,
		}
	case "motion":
		return map[string]interface{}{
			"motion":  rand.Intn(2) == 1,
			"count":   rand.Intn(100),
			"battery": 70 + rand.Intn(30),
		}
	case "air_quality":
		return map[string]interface{}{
			"co2":  400 + rand.Float64()*600,
			"pm25": rand.Float64() * 50,
			"voc":  rand.Float64() * 100,
		}
	default:
		return map[string]interface{}{
			"value": rand.Float64() * 100,
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	devices := []*IoTDevice{
		NewIoTDevice(DeviceConfig{
			DeviceID:      "temp-sensor-001",
			SensorType:    "temperature",
			Location:      "living-room",
			PublishPeriod: 3 * time.Second,
		}),
		NewIoTDevice(DeviceConfig{
			DeviceID:      "motion-sensor-001",
			SensorType:    "motion",
			Location:      "entrance",
			PublishPeriod: 2 * time.Second,
		}),
		NewIoTDevice(DeviceConfig{
			DeviceID:      "air-quality-001",
			SensorType:    "air_quality",
			Location:      "bedroom",
			PublishPeriod: 5 * time.Second,
		}),
	}

	for _, device := range devices {
		if err := device.Connect(); err != nil {
			log.Fatalf("Failed to connect device: %v", err)
		}
		device.Start()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\nIoT Devices Running:")
	fmt.Println("  - temp-sensor-001 (temperature, living-room)")
	fmt.Println("  - motion-sensor-001 (motion, entrance)")
	fmt.Println("  - air-quality-001 (air_quality, bedroom)")
	fmt.Println("\nPress Ctrl+C to stop all devices\n")

	<-sigChan

	fmt.Println("\nShutting down all devices...")
	for _, device := range devices {
		device.Stop()
	}
}
