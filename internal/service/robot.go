package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type BatteryParams struct {
	Current      int   `validate:"required"`
	Temp         int   `validate:"required"`
	Voltage      int   `validate:"required"`
	CellVoltages []int `validate:"required"`
	Percent      uint8 `validate:"required"`
	Fault        uint8 `validate:"required"`
	Health       uint8 `validate:"required"`
	Status       int   `validate:"required"`
}

type ChargeParams struct {
	CurrentLimit int `validate:"required"`
	Enabled      bool
}

type DischargeParams struct {
	CurrentLimit int `validate:"required"`
	Enabled      bool
}

type DistanceSensorParams struct {
	FrontDistance int `validate:"required"`
	BackDistance  int `validate:"required"`
	DownDistance  int `validate:"required"`
}

type LiftMotorParams struct {
	Direction model.LiftMotorDirection `validate:"required"`
	Speed     uint8                    `validate:"required"`
}

type DriveMotorParams struct {
	Direction model.DriveMotorDirection `validate:"required"`
	Speed     uint8                     `validate:"required"`
}

type UpdateRobotStateParams struct {
	Battery           BatteryParams `validate:"required_if=SetBattery true"`
	SetBattery        bool
	Charge            ChargeParams `validate:"required_if=SetCharge true"`
	SetCharge         bool
	Discharge         DischargeParams `validate:"required_if=SetDischarge true"`
	SetDischarge      bool
	DistanceSensor    DistanceSensorParams `validate:"required_if=SetDistanceSensor true"`
	SetDistanceSensor bool
	LiftMotor         LiftMotorParams `validate:"required_if=SetLiftMotor true"`
	SetLiftMotor      bool
	DriveMotor        DriveMotorParams `validate:"required_if=SetDriveMotor true"`
	SetDriveMotor     bool
}

type RobotService interface {
	UpdateRobotState(ctx context.Context, params UpdateRobotStateParams) (model.RobotState, error)
}
