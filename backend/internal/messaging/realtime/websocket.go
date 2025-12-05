// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package realtime provides real-time communication functionality.
package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RealtimeEvent represents an event sent over WebSocket.
type RealtimeEvent struct {
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
	UserID    string      `json:"user_id,omitempty"`
}

// Client represents a WebSocket client connection.
type Client struct {
	ID       string
	UserID   string
	Conn     *websocket.Conn
	Send     chan *RealtimeEvent
	Done     chan struct{}
	mu       sync.Mutex
}

// Hub manages WebSocket connections.
type Hub struct {
	mu          sync.RWMutex
	clients     map[string]*Client
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *RealtimeEvent
	userClients map[string][]*Client // Map user IDs to client connections
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		register:    make(chan *Client, 100),
		unregister:  make(chan *Client, 100),
		broadcast:   make(chan *RealtimeEvent, 1000),
		userClients: make(map[string][]*Client),
	}
}

// RegisterClient registers a new WebSocket client.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a WebSocket client.
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastEvent broadcasts an event to all connected clients.
func (h *Hub) BroadcastEvent(event *RealtimeEvent) {
	h.broadcast <- event
}

// BroadcastEventToUser sends an event to all clients of a specific user.
func (h *Hub) BroadcastEventToUser(userID string, event *RealtimeEvent) {
	h.mu.RLock()
	clients, ok := h.userClients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- event:
		default:
			log.Printf("Client send channel full for user %s", userID)
		}
	}
}

// Start starts the hub event loop.
func (h *Hub) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
			h.mu.Unlock()
			log.Printf("Client registered: %s (user: %s)", client.ID, client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)

				// Remove from user clients
				if clients, ok := h.userClients[client.UserID]; ok {
					for i, c := range clients {
						if c.ID == client.ID {
							h.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
							if len(h.userClients[client.UserID]) == 0 {
								delete(h.userClients, client.UserID)
							}
							break
						}
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s", client.ID)

		case event := <-h.broadcast:
			h.mu.RLock()
			clients := make([]*Client, 0, len(h.clients))
			for _, client := range h.clients {
				clients = append(clients, client)
			}
			h.mu.RUnlock()

			for _, client := range clients {
				select {
				case client.Send <- event:
				default:
					log.Printf("Client send channel full for client %s", client.ID)
				}
			}
		}
	}
}

// ClientHandler handles individual WebSocket connections.
type ClientHandler struct {
	hub *Hub
}

// NewClientHandler creates a new client handler.
func NewClientHandler(hub *Hub) *ClientHandler {
	return &ClientHandler{hub: hub}
}

// HandleConnection handles a WebSocket connection.
func (ch *ClientHandler) HandleConnection(conn *websocket.Conn, userID string) error {
	client := &Client{
		ID:     fmt.Sprintf("client_%d", time.Now().UnixNano()),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan *RealtimeEvent, 256),
		Done:   make(chan struct{}),
	}

	ch.hub.RegisterClient(client)

	// Start reading from client
	go ch.readPump(client)
	// Start writing to client
	go ch.writePump(client)

	return nil
}

// readPump reads messages from the WebSocket connection.
func (ch *ClientHandler) readPump(client *Client) {
	defer func() {
		ch.hub.UnregisterClient(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg json.RawMessage
		if err := client.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Echo message back with timestamp
		event := &RealtimeEvent{
			Type:      "echo",
			Timestamp: time.Now().Unix(),
			Data:      msg,
			UserID:    client.UserID,
		}

		select {
		case client.Send <- event:
		default:
			log.Printf("Client send channel full")
		}
	}
}

// writePump writes messages to the WebSocket connection.
func (ch *ClientHandler) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case <-client.Done:
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return

		case event := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteJSON(event); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("WebSocket ping error: %v", err)
				return
			}
		}
	}
}

// DashboardAdapter adapts Kafka events to WebSocket events.
type DashboardAdapter struct {
	hub    *Hub
	ctx    context.Context
	cancel context.CancelFunc
}

// NewDashboardAdapter creates a new dashboard adapter.
func NewDashboardAdapter(hub *Hub) *DashboardAdapter {
	ctx, cancel := context.WithCancel(context.Background())

	return &DashboardAdapter{
		hub:    hub,
		ctx:    ctx,
		cancel: cancel,
	}
}

// BroadcastAnalyticsUpdate broadcasts analytics metrics to all connected clients.
func (da *DashboardAdapter) BroadcastAnalyticsUpdate(metrics interface{}) {
	event := &RealtimeEvent{
		Type:      "analytics_update",
		Timestamp: time.Now().Unix(),
		Data:      metrics,
	}

	da.hub.BroadcastEvent(event)
}

// SendUserNotification sends a notification to a specific user.
func (da *DashboardAdapter) SendUserNotification(userID string, title, message string, priority string) {
	event := &RealtimeEvent{
		Type:      "notification",
		Timestamp: time.Now().Unix(),
		UserID:    userID,
		Data: map[string]interface{}{
			"title":    title,
			"message":  message,
			"priority": priority,
		},
	}

	da.hub.BroadcastEventToUser(userID, event)
}

// BroadcastCourseUpdate broadcasts course progress updates.
func (da *DashboardAdapter) BroadcastCourseUpdate(courseID string, metrics interface{}) {
	event := &RealtimeEvent{
		Type:      "course_update",
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"course_id": courseID,
			"metrics":   metrics,
		},
	}

	da.hub.BroadcastEvent(event)
}

// BroadcastProgressUpdate broadcasts student progress updates.
func (da *DashboardAdapter) BroadcastProgressUpdate(userID, lessonID string, score int, completed bool) {
	event := &RealtimeEvent{
		Type:      "progress_update",
		Timestamp: time.Now().Unix(),
		UserID:    userID,
		Data: map[string]interface{}{
			"lesson_id": lessonID,
			"score":     score,
			"completed": completed,
		},
	}

	da.hub.BroadcastEventToUser(userID, event)
}

// Close closes the adapter.
func (da *DashboardAdapter) Close() error {
	da.cancel()
	return nil
}
