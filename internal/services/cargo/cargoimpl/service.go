package cargoimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator
	publisher eventbus.Publisher

	cargoRepo           cargo.Repository
	espSerialController espserial.Controller
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	cargoRepo cargo.Repository,
	espSerialController espserial.Controller,
) cargo.Service {
	return &service{
		validator:           validator,
		publisher:           publisher,
		cargoRepo:           cargoRepo,
		espSerialController: espSerialController,
	}
}

func (s *service) GetCargo(ctx context.Context) (cargo.Cargo, error) {
	return s.cargoRepo.GetCargo(ctx)
}

func (s *service) UpdateCargoDoor(ctx context.Context, params cargo.UpdateCargoDoorParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.cargoRepo.UpdateCargoDoor(ctx, params); err != nil {
		return fmt.Errorf("update cargo door: %w", err)
	}

	s.publisher.Publish(events.CargoDoorUpdatedTopic, eventbus.NewMessage(
		events.CargoDoorUpdatedEvent{
			IsOpen: params.IsOpen,
		},
	))

	return nil
}

func (s *service) UpdateCargoQRCode(ctx context.Context, params cargo.UpdateCargoQRCodeParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.cargoRepo.UpdateCargoQRCode(ctx, params); err != nil {
		return fmt.Errorf("update cargo qr code: %w", err)
	}

	s.publisher.Publish(events.CargoQRCodeUpdatedTopic, eventbus.NewMessage(
		events.CargoQRCodeUpdatedEvent{
			QRCode: params.QRCode,
		},
	))

	return nil
}

func (s *service) UpdateCargoBottomDistance(ctx context.Context, params cargo.UpdateCargoBottomDistanceParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.cargoRepo.UpdateCargoBottomDistance(ctx, params); err != nil {
		return fmt.Errorf("update cargo bottom distance: %w", err)
	}

	s.publisher.Publish(events.CargoBottomDistanceUpdatedTopic, eventbus.NewMessage(
		events.CargoBottomDistanceUpdatedEvent{
			BottomDistance: params.BottomDistance,
		},
	))

	return nil
}

func (s *service) UpdateCargoDoorMotorState(ctx context.Context, params cargo.UpdateCargoDoorMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, params)
}

func (s *service) OpenCargoDoor(ctx context.Context, params cargo.OpenCargoDoorParams) error {
	if err := s.espSerialController.OpenCargoDoor(ctx, params.Speed); err != nil {
		return fmt.Errorf("open cargo door: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
		Direction:    cargo.DirectionOpen,
		SetDirection: true,
		Speed:        params.Speed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}

func (s *service) CloseCargoDoor(ctx context.Context, params cargo.CloseCargoDoorParams) error {
	if err := s.espSerialController.CloseCargoDoor(ctx, params.Speed); err != nil {
		return fmt.Errorf("close cargo door: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
		Direction:    cargo.DirectionClose,
		SetDirection: true,
		Speed:        params.Speed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}
