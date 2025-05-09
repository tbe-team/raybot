package commandimpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessingLockRepository(t *testing.T) {
	t.Run("Lock then Unlock", func(t *testing.T) {
		r := NewProcessingLockRepository()
		ctx := context.Background()

		err := r.Lock(ctx)
		assert.NoError(t, err)

		// Unlock should unblock WaitUntilUnlocked
		waitCh := make(chan error, 1)
		go func() {
			waitCh <- r.WaitUntilUnlocked(ctx)
		}()

		time.Sleep(10 * time.Millisecond)
		err = r.Unlock(ctx)
		assert.NoError(t, err)

		select {
		case err := <-waitCh:
			assert.NoError(t, err)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("WaitUntilUnlocked did not return after Unlock")
		}
	})

	t.Run("WaitUntilUnlocked returns immediately if not locked", func(t *testing.T) {
		r := NewProcessingLockRepository()
		ctx := context.Background()

		err := r.WaitUntilUnlocked(ctx)
		assert.NoError(t, err)
	})

	t.Run("WaitUntilUnlocked returns error if context is canceled", func(t *testing.T) {
		r := NewProcessingLockRepository()
		ctx := context.Background()
		err := r.Lock(ctx)
		assert.NoError(t, err)

		ctx2, cancel := context.WithCancel(ctx)
		waitCh := make(chan error, 1)
		go func() {
			waitCh <- r.WaitUntilUnlocked(ctx2)
		}()

		cancel()

		select {
		case err := <-waitCh:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("WaitUntilUnlocked did not return after context cancel")
		}
	})
}
