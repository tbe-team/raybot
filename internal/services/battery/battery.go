package battery

import "context"

type UpdateBatteryStateParams struct {
	Current      uint16 `validate:"min=0"`
	Temp         uint8  `validate:"min=0,max=100"`
	Voltage      uint16 `validate:"min=0"`
	CellVoltages []uint16
	Percent      uint8 `validate:"min=0,max=100"`
	Fault        uint8
	Health       uint8
}

type UpdateChargeSettingParams struct {
	CurrentLimit uint16 `validate:"min=0"`
	Enabled      bool
}

type UpdateDischargeSettingParams struct {
	CurrentLimit uint16 `validate:"min=0"`
	Enabled      bool
}

type Service interface {
	UpdateBatteryState(ctx context.Context, params UpdateBatteryStateParams) error
	UpdateChargeSetting(ctx context.Context, params UpdateChargeSettingParams) error
	UpdateDischargeSetting(ctx context.Context, params UpdateDischargeSettingParams) error
}

//nolint:revive
type BatteryStateRepository interface {
	GetBatteryState(ctx context.Context) (BatteryState, error)
	UpdateBatteryState(ctx context.Context, params UpdateBatteryStateParams) error
}

type SettingRepository interface {
	GetChargeSetting(ctx context.Context) (ChargeSetting, error)
	GetDischargeSetting(ctx context.Context) (DischargeSetting, error)
	UpdateChargeSetting(ctx context.Context, params UpdateChargeSettingParams) error
	UpdateDischargeSetting(ctx context.Context, params UpdateDischargeSettingParams) error
}
