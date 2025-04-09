package locationimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type repository struct {
	db      db.DB
	queries *sqlc.Queries
}

func NewLocationRepository(db db.DB, queries *sqlc.Queries) location.Repository {
	return &repository{
		db:      db,
		queries: queries,
	}
}

func (r repository) GetLocation(ctx context.Context) (location.Location, error) {
	row, err := r.queries.LocationGetCurrent(ctx, r.db)
	if err != nil {
		return location.Location{}, fmt.Errorf("failed to get location: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return location.Location{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	//nolint:gosec
	return location.Location{
		CurrentLocation: row.CurrentLocation,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r repository) UpdateLocation(ctx context.Context, location string) error {
	if err := r.queries.LocationUpdate(ctx, r.db, sqlc.LocationUpdateParams{
		CurrentLocation: location,
		UpdatedAt:       time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("queries update location: %w", err)
	}

	return nil
}
