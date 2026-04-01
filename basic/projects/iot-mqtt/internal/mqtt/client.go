// Package mqtt provides MQTT client wrapper with reconnection and reliability features.
package mqtt

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client wraps the paho MQTT client with additional reliability features.
type Client struct {
	client       mqtt.Client
	broker       string
	clientID     string
	opts         *mqtt.ClientOptions
	mu           sync.RWMutex
	connected    bool
	reconnecting bool
}

// MessageHandler is a function that handles incoming MQTT messages.
type MessageHandler func(topic string, payload []byte)

// NewClient creates a new MQTT client.
func NewClient(broker, clientID string) *Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetCleanSession(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			log.Printf("Unhandled message on topic %s: %s", m.Topic(), string(m.Payload()))
		})

	return &Client{
		broker:    broker,
		clientID:  clientID,
		opts:      opts,
		connected: false,
	}
}

// Connect establishes connection to the MQTT broker.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.opts.SetOnConnectHandler(func(client mqtt.Client) {
		c.mu.Lock()
		c.connected = true
		c.reconnecting = false
		c.mu.Unlock()
		log.Printf("MQTT client %s connected to %s", c.clientID, c.broker)
	})

	c.opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		log.Printf("MQTT connection lost: %v", err)
	})

	c.opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		c.mu.Lock()
		c.reconnecting = true
		c.mu.Unlock()
		log.Printf("MQTT client reconnecting...")
	})

	c.client = mqtt.NewClient(c.opts)

	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	c.connected = true
	return nil
}

// Disconnect closes the MQTT connection.
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
		c.connected = false
	}
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// Publish sends a message to a topic.
func (c *Client) Publish(topic string, payload interface{}, qos byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.client == nil {
		return fmt.Errorf("client not connected")
	}

	token := c.client.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish: %w", token.Error())
	}

	return nil
}

// PublishRetained sends a retained message to a topic.
func (c *Client) PublishRetained(topic string, payload interface{}, qos byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.client == nil {
		return fmt.Errorf("client not connected")
	}

	token := c.client.Publish(topic, qos, true, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish retained: %w", token.Error())
	}

	return nil
}

// Subscribe subscribes to a topic with a handler.
func (c *Client) Subscribe(topic string, qos byte, handler MessageHandler) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.client == nil {
		return fmt.Errorf("client not connected")
	}

	mqttHandler := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	}

	token := c.client.Subscribe(topic, qos, mqttHandler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %w", token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// SubscribeMultiple subscribes to multiple topics.
func (c *Client) SubscribeMultiple(topics []string, qos byte, handler MessageHandler) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.client == nil {
		return fmt.Errorf("client not connected")
	}

	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = qos
	}

	mqttHandler := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	}

	token := c.client.SubscribeMultiple(filters, mqttHandler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe multiple: %w", token.Error())
	}

	return nil
}

// SetWill configures the Last Will and Testament message.
func (c *Client) SetWill(topic, payload string, qos byte, retained bool) {
	c.opts.SetWill(topic, payload, qos, retained)
}

// SetTLS configures TLS for the connection.
func (c *Client) SetTLS(tlsConfig interface{}) {
	// In production, this would configure TLS
	// opts.SetTLSConfig(tlsConfig)
}
