package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

// syncStateType is the type of sync state received from the PIC
type syncStateType uint8

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *syncStateType) UnmarshalText(text []byte) error {
	n, err := strconv.ParseUint(string(text), 10, 8)
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
		return fmt.Errorf("invalid sync state type: %s", string(text))
	}
	return nil
}

const (
	syncStateTypeBattery syncStateType = iota
	syncStateTypeCharge
	syncStateTypeDischarge
	syncStateTypeDistanceSensor
	syncStateTypeLiftMotor
	syncStateTypeDriveMotor
)

type SyncStateMessage struct {
	StateType syncStateType   `json:"state_type"`
	Data      json.RawMessage `json:"data"`
}

type SyncStateHandler struct {
	robotService service.RobotService
	log          *slog.Logger
}

func NewSyncStateHandler(robotService service.RobotService) *SyncStateHandler {
	return &SyncStateHandler{
		robotService: robotService,
		log: slog.With(
			slog.String("module", "pic"),
			slog.String("handler", "SyncStateHandler"),
		),
	}
}

func (h *SyncStateHandler) Handle(ctx context.Context, msg SyncStateMessage) {
	params := service.UpdateRobotStateParams{}
	switch msg.StateType {
	case syncStateTypeBattery:
		params.SetBattery = true
		var temp struct {
			Current      uint   `json:"current"`
			Temp         uint   `json:"temp"`
			Voltage      uint   `json:"voltage"`
			CellVoltages []uint `json:"cell_voltages"`
			Percent      uint8  `json:"percent"`
			Fault        uint8  `json:"fault"`
			Health       uint8  `json:"health"`
			Status       uint8  `json:"status"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal battery data", "error", err)
			return
		}

		params.Battery.Current = temp.Current
		params.Battery.Temp = temp.Temp
		params.Battery.Voltage = temp.Voltage
		params.Battery.CellVoltages = temp.CellVoltages
		params.Battery.Percent = temp.Percent
		params.Battery.Fault = temp.Fault
		params.Battery.Health = temp.Health
		params.Battery.Status = temp.Status

	case syncStateTypeCharge:
		params.SetCharge = true
		var temp struct {
			CurrentLimit uint `json:"current_limit"`
			Enabled      bool `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal charge data", "error", err)
			return
		}

		params.Charge.CurrentLimit = temp.CurrentLimit
		params.Charge.Enabled = temp.Enabled

	case syncStateTypeDischarge:
		params.SetDischarge = true
		var temp struct {
			CurrentLimit uint `json:"current_limit"`
			Enabled      bool `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal discharge data", "error", err)
			return
		}

		params.Discharge.CurrentLimit = temp.CurrentLimit
		params.Discharge.Enabled = temp.Enabled

	case syncStateTypeDistanceSensor:
		params.SetDistanceSensor = true
		var temp struct {
			FrontDistance uint `json:"front_distance"`
			BackDistance  uint `json:"back_distance"`
			DownDistance  uint `json:"down_distance"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal distance sensor data", "error", err)
			return
		}

		params.DistanceSensor.FrontDistance = temp.FrontDistance
		params.DistanceSensor.BackDistance = temp.BackDistance
		params.DistanceSensor.DownDistance = temp.DownDistance

	case syncStateTypeLiftMotor:
		params.SetLiftMotor = true
		var temp struct {
			Direction model.LiftMotorDirection `json:"direction"`
			Speed     uint8                    `json:"speed"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal lift motor data", "error", err)
			return
		}

		params.LiftMotor.Direction = temp.Direction
		params.LiftMotor.Speed = temp.Speed

	case syncStateTypeDriveMotor:
		params.SetDriveMotor = true
		var temp struct {
			Direction model.DriveMotorDirection `json:"direction"`
			Speed     uint8                     `json:"speed"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal drive motor data", "error", err)
			return
		}

		params.DriveMotor.Direction = temp.Direction
		params.DriveMotor.Speed = temp.Speed

	default:
		h.log.Error("unknown state type", "type", msg.StateType)
		return
	}

	_, err := h.robotService.UpdateRobotState(ctx, params)
	if err != nil {
		h.log.Error("failed to update robot state", "error", err)
	}
}
