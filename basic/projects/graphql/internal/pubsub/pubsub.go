package pubsub

import (
	"sync"
)

type PubSub struct {
	subscribers map[string][]chan interface{}
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan interface{}),
	}
}

func (ps *PubSub) Subscribe(topic string) chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan interface{}, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Unsubscribe(topic string, ch chan interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if subs, ok := ps.subscribers[topic]; ok {
		for i, sub := range subs {
			if sub == ch {
				ps.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}
}

func (ps *PubSub) Publish(topic string, payload interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if subs, ok := ps.subscribers[topic]; ok {
		for _, ch := range subs {
			select {
			case ch <- payload:
			default:
			}
		}
	}
}
