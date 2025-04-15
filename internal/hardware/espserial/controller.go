package espserial

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lithammer/shortuuid/v4"
)

type Controller interface {
	OpenCargoDoor(ctx context.Context, speed uint8) error
	CloseCargoDoor(ctx context.Context, speed uint8) error
}

var _ Controller = (*DefaultClient)(nil)

func (c *DefaultClient) OpenCargoDoor(ctx context.Context, speed uint8) error {
	cmd := espCommand{
		ID:   shortuuid.New(),
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionOpen,
			Speed:     speed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal cargo door motor command: %w", err)
	}

	if err := c.Write(ctx, cmdJSON); err != nil {
		return fmt.Errorf("write cargo door motor command: %w", err)
	}

	return nil
}

func (c *DefaultClient) CloseCargoDoor(ctx context.Context, speed uint8) error {
	cmd := espCommand{
		ID:   shortuuid.New(),
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionClose,
			Speed:     speed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal cargo door motor command: %w", err)
	}

	if err := c.Write(ctx, cmdJSON); err != nil {
		return fmt.Errorf("write cargo door motor command: %w", err)
	}

	return nil
}

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

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
