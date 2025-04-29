package eventbus

import "context"

var _ EventBus = (*NoopEventBus)(nil)

type NoopEventBus struct{}

func NewNoopEventBus() *NoopEventBus {
	return &NoopEventBus{}
}

func (e *NoopEventBus) Publish(_ string, _ *Message) {}

func (e *NoopEventBus) Subscribe(_ context.Context, _ string, _ HandlerFunc) {}
