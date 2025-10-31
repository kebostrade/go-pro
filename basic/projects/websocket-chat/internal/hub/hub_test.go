package hub

import (
	"testing"
	"time"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	id       string
	username string
	roomID   string
	send     chan []byte
}

func (m *MockClient) GetID() string               { return m.id }
func (m *MockClient) GetUsername() string         { return m.username }
func (m *MockClient) GetRoomID() string           { return m.roomID }
func (m *MockClient) GetSendChannel() chan []byte { return m.send }

func newMockClient(id, username, roomID string) *MockClient {
	return &MockClient{
		id:       id,
		username: username,
		roomID:   roomID,
		send:     make(chan []byte, 256),
	}
}

func TestHubRegister(t *testing.T) {
	h := New()
	go h.Run()

	client := newMockClient("1", "Alice", "general")
	h.Register(client)

	// Give hub time to process
	time.Sleep(10 * time.Millisecond)

	stats := h.GetRoomStats("general")
	if stats["client_count"] != 1 {
		t.Errorf("Expected 1 client, got %v", stats["client_count"])
	}
}

func TestHubUnregister(t *testing.T) {
	h := New()
	go h.Run()

	client := newMockClient("1", "Alice", "general")
	h.Register(client)
	time.Sleep(10 * time.Millisecond)

	h.Unregister(client)
	time.Sleep(10 * time.Millisecond)

	stats := h.GetRoomStats("general")
	if stats["client_count"] != 0 {
		t.Errorf("Expected 0 clients, got %v", stats["client_count"])
	}
}

func TestHubBroadcast(t *testing.T) {
	h := New()
	go h.Run()

	client1 := newMockClient("1", "Alice", "general")
	client2 := newMockClient("2", "Bob", "general")

	h.Register(client1)
	h.Register(client2)
	time.Sleep(10 * time.Millisecond)

	// Clear system messages from registration
	drainChannel(client1.send)
	drainChannel(client2.send)

	// Broadcast a message
	message := []byte(`{"type":"message","content":"Hello"}`)
	h.Broadcast(message, "general")

	// Both clients should receive the message
	select {
	case msg := <-client1.send:
		if string(msg) != string(message) {
			t.Errorf("Client1 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive message")
	}

	select {
	case msg := <-client2.send:
		if string(msg) != string(message) {
			t.Errorf("Client2 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 did not receive message")
	}
}

// drainChannel drains all messages from a channel
func drainChannel(ch chan []byte) {
	for {
		select {
		case <-ch:
			// Drain
		default:
			return
		}
	}
}

func TestHubMultipleRooms(t *testing.T) {
	h := New()
	go h.Run()

	client1 := newMockClient("1", "Alice", "general")
	client2 := newMockClient("2", "Bob", "tech")

	h.Register(client1)
	h.Register(client2)
	time.Sleep(10 * time.Millisecond)

	// Clear system messages from registration
	drainChannel(client1.send)
	drainChannel(client2.send)

	// Broadcast to general room
	message := []byte(`{"type":"message","content":"Hello general"}`)
	h.Broadcast(message, "general")

	// Only client1 should receive
	select {
	case <-client1.send:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive message")
	}

	// Client2 should not receive
	select {
	case <-client2.send:
		t.Error("Client2 should not receive message from different room")
	case <-time.After(50 * time.Millisecond):
		// Expected
	}
}

func TestGetAllRooms(t *testing.T) {
	h := New()
	go h.Run()

	client1 := newMockClient("1", "Alice", "general")
	client2 := newMockClient("2", "Bob", "tech")
	client3 := newMockClient("3", "Charlie", "random")

	h.Register(client1)
	h.Register(client2)
	h.Register(client3)
	time.Sleep(10 * time.Millisecond)

	rooms := h.GetAllRooms()
	if len(rooms) != 3 {
		t.Errorf("Expected 3 rooms, got %d", len(rooms))
	}
}

func TestGetRoomStats(t *testing.T) {
	h := New()
	go h.Run()

	client1 := newMockClient("1", "Alice", "general")
	client2 := newMockClient("2", "Bob", "general")

	h.Register(client1)
	h.Register(client2)
	time.Sleep(10 * time.Millisecond)

	stats := h.GetRoomStats("general")

	if stats["client_count"] != 2 {
		t.Errorf("Expected 2 clients, got %v", stats["client_count"])
	}

	users := stats["users"].([]string)
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func BenchmarkBroadcast(b *testing.B) {
	h := New()
	go h.Run()

	// Create 100 clients
	clients := make([]*MockClient, 100)
	for i := 0; i < 100; i++ {
		client := newMockClient(string(rune(i)), "User", "general")
		clients[i] = client
		h.Register(client)
	}
	time.Sleep(50 * time.Millisecond)

	message := []byte(`{"type":"message","content":"Benchmark message"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Broadcast(message, "general")
	}
}

func BenchmarkRegister(b *testing.B) {
	h := New()
	go h.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := newMockClient(string(rune(i)), "User", "general")
		h.Register(client)
	}
}
