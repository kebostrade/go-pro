package creational

import (
	"fmt"
)

/*
FACTORY PATTERN

Purpose: Define an interface for creating objects, but let subclasses decide which class to instantiate.

Use Cases:
- Creating different types of database connections
- Payment processors (PayPal, Stripe, etc.)
- Notification systems (Email, SMS, Push)
- Document parsers (JSON, XML, YAML)

Go-Specific Implementation:
- Use interfaces for product types
- Factory function returns interface type
- Type switching or map-based factories
*/

// Notification interface defines the contract for all notifications
type Notification interface {
	Send(message string) error
	GetType() string
}

// EmailNotification implements Notification for email
type EmailNotification struct {
	To      string
	From    string
	Subject string
}

func (e *EmailNotification) Send(message string) error {
	fmt.Printf("📧 Sending EMAIL to %s: %s\n", e.To, message)
	return nil
}

func (e *EmailNotification) GetType() string {
	return "email"
}

// SMSNotification implements Notification for SMS
type SMSNotification struct {
	PhoneNumber string
	Provider    string
}

func (s *SMSNotification) Send(message string) error {
	fmt.Printf("📱 Sending SMS to %s: %s\n", s.PhoneNumber, message)
	return nil
}

func (s *SMSNotification) GetType() string {
	return "sms"
}

// PushNotification implements Notification for push notifications
type PushNotification struct {
	DeviceToken string
	Platform    string
}

func (p *PushNotification) Send(message string) error {
	fmt.Printf("🔔 Sending PUSH to %s (%s): %s\n", p.DeviceToken, p.Platform, message)
	return nil
}

func (p *PushNotification) GetType() string {
	return "push"
}

// NotificationFactory creates notifications based on type
func NotificationFactory(notificationType string) (Notification, error) {
	switch notificationType {
	case "email":
		return &EmailNotification{
			To:      "user@example.com",
			From:    "noreply@app.com",
			Subject: "Notification",
		}, nil
	case "sms":
		return &SMSNotification{
			PhoneNumber: "+1234567890",
			Provider:    "Twilio",
		}, nil
	case "push":
		return &PushNotification{
			DeviceToken: "device-token-123",
			Platform:    "iOS",
		}, nil
	default:
		return nil, fmt.Errorf("unknown notification type: %s", notificationType)
	}
}

// Abstract Factory Pattern Example
// Creates families of related objects

// PaymentProcessor interface
type PaymentProcessor interface {
	ProcessPayment(amount float64) error
	Refund(transactionID string, amount float64) error
	GetName() string
}

// PayPalProcessor implements PaymentProcessor
type PayPalProcessor struct {
	APIKey string
}

func (p *PayPalProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("💳 Processing $%.2f via PayPal\n", amount)
	return nil
}

func (p *PayPalProcessor) Refund(transactionID string, amount float64) error {
	fmt.Printf("💰 Refunding $%.2f via PayPal (Transaction: %s)\n", amount, transactionID)
	return nil
}

func (p *PayPalProcessor) GetName() string {
	return "PayPal"
}

// StripeProcessor implements PaymentProcessor
type StripeProcessor struct {
	SecretKey string
}

func (s *StripeProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("💳 Processing $%.2f via Stripe\n", amount)
	return nil
}

func (s *StripeProcessor) Refund(transactionID string, amount float64) error {
	fmt.Printf("💰 Refunding $%.2f via Stripe (Transaction: %s)\n", amount, transactionID)
	return nil
}

func (s *StripeProcessor) GetName() string {
	return "Stripe"
}

// CryptoProcessor implements PaymentProcessor
type CryptoProcessor struct {
	WalletAddress string
	Currency      string
}

func (c *CryptoProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("₿ Processing $%.2f via Crypto (%s)\n", amount, c.Currency)
	return nil
}

func (c *CryptoProcessor) Refund(transactionID string, amount float64) error {
	fmt.Printf("₿ Refunding $%.2f via Crypto (Transaction: %s)\n", amount, transactionID)
	return nil
}

func (c *CryptoProcessor) GetName() string {
	return "Crypto"
}

// PaymentProcessorFactory creates payment processors
func PaymentProcessorFactory(processorType string) (PaymentProcessor, error) {
	switch processorType {
	case "paypal":
		return &PayPalProcessor{
			APIKey: "paypal-api-key",
		}, nil
	case "stripe":
		return &StripeProcessor{
			SecretKey: "stripe-secret-key",
		}, nil
	case "crypto":
		return &CryptoProcessor{
			WalletAddress: "0x1234567890abcdef",
			Currency:      "BTC",
		}, nil
	default:
		return nil, fmt.Errorf("unknown payment processor: %s", processorType)
	}
}

// Map-based Factory (alternative approach)
type NotificationCreator func() Notification

var notificationRegistry = map[string]NotificationCreator{
	"email": func() Notification {
		return &EmailNotification{To: "user@example.com", From: "noreply@app.com"}
	},
	"sms": func() Notification {
		return &SMSNotification{PhoneNumber: "+1234567890", Provider: "Twilio"}
	},
	"push": func() Notification {
		return &PushNotification{DeviceToken: "token-123", Platform: "iOS"}
	},
}

// CreateNotification creates a notification using the registry
func CreateNotification(notificationType string) (Notification, error) {
	creator, exists := notificationRegistry[notificationType]
	if !exists {
		return nil, fmt.Errorf("unknown notification type: %s", notificationType)
	}
	return creator(), nil
}

// RegisterNotification allows registering new notification types at runtime
func RegisterNotification(notificationType string, creator NotificationCreator) {
	notificationRegistry[notificationType] = creator
}

