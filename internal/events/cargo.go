package events

const (
	CargoDoorUpdatedTopic           = "cargo:door:updated"
	CargoQRCodeUpdatedTopic         = "cargo:qrcode:updated"
	CargoBottomDistanceUpdatedTopic = "cargo:bottom_distance:updated"
)

type CargoDoorUpdatedEvent struct {
	IsOpen bool
}

type CargoQRCodeUpdatedEvent struct {
	QRCode string
}

type CargoBottomDistanceUpdatedEvent struct {
	BottomDistance uint16
}
