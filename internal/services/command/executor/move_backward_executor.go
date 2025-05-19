package executor

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

const defaultMoveBackwardSpeed = 100

type moveBackwardExecutor struct {
	driveMotorService drivemotor.Service
}

func newMoveBackwardExecutor(
	driveMotorService drivemotor.Service,
) CommandExecutor[command.MoveBackwardInputs, command.MoveBackwardOutputs] {
	return moveBackwardExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e moveBackwardExecutor) Execute(ctx context.Context, _ command.MoveBackwardInputs) (command.MoveBackwardOutputs, error) {
	if err := e.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
		Speed: defaultMoveBackwardSpeed,
	}); err != nil {
		return command.MoveBackwardOutputs{}, fmt.Errorf("failed to move backward: %w", err)
	}

	return command.MoveBackwardOutputs{}, nil
}

func (e moveBackwardExecutor) OnCancel(ctx context.Context) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return nil
}
