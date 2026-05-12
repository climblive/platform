package rest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_GetVersion(t *testing.T) {
	mux := rest.NewMux()
	rest.InstallHealthHandler(mux, healthUseCaseMock{})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/version")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	var version string
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err)

	buildInfo, ok := debug.ReadBuildInfo()
	require.True(t, ok)

	assert.Equal(t, buildInfo.GoVersion, version)
}

type healthUseCaseMock struct{}

func (m healthUseCaseMock) GetHealth(_ context.Context) ([]domain.ServiceStatus, error) {
	return nil, nil
}
