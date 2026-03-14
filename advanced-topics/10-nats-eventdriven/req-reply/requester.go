package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Request-Reply Requester")
	fmt.Println("Connected to NATS")

	// Request user by ID
	userIDs := []int{1, 2, 3}

	for _, id := range userIDs {
		// Create request
		request := fmt.Sprintf("%d", id)

		fmt.Printf("\nRequesting user ID: %d\n", id)

		// Send request with timeout
		msg, err := nc.Request("user.get", []byte(request), 2*time.Second)
		if err != nil {
			if err == nats.ErrTimeout {
				fmt.Println("Request timed out!")
			} else {
				log.Printf("Request failed: %v", err)
			}
			continue
		}

		// Parse response
		var user User
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			log.Printf("Failed to parse response: %v", err)
			continue
		}

		fmt.Printf("Received user:\n")
		fmt.Printf("  ID: %d\n", user.ID)
		fmt.Printf("  Name: %s\n", user.Name)
		fmt.Printf("  Email: %s\n", user.Email)
	}

	fmt.Println("\nAll requests completed!")
}
