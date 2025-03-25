package service

import "context"

type UpdateLocationParams struct {
	CurrentLocation string
}

type LocationService interface {
	UpdateLocation(ctx context.Context, params UpdateLocationParams) error
}
