package drivemotorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator
	publisher eventbus.Publisher

	driveMotorStateRepo drivemotor.DriveMotorStateRepository
	picSerialController picserial.Controller
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	driveMotorStateRepo drivemotor.DriveMotorStateRepository,
	picSerialController picserial.Controller,
) drivemotor.Service {
	return &service{
		validator:           validator,
		publisher:           publisher,
		driveMotorStateRepo: driveMotorStateRepo,
		picSerialController: picSerialController,
	}
}

func (s service) UpdateDriveMotorState(ctx context.Context, params drivemotor.UpdateDriveMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.driveMotorStateRepo.UpdateDriveMotorState(ctx, params); err != nil {
		return fmt.Errorf("update drive motor state: %w", err)
	}

	s.publisher.Publish(events.DriveMotorUpdatedTopic, eventbus.NewMessage(
		events.DriveMotorStateUpdatedEvent{
			Direction: params.Direction,
			Speed:     params.Speed,
			IsRunning: params.IsRunning,
			Enabled:   params.Enabled,
		},
	))

	return nil
}

func (s service) MoveForward(_ context.Context, params drivemotor.MoveForwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.MoveForward(params.Speed); err != nil {
		return fmt.Errorf("move forward: %w", err)
	}

	return nil
}

func (s service) MoveBackward(_ context.Context, params drivemotor.MoveBackwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.MoveBackward(params.Speed); err != nil {
		return fmt.Errorf("move backward: %w", err)
	}

	return nil
}

func (s service) Stop(_ context.Context) error {
	if err := s.picSerialController.StopDriveMotor(); err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
