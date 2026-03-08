package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

func main() {
	// Connect to NATS with options
	nc, err := nats.Connect(
		nats.DefaultURL,
		nats.Name("NATS-Patterns-Example"),
		nats.Timeout(10*time.Second),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("=== NATS Patterns Demo ===\n")

	// Run all pattern examples
	demoBasicPubSub(nc)
	demoQueueGroups(nc)
	demoWildcards(nc)
	demoRequestReply(nc)
	demoReplySubject(nc)

	fmt.Println("\n=== All demos completed ===")
	fmt.Println("\nNote: Run individual examples in separate terminals:")
	fmt.Println("  - cd publisher && go run main.go")
	fmt.Println("  - cd subscriber && go run main.go")
	fmt.Println("  - cd queue-group && go run worker.go worker-1")
	fmt.Println("  - cd req-reply && go run responder.go")
}

// 1. Basic Publish/Subscribe
func demoBasicPubSub(nc *nats.Conn) {
	fmt.Println("1. Basic Pub/Sub Pattern")
	fmt.Println("---")

	// Subscribe
	sub, err := nc.Subscribe("demo.basic", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
	})
	if err != nil {
		log.Printf("Subscribe error: %v", err)
		return
	}
	defer sub.Unsubscribe()

	// Publish
	msg := Message{
		Content: "Hello from NATS!",
		Time:    time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(msg)

	err = nc.Publish("demo.basic", data)
	if err != nil {
		log.Printf("Publish error: %v", err)
		return
	}

	nc.Flush()
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// 2. Queue Groups
func demoQueueGroups(nc *nats.Conn) {
	fmt.Println("2. Queue Group Pattern")
	fmt.Println("---")

	// Create 3 subscribers in same queue group
	for i := 1; i <= 3; i++ {
		workerID := fmt.Sprintf("worker-%d", i)
		nc.QueueSubscribe("demo.tasks", "task-workers", func(m *nats.Msg) {
			fmt.Printf("%s processing: %s\n", workerID, string(m.Data))
		})
	}

	// Publish 5 tasks
	for i := 1; i <= 5; i++ {
		task := fmt.Sprintf("Task #%d", i)
		nc.Publish("demo.tasks", []byte(task))
	}

	nc.Flush()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\nNote: Each task processed by ONE worker only")
	fmt.Println()
}

// 3. Wildcard Subscriptions
func demoWildcards(nc *nats.Conn) {
	fmt.Println("3. Wildcard Subscriptions")
	fmt.Println("---")

	// Subscribe to all orders
	subAll, _ := nc.Subscribe("orders.>", func(m *nats.Msg) {
		fmt.Printf("ALL - Subject: %s, Data: %s\n", m.Subject, string(m.Data))
	})
	defer subAll.Unsubscribe()

	// Subscribe to Europe orders only
	subEU, _ := nc.Subscribe("orders.europe.>", func(m *nats.Msg) {
		fmt.Printf("EU  - Subject: %s, Data: %s\n", m.Subject, string(m.Data))
	})
	defer subEU.Unsubscribe()

	// Publish to different subjects
	nc.Publish("orders.new", []byte("Order #1"))
	nc.Publish("orders.europe.de", []byte("Order #2"))
	nc.Publish("orders.europe.fr", []byte("Order #3"))
	nc.Publish("orders.asia.jp", []byte("Order #4"))

	nc.Flush()
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// 4. Request/Reply
func demoRequestReply(nc *nats.Conn) {
	fmt.Println("4. Request/Reply Pattern")
	fmt.Println("---")

	// Responder
	nc.Subscribe("demo.echo", func(m *nats.Msg) {
		response := fmt.Sprintf("Echo: %s", string(m.Data))
		m.Respond([]byte(response))
	})

	// Requester
	msg, err := nc.Request("demo.echo", []byte("Hello!"), 1*time.Second)
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}

	fmt.Printf("Request: Hello!")
	fmt.Printf("Response: %s\n", string(msg.Data))
	fmt.Println()
}

// 5. Reply Subjects
func demoReplySubject(nc *nats.Conn) {
	fmt.Println("5. Reply Subject Pattern")
	fmt.Println("---")

	// Create inbox for replies
	inbox := nats.NewInbox()

	// Subscribe to inbox for replies
	sub, _ := nc.Subscribe(inbox, func(m *nats.Msg) {
		fmt.Printf("Async reply received: %s\n", string(m.Data))
	})
	defer sub.Unsubscribe()

	// Publish request with reply subject
	nc.Subscribe("demo.process", func(m *nats.Msg) {
		fmt.Printf("Processing: %s\n", string(m.Data))
		time.Sleep(50 * time.Millisecond)
		nc.Publish(m.Reply, []byte("Processed!"))
	})

	nc.PublishRequest("demo.process", inbox, []byte("Data"))

	nc.Flush()
	time.Sleep(200 * time.Millisecond)
	fmt.Println()
}

// Helper: Connection status
func printConnectionStatus(nc *nats.Conn) {
	fmt.Printf("Status: %s\n", nc.Status())
	fmt.Printf("Connected: %v\n", nc.IsConnected())
	fmt.Printf("Servers: %v\n", nc.Servers())
	stats := nc.Stats()
	fmt.Printf("Messages Sent: %d\n", stats.OutMsgs)
	fmt.Printf("Messages Received: %d\n", stats.InMsgs)
}
