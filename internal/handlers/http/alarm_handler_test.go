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
	"github.com/tbe-team/raybot/internal/services/alarm"
	alarmmocks "github.com/tbe-team/raybot/internal/services/alarm/mocks"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/ptr"
)

func TestAlarmHandler_ListAlarms(t *testing.T) {
	t.Run("Should list active alarms successfully", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().ListActiveAlarms(mock.Anything,
			mock.MatchedBy(
				func(params alarm.ListActiveAlarmsParams) bool {
					return params.PagingParams.Page == 2 &&
						params.PagingParams.PageSize == 20
				},
			),
		).Return(paging.List[alarm.Alarm]{
			Items:      []alarm.Alarm{validActiveAlarm},
			TotalItems: 1,
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=ACTIVE&page=2&pageSize=20", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.ListAlarms200JSONResponse](t, rec.Body)
		require.Equal(t, 1, res.TotalItems)
		require.Equal(t, validActiveAlarm.ID, res.Items[0].Id)
		require.Equal(t, validActiveAlarm.Type.String(), res.Items[0].Type)
		require.Equal(t, validActiveAlarm.Data.Message(), res.Items[0].Message)
	})

	t.Run("Should list deactive alarms successfully", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().ListDeactiveAlarms(mock.Anything,
			mock.MatchedBy(
				func(params alarm.ListDeactiveAlarmsParams) bool {
					return params.PagingParams.Page == 1 &&
						params.PagingParams.PageSize == 15
				},
			),
		).Return(paging.List[alarm.Alarm]{
			Items:      []alarm.Alarm{validDeactiveAlarm},
			TotalItems: 1,
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=DEACTIVE&page=1&pageSize=15", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.ListAlarms200JSONResponse](t, rec.Body)
		require.Equal(t, 1, res.TotalItems)
		require.Equal(t, validDeactiveAlarm.ID, res.Items[0].Id)
		require.Equal(t, validDeactiveAlarm.Type.String(), res.Items[0].Type)
		require.NotNil(t, res.Items[0].DeactivatedAt)
	})

	t.Run("Should use default paging params if not provided", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().ListActiveAlarms(mock.Anything,
			mock.MatchedBy(
				func(params alarm.ListActiveAlarmsParams) bool {
					return params.PagingParams.Page == 1 && params.PagingParams.PageSize == 10
				},
			),
		).Return(paging.List[alarm.Alarm]{
			Items:      []alarm.Alarm{validActiveAlarm},
			TotalItems: 1,
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=ACTIVE", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should not able to list alarms with invalid status", func(t *testing.T) {
		h := SetupAPITestHandler(t)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=invalid", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Should not able to list active alarms if fetching failed", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().ListActiveAlarms(mock.Anything, mock.Anything).
			Return(paging.List[alarm.Alarm]{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=ACTIVE", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Should not able to list deactive alarms if fetching failed", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().ListDeactiveAlarms(mock.Anything, mock.Anything).
			Return(paging.List[alarm.Alarm]{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/alarms?status=DEACTIVE", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAlarmHandler_DeleteDeactiveAlarms(t *testing.T) {
	t.Run("Should delete deactive alarms successfully", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().DeleteDeactivatedAlarms(mock.Anything).Return(nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/alarms", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should not able to delete deactive alarms if deleting failed", func(t *testing.T) {
		alarmService := alarmmocks.NewFakeService(t)
		alarmService.EXPECT().DeleteDeactivatedAlarms(mock.Anything).
			Return(errors.New("deleting failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.alarmService = alarmService
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/alarms", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

var validActiveAlarm = alarm.Alarm{
	ID:   1,
	Type: alarm.AlarmTypeBatteryVoltageLow,
	Data: alarm.DataBatteryVoltageLow{
		Threshold: 14.0,
		Voltage:   13.5,
	},
	ActivatedAt:   time.Now().Add(-2 * time.Hour),
	DeactivatedAt: nil,
}

var validDeactiveAlarm = alarm.Alarm{
	ID:   2,
	Type: alarm.AlarmTypeBatteryCurrentHigh,
	Data: alarm.DataBatteryCurrentHigh{
		Threshold: 6.0,
		Current:   6.5,
	},
	ActivatedAt:   time.Now().Add(-12 * time.Hour),
	DeactivatedAt: ptr.New(time.Now().Add(-10 * time.Hour)),
}
