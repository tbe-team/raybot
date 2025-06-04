package limitswitchimpl

import (
	"context"
	"sync"
	"time"

	"github.com/tbe-team/raybot/internal/services/limitswitch"
)

type Repository struct {
	LimitSwitchs map[limitswitch.LimitSwitchID]limitswitch.LimitSwitch
	mu           sync.RWMutex
}

func NewRepository() limitswitch.Repository {
	return &Repository{
		LimitSwitchs: make(map[limitswitch.LimitSwitchID]limitswitch.LimitSwitch),
	}
}

func (r *Repository) GetLimitSwitchState(_ context.Context, id limitswitch.LimitSwitchID) (limitswitch.LimitSwitch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.LimitSwitchs[id], nil
}

func (r *Repository) UpdateLimitSwitchState(_ context.Context, id limitswitch.LimitSwitchID, pressed bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.LimitSwitchs[id] = limitswitch.LimitSwitch{
		Pressed:   pressed,
		UpdatedAt: time.Now(),
	}

	return nil
}
