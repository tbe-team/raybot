package swagger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/handlers/http/swagger"
)

func TestSwaggerDocsRoute(t *testing.T) {
	r := chi.NewRouter()
	swagger.Register(r)

	t.Run("Should get docs successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	})

	t.Run("Should get openapi.json successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")
	})
}
