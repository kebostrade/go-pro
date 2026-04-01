package jetstream

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Publisher handles publishing messages to NATS JetStream
type Publisher struct {
	js jetstream.JetStream
}

// NewPublisher creates a new JetStream publisher
func NewPublisher(nc *nats.Conn) (*Publisher, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &Publisher{js: js}, nil
}

// Publish sends a message to the specified subject
func (p *Publisher) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := p.js.Publish(ctx, subject, data)
	return err
}

// PublishJSON serializes and publishes a value as JSON
func (p *Publisher) PublishJSON(ctx context.Context, subject string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(ctx, subject, data)
	return err
}

// PublishAsync publishes a message asynchronously
func (p *Publisher) PublishAsync(ctx context.Context, subject string, data []byte) error {
	_, err := p.js.PublishAsync(subject, data)
	if err != nil {
		return err
	}
	// Wait for publish to complete
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
