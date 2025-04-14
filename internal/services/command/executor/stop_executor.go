package executor

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type stopExecutor struct {
	driveMotorService drivemotor.Service
}

func newStopExecutor(
	driveMotorService drivemotor.Service,
) stopExecutor {
	return stopExecutor{
		driveMotorService: driveMotorService,
	}
}

func (e stopExecutor) Execute(ctx context.Context, _ command.StopInputs) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return NewExecutorError(err, "failed to stop")
	}
	return nil
}
