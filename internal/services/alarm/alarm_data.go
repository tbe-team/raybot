package alarm

import (
	"encoding/json"
	"fmt"
)

type Data interface {
	isAlarmData()
	Message() string
}

type DataBatteryVoltageLow struct {
	Threshold float64 `json:"threshold"`
	Voltage   float64 `json:"voltage"`
}

func (a DataBatteryVoltageLow) Message() string {
	return fmt.Sprintf("Battery voltage is low: %.2f", a.Voltage)
}

func (DataBatteryVoltageLow) isAlarmData() {}

type DataBatteryVoltageHigh struct {
	Threshold float64 `json:"threshold"`
	Voltage   float64 `json:"voltage"`
}

func (DataBatteryVoltageHigh) isAlarmData() {}

func (a DataBatteryVoltageHigh) Message() string {
	return fmt.Sprintf("Battery voltage is high: %.2f", a.Voltage)
}

type DataBatteryCellVoltageHigh struct {
	Threshold          float64   `json:"threshold"`
	CellVoltages       []float64 `json:"cell_voltages"`
	OverThresholdIndex []int     `json:"over_threshold_index"`
}

func (DataBatteryCellVoltageHigh) isAlarmData() {}

func (a DataBatteryCellVoltageHigh) Message() string {
	if len(a.OverThresholdIndex) == 0 {
		return "No battery cell voltage exceeds the threshold"
	}

	msg := "Battery cell(s) over voltage threshold:"
	for _, idx := range a.OverThresholdIndex {
		if idx >= 0 && idx < len(a.CellVoltages) {
			msg += fmt.Sprintf(" [Cell %d: %.2fV]", idx, a.CellVoltages[idx])
		}
	}
	return msg
}

type DataBatteryCellVoltageLow struct {
	Threshold           float64   `json:"threshold"`
	CellVoltages        []float64 `json:"cell_voltages"`
	UnderThresholdIndex []int     `json:"under_threshold_index"`
}

func (DataBatteryCellVoltageLow) isAlarmData() {}

func (a DataBatteryCellVoltageLow) Message() string {
	if len(a.UnderThresholdIndex) == 0 {
		return "No battery cell voltage below the threshold"
	}

	msg := "Battery cell(s) under voltage threshold:"
	for _, idx := range a.UnderThresholdIndex {
		if idx >= 0 && idx < len(a.CellVoltages) {
			msg += fmt.Sprintf(" [Cell %d: %.2fV]", idx, a.CellVoltages[idx])
		}
	}
	return msg
}

type DataBatteryCellVoltageDiff struct {
	Threshold    float64   `json:"threshold"`
	CellVoltages []float64 `json:"cell_voltages"`
	DiffIndex    []int     `json:"diff_index"`
}

func (DataBatteryCellVoltageDiff) isAlarmData() {}

func (a DataBatteryCellVoltageDiff) Message() string {
	if len(a.DiffIndex) == 0 {
		return "No battery cell voltage difference exceeds the threshold"
	}

	msg := "Battery cell(s) voltage difference exceeds the threshold:"
	for _, idx := range a.DiffIndex {
		if idx >= 0 && idx < len(a.CellVoltages) {
			msg += fmt.Sprintf(" [Cell %d: %.2fV]", idx, a.CellVoltages[idx])
		}
	}
	return msg
}

type DataBatteryCurrentHigh struct {
	Threshold float64 `json:"threshold"`
	Current   float64 `json:"current"`
}

func (DataBatteryCurrentHigh) isAlarmData() {}

func (a DataBatteryCurrentHigh) Message() string {
	return fmt.Sprintf("Battery current is high: %.2f", a.Current)
}

type DataBatteryTempHigh struct {
	Threshold float64 `json:"threshold"`
	Temp      float64 `json:"temp"`
}

func (DataBatteryTempHigh) isAlarmData() {}

func (a DataBatteryTempHigh) Message() string {
	return fmt.Sprintf("Battery temperature is high: %.2f", a.Temp)
}

type DataBatteryPercentLow struct {
	Threshold float64 `json:"threshold"`
	Percent   float64 `json:"percent"`
}

func (DataBatteryPercentLow) isAlarmData() {}

func (a DataBatteryPercentLow) Message() string {
	return fmt.Sprintf("Battery percent is low: %.2f", a.Percent)
}

type DataBatteryHealthLow struct {
	Threshold float64 `json:"threshold"`
	Health    float64 `json:"health"`
}

func (DataBatteryHealthLow) isAlarmData() {}

func (a DataBatteryHealthLow) Message() string {
	return fmt.Sprintf("Battery health is low: %.2f", a.Health)
}

func UnmarshalAlarmData(alarmType AlarmType, data []byte) (Data, error) {
	var ret Data

	switch alarmType {
	case AlarmTypeBatteryVoltageLow:
		var d DataBatteryVoltageLow
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery voltage low alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryVoltageHigh:
		var d DataBatteryVoltageHigh
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery voltage high alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryCellVoltageHigh:
		var d DataBatteryCellVoltageHigh
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery cell voltage high alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryCellVoltageLow:
		var d DataBatteryCellVoltageLow
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery cell voltage low alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryCellVoltageDiff:
		var d DataBatteryCellVoltageDiff
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery cell voltage diff alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryCurrentHigh:
		var d DataBatteryCurrentHigh
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery current high alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryTempHigh:
		var d DataBatteryTempHigh
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery temp high alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryPercentLow:
		var d DataBatteryPercentLow
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery percent low alarm data: %w", err)
		}
		ret = d

	case AlarmTypeBatteryHealthLow:
		var d DataBatteryHealthLow
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battery health low alarm data: %w", err)
		}
		ret = d

	default:
		return nil, fmt.Errorf("invalid alarm type: %s", alarmType)
	}

	return ret, nil
}
