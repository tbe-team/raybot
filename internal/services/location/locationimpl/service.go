package locationimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	publisher    eventbus.Publisher
	locationRepo location.Repository
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	locationRepo location.Repository,
) location.Service {
	return &service{
		validator:    validator,
		publisher:    publisher,
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

	s.publisher.Publish(
		events.LocationUpdatedTopic,
		eventbus.NewMessage(events.UpdateLocationEvent{
			Location: params.CurrentLocation,
		}),
	)

	return nil
}
