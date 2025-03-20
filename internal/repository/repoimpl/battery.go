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

type BatteryRepository struct {
	queries *sqlc.Queries
}

func NewBatteryRepository(queries *sqlc.Queries) *BatteryRepository {
	return &BatteryRepository{queries: queries}
}

func (r BatteryRepository) GetBattery(ctx context.Context, db db.SQLDB) (model.Battery, error) {
	battery, err := r.queries.BatteryGet(ctx, db)
	if err != nil {
		return model.Battery{}, fmt.Errorf("queries get battery: %w", err)
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
		Voltage:      uint16(battery.Voltage),
		Temp:         uint8(battery.Temp),
		CellVoltages: cellVoltages,
		Percent:      uint8(battery.Percent),
		Fault:        uint8(battery.Fault),
		Health:       uint8(battery.Health),
		UpdatedAt:    updatedAt,
	}, nil
}

func (r BatteryRepository) UpdateBattery(ctx context.Context, db db.SQLDB, battery model.Battery) error {
	cellVoltages, err := json.Marshal(battery.CellVoltages)
	if err != nil {
		return fmt.Errorf("json marshal cell voltages: %w", err)
	}

	params := sqlc.BatteryUpdateParams{
		Current:      int64(battery.Current),
		Voltage:      int64(battery.Voltage),
		Temp:         int64(battery.Temp),
		CellVoltages: string(cellVoltages),
		Percent:      int64(battery.Percent),
		Fault:        int64(battery.Fault),
		Health:       int64(battery.Health),
		UpdatedAt:    battery.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.BatteryUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update battery: %w", err)
	}

	return nil
}

func (r BatteryRepository) GetBatteryCharge(ctx context.Context, db db.SQLDB) (model.BatteryCharge, error) {
	batteryCharge, err := r.queries.BatteryChargeGet(ctx, db)
	if err != nil {
		return model.BatteryCharge{}, fmt.Errorf("queries get battery charge: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, batteryCharge.UpdatedAt)
	if err != nil {
		return model.BatteryCharge{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.BatteryCharge{
		CurrentLimit: uint16(batteryCharge.CurrentLimit),
		Enabled:      batteryCharge.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r BatteryRepository) UpdateBatteryCharge(ctx context.Context, db db.SQLDB, batteryCharge model.BatteryCharge) error {
	params := sqlc.BatteryChargeUpdateParams{
		CurrentLimit: int64(batteryCharge.CurrentLimit),
		Enabled:      boolToInt64(batteryCharge.Enabled),
		UpdatedAt:    batteryCharge.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.BatteryChargeUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update battery charge: %w", err)
	}

	return nil
}

func (r BatteryRepository) GetBatteryDischarge(ctx context.Context, db db.SQLDB) (model.BatteryDischarge, error) {
	batteryDischarge, err := r.queries.BatteryDischargeGet(ctx, db)
	if err != nil {
		return model.BatteryDischarge{}, fmt.Errorf("queries get battery discharge: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, batteryDischarge.UpdatedAt)
	if err != nil {
		return model.BatteryDischarge{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.BatteryDischarge{
		CurrentLimit: uint16(batteryDischarge.CurrentLimit),
		Enabled:      batteryDischarge.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r BatteryRepository) UpdateBatteryDischarge(ctx context.Context, db db.SQLDB, batteryDischarge model.BatteryDischarge) error {
	params := sqlc.BatteryDischargeUpdateParams{
		CurrentLimit: int64(batteryDischarge.CurrentLimit),
		Enabled:      boolToInt64(batteryDischarge.Enabled),
		UpdatedAt:    batteryDischarge.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.BatteryDischargeUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update battery discharge: %w", err)
	}

	return nil
}
