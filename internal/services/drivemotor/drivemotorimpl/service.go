package drivemotorimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/controller"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator
	publisher eventbus.Publisher

	driveMotorStateRepo  drivemotor.DriveMotorStateRepository
	driveMotorController controller.DriveMotorController
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	driveMotorStateRepo drivemotor.DriveMotorStateRepository,
	driveMotorController controller.DriveMotorController,
) drivemotor.Service {
	return &service{
		validator:            validator,
		publisher:            publisher,
		driveMotorStateRepo:  driveMotorStateRepo,
		driveMotorController: driveMotorController,
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

func (s service) MoveForward(ctx context.Context, params drivemotor.MoveForwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.driveMotorController.MoveForward(ctx, params.Speed); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return drivemotor.ErrCanNotControlDriveMotor
		}
		return fmt.Errorf("move forward: %w", err)
	}

	return nil
}

func (s service) MoveBackward(ctx context.Context, params drivemotor.MoveBackwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.driveMotorController.MoveBackward(ctx, params.Speed); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return drivemotor.ErrCanNotControlDriveMotor
		}
		return fmt.Errorf("move backward: %w", err)
	}

	return nil
}

func (s service) Stop(ctx context.Context) error {
	state, err := s.driveMotorStateRepo.GetDriveMotorState(ctx)
	if err != nil {
		return fmt.Errorf("get drive motor state: %w", err)
	}

	var moveForward bool
	switch state.Direction {
	case drivemotor.DirectionForward:
		moveForward = true
	case drivemotor.DirectionBackward:
		moveForward = false
	}

	if err := s.driveMotorController.StopDriveMotor(ctx, moveForward); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return drivemotor.ErrCanNotControlDriveMotor
		}

		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
