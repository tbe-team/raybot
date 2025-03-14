package model

import (
	"fmt"
	"time"
)

// BatteryState represents the state of the battery
type BatteryState struct {
	Current      uint16
	Temp         uint8
	Voltage      uint16
	CellVoltages []uint16
	Percent      uint8
	Fault        uint8
	Health       uint8
	UpdatedAt    time.Time
}

// ChargeState represents the state of the charge
type ChargeState struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}

// DischargeState represents the state of the discharge
type DischargeState struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}

// DistanceSensorState represents the state of the distance sensors
type DistanceSensorState struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
	UpdatedAt     time.Time
}

// LiftMotorState represents the state of the lift motor
type LiftMotorState struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
	UpdatedAt       time.Time
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
	IsRunning bool
	Enabled   bool
	UpdatedAt time.Time
}

type LocationState struct {
	CurrentLocation string
	UpdatedAt       time.Time
}

type RobotState struct {
	Battery        BatteryState
	Charge         ChargeState
	Discharge      DischargeState
	DistanceSensor DistanceSensorState
	LiftMotor      LiftMotorState
	DriveMotor     DriveMotorState
	Location       LocationState
}
