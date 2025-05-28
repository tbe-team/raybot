package controller

import (
	"context"
	"fmt"
)

type DriveMotorController interface {
	MoveForward(ctx context.Context, speed uint8) error
	MoveBackward(ctx context.Context, speed uint8) error
	// When moving backward, the robot jerks on stop because the STOP command sets direction to FORWARD.
	// Fix it respecting the current direction in STOP.
	StopDriveMotor(ctx context.Context, moveForward bool) error
}

func (c *controller) MoveForward(ctx context.Context, speed uint8) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: moveDirectionForward,
			Speed:     speed,
			Enable:    true,
		},
	}

	if err := c.createPICCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}

	return nil
}

func (c *controller) MoveBackward(ctx context.Context, speed uint8) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: moveDirectionBackward,
			Speed:     speed,
			Enable:    true,
		},
	}

	if err := c.createPICCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}
	return nil
}

func (c *controller) StopDriveMotor(ctx context.Context, moveForward bool) error {
	var direction moveDirection
	if moveForward {
		direction = moveDirectionForward
	} else {
		direction = moveDirectionBackward
	}

	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: direction,
			Speed:     0,
			Enable:    false,
		},
	}

	if err := c.createPICCommandWithACK(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}

	return nil
}
