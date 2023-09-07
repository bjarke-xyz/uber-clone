package pubsub

import (
	"context"
	"sync"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
)

type InMemoryPubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan []byte
	closed bool
}

func NewInMemoryPubsub() core.Pubsub {
	return &InMemoryPubsub{
		mu:     sync.RWMutex{},
		subs:   make(map[string][]chan []byte),
		closed: false,
	}
}

// Subscribe implements core.Pubsub.
func (ps *InMemoryPubsub) Subscribe(topic string) <-chan []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan []byte, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

// Publish implements core.Pubsub.
func (ps *InMemoryPubsub) Publish(ctx context.Context, topic string, msg []byte) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

// Close implements core.Pubsub.
func (ps *InMemoryPubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
