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
	"github.com/tbe-team/raybot/internal/services/limitswitch"
	limitswitchmocks "github.com/tbe-team/raybot/internal/services/limitswitch/mocks"
)

func TestStateHandler_GetLimitSwitchState(t *testing.T) {
	t.Run("Should get limit switch state successfully", func(t *testing.T) {
		limitSwitchService := limitswitchmocks.NewFakeService(t)
		limitSwitchService.EXPECT().GetLimitSwitchState(mock.Anything).
			Return(limitswitch.GetLimitSwitchStateOutput{
				LimitSwitch1: limitswitch.LimitSwitch{
					Pressed:   true,
					UpdatedAt: time.Now(),
				},
			}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.limitSwitchService = limitSwitchService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/states/limit-switch", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetLimitSwitchState200JSONResponse](t, rec.Body)
		require.Equal(t, true, res.LimitSwitch1.Pressed)
		require.NotEmpty(t, res.LimitSwitch1.UpdatedAt)
	})

	t.Run("Should return error if getting limit switch state failed", func(t *testing.T) {
		limitSwitchService := limitswitchmocks.NewFakeService(t)
		limitSwitchService.EXPECT().GetLimitSwitchState(mock.Anything).
			Return(limitswitch.GetLimitSwitchStateOutput{}, errors.New("failed to get limit switch state"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.limitSwitchService = limitSwitchService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/states/limit-switch", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
