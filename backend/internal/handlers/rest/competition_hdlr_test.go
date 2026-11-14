package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type competitionHandlerUseCaseMock struct {
	competitionUseCase
	mock.Mock
}

func (m *competitionHandlerUseCaseMock) GetCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error) {
	args := m.Called(ctx, competitionID)
	return args.Get(0).(domain.Competition), args.Error(1)
}

func TestGetCompetitionHandler(t *testing.T) {
	useCase := new(competitionHandlerUseCaseMock)
	useCase.On("GetCompetition", mock.Anything, domain.CompetitionID(42)).Return(domain.Competition{
		ID:   42,
		Name: "World Championships",
	}, nil)

	mux := NewMux()
	InstallCompetitionHandler(mux, useCase, nil, nil, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/competitions/42", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"id":42`)
	assert.Contains(t, response.Body.String(), `"name":"World Championships"`)
	useCase.AssertExpectations(t)
}

func TestCompetitionRoutes(t *testing.T) {
	mux := NewMux()
	InstallCompetitionHandler(mux, nil, nil, nil, nil)
	InstallCompClassHandler(mux, nil)
	InstallContenderHandler(mux, nil)
	InstallProblemHandler(mux, nil)
	InstallTickHandler(mux, nil)
	InstallRaffleHandler(mux, nil)
	InstallScoreEngineHandler(mux, nil)

	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/competitions/invalid"},
		{http.MethodPatch, "/competitions/invalid"},
		{http.MethodGet, "/competitions/invalid/scoreboard"},
		{http.MethodGet, "/competitions/invalid/results"},
		{http.MethodPost, "/competitions/invalid/duplicate"},
		{http.MethodPost, "/competitions/invalid/transfer"},
		{http.MethodPost, "/competitions/invalid/archive"},
		{http.MethodPost, "/competitions/invalid/restore"},
		{http.MethodGet, "/organizers/invalid/competitions"},
		{http.MethodPost, "/organizers/invalid/competitions"},
		{http.MethodGet, "/competitions/invalid/comp-classes"},
		{http.MethodPost, "/competitions/invalid/comp-classes"},
		{http.MethodGet, "/competitions/invalid/contenders"},
		{http.MethodPost, "/competitions/invalid/contenders"},
		{http.MethodGet, "/competitions/invalid/problems"},
		{http.MethodPost, "/competitions/invalid/problems"},
		{http.MethodGet, "/competitions/invalid/ticks"},
		{http.MethodGet, "/competitions/invalid/raffles"},
		{http.MethodPost, "/competitions/invalid/raffles"},
		{http.MethodGet, "/competitions/invalid/score-engines"},
		{http.MethodPost, "/competitions/invalid/score-engines"},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, nil))
			assert.Equal(t, http.StatusBadRequest, response.Code)
		})
	}
}
