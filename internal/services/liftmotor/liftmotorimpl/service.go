package liftmotorimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

const (
	OpenCargoDoorSpeed = 100
)

type service struct {
	validator validator.Validator

	publisher           eventbus.Publisher
	liftMotorStateRepo  liftmotor.LiftMotorStateRepository
	picSerialController picserial.Controller
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	liftMotorStateRepo liftmotor.LiftMotorStateRepository,
	picSerialClient picserial.Controller,
) liftmotor.Service {
	return &service{
		validator:           validator,
		publisher:           publisher,
		liftMotorStateRepo:  liftMotorStateRepo,
		picSerialController: picSerialClient,
	}
}

func (s *service) GetLiftMotorState(ctx context.Context) (liftmotor.LiftMotorState, error) {
	return s.liftMotorStateRepo.GetLiftMotorState(ctx)
}

func (s *service) UpdateLiftMotorState(ctx context.Context, params liftmotor.UpdateLiftMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.liftMotorStateRepo.UpdateLiftMotorState(ctx, params); err != nil {
		return fmt.Errorf("update lift motor state: %w", err)
	}

	s.publisher.Publish(events.LiftMotorUpdatedTopic, eventbus.NewMessage(
		events.LiftMotorStateUpdatedEvent{
			CurrentPosition: params.CurrentPosition,
			TargetPosition:  params.TargetPosition,
			IsRunning:       params.IsRunning,
			Enabled:         params.Enabled,
		},
	))

	return nil
}

func (s *service) SetCargoPosition(ctx context.Context, params liftmotor.SetCargoPositionParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.picSerialController.SetCargoPosition(ctx, params.MotorSpeed, params.Position); err != nil {
		if errors.Is(err, picserial.ErrPICSerialNotConnected) {
			return liftmotor.ErrCanNotControlLiftMotor
		}
		return fmt.Errorf("set cargo position: %w", err)
	}

	return nil
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
