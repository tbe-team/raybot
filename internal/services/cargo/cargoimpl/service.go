package cargoimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/pkg/validator"
)

const (
	openCargoDoorSpeed  = 50
	closeCargoDoorSpeed = 50
)

type service struct {
	validator validator.Validator

	cargoRepo           cargo.Repository
	espSerialController espserial.Controller
}

func NewService(
	validator validator.Validator,
	cargoRepo cargo.Repository,
	espSerialController espserial.Controller,
) cargo.Service {
	return &service{
		validator:           validator,
		cargoRepo:           cargoRepo,
		espSerialController: espSerialController,
	}
}

func (s *service) UpdateCargoDoor(ctx context.Context, params cargo.UpdateCargoDoorParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoor(ctx, params)
}

func (s *service) UpdateCargoQRCode(ctx context.Context, params cargo.UpdateCargoQRCodeParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.cargoRepo.UpdateCargoQRCode(ctx, params)
}

func (s *service) UpdateCargoBottomDistance(ctx context.Context, params cargo.UpdateCargoBottomDistanceParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.cargoRepo.UpdateCargoBottomDistance(ctx, params)
}

func (s *service) UpdateCargoDoorMotorState(ctx context.Context, params cargo.UpdateCargoDoorMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, params)
}

func (s *service) OpenCargoDoor(ctx context.Context) error {
	if err := s.espSerialController.OpenCargoDoor(openCargoDoorSpeed); err != nil {
		return fmt.Errorf("open cargo door: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
		Direction:    cargo.DirectionOpen,
		SetDirection: true,
		Speed:        openCargoDoorSpeed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}

func (s *service) CloseCargoDoor(ctx context.Context) error {
	if err := s.espSerialController.CloseCargoDoor(closeCargoDoorSpeed); err != nil {
		return fmt.Errorf("close cargo door: %w", err)
	}

	return s.cargoRepo.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
		Direction:    cargo.DirectionClose,
		SetDirection: true,
		Speed:        closeCargoDoorSpeed,
		SetSpeed:     true,
		IsRunning:    true,
		SetIsRunning: true,
		Enabled:      true,
		SetEnabled:   true,
	})
}
