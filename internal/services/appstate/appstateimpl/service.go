package appstateimpl

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/appstate"
)

type service struct {
	appStateRepo appstate.Repository
}

func NewService(appStateRepo appstate.Repository) appstate.Service {
	return &service{
		appStateRepo: appStateRepo,
	}
}

func (s service) UpdateCloudConnection(ctx context.Context, params appstate.UpdateCloudConnectionParams) error {
	return s.appStateRepo.UpdateCloudConnection(ctx, params)
}

func (s service) UpdateESPSerialConnection(ctx context.Context, params appstate.UpdateESPSerialConnectionParams) error {
	return s.appStateRepo.UpdateESPSerialConnection(ctx, params)
}

func (s service) UpdatePICSerialConnection(ctx context.Context, params appstate.UpdatePICSerialConnectionParams) error {
	return s.appStateRepo.UpdatePICSerialConnection(ctx, params)
}

func (s service) UpdateRFIDUSBConnection(ctx context.Context, params appstate.UpdateRFIDUSBConnectionParams) error {
	return s.appStateRepo.UpdateRFIDUSBConnection(ctx, params)
}
