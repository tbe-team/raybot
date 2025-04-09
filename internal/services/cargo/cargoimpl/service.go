package cargoimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	cargoRepo cargo.Repository
}

func NewService(
	validator validator.Validator,
	cargoRepo cargo.Repository,
) cargo.Service {
	return &service{
		validator: validator,
		cargoRepo: cargoRepo,
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
