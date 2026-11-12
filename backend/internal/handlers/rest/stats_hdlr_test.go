package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type statsUseCaseMock struct {
	mock.Mock
}

func (m *statsUseCaseMock) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.GlobalStats), args.Error(1)
}

func TestGlobalStatsHandler(t *testing.T) {
	for _, tc := range []struct {
		name   string
		stats  domain.GlobalStats
		err    error
		status int
		body   string
	}{
		{
			name:   "AnonymousCrossOriginRequest",
			stats:  domain.GlobalStats{Contests: 12, RegisteredContenders: 345, Ticks: 6789},
			status: http.StatusOK,
			body:   `{"contests":12,"registeredContenders":345,"ticks":6789}`,
		},
		{
			name:   "Empty",
			status: http.StatusOK,
			body:   `{"contests":0,"registeredContenders":0,"ticks":0}`,
		},
		{
			name:   "Error",
			err:    errors.New("database unavailable"),
			status: http.StatusInternalServerError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			uc := new(statsUseCaseMock)
			uc.On("GetGlobalStats", mock.Anything).Return(tc.stats, tc.err).Once()
			mux := rest.NewMux()
			mux.RegisterMiddleware(rest.CORS)
			mux.RegisterMiddleware(authorizer.NewAuthorizer(nil, nil).Middleware)
			rest.InstallStatsHandler(mux, uc)
			request := httptest.NewRequest(http.MethodGet, "/stats/global", nil)
			request.Header.Set("Origin", "https://climblive.com")
			response := httptest.NewRecorder()

			mux.ServeHTTP(response, request)

			assert.Equal(t, tc.status, response.Code)
			assert.Equal(t, "*", response.Header().Get("Access-Control-Allow-Origin"))
			if tc.body != "" {
				assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
				assert.JSONEq(t, tc.body, response.Body.String())
			} else {
				assert.Empty(t, response.Body.String())
			}
			uc.AssertExpectations(t)
		})
	}
}
