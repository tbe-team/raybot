package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, cmd command.Command) error
}

type dispatcher struct {
	stopExecutor         stopExecutor
	moveToExecutor       moveToExecutor
	moveForwardExecutor  moveForwardExecutor
	moveBackwardExecutor moveBackwardExecutor

	cargoOpenExecutor  cargoOpenExecutor
	cargoCloseExecutor cargoCloseExecutor
	cargoLiftExecutor  cargoLiftExecutor
	cargoLowerExecutor cargoLowerExecutor
}

func NewDispatcher(
	cargoCfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	cargoService cargo.Service,
	liftMotorService liftmotor.Service,
) Dispatcher {
	return dispatcher{
		stopExecutor:         newStopExecutor(driveMotorService),
		moveToExecutor:       newMoveToExecutor(log, subscriber, driveMotorService),
		moveForwardExecutor:  newMoveForwardExecutor(driveMotorService),
		moveBackwardExecutor: newMoveBackwardExecutor(driveMotorService),

		cargoOpenExecutor:  newCargoOpenExecutor(cargoService),
		cargoCloseExecutor: newCargoCloseExecutor(cargoService),
		cargoLiftExecutor:  newCargoLiftExecutor(cargoCfg.LiftPosition, liftMotorService),
		cargoLowerExecutor: newCargoLowerExecutor(cargoCfg.LowerPosition, liftMotorService),
	}
}

func (d dispatcher) Dispatch(ctx context.Context, cmd command.Command) error {
	switch cmd.Type {
	case command.CommandTypeStop:
		i, ok := cmd.Inputs.(*command.StopInputs)
		if !ok {
			return fmt.Errorf("invalid stop inputs: %v", cmd.Inputs)
		}
		return d.stopExecutor.Execute(ctx, *i)

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

	case command.CommandTypeCargoOpen:
		i, ok := cmd.Inputs.(*command.CargoOpenInputs)
		if !ok {
			return fmt.Errorf("invalid cargo open inputs: %v", cmd.Inputs)
		}
		return d.cargoOpenExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoClose:
		i, ok := cmd.Inputs.(*command.CargoCloseInputs)
		if !ok {
			return fmt.Errorf("invalid cargo close inputs: %v", cmd.Inputs)
		}
		return d.cargoCloseExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoLift:
		i, ok := cmd.Inputs.(*command.CargoLiftInputs)
		if !ok {
			return fmt.Errorf("invalid cargo lift inputs: %v", cmd.Inputs)
		}
		return d.cargoLiftExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoLower:
		i, ok := cmd.Inputs.(*command.CargoLowerInputs)
		if !ok {
			return fmt.Errorf("invalid cargo lower inputs: %v", cmd.Inputs)
		}
		return d.cargoLowerExecutor.Execute(ctx, *i)

	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}
