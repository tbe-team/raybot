package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type moveForwardExecutor struct {
	driveMotorService drivemotor.Service
}

func newMoveForwardExecutor(
	driveMotorService drivemotor.Service,
) moveForwardExecutor {
	return moveForwardExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e moveForwardExecutor) Execute(ctx context.Context, _ command.MoveForwardInputs) error {
	if err := e.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction: drivemotor.DirectionForward,
		Speed:     100,
		Enabled:   true,
	}); err != nil {
		return NewExecutorError(err, "failed to update drive motor state, (start driving forward)")
	}
	return nil
}
