package model

import (
	"fmt"
	"strconv"
	"time"
)

type PICSerialCommandType uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *PICSerialCommandType) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}
	switch n {
	case 0:
		*s = PICSerialCommandTypeBatteryCharge
	case 1:
		*s = PICSerialCommandTypeBatteryDischarge
	case 2:
		*s = PICSerialCommandTypeLiftMotor
	case 3:
		*s = PICSerialCommandTypeDriveMotor
	default:
		return fmt.Errorf("invalid PIC serial command type: %s", string(data))
	}

	return nil
}

const (
	PICSerialCommandTypeBatteryCharge PICSerialCommandType = iota
	PICSerialCommandTypeBatteryDischarge
	PICSerialCommandTypeLiftMotor
	PICSerialCommandTypeDriveMotor
)

type PICSerialCommand struct {
	ID        string
	Type      PICSerialCommandType
	Data      PICSerialCommandData
	CreatedAt time.Time
}

type PICSerialCommandData interface {
	isPICSerialCommandData()
}

type PICSerialCommandBatteryChargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (PICSerialCommandBatteryChargeData) isPICSerialCommandData() {}

type PICSerialCommandBatteryDischargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (PICSerialCommandBatteryDischargeData) isPICSerialCommandData() {}

type PICSerialCommandBatteryLiftMotorData struct {
	TargetPosition uint16
	Enable         bool
}

func (PICSerialCommandBatteryLiftMotorData) isPICSerialCommandData() {}

type MoveDirection uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *MoveDirection) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}
	switch n {
	case 0:
		*m = MoveDirectionForward
	case 1:
		*m = MoveDirectionBackward
	default:
		return fmt.Errorf("invalid move direction: %s", string(data))
	}

	return nil
}

const (
	MoveDirectionForward MoveDirection = iota
	MoveDirectionBackward
)

type PICSerialCommandBatteryDriveMotorData struct {
	Direction MoveDirection
	Speed     uint8
	Enable    bool
}

func (PICSerialCommandBatteryDriveMotorData) isPICSerialCommandData() {}
