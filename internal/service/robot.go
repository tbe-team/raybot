package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type BatteryParams struct {
	Current      uint16
	Temp         uint8
	Voltage      uint16
	CellVoltages []uint16
	Percent      uint8
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
	Direction model.DriveMotorDirection
	Speed     uint8
	IsRunning bool
	Enabled   bool
}

type UpdateRobotStateParams struct {
	Battery           BatteryParams
	SetBattery        bool
	Charge            ChargeParams
	SetCharge         bool
	Discharge         DischargeParams
	SetDischarge      bool
	DistanceSensor    DistanceSensorParams
	SetDistanceSensor bool
	LiftMotor         LiftMotorParams
	SetLiftMotor      bool
	DriveMotor        DriveMotorParams
	SetDriveMotor     bool
}

type RobotService interface {
	GetRobotState(ctx context.Context) (model.RobotState, error)
	UpdateRobotState(ctx context.Context, params UpdateRobotStateParams) (model.RobotState, error)
}
