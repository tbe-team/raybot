package model

import "time"

type Battery struct {
	Current      uint16
	Temp         uint8
	Voltage      uint16
	CellVoltages []uint16
	Percent      uint8
	Fault        uint8
	Health       uint8
	UpdatedAt    time.Time
}

type BatteryCharge struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}

type BatteryDischarge struct {
	CurrentLimit uint16
	Enabled      bool
	UpdatedAt    time.Time
}
