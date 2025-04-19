package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

func newStopMovementExecutor(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.StopMovementInputs, command.StopMovementOutputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
			if err := driveMotorService.Stop(ctx); err != nil {
				return command.StopMovementOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
			}
			return command.StopMovementOutputs{}, nil
		},
		Hooks[command.StopMovementOutputs]{},
		log,
		commandRepository,
	)
}
