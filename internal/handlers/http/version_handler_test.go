package http

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/build"
	"github.com/tbe-team/raybot/internal/handlers/http/gen"
)

func TestVersionHandler_GetVersion(t *testing.T) {
	t.Run("Should get version successfully", func(t *testing.T) {
		h := SetupAPITestHandler(t)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/version", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetVersion200JSONResponse](t, rec.Body)
		require.Equal(t, build.Version, res.Version)
		require.Equal(t, build.Date, res.BuildDate)
		require.Equal(t, runtime.Version(), res.GoVersion)
	})
}
