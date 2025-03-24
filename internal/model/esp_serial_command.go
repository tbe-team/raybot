package model

import (
	"fmt"
	"strconv"
	"time"
)

type ESPSerialCommand struct {
	ID        string
	Type      ESPSerialCommandType
	Data      ESPSerialCommandData
	CreatedAt time.Time
}

type ESPSerialCommandData interface {
	isESPSerialCommandData()
}

type ESPSerialCommandCargoDoorMotorData struct {
	Direction CargoDoorMotorDirection `json:"direction"`
	Speed     uint8                   `json:"speed"`
	Enable    bool                    `json:"enable"`
}

func (ESPSerialCommandCargoDoorMotorData) isESPSerialCommandData() {}

type ESPSerialCommandType uint8

func (s *ESPSerialCommandType) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}
	switch n {
	case 0:
		*s = ESPSerialCommandTypeCargoDoorMotor
	default:
		return fmt.Errorf("invalid ESP serial command type: %s", string(data))
	}

	return nil
}

func (s ESPSerialCommandType) String() string {
	return []string{"CARGO_DOOR_MOTOR"}[s]
}

const (
	ESPSerialCommandTypeCargoDoorMotor ESPSerialCommandType = iota
)
