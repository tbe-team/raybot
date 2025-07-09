package systemimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/led"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/internal/services/system"
)

type service struct {
	log *slog.Logger

	commandService    command.Service
	driveMotorService drivemotor.Service
	liftMotorService  liftmotor.Service
	ledService        led.Service

	systemRepo system.Repository
}

func NewService(
	log *slog.Logger,
	commandService command.Service,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	ledService led.Service,
	systemInfoRepo system.Repository,
) system.Service {
	return &service{
		log:               log,
		commandService:    commandService,
		driveMotorService: driveMotorService,
		liftMotorService:  liftMotorService,
		ledService:        ledService,
		systemRepo:        systemInfoRepo,
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

func (s service) GetInfo(ctx context.Context) (system.Info, error) {
	return s.systemRepo.GetInfo(ctx)
}

func (s service) GetStatus(ctx context.Context) (system.Status, error) {
	return s.systemRepo.GetStatus(ctx)
}

func (s service) SetStatusError(ctx context.Context) error {
	currentStatus, err := s.systemRepo.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("get current status: %w", err)
	}

	if currentStatus == system.StatusError {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.ledService.SetAlertLedOn(ctx); err != nil {
			if errors.Is(err, led.ErrLedNotConnected) {
				s.log.Warn("alert led is not connected, skipping")
				return nil
			}
			return fmt.Errorf("set alert led on: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		return s.systemRepo.UpdateStatus(ctx, system.StatusError)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("set alert led: %w", err)
	}

	return nil
}
