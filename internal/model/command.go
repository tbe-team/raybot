package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// CommandType represents the type of command
type CommandType uint16

func (s CommandType) Validate() error {
	if s < CommandTypeMoveToLocation || s > CommandTypeCloseCargo {
		return fmt.Errorf("invalid command type: %d", s)
	}
	return nil
}

func (s CommandType) String() string {
	return []string{
		"MOVE_TO_LOCATION",
		"LIFT_CARGO",
		"DROP_CARGO",
		"OPEN_CARGO",
		"CLOSE_CARGO",
	}[s]
}

const (
	CommandTypeMoveToLocation CommandType = iota
	CommandTypeLiftCargo
	CommandTypeDropCargo
	CommandTypeOpenCargo
	CommandTypeCloseCargo
)

// CommandStatus represents the status of the command
type CommandStatus uint8

func (s CommandStatus) Validate() error {
	if s < CommandStatusInProgress || s > CommandStatusFailed {
		return fmt.Errorf("invalid command status: %d", s)
	}
	return nil
}

func (s CommandStatus) String() string {
	return []string{
		"IN_PROGRESS",
		"SUCCEEDED",
		"FAILED",
	}[s]
}

const (
	CommandStatusInProgress CommandStatus = iota
	CommandStatusSucceeded
	CommandStatusFailed
)

// CommandSource represents where the command originated from
type CommandSource uint8

const (
	CommandSourceManual CommandSource = iota
	CommandSourceCloud
)

func (s CommandSource) Validate() error {
	switch s {
	case CommandSourceManual, CommandSourceCloud:
		return nil
	default:
		return fmt.Errorf("invalid command source: %d", s)
	}
}

func (s CommandSource) String() string {
	return []string{
		"MANUAL",
		"CLOUD",
	}[s]
}

// Command represents a robot command
type Command struct {
	ID          string
	Type        CommandType
	Status      CommandStatus
	Source      CommandSource
	Inputs      CommandInputs
	Error       *string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// CommandData represents the data of the command
type CommandInputs interface {
	isCommandInputs()
}

type CommandMoveToLocationInputs struct {
	Location string `json:"location" validate:"required"`
}

func (CommandMoveToLocationInputs) isCommandInputs() {}

type CommandLiftCargoInputs struct{}

func (CommandLiftCargoInputs) isCommandInputs() {}

type CommandDropCargoInputs struct{}

func (CommandDropCargoInputs) isCommandInputs() {}

type CommandOpenCargoInputs struct{}

func (CommandOpenCargoInputs) isCommandInputs() {}

type CommandCloseCargoInputs struct{}

func (CommandCloseCargoInputs) isCommandInputs() {}

func UnmarshalCommandInputs(cmdType CommandType, data []byte) (CommandInputs, error) {
	var inputs CommandInputs

	switch cmdType {
	case CommandTypeMoveToLocation:
		var moveToLocationInputs CommandMoveToLocationInputs
		if err := json.Unmarshal(data, &moveToLocationInputs); err != nil {
			return nil, fmt.Errorf("unmarshal move to location inputs: %w", err)
		}
		inputs = moveToLocationInputs
	case CommandTypeLiftCargo:
		inputs = CommandLiftCargoInputs{}
	case CommandTypeDropCargo:
		inputs = CommandDropCargoInputs{}
	case CommandTypeOpenCargo:
		inputs = CommandOpenCargoInputs{}
	case CommandTypeCloseCargo:
		inputs = CommandCloseCargoInputs{}
	default:
		return nil, fmt.Errorf("invalid command type: %d", cmdType)
	}

	return inputs, nil
}
