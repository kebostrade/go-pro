package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	username = flag.String("username", "GoClient", "username for chat")
	room     = flag.String("room", "general", "room to join")
)

// Message represents a chat message
type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func main() {
	flag.Parse()

	// Setup interrupt handler
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to WebSocket
	url := fmt.Sprintf("ws://%s/ws?username=%s&room=%s", *addr, *username, *room)
	log.Printf("Connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	log.Printf("✓ Connected as %s in room %s", *username, *room)
	log.Println("Type messages and press Enter to send. Ctrl+C to exit.")

	// Channel for done signal
	done := make(chan struct{})

	// Start goroutine to read messages
	go func() {
		defer close(done)
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}

			// Display message
			timestamp := msg.Timestamp.Format("15:04:05")
			if msg.Type == "system" {
				fmt.Printf("[%s] 🔔 %s\n", timestamp, msg.Content)
			} else {
				fmt.Printf("[%s] %s: %s\n", timestamp, msg.Username, msg.Content)
			}
		}
	}()

	// Start goroutine to read from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}

			msg := Message{
				Type:    "message",
				Content: text,
			}

			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()

	// Wait for interrupt or done
	select {
	case <-done:
		log.Println("Connection closed")
	case <-interrupt:
		log.Println("\nInterrupt received, closing connection...")

		// Cleanly close the connection
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}

		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}
