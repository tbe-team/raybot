package system

import "context"

type Service interface {
	Reboot(ctx context.Context) error

	// StopEmergency stops all motors and cancel all queued and processing commands.
	StopEmergency(ctx context.Context) error
}
