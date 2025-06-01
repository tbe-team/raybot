package http

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/logging"
)

// APITestHandlerOption option func for customizing HTTPServer configuration
// when setting up an API test handler via SetupAPITestHandler.
type APITestHandlerOption func(hs *Service)

// SetupAPITestHandler sets up an API test handler ready to be used in tests.
// Optionally customize the service by passing in APITestHandlerOption functions.
func SetupAPITestHandler(t *testing.T, opts ...APITestHandlerOption) http.Handler {
	t.Helper()

	svc := Service{
		log: logging.NewNoopLogger(),
	}

	for _, opt := range opts {
		opt(&svc)
	}

	r := chi.NewRouter()
	svc.RegisterHandlers(r)

	return r
}

// MustDecodeJSON decodes the JSON response body into the given type.
// Useful for testing API responses.
func MustDecodeJSON[T any](t *testing.T, body io.Reader) T {
	t.Helper()

	var data T
	err := json.NewDecoder(body).Decode(&data)
	require.NoError(t, err)

	return data
}
