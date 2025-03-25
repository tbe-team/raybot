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
	battery, err := r.convertRowToBattery(row.Battery)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert battery: %w", err)
	}

	charge, err := r.convertRowToBatteryCharge(row.BatteryCharge)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert battery charge: %w", err)
	}

	discharge, err := r.convertRowToBatteryDischarge(row.BatteryDischarge)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert battery discharge: %w", err)
	}

	distanceSensor, err := r.convertRowToDistanceSensor(row.DistanceSensor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert distance sensor: %w", err)
	}

	driveMotor, err := r.convertRowToDriveMotor(row.DriveMotor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert drive motor: %w", err)
	}

	liftMotor, err := r.convertRowToLiftMotor(row.LiftMotor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert lift motor: %w", err)
	}

	location, err := r.convertRowToLocation(row.Location)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert location: %w", err)
	}

	cargo, err := r.convertRowToCargo(row.Cargo)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert cargo: %w", err)
	}

	cargoDoorMotor, err := r.convertRowToCargoDoorMotor(row.CargoDoorMotor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("convert cargo door motor: %w", err)
	}

	return model.RobotState{
		Battery:        battery,
		Charge:         charge,
		Discharge:      discharge,
		DistanceSensor: distanceSensor,
		DriveMotor:     driveMotor,
		LiftMotor:      liftMotor,
		Location:       location,
		Cargo:          cargo,
		CargoDoorMotor: cargoDoorMotor,
	}, nil
}

