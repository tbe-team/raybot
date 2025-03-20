package repoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type LocationRepository struct {
	queries *sqlc.Queries
}

func NewLocationRepository(queries *sqlc.Queries) *LocationRepository {
	return &LocationRepository{queries: queries}
}

func (r LocationRepository) UpdateLocation(ctx context.Context, db db.SQLDB, location model.Location) error {
	params := sqlc.LocationUpdateParams{
		CurrentLocation: location.CurrentLocation,
		UpdatedAt:       location.UpdatedAt.Format(time.RFC3339),
	}
	if err := r.queries.LocationUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update location: %w", err)
	}

	return nil
}

func (r LocationRepository) GetCurrentLocation(ctx context.Context, db db.SQLDB) (model.Location, error) {
	location, err := r.queries.LocationGetCurrent(ctx, db)
	if err != nil {
		return model.Location{}, fmt.Errorf("queries get current location: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, location.UpdatedAt)
	if err != nil {
		return model.Location{}, fmt.Errorf("parse updated at: %w", err)
	}

	return model.Location{
		CurrentLocation: location.CurrentLocation,
		UpdatedAt:       updatedAt,
	}, nil
}
