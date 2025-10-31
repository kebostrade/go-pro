package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Client interface defines the methods a client must implement
type Client interface {
	GetID() string
	GetUsername() string
	GetRoomID() string
	GetSendChannel() chan []byte
}

// ClientImpl wraps the actual client to implement the Client interface
type ClientImpl struct {
	ID       string
	Username string
	RoomID   string
	Send     chan []byte
}

func (c *ClientImpl) GetID() string              { return c.ID }
func (c *ClientImpl) GetUsername() string        { return c.Username }
func (c *ClientImpl) GetRoomID() string          { return c.RoomID }
func (c *ClientImpl) GetSendChannel() chan []byte { return c.Send }

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients organized by room
	rooms map[string]map[Client]bool

	// Inbound messages from the clients
	broadcast chan *BroadcastMessage

	// Register requests from the clients
	register chan Client

	// Unregister requests from clients
	unregister chan Client

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Message history per room (last 100 messages)
	history map[string][]Message
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	Message []byte
	RoomID  string
}

// Message represents a chat message
type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
}

// SystemMessage represents a system notification
type SystemMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
}

// New creates a new Hub instance
func New() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan Client),
		unregister: make(chan Client),
		rooms:      make(map[string]map[Client]bool),
		history:    make(map[string][]Message),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client Client) {
	h.unregister <- client
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(message []byte, roomID string) {
	h.broadcast <- &BroadcastMessage{
		Message: message,
		RoomID:  roomID,
	}
}

// registerClient handles client registration
func (h *Hub) registerClient(client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := client.GetRoomID()

	// Create room if it doesn't exist
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[Client]bool)
		h.history[roomID] = make([]Message, 0, 100)
	}

	// Add client to room
	h.rooms[roomID][client] = true

	log.Printf("Client %s (%s) joined room %s. Total clients in room: %d",
		client.GetID(), client.GetUsername(), roomID, len(h.rooms[roomID]))

	// Send join notification
	h.sendSystemMessage(roomID, client.GetUsername()+" joined the chat")

	// Send message history to new client
	h.sendHistory(client)
}

// unregisterClient handles client disconnection
func (h *Hub) unregisterClient(client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := client.GetRoomID()

	if clients, ok := h.rooms[roomID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.GetSendChannel())

			log.Printf("Client %s (%s) left room %s. Remaining clients: %d",
				client.GetID(), client.GetUsername(), roomID, len(clients))

			// Send leave notification
			h.sendSystemMessage(roomID, client.GetUsername()+" left the chat")

			// Clean up empty rooms
			if len(clients) == 0 {
				delete(h.rooms, roomID)
				delete(h.history, roomID)
				log.Printf("Room %s is now empty and has been removed", roomID)
			}
		}
	}
}

// broadcastMessage sends a message to all clients in a room
func (h *Hub) broadcastMessage(bm *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Parse and store message in history
	var msg Message
	if err := json.Unmarshal(bm.Message, &msg); err == nil {
		h.addToHistory(bm.RoomID, msg)
	}

	// Broadcast to all clients in the room
	if clients, ok := h.rooms[bm.RoomID]; ok {
		for client := range clients {
			select {
			case client.GetSendChannel() <- bm.Message:
			default:
				// Client's send channel is full, close it
				close(client.GetSendChannel())
				delete(clients, client)
			}
		}
	}
}

// sendSystemMessage sends a system notification to all clients in a room
func (h *Hub) sendSystemMessage(roomID, content string) {
	sysMsg := SystemMessage{
		Type:      "system",
		Content:   content,
		RoomID:    roomID,
		Timestamp: time.Now(),
	}

	message, err := json.Marshal(sysMsg)
	if err != nil {
		log.Printf("error marshaling system message: %v", err)
		return
	}

	// Broadcast without locking (already locked by caller)
	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			select {
			case client.GetSendChannel() <- message:
			default:
				close(client.GetSendChannel())
				delete(clients, client)
			}
		}
	}
}

// sendHistory sends message history to a client
func (h *Hub) sendHistory(client Client) {
	roomID := client.GetRoomID()
	if history, ok := h.history[roomID]; ok && len(history) > 0 {
		for _, msg := range history {
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			select {
			case client.GetSendChannel() <- msgBytes:
			default:
				// Channel full, skip history
				return
			}
		}
	}
}

// addToHistory adds a message to room history (keeps last 100 messages)
func (h *Hub) addToHistory(roomID string, msg Message) {
	if history, ok := h.history[roomID]; ok {
		// Keep only last 100 messages
		if len(history) >= 100 {
			history = history[1:]
		}
		h.history[roomID] = append(history, msg)
	}
}

// GetRoomStats returns statistics for a room
func (h *Hub) GetRoomStats(roomID string) map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[string]interface{})
	if clients, ok := h.rooms[roomID]; ok {
		stats["client_count"] = len(clients)
		stats["message_count"] = len(h.history[roomID])

		usernames := make([]string, 0, len(clients))
		for client := range clients {
			usernames = append(usernames, client.GetUsername())
		}
		stats["users"] = usernames
	} else {
		stats["client_count"] = 0
		stats["message_count"] = 0
		stats["users"] = []string{}
	}

	return stats
}

// GetAllRooms returns a list of all active rooms
func (h *Hub) GetAllRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rooms := make([]string, 0, len(h.rooms))
	for roomID := range h.rooms {
		rooms = append(rooms, roomID)
	}
	return rooms
}

