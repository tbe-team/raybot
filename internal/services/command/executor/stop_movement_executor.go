package executor

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type stopMovementExecutor struct {
	driveMotorService drivemotor.Service
}

func newStopMovementExecutor(
	driveMotorService drivemotor.Service,
) CommandExecutor[command.StopMovementInputs, command.StopMovementOutputs] {
	return stopMovementExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e stopMovementExecutor) Execute(ctx context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return command.StopMovementOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return command.StopMovementOutputs{}, nil
}

func (e stopMovementExecutor) OnCancel(_ context.Context) error {
	return nil
}
