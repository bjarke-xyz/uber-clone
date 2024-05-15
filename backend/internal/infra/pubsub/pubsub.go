package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
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

// SubscribeBytes implements core.Pubsub.
func (ps *InMemoryPubsub) SubscribeBytes(topic string) <-chan []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan []byte, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

// Subscribe implements core.Pubsub.
func (ps *InMemoryPubsub) Subscribe(topic string) <-chan any {
	panic("lol")
}

// PublishBytes implements core.Pubsub.
func (ps *InMemoryPubsub) PublishBytes(ctx context.Context, topic string, msg []byte) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *InMemoryPubsub) Publish(ctx context.Context, topic string, msg any) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal msg: %w", err)
	}
	ps.PublishBytes(ctx, topic, bytes)
	return nil
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
