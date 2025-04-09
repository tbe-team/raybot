package distancesensorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	distanceSensorStateRepo distancesensor.DistanceSensorStateRepository
}

func NewService(
	validator validator.Validator,
	distanceSensorStateRepo distancesensor.DistanceSensorStateRepository,
) distancesensor.Service {
	return &service{
		validator:               validator,
		distanceSensorStateRepo: distanceSensorStateRepo,
	}
}

func (s *service) UpdateDistanceSensorState(ctx context.Context, params distancesensor.UpdateDistanceSensorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.distanceSensorStateRepo.UpdateDistanceSensorState(ctx, params)
}
