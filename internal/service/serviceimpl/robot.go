package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/validator"
)

type robotService struct {
	robotStateRepo repository.RobotStateRepository
	validator      validator.Validator
}

func newRobotService(
	robotStateRepo repository.RobotStateRepository,
	validator validator.Validator,
) *robotService {
	return &robotService{
		robotStateRepo: robotStateRepo,
		validator:      validator,
	}
}

func (s *robotService) UpdateRobotState(ctx context.Context, params service.UpdateRobotStateParams) (model.RobotState, error) {
	if err := s.validator.Validate(params); err != nil {
		return model.RobotState{}, fmt.Errorf("validate params: %w", err)
	}

	state, err := s.robotStateRepo.GetRobotState(ctx)
	if err != nil {
		return model.RobotState{}, fmt.Errorf("get robot state: %w", err)
	}

	now := time.Now()
	if params.SetBattery {
		state.Battery = model.BatteryState{
			Current:      params.Battery.Current,
			Temp:         params.Battery.Temp,
			Voltage:      params.Battery.Voltage,
			CellVoltages: params.Battery.CellVoltages,
			Percent:      params.Battery.Percent,
			Fault:        params.Battery.Fault,
			Health:       params.Battery.Health,
			Status:       params.Battery.Status,
			UpdatedAt:    now,
		}
	}
	if params.SetCharge {
		state.Charge = model.ChargeState{
			CurrentLimit: params.Charge.CurrentLimit,
			Enabled:      params.Charge.Enabled,
			UpdatedAt:    now,
		}
	}

	if params.SetDischarge {
		state.Discharge = model.DischargeState{
			CurrentLimit: params.Discharge.CurrentLimit,
			Enabled:      params.Discharge.Enabled,
			UpdatedAt:    now,
		}
	}

	if params.SetDistanceSensor {
		state.DistanceSensor = model.DistanceSensorState{
			FrontDistance: params.DistanceSensor.FrontDistance,
			BackDistance:  params.DistanceSensor.BackDistance,
			DownDistance:  params.DistanceSensor.DownDistance,
			UpdatedAt:     now,
		}
	}

	if params.SetLiftMotor {
		state.LiftMotor = model.LiftMotorState{
			Direction: params.LiftMotor.Direction,
			Speed:     params.LiftMotor.Speed,
			UpdatedAt: now,
		}
	}

	if params.SetDriveMotor {
		state.DriveMotor = model.DriveMotorState{
			Direction: params.DriveMotor.Direction,
			Speed:     params.DriveMotor.Speed,
			UpdatedAt: now,
		}
	}

	return state, s.robotStateRepo.UpdateRobotState(ctx, state)
}
