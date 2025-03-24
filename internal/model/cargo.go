package model

import (
	"fmt"
	"time"
)

type Cargo struct {
	IsOpen         bool
	QRCode         string
	BottomDistance uint16
	UpdatedAt      time.Time
}

type CargoDoorMotor struct {
	Direction CargoDoorMotorDirection
	Speed     uint8
	IsRunning bool
	Enabled   bool
	UpdatedAt time.Time
}

type CargoDoorMotorDirection uint8

func (s CargoDoorMotorDirection) Validate() error {
	switch s {
	case CargoDoorDirectionClose, CargoDoorDirectionOpen:
		return nil
	default:
		return fmt.Errorf("invalid cargo door direction: %d", s)
	}
}

func (s CargoDoorMotorDirection) String() string {
	return []string{"CLOSE", "OPEN"}[s]
}

const (
	CargoDoorDirectionClose CargoDoorMotorDirection = iota
	CargoDoorDirectionOpen
)
