package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/command"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
)

func TestService_route(t *testing.T) {
	execErr := errors.New("exec error")

	testCases := []struct {
		name            string
		cmd             command.Command
		expectedOutputs command.Outputs
		expectedErr     error
	}{
		{
			name: "stop movement execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeStopMovement,
				Inputs: &command.StopMovementInputs{},
			},
			expectedOutputs: command.StopMovementOutputs{},
		},
		{
			name: "stop movement execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeStopMovement,
				Inputs: &command.StopMovementInputs{},
			},
			expectedOutputs: command.StopMovementOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "move to execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeMoveTo,
				Inputs: &command.MoveToInputs{},
			},
			expectedOutputs: command.MoveToOutputs{},
		},
		{
			name: "move to execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeMoveTo,
				Inputs: &command.MoveToInputs{},
			},
			expectedOutputs: command.MoveToOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "move forward execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeMoveForward,
				Inputs: &command.MoveForwardInputs{},
			},
			expectedOutputs: command.MoveForwardOutputs{},
		},
		{
			name: "move forward execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeMoveForward,
				Inputs: &command.MoveForwardInputs{},
			},
			expectedOutputs: command.MoveForwardOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "move backward execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeMoveBackward,
				Inputs: &command.MoveBackwardInputs{},
			},
			expectedOutputs: command.MoveBackwardOutputs{},
		},
		{
			name: "move backward execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeMoveBackward,
				Inputs: &command.MoveBackwardInputs{},
			},
			expectedOutputs: command.MoveBackwardOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "cargo open execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeCargoOpen,
				Inputs: &command.CargoOpenInputs{},
			},
			expectedOutputs: command.CargoOpenOutputs{},
		},
		{
			name: "cargo open execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeCargoOpen,
				Inputs: &command.CargoOpenInputs{},
			},
			expectedOutputs: command.CargoOpenOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "cargo close execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeCargoClose,
				Inputs: &command.CargoCloseInputs{},
			},
			expectedOutputs: command.CargoCloseOutputs{},
		},
		{
			name: "cargo close execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeCargoClose,
				Inputs: &command.CargoCloseInputs{},
			},
			expectedOutputs: command.CargoCloseOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "cargo lift execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeCargoLift,
				Inputs: &command.CargoLiftInputs{},
			},
			expectedOutputs: command.CargoLiftOutputs{},
		},
		{
			name: "cargo lift execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeCargoLift,
				Inputs: &command.CargoLiftInputs{},
			},
			expectedOutputs: command.CargoLiftOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "cargo lower execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeCargoLower,
				Inputs: &command.CargoLowerInputs{},
			},
			expectedOutputs: command.CargoLowerOutputs{},
		},
		{
			name: "cargo lower execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeCargoLower,
				Inputs: &command.CargoLowerInputs{},
			},
			expectedOutputs: command.CargoLowerOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "cargo check QR execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeCargoCheckQR,
				Inputs: &command.CargoCheckQRInputs{},
			},
			expectedOutputs: command.CargoCheckQROutputs{},
		},
		{
			name: "cargo check QR execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeCargoCheckQR,
				Inputs: &command.CargoCheckQRInputs{},
			},
			expectedOutputs: command.CargoCheckQROutputs{},
			expectedErr:     execErr,
		},
		{
			name: "scan location execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeScanLocation,
				Inputs: &command.ScanLocationInputs{},
			},
			expectedOutputs: command.ScanLocationOutputs{},
		},
		{
			name: "scan location execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeScanLocation,
				Inputs: &command.ScanLocationInputs{},
			},
			expectedOutputs: command.ScanLocationOutputs{},
			expectedErr:     execErr,
		},
		{
			name: "wait execute successfully",
			cmd: command.Command{
				Type:   command.CommandTypeWait,
				Inputs: &command.WaitInputs{},
			},
			expectedOutputs: command.WaitOutputs{},
		},
		{
			name: "wait execute with error",
			cmd: command.Command{
				Type:   command.CommandTypeWait,
				Inputs: &command.WaitInputs{},
			},
			expectedOutputs: command.WaitOutputs{},
			expectedErr:     execErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(
				logging.NewNoopLogger(),
				commandmocks.NewFakeRunningCommandRepository(t),
				commandmocks.NewFakeRepository(t),
				tc.expectedErr,
			)
			outputs, err := s.route(context.Background(), tc.cmd)
			assert.Equal(t, tc.expectedOutputs, outputs)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
