package liftmotor

import "context"

type UpdateLiftMotorStateParams struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
}

type Service interface {
	UpdateLiftMotorState(ctx context.Context, params UpdateLiftMotorStateParams) error
}

//nolint:revive
type LiftMotorStateRepository interface {
	GetLiftMotorState(ctx context.Context) (LiftMotorState, error)
	UpdateLiftMotorState(ctx context.Context, params UpdateLiftMotorStateParams) error
}
