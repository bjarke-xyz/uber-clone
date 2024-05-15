package core

import "context"

type Pubsub interface {
	SubscribeBytes(topic string) <-chan []byte
	Subscribe(topic string) <-chan any
	PublishBytes(ctx context.Context, topic string, msg []byte)
	Publish(ctx context.Context, topic string, msg any) error
	Close()
}
