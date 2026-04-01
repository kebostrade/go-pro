package mqtt

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("tcp://localhost:1883", "test-client")

	if client == nil {
		t.Error("NewClient returned nil")
	}

	if client.clientID != "test-client" {
		t.Errorf("clientID = %s, want test-client", client.clientID)
	}

	if client.broker != "tcp://localhost:1883" {
		t.Errorf("broker = %s, want tcp://localhost:1883", client.broker)
	}
}

func TestClientNotConnected(t *testing.T) {
	client := NewClient("tcp://localhost:1883", "test-client")

	if client.IsConnected() {
		t.Error("New client should not be connected")
	}

	// Publishing without connection should fail
	err := client.Publish("test/topic", "payload", 1)
	if err == nil {
		t.Error("Publish should fail when not connected")
	}

	err = client.Subscribe("test/#", 1, func(topic string, payload []byte) {})
	if err == nil {
		t.Error("Subscribe should fail when not connected")
	}
}

func TestSetWill(t *testing.T) {
	client := NewClient("tcp://localhost:1883", "test-client")

	// Setting LWT should not panic
	client.SetWill("device/status", "offline", 1, true)
}
