package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

func newMoveForwardExecutor(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.MoveForwardInputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.MoveForwardInputs) error {
			if err := driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
				Speed: 100,
			}); err != nil {
				return fmt.Errorf("failed to move forward: %w", err)
			}
			return nil
		},
		Hooks{},
		log,
		commandRepository,
	)
}
