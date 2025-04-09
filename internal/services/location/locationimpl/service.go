package locationimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	locationRepo location.Repository
}

func NewService(
	validator validator.Validator,
	locationRepo location.Repository,
) location.Service {
	return &service{
		validator:    validator,
		locationRepo: locationRepo,
	}
}

func (s *service) UpdateLocation(ctx context.Context, params location.UpdateLocationParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.locationRepo.UpdateLocation(ctx, params.CurrentLocation); err != nil {
		return fmt.Errorf("update location: %w", err)
	}

	events.UpdateLocationSignal.Emit(ctx, events.UpdateLocationEvent{
		CurrentLocation: params.CurrentLocation,
	})

	return nil
}
