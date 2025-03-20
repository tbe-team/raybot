package repoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type LiftMotorRepository struct {
	queries *sqlc.Queries
}

func NewLiftMotorRepository(queries *sqlc.Queries) *LiftMotorRepository {
	return &LiftMotorRepository{queries: queries}
}

func (r LiftMotorRepository) GetLiftMotor(ctx context.Context, db db.SQLDB) (model.LiftMotor, error) {
	liftMotor, err := r.queries.LiftMotorGet(ctx, db)
	if err != nil {
		return model.LiftMotor{}, fmt.Errorf("queries get lift motor: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, liftMotor.UpdatedAt)
	if err != nil {
		return model.LiftMotor{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.LiftMotor{
		CurrentPosition: uint16(liftMotor.CurrentPosition),
		TargetPosition:  uint16(liftMotor.TargetPosition),
		IsRunning:       liftMotor.IsRunning == 1,
		Enabled:         liftMotor.Enabled == 1,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r LiftMotorRepository) UpdateLiftMotor(ctx context.Context, db db.SQLDB, liftMotor model.LiftMotor) error {
	params := sqlc.LiftMotorUpdateParams{
		CurrentPosition: int64(liftMotor.CurrentPosition),
		TargetPosition:  int64(liftMotor.TargetPosition),
		IsRunning:       boolToInt64(liftMotor.IsRunning),
		Enabled:         boolToInt64(liftMotor.Enabled),
		UpdatedAt:       liftMotor.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.LiftMotorUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update lift motor: %w", err)
	}

	return nil
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
