package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
)

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	handler := rest.CORS([]string{"https://climblive.com"})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	r.Header.Set("Origin", "https://climblive.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, "https://climblive.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSRejectsUnconfiguredOrigin(t *testing.T) {
	handler := rest.CORS([]string{"https://climblive.com"})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	r.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSPreFlightAllowsConfiguredOrigin(t *testing.T) {
	handler := rest.CORSPreFlight([]string{"https://climblive.com"})
	r := httptest.NewRequest(http.MethodOptions, "http://localhost", nil)
	r.Header.Set("Origin", "https://climblive.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://climblive.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", w.Header().Get("Vary"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Authorization, Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestCORSPreFlightRejectsUnconfiguredOrigin(t *testing.T) {
	handler := rest.CORSPreFlight([]string{"https://climblive.com"})
	r := httptest.NewRequest(http.MethodOptions, "http://localhost", nil)
	r.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Vary"))
}
