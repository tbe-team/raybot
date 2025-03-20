package handler

import (
	"context"

	"buf.build/gen/go/tbe-team/raybot-api/grpc/go/raybot/v1/raybotv1grpc"
	raybotv1 "buf.build/gen/go/tbe-team/raybot-api/protocolbuffers/go/raybot/v1"

	"github.com/tbe-team/raybot/internal/controller/grpc/converter"
	"github.com/tbe-team/raybot/internal/service"
)

type RobotStateHandler struct {
	raybotv1grpc.UnimplementedRobotStateServiceServer

	robotService service.RobotService
}

func NewRobotStateHandler(robotService service.RobotService) *RobotStateHandler {
	return &RobotStateHandler{
		robotService: robotService,
	}
}

func (h RobotStateHandler) GetRobotState(ctx context.Context, _ *raybotv1.GetRobotStateRequest) (*raybotv1.GetRobotStateResponse, error) {
	state, err := h.robotService.GetRobotState(ctx)
	if err != nil {
		return nil, err
	}

	res := converter.ToGetRobotStateResponse(state)
	return res, nil
}
