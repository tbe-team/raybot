package drivemotor

import (
	"context"

	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrCanNotControlDriveMotor = xerror.BadRequest(picserial.ErrPICSerialNotConnected, "drivemotor.canNotControl", "can not control drive motor")

type UpdateDriveMotorStateParams struct {
	Direction    Direction `validate:"required_if=SetDirection true,enum"`
	SetDirection bool
	Speed        uint8
	SetSpeed     bool
	IsRunning    bool
	SetIsRunning bool
	Enabled      bool
	SetEnabled   bool
}

type MoveForwardParams struct {
	Speed uint8 `validate:"min=0,max=100"`
}

type MoveBackwardParams struct {
	Speed uint8 `validate:"min=0,max=100"`
}

type Service interface {
	// UpdateDriveMotorState updates the desired state of the drive motor.
	// This does not directly interact with the hardware, it just updates the internal state.
	UpdateDriveMotorState(ctx context.Context, params UpdateDriveMotorStateParams) error

	// MoveForward moves the drive motor forward.
	// This directly sends commands to the hardware.
	MoveForward(ctx context.Context, params MoveForwardParams) error

	// MoveBackward moves the drive motor backward.
	// This directly sends commands to the hardware.
	MoveBackward(ctx context.Context, params MoveBackwardParams) error

	// Stop stops the drive motor.
	// This directly sends commands to the hardware.
	Stop(ctx context.Context) error
}

//nolint:revive
type DriveMotorStateRepository interface {
	GetDriveMotorState(ctx context.Context) (DriveMotorState, error)
	UpdateDriveMotorState(ctx context.Context, params UpdateDriveMotorStateParams) error
}
