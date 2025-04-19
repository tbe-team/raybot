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
	handler := stopMovementHandler{
		log:               log,
		driveMotorService: driveMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.StopMovementOutputs]{},
		log,
		commandRepository,
	)
}

type stopMovementHandler struct {
	log               *slog.Logger
	driveMotorService drivemotor.Service
}

func (h stopMovementHandler) Handle(ctx context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
	if err := h.driveMotorService.Stop(ctx); err != nil {
		return command.StopMovementOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return command.StopMovementOutputs{}, nil
}
