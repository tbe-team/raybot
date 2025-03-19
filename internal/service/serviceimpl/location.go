package serviceimpl

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type LocationService struct {
	locationRepo repository.LocationRepository
	dbProvider   db.Provider
}

func NewLocationService(locationRepo repository.LocationRepository, dbProvider db.Provider) *LocationService {
	return &LocationService{locationRepo: locationRepo, dbProvider: dbProvider}
}

func (s LocationService) UpdateLocation(ctx context.Context, params service.UpdateLocationParams) error {
	return s.locationRepo.UpdateLocation(ctx, s.dbProvider.DB(), model.Location{
		CurrentLocation: params.CurrentLocation,
		UpdatedAt:       time.Now(),
	})
}
