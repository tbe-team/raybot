package cloud

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	motorv1 "github.com/tbe-team/raybot-api/motor/v1"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type motorHandler struct {
	motorv1.UnimplementedMotorServiceServer
	driveMotorService drivemotor.Service
	liftMotorService  liftmotor.Service
	cargoService      cargo.Service
}

func newMotorHandler(
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	cargoService cargo.Service,
) *motorHandler {
	return &motorHandler{
		driveMotorService: driveMotorService,
		liftMotorService:  liftMotorService,
		cargoService:      cargoService,
	}
}

func (h motorHandler) GetDriveMotor(ctx context.Context, _ *motorv1.GetDriveMotorRequest) (*motorv1.GetDriveMotorResponse, error) {
	state, err := h.driveMotorService.GetDriveMotorState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get drive motor state: %w", err)
	}

	direction := motorv1.GetDriveMotorResponse_DIRECTION_FORWARD
	if state.Direction == drivemotor.DirectionBackward {
		direction = motorv1.GetDriveMotorResponse_DIRECTION_BACKWARD
	}

	return &motorv1.GetDriveMotorResponse{
		Direction: direction,
		Speed:     uint32(state.Speed),
		IsRunning: state.IsRunning,
		Enabled:   state.Enabled,
		UpdatedAt: timestamppb.New(state.UpdatedAt),
	}, nil
}

func (h motorHandler) GetLiftMotor(ctx context.Context, _ *motorv1.GetLiftMotorRequest) (*motorv1.GetLiftMotorResponse, error) {
	state, err := h.liftMotorService.GetLiftMotorState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get lift motor state: %w", err)
	}

	return &motorv1.GetLiftMotorResponse{
		CurrentPosition: uint32(state.CurrentPosition),
		TargetPosition:  uint32(state.TargetPosition),
		IsRunning:       state.IsRunning,
		Enabled:         state.Enabled,
		UpdatedAt:       timestamppb.New(state.UpdatedAt),
	}, nil
}

func (h motorHandler) GetCargoDoorMotor(ctx context.Context, _ *motorv1.GetCargoDoorMotorRequest) (*motorv1.GetCargoDoorMotorResponse, error) {
	state, err := h.cargoService.GetCargoDoorMotorState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cargo door motor state: %w", err)
	}

	direction := motorv1.GetCargoDoorMotorResponse_DIRECTION_OPEN
	if state.Direction == cargo.DirectionClose {
		direction = motorv1.GetCargoDoorMotorResponse_DIRECTION_CLOSE
	}

	return &motorv1.GetCargoDoorMotorResponse{
		Direction: direction,
		Speed:     uint32(state.Speed),
		IsRunning: state.IsRunning,
		Enabled:   state.Enabled,
		UpdatedAt: timestamppb.New(state.UpdatedAt),
	}, nil
}
