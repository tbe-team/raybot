package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/system"
	systemmocks "github.com/tbe-team/raybot/internal/services/system/mocks"
)

func TestSystemHandler_GetSystemInfo(t *testing.T) {
	t.Run("Should get system info successfully", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().GetInfo(mock.Anything).
			Return(system.Info{
				LocalIP:     "192.168.1.1",
				CPUUsage:    32.5,
				MemoryUsage: 50.2,
				TotalMemory: 1234,
				Uptime:      100 * time.Second,
			}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetSystemInfo200JSONResponse](t, rec.Body)
		require.Equal(t, "192.168.1.1", res.LocalIp)
		require.Equal(t, float32(32.5), res.CpuUsage)
		require.Equal(t, float32(50.2), res.MemoryUsage)
		require.Equal(t, float32(1234), res.TotalMemory)
		require.Equal(t, float32(100), res.Uptime)
	})

	t.Run("Should return error if getting system info failed", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().GetInfo(mock.Anything).Return(system.Info{}, errors.New("failed to get system info"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSystemHandler_RebootSystem(t *testing.T) {
	t.Run("Should reboot system successfully", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().Reboot(mock.Anything).Return(nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/system/reboot", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should return error if rebooting system failed", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().Reboot(mock.Anything).Return(errors.New("failed to reboot system"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/system/reboot", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSystemHandler_StopEmergency(t *testing.T) {
	t.Run("Should stop emergency successfully", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().StopEmergency(mock.Anything).Return(nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/system/stop-emergency", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should return error if stopping emergency failed", func(t *testing.T) {
		systemService := systemmocks.NewFakeService(t)
		systemService.EXPECT().StopEmergency(mock.Anything).Return(errors.New("failed to stop emergency"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.systemService = systemService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/system/stop-emergency", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
