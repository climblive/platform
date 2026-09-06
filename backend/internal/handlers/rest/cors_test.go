package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
)

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	handler := rest.CORS([]string{"https://admin.climblive.com"})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	r.Header.Set("Origin", "https://admin.climblive.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, "https://admin.climblive.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSRejectsUnconfiguredOrigin(t *testing.T) {
	handler := rest.CORS([]string{"https://admin.climblive.com"})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	r.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
