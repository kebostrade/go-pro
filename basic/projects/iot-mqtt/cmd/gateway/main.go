package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goproject/iot-mqtt/internal/mqtt"
	"github.com/goproject/iot-mqtt/internal/processor"
)

func main() {
	log.Println("IoT Gateway Starting...")

	broker := os.Getenv("BROKER_URL")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	// Create MQTT client for gateway
	client := mqtt.NewClient(broker, "gateway")
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	defer client.Disconnect()

	// Create telemetry processor
	tp := processor.NewTelemetryProcessor(client)

	// Setup command handler for sending commands to devices
	go handleCommands(client)

	// Subscribe to all device telemetry
	ctx := context.Background()
	if err := tp.Start(ctx, "devices/+/telemetry"); err != nil {
		log.Printf("Warning: Could not subscribe to telemetry: %v", err)
	}

	// Setup HTTP API
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// List known devices (from storage)
	r.Get("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		// Return mock device list for now
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"devices": []string{},
		})
	})

	// Get device readings
	r.Get("/api/devices/{id}/readings", func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "id")
		readings, err := tp.GetReadings(deviceID, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(readings)
	})

	// Send command to device
	r.Post("/api/devices/{id}/commands", func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "id")

		var cmd struct {
			Action string `json:"action"`
			Value  string `json:"value,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Publish command to device
		cmdPayload, _ := json.Marshal(cmd)
		topic := fmt.Sprintf("devices/%s/commands", deviceID)
		if err := client.Publish(topic, cmdPayload, 1); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "sent",
			"command": cmd.Action,
		})
	})

	// Get alerts
	r.Get("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tp.GetAlerts())
	})

	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Gateway HTTP API listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")
}

func handleCommands(client *mqtt.Client) {
	// This would handle command responses from devices
	// For now, just log any command responses received
	client.Subscribe("devices/+/commands/response", 1, func(topic string, payload []byte) {
		log.Printf("Command response on %s: %s", topic, string(payload))
	})
}
