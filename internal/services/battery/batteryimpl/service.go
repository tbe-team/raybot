package batteryimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/controller"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	publisher         eventbus.Publisher
	batteryStateRepo  battery.BatteryStateRepository
	settingRepo       battery.SettingRepository
	batteryController controller.BatteryController
}

func NewService(
	validator validator.Validator,
	publisher eventbus.Publisher,
	batteryStateRepo battery.BatteryStateRepository,
	settingRepo battery.SettingRepository,
	batteryController controller.BatteryController,
) battery.Service {
	return &service{
		validator:         validator,
		publisher:         publisher,
		batteryStateRepo:  batteryStateRepo,
		settingRepo:       settingRepo,
		batteryController: batteryController,
	}
}

func (s service) GetBatteryState(ctx context.Context) (battery.BatteryState, error) {
	return s.batteryStateRepo.GetBatteryState(ctx)
}

func (s service) UpdateBatteryState(ctx context.Context, params battery.UpdateBatteryStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	state, err := s.batteryStateRepo.UpdateBatteryState(ctx, params)
	if err != nil {
		return fmt.Errorf("update battery state: %w", err)
	}

	s.publisher.Publish(events.BatteryUpdatedTopic, eventbus.NewMessage(
		events.BatteryUpdatedEvent{
			BatteryState: state,
		},
	))

	return nil
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

func (s service) DisableCharge(ctx context.Context) error {
	setting, err := s.settingRepo.GetChargeSetting(ctx)
	if err != nil {
		return fmt.Errorf("get charge setting: %w", err)
	}

	if !setting.Enabled {
		return nil
	}

	if err := s.batteryController.ConfigBatteryCharge(ctx, setting.CurrentLimit, false); err != nil {
		return fmt.Errorf("disable charge: %w", err)
	}

	if err := s.settingRepo.UpdateChargeSetting(ctx, battery.UpdateChargeSettingParams{
		CurrentLimit: setting.CurrentLimit,
		Enabled:      false,
	}); err != nil {
		return fmt.Errorf("update charge setting: %w", err)
	}

	return nil
}
