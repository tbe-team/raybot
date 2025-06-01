package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/apperrorcode"
	apperrorcodemocks "github.com/tbe-team/raybot/internal/services/apperrorcode/mocks"
)

func TestErrorCodeHandler_GetErrorCodes(t *testing.T) {
	t.Run("Should get error codes successfully", func(t *testing.T) {
		apperrorcodeService := apperrorcodemocks.NewFakeService(t)
		apperrorcodeService.EXPECT().ListErrorCodes(mock.Anything).Return([]apperrorcode.ErrorCode{
			{
				Code:    "123",
				Message: "test",
			},
			{
				Code:    "456",
				Message: "test2",
			},
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.apperrorcodeService = apperrorcodeService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/error-codes", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetErrorCodes200JSONResponse](t, rec.Body)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, 2, len(res))
		require.Equal(t, "123", res[0].Code)
		require.Equal(t, "test", res[0].Message)
		require.Equal(t, "456", res[1].Code)
		require.Equal(t, "test2", res[1].Message)
	})

	t.Run("Should not able to get error codes if fetching failed", func(t *testing.T) {
		apperrorcodeService := apperrorcodemocks.NewFakeService(t)
		apperrorcodeService.EXPECT().ListErrorCodes(mock.Anything).
			Return([]apperrorcode.ErrorCode{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.apperrorcodeService = apperrorcodeService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/error-codes", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
