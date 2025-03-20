package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type RobotStateRepository struct {
	queries *sqlc.Queries
}

func NewRobotStateRepository(queries *sqlc.Queries) *RobotStateRepository {
	return &RobotStateRepository{
		queries: queries,
	}
}

func (r RobotStateRepository) GetRobotState(ctx context.Context, db db.SQLDB) (model.RobotState, error) {
	row, err := r.queries.RobotStateGet(ctx, db)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("queries get robot state: %w", err)
	}

	return r.convertRowToRobotState(row)
}

func (r RobotStateRepository) convertRowToRobotState(row sqlc.RobotStateGetRow) (model.RobotState, error) {
	var battery sqlc.Battery
	d, ok := row.Battery.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("battery is not a string")
	}
	if err := json.Unmarshal([]byte(d), &battery); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal battery: %w", err)
	}

	var batteryCharge sqlc.BatteryCharge
	d, ok = row.BatteryCharge.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("battery charge is not a string")
	}
	if err := json.Unmarshal([]byte(d), &batteryCharge); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal battery charge: %w", err)
	}

	var batteryDischarge sqlc.BatteryDischarge
	d, ok = row.BatteryDischarge.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("battery discharge is not a string")
	}
	if err := json.Unmarshal([]byte(d), &batteryDischarge); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal battery discharge: %w", err)
	}

	var distanceSensor sqlc.DistanceSensor
	d, ok = row.DistanceSensor.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("distance sensor is not a string")
	}
	if err := json.Unmarshal([]byte(d), &distanceSensor); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal distance sensor: %w", err)
	}

	var driveMotor sqlc.DriveMotor
	d, ok = row.DriveMotor.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("drive motor is not a string")
	}
	if err := json.Unmarshal([]byte(d), &driveMotor); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal drive motor: %w", err)
	}

	var liftMotor sqlc.LiftMotor
	d, ok = row.LiftMotor.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("lift motor is not a string")
	}
	if err := json.Unmarshal([]byte(d), &liftMotor); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal lift motor: %w", err)
	}

	var location sqlc.Location
	d, ok = row.Location.(string)
	if !ok {
		return model.RobotState{}, fmt.Errorf("location is not a string")
	}
	if err := json.Unmarshal([]byte(d), &location); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal location: %w", err)
	}

	cellVoltages := make([]uint16, 0)
	if err := json.Unmarshal([]byte(battery.CellVoltages), &cellVoltages); err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal cell voltages: %w", err)
	}

	batteryUpdatedAt, err := time.Parse(time.RFC3339, battery.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	batteryChargeUpdatedAt, err := time.Parse(time.RFC3339, batteryCharge.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	batteryDischargeUpdatedAt, err := time.Parse(time.RFC3339, batteryDischarge.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	distanceSensorUpdatedAt, err := time.Parse(time.RFC3339, distanceSensor.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	driveMotorUpdatedAt, err := time.Parse(time.RFC3339, driveMotor.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	liftMotorUpdatedAt, err := time.Parse(time.RFC3339, liftMotor.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	locationUpdatedAt, err := time.Parse(time.RFC3339, location.UpdatedAt)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.RobotState{
		Battery: model.Battery{
			Current:      uint16(battery.Current),
			Temp:         uint8(battery.Temp),
			Voltage:      uint16(battery.Voltage),
			CellVoltages: cellVoltages,
			Percent:      uint8(battery.Percent),
			Fault:        uint8(battery.Fault),
			Health:       uint8(battery.Health),
			UpdatedAt:    batteryUpdatedAt,
		},
		Charge: model.BatteryCharge{
			CurrentLimit: uint16(batteryCharge.CurrentLimit),
			Enabled:      batteryCharge.Enabled == 1,
			UpdatedAt:    batteryChargeUpdatedAt,
		},
		Discharge: model.BatteryDischarge{
			CurrentLimit: uint16(batteryDischarge.CurrentLimit),
			Enabled:      batteryDischarge.Enabled == 1,
			UpdatedAt:    batteryDischargeUpdatedAt,
		},
		DistanceSensor: model.DistanceSensor{
			FrontDistance: uint16(distanceSensor.FrontDistance),
			BackDistance:  uint16(distanceSensor.BackDistance),
			DownDistance:  uint16(distanceSensor.DownDistance),
			UpdatedAt:     distanceSensorUpdatedAt,
		},
		DriveMotor: model.DriveMotor{
			Direction: model.DriveMotorDirection(driveMotor.Direction),
			Speed:     uint8(driveMotor.Speed),
			IsRunning: driveMotor.IsRunning == 1,
			Enabled:   driveMotor.Enabled == 1,
			UpdatedAt: driveMotorUpdatedAt,
		},
		LiftMotor: model.LiftMotor{
			CurrentPosition: uint16(liftMotor.CurrentPosition),
			TargetPosition:  uint16(liftMotor.TargetPosition),
			IsRunning:       liftMotor.IsRunning == 1,
			Enabled:         liftMotor.Enabled == 1,
			UpdatedAt:       liftMotorUpdatedAt,
		},
		Location: model.Location{
			CurrentLocation: location.CurrentLocation,
			UpdatedAt:       locationUpdatedAt,
		},
	}, nil
}
