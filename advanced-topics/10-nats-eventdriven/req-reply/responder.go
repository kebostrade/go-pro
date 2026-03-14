package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock database
var users = map[int]User{
	1: {ID: 1, Name: "Alice Johnson", Email: "alice@example.com"},
	2: {ID: 2, Name: "Bob Smith", Email: "bob@example.com"},
	3: {ID: 3, Name: "Charlie Brown", Email: "charlie@example.com"},
}

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Request-Reply Responder")
	fmt.Println("Connected to NATS")

	// Subscribe and respond to requests
	_, err = nc.Subscribe("user.get", func(m *nats.Msg) {
		// Parse request
		var userID int
		if _, err := fmt.Sscanf(string(m.Data), "%d", &userID); err != nil {
			errMsg := fmt.Sprintf("Invalid user ID: %s", string(m.Data))
			m.Respond([]byte(errMsg))
			fmt.Printf("Error: %s\n", errMsg)
			return
		}

		fmt.Printf("Received request for user ID: %d\n", userID)

		// Look up user
		user, exists := users[userID]
		if !exists {
			errMsg := fmt.Sprintf("User not found: %d", userID)
			m.Respond([]byte(errMsg))
			fmt.Printf("Error: %s\n", errMsg)
			return
		}

		// Serialize and respond
		response, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to marshal user: %v", err)
			m.Respond([]byte("Internal server error"))
			return
		}

		m.Respond(response)
		fmt.Printf("Sent user data for: %s\n", user.Name)
		fmt.Println("---")
	})
	if err != nil {
		log.Fatal(err)
	}

	// Keep running
	fmt.Println("Waiting for requests...")
	fmt.Println("Run requester.go in another terminal")
	select {}
}
