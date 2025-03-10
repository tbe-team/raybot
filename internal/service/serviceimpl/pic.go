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

type PICService struct {
	robotStateRepo       repository.RobotStateRepository
	picCommandSerialRepo repository.PICSerialCommandRepository
	validator            validator.Validator
}

func NewPICService(
	robotStateRepo repository.RobotStateRepository,
	picCommandSerialRepo repository.PICSerialCommandRepository,
	validator validator.Validator,
) *PICService {
	return &PICService{
		robotStateRepo:       robotStateRepo,
		picCommandSerialRepo: picCommandSerialRepo,
		validator:            validator,
	}
}

func (s PICService) ProcessSerialCommandACK(ctx context.Context, params service.ProcessSerialCommandACK) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	// Early return if command failed
	if !params.Success {
		return fmt.Errorf("pic serial command failed: %s", params.ID)
	}

	cmd, err := s.picCommandSerialRepo.GetPICSerialCommand(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("get pic serial command: %w", err)
	}

	robotState, err := s.robotStateRepo.GetRobotState(ctx)
	if err != nil {
		return fmt.Errorf("get robot state: %w", err)
	}

	// Update robot state by command type
	switch cmd.Type {
	case model.PICSerialCommandTypeBatteryCharge:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryChargeData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}
		robotState.Charge.CurrentLimit = data.CurrentLimit
		robotState.Charge.Enabled = data.Enable
		robotState.Charge.UpdatedAt = time.Now()

	case model.PICSerialCommandTypeBatteryDischarge:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryDischargeData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		robotState.Discharge.CurrentLimit = data.CurrentLimit
		robotState.Discharge.Enabled = data.Enable
		robotState.Discharge.UpdatedAt = time.Now()

	case model.PICSerialCommandTypeLiftMotor:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryLiftMotorData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		robotState.LiftMotor.TargetPosition = data.TargetPosition
		robotState.LiftMotor.Enabled = data.Enable
		robotState.LiftMotor.UpdatedAt = time.Now()

	case model.PICSerialCommandTypeDriveMotor:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryDriveMotorData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		switch data.Direction {
		case model.MoveDirectionForward:
			robotState.DriveMotor.Direction = model.DriveMotorDirectionForward
		case model.MoveDirectionBackward:
			robotState.DriveMotor.Direction = model.DriveMotorDirectionBackward
		default:
			return fmt.Errorf("invalid move direction: %d", data.Direction)
		}
		robotState.DriveMotor.Speed = data.Speed
		robotState.DriveMotor.Enabled = data.Enable
		robotState.DriveMotor.UpdatedAt = time.Now()

	default:
		return fmt.Errorf("unknown command type: %d", cmd.Type)
	}

	if err := s.robotStateRepo.UpdateRobotState(ctx, robotState); err != nil {
		return fmt.Errorf("update robot state %w", err)
	}

	if err := s.picCommandSerialRepo.DeletePICSerialCommand(ctx, params.ID); err != nil {
		return fmt.Errorf("delete pic serial command: %w", err)
	}

	return nil
}
