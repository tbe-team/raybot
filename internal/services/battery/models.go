package battery

import "time"

//nolint:revive
type BatteryState struct {
	Current      uint16 // unit: mA
	Temp         uint8
	Voltage      uint16   // unit: mV
	CellVoltages []uint16 // unit: mV
	Percent      uint8
	Fault        uint8
	Health       uint8
	UpdatedAt    time.Time
}

func (s BatteryState) IsZero() bool {
	return s.Current == 0 &&
		s.Temp == 0 &&
		s.Voltage == 0 &&
		isAllZero(s.CellVoltages) &&
		s.Percent == 0 &&
		s.Fault == 0 &&
		s.Health == 0 &&
		s.UpdatedAt.IsZero()
}

func isAllZero(arr []uint16) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

type ChargeSetting struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}

type DischargeSetting struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}
