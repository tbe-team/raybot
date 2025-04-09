package appconnectionimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/appconnection"
)

type repository struct {
	appConnection appconnection.AppConnection
	mu            sync.RWMutex
}

func NewAppConnectionRepository() appconnection.Repository {
	return &repository{}
}

func (r *repository) GetAppConnection(_ context.Context) (appconnection.AppConnection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.appConnection, nil
}

func (r *repository) UpdateCloudConnection(_ context.Context, params appconnection.UpdateConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cloudConnection := r.appConnection.CloudConnection
	if params.SetConnected {
		cloudConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		cloudConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		cloudConnection.Error = params.Error
	}
	r.appConnection.CloudConnection = cloudConnection

	return nil
}

func (r *repository) UpdateESPSerialConnection(_ context.Context, params appconnection.UpdateConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	espSerialConnection := r.appConnection.ESPSerialConnection
	if params.SetConnected {
		espSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		espSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		espSerialConnection.Error = params.Error
	}
	r.appConnection.ESPSerialConnection = espSerialConnection

	return nil
}

func (r *repository) UpdatePICSerialConnection(_ context.Context, params appconnection.UpdateConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	picSerialConnection := r.appConnection.PICSerialConnection
	if params.SetConnected {
		picSerialConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		picSerialConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		picSerialConnection.Error = params.Error
	}
	r.appConnection.PICSerialConnection = picSerialConnection

	return nil
}

func (r *repository) UpdateRFIDUSBConnection(_ context.Context, params appconnection.UpdateConnectionParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rfidUsbConnection := r.appConnection.RFIDUSBConnection
	if params.SetConnected {
		rfidUsbConnection.Connected = params.Connected
	}
	if params.SetLastConnectedAt {
		rfidUsbConnection.LastConnectedAt = params.LastConnectedAt
	}
	if params.SetError {
		rfidUsbConnection.Error = params.Error
	}
	r.appConnection.RFIDUSBConnection = rfidUsbConnection

	return nil
}
