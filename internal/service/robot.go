package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type BatteryParams struct {
	Current      uint
	Temp         uint
	Voltage      uint
	CellVoltages []uint
	Percent      uint8
	Fault        uint8
	Health       uint8
	Status       uint8
}

type ChargeParams struct {
	CurrentLimit uint
	Enabled      bool
}

type DischargeParams struct {
	CurrentLimit uint
	Enabled      bool
}

type DistanceSensorParams struct {
	FrontDistance uint
	BackDistance  uint
	DownDistance  uint
}

type LiftMotorParams struct {
	Direction model.LiftMotorDirection
	Speed     uint8
}

type DriveMotorParams struct {
	Direction model.DriveMotorDirection
	Speed     uint8
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
	UpdateRobotState(ctx context.Context, params UpdateRobotStateParams) (model.RobotState, error)
}
