package jobs

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

const (
	executeCommandInterval = 1 * time.Second
)

type executeCommandHandler struct {
	running atomic.Bool

	log            *slog.Logger
	commandService command.Service
	subscriber     eventbus.Subscriber
}

func newExecuteCommandHandler(
	log *slog.Logger,
	commandService command.Service,
	subscriber eventbus.Subscriber,
) *executeCommandHandler {
	return &executeCommandHandler{
		log:            log.With("service", "execute_command_handler"),
		commandService: commandService,
		subscriber:     subscriber,
	}
}

func (h *executeCommandHandler) Run(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)

	go h.run(ctx)

	return cancel
}

func (h *executeCommandHandler) run(ctx context.Context) {
	ch := make(chan struct{}, 1)

	h.subscriber.Subscribe(ctx, events.CommandCreatedTopic, func(_ context.Context, _ *eventbus.Message) {
		select {
		case ch <- struct{}{}:
		default:
		}
	})

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(executeCommandInterval):
			h.handle(ctx)

		case <-ch:
			h.handle(ctx)
		}
	}
}

func (h *executeCommandHandler) handle(ctx context.Context) {
	if h.running.Swap(true) {
		// if already running, do nothing
		return
	}
	defer h.running.Store(false)

	if err := h.commandService.RunNextExecutableCommand(ctx); err != nil {
		h.log.Error("failed to find next executable command and run", slog.Any("error", err))
	}
}
