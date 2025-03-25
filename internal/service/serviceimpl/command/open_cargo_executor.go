package command

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

type OpenCargoExecutor struct {
	cargoControlService service.CargoControlService
}

func NewOpenCargoExecutor(cargoControlService service.CargoControlService) *OpenCargoExecutor {
	return &OpenCargoExecutor{
		cargoControlService: cargoControlService,
	}
}

func (e OpenCargoExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeOpenCargo {
		return fmt.Errorf("command type is not open cargo")
	}

	if err := e.cargoControlService.OpenCargoDoor(ctx); err != nil {
		return fmt.Errorf("failed to open cargo: %w", err)
	}

	return nil
}
