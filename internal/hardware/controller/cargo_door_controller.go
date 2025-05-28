package controller

import (
	"context"
	"fmt"
)

type CargoDoorController interface {
	OpenCargoDoor(ctx context.Context, speed uint8) error
	CloseCargoDoor(ctx context.Context, speed uint8) error
}

func (c *controller) OpenCargoDoor(ctx context.Context, speed uint8) error {
	id := c.genIDFunc()
	cmd := espCommand{
		ID:   id,
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionOpen,
			Speed:     speed,
			Enable:    true,
		},
	}

	if err := c.createESPCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create ESP command with ACK: %w", err)
	}

	return nil
}

func (c *controller) CloseCargoDoor(ctx context.Context, speed uint8) error {
	id := c.genIDFunc()
	cmd := espCommand{
		ID:   id,
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionClose,
			Speed:     speed,
			Enable:    true,
		},
	}

	if err := c.createESPCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create ESP command with ACK: %w", err)
	}

	return nil
}
