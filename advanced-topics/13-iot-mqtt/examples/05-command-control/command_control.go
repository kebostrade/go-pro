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

type Command struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type Response struct {
	DeviceID  string                 `json:"device_id"`
	Command   string                 `json:"command"`
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

type DeviceInfo struct {
	DeviceID   string
	DeviceType string
	Location   string
	Status     string
	LastSeen   time.Time
}

type CommandController struct {
	client  mqtt.Client
	devices map[string]DeviceInfo
}

func NewCommandController(broker string) *CommandController {
	return &CommandController{
		devices: make(map[string]DeviceInfo),
	}
}

func (c *CommandController) Connect(broker string) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("command-controller")
	opts.SetCleanSession(false)

	opts.OnConnect = func(client mqtt.Client) {
		log.Println("Command controller connected")
		client.Subscribe("devices/+/status", 1, c.handleStatus)
		client.Subscribe("devices/+/response", 1, c.handleResponse)
		client.Subscribe("devices/+/telemetry", 1, c.handleTelemetry)
		log.Println("Subscribed to device topics")
	}

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *CommandController) handleStatus(client mqtt.Client, msg mqtt.Message) {
	deviceID := extractDeviceID(msg.Topic())
	status := string(msg.Payload())

	info, exists := c.devices[deviceID]
	if !exists {
		info = DeviceInfo{DeviceID: deviceID}
	}
	info.Status = status
	info.LastSeen = time.Now()
	c.devices[deviceID] = info

	log.Printf("[STATUS] %s: %s", deviceID, status)
}

func (c *CommandController) handleResponse(client mqtt.Client, msg mqtt.Message) {
	deviceID := extractDeviceID(msg.Topic())

	var resp Response
	if err := json.Unmarshal(msg.Payload(), &resp); err != nil {
		log.Printf("Failed to parse response from %s: %v", deviceID, err)
		return
	}

	status := "✓"
	if !resp.Success {
		status = "✗"
	}
	log.Printf("[RESPONSE] %s: %s %s - %s", deviceID, status, resp.Command, resp.Message)
}

func (c *CommandController) handleTelemetry(client mqtt.Client, msg mqtt.Message) {
	deviceID := extractDeviceID(msg.Topic())

	var data map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		return
	}

	info, exists := c.devices[deviceID]
	if !exists {
		info = DeviceInfo{DeviceID: deviceID}
	}
	if loc, ok := data["location"].(string); ok {
		info.Location = loc
	}
	if t, ok := data["type"].(string); ok {
		info.DeviceType = t
	}
	info.LastSeen = time.Now()
	info.Status = "online"
	c.devices[deviceID] = info
}

func (c *CommandController) SendCommand(deviceID string, cmd Command) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("devices/%s/commands", deviceID)
	token := c.client.Publish(topic, 1, false, payload)
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("[COMMAND] Sent to %s: %s", deviceID, cmd.Action)
	return nil
}

func (c *CommandController) BroadcastCommand(cmd Command, deviceType string) {
	for deviceID, info := range c.devices {
		if deviceType == "" || info.DeviceType == deviceType {
			c.SendCommand(deviceID, cmd)
		}
	}
}

func (c *CommandController) ListDevices() []DeviceInfo {
	var list []DeviceInfo
	for _, info := range c.devices {
		list = append(list, info)
	}
	return list
}

func (c *CommandController) GetOnlineDevices() []DeviceInfo {
	var list []DeviceInfo
	for _, info := range c.devices {
		if info.Status == "online" {
			list = append(list, info)
		}
	}
	return list
}

func (c *CommandController) PingDevice(deviceID string) {
	c.SendCommand(deviceID, Command{Action: "ping"})
}

func (c *CommandController) SetInterval(deviceID string, seconds int) {
	c.SendCommand(deviceID, Command{
		Action: "set_interval",
		Params: map[string]interface{}{"interval": seconds},
	})
}

func (c *CommandController) RebootDevice(deviceID string) {
	c.SendCommand(deviceID, Command{Action: "reboot"})
}

func (c *CommandController) Disconnect() {
	c.client.Disconnect(250)
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

func printHelp() {
	fmt.Println("\nCommand Controller - Interactive Mode")
	fmt.Println("======================================")
	fmt.Println("Commands:")
	fmt.Println("  list              - List all known devices")
	fmt.Println("  online            - List online devices")
	fmt.Println("  ping <device_id>  - Ping a specific device")
	fmt.Println("  interval <device_id> <seconds> - Set device interval")
	fmt.Println("  reboot <device_id> - Reboot a device")
	fmt.Println("  broadcast <action> - Send command to all devices")
	fmt.Println("  help              - Show this help")
	fmt.Println("  quit              - Exit the controller")
}

func main() {
	controller := NewCommandController("tcp://localhost:1883")

	if err := controller.Connect("tcp://localhost:1883"); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer controller.Disconnect()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	fmt.Println("\n╔══════════════════════════════════════════════╗")
	fmt.Println("║        IoT Command Controller                ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	printHelp()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\nShutting down...")
			return

		case <-ticker.C:
			devices := controller.ListDevices()
			if len(devices) > 0 {
				fmt.Printf("\n[Monitor] %d devices known\n", len(devices))
			}

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
