package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tbe-team/raybot/internal/handlers/http/apierr"
	"github.com/tbe-team/raybot/pkg/xerror"
)

func (s *Service) handleRequestError(w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	err = xerror.ValidationFailed(err, err.Error())
	res := apierr.New(err)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Warn("error encoding error request", slog.Any("error", err))
	}
}

func (s *Service) handleResponseError(w http.ResponseWriter, r *http.Request, err error) {
	res := apierr.New(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)

	logLevel := slog.LevelInfo
	if res.StatusCode >= 500 {
		logLevel = slog.LevelError
	} else if res.StatusCode >= 400 {
		logLevel = slog.LevelWarn
	}
	s.log.Log(r.Context(), logLevel, "http error", slog.Any("error", err))

	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error("error encoding error response", slog.Any("error", err))
	}
}
