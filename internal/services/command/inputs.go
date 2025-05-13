package command

import (
	"encoding/json"
	"fmt"
)

var (
	_ Inputs = (*StopMovementInputs)(nil)
	_ Inputs = (*MoveForwardInputs)(nil)
	_ Inputs = (*MoveBackwardInputs)(nil)
	_ Inputs = (*MoveToInputs)(nil)
	_ Inputs = (*CargoOpenInputs)(nil)
	_ Inputs = (*CargoCloseInputs)(nil)
	_ Inputs = (*CargoLiftInputs)(nil)
	_ Inputs = (*CargoLowerInputs)(nil)
	_ Inputs = (*CargoCheckQRInputs)(nil)
	_ Inputs = (*ScanLocationInputs)(nil)
	_ Inputs = (*WaitInputs)(nil)
)

type Inputs interface {
	isInputs()
	CommandType() CommandType
}

type StopMovementInputs struct{}

func (StopMovementInputs) CommandType() CommandType {
	return CommandTypeStopMovement
}
func (StopMovementInputs) isInputs() {}

type MoveForwardInputs struct{}

func (MoveForwardInputs) CommandType() CommandType {
	return CommandTypeMoveForward
}
func (MoveForwardInputs) isInputs() {}

type MoveBackwardInputs struct{}

func (MoveBackwardInputs) CommandType() CommandType {
	return CommandTypeMoveBackward
}
func (MoveBackwardInputs) isInputs() {}

type MoveDirection string

func (m MoveDirection) Validate() error {
	if m != MoveDirectionForward && m != MoveDirectionBackward {
		return fmt.Errorf("invalid move direction: %s", m)
	}
	return nil
}

func (m MoveDirection) String() string {
	return string(m)
}

const (
	MoveDirectionForward  MoveDirection = "FORWARD"
	MoveDirectionBackward MoveDirection = "BACKWARD"
)

type MoveToInputs struct {
	Location  string        `json:"location" validate:"required"`
	Direction MoveDirection `json:"direction" validate:"required,enum"`
}

func (MoveToInputs) CommandType() CommandType {
	return CommandTypeMoveTo
}
func (MoveToInputs) isInputs() {}

type CargoOpenInputs struct{}

func (CargoOpenInputs) CommandType() CommandType {
	return CommandTypeCargoOpen
}
func (CargoOpenInputs) isInputs() {}

type CargoCloseInputs struct{}

func (CargoCloseInputs) CommandType() CommandType {
	return CommandTypeCargoClose
}
func (CargoCloseInputs) isInputs() {}

type CargoLiftInputs struct{}

func (CargoLiftInputs) CommandType() CommandType {
	return CommandTypeCargoLift
}
func (CargoLiftInputs) isInputs() {}

type CargoLowerInputs struct{}

func (CargoLowerInputs) CommandType() CommandType {
	return CommandTypeCargoLower
}
func (CargoLowerInputs) isInputs() {}

type CargoCheckQRInputs struct {
	QRCode string `json:"qr_code" validate:"required"`
}

func (CargoCheckQRInputs) CommandType() CommandType {
	return CommandTypeCargoCheckQR
}
func (CargoCheckQRInputs) isInputs() {}

type ScanLocationInputs struct{}

func (ScanLocationInputs) CommandType() CommandType {
	return CommandTypeScanLocation
}
func (ScanLocationInputs) isInputs() {}

type WaitInputs struct {
	DurationMs int64 `json:"duration_ms" validate:"required"`
}

func (WaitInputs) CommandType() CommandType {
	return CommandTypeWait
}
func (WaitInputs) isInputs() {}

func UnmarshalInputs(cmdType CommandType, inputsBytes []byte) (Inputs, error) {
	var inputs Inputs

	switch cmdType {
	case CommandTypeStopMovement:
		inputs = &StopMovementInputs{}

	case CommandTypeMoveForward:
		inputs = &MoveForwardInputs{}

	case CommandTypeMoveBackward:
		inputs = &MoveBackwardInputs{}

	case CommandTypeMoveTo:
		i := &MoveToInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal move to inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoOpen:
		inputs = &CargoOpenInputs{}

	case CommandTypeCargoClose:
		inputs = &CargoCloseInputs{}

	case CommandTypeCargoLift:
		inputs = &CargoLiftInputs{}

	case CommandTypeCargoLower:
		inputs = &CargoLowerInputs{}

	case CommandTypeCargoCheckQR:
		i := &CargoCheckQRInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo check qr inputs: %w", err)
		}
		inputs = i

	case CommandTypeScanLocation:
		inputs = &ScanLocationInputs{}

	case CommandTypeWait:
		i := &WaitInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal wait inputs: %w", err)
		}
		inputs = i

	default:
		return nil, fmt.Errorf("invalid command type: %s", cmdType)
	}

	return inputs, nil
}
