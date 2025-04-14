package events

const (
	CargoDoorUpdatedTopic   = "cargo:door:updated"
	CargoQRCodeUpdatedTopic = "cargo:qrcode:updated"
)

type CargoDoorUpdatedEvent struct {
	IsOpen bool
}

type CargoQRCodeUpdatedEvent struct {
	QRCode string
}
