package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type DistanceSensorRepository interface {
	GetDistanceSensor(ctx context.Context, db db.SQLDB) (model.DistanceSensor, error)
	UpdateDistanceSensor(ctx context.Context, db db.SQLDB, distanceSensor model.DistanceSensor) error
}
