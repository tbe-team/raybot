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

type MoveForwardInputs struct {
	MotorSpeed uint8 `json:"motor_speed" validate:"required,max=100"`
}

func (MoveForwardInputs) CommandType() CommandType {
	return CommandTypeMoveForward
}
func (MoveForwardInputs) isInputs() {}

type MoveBackwardInputs struct {
	MotorSpeed uint8 `json:"motor_speed" validate:"required,max=100"`
}

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
	Location   string        `json:"location" validate:"required"`
	Direction  MoveDirection `json:"direction" validate:"required,enum"`
	MotorSpeed uint8         `json:"motor_speed" validate:"required,max=100"`
}

func (MoveToInputs) CommandType() CommandType {
	return CommandTypeMoveTo
}
func (MoveToInputs) isInputs() {}

type CargoOpenInputs struct {
	MotorSpeed uint8 `json:"motor_speed" validate:"required,max=100"`
}

func (CargoOpenInputs) CommandType() CommandType {
	return CommandTypeCargoOpen
}
func (CargoOpenInputs) isInputs() {}

type CargoCloseInputs struct {
	MotorSpeed uint8 `json:"motor_speed" validate:"required,max=100"`
}

func (CargoCloseInputs) CommandType() CommandType {
	return CommandTypeCargoClose
}
func (CargoCloseInputs) isInputs() {}

type CargoLiftInputs struct {
	Position   uint16 `json:"position" validate:"required"`
	MotorSpeed uint8  `json:"motor_speed" validate:"required,max=100"`
}

func (CargoLiftInputs) CommandType() CommandType {
	return CommandTypeCargoLift
}
func (CargoLiftInputs) isInputs() {}

type BottomObstacleTracking struct {
	// Start detecting obstacle when distance is below this value
	EnterDistance uint16 `json:"enter_distance" validate:"required"`
	// Stop detecting obstacle when distance is above this value
	ExitDistance uint16 `json:"exit_distance" validate:"required"`
}

type CargoLowerInputs struct {
	Position               uint16                 `json:"position" validate:"required"`
	MotorSpeed             uint8                  `json:"motor_speed" validate:"required,max=100"`
	BottomObstacleTracking BottomObstacleTracking `json:"bottom_obstacle_tracking" validate:"required"`
}

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
		i := &StopMovementInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stop movement inputs: %w", err)
		}
		inputs = i

	case CommandTypeMoveForward:
		i := &MoveForwardInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal move forward inputs: %w", err)
		}
		inputs = i

	case CommandTypeMoveBackward:
		i := &MoveBackwardInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal move backward inputs: %w", err)
		}
		inputs = i

	case CommandTypeMoveTo:
		i := &MoveToInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal move to inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoOpen:
		i := &CargoOpenInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo open inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoClose:
		i := &CargoCloseInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo close inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoLift:
		i := &CargoLiftInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo lift inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoLower:
		i := &CargoLowerInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo lower inputs: %w", err)
		}
		inputs = i

	case CommandTypeCargoCheckQR:
		i := &CargoCheckQRInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cargo check qr inputs: %w", err)
		}
		inputs = i

	case CommandTypeScanLocation:
		i := &ScanLocationInputs{}
		if err := json.Unmarshal(inputsBytes, i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal scan location inputs: %w", err)
		}
		inputs = i

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
