package cargocontrol

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	repo       repository.CargoRepository
	dbProvider db.Provider
	validator  validator.Validator
}

func New(repo repository.CargoRepository, dbProvider db.Provider, validator validator.Validator) *Service {
	return &Service{repo: repo, dbProvider: dbProvider, validator: validator}
}

func (s Service) SyncCargoDoorState(ctx context.Context, params service.SyncCargoDoorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if _, err := s.repo.UpdateCargo(ctx, s.dbProvider.DB(), repository.UpdateCargoParams{
		IsOpen:    params.IsOpen,
		SetIsOpen: true,
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("repository update cargo door state: %w", err)
	}

	return nil
}

func (s Service) SyncCargoQRCode(ctx context.Context, params service.SyncCargoQRCodeParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}
	if _, err := s.repo.UpdateCargo(ctx, s.dbProvider.DB(), repository.UpdateCargoParams{
		QRCode:    params.QRCode,
		SetQRCode: params.QRCode,
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("repository update cargo qr code: %w", err)
	}

	return nil
}

func (s Service) SyncCargoBottomDistance(ctx context.Context, params service.SyncCargoBottomDistanceParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if _, err := s.repo.UpdateCargo(ctx, s.dbProvider.DB(), repository.UpdateCargoParams{
		BottomDistance:    params.BottomDistance,
		SetBottomDistance: params.BottomDistance,
		UpdatedAt:         time.Now(),
	}); err != nil {
		return fmt.Errorf("repository update cargo bottom distance: %w", err)
	}

	return nil
}

func (s Service) SyncCargoDoorMotorState(ctx context.Context, params service.SyncCargoDoorMotorStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if _, err := s.repo.UpdateCargoDoorMotor(ctx, s.dbProvider.DB(), repository.UpdateCargoDoorMotorParams{
		Direction:    params.Direction,
		SetDirection: true,
		Speed:        params.Speed,
		SetSpeed:     true,
		IsRunning:    params.IsRunning,
		SetIsRunning: true,
		Enabled:      params.Enabled,
		SetEnabled:   true,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return fmt.Errorf("repository update cargo door motor state: %w", err)
	}

	return nil
}
