package controller

import (
	"context"
	"fmt"
)

type BatteryController interface {
	ConfigBatteryCharge(ctx context.Context, currentLimit uint16, enable bool) error
	ConfigBatteryDischarge(ctx context.Context, currentLimit uint16, enable bool) error
}

func (c *controller) ConfigBatteryCharge(ctx context.Context, currentLimit uint16, enable bool) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeBatteryCharge,
		Data: picCommandBatteryChargeData{
			CurrentLimit: currentLimit,
			Enable:       enable,
		},
	}

	if err := c.createPICCommand(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}
	return nil
}

func (c *controller) ConfigBatteryDischarge(ctx context.Context, currentLimit uint16, enable bool) error {
	id := c.genIDFunc()
	cmd := picCommand{
		ID:   id,
		Type: picCommandTypeBatteryDischarge,
		Data: picCommandBatteryDischargeData{
			CurrentLimit: currentLimit,
			Enable:       enable,
		},
	}

	if err := c.createPICCommand(ctx, cmd); err != nil {
		return fmt.Errorf("create PIC command with ACK: %w", err)
	}

	return nil
}
