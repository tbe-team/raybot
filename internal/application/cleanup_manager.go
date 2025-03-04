package application

import (
	"context"
	"sync"
)

type CleanupManager struct {
	mu           sync.Mutex
	cleanupFuncs []func(context.Context) error
}

func NewCleanupManager() *CleanupManager {
	return &CleanupManager{}
}

func (c *CleanupManager) Add(cleanupFunc func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanupFuncs = append(c.cleanupFuncs, cleanupFunc)
}

func (c *CleanupManager) Cleanup(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// reverse order cleanup
	for i := len(c.cleanupFuncs) - 1; i >= 0; i-- {
		if err := c.cleanupFuncs[i](ctx); err != nil {
			return err
		}
	}

	return nil
}
