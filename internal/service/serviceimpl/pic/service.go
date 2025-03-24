package pic

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

type Service struct {
	robotStateRepo       repository.RobotStateRepository
	picCommandSerialRepo repository.PICSerialCommandRepository
	batteryRepo          repository.BatteryRepository
	distanceSensorRepo   repository.DistanceSensorRepository
	liftMotorRepo        repository.LiftMotorRepository
	driveMotorRepo       repository.DriveMotorRepository
	locationRepo         repository.LocationRepository
	picSerialClient      serial.Client
	dbProvider           db.Provider
	validator            validator.Validator
}

func NewService(
	robotStateRepo repository.RobotStateRepository,
	picCommandSerialRepo repository.PICSerialCommandRepository,
	batteryRepo repository.BatteryRepository,
	distanceSensorRepo repository.DistanceSensorRepository,
	liftMotorRepo repository.LiftMotorRepository,
	driveMotorRepo repository.DriveMotorRepository,
	locationRepo repository.LocationRepository,
	picSerialClient serial.Client,
	dbProvider db.Provider,
	validator validator.Validator,
) *Service {
	return &Service{
		robotStateRepo:       robotStateRepo,
		picCommandSerialRepo: picCommandSerialRepo,
		batteryRepo:          batteryRepo,
		distanceSensorRepo:   distanceSensorRepo,
		liftMotorRepo:        liftMotorRepo,
		driveMotorRepo:       driveMotorRepo,
		locationRepo:         locationRepo,
		picSerialClient:      picSerialClient,
		dbProvider:           dbProvider,
		validator:            validator,
	}
}

func (s Service) ProcessSyncState(ctx context.Context, params service.ProcessSyncStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if params.SetBattery {
		m := model.Battery{
			Current:      params.Battery.Current,
			Temp:         params.Battery.Temp,
			Voltage:      params.Battery.Voltage,
			CellVoltages: params.Battery.CellVoltages,
			Percent:      params.Battery.Percent,
			Fault:        params.Battery.Fault,
			Health:       params.Battery.Health,
		}
		return s.batteryRepo.UpdateBattery(ctx, s.dbProvider.DB(), m)
	}

	if params.SetCharge {
		m := model.BatteryCharge{
			CurrentLimit: params.Charge.CurrentLimit,
			Enabled:      params.Charge.Enabled,
		}
		return s.batteryRepo.UpdateBatteryCharge(ctx, s.dbProvider.DB(), m)
	}

	if params.SetDischarge {
		m := model.BatteryDischarge{
			CurrentLimit: params.Discharge.CurrentLimit,
			Enabled:      params.Discharge.Enabled,
		}
		return s.batteryRepo.UpdateBatteryDischarge(ctx, s.dbProvider.DB(), m)
	}

	if params.SetDistanceSensor {
		m := model.DistanceSensor{
			FrontDistance: params.DistanceSensor.FrontDistance,
			BackDistance:  params.DistanceSensor.BackDistance,
			DownDistance:  params.DistanceSensor.DownDistance,
		}
		return s.distanceSensorRepo.UpdateDistanceSensor(ctx, s.dbProvider.DB(), m)
	}

	if params.SetLiftMotor {
		m := model.LiftMotor{
			TargetPosition: params.LiftMotor.TargetPosition,
			Enabled:        params.LiftMotor.Enabled,
		}
		return s.liftMotorRepo.UpdateLiftMotor(ctx, s.dbProvider.DB(), m)
	}

	if params.SetDriveMotor {
		m := model.DriveMotor{
			Direction: params.DriveMotor.Direction,
			Speed:     params.DriveMotor.Speed,
			Enabled:   params.DriveMotor.Enabled,
		}
		return s.driveMotorRepo.UpdateDriveMotor(ctx, s.dbProvider.DB(), m)
	}

	return nil
}

func (s Service) CreateSerialCommand(ctx context.Context, params service.CreateSerialCommandParams) error {
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

func (s Service) ProcessSerialCommandACK(ctx context.Context, params service.ProcessSerialCommandACKParams) error {
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

	// Update robot state by command type
	switch cmd.Type {
	case model.PICSerialCommandTypeBatteryCharge:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryChargeData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		m, err := s.batteryRepo.GetBatteryCharge(ctx, s.dbProvider.DB())
		if err != nil {
			return fmt.Errorf("get battery charge: %w", err)
		}

		m.CurrentLimit = data.CurrentLimit
		m.Enabled = data.Enable
		m.UpdatedAt = time.Now()
		if err := s.batteryRepo.UpdateBatteryCharge(ctx, s.dbProvider.DB(), m); err != nil {
			return fmt.Errorf("update battery charge: %w", err)
		}

	case model.PICSerialCommandTypeBatteryDischarge:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryDischargeData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		m, err := s.batteryRepo.GetBatteryDischarge(ctx, s.dbProvider.DB())
		if err != nil {
			return fmt.Errorf("get battery discharge: %w", err)
		}

		m.CurrentLimit = data.CurrentLimit
		m.Enabled = data.Enable
		m.UpdatedAt = time.Now()
		if err := s.batteryRepo.UpdateBatteryDischarge(ctx, s.dbProvider.DB(), m); err != nil {
			return fmt.Errorf("update battery discharge: %w", err)
		}

	case model.PICSerialCommandTypeLiftMotor:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryLiftMotorData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		m, err := s.liftMotorRepo.GetLiftMotor(ctx, s.dbProvider.DB())
		if err != nil {
			return fmt.Errorf("get lift motor: %w", err)
		}

		m.TargetPosition = data.TargetPosition
		m.Enabled = data.Enable
		m.UpdatedAt = time.Now()
		if err := s.liftMotorRepo.UpdateLiftMotor(ctx, s.dbProvider.DB(), m); err != nil {
			return fmt.Errorf("update lift motor: %w", err)
		}

	case model.PICSerialCommandTypeDriveMotor:
		data, ok := cmd.Data.(model.PICSerialCommandBatteryDriveMotorData)
		if !ok {
			return fmt.Errorf("invalid command data type: %T", cmd.Data)
		}

		m, err := s.driveMotorRepo.GetDriveMotor(ctx, s.dbProvider.DB())
		if err != nil {
			return fmt.Errorf("get drive motor: %w", err)
		}

		switch data.Direction {
		case model.MoveDirectionForward:
			m.Direction = model.DriveMotorDirectionForward
		case model.MoveDirectionBackward:
			m.Direction = model.DriveMotorDirectionBackward
		default:
			return fmt.Errorf("invalid move direction: %d", data.Direction)
		}

		m.Speed = data.Speed
		m.Enabled = data.Enable
		m.UpdatedAt = time.Now()
		if err := s.driveMotorRepo.UpdateDriveMotor(ctx, s.dbProvider.DB(), m); err != nil {
			return fmt.Errorf("update drive motor: %w", err)
		}

	default:
		return fmt.Errorf("unknown command type: %d", cmd.Type)
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
