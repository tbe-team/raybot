package systemimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/internal/services/system"
)

type service struct {
	log *slog.Logger

	commandService    command.Service
	driveMotorService drivemotor.Service
	liftMotorService  liftmotor.Service
}

func NewService(
	log *slog.Logger,
	commandService command.Service,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
) system.Service {
	return &service{
		log:               log,
		commandService:    commandService,
		driveMotorService: driveMotorService,
		liftMotorService:  liftMotorService,
	}
}

func (s service) Reboot(_ context.Context) error {
	go func() {
		time.Sleep(1 * time.Second)
		cmd := exec.Command("reboot")
		if err := cmd.Run(); err != nil {
			s.log.Error("failed to reboot", slog.Any("error", err))
		}
	}()

	return nil
}

func (s service) StopEmergency(ctx context.Context) error {
	if err := s.commandService.CancelAllRunningCommands(ctx); err != nil {
		return fmt.Errorf("cancel all running commands: %w", err)
	}

	if err := s.driveMotorService.Stop(ctx); err != nil {
		if !errors.Is(err, drivemotor.ErrCanNotControlDriveMotor) {
			return fmt.Errorf("failed to stop drive motor: %w", err)
		}
	}

	if err := s.liftMotorService.Stop(ctx); err != nil {
		if !errors.Is(err, liftmotor.ErrCanNotControlLiftMotor) {
			return fmt.Errorf("failed to stop lift motor: %w", err)
		}
	}

	return nil
}
