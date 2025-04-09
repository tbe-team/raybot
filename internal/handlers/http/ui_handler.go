package http

import (
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/tbe-team/raybot/ui"
)

func (s *Service) RegisterUIHandler(r chi.Router) {
	uiFS, err := fs.Sub(ui.GetDist(), "dist")
	if err != nil {
		s.log.Error("failed to get UI file system", slog.Any("error", err))
		return
	}

	// Create a file server using the embedded filesystem
	fileServer := http.FileServer(http.FS(uiFS))

	// Handle the root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusFound)
	})

	// Handle the /ui path
	r.Get("/ui", func(w http.ResponseWriter, _ *http.Request) {
		content, err := fs.ReadFile(uiFS, "index.html")
		if err != nil {
			http.Error(w, "Could not read index.html", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(content); err != nil {
			s.log.Error("Failed to write index.html", slog.Any("error", err))
		}
	})

	// Handle static assets
	r.Get("/ui/assets/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/ui")
		fileServer.ServeHTTP(w, r)
	})

	// Handle all other /ui/* paths for SPA routing
	r.Get("/ui/*", func(w http.ResponseWriter, _ *http.Request) {
		content, err := fs.ReadFile(uiFS, "index.html")
		if err != nil {
			http.Error(w, "Could not read index.html", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(content); err != nil {
			s.log.Error("Failed to write index.html", slog.Any("error", err))
		}
	})

	// Not Found handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/404?path="+r.URL.Path, http.StatusFound)
	})
}
