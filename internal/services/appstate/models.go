package appstate

import "time"

type AppState struct {
	CloudConnection     CloudConnection
	ESPSerialConnection ESPSerialConnection
	PICSerialConnection PICSerialConnection
	RFIDUSBConnection   RFIDUSBConnection
}

type CloudConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

func (c CloudConnection) ServiceInitialized() bool {
	return c.LastConnectedAt != nil || c.Error != nil
}

type ESPSerialConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

func (c ESPSerialConnection) ServiceInitialized() bool {
	return c.LastConnectedAt != nil || c.Error != nil
}

type PICSerialConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

func (c PICSerialConnection) ServiceInitialized() bool {
	return c.LastConnectedAt != nil || c.Error != nil
}

type RFIDUSBConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

func (c RFIDUSBConnection) ServiceInitialized() bool {
	return c.LastConnectedAt != nil || c.Error != nil
}
