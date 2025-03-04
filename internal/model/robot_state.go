package model

import (
	"fmt"
	"time"
)

// BatteryState represents the state of the battery
type BatteryState struct {
	Current      uint
	Temp         uint
	Voltage      uint
	CellVoltages []uint
	Percent      uint8
	Fault        uint8
	Health       uint8
	Status       uint8
	UpdatedAt    time.Time
}

// ChargeState represents the state of the charge
type ChargeState struct {
	CurrentLimit uint
	Enabled      bool
	UpdatedAt    time.Time
}

// DischargeState represents the state of the discharge
type DischargeState struct {
	CurrentLimit uint
	Enabled      bool
	UpdatedAt    time.Time
}

// DistanceSensorState represents the state of the distance sensors
type DistanceSensorState struct {
	FrontDistance uint
	BackDistance  uint
	DownDistance  uint
	UpdatedAt     time.Time
}

type LiftMotorDirection uint8

func (s LiftMotorDirection) Validate() error {
	switch s {
	case LiftMotorDirectionDown, LiftMotorDirectionUp:
		return nil
	default:
		return fmt.Errorf("invalid lift motor direction: %d", s)
	}
}

const (
	LiftMotorDirectionUp LiftMotorDirection = iota
	LiftMotorDirectionDown
)

// LiftMotorState represents the state of the lift motor
type LiftMotorState struct {
	Direction LiftMotorDirection
	Speed     uint8
	UpdatedAt time.Time
}

type DriveMotorDirection uint8

func (s DriveMotorDirection) Validate() error {
	switch s {
	case DriveMotorDirectionForward, DriveMotorDirectionBackward:
		return nil
	default:
		return fmt.Errorf("invalid drive motor direction: %d", s)
	}
}

const (
	DriveMotorDirectionForward DriveMotorDirection = iota
	DriveMotorDirectionBackward
)

// DriveMotorState represents the state of the drive motor
type DriveMotorState struct {
	Direction DriveMotorDirection
	Speed     uint8
	UpdatedAt time.Time
}

type RobotState struct {
	Battery        BatteryState
	Charge         ChargeState
	Discharge      DischargeState
	DistanceSensor DistanceSensorState
	LiftMotor      LiftMotorState
	DriveMotor     DriveMotorState
}
