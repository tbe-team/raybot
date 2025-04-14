package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type cargoLiftExecutor struct {
	liftPosition     uint16
	liftMotorService liftmotor.Service
}

func newCargoLiftExecutor(
	liftPosition uint16,
	liftMotorService liftmotor.Service,
) cargoLiftExecutor {
	return cargoLiftExecutor{
		liftPosition:     liftPosition,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLiftExecutor) Execute(ctx context.Context, _ command.CargoLiftInputs) error {
	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: e.liftPosition,
	}); err != nil {
		return NewExecutorError(err, "failed to set cargo position")
	}

	return nil
}
