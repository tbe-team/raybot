package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/peripheral"
	peripheralmocks "github.com/tbe-team/raybot/internal/services/peripheral/mocks"
)

func TestPeripheralHandler_ListAvailableSerialPorts(t *testing.T) {
	t.Run("Should list available serial ports successfully", func(t *testing.T) {
		peripheralService := peripheralmocks.NewFakeService(t)
		peripheralService.EXPECT().ListAvailableSerialPorts(mock.Anything).
			Return([]peripheral.SerialPort{
				{
					Port: "/dev/ttyUSB0",
				},
				{
					Port: "/dev/ttyUSB1",
				},
			}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.peripheralService = peripheralService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/peripherals/serials", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.ListAvailableSerialPorts200JSONResponse](t, rec.Body)
		require.Equal(t, 2, len(res.Items))
		require.Equal(t, "/dev/ttyUSB0", res.Items[0].Port)
		require.Equal(t, "/dev/ttyUSB1", res.Items[1].Port)
	})

	t.Run("Should return error if listing available serial ports failed", func(t *testing.T) {
		peripheralService := peripheralmocks.NewFakeService(t)
		peripheralService.EXPECT().ListAvailableSerialPorts(mock.Anything).Return(nil, errors.New("failed to list available serial ports"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.peripheralService = peripheralService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/peripherals/serials", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
