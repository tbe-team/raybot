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
	if err := e.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
		Speed: 100,
	}); err != nil {
		return NewExecutorError(err, "failed to move backward")
	}
	return nil
}
