package appconnectionimpl

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/appconnection"
)

type service struct {
	appConnectionRepo appconnection.Repository
}

func NewService(appConnectionRepo appconnection.Repository) appconnection.Service {
	return &service{
		appConnectionRepo: appConnectionRepo,
	}
}

func (s service) UpdateCloudConnection(ctx context.Context, params appconnection.UpdateConnectionParams) error {
	return s.appConnectionRepo.UpdateCloudConnection(ctx, params)
}

func (s service) UpdateESPSerialConnection(ctx context.Context, params appconnection.UpdateConnectionParams) error {
	return s.appConnectionRepo.UpdateESPSerialConnection(ctx, params)
}

func (s service) UpdatePICSerialConnection(ctx context.Context, params appconnection.UpdateConnectionParams) error {
	return s.appConnectionRepo.UpdatePICSerialConnection(ctx, params)
}

func (s service) UpdateRFIDUSBConnection(ctx context.Context, params appconnection.UpdateConnectionParams) error {
	return s.appConnectionRepo.UpdateRFIDUSBConnection(ctx, params)
}
