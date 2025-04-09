package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Logging returns a logger middleware for chi, that implements the http.Handler interface.
func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				logLevel := slog.LevelInfo
				if ww.Status() >= 500 {
					logLevel = slog.LevelError
				} else if ww.Status() >= 400 {
					logLevel = slog.LevelWarn
				}

				log.Log(r.Context(), logLevel, "http request",
					slog.Duration("latency", time.Since(t1)), // Duration
					slog.Int("status", ww.Status()),          // Status code
					slog.String("method", r.Method),          // HTTP method
					slog.String("path", r.URL.Path),          // Request URI
					slog.String("query", r.URL.RawQuery),     // Request query string
					slog.String("remote_ip", r.RemoteAddr),   // IP address
					slog.String("host", r.Host),              // Host
					slog.String("user_agent", r.UserAgent()), // User agent
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
