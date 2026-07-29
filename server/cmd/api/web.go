package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// The built SPA (web/dist/client) is copied into ./dist at image-build time and
// embedded. `all:` is required so files starting with `_` (e.g. _shell.html)
// and `.` are included.
//
//go:embed all:dist
var distFS embed.FS

// registerWeb serves the SPA at "/" with a client-routing fallback: existing
// files are served directly; anything else returns the shell.
func (app *app) registerWeb() {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		app.logger.Error("web: sub fs", "err", err)
		return
	}
	shell, err := fs.ReadFile(sub, "_shell.html")
	if err != nil {
		// No built frontend (e.g. local `go run` without the copy step). The API
		// still works; we just don't serve static files.
		app.logger.Warn("web: no embedded frontend (_shell.html missing), skipping static serving")
		return
	}
	fileServer := http.FileServer(http.FS(sub))
	app.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p != "" {
			if f, openErr := sub.Open(p); openErr == nil {
				f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(shell)
	})
}
