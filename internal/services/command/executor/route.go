package executor

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
)

func (s *service) route(ctx context.Context, cmd command.Command) (command.Outputs, error) {
	var (
		outputs command.Outputs
		err     error
	)

	switch cmd.Type {
	case command.CommandTypeStopMovement:
		i, ok := cmd.Inputs.(*command.StopMovementInputs)
		if !ok {
			return nil, fmt.Errorf("invalid stop movement inputs: %v", cmd.Inputs)
		}
		outputs, err = s.stopMovementExecutor.Execute(ctx, *i)

	case command.CommandTypeMoveTo:
		i, ok := cmd.Inputs.(*command.MoveToInputs)
		if !ok {
			return nil, fmt.Errorf("invalid move to inputs: %v", cmd.Inputs)
		}
		outputs, err = s.moveToExecutor.Execute(ctx, *i)

	case command.CommandTypeMoveForward:
		i, ok := cmd.Inputs.(*command.MoveForwardInputs)
		if !ok {
			return nil, fmt.Errorf("invalid move forward inputs: %v", cmd.Inputs)
		}
		outputs, err = s.moveForwardExecutor.Execute(ctx, *i)

	case command.CommandTypeMoveBackward:
		i, ok := cmd.Inputs.(*command.MoveBackwardInputs)
		if !ok {
			return nil, fmt.Errorf("invalid move backward inputs: %v", cmd.Inputs)
		}
		outputs, err = s.moveBackwardExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoOpen:
		i, ok := cmd.Inputs.(*command.CargoOpenInputs)
		if !ok {
			return nil, fmt.Errorf("invalid cargo open inputs: %v", cmd.Inputs)
		}
		outputs, err = s.cargoOpenExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoClose:
		i, ok := cmd.Inputs.(*command.CargoCloseInputs)
		if !ok {
			return nil, fmt.Errorf("invalid cargo close inputs: %v", cmd.Inputs)
		}
		outputs, err = s.cargoCloseExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoLift:
		i, ok := cmd.Inputs.(*command.CargoLiftInputs)
		if !ok {
			return nil, fmt.Errorf("invalid cargo lift inputs: %v", cmd.Inputs)
		}
		outputs, err = s.cargoLiftExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoLower:
		i, ok := cmd.Inputs.(*command.CargoLowerInputs)
		if !ok {
			return nil, fmt.Errorf("invalid cargo lower inputs: %v", cmd.Inputs)
		}
		outputs, err = s.cargoLowerExecutor.Execute(ctx, *i)

	case command.CommandTypeCargoCheckQR:
		i, ok := cmd.Inputs.(*command.CargoCheckQRInputs)
		if !ok {
			return nil, fmt.Errorf("invalid cargo check QR inputs: %v", cmd.Inputs)
		}
		outputs, err = s.cargoCheckQRExecutor.Execute(ctx, *i)

	case command.CommandTypeScanLocation:
		i, ok := cmd.Inputs.(*command.ScanLocationInputs)
		if !ok {
			return nil, fmt.Errorf("invalid scan location inputs: %v", cmd.Inputs)
		}
		outputs, err = s.scanLocationExecutor.Execute(ctx, *i)

	case command.CommandTypeWait:
		i, ok := cmd.Inputs.(*command.WaitInputs)
		if !ok {
			return nil, fmt.Errorf("invalid wait inputs: %v", cmd.Inputs)
		}
		outputs, err = s.waitExecutor.Execute(ctx, *i)

	default:
		return nil, fmt.Errorf("invalid command type: %v", cmd.Type)
	}

	return outputs, err
}
