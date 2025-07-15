package alarmimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	log       *slog.Logger
	validator validator.Validator
	alarmRepo alarm.Repository
}

func NewService(
	log *slog.Logger,
	validator validator.Validator,
	alarmRepo alarm.Repository,
) *Service {
	s := &Service{
		log:       log,
		validator: validator,
		alarmRepo: alarmRepo,
	}

	go s.deactivateAllAlarms(context.TODO())

	return s
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

func (s Service) deactivateAllAlarms(ctx context.Context) {
	if err := s.alarmRepo.DeactivateAllAlarms(ctx); err != nil {
		s.log.Error("failed to deactivate all alarms", "error", err)
	}
}
