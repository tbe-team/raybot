package drivemotorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	driveMotorStateRepo drivemotor.DriveMotorStateRepository
}

func NewService(
	validator validator.Validator,
	driveMotorStateRepo drivemotor.DriveMotorStateRepository,
) drivemotor.Service {
	return &service{
		validator:           validator,
		driveMotorStateRepo: driveMotorStateRepo,
	}
}

func (s service) UpdateDriveMotorState(ctx context.Context, params drivemotor.UpdateDriveMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.driveMotorStateRepo.UpdateDriveMotorState(ctx, params)
}
