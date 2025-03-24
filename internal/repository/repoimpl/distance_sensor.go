package repoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type DistanceSensorRepository struct {
	queries *sqlc.Queries
}

func NewDistanceSensorRepository(queries *sqlc.Queries) *DistanceSensorRepository {
	return &DistanceSensorRepository{queries: queries}
}

func (r DistanceSensorRepository) GetDistanceSensor(ctx context.Context, db db.SQLDB) (model.DistanceSensor, error) {
	distanceSensor, err := r.queries.DistanceSensorGet(ctx, db)
	if err != nil {
		return model.DistanceSensor{}, fmt.Errorf("queries get distance sensor: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, distanceSensor.UpdatedAt)
	if err != nil {
		return model.DistanceSensor{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.DistanceSensor{
		FrontDistance: uint16(distanceSensor.FrontDistance),
		BackDistance:  uint16(distanceSensor.BackDistance),
		DownDistance:  uint16(distanceSensor.DownDistance),
		UpdatedAt:     updatedAt,
	}, nil
}

func (r DistanceSensorRepository) UpdateDistanceSensor(ctx context.Context, db db.SQLDB, distanceSensor model.DistanceSensor) error {
	params := sqlc.DistanceSensorUpdateParams{
		FrontDistance: int64(distanceSensor.FrontDistance),
		BackDistance:  int64(distanceSensor.BackDistance),
		DownDistance:  int64(distanceSensor.DownDistance),
		UpdatedAt:     distanceSensor.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.DistanceSensorUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update distance sensor: %w", err)
	}

	return nil
}
