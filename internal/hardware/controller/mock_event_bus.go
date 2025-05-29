package controller

import (
	"context"

	"github.com/tbe-team/raybot/pkg/eventbus"
)

var _ eventbus.EventBus = (*fakeEventBus)(nil)

type fakeEventBus struct {
	expectedPayload eventbus.Payload
}

func (e *fakeEventBus) Publish(_ string, _ *eventbus.Message) {}

func (e *fakeEventBus) Subscribe(_ context.Context, _ string, handler eventbus.HandlerFunc) {
	if e.expectedPayload != nil {
		handler(&eventbus.Message{
			Payload: e.expectedPayload,
		})
	}
}
