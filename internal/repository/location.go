package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type LocationRepository interface {
	UpdateLocation(ctx context.Context, db db.SQLDB, location model.Location) error
	GetCurrentLocation(ctx context.Context, db db.SQLDB) (model.Location, error)
}
