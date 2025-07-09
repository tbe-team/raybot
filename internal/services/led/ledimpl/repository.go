package ledimpl

import (
	"context"
	"sync"
	"time"

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
			Mode:      led.ModeOff,
			UpdatedAt: time.Now(),
		},
		systemLedConnection: led.Connection{
			Connected: false,
		},
		alertLedState: led.State{
			Mode:      led.ModeOff,
			UpdatedAt: time.Now(),
		},
		alertLedConnection: led.Connection{
			Connected: false,
		},
	}
}

func (r *Repository) GetLeds(_ context.Context) (led.LedsOutput, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return led.LedsOutput{
		SystemLedState:      r.systemLedState,
		SystemLedConnection: r.systemLedConnection,
		AlertLedState:       r.alertLedState,
		AlertLedConnection:  r.alertLedConnection,
	}, nil
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
