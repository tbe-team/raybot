package liftmotorimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type liftMotorStateRepository struct {
	liftMotor liftmotor.LiftMotorState
	mu        sync.RWMutex
}

func NewLiftMotorStateRepository() liftmotor.LiftMotorStateRepository {
	return &liftMotorStateRepository{
		liftMotor: liftmotor.LiftMotorState{},
	}
}

func (r *liftMotorStateRepository) GetLiftMotorState(_ context.Context) (liftmotor.LiftMotorState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.liftMotor, nil
}

func (r *liftMotorStateRepository) UpdateLiftMotorState(_ context.Context, params liftmotor.UpdateLiftMotorStateParams) error {
	r.mu.RLock()
	state := r.liftMotor
	r.mu.RUnlock()

	if params.SetCurrentPosition {
		state.CurrentPosition = params.CurrentPosition
	}

	if params.SetTargetPosition {
		state.TargetPosition = params.TargetPosition
	}

	if params.SetIsRunning {
		state.IsRunning = params.IsRunning
	}

	if params.SetEnabled {
		state.Enabled = params.Enabled
	}

	r.mu.Lock()
	r.liftMotor = state
	r.mu.Unlock()

	return nil
}
