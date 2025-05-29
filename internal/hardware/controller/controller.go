package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

var (
	ErrCommandACKTimeout = errors.New("command ACK timeout")
)

type Controller interface {
	LiftMotorController
	DriveMotorController
	BatteryController
	CargoDoorController
}

type controller struct {
	cfg             config.Hardware
	log             *slog.Logger
	subscriber      eventbus.Subscriber
	picSerialClient picserial.Client
	espSerialClient espserial.Client

	genIDFunc func() string
}

func New(
	cfg config.Hardware,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	picSerialClient picserial.Client,
	espSerialClient espserial.Client,
	opts ...OptionFunc,
) Controller {
	c := &controller{
		cfg:             cfg,
		log:             log,
		subscriber:      subscriber,
		picSerialClient: picSerialClient,
		espSerialClient: espSerialClient,
		genIDFunc:       newShortID,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *controller) createPICCommandWithACK(ctx context.Context, cmd picCommand) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // ensure cleanup other goroutines

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.trackingPICCommandACK(ctx, cmd.ID)
	}()

	if err := c.picSerialClient.Write(ctx, cmdJSON); err != nil {
		return fmt.Errorf("write command: %w", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("tracking PIC command ack: %w", err)
		}
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *controller) trackingPICCommandACK(ctx context.Context, id string) error {
	log := c.log.With(slog.String("id", id))
	log.Debug("start tracking PIC command ack")

	doneCh := make(chan struct{})
	c.subscriber.Subscribe(ctx, events.PICCmdAckTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.PICCmdAckEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.ID == id {
			if ev.Success {
				log.Debug("PIC command ack received")
			} else {
				log.Error("PIC command ack failed")
			}
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
		log.Debug("stop tracking PIC command ack")
		return nil

	case <-time.After(c.cfg.PIC.CommandACKTimeout):
		log.Error("PIC command ack timeout")
		return ErrCommandACKTimeout

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *controller) createESPCommandWithACK(ctx context.Context, cmd espCommand) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // ensure cleanup other goroutines

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.trackingESPCommandACK(ctx, cmd.ID)
	}()

	if err := c.espSerialClient.Write(ctx, cmdJSON); err != nil {
		return fmt.Errorf("write command: %w", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("tracking ESP command ack: %w", err)
		}
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *controller) trackingESPCommandACK(ctx context.Context, id string) error {
	log := c.log.With(slog.String("id", id))
	log.Debug("start tracking ESP command ack")

	doneCh := make(chan struct{})
	c.subscriber.Subscribe(ctx, events.ESPCmdAckTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.ESPCmdAckEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.ID == id {
			if ev.Success {
				log.Debug("ESP command ack received")
			} else {
				log.Error("ESP command ack failed")
			}
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
		log.Debug("stop tracking ESP command ack")
		return nil

	case <-time.After(c.cfg.ESP.CommandACKTimeout):
		log.Error("ESP command ack timeout")
		return ErrCommandACKTimeout

	case <-ctx.Done():
		return ctx.Err()
	}
}
