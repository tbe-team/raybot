package events

import "github.com/maniartech/signals"

type CloudConnectedEvent struct {
}

type CloudDisconnectedEvent struct {
	Error error
}

type ESPSerialConnectedEvent struct {
}

type ESPSerialDisconnectedEvent struct {
	Error error
}

type PICSerialConnectedEvent struct {
}

type PICSerialDisconnectedEvent struct {
	Error error
}

type RFIDUSBConnectedEvent struct {
}

type RFIDUSBDisconnectedEvent struct {
	Error error
}

var CloudConnectedSignal = signals.New[CloudConnectedEvent]()
var CloudDisconnectedSignal = signals.New[CloudDisconnectedEvent]()

var ESPSerialConnectedSignal = signals.New[ESPSerialConnectedEvent]()
var ESPSerialDisconnectedSignal = signals.New[ESPSerialDisconnectedEvent]()

var PICSerialConnectedSignal = signals.New[PICSerialConnectedEvent]()
var PICSerialDisconnectedSignal = signals.New[PICSerialDisconnectedEvent]()

var RFIDUSBConnectedSignal = signals.New[RFIDUSBConnectedEvent]()
var RFIDUSBDisconnectedSignal = signals.New[RFIDUSBDisconnectedEvent]()
