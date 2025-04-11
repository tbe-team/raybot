package eventbus

import (
	"context"
)

type HandlerFunc func(ctx context.Context, msg *Message)

type Subscriber interface {
	// Subscribe subscribes to a topic and returns a channel for receiving messages.
	//
	// When the context is canceled, the subscriber will be removed from the topic.
	Subscribe(ctx context.Context, topic string, handler HandlerFunc)
}

type Publisher interface {
	// Publish publishes a message to a topic asynchronously.
	Publish(topic string, message *Message)
}

type EventBus interface {
	Subscriber
	Publisher
}
