package emergencyimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/emergency"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type Service struct {
	commandService    command.Service
	driveMotorService drivemotor.Service
	liftMotorService  liftmotor.Service

	emergencyStateRepository emergency.Repository
}

func NewService(
	commandService command.Service,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	emergencyStateRepository emergency.Repository,
) emergency.Service {
	return &Service{
		commandService:           commandService,
		driveMotorService:        driveMotorService,
		liftMotorService:         liftMotorService,
		emergencyStateRepository: emergencyStateRepository,
	}
}

func (s Service) GetEmergencyState(ctx context.Context) (emergency.State, error) {
	return s.emergencyStateRepository.GetEmergencyState(ctx)
}

func (s Service) StopOperation(ctx context.Context) error {
	if err := s.commandService.LockProcessingCommand(ctx); err != nil {
		return fmt.Errorf("lock processing command: %w", err)
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

	if err := s.emergencyStateRepository.UpdateEmergencyState(ctx, emergency.State{Locked: true}); err != nil {
		return fmt.Errorf("failed to update emergency state: %w", err)
	}

	return nil
}

func (s Service) ResumeOperation(ctx context.Context) error {
	if err := s.commandService.UnlockProcessingCommand(ctx); err != nil {
		return fmt.Errorf("unlock processing command: %w", err)
	}

	if err := s.emergencyStateRepository.UpdateEmergencyState(ctx, emergency.State{Locked: false}); err != nil {
		return fmt.Errorf("failed to update emergency state: %w", err)
	}

	return nil
}
