package drivemotorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
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

	if err := s.driveMotorStateRepo.UpdateDriveMotorState(ctx, params); err != nil {
		return fmt.Errorf("update drive motor state: %w", err)
	}

	ev := events.UpdateDriveMotorStateEvent{
		Direction: events.MoveDirectionForward,
		Speed:     params.Speed,
		Enable:    params.Enabled,
	}
	if params.Direction == drivemotor.DirectionBackward {
		ev.Direction = events.MoveDirectionBackward
	}
	events.UpdateDriveMotorStateSignal.Emit(ctx, ev)

	return nil
}
