package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type moveBackwardExecutor struct {
	driveMotorService drivemotor.Service
}

func newMoveBackwardExecutor(
	driveMotorService drivemotor.Service,
) moveBackwardExecutor {
	return moveBackwardExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e moveBackwardExecutor) Execute(ctx context.Context, _ command.MoveBackwardInputs) error {
	if err := e.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction: drivemotor.DirectionBackward,
		Speed:     100,
		Enabled:   true,
	}); err != nil {
		return NewExecutorError(err, "failed to update drive motor state, (start driving backward)")
	}
	return nil
}
