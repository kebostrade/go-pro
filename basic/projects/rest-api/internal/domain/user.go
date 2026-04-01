// Package domain contains the core domain entities
package domain

// User represents a user entity
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
