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
) *commandExecutor[command.MoveForwardInputs, command.MoveForwardOutputs] {
	handler := moveForwardHandler{
		log:               log,
		driveMotorService: driveMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.MoveForwardOutputs]{},
		log,
		commandRepository,
	)
}

type moveForwardHandler struct {
	log               *slog.Logger
	driveMotorService drivemotor.Service
}

func (h moveForwardHandler) Handle(ctx context.Context, _ command.MoveForwardInputs) (command.MoveForwardOutputs, error) {
	if err := h.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
		Speed: 100,
	}); err != nil {
		return command.MoveForwardOutputs{}, fmt.Errorf("failed to move forward: %w", err)
	}

	return command.MoveForwardOutputs{}, nil
}
