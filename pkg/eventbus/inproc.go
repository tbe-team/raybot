package eventbus

import (
	"context"
	"sync"
)

var _ EventBus = (*InProcEventBus)(nil)

type InProcEventBus struct {
	subscribers map[string][]chan *Message
	mu          sync.RWMutex
}

func NewInProcEventBus() *InProcEventBus {
	return &InProcEventBus{
		subscribers: make(map[string][]chan *Message),
	}
}

func (e *InProcEventBus) Subscribe(ctx context.Context, topic string) (<-chan *Message, error) {
	ch := make(chan *Message)

	e.mu.Lock()
	e.subscribers[topic] = append(e.subscribers[topic], ch)
	e.mu.Unlock()

	go func() {
		<-ctx.Done()

		close(ch)

		e.removeSubscriber(topic, ch)
	}()

	return ch, nil
}

func (e *InProcEventBus) Unsubscribe(topic string, ch chan *Message) error {
	e.mu.Lock()
	e.removeSubscriber(topic, ch)
	e.mu.Unlock()
	return nil
}

func (e *InProcEventBus) Publish(topic string, message *Message) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, ch := range e.subscribers[topic] {
		ch <- message
	}

	return nil
}

func (e *InProcEventBus) PublishAsync(topic string, message *Message) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, ch := range e.subscribers[topic] {
		go func(ch chan *Message) {
			ch <- message
		}(ch)
	}

	return nil
}

func (e *InProcEventBus) removeSubscriber(topic string, toRemove chan *Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, c := range e.subscribers[topic] {
		if c == toRemove {
			e.subscribers[topic] = append(e.subscribers[topic][:i], e.subscribers[topic][i+1:]...)
			break
		}
	}
}
