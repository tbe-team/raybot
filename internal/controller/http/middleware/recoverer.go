package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/tbe-team/raybot/internal/controller/http/apierr"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Recoverer prints a stack trace of the last function call.
func Recoverer(next http.Handler) http.Handler {
	res := apierr.ErrorResponse{
		ErrorResponse: gen.ErrorResponse{
			Code:    "internal_server_error",
			Message: "Internal Server Error",
		},
	}
	errorMsg, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				slog.Error("panic", slog.Any("recover", rvr),
					slog.String("stack", string(debug.Stack())))

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
					//nolint:errcheck
					w.Write(errorMsg)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
