package processinglockimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/command"
)

type processingLock struct {
	cond   *sync.Cond
	locked bool
}

func New() command.ProcessingLock {
	return &processingLock{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (r *processingLock) WithLock(fn func() error) error {
	r.cond.L.Lock()
	r.locked = true
	r.cond.L.Unlock()

	err := fn()

	r.cond.L.Lock()
	r.locked = false
	r.cond.Broadcast()
	r.cond.L.Unlock()

	return err
}

func (r *processingLock) WaitUntilUnlocked(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		r.cond.L.Lock()
		defer r.cond.L.Unlock()

		for r.locked {
			r.cond.Wait()
		}

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
