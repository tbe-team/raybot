package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

func newMoveBackwardExecutor(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.MoveBackwardInputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.MoveBackwardInputs) error {
			if err := driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
				Speed: 100,
			}); err != nil {
				return fmt.Errorf("failed to move backward: %w", err)
			}
			return nil
		},
		Hooks{},
		log,
		commandRepository,
	)
}
