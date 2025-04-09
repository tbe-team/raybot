package liftmotorimpl

import (
	"context"
	"sync"
	"time"

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
	r.mu.Lock()
	defer r.mu.Unlock()

	r.liftMotor = liftmotor.LiftMotorState{
		CurrentPosition: params.CurrentPosition,
		TargetPosition:  params.TargetPosition,
		IsRunning:       params.IsRunning,
		Enabled:         params.Enabled,
		UpdatedAt:       time.Now(),
	}
	return nil
}
