package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/foo", nil)
	w := httptest.NewRecorder()

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
	}

	handler := rest.CORS(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
