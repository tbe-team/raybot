package liftmotor

import (
	"context"

	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrCanNotControlLiftMotor = xerror.BadRequest(picserial.ErrPICSerialNotConnected, "liftmotor.canNotControl", "can not control lift motor")

type UpdateLiftMotorStateParams struct {
	CurrentPosition    uint16
	SetCurrentPosition bool
	TargetPosition     uint16
	SetTargetPosition  bool
	IsRunning          bool
	SetIsRunning       bool
	Enabled            bool
	SetEnabled         bool
}

type SetCargoPositionParams struct {
	Position uint16
}

type Service interface {
	// UpdateLiftMotorState updates the desired state of the lift motor.
	// This does not directly interact with the hardware, it just updates the internal state.
	UpdateLiftMotorState(ctx context.Context, params UpdateLiftMotorStateParams) error

	// SetCargoPosition moves the cargo to the specified position using hardware control.
	// This directly sends commands to the hardware.
	SetCargoPosition(ctx context.Context, params SetCargoPositionParams) error

	// Stop stops the cargo motor using hardware control.
	Stop(ctx context.Context) error
}

//nolint:revive
type LiftMotorStateRepository interface {
	GetLiftMotorState(ctx context.Context) (LiftMotorState, error)
	UpdateLiftMotorState(ctx context.Context, params UpdateLiftMotorStateParams) error
}
