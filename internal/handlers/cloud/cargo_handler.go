package cloud

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	cargov1 "github.com/tbe-team/raybot-api/cargo/v1"
	"github.com/tbe-team/raybot/internal/services/cargo"
)

type cargoHandler struct {
	cargov1.UnimplementedCargoServiceServer
	cargoService cargo.Service
}

func newCargoHandler(cargoService cargo.Service) cargov1.CargoServiceServer {
	return &cargoHandler{
		cargoService: cargoService,
	}
}

func (h cargoHandler) GetCargo(ctx context.Context, _ *cargov1.GetCargoRequest) (*cargov1.GetCargoResponse, error) {
	state, err := h.cargoService.GetCargo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cargo state: %w", err)
	}

	return &cargov1.GetCargoResponse{
		IsOpen:         state.IsOpen,
		QrCode:         state.QRCode,
		BottomDistance: uint32(state.BottomDistance),
		HasItem:        state.HasItem,
		UpdatedAt:      timestamppb.New(state.UpdatedAt),
	}, nil
}
