package handler

import (
	"context"

	"buf.build/gen/go/tbe-team/raybot-api/grpc/go/raybot/v1/raybotv1grpc"
	raybotv1 "buf.build/gen/go/tbe-team/raybot-api/protocolbuffers/go/raybot/v1"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type DriveMotorHandler struct {
	raybotv1grpc.UnimplementedDriveMotorServiceServer

	picService service.PICService
}

func NewDriveMotorHandler(picService service.PICService) *DriveMotorHandler {
	return &DriveMotorHandler{
		picService: picService,
	}
}

func (h DriveMotorHandler) SetDriveMotorConfiguration(ctx context.Context, req *raybotv1.SetDriveMotorConfigurationRequest) (*raybotv1.SetDriveMotorConfigurationResponse, error) {
	if req.Speed > 100 {
		return nil, xerror.ValidationFailed(nil, "speed must be less than or equal to 100")
	}

	direction := model.MoveDirectionForward
	if req.Direction == raybotv1.SetDriveMotorConfigurationRequest_DIRECTION_BACKWARD {
		direction = model.MoveDirectionBackward
	}

	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryDriveMotorData{
			Direction: direction,
			Speed:     uint8(req.Speed),
			Enable:    req.Enabled,
		},
	}
	if err := h.picService.CreateSerialCommand(ctx, params); err != nil {
		return nil, err
	}

	return &raybotv1.SetDriveMotorConfigurationResponse{}, nil
}
