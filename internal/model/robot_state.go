package model

import "time"

// BatteryState represents the state of the battery
type BatteryState struct {
	Current      int
	Temp         int
	Voltage      int
	CellVoltages []int
	Percent      uint8
	Fault        uint8
	Health       uint8
	Status       int
	UpdatedAt    time.Time
}

// ChargeState represents the state of the charge
type ChargeState struct {
	CurrentLimit int
	Enabled      bool
	UpdatedAt    time.Time
}

// DischargeState represents the state of the discharge
type DischargeState struct {
	CurrentLimit int
	Enabled      bool
	UpdatedAt    time.Time
}

// DistanceSensorState represents the state of the distance sensors
type DistanceSensorState struct {
	FrontDistance int
	BackDistance  int
	DownDistance  int
	UpdatedAt     time.Time
}

type LiftMotorDirection int8

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

type DriveMotorDirection int8

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
