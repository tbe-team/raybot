package batteryimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type batterySettingRepository struct {
	db      db.DB
	queries *sqlc.Queries
}

func NewBatterySettingRepository(db db.DB, queries *sqlc.Queries) battery.SettingRepository {
	return &batterySettingRepository{
		db:      db,
		queries: queries,
	}
}

func (r batterySettingRepository) GetChargeSetting(ctx context.Context) (battery.ChargeSetting, error) {
	row, err := r.queries.BatteryChargeSettingGet(ctx, r.db)
	if err != nil {
		return battery.ChargeSetting{}, fmt.Errorf("failed to get charge setting: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return battery.ChargeSetting{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	//nolint:gosec
	return battery.ChargeSetting{
		CurrentLimit: uint16(row.CurrentLimit),
		Enabled:      row.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r batterySettingRepository) GetDischargeSetting(ctx context.Context) (battery.DischargeSetting, error) {
	row, err := r.queries.BatteryDischargeSettingGet(ctx, r.db)
	if err != nil {
		return battery.DischargeSetting{}, fmt.Errorf("failed to get discharge setting: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return battery.DischargeSetting{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	//nolint:gosec
	return battery.DischargeSetting{
		CurrentLimit: uint16(row.CurrentLimit),
		Enabled:      row.Enabled == 1,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r batterySettingRepository) UpdateChargeSetting(ctx context.Context, params battery.UpdateChargeSettingParams) error {
	if err := r.queries.BatteryChargeSettingUpdate(ctx, r.db, sqlc.BatteryChargeSettingUpdateParams{
		CurrentLimit: int64(params.CurrentLimit),
		Enabled:      boolToInt64(params.Enabled),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("queriesupdate charge setting: %w", err)
	}

	return nil
}

func (r batterySettingRepository) UpdateDischargeSetting(ctx context.Context, params battery.UpdateDischargeSettingParams) error {
	if err := r.queries.BatteryDischargeSettingUpdate(ctx, r.db, sqlc.BatteryDischargeSettingUpdateParams{
		CurrentLimit: int64(params.CurrentLimit),
		Enabled:      boolToInt64(params.Enabled),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("queries update discharge setting: %w", err)
	}

	return nil
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
