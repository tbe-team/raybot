package command

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

type CloseCargoExecutor struct {
	cargoControlService service.CargoControlService
}

func NewCloseCargoExecutor(cargoControlService service.CargoControlService) *CloseCargoExecutor {
	return &CloseCargoExecutor{cargoControlService: cargoControlService}
}

func (e CloseCargoExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeCloseCargo {
		return fmt.Errorf("command type is not close cargo")
	}

	if err := e.cargoControlService.CloseCargoDoor(ctx); err != nil {
		return fmt.Errorf("failed to close cargo: %w", err)
	}

	return nil
}
