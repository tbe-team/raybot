package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
)

type cargoCloseExecutor struct {
	cargoService cargo.Service
}

func newCargoCloseExecutor(
	cargoService cargo.Service,
) cargoCloseExecutor {
	return cargoCloseExecutor{
		cargoService: cargoService,
	}
}

func (e cargoCloseExecutor) Execute(ctx context.Context, _ command.CargoCloseInputs) error {
	if err := e.cargoService.CloseCargoDoor(ctx); err != nil {
		return NewExecutorError(err, "failed to close cargo")
	}
	return nil
}
