package distancesensorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator
	publisher eventbus.Publisher

	distanceSensorStateRepo distancesensor.DistanceSensorStateRepository
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	distanceSensorStateRepo distancesensor.DistanceSensorStateRepository,
) distancesensor.Service {
	return &service{
		validator:               validator,
		publisher:               publisher,
		distanceSensorStateRepo: distanceSensorStateRepo,
	}
}

func (s *service) UpdateDistanceSensorState(ctx context.Context, params distancesensor.UpdateDistanceSensorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.distanceSensorStateRepo.UpdateDistanceSensorState(ctx, params); err != nil {
		return fmt.Errorf("update distance sensor state: %w", err)
	}

	s.publisher.Publish(events.DistanceSensorUpdatedTopic, eventbus.NewMessage(events.UpdateDistanceSensorEvent{
		FrontDistance: params.FrontDistance,
		BackDistance:  params.BackDistance,
		DownDistance:  params.DownDistance,
	}))

	return nil
}
