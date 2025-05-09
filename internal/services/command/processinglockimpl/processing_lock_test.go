package processinglockimpl

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessingLock(t *testing.T) {
	t.Run("WithLock should return original error", func(t *testing.T) {
		l := New()
		myErr := fmt.Errorf("myerr")

		err := l.WithLock(func() error {
			return myErr
		})
		assert.ErrorIs(t, err, myErr)
	})

	t.Run("WaitUntilUnlocked should block until the lock is released", func(t *testing.T) {
		l := New()
		waitCh := make(chan struct{}, 1)
		go func() {
			err := l.WaitUntilUnlocked(context.Background())
			assert.NoError(t, err)
			close(waitCh)
		}()

		select {
		case <-waitCh:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("WaitUntilUnlocked did not wait")
		}
	})

	t.Run("WaitUntilUnlocked should not block if the lock is already released", func(t *testing.T) {
		l := New()
		err := l.WaitUntilUnlocked(context.Background())
		assert.NoError(t, err)
	})

	t.Run("WaitUntilUnlocked should return the context error if the context is canceled", func(t *testing.T) {
		l := New()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := l.WaitUntilUnlocked(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Multiple goroutines WaitUntilUnblocked", func(t *testing.T) {
		l := New()

		var counter atomic.Uint32
		var wg sync.WaitGroup

		wg.Add(10)
		for range 10 {
			go func() {
				wg.Done()
				err := l.WaitUntilUnlocked(context.Background())
				assert.NoError(t, err)
				counter.Add(1)
			}()
		}

		wg.Wait()

		err := l.WithLock(func() error {
			return nil
		})
		assert.NoError(t, err)

		assert.Equal(t, 10, int(counter.Load()))
	})
}
