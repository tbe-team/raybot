package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	t.Run("Should get health successfully", func(t *testing.T) {
		h := SetupAPITestHandler(t)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetHealth200JSONResponse](t, rec.Body)
		require.Equal(t, "ok", res.Status)
	})
}
