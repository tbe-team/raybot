package drivemotorimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type driveMotorStateRepository struct {
	driveMotor drivemotor.DriveMotorState
	mu         sync.RWMutex
}

func NewDriveMotorStateRepository() drivemotor.DriveMotorStateRepository {
	return &driveMotorStateRepository{
		driveMotor: drivemotor.DriveMotorState{},
	}
}

func (r *driveMotorStateRepository) GetDriveMotorState(_ context.Context) (drivemotor.DriveMotorState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.driveMotor, nil
}

func (r *driveMotorStateRepository) UpdateDriveMotorState(_ context.Context, params drivemotor.UpdateDriveMotorStateParams) error {
	r.mu.RLock()
	state := r.driveMotor
	r.mu.RUnlock()

	if params.SetDirection {
		state.Direction = params.Direction
	}

	if params.SetSpeed {
		state.Speed = params.Speed
	}

	if params.SetIsRunning {
		state.IsRunning = params.IsRunning
	}

	if params.SetEnabled {
		state.Enabled = params.Enabled
	}

	r.mu.Lock()
	r.driveMotor = state
	r.mu.Unlock()

	return nil
}
