package ledimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/tbe-team/raybot/internal/config"
	ledh "github.com/tbe-team/raybot/internal/hardware/led"
	"github.com/tbe-team/raybot/internal/services/led"
	"github.com/tbe-team/raybot/pkg/ptr"
)

const defaultBlinkInterval = 500 * time.Millisecond

type Service struct {
	cfg        config.Leds
	log        *slog.Logger
	repository led.Repository
	systemLed  *ledh.Led
	alertLed   *ledh.Led
}

func NewService(
	cfg config.Leds,
	log *slog.Logger,
	repository led.Repository,
) *Service {
	return &Service{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}

func (s *Service) Start(ctx context.Context) {
	systemLed, err := ledh.New(s.cfg.System.Pin)
	if err != nil {
		s.log.Error("failed to create system led", slog.Any("error", err))
		if err := s.repository.UpdateSystemLedConnection(ctx, led.Connection{
			Connected: false,
			Error:     ptr.New(err.Error()),
		}); err != nil {
			s.log.Error("failed to update system led connection", slog.Any("error", err))
		}
	} else {
		s.systemLed = systemLed
		if err := s.repository.UpdateSystemLedConnection(ctx, led.Connection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
		}); err != nil {
			s.log.Error("failed to update system led connection", slog.Any("error", err))
		}
	}

	alertLed, err := ledh.New(s.cfg.Alert.Pin)
	if err != nil {
		s.log.Error("failed to create alert led", slog.Any("error", err))
		if err := s.repository.UpdateAlertLedConnection(ctx, led.Connection{
			Connected: false,
			Error:     ptr.New(err.Error()),
		}); err != nil {
			s.log.Error("failed to update alert led connection", slog.Any("error", err))
		}
	} else {
		s.alertLed = alertLed
		if err := s.repository.UpdateAlertLedConnection(ctx, led.Connection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
		}); err != nil {
			s.log.Error("failed to update alert led connection", slog.Any("error", err))
		}
	}

	if err := s.SetSystemLedOn(ctx); err != nil {
		s.log.Error("failed to set system led on", slog.Any("error", err))
	}

	if err = s.SetAlertLedOff(ctx); err != nil {
		s.log.Error("failed to set alert led on", slog.Any("error", err))
	}
}

func (s *Service) Stop() error {
	var errs []error
	if s.systemLed != nil {
		if err := s.systemLed.Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.alertLed != nil {
		if err := s.alertLed.Stop(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to stop leds: %w", errors.Join(errs...))
	}
	return nil
}

func (s *Service) SetSystemLedOn(ctx context.Context) error {
	if s.systemLed == nil {
		return led.ErrLedNotConnected
	}

	if err := s.systemLed.On(); err != nil {
		return fmt.Errorf("failed to set system led on: %w", err)
	}
	return s.repository.UpdateSystemLedState(ctx, led.State{
		Mode:      led.ModeOn,
		UpdatedAt: time.Now(),
	})
}

func (s *Service) SetSystemLedOff(ctx context.Context) error {
	if s.systemLed == nil {
		return led.ErrLedNotConnected
	}

	if err := s.systemLed.Off(); err != nil {
		return fmt.Errorf("failed to set system led off: %w", err)
	}
	return s.repository.UpdateSystemLedState(ctx, led.State{
		Mode:      led.ModeOff,
		UpdatedAt: time.Now(),
	})
}

func (s *Service) BlinkSystemLed(ctx context.Context, params led.BlinkSystemLedParams) error {
	if s.systemLed == nil {
		return led.ErrLedNotConnected
	}

	defer func() {
		if err := s.repository.UpdateSystemLedState(ctx, led.State{
			Mode:      led.ModeOff,
			UpdatedAt: time.Now(),
		}); err != nil {
			s.log.Error("failed to update system led state", slog.Any("error", err))
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		blinkCtx, cancel := context.WithTimeout(ctx, params.Duration)
		defer cancel()

		if err := s.systemLed.Blink(blinkCtx, defaultBlinkInterval); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return fmt.Errorf("failed to blink system led: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		return s.repository.UpdateSystemLedState(ctx, led.State{
			Mode:      led.ModeBlink,
			UpdatedAt: time.Now(),
		})
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to blink system led: %w", err)
	}

	return nil
}

func (s *Service) SetAlertLedOn(ctx context.Context) error {
	if s.alertLed == nil {
		return led.ErrLedNotConnected
	}

	if err := s.alertLed.On(); err != nil {
		return fmt.Errorf("failed to set alert led on: %w", err)
	}
	return s.repository.UpdateAlertLedState(ctx, led.State{
		Mode:      led.ModeOn,
		UpdatedAt: time.Now(),
	})
}

func (s *Service) SetAlertLedOff(ctx context.Context) error {
	if s.alertLed == nil {
		return led.ErrLedNotConnected
	}

	if err := s.alertLed.Off(); err != nil {
		return fmt.Errorf("failed to set alert led off: %w", err)
	}
	return s.repository.UpdateAlertLedState(ctx, led.State{
		Mode:      led.ModeOff,
		UpdatedAt: time.Now(),
	})
}

func (s *Service) BlinkAlertLed(ctx context.Context, params led.BlinkAlertLedParams) error {
	if s.alertLed == nil {
		return led.ErrLedNotConnected
	}

	defer func() {
		if err := s.repository.UpdateAlertLedState(ctx, led.State{
			Mode:      led.ModeOff,
			UpdatedAt: time.Now(),
		}); err != nil {
			s.log.Error("failed to update alert led state", slog.Any("error", err))
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		blinkCtx, cancel := context.WithTimeout(ctx, params.Duration)
		defer cancel()

		if err := s.alertLed.Blink(blinkCtx, defaultBlinkInterval); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return fmt.Errorf("failed to blink alert led: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		return s.repository.UpdateAlertLedState(ctx, led.State{
			Mode:      led.ModeBlink,
			UpdatedAt: time.Now(),
		})
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to blink alert led: %w", err)
	}

	return nil
}
