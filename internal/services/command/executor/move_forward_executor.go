package executor

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type moveForwardExecutor struct {
	driveMotorService drivemotor.Service
}

func newMoveForwardExecutor(
	driveMotorService drivemotor.Service,
) CommandExecutor[command.MoveForwardInputs, command.MoveForwardOutputs] {
	return moveForwardExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e moveForwardExecutor) Execute(ctx context.Context, inputs command.MoveForwardInputs) (command.MoveForwardOutputs, error) {
	if err := e.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
		Speed: inputs.MotorSpeed,
	}); err != nil {
		return command.MoveForwardOutputs{}, fmt.Errorf("failed to move forward: %w", err)
	}

	return command.MoveForwardOutputs{}, nil
}

func (e moveForwardExecutor) OnCancel(ctx context.Context) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return nil
}
