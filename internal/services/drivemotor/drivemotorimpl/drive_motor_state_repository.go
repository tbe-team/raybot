package drivemotorimpl

import (
	"context"
	"sync"
	"time"

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
	r.mu.Lock()
	defer r.mu.Unlock()

	r.driveMotor = drivemotor.DriveMotorState{
		Direction: params.Direction,
		Speed:     params.Speed,
		IsRunning: params.IsRunning,
		Enabled:   params.Enabled,
		UpdatedAt: time.Now(),
	}
	return nil
}
