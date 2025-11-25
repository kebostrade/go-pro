package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/projects/design-patterns/behavioral"
	"github.com/DimaJoyti/go-pro/basic/projects/design-patterns/creational"
	"github.com/DimaJoyti/go-pro/basic/projects/design-patterns/structural"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║           🎨 Design Patterns in Go - Demo                   ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// CREATIONAL PATTERNS
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("🏗️  CREATIONAL PATTERNS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	
	demonstrateSingleton()
	demonstrateFactory()
	demonstrateBuilder()

	// STRUCTURAL PATTERNS
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("🔧 STRUCTURAL PATTERNS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	
	demonstrateAdapter()
	demonstrateDecorator()

	// BEHAVIORAL PATTERNS
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("🎭 BEHAVIORAL PATTERNS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	
	demonstrateStrategy()
	demonstrateObserver()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    ✅ Demo Complete!                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func demonstrateSingleton() {
	fmt.Println("\n1️⃣  SINGLETON PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Get database instances
	db1 := creational.GetDatabase()
	db2 := creational.GetDatabase()
	
	fmt.Printf("db1 == db2: %v\n", db1 == db2)
	
	db1.Connect()
	db2.Connect()
	
	// Get config instances
	config := creational.GetConfig()
	config.Set("theme", "dark")
	
	value, _ := config.Get("theme")
	fmt.Printf("Config theme: %s\n", value)
}

func demonstrateFactory() {
	fmt.Println("\n2️⃣  FACTORY PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Create different notification types
	email, _ := creational.NotificationFactory("email")
	sms, _ := creational.NotificationFactory("sms")
	push, _ := creational.NotificationFactory("push")
	
	email.Send("Welcome to our platform!")
	sms.Send("Your verification code is 123456")
	push.Send("You have a new message")
	
	// Payment processors
	fmt.Println()
	paypal, _ := creational.PaymentProcessorFactory("paypal")
	stripe, _ := creational.PaymentProcessorFactory("stripe")
	
	paypal.ProcessPayment(99.99)
	stripe.ProcessPayment(149.99)
}

func demonstrateBuilder() {
	fmt.Println("\n3️⃣  BUILDER PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Build HTTP request
	request := creational.NewHTTPRequestBuilder().
		Method("POST").
		URL("https://api.example.com/users").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer token123").
		Body(`{"name":"John Doe"}`).
		Timeout(60).
		Build()
	
	fmt.Println(request)
	
	// Build SQL query
	query := creational.NewSQLQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > 18").
		Where("active = true").
		OrderBy("created_at DESC").
		Limit(10).
		Build()
	
	fmt.Printf("\nSQL Query:\n%s\n", query)
	
	// Build user
	user := creational.NewUserBuilder().
		ID("user-123").
		Username("johndoe").
		Email("john@example.com").
		Name("John", "Doe").
		Age(30).
		AddRole("admin").
		Build()
	
	fmt.Printf("\nUser: %s (%s)\n", user.Username, user.Email)
}

func demonstrateAdapter() {
	fmt.Println("\n4️⃣  ADAPTER PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Use different players through same interface
	var player structural.MediaPlayer
	
	player = structural.NewVLCAdapter()
	player.Play("movie.vlc")
	
	player = structural.NewMP4Adapter()
	player.Play("video.mp4")
}

func demonstrateDecorator() {
	fmt.Println("\n5️⃣  DECORATOR PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Build coffee with decorators
	coffee := &structural.SimpleCoffee{}
	fmt.Printf("%s: $%.2f\n", coffee.Description(), coffee.Cost())
	
	coffeeWithMilk := structural.NewMilkDecorator(coffee)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilk.Description(), coffeeWithMilk.Cost())
	
	coffeeWithMilkAndSugar := structural.NewSugarDecorator(coffeeWithMilk)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilkAndSugar.Description(), coffeeWithMilkAndSugar.Cost())
	
	fancyCoffee := structural.NewWhipDecorator(coffeeWithMilkAndSugar)
	fmt.Printf("%s: $%.2f\n", fancyCoffee.Description(), fancyCoffee.Cost())
}

func demonstrateStrategy() {
	fmt.Println("\n6️⃣  STRATEGY PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Payment strategies
	payment := behavioral.NewPaymentContext(&behavioral.CreditCardStrategy{
		CardNumber: "1234567890123456",
		CVV:        "123",
	})
	payment.ExecutePayment(100.00)
	
	payment.SetStrategy(&behavioral.PayPalStrategy{
		Email: "user@example.com",
	})
	payment.ExecutePayment(50.00)
	
	payment.SetStrategy(&behavioral.CryptoStrategy{
		WalletAddress: "0x1234567890abcdef",
		Currency:      "BTC",
	})
	payment.ExecutePayment(200.00)
	
	// Sorting strategies
	fmt.Println()
	data := []int{64, 34, 25, 12, 22, 11, 90}
	
	sorter := behavioral.NewSortContext(&behavioral.BubbleSortStrategy{})
	sorted := sorter.Sort(data)
	fmt.Printf("Sorted: %v\n", sorted)
	
	sorter.SetStrategy(&behavioral.QuickSortStrategy{})
	sorted = sorter.Sort(data)
	fmt.Printf("Sorted: %v\n", sorted)
}

func demonstrateObserver() {
	fmt.Println("\n7️⃣  OBSERVER PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")
	
	// Create event manager
	eventManager := behavioral.NewEventManager()
	
	// Create observers
	emailObserver := &behavioral.EmailObserver{
		ID:    "email-1",
		Email: "user@example.com",
	}
	
	smsObserver := &behavioral.SMSObserver{
		ID:          "sms-1",
		PhoneNumber: "+1234567890",
	}
	
	logObserver := &behavioral.LogObserver{
		ID: "log-1",
	}
	
	// Attach observers
	eventManager.Attach(emailObserver)
	eventManager.Attach(smsObserver)
	eventManager.Attach(logObserver)
	
	// Trigger events
	eventManager.Notify("user.registered", map[string]string{
		"username": "johndoe",
		"email":    "john@example.com",
	})
	
	// Detach one observer
	eventManager.Detach("sms-1")
	
	eventManager.Notify("order.placed", map[string]interface{}{
		"order_id": "ORD-123",
		"total":    99.99,
	})
}

