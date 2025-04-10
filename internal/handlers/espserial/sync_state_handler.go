package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/services/cargo"
)

func (s *Service) HandleSyncState(ctx context.Context, msg syncStateMessage) error {
	switch msg.StateType {
	case syncStateTypeDoor:
		var temp struct {
			IsOpen bool `json:"is_open"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal door state: %w", err)
		}

		if err := s.cargoService.UpdateCargoDoor(ctx, cargo.UpdateCargoDoorParams{
			IsOpen: temp.IsOpen,
		}); err != nil {
			return fmt.Errorf("failed to sync door state: %w", err)
		}

	case syncStateTypeMotor:
		var temp struct {
			State     uint8 `json:"state"` // 0: close, 1: open
			Enabled   uint8 `json:"enabled"`
			Speed     uint8 `json:"speed"` // 0-100
			IsRunning uint8 `json:"is_running"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal motor state: %w", err)
		}

		var direction cargo.DoorDirection
		switch temp.State {
		case 0:
			direction = cargo.DirectionClose
		case 1:
			direction = cargo.DirectionOpen
		default:
			return fmt.Errorf("invalid motor state: %d", temp.State)
		}

		if err := s.cargoService.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
			Direction:    direction,
			SetDirection: true,
			Speed:        temp.Speed,
			SetSpeed:     true,
			IsRunning:    temp.IsRunning == 1,
			SetIsRunning: true,
			Enabled:      temp.Enabled == 1,
			SetEnabled:   true,
		}); err != nil {
			return fmt.Errorf("failed to sync motor state: %w", err)
		}

	case syncStateTypeQRScanner:
		var temp struct {
			Code string `json:"code"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal qr code: %w", err)
		}

		if err := s.cargoService.UpdateCargoQRCode(ctx, cargo.UpdateCargoQRCodeParams{
			QRCode: temp.Code,
		}); err != nil {
			return fmt.Errorf("failed to sync qr code: %w", err)
		}

	case syncStateTypeBottomDistanceSensor:
		var temp struct {
			Distance uint16 `json:"under"`
		}
		if err := json.Unmarshal(msg.Data, &temp); err != nil {
			return fmt.Errorf("failed to unmarshal bottom distance: %w", err)
		}

		if err := s.cargoService.UpdateCargoBottomDistance(ctx, cargo.UpdateCargoBottomDistanceParams{
			BottomDistance: temp.Distance,
		}); err != nil {
			return fmt.Errorf("failed to sync bottom distance: %w", err)
		}

	default:
		return fmt.Errorf("invalid sync state type: %d", msg.StateType)
	}

	return nil
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
