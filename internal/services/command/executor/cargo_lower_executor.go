package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type cargoLowerExecutor struct {
	lowerPosition    uint16
	liftMotorService liftmotor.Service
}

func newCargoLowerExecutor(
	lowerPosition uint16,
	liftMotorService liftmotor.Service,
) cargoLowerExecutor {
	return cargoLowerExecutor{
		lowerPosition:    lowerPosition,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLowerExecutor) Execute(ctx context.Context, _ command.CargoLowerInputs) error {
	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: e.lowerPosition,
	}); err != nil {
		return NewExecutorError(err, "failed to set cargo position")
	}
	return nil
}
