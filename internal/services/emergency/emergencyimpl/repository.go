package emergencyimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/emergency"
)

type emergencyStateRepository struct {
	emergencyState emergency.State
	mu             sync.RWMutex
}

func NewEmergencyStateRepository() emergency.Repository {
	return &emergencyStateRepository{
		emergencyState: emergency.State{Locked: false},
	}
}

func (r *emergencyStateRepository) GetEmergencyState(_ context.Context) (emergency.State, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.emergencyState, nil
}

func (r *emergencyStateRepository) UpdateEmergencyState(_ context.Context, state emergency.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.emergencyState = state
	return nil
}
