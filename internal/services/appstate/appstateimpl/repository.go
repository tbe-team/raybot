package appstateimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/appstate"
)

type repository struct {
	appState appstate.AppState
	mu       sync.RWMutex
}

func NewAppStateRepository() appstate.Repository {
	return &repository{}
}

func (r *repository) GetAppState(_ context.Context) (appstate.AppState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.appState, nil
}

func (r *repository) UpdateCloudConnection(_ context.Context, params appstate.UpdateCloudConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cloudConnection := r.appState.CloudConnection
	if params.SetConnected {
		cloudConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		cloudConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		cloudConnection.Error = params.Error
	}
	r.appState.CloudConnection = cloudConnection

	return nil
}

func (r *repository) UpdateESPSerialConnection(_ context.Context, params appstate.UpdateESPSerialConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	espSerialConnection := r.appState.ESPSerialConnection
	if params.SetConnected {
		espSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		espSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		espSerialConnection.Error = params.Error
	}
	r.appState.ESPSerialConnection = espSerialConnection

	return nil
}

func (r *repository) UpdatePICSerialConnection(_ context.Context, params appstate.UpdatePICSerialConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	picSerialConnection := r.appState.PICSerialConnection
	if params.SetConnected {
		picSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		picSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		picSerialConnection.Error = params.Error
	}
	r.appState.PICSerialConnection = picSerialConnection

	return nil
}

func (r *repository) UpdateRFIDUSBConnection(_ context.Context, params appstate.UpdateRFIDUSBConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rfidUsbConnection := r.appState.RFIDUSBConnection
	if params.SetConnected {
		rfidUsbConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		rfidUsbConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		rfidUsbConnection.Error = params.Error
	}
	r.appState.RFIDUSBConnection = rfidUsbConnection

	return nil
}
