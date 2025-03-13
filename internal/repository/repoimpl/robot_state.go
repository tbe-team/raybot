package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"

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

	return r.scanRobotState(row)
}

func (r RobotStateRepository) UpdateRobotState(ctx context.Context, db db.SQLDB, state model.RobotState) error {
	batteryState, err := json.Marshal(state.Battery)
	if err != nil {
		return fmt.Errorf("marshal battery state: %w", err)
	}

	chargeState, err := json.Marshal(state.Charge)
	if err != nil {
		return fmt.Errorf("marshal charge state: %w", err)
	}

	dischargeState, err := json.Marshal(state.Discharge)
	if err != nil {
		return fmt.Errorf("marshal discharge state: %w", err)
	}

	distanceSensorState, err := json.Marshal(state.DistanceSensor)
	if err != nil {
		return fmt.Errorf("marshal distance sensor state: %w", err)
	}

	liftMotorState, err := json.Marshal(state.LiftMotor)
	if err != nil {
		return fmt.Errorf("marshal lift motor state: %w", err)
	}

	driveMotorState, err := json.Marshal(state.DriveMotor)
	if err != nil {
		return fmt.Errorf("marshal drive motor state: %w", err)
	}

	params := sqlc.RobotStateUpdateParams{
		BatteryState:        string(batteryState),
		ChargeState:         string(chargeState),
		DischargeState:      string(dischargeState),
		DistanceSensorState: string(distanceSensorState),
		LiftMotorState:      string(liftMotorState),
		DriveMotorState:     string(driveMotorState),
	}
	err = r.queries.RobotStateUpdate(ctx, db, params)
	if err != nil {
		return fmt.Errorf("queries update robot state: %w", err)
	}

	return nil
}

func (r RobotStateRepository) scanRobotState(row sqlc.RobotState) (model.RobotState, error) {
	var state model.RobotState

	err := json.Unmarshal([]byte(row.BatteryState), &state.Battery)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal battery state: %w", err)
	}

	err = json.Unmarshal([]byte(row.ChargeState), &state.Charge)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal charge state: %w", err)
	}

	err = json.Unmarshal([]byte(row.DischargeState), &state.Discharge)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal discharge state: %w", err)
	}

	err = json.Unmarshal([]byte(row.DistanceSensorState), &state.DistanceSensor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal distance sensor state: %w", err)
	}

	err = json.Unmarshal([]byte(row.LiftMotorState), &state.LiftMotor)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("unmarshal lift motor state: %w", err)
	}

	return state, nil
}
