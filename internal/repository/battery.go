package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type BatteryRepository interface {
	GetBattery(ctx context.Context, db db.SQLDB) (model.Battery, error)
	UpdateBattery(ctx context.Context, db db.SQLDB, battery model.Battery) error

	GetBatteryCharge(ctx context.Context, db db.SQLDB) (model.BatteryCharge, error)
	UpdateBatteryCharge(ctx context.Context, db db.SQLDB, batteryCharge model.BatteryCharge) error

	GetBatteryDischarge(ctx context.Context, db db.SQLDB) (model.BatteryDischarge, error)
	UpdateBatteryDischarge(ctx context.Context, db db.SQLDB, batteryDischarge model.BatteryDischarge) error
}
