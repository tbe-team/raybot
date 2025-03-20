package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type PICService struct {
	robotStateRepo       repository.RobotStateRepository
	picCommandSerialRepo repository.PICSerialCommandRepository
	picSerialClient      serial.Client
	dbProvider           db.Provider
	validator            validator.Validator
}

func NewPICService(
	robotStateRepo repository.RobotStateRepository,
	picCommandSerialRepo repository.PICSerialCommandRepository,
	picSerialClient serial.Client,
	dbProvider db.Provider,
	validator validator.Validator,
) *PICService {
	return &PICService{
		robotStateRepo:       robotStateRepo,
		picCommandSerialRepo: picCommandSerialRepo,
		picSerialClient:      picSerialClient,
		dbProvider:           dbProvider,
		validator:            validator,
	}
}

func (s PICService) CreateSerialCommand(ctx context.Context, params service.CreateSerialCommandParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	var cmdType model.PICSerialCommandType
	switch params.Data.(type) {
	case model.PICSerialCommandBatteryChargeData:
		cmdType = model.PICSerialCommandTypeBatteryCharge
	case model.PICSerialCommandBatteryDischargeData:
		cmdType = model.PICSerialCommandTypeBatteryDischarge
	case model.PICSerialCommandBatteryLiftMotorData:
		cmdType = model.PICSerialCommandTypeLiftMotor
	case model.PICSerialCommandBatteryDriveMotorData:
		cmdType = model.PICSerialCommandTypeDriveMotor
	}

	cmd := model.PICSerialCommand{
		ID:        shortuuid.New(),
		Type:      cmdType,
		Data:      params.Data,
		CreatedAt: time.Now(),
	}
	if err := s.picCommandSerialRepo.CreatePICSerialCommand(ctx, cmd); err != nil {
		return fmt.Errorf("create pic serial command: %w", err)
	}

	// Send command to PIC serial
	cmdData, err := buildCommandData(params.Data)
	if err != nil {
		return fmt.Errorf("marshal command data: %w", err)
	}

	msg := struct {
		ID   string                     `json:"id"`
		Type model.PICSerialCommandType `json:"type"`
		Data any                        `json:"data"`
	}{
		ID:   cmd.ID,
		Type: cmd.Type,
		Data: cmdData,
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal command message: %w", err)
	}

	if err := s.picSerialClient.Write(msgJSON); err != nil {
		return fmt.Errorf("pic serial client write: %w", err)
	}

	return nil
}

func (s PICService) ProcessSerialCommandACK(ctx context.Context, params service.ProcessSerialCommandACKParams) error {
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

	robotState, err := s.robotStateRepo.GetRobotState(ctx, s.dbProvider.DB())
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

	if err := s.robotStateRepo.UpdateRobotState(ctx, s.dbProvider.DB(), robotState); err != nil {
		return fmt.Errorf("update robot state %w", err)
	}

	if err := s.picCommandSerialRepo.DeletePICSerialCommand(ctx, params.ID); err != nil {
		return fmt.Errorf("delete pic serial command: %w", err)
	}

	return nil
}

func buildCommandData(data model.PICSerialCommandData) (any, error) {
	boolToUint8 := func(b bool) uint8 {
		if b {
			return 1
		}
		return 0
	}

	switch data := data.(type) {
	case model.PICSerialCommandBatteryChargeData:
		return struct {
			CurrentLimit uint16 `json:"current_limit"`
			Enable       uint8  `json:"enable"`
		}{
			CurrentLimit: data.CurrentLimit,
			Enable:       boolToUint8(data.Enable),
		}, nil
	case model.PICSerialCommandBatteryDischargeData:
		return struct {
			CurrentLimit uint16 `json:"current_limit"`
			Enable       uint8  `json:"enable"`
		}{
			CurrentLimit: data.CurrentLimit,
			Enable:       boolToUint8(data.Enable),
		}, nil
	case model.PICSerialCommandBatteryLiftMotorData:
		return struct {
			TargetPosition uint16 `json:"target_position"`
			Enable         uint8  `json:"enable"`
		}{
			TargetPosition: data.TargetPosition,
			Enable:         boolToUint8(data.Enable),
		}, nil
	case model.PICSerialCommandBatteryDriveMotorData:
		return struct {
			Direction uint8 `json:"direction"`
			Speed     uint8 `json:"speed"`
			Enable    uint8 `json:"enable"`
		}{
			Direction: uint8(data.Direction),
			Speed:     data.Speed,
			Enable:    boolToUint8(data.Enable),
		}, nil
	default:
		return nil, fmt.Errorf("unknown command data type: %T", data)
	}
}
