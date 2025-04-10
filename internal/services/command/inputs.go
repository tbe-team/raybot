package command

import (
	"encoding/json"
	"fmt"
)

type Inputs interface {
	isInputs()
	CommandType() CommandType
}

type StopInputs struct{}

func (StopInputs) CommandType() CommandType {
	return CommandTypeStop
}
func (StopInputs) isInputs() {}

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

type MoveToInputs struct {
	Location string `json:"location" validate:"required"`
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

func UnmarshalInputs(cmdType CommandType, inputsBytes []byte) (Inputs, error) {
	var inputs Inputs

	switch cmdType {
	case CommandTypeStop:
		inputs = &StopInputs{}

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

	default:
		return nil, fmt.Errorf("invalid command type: %s", cmdType)
	}

	return inputs, nil
}
