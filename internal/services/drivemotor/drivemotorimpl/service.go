package drivemotorimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	driveMotorStateRepo drivemotor.DriveMotorStateRepository
	picSerialController picserial.Controller
}

func NewService(
	validator validator.Validator,
	driveMotorStateRepo drivemotor.DriveMotorStateRepository,
	picSerialController picserial.Controller,
) drivemotor.Service {
	return &service{
		validator:           validator,
		driveMotorStateRepo: driveMotorStateRepo,
		picSerialController: picSerialController,
	}
}

func (s service) UpdateDriveMotorState(ctx context.Context, params drivemotor.UpdateDriveMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.driveMotorStateRepo.UpdateDriveMotorState(ctx, params)
}

func (s service) MoveForward(ctx context.Context, params drivemotor.MoveForwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.MoveForward(params.Speed); err != nil {
		return fmt.Errorf("move forward: %w", err)
	}

	return s.driveMotorStateRepo.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction:    drivemotor.DirectionForward,
		SetDirection: true,
		Speed:        params.Speed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}

func (s service) MoveBackward(ctx context.Context, params drivemotor.MoveBackwardParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.MoveBackward(params.Speed); err != nil {
		return fmt.Errorf("move backward: %w", err)
	}

	return s.driveMotorStateRepo.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction:    drivemotor.DirectionBackward,
		SetDirection: true,
		Speed:        params.Speed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}

func (s service) Stop(ctx context.Context) error {
	if err := s.picSerialController.StopDriveMotor(); err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return s.driveMotorStateRepo.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		IsRunning:    false,
		SetIsRunning: true,
		Enabled:      false,
		SetEnabled:   true,
	})
}
