package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
)

type cargoOpenExecutor struct {
	cargoService cargo.Service
}

func newCargoOpenExecutor(
	cargoService cargo.Service,
) cargoOpenExecutor {
	return cargoOpenExecutor{
		cargoService: cargoService,
	}
}

func (e cargoOpenExecutor) Execute(ctx context.Context, _ command.CargoOpenInputs) error {
	if err := e.cargoService.OpenCargoDoor(ctx); err != nil {
		return NewExecutorError(err, "failed to open cargo")
	}
	return nil
}
