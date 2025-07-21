package alarmimpl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/internal/services/led"
	"github.com/tbe-team/raybot/internal/services/system"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	log        *slog.Logger
	validator  validator.Validator
	alarmRepo  alarm.Repository
	systemRepo system.Repository
	ledService led.Service
}

func NewService(
	log *slog.Logger,
	validator validator.Validator,
	alarmRepo alarm.Repository,
	systemRepo system.Repository,
	ledService led.Service,
) *Service {
	s := &Service{
		log:        log,
		validator:  validator,
		alarmRepo:  alarmRepo,
		systemRepo: systemRepo,
		ledService: ledService,
	}

	go s.deactivateAllActivatedAlarms(context.TODO())

	return s
}

func (s Service) ListActiveAlarms(ctx context.Context, params alarm.ListActiveAlarmsParams) (paging.List[alarm.Alarm], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[alarm.Alarm]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.alarmRepo.ListActiveAlarms(ctx, params.PagingParams)
}

func (s Service) ListDeactiveAlarms(ctx context.Context, params alarm.ListDeactiveAlarmsParams) (paging.List[alarm.Alarm], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[alarm.Alarm]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.alarmRepo.ListDeactiveAlarms(ctx, params.PagingParams)
}

func (s Service) DeleteDeactivatedAlarms(ctx context.Context) error {
	return s.alarmRepo.DeleteDeactivatedAlarms(ctx)
}

func (s Service) DeleteDeactivatedAlarmsByThreshold(ctx context.Context, params alarm.DeleteDeactivatedAlarmsByThresholdParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.alarmRepo.DeleteDeactivatedAlarmsByThreshold(ctx, params.Threshold)
}

func (s Service) deactivateAllActivatedAlarms(ctx context.Context) {
	if err := s.alarmRepo.DeactivateAllAlarms(ctx); err != nil {
		s.log.Error("failed to deactivate all activated alarms", slog.Any("error", err))
	}
}

func (s Service) DeactivateAlarm(ctx context.Context, params alarm.DeactivateAlarmParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.alarmRepo.DeactivateAlarm(ctx, params.ID, time.Now()); err != nil {
		return fmt.Errorf("failed to deactivate alarm: %w", err)
	}

	countActivatedAlarms, err := s.alarmRepo.CountActivatedAlarms(ctx)
	if err != nil {
		return fmt.Errorf("failed to count activated alarms: %w", err)
	}

	if countActivatedAlarms == 0 {
		if err := s.systemRepo.UpdateStatus(ctx, system.StatusNormal); err != nil {
			return fmt.Errorf("failed to update system status: %w", err)
		}

		if err := s.ledService.SetAlertLedOff(ctx); err != nil {
			return fmt.Errorf("failed to set alert led off: %w", err)
		}
	}

	return nil
}
