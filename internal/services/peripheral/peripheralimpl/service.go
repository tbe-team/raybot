package peripheralimpl

import (
	"context"

	libserial "go.bug.st/serial"

	"github.com/tbe-team/raybot/internal/services/peripheral"
)

type service struct {
}

func NewService() peripheral.Service {
	return &service{}
}

func (s service) ListAvailableSerialPorts(_ context.Context) ([]peripheral.SerialPort, error) {
	ports, err := libserial.GetPortsList()
	if err != nil {
		return nil, err
	}

	serialPorts := make([]peripheral.SerialPort, len(ports))
	for i, port := range ports {
		serialPorts[i] = peripheral.SerialPort{Port: port}
	}

	return serialPorts, nil
}
