package rest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMuxWithMiddlewares(t *testing.T) {
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("A\n"))

			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("B\n"))

			next.ServeHTTP(w, r)
		})
	}

	mw3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("C\n"))

			next.ServeHTTP(w, r)
		})
	}

	r := httptest.NewRequest("GET", "http://localhost/", nil)
	w := httptest.NewRecorder()

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Mux!"))
	}

	mux := rest.NewMux()

	mux.RegisterMiddleware(mw1)
	mux.RegisterMiddleware(mw2)
	mux.RegisterMiddleware(mw3)

	mux.HandleFunc("/", dummyHandler)

	mux.ServeHTTP(w, r)

	result := w.Result()

	buf := new(strings.Builder)
	_, err := io.Copy(buf, result.Body)

	require.NoError(t, err)
	assert.Equal(t, `A
B
C
Hello, Mux!`, buf.String())
}

func TestMuxWithoutMiddlewares(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/", nil)
	w := httptest.NewRecorder()

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Mux!"))
	}

	mux := rest.NewMux()

	mux.HandleFunc("/", dummyHandler)

	mux.ServeHTTP(w, r)

	result := w.Result()

	buf := new(strings.Builder)
	_, err := io.Copy(buf, result.Body)

	require.NoError(t, err)
	assert.Equal(t, "Hello, Mux!", buf.String())
}
