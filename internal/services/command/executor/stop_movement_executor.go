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
) *commandExecutor[command.StopMovementInputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.StopMovementInputs) error {
			if err := driveMotorService.Stop(ctx); err != nil {
				return fmt.Errorf("failed to stop drive motor: %w", err)
			}
			return nil
		},
		Hooks{},
		log,
		commandRepository,
	)
}
