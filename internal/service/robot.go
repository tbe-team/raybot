package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type BatteryParams struct {
	Current      uint
	Temp         uint `validate:"gt=0,lte=100"`
	Voltage      uint
	CellVoltages []uint
	Percent      uint8 `validate:"gte=0,lte=100"`
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
	Direction model.LiftMotorDirection `validate:"enum"`
	Speed     uint8                    `validate:"gte=0,lte=100"`
}

type DriveMotorParams struct {
	Direction model.DriveMotorDirection `validate:"enum"`
	Speed     uint8                     `validate:"gte=0,lte=100"`
}

type UpdateRobotStateParams struct {
	Battery           BatteryParams `validate:"required_if=SetBattery true,omitempty"`
	SetBattery        bool
	Charge            ChargeParams `validate:"required_if=SetCharge true,omitempty"`
	SetCharge         bool
	Discharge         DischargeParams `validate:"required_if=SetDischarge true,omitempty"`
	SetDischarge      bool
	DistanceSensor    DistanceSensorParams `validate:"required_if=SetDistanceSensor true,omitempty"`
	SetDistanceSensor bool
	LiftMotor         LiftMotorParams `validate:"required_if=SetLiftMotor true,omitempty"`
	SetLiftMotor      bool
	DriveMotor        DriveMotorParams `validate:"required_if=SetDriveMotor true,omitempty"`
	SetDriveMotor     bool
}

type RobotService interface {
	UpdateRobotState(ctx context.Context, params UpdateRobotStateParams) (model.RobotState, error)
}
