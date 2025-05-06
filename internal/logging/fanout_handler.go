package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
)

var _ slog.Handler = (*FanoutHandler)(nil)

type FanoutHandler struct {
	handlers []slog.Handler
}

// Fanout distributes records to multiple slog.Handler.
func Fanout(handlers ...slog.Handler) slog.Handler {
	return &FanoutHandler{
		handlers: handlers,
	}
}

func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := tryRecover(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// If errs is empty, or contains only nil errors, this returns nil
	return errors.Join(errs...)
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i := range h.handlers {
		handlers[i] = h.handlers[i].WithAttrs(slices.Clone(attrs))
	}
	return Fanout(handlers...)
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	handlers := make([]slog.Handler, len(h.handlers))
	for i := range h.handlers {
		handlers[i] = h.handlers[i].WithGroup(name)
	}
	return Fanout(handlers...)
}

func tryRecover(callback func() error) (err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			if e, ok := rvr.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unexpected error: %+v", rvr)
			}
		}
	}()

	err = callback()

	return
}
