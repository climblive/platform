package rest

import (
	"io/fs"
	"net/http"
	"net/url"
	"path"
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
	file.Close()

	h.setHeaders(w, urlPath)

	modifiedReq := r.Clone(r.Context())
	modifiedReq.URL = &url.URL{Path: "/" + urlPath, RawQuery: r.URL.RawQuery}
	h.fileServer.ServeHTTP(w, modifiedReq)
}

func (h *spaHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, "index.html")

	content, err := fs.ReadFile(h.fsys, "index.html")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content)
}

func (h *spaHandler) setHeaders(w http.ResponseWriter, filePath string) {
	ext := path.Ext(filePath)

	if filePath == "index.html" || ext == ".html" {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		if h.csp != "" {
			w.Header().Set("Content-Security-Policy", h.csp)
		}
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}
}
