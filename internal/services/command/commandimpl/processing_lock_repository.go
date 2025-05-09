package commandimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/command"
)

type processingLockRepository struct {
	mu       sync.RWMutex
	unlockCh chan struct{}
}

func NewProcessingLockRepository() command.ProcessingLockRepository {
	return &processingLockRepository{}
}

func (r *processingLockRepository) Lock(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.unlockCh = make(chan struct{})
	return nil
}

func (r *processingLockRepository) Unlock(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.unlockCh != nil {
		close(r.unlockCh)
		r.unlockCh = nil
	}

	return nil
}

func (r *processingLockRepository) WaitUntilUnlocked(ctx context.Context) error {
	r.mu.RLock()
	ch := r.unlockCh
	r.mu.RUnlock()

	if ch == nil {
		return nil
	}

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
