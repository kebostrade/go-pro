//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Task: Work with JSON data in Go - marshal and unmarshal structs.
// Create a User management system that can convert between Go structs and JSON.

// User represents a user in the system
type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Age      int      `json:"age"`
	Active   bool     `json:"active"`
	Tags     []string `json:"tags,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata contains additional user information
type Metadata struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Role      string `json:"role"`
}

// UserList represents a collection of users
type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// toJSON converts a User to JSON string
func (u *User) toJSON() (string, error) {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// fromJSON creates a User from JSON string
func fromJSON(jsonStr string) (*User, error) {
	var user User
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// filterActiveUsers returns only active users
func filterActiveUsers(users []User) []User {
	var active []User
	for _, user := range users {
		if user.Active {
			active = append(active, user)
		}
	}
	return active
}

// findUserByID finds a user by their ID
func findUserByID(users []User, id int) *User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func main() {
	// Create sample users
	users := []User{
		{
			ID:     1,
			Name:   "Alice Johnson",
			Email:  "alice@example.com",
			Age:    28,
			Active: true,
			Tags:   []string{"developer", "golang", "backend"},
			Metadata: Metadata{
				CreatedAt: "2024-01-15",
				UpdatedAt: "2024-03-20",
				Role:      "Senior Developer",
			},
		},
		{
			ID:     2,
			Name:   "Bob Smith",
			Email:  "bob@example.com",
			Age:    35,
			Active: false,
			Tags:   []string{"manager", "agile"},
			Metadata: Metadata{
				CreatedAt: "2023-06-10",
				UpdatedAt: "2024-02-14",
				Role:      "Project Manager",
			},
		},
		{
			ID:     3,
			Name:   "Carol White",
			Email:  "carol@example.com",
			Age:    31,
			Active: true,
			Tags:   []string{"developer", "frontend", "react"},
			Metadata: Metadata{
				CreatedAt: "2024-02-01",
				UpdatedAt: "2024-03-25",
				Role:      "Frontend Developer",
			},
		},
	}

	// Create user list
	userList := UserList{
		Users: users,
		Total: len(users),
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(userList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User List as JSON:")
	fmt.Println(string(jsonData))

	// Parse JSON back to struct
	var parsedList UserList
	err = json.Unmarshal(jsonData, &parsedList)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n" + "=".repeat(60))
	fmt.Printf("Parsed %d users from JSON\n", parsedList.Total)

	// Filter active users
	activeUsers := filterActiveUsers(parsedList.Users)
	fmt.Printf("\nActive users: %d\n", len(activeUsers))
	for _, user := range activeUsers {
		fmt.Printf("  - %s (%s)\n", user.Name, user.Email)
	}

	// Find specific user
	fmt.Println("\n" + "=".repeat(60))
	userID := 2
	foundUser := findUserByID(parsedList.Users, userID)
	if foundUser != nil {
		fmt.Printf("Found user with ID %d:\n", userID)
		userJSON, _ := foundUser.toJSON()
		fmt.Println(userJSON)
	}

	// Demonstrate parsing from JSON string
	fmt.Println("\n" + "=".repeat(60))
	jsonString := `{
		"id": 4,
		"name": "David Brown",
		"email": "david@example.com",
		"age": 29,
		"active": true,
		"tags": ["developer", "devops", "kubernetes"]
	}`

	newUser, err := fromJSON(jsonString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsed new user from JSON string:")
	fmt.Printf("Name: %s, Email: %s, Tags: %v\n", newUser.Name, newUser.Email, newUser.Tags)
}

// Helper function to repeat strings (Go doesn't have built-in string repeat in older versions)
func (s string) repeat(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
