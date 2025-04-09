package location

import (
	"context"
)

type UpdateLocationParams struct {
	CurrentLocation string
}

type Service interface {
	UpdateLocation(ctx context.Context, params UpdateLocationParams) error
}

type Repository interface {
	GetLocation(ctx context.Context) (Location, error)
	UpdateLocation(ctx context.Context, location string) error
}
