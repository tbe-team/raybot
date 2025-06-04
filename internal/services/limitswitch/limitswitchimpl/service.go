package limitswitchimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/limitswitch"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	log       *slog.Logger
	validator validator.Validator
	publisher eventbus.EventBus
	repo      limitswitch.Repository
}

func NewService(
	log *slog.Logger,
	validator validator.Validator,
	publisher eventbus.EventBus,
	repo limitswitch.Repository,
) limitswitch.Service {
	return &Service{
		log:       log.With("service", "limitswitch"),
		validator: validator,
		publisher: publisher,
		repo:      repo,
	}
}

func (s Service) UpdateLimitSwitchState(ctx context.Context, params limitswitch.UpdateLimitSwitchStateParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	cur, err := s.repo.GetLimitSwitchState(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("get limit switch state: %w", err)
	}

	// No change, do nothing
	if cur.Pressed == params.Pressed {
		return nil
	}

	if err := s.repo.UpdateLimitSwitchState(ctx, params.ID, params.Pressed); err != nil {
		return fmt.Errorf("update limit switch state: %w", err)
	}

	if params.Pressed {
		s.publishLimitSwitchPressedEvent(ctx, params.ID)
	}

	return nil
}

func (s Service) publishLimitSwitchPressedEvent(_ context.Context, id limitswitch.LimitSwitchID) {
	switch id {
	case limitswitch.LimitSwitchID1:
		s.publisher.Publish(events.LimitSwitch1PressedTopic, &eventbus.Message{
			Payload: events.LimitSwitch1PressedEvent{},
		})

	default:
		s.log.Error("invalid limit switch id", slog.Any("id", id))
	}
}
