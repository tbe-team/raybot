package repoimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type DriveMotorRepository struct {
	queries *sqlc.Queries
}

func NewDriveMotorRepository(queries *sqlc.Queries) *DriveMotorRepository {
	return &DriveMotorRepository{queries: queries}
}

func (r DriveMotorRepository) GetDriveMotor(ctx context.Context, db db.SQLDB) (model.DriveMotor, error) {
	driveMotor, err := r.queries.DriveMotorGet(ctx, db)
	if err != nil {
		return model.DriveMotor{}, fmt.Errorf("queries get drive motor: %w", err)
	}

	//nolint:gosec
	return model.DriveMotor{
		Direction: model.DriveMotorDirection(driveMotor.Direction),
		Speed:     uint8(driveMotor.Speed),
		IsRunning: driveMotor.IsRunning == 1,
		Enabled:   driveMotor.Enabled == 1,
		UpdatedAt: driveMotor.UpdatedAt,
	}, nil
}

func (r DriveMotorRepository) UpdateDriveMotor(ctx context.Context, db db.SQLDB, driveMotor model.DriveMotor) error {
	params := sqlc.DriveMotorUpdateParams{
		Direction: int64(driveMotor.Direction),
		Speed:     int64(driveMotor.Speed),
		IsRunning: boolToInt64(driveMotor.IsRunning),
		Enabled:   boolToInt64(driveMotor.Enabled),
		UpdatedAt: driveMotor.UpdatedAt,
	}
	if err := r.queries.DriveMotorUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update drive motor: %w", err)
	}

	return nil
}
