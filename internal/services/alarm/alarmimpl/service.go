package alarmimpl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	log       *slog.Logger
	validator validator.Validator
	alarmRepo alarm.Repository

	stopCh     chan struct{}
	stopDoneCh chan struct{}
}

func NewService(
	log *slog.Logger,
	validator validator.Validator,
	alarmRepo alarm.Repository,
) *Service {
	return &Service{
		log:        log,
		validator:  validator,
		alarmRepo:  alarmRepo,
		stopCh:     make(chan struct{}),
		stopDoneCh: make(chan struct{}),
	}
}

func (s Service) Start(ctx context.Context) {
	go func() {
		if err := s.alarmRepo.DeactivateAllAlarms(ctx); err != nil {
			s.log.Error("failed to deactivate all alarms", "error", err)
		}
	}()

	go s.startDeactivatedAlarmsCleanupJob(ctx)
}

func (s Service) Stop() {
	close(s.stopCh)
	<-s.stopDoneCh
}

const cleanupAlarmInterval = 24 * 7 * time.Hour

func (s Service) startDeactivatedAlarmsCleanupJob(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-s.stopCh:
				s.stopDoneCh <- struct{}{}
				return

			case <-time.After(cleanupAlarmInterval):
				if err := s.alarmRepo.DeactivateAllAlarms(ctx); err != nil {
					s.log.Error("failed to deactivate all alarms", slog.Any("error", err))
				}
			}
		}
	}()
}

func (s Service) ListActiveAlarms(ctx context.Context, params alarm.ListActiveAlarmsParams) (paging.List[alarm.Alarm], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[alarm.Alarm]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.alarmRepo.ListActiveAlarms(ctx, params)
}

func (s Service) ListDeactiveAlarms(ctx context.Context, params alarm.ListDeactiveAlarmsParams) (paging.List[alarm.Alarm], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[alarm.Alarm]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.alarmRepo.ListDeactiveAlarms(ctx, params)
}

func (s Service) DeleteDeactiveAlarms(ctx context.Context) error {
	return s.alarmRepo.DeleteDeactiveAlarms(ctx)
}
