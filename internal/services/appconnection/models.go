package appconnection

import "time"

type AppConnection struct {
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

type ESPSerialConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

type PICSerialConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}

type RFIDUSBConnection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}
