package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type RobotService struct {
	robotStateRepo repository.RobotStateRepository
	dbProvider     db.Provider
	validator      validator.Validator
}

func NewRobotService(
	robotStateRepo repository.RobotStateRepository,
	dbProvider db.Provider,
	validator validator.Validator,
) *RobotService {
	return &RobotService{
		robotStateRepo: robotStateRepo,
		dbProvider:     dbProvider,
		validator:      validator,
	}
}

func (s RobotService) GetRobotState(ctx context.Context) (model.RobotState, error) {
	return s.robotStateRepo.GetRobotState(ctx, s.dbProvider.DB())
}

func (s RobotService) UpdateRobotState(ctx context.Context, params service.UpdateRobotStateParams) (model.RobotState, error) {
	if err := s.validator.Validate(params); err != nil {
		return model.RobotState{}, fmt.Errorf("validate params: %w", err)
	}

	state, err := s.robotStateRepo.GetRobotState(ctx, s.dbProvider.DB())
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
			CurrentPosition: params.LiftMotor.CurrentPosition,
			TargetPosition:  params.LiftMotor.TargetPosition,
			IsRunning:       params.LiftMotor.IsRunning,
			Enabled:         params.LiftMotor.Enabled,
			UpdatedAt:       now,
		}
	}

	if params.SetDriveMotor {
		state.DriveMotor = model.DriveMotorState{
			Direction: params.DriveMotor.Direction,
			Speed:     params.DriveMotor.Speed,
			IsRunning: params.DriveMotor.IsRunning,
			Enabled:   params.DriveMotor.Enabled,
			UpdatedAt: now,
		}
	}

	if params.SetLocation {
		state.Location = model.LocationState{
			CurrentLocation: params.Location.CurrentLocation,
			UpdatedAt:       now,
		}
	}

	return state, s.robotStateRepo.UpdateRobotState(ctx, s.dbProvider.DB(), state)
}
