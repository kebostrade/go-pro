package greeting

import (
	"fmt"
	"sync"
)

// Template for greeting messages
const template = "Hello, %s!\n"

// GreetingService provides thread-safe greeting operations
type GreetingService struct {
	mu  sync.Mutex
	msg string
}

// NewGreetingService creates a new GreetingService
func NewGreetingService() *GreetingService {
	return &GreetingService{
		msg: GetGreetingTemplate(),
	}
}

// Greet returns a greeting message repeated n times
func Greet(name string, times int) string {
	if times <= 0 {
		times = 1
	}
	if name == "" {
		name = "World"
	}

	result := ""
	for i := 0; i < times; i++ {
		result += fmt.Sprintf(template, name)
	}
	return result
}

// GetGreetingTemplate returns the greeting template
func GetGreetingTemplate() string {
	return template
}

// GreetWithLock is a thread-safe version of Greet
func GreetWithLock(svc *GreetingService, name string, times int) string {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	return Greet(name, times)
}
