package picserial

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

func (s *Service) HandleSyncState(ctx context.Context, msg syncStateMessage) error {
	switch msg.StateType {
	case syncStateTypeBattery:
		var temp struct {
			Current      uint16   `json:"current"`
			Temp         uint8    `json:"temp"`
			Voltage      uint16   `json:"voltage"`
			CellVoltages []uint16 `json:"cell_voltages"`
			Percent      uint8    `json:"percent"`
			Fault        uint8    `json:"fault"`
			Health       uint8    `json:"health"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal battery data: %w", err)
		}

		if err := s.batteryService.UpdateBatteryState(ctx, battery.UpdateBatteryStateParams{
			Current:      temp.Current,
			Temp:         temp.Temp,
			Voltage:      temp.Voltage,
			CellVoltages: temp.CellVoltages,
			Percent:      temp.Percent,
			Fault:        temp.Fault,
			Health:       temp.Health,
		}); err != nil {
			return fmt.Errorf("failed to update battery state: %w", err)
		}

	case syncStateTypeCharge:
		var temp struct {
			CurrentLimit uint16 `json:"current_limit"`
			Enabled      uint8  `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal charge data: %w", err)
		}

		if err := s.batteryService.UpdateChargeSetting(ctx, battery.UpdateChargeSettingParams{
			CurrentLimit: temp.CurrentLimit,
			Enabled:      temp.Enabled == 1,
		}); err != nil {
			return fmt.Errorf("failed to update charge setting: %w", err)
		}

	case syncStateTypeDischarge:
		var temp struct {
			CurrentLimit uint16 `json:"current_limit"`
			Enabled      uint8  `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal discharge data: %w", err)
		}

		if err := s.batteryService.UpdateDischargeSetting(ctx, battery.UpdateDischargeSettingParams{
			CurrentLimit: temp.CurrentLimit,
			Enabled:      temp.Enabled == 1,
		}); err != nil {
			return fmt.Errorf("failed to update discharge setting: %w", err)
		}

	case syncStateTypeDistanceSensor:
		var temp struct {
			FrontDistance uint16 `json:"front"`
			BackDistance  uint16 `json:"back"`
			DownDistance  uint16 `json:"down"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal distance sensor data: %w", err)
		}

		if err := s.distanceSensorService.UpdateDistanceSensorState(ctx, distancesensor.UpdateDistanceSensorStateParams{
			FrontDistance: temp.FrontDistance,
			BackDistance:  temp.BackDistance,
			DownDistance:  temp.DownDistance,
		}); err != nil {
			return fmt.Errorf("failed to update distance sensor state: %w", err)
		}

	case syncStateTypeLiftMotor:
		var temp struct {
			CurrentPosition uint16 `json:"current_position"`
			TargetPosition  uint16 `json:"target_position"`
			IsRunning       uint8  `json:"is_running"`
			Enabled         uint8  `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal lift motor data: %w", err)
		}

		if err := s.liftMotorService.UpdateLiftMotorState(ctx, liftmotor.UpdateLiftMotorStateParams{
			CurrentPosition: temp.CurrentPosition,
			TargetPosition:  temp.TargetPosition,
			IsRunning:       temp.IsRunning == 1,
			Enabled:         temp.Enabled == 1,
		}); err != nil {
			return fmt.Errorf("failed to update lift motor state: %w", err)
		}

	case syncStateTypeDriveMotor:
		var temp struct {
			Direction uint8 `json:"direction"`
			Speed     uint8 `json:"speed"`
			IsRunning uint8 `json:"is_running"`
			Enabled   uint8 `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal drive motor data: %w", err)
		}

		var direction drivemotor.Direction
		switch temp.Direction {
		case 0:
			direction = drivemotor.DirectionForward
		case 1:
			direction = drivemotor.DirectionBackward
		default:
			return fmt.Errorf("invalid drive motor direction: %d", temp.Direction)
		}

		if err := s.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
			Direction: direction,
			Speed:     temp.Speed,
			IsRunning: temp.IsRunning == 1,
			Enabled:   temp.Enabled == 1,
		}); err != nil {
			return fmt.Errorf("failed to update drive motor state: %w", err)
		}

	default:
		return fmt.Errorf("invalid sync state type: %s", string(msg.Data))
	}

	return nil
}

// syncStateType is the type of sync state received from the PIC
type syncStateType uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *syncStateType) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}

	switch n {
	case 0:
		*s = syncStateTypeBattery
	case 1:
		*s = syncStateTypeCharge
	case 2:
		*s = syncStateTypeDischarge
	case 3:
		*s = syncStateTypeDistanceSensor
	case 4:
		*s = syncStateTypeLiftMotor
	case 5:
		*s = syncStateTypeDriveMotor
	default:
		return fmt.Errorf("invalid sync state type: %s", string(data))
	}
	return nil
}

const (
	syncStateTypeBattery        syncStateType = 0
	syncStateTypeCharge         syncStateType = 1
	syncStateTypeDischarge      syncStateType = 2
	syncStateTypeDistanceSensor syncStateType = 3
	syncStateTypeLiftMotor      syncStateType = 4
	syncStateTypeDriveMotor     syncStateType = 5
)

type syncStateMessage struct {
	StateType syncStateType   `json:"state_type"`
	Data      json.RawMessage `json:"data"`
}
