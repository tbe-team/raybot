package controller

import (
	"context"
	"fmt"
)

type LiftMotorController interface {
	SetCargoPosition(ctx context.Context, motorSpeed uint8, targetPosition uint16) error
	StopLiftCargoMotor(ctx context.Context) error
}

func (c *controller) SetCargoPosition(ctx context.Context, motorSpeed uint8, targetPosition uint16) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeLiftMotor,
		Data: picCommandLiftMotorData{
			TargetPosition: targetPosition,
			Speed:          motorSpeed,
			Enable:         true,
		},
	}

	if err := c.createPICCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}

	return nil
}

func (c *controller) StopLiftCargoMotor(ctx context.Context) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeLiftMotor,
		Data: picCommandLiftMotorData{
			Enable: false,
		},
	}

	if err := c.createPICCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}

	return nil
}
