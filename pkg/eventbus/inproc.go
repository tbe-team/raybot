package eventbus

import (
	"context"
	"log/slog"
	"sync"
)

var _ EventBus = (*InProcEventBus)(nil)

type InProcEventBus struct {
	log *slog.Logger

	subscribers map[string][]*subscriber
	mu          sync.RWMutex
}

func NewInProcEventBus(log *slog.Logger) *InProcEventBus {
	return &InProcEventBus{
		log:         log.With("component", "inproc_event_bus"),
		subscribers: make(map[string][]*subscriber),
	}
}

func (e *InProcEventBus) Subscribe(ctx context.Context, topic string, handler HandlerFunc) {
	sub := &subscriber{
		handler: handler,
	}

	e.mu.Lock()
	e.subscribers[topic] = append(e.subscribers[topic], sub)
	e.mu.Unlock()

	go func() {
		<-ctx.Done()

		e.removeSubscriber(topic, sub)
	}()
}

func (e *InProcEventBus) Publish(topic string, message *Message) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, sub := range e.subscribers[topic] {
		go func(sub *subscriber) {
			defer func() {
				if r := recover(); r != nil {
					e.log.Error("recovered panic in subscriber", slog.Any("error", r))
				}
			}()

			sub.handle(context.TODO(), message)
		}(sub)
	}
}

func (e *InProcEventBus) removeSubscriber(topic string, toRemove *subscriber) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, c := range e.subscribers[topic] {
		if c == toRemove {
			e.subscribers[topic] = append(e.subscribers[topic][:i], e.subscribers[topic][i+1:]...)
			break
		}
	}
}

type subscriber struct {
	handler HandlerFunc
}

func (s subscriber) handle(ctx context.Context, msg *Message) {
	s.handler(ctx, msg)
}