func (r RobotStateRepository) convertRowToBattery(data any) (model.Battery, error) {
	d, ok := data.(string)
	if !ok {
		return model.Battery{}, fmt.Errorf("battery is not a string")
	}

	var battery sqlc.Battery
	if err := json.Unmarshal([]byte(d), &battery); err != nil {
		return model.Battery{}, fmt.Errorf("unmarshal battery: %w", err)
	}

	cellVoltages := make([]uint16, 0)
	if err := json.Unmarshal([]byte(battery.CellVoltages), &cellVoltages); err != nil {
		return model.Battery{}, fmt.Errorf("unmarshal cell voltages: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, battery.UpdatedAt)
	if err != nil {
		return model.Battery{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.Battery{
		Current:      uint16(battery.Current),
		Temp:         uint8(battery.Temp),
		Voltage:      uint16(battery.Voltage),
		CellVoltages: cellVoltages,
		Percent:      uint8(battery.Percent),
		Fault:        uint8(battery.Fault),
		Health:       uint8(battery.Health),
		UpdatedAt:    updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToBatteryCharge(data any) (model.BatteryCharge, error) {
	d, ok := data.(string)
	if !ok {
		return model.BatteryCharge{}, fmt.Errorf("battery charge is not a string")
	}

	var charge sqlc.BatteryCharge
	if err := json.Unmarshal([]byte(d), &charge); err != nil {
		return model.BatteryCharge{}, fmt.Errorf("unmarshal battery charge: %w", err)
	}

	updatedAt, err := parseTime(charge.UpdatedAt)
	if err != nil {
		return model.BatteryCharge{}, err
	}

	//nolint:gosec
	return model.BatteryCharge{
		CurrentLimit: uint16(charge.CurrentLimit),
		Enabled:      charge.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToDriveMotor(data any) (model.DriveMotor, error) {
	d, ok := data.(string)
	if !ok {
		return model.DriveMotor{}, fmt.Errorf("drive motor is not a string")
	}

	var motor sqlc.DriveMotor
	if err := json.Unmarshal([]byte(d), &motor); err != nil {
		return model.DriveMotor{}, fmt.Errorf("unmarshal drive motor: %w", err)
	}

	updatedAt, err := parseTime(motor.UpdatedAt)
	if err != nil {
		return model.DriveMotor{}, err
	}

	//nolint:gosec
	return model.DriveMotor{
		Direction: model.DriveMotorDirection(motor.Direction),
		Speed:     uint8(motor.Speed),
		IsRunning: motor.IsRunning == 1,
		Enabled:   motor.Enabled == 1,
		UpdatedAt: updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToBatteryDischarge(data any) (model.BatteryDischarge, error) {
	d, ok := data.(string)
	if !ok {
		return model.BatteryDischarge{}, fmt.Errorf("battery discharge is not a string")
	}

	var discharge sqlc.BatteryDischarge
	if err := json.Unmarshal([]byte(d), &discharge); err != nil {
		return model.BatteryDischarge{}, fmt.Errorf("unmarshal battery discharge: %w", err)
	}

	updatedAt, err := parseTime(discharge.UpdatedAt)
	if err != nil {
		return model.BatteryDischarge{}, err
	}

	//nolint:gosec
	return model.BatteryDischarge{
		CurrentLimit: uint16(discharge.CurrentLimit),
		Enabled:      discharge.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToDistanceSensor(data any) (model.DistanceSensor, error) {
	d, ok := data.(string)
	if !ok {
		return model.DistanceSensor{}, fmt.Errorf("distance sensor is not a string")
	}

	var sensor sqlc.DistanceSensor
	if err := json.Unmarshal([]byte(d), &sensor); err != nil {
		return model.DistanceSensor{}, fmt.Errorf("unmarshal distance sensor: %w", err)
	}

	updatedAt, err := parseTime(sensor.UpdatedAt)
	if err != nil {
		return model.DistanceSensor{}, err
	}

	//nolint:gosec
	return model.DistanceSensor{
		FrontDistance: uint16(sensor.FrontDistance),
		BackDistance:  uint16(sensor.BackDistance),
		DownDistance:  uint16(sensor.DownDistance),
		UpdatedAt:     updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToLiftMotor(data any) (model.LiftMotor, error) {
	d, ok := data.(string)
	if !ok {
		return model.LiftMotor{}, fmt.Errorf("lift motor is not a string")
	}

	var motor sqlc.LiftMotor
	if err := json.Unmarshal([]byte(d), &motor); err != nil {
		return model.LiftMotor{}, fmt.Errorf("unmarshal lift motor: %w", err)
	}

	updatedAt, err := parseTime(motor.UpdatedAt)
	if err != nil {
		return model.LiftMotor{}, err
	}

	//nolint:gosec
	return model.LiftMotor{
		CurrentPosition: uint16(motor.CurrentPosition),
		TargetPosition:  uint16(motor.TargetPosition),
		IsRunning:       motor.IsRunning == 1,
		Enabled:         motor.Enabled == 1,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToLocation(data any) (model.Location, error) {
	d, ok := data.(string)
	if !ok {
		return model.Location{}, fmt.Errorf("location is not a string")
	}

	var location sqlc.Location
	if err := json.Unmarshal([]byte(d), &location); err != nil {
		return model.Location{}, fmt.Errorf("unmarshal location: %w", err)
	}

	updatedAt, err := parseTime(location.UpdatedAt)
	if err != nil {
		return model.Location{}, err
	}

	//nolint:gosec
	return model.Location{
		CurrentLocation: location.CurrentLocation,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToCargo(data any) (model.Cargo, error) {
	d, ok := data.(string)
	if !ok {
		return model.Cargo{}, fmt.Errorf("cargo is not a string")
	}

	var cargo sqlc.Cargo
	if err := json.Unmarshal([]byte(d), &cargo); err != nil {
		return model.Cargo{}, fmt.Errorf("unmarshal cargo: %w", err)
	}

	updatedAt, err := parseTime(cargo.UpdatedAt)
	if err != nil {
		return model.Cargo{}, err
	}

	//nolint:gosec
	return model.Cargo{
		IsOpen:         cargo.IsOpen == 1,
		QRCode:         cargo.QrCode,
		BottomDistance: uint16(cargo.BottomDistance),
		UpdatedAt:      updatedAt,
	}, nil
}

func (r RobotStateRepository) convertRowToCargoDoorMotor(data any) (model.CargoDoorMotor, error) {
	d, ok := data.(string)
	if !ok {
		return model.CargoDoorMotor{}, fmt.Errorf("cargo door motor is not a string")
	}

	var motor sqlc.CargoDoorMotor
	if err := json.Unmarshal([]byte(d), &motor); err != nil {
		return model.CargoDoorMotor{}, fmt.Errorf("unmarshal cargo door motor: %w", err)
	}

	updatedAt, err := parseTime(motor.UpdatedAt)
	if err != nil {
		return model.CargoDoorMotor{}, err
	}

	//nolint:gosec
	return model.CargoDoorMotor{
		Direction: model.CargoDoorMotorDirection(motor.Direction),
		Speed:     uint8(motor.Speed),
		IsRunning: motor.IsRunning == 1,
		Enabled:   motor.Enabled == 1,
		UpdatedAt: updatedAt,
	}, nil
}
