package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/peripheral"
)

type peripheralHandler struct {
	peripheralService peripheral.Service
}

func newPeripheralHandler(peripheralService peripheral.Service) *peripheralHandler {
	return &peripheralHandler{
		peripheralService: peripheralService,
	}
}

func (h peripheralHandler) ListAvailableSerialPorts(ctx context.Context, _ gen.ListAvailableSerialPortsRequestObject) (gen.ListAvailableSerialPortsResponseObject, error) {
	ports, err := h.peripheralService.ListAvailableSerialPorts(ctx)
	if err != nil {
		return nil, fmt.Errorf("peripheral service list available serial ports: %w", err)
	}

	items := make([]gen.SerialPort, len(ports))
	for i, port := range ports {
		items[i] = gen.SerialPort{
			Port: port.Port,
		}
	}

	return gen.ListAvailableSerialPorts200JSONResponse(gen.SerialPortListResponse{
		Items: items,
	}), nil
}
