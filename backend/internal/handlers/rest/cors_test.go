package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	tests := map[string]struct {
		origin        string
		allowedOrigin string
	}{
		"AllowedOrigin": {
			origin:        "https://admin.climblive.com",
			allowedOrigin: "https://admin.climblive.com",
		},
		"DisallowedOrigin": {
			origin: "https://example.com",
		},
		"MissingOrigin": {},
	}

	handler := rest.CORSWithOrigins([]string{"https://admin.climblive.com"})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			r.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.allowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}
