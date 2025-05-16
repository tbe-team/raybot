package appstateimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/appstate"
)

type repository struct {
	appState appstate.AppState
	mu       sync.RWMutex

	subscribers      map[chan appstate.AppState]struct{}
	appStateChangeCh chan appstate.AppState

	closeCh chan struct{}
}

func NewAppStateRepository() appstate.Repository {
	r := &repository{
		subscribers:      make(map[chan appstate.AppState]struct{}),
		appStateChangeCh: make(chan appstate.AppState),
		closeCh:          make(chan struct{}),
	}

	go r.notifySubscribersLoop()

	return r
}

func (r *repository) Cleanup() {
	close(r.closeCh)
}

func (r *repository) GetAppState(_ context.Context) (appstate.AppState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.appState, nil
}

func (r *repository) UpdateCloudConnection(_ context.Context, params appstate.UpdateCloudConnectionParams) error {
	r.mu.RLock()
	cloudConnection := r.appState.CloudConnection
	r.mu.RUnlock()
	if params.SetConnected {
		cloudConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		cloudConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		cloudConnection.Error = params.Error
	}

	r.mu.Lock()
	r.appState.CloudConnection = cloudConnection
	r.mu.Unlock()

	r.appStateChangeCh <- r.appState

	return nil
}

func (r *repository) UpdateESPSerialConnection(_ context.Context, params appstate.UpdateESPSerialConnectionParams) error {
	r.mu.RLock()
	espSerialConnection := r.appState.ESPSerialConnection
	r.mu.RUnlock()

	if params.SetConnected {
		espSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		espSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		espSerialConnection.Error = params.Error
	}

	r.mu.Lock()
	r.appState.ESPSerialConnection = espSerialConnection
	r.mu.Unlock()

	r.appStateChangeCh <- r.appState

	return nil
}

func (r *repository) UpdatePICSerialConnection(_ context.Context, params appstate.UpdatePICSerialConnectionParams) error {
	r.mu.RLock()
	picSerialConnection := r.appState.PICSerialConnection
	r.mu.RUnlock()

	if params.SetConnected {
		picSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		picSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		picSerialConnection.Error = params.Error
	}

	r.mu.Lock()
	r.appState.PICSerialConnection = picSerialConnection
	r.mu.Unlock()

	r.appStateChangeCh <- r.appState

	return nil
}

func (r *repository) UpdateRFIDUSBConnection(_ context.Context, params appstate.UpdateRFIDUSBConnectionParams) error {
	r.mu.RLock()
	rfidUsbConnection := r.appState.RFIDUSBConnection
	r.mu.RUnlock()

	if params.SetConnected {
		rfidUsbConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		rfidUsbConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		rfidUsbConnection.Error = params.Error
	}

	r.mu.Lock()
	r.appState.RFIDUSBConnection = rfidUsbConnection
	r.mu.Unlock()

	r.appStateChangeCh <- r.appState

	return nil
}

func (r *repository) ListenForAppStateChanges(ctx context.Context) <-chan appstate.AppState {
	ch := make(chan appstate.AppState, 1)

	r.mu.Lock()
	r.subscribers[ch] = struct{}{}
	r.mu.Unlock()

	go func() {
		<-ctx.Done()

		r.mu.Lock()
		delete(r.subscribers, ch)
		r.mu.Unlock()

		close(ch)
	}()

	return ch
}

func (r *repository) notifySubscribersLoop() {
	for {
		select {
		case <-r.closeCh:
			return

		case state := <-r.appStateChangeCh:
			r.mu.RLock()
			for ch := range r.subscribers {
				select {
				case ch <- state:
				default:
				}
			}
			r.mu.RUnlock()
		}
	}
}
