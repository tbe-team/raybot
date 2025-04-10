package picserial

import (
	"encoding/json"
	"fmt"

	"github.com/lithammer/shortuuid/v4"
)

type Controller interface {
	SetCargoPosition(targetPosition uint16) error
	MoveForward(speed uint8) error
	MoveBackward(speed uint8) error
	StopDriveMotor() error
	ConfigBatteryCharge(currentLimit uint16, enable bool) error
	ConfigBatteryDischarge(currentLimit uint16, enable bool) error
}

var _ Controller = (*DefaultClient)(nil)

func (c *DefaultClient) SetCargoPosition(targetPosition uint16) error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeLiftMotor,
		Data: picCommandLiftMotorData{
			TargetPosition: targetPosition,
			Enable:         true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal lift cargo command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write lift cargo command: %w", err)
	}

	return nil
}

func (c *DefaultClient) MoveForward(speed uint8) error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: moveDirectionForward,
			Speed:     speed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal move forward command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write move forward command: %w", err)
	}

	return nil
}

func (c *DefaultClient) MoveBackward(speed uint8) error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: moveDirectionBackward,
			Speed:     speed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal move backward command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write move backward command: %w", err)
	}

	return nil
}

func (c *DefaultClient) StopDriveMotor() error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: moveDirectionForward,
			Speed:     0,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal stop drive motor command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write stop drive motor command: %w", err)
	}

	return nil
}

func (c *DefaultClient) ConfigBatteryCharge(currentLimit uint16, enable bool) error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeBatteryCharge,
		Data: picCommandBatteryChargeData{
			CurrentLimit: currentLimit,
			Enable:       enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal config battery charge command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write config battery charge command: %w", err)
	}

	return nil
}

func (c *DefaultClient) ConfigBatteryDischarge(currentLimit uint16, enable bool) error {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeBatteryDischarge,
		Data: picCommandBatteryDischargeData{
			CurrentLimit: currentLimit,
			Enable:       enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal config battery discharge command: %w", err)
	}

	if err := c.Write(cmdJSON); err != nil {
		return fmt.Errorf("write config battery discharge command: %w", err)
	}

	return nil
}

type picCommand struct {
	ID   string         `json:"id"`
	Type picCommandType `json:"type"`
	Data picCommandData `json:"data"`
}

type picCommandType uint8

func (t picCommandType) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(t))
}

const (
	picCommandTypeBatteryCharge    picCommandType = 0
	picCommandTypeBatteryDischarge picCommandType = 1
	picCommandTypeLiftMotor        picCommandType = 2
	picCommandTypeDriveMotor       picCommandType = 3
)

type picCommandData interface {
	isPICCommandData()
}

type picCommandBatteryChargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (d picCommandBatteryChargeData) MarshalJSON() ([]byte, error) {
	var temp struct {
		CurrentLimit uint16 `json:"current_limit"`
		Enable       uint8  `json:"enable"`
	}

	temp.CurrentLimit = d.CurrentLimit
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandBatteryChargeData) isPICCommandData() {}

type picCommandBatteryDischargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (d picCommandBatteryDischargeData) MarshalJSON() ([]byte, error) {
	var temp struct {
		CurrentLimit uint16 `json:"current_limit"`
		Enable       uint8  `json:"enable"`
	}

	temp.CurrentLimit = d.CurrentLimit
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandBatteryDischargeData) isPICCommandData() {}

type picCommandLiftMotorData struct {
	TargetPosition uint16
	Enable         bool
}

func (d picCommandLiftMotorData) MarshalJSON() ([]byte, error) {
	var temp struct {
		TargetPosition uint16 `json:"target_position"`
		Enable         uint8  `json:"enable"`
	}

	temp.TargetPosition = d.TargetPosition
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandLiftMotorData) isPICCommandData() {}

type picCommandDriveMotorData struct {
	Direction moveDirection
	Speed     uint8
	Enable    bool
}

func (d picCommandDriveMotorData) MarshalJSON() ([]byte, error) {
	var temp struct {
		Direction uint8 `json:"direction"`
		Speed     uint8 `json:"speed"`
		Enable    uint8 `json:"enable"`
	}

	temp.Direction = uint8(d.Direction)
	temp.Speed = d.Speed
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandDriveMotorData) isPICCommandData() {}

type moveDirection uint8

const (
	moveDirectionForward  moveDirection = 0
	moveDirectionBackward moveDirection = 1
)

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
