# Building Real-time Applications with Go and WebSockets

Create real-time bidirectional communication applications.

## Learning Objectives

- Understand WebSocket protocol
- Implement WebSocket servers
- Handle concurrent connections
- Build chat applications
- Implement real-time updates
- Scale WebSocket servers

## Theory

### Basic WebSocket Server

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("upgrade error: %v", err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("read error: %v", err)
            break
        }

        log.Printf("received: %s", message)

        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Printf("write error: %v", err)
            break
        }
    }
}
```

### Chat Server with Rooms

```go
type ChatServer struct {
    clients   map[*Client]bool
    broadcast chan *Message
    register  chan *Client
    unregister chan *Client
    rooms     map[string]*Room
    mu        sync.RWMutex
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
    room string
}

type Message struct {
    Content string `json:"content"`
    Room    string `json:"room"`
    Sender  string `json:"sender"`
}

func (s *ChatServer) Run() {
    for {
        select {
        case client := <-s.register:
            s.mu.Lock()
            s.clients[client] = true
            s.mu.Unlock()

        case client := <-s.unregister:
            s.mu.Lock()
            if _, ok := s.clients[client]; ok {
                delete(s.clients, client)
                close(client.send)
            }
            s.mu.Unlock()

        case message := <-s.broadcast:
            s.mu.RLock()
            for client := range s.clients {
                if client.room == message.Room {
                    select {
                    case client.send <- []byte(message.Content):
                    default:
                        close(client.send)
                        delete(s.clients, client)
                    }
                }
            }
            s.mu.RUnlock()
        }
    }
}

func (c *Client) readPump(broadcast chan<- *Message) {
    defer func() {
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        broadcast <- &msg
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
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.conn.WriteMessage(websocket.TextMessage, message)

        case <-ticker.C:
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

### Real-time Updates Pattern

```go
type SubscriptionManager struct {
    subscriptions map[string]map[*websocket.Conn]bool
    mu            sync.RWMutex
}

func (sm *SubscriptionManager) Subscribe(topic string, conn *websocket.Conn) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.subscriptions[topic] == nil {
        sm.subscriptions[topic] = make(map[*websocket.Conn]bool)
    }
    sm.subscriptions[topic][conn] = true
}

func (sm *SubscriptionManager) Publish(topic string, message []byte) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    for conn := range sm.subscriptions[topic] {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            conn.Close()
            delete(sm.subscriptions[topic], conn)
        }
    }
}
```

## Security Considerations

```go
func websocketAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.URL.Query().Get("token")
        if token == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func rateLimitConnections(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "too many connections", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Performance Tips

```go
type ConnectionPool struct {
    conns   chan *websocket.Conn
    maxSize int
}

func NewConnectionPool(maxSize int) *ConnectionPool {
    return &ConnectionPool{
        conns:   make(chan *websocket.Conn, maxSize),
        maxSize: maxSize,
    }
}

func (p *ConnectionPool) Acquire() (*websocket.Conn, error) {
    select {
    case conn := <-p.conns:
        return conn, nil
    default:
        return nil, errors.New("connection pool exhausted")
    }
}
```

## Exercises

1. Build a multi-room chat server
2. Implement real-time notifications
3. Add reconnection handling
4. Create a collaborative editor

## Validation

```bash
cd exercises
go test -v ./...
wscat -c ws://localhost:8080/ws
```

## Key Takeaways

- Use channels for concurrent message handling
- Implement ping/pong for connection health
- Handle disconnects gracefully
- Rate limit connections
- Use proper message framing

## Next Steps

**[AT-08: NATS Event-Driven](../AT-08-nats-event-driven/README.md)**

---

Real-time is real powerful. ⚡
