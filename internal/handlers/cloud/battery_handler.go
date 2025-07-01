package cloud

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	batteryv1 "github.com/tbe-team/raybot-api/battery/v1"
	"github.com/tbe-team/raybot/internal/services/battery"
)

type batteryHandler struct {
	batteryv1.UnimplementedBatteryServiceServer
	batteryService battery.Service
}

func newBatteryHandler(batteryService battery.Service) batteryv1.BatteryServiceServer {
	return &batteryHandler{
		batteryService: batteryService,
	}
}

func (h batteryHandler) GetBattery(ctx context.Context, _ *batteryv1.GetBatteryRequest) (*batteryv1.GetBatteryResponse, error) {
	state, err := h.batteryService.GetBatteryState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get battery state: %w", err)
	}

	cellVoltages := make([]uint32, len(state.CellVoltages))
	for i, cell := range state.CellVoltages {
		cellVoltages[i] = uint32(cell)
	}

	return &batteryv1.GetBatteryResponse{
		Current:      uint32(state.Current),
		Temp:         uint32(state.Temp),
		Voltage:      uint32(state.Voltage),
		CellVoltages: cellVoltages,
		Percent:      uint32(state.Percent),
		Fault:        uint32(state.Fault),
		Health:       uint32(state.Health),
		UpdatedAt:    timestamppb.New(state.UpdatedAt),
	}, nil
}
