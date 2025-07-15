package config

import "fmt"

type Monitoring struct {
	Battery BatteryMonitoring `yaml:"battery"`
}

func (m *Monitoring) Validate() error {
	if err := m.Battery.Validate(); err != nil {
		return fmt.Errorf("battery: %w", err)
	}

	return nil
}

type BatteryMonitoring struct {
	BatteryVoltageLow      BatteryVoltageLow      `yaml:"voltage_low"`
	BatteryVoltageHigh     BatteryVoltageHigh     `yaml:"voltage_high"`
	BatteryCellVoltageHigh BatteryCellVoltageHigh `yaml:"cell_voltage_high"`
	BatteryCellVoltageLow  BatteryCellVoltageLow  `yaml:"cell_voltage_low"`
	BatteryCellVoltageDiff BatteryCellVoltageDiff `yaml:"cell_voltage_diff"`
	BatteryCurrentHigh     BatteryCurrentHigh     `yaml:"current_high"`
	BatteryTempHigh        BatteryTempHigh        `yaml:"temp_high"`
	BatteryPercentLow      BatteryPercentLow      `yaml:"percent_low"`
	BatteryHealthLow       BatteryHealthLow       `yaml:"health_low"`
}

func (m *BatteryMonitoring) Validate() error {
	if err := m.BatteryVoltageLow.Validate(); err != nil {
		return fmt.Errorf("voltage_low: %w", err)
	}

	if err := m.BatteryVoltageHigh.Validate(); err != nil {
		return fmt.Errorf("battery_voltage_high: %w", err)
	}

	if err := m.BatteryCellVoltageHigh.Validate(); err != nil {
		return fmt.Errorf("battery_cell_voltage_high: %w", err)
	}

	if err := m.BatteryCellVoltageLow.Validate(); err != nil {
		return fmt.Errorf("battery_cell_voltage_low: %w", err)
	}

	if err := m.BatteryCellVoltageDiff.Validate(); err != nil {
		return fmt.Errorf("battery_cell_voltage_diff: %w", err)
	}

	if err := m.BatteryCurrentHigh.Validate(); err != nil {
		return fmt.Errorf("battery_current_high: %w", err)
	}

	if err := m.BatteryTempHigh.Validate(); err != nil {
		return fmt.Errorf("battery_temp_high: %w", err)
	}

	if err := m.BatteryPercentLow.Validate(); err != nil {
		return fmt.Errorf("battery_percent_low: %w", err)
	}

	if err := m.BatteryHealthLow.Validate(); err != nil {
		return fmt.Errorf("battery_health_low: %w", err)
	}

	return nil
}

type BatteryVoltageLow struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryVoltageLow) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryVoltageHigh struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryVoltageHigh) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryCellVoltageHigh struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryCellVoltageHigh) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryCellVoltageLow struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryCellVoltageLow) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryCellVoltageDiff struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryCellVoltageDiff) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryCurrentHigh struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryCurrentHigh) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryTempHigh struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryTempHigh) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryPercentLow struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryPercentLow) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}

type BatteryHealthLow struct {
	Enable    bool    `yaml:"enable"`
	Threshold float64 `yaml:"threshold"`
}

func (b BatteryHealthLow) Validate() error {
	if b.Threshold < 0 {
		return fmt.Errorf("threshold must be greater than 0")
	}

	return nil
}
