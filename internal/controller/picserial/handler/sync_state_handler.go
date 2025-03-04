package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

// SyncStateType is the type of sync state received from the PIC
type SyncStateType int8

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *SyncStateType) UnmarshalText(text []byte) error {
	switch int(text[0]) {
	case 0:
		*s = SyncStateTypeBattery
	case 1:
		*s = SyncStateTypeCharge
	case 2:
		*s = SyncStateTypeDischarge
	case 3:
		*s = SyncStateTypeDistanceSensor
	case 4:
		*s = SyncStateTypeLiftMotor
	case 5:
		*s = SyncStateTypeDriveMotor
	default:
		return fmt.Errorf("invalid sync state type: %s", string(text))
	}
	return nil
}

const (
	SyncStateTypeBattery SyncStateType = iota
	SyncStateTypeCharge
	SyncStateTypeDischarge
	SyncStateTypeDistanceSensor
	SyncStateTypeLiftMotor
	SyncStateTypeDriveMotor
)

type SyncStateMessage struct {
	StateType SyncStateType   `json:"state_type"`
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
	case SyncStateTypeBattery:
		params.SetBattery = true
		var temp struct {
			Current      int   `json:"current"`
			Temp         int   `json:"temp"`
			Voltage      int   `json:"voltage"`
			CellVoltages []int `json:"cell_voltages"`
			Percent      uint8 `json:"percent"`
			Fault        uint8 `json:"fault"`
			Health       uint8 `json:"health"`
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
	case SyncStateTypeCharge:
		params.SetCharge = true
		var temp struct {
			CurrentLimit int  `json:"current_limit"`
			Enabled      bool `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal charge data", "error", err)
			return
		}

		params.Charge.CurrentLimit = temp.CurrentLimit
		params.Charge.Enabled = temp.Enabled
	case SyncStateTypeDischarge:
		params.SetDischarge = true
		var temp struct {
			CurrentLimit int  `json:"current_limit"`
			Enabled      bool `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal discharge data", "error", err)
			return
		}

		params.Discharge.CurrentLimit = temp.CurrentLimit
		params.Discharge.Enabled = temp.Enabled
	case SyncStateTypeDistanceSensor:
		params.SetDistanceSensor = true
		var temp struct {
			FrontDistance int `json:"front_distance"`
			BackDistance  int `json:"back_distance"`
			DownDistance  int `json:"down_distance"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal distance sensor data", "error", err)
			return
		}

		params.DistanceSensor.FrontDistance = temp.FrontDistance
		params.DistanceSensor.BackDistance = temp.BackDistance
		params.DistanceSensor.DownDistance = temp.DownDistance
	case SyncStateTypeLiftMotor:
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
	case SyncStateTypeDriveMotor:
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
