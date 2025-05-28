package controller

import "encoding/json"

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
	Speed          uint8
	Enable         bool
}

func (d picCommandLiftMotorData) MarshalJSON() ([]byte, error) {
	var temp struct {
		TargetPosition uint16 `json:"target_position"`
		MaxOutput      uint16 `json:"max_output"`
		Enable         uint8  `json:"enable"`
	}

	temp.TargetPosition = d.TargetPosition
	temp.MaxOutput = uint16(d.Speed)
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
