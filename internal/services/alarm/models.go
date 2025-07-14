package alarm

import (
	"fmt"
	"time"
)

//nolint:revive
type AlarmType string

func (t AlarmType) String() string {
	return string(t)
}

func (t AlarmType) Validate() error {
	if _, ok := alarmTypeMap[t]; !ok {
		return fmt.Errorf("invalid alarm type: %s", t)
	}
	return nil
}

var alarmTypeMap = map[AlarmType]struct{}{
	AlarmTypeBatteryVoltageLow:      {},
	AlarmTypeBatteryVoltageHigh:     {},
	AlarmTypeBatteryCellVoltageHigh: {},
	AlarmTypeBatteryCellVoltageLow:  {},
	AlarmTypeBatteryCellVoltageDiff: {},
	AlarmTypeBatteryCurrentHigh:     {},
	AlarmTypeBatteryTempHigh:        {},
	AlarmTypeBatteryPercentLow:      {},
	AlarmTypeBatteryHealthLow:       {},
}

const (
	AlarmTypeBatteryVoltageLow      AlarmType = "battery_voltage_low"
	AlarmTypeBatteryVoltageHigh     AlarmType = "battery_voltage_high"
	AlarmTypeBatteryCellVoltageHigh AlarmType = "battery_cell_voltage_high"
	AlarmTypeBatteryCellVoltageLow  AlarmType = "battery_cell_voltage_low"
	AlarmTypeBatteryCellVoltageDiff AlarmType = "battery_cell_voltage_diff"
	AlarmTypeBatteryCurrentHigh     AlarmType = "battery_current_high"
	AlarmTypeBatteryTempHigh        AlarmType = "battery_temp_high"
	AlarmTypeBatteryPercentLow      AlarmType = "battery_percent_low"
	AlarmTypeBatteryHealthLow       AlarmType = "battery_health_low"
)

type Alarm struct {
	ID            int64
	Type          AlarmType
	Data          Data
	ActivatedAt   time.Time
	DeactivatedAt *time.Time
}
