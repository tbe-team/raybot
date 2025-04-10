package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, cmd command.Command) error
}

type dispatcher struct {
	moveToExecutor       moveToExecutor
	moveForwardExecutor  moveForwardExecutor
	moveBackwardExecutor moveBackwardExecutor
}

func NewDispatcher(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
) Dispatcher {
	return dispatcher{
		moveToExecutor:       newMoveToExecutor(log, driveMotorService),
		moveForwardExecutor:  newMoveForwardExecutor(driveMotorService),
		moveBackwardExecutor: newMoveBackwardExecutor(driveMotorService),
	}
}

func (d dispatcher) Dispatch(ctx context.Context, cmd command.Command) error {
	switch cmd.Type {
	case command.CommandTypeMoveTo:
		i, ok := cmd.Inputs.(*command.MoveToInputs)
		if !ok {
			return fmt.Errorf("invalid move to inputs: %v", cmd.Inputs)
		}
		return d.moveToExecutor.Execute(ctx, *i)

	case command.CommandTypeMoveForward:
		i, ok := cmd.Inputs.(*command.MoveForwardInputs)
		if !ok {
			return fmt.Errorf("invalid move forward inputs: %v", cmd.Inputs)
		}
		return d.moveForwardExecutor.Execute(ctx, *i)

	case command.CommandTypeMoveBackward:
		i, ok := cmd.Inputs.(*command.MoveBackwardInputs)
		if !ok {
			return fmt.Errorf("invalid move backward inputs: %v", cmd.Inputs)
		}
		return d.moveBackwardExecutor.Execute(ctx, *i)

	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}
