package behavioral

import (
	"fmt"
	"sync"
)

/*
OBSERVER PATTERN

Purpose: Define a one-to-many dependency between objects so that when one object changes state,
all its dependents are notified and updated automatically.

Use Cases:
- Event handling systems
- Pub/Sub messaging
- Model-View updates
- Real-time notifications

Go-Specific Implementation:
- Channel-based observers
- Interface-based observers
- Thread-safe with mutexes
*/

// Observer interface
type Observer interface {
	Update(event string, data interface{})
	GetID() string
}

// Subject interface
type Subject interface {
	Attach(observer Observer)
	Detach(observerID string)
	Notify(event string, data interface{})
}

// EmailObserver observes events and sends emails
type EmailObserver struct {
	ID    string
	Email string
}

func (e *EmailObserver) Update(event string, data interface{}) {
	fmt.Printf("📧 Email to %s: Event '%s' occurred with data: %v\n", 
		e.Email, event, data)
}

func (e *EmailObserver) GetID() string {
	return e.ID
}

// SMSObserver observes events and sends SMS
type SMSObserver struct {
	ID          string
	PhoneNumber string
}

func (s *SMSObserver) Update(event string, data interface{}) {
	fmt.Printf("📱 SMS to %s: Event '%s' occurred with data: %v\n", 
		s.PhoneNumber, event, data)
}

func (s *SMSObserver) GetID() string {
	return s.ID
}

// LogObserver logs events
type LogObserver struct {
	ID string
}

func (l *LogObserver) Update(event string, data interface{}) {
	fmt.Printf("📝 LOG: Event '%s' occurred with data: %v\n", event, data)
}

func (l *LogObserver) GetID() string {
	return l.ID
}

// EventManager is the subject that manages observers
type EventManager struct {
	observers map[string]Observer
	mu        sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		observers: make(map[string]Observer),
	}
}

func (e *EventManager) Attach(observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observers[observer.GetID()] = observer
	fmt.Printf("✅ Observer %s attached\n", observer.GetID())
}

func (e *EventManager) Detach(observerID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.observers, observerID)
	fmt.Printf("❌ Observer %s detached\n", observerID)
}

func (e *EventManager) Notify(event string, data interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	fmt.Printf("\n🔔 Notifying %d observers about event: %s\n", len(e.observers), event)
	for _, observer := range e.observers {
		observer.Update(event, data)
	}
}

// Channel-based Observer Pattern
type ChannelObserver struct {
	ID      string
	Channel chan Event
}

type Event struct {
	Type string
	Data interface{}
}

func NewChannelObserver(id string) *ChannelObserver {
	return &ChannelObserver{
		ID:      id,
		Channel: make(chan Event, 10),
	}
}

func (c *ChannelObserver) Start() {
	go func() {
		for event := range c.Channel {
			fmt.Printf("📢 Observer %s received event: %s with data: %v\n", 
				c.ID, event.Type, event.Data)
		}
	}()
}

func (c *ChannelObserver) Stop() {
	close(c.Channel)
}

// ChannelEventManager manages channel-based observers
type ChannelEventManager struct {
	observers map[string]*ChannelObserver
	mu        sync.RWMutex
}

func NewChannelEventManager() *ChannelEventManager {
	return &ChannelEventManager{
		observers: make(map[string]*ChannelObserver),
	}
}

func (c *ChannelEventManager) Subscribe(observer *ChannelObserver) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.observers[observer.ID] = observer
	observer.Start()
}

func (c *ChannelEventManager) Unsubscribe(observerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if observer, exists := c.observers[observerID]; exists {
		observer.Stop()
		delete(c.observers, observerID)
	}
}

func (c *ChannelEventManager) Publish(event Event) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	for _, observer := range c.observers {
		select {
		case observer.Channel <- event:
		default:
			fmt.Printf("⚠️  Observer %s channel full, skipping event\n", observer.ID)
		}
	}
}

