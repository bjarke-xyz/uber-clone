package core

import "context"

type Pubsub interface {
	Subscribe(topic string) <-chan []byte
	Publish(ctx context.Context, topic string, msg []byte)
	Close()
}
