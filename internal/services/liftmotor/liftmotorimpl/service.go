package liftmotorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	liftMotorStateRepo liftmotor.LiftMotorStateRepository
}

func NewService(
	validator validator.Validator,
	liftMotorStateRepo liftmotor.LiftMotorStateRepository,
) liftmotor.Service {
	return &service{
		validator:          validator,
		liftMotorStateRepo: liftMotorStateRepo,
	}
}

func (s *service) UpdateLiftMotorState(ctx context.Context, params liftmotor.UpdateLiftMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.liftMotorStateRepo.UpdateLiftMotorState(ctx, params)
}
