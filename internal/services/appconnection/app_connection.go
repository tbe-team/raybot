package appconnection

import (
	"context"
	"time"
)

type UpdateConnectionParams struct {
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
	UpdateCloudConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdateESPSerialConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdatePICSerialConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdateRFIDUSBConnection(ctx context.Context, params UpdateConnectionParams) error
}

type Repository interface {
	GetAppConnection(ctx context.Context) (AppConnection, error)
	UpdateCloudConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdateESPSerialConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdatePICSerialConnection(ctx context.Context, params UpdateConnectionParams) error
	UpdateRFIDUSBConnection(ctx context.Context, params UpdateConnectionParams) error
}
