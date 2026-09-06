package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

type deadlineWriter struct {
	*httptest.ResponseRecorder
	deadline time.Time
}

func (w *deadlineWriter) SetWriteDeadline(deadline time.Time) error {
	w.deadline = deadline
	return nil
}

func TestWriteDeadlines(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mux := NewMux()
		mux.HandleFunc("GET /regular", func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, time.Now().Add(writeTimeout), w.(*deadlineWriter).deadline)
			w.WriteHeader(http.StatusOK)
		})
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			require.True(t, write(w, "retry: 5000\n\n"))
			require.True(t, w.(*deadlineWriter).deadline.IsZero())
			time.Sleep(2 * writeTimeout)
			require.True(t, write(w, ":\n\n"))
			require.True(t, w.(*deadlineWriter).deadline.IsZero())
		})
		for _, path := range []string{"/regular", "/missing", "/events"} {
			w := &deadlineWriter{ResponseRecorder: httptest.NewRecorder()}
			mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			if path != "/events" {
				require.Equal(t, time.Now().Add(writeTimeout), w.deadline)
			}
		}
	})
}
