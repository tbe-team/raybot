package eventbus

import (
	"context"
)

type Subscriber interface {
	// Subscribe subscribes to a topic and returns a channel for receiving messages.
	// When the context is canceled,
	// the channel will be closed and the subscriber will be removed from the topic.
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)

	// Unsubscribe unsubscribes from a topic.
	Unsubscribe(topic string, ch chan *Message) error
}

type Publisher interface {
	// Publish publishes a message to a topic.
	Publish(topic string, message *Message) error

	// PublishAsync publishes a message to a topic asynchronously.
	PublishAsync(topic string, message *Message) error
}

type EventBus interface {
	Subscriber
	Publisher
}
