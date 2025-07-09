package ledimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/led"
)

type Repository struct {
	systemLedState      led.State
	systemLedConnection led.Connection
	alertLedState       led.State
	alertLedConnection  led.Connection
	mu                  sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		systemLedState: led.State{
			Mode: led.ModeOff,
		},
		systemLedConnection: led.Connection{
			Connected: false,
		},
		alertLedState: led.State{
			Mode: led.ModeOff,
		},
		alertLedConnection: led.Connection{
			Connected: false,
		},
	}
}

func (r *Repository) UpdateSystemLedState(_ context.Context, state led.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.systemLedState = state

	return nil
}

func (r *Repository) UpdateAlertLedState(_ context.Context, state led.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.alertLedState = state

	return nil
}

func (r *Repository) UpdateAlertLedConnection(_ context.Context, connection led.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.alertLedConnection = connection

	return nil
}

func (r *Repository) UpdateSystemLedConnection(_ context.Context, connection led.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.systemLedConnection = connection

	return nil
}
