package cloud

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	distanceSensorv1 "github.com/tbe-team/raybot-api/distancesensor/v1"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
)

type distanceSensorHandler struct {
	distanceSensorv1.UnimplementedDistanceSensorServiceServer
	distanceSensorService distancesensor.Service
}

func newDistanceSensorHandler(distanceSensorService distancesensor.Service) distanceSensorv1.DistanceSensorServiceServer {
	return &distanceSensorHandler{
		distanceSensorService: distanceSensorService,
	}
}

func (h distanceSensorHandler) GetDistanceSensor(ctx context.Context, _ *distanceSensorv1.GetDistanceSensorRequest) (*distanceSensorv1.GetDistanceSensorResponse, error) {
	state, err := h.distanceSensorService.GetDistanceSensorState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get distance sensor state: %w", err)
	}

	return &distanceSensorv1.GetDistanceSensorResponse{
		FrontDistance: uint32(state.FrontDistance),
		BackDistance:  uint32(state.BackDistance),
		DownDistance:  uint32(state.DownDistance),
		UpdatedAt:     timestamppb.New(state.UpdatedAt),
	}, nil
}
