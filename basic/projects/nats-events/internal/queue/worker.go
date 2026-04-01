package queue

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Worker handles processing messages from a queue group
type Worker struct {
	nc      *nats.Conn
	group   string
	subject string
	handler func([]byte) error
}

// NewWorker creates a new queue worker
func NewWorker(nc *nats.Conn, group, subject string, handler func([]byte) error) *Worker {
	return &Worker{
		nc:      nc,
		group:   group,
		subject: subject,
		handler: handler,
	}
}

// Start begins processing messages from the queue group
func (w *Worker) Start(ctx context.Context) error {
	_, err := w.nc.QueueSubscribe(w.subject, w.group, func(msg *nats.Msg) {
		if err := w.handler(msg.Data); err != nil {
			// On failure, NAK and requeue for retry
			log.Printf("Worker: failed to process message: %v", err)
			if nakErr := msg.Nak(); nakErr != nil {
				log.Printf("Worker: failed to NAK message: %v", nakErr)
			}
		} else {
			// On success, ACK the message
			if ackErr := msg.Ack(); ackErr != nil {
				log.Printf("Worker: failed to ACK message: %v", ackErr)
			}
		}
	})
	return err
}

// StartWithRetry starts the worker with configurable retry settings
func (w *Worker) StartWithRetry(ctx context.Context, maxRetries int, retryDelay time.Duration) error {
	_, err := w.nc.QueueSubscribe(w.subject, w.group, func(msg *nats.Msg) {
		if err := w.handler(msg.Data); err != nil {
			log.Printf("Worker: failed to process message: %v", err)
			// Get retry count from header if available
			retries := getRetryCount(msg)
			if retries < maxRetries {
				log.Printf("Worker: retrying message (attempt %d/%d)", retries+1, maxRetries)
				if nakErr := msg.NakWithDelay(retryDelay); nakErr != nil {
					log.Printf("Worker: failed to NAK message: %v", nakErr)
				}
			} else {
				log.Printf("Worker: max retries reached, dropping message")
				msg.Term()
			}
		} else {
			if ackErr := msg.Ack(); ackErr != nil {
				log.Printf("Worker: failed to ACK message: %v", ackErr)
			}
		}
	})
	return err
}

// getRetryCount extracts the retry count from message headers
func getRetryCount(msg *nats.Msg) int {
	if msg.Header == nil {
		return 0
	}
	retryStr := msg.Header.Get("Nats-Retries")
	if retryStr == "" {
		return 0
	}
	var count int
	for _, c := range retryStr {
		if c >= '0' && c <= '9' {
			count = count*10 + int(c-'0')
		}
	}
	return count
}
