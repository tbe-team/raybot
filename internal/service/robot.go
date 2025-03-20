package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type BatteryParams struct {
	Current      uint16 `validate:"min=0"`
	Temp         uint8  `validate:"min=0,max=100"`
	Voltage      uint16 `validate:"min=0"`
	CellVoltages []uint16
	Percent      uint8 `validate:"min=0,max=100"`
	Fault        uint8
	Health       uint8
}

type ChargeParams struct {
	CurrentLimit uint16
	Enabled      bool
}

type DischargeParams struct {
	CurrentLimit uint16
	Enabled      bool
}

type DistanceSensorParams struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
}

type LiftMotorParams struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
}

type DriveMotorParams struct {
	Direction model.DriveMotorDirection `validate:"enum"`
	Speed     uint8                     `validate:"min=0,max=100"`
	IsRunning bool
	Enabled   bool
}

type LocationParams struct {
	CurrentLocation string
}

type UpdateRobotStateParams struct {
	Battery           BatteryParams `validate:"omitempty,required_if=SetBattery true"`
	SetBattery        bool
	Charge            ChargeParams `validate:"omitempty,required_if=SetCharge true"`
	SetCharge         bool
	Discharge         DischargeParams `validate:"omitempty,required_if=SetDischarge true"`
	SetDischarge      bool
	DistanceSensor    DistanceSensorParams `validate:"omitempty,required_if=SetDistanceSensor true"`
	SetDistanceSensor bool
	LiftMotor         LiftMotorParams `validate:"omitempty,required_if=SetLiftMotor true"`
	SetLiftMotor      bool
	DriveMotor        DriveMotorParams `validate:"omitempty,required_if=SetDriveMotor true"`
	SetDriveMotor     bool
	Location          LocationParams `validate:"omitempty,required_if=SetLocation true"`
	SetLocation       bool
}

type RobotService interface {
	GetRobotState(ctx context.Context) (model.RobotState, error)
	UpdateRobotState(ctx context.Context, params UpdateRobotStateParams) (model.RobotState, error)
}
