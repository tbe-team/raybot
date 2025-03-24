package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

type SyncStateHandler struct {
	cargoService service.CargoControlService
	log          *slog.Logger
}

func NewSyncStateHandler(cargoService service.CargoControlService, log *slog.Logger) *SyncStateHandler {
	return &SyncStateHandler{
		cargoService: cargoService,
		log:          log,
	}
}

func (h SyncStateHandler) Handle(ctx context.Context, msg syncStateMessage) {
	switch msg.StateType {
	case syncStateTypeDoor:
		var temp struct {
			IsOpen bool `json:"is_open"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal door state", slog.Any("error", err), slog.Any("data", msg.Data))
			return
		}

		if err := h.cargoService.SyncCargoDoorState(ctx, service.SyncCargoDoorStateParams{
			IsOpen: temp.IsOpen,
		}); err != nil {
			h.log.Error("failed to sync door state", slog.Any("error", err))
		}

	case syncStateTypeMotor:
		var temp struct {
			Direction uint8 `json:"state"` // 0: close, 1: open
			Speed     uint8 `json:"speed"` // 0-100
			IsRunning bool  `json:"is_running"`
			Enabled   bool  `json:"enabled"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal motor state", slog.Any("error", err), slog.Any("data", msg.Data))
			return
		}

		if err := h.cargoService.SyncCargoDoorMotorState(ctx, service.SyncCargoDoorMotorStateParams{
			Direction: model.CargoDoorMotorDirection(temp.Direction),
			Speed:     temp.Speed,
			IsRunning: temp.IsRunning,
			Enabled:   temp.Enabled,
		}); err != nil {
			h.log.Error("failed to sync motor state", slog.Any("error", err))
		}

	case syncStateTypeQRScanner:
		var temp struct {
			Code string `json:"code"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal qr code", slog.Any("error", err), slog.Any("data", msg.Data))
			return
		}

		if err := h.cargoService.SyncCargoQRCode(ctx, service.SyncCargoQRCodeParams{
			QRCode: temp.Code,
		}); err != nil {
			h.log.Error("failed to sync qr code", slog.Any("error", err))
		}

	case syncStateTypeBottomDistanceSensor:
		var temp struct {
			Distance uint16 `json:"under_distance"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			h.log.Error("failed to unmarshal bottom distance", slog.Any("error", err), slog.Any("data", msg.Data))
			return
		}

		if err := h.cargoService.SyncCargoBottomDistance(ctx, service.SyncCargoBottomDistanceParams{
			BottomDistance: temp.Distance,
		}); err != nil {
			h.log.Error("failed to sync bottom distance", slog.Any("error", err))
		}

	default:
		h.log.Error("invalid sync state type", slog.Any("type", msg.StateType))
	}
}

type syncStateMessage struct {
	StateType syncStateType   `json:"state_type"`
	Data      json.RawMessage `json:"data"`
}

type syncStateType uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *syncStateType) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}

	switch n {
	case 0:
		*s = syncStateTypeDoor
	case 1:
		*s = syncStateTypeMotor
	case 2:
		*s = syncStateTypeQRScanner
	case 3:
		*s = syncStateTypeBottomDistanceSensor
	default:
		return fmt.Errorf("invalid sync state type: %d", n)
	}

	return nil
}

const (
	syncStateTypeDoor syncStateType = iota
	syncStateTypeMotor
	syncStateTypeQRScanner
	syncStateTypeBottomDistanceSensor
)
