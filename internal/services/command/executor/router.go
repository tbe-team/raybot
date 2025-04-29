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

type Router interface {
	Route(ctx context.Context, cmd command.Command) error
}

type executorRouter struct {
	stopMovementExecutor *commandExecutor[command.StopMovementInputs, command.StopMovementOutputs]
	moveToExecutor       *commandExecutor[command.MoveToInputs, command.MoveToOutputs]
	moveForwardExecutor  *commandExecutor[command.MoveForwardInputs, command.MoveForwardOutputs]
	moveBackwardExecutor *commandExecutor[command.MoveBackwardInputs, command.MoveBackwardOutputs]

	cargoOpenExecutor    *commandExecutor[command.CargoOpenInputs, command.CargoOpenOutputs]
	cargoCloseExecutor   *commandExecutor[command.CargoCloseInputs, command.CargoCloseOutputs]
	cargoLiftExecutor    *commandExecutor[command.CargoLiftInputs, command.CargoLiftOutputs]
	cargoLowerExecutor   *commandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs]
	cargoCheckQRExecutor *commandExecutor[command.CargoCheckQRInputs, command.CargoCheckQROutputs]

	scanLocationExecutor *commandExecutor[command.ScanLocationInputs, command.ScanLocationOutputs]
	waitExecutor         *commandExecutor[command.WaitInputs, command.WaitOutputs]
}

func NewRouter(
	cargoCfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	cargoService cargo.Service,
	commandRepository command.Repository,
) Router {
	return &executorRouter{
		stopMovementExecutor: newStopMovementExecutor(log, driveMotorService, commandRepository),
		moveToExecutor:       newMoveToExecutor(log, subscriber, driveMotorService, commandRepository),
		moveForwardExecutor:  newMoveForwardExecutor(log, driveMotorService, commandRepository),
		moveBackwardExecutor: newMoveBackwardExecutor(log, driveMotorService, commandRepository),

		cargoOpenExecutor:    newCargoOpenExecutor(log, subscriber, cargoService, commandRepository),
		cargoCloseExecutor:   newCargoCloseExecutor(log, subscriber, cargoService, commandRepository),
		cargoLiftExecutor:    newCargoLiftExecutor(cargoCfg, log, subscriber, liftMotorService, commandRepository),
		cargoLowerExecutor:   newCargoLowerExecutor(cargoCfg, log, subscriber, liftMotorService, commandRepository),
		cargoCheckQRExecutor: newCargoCheckQRExecutor(log, subscriber, commandRepository),

		scanLocationExecutor: newScanLocationExecutor(log, subscriber, driveMotorService, commandRepository),
		waitExecutor:         newWaitExecutor(log, commandRepository),
	}
}

func (r *executorRouter) Route(ctx context.Context, cmd command.Command) error {
	switch cmd.Type {
	case command.CommandTypeStopMovement:
		i, ok := cmd.Inputs.(*command.StopMovementInputs)
		if !ok {
			return fmt.Errorf("invalid stop movement inputs: %v", cmd.Inputs)
		}
		r.stopMovementExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeMoveTo:
		i, ok := cmd.Inputs.(*command.MoveToInputs)
		if !ok {
			return fmt.Errorf("invalid move to inputs: %v", cmd.Inputs)
		}
		r.moveToExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeMoveForward:
		i, ok := cmd.Inputs.(*command.MoveForwardInputs)
		if !ok {
			return fmt.Errorf("invalid move forward inputs: %v", cmd.Inputs)
		}
		r.moveForwardExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeMoveBackward:
		i, ok := cmd.Inputs.(*command.MoveBackwardInputs)
		if !ok {
			return fmt.Errorf("invalid move backward inputs: %v", cmd.Inputs)
		}
		r.moveBackwardExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeCargoOpen:
		i, ok := cmd.Inputs.(*command.CargoOpenInputs)
		if !ok {
			return fmt.Errorf("invalid cargo open inputs: %v", cmd.Inputs)
		}
		r.cargoOpenExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeCargoClose:
		i, ok := cmd.Inputs.(*command.CargoCloseInputs)
		if !ok {
			return fmt.Errorf("invalid cargo close inputs: %v", cmd.Inputs)
		}
		r.cargoCloseExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeCargoLift:
		i, ok := cmd.Inputs.(*command.CargoLiftInputs)
		if !ok {
			return fmt.Errorf("invalid cargo lift inputs: %v", cmd.Inputs)
		}
		r.cargoLiftExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeCargoLower:
		i, ok := cmd.Inputs.(*command.CargoLowerInputs)
		if !ok {
			return fmt.Errorf("invalid cargo lower inputs: %v", cmd.Inputs)
		}
		r.cargoLowerExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeCargoCheckQR:
		i, ok := cmd.Inputs.(*command.CargoCheckQRInputs)
		if !ok {
			return fmt.Errorf("invalid cargo check qr inputs: %v", cmd.Inputs)
		}
		r.cargoCheckQRExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeScanLocation:
		i, ok := cmd.Inputs.(*command.ScanLocationInputs)
		if !ok {
			return fmt.Errorf("invalid scan location inputs: %v", cmd.Inputs)
		}
		r.scanLocationExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	case command.CommandTypeWait:
		i, ok := cmd.Inputs.(*command.WaitInputs)
		if !ok {
			return fmt.Errorf("invalid wait inputs: %v", cmd.Inputs)
		}
		r.waitExecutor.Execute(ctx, cmd.ID, *i)
		return nil

	default:
		return fmt.Errorf("invalid command type: %v", cmd.Type)
	}
}

type NoopRouter struct{}

func NewNoopRouter() *NoopRouter {
	return &NoopRouter{}
}

func (r *NoopRouter) Route(_ context.Context, _ command.Command) error {
	return nil
}
