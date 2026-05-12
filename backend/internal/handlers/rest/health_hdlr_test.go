package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type healthUseCaseStub struct{}

func (healthUseCaseStub) GetHealth(_ context.Context) ([]domain.ServiceStatus, error) {
	return nil, nil
}

func TestGetVersion(t *testing.T) {
	previousVersion := Version
	Version = "v1.2.3"
	t.Cleanup(func() {
		Version = previousVersion
	})

	mux := NewMux()
	InstallHealthHandler(mux, healthUseCaseStub{})

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	var version string
	err := json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3", version)
}
