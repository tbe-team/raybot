package appstate

import (
	"context"
	"time"
)

type UpdateCloudConnectionParams struct {
	Connected          bool
	SetConnected       bool
	LastConnectedAt    *time.Time
	SetLastConnectedAt bool
	Error              *string
	SetError           bool
}

type UpdateESPSerialConnectionParams struct {
	Connected          bool
	SetConnected       bool
	LastConnectedAt    *time.Time
	SetLastConnectedAt bool
	Error              *string
	SetError           bool
}

type UpdatePICSerialConnectionParams struct {
	Connected          bool
	SetConnected       bool
	LastConnectedAt    *time.Time
	SetLastConnectedAt bool
	Error              *string
	SetError           bool
}

type UpdateRFIDUSBConnectionParams struct {
	Connected          bool
	SetConnected       bool
	LastConnectedAt    *time.Time
	SetLastConnectedAt bool
	Error              *string
	SetError           bool
}

type Service interface {
	UpdateCloudConnection(ctx context.Context, params UpdateCloudConnectionParams) error
	UpdateESPSerialConnection(ctx context.Context, params UpdateESPSerialConnectionParams) error
	UpdatePICSerialConnection(ctx context.Context, params UpdatePICSerialConnectionParams) error
	UpdateRFIDUSBConnection(ctx context.Context, params UpdateRFIDUSBConnectionParams) error
}

type Repository interface {
	GetAppState(ctx context.Context) (AppState, error)
	UpdateCloudConnection(ctx context.Context, params UpdateCloudConnectionParams) error
	UpdateESPSerialConnection(ctx context.Context, params UpdateESPSerialConnectionParams) error
	UpdatePICSerialConnection(ctx context.Context, params UpdatePICSerialConnectionParams) error
	UpdateRFIDUSBConnection(ctx context.Context, params UpdateRFIDUSBConnectionParams) error
}
