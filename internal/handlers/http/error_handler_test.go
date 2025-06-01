package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/logging"
)

func TestErrorHandler_handleRequestError(t *testing.T) {
	t.Run("Should handle request error successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		service := &Service{
			log: logging.NewNoopLogger(),
		}

		service.handleRequestError(rec, req, errors.New("should be bad request"))

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestErrorHandler_handleResponseError(t *testing.T) {
	t.Run("Should handle response error successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		service := &Service{
			log: logging.NewNoopLogger(),
		}

		service.handleResponseError(rec, req, errors.New("should be internal server error"))

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
