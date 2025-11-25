# 🚀 Quick Start Guide - WebSocket Chat

Get up and running with the WebSocket chat application in 5 minutes!

## ⚡ 1-Minute Setup

```bash
# Navigate to project
cd basic/projects/websocket-chat

# Install dependencies
make deps

# Run the server
make run
```

**Open your browser**: http://localhost:8080

That's it! 🎉

## 📱 Using the Web Interface

1. **Enter your username** (e.g., "Alice")
2. **Enter a room name** (e.g., "general")
3. **Click "Join Chat"**
4. **Start messaging!**

### Test with Multiple Users

Open 3 browser windows:
- Window 1: Username "Alice", Room "general"
- Window 2: Username "Bob", Room "general"  
- Window 3: Username "Charlie", Room "tech"

Alice and Bob can chat in "general", Charlie is in a separate "tech" room.

## 💻 Command Line Client

```bash
# Terminal 1: Start server
make run

# Terminal 2: Connect as Alice
go run examples/client.go -username Alice -room general

# Terminal 3: Connect as Bob
go run examples/client.go -username Bob -room general
```

Type messages and press Enter to send!

## 🔧 Common Commands

```bash
# Build
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench

# Clean build artifacts
make clean

# View help
make help
```

## 🌐 API Endpoints

### WebSocket Connection
```
ws://localhost:8080/ws?username=<name>&room=<room>
```

### REST API
```bash
# List all rooms
curl http://localhost:8080/api/rooms

# Get room stats
curl http://localhost:8080/api/rooms/general
```

## 📊 Message Format

### Send Message
```json
{
  "type": "message",
  "content": "Hello, World!"
}
```

### Receive Message
```json
{
  "type": "message",
  "username": "Alice",
  "content": "Hello, World!",
  "room_id": "general",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### System Message
```json
{
  "type": "system",
  "content": "Alice joined the chat",
  "room_id": "general",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🐳 Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Or use docker-compose
docker-compose up
```

## 🧪 Testing

```bash
# Run all tests
make test

# Expected output:
# ✓ TestHubRegister
# ✓ TestHubUnregister
# ✓ TestHubBroadcast
# ✓ TestHubMultipleRooms
# ✓ TestGetAllRooms
# ✓ TestGetRoomStats
# PASS
```

## 🎯 Quick Challenges

Try these to learn the system:

1. **Multi-Room Chat**: Create 3 rooms and switch between them
2. **Load Test**: Open 10+ browser windows and chat simultaneously
3. **Message History**: Join a room with existing messages
4. **Connection Test**: Close and reopen browser to test reconnection
5. **API Test**: Use curl to monitor room statistics

## 🐛 Troubleshooting

### Server won't start
```bash
# Check if port 8080 is in use
lsof -i :8080

# Use different port
./bin/chat-server -addr :8081
```

### Can't connect
```bash
# Check server is running
curl http://localhost:8080/

# Check WebSocket endpoint
wscat -c ws://localhost:8080/ws?username=test&room=general
```

### Messages not appearing
- Verify you're in the same room
- Check browser console for errors
- Refresh the page

## 📚 Next Steps

1. **Read the full tutorial**: [Tutorial 12 in TUTORIALS.md](../../docs/TUTORIALS.md#-tutorial-12-websocket-real-time-communication)
2. **Explore the code**: Start with `cmd/main.go`
3. **Customize**: Add features like private messaging
4. **Deploy**: Use Docker for production deployment

## 🔗 Resources

- [Full README](README.md)
- [Tutorial 12](../../docs/TUTORIALS.md#-tutorial-12-websocket-real-time-communication)
- [Gorilla WebSocket Docs](https://pkg.go.dev/github.com/gorilla/websocket)
- [WebSocket RFC](https://tools.ietf.org/html/rfc6455)

---

**Happy Chatting! 💬**

