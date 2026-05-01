package rest

import (
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

func InstallStaticHandler(mux *http.ServeMux, basePath string, fsys fs.FS, csp string) {
	handler := &spaHandler{
		fsys:       fsys,
		fileServer: http.FileServer(http.FS(fsys)),
		csp:        csp,
	}

	pattern := basePath
	if !strings.HasSuffix(pattern, "/") {
		pattern += "/"
	}

	mux.Handle(pattern, http.StripPrefix(basePath, handler))
}

type spaHandler struct {
	fsys       fs.FS
	fileServer http.Handler
	csp        string
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/")

	if urlPath == "" {
		h.serveIndex(w, r)
		return
	}

	file, err := h.fsys.Open(urlPath)
	if err != nil {
		h.serveIndex(w, r)
		return
	}

	if err := file.Close(); err != nil {
		h.serveIndex(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")

	r.URL.Path = "/" + urlPath
	h.fileServer.ServeHTTP(w, r)
}

func (h *spaHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	content, err := fs.ReadFile(h.fsys, "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")

	if h.csp != "" {
		w.Header().Set("Content-Security-Policy", h.csp)
	}

	_, err = w.Write(content)
	if err != nil {
		slog.Error("failed to write response", "error", err)
	}
}
