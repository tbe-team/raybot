package events

const (
	CloudConnectedTopic    = "cloud:connected"
	CloudDisconnectedTopic = "cloud:disconnected"

	ESPSerialConnectedTopic    = "espserial:connected"
	ESPSerialDisconnectedTopic = "espserial:disconnected"

	PICSerialConnectedTopic    = "picserial:connected"
	PICSerialDisconnectedTopic = "picserial:disconnected"

	RFIDUSBConnectedTopic    = "rfidusb:connected"
	RFIDUSBDisconnectedTopic = "rfidusb:disconnected"
)

type CloudConnectedEvent struct{}

type CloudDisconnectedEvent struct {
	Error error
}

type ESPSerialConnectedEvent struct{}

type ESPSerialDisconnectedEvent struct {
	Error error
}

type PICSerialConnectedEvent struct{}

type PICSerialDisconnectedEvent struct {
	Error error
}

type RFIDUSBConnectedEvent struct{}

type RFIDUSBDisconnectedEvent struct {
	Error error
}
