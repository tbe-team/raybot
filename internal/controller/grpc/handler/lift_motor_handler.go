package handler

import (
	"context"
	"fmt"
	"math"

	"buf.build/gen/go/tbe-team/raybot-api/grpc/go/raybot/v1/raybotv1grpc"
	raybotv1 "buf.build/gen/go/tbe-team/raybot-api/protocolbuffers/go/raybot/v1"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type LiftMotorHandler struct {
	raybotv1grpc.UnimplementedLiftMotorServiceServer

	picService service.PICService
}

func NewLiftMotorHandler(picService service.PICService) *LiftMotorHandler {
	return &LiftMotorHandler{
		picService: picService,
	}
}

func (h LiftMotorHandler) SetLiftMotorConfiguration(ctx context.Context, req *raybotv1.SetLiftMotorConfigurationRequest) (*raybotv1.SetLiftMotorConfigurationResponse, error) {
	if req.TargetPosition > math.MaxUint16 {
		return nil, xerror.ValidationFailed(nil, fmt.Sprintf("target position must be less than or equal to %d", math.MaxUint16))
	}

	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryLiftMotorData{
			TargetPosition: uint16(req.TargetPosition),
			Enable:         req.Enabled,
		},
	}
	if err := h.picService.CreateSerialCommand(ctx, params); err != nil {
		return nil, err
	}

	return &raybotv1.SetLiftMotorConfigurationResponse{}, nil
}
