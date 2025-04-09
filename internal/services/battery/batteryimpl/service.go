package batteryimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	batteryStateRepo battery.BatteryStateRepository
	settingRepo      battery.SettingRepository
}

func NewService(
	validator validator.Validator,
	repo battery.BatteryStateRepository,
	settingRepo battery.SettingRepository,
) battery.Service {
	return &service{
		validator:        validator,
		batteryStateRepo: repo,
		settingRepo:      settingRepo,
	}
}

func (s service) UpdateBatteryState(ctx context.Context, params battery.UpdateBatteryStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.batteryStateRepo.UpdateBatteryState(ctx, params)
}

func (s service) UpdateChargeSetting(ctx context.Context, params battery.UpdateChargeSettingParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.settingRepo.UpdateChargeSetting(ctx, params)
}

func (s service) UpdateDischargeSetting(ctx context.Context, params battery.UpdateDischargeSettingParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.settingRepo.UpdateDischargeSetting(ctx, params)
}
