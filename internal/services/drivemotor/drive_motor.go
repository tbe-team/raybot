package drivemotor

import "context"

type UpdateDriveMotorStateParams struct {
	Direction Direction `validate:"enum"`
	Speed     uint8     `validate:"min=0,max=100"`
	IsRunning bool
	Enabled   bool
}

type Service interface {
	UpdateDriveMotorState(ctx context.Context, params UpdateDriveMotorStateParams) error
}

//nolint:revive
type DriveMotorStateRepository interface {
	GetDriveMotorState(ctx context.Context) (DriveMotorState, error)
	UpdateDriveMotorState(ctx context.Context, params UpdateDriveMotorStateParams) error
}
