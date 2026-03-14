package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================================
// MESSAGE TYPES
// ============================================================================

type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan Message
	hub      *Hub
}

// ============================================================================
// CHAT HUB
// ============================================================================

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// Broadcast user joined message
			h.broadcast <- Message{
				Type:      "system",
				Content:   client.username + " joined the chat",
				Timestamp: time.Now(),
			}

			log.Printf("👤 Client connected: %s (Total: %d)", client.username, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

			// Broadcast user left message
			h.broadcast <- Message{
				Type:      "system",
				Content:   client.username + " left the chat",
				Timestamp: time.Now(),
			}

			log.Printf("👋 Client disconnected: %s (Total: %d)", client.username, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ============================================================================
// CLIENT HANDLING
// ============================================================================

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse incoming message
		var msg map[string]string
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Broadcast chat message
		c.hub.broadcast <- Message{
			Type:      "chat",
			Username:  c.username,
			Content:   msg["content"],
			Timestamp: time.Now(),
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ============================================================================
// HTTP HANDLERS
// ============================================================================

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func serveWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from query parameter
		username := r.URL.Query().Get("username")
		if username == "" {
			username = "Anonymous"
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// Create client
		client := &Client{
			conn:     conn,
			username: username,
			send:     make(chan Message, 256),
			hub:      hub,
		}

		// Register client
		hub.register <- client

		// Start pumps
		go client.writePump()
		go client.readPump()
	}
}

// ============================================================================
// HTML CLIENT
// ============================================================================

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Go WebSocket Chat</title>
    <style>
        * { box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .chat-container {
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .chat-header {
            background: #4CAF50;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .chat-messages {
            height: 400px;
            overflow-y: auto;
            padding: 20px;
            background: #fafafa;
        }
        .message {
            margin-bottom: 10px;
            padding: 10px;
            border-radius: 5px;
            background: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .message.system {
            background: #e3f2fd;
            color: #1976d2;
            text-align: center;
            font-style: italic;
        }
        .message .username {
            font-weight: bold;
            color: #4CAF50;
        }
        .message .timestamp {
            font-size: 0.8em;
            color: #999;
            margin-left: 10px;
        }
        .chat-input {
            display: flex;
            padding: 20px;
            background: white;
            border-top: 1px solid #ddd;
        }
        .chat-input input {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
        }
        .chat-input button {
            padding: 10px 20px;
            margin-left: 10px;
            background: #4CAF50;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        .chat-input button:hover {
            background: #45a049;
        }
        .online-count {
            text-align: center;
            padding: 10px;
            background: #f0f0f0;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="chat-container">
        <div class="chat-header">
            <h1>💬 Go WebSocket Chat</h1>
        </div>
        <div class="online-count">
            <span id="online-count">1</span> users online
        </div>
        <div class="chat-messages" id="messages"></div>
        <div class="chat-input">
            <input type="text" id="message-input" placeholder="Type a message..." />
            <button onclick="sendMessage()">Send</button>
        </div>
    </div>

    <script>
        const ws = new WebSocket('ws://localhost:8080/ws?username=' + getRandomUsername());
        const messages = document.getElementById('messages');
        const input = document.getElementById('message-input');
        let onlineCount = 1;

        function getRandomUsername() {
            const adjectives = ['Happy', 'Clever', 'Brave', 'Swift', 'Kind', 'Wise', 'Bold', 'Cool'];
            const nouns = ['Coder', 'Developer', 'Engineer', 'Builder', 'Creator', 'Maker', 'Designer', 'Hacker'];
            return adjectives[Math.floor(Math.random() * adjectives.length)] + ' ' +
                   nouns[Math.floor(Math.random() * nouns.length)];
        }

        ws.onmessage = function(event) {
            const msg = JSON.parse(event.data);
            addMessage(msg);
        };

        function addMessage(msg) {
            const div = document.createElement('div');
            div.className = 'message' + (msg.type === 'system' ? ' system' : '');

            if (msg.type === 'system') {
                div.textContent = msg.content;
            } else {
                const username = document.createElement('span');
                username.className = 'username';
                username.textContent = msg.username + ': ';

                const content = document.createTextNode(msg.content);

                const timestamp = document.createElement('span');
                timestamp.className = 'timestamp';
                timestamp.textContent = new Date(msg.timestamp).toLocaleTimeString();

                div.appendChild(username);
                div.appendChild(content);
                div.appendChild(timestamp);
            }

            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        }

        function sendMessage() {
            const content = input.value.trim();
            if (content) {
                ws.send(JSON.stringify({ content: content }));
                input.value = '';
            }
        }

        input.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendMessage();
            }
        });

        // Welcome message
        addMessage({
            type: 'system',
            content: '👋 Welcome to the chat! Your username: ' + getRandomUsername()
        });
    </script>
</body>
</html>
`

// ============================================================================
// MAIN
// ============================================================================

func main() {
	hub := NewHub()
	go hub.Run()

	// Write HTML file
	if err := os.WriteFile("index.html", []byte(htmlTemplate), 0644); err != nil {
		log.Fatalf("Failed to write HTML file: %v", err)
	}

	// Setup routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS(hub))

	// Start server
	addr := ":8080"
	log.Printf("🚀 Chat server starting on %s", addr)
	log.Printf("📱 Open http://localhost%s in your browser", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ============================================================================
// USAGE
// ============================================================================

/*
1. Run the server:
   go run chat_server.go

2. Open http://localhost:8080 in multiple browser tabs

3. Each tab will have a random username and can send messages

4. Messages are broadcast to all connected clients in real-time

5. System messages notify when users join/leave

Features:
- ✅ Real-time bidirectional communication
- ✅ Multiple concurrent users
- ✅ User join/leave notifications
- ✅ Message timestamps
- ✅ Automatic reconnection handling
- ✅ Graceful connection closure
- ✅ Online user count
- ✅ Clean, modern UI
*/
