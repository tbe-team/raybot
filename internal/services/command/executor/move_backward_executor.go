package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

const defaultMoveBackwardSpeed = 100

func newMoveBackwardExecutor(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.MoveBackwardInputs, command.MoveBackwardOutputs] {
	handler := moveBackwardHandler{
		log:               log,
		driveMotorService: driveMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.MoveBackwardOutputs]{},
		log,
		commandRepository,
	)
}

type moveBackwardHandler struct {
	log               *slog.Logger
	driveMotorService drivemotor.Service
}

func (h moveBackwardHandler) Handle(ctx context.Context, _ command.MoveBackwardInputs) (command.MoveBackwardOutputs, error) {
	if err := h.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
		Speed: defaultMoveBackwardSpeed,
	}); err != nil {
		return command.MoveBackwardOutputs{}, fmt.Errorf("failed to move backward: %w", err)
	}

	return command.MoveBackwardOutputs{}, nil
}
