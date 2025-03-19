package repoimpl

import (
	"context"
	"fmt"

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
		UpdatedAt:       location.UpdatedAt,
	}
	if err := r.queries.LocationUpdate(ctx, db, params); err != nil {
		return fmt.Errorf("queries update location: %w", err)
	}

	return nil
}
