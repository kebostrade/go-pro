package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/websocket-chat/internal/client"
	"github.com/DimaJoyti/go-pro/basic/projects/websocket-chat/internal/hub"
	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo purposes
		},
	}
)

// ClientAdapter adapts client.Client to hub.Client interface
type ClientAdapter struct {
	*client.Client
}

func (ca *ClientAdapter) GetID() string              { return ca.Client.ID }
func (ca *ClientAdapter) GetUsername() string        { return ca.Client.Username }
func (ca *ClientAdapter) GetRoomID() string          { return ca.Client.RoomID }
func (ca *ClientAdapter) GetSendChannel() chan []byte { return ca.Client.Send }

// HubAdapter adapts hub.Hub to client.Hub interface
type HubAdapter struct {
	*hub.Hub
}

func (ha *HubAdapter) Unregister(c *client.Client) {
	ha.Hub.Unregister(&ClientAdapter{c})
}

func (ha *HubAdapter) Broadcast(message []byte, roomID string) {
	ha.Hub.Broadcast(message, roomID)
}

func main() {
	flag.Parse()

	// Create hub
	h := hub.New()
	hubAdapter := &HubAdapter{h}

	// Start hub in goroutine
	go h.Run()

	// Setup HTTP routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hubAdapter, h, w, r)
	})
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		getRooms(h, w, r)
	})
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		getRoomStats(h, w, r)
	})

	// Print startup message
	printStartupMessage()

	// Start server
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(homeHTML))
}

func serveWs(hubAdapter *HubAdapter, h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	username := r.URL.Query().Get("username")
	roomID := r.URL.Query().Get("room")

	if username == "" {
		username = "Anonymous"
	}
	if roomID == "" {
		roomID = "general"
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create client
	clientID := fmt.Sprintf("%s-%d", username, time.Now().UnixNano())
	c := client.New(clientID, username, roomID, conn, hubAdapter)
	adapter := &ClientAdapter{c}

	// Register client
	h.Register(adapter)

	// Start client goroutines
	go c.WritePump()
	go c.ReadPump()
}

func getRooms(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms := h.GetAllRooms()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rooms": rooms,
		"count": len(rooms),
	})
}

func getRoomStats(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract room ID from path
	roomID := r.URL.Path[len("/api/rooms/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	stats := h.GetRoomStats(roomID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func printStartupMessage() {
	fmt.Println(`
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║           💬 WebSocket Chat Server                          ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

✓ Server starting on http://localhost` + *addr + `
✓ WebSocket endpoint: ws://localhost` + *addr + `/ws
✓ Ready to accept connections

📚 API Endpoints:
  GET  /                    - Web interface
  WS   /ws?username=X&room=Y - WebSocket connection
  GET  /api/rooms           - List all rooms
  GET  /api/rooms/{id}      - Get room statistics

🎯 Quick Start:
  1. Open http://localhost` + *addr + ` in your browser
  2. Enter your username and room name
  3. Start chatting!

Press Ctrl+C to stop the server
`)
}

const homeHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WebSocket Chat</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 10px; margin-bottom: 20px; text-align: center; }
        .header h1 { font-size: 2em; margin-bottom: 10px; }
        .login-form { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .login-form input { width: 100%; padding: 12px; margin: 10px 0; border: 2px solid #e0e0e0; border-radius: 5px; font-size: 16px; }
        .login-form button { width: 100%; padding: 12px; background: #667eea; color: white; border: none; border-radius: 5px; font-size: 16px; cursor: pointer; margin-top: 10px; }
        .login-form button:hover { background: #5568d3; }
        .chat-container { background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); display: none; }
        .chat-header { background: #667eea; color: white; padding: 20px; border-radius: 10px 10px 0 0; }
        .messages { height: 400px; overflow-y: auto; padding: 20px; border-bottom: 2px solid #e0e0e0; }
        .message { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .message.user { background: #e3f2fd; }
        .message.system { background: #fff3e0; font-style: italic; text-align: center; }
        .message .username { font-weight: bold; color: #667eea; margin-right: 10px; }
        .message .time { font-size: 0.8em; color: #999; margin-left: 10px; }
        .input-area { padding: 20px; display: flex; gap: 10px; }
        .input-area input { flex: 1; padding: 12px; border: 2px solid #e0e0e0; border-radius: 5px; font-size: 16px; }
        .input-area button { padding: 12px 30px; background: #667eea; color: white; border: none; border-radius: 5px; cursor: pointer; }
        .input-area button:hover { background: #5568d3; }
        .status { padding: 10px; text-align: center; font-size: 0.9em; color: #666; }
        .status.connected { color: #4caf50; }
        .status.disconnected { color: #f44336; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💬 WebSocket Chat</h1>
            <p>Real-time messaging with Go and WebSockets</p>
        </div>

        <div class="login-form" id="loginForm">
            <h2>Join a Chat Room</h2>
            <input type="text" id="username" placeholder="Enter your username" value="User">
            <input type="text" id="room" placeholder="Enter room name" value="general">
            <button onclick="connect()">Join Chat</button>
        </div>

        <div class="chat-container" id="chatContainer">
            <div class="chat-header">
                <h2 id="roomName">Room: general</h2>
                <div class="status" id="status">Connecting...</div>
            </div>
            <div class="messages" id="messages"></div>
            <div class="input-area">
                <input type="text" id="messageInput" placeholder="Type a message..." onkeypress="handleKeyPress(event)">
                <button onclick="sendMessage()">Send</button>
            </div>
        </div>
    </div>

    <script>
        let ws;
        let username;
        let room;

        function connect() {
            username = document.getElementById('username').value || 'Anonymous';
            room = document.getElementById('room').value || 'general';

            document.getElementById('loginForm').style.display = 'none';
            document.getElementById('chatContainer').style.display = 'block';
            document.getElementById('roomName').textContent = 'Room: ' + room;

            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws?username=' + encodeURIComponent(username) + '&room=' + encodeURIComponent(room);

            ws = new WebSocket(wsUrl);

            ws.onopen = function() {
                updateStatus('Connected', true);
            };

            ws.onmessage = function(event) {
                const msg = JSON.parse(event.data);
                displayMessage(msg);
            };

            ws.onclose = function() {
                updateStatus('Disconnected', false);
            };

            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
                updateStatus('Error', false);
            };
        }

        function sendMessage() {
            const input = document.getElementById('messageInput');
            const content = input.value.trim();

            if (content && ws && ws.readyState === WebSocket.OPEN) {
                const msg = {
                    type: 'message',
                    content: content
                };
                ws.send(JSON.stringify(msg));
                input.value = '';
            }
        }

        function displayMessage(msg) {
            const messagesDiv = document.getElementById('messages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + msg.type;

            const time = new Date(msg.timestamp).toLocaleTimeString();

            if (msg.type === 'system') {
                messageDiv.innerHTML = msg.content + '<span class="time">' + time + '</span>';
            } else {
                messageDiv.innerHTML = '<span class="username">' + msg.username + ':</span>' + msg.content + '<span class="time">' + time + '</span>';
            }

            messagesDiv.appendChild(messageDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        function updateStatus(text, connected) {
            const status = document.getElementById('status');
            status.textContent = text;
            status.className = 'status ' + (connected ? 'connected' : 'disconnected');
        }

        function handleKeyPress(event) {
            if (event.key === 'Enter') {
                sendMessage();
            }
        }
    </script>
</body>
</html>`

