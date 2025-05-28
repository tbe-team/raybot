package controller

import "encoding/json"

type espCommand struct {
	ID   string         `json:"id"`
	Type espCommandType `json:"type"`
	Data espData        `json:"data"`
}

type espCommandType uint8

func (t espCommandType) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(t))
}

const (
	espCommandTypeCargoDoorMotor espCommandType = 0
)

type espData interface {
	isEspData()
}

type espCargoDoorMotorData struct {
	Direction doorDirection
	Speed     uint8
	Enable    bool
}

func (espCargoDoorMotorData) isEspData() {}

func (d espCargoDoorMotorData) MarshalJSON() ([]byte, error) {
	var data struct {
		State  uint8 `json:"state"`
		Speed  uint8 `json:"speed"`
		Enable uint8 `json:"enable"`
	}

	data.State = uint8(d.Direction)
	data.Speed = d.Speed
	data.Enable = boolToUint8(d.Enable)
	return json.Marshal(data)
}

type doorDirection uint8

const (
	doorDirectionClose doorDirection = 0
	doorDirectionOpen  doorDirection = 1
)
