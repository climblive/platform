package rest

import (
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type SPAHandler struct {
	fsys     fs.FS
	basePath string
}

func NewSPAHandler(fsys fs.FS, basePath string) *SPAHandler {
	return &SPAHandler{
		fsys:     fsys,
		basePath: basePath,
	}
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, h.basePath)
	urlPath = strings.TrimPrefix(urlPath, "/")

	if urlPath == "" {
		urlPath = "index.html"
	}

	file, err := h.fsys.Open(urlPath)
	if err != nil {
		h.serveIndex(w, r)
		return
	}
	file.Close()

	h.setHeaders(w, urlPath)

	modifiedReq := r.Clone(r.Context())
	modifiedReq.URL = &url.URL{Path: "/" + urlPath}
	http.FileServer(http.FS(h.fsys)).ServeHTTP(w, modifiedReq)
}

func (h *SPAHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, "index.html")

	content, err := fs.ReadFile(h.fsys, "index.html")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content)
}

func (h *SPAHandler) setHeaders(w http.ResponseWriter, filePath string) {
	ext := path.Ext(filePath)

	if filePath == "index.html" || ext == ".html" {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}
}

func InstallStaticHandler(mux *http.ServeMux, basePath string, fsys fs.FS) {
	handler := NewSPAHandler(fsys, basePath)

	pattern := basePath
	if !strings.HasSuffix(pattern, "/") {
		pattern += "/"
	}

	mux.Handle(pattern, http.StripPrefix(basePath, handler))
}
