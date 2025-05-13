package liftmotorimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/validator"
)

const (
	OpenCargoDoorSpeed = 100
)

type service struct {
	validator validator.Validator

	liftMotorStateRepo  liftmotor.LiftMotorStateRepository
	picSerialController picserial.Controller
}

func NewService(
	validator validator.Validator,
	liftMotorStateRepo liftmotor.LiftMotorStateRepository,
	picSerialClient picserial.Controller,
) liftmotor.Service {
	return &service{
		validator:           validator,
		liftMotorStateRepo:  liftMotorStateRepo,
		picSerialController: picSerialClient,
	}
}

func (s *service) UpdateLiftMotorState(ctx context.Context, params liftmotor.UpdateLiftMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.liftMotorStateRepo.UpdateLiftMotorState(ctx, params)
}

func (s *service) SetCargoPosition(ctx context.Context, params liftmotor.SetCargoPositionParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.SetCargoPosition(ctx, params.Position); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return liftmotor.ErrCanNotControlLiftMotor
		}
		return fmt.Errorf("set cargo position: %w", err)
	}

	return s.liftMotorStateRepo.UpdateLiftMotorState(ctx, liftmotor.UpdateLiftMotorStateParams{
		TargetPosition:    params.Position,
		SetTargetPosition: true,
		SetIsRunning:      true,
		IsRunning:         true,
		SetEnabled:        true,
		Enabled:           true,
	})
}

func (s *service) Stop(ctx context.Context) error {
	if err := s.picSerialController.StopLiftCargoMotor(ctx); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return liftmotor.ErrCanNotControlLiftMotor
		}
		return fmt.Errorf("stop cargo motor: %w", err)
	}

	return nil
}
