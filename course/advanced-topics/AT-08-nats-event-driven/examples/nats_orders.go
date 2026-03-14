package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderService struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewOrderService(url string) (*OrderService, error) {
	nc, err := nats.Connect(url,
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
		return nil, fmt.Errorf("connect: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("jetstream: %w", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"orders.*"},
		Retention: nats.LimitsPolicy,
		MaxMsgs:   10000,
		MaxBytes:  100 * 1024 * 1024,
	})
	if err != nil {
		log.Printf("Stream may already exist: %v", err)
	}

	return &OrderService{nc: nc, js: js}, nil
}

func (s *OrderService) PublishOrderCreated(order *Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	_, err = s.js.Publish("orders.created", data)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	log.Printf("Published order.created: %s", order.ID)
	return nil
}

func (s *OrderService) SubscribeToOrders(durable string, handler func(order *Order) error) (*nats.Subscription, error) {
	return s.js.Subscribe("orders.*", func(msg *nats.Msg) {
		var order Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Unmarshal error: %v", err)
			msg.Nak()
			return
		}

		log.Printf("Processing order: %s (subject: %s)", order.ID, msg.Subject)

		if err := handler(&order); err != nil {
			log.Printf("Handler error: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	}, nats.Durable(durable), nats.ManualAck(), nats.DeliverAll())
}

func (s *OrderService) Close() {
	s.nc.Close()
}

func main() {
	service, err := NewOrderService("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	sub, err := service.SubscribeToOrders("order-processor", func(order *Order) error {
		log.Printf("Processing order %s: %s x%d", order.ID, order.Product, order.Quantity)
		time.Sleep(100 * time.Millisecond)
		order.Status = "processed"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	for i := 0; i < 5; i++ {
		order := &Order{
			ID:        fmt.Sprintf("ORD-%d", i+1),
			UserID:    "user-123",
			Product:   "Widget",
			Quantity:  i + 1,
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		if err := service.PublishOrderCreated(order); err != nil {
			log.Printf("Publish error: %v", err)
		}

		time.Sleep(time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()
	log.Println("Shutting down")
}
